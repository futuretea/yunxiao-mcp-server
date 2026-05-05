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
		"projectId":      "project-1",
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
		"projectId":      "project-1",
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
		"projectId":      "project-1",
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

func TestHandleGetProjectMemberTaskStatusRejectsEmptyUserIdsFromMembers(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, _ = w.Write([]byte(`[{"name":"no-id"}]`))
			return
		}
		_, _ = w.Write([]byte(`[]`))
	})

	if _, err := handleGetProjectMemberTaskStatus(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	}); err == nil {
		t.Fatal("handleGetProjectMemberTaskStatus() expected empty assigneeIds error")
	}
}

func TestHandleGetProjectRiskDashboardRejectsEmptyCategories(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without categories")
	})

	if _, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     ", ,",
	}); err == nil {
		t.Fatal("handleGetProjectRiskDashboard() expected missing categories error")
	}
}

func TestHandleGetProjectMemberTaskStatusReturnsMembersError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(`[]`))
	})

	if _, err := handleGetProjectMemberTaskStatus(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	}); err == nil {
		t.Fatal("handleGetProjectMemberTaskStatus() expected members error")
	}
}

func TestHandleGetProjectMemberTaskStatusReturnsOverdueSearchError(t *testing.T) {
	requestCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			_, _ = w.Write([]byte(`[{"id":"wi-1"}]`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	})

	if _, err := handleGetProjectMemberTaskStatus(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"assigneeIds":    "user-1",
		"categories":     "Task",
		"overdueBefore":  "2026-05-01",
	}); err == nil {
		t.Fatal("handleGetProjectMemberTaskStatus() expected overdue search error")
	}
}

func TestHandleGetProjectRiskDashboardReturnsOverdueSearchError(t *testing.T) {
	requestCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(`[]`))
	})

	if _, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Risk",
		"overdueBefore":  "2026-05-01",
	}); err == nil {
		t.Fatal("handleGetProjectRiskDashboard() expected overdue search error")
	}
}

func TestHandleGetProjectMemberTaskStatusRejectsInvalidStatusGroups(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"userId":"user-1"}]`))
	})

	if _, err := handleGetProjectMemberTaskStatus(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"assigneeIds":    "user-1",
		"statusGroups":   "not-json",
	}); err == nil {
		t.Fatal("handleGetProjectMemberTaskStatus() expected invalid statusGroups error")
	}
}

func TestProjexInsightsHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetProjectRiskDashboard(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "projectId": "p-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetProjectMemberTaskStatus(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetProjectMemberTaskStatus(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "projectId": "p-1", "assigneeIds": "u-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}

func TestParseStatusGroups(t *testing.T) {
	if got, err := parseStatusGroups(map[string]any{}); err != nil || got != nil {
		t.Fatalf("parseStatusGroups(empty) = %v, %v", got, err)
	}
	if _, err := parseStatusGroups(map[string]any{"statusGroups": "not-json"}); err == nil {
		t.Fatal("expected invalid statusGroups error")
	}
	got, err := parseStatusGroups(map[string]any{"statusGroups": `{"todo":"todo-id","doing":"doing-id"}`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["todo"] != "todo-id" || got["doing"] != "doing-id" {
		t.Fatalf("unexpected groups: %v", got)
	}
}

func TestCopyParams(t *testing.T) {
	original := map[string]any{"a": 1, "b": "two"}
	copied := copyParams(original)
	copied["a"] = 999
	if original["a"] != 1 {
		t.Fatal("copyParams did not create independent copy")
	}
}

func TestTodayDate(t *testing.T) {
	if todayDate() == "" {
		t.Fatal("todayDate() returned empty")
	}
}

func TestProjectMembersFromResponse(t *testing.T) {
	resp := &Response{Body: []byte(`[{"userId":" u-1 ","name":"Alice"},{"userId":"","name":"Bob"},{"userId":"u-2","name":"Charlie"}]`)}
	members, ids, err := projectMembersFromResponse(resp, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 || ids[0] != "u-1" || ids[1] != "u-2" {
		t.Fatalf("ids = %v", ids)
	}
	if members["u-1"] == nil || members["u-2"] == nil {
		t.Fatal("missing members")
	}

	_, _, err = projectMembersFromResponse(&Response{Body: []byte(`{invalid`)}, 10)
	if err == nil {
		t.Fatal("expected decode error")
	}
}

func TestSearchProjectWorkitemsReturnsError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	c, _ := getClient(client)
	if _, err := searchProjectWorkitems(context.Background(), c, "org-1", "project-1", "Task", map[string]any{}); err == nil {
		t.Fatal("expected search error")
	}
}

func TestProjectTaskStatusMembersReturnsError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	c, _ := getClient(client)
	if _, _, err := projectTaskStatusMembers(context.Background(), c, "org-1", "project-1", map[string]any{}); err == nil {
		t.Fatal("expected members fetch error")
	}
}

func TestProjectMembersFromResponseNegativeLimitIncludesAll(t *testing.T) {
	resp := &Response{Body: []byte(`[{"userId":"u-1"},{"userId":"u-2"}]`)}
	_, ids, err := projectMembersFromResponse(resp, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("ids = %v, want 2", ids)
	}
}

func TestParseStatusGroupsTrimsEmptyEntries(t *testing.T) {
	got, err := parseStatusGroups(map[string]any{"statusGroups": `{"todo":"todo-id","":"","  ":"  "}`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got["todo"] != "todo-id" {
		t.Fatalf("unexpected groups: %v", got)
	}
}

func TestHandleGetProjectRiskDashboardWithHighPriorityAndStale(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		conditions, _ := body["conditions"].(string)
		if strings.Contains(conditions, `"fieldIdentifier":"priority"`) {
			seen["highPriority"] = true
		}
		if strings.Contains(conditions, `"fieldIdentifier":"updateStatusAt"`) {
			seen["stale"] = true
		}
		if strings.Contains(conditions, `"fieldIdentifier":"finishTime"`) {
			seen["overdue"] = true
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[]`))
	})

	_, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
		"highPriority":   "p0",
		"staleBefore":    "2024-01-01",
	})
	if err != nil {
		t.Fatalf("handleGetProjectRiskDashboard() error = %v", err)
	}
	if !seen["overdue"] {
		t.Fatal("expected overdue search")
	}
	if !seen["highPriority"] {
		t.Fatal("expected highPriority search")
	}
	if !seen["stale"] {
		t.Fatal("expected stale search")
	}
}

func TestHandleGetProjectRiskDashboardReturnsSearchError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	}); err == nil {
		t.Fatal("expected search error")
	}
}

func TestHandleGetProjectMemberTaskStatusReturnsSearchError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := handleGetProjectMemberTaskStatus(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"assigneeIds":    "user-1",
		"categories":     "Task",
	}); err == nil {
		t.Fatal("expected search error")
	}
}

func TestHandleGetProjectRiskDashboardReturnsFocusSearchError(t *testing.T) {
	requests := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		conditions, _ := body["conditions"].(string)
		if strings.Contains(conditions, `"fieldIdentifier":"priority"`) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[]`))
	})
	if _, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
		"highPriority":   "p0",
	}); err == nil {
		t.Fatal("expected highPriority search error")
	}
}

func TestHandleGetProjectRiskDashboardReturnsStaleSearchError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		conditions, _ := body["conditions"].(string)
		if strings.Contains(conditions, `"fieldIdentifier":"updateStatusAt"`) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[]`))
	})
	if _, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
		"staleBefore":    "2024-01-01",
	}); err == nil {
		t.Fatal("expected stale search error")
	}
}

func TestHandleGetProjectMemberTaskStatusReturnsStatusGroupError(t *testing.T) {
	requests := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		conditions, _ := body["conditions"].(string)
		if strings.Contains(conditions, `"fieldIdentifier":"status"`) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[]`))
	})
	if _, err := handleGetProjectMemberTaskStatus(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"assigneeIds":    "user-1",
		"categories":     "Task",
		"statusGroups":   `{"todo":"todo-id"}`,
	}); err == nil {
		t.Fatal("expected statusGroup search error")
	}
}

// --- Tests for get_sprint_velocity ---

func TestHandleGetSprintVelocityBuildsCorrectRequests(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/sprints" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			seen["sprints"] = true
			w.Header().Set("x-total", "2")
			_, _ = w.Write([]byte(`[{"id":"sp-1","name":"Sprint 1"},{"id":"sp-2","name":"Sprint 2"}]`))
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
			conditions, _ := body["conditions"].(string)
			if !strings.Contains(conditions, `"fieldIdentifier":"sprint"`) {
				t.Fatalf("missing sprint condition: %q", conditions)
			}
			w.Header().Set("x-total", "3")
			_, _ = w.Write([]byte(`[{"id":"wi-1","status":{"stage":"DONE"}},{"id":"wi-2","status":{"stage":"DOING"}},{"id":"wi-3","status":{"stage":"DONE"}}]`))
		default:
			t.Fatalf("method = %s", r.Method)
		}
	})

	result, err := handleGetSprintVelocity(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
		"sprintCount":    float64(2),
	})
	if err != nil {
		t.Fatalf("handleGetSprintVelocity() error = %v", err)
	}
	if !seen["sprints"] {
		t.Fatal("missing sprints request")
	}
	if !seen["cat:Task"] {
		t.Fatal("missing Task search request")
	}
	if !strings.Contains(result, `"sprints"`) {
		t.Fatalf("result missing sprints: %q", result)
	}
	if !strings.Contains(result, `"rate"`) {
		t.Fatalf("result missing rate: %q", result)
	}
}

func TestHandleGetSprintVelocityRejectsEmptyCategories(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without categories")
	})
	if _, err := handleGetSprintVelocity(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     ", ,",
	}); err == nil {
		t.Fatal("expected missing categories error")
	}
}

func TestHandleGetSprintVelocityReturnsSprintError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := handleGetSprintVelocity(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	}); err == nil {
		t.Fatal("expected sprint list error")
	}
}

// --- Tests for get_workitem_status_timeline ---

func TestHandleGetWorkitemStatusTimelineBuildsCorrectRequests(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			switch r.URL.Path {
			case "/oapi/v1/projex/organizations/org-1/workitems/wi-1":
				seen["workitem"] = true
				_, _ = w.Write([]byte(`{"id":"wi-1","subject":"Test"}`))
			case "/oapi/v1/projex/organizations/org-1/workitems/wi-1/activities":
				seen["activities"] = true
				_, _ = w.Write([]byte(`[{"action":"UPDATE","field":"status","gmtCreate":1234567890,"oldValue":{"name":"backlog"},"newValue":{"name":"doing"}}]`))
			default:
				t.Fatalf("unexpected path = %q", r.URL.Path)
			}
		default:
			t.Fatalf("method = %s", r.Method)
		}
	})

	result, err := handleGetWorkitemStatusTimeline(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
	})
	if err != nil {
		t.Fatalf("handleGetWorkitemStatusTimeline() error = %v", err)
	}
	if !seen["workitem"] {
		t.Fatal("missing workitem request")
	}
	if !seen["activities"] {
		t.Fatal("missing activities request")
	}
	if !strings.Contains(result, `"timeline"`) {
		t.Fatalf("result missing timeline: %q", result)
	}
	if !strings.Contains(result, `"summary"`) {
		t.Fatalf("result missing summary: %q", result)
	}
}

func TestHandleGetWorkitemStatusTimelineSkipsWorkitemWhenDisabled(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/projex/organizations/org-1/workitems/wi-1" {
			seen["workitem"] = true
		}
		if r.URL.Path == "/oapi/v1/projex/organizations/org-1/workitems/wi-1/activities" {
			seen["activities"] = true
			_, _ = w.Write([]byte(`[]`))
		}
	})

	_, err := handleGetWorkitemStatusTimeline(context.Background(), client, map[string]any{
		"organizationId":  "org-1",
		"workitemId":      "wi-1",
		"includeWorkitem": false,
	})
	if err != nil {
		t.Fatalf("handleGetWorkitemStatusTimeline() error = %v", err)
	}
	if seen["workitem"] {
		t.Fatal("should not fetch workitem when includeWorkitem=false")
	}
	if !seen["activities"] {
		t.Fatal("missing activities request")
	}
}

func TestHandleGetWorkitemStatusTimelineRequiresWorkitemId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without workitemId")
	})
	if _, err := handleGetWorkitemStatusTimeline(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	}); err == nil {
		t.Fatal("expected missing workitemId error")
	}
}

// --- Tests for get_blocker_analysis ---

func TestHandleGetBlockerAnalysisBuildsCorrectRequests(t *testing.T) {
	requestCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		switch r.Method {
		case http.MethodPost:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			w.Header().Set("x-total", "2")
			_, _ = w.Write([]byte(`[{"id":"wi-1","status":{"name":"doing"}},{"id":"wi-2","status":{"name":"done"}}]`))
		case http.MethodGet:
			if strings.HasSuffix(r.URL.Path, "/relationRecords") {
				if strings.Contains(r.URL.Path, "wi-1") {
					_, _ = w.Write([]byte(`[{"relationType":"DEPEND_ON","target":{"status":{"stage":"DOING"}}}]`))
				} else {
					_, _ = w.Write([]byte(`[]`))
				}
			} else {
				t.Fatalf("unexpected path = %q", r.URL.Path)
			}
		default:
			t.Fatalf("method = %s", r.Method)
		}
	})

	result, err := handleGetBlockerAnalysis(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
		"sampleLimit":    float64(2),
	})
	if err != nil {
		t.Fatalf("handleGetBlockerAnalysis() error = %v", err)
	}
	if requestCount != 3 {
		t.Fatalf("requests = %d, want 3", requestCount)
	}
	if !strings.Contains(result, `"blocked"`) {
		t.Fatalf("result missing blocked: %q", result)
	}
	if !strings.Contains(result, `"blocking"`) {
		t.Fatalf("result missing blocking: %q", result)
	}
	if !strings.Contains(result, `"summary"`) {
		t.Fatalf("result missing summary: %q", result)
	}
}

func TestHandleGetBlockerAnalysisRejectsEmptyCategories(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without categories")
	})
	if _, err := handleGetBlockerAnalysis(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     ", ,",
	}); err == nil {
		t.Fatal("expected missing categories error")
	}
}

// --- Tests for get_member_workload_trend ---

func TestHandleGetMemberWorkloadTrendBuildsCorrectRequests(t *testing.T) {
	requestCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/members" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			_, _ = w.Write([]byte(`[{"userId":"user-1","userName":"Alice"}]`))
		case http.MethodPost:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			conditions, _ := body["conditions"].(string)
			if !strings.Contains(conditions, `"fieldIdentifier":"assignedTo"`) {
				t.Fatalf("missing assignedTo condition: %q", conditions)
			}
			w.Header().Set("x-total", "1")
			_, _ = w.Write([]byte(`[{"id":"wi-1","status":{"name":"doing"},"finishTime":"2099-01-01"}]`))
		default:
			t.Fatalf("method = %s", r.Method)
		}
	})

	result, err := handleGetMemberWorkloadTrend(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
		"memberLimit":    float64(1),
	})
	if err != nil {
		t.Fatalf("handleGetMemberWorkloadTrend() error = %v", err)
	}
	if requestCount != 3 {
		t.Fatalf("requests = %d, want 3", requestCount)
	}
	if !strings.Contains(result, `"members"`) {
		t.Fatalf("result missing members: %q", result)
	}
	if !strings.Contains(result, `"summary"`) {
		t.Fatalf("result missing summary: %q", result)
	}
}

func TestHandleGetMemberWorkloadTrendRejectsEmptyCategories(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without categories")
	})
	if _, err := handleGetMemberWorkloadTrend(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     ", ,",
	}); err == nil {
		t.Fatal("expected missing categories error")
	}
}

func TestHandleGetMemberWorkloadTrendReturnsMembersError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(`[]`))
	})
	if _, err := handleGetMemberWorkloadTrend(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	}); err == nil {
		t.Fatal("expected members fetch error")
	}
}

// --- Tests for get_team_workload_breakdown ---

func TestHandleGetTeamWorkloadBreakdownBuildsCorrectRequests(t *testing.T) {
	requestCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/members" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			_, _ = w.Write([]byte(`[{"userId":"user-1","userName":"Alice"}]`))
		case http.MethodPost:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["perPage"].(float64) != 5 {
				t.Fatalf("taskLimit not applied: perPage = %v", body["perPage"])
			}
			w.Header().Set("x-total", "1")
			_, _ = w.Write([]byte(`[{"id":"wi-1","serialNumber":"ZGPQ-1","subject":"Test task","status":{"name":"doing"},"labels":[{"name":"area/test"}],"gmtCreate":1234567890}]`))
		default:
			t.Fatalf("method = %s", r.Method)
		}
	})

	result, err := handleGetTeamWorkloadBreakdown(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
		"memberLimit":    float64(1),
		"taskLimit":      float64(5),
	})
	if err != nil {
		t.Fatalf("handleGetTeamWorkloadBreakdown() error = %v", err)
	}
	if requestCount != 2 {
		t.Fatalf("requests = %d, want 2", requestCount)
	}
	if !strings.Contains(result, `"tasks"`) {
		t.Fatalf("result missing tasks: %q", result)
	}
	if !strings.Contains(result, `"ZGPQ-1"`) {
		t.Fatalf("result missing task serialNumber: %q", result)
	}
	if !strings.Contains(result, `"area/test"`) {
		t.Fatalf("result missing label: %q", result)
	}
}

func TestHandleGetTeamWorkloadBreakdownRejectsEmptyCategories(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without categories")
	})
	if _, err := handleGetTeamWorkloadBreakdown(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     ", ,",
	}); err == nil {
		t.Fatal("expected missing categories error")
	}
}

func TestHandleGetTeamWorkloadBreakdownReturnsMembersError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(`[]`))
	})
	if _, err := handleGetTeamWorkloadBreakdown(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	}); err == nil {
		t.Fatal("expected members fetch error")
	}
}
