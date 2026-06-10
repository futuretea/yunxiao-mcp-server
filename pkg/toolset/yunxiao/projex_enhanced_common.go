package yunxiao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

func searchProjectWorkitems(ctx context.Context, c *Client, organizationID, projectID, category string, params map[string]any) (any, error) {
	body := projectWorkitemSummaryBody(projectID, category, params)
	path := projexOrganizationPath(organizationID) + "/workitems:search"
	resp, err := c.Request(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", category, err)
	}
	return responsePayload(resp), nil
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
	if result, ok := mergeConditionsAsArrays(existing, extra); ok {
		return result
	}
	return mergeConditionGroups(existing, extra)
}

func mergeConditionsAsArrays(existing, extra string) (string, bool) {
	var existingArr, extraArr []map[string]any
	if err := json.Unmarshal([]byte(existing), &existingArr); err != nil {
		return "", false
	}
	if err := json.Unmarshal([]byte(extra), &extraArr); err != nil {
		return "", false
	}
	mergedBytes, err := json.Marshal(append(existingArr, extraArr...))
	if err != nil {
		return "", false
	}
	return string(mergedBytes), true
}

func mergeConditionGroups(existing, extra string) string {
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

	mergedBytes, err := json.Marshal(existingObj)
	if err != nil {
		return existing
	}
	return string(mergedBytes)
}

func sprintOverviewFilters(params map[string]any, categories []string) map[string]any {
	return map[string]any{
		"categories":  categories,
		"subject":     optionalStringDefault(params, "subject", ""),
		"status":      optionalStringDefault(params, "status", ""),
		"assignedTo":  optionalStringDefault(params, "assignedTo", ""),
		"creator":     optionalStringDefault(params, "creator", ""),
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
			total = intValue(pagination["total"])
		}
		return data, total, nil
	default:
		return nil, 0, fmt.Errorf("unexpected payload type %T", payload)
	}
}

func intValue(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	}
	return 0
}

func groupWorkitemsByStatus(data []any) (map[string]any, map[string]int) {
	columns := map[string]any{}
	counts := map[string]int{}
	for _, item := range data {
		itemMap, _ := item.(map[string]any)
		statusName := extractWorkitemStatusName(itemMap)
		col, _ := columns[statusName].([]any)
		columns[statusName] = append(col, item)
		counts[statusName]++
	}
	return columns, counts
}

func extractWorkitemStatusName(itemMap map[string]any) string {
	if status, ok := itemMap["status"].(map[string]any); ok {
		if name, ok := status["name"].(string); ok {
			return name
		}
	}
	return "Unknown"
}

// parseListData extracts a []any from either a raw slice or a map["data"] wrapper.
func parseListData(data any) []any {
	switch d := data.(type) {
	case []any:
		return d
	case map[string]any:
		if data, ok := d["data"].([]any); ok {
			return data
		}
	}
	return nil
}
