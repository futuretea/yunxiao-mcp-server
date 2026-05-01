package yunxiao

import (
	"context"
	"encoding/json"
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

func TestHandleListTestcaseRepositoriesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/repo/list" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "2" || r.URL.Query().Get("perPage") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"repo-1"}]`))
	})

	result, err := handleListTestcaseRepositories(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"page":           float64(2),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleListTestcaseRepositories() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
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

func TestHandleSearchTestcasesBuildsBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/testRepos/repo-1/testcases:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["directoryId"] != "dir-1" || body["orderBy"] != "name" || body["sort"] != "asc" || body["page"].(float64) != 2 {
			t.Fatalf("body = %#v", body)
		}
		conditions, _ := body["conditions"].(string)
		if !strings.Contains(conditions, `"fieldIdentifier":"subject"`) ||
			!strings.Contains(conditions, `"value":["login"]`) {
			t.Fatalf("conditions = %q", conditions)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"testcase-1"}]`))
	})

	result, err := handleSearchTestcases(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"testRepoId":     "repo-1",
		"directoryId":    "dir-1",
		"subject":        "login",
		"orderBy":        "name",
		"sort":           "asc",
		"page":           float64(2),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleSearchTestcases() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleSearchTestcasesAdvancedConditionsOverrideSimpleFilters(t *testing.T) {
	const advanced = `{"conditionGroups":[[{"fieldIdentifier":"custom"}]]}`
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["conditions"] != advanced {
			t.Fatalf("conditions = %#v", body["conditions"])
		}
		_, _ = w.Write([]byte(`[]`))
	})

	if _, err := handleSearchTestcases(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"testRepoId":     "repo-1",
		"subject":        "login",
		"conditions":     advanced,
	}); err != nil {
		t.Fatalf("handleSearchTestcases() error = %v", err)
	}
}

func TestHandleListTestPlansBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/testPlan/list" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"id":"plan-1"}]`))
	})

	if _, err := handleListTestPlans(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	}); err != nil {
		t.Fatalf("handleListTestPlans() error = %v", err)
	}
}

func TestHandleGetTestResultListBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/plan%2F1/result/list/dir%2F1" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"testcaseId":"testcase-1"}]`))
	})

	if _, err := handleGetTestResultList(context.Background(), client, map[string]any{
		"organizationId":      "org-1",
		"testPlanIdentifier":  "plan/1",
		"directoryIdentifier": "dir/1",
	}); err != nil {
		t.Fatalf("handleGetTestResultList() error = %v", err)
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
