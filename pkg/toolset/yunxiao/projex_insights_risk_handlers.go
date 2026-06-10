package yunxiao

import (
	"context"
	"strings"
)

func handleGetProjectRiskDashboard(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
	if err != nil {
		return "", err
	}

	categories := splitCSV(optionalStringDefault(params, "categories", "Risk,Bug,Task"))
	if len(categories) == 0 {
		return "", errNoCategories
	}

	dashboard := map[string]any{
		"filters":      riskDashboardFilters(params, categories),
		"byCategory":   map[string]any{},
		"overdue":      nil,
		"highPriority": nil,
		"stale":        nil,
	}
	if err := addRiskCategorySections(ctx, c, dashboard["byCategory"].(map[string]any), organizationID, projectID, categories, params); err != nil {
		return "", err
	}
	if err := addRiskFocusSections(ctx, c, dashboard, organizationID, projectID, categories, params); err != nil {
		return "", err
	}
	return marshalPretty(dashboard)
}

func addRiskCategorySections(ctx context.Context, c *Client, target map[string]any, organizationID, projectID string, categories []string, params map[string]any) error {
	for _, category := range categories {
		payload, err := searchProjectWorkitems(ctx, c, organizationID, projectID, category, params)
		if err != nil {
			return err
		}
		target[category] = payload
	}
	return nil
}

func addRiskFocusSections(ctx context.Context, c *Client, dashboard map[string]any, organizationID, projectID string, categories []string, params map[string]any) error {
	categoryList := strings.Join(categories, ",")

	overdueParams := copyParams(params)
	overdueParams["finishTimeBefore"] = optionalStringDefault(params, "overdueBefore", todayDate())
	overdue, err := searchProjectWorkitems(ctx, c, organizationID, projectID, categoryList, overdueParams)
	if err != nil {
		return err
	}
	dashboard["overdue"] = overdue

	if highPriority := optionalStringDefault(params, "highPriority", ""); highPriority != "" {
		priorityParams := copyParams(params)
		priorityParams["priority"] = highPriority
		highPriorityPayload, err := searchProjectWorkitems(ctx, c, organizationID, projectID, categoryList, priorityParams)
		if err != nil {
			return err
		}
		dashboard["highPriority"] = highPriorityPayload
	}
	if staleBefore := optionalStringDefault(params, "staleBefore", ""); staleBefore != "" {
		staleParams := copyParams(params)
		staleParams["updateStatusAtBefore"] = staleBefore
		stalePayload, err := searchProjectWorkitems(ctx, c, organizationID, projectID, categoryList, staleParams)
		if err != nil {
			return err
		}
		dashboard["stale"] = stalePayload
	}
	return nil
}

func riskDashboardFilters(params map[string]any, categories []string) map[string]any {
	return map[string]any{
		"categories":    categories,
		"subject":       optionalStringDefault(params, "subject", ""),
		"status":        optionalStringDefault(params, "status", ""),
		"statusStage":   optionalStringDefault(params, "statusStage", ""),
		"assignedTo":    optionalStringDefault(params, "assignedTo", ""),
		"creator":       optionalStringDefault(params, "creator", ""),
		"sprint":        optionalStringDefault(params, "sprint", ""),
		"workitemType":  optionalStringDefault(params, "workitemType", ""),
		"tag":           optionalStringDefault(params, "tag", ""),
		"overdueBefore": optionalStringDefault(params, "overdueBefore", todayDate()),
		"highPriority":  optionalStringDefault(params, "highPriority", ""),
		"staleBefore":   optionalStringDefault(params, "staleBefore", ""),
		"sampleLimit":   normalizedSampleLimit(params),
	}
}
