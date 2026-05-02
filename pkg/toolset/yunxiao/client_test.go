package yunxiao

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestClientGetJSONAddsAuthHeaderAndAPIBasePath(t *testing.T) {
	var gotPath string
	var gotToken string
	var gotQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotToken = r.Header.Get("x-yunxiao-token")
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	query := url.Values{}
	query.Set("page", "2")
	result, err := client.GetJSON(context.Background(), "/platform/users:me", query)
	if err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}

	if gotPath != "/oapi/v1/platform/users:me" {
		t.Fatalf("path = %q", gotPath)
	}
	if gotToken != "token-1" {
		t.Fatalf("token = %q", gotToken)
	}
	if gotQuery != "page=2" {
		t.Fatalf("query = %q", gotQuery)
	}
	if result != "{\n  \"ok\": true\n}" {
		t.Fatalf("result = %q", result)
	}
}

func TestClientUsesAccessTokenFromContext(t *testing.T) {
	var gotToken string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.Header.Get(AccessTokenHeader)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "default-token", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx := WithAccessToken(context.Background(), "request-token")
	if _, err := client.GetJSON(ctx, "/platform/users:me", nil); err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}

	if gotToken != "request-token" {
		t.Fatalf("token = %q", gotToken)
	}
}

func TestClientUsesContextAccessTokenWithoutDefaultToken(t *testing.T) {
	var gotToken string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.Header.Get(AccessTokenHeader)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx := WithAccessToken(context.Background(), "request-token")
	if _, err := client.GetJSON(ctx, "/platform/users:me", nil); err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}

	if gotToken != "request-token" {
		t.Fatalf("token = %q", gotToken)
	}
}

func TestClientGetJSONWithMetadataIncludesPaginationHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("x-page", "2")
		w.Header().Set("x-per-page", "20")
		w.Header().Set("x-next-page", "3")
		w.Header().Set("x-total", "45")
		w.Header().Set("x-total-pages", "3")
		w.Header().Set("x-next-token", "next-token-1")
		w.Header().Set("x-request-id", "request-1")
		_, _ = w.Write([]byte(`[{"name":"repo"}]`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.GetJSONWithMetadata(context.Background(), "/codeup/organizations/org/repositories", nil)
	if err != nil {
		t.Fatalf("GetJSONWithMetadata() error = %v", err)
	}

	for _, want := range []string{`"data"`, `"pagination"`, `"page": 2`, `"nextPage": 3`, `"nextToken": "next-token-1"`, `"requestId": "request-1"`} {
		if !strings.Contains(result, want) {
			t.Fatalf("result %q does not contain %s", result, want)
		}
	}
}

func TestClientDoesNotDuplicateAPIBasePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/custom/oapi/v1/platform/users:me" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL+"/custom/oapi/v1", "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if _, err := client.GetJSON(context.Background(), "platform/users:me", nil); err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}
}

func TestClientPreservesEscapedRepositoryPath(t *testing.T) {
	var gotRequestURI string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequestURI = r.RequestURI
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if _, err := client.GetJSON(context.Background(), "/codeup/organizations/org/repositories/"+EncodeRepositoryID("group/repo"), nil); err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}

	if !strings.Contains(gotRequestURI, "/repositories/group%2Frepo") {
		t.Fatalf("RequestURI = %q, want escaped repository id", gotRequestURI)
	}
	if strings.Contains(gotRequestURI, "%252F") {
		t.Fatalf("RequestURI = %q, contains double-encoded slash", gotRequestURI)
	}
}

func TestClientRequiresAccessToken(t *testing.T) {
	client, err := NewClient("https://example.com", "", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if _, err := client.GetJSON(context.Background(), "/platform/users:me", nil); err == nil {
		t.Fatal("GetJSON() expected missing access token error")
	}
}

func TestClientReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"denied"}`, http.StatusForbidden)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if _, err := client.GetJSON(context.Background(), "/platform/users:me", nil); err == nil {
		t.Fatal("GetJSON() expected API error")
	}
}

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
	if got := AccessTokenFromContext(nil); got != "" {
		t.Fatalf("AccessTokenFromContext(nil) = %q, want empty", got)
	}
	ctx := WithAccessToken(context.Background(), "scoped-token")
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
