package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIRepoFileViewPrintsJSONWithDefaultOrganization(t *testing.T) {
	requests := map[string]int{}
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/files/src/main.go":
			if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/files/src%2Fmain.go") {
				t.Fatalf("RequestURI = %q", r.RequestURI)
			}
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"filePath":"src/main.go","content":"package main"}`))
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
		"repo", "file", "view", "src/main.go",
		"--repository-id", "group/repo",
		"--ref", "main",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out.String())
	}
	if payload["filePath"] != "src/main.go" || payload["content"] != "package main" {
		t.Fatalf("payload = %#v", payload)
	}
	if gotQuery != "ref=main" {
		t.Fatalf("query = %q", gotQuery)
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/files/src/main.go"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIRepoFileContentAliasUsesExplicitOrganizationAndSince(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/codeup/organizations/org-2/repositories/repo-1/files/README.md":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"filePath":"README.md"}`))
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
		"repo", "file", "content", "README.md",
		"--organization-id", "org-2",
		"--repository-id", "repo-1",
		"--ref", "main",
		"--since", "2026-05-01T00:00:00Z",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"filePath": "README.md"`) {
		t.Fatalf("stdout = %q", out.String())
	}
	for _, want := range []string{"ref=main", "since=2026-05-01T00%3A00%3A00Z"} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %q", gotQuery, want)
		}
	}
}

func TestYunxiaoCLIRepoFileViewRequiresIDs(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "file", "view", "README.md"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected repository-id error")
	}
	if !strings.Contains(err.Error(), "repository-id is required") {
		t.Fatalf("error = %v", err)
	}

	command = NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "file", "view", "README.md", "--repository-id", "repo-1"})
	err = command.Execute()
	if err == nil {
		t.Fatal("Execute() expected ref error")
	}
	if !strings.Contains(err.Error(), "ref is required") {
		t.Fatalf("error = %v", err)
	}

	command = NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "file", "view", "--repository-id", "repo-1", "--ref", "main"})
	err = command.Execute()
	if err == nil {
		t.Fatal("Execute() expected missing argument error")
	}
}

func TestYunxiaoCLIRepoFileViewReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "codeup", "repo", "file", "view", "README.md", "--repository-id", "repo-1", "--ref", "main"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "get_file_blobs"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestRepoFileViewOptionsParamsIncludesOrganization(t *testing.T) {
	params, err := (repoFileViewOptions{
		OrganizationID: " org-1 ",
		RepositoryID:   " group/repo ",
		Ref:            " main ",
		Since:          " 2026-05-01T00:00:00Z ",
	}).params(" src/main.go ")
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"filePath":       "src/main.go",
		"ref":            "main",
		"since":          "2026-05-01T00:00:00Z",
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestRepoFileViewOptionsParamsRequiresIDs(t *testing.T) {
	if _, err := (repoFileViewOptions{}).params("README.md"); err == nil {
		t.Fatal("params() expected repository-id error")
	}
	if _, err := (repoFileViewOptions{RepositoryID: "repo-1", Ref: "main"}).params(" "); err == nil {
		t.Fatal("params() expected path error")
	}
	if _, err := (repoFileViewOptions{RepositoryID: "repo-1"}).params("README.md"); err == nil {
		t.Fatal("params() expected ref error")
	}
}
