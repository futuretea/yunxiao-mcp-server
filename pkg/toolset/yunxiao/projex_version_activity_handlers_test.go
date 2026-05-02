package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListVersionsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/versions" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("status") != "TODO,DOING" ||
			r.URL.Query().Get("name") != "release" ||
			r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"version-1"}]`))
	})

	result, err := handleListVersions(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"status":         "TODO,DOING",
		"name":           "release",
		"page":           float64(2),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleListVersions() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleListWorkitemActivitiesBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/workitems/workitem-1/activities" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"action":"update"}]`))
	})

	if _, err := handleListWorkitemActivities(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "workitem-1",
	}); err != nil {
		t.Fatalf("handleListWorkitemActivities() error = %v", err)
	}
}

func TestHandleListVersionsRequiresOrganizationId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleListVersions(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
}

func TestHandleListWorkitemActivitiesRequiresOrganizationId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleListWorkitemActivities(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
}

func TestProjexVersionActivityHandlersRequireParams(t *testing.T) {
	if _, err := handleListVersions(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "id": "project-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListWorkitemActivities(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "id": "wi-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}
