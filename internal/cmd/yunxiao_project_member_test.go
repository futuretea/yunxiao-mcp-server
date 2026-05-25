package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIProjectMemberListPrintsTableWithDefaultOrganization(t *testing.T) {
	var gotQuery string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/members":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[{"userId":"user-1","userName":"Alice","roleId":"project.admin","roleName":"Admin"}]`))
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
		"project", "member", "list",
		"--project-id", "project-1",
		"--name", "Alice",
		"--role-id", "project.admin",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"USER_ID", "NAME", "ROLE_ID", "ROLE", "user-1", "Alice", "project.admin", "Admin"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	for _, want := range []string{"name=Alice", "roleId=project.admin"} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %q", gotQuery, want)
		}
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/projex/organizations/org-1/projects/project-1/members"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIProjectsMembersAliasListPrintsJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/members":
			_, _ = w.Write([]byte(`[{"userId":"user-1","displayName":"Alice"}]`))
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
		"projects", "members", "list",
		"--project-id", "project-1",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"userId": "user-1"`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIProjectMemberListSkipsDefaultOrgWhenOrganizationProvided(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/projex/organizations/org-2/projects/project-1/members":
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
		"project", "member", "list",
		"--organization-id", "org-2",
		"--project-id", "project-1",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if gotPath != "/oapi/v1/projex/organizations/org-2/projects/project-1/members" {
		t.Fatalf("path = %q", gotPath)
	}
}

func TestYunxiaoCLIProjectMemberListRequiresProjectID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"project", "member", "list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected project-id error")
	}
	if !strings.Contains(err.Error(), "project-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIProjectMemberListReturnsConfigError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--base-url", "://invalid-url", "project", "member", "list", "--project-id", "project-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected config error")
	}
	if !strings.Contains(err.Error(), "load config") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIProjectMemberListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "projex", "project", "member", "list", "--project-id", "project-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_project_members"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestProjectMemberListOptionsParamsIncludesFilters(t *testing.T) {
	params, err := (projectMemberListOptions{
		OrganizationID: " org-1 ",
		ProjectID:      " project-1 ",
		Name:           " Alice ",
		RoleID:         " project.admin ",
	}).params()
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"name":           "Alice",
		"roleId":         "project.admin",
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPrintProjectMemberListShowsNoResultsWhenRowsEmpty(t *testing.T) {
	var out bytes.Buffer
	raw := "No results found."
	if err := printProjectMemberList(&out, raw); err != nil {
		t.Fatalf("printProjectMemberList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want \"No results found.\"", out.String())
	}
}

func TestProjectMemberRowsFromJSONExtractsNestedValues(t *testing.T) {
	rows := projectMemberRowsFromJSON(`{"result":{"items":[{"accountId":123,"realName":"Alice","projectRoleId":"project.maintainer","projectRole":{"displayName":"Maintainer"}}]}}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := projectMemberRow{UserID: "123", Name: "Alice", RoleID: "project.maintainer", RoleName: "Maintainer"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestProjectMemberRowsFromJSONSkipsNonObjectRows(t *testing.T) {
	rows := projectMemberRowsFromJSON(`{"data":["skip",{"id":"member-1","userId":"user-1"}]}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	if rows[0].UserID != "user-1" {
		t.Fatalf("row = %#v", rows[0])
	}
}

func TestProjectMemberRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows := projectMemberRowsFromJSON(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}
