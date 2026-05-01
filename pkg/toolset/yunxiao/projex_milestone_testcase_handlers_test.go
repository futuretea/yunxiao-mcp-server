package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListMilestonesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/milestones" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("status") != "TODO,DOING" ||
			r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"milestone-1"}]`))
	})

	result, err := handleListMilestones(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"status":         "TODO,DOING",
		"page":           float64(2),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleListMilestones() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleListDirectoriesBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/testRepos/repo-1/directories" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"id":"directory-1"}]`))
	})

	if _, err := handleListDirectories(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "repo-1",
	}); err != nil {
		t.Fatalf("handleListDirectories() error = %v", err)
	}
}

func TestHandleGetTestcaseFieldConfigBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/testRepos/repo-1/testcases/fields" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"fieldIdentifier":"name"}]`))
	})

	if _, err := handleGetTestcaseFieldConfig(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "repo-1",
	}); err != nil {
		t.Fatalf("handleGetTestcaseFieldConfig() error = %v", err)
	}
}

func TestHandleGetTestcaseBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/testRepos/repo-1/testcases/testcase-1" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"id":"testcase-1"}`))
	})

	if _, err := handleGetTestcase(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"testRepoId":     "repo-1",
		"id":             "testcase-1",
	}); err != nil {
		t.Fatalf("handleGetTestcase() error = %v", err)
	}
}
