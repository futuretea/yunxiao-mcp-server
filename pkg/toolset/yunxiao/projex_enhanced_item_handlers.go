package yunxiao

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func handleGetSprintOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
	if err != nil {
		return "", err
	}
	sprintID, err := requiredString(params, "sprintId")
	if err != nil {
		return "", err
	}

	categories := splitCSV(optionalStringDefault(params, "categories", "Task,Bug"))
	if len(categories) == 0 {
		return "", errNoCategories
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

func handleGetProjectWorkitemDetail(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	workitemID, err := requiredString(params, "workitemId")
	if err != nil {
		return "", err
	}

	workitemPath := projexWorkitemPath(organizationID, workitemID)
	workitem, err := getProjectOverviewSection(ctx, c, "workitem", workitemPath, nil)
	if err != nil {
		return "", err
	}

	detail := map[string]any{
		"workitem": workitem,
		"filters":  workitemDetailFilters(params),
	}

	for _, section := range workitemDetailSections(workitemPath, params) {
		if err := addOverviewSection(ctx, c, detail, params, section); err != nil {
			return "", err
		}
	}

	return marshalPretty(detail)
}

func handleGetWorkItemTypeOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, projectID, workItemTypeID, err := requiredOrganizationProjectAndWorkItemType(params)
	if err != nil {
		return "", err
	}

	typePath := projexOrganizationPath(organizationID) + "/workitemTypes/" + url.PathEscape(workItemTypeID)
	workItemType, err := getProjectOverviewSection(ctx, c, "workItemType", typePath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"workItemType": workItemType,
		"filters":      workItemTypeOverviewFilters(params),
	}

	projectTypePath := workItemTypeProjectPath(organizationID, projectID, workItemTypeID)

	if optionalBoolDefault(params, "includeFieldConfig", true) {
		fields, err := getProjectOverviewSection(ctx, c, "fieldConfig", projectTypePath+"/fields", nil)
		if err != nil {
			return "", err
		}
		overview["fieldConfig"] = fields
	}

	if optionalBoolDefault(params, "includeWorkflow", true) {
		workflow, err := getProjectOverviewSection(ctx, c, "workflow", projectTypePath+"/workflows", nil)
		if err != nil {
			return "", err
		}
		overview["workflow"] = workflow
	}

	return marshalPretty(overview)
}

func workItemTypeOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeFieldConfig": optionalBoolDefault(params, "includeFieldConfig", true),
		"includeWorkflow":    optionalBoolDefault(params, "includeWorkflow", true),
	}
}
