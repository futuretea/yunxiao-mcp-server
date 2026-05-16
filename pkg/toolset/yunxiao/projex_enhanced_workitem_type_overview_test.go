package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetWorkItemTypeOverviewRequiresClient(t *testing.T) {
	_, err := handleGetWorkItemTypeOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"workItemTypeId": "type-1",
	})
	if err == nil || !strings.Contains(err.Error(), "yunxiao client is not configured") {
		t.Fatalf("expected client error, got %v", err)
	}
}

func TestHandleGetWorkItemTypeOverviewRequiresOrganizationId(t *testing.T) {
	_, err := handleGetWorkItemTypeOverview(context.Background(), &Client{}, map[string]any{
		"projectId":      "project-1",
		"workItemTypeId": "type-1",
	})
	if err == nil || !strings.Contains(err.Error(), "organizationId is required") {
		t.Fatalf("expected organizationId required error, got %v", err)
	}
}

func TestHandleGetWorkItemTypeOverviewRequiresProjectId(t *testing.T) {
	_, err := handleGetWorkItemTypeOverview(context.Background(), &Client{}, map[string]any{
		"organizationId": "org-1",
		"workItemTypeId": "type-1",
	})
	if err == nil || !strings.Contains(err.Error(), "projectId is required") {
		t.Fatalf("expected projectId required error, got %v", err)
	}
}

func TestHandleGetWorkItemTypeOverviewRequiresWorkItemTypeId(t *testing.T) {
	_, err := handleGetWorkItemTypeOverview(context.Background(), &Client{}, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
	})
	if err == nil || !strings.Contains(err.Error(), "workItemTypeId is required") {
		t.Fatalf("expected workItemTypeId required error, got %v", err)
	}
}

func TestHandleGetWorkItemTypeOverviewReturnsErrorOnTypeFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleGetWorkItemTypeOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"workItemTypeId": "type-1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetWorkItemTypeOverviewReturnsErrorOnFieldConfigFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/fields") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"fields boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"type-1"}`))
	})
	_, err := handleGetWorkItemTypeOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"workItemTypeId": "type-1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetWorkItemTypeOverviewReturnsErrorOnWorkflowFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/workflows") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"workflow boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"type-1"}`))
	})
	_, err := handleGetWorkItemTypeOverview(context.Background(), client, map[string]any{
		"organizationId":     "org-1",
		"projectId":          "project-1",
		"workItemTypeId":     "type-1",
		"includeFieldConfig": false,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetWorkItemTypeOverviewSuccessAllSections(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.EscapedPath() {
		case "/oapi/v1/projex/organizations/org-1/workitemTypes/type-1":
			_, _ = w.Write([]byte(`{"id":"type-1"}`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/workitemTypes/type-1/fields":
			_, _ = w.Write([]byte(`["field-1"]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/workitemTypes/type-1/workflows":
			_, _ = w.Write([]byte(`["status-1"]`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})

	result, err := handleGetWorkItemTypeOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"workItemTypeId": "type-1",
	})
	if err != nil {
		t.Fatalf("handleGetWorkItemTypeOverview() error = %v", err)
	}
	if !strings.Contains(result, `"workItemType"`) {
		t.Fatalf("result missing workItemType: %q", result)
	}
	if !strings.Contains(result, `"fieldConfig"`) {
		t.Fatalf("result missing fieldConfig: %q", result)
	}
	if !strings.Contains(result, `"workflow"`) {
		t.Fatalf("result missing workflow: %q", result)
	}
}

func TestHandleGetWorkItemTypeOverviewSkipsSectionsWhenDisabled(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == "/oapi/v1/projex/organizations/org-1/workitemTypes/type-1" {
			_, _ = w.Write([]byte(`{"id":"type-1"}`))
			return
		}
		t.Fatalf("unexpected request to %q", r.URL.Path)
	})

	result, err := handleGetWorkItemTypeOverview(context.Background(), client, map[string]any{
		"organizationId":     "org-1",
		"projectId":          "project-1",
		"workItemTypeId":     "type-1",
		"includeFieldConfig": false,
		"includeWorkflow":    false,
	})
	if err != nil {
		t.Fatalf("handleGetWorkItemTypeOverview() error = %v", err)
	}
	if strings.Contains(result, `"fieldConfig"`) {
		t.Fatalf("result should not contain fieldConfig: %q", result)
	}
	if strings.Contains(result, `"workflow"`) {
		t.Fatalf("result should not contain workflow: %q", result)
	}
}

func TestWorkItemTypeOverviewFilters(t *testing.T) {
	params := map[string]any{
		"includeFieldConfig": false,
		"includeWorkflow":    false,
	}
	filters := workItemTypeOverviewFilters(params)
	if filters["includeFieldConfig"].(bool) != false {
		t.Fatalf("includeFieldConfig = %v", filters["includeFieldConfig"])
	}
	if filters["includeWorkflow"].(bool) != false {
		t.Fatalf("includeWorkflow = %v", filters["includeWorkflow"])
	}
}
