package yunxiao

import (
	"net/url"
	"strconv"
	"strings"
)

func optionalStringDefault(params map[string]any, key, defaultValue string) string {
	if value, _ := params[key].(string); value != "" {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return defaultValue
}

func optionalBoolDefault(params map[string]any, key string, defaultValue bool) bool {
	if value, ok := params[key].(bool); ok {
		return value
	}
	return defaultValue
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

func trimmedNonEmptyStrings(value any) []string {
	switch v := value.(type) {
	case []any:
		return filterNonEmptyStrings(v, func(item any) string {
			if s, ok := item.(string); ok {
				return s
			}
			return ""
		})
	case []string:
		return filterNonEmptyStrings(v, func(item string) string { return item })
	}
	return nil
}

func filterNonEmptyStrings[T any](items []T, get func(T) string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		if s := strings.TrimSpace(get(item)); s != "" {
			result = append(result, s)
		}
	}
	return result
}

func setOptionalStringArrayBody(body map[string]any, params map[string]any, key string) {
	if values := trimmedNonEmptyStrings(params[key]); len(values) > 0 {
		body[key] = values
	}
}

func setOptionalStringArrayQuery(query url.Values, params map[string]any, key string) {
	switch value := params[key].(type) {
	case []any, []string:
		for _, item := range trimmedNonEmptyStrings(value) {
			query.Add(key, item)
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

func pageOneLimitQuery(params map[string]any, limitKey string, defaultLimit int) url.Values {
	query := url.Values{}
	query.Set("page", "1")
	query.Set("perPage", strconv.Itoa(optionalIntDefault(params, limitKey, defaultLimit)))
	return query
}
