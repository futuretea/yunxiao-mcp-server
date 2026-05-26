package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"fmt"
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
			if options.JSONOutput {
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
