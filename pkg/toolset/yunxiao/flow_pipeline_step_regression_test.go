package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestPipelineJobStepHandlersRejectInvalidQueriesBeforeRequest(t *testing.T) {
	requests := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
	})

	tests := []struct {
		name    string
		handler func(context.Context, any, map[string]any) (string, error)
		params  func() map[string]any
		mutate  func(map[string]any)
	}{
		{"steps missing job", handleGetPipelineJobSteps, pipelineStepParams, func(params map[string]any) { delete(params, "jobId") }},
		{"log missing step index", handleGetPipelineJobStepLog, pipelineStepLogParams, func(params map[string]any) { delete(params, "stepIndex") }},
		{"log fractional step index", handleGetPipelineJobStepLog, pipelineStepLogParams, func(params map[string]any) { params["stepIndex"] = 0.5 }},
		{"log malformed step index", handleGetPipelineJobStepLog, pipelineStepLogParams, func(params map[string]any) { params["stepIndex"] = "two" }},
		{"log missing offset", handleGetPipelineJobStepLog, pipelineStepLogParams, func(params map[string]any) { delete(params, "offset") }},
		{"log fractional offset", handleGetPipelineJobStepLog, pipelineStepLogParams, func(params map[string]any) { params["offset"] = 0.5 }},
		{"log malformed offset", handleGetPipelineJobStepLog, pipelineStepLogParams, func(params map[string]any) { params["offset"] = "zero" }},
		{"log missing limit", handleGetPipelineJobStepLog, pipelineStepLogParams, func(params map[string]any) { delete(params, "limit") }},
		{"log fractional limit", handleGetPipelineJobStepLog, pipelineStepLogParams, func(params map[string]any) { params["limit"] = 1.5 }},
		{"log malformed limit", handleGetPipelineJobStepLog, pipelineStepLogParams, func(params map[string]any) { params["limit"] = "ten" }},
		{"log missing build ID", handleGetPipelineJobStepLog, pipelineStepLogParams, func(params map[string]any) { delete(params, "buildId") }},
		{"log malformed build ID", handleGetPipelineJobStepLog, pipelineStepLogParams, func(params map[string]any) { params["buildId"] = "12x" }},
		{"URL missing step index", handleGetPipelineJobStepLogURL, pipelineStepLogURLParams, func(params map[string]any) { delete(params, "stepIndex") }},
		{"URL fractional step index", handleGetPipelineJobStepLogURL, pipelineStepLogURLParams, func(params map[string]any) { params["stepIndex"] = 0.5 }},
		{"URL malformed step index", handleGetPipelineJobStepLogURL, pipelineStepLogURLParams, func(params map[string]any) { params["stepIndex"] = "two" }},
		{"URL missing build ID", handleGetPipelineJobStepLogURL, pipelineStepLogURLParams, func(params map[string]any) { delete(params, "buildId") }},
		{"URL malformed build ID", handleGetPipelineJobStepLogURL, pipelineStepLogURLParams, func(params map[string]any) { params["buildId"] = "-1" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := tt.params()
			tt.mutate(params)
			if _, err := tt.handler(context.Background(), client, params); err == nil {
				t.Fatal("expected validation error")
			}
			if requests != 0 {
				t.Fatalf("invalid input issued %d requests, want 0", requests)
			}
		})
	}
}

func TestPipelineJobStepHandlersAcceptDecimalStringQueries(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("stepIndex") != "2" || r.URL.Query().Get("buildId") != "9007199254740993" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		if r.URL.Path == "/oapi/v1/flow/organizations/org-1/pipelines/pipeline-1/pipelineRuns/run-1/jobs/job-1/step/log" {
			if r.URL.Query().Get("offset") != "0" || r.URL.Query().Get("limit") != "200" {
				t.Fatalf("log query = %q", r.URL.RawQuery)
			}
		}
		_, _ = w.Write([]byte(`{}`))
	})

	logParams := pipelineStepLogParams()
	logParams["stepIndex"], logParams["offset"], logParams["limit"] = "2", "0", "200"
	if _, err := handleGetPipelineJobStepLog(context.Background(), client, logParams); err != nil {
		t.Fatalf("handleGetPipelineJobStepLog() error = %v", err)
	}

	urlParams := pipelineStepLogURLParams()
	urlParams["stepIndex"] = "2"
	if _, err := handleGetPipelineJobStepLogURL(context.Background(), client, urlParams); err != nil {
		t.Fatalf("handleGetPipelineJobStepLogURL() error = %v", err)
	}
}

func TestPipelineRunOverviewPreservesMetadataAndLargeJobIDHandoff(t *testing.T) {
	const largeJobID = "9007199254740993"
	requests := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		switch r.URL.Path {
		case "/oapi/v1/flow/organizations/org-1/pipelines/pipeline-1/runs/run-1":
			w.Header().Set("x-page", "1")
			w.Header().Set("x-per-page", "20")
			w.Header().Set("x-total", "1")
			w.Header().Set("x-total-pages", "1")
			w.Header().Set("x-next-token", "next-token")
			w.Header().Set("x-request-id", "request-id")
			_, _ = w.Write([]byte(`{"stages":[{"stageInfo":{"jobs":[{"id":9007199254740993}]}}]}`))
		case "/oapi/v1/flow/organizations/org-1/pipelines/pipeline-1/pipelineRuns/run-1/jobs/9007199254740993/steps":
			_, _ = w.Write([]byte(`[]`))
		default:
			t.Fatalf("unexpected request to %q", r.URL.Path)
		}
	})

	result, err := handleGetPipelineRunOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
		"pipelineRunId":  "run-1",
		"includeJobs":    false,
	})
	if err != nil {
		t.Fatalf("handleGetPipelineRunOverview() error = %v", err)
	}

	var overview struct {
		Run struct {
			Data struct {
				Stages []struct {
					StageInfo struct {
						Jobs []struct {
							ID json.Number `json:"id"`
						} `json:"jobs"`
					} `json:"stageInfo"`
				} `json:"stages"`
			} `json:"data"`
			Pagination Pagination `json:"pagination"`
			NextToken  string     `json:"nextToken"`
			RequestID  string     `json:"requestId"`
		} `json:"run"`
	}
	decoder := json.NewDecoder(strings.NewReader(result))
	decoder.UseNumber()
	if err := decoder.Decode(&overview); err != nil {
		t.Fatalf("decode overview: %v", err)
	}
	if overview.Run.Pagination.Page != 1 || overview.Run.NextToken != "next-token" || overview.Run.RequestID != "request-id" {
		t.Fatalf("run metadata = %#v", overview.Run)
	}
	jobID := overview.Run.Data.Stages[0].StageInfo.Jobs[0].ID.String()
	if jobID != largeJobID {
		t.Fatalf("job ID = %q, want %q", jobID, largeJobID)
	}
	if _, err := handleGetPipelineJobSteps(context.Background(), client, map[string]any{
		"organizationId": "org-1", "pipelineId": "pipeline-1", "pipelineRunId": "run-1", "jobId": jobID,
	}); err != nil {
		t.Fatalf("handleGetPipelineJobSteps() error = %v", err)
	}
	if requests != 2 {
		t.Fatalf("request count = %d, want 2", requests)
	}
}
