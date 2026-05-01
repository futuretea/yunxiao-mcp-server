package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListAppReleaseWorkflowsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/releaseWorkflows" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"sn":"rw-1"}]`))
	})

	if _, err := handleListAppReleaseWorkflows(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	}); err != nil {
		t.Fatalf("handleListAppReleaseWorkflows() error = %v", err)
	}
}

func TestHandleListAppReleaseWorkflowBriefsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/releaseWorkflowBriefs" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"sn":"rw-1"}]`))
	})

	if _, err := handleListAppReleaseWorkflowBriefs(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	}); err != nil {
		t.Fatalf("handleListAppReleaseWorkflowBriefs() error = %v", err)
	}
}

func TestHandleGetAppReleaseWorkflowStageBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/apps/app%2F1/releaseWorkflow/rw%2F1/releaseStage/rs%2F1") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"sn":"rs-1"}`))
	})

	if _, err := handleGetAppReleaseWorkflowStage(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app/1",
		"releaseWorkflowSn": "rw/1",
		"releaseStageSn":    "rs/1",
	}); err != nil {
		t.Fatalf("handleGetAppReleaseWorkflowStage() error = %v", err)
	}
}

func TestHandleListAppReleaseStageBriefsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/releaseWorkflow/rw-1/releaseStageBriefs" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"sn":"rs-1"}]`))
	})

	if _, err := handleListAppReleaseStageBriefs(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"appName":           "app-1",
		"releaseWorkflowSn": "rw-1",
	}); err != nil {
		t.Fatalf("handleListAppReleaseStageBriefs() error = %v", err)
	}
}
