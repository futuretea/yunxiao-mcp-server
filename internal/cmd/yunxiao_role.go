package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type roleListOptions struct {
	OrganizationID string
	JSONOutput     bool
	OutputFormat string
}

func newYunxiaoRoleCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "role",
		Aliases: []string{"roles"},
		Short:   "work with Yunxiao organization roles",
	}
	command.AddCommand(newYunxiaoRoleListCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoRoleListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options roleListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Yunxiao organization roles",
		Example: `  # List organization roles
  yunxiao role list

  # Output as JSON
  yunxiao role list --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_organization_roles", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput || options.OutputFormat == "json" {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printRoleList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	flags.StringVar(&options.OutputFormat, "output", "", "output format: table, json, or csv")
	return command
}

func (o roleListOptions) params() map[string]any {
	params := map[string]any{}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	return params
}

func printRoleList(out anyWriter, raw string) error {
	rows, ok := roleRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tNAME\tDESCRIPTION"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\n", row.ID, row.Name, row.Description)
	}
	return writer.Flush()
}

type roleRow struct {
	ID          string
	Name        string
	Description string
}

func roleRowsFromJSONForPrint(raw string) ([]roleRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]roleRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, roleRow{
			ID:          firstStringValue(m, "id", "roleId"),
			Name:        firstStringValue(m, "name", "displayName", "roleName"),
			Description: firstStringValue(m, "description", "desc"),
		})
	}
	return rows, true
}
