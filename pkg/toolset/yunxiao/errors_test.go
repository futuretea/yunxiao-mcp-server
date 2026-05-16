package yunxiao

import (
	"errors"
	"net"
	"net/http"
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
