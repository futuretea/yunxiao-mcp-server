package yunxiao

import (
	"context"
	"net/http"
	"testing"
)

func TestHandleListApplicationsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("pagination") != "keyset" ||
			r.URL.Query().Get("perPage") != "20" ||
			r.URL.Query().Get("orderBy") != "id" ||
			r.URL.Query().Get("sort") != "asc" ||
			r.URL.Query().Get("nextToken") != "token-1" ||
			r.URL.Query().Get("page") != "2" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"data":[{"name":"app"}],"nextToken":"token-2"}`))
	})

	if _, err := handleListApplications(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pagination":     "keyset",
		"perPage":        float64(20),
		"orderBy":        "id",
		"sort":           "asc",
		"nextToken":      "token-1",
		"page":           float64(2),
	}); err != nil {
		t.Fatalf("handleListApplications() error = %v", err)
	}
}

func TestHandleGetApplicationBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"name":"app-1"}`))
	})

	if _, err := handleGetApplication(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	}); err != nil {
		t.Fatalf("handleGetApplication() error = %v", err)
	}
}

func TestAppstackHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleListApplications(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListApplications(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetApplication(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetApplication(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing appName error")
	}
	if _, err := handleGetApplication(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "appName": "app-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}
