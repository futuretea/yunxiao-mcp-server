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

func TestAppstackSystemReleaseHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleListSystemReleaseWorkflows(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListSystemReleaseWorkflows(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "systemName": "sys-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetRelease(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetRelease(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "systemName": "sys-1", "sn": "rel-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListReleaseMembers(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListReleaseMembers(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "systemName": "sys-1", "sn": "rel-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListReleaseProducts(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListReleaseProducts(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "systemName": "sys-1", "sn": "rel-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListAttachedChangeRequests(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListAttachedChangeRequests(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "systemName": "sys-1", "releaseSn": "rel-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListReleaseExecutions(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListReleaseExecutions(context.Background(), client, map[string]any{"organizationId": "org-1", "systemName": "sys-1", "sn": "rel-1"}); err == nil {
		t.Fatal("expected missing releaseWorkflowSn error")
	}
	if _, err := handleListReleaseExecutions(context.Background(), client, map[string]any{"organizationId": "org-1", "systemName": "sys-1", "sn": "rel-1", "releaseWorkflowSn": "rw-1"}); err == nil {
		t.Fatal("expected missing releaseStageSn error")
	}
	if _, err := handleListReleaseExecutions(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "systemName": "sys-1", "sn": "rel-1", "releaseWorkflowSn": "rw-1", "releaseStageSn": "rs-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}

func TestRequiredOrganizationAndSystem(t *testing.T) {
	if _, _, err := requiredOrganizationAndSystem(map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, _, err := requiredOrganizationAndSystem(map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing systemName error")
	}
	org, sys, err := requiredOrganizationAndSystem(map[string]any{"organizationId": "org-1", "systemName": "sys/1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if org != "org-1" || sys != "sys/1" {
		t.Fatalf("unexpected values: %q %q", org, sys)
	}
}

func TestRequiredSystemRelease(t *testing.T) {
	if _, _, _, err := requiredSystemRelease(map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, _, _, err := requiredSystemRelease(map[string]any{"organizationId": "org-1", "systemName": "sys-1"}); err == nil {
		t.Fatal("expected missing sn error")
	}
	org, sys, sn, err := requiredSystemRelease(map[string]any{"organizationId": "org-1", "systemName": "sys-1", "sn": "rel-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if org != "org-1" || sys != "sys-1" || sn != "rel-1" {
		t.Fatalf("unexpected values: %q %q %q", org, sys, sn)
	}
}

func TestRequiredSystemReleaseWithKey(t *testing.T) {
	if _, _, _, err := requiredSystemReleaseWithKey(map[string]any{"organizationId": "org-1", "systemName": "sys-1"}, "releaseSn"); err == nil {
		t.Fatal("expected missing releaseSn error")
	}
	org, sys, sn, err := requiredSystemReleaseWithKey(map[string]any{"organizationId": "org-1", "systemName": "sys-1", "releaseSn": "rel-1"}, "releaseSn")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if org != "org-1" || sys != "sys-1" || sn != "rel-1" {
		t.Fatalf("unexpected values: %q %q %q", org, sys, sn)
	}
}

func TestAppstackSystemReleasePaths(t *testing.T) {
	if got := appstackSystemPath("org-1", "sys/1"); got != "/appstack/organizations/org-1/systems/sys%2F1" {
		t.Fatalf("appstackSystemPath() = %q", got)
	}
	if got := appstackSystemReleasePath("org-1", "sys/1", "rel/1"); got != "/appstack/organizations/org-1/systems/sys%2F1/releases/rel%2F1" {
		t.Fatalf("appstackSystemReleasePath() = %q", got)
	}
}
