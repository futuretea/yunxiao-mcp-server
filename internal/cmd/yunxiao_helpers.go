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


var cliOutputFormat string

// SetCLIOutputFormat sets the global output format for error formatting.
func SetCLIOutputFormat(v string) { cliOutputFormat = v }

// FormatCLIError formats an error for CLI output. When output is json,
// returns a JSON object; otherwise returns the plain error message.
func FormatCLIError(err error) string {
	if cliOutputFormat == "json" {
		data, _ := json.Marshal(map[string]any{
			"error": err.Error(),
		})
		return string(data)
	}
	return "Error: " + err.Error()
}

// ExitCodeFromError maps an error to an exit code based on its category.
func ExitCodeFromError(err error) int {
	cat := classifyCLIError(err)
	switch cat {
	case "auth":
		return 2
	case "permission":
		return 3
	case "validation":
		return 4
	case "rate_limit":
		return 5
	case "network":
		return 6
	case "server":
		return 7
	default:
		return 1
	}
}

func classifyCLIError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	if len(msg) > 0 {
		switch msg[0] {
		case '[':
			if end := strings.IndexByte(msg, ']'); end > 0 {
				return msg[1:end]
			}
		}
	}
	return ""
}
