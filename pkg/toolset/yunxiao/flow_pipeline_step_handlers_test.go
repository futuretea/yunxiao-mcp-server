package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

type pipelineStepToolExpectation struct {
	name     string
	required []string
	types    map[string]string
}

var pipelineStepToolExpectations = []pipelineStepToolExpectation{
	{
		name:     "get_pipeline_job_steps",
		required: []string{"organizationId", "pipelineId", "pipelineRunId", "jobId"},
		types:    map[string]string{"organizationId": "string", "pipelineId": "string", "pipelineRunId": "string", "jobId": "string"},
	},
	{
		name:     "get_pipeline_job_step_log",
		required: []string{"organizationId", "pipelineId", "pipelineRunId", "jobId", "stepIndex", "offset", "limit", "buildId"},
		types:    map[string]string{"stepIndex": "number", "offset": "number", "limit": "number", "buildId": "string"},
	},
	{
		name:     "get_pipeline_job_step_log_url",
		required: []string{"organizationId", "pipelineId", "pipelineRunId", "jobId", "stepIndex", "buildId"},
		types:    map[string]string{"stepIndex": "number", "buildId": "string"},
	},
}

func TestPipelineStepToolCounts(t *testing.T) {
	all := (&Toolset{ReadOnly: false}).GetTools(nil)
	if len(all) != 196 {
		t.Fatalf("tool count = %d, want 196", len(all))
	}

	readOnly := (&Toolset{ReadOnly: true}).GetTools(nil)
	if len(readOnly) != 180 {
		t.Fatalf("read-only tool count = %d, want 180", len(readOnly))
	}

	flowCount := 0
	for _, tool := range all {
		if tool.Domain == "flow" {
			flowCount++
		}
	}
	if flowCount != 21 {
		t.Fatalf("flow tool count = %d, want 21", flowCount)
	}
	flowReadOnlyCount := 0
	for _, tool := range readOnly {
		if tool.Domain == "flow" {
			flowReadOnlyCount++
		}
	}
	if flowReadOnlyCount != 19 {
		t.Fatalf("read-only flow tool count = %d, want 19", flowReadOnlyCount)
	}
	if len(compactHiddenTools) != 62 {
		t.Fatalf("compact hidden tool count = %d, want 62", len(compactHiddenTools))
	}
}

func TestPipelineStepToolSchemasAreReadOnly(t *testing.T) {
	tools := (&Toolset{ReadOnly: true}).GetTools(nil)
	for _, tt := range pipelineStepToolExpectations {
		t.Run(tt.name, func(t *testing.T) {
			tool := pipelineStepToolByName(t, tools, tt.name)
			if tool.domain != "flow" {
				t.Fatalf("domain = %q, want flow", tool.domain)
			}
			if tool.read == nil || !*tool.read {
				t.Fatal("tool must be read-only")
			}

			var schema struct {
				Properties map[string]struct {
					Type string `json:"type"`
				} `json:"properties"`
				Required []string `json:"required"`
			}
			if err := json.Unmarshal(tool.raw, &schema); err != nil {
				t.Fatalf("decode schema: %v", err)
			}
			for _, required := range tt.required {
				if !slices.Contains(schema.Required, required) {
					t.Fatalf("required fields = %v, missing %q", schema.Required, required)
				}
			}
			for field, wantType := range tt.types {
				if got := schema.Properties[field].Type; got != wantType {
					t.Fatalf("%s type = %q, want %q", field, got, wantType)
				}
			}
		})
	}
}

func TestPipelineStepToolsAppearInCompactCatalog(t *testing.T) {
	readOnly := (&Toolset{ReadOnly: true}).GetTools(nil)
	compact := (&Toolset{ReadOnly: true}).GetCompactTools(readOnly)
	compactNames := make(map[string]bool, len(compact))
	for _, tool := range compact {
		compactNames[tool.Tool.Name] = true
	}
	for _, name := range []string{
		"get_pipeline_run_overview",
		"get_pipeline_job_steps",
		"get_pipeline_job_step_log",
		"get_pipeline_job_step_log_url",
	} {
		if !compactNames[name] {
			t.Fatalf("compact tools missing %q", name)
		}
	}
	if compactNames["get_pipeline_run"] {
		t.Fatal("compact tools must keep get_pipeline_run hidden")
	}
}

func pipelineStepToolByName(t *testing.T, tools []toolset.ServerTool, name string) struct {
	domain string
	raw    json.RawMessage
	read   *bool
} {
	t.Helper()
	for _, tool := range tools {
		if tool.Tool.Name == name {
			raw, err := json.Marshal(tool.Tool.InputSchema)
			if err != nil {
				t.Fatalf("encode schema for %q: %v", name, err)
			}
			return struct {
				domain string
				raw    json.RawMessage
				read   *bool
			}{tool.Domain, raw, tool.Tool.Annotations.ReadOnlyHint}
		}
	}
	t.Fatalf("missing tool %q", name)
	return struct {
		domain string
		raw    json.RawMessage
		read   *bool
	}{}
}

func TestPipelineStepToolCountsMatchPublicDocs(t *testing.T) {
	tests := []struct {
		path    string
		matches []string
	}{
		{
			path: "../../../README.md",
			matches: []string{
				"**196 MCP tools** across 9 domains (180 read-only, 16 write)",
				"| **Flow** | 21 |",
			},
		},
		{
			path: "../../../README.zh.md",
			matches: []string{
				"196 个工具中 180 个为只读查询，16 个写操作",
				"| **Flow** | 21 | 19 只读 + 2 可写 |",
			},
		},
		{
			path: "../../../docs/ga-readiness.md",
			matches: []string{
				"180 read-only MCP tools",
				"full catalog contains 196 tools: 180 read-only tools plus 16 write-capable tools",
			},
		},
		{
			path: "../../../docs/flow-tools.md",
			matches: []string{
				"This document describes the 21 MCP tools in the flow domain.",
				"Access summary: 19 read-only, 2 write-capable.",
				"`get_pipeline_job_steps`",
				"`get_pipeline_job_step_log`",
				"`get_pipeline_job_step_log_url`",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			contents, err := os.ReadFile(tt.path)
			if err != nil {
				t.Fatalf("read %s: %v", tt.path, err)
			}
			for _, match := range tt.matches {
				if !strings.Contains(string(contents), match) {
					t.Fatalf("%s missing %q", tt.path, match)
				}
			}
		})
	}
}

func TestHandleGetPipelineJobStepsBuildsEscapedPathAndPreservesNumbers(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		wantURI := "/oapi/v1/flow/organizations/org%2F1/pipelines/pipe%2F1/pipelineRuns/run%2F1/jobs/job%2F1/steps"
		if r.RequestURI != wantURI {
			t.Fatalf("request URI = %q, want %q", r.RequestURI, wantURI)
		}
		_, _ = w.Write([]byte(`{"steps":[{"buildId":9007199254740993}]}`))
	})

	params := pipelineStepParams()
	params["organizationId"] = "org/1"
	params["pipelineId"] = "pipe/1"
	params["pipelineRunId"] = "run/1"
	params["jobId"] = "job/1"
	result, err := handleGetPipelineJobSteps(context.Background(), client, params)
	if err != nil {
		t.Fatalf("handleGetPipelineJobSteps() error = %v", err)
	}
	if !strings.Contains(result, `9007199254740993`) {
		t.Fatalf("result lost buildId precision: %q", result)
	}
}

func TestHandlePipelineJobStepLogBuildsExactQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		wantPath := "/oapi/v1/flow/organizations/org-1/pipelines/pipeline-1/pipelineRuns/run-1/jobs/job-1/step/log"
		if r.URL.Path != wantPath {
			t.Fatalf("path = %q, want %q", r.URL.Path, wantPath)
		}
		wantQuery := url.Values{"stepIndex": {"2"}, "offset": {"0"}, "limit": {"200"}, "buildId": {"9007199254740993"}}
		if got := r.URL.Query(); !slices.Equal(got["stepIndex"], wantQuery["stepIndex"]) ||
			!slices.Equal(got["offset"], wantQuery["offset"]) ||
			!slices.Equal(got["limit"], wantQuery["limit"]) ||
			!slices.Equal(got["buildId"], wantQuery["buildId"]) || len(got) != len(wantQuery) {
			t.Fatalf("query = %q, want %q", got, wantQuery)
		}
		_, _ = w.Write([]byte(`{"content":"ok"}`))
	})

	params := pipelineStepParams()
	params["stepIndex"] = float64(2)
	params["offset"] = float64(0)
	params["limit"] = float64(200)
	params["buildId"] = "9007199254740993"
	if _, err := handleGetPipelineJobStepLog(context.Background(), client, params); err != nil {
		t.Fatalf("handleGetPipelineJobStepLog() error = %v", err)
	}
}

func TestHandlePipelineJobStepLogURLBuildsExactQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		wantPath := "/oapi/v1/flow/organizations/org-1/pipelines/pipeline-1/pipelineRuns/run-1/jobs/job-1/step/log/url"
		if r.URL.Path != wantPath {
			t.Fatalf("path = %q, want %q", r.URL.Path, wantPath)
		}
		wantQuery := url.Values{"stepIndex": {"2"}, "buildId": {"9007199254740993"}}
		if got := r.URL.Query(); !slices.Equal(got["stepIndex"], wantQuery["stepIndex"]) ||
			!slices.Equal(got["buildId"], wantQuery["buildId"]) || len(got) != len(wantQuery) {
			t.Fatalf("query = %q, want %q", got, wantQuery)
		}
		_, _ = w.Write([]byte(`"https://example.test/log"`))
	})

	params := pipelineStepParams()
	params["stepIndex"] = float64(2)
	params["buildId"] = "9007199254740993"
	result, err := handleGetPipelineJobStepLogURL(context.Background(), client, params)
	if err != nil {
		t.Fatalf("handleGetPipelineJobStepLogURL() error = %v", err)
	}
	if result != `"https://example.test/log"` {
		t.Fatalf("result = %q, want JSON string URL", result)
	}
}

func TestPipelineStepHandlersWrapUpstreamErrors(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"upstream failure"}`))
	})

	tests := []struct {
		name    string
		handler func(context.Context, any, map[string]any) (string, error)
		params  map[string]any
	}{
		{"steps", handleGetPipelineJobSteps, pipelineStepParams()},
		{"log", handleGetPipelineJobStepLog, pipelineStepLogParams()},
		{"log URL", handleGetPipelineJobStepLogURL, pipelineStepLogURLParams()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.handler(context.Background(), client, tt.params); err == nil || !strings.Contains(err.Error(), "Suggestion:") {
				t.Fatalf("expected wrapped upstream error, got %v", err)
			}
		})
	}
}

func pipelineStepParams() map[string]any {
	return map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipeline-1",
		"pipelineRunId":  "run-1",
		"jobId":          "job-1",
	}
}

func pipelineStepLogParams() map[string]any {
	params := pipelineStepParams()
	params["stepIndex"] = float64(2)
	params["offset"] = float64(0)
	params["limit"] = float64(200)
	params["buildId"] = "9007199254740993"
	return params
}

func pipelineStepLogURLParams() map[string]any {
	params := pipelineStepParams()
	params["stepIndex"] = float64(2)
	params["buildId"] = "9007199254740993"
	return params
}
