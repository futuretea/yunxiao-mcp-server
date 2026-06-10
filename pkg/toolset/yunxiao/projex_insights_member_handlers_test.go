package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

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

func TestProjexInsightsMemberHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetProjectMemberTaskStatus(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetProjectMemberTaskStatus(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "projectId": "p-1", "assigneeIds": "u-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}
