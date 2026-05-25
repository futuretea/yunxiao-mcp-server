package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIRepoChangeRequestListPrintsTableWithDefaultOrganization(t *testing.T) {
	var gotQuery string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/changeRequests":
			gotQuery = r.URL.RawQuery
			w.Header().Set("x-total", "1")
			_, _ = w.Write([]byte(`{"data":[{"localId":12,"repositoryId":"group/repo","title":"Add CLI","state":"opened","author":{"name":"Ada"},"sourceBranch":"feature","targetBranch":"main"}],"pagination":{"total":1}}`))
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
		"repo", "change-request", "list",
		"--project-ids", "group/repo,repo-2",
		"--author-ids", "user-1",
		"--reviewer-ids", "user-2",
		"--state", "opened",
		"--search", "CLI",
		"--order-by", "updated_at",
		"--sort", "desc",
		"--created-after", "2026-05-01T00:00:00Z",
		"--created-before", "2026-05-25T00:00:00Z",
		"--page", "2",
		"--per-page", "10",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"LOCAL_ID", "REPOSITORY", "TITLE", "STATE", "AUTHOR", "SOURCE", "TARGET", "12", "group/repo", "Add CLI", "opened", "Ada", "feature", "main"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	for _, want := range []string{
		"projectIds=group%2Frepo%2Crepo-2",
		"authorIds=user-1",
		"reviewerIds=user-2",
		"state=opened",
		"search=CLI",
		"orderBy=updated_at",
		"sort=desc",
		"createdAfter=2026-05-01T00%3A00%3A00Z",
		"createdBefore=2026-05-25T00%3A00%3A00Z",
		"page=2",
		"perPage=10",
	} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %q", gotQuery, want)
		}
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/codeup/organizations/org-1/changeRequests"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIRepoCRAliasListPrintsJSONWithExplicitOrganization(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/codeup/organizations/org-2/changeRequests":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[{"localId":"12","title":"Add CLI"}]`))
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
		"repo", "cr", "list",
		"--organization-id", "org-2",
		"--state", "merged",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"localId": "12"`) {
		t.Fatalf("stdout = %q", out.String())
	}
	if gotQuery != "state=merged" {
		t.Fatalf("query = %q, want state=merged", gotQuery)
	}
}

func TestYunxiaoCLIRepoChangeRequestsAliasResolves(t *testing.T) {
	var out bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"repo", "change-requests", "list", "--help"})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), "list CodeUp change requests") {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIRepoChangeRequestListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "codeup", "repo", "change-request", "list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_change_requests"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestRepoChangeRequestListOptionsParamsIncludesFilters(t *testing.T) {
	params := (repoChangeRequestListOptions{
		OrganizationID: " org-1 ",
		ProjectIDs:     " group/repo ",
		AuthorIDs:      " user-1 ",
		ReviewerIDs:    " user-2 ",
		State:          " opened ",
		Search:         " CLI ",
		OrderBy:        " updated_at ",
		Sort:           " desc ",
		CreatedBefore:  " 2026-05-25T00:00:00Z ",
		CreatedAfter:   " 2026-05-01T00:00:00Z ",
		Page:           2,
		PerPage:        10,
	}).params()

	wants := map[string]any{
		"organizationId": "org-1",
		"projectIds":     "group/repo",
		"authorIds":      "user-1",
		"reviewerIds":    "user-2",
		"state":          "opened",
		"search":         "CLI",
		"orderBy":        "updated_at",
		"sort":           "desc",
		"createdBefore":  "2026-05-25T00:00:00Z",
		"createdAfter":   "2026-05-01T00:00:00Z",
		"page":           2,
		"perPage":        10,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPrintRepoChangeRequestListFallsBackToRawJSON(t *testing.T) {
	var out bytes.Buffer
	raw := `{"data":{"total":0}}`
	if err := printRepoChangeRequestList(&out, raw); err != nil {
		t.Fatalf("printRepoChangeRequestList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want raw JSON", out.String())
	}
}

func TestRepoChangeRequestRowsFromJSONExtractsAlternateFields(t *testing.T) {
	rows := repoChangeRequestRowsFromJSON(`{"result":{"items":[{"iid":"12","project":{"displayName":"Demo Repo"},"subject":"Add CLI","status":"opened","creator":{"displayName":"Ada"},"source_branch":"feature","target_branch":"main"}]}}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := repoChangeRequestRow{
		LocalID:      "12",
		Repository:   "Demo Repo",
		Title:        "Add CLI",
		State:        "opened",
		Author:       "Ada",
		SourceBranch: "feature",
		TargetBranch: "main",
	}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}
