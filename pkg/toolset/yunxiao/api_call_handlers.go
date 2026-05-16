package yunxiao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type fallbackAPICallRequest struct {
	path   string
	method string
	query  url.Values
	body   any
}

func handleCallYunxiaoAPI(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	request, err := parseFallbackAPICallRequest(params)
	if err != nil {
		return "", err
	}

	resp, err := c.Request(ctx, request.method, request.path, request.query, request.body)
	if err != nil {
		return "", err
	}
	return prettyResponseJSON(resp), nil
}

func parseFallbackAPICallRequest(params map[string]any) (*fallbackAPICallRequest, error) {
	path, err := requiredString(params, "path")
	if err != nil {
		return nil, err
	}

	method, err := fallbackAPICallMethod(params)
	if err != nil {
		return nil, err
	}
	if err := validateFallbackAPICallPath(path, method); err != nil {
		return nil, err
	}

	query, err := fallbackAPICallQuery(params)
	if err != nil {
		return nil, err
	}
	body, err := fallbackAPICallBody(params)
	if err != nil {
		return nil, err
	}

	return &fallbackAPICallRequest{
		path:   path,
		method: method,
		query:  query,
		body:   body,
	}, nil
}

func fallbackAPICallMethod(params map[string]any) (string, error) {
	method := strings.ToUpper(optionalStringDefault(params, "method", "GET"))
	if method != "GET" && method != "POST" {
		return "", fmt.Errorf("method must be GET or POST, got %s", method)
	}
	return method, nil
}

func validateFallbackAPICallPath(path, method string) error {
	canonicalPath, err := canonicalAPICallPath(path)
	if err != nil {
		return err
	}
	if isBlockedPath(canonicalPath) {
		return fmt.Errorf("path %q is blocked: potential mutation endpoint", path)
	}
	if method == "POST" && !isAllowedReadOnlyPOST(canonicalPath) {
		return fmt.Errorf("POST path %q is blocked: fallback API calls only allow read-only search/list endpoints", path)
	}
	return nil
}

func fallbackAPICallQuery(params map[string]any) (url.Values, error) {
	if q, ok := params["queryParams"].(string); ok && q != "" {
		var qmap map[string]string
		if err := json.Unmarshal([]byte(q), &qmap); err != nil {
			return nil, fmt.Errorf("invalid queryParams JSON: %w", err)
		}
		query := url.Values{}
		for k, v := range qmap {
			query.Set(k, v)
		}
		return query, nil
	}
	return nil, nil
}

func fallbackAPICallBody(params map[string]any) (any, error) {
	if b, ok := params["body"].(string); ok && b != "" {
		var body any
		if err := json.Unmarshal([]byte(b), &body); err != nil {
			return nil, fmt.Errorf("invalid body JSON: %w", err)
		}
		return body, nil
	}
	return nil, nil
}

func canonicalAPICallPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	decoded, err := url.PathUnescape(path)
	if err != nil {
		return "", fmt.Errorf("invalid path escape: %w", err)
	}
	if hasDotPathSegment(decoded) {
		return "", fmt.Errorf("path %q is blocked: dot segments are not allowed", path)
	}
	if containsEscapedOctet(decoded) {
		return "", fmt.Errorf("path %q is blocked: nested path escapes are not allowed", path)
	}
	return decoded, nil
}

func hasDotPathSegment(path string) bool {
	for _, segment := range strings.Split(path, "/") {
		if segment == "." || segment == ".." {
			return true
		}
	}
	return false
}

func containsEscapedOctet(path string) bool {
	for i := 0; i+2 < len(path); i++ {
		if path[i] == '%' && isHex(path[i+1]) && isHex(path[i+2]) {
			return true
		}
	}
	return false
}

func isHex(b byte) bool {
	return (b >= '0' && b <= '9') ||
		(b >= 'a' && b <= 'f') ||
		(b >= 'A' && b <= 'F')
}

func isBlockedPath(path string) bool {
	if hasUnsafeColonAction(path) {
		return true
	}

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

func hasUnsafeColonAction(path string) bool {
	for _, segment := range strings.Split(path, "/") {
		actionStart := strings.LastIndex(segment, ":")
		if actionStart == -1 {
			continue
		}
		switch segment[actionStart+1:] {
		case "me":
			if path == "/platform/users:me" {
				continue
			}
		case "search":
			continue
		}
		return true
	}
	return false
}

func isAllowedReadOnlyPOST(path string) bool {
	lower := strings.ToLower(strings.TrimSpace(path))
	return strings.HasSuffix(lower, ":search") ||
		strings.HasSuffix(lower, "/list") ||
		strings.Contains(lower, "/result/list/")
}
