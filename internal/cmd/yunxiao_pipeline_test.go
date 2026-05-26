package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIPipelineListPrintsTableWithDefaultOrganization(t *testing.T) {
	var gotQuery string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/flow/organizations/org-1/pipelines":
			gotQuery = r.URL.RawQuery
			w.Header().Set("x-total", "1")
			_, _ = w.Write([]byte(`{"data":[{"pipelineId":123,"pipelineName":"Deploy","status":"SUCCESS","latestRun":{"pipelineRunId":"run-1"}}],"pagination":{"total":1}}`))
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
		"pipeline", "list",
		"--name", "deploy",
		"--status", "RUNNING,SUCCESS",
		"--create-start-time", "1704067200000",
		"--create-end-time", "1704153600000",
		"--execute-start-time", "1704240000000",
		"--execute-end-time", "1704326400000",
		"--page", "2",
		"--per-page", "30",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"ID", "NAME", "STATUS", "LAST_RUN", "123", "Deploy", "SUCCESS", "run-1"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	for _, want := range []string{
		"pipelineName=deploy",
		"statusList=RUNNING%2CSUCCESS",
		"createStartTime=1704067200000",
		"createEndTime=1704153600000",
		"executeStartTime=1704240000000",
		"executeEndTime=1704326400000",
		"page=2",
		"perPage=30",
	} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %q", gotQuery, want)
		}
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/flow/organizations/org-1/pipelines"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIPipelinesAliasListPrintsJSONWithExplicitOrganization(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/flow/organizations/org-2/pipelines":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[{"pipelineId":"pipe-1","pipelineName":"Deploy"}]`))
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
		"pipelines", "list",
		"--organization-id", "org-2",
		"--status", "SUCCESS",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"pipelineId": "pipe-1"`) {
		t.Fatalf("stdout = %q", out.String())
	}
	if gotQuery != "statusList=SUCCESS" {
		t.Fatalf("query = %q, want statusList=SUCCESS", gotQuery)
	}
}

func TestYunxiaoCLIPipelineListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "flow", "pipeline", "list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_pipelines"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestPipelineListOptionsParamsIncludesFilters(t *testing.T) {
	params := (pipelineListOptions{
		OrganizationID:   " org-1 ",
		PipelineName:     " deploy ",
		StatusList:       " RUNNING,SUCCESS ",
		CreateStartTime:  1704067200000,
		CreateEndTime:    1704153600000,
		ExecuteStartTime: 1704240000000,
		ExecuteEndTime:   1704326400000,
		Page:             2,
		PerPage:          30,
	}).params()

	wants := map[string]any{
		"organizationId":   "org-1",
		"pipelineName":     "deploy",
		"statusList":       "RUNNING,SUCCESS",
		"createStartTime":  int64(1704067200000),
		"createEndTime":    int64(1704153600000),
		"executeStartTime": int64(1704240000000),
		"executeEndTime":   int64(1704326400000),
		"page":             2,
		"perPage":          30,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPrintPipelineListShowsNoResultsWhenRowsEmpty(t *testing.T) {
	var out bytes.Buffer
	raw := "No results found."
	if err := printPipelineList(&out, raw); err != nil {
		t.Fatalf("printPipelineList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want \"No results found.\"", out.String())
	}
}

func TestPrintPipelineListPrintsHeaderForEmptyList(t *testing.T) {
	var out bytes.Buffer
	if err := printPipelineList(&out, `{"data":[]}`); err != nil {
		t.Fatalf("printPipelineList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != "[1mID  NAME  STATUS  LAST_RUN[0m" {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestPipelineRowsFromJSONExtractsAlternateFields(t *testing.T) {
	rows, _ := pipelineRowsFromJSONForPrint(`{"result":{"items":[{"id":"pipe-1","name":"Deploy","lastRunStatus":"FAIL","pipelineRun":{"runId":"run-2"}}]}}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := pipelineRow{ID: "pipe-1", Name: "Deploy", Status: "FAIL", LastRun: "run-2"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestPipelineRowsFromJSONSkipsNonObjectRows(t *testing.T) {
	rows, _ := pipelineRowsFromJSONForPrint(`{"data":["skip",{"pipelineID":"pipe-1","displayName":"Deploy","latestRunId":"run-3"}]}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := pipelineRow{ID: "pipe-1", Name: "Deploy", LastRun: "run-3"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestPipelineRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows, _ := pipelineRowsFromJSONForPrint(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}
