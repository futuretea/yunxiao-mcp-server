package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetOrganizationDepartmentOverviewRequiresOrganizationId(t *testing.T) {
	_, err := handleGetOrganizationDepartmentOverview(context.Background(), nil, map[string]any{
		"departmentId": "dept-1",
	})
	if err == nil || !strings.Contains(err.Error(), "organizationId is required") {
		t.Fatalf("expected organizationId required error, got %v", err)
	}
}

func TestHandleGetOrganizationDepartmentOverviewRequiresDepartmentId(t *testing.T) {
	_, err := handleGetOrganizationDepartmentOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil || !strings.Contains(err.Error(), "departmentId is required") {
		t.Fatalf("expected departmentId required error, got %v", err)
	}
}

func TestHandleGetOrganizationDepartmentOverviewRequiresClient(t *testing.T) {
	_, err := handleGetOrganizationDepartmentOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"departmentId":   "dept-1",
	})
	if err == nil || !strings.Contains(err.Error(), "yunxiao client is not configured") {
		t.Fatalf("expected client error, got %v", err)
	}
}

func TestHandleGetOrganizationDepartmentOverviewReturnsErrorOnDeptFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleGetOrganizationDepartmentOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"departmentId":   "dept-1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetOrganizationDepartmentOverviewReturnsErrorOnAncestorsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/ancestors") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"ancestors boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"dept-1"}`))
	})
	_, err := handleGetOrganizationDepartmentOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"departmentId":   "dept-1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetOrganizationDepartmentOverviewSuccessAllSections(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations/org-1/departments/dept-1":
			_, _ = w.Write([]byte(`{"id":"dept-1"}`))
		case "/oapi/v1/platform/organizations/org-1/departments/dept-1/ancestors":
			_, _ = w.Write([]byte(`["ancestor-1"]`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})

	result, err := handleGetOrganizationDepartmentOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"departmentId":   "dept-1",
	})
	if err != nil {
		t.Fatalf("handleGetOrganizationDepartmentOverview() error = %v", err)
	}
	if !strings.Contains(result, `"department"`) {
		t.Fatalf("result missing department: %q", result)
	}
	if !strings.Contains(result, `"ancestors"`) {
		t.Fatalf("result missing ancestors: %q", result)
	}
}

func TestHandleGetOrganizationDepartmentOverviewSkipsAncestorsWhenDisabled(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/platform/organizations/org-1/departments/dept-1" {
			_, _ = w.Write([]byte(`{"id":"dept-1"}`))
			return
		}
		t.Fatalf("unexpected request to %q", r.URL.Path)
	})

	result, err := handleGetOrganizationDepartmentOverview(context.Background(), client, map[string]any{
		"organizationId":   "org-1",
		"departmentId":     "dept-1",
		"includeAncestors": false,
	})
	if err != nil {
		t.Fatalf("handleGetOrganizationDepartmentOverview() error = %v", err)
	}
	if strings.Contains(result, `"ancestors"`) {
		t.Fatalf("result should not contain ancestors: %q", result)
	}
}

func TestOrganizationDepartmentOverviewFilters(t *testing.T) {
	params := map[string]any{
		"includeAncestors": false,
	}
	filters := organizationDepartmentOverviewFilters(params)
	if filters["includeAncestors"].(bool) != false {
		t.Fatalf("includeAncestors = %v", filters["includeAncestors"])
	}
}

func TestHandleGetOrganizationGroupOverviewRequiresOrganizationId(t *testing.T) {
	_, err := handleGetOrganizationGroupOverview(context.Background(), nil, map[string]any{
		"groupId": "group-1",
	})
	if err == nil || !strings.Contains(err.Error(), "organizationId is required") {
		t.Fatalf("expected organizationId required error, got %v", err)
	}
}

func TestHandleGetOrganizationGroupOverviewRequiresGroupId(t *testing.T) {
	_, err := handleGetOrganizationGroupOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil || !strings.Contains(err.Error(), "groupId is required") {
		t.Fatalf("expected groupId required error, got %v", err)
	}
}

func TestHandleGetOrganizationGroupOverviewRequiresClient(t *testing.T) {
	_, err := handleGetOrganizationGroupOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"groupId":        "group-1",
	})
	if err == nil || !strings.Contains(err.Error(), "yunxiao client is not configured") {
		t.Fatalf("expected client error, got %v", err)
	}
}

func TestHandleGetOrganizationGroupOverviewReturnsErrorOnGroupFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleGetOrganizationGroupOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"groupId":        "group-1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetOrganizationGroupOverviewReturnsErrorOnMembersFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/members") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"members boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"group-1"}`))
	})
	_, err := handleGetOrganizationGroupOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"groupId":        "group-1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetOrganizationGroupOverviewSuccessAllSections(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations/org-1/groups/group-1":
			_, _ = w.Write([]byte(`{"id":"group-1"}`))
		case "/oapi/v1/platform/organizations/org-1/groups/group-1/members":
			if r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("members perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`["member-a"]`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})

	result, err := handleGetOrganizationGroupOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"groupId":        "group-1",
	})
	if err != nil {
		t.Fatalf("handleGetOrganizationGroupOverview() error = %v", err)
	}
	if !strings.Contains(result, `"group"`) {
		t.Fatalf("result missing group: %q", result)
	}
	if !strings.Contains(result, `"members"`) {
		t.Fatalf("result missing members: %q", result)
	}
}

func TestHandleGetOrganizationGroupOverviewSkipsMembersWhenDisabled(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/platform/organizations/org-1/groups/group-1" {
			_, _ = w.Write([]byte(`{"id":"group-1"}`))
			return
		}
		t.Fatalf("unexpected request to %q", r.URL.Path)
	})

	result, err := handleGetOrganizationGroupOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"groupId":        "group-1",
		"includeMembers": false,
	})
	if err != nil {
		t.Fatalf("handleGetOrganizationGroupOverview() error = %v", err)
	}
	if strings.Contains(result, `"members"`) {
		t.Fatalf("result should not contain members: %q", result)
	}
}

func TestHandleGetOrganizationGroupOverviewUsesCustomMemberLimit(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/members") {
			if r.URL.Query().Get("perPage") != "3" {
				t.Fatalf("members perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[]`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"group-1"}`))
	})

	_, err := handleGetOrganizationGroupOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"groupId":        "group-1",
		"memberLimit":    float64(3),
	})
	if err != nil {
		t.Fatalf("handleGetOrganizationGroupOverview() error = %v", err)
	}
}

func TestOrganizationGroupOverviewFilters(t *testing.T) {
	params := map[string]any{
		"includeMembers": false,
		"memberLimit":    float64(10),
	}
	filters := organizationGroupOverviewFilters(params)
	if filters["includeMembers"].(bool) != false {
		t.Fatalf("includeMembers = %v", filters["includeMembers"])
	}
	if filters["memberLimit"].(int) != 10 {
		t.Fatalf("memberLimit = %v", filters["memberLimit"])
	}
}
