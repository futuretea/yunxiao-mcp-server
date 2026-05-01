package yunxiao

import (
	"context"
	"net/http"
	"testing"
)

func TestHandleListChangeOrderVersionsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/changeOrders/versions" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("envNames") != "dev,test" ||
			r.URL.Query().Get("creators") != "user-1" ||
			r.URL.Query().Get("current") != "2" ||
			r.URL.Query().Get("pageSize") != "10" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	if _, err := handleListChangeOrderVersions(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"envNames":       "dev,test",
		"creators":       "user-1",
		"current":        float64(2),
		"pageSize":       float64(10),
	}); err != nil {
		t.Fatalf("handleListChangeOrderVersions() error = %v", err)
	}
}

func TestHandleListChangeOrderVersionsDefaultsPagination(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("current") != "1" || r.URL.Query().Get("pageSize") != "10" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	if _, err := handleListChangeOrderVersions(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	}); err != nil {
		t.Fatalf("handleListChangeOrderVersions() error = %v", err)
	}
}

func TestHandleGetChangeOrderBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/changeOrders/co-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"sn":"co-1"}`))
	})

	if _, err := handleGetChangeOrder(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"changeOrderSn":  "co-1",
	}); err != nil {
		t.Fatalf("handleGetChangeOrder() error = %v", err)
	}
}

func TestHandleListChangeOrderJobLogsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/changeOrders/co-1/jobs/job-1/logs" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("current") != "3" || r.URL.Query().Get("pageSize") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	if _, err := handleListChangeOrderJobLogs(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"changeOrderSn":  "co-1",
		"jobSn":          "job-1",
		"current":        float64(3),
		"pageSize":       float64(20),
	}); err != nil {
		t.Fatalf("handleListChangeOrderJobLogs() error = %v", err)
	}
}

func TestHandleListChangeOrderJobLogsDefaultsPagination(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("current") != "1" || r.URL.Query().Get("pageSize") != "10" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	if _, err := handleListChangeOrderJobLogs(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"changeOrderSn":  "co-1",
		"jobSn":          "job-1",
	}); err != nil {
		t.Fatalf("handleListChangeOrderJobLogs() error = %v", err)
	}
}

func TestHandleFindTaskOperationLogBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/changeOrders/co-1/jobs/job-1/stages/stage-1/tasks/task-1/operationLog" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`"ok"`))
	})

	if _, err := handleFindTaskOperationLog(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"changeOrderSn":  "co-1",
		"jobSn":          "job-1",
		"stageSn":        "stage-1",
		"taskSn":         "task-1",
	}); err != nil {
		t.Fatalf("handleFindTaskOperationLog() error = %v", err)
	}
}

func TestHandleListChangeOrdersByOriginBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/changeOrders:byOrigin" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("originType") != "FLOW" ||
			r.URL.Query().Get("originId") != "123+1" ||
			r.URL.Query().Get("appName") != "app-1" ||
			r.URL.Query().Get("envName") != "dev" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[{"sn":"co-1"}]`))
	})

	if _, err := handleListChangeOrdersByOrigin(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"originType":     "FLOW",
		"originId":       "123+1",
		"appName":        "app-1",
		"envName":        "dev",
	}); err != nil {
		t.Fatalf("handleListChangeOrdersByOrigin() error = %v", err)
	}
}
