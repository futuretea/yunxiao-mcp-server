package cmd

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIProjectListPrintsTable(t *testing.T) {
	var gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects:search":
			body, _ := io.ReadAll(r.Body)
			gotBody = string(body)
			_, _ = w.Write([]byte(`{"data":[{"id":"project-1","name":"Alpha","status":{"name":"Doing"},"creator":{"displayName":"Ada"}}]}`))
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
		"project", "list",
		"--name", "Alpha",
		"--page", "2",
		"--per-page", "50",
		"--order-by", "name",
		"--sort", "asc",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"ID", "NAME", "STATUS", "CREATOR", "project-1", "Alpha", "Doing", "Ada"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	for _, want := range []string{`"page":2`, `"perPage":50`, `"orderBy":"name"`, `"sort":"asc"`, `"conditions"`} {
		if !strings.Contains(gotBody, want) {
			t.Fatalf("body = %q, missing %q", gotBody, want)
		}
	}
}

func TestYunxiaoCLIProjectsAliasListPrintsJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects:search":
			_, _ = w.Write([]byte(`{"data":[{"id":"project-1","name":"Alpha"}]}`))
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
		"projects", "list",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"id": "project-1"`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIProjectListSkipsDefaultOrgWhenOrganizationProvided(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/projex/organizations/org-2/projects:search":
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"data":[]}`))
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
		"project", "list",
		"--organization-id", "org-2",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if gotPath != "/oapi/v1/projex/organizations/org-2/projects:search" {
		t.Fatalf("path = %q", gotPath)
	}
}

func TestProjectListOptionsParamsIncludesFiltersAndPagination(t *testing.T) {
	params := (projectListOptions{
		OrganizationID: " org-1 ",
		Name:           " Alpha ",
		Status:         " open,closed ",
		Creator:        " user-1 ",
		OrderBy:        " name ",
		Sort:           " desc ",
		Page:           2,
		PerPage:        50,
	}).params()

	wants := map[string]any{
		"organizationId": "org-1",
		"name":           "Alpha",
		"status":         "open,closed",
		"creator":        "user-1",
		"orderBy":        "name",
		"sort":           "desc",
		"page":           2,
		"perPage":        50,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPrintProjectListShowsNoResultsWhenRowsEmpty(t *testing.T) {
	var out bytes.Buffer
	raw := "No results found."
	if err := printProjectList(&out, raw); err != nil {
		t.Fatalf("printProjectList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want \"No results found.\"", out.String())
	}
}

func TestProjectRowsFromJSONExtractsNestedValues(t *testing.T) {
	rows := projectRowsFromJSON(`{"result":{"items":[{"projectId":123,"title":"Alpha","status":{"displayName":"Open"},"owner":{"name":"Ada"}}]}}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := projectRow{ID: "123", Name: "Alpha", Status: "Open", Creator: "Ada"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestProjectRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows := projectRowsFromJSON(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}

func TestYunxiaoCLIProjectViewReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "projex", "project", "view", "123"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "get_project_overview"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestProjectViewOptionsParamsIncludesFlags(t *testing.T) {
	params := (projectViewOptions{
		OrganizationID:    " org-1 ",
		ProjectID:         "123",
		IncludeMembers:    true,
		IncludeSprints:    false,
		ActiveOnly:        true,
		Status:            " DOING ",
		Page:              2,
		PerPage:           10,
	}).params()

	wants := map[string]any{
		"organizationId":  "org-1",
		"projectId":       "123",
		"includeMembers":  true,
		"includeSprints":  false,
		"activeOnly":      true,
		"status":          "DOING",
		"page":            2,
		"perPage":         10,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}
