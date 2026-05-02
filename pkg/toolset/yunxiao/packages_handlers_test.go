package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListPackageRepositoriesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/packages/organizations/org-1/repositories" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("repoTypes") != "MAVEN,NPM" ||
			r.URL.Query().Get("repoCategories") != "Hybrid" ||
			r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "8" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"repoId":"repo-1"}]`))
	})

	result, err := handleListPackageRepositories(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repoTypes":      "MAVEN,NPM",
		"repoCategories": "Hybrid",
		"page":           float64(2),
		"perPage":        float64(8),
	})
	if err != nil {
		t.Fatalf("handleListPackageRepositories() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleListArtifactsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/packages/organizations/org-1/repositories/repo-1/artifacts" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("repoType") != "MAVEN" ||
			r.URL.Query().Get("page") != "3" ||
			r.URL.Query().Get("perPage") != "10" ||
			r.URL.Query().Get("search") != "junit" ||
			r.URL.Query().Get("orderBy") != "latestUpdate" ||
			r.URL.Query().Get("sort") != "desc" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":123}]`))
	})

	result, err := handleListArtifacts(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repoId":         "repo-1",
		"repoType":       "MAVEN",
		"page":           float64(3),
		"perPage":        float64(10),
		"search":         "junit",
		"orderBy":        "latestUpdate",
		"sort":           "desc",
	})
	if err != nil {
		t.Fatalf("handleListArtifacts() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetArtifactBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/packages/organizations/org-1/repositories/repo-1/artifacts/123" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("repoType") != "MAVEN" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"id":123}`))
	})

	if _, err := handleGetArtifact(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repoId":         "repo-1",
		"id":             float64(123),
		"repoType":       "MAVEN",
	}); err != nil {
		t.Fatalf("handleGetArtifact() error = %v", err)
	}
}

func TestRequiredNumberPathStringAcceptsTypes(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want string
	}{
		{"float64", float64(123), "123"},
		{"int", int(456), "456"},
		{"int64", int64(789), "789"},
		{"string", "abc", "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := requiredNumberPathString(map[string]any{"k": tt.val}, "k")
			if err != nil {
				t.Fatalf("requiredNumberPathString() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("requiredNumberPathString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRequiredNumberPathStringRejectsEmptyAndMissing(t *testing.T) {
	_, err := requiredNumberPathString(map[string]any{"k": ""}, "k")
	if err == nil {
		t.Fatal("expected error for empty string")
	}
	_, err = requiredNumberPathString(map[string]any{}, "k")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRequiredOrganizationAndPackageRepoRequiresRepoId(t *testing.T) {
	_, _, err := requiredOrganizationAndPackageRepo(map[string]any{"organizationId": "org-1"})
	if err == nil {
		t.Fatal("expected missing repoId error")
	}
}
