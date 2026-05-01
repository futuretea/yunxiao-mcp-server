package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetEnvVariableGroupsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/apps/app%2F1/envs/dev%2F1/variableGroups") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"name":"vg-1"}]`))
	})

	if _, err := handleGetEnvVariableGroups(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app/1",
		"envName":        "dev/1",
	}); err != nil {
		t.Fatalf("handleGetEnvVariableGroups() error = %v", err)
	}
}

func TestHandleGetVariableGroupBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/apps/app%2F1/variableGroup/vg%2F1") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"name":"vg-1"}`))
	})

	if _, err := handleGetVariableGroup(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app/1",
		"variableGroupName": "vg/1",
	}); err != nil {
		t.Fatalf("handleGetVariableGroup() error = %v", err)
	}
}

func TestHandleGetAppVariableGroupsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/variableGroups" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"name":"vg-1"}]`))
	})

	if _, err := handleGetAppVariableGroups(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	}); err != nil {
		t.Fatalf("handleGetAppVariableGroups() error = %v", err)
	}
}

func TestHandleGetAppVariableGroupsRevisionBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/variableGroups:revision" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"sha":"rev-1"}`))
	})

	if _, err := handleGetAppVariableGroupsRevision(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	}); err != nil {
		t.Fatalf("handleGetAppVariableGroupsRevision() error = %v", err)
	}
}
