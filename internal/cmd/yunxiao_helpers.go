package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
)

func setCLIStringParam(params map[string]any, key, value string) {
	if value = strings.TrimSpace(value); value != "" {
		params[key] = value
	}
}

func rowsFromJSON(raw string) []any {
	var payload any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil
	}
	return firstJSONArray(payload)
}

func firstJSONArray(value any) []any {
	switch typed := value.(type) {
	case []any:
		return typed
	case map[string]any:
		for _, key := range []string{"data", "result", "items", "workitems", "workItems", "list", "records", "content"} {
			if nested, ok := typed[key]; ok {
				if items := firstJSONArray(nested); len(items) > 0 {
					return items
				}
			}
		}
	}
	return nil
}

func firstStringValue(m map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := m[key]; ok {
			if s := stringifyCLIValue(value); s != "" {
				return s
			}
		}
	}
	return ""
}

func stringifyCLIValue(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return fmt.Sprintf("%.0f", typed)
	case map[string]any:
		return firstStringValue(typed, "name", "displayName", "nickName", "id", "identifier")
	default:
		return ""
	}
}
