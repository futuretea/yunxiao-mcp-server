package yunxiao

import (
	"context"
	"net/http"
	"testing"
)

func TestHandleListAllWorkItemTypesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitemTypes" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("categories") != "Req,Bug" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[{"id":"type-1"}]`))
	})

	if _, err := handleListAllWorkItemTypes(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"categories":     "Req,Bug",
	}); err != nil {
		t.Fatalf("handleListAllWorkItemTypes() error = %v", err)
	}
}

func TestHandleListWorkItemTypesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/workitemTypes" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("category") != "Task" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[{"id":"type-1"}]`))
	})

	if _, err := handleListWorkItemTypes(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"category":       "Task",
	}); err != nil {
		t.Fatalf("handleListWorkItemTypes() error = %v", err)
	}
}

func TestHandleGetWorkItemTypeBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitemTypes/type-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"type-1"}`))
	})

	if _, err := handleGetWorkItemType(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "type-1",
	}); err != nil {
		t.Fatalf("handleGetWorkItemType() error = %v", err)
	}
}

func TestHandleListWorkItemRelationWorkItemTypesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitemTypes/type-1/relationWorkitemTypes" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("relationType") != "PARENT" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[{"id":"type-2"}]`))
	})

	if _, err := handleListWorkItemRelationWorkItemTypes(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workItemTypeId": "type-1",
		"relationType":   "PARENT",
	}); err != nil {
		t.Fatalf("handleListWorkItemRelationWorkItemTypes() error = %v", err)
	}
}

func TestHandleGetWorkItemTypeFieldConfigBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/workitemTypes/type-1/fields" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"fieldId":"subject"}]`))
	})

	if _, err := handleGetWorkItemTypeFieldConfig(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"workItemTypeId": "type-1",
	}); err != nil {
		t.Fatalf("handleGetWorkItemTypeFieldConfig() error = %v", err)
	}
}

func TestHandleGetWorkItemWorkflowBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/workitemTypes/type-1/workflows" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"workflow-1"}`))
	})

	if _, err := handleGetWorkItemWorkflow(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"workItemTypeId": "type-1",
	}); err != nil {
		t.Fatalf("handleGetWorkItemWorkflow() error = %v", err)
	}
}

func TestHandleListAllWorkItemTypesRequiresOrganizationId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleListAllWorkItemTypes(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
}

func TestHandleListWorkItemTypesRequiresCategory(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleListWorkItemTypes(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
	}); err == nil {
		t.Fatal("expected missing category error")
	}
}

func TestHandleGetWorkItemTypeRequiresID(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleGetWorkItemType(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing id error")
	}
}

func TestHandleListWorkItemRelationWorkItemTypesRequiresWorkItemTypeId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleListWorkItemRelationWorkItemTypes(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing workItemTypeId error")
	}
}

func TestHandleGetWorkItemTypeFieldConfigRequiresWorkItemTypeId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleGetWorkItemTypeFieldConfig(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
	}); err == nil {
		t.Fatal("expected missing workItemTypeId error")
	}
}

func TestHandleGetWorkItemWorkflowRequiresWorkItemTypeId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleGetWorkItemWorkflow(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
	}); err == nil {
		t.Fatal("expected missing workItemTypeId error")
	}
}

func TestRequiredOrganizationAndProjectRequiresProjectId(t *testing.T) {
	_, _, err := requiredOrganizationAndProject(map[string]any{"organizationId": "org-1"})
	if err == nil {
		t.Fatal("expected missing projectId error")
	}
}

func TestRequiredOrganizationProjectAndWorkItemTypeRequiresWorkItemTypeId(t *testing.T) {
	_, _, _, err := requiredOrganizationProjectAndWorkItemType(map[string]any{"organizationId": "org-1", "projectId": "project-1"})
	if err == nil {
		t.Fatal("expected missing workItemTypeId error")
	}
}

func TestProjexWorkitemTypeHandlersRequireParams(t *testing.T) {
	if _, err := handleListAllWorkItemTypes(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListWorkItemTypes(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "projectId": "project-1", "category": "Task"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetWorkItemType(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "id": "type-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListWorkItemRelationWorkItemTypes(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "workItemTypeId": "type-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetWorkItemTypeFieldConfig(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "projectId": "project-1", "workItemTypeId": "type-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetWorkItemWorkflow(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "projectId": "project-1", "workItemTypeId": "type-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}
