package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleCreateChangeOrderBuildsPathAndBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/changeOrders" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"sn":"co-1","status":"RUNNING"}`))
	})

	result, err := handleCreateChangeOrder(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"changeOrder":    `{"changeOrderName":"release-1","type":"Deploy"}`,
	})
	if err != nil {
		t.Fatalf("handleCreateChangeOrder() error = %v", err)
	}
	if !strings.Contains(result, "co-1") {
		t.Fatalf("result = %q, want sn:co-1", result)
	}
}

func TestHandleCreateChangeOrderRequiresParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleCreateChangeOrder(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleCreateChangeOrder(context.Background(), client, map[string]any{
		"organizationId": "org-1", "appName": "app-1",
	}); err == nil {
		t.Fatal("expected missing changeOrder error")
	}
	if _, err := handleCreateChangeOrder(context.Background(), client, map[string]any{
		"organizationId": "org-1", "appName": "app-1", "changeOrder": "not-json",
	}); err == nil {
		t.Fatal("expected invalid JSON error")
	}
	if _, err := handleCreateChangeOrder(context.Background(), "invalid-client", map[string]any{
		"organizationId": "org-1", "appName": "app-1", "changeOrder": `{}`,
	}); err == nil {
		t.Fatal("expected getClient error")
	}
}

func TestHandleExecuteJobActionBuildsPathAndBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, ":execute") {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"sn":"co-1"}`))
	})

	result, err := handleExecuteJobAction(context.Background(), client, map[string]any{
		"organizationId": "org-1", "appName": "app-1", "changeOrderSn": "co-1", "jobSn": "job-1",
		"action": `{"actionType":"RESUME"}`,
	})
	if err != nil {
		t.Fatalf("handleExecuteJobAction() error = %v", err)
	}
	if !strings.Contains(result, "co-1") {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleExecuteJobActionRequiresParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request")
	})

	if _, err := handleExecuteJobAction(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleExecuteJobAction(context.Background(), "invalid-client", map[string]any{
		"organizationId": "org-1", "appName": "app-1", "changeOrderSn": "co-1", "jobSn": "job-1", "action": `{}`,
	}); err == nil {
		t.Fatal("expected getClient error")
	}
}

func TestHandleExecuteSystemReleaseStageBuildsPathAndBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, ":execute") {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"pipelineRunId":123}`))
	})

	result, err := handleExecuteSystemReleaseStage(context.Background(), client, map[string]any{
		"organizationId": "org-1", "systemName": "sys-1", "releaseWorkflowSn": "wf-1", "releaseStageSn": "stage-1",
		"execution": `{"params":{"key":"val"}}`,
	})
	if err != nil {
		t.Fatalf("handleExecuteSystemReleaseStage() error = %v", err)
	}
	if !strings.Contains(result, "pipelineRunId") {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleExecuteJobActionReturnsAPIError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	_, err := handleExecuteJobAction(context.Background(), client, map[string]any{
		"organizationId": "org-1", "appName": "app-1", "changeOrderSn": "co-1", "jobSn": "job-1",
		"action": `{"actionType":"RESUME"}`,
	})
	if err == nil {
		t.Fatal("expected API error")
	}
}

func TestHandleExecuteAppReleaseStageBuildsPathAndBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, ":execute") {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"pipelineRunId":456}`))
	})

	result, err := handleExecuteAppReleaseStage(context.Background(), client, map[string]any{
		"organizationId": "org-1", "appName": "app-1", "releaseWorkflowSn": "wf-1", "releaseStageSn": "stage-1",
		"execution": `{}`,
	})
	if err != nil {
		t.Fatalf("handleExecuteAppReleaseStage() error = %v", err)
	}
	if !strings.Contains(result, "pipelineRunId") {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleExecuteSystemReleaseStageRequiresParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request")
	})

	if _, err := handleExecuteSystemReleaseStage(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleExecuteSystemReleaseStage(context.Background(), client, map[string]any{
		"organizationId": "org-1", "systemName": "sys-1", "releaseWorkflowSn": "wf-1", "releaseStageSn": "stage-1",
	}); err == nil {
		t.Fatal("expected missing execution error")
	}
	if _, err := handleExecuteSystemReleaseStage(context.Background(), client, map[string]any{
		"organizationId": "org-1", "systemName": "sys-1", "releaseWorkflowSn": "wf-1", "releaseStageSn": "stage-1",
		"execution": "not-json",
	}); err == nil {
		t.Fatal("expected invalid JSON error")
	}
	if _, err := handleExecuteSystemReleaseStage(context.Background(), "invalid-client", map[string]any{
		"organizationId": "org-1", "systemName": "sys-1", "releaseWorkflowSn": "wf-1", "releaseStageSn": "stage-1",
		"execution": `{}`,
	}); err == nil {
		t.Fatal("expected getClient error")
	}
}

func TestHandleExecuteAppReleaseStageRequiresParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request")
	})

	if _, err := handleExecuteAppReleaseStage(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleExecuteAppReleaseStage(context.Background(), client, map[string]any{
		"organizationId": "org-1", "appName": "app-1", "releaseWorkflowSn": "wf-1", "releaseStageSn": "stage-1",
	}); err == nil {
		t.Fatal("expected missing execution error")
	}
	if _, err := handleExecuteAppReleaseStage(context.Background(), "invalid-client", map[string]any{
		"organizationId": "org-1", "appName": "app-1", "releaseWorkflowSn": "wf-1", "releaseStageSn": "stage-1",
		"execution": `{}`,
	}); err == nil {
		t.Fatal("expected getClient error")
	}
}
