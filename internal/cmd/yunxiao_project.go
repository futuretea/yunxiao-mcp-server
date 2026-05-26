package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type projectListOptions struct {
	OrganizationID string
	Name           string
	Status         string
	Creator        string
	OrderBy        string
	Sort           string
	Page           int
	PerPage        int
	JSONOutput     bool
}

func newYunxiaoProjectCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "project",
		Aliases: []string{"projects"},
		Short:   "work with Projex projects",
	}
	command.AddCommand(newYunxiaoProjectListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoProjectMemberCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoProjectRoleCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoProjectViewCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoProjectSummaryCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoProjectContextCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoProjectRiskCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoProjectBlockersCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoProjectWorkloadCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoProjectBoardCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoProjectLabelsCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoProjectMilestonesCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoProjectMemberTasksCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoProjectListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options projectListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Projex projects",
		Example: `  # List all projects
  yunxiao project list

  # Search by name keyword
  yunxiao project list --name demo

  # Output as JSON
  yunxiao project list --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "search_projects", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printProjectList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.Name, "name", "", "project name keyword")
	flags.StringVar(&options.Status, "status", "", "comma-separated project status IDs")
	flags.StringVar(&options.Creator, "creator", "", "comma-separated creator user IDs")
	flags.StringVar(&options.OrderBy, "order-by", "", "sort field")
	flags.StringVar(&options.Sort, "sort", "", "sort direction, e.g. asc or desc")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o projectListOptions) params() map[string]any {
	params := map[string]any{}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "name", o.Name)
	setCLIStringParam(params, "status", o.Status)
	setCLIStringParam(params, "creator", o.Creator)
	setCLIStringParam(params, "orderBy", o.OrderBy)
	setCLIStringParam(params, "sort", o.Sort)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params
}

func printProjectList(out anyWriter, raw string) error {
	rows := projectRowsFromJSON(raw)
	if len(rows) == 0 {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tNAME\tSTATUS\tCREATOR"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", row.ID, row.Name, row.Status, row.Creator)
	}
	return writer.Flush()
}

type projectRow struct {
	ID      string
	Name    string
	Status  string
	Creator string
}

func projectRowsFromJSON(raw string) []projectRow {
	items := rowsFromJSON(raw)
	rows := make([]projectRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, projectRow{
			ID:      firstStringValue(m, "id", "identifier", "projectId", "spaceId"),
			Name:    firstStringValue(m, "name", "title"),
			Status:  firstStringValue(m, "status", "statusName", "statusIdentifier"),
			Creator: firstStringValue(m, "creator", "owner", "creatorName", "ownerName"),
		})
	}
	return rows
}

type projectViewOptions struct {
	OrganizationID    string
	ProjectID         string
	IncludeMembers    bool
	IncludeSprints    bool
	IncludeMilestones bool
	IncludeVersions   bool
	IncludeLabels     bool
	ActiveOnly        bool
	Status            string
	Page              int
	PerPage           int
}

func newYunxiaoProjectViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	options := projectViewOptions{
		IncludeMembers:    true,
		IncludeSprints:    true,
		IncludeMilestones: true,
		IncludeVersions:   true,
		IncludeLabels:     true,
		ActiveOnly:        true,
	}
	command := &cobra.Command{
		Use:     "view <project-id>",
		Aliases: []string{"overview", "info"},
		Short:   "view a Projex project overview as JSON",
		Example: `  # View project overview
  yunxiao project view 123

  # View without sprint details
  yunxiao project view 123 --include-sprints=false`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ProjectID = args[0]
			result, err := callYunxiaoTool(cmd, cfg, "get_project_overview", options.params())
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.BoolVar(&options.IncludeMembers, "include-members", true, "include project members")
	flags.BoolVar(&options.IncludeSprints, "include-sprints", true, "include sprints list")
	flags.BoolVar(&options.IncludeMilestones, "include-milestones", true, "include milestones list")
	flags.BoolVar(&options.IncludeVersions, "include-versions", true, "include versions list")
	flags.BoolVar(&options.IncludeLabels, "include-labels", true, "include labels list")
	flags.BoolVar(&options.ActiveOnly, "active-only", true, "show only active sprints, milestones, and versions")
	flags.StringVar(&options.Status, "status", "", "comma-separated statuses for active-only filter; defaults to TODO,DOING")
	flags.IntVar(&options.Page, "page", 0, "page number for paged sections")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size for paged sections")
	return command
}

func (o projectViewOptions) params() map[string]any {
	params := map[string]any{
		"projectId":         o.ProjectID,
		"includeMembers":    o.IncludeMembers,
		"includeSprints":    o.IncludeSprints,
		"includeMilestones": o.IncludeMilestones,
		"includeVersions":   o.IncludeVersions,
		"includeLabels":     o.IncludeLabels,
		"activeOnly":        o.ActiveOnly,
	}
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

type projectSummaryOptions struct {
	OrganizationID string
	ProjectID      string
	Categories     string
	Subject        string
	Status         string
	AssignedTo     string
	Creator        string
	Tag            string
	OrderBy        string
	Sort           string
	SampleLimit    int
}

func newYunxiaoProjectSummaryCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options projectSummaryOptions
	command := &cobra.Command{
		Use:     "summary <project-id>",
		Aliases: []string{"stats"},
		Short:   "summarize work items by category for a project as JSON",
		Example: `  # Summarize work items
  yunxiao project summary 123

  # Filter by specific categories
  yunxiao project summary 123 --categories "Task,Bug"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ProjectID = args[0]
			result, err := callYunxiaoTool(cmd, cfg, "get_project_workitem_summary", options.params())
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.Categories, "categories", "", "comma-separated categories; defaults to Req,Task,Bug,Risk")
	flags.StringVar(&options.Subject, "subject", "", "subject/title keyword")
	flags.StringVar(&options.Status, "status", "", "comma-separated status IDs")
	flags.StringVar(&options.AssignedTo, "assigned-to", "", "comma-separated assignee user IDs")
	flags.StringVar(&options.Creator, "creator", "", "comma-separated creator user IDs")
	flags.StringVar(&options.Tag, "tag", "", "comma-separated tag IDs")
	flags.StringVar(&options.OrderBy, "order-by", "", "sort field")
	flags.StringVar(&options.Sort, "sort", "", "sort direction, e.g. asc or desc")
	flags.IntVar(&options.SampleLimit, "sample-limit", 0, "samples returned per category")
	return command
}

func (o projectSummaryOptions) params() map[string]any {
	params := map[string]any{"projectId": o.ProjectID}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "categories", o.Categories)
	setCLIStringParam(params, "subject", o.Subject)
	setCLIStringParam(params, "status", o.Status)
	setCLIStringParam(params, "assignedTo", o.AssignedTo)
	setCLIStringParam(params, "creator", o.Creator)
	setCLIStringParam(params, "tag", o.Tag)
	setCLIStringParam(params, "orderBy", o.OrderBy)
	setCLIStringParam(params, "sort", o.Sort)
	if o.SampleLimit > 0 {
		params["sampleLimit"] = o.SampleLimit
	}
	return params
}

type projectContextOptions struct {
	OrganizationID  string
	ProjectID       string
	Category        string
	WorkItemTypeID  string
	IncludeMembers  bool
	IncludeLabels   bool
	IncludeFields   bool
	IncludeWorkflow bool
	Page            int
	PerPage         int
}

func newYunxiaoProjectContextCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	options := projectContextOptions{
		IncludeMembers:  true,
		IncludeLabels:   true,
		IncludeFields:   true,
		IncludeWorkflow: true,
	}
	command := &cobra.Command{
		Use:     "context <project-id>",
		Aliases: []string{"ctx", "meta"},
		Short:   "get project work item metadata context as JSON",
		Example: `  # Get context for tasks
  yunxiao project context 123 --category Task

  # Get context for a specific work item type
  yunxiao project ctx 123 --category Task --work-item-type-id 456`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ProjectID = args[0]
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "get_project_workitem_context", params)
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.Category, "category", "", "work item category, e.g. Task, Bug, Req, Risk")
	flags.StringVar(&options.WorkItemTypeID, "work-item-type-id", "", "optional work item type ID for field/workflow metadata")
	flags.BoolVar(&options.IncludeMembers, "include-members", true, "include project members")
	flags.BoolVar(&options.IncludeLabels, "include-labels", true, "include project labels")
	flags.BoolVar(&options.IncludeFields, "include-fields", true, "include field configuration when type is set")
	flags.BoolVar(&options.IncludeWorkflow, "include-workflow", true, "include workflow metadata when type is set")
	flags.IntVar(&options.Page, "page", 0, "labels page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "labels page size")
	return command
}

func (o projectContextOptions) params() (map[string]any, error) {
	params := map[string]any{
		"projectId":        o.ProjectID,
		"category":         o.Category,
		"includeMembers":   o.IncludeMembers,
		"includeLabels":    o.IncludeLabels,
		"includeFields":    o.IncludeFields,
		"includeWorkflow":  o.IncludeWorkflow,
	}
	if params["category"] == "" {
		return nil, fmt.Errorf("category is required, e.g. Task, Bug, Req, Risk")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "workItemTypeId", o.WorkItemTypeID)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params, nil
}

type projectRiskOptions struct {
	OrganizationID string
	ProjectID      string
	Categories     string
	Subject        string
	Status         string
	StatusStage    string
	AssignedTo     string
	Creator        string
	Sprint         string
	WorkItemType   string
	Tag            string
	SampleLimit    int
}

func newYunxiaoProjectRiskCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options projectRiskOptions
	command := &cobra.Command{
		Use:     "risk <project-id>",
		Aliases: []string{"health"},
		Short:   "show project risk dashboard as JSON",
		Example: `  # View risk dashboard
  yunxiao project risk 123

  # Focus on specific categories
  yunxiao project risk 123 --categories "Bug,Risk"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ProjectID = args[0]
			result, err := callYunxiaoTool(cmd, cfg, "get_project_risk_dashboard", options.params())
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.Categories, "categories", "", "comma-separated categories; defaults to Risk,Bug,Task")
	flags.StringVar(&options.Subject, "subject", "", "subject/title keyword")
	flags.StringVar(&options.Status, "status", "", "comma-separated active status IDs")
	flags.StringVar(&options.StatusStage, "status-stage", "", "comma-separated status stage IDs")
	flags.StringVar(&options.AssignedTo, "assigned-to", "", "comma-separated assignee user IDs")
	flags.StringVar(&options.Creator, "creator", "", "comma-separated creator user IDs")
	flags.StringVar(&options.Sprint, "sprint", "", "comma-separated sprint IDs")
	flags.StringVar(&options.WorkItemType, "workitem-type", "", "comma-separated work item type IDs")
	flags.StringVar(&options.Tag, "tag", "", "comma-separated tag IDs")
	flags.IntVar(&options.SampleLimit, "sample-limit", 0, "samples per section")
	return command
}

func (o projectRiskOptions) params() map[string]any {
	params := map[string]any{"projectId": o.ProjectID}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "categories", o.Categories)
	setCLIStringParam(params, "subject", o.Subject)
	setCLIStringParam(params, "status", o.Status)
	setCLIStringParam(params, "statusStage", o.StatusStage)
	setCLIStringParam(params, "assignedTo", o.AssignedTo)
	setCLIStringParam(params, "creator", o.Creator)
	setCLIStringParam(params, "sprint", o.Sprint)
	setCLIStringParam(params, "workitemType", o.WorkItemType)
	setCLIStringParam(params, "tag", o.Tag)
	if o.SampleLimit > 0 {
		params["sampleLimit"] = o.SampleLimit
	}
	return params
}

type projectBlockersOptions struct {
	OrganizationID string
	ProjectID      string
	Categories     string
	SampleLimit    int
}

func newYunxiaoProjectBlockersCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options projectBlockersOptions
	command := &cobra.Command{
		Use:     "blockers <project-id>",
		Aliases: []string{"dependencies"},
		Short:   "show dependency blocker analysis as JSON",
		Example: `  # View blocker analysis
  yunxiao project blockers 123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ProjectID = args[0]
			result, err := callYunxiaoTool(cmd, cfg, "get_blocker_analysis", options.params())
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.Categories, "categories", "", "comma-separated categories; defaults to Task,Bug")
	flags.IntVar(&options.SampleLimit, "sample-limit", 0, "max work items per category")
	return command
}

func (o projectBlockersOptions) params() map[string]any {
	params := map[string]any{"projectId": o.ProjectID}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "categories", o.Categories)
	if o.SampleLimit > 0 {
		params["sampleLimit"] = o.SampleLimit
	}
	return params
}

type projectWorkloadOptions struct {
	OrganizationID string
	ProjectID      string
	AssigneeIDs    string
	Categories     string
	Status         string
	MemberLimit    int
	TaskLimit      int
}

func newYunxiaoProjectWorkloadCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options projectWorkloadOptions
	command := &cobra.Command{
		Use:     "workload <project-id>",
		Aliases: []string{"capacity"},
		Short:   "show per-member workload breakdown as JSON",
		Example: `  # View team workload
  yunxiao project workload 123

  # Filter by assignees
  yunxiao project workload 123 --assignee-ids "user1,user2"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ProjectID = args[0]
			result, err := callYunxiaoTool(cmd, cfg, "get_team_workload_breakdown", options.params())
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
	flags.StringVar(&options.Status, "status", "", "comma-separated status IDs to include")
	flags.IntVar(&options.MemberLimit, "member-limit", 0, "max members to inspect")
	flags.IntVar(&options.TaskLimit, "task-limit", 0, "max tasks per member")
	return command
}

func (o projectWorkloadOptions) params() map[string]any {
	params := map[string]any{"projectId": o.ProjectID}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "assigneeIds", o.AssigneeIDs)
	setCLIStringParam(params, "categories", o.Categories)
	setCLIStringParam(params, "status", o.Status)
	if o.MemberLimit > 0 {
		params["memberLimit"] = o.MemberLimit
	}
	if o.TaskLimit > 0 {
		params["taskLimit"] = o.TaskLimit
	}
	return params
}

type projectBoardOptions struct {
	OrganizationID string
	ProjectID      string
	Category       string
	Sprint         string
	Subject        string
	Status         string
	AssignedTo     string
	Creator        string
	SampleLimit    int
}

func newYunxiaoProjectBoardCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options projectBoardOptions
	command := &cobra.Command{
		Use:     "board <project-id>",
		Aliases: []string{"kanban"},
		Short:   "show Kanban board grouped by status as JSON",
		Example: `  # View task board
  yunxiao project board 123 --category Task

  # Filter by sprint
  yunxiao project board 123 --category Task --sprint sprint-456`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ProjectID = args[0]
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "get_project_workitem_board", params)
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.Category, "category", "", "work item category, e.g. Task or Bug")
	flags.StringVar(&options.Sprint, "sprint", "", "sprint ID filter")
	flags.StringVar(&options.Subject, "subject", "", "subject/title keyword")
	flags.StringVar(&options.Status, "status", "", "comma-separated status IDs")
	flags.StringVar(&options.AssignedTo, "assigned-to", "", "comma-separated assignee user IDs")
	flags.StringVar(&options.Creator, "creator", "", "comma-separated creator user IDs")
	flags.IntVar(&options.SampleLimit, "sample-limit", 0, "max work items returned")
	return command
}

func (o projectBoardOptions) params() (map[string]any, error) {
	params := map[string]any{
		"projectId": o.ProjectID,
		"category":  o.Category,
	}
	if params["category"] == "" {
		return nil, fmt.Errorf("category is required, e.g. Task or Bug")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "sprint", o.Sprint)
	setCLIStringParam(params, "subject", o.Subject)
	setCLIStringParam(params, "status", o.Status)
	setCLIStringParam(params, "assignedTo", o.AssignedTo)
	setCLIStringParam(params, "creator", o.Creator)
	if o.SampleLimit > 0 {
		params["sampleLimit"] = o.SampleLimit
	}
	return params, nil
}

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
