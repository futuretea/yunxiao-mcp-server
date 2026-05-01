package yunxiao

import (
	"context"
	"net/http"
	"testing"
)

func TestHandleListSystemsBuildsPathAndDefaultQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/systems" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("current") != "1" || r.URL.Query().Get("pageSize") != "10" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	if _, err := handleListSystems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	}); err != nil {
		t.Fatalf("handleListSystems() error = %v", err)
	}
}

func TestHandleListAttachedAppsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/systems/system%2F1/apps?current=2&pageSize=20" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	if _, err := handleListAttachedApps(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"systemName":     "system/1",
		"current":        float64(2),
		"pageSize":       float64(20),
	}); err != nil {
		t.Fatalf("handleListAttachedApps() error = %v", err)
	}
}

func TestHandleListSystemMembersBuildsPathAndDefaultQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/systems/system%2F1/members?current=1&pageSize=10" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	if _, err := handleListSystemMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"systemName":     "system/1",
	}); err != nil {
		t.Fatalf("handleListSystemMembers() error = %v", err)
	}
}
