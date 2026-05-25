package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIProjectRoleListPrintsTableWithDefaultOrganization(t *testing.T) {
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/roles":
			_, _ = w.Write([]byte(`[{"id":"project.admin","name":"Admin","description":"Project administrator"}]`))
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
		"project", "role", "list",
		"--project-id", "project-1",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"ID", "NAME", "DESCRIPTION", "project.admin", "Admin", "Project administrator"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/projex/organizations/org-1/projects/project-1/roles"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIProjectsRolesAliasListPrintsJSONWithExplicitOrganization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/projex/organizations/org-2/projects/project-1/roles":
			_, _ = w.Write([]byte(`[{"id":"project.admin","name":"Admin"}]`))
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
		"projects", "roles", "list",
		"--organization-id", "org-2",
		"--project-id", "project-1",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"id": "project.admin"`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIProjectRoleListRequiresProjectID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"project", "role", "list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected project-id error")
	}
	if !strings.Contains(err.Error(), "project-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIProjectRoleListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "projex", "project", "role", "list", "--project-id", "project-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_project_roles"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestProjectRoleListOptionsParamsIncludesOrganization(t *testing.T) {
	params, err := (projectRoleListOptions{
		OrganizationID: " org-1 ",
		ProjectID:      " project-1 ",
	}).params()
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPrintProjectRoleListShowsNoResultsWhenRowsEmpty(t *testing.T) {
	var out bytes.Buffer
	raw := "No results found."
	if err := printProjectRoleList(&out, raw); err != nil {
		t.Fatalf("printProjectRoleList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want \"No results found.\"", out.String())
	}
}

func TestProjectRoleRowsFromJSONExtractsNestedValues(t *testing.T) {
	rows := projectRoleRowsFromJSON(`{"result":{"items":[{"projectRoleId":"project.maintainer","projectRoleName":"Maintainer","desc":"Can maintain"}]}}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := projectRoleRow{ID: "project.maintainer", Name: "Maintainer", Description: "Can maintain"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestProjectRoleRowsFromJSONSkipsNonObjectRows(t *testing.T) {
	rows := projectRoleRowsFromJSON(`{"data":["skip",{"id":"project.admin","name":"Admin"}]}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	if rows[0].ID != "project.admin" {
		t.Fatalf("row = %#v", rows[0])
	}
}

func TestProjectRoleRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows := projectRoleRowsFromJSON(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}
