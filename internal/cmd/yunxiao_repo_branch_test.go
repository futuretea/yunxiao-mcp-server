package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIRepoBranchListPrintsTableWithDefaultOrganization(t *testing.T) {
	var gotQuery string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/branches":
			if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/branches") {
				t.Fatalf("RequestURI = %q", r.RequestURI)
			}
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"data":[{"name":"main","default":true,"protected":false,"commit":{"id":"abc123"}}],"pagination":{"total":1}}`))
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
		"repo", "branch", "list",
		"--repository-id", "group/repo",
		"--page", "2",
		"--per-page", "10",
		"--sort", "updated_desc",
		"--search", "main",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"NAME", "DEFAULT", "PROTECTED", "LAST_COMMIT", "main", "true", "false", "abc123"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	for _, want := range []string{"page=2", "perPage=10", "sort=updated_desc", "search=main"} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %q", gotQuery, want)
		}
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/branches"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIRepoBranchesAliasListPrintsJSONWithExplicitOrganization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/codeup/organizations/org-2/repositories/repo-1/branches":
			_, _ = w.Write([]byte(`[{"name":"main"}]`))
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
		"repo", "branches", "list",
		"--organization-id", "org-2",
		"--repository-id", "repo-1",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"name": "main"`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIRepoBranchListRequiresRepositoryID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "branch", "list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected repository-id error")
	}
	if !strings.Contains(err.Error(), "repository-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIRepoBranchListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "codeup", "repo", "branch", "list", "--repository-id", "repo-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_branches"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestRepoBranchListOptionsParamsIncludesFilters(t *testing.T) {
	params, err := (repoBranchListOptions{
		OrganizationID: " org-1 ",
		RepositoryID:   " group/repo ",
		Page:           2,
		PerPage:        10,
		Sort:           " updated_desc ",
		Search:         " main ",
	}).params()
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"page":           2,
		"perPage":        10,
		"sort":           "updated_desc",
		"search":         "main",
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPrintRepoBranchListShowsNoResultsWhenRowsEmpty(t *testing.T) {
	var out bytes.Buffer
	raw := "No results found."
	if err := printRepoBranchList(&out, raw); err != nil {
		t.Fatalf("printRepoBranchList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want \"No results found.\"", out.String())
	}
}

func TestPrintRepoBranchListPrintsHeaderForEmptyList(t *testing.T) {
	var out bytes.Buffer
	if err := printRepoBranchList(&out, `{"data":[]}`); err != nil {
		t.Fatalf("printRepoBranchList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != "[1mNAME  DEFAULT  PROTECTED  LAST_COMMIT[0m" {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestRepoBranchRowsFromJSONExtractsAlternateFields(t *testing.T) {
	rows, _ := repoBranchRowsFromJSONForPrint(`{"result":{"items":[{"branchName":"dev","isDefault":false,"isProtected":true,"latestCommit":{"sha":"def456"}}]}}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := repoBranchRow{Name: "dev", Default: "false", Protected: "true", LastCommit: "def456"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestRepoBranchRowsFromJSONSkipsNonObjectRows(t *testing.T) {
	rows, _ := repoBranchRowsFromJSONForPrint(`{"data":["skip",{"name":"main"}]}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	if rows[0].Name != "main" {
		t.Fatalf("row = %#v", rows[0])
	}
}

func TestRepoBranchRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows, _ := repoBranchRowsFromJSONForPrint(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}
