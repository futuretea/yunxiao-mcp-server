package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

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
		"projectId":      "project-1",
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
		"projectId":      "project-1",
		"categories":     ", ,",
	}); err == nil {
		t.Fatal("handleGetProjectWorkitemSummary() expected missing categories error")
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
		"projectId":      "project-1",
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
		"projectId":      "project-1",
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
		"projectId":      "project-1",
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
		"projectId":      "project-1",
	}); err == nil {
		t.Fatal("handleGetMyProjectWorkitems() expected missing userId error")
	}
}

func TestHandleGetMyProjectWorkitemsRejectsEmptyCategories(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without categories")
	})

	if _, err := handleGetMyProjectWorkitems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"userId":         "user-1",
		"categories":     ", ,",
	}); err == nil {
		t.Fatal("handleGetMyProjectWorkitems() expected missing categories error")
	}
}

func TestProjexEnhancedHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetProjectOverview(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetProjectOverview(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "projectId": "project-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetProjectWorkitemSummary(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetProjectWorkitemSummary(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "projectId": "project-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetProjectWorkitemContext(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetProjectWorkitemContext(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "projectId": "project-1", "category": "Task"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetMyProjectWorkitems(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetMyProjectWorkitems(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "projectId": "project-1", "userId": "user-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetSprintOverview(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetSprintOverview(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "projectId": "project-1", "sprintId": "sp-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetProjectWorkitemBoard(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetProjectWorkitemBoard(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "projectId": "project-1", "category": "Task"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetProjectWorkitemDetail(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetProjectWorkitemDetail(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "workitemId": "wi-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}

func TestHandleGetProjectOverviewAPIErrors(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/projex/organizations/org-1/projects/project-1" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})
	if _, err := handleGetProjectOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
	}); err == nil {
		t.Fatal("expected project fetch error")
	}
}

func TestHandleGetProjectWorkitemDetailAPIErrors(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/projex/organizations/org-1/workitems/wi-1" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})
	if _, err := handleGetProjectWorkitemDetail(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
	}); err == nil {
		t.Fatal("expected workitem fetch error")
	}
}

func TestHandleGetProjectWorkitemDetailReturnsSectionError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/projex/organizations/org-1/workitems/wi-1" {
			_, _ = w.Write([]byte(`{"id":"wi-1"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := handleGetProjectWorkitemDetail(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
	}); err == nil {
		t.Fatal("expected section error")
	}
}

func TestHandleGetProjectWorkitemSummaryAPIErrors(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := handleGetProjectWorkitemSummary(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	}); err == nil {
		t.Fatal("expected search error")
	}
}

func TestHandleGetProjectWorkitemContextAPIErrors(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := handleGetProjectWorkitemContext(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"category":       "Task",
	}); err == nil {
		t.Fatal("expected context fetch error")
	}
}

func TestHandleGetMyProjectWorkitemsAPIErrors(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := handleGetMyProjectWorkitems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"userId":         "user-1",
		"categories":     "Task",
	}); err == nil {
		t.Fatal("expected search error")
	}
}

func TestHandleGetSprintOverviewAPIErrors(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := handleGetSprintOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"sprintId":       "sp-1",
		"categories":     "Task",
	}); err == nil {
		t.Fatal("expected sprint fetch error")
	}
}

func TestHandleGetProjectWorkitemBoardAPIErrors(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := handleGetProjectWorkitemBoard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"category":       "Task",
	}); err == nil {
		t.Fatal("expected search error")
	}
}
