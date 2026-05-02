package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListAppReleaseWorkflowsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/releaseWorkflows" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"sn":"rw-1"}]`))
	})

	if _, err := handleListAppReleaseWorkflows(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	}); err != nil {
		t.Fatalf("handleListAppReleaseWorkflows() error = %v", err)
	}
}

func TestHandleListAppReleaseWorkflowBriefsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/releaseWorkflowBriefs" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"sn":"rw-1"}]`))
	})

	if _, err := handleListAppReleaseWorkflowBriefs(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	}); err != nil {
		t.Fatalf("handleListAppReleaseWorkflowBriefs() error = %v", err)
	}
}

func TestHandleGetAppReleaseWorkflowStageBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/apps/app%2F1/releaseWorkflow/rw%2F1/releaseStage/rs%2F1") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"sn":"rs-1"}`))
	})

	if _, err := handleGetAppReleaseWorkflowStage(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app/1",
		"releaseWorkflowSn": "rw/1",
		"releaseStageSn":    "rs/1",
	}); err != nil {
		t.Fatalf("handleGetAppReleaseWorkflowStage() error = %v", err)
	}
}

func TestHandleListAppReleaseStageBriefsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/releaseWorkflow/rw-1/releaseStageBriefs" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"sn":"rs-1"}]`))
	})

	if _, err := handleListAppReleaseStageBriefs(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app-1",
		"releaseWorkflowSn": "rw-1",
	}); err != nil {
		t.Fatalf("handleListAppReleaseStageBriefs() error = %v", err)
	}
}

func TestHandleListAppReleaseStageRunsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/releaseWorkflows/rw-1/releaseStages/rs-1/executions" {
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
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	if _, err := handleListAppReleaseStageRuns(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app-1",
		"releaseWorkflowSn": "rw-1",
		"releaseStageSn":    "rs-1",
		"pagination":        "keyset",
		"perPage":           float64(20),
		"orderBy":           "id",
		"sort":              "asc",
		"nextToken":         "token-1",
		"page":              float64(2),
	}); err != nil {
		t.Fatalf("handleListAppReleaseStageRuns() error = %v", err)
	}
}

func TestHandleListAppReleaseStageExecMetadataBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/releaseWorkflows/rw-1/releaseStages/rs-1/executions/1/integratedMetadata" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"changeRequestSn":"cr-1"}]`))
	})

	if _, err := handleListAppReleaseStageExecMetadata(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app-1",
		"releaseWorkflowSn": "rw-1",
		"releaseStageSn":    "rs-1",
		"executionNumber":   "1",
	}); err != nil {
		t.Fatalf("handleListAppReleaseStageExecMetadata() error = %v", err)
	}
}

func TestHandleGetAppReleaseStagePipelineRunBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/releaseWorkflows/rw-1/releaseStages/rs-1/executions/1:getPipelineRun" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"pipelineId":123}`))
	})

	if _, err := handleGetAppReleaseStagePipelineRun(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app-1",
		"releaseWorkflowSn": "rw-1",
		"releaseStageSn":    "rs-1",
		"executionNumber":   "1",
	}); err != nil {
		t.Fatalf("handleGetAppReleaseStagePipelineRun() error = %v", err)
	}
}

func TestHandleGetAppReleaseStagePipelineJobLogBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/apps/app%2F1/releaseWorkflows/rw%2F1/releaseStages/rs%2F1/executions/1%2F2:pipelineJobLog") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		if r.URL.Query().Get("jobId") != "job/1" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"log":"ok"}`))
	})

	if _, err := handleGetAppReleaseStagePipelineJobLog(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app/1",
		"releaseWorkflowSn": "rw/1",
		"releaseStageSn":    "rs/1",
		"executionNumber":   "1/2",
		"jobId":             "job/1",
	}); err != nil {
		t.Fatalf("handleGetAppReleaseStagePipelineJobLog() error = %v", err)
	}
}

func TestHandleGetAppReleaseStagePipelineJobLogRequiresJobId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleGetAppReleaseStagePipelineJobLog(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app-1",
		"releaseWorkflowSn": "rw-1",
		"releaseStageSn":    "rs-1",
		"executionNumber":   "1",
	}); err == nil {
		t.Fatal("expected missing jobId error")
	}
}

func TestRequiredAppReleaseWorkflowRequiresReleaseWorkflowSn(t *testing.T) {
	_, _, _, err := requiredAppReleaseWorkflow(map[string]any{"organizationId": "org-1", "appName": "app-1"})
	if err == nil {
		t.Fatal("expected missing releaseWorkflowSn error")
	}
}

func TestRequiredAppReleaseStageRequiresReleaseStageSn(t *testing.T) {
	_, _, _, _, err := requiredAppReleaseStage(map[string]any{
		"organizationId":    "org-1",
		"appName":           "app-1",
		"releaseWorkflowSn": "rw-1",
	})
	if err == nil {
		t.Fatal("expected missing releaseStageSn error")
	}
}

func TestRequiredAppReleaseStageExecutionRequiresExecutionNumber(t *testing.T) {
	_, _, _, _, _, err := requiredAppReleaseStageExecution(map[string]any{
		"organizationId":    "org-1",
		"appName":           "app-1",
		"releaseWorkflowSn": "rw-1",
		"releaseStageSn":    "rs-1",
	})
	if err == nil {
		t.Fatal("expected missing executionNumber error")
	}
}
