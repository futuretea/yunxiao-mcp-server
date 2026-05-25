package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIRepoCommitViewPrintsJSONWithDefaultOrganization(t *testing.T) {
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/commits/abc123":
			if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/commits/abc123") {
				t.Fatalf("RequestURI = %q", r.RequestURI)
			}
			_, _ = w.Write([]byte(`{"sha":"abc123","title":"Fix login"}`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/commits/abc123/statuses":
			if r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("statuses query = %q", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`["success"]`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/checkRuns":
			if r.URL.Query().Get("ref") != "abc123" || r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("check runs query = %q", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`["check-1"]`))
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
		"repo", "commit", "view", "abc123",
		"--repository-id", "group/repo",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out.String())
	}
	commit, ok := payload["commit"].(map[string]any)
	if !ok || commit["sha"] != "abc123" || commit["title"] != "Fix login" {
		t.Fatalf("commit = %#v", payload["commit"])
	}
	if _, ok := payload["statuses"]; !ok {
		t.Fatalf("payload missing statuses: %#v", payload)
	}
	if _, ok := payload["checkRuns"]; !ok {
		t.Fatalf("payload missing checkRuns: %#v", payload)
	}
	if requests["/oapi/v1/platform/organizations"] != 1 ||
		requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/commits/abc123"] != 1 ||
		requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/commits/abc123/statuses"] != 1 ||
		requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/checkRuns"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIRepoCommitDetailAliasUsesExplicitOrganization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/codeup/organizations/org-2/repositories/repo-1/commits/abc123":
			_, _ = w.Write([]byte(`{"sha":"abc123"}`))
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
		"repo", "commit", "detail", "abc123",
		"--organization-id", "org-2",
		"--repository-id", "repo-1",
		"--include-statuses=false",
		"--include-check-runs=false",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"sha": "abc123"`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIRepoCommitViewRequiresIDs(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "commit", "view", "abc123"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected repository-id error")
	}
	if !strings.Contains(err.Error(), "repository-id is required") {
		t.Fatalf("error = %v", err)
	}

	command = NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "commit", "view", "--repository-id", "repo-1"})
	err = command.Execute()
	if err == nil {
		t.Fatal("Execute() expected missing argument error")
	}
}

func TestYunxiaoCLIRepoCommitViewReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "codeup", "repo", "commit", "view", "abc123", "--repository-id", "repo-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "get_commit_overview"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestRepoCommitViewOptionsParamsIncludesOrganization(t *testing.T) {
	params, err := (repoCommitViewOptions{
		OrganizationID:   " org-1 ",
		RepositoryID:     " group/repo ",
		IncludeStatuses:  true,
		IncludeCheckRuns: false,
		StatusLimit:      3,
		CheckRunLimit:    2,
	}).params(" abc123 ")
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId":   "org-1",
		"repositoryId":     "group/repo",
		"sha":              "abc123",
		"includeStatuses":  true,
		"includeCheckRuns": false,
		"statusLimit":      3,
		"checkRunLimit":    2,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestRepoCommitViewOptionsParamsRequiresIDs(t *testing.T) {
	if _, err := (repoCommitViewOptions{}).params("abc123"); err == nil {
		t.Fatal("params() expected repository-id error")
	}
	if _, err := (repoCommitViewOptions{RepositoryID: "repo-1"}).params(" "); err == nil {
		t.Fatal("params() expected sha error")
	}
}
