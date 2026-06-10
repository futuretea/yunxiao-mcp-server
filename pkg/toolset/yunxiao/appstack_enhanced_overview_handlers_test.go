package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
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
		if callCount > 1 {
			t.Fatalf("unexpected request %d", callCount)
		}
		_, _ = w.Write([]byte(`{"sn":"co-1"}`))
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

func TestHandleGetSystemOverviewWithoutApps(t *testing.T) {
	callCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 { //nolint:staticcheck
			_, _ = w.Write([]byte(`{"name":"sys-1"}`))
		} else if callCount == 2 {
			if !strings.Contains(r.URL.Path, "/members") {
				t.Fatalf("expected members path, got %s", r.URL.Path)
			}
			_, _ = w.Write([]byte(`[{"userId":"u1"}]`))
		} else {
			t.Fatalf("unexpected request %d", callCount)
		}
	})

	result, err := handleGetSystemOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"systemName":     "sys-1",
		"includeApps":    false,
	})
	if err != nil {
		t.Fatalf("handleGetSystemOverview() error = %v", err)
	}
	if callCount != 2 {
		t.Fatalf("callCount = %d, want 2 (system+members, no apps)", callCount)
	}

	var overview map[string]any
	if err := json.Unmarshal([]byte(result), &overview); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if _, ok := overview["apps"]; ok {
		t.Fatal("overview should not include apps when includeApps=false")
	}
}

func TestHandleGetSystemOverviewAPIErrorOnSystem(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := handleGetSystemOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"systemName":     "sys-1",
	})
	if err == nil {
		t.Fatal("expected API error on system info")
	}
}

func TestHandleGetChangeOrderOverviewAPIErrorOnChangeOrder(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := handleGetChangeOrderOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"changeOrderSn":  "co-1",
	})
	if err == nil {
		t.Fatal("expected API error on change order")
	}
}

func TestHandleGetChangeOrderOverviewRequiresAppName(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetChangeOrderOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	}); err == nil {
		t.Fatal("expected missing appName error")
	}
	if _, err := handleGetChangeOrderOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	}); err == nil {
		t.Fatal("expected missing changeOrderSn error")
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

func TestHandleGetAppReleaseWorkflowOverviewCombinesResponses(t *testing.T) {
	callCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		switch callCount {
		case 1:
			_, _ = w.Write([]byte(`{"name":"wf-1"}`))
		case 2:
			if !strings.Contains(r.URL.Path, "/releaseStageBriefs") {
				t.Fatalf("path = %q", r.URL.Path)
			}
			_, _ = w.Write([]byte(`[{"sn":"stage-1"}]`))
		default:
			t.Fatalf("unexpected request %d", callCount)
		}
	})

	result, err := handleGetAppReleaseWorkflowOverview(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app-1",
		"releaseWorkflowSn": "wf-1",
	})
	if err != nil {
		t.Fatalf("handleGetAppReleaseWorkflowOverview() error = %v", err)
	}

	var overview map[string]any
	if err := json.Unmarshal([]byte(result), &overview); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if _, ok := overview["workflow"]; !ok {
		t.Fatal("overview missing workflow")
	}
	if _, ok := overview["stageBriefs"]; !ok {
		t.Fatal("overview missing stageBriefs")
	}
}

func TestHandleGetAppReleaseStageOverviewCombinesResponses(t *testing.T) {
	callCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		switch callCount {
		case 1:
			if !strings.Contains(r.URL.Path, "/releaseStages/stage-1") {
				t.Fatalf("unexpected stage path: %s", r.URL.Path)
			}
			_, _ = w.Write([]byte(`{"name":"stage-1"}`))
		case 2:
			if !strings.Contains(r.RequestURI, ":getPipelineRun") {
				t.Fatalf("unexpected pipelineRun path: %s", r.URL.Path)
			}
			_, _ = w.Write([]byte(`{"id":"run-1"}`))
		case 3:
			if !strings.Contains(r.URL.Path, "/integratedMetadata") {
				t.Fatalf("unexpected metadata path: %s", r.URL.Path)
			}
			_, _ = w.Write([]byte(`{"key":"value"}`))
		default:
			t.Fatalf("unexpected request %d", callCount)
		}
	})

	result, err := handleGetAppReleaseStageOverview(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app-1",
		"releaseWorkflowSn": "wf-1",
		"releaseStageSn":    "stage-1",
		"executionNumber":   "1",
	})
	if err != nil {
		t.Fatalf("handleGetAppReleaseStageOverview() error = %v", err)
	}
	if callCount != 3 {
		t.Fatalf("callCount = %d, want 3", callCount)
	}

	var overview map[string]any
	if err := json.Unmarshal([]byte(result), &overview); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	for _, key := range []string{"stage", "pipelineRun", "metadata", "filters"} {
		if _, ok := overview[key]; !ok {
			t.Fatalf("overview missing key %q", key)
		}
	}
}

func TestHandleGetAppReleaseStageOverviewWithoutOptionalSections(t *testing.T) {
	callCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount > 1 {
			t.Fatalf("unexpected request %d", callCount)
		}
		_, _ = w.Write([]byte(`{"name":"stage-1"}`))
	})

	result, err := handleGetAppReleaseStageOverview(context.Background(), client, map[string]any{
		"organizationId":     "org-1",
		"appName":            "app-1",
		"releaseWorkflowSn":  "wf-1",
		"releaseStageSn":     "stage-1",
		"executionNumber":    "1",
		"includePipelineRun": false,
		"includeMetadata":    false,
	})
	if err != nil {
		t.Fatalf("handleGetAppReleaseStageOverview() error = %v", err)
	}
	if callCount != 1 {
		t.Fatalf("callCount = %d, want 1", callCount)
	}

	var overview map[string]any
	if err := json.Unmarshal([]byte(result), &overview); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if _, ok := overview["pipelineRun"]; ok {
		t.Fatal("overview should not include pipelineRun when disabled")
	}
	if _, ok := overview["metadata"]; ok {
		t.Fatal("overview should not include metadata when disabled")
	}
}

func TestHandleGetAppReleaseStageOverviewSoftErrors(t *testing.T) {
	callCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		switch callCount {
		case 1:
			_, _ = w.Write([]byte(`{"name":"stage-1"}`))
		case 2:
			w.WriteHeader(http.StatusInternalServerError)
		case 3:
			w.WriteHeader(http.StatusBadGateway)
		default:
			t.Fatalf("unexpected request %d", callCount)
		}
	})

	result, err := handleGetAppReleaseStageOverview(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app-1",
		"releaseWorkflowSn": "wf-1",
		"releaseStageSn":    "stage-1",
		"executionNumber":   "1",
	})
	if err != nil {
		t.Fatalf("handleGetAppReleaseStageOverview() error = %v", err)
	}

	var overview map[string]any
	if err := json.Unmarshal([]byte(result), &overview); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if _, ok := overview["pipelineRunError"]; !ok {
		t.Fatal("overview should include pipelineRunError for soft failure")
	}
	if _, ok := overview["metadataError"]; !ok {
		t.Fatal("overview should include metadataError for soft failure")
	}
}

func TestHandleGetAppReleaseStageOverviewRequiresParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetAppReleaseStageOverview(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetAppReleaseStageOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	}); err == nil {
		t.Fatal("expected missing appName error")
	}
	if _, err := handleGetAppReleaseStageOverview(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app-1",
		"releaseWorkflowSn": "wf-1",
		"releaseStageSn":    "stage-1",
	}); err == nil {
		t.Fatal("expected missing executionNumber error")
	}
}

func TestStageOverviewFilters(t *testing.T) {
	params := map[string]any{
		"includeStageInfo":   false,
		"includePipelineRun": false,
		"includeMetadata":    false,
	}
	filters := stageOverviewFilters(params)
	if filters["includeStageInfo"] != false {
		t.Fatalf("includeStageInfo = %v, want false", filters["includeStageInfo"])
	}
	if filters["includePipelineRun"] != false {
		t.Fatalf("includePipelineRun = %v, want false", filters["includePipelineRun"])
	}
	if filters["includeMetadata"] != false {
		t.Fatalf("includeMetadata = %v, want false", filters["includeMetadata"])
	}
}
