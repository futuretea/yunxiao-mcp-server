package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type organizationListOptions struct {
	Page         int
	PerPage      int
	JSONOutput   bool
	OutputFormat string
}

func newYunxiaoOrganizationCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "organization",
		Aliases: []string{"organizations"},
		Short:   "work with Yunxiao organizations",
	}
	command.AddCommand(newYunxiaoOrganizationListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoOrganizationViewCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoOrganizationInfoCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoOrganizationListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options organizationListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Yunxiao organizations visible to the current user",
		Example: `  # List all organizations
  yunxiao organization list

  # Output as JSON
  yunxiao organization list --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_organizations", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput || options.OutputFormat == "json" {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printOrganizationList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	flags.StringVar(&options.OutputFormat, "output", "", "output format: table, json, or csv")
	return command
}

func (o organizationListOptions) params() map[string]any {
	params := map[string]any{}
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params
}

func printOrganizationList(out anyWriter, raw string) error {
	rows := organizationRowsFromJSON(raw)
	if len(rows) == 0 {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tNAME"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\n", row.ID, row.Name)
	}
	return writer.Flush()
}

type organizationRow struct {
	ID   string
	Name string
}

func organizationRowsFromJSON(raw string) []organizationRow {
	items := rowsFromJSON(raw)
	rows := make([]organizationRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, organizationRow{
			ID:   firstStringValue(m, "id", "identifier", "organizationId"),
			Name: firstStringValue(m, "name", "displayName", "title"),
		})
	}
	return rows
}

type orgViewOptions struct {
	OrganizationID     string
	IncludeDepartments bool
	IncludeMembers     bool
	IncludeGroups      bool
	IncludeRoles       bool
	DepartmentLimit    int
	MemberLimit        int
	GroupLimit         int
}

func newYunxiaoOrganizationViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	options := orgViewOptions{
		IncludeDepartments: true,
		IncludeMembers:     true,
		IncludeGroups:      true,
		IncludeRoles:       true,
	}
	command := &cobra.Command{
		Use:     "view",
		Aliases: []string{"overview", "info"},
		Short:   "view a Yunxiao organization overview as JSON",
		Example: `  # View default organization
  yunxiao organization view

  # View specific organization
  yunxiao organization view --organization-id org-abc

  # View without member details
  yunxiao organization view --include-members=false`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "get_organization_overview", options.params())
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.BoolVar(&options.IncludeDepartments, "include-departments", true, "include departments list")
	flags.BoolVar(&options.IncludeMembers, "include-members", true, "include members list")
	flags.BoolVar(&options.IncludeGroups, "include-groups", true, "include groups list")
	flags.BoolVar(&options.IncludeRoles, "include-roles", true, "include roles list")
	flags.IntVar(&options.DepartmentLimit, "department-limit", 0, "max departments to include")
	flags.IntVar(&options.MemberLimit, "member-limit", 0, "max members to include")
	flags.IntVar(&options.GroupLimit, "group-limit", 0, "max groups to include")
	return command
}

func (o orgViewOptions) params() map[string]any {
	params := map[string]any{
		"includeDepartments": o.IncludeDepartments,
		"includeMembers":     o.IncludeMembers,
		"includeGroups":      o.IncludeGroups,
		"includeRoles":       o.IncludeRoles,
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	if o.DepartmentLimit > 0 {
		params["departmentLimit"] = o.DepartmentLimit
	}
	if o.MemberLimit > 0 {
		params["memberLimit"] = o.MemberLimit
	}
	if o.GroupLimit > 0 {
		params["groupLimit"] = o.GroupLimit
	}
	return params
}

func newYunxiaoOrganizationInfoCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:   "info",
		Short: "show current default organization info as JSON",
		Example: `  # Show default organization
  yunxiao organization info`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "get_current_organization_info", map[string]any{})
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	return command
}
