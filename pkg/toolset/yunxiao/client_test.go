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

func TestClientGetJSONWithMetadataIncludesPaginationHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("x-page", "2")
		w.Header().Set("x-per-page", "20")
		w.Header().Set("x-next-page", "3")
		w.Header().Set("x-total", "45")
		w.Header().Set("x-total-pages", "3")
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

	for _, want := range []string{`"data"`, `"pagination"`, `"page": 2`, `"nextPage": 3`, `"requestId": "request-1"`} {
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
