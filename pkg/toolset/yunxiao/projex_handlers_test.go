package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestHandleSearchProjectsBuildsBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["orderBy"] != "name" || body["sort"] != "asc" || body["page"].(float64) != 2 {
			t.Fatalf("body = %#v", body)
		}
		conditions, _ := body["conditions"].(string)
		if !strings.Contains(conditions, `"fieldIdentifier":"name"`) ||
			!strings.Contains(conditions, `"value":["demo"]`) {
			t.Fatalf("conditions = %q", conditions)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"project-1"}]`))
	})

	result, err := handleSearchProjects(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"name":           "demo",
		"orderBy":        "name",
		"sort":           "asc",
		"page":           float64(2),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleSearchProjects() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleSearchProjectsAdvancedConditionsOverrideSimpleFilters(t *testing.T) {
	const advanced = `{"conditionGroups":[[{"fieldIdentifier":"custom"}]]}`
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["conditions"] != advanced {
			t.Fatalf("conditions = %#v", body["conditions"])
		}
		_, _ = w.Write([]byte(`[]`))
	})

	if _, err := handleSearchProjects(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"name":           "demo",
		"conditions":     advanced,
	}); err != nil {
		t.Fatalf("handleSearchProjects() error = %v", err)
	}
}

func TestHandleGetProjectBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"project-1"}`))
	})

	if _, err := handleGetProject(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
	}); err != nil {
		t.Fatalf("handleGetProject() error = %v", err)
	}
}

func TestHandleListSprintsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/sprints" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("status") != "TODO,DOING" ||
			r.URL.Query().Get("name") != "release" ||
			r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"sprint-1"}]`))
	})

	result, err := handleListSprints(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"status":         "TODO,DOING",
		"name":           "release",
		"page":           float64(2),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleListSprints() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetSprintBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/sprints/sprint-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"sprint-1"}`))
	})

	if _, err := handleGetSprint(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"id":             "sprint-1",
	}); err != nil {
		t.Fatalf("handleGetSprint() error = %v", err)
	}
}

func TestHandleSearchWorkitemsBuildsBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["category"] != "Task" || body["spaceId"] != "project-1" || body["sort"] != "desc" {
			t.Fatalf("body = %#v", body)
		}
		conditions, _ := body["conditions"].(string)
		if !strings.Contains(conditions, `"fieldIdentifier":"subject"`) ||
			!strings.Contains(conditions, `"fieldIdentifier":"assignedTo"`) {
			t.Fatalf("conditions = %q", conditions)
		}
		var parsedConditions struct {
			ConditionGroups [][]map[string]any `json:"conditionGroups"`
		}
		if err := json.Unmarshal([]byte(conditions), &parsedConditions); err != nil {
			t.Fatalf("unmarshal conditions: %v", err)
		}
		assertConditionFormat(t, parsedConditions.ConditionGroups, "tag", "multiList")
		assertConditionFormat(t, parsedConditions.ConditionGroups, "priority", "list")
		assertConditionFormat(t, parsedConditions.ConditionGroups, "finishTime", "input")
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"workitem-1"}]`))
	})

	if _, err := handleSearchWorkitems(context.Background(), client, map[string]any{
		"organizationId":   "org-1",
		"category":         "Task",
		"projectId":        "project-1",
		"subject":          "demo",
		"assignedTo":       "user-1,user-2",
		"tag":              "tag-1,tag-2",
		"priority":         "high",
		"finishTimeBefore": "2026-05-01",
		"sort":             "desc",
	}); err != nil {
		t.Fatalf("handleSearchWorkitems() error = %v", err)
	}
}

func assertConditionFormat(t *testing.T, conditionGroups [][]map[string]any, fieldIdentifier, format string) {
	t.Helper()

	for _, group := range conditionGroups {
		for _, condition := range group {
			if condition["fieldIdentifier"] == fieldIdentifier {
				if condition["format"] != format {
					t.Fatalf("%s format = %v, want %s", fieldIdentifier, condition["format"], format)
				}
				return
			}
		}
	}
	t.Fatalf("condition %s not found in %#v", fieldIdentifier, conditionGroups)
}

func TestHandleGetWorkitemBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems/workitem-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"workitem-1"}`))
	})

	if _, err := handleGetWorkitem(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "workitem-1",
	}); err != nil {
		t.Fatalf("handleGetWorkitem() error = %v", err)
	}
}

func TestHandleListWorkItemCommentsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems/workitem-1/comments" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "2" || r.URL.Query().Get("perPage") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"comment-1"}]`))
	})

	result, err := handleListWorkItemComments(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workItemId":     "workitem-1",
		"page":           float64(2),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleListWorkItemComments() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestSearchWorkitemsRequiresCategoryAndSpace(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without required params")
	})

	if _, err := handleSearchWorkitems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"spaceId":        "project-1",
	}); err == nil {
		t.Fatal("handleSearchWorkitems() expected missing category error")
	}
	if _, err := handleSearchWorkitems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"category":       "Task",
	}); err == nil {
		t.Fatal("handleSearchWorkitems() expected missing projectId error")
	}
}

func TestProjexHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleSearchProjects(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleSearchProjects(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetProject(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetProject(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing id error")
	}
	if _, err := handleGetProject(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "id": "p-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListSprints(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListSprints(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "id": "p-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetSprint(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetSprint(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing projectId error")
	}
	if _, err := handleGetSprint(context.Background(), client, map[string]any{"organizationId": "org-1", "projectId": "p-1"}); err == nil {
		t.Fatal("expected missing id error")
	}
	if _, err := handleGetSprint(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "projectId": "p-1", "id": "s-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleSearchWorkitems(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleSearchWorkitems(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "category": "Task", "projectId": "p-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetWorkitem(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetWorkitem(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing id error")
	}
	if _, err := handleGetWorkitem(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "id": "wi-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListWorkItemComments(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListWorkItemComments(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "workItemId": "wi-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}
