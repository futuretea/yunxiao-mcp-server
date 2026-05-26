package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIRepoMrListPrintsTableWithDefaultOrganization(t *testing.T) {
	var gotQuery string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/mergeRequests":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"data":[{"id":"1","iid":"42","title":"Fix bug","state":"opened","author":{"username":"alice"},"targetBranch":"main"}],"pagination":{"total":1}}`))
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
		"repo", "mr", "list",
		"--state", "opened",
		"--page", "2",
		"--per-page", "20",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"ID", "IID", "TITLE", "STATE", "AUTHOR", "TARGET", "1", "42", "Fix bug", "opened", "alice", "main"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	if !strings.Contains(gotQuery, "state=opened") || !strings.Contains(gotQuery, "page=2") || !strings.Contains(gotQuery, "perPage=20") {
		t.Fatalf("query = %q", gotQuery)
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/codeup/organizations/org-1/mergeRequests"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIRepoMergeRequestAliasListPrintsJSONWithExplicitOrganization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/codeup/organizations/org-2/mergeRequests":
			_, _ = w.Write([]byte(`[{"id":"1","title":"MR 1"}]`))
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
		"repos", "merge-request", "list",
		"--organization-id", "org-2",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"title": "MR 1"`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIRepoMrListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "codeup", "repo", "mr", "list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_merge_requests"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestMRListOptionsParamsIncludesFilters(t *testing.T) {
	params := (mrListOptions{
		OrganizationID:  " org-1 ",
		State:           " opened ",
		Search:          " demo ",
		AuthorUserIDs:   " alice ",
		AssigneeUserIDs: " bob ",
		OrderBy:         " updated_at ",
		TargetBranch:    " main ",
		CreatedAfter:    " 2026-01-01 ",
		CreatedBefore:   " 2026-02-01 ",
		Page:            2,
		PerPage:         20,
	}).params()

	wants := map[string]any{
		"organizationId":  "org-1",
		"state":           "opened",
		"search":          "demo",
		"authorUserIds":   []string{"alice"},
		"assigneeUserIds": []string{"bob"},
		"orderBy":         "updated_at",
		"targetBranch":    "main",
		"createdAfter":    "2026-01-01",
		"createdBefore":   "2026-02-01",
		"page":            2,
		"perPage":         20,
	}
	for key, want := range wants {
		got := params[key]
		if key == "authorUserIds" || key == "assigneeUserIds" {
			gotSlice, _ := got.([]string)
			wantSlice, _ := want.([]string)
			if len(gotSlice) != len(wantSlice) || (len(gotSlice) > 0 && gotSlice[0] != wantSlice[0]) {
				t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
			}
			continue
		}
		if got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestMRListOptionsParamsOmitEmptyArrays(t *testing.T) {
	params := (mrListOptions{OrganizationID: "org-1"}).params()
	if _, ok := params["authorUserIds"]; ok {
		t.Fatal("params should not include authorUserIds when empty")
	}
	if _, ok := params["assigneeUserIds"]; ok {
		t.Fatal("params should not include assigneeUserIds when empty")
	}
}

func TestPrintMRListShowsNoResultsWhenRowsEmpty(t *testing.T) {
	var out bytes.Buffer
	raw := "No results found."
	if err := printMRList(&out, raw); err != nil {
		t.Fatalf("printMRList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want \"No results found.\"", out.String())
	}
}

func TestPrintMRListPrintsHeaderForEmptyList(t *testing.T) {
	var out bytes.Buffer
	if err := printMRList(&out, `{"data":[]}`); err != nil {
		t.Fatalf("printMRList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != "\033[1mID  IID  TITLE  STATE  AUTHOR  TARGET\033[0m" {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestMRRowsFromJSONExtractsAlternateFields(t *testing.T) {
	rows, _ := mrRowsFromJSONForPrint(`{"result":{"items":[{"mergeRequestId":"mr-1","localId":"12","name":"Fix login","status":"merged","authorName":"carol","target":"develop"}]}}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := mrRow{ID: "mr-1", IID: "12", Title: "Fix login", State: "merged", Author: "carol", Target: "develop"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestMRRowsFromJSONSkipsNonObjectRows(t *testing.T) {
	rows, _ := mrRowsFromJSONForPrint(`{"data":["skip",{"id":"2","iid":"43","title":"Update docs","state":"closed","authorUsername":"dan","targetBranchName":"staging"}]}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := mrRow{ID: "2", IID: "43", Title: "Update docs", State: "closed", Author: "dan", Target: "staging"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestYunxiaoCLIRepoMrViewPrintsJSONWithDefaultOrganization(t *testing.T) {
	var gotPath string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/mergeRequests/42":
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"iid":42,"title":"Fix bug","state":"merged"}`))
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
		"repo", "mr", "view", "42",
		"--repository-id", "group/repo",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"iid":`) {
		t.Fatalf("stdout = %q", out.String())
	}
	if gotPath != "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/mergeRequests/42" {
		t.Fatalf("path = %q", gotPath)
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/mergeRequests/42"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIRepoMrViewRequiresRepositoryID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "mr", "view", "42"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected repository-id error")
	}
	if !strings.Contains(err.Error(), "repository-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIRepoMrViewReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "codeup", "repo", "mr", "view", "42", "--repository-id", "group/repo"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "get_merge_request"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestMRViewOptionsParamsIncludesBothKeys(t *testing.T) {
	params, err := (mrViewOptions{
		OrganizationID: " org-1 ",
		RepositoryID:   " group/repo ",
		MergeRequestID: " 42 ",
	}).params()
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}
	if params["mergeRequestId"] != "42" {
		t.Fatalf("mergeRequestId = %q", params["mergeRequestId"])
	}
	if params["iid"] != "42" {
		t.Fatalf("iid = %q", params["iid"])
	}
	if params["repositoryId"] != "group/repo" {
		t.Fatalf("repositoryId = %q", params["repositoryId"])
	}
}

func TestMRViewOptionsParamsRequiresRepositoryID(t *testing.T) {
	if _, err := (mrViewOptions{MergeRequestID: "42"}).params(); err == nil {
		t.Fatal("params() expected repository-id error")
	}
}

func TestMRRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows, _ := mrRowsFromJSONForPrint(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}
