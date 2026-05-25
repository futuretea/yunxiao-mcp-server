package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIRepoListPrintsTableWithDefaultOrganization(t *testing.T) {
	var gotQuery string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/codeup/organizations/org-1/repositories":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"data":[{"id":123,"name":"Demo","path":"group/demo","archived":false,"lastActivityAt":"2026-05-24T10:00:00Z"}],"pagination":{"total":1}}`))
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
		"repo", "list",
		"--page", "2",
		"--per-page", "10",
		"--order-by", "name",
		"--sort", "asc",
		"--search", "demo",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"ID", "NAME", "PATH", "ARCHIVED", "LAST_ACTIVITY", "123", "Demo", "group/demo", "false", "2026-05-24T10:00:00Z"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	for _, want := range []string{"page=2", "perPage=10", "orderBy=name", "sort=asc", "search=demo"} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %q", gotQuery, want)
		}
	}
	if strings.Contains(gotQuery, "archived=") {
		t.Fatalf("query = %q, should not include archived when flag is absent", gotQuery)
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/codeup/organizations/org-1/repositories"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIRepositoriesAliasListPrintsJSONWithExplicitOrganizationAndArchivedFalse(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/codeup/organizations/org-2/repositories":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[{"id":"repo-1","name":"Demo"}]`))
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
		"repositories", "list",
		"--organization-id", "org-2",
		"--archived=false",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"id": "repo-1"`) {
		t.Fatalf("stdout = %q", out.String())
	}
	if gotQuery != "archived=false" {
		t.Fatalf("query = %q, want archived=false", gotQuery)
	}
}

func TestYunxiaoCLIRepoListAliasesResolve(t *testing.T) {
	for _, alias := range []string{"repos", "repository"} {
		t.Run(alias, func(t *testing.T) {
			var out bytes.Buffer
			command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
			command.SetArgs([]string{alias, "list", "--help"})

			if err := command.Execute(); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if !strings.Contains(out.String(), "list CodeUp repositories") {
				t.Fatalf("stdout = %q", out.String())
			}
		})
	}
}

func TestYunxiaoCLIRepoListPassesArchivedTrue(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/codeup/organizations/org-1/repositories":
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
		"repo", "list",
		"--organization-id", "org-1",
		"--archived",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if gotQuery != "archived=true" {
		t.Fatalf("query = %q, want archived=true", gotQuery)
	}
}

func TestYunxiaoCLIRepoListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "codeup", "repo", "list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_repositories"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestRepoListOptionsParamsIncludesFilters(t *testing.T) {
	params := (repoListOptions{
		OrganizationID: " org-1 ",
		Page:           2,
		PerPage:        10,
		OrderBy:        " name ",
		Sort:           " asc ",
		Search:         " demo ",
		Archived:       false,
		ArchivedSet:    true,
	}).params()

	wants := map[string]any{
		"organizationId": "org-1",
		"page":           2,
		"perPage":        10,
		"orderBy":        "name",
		"sort":           "asc",
		"search":         "demo",
		"archived":       false,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestRepoListOptionsParamsOmitsAbsentArchived(t *testing.T) {
	params := (repoListOptions{Archived: false}).params()
	if _, ok := params["archived"]; ok {
		t.Fatalf("params = %#v, should omit archived when flag is absent", params)
	}
}

func TestPrintRepoListFallsBackToRawJSON(t *testing.T) {
	var out bytes.Buffer
	raw := `{"data":{"total":0}}`
	if err := printRepoList(&out, raw); err != nil {
		t.Fatalf("printRepoList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want raw JSON", out.String())
	}
}

func TestRepoRowsFromJSONExtractsAlternateFields(t *testing.T) {
	rows := repoRowsFromJSON(`{"result":{"items":[{"repositoryId":"repo-1","displayName":"Demo","pathWithNamespace":"group/demo","isArchived":true,"last_activity_at":"2026-05-24"}]}}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := repoRow{ID: "repo-1", Name: "Demo", Path: "group/demo", Archived: "true", LastActivity: "2026-05-24"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestRepoRowsFromJSONSkipsNonObjectRows(t *testing.T) {
	rows := repoRowsFromJSON(`{"data":["skip",{"id":"repo-1","name":"Demo"}]}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	if rows[0].ID != "repo-1" {
		t.Fatalf("row = %#v", rows[0])
	}
}

func TestRepoRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows := repoRowsFromJSON(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}
