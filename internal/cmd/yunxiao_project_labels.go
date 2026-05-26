package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
type projectLabelsOptions struct {
	OrganizationID string
	ProjectID      string
	Page           int
	PerPage        int
	JSONOutput     bool
}

func newYunxiaoProjectLabelsCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options projectLabelsOptions
	command := &cobra.Command{
		Use:     "labels <project-id>",
		Aliases: []string{"tags"},
		Short:   "list Projex project labels",
		Example: `  # List project labels
  yunxiao project labels 123

  # Output as JSON
  yunxiao project labels 123 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ProjectID = args[0]
			result, err := callYunxiaoTool(cmd, cfg, "list_labels", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printLabelList(streams.Out, result)
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

func (o projectLabelsOptions) params() map[string]any {
	params := map[string]any{"projectId": o.ProjectID}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params
}

func printLabelList(out anyWriter, raw string) error {
	rows, ok := labelRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tNAME\tCOLOR"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\n", row.ID, row.Name, row.Color)
	}
	return writer.Flush()
}

type labelRow struct {
	ID    string
	Name  string
	Color string
}

func labelRowsFromJSONForPrint(raw string) ([]labelRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]labelRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, labelRow{
			ID:    firstStringValue(m, "id", "labelId"),
			Name:  firstStringValue(m, "name", "displayName"),
			Color: firstStringValue(m, "color", "colour", "backgroundColor", "foregroundColor"),
		})
	}
	return rows, true
}

type projectMilestonesOptions struct {
	OrganizationID string
	ProjectID      string
	Status         string
	Page           int
	PerPage        int
	JSONOutput     bool
}

func newYunxiaoProjectMilestonesCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options projectMilestonesOptions
	command := &cobra.Command{
		Use:     "milestones <project-id>",
		Aliases: []string{"milestone"},
		Short:   "list Projex project milestones",
		Example: `  # List milestones
  yunxiao project milestones 123

  # Output as JSON
  yunxiao project milestones 123 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ProjectID = args[0]
			result, err := callYunxiaoTool(cmd, cfg, "list_milestones", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printMilestoneList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.Status, "status", "", "comma-separated milestone status IDs")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o projectMilestonesOptions) params() map[string]any {
	params := map[string]any{"projectId": o.ProjectID}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "status", o.Status)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params
}

func printMilestoneList(out anyWriter, raw string) error {
	rows, ok := milestoneRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tNAME\tSTATUS\tDUE_DATE"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", row.ID, row.Name, row.Status, row.DueDate)
	}
	return writer.Flush()
}

type milestoneRow struct {
	ID      string
	Name    string
	Status  string
	DueDate string
}

func milestoneRowsFromJSONForPrint(raw string) ([]milestoneRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]milestoneRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, milestoneRow{
			ID:      firstStringValue(m, "id", "milestoneId"),
			Name:    firstStringValue(m, "name", "displayName", "title"),
			Status:  firstStringValue(m, "status", "statusName"),
			DueDate: firstStringValue(m, "dueDate", "endDate", "deadline"),
		})
	}
	return rows, true
}

type projectMemberTasksOptions struct {
	OrganizationID string
	ProjectID      string
	AssigneeIDs    string
	Categories     string
	Status         string
	MemberLimit    int
	SampleLimit    int
}

func newYunxiaoProjectMemberTasksCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options projectMemberTasksOptions
	command := &cobra.Command{
		Use:     "member-tasks <project-id>",
		Aliases: []string{"member-status"},
		Short:   "show per-member task status as JSON",
		Example: `  # View per-member task status
  yunxiao project member-tasks 123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ProjectID = args[0]
			result, err := callYunxiaoTool(cmd, cfg, "get_project_member_task_status", options.params())
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.AssigneeIDs, "assignee-ids", "", "comma-separated assignee user IDs")
	flags.StringVar(&options.Categories, "categories", "", "comma-separated categories; defaults to Task,Bug")
	flags.StringVar(&options.Status, "status", "", "comma-separated status IDs")
	flags.IntVar(&options.MemberLimit, "member-limit", 0, "max members to inspect")
	flags.IntVar(&options.SampleLimit, "sample-limit", 0, "samples per member section")
	return command
}

func (o projectMemberTasksOptions) params() map[string]any {
	params := map[string]any{"projectId": o.ProjectID}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "assigneeIds", o.AssigneeIDs)
	setCLIStringParam(params, "categories", o.Categories)
	setCLIStringParam(params, "status", o.Status)
	if o.MemberLimit > 0 {
		params["memberLimit"] = o.MemberLimit
	}
	if o.SampleLimit > 0 {
		params["sampleLimit"] = o.SampleLimit
	}
	return params
}
