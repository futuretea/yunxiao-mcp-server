package yunxiao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func normalizedSampleLimit(params map[string]any) int {
	limit := optionalIntDefault(params, "sampleLimit", 5)
	if limit < 0 {
		return 0
	}
	if limit > 200 {
		return 200
	}
	return limit
}

func projectWorkitemSummaryBody(projectID, category string, params map[string]any) map[string]any {
	body := map[string]any{
		"category": category,
		"spaceId":  projectID,
		"page":     1,
		"perPage":  normalizedSampleLimit(params),
	}
	setOptionalStringBody(body, params, "conditions")
	setOptionalStringBody(body, params, "orderBy")
	setOptionalStringBody(body, params, "sort")
	if body["conditions"] == nil {
		if conditions := buildWorkitemConditions(params); conditions != "" {
			body["conditions"] = conditions
		}
	}
	return body
}

func searchProjectWorkitemSummaryCategory(ctx context.Context, c *Client, organizationID, projectID, category string, params map[string]any) (any, error) {
	body := projectWorkitemSummaryBody(projectID, category, params)
	path := projexOrganizationPath(organizationID) + "/workitems:search"
	resp, err := c.Request(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", category, err)
	}
	return responsePayload(resp), nil
}

func projectWorkitemSummaryFilters(params map[string]any, categories []string) map[string]any {
	return map[string]any{
		"categories":  categories,
		"subject":     optionalStringDefault(params, "subject", ""),
		"status":      optionalStringDefault(params, "status", ""),
		"assignedTo":  optionalStringDefault(params, "assignedTo", ""),
		"creator":     optionalStringDefault(params, "creator", ""),
		"tag":         optionalStringDefault(params, "tag", ""),
		"orderBy":     optionalStringDefault(params, "orderBy", ""),
		"sort":        optionalStringDefault(params, "sort", ""),
		"sampleLimit": normalizedSampleLimit(params),
	}
}

func projectWorkitemContextFilters(params map[string]any, category string) map[string]any {
	return map[string]any{
		"category":       category,
		"workItemTypeId": optionalStringDefault(params, "workItemTypeId", ""),
		"page":           optionalIntDefault(params, "page", 1),
		"perPage":        optionalIntDefault(params, "perPage", 20),
	}
}

func buildCategoryResult(ctx context.Context, categories []string, filters map[string]any, searchFn func(string) (any, error)) (map[string]any, error) {
	result := map[string]any{
		"filters":    filters,
		"categories": map[string]any{},
	}
	payloads := result["categories"].(map[string]any)
	for _, category := range categories {
		payload, err := searchFn(category)
		if err != nil {
			return nil, err
		}
		payloads[category] = payload
	}
	return result, nil
}

func searchSprintWorkitems(ctx context.Context, c *Client, organizationID, projectID, sprintID, category string, params map[string]any) (any, error) {
	body := projectWorkitemSummaryBody(projectID, category, params)
	sprintConditions := buildWorkitemConditions(map[string]any{"sprint": sprintID})
	if sprintConditions != "" {
		existing, _ := body["conditions"].(string)
		body["conditions"] = mergeConditions(existing, sprintConditions)
	}
	path := projexOrganizationPath(organizationID) + "/workitems:search"
	resp, err := c.Request(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", category, err)
	}
	return responsePayload(resp), nil
}

func mergeConditions(existing, extra string) string {
	if existing == "" {
		return extra
	}
	if extra == "" {
		return existing
	}
	var existingArr, extraArr []map[string]any
	existingArrErr := json.Unmarshal([]byte(existing), &existingArr)
	extraArrErr := json.Unmarshal([]byte(extra), &extraArr)
	if existingArrErr == nil && extraArrErr == nil {
		merged := append(existingArr, extraArr...)
		mergedBytes, _ := json.Marshal(merged)
		return string(mergedBytes)
	}

	var existingObj, extraObj map[string]any
	if err := json.Unmarshal([]byte(existing), &existingObj); err != nil {
		return existing
	}
	if err := json.Unmarshal([]byte(extra), &extraObj); err != nil {
		return existing
	}

	existingGroups, _ := existingObj["conditionGroups"].([]any)
	extraGroups, _ := extraObj["conditionGroups"].([]any)
	if len(existingGroups) > 0 && len(extraGroups) > 0 {
		existingGroup, _ := existingGroups[0].([]any)
		extraGroup, _ := extraGroups[0].([]any)
		existingGroups[0] = append(existingGroup, extraGroup...)
	} else if len(extraGroups) > 0 {
		existingObj["conditionGroups"] = extraGroups
	}

	mergedBytes, _ := json.Marshal(existingObj)
	return string(mergedBytes)
}

func sprintOverviewFilters(params map[string]any, categories []string) map[string]any {
	return map[string]any{
		"categories":  categories,
		"status":      optionalStringDefault(params, "status", ""),
		"sampleLimit": normalizedSampleLimit(params),
	}
}

func myProjectWorkitemFilters(params map[string]any, userID, relation string, categories []string) map[string]any {
	return map[string]any{
		"userId":      userID,
		"relation":    relation,
		"categories":  categories,
		"status":      optionalStringDefault(params, "status", ""),
		"sampleLimit": normalizedSampleLimit(params),
	}
}

func extractWorkitemData(payload any) ([]any, int, error) {
	switch p := payload.(type) {
	case []any:
		return p, len(p), nil
	case map[string]any:
		data, _ := p["data"].([]any)
		total := 0
		if pagination, ok := p["pagination"].(map[string]any); ok {
			if t, ok := pagination["total"].(float64); ok {
				total = int(t)
			}
		}
		return data, total, nil
	default:
		return nil, 0, fmt.Errorf("unexpected payload type %T", payload)
	}
}

func extractStatusName(item any) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}
	status, ok := m["status"].(map[string]any)
	if !ok {
		return ""
	}
	name, _ := status["name"].(string)
	return name
}

func marshalPretty(value any) (string, error) {
	formatted, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}

func optionalIntDefault(params map[string]any, key string, defaultValue int) int {
	switch value := params[key].(type) {
	case float64:
		return int(value)
	case int:
		return value
	case int64:
		return int(value)
	case string:
		if parsed, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func responsePayload(resp *Response) any {
	var data any
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		data = string(resp.Body)
	}

	if resp.Pagination == nil && resp.NextToken == "" && resp.RequestID == "" {
		return data
	}

	payload := map[string]any{"data": data}
	if resp.Pagination != nil {
		payload["pagination"] = resp.Pagination
	}
	if resp.NextToken != "" {
		payload["nextToken"] = resp.NextToken
	}
	if resp.RequestID != "" {
		payload["requestId"] = resp.RequestID
	}
	return payload
}

type projectOverviewSection struct {
	flag  string
	name  string
	path  string
	query url.Values
}

func projectOverviewSections(projectPath string, params map[string]any) []projectOverviewSection {
	return []projectOverviewSection{
		{flag: "includeMembers", name: "members", path: projectPath + "/members"},
		{flag: "includeSprints", name: "sprints", path: projectPath + "/sprints", query: projectOverviewListQuery(params, true)},
		{flag: "includeMilestones", name: "milestones", path: projectPath + "/milestones", query: projectOverviewListQuery(params, true)},
		{flag: "includeVersions", name: "versions", path: projectPath + "/versions", query: projectOverviewListQuery(params, true)},
		{flag: "includeLabels", name: "labels", path: projectPath + "/labels", query: projectOverviewListQuery(params, false)},
	}
}

func addProjectOverviewSection(ctx context.Context, c *Client, overview map[string]any, params map[string]any, section projectOverviewSection) error {
	if !optionalBoolDefault(params, section.flag, true) {
		return nil
	}
	payload, err := getProjectOverviewSection(ctx, c, section.name, section.path, section.query)
	if err != nil {
		return err
	}
	overview[section.name] = payload
	return nil
}

func getProjectOverviewSection(ctx context.Context, c *Client, name, path string, query url.Values) (any, error) {
	resp, err := c.Request(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	return responsePayload(resp), nil
}

func projectOverviewListQuery(params map[string]any, withStatus bool) url.Values {
	query := url.Values{}
	query.Set("page", strconv.Itoa(optionalIntDefault(params, "page", 1)))
	query.Set("perPage", strconv.Itoa(optionalIntDefault(params, "perPage", 20)))
	if withStatus && optionalBoolDefault(params, "activeOnly", true) {
		status := optionalStringDefault(params, "status", "TODO,DOING")
		if status != "" {
			query.Set("status", status)
		}
	}
	return query
}

func projectOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"activeOnly": optionalBoolDefault(params, "activeOnly", true),
		"status":     optionalStringDefault(params, "status", "TODO,DOING"),
		"page":       optionalIntDefault(params, "page", 1),
		"perPage":    optionalIntDefault(params, "perPage", 20),
	}
}
