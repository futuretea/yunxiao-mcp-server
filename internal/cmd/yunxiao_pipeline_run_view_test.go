package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIPipelineRunViewPrintsJSONWithDefaultOrganization(t *testing.T) {
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs/run-1":
			_, _ = w.Write([]byte(`{"pipelineRunId":"run-1","status":"SUCCESS"}`))
		case "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/listTasksByCategory/TEST":
			_, _ = w.Write([]byte(`[{"id":"job-1"}]`))
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
		"pipeline", "run", "view", "run-1",
		"--pipeline-id", "pipe-1",
		"--category", "TEST",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out.String())
	}
	for _, section := range []string{"run", "jobs", "filters"} {
		if _, ok := payload[section]; !ok {
			t.Fatalf("payload missing %s: %#v", section, payload)
		}
	}
	if requests["/oapi/v1/platform/organizations"] != 1 ||
		requests["/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs/run-1"] != 1 ||
		requests["/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/listTasksByCategory/TEST"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIPipelineRunOverviewAliasUsesExplicitOrganizationAndSkipsJobs(t *testing.T) {
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/flow/organizations/org-2/pipelines/pipe-1/runs/run-1":
			_, _ = w.Write([]byte(`{"pipelineRunId":"run-1"}`))
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
		"pipelines", "runs", "overview", "run-1",
		"--organization-id", "org-2",
		"--pipeline-id", "pipe-1",
		"--include-jobs=false",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out.String())
	}
	if _, ok := payload["jobs"]; ok {
		t.Fatalf("payload should omit jobs: %#v", payload)
	}
	if requests["/oapi/v1/flow/organizations/org-2/pipelines/pipe-1/listTasksByCategory/DEPLOY"] != 0 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIPipelineRunViewRequiresIDs(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"pipeline", "run", "view", "run-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected pipeline-id error")
	}
	if !strings.Contains(err.Error(), "pipeline-id is required") {
		t.Fatalf("error = %v", err)
	}

	command = NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"pipeline", "run", "view", "--pipeline-id", "pipe-1"})
	err = command.Execute()
	if err == nil {
		t.Fatal("Execute() expected missing argument error")
	}
}

func TestYunxiaoCLIPipelineRunViewRejectsBlankRunIDBeforeNetwork(t *testing.T) {
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
		"pipeline", "run", "view", " ",
		"--pipeline-id", "pipe-1",
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected pipeline-run-id error")
	}
	if !strings.Contains(err.Error(), "pipeline-run-id is required") {
		t.Fatalf("error = %v", err)
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0", requests)
	}
}

func TestYunxiaoCLIPipelineRunViewReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "flow", "pipeline", "run", "view", "run-1", "--pipeline-id", "pipe-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "get_pipeline_run_overview"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestPipelineRunViewOptionsParamsIncludesOverviewFilters(t *testing.T) {
	params, err := (pipelineRunViewOptions{
		OrganizationID: " org-1 ",
		PipelineID:     " pipe-1 ",
		IncludeJobs:    false,
		Category:       " TEST ",
	}).params(" run-1 ")
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
		"pipelineRunId":  "run-1",
		"includeJobs":    false,
		"category":       "TEST",
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPipelineRunViewOptionsParamsRequiresIDs(t *testing.T) {
	if _, err := (pipelineRunViewOptions{}).params("run-1"); err == nil {
		t.Fatal("params() expected pipeline-id error")
	}
	if _, err := (pipelineRunViewOptions{PipelineID: "pipe-1"}).params(" "); err == nil {
		t.Fatal("params() expected pipeline-run-id error")
	}
}
