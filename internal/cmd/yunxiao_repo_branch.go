package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type repoBranchListOptions struct {
	OrganizationID string
	RepositoryID   string
	Page           int
	PerPage        int
	Sort           string
	Search         string
	JSONOutput     bool
}

func newYunxiaoRepoBranchCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "branch",
		Aliases: []string{"branches"},
		Short:   "work with CodeUp repository branches",
	}
	command.AddCommand(newYunxiaoRepoBranchListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoRepoBranchViewCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoRepoBranchListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options repoBranchListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list CodeUp repository branches",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_branches", params)
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printRepoBranchList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.RepositoryID, "repository-id", "", "CodeUp repository ID or full path")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.StringVar(&options.Sort, "sort", "", "sort mode, e.g. name_asc, name_desc, updated_asc, updated_desc")
	flags.StringVar(&options.Search, "search", "", "branch search keyword")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o repoBranchListOptions) params() (map[string]any, error) {
	params := map[string]any{"repositoryId": strings.TrimSpace(o.RepositoryID)}
	if params["repositoryId"] == "" {
		return nil, fmt.Errorf("repository-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "sort", o.Sort)
	setCLIStringParam(params, "search", o.Search)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params, nil
}

func printRepoBranchList(out anyWriter, raw string) error {
	rows, ok := repoBranchRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, raw)
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, "NAME\tDEFAULT\tPROTECTED\tLAST_COMMIT")
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", row.Name, row.Default, row.Protected, row.LastCommit)
	}
	return writer.Flush()
}

type repoBranchRow struct {
	Name       string
	Default    string
	Protected  string
	LastCommit string
}

func repoBranchRowsFromJSON(raw string) []repoBranchRow {
	rows, _ := repoBranchRowsFromJSONForPrint(raw)
	return rows
}

func repoBranchRowsFromJSONForPrint(raw string) ([]repoBranchRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]repoBranchRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, repoBranchRow{
			Name:       firstStringValue(m, "name", "branchName"),
			Default:    firstStringValue(m, "default", "isDefault", "defaultBranch"),
			Protected:  firstStringValue(m, "protected", "isProtected"),
			LastCommit: repoBranchLastCommitValue(m),
		})
	}
	return rows, true
}

func repoBranchLastCommitValue(m map[string]any) string {
	for _, key := range []string{"commit", "latestCommit", "lastCommit"} {
		if nested, ok := m[key].(map[string]any); ok {
			if value := firstStringValue(nested, "id", "sha", "commitId", "shortId"); value != "" {
				return value
			}
		}
	}
	return firstStringValue(m, "commitId", "sha")
}
