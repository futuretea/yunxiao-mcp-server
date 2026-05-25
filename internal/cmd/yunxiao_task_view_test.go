package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLITaskViewPrintsJSONWithDefaultSections(t *testing.T) {
	requests := map[string]int{}
	relationTypes := map[string]bool{}
	var commentsQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/workitems/wi-1":
			_, _ = w.Write([]byte(`{"id":"wi-1","subject":"Fix login"}`))
		case "/oapi/v1/projex/organizations/org-1/workitems/wi-1/activities":
			_, _ = w.Write([]byte(`[{"id":"act-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/workitems/wi-1/attachments":
			_, _ = w.Write([]byte(`[{"id":"att-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/workitems/wi-1/comments":
			commentsQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[{"id":"comment-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/workitems/wi-1/relationRecords":
			relationTypes[r.URL.Query().Get("relationType")] = true
			_, _ = w.Write([]byte(`[]`))
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
		"task", "view", "wi-1",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out.String())
	}
	workitem, ok := payload["workitem"].(map[string]any)
	if !ok || workitem["id"] != "wi-1" || workitem["subject"] != "Fix login" {
		t.Fatalf("workitem = %#v", payload["workitem"])
	}
	for _, section := range []string{"activities", "attachments", "comments", "relations_associated", "relations_sub"} {
		if _, ok := payload[section]; !ok {
			t.Fatalf("payload missing %q: %#v", section, payload)
		}
	}
	if !relationTypes["ASSOCIATED"] || !relationTypes["SUB"] {
		t.Fatalf("relationTypes = %#v", relationTypes)
	}
	for _, want := range []string{"page=1", "perPage=20"} {
		if !strings.Contains(commentsQuery, want) {
			t.Fatalf("comments query = %q, missing %q", commentsQuery, want)
		}
	}
	if requests["/oapi/v1/platform/organizations"] != 1 ||
		requests["/oapi/v1/projex/organizations/org-1/workitems/wi-1"] != 1 ||
		requests["/oapi/v1/projex/organizations/org-1/workitems/wi-1/relationRecords"] != 2 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLITaskDetailAliasUsesExplicitOrganizationAndSectionFlags(t *testing.T) {
	requests := map[string]int{}
	relationTypes := map[string]bool{}
	var commentsQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			t.Fatal("should not resolve default organization when organizationId is provided")
		case "/oapi/v1/projex/organizations/org-2/workitems/wi-1":
			_, _ = w.Write([]byte(`{"id":"wi-1"}`))
		case "/oapi/v1/projex/organizations/org-2/workitems/wi-1/comments":
			commentsQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[{"id":"comment-1"}]`))
		case "/oapi/v1/projex/organizations/org-2/workitems/wi-1/relationRecords":
			relationTypes[r.URL.Query().Get("relationType")] = true
			_, _ = w.Write([]byte(`[]`))
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
		"task", "detail", "wi-1",
		"--organization-id", "org-2",
		"--include-activities=false",
		"--include-attachments=false",
		"--include-comments=true",
		"--include-relations=true",
		"--relation-types", "PARENT,SUB",
		"--page", "2",
		"--per-page", "50",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out.String())
	}
	if _, ok := payload["comments"]; !ok {
		t.Fatalf("payload missing comments: %#v", payload)
	}
	for _, want := range []string{"page=2", "perPage=50"} {
		if !strings.Contains(commentsQuery, want) {
			t.Fatalf("comments query = %q, missing %q", commentsQuery, want)
		}
	}
	if !relationTypes["PARENT"] || !relationTypes["SUB"] {
		t.Fatalf("relationTypes = %#v", relationTypes)
	}
	if requests["/oapi/v1/projex/organizations/org-2/workitems/wi-1/relationRecords"] != 2 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLITaskViewRejectsMissingWorkitemIDBeforeNetwork(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		http.Error(w, "unexpected request", http.StatusInternalServerError)
	}))
	defer server.Close()

	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"task", "view",
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected missing argument error")
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0", requests)
	}
}

func TestYunxiaoCLITaskViewReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "projex", "task", "view", "wi-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "get_project_workitem_detail"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestTaskViewOptionsParamsIncludesFlags(t *testing.T) {
	params, err := (taskViewOptions{
		OrganizationID:     " org-1 ",
		IncludeActivities:  true,
		IncludeRelations:   false,
		RelationTypes:      " PARENT,SUB ",
		IncludeAttachments: false,
		IncludeComments:    true,
		Page:               2,
		PerPage:            50,
	}).params(" wi-1 ")
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId":     "org-1",
		"workitemId":         "wi-1",
		"includeActivities":  true,
		"includeRelations":   false,
		"relationTypes":      "PARENT,SUB",
		"includeAttachments": false,
		"includeComments":    true,
		"page":               2,
		"perPage":            50,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestTaskViewOptionsParamsRequiresWorkitemID(t *testing.T) {
	_, err := (taskViewOptions{}).params(" ")
	if err == nil {
		t.Fatal("params() expected workitem-id error")
	}
	if !strings.Contains(err.Error(), "workitem-id is required") {
		t.Fatalf("error = %v", err)
	}
}
