package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLISprintViewPrintsJSONWithDefaultOrganization(t *testing.T) {
	requests := map[string]int{}
	categories := map[string]bool{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/sprints/sprint-1":
			_, _ = w.Write([]byte(`{"id":"sprint-1","name":"Release 1"}`))
		case "/oapi/v1/projex/organizations/org-1/workitems:search":
			body, _ := io.ReadAll(r.Body)
			var payload map[string]any
			if err := json.Unmarshal(body, &payload); err != nil {
				t.Fatalf("decode body: %v\n%s", err, string(body))
			}
			category, _ := payload["category"].(string)
			categories[category] = true
			if payload["spaceId"] != "project-1" || payload["perPage"] != float64(5) {
				t.Fatalf("body = %#v", payload)
			}
			conditions, _ := payload["conditions"].(string)
			if !strings.Contains(conditions, `"fieldIdentifier":"sprint"`) || !strings.Contains(conditions, `"value":["sprint-1"]`) {
				t.Fatalf("conditions = %q", conditions)
			}
			_, _ = w.Write([]byte(`[{"id":"wi-1"}]`))
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
		"sprint", "view", "sprint-1",
		"--project-id", "project-1",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out.String())
	}
	sprint, ok := payload["sprint"].(map[string]any)
	if !ok || sprint["id"] != "sprint-1" {
		t.Fatalf("sprint = %#v", payload["sprint"])
	}
	if !categories["Task"] || !categories["Bug"] {
		t.Fatalf("categories = %#v", categories)
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/projex/organizations/org-1/workitems:search"] != 2 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLISprintOverviewAliasUsesExplicitOrganizationAndFilters(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			t.Fatal("should not resolve default organization when organizationId is provided")
		case "/oapi/v1/projex/organizations/org-2/projects/project-1/sprints/sprint-1":
			_, _ = w.Write([]byte(`{"id":"sprint-1"}`))
		case "/oapi/v1/projex/organizations/org-2/workitems:search":
			body, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(body, &gotBody); err != nil {
				t.Fatalf("decode body: %v\n%s", err, string(body))
			}
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
		"sprint", "overview", "sprint-1",
		"--organization-id", "org-2",
		"--project-id", "project-1",
		"--categories", "Task",
		"--subject", "login",
		"--status", "open",
		"--assigned-to", "user-1",
		"--creator", "user-2",
		"--sample-limit", "7",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if err := json.Unmarshal(out.Bytes(), &map[string]any{}); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out.String())
	}
	if gotBody["category"] != "Task" || gotBody["spaceId"] != "project-1" || gotBody["perPage"] != float64(7) {
		t.Fatalf("body = %#v", gotBody)
	}
	conditions, _ := gotBody["conditions"].(string)
	for _, want := range []string{`"fieldIdentifier":"subject"`, `"fieldIdentifier":"status"`, `"fieldIdentifier":"assignedTo"`, `"fieldIdentifier":"creator"`, `"fieldIdentifier":"sprint"`} {
		if !strings.Contains(conditions, want) {
			t.Fatalf("conditions = %q, missing %q", conditions, want)
		}
	}
}

func TestYunxiaoCLISprintViewPassesExplicitZeroSampleLimit(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/sprints/sprint-1":
			_, _ = w.Write([]byte(`{"id":"sprint-1"}`))
		case "/oapi/v1/projex/organizations/org-1/workitems:search":
			body, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(body, &gotBody); err != nil {
				t.Fatalf("decode body: %v\n%s", err, string(body))
			}
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
		"sprint", "view", "sprint-1",
		"--organization-id", "org-1",
		"--project-id", "project-1",
		"--categories", "Task",
		"--sample-limit", "0",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if err := json.Unmarshal(out.Bytes(), &map[string]any{}); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out.String())
	}
	if gotBody["perPage"] != float64(0) {
		t.Fatalf("perPage = %#v, want 0", gotBody["perPage"])
	}
}

func TestYunxiaoCLISprintViewRequiresProjectID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"sprint", "view", "sprint-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected project-id error")
	}
	if !strings.Contains(err.Error(), "project-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLISprintViewRejectsMissingSprintIDBeforeNetwork(t *testing.T) {
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
		"sprint", "view",
		"--project-id", "project-1",
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected missing argument error")
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0", requests)
	}
}

func TestYunxiaoCLISprintViewReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "projex", "sprint", "view", "sprint-1", "--project-id", "project-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "get_sprint_overview"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestSprintViewOptionsParamsIncludesFilters(t *testing.T) {
	params, err := (sprintViewOptions{
		OrganizationID: " org-1 ",
		ProjectID:      " project-1 ",
		Categories:     " Task,Bug ",
		Subject:        " login ",
		Status:         " open ",
		AssignedTo:     " user-1 ",
		Creator:        " user-2 ",
		SampleLimit:    7,
		SampleLimitSet: true,
	}).params(" sprint-1 ")
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"sprintId":       "sprint-1",
		"categories":     "Task,Bug",
		"subject":        "login",
		"status":         "open",
		"assignedTo":     "user-1",
		"creator":        "user-2",
		"sampleLimit":    7,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestSprintViewOptionsParamsRequiresIDs(t *testing.T) {
	if _, err := (sprintViewOptions{}).params("sprint-1"); err == nil {
		t.Fatal("params() expected project-id error")
	}
	if _, err := (sprintViewOptions{ProjectID: "project-1"}).params(" "); err == nil {
		t.Fatal("params() expected sprint-id error")
	}
}
