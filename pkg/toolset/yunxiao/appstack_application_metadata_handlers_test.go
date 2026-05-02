package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleSearchAppTemplatesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/appTemplates:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("pagination") != "keyset" ||
			r.URL.Query().Get("perPage") != "20" ||
			r.URL.Query().Get("orderBy") != "id" ||
			r.URL.Query().Get("sort") != "asc" ||
			r.URL.Query().Get("nextToken") != "token-1" ||
			r.URL.Query().Get("displayNameKeyword") != "demo" ||
			r.URL.Query().Get("page") != "2" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	if _, err := handleSearchAppTemplates(context.Background(), client, map[string]any{
		"organizationId":     "org-1",
		"pagination":         "keyset",
		"perPage":            float64(20),
		"orderBy":            "id",
		"sort":               "asc",
		"nextToken":          "token-1",
		"displayNameKeyword": "demo",
		"page":               float64(2),
	}); err != nil {
		t.Fatalf("handleSearchAppTemplates() error = %v", err)
	}
}

func TestHandleListEnvironmentsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/apps/app%2F1/envs?") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		if r.URL.Query().Get("pagination") != "keyset" ||
			r.URL.Query().Get("perPage") != "10" ||
			r.URL.Query().Get("orderBy") != "id" ||
			r.URL.Query().Get("sort") != "desc" ||
			r.URL.Query().Get("nextToken") != "next" ||
			r.URL.Query().Get("page") != "3" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	if _, err := handleListEnvironments(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app/1",
		"pagination":     "keyset",
		"perPage":        float64(10),
		"orderBy":        "id",
		"sort":           "desc",
		"nextToken":      "next",
		"page":           float64(3),
	}); err != nil {
		t.Fatalf("handleListEnvironments() error = %v", err)
	}
}

func TestHandleGetEnvironmentBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/apps/app%2F1/envs/env%2F1" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"name":"env-1"}`))
	})

	if _, err := handleGetEnvironment(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app/1",
		"envName":        "env/1",
	}); err != nil {
		t.Fatalf("handleGetEnvironment() error = %v", err)
	}
}

func TestHandleListApplicationMembersBuildsPathAndDefaultQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/apps/app%2F1/members?current=1&pageSize=10" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	if _, err := handleListApplicationMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app/1",
	}); err != nil {
		t.Fatalf("handleListApplicationMembers() error = %v", err)
	}
}

func TestHandleListApplicationSourcesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/sources" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("pagination") != "keyset" ||
			r.URL.Query().Get("perPage") != "20" ||
			r.URL.Query().Get("orderBy") != "gmtCreate" ||
			r.URL.Query().Get("sort") != "asc" ||
			r.URL.Query().Get("nextToken") != "token-1" ||
			r.URL.Query().Get("page") != "2" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	if _, err := handleListApplicationSources(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"pagination":     "keyset",
		"perPage":        float64(20),
		"orderBy":        "gmtCreate",
		"sort":           "asc",
		"nextToken":      "token-1",
		"page":           float64(2),
	}); err != nil {
		t.Fatalf("handleListApplicationSources() error = %v", err)
	}
}

func TestHandleSearchAppTemplatesRequiresOrganizationId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleSearchAppTemplates(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
}

func TestHandleListEnvironmentsRequiresAppName(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleListEnvironments(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing appName error")
	}
}

func TestHandleGetEnvironmentRequiresEnvName(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleGetEnvironment(context.Background(), client, map[string]any{"organizationId": "org-1", "appName": "app-1"}); err == nil {
		t.Fatal("expected missing envName error")
	}
}

func TestHandleListApplicationMembersRequiresAppName(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleListApplicationMembers(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing appName error")
	}
}

func TestHandleListApplicationSourcesRequiresAppName(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleListApplicationSources(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing appName error")
	}
}

func TestRequiredAppEnvironmentRequiresEnvName(t *testing.T) {
	_, _, _, err := requiredAppEnvironment(map[string]any{"organizationId": "org-1", "appName": "app-1"})
	if err == nil {
		t.Fatal("expected missing envName error")
	}
}

func TestRequiredAppEnvironmentRequiresOrganizationId(t *testing.T) {
	_, _, _, err := requiredAppEnvironment(map[string]any{"appName": "app-1", "envName": "env-1"})
	if err == nil {
		t.Fatal("expected missing organizationId error")
	}
}
