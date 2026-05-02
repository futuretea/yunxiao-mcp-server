package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetProjectOverviewBuildsCommonRequests(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seen[r.URL.Path] = true
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}

		switch r.URL.Path {
		case "/oapi/v1/projex/organizations/org-1/projects/project-1":
			_, _ = w.Write([]byte(`{"id":"project-1"}`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/members":
			_, _ = w.Write([]byte(`[{"id":"user-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/sprints",
			"/oapi/v1/projex/organizations/org-1/projects/project-1/milestones",
			"/oapi/v1/projex/organizations/org-1/projects/project-1/versions":
			if r.URL.Query().Get("status") != "TODO,DOING" ||
				r.URL.Query().Get("page") != "2" ||
				r.URL.Query().Get("perPage") != "3" {
				t.Fatalf("query = %q", r.URL.RawQuery)
			}
			w.Header().Set("x-total", "1")
			_, _ = w.Write([]byte(`[{"id":"item-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/labels":
			if r.URL.Query().Get("status") != "" ||
				r.URL.Query().Get("page") != "2" ||
				r.URL.Query().Get("perPage") != "3" {
				t.Fatalf("query = %q", r.URL.RawQuery)
			}
			w.Header().Set("x-total", "1")
			_, _ = w.Write([]byte(`[{"id":"label-1"}]`))
		default:
			t.Fatalf("unexpected path = %q", r.URL.Path)
		}
	})

	result, err := handleGetProjectOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"page":           float64(2),
		"perPage":        float64(3),
	})
	if err != nil {
		t.Fatalf("handleGetProjectOverview() error = %v", err)
	}
	for _, path := range []string{
		"/oapi/v1/projex/organizations/org-1/projects/project-1",
		"/oapi/v1/projex/organizations/org-1/projects/project-1/members",
		"/oapi/v1/projex/organizations/org-1/projects/project-1/sprints",
		"/oapi/v1/projex/organizations/org-1/projects/project-1/milestones",
		"/oapi/v1/projex/organizations/org-1/projects/project-1/versions",
		"/oapi/v1/projex/organizations/org-1/projects/project-1/labels",
	} {
		if !seen[path] {
			t.Fatalf("missing request to %s", path)
		}
	}
	if !strings.Contains(result, `"sprints"`) || !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetProjectOverviewHonorsIncludeFlags(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1" {
			t.Fatalf("unexpected path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"project-1"}`))
	})

	result, err := handleGetProjectOverview(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"id":                "project-1",
		"includeMembers":    false,
		"includeSprints":    false,
		"includeMilestones": false,
		"includeVersions":   false,
		"includeLabels":     false,
	})
	if err != nil {
		t.Fatalf("handleGetProjectOverview() error = %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	for _, absent := range []string{"members", "sprints", "milestones", "versions", "labels"} {
		if _, ok := payload[absent]; ok {
			t.Fatalf("section %q should be absent in %#v", absent, payload)
		}
	}
}

func TestHandleGetProjectWorkitemSummarySearchesCategories(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		category, _ := body["category"].(string)
		if category != "Task" && category != "Bug" {
			t.Fatalf("category = %#v", body["category"])
		}
		seen[category] = true
		if body["spaceId"] != "project-1" || body["page"].(float64) != 1 || body["perPage"].(float64) != 3 {
			t.Fatalf("body = %#v", body)
		}
		conditions, _ := body["conditions"].(string)
		if !strings.Contains(conditions, `"fieldIdentifier":"status"`) ||
			!strings.Contains(conditions, `"fieldIdentifier":"assignedTo"`) {
			t.Fatalf("conditions = %q", conditions)
		}
		w.Header().Set("x-total", "2")
		_, _ = w.Write([]byte(`[{"id":"workitem-1"}]`))
	})

	result, err := handleGetProjectWorkitemSummary(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"categories":     "Task,Bug",
		"status":         "open,doing",
		"assignedTo":     "user-1",
		"sampleLimit":    float64(3),
	})
	if err != nil {
		t.Fatalf("handleGetProjectWorkitemSummary() error = %v", err)
	}
	if !seen["Task"] || !seen["Bug"] {
		t.Fatalf("seen = %#v", seen)
	}
	if !strings.Contains(result, `"Task"`) || !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetProjectWorkitemSummaryRejectsEmptyCategories(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without categories")
	})

	if _, err := handleGetProjectWorkitemSummary(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"categories":     ", ,",
	}); err == nil {
		t.Fatal("handleGetProjectWorkitemSummary() expected missing categories error")
	}
}

func TestHandleGetProjectWorkitemContextBuildsMetadataRequests(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		requestPath := r.URL.EscapedPath()
		seen[requestPath] = true

		switch requestPath {
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/workitemTypes":
			if r.URL.Query().Get("category") != "Task" {
				t.Fatalf("query = %q", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`[{"id":"type-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/members":
			_, _ = w.Write([]byte(`[{"id":"user-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/labels":
			if r.URL.Query().Get("page") != "2" || r.URL.Query().Get("perPage") != "10" {
				t.Fatalf("query = %q", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`[{"id":"label-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/workitemTypes/type%2F1/fields":
			_, _ = w.Write([]byte(`[{"fieldIdentifier":"status"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/workitemTypes/type%2F1/workflows":
			_, _ = w.Write([]byte(`{"id":"workflow-1"}`))
		default:
			t.Fatalf("unexpected path = %q", requestPath)
		}
	})

	result, err := handleGetProjectWorkitemContext(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"category":       "Task",
		"workItemTypeId": "type/1",
		"page":           float64(2),
		"perPage":        float64(10),
	})
	if err != nil {
		t.Fatalf("handleGetProjectWorkitemContext() error = %v", err)
	}
	for _, path := range []string{
		"/oapi/v1/projex/organizations/org-1/projects/project-1/workitemTypes",
		"/oapi/v1/projex/organizations/org-1/projects/project-1/members",
		"/oapi/v1/projex/organizations/org-1/projects/project-1/labels",
		"/oapi/v1/projex/organizations/org-1/projects/project-1/workitemTypes/type%2F1/fields",
		"/oapi/v1/projex/organizations/org-1/projects/project-1/workitemTypes/type%2F1/workflows",
	} {
		if !seen[path] {
			t.Fatalf("missing request to %s", path)
		}
	}
	if !strings.Contains(result, `"workItemTypes"`) || !strings.Contains(result, `"workflow"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetProjectWorkitemContextRequiresCategory(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without category")
	})

	if _, err := handleGetProjectWorkitemContext(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
	}); err == nil {
		t.Fatal("handleGetProjectWorkitemContext() expected missing category error")
	}
}

func TestHandleGetSprintOverviewBuildsSprintAndSearchRequests(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/sprints/sp-1" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			seen["sprint"] = true
			_, _ = w.Write([]byte(`{"id":"sp-1","name":"Sprint 1"}`))
		case http.MethodPost:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			category, _ := body["category"].(string)
			seen["cat:"+category] = true
			if body["spaceId"] != "project-1" {
				t.Fatalf("body = %#v", body)
			}
			conditions, _ := body["conditions"].(string)
			if !strings.Contains(conditions, `"fieldIdentifier":"sprint"`) {
				t.Fatalf("conditions = %q", conditions)
			}
			w.Header().Set("x-total", "2")
			_, _ = w.Write([]byte(`[{"id":"wi-1"}]`))
		default:
			t.Fatalf("method = %s", r.Method)
		}
	})

	result, err := handleGetSprintOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"sprintId":       "sp-1",
		"categories":     "Task,Bug",
		"sampleLimit":    float64(3),
	})
	if err != nil {
		t.Fatalf("handleGetSprintOverview() error = %v", err)
	}
	for _, key := range []string{"sprint", "cat:Task", "cat:Bug"} {
		if !seen[key] {
			t.Fatalf("missing request for %s", key)
		}
	}
	if !strings.Contains(result, `"sprint"`) || !strings.Contains(result, `"categories"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetSprintOverviewRejectsMissingSprintId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without sprintId")
	})

	if _, err := handleGetSprintOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
	}); err == nil {
		t.Fatal("handleGetSprintOverview() expected missing sprintId error")
	}
}

func TestHandleGetSprintOverviewRejectsEmptyCategories(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"id":"sp-1"}`))
	})

	if _, err := handleGetSprintOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"sprintId":       "sp-1",
		"categories":     ", ,",
	}); err == nil {
		t.Fatal("handleGetSprintOverview() expected missing categories error")
	}
}

func TestHandleGetMyProjectWorkitemsBuildsAssignedSearch(t *testing.T) {
	seen := map[string]int{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		category, _ := body["category"].(string)
		seen["cat:"+category]++
		conditions, _ := body["conditions"].(string)
		if !strings.Contains(conditions, `"fieldIdentifier":"assignedTo"`) {
			t.Fatalf("conditions = %q", conditions)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"wi-1"}]`))
	})

	result, err := handleGetMyProjectWorkitems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"userId":         "user-1",
		"relation":       "assigned",
		"categories":     "Task,Bug",
		"sampleLimit":    float64(3),
	})
	if err != nil {
		t.Fatalf("handleGetMyProjectWorkitems() error = %v", err)
	}
	if seen["cat:Task"] != 1 || seen["cat:Bug"] != 1 {
		t.Fatalf("seen = %#v", seen)
	}
	if !strings.Contains(result, `"assigned"`) || !strings.Contains(result, `"userId"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetMyProjectWorkitemsBuildsCreatedSearch(t *testing.T) {
	seen := map[string]int{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		conditions, _ := body["conditions"].(string)
		if !strings.Contains(conditions, `"fieldIdentifier":"creator"`) {
			t.Fatalf("conditions = %q", conditions)
		}
		seen["search"]++
		_, _ = w.Write([]byte(`[]`))
	})

	_, err := handleGetMyProjectWorkitems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"userId":         "user-1",
		"relation":       "created",
		"categories":     "Task",
	})
	if err != nil {
		t.Fatalf("handleGetMyProjectWorkitems() error = %v", err)
	}
	if seen["search"] != 1 {
		t.Fatalf("seen = %#v", seen)
	}
}

func TestHandleGetMyProjectWorkitemsRejectsInvalidRelation(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request with invalid relation")
	})

	if _, err := handleGetMyProjectWorkitems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"userId":         "user-1",
		"relation":       "invalid",
	}); err == nil {
		t.Fatal("handleGetMyProjectWorkitems() expected invalid relation error")
	}
}

func TestHandleGetMyProjectWorkitemsRejectsMissingUserId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without userId")
	})

	if _, err := handleGetMyProjectWorkitems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
	}); err == nil {
		t.Fatal("handleGetMyProjectWorkitems() expected missing userId error")
	}
}

func TestHandleGetProjectWorkitemBoardGroupsByStatus(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["category"] != "Task" || body["spaceId"] != "project-1" {
			t.Fatalf("body = %#v", body)
		}
		w.Header().Set("x-total", "3")
		_, _ = w.Write([]byte(`[
			{"id":"wi-1","status":{"name":"Doing"}},
			{"id":"wi-2","status":{"name":"Done"}},
			{"id":"wi-3","status":{"name":"Doing"}}
		]`))
	})

	result, err := handleGetProjectWorkitemBoard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"category":       "Task",
		"sampleLimit":    float64(5),
	})
	if err != nil {
		t.Fatalf("handleGetProjectWorkitemBoard() error = %v", err)
	}
	if !strings.Contains(result, `"Doing"`) || !strings.Contains(result, `"Done"`) {
		t.Fatalf("result = %q", result)
	}
	if !strings.Contains(result, `"total"`) {
		t.Fatalf("result missing total: %q", result)
	}
}

func TestHandleGetProjectWorkitemBoardRequiresCategory(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without category")
	})

	if _, err := handleGetProjectWorkitemBoard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
	}); err == nil {
		t.Fatal("handleGetProjectWorkitemBoard() expected missing category error")
	}
}

func TestHandleGetProjectWorkitemBoardWithSprintFilter(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		conditions, _ := body["conditions"].(string)
		if !strings.Contains(conditions, `"fieldIdentifier":"sprint"`) {
			t.Fatalf("conditions = %q", conditions)
		}
		_, _ = w.Write([]byte(`[]`))
	})

	_, err := handleGetProjectWorkitemBoard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"category":       "Bug",
		"sprint":         "sp-1",
	})
	if err != nil {
		t.Fatalf("handleGetProjectWorkitemBoard() error = %v", err)
	}
}
