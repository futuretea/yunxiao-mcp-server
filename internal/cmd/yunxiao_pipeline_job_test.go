package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIPipelineJobListPrintsTableWithDefaultOrganization(t *testing.T) {
	var gotPath string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/listTasksByCategory/DEPLOY":
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`[{"identifier":"deploy-prod","name":"Deploy Production","category":"DEPLOY","status":"SUCCESS"}]`))
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
		"pipeline", "job", "list",
		"--pipeline-id", "pipe-1",
		"--category", "DEPLOY",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"IDENTIFIER", "NAME", "CATEGORY", "STATUS", "deploy-prod", "Deploy Production", "DEPLOY", "SUCCESS"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	if gotPath != "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/listTasksByCategory/DEPLOY" {
		t.Fatalf("path = %q", gotPath)
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/listTasksByCategory/DEPLOY"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIPipelineJobsAliasListPrintsJSONWithExplicitOrganization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/flow/organizations/org-2/pipelines/pipe-1/listTasksByCategory/DEPLOY":
			_, _ = w.Write([]byte(`[{"identifier":"job-1","name":"Job 1"}]`))
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
		"pipelines", "jobs", "list",
		"--organization-id", "org-2",
		"--pipeline-id", "pipe-1",
		"--category", "DEPLOY",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"identifier": "job-1"`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIPipelineJobListRequiresPipelineID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"pipeline", "job", "list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected pipeline-id error")
	}
	if !strings.Contains(err.Error(), "pipeline-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIPipelineJobListRequiresCategory(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"pipeline", "job", "list", "--pipeline-id", "pipe-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected category error")
	}
	if !strings.Contains(err.Error(), "category is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIPipelineJobListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "flow", "pipeline", "job", "list", "--pipeline-id", "pipe-1", "--category", "DEPLOY"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_pipeline_jobs_by_category"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestPipelineJobListOptionsParamsIncludesFilters(t *testing.T) {
	params, err := (pipelineJobListOptions{
		OrganizationID: " org-1 ",
		PipelineID:     " pipe-1 ",
		Category:       " DEPLOY ",
	}).params()
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
		"category":       "DEPLOY",
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPipelineJobListOptionsParamsRequiresPipelineID(t *testing.T) {
	if _, err := (pipelineJobListOptions{Category: "DEPLOY"}).params(); err == nil {
		t.Fatal("params() expected pipeline-id error")
	}
}

func TestPipelineJobListOptionsParamsRequiresCategory(t *testing.T) {
	if _, err := (pipelineJobListOptions{PipelineID: "pipe-1"}).params(); err == nil {
		t.Fatal("params() expected category error")
	}
}

func TestPrintPipelineJobListShowsNoResultsWhenRowsEmpty(t *testing.T) {
	var out bytes.Buffer
	raw := "No results found."
	if err := printPipelineJobList(&out, raw); err != nil {
		t.Fatalf("printPipelineJobList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want \"No results found.\"", out.String())
	}
}

func TestPrintPipelineJobListPrintsHeaderForEmptyList(t *testing.T) {
	var out bytes.Buffer
	if err := printPipelineJobList(&out, `[]`); err != nil {
		t.Fatalf("printPipelineJobList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != "[1mIDENTIFIER  NAME  CATEGORY  STATUS[0m" {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestPipelineJobRowsFromJSONExtractsAlternateFields(t *testing.T) {
	rows, _ := pipelineJobRowsFromJSONForPrint(`{"result":{"items":[{"id":"job-1","displayName":"Job 1","taskCategory":"BUILD","taskStatus":"RUNNING"}]}}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := pipelineJobRow{
		Identifier: "job-1",
		Name:       "Job 1",
		Category:   "BUILD",
		Status:     "RUNNING",
	}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestPipelineJobRowsFromJSONSkipsNonObjectRows(t *testing.T) {
	rows, _ := pipelineJobRowsFromJSONForPrint(`{"data":["skip",{"jobId":"j-1","taskName":"Deploy","category":"DEPLOY","state":"SUCCESS"}]}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := pipelineJobRow{Identifier: "j-1", Name: "Deploy", Category: "DEPLOY", Status: "SUCCESS"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestPipelineJobRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows, _ := pipelineJobRowsFromJSONForPrint(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}

func TestYunxiaoCLIPipelineJobLogPrintsOutputWithDefaultOrganization(t *testing.T) {
	var gotPath string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs/run-1/job/job-1/log":
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"content":"build log output"}`))
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
		"pipeline", "job", "log",
		"--pipeline-id", "pipe-1",
		"--run-id", "run-1",
		"--job-id", "job-1",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `build log output`) {
		t.Fatalf("stdout = %q", out.String())
	}
	if gotPath != "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs/run-1/job/job-1/log" {
		t.Fatalf("path = %q", gotPath)
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs/run-1/job/job-1/log"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIPipelineJobLogPrintsOutputWithExplicitOrganization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/flow/organizations/org-2/pipelines/pipe-1/runs/run-1/job/job-1/log":
			_, _ = w.Write([]byte(`{"content":"log text"}`))
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
		"pipelines", "jobs", "log",
		"--organization-id", "org-2",
		"--pipeline-id", "pipe-1",
		"--run-id", "run-1",
		"--job-id", "job-1",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `log text`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIPipelineJobLogRequiresPipelineID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"pipeline", "job", "log"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected pipeline-id error")
	}
	if !strings.Contains(err.Error(), "pipeline-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIPipelineJobLogRequiresRunID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"pipeline", "job", "log", "--pipeline-id", "pipe-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected run-id error")
	}
	if !strings.Contains(err.Error(), "run-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIPipelineJobLogRequiresJobID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"pipeline", "job", "log", "--pipeline-id", "pipe-1", "--run-id", "run-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected job-id error")
	}
	if !strings.Contains(err.Error(), "job-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIPipelineJobLogReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "flow", "pipeline", "job", "log", "--pipeline-id", "pipe-1", "--run-id", "run-1", "--job-id", "job-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "get_pipeline_job_run_log"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestPipelineJobLogOptionsParamsIncludesFilters(t *testing.T) {
	params, err := (pipelineJobLogOptions{
		OrganizationID: " org-1 ",
		PipelineID:     " pipe-1 ",
		PipelineRunID:  " run-1 ",
		JobID:          " job-1 ",
	}).params()
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
		"pipelineRunId":  "run-1",
		"jobId":          "job-1",
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPipelineJobLogOptionsParamsRequiresPipelineID(t *testing.T) {
	if _, err := (pipelineJobLogOptions{PipelineRunID: "run-1", JobID: "job-1"}).params(); err == nil {
		t.Fatal("params() expected pipeline-id error")
	}
}

func TestPipelineJobLogOptionsParamsRequiresRunID(t *testing.T) {
	if _, err := (pipelineJobLogOptions{PipelineID: "pipe-1", JobID: "job-1"}).params(); err == nil {
		t.Fatal("params() expected run-id error")
	}
}

func TestPipelineJobLogOptionsParamsRequiresJobID(t *testing.T) {
	if _, err := (pipelineJobLogOptions{PipelineID: "pipe-1", PipelineRunID: "run-1"}).params(); err == nil {
		t.Fatal("params() expected job-id error")
	}
}
