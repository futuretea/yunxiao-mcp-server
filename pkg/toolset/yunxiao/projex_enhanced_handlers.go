package yunxiao

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func handleGetProjectOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndID(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	projectPath := projexProjectPath(organizationID, projectID)
	overview := map[string]any{
		"project": nil,
		"filters": projectOverviewFilters(params),
	}

	project, err := getProjectOverviewSection(ctx, c, "project", projectPath, nil)
	if err != nil {
		return "", err
	}
	overview["project"] = project

	for _, section := range projectOverviewSections(projectPath, params) {
		if err := addProjectOverviewSection(ctx, c, overview, params, section); err != nil {
			return "", err
		}
	}

	return marshalPretty(overview)
}

func handleGetProjectWorkitemSummary(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndID(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	categories := splitCSV(optionalStringDefault(params, "categories", "Req,Task,Bug,Risk"))
	if len(categories) == 0 {
		return "", fmt.Errorf("categories must include at least one category")
	}

	result, err := buildCategoryResult(ctx, categories, projectWorkitemSummaryFilters(params, categories),
		func(cat string) (any, error) {
			return searchProjectWorkitemSummaryCategory(ctx, c, organizationID, projectID, cat, params)
		})
	if err != nil {
		return "", err
	}
	return marshalPretty(result)
}

func handleGetProjectWorkitemContext(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndID(params)
	if err != nil {
		return "", err
	}
	category, err := requiredString(params, "category")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	projectPath := projexProjectPath(organizationID, projectID)
	payload := map[string]any{"filters": projectWorkitemContextFilters(params, category)}
	if err := addProjectWorkitemContextBaseSections(ctx, c, payload, params, projectPath, category); err != nil {
		return "", err
	}
	if err := addProjectWorkitemTypeContext(ctx, c, payload, params, projectPath); err != nil {
		return "", err
	}
	return marshalPretty(payload)
}

func addProjectWorkitemContextBaseSections(ctx context.Context, c *Client, payload map[string]any, params map[string]any, projectPath, category string) error {
	typeQuery := url.Values{}
	typeQuery.Set("category", category)
	types, err := getProjectOverviewSection(ctx, c, "workItemTypes", projectPath+"/workitemTypes", typeQuery)
	if err != nil {
		return err
	}
	payload["workItemTypes"] = types

	if optionalBoolDefault(params, "includeMembers", true) {
		members, err := getProjectOverviewSection(ctx, c, "members", projectPath+"/members", nil)
		if err != nil {
			return err
		}
		payload["members"] = members
	}
	if optionalBoolDefault(params, "includeLabels", true) {
		labels, err := getProjectOverviewSection(ctx, c, "labels", projectPath+"/labels", projectOverviewListQuery(params, false))
		if err != nil {
			return err
		}
		payload["labels"] = labels
	}
	return nil
}

func addProjectWorkitemTypeContext(ctx context.Context, c *Client, payload map[string]any, params map[string]any, projectPath string) error {
	workItemTypeID, _ := params["workItemTypeId"].(string)
	if strings.TrimSpace(workItemTypeID) == "" {
		return nil
	}
	typePath := projectPath + "/workitemTypes/" + url.PathEscape(strings.TrimSpace(workItemTypeID))
	if optionalBoolDefault(params, "includeFields", true) {
		fields, err := getProjectOverviewSection(ctx, c, "fields", typePath+"/fields", nil)
		if err != nil {
			return err
		}
		payload["fields"] = fields
	}
	if optionalBoolDefault(params, "includeWorkflow", true) {
		workflow, err := getProjectOverviewSection(ctx, c, "workflow", typePath+"/workflows", nil)
		if err != nil {
			return err
		}
		payload["workflow"] = workflow
	}
	return nil
}

func handleGetMyProjectWorkitems(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndID(params)
	if err != nil {
		return "", err
	}
	userID, err := requiredString(params, "userId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	categories := splitCSV(optionalStringDefault(params, "categories", "Task,Bug"))
	if len(categories) == 0 {
		return "", fmt.Errorf("categories must include at least one category")
	}

	relation := optionalStringDefault(params, "relation", "assigned")
	searchParams := copyParams(params)
	switch relation {
	case "assigned":
		searchParams["assignedTo"] = userID
	case "created":
		searchParams["creator"] = userID
	default:
		return "", fmt.Errorf("relation must be assigned or created")
	}

	result, err := buildCategoryResult(ctx, categories, myProjectWorkitemFilters(params, userID, relation, categories),
		func(cat string) (any, error) {
			return searchProjectWorkitemSummaryCategory(ctx, c, organizationID, projectID, cat, searchParams)
		})
	if err != nil {
		return "", err
	}
	return marshalPretty(result)
}

func handleGetSprintOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndID(params)
	if err != nil {
		return "", err
	}
	sprintID, err := requiredString(params, "sprintId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	categories := splitCSV(optionalStringDefault(params, "categories", "Task,Bug"))
	if len(categories) == 0 {
		return "", fmt.Errorf("categories must include at least one category")
	}

	sprintPath := projexProjectPath(organizationID, projectID) + "/sprints/" + url.PathEscape(strings.TrimSpace(sprintID))
	sprintResp, err := c.Request(ctx, http.MethodGet, sprintPath, nil, nil)
	if err != nil {
		return "", fmt.Errorf("sprint: %w", err)
	}

	result, err := buildCategoryResult(ctx, categories, sprintOverviewFilters(params, categories),
		func(cat string) (any, error) {
			return searchSprintWorkitems(ctx, c, organizationID, projectID, sprintID, cat, params)
		})
	if err != nil {
		return "", err
	}
	result["sprint"] = responsePayload(sprintResp)
	return marshalPretty(result)
}

func handleGetProjectWorkitemBoard(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndID(params)
	if err != nil {
		return "", err
	}
	category, err := requiredString(params, "category")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	body := projectWorkitemSummaryBody(projectID, category, params)
	sprintID := optionalStringDefault(params, "sprint", "")

	path := projexOrganizationPath(organizationID) + "/workitems:search"
	resp, err := c.Request(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return "", fmt.Errorf("search: %w", err)
	}

	payload := responsePayload(resp)
	board := map[string]any{
		"filters": map[string]any{
			"category":    category,
			"sprint":      sprintID,
			"subject":     optionalStringDefault(params, "subject", ""),
			"status":      optionalStringDefault(params, "status", ""),
			"assignedTo":  optionalStringDefault(params, "assignedTo", ""),
			"creator":     optionalStringDefault(params, "creator", ""),
			"sampleLimit": normalizedSampleLimit(params),
		},
		"columns": map[string]any{},
	}

	data, total, err := extractWorkitemData(payload)
	if err != nil {
		return "", err
	}
	board["total"] = total

	columns := board["columns"].(map[string]any)
	for _, item := range data {
		statusName := extractStatusName(item)
		if statusName == "" {
			statusName = "Unknown"
		}
		if columns[statusName] == nil {
			columns[statusName] = []any{}
		}
		columns[statusName] = append(columns[statusName].([]any), item)
	}

	return marshalPretty(board)
}
