package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestPlatformHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetCurrentUser(context.Background(), "invalid-client", nil); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetCurrentUser(context.Background(), "invalid-client", nil); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListOrganizations(context.Background(), "invalid-client", nil); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListOrganizations(context.Background(), "invalid-client", nil); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetOrganization(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing id error")
	}
	if _, err := handleGetOrganization(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListOrganizationDepartments(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListOrganizationDepartments(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetOrganizationDepartmentInfo(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetOrganizationDepartmentInfo(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing id error")
	}
	if _, err := handleGetOrganizationDepartmentInfo(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "departmentId": "dept-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetOrganizationDepartmentAncestors(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetOrganizationDepartmentAncestors(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "departmentId": "dept-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListOrganizationMembers(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListOrganizationMembers(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetOrganizationMemberInfo(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetOrganizationMemberInfo(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing id and memberId error")
	}
	if _, err := handleGetOrganizationMemberInfo(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "memberId": "m-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetOrganizationMemberInfoByUserID(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetOrganizationMemberInfoByUserID(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing userId error")
	}
	if _, err := handleGetOrganizationMemberInfoByUserID(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "userId": "u-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleSearchOrganizationMembers(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleSearchOrganizationMembers(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListOrganizationRoles(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListOrganizationRoles(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetOrganizationRole(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetOrganizationRole(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing id and roleId error")
	}
	if _, err := handleGetOrganizationRole(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "roleId": "r-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListUsers(context.Background(), "invalid-client", nil); err == nil {
		t.Fatal("expected getClient error")
	}
}

func TestOrganizationPath(t *testing.T) {
	if got := organizationPath("org-1"); got != "/platform/organizations/org-1" {
		t.Fatalf("organizationPath() = %q", got)
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
