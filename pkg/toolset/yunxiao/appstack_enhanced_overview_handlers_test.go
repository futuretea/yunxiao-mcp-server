package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestHandleGetSystemOverviewBuildsAndCombinesResponses(t *testing.T) {
	callCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		switch callCount {
		case 1:
			if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/systems/sys-1" {
				t.Fatalf("unexpected system path: %s", r.URL.Path)
			}
			_, _ = w.Write([]byte(`{"name":"sys-1","displayName":"System 1"}`))
		case 2:
			if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/systems/sys-1/apps" {
				t.Fatalf("unexpected apps path: %s", r.URL.Path)
			}
			_, _ = w.Write([]byte(`[{"appName":"app-1"}]`))
		case 3:
			if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/systems/sys-1/members" {
				t.Fatalf("unexpected members path: %s", r.URL.Path)
			}
			_, _ = w.Write([]byte(`[{"userId":"user-1"}]`))
		default:
			t.Fatalf("unexpected request %d", callCount)
		}
	})

	result, err := handleGetSystemOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"systemName":     "sys-1",
	})
	if err != nil {
		t.Fatalf("handleGetSystemOverview() error = %v", err)
	}
	if callCount != 3 {
		t.Fatalf("callCount = %d, want 3", callCount)
	}

	var overview map[string]any
	if err := json.Unmarshal([]byte(result), &overview); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	for _, key := range []string{"system", "apps", "members", "filters"} {
		if _, ok := overview[key]; !ok {
			t.Fatalf("overview missing key %q", key)
		}
	}
}

func TestHandleGetChangeOrderOverviewBuildsAndCombinesResponses(t *testing.T) {
	callCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		switch callCount {
		case 1:
			if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/changeOrders/co-1" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			_, _ = w.Write([]byte(`{"sn":"co-1","status":"SUCCESS"}`))
		case 2:
			_, _ = w.Write([]byte(`[{"jobSn":"job-1"}]`))
		default:
			t.Fatalf("unexpected request %d", callCount)
		}
	})

	result, err := handleGetChangeOrderOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"changeOrderSn":  "co-1",
	})
	if err != nil {
		t.Fatalf("handleGetChangeOrderOverview() error = %v", err)
	}
	if callCount != 2 {
		t.Fatalf("callCount = %d, want 2", callCount)
	}

	var overview map[string]any
	if err := json.Unmarshal([]byte(result), &overview); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	for _, key := range []string{"changeOrder", "jobs", "filters"} {
		if _, ok := overview[key]; !ok {
			t.Fatalf("overview missing key %q", key)
		}
	}
}

func TestHandleGetChangeOrderOverviewWithoutJobs(t *testing.T) {
	callCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			_, _ = w.Write([]byte(`{"sn":"co-1"}`))
		} else {
			t.Fatalf("unexpected request %d", callCount)
		}
	})

	result, err := handleGetChangeOrderOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"changeOrderSn":  "co-1",
		"includeJobLogs": false,
	})
	if err != nil {
		t.Fatalf("handleGetChangeOrderOverview() error = %v", err)
	}

	var overview map[string]any
	if err := json.Unmarshal([]byte(result), &overview); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if _, ok := overview["jobs"]; ok {
		t.Fatal("overview should not include jobs when includeJobLogs=false")
	}
}

func TestEnhancedOverviewHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetSystemOverview(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetSystemOverview(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing systemName error")
	}
	if _, err := handleGetChangeOrderOverview(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetChangeOrderOverview(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "appName": "app-1", "changeOrderSn": "co-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}
