package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIRepoChangeRequestViewPrintsJSONWithDefaultOrganization(t *testing.T) {
	requests := map[string]int{}
	var commentBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/changeRequests/12":
			if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/changeRequests/12") {
				t.Fatalf("RequestURI = %q", r.RequestURI)
			}
			_, _ = w.Write([]byte(`{"localId":"12","title":"Add CLI"}`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/changeRequests/12/diffs/patches":
			_, _ = w.Write([]byte(`[{"patchSetBizId":"p1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/repositories/group/repo/changeRequests/12/comments/list":
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s", r.Method)
			}
			if err := json.NewDecoder(r.Body).Decode(&commentBody); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			_, _ = w.Write([]byte(`[{"commentBizId":"c1"}]`))
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
		"repo", "change-request", "view", "12",
		"--repository-id", "group/repo",
		"--comment-state", "RESOLVED",
		"--comment-resolved",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out.String())
	}
	for _, section := range []string{"changeRequest", "patchSets", "comments"} {
		if _, ok := payload[section]; !ok {
			t.Fatalf("payload missing %s: %#v", section, payload)
		}
	}
	if commentBody["state"] != "RESOLVED" || commentBody["resolved"] != true || commentBody["comment_type"] != "GLOBAL_COMMENT" {
		t.Fatalf("comment body = %#v", commentBody)
	}
	if requests["/oapi/v1/platform/organizations"] != 1 ||
		requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/changeRequests/12"] != 1 ||
		requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/changeRequests/12/diffs/patches"] != 1 ||
		requests["/oapi/v1/codeup/organizations/org-1/repositories/group/repo/changeRequests/12/comments/list"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIRepoCRDetailAliasUsesExplicitOrganizationAndSkipsSections(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/codeup/organizations/org-2/repositories/repo-1/changeRequests/12":
			_, _ = w.Write([]byte(`{"localId":"12"}`))
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
		"repo", "cr", "detail", "12",
		"--organization-id", "org-2",
		"--repository-id", "repo-1",
		"--include-patch-sets=false",
		"--include-comments=false",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"changeRequest"`) {
		t.Fatalf("stdout = %q", out.String())
	}
	if strings.Contains(out.String(), `"patchSets"`) || strings.Contains(out.String(), `"comments"`) {
		t.Fatalf("stdout should omit optional sections: %q", out.String())
	}
}

func TestYunxiaoCLIRepoChangeRequestViewRequiresInputs(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "change-request", "view", "12"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected repository-id error")
	}
	if !strings.Contains(err.Error(), "repository-id is required") {
		t.Fatalf("error = %v", err)
	}

	command = NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "change-request", "view", "--repository-id", "repo-1"})
	err = command.Execute()
	if err == nil {
		t.Fatal("Execute() expected missing argument error")
	}
}

func TestYunxiaoCLIRepoChangeRequestViewReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "codeup", "repo", "change-request", "view", "12", "--repository-id", "repo-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "get_change_request_overview"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestRepoChangeRequestViewOptionsParamsIncludesOverviewFilters(t *testing.T) {
	params, err := (repoChangeRequestViewOptions{
		OrganizationID:   " org-1 ",
		RepositoryID:     " group/repo ",
		IncludePatchSets: true,
		IncludeComments:  false,
		CommentState:     " RESOLVED ",
		CommentResolved:  true,
	}).params(" 12 ")
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId":   "org-1",
		"repositoryId":     "group/repo",
		"localId":          "12",
		"includePatchSets": true,
		"includeComments":  false,
		"commentState":     "RESOLVED",
		"commentResolved":  true,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestRepoChangeRequestViewOptionsParamsRequiresInputs(t *testing.T) {
	if _, err := (repoChangeRequestViewOptions{}).params("12"); err == nil {
		t.Fatal("params() expected repository-id error")
	}
	if _, err := (repoChangeRequestViewOptions{RepositoryID: "repo-1"}).params(" "); err == nil {
		t.Fatal("params() expected local-id error")
	}
}
