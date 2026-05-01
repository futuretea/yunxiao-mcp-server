package yunxiao

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
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
