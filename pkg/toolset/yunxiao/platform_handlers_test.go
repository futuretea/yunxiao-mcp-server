package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListOrganizationDepartmentsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/platform/organizations/org-1/departments" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("parentId") != "dept-1" ||
			r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "50" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"dept-2"}]`))
	})

	result, err := handleListOrganizationDepartments(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"parentId":       "dept-1",
		"page":           float64(2),
		"perPage":        float64(50),
	})
	if err != nil {
		t.Fatalf("handleListOrganizationDepartments() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetOrganizationDepartmentInfoBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/platform/organizations/org-1/departments/dept-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"dept-1"}`))
	})

	if _, err := handleGetOrganizationDepartmentInfo(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "dept-1",
	}); err != nil {
		t.Fatalf("handleGetOrganizationDepartmentInfo() error = %v", err)
	}
}

func TestHandleGetOrganizationDepartmentAncestorsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/platform/organizations/org-1/departments/dept-1/ancestors" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"id":"root"}]`))
	})

	if _, err := handleGetOrganizationDepartmentAncestors(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "dept-1",
	}); err != nil {
		t.Fatalf("handleGetOrganizationDepartmentAncestors() error = %v", err)
	}
}

func TestHandleListOrganizationMembersBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/platform/organizations/org-1/members" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "3" || r.URL.Query().Get("perPage") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"member-1"}]`))
	})

	result, err := handleListOrganizationMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"page":           float64(3),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleListOrganizationMembers() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetOrganizationMemberInfoBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/platform/organizations/org-1/members/member-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"member-1"}`))
	})

	if _, err := handleGetOrganizationMemberInfo(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"memberId":       "member-1",
	}); err != nil {
		t.Fatalf("handleGetOrganizationMemberInfo() error = %v", err)
	}
}

func TestHandleGetOrganizationMemberInfoByUserIDBuildsQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/platform/organizations/org-1/members:readByUser" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("userId") != "user-1" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"userId":"user-1"}`))
	})

	if _, err := handleGetOrganizationMemberInfoByUserID(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"userId":         "user-1",
	}); err != nil {
		t.Fatalf("handleGetOrganizationMemberInfoByUserID() error = %v", err)
	}
}

func TestHandleSearchOrganizationMembersBuildsBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/platform/organizations/org-1/members:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["query"] != "alice" || body["includeChildren"] != true || body["page"] != float64(1) {
			t.Fatalf("body = %#v", body)
		}
		if got := body["deptIds"].([]any); len(got) != 2 || got[0] != "dept-1" || got[1] != "dept-2" {
			t.Fatalf("deptIds = %#v", body["deptIds"])
		}
		if got := body["roleIds"].([]any); len(got) != 1 || got[0] != "role-1" {
			t.Fatalf("roleIds = %#v", body["roleIds"])
		}
		if got := body["statuses"].([]any); len(got) != 1 || got[0] != "ENABLED" {
			t.Fatalf("statuses = %#v", body["statuses"])
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"member-1"}]`))
	})

	result, err := handleSearchOrganizationMembers(context.Background(), client, map[string]any{
		"organizationId":  "org-1",
		"deptIds":         []any{"dept-1", "dept-2"},
		"query":           "alice",
		"includeChildren": true,
		"roleIds":         []any{"role-1"},
		"statuses":        []any{"ENABLED"},
		"page":            float64(1),
		"perPage":         float64(100),
	})
	if err != nil {
		t.Fatalf("handleSearchOrganizationMembers() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleOrganizationRoleToolsBuildPaths(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations/org-1/roles",
			"/oapi/v1/platform/organizations/org-1/roles/role-1":
			_, _ = w.Write([]byte(`{"ok":true}`))
		default:
			t.Fatalf("path = %q", r.URL.Path)
		}
	})

	if _, err := handleListOrganizationRoles(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	}); err != nil {
		t.Fatalf("handleListOrganizationRoles() error = %v", err)
	}
	if _, err := handleGetOrganizationRole(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"roleId":         "role-1",
	}); err != nil {
		t.Fatalf("handleGetOrganizationRole() error = %v", err)
	}
}

func TestHandleListUsersBuildsQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/platform/users" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("filter") != "alice" ||
			r.URL.Query().Get("status") != "enabled" ||
			r.URL.Query().Get("deptId") != "dept-1" ||
			r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "30" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"user-1"}]`))
	})

	result, err := handleListUsers(context.Background(), client, map[string]any{
		"filter":  "alice",
		"status":  "enabled",
		"deptId":  "dept-1",
		"page":    float64(2),
		"perPage": float64(30),
	})
	if err != nil {
		t.Fatalf("handleListUsers() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}
