package yunxiao

import (
	"context"
	"net/http"
	"testing"
)

func TestHandleListProjectMembersBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/members" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("name") != "demo" ||
			r.URL.Query().Get("roleId") != "project.admin" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[{"userId":"u1"}]`))
	})

	if _, err := handleListProjectMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"name":           "demo",
		"roleId":         "project.admin",
	}); err != nil {
		t.Fatalf("handleListProjectMembers() error = %v", err)
	}
}

func TestHandleListProjectTemplatesBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/projectTemplates" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"id":"template-1"}]`))
	})

	if _, err := handleListProjectTemplates(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	}); err != nil {
		t.Fatalf("handleListProjectTemplates() error = %v", err)
	}
}

func TestHandleGetProjectTemplateFieldConfigBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/projectTemplates/template-1/fields" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"fieldIdentifier":"name"}]`))
	})

	if _, err := handleGetProjectTemplateFieldConfig(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "template-1",
	}); err != nil {
		t.Fatalf("handleGetProjectTemplateFieldConfig() error = %v", err)
	}
}

func TestHandleListProjectProgramBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/program-1/binding/project/list" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"id":"project-1"}]`))
	})

	if _, err := handleListProjectProgram(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"programIdentifier": "program-1",
	}); err != nil {
		t.Fatalf("handleListProjectProgram() error = %v", err)
	}
}

func TestHandleListProjectRolesBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/projects/project-1/roles" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"id":"project.admin"}]`))
	})

	if _, err := handleListProjectRoles(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
	}); err != nil {
		t.Fatalf("handleListProjectRoles() error = %v", err)
	}
}

func TestHandleListAllProjectRolesBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/roles" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"id":"project.admin"}]`))
	})

	if _, err := handleListAllProjectRoles(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	}); err != nil {
		t.Fatalf("handleListAllProjectRoles() error = %v", err)
	}
}
