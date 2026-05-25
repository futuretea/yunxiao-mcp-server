package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIMemberListPrintsTableWithDefaultOrganization(t *testing.T) {
	var gotQuery string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/platform/organizations/org-1/members":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[{"id":"member-1","userId":"user-1","displayName":"Alice","email":"alice@example.com"}]`))
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
		"member", "list",
		"--page", "2",
		"--per-page", "50",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"MEMBER_ID", "USER_ID", "NAME", "EMAIL", "member-1", "user-1", "Alice", "alice@example.com"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	for _, want := range []string{"page=2", "perPage=50"} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %q", gotQuery, want)
		}
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/platform/organizations/org-1/members"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIMembersAliasListPrintsJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/platform/organizations/org-1/members":
			_, _ = w.Write([]byte(`[{"userId":"user-1","displayName":"Alice"}]`))
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
		"members", "list",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"userId": "user-1"`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIMemberListSkipsDefaultOrgWhenOrganizationProvided(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations/org-2/members":
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`[]`))
		case "/oapi/v1/platform/organizations":
			t.Fatal("should not resolve default organization when organizationId is provided")
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"member", "list",
		"--organization-id", "org-2",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if gotPath != "/oapi/v1/platform/organizations/org-2/members" {
		t.Fatalf("path = %q", gotPath)
	}
}

func TestMemberListOptionsParamsIncludesOrganizationAndPagination(t *testing.T) {
	params := (memberListOptions{
		OrganizationID: " org-1 ",
		Page:           2,
		PerPage:        50,
	}).params()
	wants := map[string]any{
		"organizationId": "org-1",
		"page":           2,
		"perPage":        50,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPrintMemberListShowsNoResultsWhenRowsEmpty(t *testing.T) {
	var out bytes.Buffer
	raw := "No results found."
	if err := printMemberList(&out, raw); err != nil {
		t.Fatalf("printMemberList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want \"No results found.\"", out.String())
	}
}

func TestMemberRowsFromJSONExtractsNestedValues(t *testing.T) {
	rows := memberRowsFromJSON(`{"data":[{"id":"member-1","accountId":123,"realName":"Alice","mail":"alice@example.com"},"skip"]}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := memberRow{MemberID: "member-1", UserID: "123", Name: "Alice", Email: "alice@example.com"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestMemberRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows := memberRowsFromJSON(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}
