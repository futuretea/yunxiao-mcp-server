package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIRepoViewPrintsJSONWithDefaultOrganization(t *testing.T) {
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo":
			if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo") {
				t.Fatalf("RequestURI = %q", r.RequestURI)
			}
			_, _ = w.Write([]byte(`{"id":"repo-1","name":"demo","defaultBranch":"main"}`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/branches":
			if r.URL.Query().Get("perPage") != "3" {
				t.Fatalf("branches query = %q", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`[{"name":"main"}]`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/commits":
			query := r.URL.Query()
			if query.Get("refName") != "main" || query.Get("perPage") != "2" {
				t.Fatalf("commits query = %q", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`[{"id":"sha-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/mergeRequests":
			query := r.URL.Query()
			if query.Get("repositoryIds") != "group/repo" || query.Get("state") != "merged" || query.Get("perPage") != "4" {
				t.Fatalf("merge requests query = %q", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`[{"localId":"1"}]`))
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
		"repo", "view", "group/repo",
		"--branch-limit", "3",
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
	if _, ok := payload["repository"]; !ok {
		t.Fatalf("payload missing repository: %#v", payload)
	}
	for _, section := range []string{"branches", "commits", "mergeRequests"} {
		if _, ok := payload[section]; !ok {
			t.Fatalf("payload missing %s: %#v", section, payload)
		}
	}
	if requests["/oapi/v1/platform/organizations"] != 1 ||
		requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo"] != 1 ||
		requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/branches"] != 1 ||
		requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/commits"] != 1 ||
		requests["/oapi/v1/codeup/organizations/org-1/mergeRequests"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIRepoOverviewAliasUsesExplicitOrganizationAndSkipsSections(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/codeup/organizations/org-2/repositories/repo-1":
			_, _ = w.Write([]byte(`{"id":"repo-1","name":"demo","defaultBranch":"main"}`))
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
		"repo", "overview", "repo-1",
		"--organization-id", "org-2",
		"--include-branches=false",
		"--include-commits=false",
		"--include-merge-requests=false",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"repository"`) {
		t.Fatalf("stdout = %q", out.String())
	}
	if strings.Contains(out.String(), `"branches"`) || strings.Contains(out.String(), `"commits"`) || strings.Contains(out.String(), `"mergeRequests"`) {
		t.Fatalf("stdout should omit optional sections: %q", out.String())
	}
}

func TestYunxiaoCLIRepoViewRequiresID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "view"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected missing argument error")
	}
}

func TestYunxiaoCLIRepoViewReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "codeup", "repo", "view", "repo-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "get_repository_overview"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestRepoViewOptionsParamsIncludesOverviewFilters(t *testing.T) {
	params, err := (repoViewOptions{
		OrganizationID:       " org-1 ",
		IncludeBranches:      true,
		IncludeCommits:       false,
		IncludeMergeRequests: true,
		Ref:                  " main ",
		BranchLimit:          3,
		CommitLimit:          2,
		MergeRequestLimit:    4,
		MergeRequestState:    " merged ",
	}).params(" group/repo ")
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId":       "org-1",
		"repositoryId":         "group/repo",
		"includeBranches":      true,
		"includeCommits":       false,
		"includeMergeRequests": true,
		"refName":              "main",
		"branchLimit":          3,
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

func TestRepoViewOptionsParamsRequiresID(t *testing.T) {
	if _, err := (repoViewOptions{}).params(" "); err == nil {
		t.Fatal("params() expected repository-id error")
	}
}
