package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIRepoComparePrintsJSONWithDefaultOrganization(t *testing.T) {
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/compares":
			if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/compares") {
				t.Fatalf("RequestURI = %q", r.RequestURI)
			}
			query := r.URL.Query()
			wants := map[string]string{
				"from":       "main",
				"to":         "feature",
				"sourceType": "branch",
				"targetType": "branch",
				"straight":   "true",
			}
			for key, want := range wants {
				if got := query.Get(key); got != want {
					t.Fatalf("query[%q] = %q, want %q; raw=%q", key, got, want, r.URL.RawQuery)
				}
			}
			_, _ = w.Write([]byte(`{"commits":[{"id":"sha-1"}],"diffs":[{"newPath":"README.md"}]}`))
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
		"repo", "compare", "main", "feature",
		"--repository-id", "group/repo",
		"--source-type", "branch",
		"--target-type", "branch",
		"--straight", "true",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out.String())
	}
	if _, ok := payload["commits"]; !ok {
		t.Fatalf("payload = %#v", payload)
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/compares"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIRepoDiffAliasUsesExplicitOrganization(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/codeup/organizations/org-2/repositories/repo-1/compares":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"diffs":[]}`))
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
		"repo", "diff", "sha-a", "sha-b",
		"--organization-id", "org-2",
		"--repository-id", "repo-1",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"diffs": []`) {
		t.Fatalf("stdout = %q", out.String())
	}
	if gotQuery != "from=sha-a&to=sha-b" {
		t.Fatalf("query = %q", gotQuery)
	}
}

func TestYunxiaoCLIRepoCompareRequiresInputs(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "compare", "main", "feature"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected repository-id error")
	}
	if !strings.Contains(err.Error(), "repository-id is required") {
		t.Fatalf("error = %v", err)
	}

	command = NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "compare", "main", "--repository-id", "repo-1"})
	err = command.Execute()
	if err == nil {
		t.Fatal("Execute() expected missing argument error")
	}
}

func TestYunxiaoCLIRepoCompareReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "codeup", "repo", "compare", "main", "feature", "--repository-id", "repo-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "compare"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestRepoCompareOptionsParamsIncludesOptionalFilters(t *testing.T) {
	params, err := (repoCompareOptions{
		OrganizationID: " org-1 ",
		RepositoryID:   " group/repo ",
		SourceType:     " branch ",
		TargetType:     " commit ",
		Straight:       " false ",
	}).params(" main ", " sha-b ")
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"from":           "main",
		"to":             "sha-b",
		"sourceType":     "branch",
		"targetType":     "commit",
		"straight":       "false",
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestRepoCompareOptionsParamsRequiresInputs(t *testing.T) {
	if _, err := (repoCompareOptions{}).params("main", "feature"); err == nil {
		t.Fatal("params() expected repository-id error")
	}
	if _, err := (repoCompareOptions{RepositoryID: "repo-1"}).params(" ", "feature"); err == nil {
		t.Fatal("params() expected from error")
	}
	if _, err := (repoCompareOptions{RepositoryID: "repo-1"}).params("main", " "); err == nil {
		t.Fatal("params() expected to error")
	}
}
