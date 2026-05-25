package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIRepoCommitListPrintsTableWithDefaultOrganization(t *testing.T) {
	var gotQuery string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/commits":
			if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/commits") {
				t.Fatalf("RequestURI = %q", r.RequestURI)
			}
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"data":[{"id":"abcdef123456","shortId":"abcdef1","title":"Fix login","authorName":"Alice","authoredDate":"2026-05-24T10:00:00Z"}],"pagination":{"total":1}}`))
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
		"repo", "commit", "list",
		"--repository-id", "group/repo",
		"--ref", "main",
		"--since", "2026-05-01T00:00:00Z",
		"--until", "2026-05-25T00:00:00Z",
		"--page", "2",
		"--per-page", "10",
		"--path", "cmd/main.go",
		"--search", "login",
		"--committer-ids", "user-1,user-2",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"SHA", "SHORT_ID", "TITLE", "AUTHOR", "DATE", "abcdef123456", "abcdef1", "Fix login", "Alice"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	for _, want := range []string{"refName=main", "since=2026-05-01T00%3A00%3A00Z", "until=2026-05-25T00%3A00%3A00Z", "page=2", "perPage=10", "path=cmd%2Fmain.go", "search=login", "committerIds=user-1%2Cuser-2"} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %q", gotQuery, want)
		}
	}
	if strings.Contains(gotQuery, "showSignature=") {
		t.Fatalf("query = %q, should not include showSignature when flag is absent", gotQuery)
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/commits"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIRepoCommitsAliasListPrintsJSONWithExplicitOrganizationAndSignatureFalse(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/codeup/organizations/org-2/repositories/repo-1/commits":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[{"id":"abcdef123456"}]`))
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
		"repo", "commits", "list",
		"--organization-id", "org-2",
		"--repository-id", "repo-1",
		"--ref", "main",
		"--show-signature=false",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"id": "abcdef123456"`) {
		t.Fatalf("stdout = %q", out.String())
	}
	for _, want := range []string{"refName=main", "showSignature=false"} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %q", gotQuery, want)
		}
	}
}

func TestYunxiaoCLIRepoCommitListPassesShowSignatureTrue(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/codeup/organizations/org-1/repositories/repo-1/commits":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[]`))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"repo", "commit", "list",
		"--organization-id", "org-1",
		"--repository-id", "repo-1",
		"--ref", "main",
		"--show-signature",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"refName=main", "showSignature=true"} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %q", gotQuery, want)
		}
	}
}

func TestYunxiaoCLIRepoCommitListRequiresIDs(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "commit", "list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected repository-id error")
	}
	if !strings.Contains(err.Error(), "repository-id is required") {
		t.Fatalf("error = %v", err)
	}

	command = NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "commit", "list", "--repository-id", "repo-1"})
	err = command.Execute()
	if err == nil {
		t.Fatal("Execute() expected ref error")
	}
	if !strings.Contains(err.Error(), "ref is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIRepoCommitListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "codeup", "repo", "commit", "list", "--repository-id", "repo-1", "--ref", "main"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_commits"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestRepoCommitListOptionsParamsIncludesFilters(t *testing.T) {
	params, err := (repoCommitListOptions{
		OrganizationID:   " org-1 ",
		RepositoryID:     " group/repo ",
		Ref:              " main ",
		Since:            " 2026-05-01T00:00:00Z ",
		Until:            " 2026-05-25T00:00:00Z ",
		Page:             2,
		PerPage:          10,
		Path:             " cmd/main.go ",
		Search:           " login ",
		ShowSignature:    false,
		ShowSignatureSet: true,
		CommitterIDs:     " user-1,user-2 ",
	}).params()
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"refName":        "main",
		"since":          "2026-05-01T00:00:00Z",
		"until":          "2026-05-25T00:00:00Z",
		"page":           2,
		"perPage":        10,
		"path":           "cmd/main.go",
		"search":         "login",
		"showSignature":  false,
		"committerIds":   "user-1,user-2",
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPrintRepoCommitListPrintsHeaderForEmptyList(t *testing.T) {
	var out bytes.Buffer
	if err := printRepoCommitList(&out, `{"data":[]}`); err != nil {
		t.Fatalf("printRepoCommitList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != "SHA  SHORT_ID  TITLE  AUTHOR  DATE" {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestPrintRepoCommitListFallsBackToRawJSON(t *testing.T) {
	var out bytes.Buffer
	raw := `{"data":{"total":0}}`
	if err := printRepoCommitList(&out, raw); err != nil {
		t.Fatalf("printRepoCommitList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want raw JSON", out.String())
	}
}

func TestRepoCommitRowsFromJSONExtractsAlternateFields(t *testing.T) {
	rows := repoCommitRowsFromJSON(`{"result":{"items":[{"sha":"abcdef123456","short_id":"abcdef1","message":"Fix login","author":{"displayName":"Alice"},"committed_date":"2026-05-24"}]}}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := repoCommitRow{SHA: "abcdef123456", ShortID: "abcdef1", Title: "Fix login", Author: "Alice", Date: "2026-05-24"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestRepoCommitRowsFromJSONSkipsNonObjectRows(t *testing.T) {
	rows := repoCommitRowsFromJSON(`{"data":["skip",{"id":"abcdef123456"}]}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	if rows[0].SHA != "abcdef123456" {
		t.Fatalf("row = %#v", rows[0])
	}
}

func TestRepoCommitRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows := repoCommitRowsFromJSON(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}
