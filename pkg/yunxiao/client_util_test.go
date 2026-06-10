package yunxiao

import (
	"context"
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
	got := PrettyJSON([]byte(`{"key":"value"}`))
	if !strings.Contains(got, `"key"`) || !strings.Contains(got, `"value"`) {
		t.Fatalf("PrettyJSON() = %q, want formatted", got)
	}
}

func TestPrettyJSONReturnsRawForInvalidJSON(t *testing.T) {
	invalid := []byte(`{invalid`)
	if got := PrettyJSON(invalid); got != string(invalid) {
		t.Fatalf("PrettyJSON() = %q, want raw", got)
	}
}

func TestPrettyResponseJSONWrapsInvalidBodyAsString(t *testing.T) {
	resp := &Response{Body: []byte(`{invalid`)}
	got := PrettyResponseJSON(resp)
	if !strings.Contains(got, `"data"`) || !strings.Contains(got, "{invalid") {
		t.Fatalf("PrettyResponseJSON() = %q", got)
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

func TestClientPostJSONWithMetadataReturnsResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		w.Header().Set("x-request-id", "request-1")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.PostJSONWithMetadata(context.Background(), "/platform/users:search", map[string]any{"query": "alice"})
	if err != nil {
		t.Fatalf("PostJSONWithMetadata() error = %v", err)
	}
	for _, want := range []string{`"ok": true`, `"requestId": "request-1"`} {
		if !strings.Contains(result, want) {
			t.Fatalf("result = %q, missing %s", result, want)
		}
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

func TestClientPutJSONWithMetadataReturnsResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s, want PUT", r.Method)
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.PutJSONWithMetadata(context.Background(), "/platform/users:me", map[string]any{"key": "value"})
	if err != nil {
		t.Fatalf("PutJSONWithMetadata() error = %v", err)
	}
	if !strings.Contains(result, `"ok": true`) {
		t.Fatalf("result = %q", result)
	}
}
