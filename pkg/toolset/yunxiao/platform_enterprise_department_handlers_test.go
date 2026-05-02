package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListEnterpriseDepartmentsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/platform/departments" {
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

	result, err := handleListEnterpriseDepartments(context.Background(), client, map[string]any{
		"parentId": "dept-1",
		"page":     float64(2),
		"perPage":  float64(50),
	})
	if err != nil {
		t.Fatalf("handleListEnterpriseDepartments() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetEnterpriseDepartmentPreservesEncodedID(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/platform/departments/dept%2F1" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"id":"dept-1"}`))
	})

	if _, err := handleGetEnterpriseDepartment(context.Background(), client, map[string]any{
		"id": "dept%2F1",
	}); err != nil {
		t.Fatalf("handleGetEnterpriseDepartment() error = %v", err)
	}
}

func TestHandleGetEnterpriseDepartmentRequiresID(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetEnterpriseDepartment(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("handleGetEnterpriseDepartment() expected missing id error")
	}
}

func TestPlatformEnterpriseDepartmentHandlersRequireParams(t *testing.T) {
	if _, err := handleListEnterpriseDepartments(context.Background(), "invalid-client", map[string]any{}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetEnterpriseDepartment(context.Background(), "invalid-client", map[string]any{"id": "dept-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}
