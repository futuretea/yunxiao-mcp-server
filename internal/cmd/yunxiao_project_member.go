package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type projectMemberListOptions struct {
	OrganizationID string
	ProjectID      string
	Name           string
	RoleID         string
	JSONOutput     bool
	OutputFormat   string
}

func newYunxiaoProjectMemberCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "member",
		Aliases: []string{"members"},
		Short:   "work with Projex project members",
	}
	command.AddCommand(newYunxiaoProjectMemberListCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoProjectMemberListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options projectMemberListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Projex project members",
		Example: `  # List project members
  yunxiao project member list --project-id 123

  # Filter by name
  yunxiao project member list --project-id 123 --name alice

  # Output as JSON
  yunxiao project member list --project-id 123 --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_project_members", params)
			if err != nil {
				return err
			}
			if options.JSONOutput || options.OutputFormat == "json" {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printProjectMemberList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.ProjectID, "project-id", "", "Projex project ID")
	flags.StringVar(&options.Name, "name", "", "project member name keyword")
	flags.StringVar(&options.RoleID, "role-id", "", "project role ID")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	flags.StringVar(&options.OutputFormat, "output", "", "output format: table, json, or csv")
	return command
}

func (o projectMemberListOptions) params() (map[string]any, error) {
	params := map[string]any{"projectId": strings.TrimSpace(o.ProjectID)}
	if params["projectId"] == "" {
		return nil, fmt.Errorf("project-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "name", o.Name)
	setCLIStringParam(params, "roleId", o.RoleID)
	return params, nil
}

func printProjectMemberList(out anyWriter, raw string) error {
	rows := projectMemberRowsFromJSON(raw)
	if len(rows) == 0 {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("USER_ID\tNAME\tROLE_ID\tROLE"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", row.UserID, row.Name, row.RoleID, row.RoleName)
	}
	return writer.Flush()
}

type projectMemberRow struct {
	UserID   string
	Name     string
	RoleID   string
	RoleName string
}

func projectMemberRowsFromJSON(raw string) []projectMemberRow {
	items := rowsFromJSON(raw)
	rows := make([]projectMemberRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, projectMemberRow{
			UserID:   firstStringValue(m, "userId", "accountId", "id", "identifier"),
			Name:     firstStringValue(m, "userName", "displayName", "name", "username", "nickName", "realName"),
			RoleID:   firstStringValue(m, "roleId", "projectRoleId"),
			RoleName: firstStringValue(m, "roleName", "projectRoleName", "role", "projectRole"),
		})
	}
	return rows
}
