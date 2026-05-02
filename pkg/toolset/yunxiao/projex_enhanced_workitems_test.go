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
