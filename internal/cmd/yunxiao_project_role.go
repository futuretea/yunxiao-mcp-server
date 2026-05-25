package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type projectRoleListOptions struct {
	OrganizationID string
	ProjectID      string
	JSONOutput     bool
}

func newYunxiaoProjectRoleCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "role",
		Aliases: []string{"roles"},
		Short:   "work with Projex project roles",
	}
	command.AddCommand(newYunxiaoProjectRoleListCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoProjectRoleListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options projectRoleListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Projex project roles",
		Example: `  # List project roles
  yunxiao project role list --project-id 123

  # Output as JSON
  yunxiao project role list --project-id 123 --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_project_roles", params)
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printProjectRoleList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.ProjectID, "project-id", "", "Projex project ID")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o projectRoleListOptions) params() (map[string]any, error) {
	params := map[string]any{"projectId": strings.TrimSpace(o.ProjectID)}
	if params["projectId"] == "" {
		return nil, fmt.Errorf("project-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	return params, nil
}

func printProjectRoleList(out anyWriter, raw string) error {
	rows := projectRoleRowsFromJSON(raw)
	if len(rows) == 0 {
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

type projectRoleRow struct {
	ID          string
	Name        string
	Description string
}

func projectRoleRowsFromJSON(raw string) []projectRoleRow {
	items := rowsFromJSON(raw)
	rows := make([]projectRoleRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, projectRoleRow{
			ID:          firstStringValue(m, "id", "identifier", "roleId", "projectRoleId"),
			Name:        firstStringValue(m, "name", "displayName", "roleName", "projectRoleName"),
			Description: firstStringValue(m, "description", "desc"),
		})
	}
	return rows
}
