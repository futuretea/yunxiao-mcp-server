package yunxiao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func handleCallYunxiaoAPI(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path, err := requiredString(params, "path")
	if err != nil {
		return "", err
	}

	method := optionalStringDefault(params, "method", "GET")
	method = strings.ToUpper(method)
	if method != "GET" && method != "POST" {
		return "", fmt.Errorf("method must be GET or POST, got %s", method)
	}

	if isBlockedPath(path) {
		return "", fmt.Errorf("path %q is blocked: potential mutation endpoint", path)
	}

	var query url.Values
	if q, ok := params["queryParams"].(string); ok && q != "" {
		var qmap map[string]string
		if err := json.Unmarshal([]byte(q), &qmap); err != nil {
			return "", fmt.Errorf("invalid queryParams JSON: %w", err)
		}
		query = url.Values{}
		for k, v := range qmap {
			query.Set(k, v)
		}
	}

	var body any
	if b, ok := params["body"].(string); ok && b != "" {
		if err := json.Unmarshal([]byte(b), &body); err != nil {
			return "", fmt.Errorf("invalid body JSON: %w", err)
		}
	}

	resp, err := c.Request(ctx, method, path, query, body)
	if err != nil {
		return "", err
	}
	return prettyResponseJSON(resp), nil
}

func isBlockedPath(path string) bool {
	lower := strings.ToLower(path)
	blocked := []string{
		"/delete", "/update", "/pass", "/refuse", "/approve",
		"/create", "/execute", "/cancel", "/abort", "/reject",
	}
	for _, b := range blocked {
		if strings.Contains(lower, b) {
			return true
		}
	}
	return false
}
