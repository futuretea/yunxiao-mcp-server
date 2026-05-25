package cmd

import (
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
		_, _ = fmt.Fprintln(out, raw)
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, "ID\tSUBJECT\tSTATUS\tASSIGNEE")
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
