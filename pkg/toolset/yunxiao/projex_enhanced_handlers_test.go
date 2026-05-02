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

func TestHandleGetSprintOverviewWithAssigneeAndSubject(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, _ = w.Write([]byte(`{"id":"sp-1"}`))
			return
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		conditions, _ := body["conditions"].(string)
		if !strings.Contains(conditions, `"fieldIdentifier":"assignedTo"`) {
			t.Fatalf("missing assignedTo in conditions = %q", conditions)
		}
		if !strings.Contains(conditions, `"fieldIdentifier":"subject"`) {
			t.Fatalf("missing subject in conditions = %q", conditions)
		}
		if !strings.Contains(conditions, `"fieldIdentifier":"sprint"`) {
			t.Fatalf("missing sprint in conditions = %q", conditions)
		}
		_, _ = w.Write([]byte(`[]`))
	})

	_, err := handleGetSprintOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"sprintId":       "sp-1",
		"categories":     "Task",
		"assignedTo":     "user-1",
		"subject":        "auth",
	})
	if err != nil {
		t.Fatalf("handleGetSprintOverview() error = %v", err)
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
	if !strings.Contains(result, `"columnCounts"`) {
		t.Fatalf("result missing columnCounts: %q", result)
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

func TestHandleGetProjectWorkitemBoardWithAssigneeAndSubject(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		conditions, _ := body["conditions"].(string)
		if !strings.Contains(conditions, `"fieldIdentifier":"assignedTo"`) {
			t.Fatalf("missing assignedTo in conditions = %q", conditions)
		}
		if !strings.Contains(conditions, `"fieldIdentifier":"subject"`) {
			t.Fatalf("missing subject in conditions = %q", conditions)
		}
		_, _ = w.Write([]byte(`[]`))
	})

	_, err := handleGetProjectWorkitemBoard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"category":       "Task",
		"assignedTo":     "user-1",
		"subject":        "login",
	})
	if err != nil {
		t.Fatalf("handleGetProjectWorkitemBoard() error = %v", err)
	}
}

func TestMergeConditions(t *testing.T) {
	tests := []struct {
		name     string
		existing string
		extra    string
		want     string
	}{
		{"both empty", "", "", ""},
		{"existing empty", "", `[{"field":"a"}]`, `[{"field":"a"}]`},
		{"extra empty", `[{"field":"a"}]`, "", `[{"field":"a"}]`},
		{"both arrays", `[{"field":"a"}]`, `[{"field":"b"}]`, `[{"field":"a"},{"field":"b"}]`},
		{"both objects", `{"conditionGroups":[[{"field":"a"}]]}`, `{"conditionGroups":[[{"field":"b"}]]}`, `{"conditionGroups":[[{"field":"a"},{"field":"b"}]]}`},
		{"existing object extra array", `{"conditionGroups":[[{"field":"a"}]]}`, `[{"field":"b"}]`, `{"conditionGroups":[[{"field":"a"}]]}`},
		{"existing array extra object", `[{"field":"a"}]`, `{"conditionGroups":[[{"field":"b"}]]}`, `[{"field":"a"}]`},
		{"existing invalid json", `not-json`, `[{"field":"b"}]`, `not-json`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeConditions(tt.existing, tt.extra)
			if got != tt.want {
				t.Fatalf("mergeConditions(%q, %q) = %q, want %q", tt.existing, tt.extra, got, tt.want)
			}
		})
	}
}

func TestHandleGetProjectWorkitemDetailBuildsCommonRequests(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seen[r.URL.Path] = true
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}

		switch r.URL.Path {
		case "/oapi/v1/projex/organizations/org-1/workitems/wi-1":
			_, _ = w.Write([]byte(`{"id":"wi-1"}`))
		case "/oapi/v1/projex/organizations/org-1/workitems/wi-1/activities":
			_, _ = w.Write([]byte(`[{"id":"act-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/workitems/wi-1/attachments":
			_, _ = w.Write([]byte(`[{"id":"att-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/workitems/wi-1/comments":
			if r.URL.Query().Get("page") != "1" || r.URL.Query().Get("perPage") != "20" {
				t.Fatalf("query = %q", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`[{"id":"cmt-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/workitems/wi-1/relationRecords":
			rt := r.URL.Query().Get("relationType")
			if rt != "ASSOCIATED" && rt != "SUB" {
				t.Fatalf("relationType = %q", rt)
			}
			_, _ = w.Write([]byte(`[]`))
		default:
			t.Fatalf("unexpected path = %q", r.URL.Path)
		}
	})

	result, err := handleGetProjectWorkitemDetail(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
	})
	if err != nil {
		t.Fatalf("handleGetProjectWorkitemDetail() error = %v", err)
	}

	if !seen["/oapi/v1/projex/organizations/org-1/workitems/wi-1"] {
		t.Fatal("expected workitem request")
	}
	if !seen["/oapi/v1/projex/organizations/org-1/workitems/wi-1/activities"] {
		t.Fatal("expected activities request")
	}
	if !seen["/oapi/v1/projex/organizations/org-1/workitems/wi-1/attachments"] {
		t.Fatal("expected attachments request")
	}
	if !seen["/oapi/v1/projex/organizations/org-1/workitems/wi-1/comments"] {
		t.Fatal("expected comments request")
	}
	if !strings.Contains(result, `"workitem"`) {
		t.Fatalf("result missing workitem: %q", result)
	}
}

func TestHandleGetProjectWorkitemDetailSkipsDisabledSections(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seen[r.URL.Path] = true
		if r.URL.Path == "/oapi/v1/projex/organizations/org-1/workitems/wi-1" {
			_, _ = w.Write([]byte(`{"id":"wi-1"}`))
			return
		}
		_, _ = w.Write([]byte(`[]`))
	})

	_, err := handleGetProjectWorkitemDetail(context.Background(), client, map[string]any{
		"organizationId":     "org-1",
		"workitemId":         "wi-1",
		"includeActivities":  false,
		"includeRelations":   false,
		"includeAttachments": false,
		"includeComments":    false,
	})
	if err != nil {
		t.Fatalf("handleGetProjectWorkitemDetail() error = %v", err)
	}

	if len(seen) != 1 {
		t.Fatalf("expected 1 request, got %d: %v", len(seen), seen)
	}
}

func TestHandleGetProjectWorkitemDetailRequiresWorkitemId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without workitemId")
	})

	if _, err := handleGetProjectWorkitemDetail(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	}); err == nil {
		t.Fatal("handleGetProjectWorkitemDetail() expected missing workitemId error")
	}
}

func TestOptionalIntDefault(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want int
	}{
		{"float64", float64(42), 42},
		{"int", int(7), 7},
		{"int64", int64(99), 99},
		{"string valid", "  123  ", 123},
		{"string invalid", "abc", 10},
		{"nil", nil, 10},
		{"bool", true, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := optionalIntDefault(map[string]any{"k": tt.val}, "k", 10)
			if got != tt.want {
				t.Fatalf("optionalIntDefault(%v) = %d, want %d", tt.val, got, tt.want)
			}
		})
	}
}

func TestNormalizedSampleLimit(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want int
	}{
		{"default", nil, 5},
		{"within range", float64(50), 50},
		{"negative", float64(-3), 0},
		{"overflow", float64(500), 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizedSampleLimit(map[string]any{"sampleLimit": tt.val})
			if got != tt.want {
				t.Fatalf("normalizedSampleLimit(%v) = %d, want %d", tt.val, got, tt.want)
			}
		})
	}
}

func TestExtractWorkitemData(t *testing.T) {
	tests := []struct {
		name      string
		payload   any
		wantLen   int
		wantTotal int
		wantErr   bool
	}{
		{"array", []any{map[string]any{"id": "1"}}, 1, 1, false},
		{"map with pagination", map[string]any{"data": []any{map[string]any{"id": "1"}}, "pagination": map[string]any{"total": float64(10)}}, 1, 10, false},
		{"map without pagination", map[string]any{"data": []any{map[string]any{"id": "1"}}}, 1, 0, false},
		{"invalid type", "string", 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, total, err := extractWorkitemData(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Fatalf("extractWorkitemData() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if len(data) != tt.wantLen {
					t.Fatalf("len(data) = %d, want %d", len(data), tt.wantLen)
				}
				if total != tt.wantTotal {
					t.Fatalf("total = %d, want %d", total, tt.wantTotal)
				}
			}
		})
	}
}
