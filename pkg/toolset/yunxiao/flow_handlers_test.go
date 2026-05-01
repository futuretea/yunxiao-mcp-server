package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListPipelinesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/flow/organizations/org-1/pipelines" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("pipelineName") != "deploy" ||
			r.URL.Query().Get("statusList") != "RUNNING,SUCCESS" ||
			r.URL.Query().Get("executeStartTime") != "1000" ||
			r.URL.Query().Get("perPage") != "30" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "2")
		_, _ = w.Write([]byte(`[{"pipelineId":1}]`))
	})

	result, err := handleListPipelines(context.Background(), client, map[string]any{
		"organizationId":   "org-1",
		"pipelineName":     "deploy",
		"statusList":       "RUNNING,SUCCESS",
		"executeStartTime": float64(1000),
		"executeEndTime":   float64(2000),
		"page":             float64(1),
		"perPage":          float64(30),
	})
	if err != nil {
		t.Fatalf("handleListPipelines() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetPipelineBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"pipelineId":1}`))
	})

	if _, err := handleGetPipeline(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
	}); err != nil {
		t.Fatalf("handleGetPipeline() error = %v", err)
	}
}

func TestHandleListPipelineRunsUsesEndTmeQueryName(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("startTime") != "1000" ||
			r.URL.Query().Get("endTme") != "2000" ||
			r.URL.Query().Get("endTime") != "" ||
			r.URL.Query().Get("triggerMode") != "1" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[{"pipelineRunId":1}]`))
	})

	if _, err := handleListPipelineRuns(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
		"startTime":      float64(1000),
		"endTime":        float64(2000),
		"status":         "SUCCESS",
		"triggerMode":    float64(1),
	}); err != nil {
		t.Fatalf("handleListPipelineRuns() error = %v", err)
	}
}

func TestHandleGetPipelineRunBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs/run-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"pipelineRunId":1}`))
	})

	if _, err := handleGetPipelineRun(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
		"pipelineRunId":  "run-1",
	}); err != nil {
		t.Fatalf("handleGetPipelineRun() error = %v", err)
	}
}

func TestHandleGetLatestPipelineRunBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs/latestPipelineRun" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"pipelineRunId":1}`))
	})

	if _, err := handleGetLatestPipelineRun(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
	}); err != nil {
		t.Fatalf("handleGetLatestPipelineRun() error = %v", err)
	}
}

func TestFlowHandlersRequirePipelineParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without required params")
	})

	_, err := handleGetPipeline(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil {
		t.Fatal("handleGetPipeline() expected missing pipelineId error")
	}

	_, err = handleGetPipelineRun(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
	})
	if err == nil {
		t.Fatal("handleGetPipelineRun() expected missing pipelineRunId error")
	}
}
