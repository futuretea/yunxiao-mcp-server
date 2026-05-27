package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type repoListOptions struct {
	OrganizationID string
	Page           int
	PerPage        int
	OrderBy        string
	Sort           string
	Search         string
	Archived       bool
	ArchivedSet    bool
	JSONOutput     bool
	OutputFormat string
}

func newYunxiaoRepoCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "repo",
		Aliases: []string{"repos", "repository", "repositories"},
		Short:   "work with CodeUp repositories",
	}
	command.AddCommand(newYunxiaoRepoListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoRepoViewCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoRepoBranchCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoRepoChangeRequestCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoRepoCommitCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoRepoMrCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoRepoCompareCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoRepoFileCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoRepoListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options repoListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list CodeUp repositories",
		Example: `  # List repositories
  yunxiao repo list

  # Search by path keyword
  yunxiao repo list --search demo

  # Show only non-archived repos
  yunxiao repo list --archived=false

  # Output as JSON
  yunxiao repo list --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ArchivedSet = cmd.Flags().Changed("archived")
			result, err := callYunxiaoTool(cmd, cfg, "list_repositories", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput || options.OutputFormat == "json" {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printRepoList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.StringVar(&options.OrderBy, "order-by", "", "sort field, e.g. created_at, name, path, last_activity_at")
	flags.StringVar(&options.Sort, "sort", "", "sort direction, e.g. asc or desc")
	flags.StringVar(&options.Search, "search", "", "repository path search keyword")
	flags.BoolVar(&options.Archived, "archived", false, "filter archived repositories; use --archived=false for non-archived only")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	flags.StringVar(&options.OutputFormat, "output", "", "output format: table, json, or csv")
	return command
}

func (o repoListOptions) params() map[string]any {
	params := map[string]any{}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "orderBy", o.OrderBy)
	setCLIStringParam(params, "sort", o.Sort)
	setCLIStringParam(params, "search", o.Search)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	if o.ArchivedSet {
		params["archived"] = o.Archived
	}
	return params
}

func printRepoList(out anyWriter, raw string) error {
	rows := repoRowsFromJSON(raw)
	if len(rows) == 0 {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tNAME\tPATH\tARCHIVED\tLAST_ACTIVITY"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n", row.ID, row.Name, row.Path, row.Archived, row.LastActivity)
	}
	return writer.Flush()
}

type repoRow struct {
	ID           string
	Name         string
	Path         string
	Archived     string
	LastActivity string
}

func repoRowsFromJSON(raw string) []repoRow {
	items := rowsFromJSON(raw)
	rows := make([]repoRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, repoRow{
			ID:           firstStringValue(m, "id", "identifier", "repositoryId"),
			Name:         firstStringValue(m, "name", "displayName"),
			Path:         firstStringValue(m, "path", "pathWithNamespace", "fullPath", "sshUrlToRepo", "webUrl"),
			Archived:     firstStringValue(m, "archived", "isArchived"),
			LastActivity: firstStringValue(m, "lastActivityAt", "last_activity_at", "updatedAt", "gmtModified"),
		})
	}
	return rows
}
