package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListOrganizationGroupsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/platform/organizations/org-1/groups" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "2" || r.URL.Query().Get("perPage") != "50" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"group-1"}]`))
	})

	result, err := handleListOrganizationGroups(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"page":           float64(2),
		"perPage":        float64(50),
	})
	if err != nil {
		t.Fatalf("handleListOrganizationGroups() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetOrganizationGroupBuildsEscapedPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/platform/organizations/org-1/groups/group%2F1" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"id":"group-1"}`))
	})

	if _, err := handleGetOrganizationGroup(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "group/1",
	}); err != nil {
		t.Fatalf("handleGetOrganizationGroup() error = %v", err)
	}
}

func TestHandleListOrganizationGroupMembersPreservesEncodedID(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/platform/organizations/org-1/groups/group%2F1/members?page=1&perPage=20" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"member-1"}]`))
	})

	result, err := handleListOrganizationGroupMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "group%2F1",
		"page":           float64(1),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleListOrganizationGroupMembers() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestPlatformGroupHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleListOrganizationGroups(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListOrganizationGroups(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetOrganizationGroup(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	_, err := handleGetOrganizationGroup(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil {
		t.Fatal("handleGetOrganizationGroup() expected missing id error")
	}
	if _, err := handleGetOrganizationGroup(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "id": "g-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListOrganizationGroupMembers(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListOrganizationGroupMembers(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing id error")
	}
	if _, err := handleListOrganizationGroupMembers(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "id": "g-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}
