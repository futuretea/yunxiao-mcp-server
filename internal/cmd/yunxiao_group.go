package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type groupListOptions struct {
	OrganizationID string
	Page           int
	PerPage        int
	JSONOutput     bool
}

func newYunxiaoGroupCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "group",
		Aliases: []string{"groups"},
		Short:   "work with Yunxiao organization groups",
	}
	command.AddCommand(newYunxiaoGroupListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoGroupViewCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoGroupListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options groupListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Yunxiao organization groups",
		Example: `  # List groups
  yunxiao group list

  # Output as JSON
  yunxiao group list --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_organization_groups", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printGroupList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o groupListOptions) params() map[string]any {
	params := map[string]any{}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params
}

func printGroupList(out anyWriter, raw string) error {
	rows, ok := groupRowsFromJSONForPrint(raw)
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

type groupRow struct {
	ID          string
	Name        string
	Description string
}

func groupRowsFromJSONForPrint(raw string) ([]groupRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]groupRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, groupRow{
			ID:          firstStringValue(m, "id", "groupId"),
			Name:        firstStringValue(m, "name", "displayName"),
			Description: firstStringValue(m, "description", "desc"),
		})
	}
	return rows, true
}

type groupViewOptions struct {
	OrganizationID string
	GroupID        string
	IncludeMembers bool
	MemberLimit    int
}

func newYunxiaoGroupViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	options := groupViewOptions{IncludeMembers: true}
	command := &cobra.Command{
		Use:     "view <group-id>",
		Aliases: []string{"overview"},
		Short:   "view a group overview as JSON",
		Example: `  # View group
  yunxiao group view group-123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.GroupID = args[0]
			result, err := callYunxiaoTool(cmd, cfg, "get_organization_group_overview", options.params())
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.BoolVar(&options.IncludeMembers, "include-members", true, "include group members")
	flags.IntVar(&options.MemberLimit, "member-limit", 0, "max members to include")
	return command
}

func (o groupViewOptions) params() map[string]any {
	params := map[string]any{
		"groupId":        o.GroupID,
		"includeMembers": o.IncludeMembers,
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	if o.MemberLimit > 0 {
		params["memberLimit"] = o.MemberLimit
	}
	return params
}
