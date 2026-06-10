package yunxiao

import (
	"context"
	"fmt"
	"net/http"
)

func handleGetProjectOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
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
		if err := addOverviewSection(ctx, c, overview, params, section); err != nil {
			return "", err
		}
	}

	return marshalPretty(overview)
}

func handleGetProjectWorkitemSummary(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
	if err != nil {
		return "", err
	}

	categories := splitCSV(optionalStringDefault(params, "categories", "Req,Task,Bug,Risk"))
	if len(categories) == 0 {
		return "", errNoCategories
	}

	result, err := buildCategoryResult(ctx, categories, projectWorkitemSummaryFilters(params, categories),
		func(cat string) (any, error) {
			return searchProjectWorkitems(ctx, c, organizationID, projectID, cat, params)
		})
	if err != nil {
		return "", err
	}
	return marshalPretty(result)
}

func handleGetProjectWorkitemContext(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
	if err != nil {
		return "", err
	}
	category, err := requiredString(params, "category")
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

func handleGetMyProjectWorkitems(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
	if err != nil {
		return "", err
	}
	userID, err := requiredString(params, "userId")
	if err != nil {
		return "", err
	}

	categories := splitCSV(optionalStringDefault(params, "categories", "Task,Bug"))
	if len(categories) == 0 {
		return "", errNoCategories
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
			return searchProjectWorkitems(ctx, c, organizationID, projectID, cat, searchParams)
		})
	if err != nil {
		return "", err
	}
	return marshalPretty(result)
}

func handleGetProjectWorkitemBoard(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
	if err != nil {
		return "", err
	}
	category, err := requiredString(params, "category")
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
	}

	data, total, err := extractWorkitemData(payload)
	if err != nil {
		return "", err
	}
	board["total"] = total

	columns, counts := groupWorkitemsByStatus(data)
	board["columns"] = columns
	board["columnCounts"] = counts

	return marshalPretty(board)
}
