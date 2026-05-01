package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetProjectRiskDashboardBuildsRiskQueries(t *testing.T) {
	seenCategories := map[string]bool{}
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
		seenCategories[category] = true
		if body["spaceId"] != "project-1" || body["perPage"].(float64) != 2 {
			t.Fatalf("body = %#v", body)
		}
		conditions, _ := body["conditions"].(string)
		if !strings.Contains(conditions, `"fieldIdentifier":"status"`) {
			t.Fatalf("conditions = %q", conditions)
		}
		if category == "Risk,Bug" && !strings.Contains(conditions, `"fieldIdentifier":"finishTime"`) {
			t.Fatalf("overdue conditions = %q", conditions)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"workitem-1"}]`))
	})

	result, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"categories":     "Risk,Bug",
		"status":         "open,doing",
		"overdueBefore":  "2026-05-01",
		"sampleLimit":    float64(2),
	})
	if err != nil {
		t.Fatalf("handleGetProjectRiskDashboard() error = %v", err)
	}
	for _, category := range []string{"Risk", "Bug", "Risk,Bug"} {
		if !seenCategories[category] {
			t.Fatalf("missing category query %q in %#v", category, seenCategories)
		}
	}
	if !strings.Contains(result, `"overdue"`) || !strings.Contains(result, `"byCategory"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetProjectMemberTaskStatusUsesProvidedAssigneesAndGroups(t *testing.T) {
	requests := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
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
		if body["category"] != "Task" || body["spaceId"] != "project-1" || body["perPage"].(float64) != 1 {
			t.Fatalf("body = %#v", body)
		}
		conditions, _ := body["conditions"].(string)
		if !strings.Contains(conditions, `"fieldIdentifier":"assignedTo"`) {
			t.Fatalf("conditions = %q", conditions)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"workitem-1"}]`))
	})

	result, err := handleGetProjectMemberTaskStatus(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"assigneeIds":    "user-1,user-2",
		"categories":     "Task",
		"overdueBefore":  "2026-05-01",
		"statusGroups":   `{"todo":"todo-id","doing":"doing-id"}`,
		"sampleLimit":    float64(1),
	})
	if err != nil {
		t.Fatalf("handleGetProjectMemberTaskStatus() error = %v", err)
	}
	if requests != 8 {
		t.Fatalf("requests = %d, want 8", requests)
	}
	if !strings.Contains(result, `"user-1"`) || !strings.Contains(result, `"statusGroups"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetProjectMemberTaskStatusLoadsProjectMembers(t *testing.T) {
	requests := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/members" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			_, _ = w.Write([]byte(`[{"userId":"user-1","userName":"Ada"},{"userId":"user-2","userName":"Grace"}]`))
		case http.MethodPost:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			_, _ = w.Write([]byte(`[]`))
		default:
			t.Fatalf("method = %s", r.Method)
		}
	})

	result, err := handleGetProjectMemberTaskStatus(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"memberLimit":    float64(1),
		"sampleLimit":    float64(0),
	})
	if err != nil {
		t.Fatalf("handleGetProjectMemberTaskStatus() error = %v", err)
	}
	if requests != 3 {
		t.Fatalf("requests = %d, want 3", requests)
	}
	if !strings.Contains(result, `"Ada"`) || strings.Contains(result, `"Grace"`) {
		t.Fatalf("result = %q", result)
	}
}
