package yunxiao

import (
	"context"
	"net/http"
	"testing"
)

func TestHandleGetAppStackChangeRequestAuditItemsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/apps/app%2F1/changeRequests/cr%2F1/auditItems?refType=CR" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"refType":"CR"}]`))
	})

	if _, err := handleGetAppStackChangeRequestAuditItems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app/1",
		"sn":             "cr/1",
		"refType":        "CR",
	}); err != nil {
		t.Fatalf("handleGetAppStackChangeRequestAuditItems() error = %v", err)
	}
}

func TestHandleListAppStackChangeRequestExecutionsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/changeRequests/cr-1/executions" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("releaseWorkflowSn") != "rw-1" ||
			r.URL.Query().Get("releaseStageSn") != "rs-1" ||
			r.URL.Query().Get("perPage") != "20" ||
			r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("orderBy") != "id" ||
			r.URL.Query().Get("sort") != "asc" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	if _, err := handleListAppStackChangeRequestExecutions(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app-1",
		"sn":                "cr-1",
		"releaseWorkflowSn": "rw-1",
		"releaseStageSn":    "rs-1",
		"perPage":           float64(20),
		"page":              float64(2),
		"orderBy":           "id",
		"sort":              "asc",
	}); err != nil {
		t.Fatalf("handleListAppStackChangeRequestExecutions() error = %v", err)
	}
}

func TestHandleListAppStackChangeRequestWorkItemsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/apps/app%2F1/changeRequests/cr%2F1/workItems" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"name":"work-1"}]`))
	})

	if _, err := handleListAppStackChangeRequestWorkItems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app/1",
		"sn":             "cr/1",
	}); err != nil {
		t.Fatalf("handleListAppStackChangeRequestWorkItems() error = %v", err)
	}
}

func TestHandleGetAppStackChangeRequestAuditItemsRequiresRefType(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleGetAppStackChangeRequestAuditItems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"sn":             "cr-1",
	}); err == nil {
		t.Fatal("expected missing refType error")
	}
}

func TestHandleListAppStackChangeRequestExecutionsRequiresReleaseWorkflowSn(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleListAppStackChangeRequestExecutions(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"sn":             "cr-1",
	}); err == nil {
		t.Fatal("expected missing releaseWorkflowSn error")
	}
}

func TestHandleListAppStackChangeRequestExecutionsRequiresReleaseStageSn(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleListAppStackChangeRequestExecutions(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app-1",
		"sn":                "cr-1",
		"releaseWorkflowSn": "rw-1",
	}); err == nil {
		t.Fatal("expected missing releaseStageSn error")
	}
}

func TestHandleListAppStackChangeRequestWorkItemsRequiresSn(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleListAppStackChangeRequestWorkItems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	}); err == nil {
		t.Fatal("expected missing sn error")
	}
}

func TestRequiredAppStackChangeRequestRequiresSn(t *testing.T) {
	_, _, _, err := requiredAppStackChangeRequest(map[string]any{"organizationId": "org-1", "appName": "app-1"})
	if err == nil {
		t.Fatal("expected missing sn error")
	}
}

func TestAppstackChangeRequestHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetAppStackChangeRequestAuditItems(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleGetAppStackChangeRequestAuditItems(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "appName": "app-1", "sn": "cr-1", "refType": "CR"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListAppStackChangeRequestExecutions(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleListAppStackChangeRequestExecutions(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "appName": "app-1", "sn": "cr-1", "releaseWorkflowSn": "rw-1", "releaseStageSn": "rs-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListAppStackChangeRequestWorkItems(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleListAppStackChangeRequestWorkItems(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "appName": "app-1", "sn": "cr-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}
