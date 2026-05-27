package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
)

const ansiBold = "\033[1m"
const ansiReset = "\033[0m"

var cliNoColor bool

// SetCLINoColor enables or disables ANSI color output.
func SetCLINoColor(v bool) { cliNoColor = v }

func boldTableHeader(line string) string {
	if cliNoColor {
		return line
	}
	return ansiBold + line + ansiReset
}

func setCLIStringParam(params map[string]any, key, value string) {
	if value = strings.TrimSpace(value); value != "" {
		params[key] = value
	}
}

func rowsFromJSON(raw string) []any {
	items, _ := rowsFromJSONWithPresence(raw)
	return items
}

func rowsFromJSONWithPresence(raw string) ([]any, bool) {
	var payload any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil, false
	}
	return firstJSONArrayWithPresence(payload)
}

func firstJSONArrayWithPresence(value any) ([]any, bool) {
	switch typed := value.(type) {
	case []any:
		return typed, true
	case map[string]any:
		for _, key := range []string{"data", "result", "items", "workitems", "workItems", "list", "records", "content"} {
			if nested, ok := typed[key]; ok {
				if items, ok := firstJSONArrayWithPresence(nested); ok {
					return items, true
				}
			}
		}
	}
	return nil, false
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

func printCLIJSON(out anyWriter, raw string) {
	var obj any
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		_, _ = fmt.Fprintln(out, raw)
		return
	}
	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		_, _ = fmt.Fprintln(out, raw)
		return
	}
	_, _ = fmt.Fprintln(out, string(pretty))
}

func stringifyCLIValue(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return fmt.Sprintf("%.0f", typed)
	case bool:
		return fmt.Sprintf("%t", typed)
	case map[string]any:
		return firstStringValue(typed, "name", "displayName", "nickName", "id", "identifier")
	default:
		return ""
	}
}
