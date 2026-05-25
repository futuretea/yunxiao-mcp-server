package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIPipelineRunListPrintsTableWithDefaultOrganization(t *testing.T) {
	var gotQuery string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"data":[{"pipelineRunId":"run-1","status":"SUCCESS","result":"PASS","startTime":1704067200000,"endTime":1704067800000,"triggerMode":1}],"pagination":{"total":1}}`))
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
		"pipeline", "run", "list",
		"--pipeline-id", "pipe-1",
		"--page", "2",
		"--per-page", "30",
		"--start-time", "1704067200000",
		"--end-time", "1704153600000",
		"--status", "SUCCESS",
		"--trigger-mode", "1",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"ID", "STATUS", "RESULT", "START", "END", "TRIGGER", "run-1", "SUCCESS", "PASS", "1704067200000", "1"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	for _, want := range []string{
		"page=2",
		"perPage=30",
		"startTime=1704067200000",
		"endTme=1704153600000",
		"status=SUCCESS",
		"triggerMode=1",
	} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %q", gotQuery, want)
		}
	}
	if strings.Contains(gotQuery, "endTime=") {
		t.Fatalf("query = %q, should use Yunxiao endTme spelling", gotQuery)
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIPipelineRunsAliasListPrintsJSONWithExplicitOrganization(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/flow/organizations/org-2/pipelines/pipe-1/runs":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[{"pipelineRunId":"run-1"}]`))
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
		"pipelines", "runs", "list",
		"--organization-id", "org-2",
		"--pipeline-id", "pipe-1",
		"--status", "RUNNING",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"pipelineRunId": "run-1"`) {
		t.Fatalf("stdout = %q", out.String())
	}
	if gotQuery != "status=RUNNING" {
		t.Fatalf("query = %q, want status=RUNNING", gotQuery)
	}
}

func TestYunxiaoCLIPipelineRunListRequiresPipelineID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"pipeline", "run", "list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected pipeline-id error")
	}
	if !strings.Contains(err.Error(), "pipeline-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIPipelineRunListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "flow", "pipeline", "run", "list", "--pipeline-id", "pipe-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_pipeline_runs"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestPipelineRunListOptionsParamsIncludesFilters(t *testing.T) {
	params, err := (pipelineRunListOptions{
		OrganizationID: " org-1 ",
		PipelineID:     " pipe-1 ",
		Page:           2,
		PerPage:        30,
		StartTime:      1704067200000,
		EndTime:        1704153600000,
		Status:         " SUCCESS ",
		TriggerMode:    1,
	}).params()
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
		"page":           2,
		"perPage":        30,
		"startTime":      int64(1704067200000),
		"endTime":        int64(1704153600000),
		"status":         "SUCCESS",
		"triggerMode":    1,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPipelineRunListOptionsParamsRequiresPipelineID(t *testing.T) {
	if _, err := (pipelineRunListOptions{}).params(); err == nil {
		t.Fatal("params() expected pipeline-id error")
	}
}

func TestPrintPipelineRunListShowsNoResultsWhenRowsEmpty(t *testing.T) {
	var out bytes.Buffer
	raw := "No results found."
	if err := printPipelineRunList(&out, raw); err != nil {
		t.Fatalf("printPipelineRunList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want \"No results found.\"", out.String())
	}
}

func TestPrintPipelineRunListPrintsHeaderForEmptyList(t *testing.T) {
	var out bytes.Buffer
	if err := printPipelineRunList(&out, `{"data":[]}`); err != nil {
		t.Fatalf("printPipelineRunList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != "ID  STATUS  RESULT  START  END  TRIGGER" {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestPipelineRunRowsFromJSONExtractsAlternateFields(t *testing.T) {
	rows := pipelineRunRowsFromJSON(`{"result":{"items":[{"runId":"run-1","runStatus":"RUNNING","runResult":"UNKNOWN","startedAt":"2026-05-25T10:00:00Z","finishedAt":"2026-05-25T10:10:00Z","triggerUser":"alice"}]}}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := pipelineRunRow{
		ID:      "run-1",
		Status:  "RUNNING",
		Result:  "UNKNOWN",
		Start:   "2026-05-25T10:00:00Z",
		End:     "2026-05-25T10:10:00Z",
		Trigger: "alice",
	}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestPipelineRunRowsFromJSONSkipsNonObjectRows(t *testing.T) {
	rows := pipelineRunRowsFromJSON(`{"data":["skip",{"id":"run-1","state":"SUCCESS","executeResult":"PASS","gmtStarted":"start","gmtFinished":"end","creator":"bob"}]}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := pipelineRunRow{ID: "run-1", Status: "SUCCESS", Result: "PASS", Start: "start", End: "end", Trigger: "bob"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestPipelineRunRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows := pipelineRunRowsFromJSON(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}
