package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIOrganizationListPrintsTable(t *testing.T) {
	var gotQuery string
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/platform/organizations" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		requests++
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`[{"id":"org-1","name":"FutureTea"}]`))
	}))
	defer server.Close()

	var out bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"organization", "list",
		"--page", "2",
		"--per-page", "50",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"ID", "NAME", "org-1", "FutureTea"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	for _, want := range []string{"page=2", "perPage=50"} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %q", gotQuery, want)
		}
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}
}

func TestYunxiaoCLIOrganizationsAliasListPrintsJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/platform/organizations" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"id":"org-1","name":"FutureTea"}]`))
	}))
	defer server.Close()

	var out bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"organizations", "list",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"id": "org-1"`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIOrganizationListReturnsConfigError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--config", t.TempDir() + "/missing.yaml",
		"organization", "list",
	})

	if err := command.Execute(); err == nil {
		t.Fatal("Execute() expected config error")
	}
}

func TestYunxiaoCLIOrganizationListReturnsToolError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/platform/organizations" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer server.Close()

	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"organization", "list",
	})

	if err := command.Execute(); err == nil {
		t.Fatal("Execute() expected tool error")
	}
}

func TestOrganizationListOptionsParamsIncludesPagination(t *testing.T) {
	params := (organizationListOptions{Page: 2, PerPage: 50}).params()
	wants := map[string]any{
		"page":    2,
		"perPage": 50,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}

	if params := (organizationListOptions{}).params(); len(params) != 0 {
		t.Fatalf("params = %#v, want empty map", params)
	}
}

func TestPrintOrganizationListShowsNoResultsWhenRowsEmpty(t *testing.T) {
	var out bytes.Buffer
	raw := "No results found."
	if err := printOrganizationList(&out, raw); err != nil {
		t.Fatalf("printOrganizationList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want \"No results found.\"", out.String())
	}
}

func TestOrganizationRowsFromJSONExtractsNestedValues(t *testing.T) {
	rows := organizationRowsFromJSON(`{"data":[{"organizationId":123,"displayName":"FutureTea"},"skip"]}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := organizationRow{ID: "123", Name: "FutureTea"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestOrganizationRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows := organizationRowsFromJSON(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}

func TestYunxiaoCLIOrganizationViewReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "platform", "organization", "view"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "get_organization_overview"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestOrgViewOptionsParamsIncludesFlags(t *testing.T) {
	params := (orgViewOptions{
		OrganizationID:     " org-1 ",
		IncludeDepartments: true,
		IncludeMembers:     false,
		IncludeGroups:      true,
		DepartmentLimit:    3,
	}).params()

	wants := map[string]any{
		"organizationId":     "org-1",
		"includeDepartments": true,
		"includeMembers":     false,
		"includeGroups":      true,
		"includeRoles":       false,
		"departmentLimit":    3,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}
