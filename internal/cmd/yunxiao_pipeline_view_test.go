package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIPipelineViewPrintsJSONWithDefaultOrganization(t *testing.T) {
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1":
			_, _ = w.Write([]byte(`{"pipelineId":"pipe-1","name":"Deploy"}`))
		case "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs/latestPipelineRun":
			_, _ = w.Write([]byte(`{"pipelineRunId":"run-latest"}`))
		case "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs":
			query := r.URL.Query()
			if query.Get("page") != "1" || query.Get("perPage") != "3" {
				t.Fatalf("runs query = %q", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`[{"pipelineRunId":"run-1"}]`))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	var out bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"pipeline", "view", "pipe-1",
		"--run-limit", "3",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out.String())
	}
	for _, section := range []string{"pipeline", "latestRun", "runs", "filters"} {
		if _, ok := payload[section]; !ok {
			t.Fatalf("payload missing %s: %#v", section, payload)
		}
	}
	if requests["/oapi/v1/platform/organizations"] != 1 ||
		requests["/oapi/v1/flow/organizations/org-1/pipelines/pipe-1"] != 1 ||
		requests["/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs/latestPipelineRun"] != 1 ||
		requests["/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIPipelineOverviewAliasUsesExplicitOrganizationAndSkipsRuns(t *testing.T) {
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/flow/organizations/org-2/pipelines/pipe-1":
			_, _ = w.Write([]byte(`{"pipelineId":"pipe-1","name":"Deploy"}`))
		case "/oapi/v1/flow/organizations/org-2/pipelines/pipe-1/runs/latestPipelineRun":
			_, _ = w.Write([]byte(`{"pipelineRunId":"run-latest"}`))
		case "/oapi/v1/platform/organizations":
			t.Fatal("should not resolve default organization when organizationId is provided")
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	var out bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"pipelines", "overview", "pipe-1",
		"--organization-id", "org-2",
		"--include-runs=false",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out.String())
	}
	if _, ok := payload["runs"]; ok {
		t.Fatalf("payload should omit runs: %#v", payload)
	}
	if requests["/oapi/v1/flow/organizations/org-2/pipelines/pipe-1/runs"] != 0 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIPipelineViewRequiresID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"pipeline", "view"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected missing argument error")
	}
}

func TestYunxiaoCLIPipelineViewRejectsBlankIDBeforeNetwork(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		http.Error(w, "unexpected request", http.StatusInternalServerError)
	}))
	defer server.Close()

	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"pipeline", "view", " ",
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected pipeline-id error")
	}
	if !strings.Contains(err.Error(), "pipeline-id is required") {
		t.Fatalf("error = %v", err)
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0", requests)
	}
}

func TestYunxiaoCLIPipelineViewReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "flow", "pipeline", "view", "pipe-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "get_pipeline_overview"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestPipelineViewOptionsParamsIncludesOverviewFilters(t *testing.T) {
	params, err := (pipelineViewOptions{
		OrganizationID: " org-1 ",
		IncludeRuns:    false,
		RunLimit:       3,
	}).params(" pipe-1 ")
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
		"includeRuns":    false,
		"runLimit":       3,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPipelineViewOptionsParamsRequiresID(t *testing.T) {
	if _, err := (pipelineViewOptions{}).params(" "); err == nil {
		t.Fatal("params() expected pipeline-id error")
	}
}
