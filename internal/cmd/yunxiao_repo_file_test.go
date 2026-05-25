package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIRepoFileListPrintsTableWithDefaultOrganization(t *testing.T) {
	var gotQuery string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/files/tree":
			if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/files/tree") {
				t.Fatalf("RequestURI = %q", r.RequestURI)
			}
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"data":[{"path":"cmd/main.go","type":"blob","size":123,"mode":"100644"}]}`))
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
		"repo", "file", "list",
		"--repository-id", "group/repo",
		"--path", "cmd",
		"--ref", "main",
		"--type", "DIRECT",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"PATH", "TYPE", "SIZE", "MODE", "cmd/main.go", "blob", "123", "100644"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	for _, want := range []string{"path=cmd", "ref=main", "type=DIRECT"} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %q", gotQuery, want)
		}
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/files/tree"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIRepoFilesAliasListPrintsJSONWithExplicitOrganization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/codeup/organizations/org-2/repositories/repo-1/files/tree":
			_, _ = w.Write([]byte(`[{"path":"README.md"}]`))
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
		"repo", "files", "list",
		"--organization-id", "org-2",
		"--repository-id", "repo-1",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"path": "README.md"`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIRepoFileListRequiresRepositoryID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "file", "list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected repository-id error")
	}
	if !strings.Contains(err.Error(), "repository-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIRepoFileListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "codeup", "repo", "file", "list", "--repository-id", "repo-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_files"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestRepoFileListOptionsParamsIncludesFilters(t *testing.T) {
	params, err := (repoFileListOptions{
		OrganizationID: " org-1 ",
		RepositoryID:   " group/repo ",
		Path:           " cmd ",
		Ref:            " main ",
		TreeType:       " DIRECT ",
	}).params()
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"path":           "cmd",
		"ref":            "main",
		"type":           "DIRECT",
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPrintRepoFileListPrintsHeaderForEmptyList(t *testing.T) {
	var out bytes.Buffer
	if err := printRepoFileList(&out, `{"data":[]}`); err != nil {
		t.Fatalf("printRepoFileList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != "\033[1mPATH  TYPE  SIZE  MODE\033[0m" {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestPrintRepoFileListShowsNoResultsWhenRowsEmpty(t *testing.T) {
	var out bytes.Buffer
	raw := "No results found."
	if err := printRepoFileList(&out, raw); err != nil {
		t.Fatalf("printRepoFileList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want \"No results found.\"", out.String())
	}
}

func TestRepoFileRowsFromJSONExtractsAlternateFields(t *testing.T) {
	rows := repoFileRowsFromJSON(`{"result":{"items":[{"filePath":"README.md","kind":"blob","fileSize":456,"mode":"100644"}]}}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := repoFileRow{Path: "README.md", Type: "blob", Size: "456", Mode: "100644"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestRepoFileRowsFromJSONSkipsNonObjectRows(t *testing.T) {
	rows := repoFileRowsFromJSON(`{"data":["skip",{"path":"README.md"}]}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	if rows[0].Path != "README.md" {
		t.Fatalf("row = %#v", rows[0])
	}
}

func TestRepoFileRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows := repoFileRowsFromJSON(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}
