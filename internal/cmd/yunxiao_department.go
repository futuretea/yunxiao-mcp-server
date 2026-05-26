package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type departmentListOptions struct {
	OrganizationID string
	ParentID       string
	Page           int
	PerPage        int
	JSONOutput     bool
}

func newYunxiaoDepartmentCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "department",
		Aliases: []string{"departments", "dept"},
		Short:   "work with Yunxiao organization departments",
	}
	command.AddCommand(newYunxiaoDepartmentListCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoDepartmentListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options departmentListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Yunxiao organization departments",
		Example: `  # List top-level departments
  yunxiao department list

  # List child departments
  yunxiao department list --parent-id dept-123

  # Output as JSON
  yunxiao department list --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_organization_departments", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printDepartmentList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.ParentID, "parent-id", "", "parent department ID for child listing")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o departmentListOptions) params() map[string]any {
	params := map[string]any{}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "parentId", o.ParentID)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params
}

func printDepartmentList(out anyWriter, raw string) error {
	rows, ok := departmentRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tNAME\tPARENT_ID"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\n", row.ID, row.Name, row.ParentID)
	}
	return writer.Flush()
}

type departmentRow struct {
	ID       string
	Name     string
	ParentID string
}

func departmentRowsFromJSONForPrint(raw string) ([]departmentRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]departmentRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, departmentRow{
			ID:       firstStringValue(m, "id", "departmentId"),
			Name:     firstStringValue(m, "name", "displayName", "departmentName"),
			ParentID: firstStringValue(m, "parentId", "parentDepartmentId"),
		})
	}
	return rows, true
}
