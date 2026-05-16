package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandlePassPipelineValidate(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/oapi/v1/flow/organizations/org-1/pipelines/pipeline-1/pipelineRuns/run-1/jobs/job-1/pass" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"result":true}`))
	})

	result, err := handlePassPipelineValidate(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
		"pipelineRunId":  "run-1",
		"jobId":          "job-1",
	})
	if err != nil {
		t.Fatalf("handlePassPipelineValidate() error = %v", err)
	}
	if !strings.Contains(result, "true") {
		t.Fatalf("result = %q", result)
	}
}

func TestHandlePassPipelineValidateMissingOrganizationId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handlePassPipelineValidate(context.Background(), client, map[string]any{
		"pipelineId":    "pipeline-1",
		"pipelineRunId": "run-1",
		"jobId":         "job-1",
	})
	if err == nil {
		t.Fatal("expected error for missing organizationId")
	}
}

func TestHandlePassPipelineValidateMissingPipelineId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handlePassPipelineValidate(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineRunId":  "run-1",
		"jobId":          "job-1",
	})
	if err == nil {
		t.Fatal("expected error for missing pipelineId")
	}
}

func TestHandlePassPipelineValidateMissingPipelineRunId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handlePassPipelineValidate(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
		"jobId":          "job-1",
	})
	if err == nil {
		t.Fatal("expected error for missing pipelineRunId")
	}
}

func TestHandlePassPipelineValidateMissingJobId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handlePassPipelineValidate(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
		"pipelineRunId":  "run-1",
	})
	if err == nil {
		t.Fatal("expected error for missing jobId")
	}
}

func TestHandlePassPipelineValidateNilClient(t *testing.T) {
	_, err := handlePassPipelineValidate(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
		"pipelineRunId":  "run-1",
		"jobId":          "job-1",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestHandlePassPipelineValidateAPIError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handlePassPipelineValidate(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
		"pipelineRunId":  "run-1",
		"jobId":          "job-1",
	})
	if err == nil {
		t.Fatal("expected error for API failure")
	}
}

func TestHandleRefusePipelineValidate(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/oapi/v1/flow/organizations/org-1/pipelines/pipeline-1/pipelineRuns/run-1/jobs/job-1/refuse" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"result":true}`))
	})

	result, err := handleRefusePipelineValidate(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
		"pipelineRunId":  "run-1",
		"jobId":          "job-1",
	})
	if err != nil {
		t.Fatalf("handleRefusePipelineValidate() error = %v", err)
	}
	if !strings.Contains(result, "true") {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleRefusePipelineValidateMissingOrganizationId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleRefusePipelineValidate(context.Background(), client, map[string]any{
		"pipelineId":    "pipeline-1",
		"pipelineRunId": "run-1",
		"jobId":         "job-1",
	})
	if err == nil {
		t.Fatal("expected error for missing organizationId")
	}
}

func TestHandleRefusePipelineValidateMissingPipelineId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleRefusePipelineValidate(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineRunId":  "run-1",
		"jobId":          "job-1",
	})
	if err == nil {
		t.Fatal("expected error for missing pipelineId")
	}
}

func TestHandleRefusePipelineValidateMissingPipelineRunId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleRefusePipelineValidate(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
		"jobId":          "job-1",
	})
	if err == nil {
		t.Fatal("expected error for missing pipelineRunId")
	}
}

func TestHandleRefusePipelineValidateMissingJobId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleRefusePipelineValidate(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
		"pipelineRunId":  "run-1",
	})
	if err == nil {
		t.Fatal("expected error for missing jobId")
	}
}

func TestHandleRefusePipelineValidateNilClient(t *testing.T) {
	_, err := handleRefusePipelineValidate(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
		"pipelineRunId":  "run-1",
		"jobId":          "job-1",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestHandleRefusePipelineValidateAPIError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleRefusePipelineValidate(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
		"pipelineRunId":  "run-1",
		"jobId":          "job-1",
	})
	if err == nil {
		t.Fatal("expected error for API failure")
	}
}
