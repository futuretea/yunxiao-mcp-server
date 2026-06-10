package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetDepartmentUsageBuildsQueryWithMetadata(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/lingma/organizations/org-1/departmentUsage" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("departmentId") != "dept-1" ||
			r.URL.Query().Get("startTime") != "2026-04-01" ||
			r.URL.Query().Get("endTime") != "2026-04-30" ||
			r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "100" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"date":"2026-04-01"}]`))
	})

	result, err := handleGetDepartmentUsage(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"departmentId":   "dept-1",
		"startTime":      "2026-04-01",
		"endTime":        "2026-04-30",
		"page":           float64(2),
		"perPage":        float64(100),
	})
	if err != nil {
		t.Fatalf("handleGetDepartmentUsage() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetDepartmentUsageRequiresDepartmentId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	_, err := handleGetDepartmentUsage(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"startTime":      "2026-04-01",
		"endTime":        "2026-04-30",
	})
	if err == nil {
		t.Fatal("handleGetDepartmentUsage() expected missing departmentId error")
	}
}

func TestHandleListDeveloperMembersBuildsQueryWithMetadata(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/lingma/organizations/org-1/developer/members" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("departmentId") != "dept-1" ||
			r.URL.Query().Get("userId") != "user-1" ||
			r.URL.Query().Get("page") != "1" ||
			r.URL.Query().Get("perPage") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"userId":"user-1"}]`))
	})

	result, err := handleListDeveloperMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"departmentId":   "dept-1",
		"userId":         "user-1",
		"page":           float64(1),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleListDeveloperMembers() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetDeveloperUsageRequiresUserOrDepartment(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	_, err := handleGetDeveloperUsage(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"startTime":      "2026-04-01",
		"endTime":        "2026-04-30",
	})
	if err == nil {
		t.Fatal("handleGetDeveloperUsage() expected missing userId/departmentId error")
	}
}

func TestLingmaUsageHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetDepartmentUsage(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetDepartmentUsage(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "departmentId": "dept-1", "startTime": "2026-04-01", "endTime": "2026-04-30"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListDeveloperMembers(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListDeveloperMembers(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetDeveloperUsage(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetDeveloperUsage(context.Background(), client, map[string]any{"organizationId": "org-1", "userId": "user-1"}); err == nil {
		t.Fatal("expected missing startTime error")
	}
	if _, err := handleGetDeveloperUsage(context.Background(), client, map[string]any{"organizationId": "org-1", "userId": "user-1", "startTime": "2026-04-01"}); err == nil {
		t.Fatal("expected missing endTime error")
	}
	if _, err := handleGetDeveloperUsage(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "userId": "user-1", "startTime": "2026-04-01", "endTime": "2026-04-30"}); err == nil {
		t.Fatal("expected getClient error")
	}
}

func TestRequiredLingmaUsageQueryRequiresTimeRange(t *testing.T) {
	if _, err := requiredLingmaUsageQuery(map[string]any{}); err == nil {
		t.Fatal("expected missing startTime error")
	}
	if _, err := requiredLingmaUsageQuery(map[string]any{"startTime": "2026-04-01"}); err == nil {
		t.Fatal("expected missing endTime error")
	}
}

func TestHandleGetDeveloperUsageBuildsQueryWithUserID(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/lingma/organizations/org-1/developerUsage" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("userId") != "user-1" ||
			r.URL.Query().Get("departmentId") != "" ||
			r.URL.Query().Get("startTime") != "2026-04-01" ||
			r.URL.Query().Get("endTime") != "2026-04-30" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"userId":"user-1"}]`))
	})

	result, err := handleGetDeveloperUsage(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"userId":         "user-1",
		"startTime":      "2026-04-01",
		"endTime":        "2026-04-30",
	})
	if err != nil {
		t.Fatalf("handleGetDeveloperUsage() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetDeveloperUsageBuildsQueryWithDepartmentID(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/lingma/organizations/org-1/developerUsage" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("userId") != "" ||
			r.URL.Query().Get("departmentId") != "dept-1" ||
			r.URL.Query().Get("startTime") != "2026-04-01" ||
			r.URL.Query().Get("endTime") != "2026-04-30" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"departmentId":"dept-1"}]`))
	})

	result, err := handleGetDeveloperUsage(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"departmentId":   "dept-1",
		"startTime":      "2026-04-01",
		"endTime":        "2026-04-30",
	})
	if err != nil {
		t.Fatalf("handleGetDeveloperUsage() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestLingmaUsageHandlersReturnAPIError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := handleGetDeveloperUsage(context.Background(), client, map[string]any{
		"organizationId": "org-1", "userId": "user-1", "startTime": "2026-04-01", "endTime": "2026-04-30",
	}); err == nil {
		t.Fatal("expected API error")
	}
}

func TestLingmaOrganizationPathEncodesValue(t *testing.T) {
	if got := lingmaOrganizationPath("org/1"); got != "/lingma/organizations/org%2F1" {
		t.Fatalf("lingmaOrganizationPath() = %q", got)
	}
	if got := lingmaOrganizationPath("org%2F1"); got != "/lingma/organizations/org%2F1" {
		t.Fatalf("lingmaOrganizationPath() = %q", got)
	}
}
