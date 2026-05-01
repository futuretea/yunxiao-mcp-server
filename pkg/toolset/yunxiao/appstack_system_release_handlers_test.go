package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListSystemReleaseWorkflowsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/systems/system%2F1/releaseWorkflows" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"sn":"rw-1"}]`))
	})

	if _, err := handleListSystemReleaseWorkflows(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"systemName":     "system/1",
	}); err != nil {
		t.Fatalf("handleListSystemReleaseWorkflows() error = %v", err)
	}
}

func TestHandleGetReleaseBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/systems/system%2F1/releases/release%2F1") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"sn":"release/1"}`))
	})

	if _, err := handleGetRelease(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"systemName":     "system/1",
		"sn":             "release/1",
	}); err != nil {
		t.Fatalf("handleGetRelease() error = %v", err)
	}
}

func TestHandleListReleaseMembersBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/systems/system-1/releases/release-1/members" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	if _, err := handleListReleaseMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"systemName":     "system-1",
		"sn":             "release-1",
	}); err != nil {
		t.Fatalf("handleListReleaseMembers() error = %v", err)
	}
}

func TestHandleListReleaseProductsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/systems/system%2F1/releases/release%2F1/products" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"name":"artifact"}]`))
	})

	if _, err := handleListReleaseProducts(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"systemName":     "system/1",
		"sn":             "release/1",
	}); err != nil {
		t.Fatalf("handleListReleaseProducts() error = %v", err)
	}
}

func TestHandleListAttachedChangeRequestsBuildsPathAndDefaultQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/systems/system%2F1/releases/release%2F1/changeRequests?current=1&pageSize=10" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	if _, err := handleListAttachedChangeRequests(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"systemName":     "system/1",
		"releaseSn":      "release/1",
	}); err != nil {
		t.Fatalf("handleListAttachedChangeRequests() error = %v", err)
	}
}

func TestHandleListReleaseExecutionsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/systems/system%2F1/releases/release%2F1/executions?") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		if r.URL.Query().Get("releaseWorkflowSn") != "rw/1" ||
			r.URL.Query().Get("releaseStageSn") != "rs/1" ||
			r.URL.Query().Get("perPage") != "20" ||
			r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("orderBy") != "id" ||
			r.URL.Query().Get("sort") != "asc" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	if _, err := handleListReleaseExecutions(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"systemName":        "system/1",
		"sn":                "release/1",
		"releaseWorkflowSn": "rw/1",
		"releaseStageSn":    "rs/1",
		"perPage":           float64(20),
		"page":              float64(2),
		"orderBy":           "id",
		"sort":              "asc",
	}); err != nil {
		t.Fatalf("handleListReleaseExecutions() error = %v", err)
	}
}
