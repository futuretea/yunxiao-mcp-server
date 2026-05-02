package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

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
