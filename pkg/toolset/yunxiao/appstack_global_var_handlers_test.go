package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestHandleGetGlobalVarBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/globalVars/group%2F1?revisionSha=rev%2F1" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"name":"group-1"}`))
	})

	if _, err := handleGetGlobalVar(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"name":           "group/1",
		"revisionSha":    "rev/1",
	}); err != nil {
		t.Fatalf("handleGetGlobalVar() error = %v", err)
	}
}

func TestHandleListGlobalVarsBuildsPathQueryAndBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/globalVars:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("current") != "2" || r.URL.Query().Get("pageSize") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["search"] != "demo" {
			t.Fatalf("body = %#v", body)
		}
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	if _, err := handleListGlobalVars(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"current":        float64(2),
		"pageSize":       float64(20),
		"search":         "demo",
	}); err != nil {
		t.Fatalf("handleListGlobalVars() error = %v", err)
	}
}

func TestHandleListGlobalVarsUsesDefaultPagination(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Query().Get("current") != "1" || r.URL.Query().Get("pageSize") != "10" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	if _, err := handleListGlobalVars(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	}); err != nil {
		t.Fatalf("handleListGlobalVars() error = %v", err)
	}
}

func TestAppstackGlobalVarHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetGlobalVar(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetGlobalVar(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing name error")
	}
	if _, err := handleGetGlobalVar(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "name": "g-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListGlobalVars(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListGlobalVars(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}

func TestAppstackGlobalVarPath(t *testing.T) {
	if got := appstackGlobalVarPath("org-1", "group/1"); got != "/appstack/organizations/org-1/globalVars/group%2F1" {
		t.Fatalf("appstackGlobalVarPath() = %q", got)
	}
}
