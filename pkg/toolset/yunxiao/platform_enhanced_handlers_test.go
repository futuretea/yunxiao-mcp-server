package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetOrganizationOverviewRequiresOrganizationId(t *testing.T) {
	_, err := handleGetOrganizationOverview(context.Background(), nil, map[string]any{})
	if err == nil || !strings.Contains(err.Error(), "organizationId is required") {
		t.Fatalf("expected organizationId required error, got %v", err)
	}
}

func TestHandleGetOrganizationOverviewRequiresClient(t *testing.T) {
	_, err := handleGetOrganizationOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil || !strings.Contains(err.Error(), "yunxiao client is not configured") {
		t.Fatalf("expected client error, got %v", err)
	}
}

func TestHandleGetOrganizationOverviewReturnsErrorOnOrgFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleGetOrganizationOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetOrganizationOverviewReturnsErrorOnDepartmentsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/departments") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"dept boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"org-1"}`))
	})
	_, err := handleGetOrganizationOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetOrganizationOverviewReturnsErrorOnMembersFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/members") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"member boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"org-1"}`))
	})
	_, err := handleGetOrganizationOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetOrganizationOverviewReturnsErrorOnGroupsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/groups") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"group boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"org-1"}`))
	})
	_, err := handleGetOrganizationOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetOrganizationOverviewReturnsErrorOnRolesFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/roles") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"role boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"org-1"}`))
	})
	_, err := handleGetOrganizationOverview(context.Background(), client, map[string]any{
		"organizationId":     "org-1",
		"includeDepartments": false,
		"includeMembers":     false,
		"includeGroups":      false,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetOrganizationOverviewSuccessAllSections(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oapi/v1/platform/organizations/org-1":
			_, _ = w.Write([]byte(`{"id":"org-1"}`))
		case strings.HasSuffix(r.URL.Path, "/departments"):
			if r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("depts perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`["dept-a"]`))
		case strings.HasSuffix(r.URL.Path, "/members"):
			if r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("members perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`["member-a"]`))
		case strings.HasSuffix(r.URL.Path, "/groups"):
			if r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("groups perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`["group-a"]`))
		case strings.HasSuffix(r.URL.Path, "/roles"):
			_, _ = w.Write([]byte(`["role-a"]`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})

	result, err := handleGetOrganizationOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	})
	if err != nil {
		t.Fatalf("handleGetOrganizationOverview() error = %v", err)
	}
	if !strings.Contains(result, `"organization"`) {
		t.Fatalf("result missing organization: %q", result)
	}
	if !strings.Contains(result, `"departments"`) {
		t.Fatalf("result missing departments: %q", result)
	}
	if !strings.Contains(result, `"members"`) {
		t.Fatalf("result missing members: %q", result)
	}
	if !strings.Contains(result, `"groups"`) {
		t.Fatalf("result missing groups: %q", result)
	}
	if !strings.Contains(result, `"roles"`) {
		t.Fatalf("result missing roles: %q", result)
	}
}

func TestHandleGetOrganizationOverviewSkipsSectionsWhenDisabled(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/platform/organizations/org-1" {
			_, _ = w.Write([]byte(`{"id":"org-1"}`))
			return
		}
		t.Fatalf("unexpected request to %q", r.URL.Path)
	})

	result, err := handleGetOrganizationOverview(context.Background(), client, map[string]any{
		"organizationId":     "org-1",
		"includeDepartments": false,
		"includeMembers":     false,
		"includeGroups":      false,
		"includeRoles":       false,
	})
	if err != nil {
		t.Fatalf("handleGetOrganizationOverview() error = %v", err)
	}
	if strings.Contains(result, `"departments"`) {
		t.Fatalf("result should not contain departments: %q", result)
	}
	if strings.Contains(result, `"members"`) {
		t.Fatalf("result should not contain members: %q", result)
	}
	if strings.Contains(result, `"groups"`) {
		t.Fatalf("result should not contain groups: %q", result)
	}
	if strings.Contains(result, `"roles"`) {
		t.Fatalf("result should not contain roles: %q", result)
	}
}

func TestHandleGetOrganizationOverviewUsesCustomLimits(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/departments"):
			if r.URL.Query().Get("perPage") != "3" {
				t.Fatalf("depts perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[]`))
		case strings.HasSuffix(r.URL.Path, "/members"):
			if r.URL.Query().Get("perPage") != "2" {
				t.Fatalf("members perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[]`))
		case strings.HasSuffix(r.URL.Path, "/groups"):
			if r.URL.Query().Get("perPage") != "4" {
				t.Fatalf("groups perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[]`))
		default:
			_, _ = w.Write([]byte(`{"id":"org-1"}`))
		}
	})

	_, err := handleGetOrganizationOverview(context.Background(), client, map[string]any{
		"organizationId":  "org-1",
		"departmentLimit": float64(3),
		"memberLimit":     float64(2),
		"groupLimit":      float64(4),
		"includeRoles":    false,
	})
	if err != nil {
		t.Fatalf("handleGetOrganizationOverview() error = %v", err)
	}
}

func TestOrganizationOverviewFilters(t *testing.T) {
	params := map[string]any{
		"includeDepartments": false,
		"includeMembers":     false,
		"includeGroups":      false,
		"includeRoles":       false,
		"departmentLimit":    float64(10),
		"memberLimit":        float64(20),
		"groupLimit":         float64(30),
	}
	filters := organizationOverviewFilters(params)
	if filters["includeDepartments"].(bool) != false {
		t.Fatalf("includeDepartments = %v", filters["includeDepartments"])
	}
	if filters["includeMembers"].(bool) != false {
		t.Fatalf("includeMembers = %v", filters["includeMembers"])
	}
	if filters["includeGroups"].(bool) != false {
		t.Fatalf("includeGroups = %v", filters["includeGroups"])
	}
	if filters["includeRoles"].(bool) != false {
		t.Fatalf("includeRoles = %v", filters["includeRoles"])
	}
	if filters["departmentLimit"].(int) != 10 {
		t.Fatalf("departmentLimit = %v", filters["departmentLimit"])
	}
	if filters["memberLimit"].(int) != 20 {
		t.Fatalf("memberLimit = %v", filters["memberLimit"])
	}
	if filters["groupLimit"].(int) != 30 {
		t.Fatalf("groupLimit = %v", filters["groupLimit"])
	}
}
