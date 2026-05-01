package yunxiao

import (
	"context"
	"net/http"
	"testing"
)

func TestHandleGetLatestOrchestrationBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/apps/app%2F1/envs/dev%2F1/orchestration:latestAvailable" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"name":"orch-1"}`))
	})

	if _, err := handleGetLatestOrchestration(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app/1",
		"envName":        "dev/1",
	}); err != nil {
		t.Fatalf("handleGetLatestOrchestration() error = %v", err)
	}
}

func TestHandleListAppOrchestrationBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/orchestrations" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	if _, err := handleListAppOrchestration(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	}); err != nil {
		t.Fatalf("handleListAppOrchestration() error = %v", err)
	}
}

func TestHandleGetAppOrchestrationBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/apps/app%2F1/orchestrations/orch%2F1?sha=abc123&tagName=stable" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if r.URL.Query().Get("tagName") != "stable" || r.URL.Query().Get("sha") != "abc123" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"sn":"orch-1"}`))
	})

	if _, err := handleGetAppOrchestration(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app/1",
		"sn":             "orch/1",
		"tagName":        "stable",
		"sha":            "abc123",
	}); err != nil {
		t.Fatalf("handleGetAppOrchestration() error = %v", err)
	}
}
