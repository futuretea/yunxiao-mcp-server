package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListTagsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/tags?") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		if r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "20" ||
			r.URL.Query().Get("search") != "v1" ||
			r.URL.Query().Get("sort") != "desc" ||
			r.URL.Query().Get("orderBy") != "create" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"name":"v1"}]`))
	})

	result, err := handleListTags(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"page":           float64(2),
		"perPage":        float64(20),
		"search":         "v1",
		"sort":           "desc",
		"orderBy":        "create",
	})
	if err != nil {
		t.Fatalf("handleListTags() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleListRepositoryMembersBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/repositories/group%2Frepo/members?accessLevel=40" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"userId":"u1"}]`))
	})

	if _, err := handleListRepositoryMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"accessLevel":    float64(40),
	}); err != nil {
		t.Fatalf("handleListRepositoryMembers() error = %v", err)
	}
}

func TestHandleListProtectedBranchesBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/repositories/group%2Frepo/protectedBranches" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"id":1}]`))
	})

	if _, err := handleListProtectedBranches(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
	}); err != nil {
		t.Fatalf("handleListProtectedBranches() error = %v", err)
	}
}

func TestHandleGetProtectedBranchBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/repositories/group%2Frepo/protectedBranches/12" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"id":12}`))
	})

	if _, err := handleGetProtectedBranch(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"id":             float64(12),
	}); err != nil {
		t.Fatalf("handleGetProtectedBranch() error = %v", err)
	}
}

func TestHandleListPushRulesBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/repositories/group%2Frepo/pushRules" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"id":1}]`))
	})

	if _, err := handleListPushRules(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
	}); err != nil {
		t.Fatalf("handleListPushRules() error = %v", err)
	}
}

func TestCodeupRepositoryMetadataHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleListTags(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleListRepositoryMembers(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleListProtectedBranches(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleGetProtectedBranch(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleListPushRules(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleGetPushRule(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
}

func TestHandleGetPushRuleBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/repositories/group%2Frepo/pushRules/7" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"id":7}`))
	})

	if _, err := handleGetPushRule(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"pushRuleId":     float64(7),
	}); err != nil {
		t.Fatalf("handleGetPushRule() error = %v", err)
	}
}
