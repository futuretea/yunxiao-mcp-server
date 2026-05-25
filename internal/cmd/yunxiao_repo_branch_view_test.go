package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIRepoBranchViewPrintsJSONWithDefaultOrganization(t *testing.T) {
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/branches/feature/demo":
			if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/branches/feature%2Fdemo") {
				t.Fatalf("RequestURI = %q", r.RequestURI)
			}
			_, _ = w.Write([]byte(`{"name":"feature/demo"}`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/commits":
			query := r.URL.Query()
			if query.Get("refName") != "feature/demo" || query.Get("perPage") != "2" {
				t.Fatalf("commits query = %q", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`[{"id":"sha-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/mergeRequests":
			query := r.URL.Query()
			if query.Get("repositoryIds") != "group/repo" ||
				query.Get("targetBranch") != "feature/demo" ||
				query.Get("state") != "merged" ||
				query.Get("perPage") != "4" {
				t.Fatalf("merge requests query = %q", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`[{"localId":"12"}]`))
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
		"repo", "branch", "view", "feature/demo",
		"--repository-id", "group/repo",
		"--commit-limit", "2",
		"--merge-request-limit", "4",
		"--merge-request-state", "merged",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out.String())
	}
	for _, section := range []string{"branch", "commits", "mergeRequests"} {
		if _, ok := payload[section]; !ok {
			t.Fatalf("payload missing %s: %#v", section, payload)
		}
	}
	if requests["/oapi/v1/platform/organizations"] != 1 ||
		requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/branches/feature/demo"] != 1 ||
		requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/commits"] != 1 ||
		requests["/oapi/v1/codeup/organizations/org-1/mergeRequests"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIRepoBranchOverviewAliasUsesExplicitOrganizationAndSkipsSections(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/codeup/organizations/org-2/repositories/repo-1/branches/main":
			_, _ = w.Write([]byte(`{"name":"main"}`))
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
		"repo", "branches", "overview", "main",
		"--organization-id", "org-2",
		"--repository-id", "repo-1",
		"--include-commits=false",
		"--include-merge-requests=false",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"branch"`) {
		t.Fatalf("stdout = %q", out.String())
	}
	if strings.Contains(out.String(), `"commits"`) || strings.Contains(out.String(), `"mergeRequests"`) {
		t.Fatalf("stdout should omit optional sections: %q", out.String())
	}
}

func TestYunxiaoCLIRepoBranchViewRequiresInputs(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "branch", "view", "main"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected repository-id error")
	}
	if !strings.Contains(err.Error(), "repository-id is required") {
		t.Fatalf("error = %v", err)
	}

	command = NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "branch", "view", "--repository-id", "repo-1"})
	err = command.Execute()
	if err == nil {
		t.Fatal("Execute() expected missing argument error")
	}
}

func TestYunxiaoCLIRepoBranchViewReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "codeup", "repo", "branch", "view", "main", "--repository-id", "repo-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "get_branch_overview"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestRepoBranchViewOptionsParamsIncludesOverviewFilters(t *testing.T) {
	params, err := (repoBranchViewOptions{
		OrganizationID:       " org-1 ",
		RepositoryID:         " group/repo ",
		IncludeCommits:       true,
		IncludeMergeRequests: false,
		CommitLimit:          2,
		MergeRequestLimit:    4,
		MergeRequestState:    " merged ",
	}).params(" feature/demo ")
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId":       "org-1",
		"repositoryId":         "group/repo",
		"branchName":           "feature/demo",
		"includeCommits":       true,
		"includeMergeRequests": false,
		"commitLimit":          2,
		"mrLimit":              4,
		"mrState":              "merged",
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestRepoBranchViewOptionsParamsRequiresInputs(t *testing.T) {
	if _, err := (repoBranchViewOptions{}).params("main"); err == nil {
		t.Fatal("params() expected repository-id error")
	}
	if _, err := (repoBranchViewOptions{RepositoryID: "repo-1"}).params(" "); err == nil {
		t.Fatal("params() expected branch-name error")
	}
}
