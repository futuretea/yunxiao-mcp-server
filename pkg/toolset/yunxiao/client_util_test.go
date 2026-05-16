package yunxiao

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientResolveDefaultOrgIDWithSingleOrganization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/platform/organizations" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"id":"org-1","name":"My Org"}]`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.ResolveDefaultOrgID(context.Background()); err != nil {
		t.Fatalf("ResolveDefaultOrgID() error = %v", err)
	}
	if client.DefaultOrgID != "org-1" {
		t.Fatalf("DefaultOrgID = %q, want org-1", client.DefaultOrgID)
	}
}

func TestClientResolveDefaultOrgIDWithWrappedData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":[{"id":"org-wrapped","name":"Wrapped Org"}]}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.ResolveDefaultOrgID(context.Background()); err != nil {
		t.Fatalf("ResolveDefaultOrgID() error = %v", err)
	}
	if client.DefaultOrgID != "org-wrapped" {
		t.Fatalf("DefaultOrgID = %q, want org-wrapped", client.DefaultOrgID)
	}
}

func TestClientResolveDefaultOrgIDWithMultipleOrganizations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"org-1","name":"Org 1"},{"id":"org-2","name":"Org 2"}]`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.ResolveDefaultOrgID(context.Background()); err != nil {
		t.Fatalf("ResolveDefaultOrgID() error = %v", err)
	}
	if client.DefaultOrgID != "" {
		t.Fatalf("DefaultOrgID = %q, want empty for multiple orgs", client.DefaultOrgID)
	}
}

func TestClientResolveDefaultOrgIDReturnsErrorOnRequestFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.ResolveDefaultOrgID(context.Background()); err == nil {
		t.Fatal("ResolveDefaultOrgID() expected error")
	}
}

func TestClientResolveDefaultOrgIDWithZeroOrganizations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.ResolveDefaultOrgID(context.Background()); err != nil {
		t.Fatalf("ResolveDefaultOrgID() error = %v", err)
	}
	if client.DefaultOrgID != "" {
		t.Fatalf("DefaultOrgID = %q, want empty for zero orgs", client.DefaultOrgID)
	}
}

func TestAPIErrorFormat(t *testing.T) {
	errNoBody := &APIError{Method: "GET", URL: "/path", StatusCode: 404}
	if got := errNoBody.Error(); got != "Yunxiao API error: GET /path returned status 404" {
		t.Fatalf("Error() = %q", got)
	}

	errWithBody := &APIError{Method: "POST", URL: "/path", StatusCode: 500, Body: `{"error":"fail"}`}
	if got := errWithBody.Error(); got != `Yunxiao API error: POST /path returned status 500: {"error":"fail"}` {
		t.Fatalf("Error() = %q", got)
	}
}

func TestAccessTokenFromContext(t *testing.T) {
	if got := AccessTokenFromContext(context.TODO()); got != "" {
		t.Fatalf("AccessTokenFromContext(TODO) = %q, want empty", got)
	}
	ctx := WithAccessToken(context.Background(), "  scoped-token  ")
	if got := AccessTokenFromContext(ctx); got != "scoped-token" {
		t.Fatalf("AccessTokenFromContext(ctx) = %q, want scoped-token", got)
	}
}

func TestWithAccessTokenReturnsOriginalContextForEmptyToken(t *testing.T) {
	ctx := context.Background()
	got := WithAccessToken(ctx, "")
	if got != ctx {
		t.Fatal("WithAccessToken(ctx, \"\") should return original context")
	}
}

func TestPrettyJSONReturnsFormattedForValidJSON(t *testing.T) {
	got := prettyJSON([]byte(`{"key":"value"}`))
	if !strings.Contains(got, `"key"`) || !strings.Contains(got, `"value"`) {
		t.Fatalf("prettyJSON() = %q, want formatted", got)
	}
}

func TestPrettyJSONReturnsRawForInvalidJSON(t *testing.T) {
	invalid := []byte(`{invalid`)
	if got := prettyJSON(invalid); got != string(invalid) {
		t.Fatalf("prettyJSON() = %q, want raw", got)
	}
}

func TestPrettyResponseJSONWrapsInvalidBodyAsString(t *testing.T) {
	resp := &Response{Body: []byte(`{invalid`)}
	got := prettyResponseJSON(resp)
	if !strings.Contains(got, `"data"`) || !strings.Contains(got, "{invalid") {
		t.Fatalf("prettyResponseJSON() = %q", got)
	}
}

func TestNewClientReturnsErrorForInvalidBaseURL(t *testing.T) {
	_, err := NewClient("://invalid-url", "token", time.Second)
	if err == nil {
		t.Fatal("NewClient() expected error for invalid base URL")
	}
}

func TestClientGetJSONWithMetadataReturnsErrorOnRequestFailure(t *testing.T) {
	client, err := NewClient("https://example.com", "token-1", time.Millisecond)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetJSONWithMetadata(context.Background(), "/platform/users:me", nil)
	if err == nil {
		t.Fatal("GetJSONWithMetadata() expected request error")
	}
}

func TestClientPostJSONWithMetadataReturnsErrorOnRequestFailure(t *testing.T) {
	client, err := NewClient("https://example.com", "token-1", time.Millisecond)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.PostJSONWithMetadata(context.Background(), "/platform/users:me", map[string]any{"key": "value"})
	if err == nil {
		t.Fatal("PostJSONWithMetadata() expected request error")
	}
}

func TestClientPutJSONWithMetadataReturnsErrorOnRequestFailure(t *testing.T) {
	client, err := NewClient("https://example.com", "token-1", time.Millisecond)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.PutJSONWithMetadata(context.Background(), "/platform/users:me", map[string]any{"key": "value"})
	if err == nil {
		t.Fatal("PutJSONWithMetadata() expected request error")
	}
}

func TestEncodeRepositoryID(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "numeric", in: "2813489", want: "2813489"},
		{name: "path", in: "org/Demo Repo", want: "org%2FDemo%20Repo"},
		{name: "already encoded", in: "org%2FDemoRepo", want: "org%2FDemoRepo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeRepositoryID(tt.in); got != tt.want {
				t.Fatalf("EncodeRepositoryID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestJoinEscapedPath(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		path     string
		want     string
	}{
		{"both non-empty", "/oapi/v1", "platform/users", "/oapi/v1/platform/users"},
		{"empty base", "", "platform/users", "/platform/users"},
		{"empty path", "/oapi/v1", "", "/oapi/v1"},
		{"both empty", "", "", "/"},
		{"trailing slash", "/oapi/v1/", "/platform/users", "/oapi/v1/platform/users"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinEscapedPath(tt.basePath, tt.path); got != tt.want {
				t.Fatalf("joinEscapedPath(%q, %q) = %q, want %q", tt.basePath, tt.path, got, tt.want)
			}
		})
	}
}

func TestEncodePathValue(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"already encoded", "org%2Frepo", "org%2Frepo"},
		{"normal", "org/repo", "org%2Frepo"},
		{"whitespace", "  org/repo  ", "org%2Frepo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodePathValue(tt.in); got != tt.want {
				t.Fatalf("encodePathValue(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseHeaderInt(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  int
	}{
		{"empty", "", 0},
		{"valid", "42", 42},
		{"invalid", "abc", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.Header{}
			if tt.value != "" {
				h.Set("x-test", tt.value)
			}
			if got := parseHeaderInt(h, "x-test"); got != tt.want {
				t.Fatalf("parseHeaderInt() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestEncodeRepositoryIDEmpty(t *testing.T) {
	if got := EncodeRepositoryID(""); got != "" {
		t.Fatalf("EncodeRepositoryID(\"\") = %q, want empty", got)
	}
}

func TestResolveURLWithInvalidEscapeSequence(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.GetJSON(context.Background(), "/bad%ZZ", nil)
	if err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}
	if result != "{}" {
		t.Fatalf("result = %q, want {}", result)
	}
}

func TestClientPostJSONWithMetadataReturnsMarshalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.PostJSONWithMetadata(context.Background(), "/platform/users:me", map[string]any{"key": make(chan int)})
	if err == nil {
		t.Fatal("PostJSONWithMetadata() expected marshal error")
	}
}

type errorReader struct{}

func (e errorReader) Read(_ []byte) (int, error) {
	return 0, fmt.Errorf("read error")
}
func (e errorReader) Close() error { return nil }

type errorBodyTransport struct{}

func (t errorBodyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       errorReader{},
		Header:     http.Header{},
		Request:    req,
	}, nil
}

func TestClientRequestReturnsReadError(t *testing.T) {
	client, err := NewClient("https://example.com", "token", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	client.httpClient = &http.Client{Transport: errorBodyTransport{}}

	_, err = client.Request(context.Background(), http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("Request() expected read error")
	}
	if !strings.Contains(err.Error(), "read response body") {
		t.Fatalf("error = %v", err)
	}
}

type blockingTransport struct {
	blockUntil context.Context
}

func (t *blockingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	select {
	case <-t.blockUntil.Done():
		return nil, t.blockUntil.Err()
	case <-req.Context().Done():
		return nil, req.Context().Err()
	}
}

func TestClientRequestRespectsContextCancellation(t *testing.T) {
	client, err := NewClient("https://example.com", "token", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	blockCtx, blockCancel := context.WithCancel(context.Background())
	defer blockCancel()
	client.httpClient = &http.Client{Transport: &blockingTransport{blockUntil: blockCtx}}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = client.Request(ctx, http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("Request() expected context cancellation error")
	}
	if !strings.Contains(err.Error(), "request failed") {
		t.Fatalf("error = %v", err)
	}
}
