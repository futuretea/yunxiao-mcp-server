package yunxiao

import (
	"encoding/json"
	"strings"
)

func buildProjectConditions(params map[string]any) string {
	var filterConditions []map[string]any
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
	var filterConditions []map[string]any
	if subject, _ := params["subject"].(string); subject != "" {
		filterConditions = append(filterConditions, stringContainsCondition("subject", subject))
	}

	type listCondition struct {
		key, field, className string
	}
	listConditions := []listCondition{
		{"status", "status", "status"},
		{"assignedTo", "assignedTo", "user"},
		{"creator", "creator", "user"},
		{"sprint", "sprint", "sprint"},
		{"workitemType", "workitemType", "workitemType"},
		{"statusStage", "statusStage", "statusStage"},
		{"priority", "priority", "option"},
	}
	for _, lc := range listConditions {
		if value, _ := params[lc.key].(string); value != "" {
			filterConditions = append(filterConditions, listContainsCondition(lc.field, lc.className, splitCSV(value)))
		}
	}

	if tag, _ := params["tag"].(string); tag != "" {
		filterConditions = append(filterConditions, containsCondition("tag", "tag", "multiList", splitCSV(tag)))
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
