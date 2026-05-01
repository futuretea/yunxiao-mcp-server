package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListTemplateRepositoriesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/codeup/organizations/org-1/repositories/templates" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("templateType") != "2" ||
			r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "3")
		_, _ = w.Write([]byte(`[{"name":"go-template"}]`))
	})

	result, err := handleListTemplateRepositories(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"templateType":   float64(2),
		"page":           float64(2),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleListTemplateRepositories() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleListTemplateRepositoriesRequiresTemplateType(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	_, err := handleListTemplateRepositories(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil || !strings.Contains(err.Error(), "templateType is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestHandleListNamespacesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/codeup/organizations/org-1/namespaces" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("parentId") != "1" ||
			r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "20" ||
			r.URL.Query().Get("search") != "demo" ||
			r.URL.Query().Get("orderBy") != "updated_at" ||
			r.URL.Query().Get("sort") != "desc" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":1}]`))
	})

	result, err := handleListNamespaces(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"parentId":       float64(1),
		"page":           float64(2),
		"perPage":        float64(20),
		"search":         "demo",
		"orderBy":        "updated_at",
		"sort":           "desc",
	})
	if err != nil {
		t.Fatalf("handleListNamespaces() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetNamespaceBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/namespaces/group%2Fsub" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"id":1}`))
	})

	if _, err := handleGetNamespace(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"namespaceId":    "group/sub",
	}); err != nil {
		t.Fatalf("handleGetNamespace() error = %v", err)
	}
}

func TestHandleGetOrgNamespaceBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/orgNamespace" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"id":1}`))
	})

	if _, err := handleGetOrgNamespace(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	}); err != nil {
		t.Fatalf("handleGetOrgNamespace() error = %v", err)
	}
}
