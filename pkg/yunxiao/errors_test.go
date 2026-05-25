package yunxiao

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"testing"
)

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorCategory
	}{
		{"auth 401", &APIError{StatusCode: http.StatusUnauthorized}, ErrAuth},
		{"permission 403", &APIError{StatusCode: http.StatusForbidden}, ErrPermission},
		{"validation 400", &APIError{StatusCode: http.StatusBadRequest}, ErrValidation},
		{"rate limit 429", &APIError{StatusCode: http.StatusTooManyRequests}, ErrRateLimit},
		{"server 500", &APIError{StatusCode: http.StatusInternalServerError}, ErrServer},
		{"server 502", &APIError{StatusCode: http.StatusBadGateway}, ErrServer},
		{"server 503", &APIError{StatusCode: http.StatusServiceUnavailable}, ErrServer},
		{"uncategorized 404", &APIError{StatusCode: http.StatusNotFound}, ""},
		{"uncategorized 302", &APIError{StatusCode: http.StatusFound}, ""},
		{"non-API error", errors.New("something failed"), ""},
		{"nil error", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyError(tt.err)
			if got != tt.expected {
				t.Fatalf("ClassifyError() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestClassifyErrorNetwork(t *testing.T) {
	netErr := &net.OpError{Op: "dial", Err: errors.New("connection refused")}
	category := ClassifyError(netErr)
	if category != ErrNetwork {
		t.Fatalf("ClassifyError(netError) = %q, want %q", category, ErrNetwork)
	}
}

func TestClassifyErrorWrapped(t *testing.T) {
	apiErr := &APIError{StatusCode: http.StatusUnauthorized}
	wrapped := friendlyAPIError(apiErr)
	category := ClassifyError(wrapped)
	if category != ErrAuth {
		t.Fatalf("ClassifyError(wrapped APIError) = %q, want %q", category, ErrAuth)
	}
}

func TestClassifyErrorValidationError(t *testing.T) {
	valErr := &ValidationError{Msg: "name is required"}
	category := ClassifyError(valErr)
	if category != ErrValidation {
		t.Fatalf("ClassifyError(ValidationError) = %q, want %q", category, ErrValidation)
	}
}

func TestWrapErrorValidationError(t *testing.T) {
	valErr := &ValidationError{Msg: "name is required"}
	got := WrapError(valErr)
	if got == nil || got == valErr {
		t.Fatal("WrapError should wrap ValidationError")
	}
	if !strings.Contains(got.Error(), "[validation]") {
		t.Fatalf("WrapError(ValidationError) = %q, want prefix [validation]", got.Error())
	}
}

func TestFriendlyAPIErrorReturnsNilOnNonAPIError(t *testing.T) {
	err := errors.New("plain error")
	got := friendlyAPIError(err)
	if got != err {
		t.Fatal("friendlyAPIError should return non-API errors unchanged")
	}
}

func TestFriendlyAPIErrorReturnsUncategorizedOnNonStandardStatus(t *testing.T) {
	apiErr := &APIError{StatusCode: http.StatusFound, Body: "redirect"}
	got := friendlyAPIError(apiErr)
	if got != apiErr {
		t.Fatal("friendlyAPIError should return uncategorized APIErrors unchanged")
	}
}

func TestFriendlyAPIErrorAddsSuggestions(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		contains string
	}{
		{"auth", http.StatusUnauthorized, "Authentication failed"},
		{"permission", http.StatusForbidden, "Access denied"},
		{"not found", http.StatusNotFound, "Resource not found"},
		{"bad request", http.StatusBadRequest, "Invalid request parameters"},
		{"rate limit", http.StatusTooManyRequests, "Rate limit exceeded"},
		{"server", http.StatusServiceUnavailable, "temporarily unavailable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FriendlyAPIError(&APIError{StatusCode: tt.status})
			if !strings.Contains(got.Error(), tt.contains) {
				t.Fatalf("FriendlyAPIError() = %q, want %q", got.Error(), tt.contains)
			}
		})
	}
}

func TestTaggedErrorUnwrap(t *testing.T) {
	err := &ValidationError{Msg: "name is required"}
	wrapped := WrapError(err)
	if !errors.Is(wrapped, err) {
		t.Fatalf("WrapError() should unwrap to original error")
	}
}

func TestWrapErrorAddsCategoryPrefix(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		contains string
	}{
		{"auth", &APIError{StatusCode: http.StatusUnauthorized}, "[auth]"},
		{"permission", &APIError{StatusCode: http.StatusForbidden}, "[permission]"},
		{"validation", &APIError{StatusCode: http.StatusBadRequest}, "[validation]"},
		{"rate_limit", &APIError{StatusCode: http.StatusTooManyRequests}, "[rate_limit]"},
		{"server", &APIError{StatusCode: http.StatusInternalServerError}, "[server]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapError(tt.err)
			if got == nil || got == tt.err {
				t.Fatal("WrapError should wrap categorized errors")
			}
			if !strings.Contains(got.Error(), tt.contains) {
				t.Fatalf("WrapError() = %q, want prefix %q", got.Error(), tt.contains)
			}
		})
	}
}

func TestWrapErrorPassthroughUncategorized(t *testing.T) {
	err := errors.New("something failed")
	got := WrapError(err)
	if got != err {
		t.Fatal("WrapError should pass through uncategorized errors unchanged")
	}
}

func TestWrapErrorPassthroughNil(t *testing.T) {
	got := WrapError(nil)
	if got != nil {
		t.Fatal("WrapError(nil) should return nil")
	}
}
