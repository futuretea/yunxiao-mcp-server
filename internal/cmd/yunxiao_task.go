package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"text/tabwriter"
)

type taskListOptions struct {
	OrganizationID string
	ProjectID      string
	Category       string
	Subject        string
	Status         string
	AssignedTo     string
	Creator        string
	Sprint         string
	OrderBy        string
	Sort           string
	Page           int
	PerPage        int
	JSONOutput     bool
	OutputFormat   string
}

func (o taskListOptions) params() (map[string]any, error) {
	params := map[string]any{
		"category":  strings.TrimSpace(o.Category),
		"projectId": strings.TrimSpace(o.ProjectID),
	}
	if params["projectId"] == "" {
		return nil, fmt.Errorf("project-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "subject", o.Subject)
	setCLIStringParam(params, "status", o.Status)
	setCLIStringParam(params, "assignedTo", o.AssignedTo)
	setCLIStringParam(params, "creator", o.Creator)
	setCLIStringParam(params, "sprint", o.Sprint)
	setCLIStringParam(params, "orderBy", o.OrderBy)
	setCLIStringParam(params, "sort", o.Sort)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params, nil
}

func printTaskList(out anyWriter, raw string) error {
	rows := taskRowsFromJSON(raw)
	if len(rows) == 0 {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tSUBJECT\tSTATUS\tASSIGNEE"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", row.ID, row.Subject, row.Status, row.Assignee)
	}
	return writer.Flush()
}

type taskRow struct {
	ID       string
	Subject  string
	Status   string
	Assignee string
}

func taskRowsFromJSON(raw string) []taskRow {
	items := rowsFromJSON(raw)
	rows := make([]taskRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, taskRow{
			ID:       firstStringValue(m, "id", "identifier", "workitemId", "workItemId"),
			Subject:  firstStringValue(m, "subject", "title", "name"),
			Status:   firstStringValue(m, "status", "statusName", "statusIdentifier"),
			Assignee: firstStringValue(m, "assignedTo", "assignee", "assignedToName", "assigneeName"),
		})
	}
	return rows
}

type taskTypeListOptions struct {
	OrganizationID string
	ProjectID      string
	Category       string
	JSONOutput     bool
	OutputFormat   string
}

func newYunxiaoTaskTypeListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options taskTypeListOptions
	command := &cobra.Command{
		Use:     "type-list",
		Aliases: []string{"types"},
		Short:   "list work item types in a project",
		Example: `  # List task types
  yunxiao task type-list --project-id 123 --category Task

  # List bug types
  yunxiao task types --project-id 123 --category Bug

  # Output as JSON
  yunxiao task type-list --project-id 123 --category Task --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_work_item_types", params)
			if err != nil {
				return err
			}
			if options.JSONOutput || options.OutputFormat == "json" {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printTaskTypeList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.ProjectID, "project-id", "", "Projex project ID")
	flags.StringVar(&options.Category, "category", "", "work item category, e.g. Task, Bug, Req, Risk")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	flags.StringVar(&options.OutputFormat, "output", "", "output format: table, json, or csv")
	return command
}

func (o taskTypeListOptions) params() (map[string]any, error) {
	params := map[string]any{
		"projectId": o.ProjectID,
		"category":  o.Category,
	}
	if params["projectId"] == "" {
		return nil, fmt.Errorf("project-id is required")
	}
	if params["category"] == "" {
		return nil, fmt.Errorf("category is required, e.g. Task, Bug, Req, Risk")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	return params, nil
}

func printTaskTypeList(out anyWriter, raw string) error {
	rows, ok := taskTypeRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tNAME\tCATEGORY"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\n", row.ID, row.Name, row.Category)
	}
	return writer.Flush()
}

type taskTypeRow struct {
	ID       string
	Name     string
	Category string
}

func taskTypeRowsFromJSONForPrint(raw string) ([]taskTypeRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]taskTypeRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, taskTypeRow{
			ID:       firstStringValue(m, "id", "workItemTypeId"),
			Name:     firstStringValue(m, "name", "displayName"),
			Category: firstStringValue(m, "category"),
		})
	}
	return rows, true
}

type taskTimelineOptions struct {
	OrganizationID  string
	IncludeWorkitem bool
}

func newYunxiaoTaskTimelineCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	options := taskTimelineOptions{IncludeWorkitem: true}
	command := &cobra.Command{
		Use:     "timeline <workitem-id>",
		Aliases: []string{"history"},
		Short:   "show status change timeline for a work item as JSON",
		Example: `  # View task timeline
  yunxiao task timeline wi-12345`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params := options.params(args[0])
			result, err := callYunxiaoTool(cmd, cfg, "get_workitem_status_timeline", params)
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.BoolVar(&options.IncludeWorkitem, "include-workitem", true, "include basic work item info")
	return command
}

func (o taskTimelineOptions) params(workitemID string) map[string]any {
	params := map[string]any{
		"workitemId":      workitemID,
		"includeWorkitem": o.IncludeWorkitem,
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	return params
}

type taskMyOptions struct {
	OrganizationID string
	ProjectID      string
	Relation       string
	Status         string
	Categories     string
	SampleLimit    int
}

func newYunxiaoTaskMyCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options taskMyOptions
	command := &cobra.Command{
		Use:     "my <project-id>",
		Aliases: []string{"mine"},
		Short:   "show my work items in a project as JSON",
		Example: `  # Show my assigned tasks
  yunxiao task my 123

  # Show tasks I created
  yunxiao task my 123 --relation created

  # Filter by status
  yunxiao task my 123 --status "处理中"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ProjectID = args[0]
			result, err := callYunxiaoTool(cmd, cfg, "get_my_project_workitems", options.params())
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.Relation, "relation", "", "filter relation: assigned or created; defaults to assigned")
	flags.StringVar(&options.Status, "status", "", "comma-separated status IDs")
	flags.StringVar(&options.Categories, "categories", "", "comma-separated categories; defaults to Task,Bug")
	flags.IntVar(&options.SampleLimit, "sample-limit", 0, "samples per category")
	return command
}

func (o taskMyOptions) params() map[string]any {
	params := map[string]any{"projectId": o.ProjectID}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "relation", o.Relation)
	setCLIStringParam(params, "status", o.Status)
	setCLIStringParam(params, "categories", o.Categories)
	if o.SampleLimit > 0 {
		params["sampleLimit"] = o.SampleLimit
	}
	return params
}

type taskTypeViewOptions struct {
	OrganizationID     string
	ProjectID          string
	WorkItemTypeID     string
	IncludeFieldConfig bool
	IncludeWorkflow    bool
}

func newYunxiaoTaskTypeViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	options := taskTypeViewOptions{
		IncludeFieldConfig: true,
		IncludeWorkflow:    true,
	}
	command := &cobra.Command{
		Use:     "type-view <type-id>",
		Aliases: []string{"type-info"},
		Short:   "show work item type details as JSON",
		Example: `  # View work item type details
  yunxiao task type-view 456 --project-id 123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.WorkItemTypeID = args[0]
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "get_work_item_type_overview", params)
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.ProjectID, "project-id", "", "Projex project ID")
	flags.BoolVar(&options.IncludeFieldConfig, "include-field-config", true, "include field configuration")
	flags.BoolVar(&options.IncludeWorkflow, "include-workflow", true, "include workflow metadata")
	return command
}

func (o taskTypeViewOptions) params() (map[string]any, error) {
	params := map[string]any{
		"projectId":          o.ProjectID,
		"workItemTypeId":     o.WorkItemTypeID,
		"includeFieldConfig": o.IncludeFieldConfig,
		"includeWorkflow":    o.IncludeWorkflow,
	}
	if params["projectId"] == "" {
		return nil, fmt.Errorf("project-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	return params, nil
}

type taskTypeAllOptions struct {
	OrganizationID string
	Categories     string
	JSONOutput     bool
	OutputFormat   string
}

func newYunxiaoTaskTypeAllCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options taskTypeAllOptions
	command := &cobra.Command{
		Use:     "type-all",
		Aliases: []string{"all-types"},
		Short:   "list all work item types across the organization",
		Example: `  # List all task types
  yunxiao task type-all --categories Task

  # Output as JSON
  yunxiao task type-all --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_all_work_item_types", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput || options.OutputFormat == "json" {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printTaskTypeList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.Categories, "categories", "", "comma-separated categories, e.g. Req,Bug,Task,Risk")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	flags.StringVar(&options.OutputFormat, "output", "", "output format: table, json, or csv")
	return command
}

func (o taskTypeAllOptions) params() map[string]any {
	params := map[string]any{}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "categories", o.Categories)
	return params
}

type taskRelationTypesOptions struct {
	OrganizationID string
	WorkItemTypeID string
	RelationType   string
	JSONOutput     bool
	OutputFormat   string
}

func newYunxiaoTaskRelationTypesCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options taskRelationTypesOptions
	command := &cobra.Command{
		Use:     "relation-types <type-id>",
		Aliases: []string{"related-types"},
		Short:   "list work item types that can be related to a given type",
		Example: `  # List types that can be related
  yunxiao task relation-types 456

  # Filter by relation type
  yunxiao task relation-types 456 --relation-type PARENT

  # Output as JSON
  yunxiao task relation-types 456 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.WorkItemTypeID = args[0]
			result, err := callYunxiaoTool(cmd, cfg, "list_work_item_relation_work_item_types", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput || options.OutputFormat == "json" {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printTaskTypeList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.RelationType, "relation-type", "", "relation type: PARENT, SUB, ASSOCIATED, DEPEND_ON, DEPENDED_BY")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	flags.StringVar(&options.OutputFormat, "output", "", "output format: table, json, or csv")
	return command
}

func (o taskRelationTypesOptions) params() map[string]any {
	params := map[string]any{"workItemTypeId": o.WorkItemTypeID}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "relationType", o.RelationType)
	return params
}
