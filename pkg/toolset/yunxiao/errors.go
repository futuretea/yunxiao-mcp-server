package yunxiao

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

// ErrorCategory classifies an API or network error for MCP consumer guidance.
type ErrorCategory string

const (
	ErrAuth       ErrorCategory = "auth"
	ErrPermission ErrorCategory = "permission"
	ErrValidation ErrorCategory = "validation"
	ErrRateLimit  ErrorCategory = "rate_limit"
	ErrServer     ErrorCategory = "server"
	ErrNetwork    ErrorCategory = "network"
)

// ClassifyError returns the error category for an error.
// Uncategorized errors return the zero value (empty string).
func ClassifyError(err error) ErrorCategory {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		switch {
		case apiErr.StatusCode == http.StatusUnauthorized:
			return ErrAuth
		case apiErr.StatusCode == http.StatusForbidden:
			return ErrPermission
		case apiErr.StatusCode == http.StatusBadRequest:
			return ErrValidation
		case apiErr.StatusCode == http.StatusTooManyRequests:
			return ErrRateLimit
		case apiErr.StatusCode >= 500:
			return ErrServer
		}
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return ErrNetwork
	}

	return ""
}

// WrapError prepends an error category tag for MCP consumer pattern matching.
// Uncategorized errors are returned unchanged.
func WrapError(err error) error {
	cat := ClassifyError(err)
	if cat == "" {
		return err
	}
	return fmt.Errorf("[%s] %w", cat, err)
}

// friendlyAPIError wraps an APIError with actionable guidance for LLM consumers.
// Non-API errors are returned unchanged.
func friendlyAPIError(err error) error {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return err
	}

	var suggestion string
	switch apiErr.StatusCode {
	case http.StatusUnauthorized:
		suggestion = "Authentication failed. Verify that your access token is valid and not expired."
	case http.StatusForbidden:
		suggestion = "Access denied. Your token may not have permission for this resource."
	case http.StatusNotFound:
		suggestion = "Resource not found. Verify that the project ID, work item ID, pipeline ID, or other identifiers are correct. Use search_projects, search_workitems, or list_pipelines to find valid IDs."
	case http.StatusBadRequest:
		suggestion = "Invalid request parameters. Check that required fields are present, IDs are correct, and enum values are valid. Use the corresponding get_*_context or list_* tools to discover valid values."
	case http.StatusTooManyRequests:
		suggestion = "Rate limit exceeded. Wait a moment before retrying."
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		suggestion = "Yunxiao service temporarily unavailable. Retry the request later."
	default:
		return err
	}

	return fmt.Errorf("%w\n\nSuggestion: %s", apiErr, suggestion)
}
