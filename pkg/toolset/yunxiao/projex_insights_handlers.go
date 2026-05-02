package yunxiao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

func handleGetProjectRiskDashboard(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	categories := splitCSV(optionalStringDefault(params, "categories", "Risk,Bug,Task"))
	if len(categories) == 0 {
		return "", fmt.Errorf("categories must include at least one category")
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

func handleGetProjectMemberTaskStatus(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	members, assigneeIDs, err := projectTaskStatusMembers(ctx, c, organizationID, projectID, params)
	if err != nil {
		return "", err
	}
	if len(assigneeIDs) == 0 {
		return "", fmt.Errorf("assigneeIds is empty and no project members with userId were returned")
	}
	groups, err := parseStatusGroups(params)
	if err != nil {
		return "", err
	}

	status := map[string]any{
		"filters": projectTaskStatusFilters(params, assigneeIDs, groups),
		"members": map[string]any{},
	}
	memberPayloads := status["members"].(map[string]any)
	for _, assigneeID := range assigneeIDs {
		payload, err := projectMemberTaskPayload(ctx, c, organizationID, projectID, assigneeID, members[assigneeID], groups, params)
		if err != nil {
			return "", err
		}
		memberPayloads[assigneeID] = payload
	}
	return marshalPretty(status)
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
	overdueParams := copyParams(params)
	overdueParams["finishTimeBefore"] = optionalStringDefault(params, "overdueBefore", todayDate())
	overdue, err := searchProjectWorkitems(ctx, c, organizationID, projectID, strings.Join(categories, ","), overdueParams)
	if err != nil {
		return err
	}
	dashboard["overdue"] = overdue

	if highPriority := optionalStringDefault(params, "highPriority", ""); highPriority != "" {
		priorityParams := copyParams(params)
		priorityParams["priority"] = highPriority
		highPriorityPayload, err := searchProjectWorkitems(ctx, c, organizationID, projectID, strings.Join(categories, ","), priorityParams)
		if err != nil {
			return err
		}
		dashboard["highPriority"] = highPriorityPayload
	}
	if staleBefore := optionalStringDefault(params, "staleBefore", ""); staleBefore != "" {
		staleParams := copyParams(params)
		staleParams["updateStatusAtBefore"] = staleBefore
		stalePayload, err := searchProjectWorkitems(ctx, c, organizationID, projectID, strings.Join(categories, ","), staleParams)
		if err != nil {
			return err
		}
		dashboard["stale"] = stalePayload
	}
	return nil
}

func projectMemberTaskPayload(ctx context.Context, c *Client, organizationID, projectID, assigneeID string, member any, groups map[string]string, params map[string]any) (map[string]any, error) {
	memberParams := copyParams(params)
	memberParams["assignedTo"] = assigneeID
	categories := optionalStringDefault(params, "categories", "Task,Bug")

	assigned, err := searchProjectWorkitems(ctx, c, organizationID, projectID, categories, memberParams)
	if err != nil {
		return nil, err
	}
	overdueParams := copyParams(memberParams)
	overdueParams["finishTimeBefore"] = optionalStringDefault(params, "overdueBefore", todayDate())
	overdue, err := searchProjectWorkitems(ctx, c, organizationID, projectID, categories, overdueParams)
	if err != nil {
		return nil, err
	}

	payload := map[string]any{
		"member":   member,
		"assigned": assigned,
		"overdue":  overdue,
	}
	if len(groups) > 0 {
		groupPayloads, err := projectMemberStatusGroups(ctx, c, organizationID, projectID, categories, memberParams, groups)
		if err != nil {
			return nil, err
		}
		payload["statusGroups"] = groupPayloads
	}
	return payload, nil
}

func projectMemberStatusGroups(ctx context.Context, c *Client, organizationID, projectID, categories string, baseParams map[string]any, groups map[string]string) (map[string]any, error) {
	names := make([]string, 0, len(groups))
	for name := range groups {
		names = append(names, name)
	}
	sort.Strings(names)
	payloads := make(map[string]any, len(groups))
	for _, name := range names {
		groupParams := copyParams(baseParams)
		groupParams["status"] = groups[name]
		payload, err := searchProjectWorkitems(ctx, c, organizationID, projectID, categories, groupParams)
		if err != nil {
			return nil, err
		}
		payloads[name] = payload
	}
	return payloads, nil
}

func projectTaskStatusMembers(ctx context.Context, c *Client, organizationID, projectID string, params map[string]any) (map[string]any, []string, error) {
	if assigneeIDs := splitCSV(optionalStringDefault(params, "assigneeIds", "")); len(assigneeIDs) > 0 {
		members := make(map[string]any, len(assigneeIDs))
		for _, assigneeID := range assigneeIDs {
			members[assigneeID] = map[string]any{"userId": assigneeID}
		}
		return members, assigneeIDs, nil
	}

	resp, err := c.Request(ctx, http.MethodGet, projexProjectPath(organizationID, projectID)+"/members", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("members: %w", err)
	}
	return projectMembersFromResponse(resp, optionalIntDefault(params, "memberLimit", 20))
}

func projectMembersFromResponse(resp *Response, limit int) (map[string]any, []string, error) {
	var members []map[string]any
	if err := json.Unmarshal(resp.Body, &members); err != nil {
		return nil, nil, fmt.Errorf("decode members: %w", err)
	}
	if limit < 0 {
		limit = 0
	}
	result := make(map[string]any, len(members))
	ids := make([]string, 0, len(members))
	for _, member := range members {
		userID, _ := member["userId"].(string)
		if strings.TrimSpace(userID) == "" {
			continue
		}
		if limit > 0 && len(ids) >= limit {
			break
		}
		userID = strings.TrimSpace(userID)
		result[userID] = member
		ids = append(ids, userID)
	}
	return result, ids, nil
}

func searchProjectWorkitems(ctx context.Context, c *Client, organizationID, projectID, category string, params map[string]any) (any, error) {
	body := projectInsightWorkitemBody(projectID, category, params)
	path := projexOrganizationPath(organizationID) + "/workitems:search"
	resp, err := c.Request(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", category, err)
	}
	return responsePayload(resp), nil
}

func projectInsightWorkitemBody(projectID, category string, params map[string]any) map[string]any {
	body := projectWorkitemSummaryBody(projectID, category, params)
	body["perPage"] = normalizedSampleLimit(params)
	return body
}

func parseStatusGroups(params map[string]any) (map[string]string, error) {
	raw := optionalStringDefault(params, "statusGroups", "")
	if raw == "" {
		return nil, nil
	}
	var groups map[string]string
	if err := json.Unmarshal([]byte(raw), &groups); err != nil {
		return nil, fmt.Errorf("statusGroups must be a JSON object of name to comma-separated status IDs: %w", err)
	}
	for name, value := range groups {
		if strings.TrimSpace(name) == "" || strings.TrimSpace(value) == "" {
			delete(groups, name)
			continue
		}
		groups[name] = strings.TrimSpace(value)
	}
	return groups, nil
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

func projectTaskStatusFilters(params map[string]any, assigneeIDs []string, groups map[string]string) map[string]any {
	return map[string]any{
		"assigneeIds":   assigneeIDs,
		"categories":    splitCSV(optionalStringDefault(params, "categories", "Task,Bug")),
		"subject":       optionalStringDefault(params, "subject", ""),
		"status":        optionalStringDefault(params, "status", ""),
		"statusStage":   optionalStringDefault(params, "statusStage", ""),
		"assignedTo":    optionalStringDefault(params, "assignedTo", ""),
		"creator":       optionalStringDefault(params, "creator", ""),
		"sprint":        optionalStringDefault(params, "sprint", ""),
		"workitemType":  optionalStringDefault(params, "workitemType", ""),
		"tag":           optionalStringDefault(params, "tag", ""),
		"overdueBefore": optionalStringDefault(params, "overdueBefore", todayDate()),
		"statusGroups":  groups,
		"memberLimit":   optionalIntDefault(params, "memberLimit", 20),
		"sampleLimit":   normalizedSampleLimit(params),
	}
}

func copyParams(params map[string]any) map[string]any {
	copied := make(map[string]any, len(params))
	for key, value := range params {
		copied[key] = value
	}
	return copied
}

func todayDate() string {
	return time.Now().Format("2006-01-02")
}
