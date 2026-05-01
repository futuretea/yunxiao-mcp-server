package yunxiao

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func getClient(client any) (*Client, error) {
	c, ok := client.(*Client)
	if !ok || c == nil {
		return nil, fmt.Errorf("Yunxiao client is not configured")
	}
	return c, nil
}

func requiredString(params map[string]any, key string) (string, error) {
	value, _ := params[key].(string)
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

func requiredOrganizationAndRepository(params map[string]any) (string, string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", "", err
	}
	repositoryID, err := requiredString(params, "repositoryId")
	if err != nil {
		return "", "", err
	}
	return organizationID, repositoryID, nil
}

func requiredOrganizationAndID(params map[string]any) (string, string, error) {
	return requiredOrganizationAndNamedID(params, "id")
}

func requiredOrganizationAndNamedID(params map[string]any, key string) (string, string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", "", err
	}
	id, err := requiredString(params, key)
	if err != nil {
		return "", "", err
	}
	return organizationID, id, nil
}

func requiredOrganizationRepositoryAndLocalID(params map[string]any) (string, string, string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", "", "", err
	}
	localID, err := requiredString(params, "localId")
	if err != nil {
		return "", "", "", err
	}
	return organizationID, repositoryID, localID, nil
}

func buildProjectConditions(params map[string]any) string {
	filterConditions := make([]map[string]any, 0)
	if name, _ := params["name"].(string); name != "" {
		filterConditions = append(filterConditions, stringContainsCondition("name", name))
	}
	if status, _ := params["status"].(string); status != "" {
		filterConditions = append(filterConditions, listContainsCondition("status", "status", splitCSV(status)))
	}
	if creator, _ := params["creator"].(string); creator != "" {
		filterConditions = append(filterConditions, listContainsCondition("creator", "user", splitCSV(creator)))
	}
	if len(filterConditions) == 0 {
		return ""
	}
	return marshalConditions(filterConditions)
}

func buildWorkitemConditions(params map[string]any) string {
	filterConditions := make([]map[string]any, 0)
	if subject, _ := params["subject"].(string); subject != "" {
		filterConditions = append(filterConditions, stringContainsCondition("subject", subject))
	}
	if status, _ := params["status"].(string); status != "" {
		filterConditions = append(filterConditions, listContainsCondition("status", "status", splitCSV(status)))
	}
	if assignedTo, _ := params["assignedTo"].(string); assignedTo != "" {
		filterConditions = append(filterConditions, listContainsCondition("assignedTo", "user", splitCSV(assignedTo)))
	}
	if creator, _ := params["creator"].(string); creator != "" {
		filterConditions = append(filterConditions, listContainsCondition("creator", "user", splitCSV(creator)))
	}
	if sprint, _ := params["sprint"].(string); sprint != "" {
		filterConditions = append(filterConditions, listContainsCondition("sprint", "sprint", splitCSV(sprint)))
	}
	if workitemType, _ := params["workitemType"].(string); workitemType != "" {
		filterConditions = append(filterConditions, listContainsCondition("workitemType", "workitemType", splitCSV(workitemType)))
	}
	if statusStage, _ := params["statusStage"].(string); statusStage != "" {
		filterConditions = append(filterConditions, listContainsCondition("statusStage", "statusStage", splitCSV(statusStage)))
	}
	if tag, _ := params["tag"].(string); tag != "" {
		filterConditions = append(filterConditions, containsCondition("tag", "tag", "multiList", splitCSV(tag)))
	}
	if priority, _ := params["priority"].(string); priority != "" {
		filterConditions = append(filterConditions, listContainsCondition("priority", "option", splitCSV(priority)))
	}
	if subjectDescription, _ := params["subjectDescription"].(string); subjectDescription != "" {
		filterConditions = append(filterConditions, stringContainsCondition("subject-description", subjectDescription))
	}
	filterConditions = appendDateRangeCondition(filterConditions, "gmtCreate", "dateTime", params, "createdAfter", "createdBefore")
	filterConditions = appendDateRangeCondition(filterConditions, "gmtModified", "dateTime", params, "updatedAfter", "updatedBefore")
	filterConditions = appendDateRangeCondition(filterConditions, "finishTime", "date", params, "finishTimeAfter", "finishTimeBefore")
	filterConditions = appendDateRangeCondition(filterConditions, "updateStatusAt", "date", params, "updateStatusAtAfter", "updateStatusAtBefore")
	if len(filterConditions) == 0 {
		return ""
	}
	return marshalConditions(filterConditions)
}

func stringContainsCondition(fieldIdentifier, value string) map[string]any {
	return map[string]any{
		"className":       "string",
		"fieldIdentifier": fieldIdentifier,
		"format":          "input",
		"operator":        "CONTAINS",
		"toValue":         nil,
		"value":           []string{value},
	}
}

func listContainsCondition(fieldIdentifier, className string, values []string) map[string]any {
	return containsCondition(fieldIdentifier, className, "list", values)
}

func containsCondition(fieldIdentifier, className, format string, values []string) map[string]any {
	return map[string]any{
		"className":       className,
		"fieldIdentifier": fieldIdentifier,
		"format":          format,
		"operator":        "CONTAINS",
		"toValue":         nil,
		"value":           values,
	}
}

func appendDateRangeCondition(filterConditions []map[string]any, fieldIdentifier, className string, params map[string]any, afterKey, beforeKey string) []map[string]any {
	after, _ := params[afterKey].(string)
	before, _ := params[beforeKey].(string)
	after = strings.TrimSpace(after)
	before = strings.TrimSpace(before)
	if after == "" && before == "" {
		return filterConditions
	}
	if after == "" {
		after = "1970-01-01"
	}
	var toValue any
	if before != "" {
		toValue = endOfDay(before)
	}
	return append(filterConditions, map[string]any{
		"className":       className,
		"fieldIdentifier": fieldIdentifier,
		"format":          "input",
		"operator":        "BETWEEN",
		"toValue":         toValue,
		"value":           []string{startOfDay(after)},
	})
}

func startOfDay(value string) string {
	if len(value) == len("2006-01-02") {
		return value + " 00:00:00"
	}
	return value
}

func endOfDay(value string) string {
	if len(value) == len("2006-01-02") {
		return value + " 23:59:59"
	}
	return value
}

func marshalConditions(filterConditions []map[string]any) string {
	conditions := map[string]any{
		"conditionGroups": []any{filterConditions},
	}
	data, err := json.Marshal(conditions)
	if err != nil {
		return ""
	}
	return string(data)
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			values = append(values, part)
		}
	}
	return values
}

func requiredOrganizationAndPipeline(params map[string]any) (string, string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", "", err
	}
	pipelineID, err := requiredString(params, "pipelineId")
	if err != nil {
		return "", "", err
	}
	return organizationID, pipelineID, nil
}

func optionalStringDefault(params map[string]any, key, defaultValue string) string {
	if value, _ := params[key].(string); strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return defaultValue
}

func optionalBoolDefault(params map[string]any, key string, defaultValue bool) bool {
	if value, ok := params[key].(bool); ok {
		return value
	}
	return defaultValue
}

func setOptionalStringBody(body map[string]any, params map[string]any, key string) {
	value, _ := params[key].(string)
	if value != "" {
		body[key] = value
	}
}

func setOptionalIntBody(body map[string]any, params map[string]any, key string) {
	switch value := params[key].(type) {
	case float64:
		body[key] = int(value)
	case int:
		body[key] = value
	case int64:
		body[key] = value
	case string:
		if value != "" {
			body[key] = value
		}
	}
}

func setOptionalBoolBody(body map[string]any, params map[string]any, key string) {
	value, ok := params[key].(bool)
	if ok {
		body[key] = value
	}
}

func setOptionalStringArrayBody(body map[string]any, params map[string]any, key string) {
	switch value := params[key].(type) {
	case []any:
		values := make([]string, 0, len(value))
		for _, item := range value {
			if item, ok := item.(string); ok && strings.TrimSpace(item) != "" {
				values = append(values, strings.TrimSpace(item))
			}
		}
		if len(values) > 0 {
			body[key] = values
		}
	case []string:
		values := make([]string, 0, len(value))
		for _, item := range value {
			if strings.TrimSpace(item) != "" {
				values = append(values, strings.TrimSpace(item))
			}
		}
		if len(values) > 0 {
			body[key] = values
		}
	}
}

func setOptionalStringArrayQuery(query url.Values, params map[string]any, key string) {
	switch value := params[key].(type) {
	case []any:
		for _, item := range value {
			if item, ok := item.(string); ok && strings.TrimSpace(item) != "" {
				query.Add(key, strings.TrimSpace(item))
			}
		}
	case []string:
		for _, item := range value {
			if strings.TrimSpace(item) != "" {
				query.Add(key, strings.TrimSpace(item))
			}
		}
	case string:
		for _, item := range splitCSV(value) {
			query.Add(key, item)
		}
	}
}

func setOptionalInt(query url.Values, params map[string]any, key string) {
	setOptionalIntAs(query, params, key, key)
}

func setOptionalIntAs(query url.Values, params map[string]any, fromKey, toKey string) {
	switch value := params[fromKey].(type) {
	case float64:
		query.Set(toKey, strconv.Itoa(int(value)))
	case int:
		query.Set(toKey, strconv.Itoa(value))
	case int64:
		query.Set(toKey, strconv.FormatInt(value, 10))
	case string:
		if value != "" {
			query.Set(toKey, value)
		}
	}
}

func setOptionalString(query url.Values, params map[string]any, key string) {
	value, _ := params[key].(string)
	if value != "" {
		query.Set(key, value)
	}
}

func setOptionalBool(query url.Values, params map[string]any, key string) {
	value, ok := params[key].(bool)
	if ok {
		query.Set(key, strconv.FormatBool(value))
	}
}
