package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetPipelineOverviewRequiresOrganizationId(t *testing.T) {
	_, err := handleGetPipelineOverview(context.Background(), nil, map[string]any{
		"pipelineId": "pipeline-1",
	})
	if err == nil || !strings.Contains(err.Error(), "organizationId is required") {
		t.Fatalf("expected organizationId required error, got %v", err)
	}
}

func TestHandleGetPipelineOverviewRequiresPipelineId(t *testing.T) {
	_, err := handleGetPipelineOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil || !strings.Contains(err.Error(), "pipelineId is required") {
		t.Fatalf("expected pipelineId required error, got %v", err)
	}
}

func TestHandleGetPipelineOverviewRequiresClient(t *testing.T) {
	_, err := handleGetPipelineOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
	})
	if err == nil || !strings.Contains(err.Error(), "yunxiao client is not configured") {
		t.Fatalf("expected client error, got %v", err)
	}
}

func TestHandleGetPipelineOverviewReturnsErrorOnPipelineFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleGetPipelineOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
	})
	if err == nil || !strings.Contains(err.Error(), "pipeline:") {
		t.Fatalf("expected pipeline error, got %v", err)
	}
}

func TestHandleGetPipelineOverviewReturnsErrorOnLatestRunFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/latestPipelineRun") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"latest run boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"pipeline-1"}`))
	})
	_, err := handleGetPipelineOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
	})
	if err == nil || !strings.Contains(err.Error(), "latestRun:") {
		t.Fatalf("expected latestRun error, got %v", err)
	}
}

func TestHandleGetPipelineOverviewReturnsErrorOnRunsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/runs") && !strings.HasSuffix(r.URL.Path, "/latestPipelineRun") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"runs boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"run-1"}`))
	})
	_, err := handleGetPipelineOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
	})
	if err == nil || !strings.Contains(err.Error(), "runs:") {
		t.Fatalf("expected runs error, got %v", err)
	}
}

func TestHandleGetPipelineOverviewSuccessAllSections(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oapi/v1/flow/organizations/org-1/pipelines/pipeline-1":
			_, _ = w.Write([]byte(`{"id":"pipeline-1"}`))
		case strings.HasSuffix(r.URL.Path, "/latestPipelineRun"):
			_, _ = w.Write([]byte(`{"id":"run-latest"}`))
		case strings.HasSuffix(r.URL.Path, "/runs"):
			if r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("runs perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[{"id":"run-1"}]`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})

	result, err := handleGetPipelineOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
	})
	if err != nil {
		t.Fatalf("handleGetPipelineOverview() error = %v", err)
	}
	if !strings.Contains(result, `"pipeline"`) {
		t.Fatalf("result missing pipeline: %q", result)
	}
	if !strings.Contains(result, `"latestRun"`) {
		t.Fatalf("result missing latestRun: %q", result)
	}
	if !strings.Contains(result, `"runs"`) {
		t.Fatalf("result missing runs: %q", result)
	}
}

func TestHandleGetPipelineOverviewSkipsRunsWhenDisabled(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/runs") && !strings.HasSuffix(r.URL.Path, "/latestPipelineRun") {
			t.Fatalf("unexpected request to %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"x"}`))
	})

	result, err := handleGetPipelineOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
		"includeRuns":    false,
	})
	if err != nil {
		t.Fatalf("handleGetPipelineOverview() error = %v", err)
	}
	if strings.Contains(result, `"runs"`) {
		t.Fatalf("result should not contain runs: %q", result)
	}
}

func TestHandleGetPipelineOverviewUsesCustomRunLimit(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/runs") && !strings.HasSuffix(r.URL.Path, "/latestPipelineRun") {
			if r.URL.Query().Get("perPage") != "3" {
				t.Fatalf("runs perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[]`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"x"}`))
	})

	_, err := handleGetPipelineOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
		"runLimit":       float64(3),
	})
	if err != nil {
		t.Fatalf("handleGetPipelineOverview() error = %v", err)
	}
}

func TestPipelineOverviewFilters(t *testing.T) {
	params := map[string]any{
		"includeRuns": false,
		"runLimit":    float64(10),
	}
	filters := pipelineOverviewFilters(params)
	if filters["includeRuns"].(bool) != false {
		t.Fatalf("includeRuns = %v", filters["includeRuns"])
	}
	if filters["runLimit"].(int) != 10 {
		t.Fatalf("runLimit = %v", filters["runLimit"])
	}
}

func TestPipelineRunQuery(t *testing.T) {
	q := pipelineRunQuery(map[string]any{"runLimit": float64(7)})
	if q.Get("page") != "1" {
		t.Fatalf("page = %q", q.Get("page"))
	}
	if q.Get("perPage") != "7" {
		t.Fatalf("perPage = %q", q.Get("perPage"))
	}

	q2 := pipelineRunQuery(map[string]any{})
	if q2.Get("perPage") != "5" {
		t.Fatalf("default perPage = %q", q2.Get("perPage"))
	}
}
