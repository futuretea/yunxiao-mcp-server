package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type repoChangeRequestListOptions struct {
	OrganizationID string
	ProjectIDs     string
	AuthorIDs      string
	ReviewerIDs    string
	State          string
	Search         string
	OrderBy        string
	Sort           string
	CreatedBefore  string
	CreatedAfter   string
	Page           int
	PerPage        int
	JSONOutput     bool
}

func newYunxiaoRepoChangeRequestCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "change-request",
		Aliases: []string{"change-requests", "cr"},
		Short:   "work with CodeUp change requests",
	}
	command.AddCommand(newYunxiaoRepoChangeRequestListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoRepoChangeRequestViewCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoRepoCRPatchSetCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoRepoChangeRequestListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options repoChangeRequestListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list CodeUp change requests",
		Example: `  # List change requests
  yunxiao repo change-request list --project-ids group/repo

  # Filter by state and author
  yunxiao repo change-request list --project-ids group/repo --state opened --author-ids user1

  # Output as JSON
  yunxiao repo change-request list --project-ids group/repo --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_change_requests", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printRepoChangeRequestList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.ProjectIDs, "project-ids", "", "comma-separated repository IDs or full paths")
	flags.StringVar(&options.AuthorIDs, "author-ids", "", "comma-separated author user IDs")
	flags.StringVar(&options.ReviewerIDs, "reviewer-ids", "", "comma-separated reviewer user IDs")
	flags.StringVar(&options.State, "state", "", "change request state, e.g. opened, merged, or closed")
	flags.StringVar(&options.Search, "search", "", "title search keyword")
	flags.StringVar(&options.OrderBy, "order-by", "", "sort field, e.g. created_at or updated_at")
	flags.StringVar(&options.Sort, "sort", "", "sort direction, e.g. asc or desc")
	flags.StringVar(&options.CreatedBefore, "created-before", "", "created-before time in ISO 8601 format")
	flags.StringVar(&options.CreatedAfter, "created-after", "", "created-after time in ISO 8601 format")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o repoChangeRequestListOptions) params() map[string]any {
	params := map[string]any{}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "projectIds", o.ProjectIDs)
	setCLIStringParam(params, "authorIds", o.AuthorIDs)
	setCLIStringParam(params, "reviewerIds", o.ReviewerIDs)
	setCLIStringParam(params, "state", o.State)
	setCLIStringParam(params, "search", o.Search)
	setCLIStringParam(params, "orderBy", o.OrderBy)
	setCLIStringParam(params, "sort", o.Sort)
	setCLIStringParam(params, "createdBefore", o.CreatedBefore)
	setCLIStringParam(params, "createdAfter", o.CreatedAfter)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params
}

func printRepoChangeRequestList(out anyWriter, raw string) error {
	rows, ok := repoChangeRequestRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("LOCAL_ID\tREPOSITORY\tTITLE\tSTATE\tAUTHOR\tSOURCE\tTARGET"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", row.LocalID, row.Repository, row.Title, row.State, row.Author, row.SourceBranch, row.TargetBranch)
	}
	return writer.Flush()
}

type repoChangeRequestRow struct {
	LocalID      string
	Repository   string
	Title        string
	State        string
	Author       string
	SourceBranch string
	TargetBranch string
}

type crPatchSetListOptions struct {
	OrganizationID string
	RepositoryID   string
	LocalID        string
	JSONOutput     bool
}

func newYunxiaoRepoCRPatchSetCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "patch-set",
		Aliases: []string{"patch-sets", "patches"},
		Short:   "list patch sets for a change request",
	}
	command.AddCommand(newYunxiaoRepoCRPatchSetListCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoRepoCRPatchSetListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options crPatchSetListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list patch sets (diff iterations) for a change request",
		Example: `  # List patch sets for a CR
  yunxiao repo cr patches list --repository-id group/repo --local-id 42`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_change_request_patch_sets", params)
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printCRPatchSetList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.RepositoryID, "repository-id", "", "repository numeric ID or full path, e.g. org/repo")
	flags.StringVar(&options.LocalID, "local-id", "", "change request local ID")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o crPatchSetListOptions) params() (map[string]any, error) {
	params := map[string]any{
		"repositoryId": strings.TrimSpace(o.RepositoryID),
		"localId":      strings.TrimSpace(o.LocalID),
	}
	if params["repositoryId"] == "" {
		return nil, fmt.Errorf("repository-id is required")
	}
	if params["localId"] == "" {
		return nil, fmt.Errorf("local-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	return params, nil
}

func printCRPatchSetList(out anyWriter, raw string) error {
	rows, ok := crPatchSetRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tCOMMIT\tDATE\tMESSAGE"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", row.ID, row.Commit, row.Date, row.Message)
	}
	return writer.Flush()
}

type crPatchSetRow struct {
	ID      string
	Commit  string
	Date    string
	Message string
}

func crPatchSetRowsFromJSON(raw string) []crPatchSetRow {
	rows, _ := crPatchSetRowsFromJSONForPrint(raw)
	return rows
}

func crPatchSetRowsFromJSONForPrint(raw string) ([]crPatchSetRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]crPatchSetRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, crPatchSetRow{
			ID:      firstStringValue(m, "id", "patchSetId", "patchSetBizId"),
			Commit:  firstStringValue(m, "commitId", "commitSha", "commit", "sha"),
			Date:    firstStringValue(m, "createdAt", "gmtCreated", "createdDate", "date"),
			Message: firstStringValue(m, "commitMessage", "message", "commitInfo", "description"),
		})
	}
	return rows, true
}

func repoChangeRequestRowsFromJSON(raw string) []repoChangeRequestRow {
	rows, _ := repoChangeRequestRowsFromJSONForPrint(raw)
	return rows
}

func repoChangeRequestRowsFromJSONForPrint(raw string) ([]repoChangeRequestRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]repoChangeRequestRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, repoChangeRequestRow{
			LocalID:      firstStringValue(m, "localId", "local_id", "iid", "id"),
			Repository:   firstStringValue(m, "repositoryId", "projectId", "repository", "project", "repositoryPath", "pathWithNamespace"),
			Title:        firstStringValue(m, "title", "subject", "name"),
			State:        firstStringValue(m, "state", "status", "mergeStatus"),
			Author:       firstStringValue(m, "authorName", "author_name", "author", "creatorName", "creator"),
			SourceBranch: firstStringValue(m, "sourceBranch", "source_branch", "sourceBranchName"),
			TargetBranch: firstStringValue(m, "targetBranch", "target_branch", "targetBranchName"),
		})
	}
	return rows, true
}
