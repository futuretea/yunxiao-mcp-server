package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type mrListOptions struct {
	OrganizationID  string
	State           string
	Search          string
	AuthorUserIDs   string
	AssigneeUserIDs string
	OrderBy         string
	TargetBranch    string
	CreatedAfter    string
	CreatedBefore   string
	Page            int
	PerPage         int
	JSONOutput      bool
	OutputFormat string
}

func newYunxiaoRepoMrCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "mr",
		Aliases: []string{"merge-request"},
		Short:   "work with legacy CodeUp merge requests",
	}
	command.AddCommand(newYunxiaoRepoMrListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoRepoMrViewCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoRepoMrListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options mrListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list legacy CodeUp merge requests",
		Example: `  # List merge requests
  yunxiao repo mr list

  # Filter by state
  yunxiao repo mr list --state opened

  # Filter by author
  yunxiao repo mr list --author user1

  # Output as JSON
  yunxiao repo mr list --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_merge_requests", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput || options.OutputFormat == "json" {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printMRList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.State, "state", "", "merge request state: opened, merged, closed, or all")
	flags.StringVar(&options.Search, "search", "", "title search keyword")
	flags.StringVar(&options.AuthorUserIDs, "author", "", "author user ID")
	flags.StringVar(&options.AssigneeUserIDs, "assignee", "", "assignee user ID")
	flags.StringVar(&options.OrderBy, "order-by", "", "sort field: id or updated_at")
	flags.StringVar(&options.TargetBranch, "target-branch", "", "target branch filter")
	flags.StringVar(&options.CreatedAfter, "created-after", "", "created-after date in yyyy-MM-dd format")
	flags.StringVar(&options.CreatedBefore, "created-before", "", "created-before date in yyyy-MM-dd format")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	flags.StringVar(&options.OutputFormat, "output", "", "output format: table, json, or csv")
	return command
}

func (o mrListOptions) params() map[string]any {
	params := map[string]any{}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "state", o.State)
	setCLIStringParam(params, "search", o.Search)
	setCLIStringParam(params, "orderBy", o.OrderBy)
	setCLIStringParam(params, "targetBranch", o.TargetBranch)
	setCLIStringParam(params, "createdAfter", o.CreatedAfter)
	setCLIStringParam(params, "createdBefore", o.CreatedBefore)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	if v := strings.TrimSpace(o.AuthorUserIDs); v != "" {
		params["authorUserIds"] = []string{v}
	}
	if v := strings.TrimSpace(o.AssigneeUserIDs); v != "" {
		params["assigneeUserIds"] = []string{v}
	}
	return params
}

func printMRList(out anyWriter, raw string) error {
	rows, ok := mrRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tIID\tTITLE\tSTATE\tAUTHOR\tTARGET"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n", row.ID, row.IID, row.Title, row.State, row.Author, row.Target)
	}
	return writer.Flush()
}

type mrRow struct {
	ID     string
	IID    string
	Title  string
	State  string
	Author string
	Target string
}


func mrRowsFromJSONForPrint(raw string) ([]mrRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]mrRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, mrRow{
			ID:     firstStringValue(m, "id", "mergeRequestId"),
			IID:    firstStringValue(m, "iid", "localId", "localID"),
			Title:  firstStringValue(m, "title", "name"),
			State:  firstStringValue(m, "state", "status"),
			Author: mrAuthorValue(m),
			Target: firstStringValue(m, "targetBranch", "target", "targetBranchName"),
		})
	}
	return rows, true
}

type mrViewOptions struct {
	OrganizationID string
	RepositoryID   string
	MergeRequestID string
}

func newYunxiaoRepoMrViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options mrViewOptions
	command := &cobra.Command{
		Use:   "view <mr-id>",
		Short: "view a legacy CodeUp merge request as JSON",
		Example: `  # View merge request by ID
  yunxiao repo mr view 12345 --repository-id group/repo`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.MergeRequestID = args[0]
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "get_merge_request", params)
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.RepositoryID, "repository-id", "", "repository numeric ID or full path, e.g. org/repo")
	return command
}

func (o mrViewOptions) params() (map[string]any, error) {
	params := map[string]any{
		"mergeRequestId": strings.TrimSpace(o.MergeRequestID),
		"iid":            strings.TrimSpace(o.MergeRequestID),
		"repositoryId":   strings.TrimSpace(o.RepositoryID),
	}
	if params["mergeRequestId"] == "" {
		return nil, fmt.Errorf("merge request ID argument is required")
	}
	if params["repositoryId"] == "" {
		return nil, fmt.Errorf("repository-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	return params, nil
}

func mrAuthorValue(m map[string]any) string {
	if author, ok := m["author"].(map[string]any); ok {
		return firstStringValue(author, "username", "name", "displayName", "id")
	}
	return firstStringValue(m, "authorName", "authorUsername", "authorId", "authorUserID")
}
