package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLISprintListPrintsTableWithDefaultOrganization(t *testing.T) {
	var gotQuery string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/sprints":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[{"id":"sprint-1","name":"Release 1","status":{"name":"DOING"}}]`))
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
		"sprint", "list",
		"--project-id", "project-1",
		"--status", "TODO,DOING",
		"--name", "Release",
		"--page", "2",
		"--per-page", "20",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"ID", "NAME", "STATUS", "sprint-1", "Release 1", "DOING"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	for _, want := range []string{"status=TODO%2CDOING", "name=Release", "page=2", "perPage=20"} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %q", gotQuery, want)
		}
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/projex/organizations/org-1/projects/project-1/sprints"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLISprintsAliasListPrintsJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/sprints":
			_, _ = w.Write([]byte(`[{"id":"sprint-1","name":"Release 1"}]`))
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
		"sprints", "list",
		"--project-id", "project-1",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"id": "sprint-1"`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLISprintListSkipsDefaultOrgWhenOrganizationProvided(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/projex/organizations/org-2/projects/project-1/sprints":
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`[]`))
		case "/oapi/v1/platform/organizations":
			t.Fatal("should not resolve default organization when organizationId is provided")
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"sprint", "list",
		"--organization-id", "org-2",
		"--project-id", "project-1",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if gotPath != "/oapi/v1/projex/organizations/org-2/projects/project-1/sprints" {
		t.Fatalf("path = %q", gotPath)
	}
}

func TestYunxiaoCLISprintListRequiresProjectID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--access-token", "token-1",
		"sprint", "list",
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected project-id error")
	}
	if !strings.Contains(err.Error(), "project-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestSprintListOptionsParamsIncludesFiltersAndPagination(t *testing.T) {
	params, err := (sprintListOptions{
		OrganizationID: " org-1 ",
		ProjectID:      " project-1 ",
		Status:         " TODO,DOING ",
		Name:           " Release ",
		Page:           2,
		PerPage:        20,
	}).params()
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"status":         "TODO,DOING",
		"name":           "Release",
		"page":           2,
		"perPage":        20,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPrintSprintListFallsBackToRawJSON(t *testing.T) {
	var out bytes.Buffer
	raw := `{"data":{"total":0}}`
	if err := printSprintList(&out, raw); err != nil {
		t.Fatalf("printSprintList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want raw JSON", out.String())
	}
}

func TestSprintRowsFromJSONExtractsNestedValues(t *testing.T) {
	rows := sprintRowsFromJSON(`{"data":[{"sprintId":123,"title":"Release 1","status":{"displayName":"Open"}},"skip"]}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := sprintRow{ID: "123", Name: "Release 1", Status: "Open"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestSprintRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows := sprintRowsFromJSON(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}
