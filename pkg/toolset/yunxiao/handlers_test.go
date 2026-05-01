package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newHandlerTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL, "token", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	return client
}

func TestHandleListRepositoriesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/codeup/organizations/org-1/repositories" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "10" ||
			r.URL.Query().Get("orderBy") != "name" ||
			r.URL.Query().Get("sort") != "asc" ||
			r.URL.Query().Get("search") != "demo" ||
			r.URL.Query().Get("archived") != "false" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"name":"repo"}]`))
	})

	result, err := handleListRepositories(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"page":           float64(2),
		"perPage":        float64(10),
		"orderBy":        "name",
		"sort":           "asc",
		"search":         "demo",
		"archived":       false,
	})
	if err != nil {
		t.Fatalf("handleListRepositories() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetRepositoryEncodesRepositoryPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/repositories/group%2FDemo%20Repo") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"name":"repo"}`))
	})

	if _, err := handleGetRepository(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/Demo Repo",
	}); err != nil {
		t.Fatalf("handleGetRepository() error = %v", err)
	}
}

func TestHandleListBranchesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/branches") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		if r.URL.Query().Get("sort") != "updated_desc" || r.URL.Query().Get("search") != "main" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[{"name":"main"}]`))
	})

	if _, err := handleListBranches(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"sort":           "updated_desc",
		"search":         "main",
	}); err != nil {
		t.Fatalf("handleListBranches() error = %v", err)
	}
}

func TestHandleListPipelinesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/flow/organizations/org-1/pipelines" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("pipelineName") != "deploy" ||
			r.URL.Query().Get("statusList") != "RUNNING,SUCCESS" ||
			r.URL.Query().Get("executeStartTime") != "1000" ||
			r.URL.Query().Get("perPage") != "30" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "2")
		_, _ = w.Write([]byte(`[{"pipelineId":1}]`))
	})

	result, err := handleListPipelines(context.Background(), client, map[string]any{
		"organizationId":   "org-1",
		"pipelineName":     "deploy",
		"statusList":       "RUNNING,SUCCESS",
		"executeStartTime": float64(1000),
		"executeEndTime":   float64(2000),
		"page":             float64(1),
		"perPage":          float64(30),
	})
	if err != nil {
		t.Fatalf("handleListPipelines() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetPipelineBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"pipelineId":1}`))
	})

	if _, err := handleGetPipeline(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
	}); err != nil {
		t.Fatalf("handleGetPipeline() error = %v", err)
	}
}

func TestHandleListPipelineRunsUsesEndTmeQueryName(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("startTime") != "1000" ||
			r.URL.Query().Get("endTme") != "2000" ||
			r.URL.Query().Get("endTime") != "" ||
			r.URL.Query().Get("triggerMode") != "1" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[{"pipelineRunId":1}]`))
	})

	if _, err := handleListPipelineRuns(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
		"startTime":      float64(1000),
		"endTime":        float64(2000),
		"status":         "SUCCESS",
		"triggerMode":    float64(1),
	}); err != nil {
		t.Fatalf("handleListPipelineRuns() error = %v", err)
	}
}

func TestHandleGetPipelineRunBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs/run-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"pipelineRunId":1}`))
	})

	if _, err := handleGetPipelineRun(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
		"pipelineRunId":  "run-1",
	}); err != nil {
		t.Fatalf("handleGetPipelineRun() error = %v", err)
	}
}

func TestHandleGetLatestPipelineRunBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/flow/organizations/org-1/pipelines/pipe-1/runs/latestPipelineRun" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"pipelineRunId":1}`))
	})

	if _, err := handleGetLatestPipelineRun(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
	}); err != nil {
		t.Fatalf("handleGetLatestPipelineRun() error = %v", err)
	}
}

func TestFlowHandlersRequirePipelineParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without required params")
	})

	_, err := handleGetPipeline(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil {
		t.Fatal("handleGetPipeline() expected missing pipelineId error")
	}

	_, err = handleGetPipelineRun(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe-1",
	})
	if err == nil {
		t.Fatal("handleGetPipelineRun() expected missing pipelineRunId error")
	}
}

func TestHandleSearchProjectsBuildsBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["orderBy"] != "name" || body["sort"] != "asc" || body["page"].(float64) != 2 {
			t.Fatalf("body = %#v", body)
		}
		conditions, _ := body["conditions"].(string)
		if !strings.Contains(conditions, `"fieldIdentifier":"name"`) ||
			!strings.Contains(conditions, `"value":["demo"]`) {
			t.Fatalf("conditions = %q", conditions)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"project-1"}]`))
	})

	result, err := handleSearchProjects(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"name":           "demo",
		"orderBy":        "name",
		"sort":           "asc",
		"page":           float64(2),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleSearchProjects() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleSearchProjectsAdvancedConditionsOverrideSimpleFilters(t *testing.T) {
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

	if _, err := handleSearchProjects(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"name":           "demo",
		"conditions":     advanced,
	}); err != nil {
		t.Fatalf("handleSearchProjects() error = %v", err)
	}
}

func TestHandleGetProjectBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"project-1"}`))
	})

	if _, err := handleGetProject(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
	}); err != nil {
		t.Fatalf("handleGetProject() error = %v", err)
	}
}

func TestHandleSearchWorkitemsBuildsBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["category"] != "Task" || body["spaceId"] != "project-1" || body["sort"] != "desc" {
			t.Fatalf("body = %#v", body)
		}
		conditions, _ := body["conditions"].(string)
		if !strings.Contains(conditions, `"fieldIdentifier":"subject"`) ||
			!strings.Contains(conditions, `"fieldIdentifier":"assignedTo"`) {
			t.Fatalf("conditions = %q", conditions)
		}
		var parsedConditions struct {
			ConditionGroups [][]map[string]any `json:"conditionGroups"`
		}
		if err := json.Unmarshal([]byte(conditions), &parsedConditions); err != nil {
			t.Fatalf("unmarshal conditions: %v", err)
		}
		assertConditionFormat(t, parsedConditions.ConditionGroups, "tag", "multiList")
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"workitem-1"}]`))
	})

	if _, err := handleSearchWorkitems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"category":       "Task",
		"spaceId":        "project-1",
		"subject":        "demo",
		"assignedTo":     "user-1,user-2",
		"tag":            "tag-1,tag-2",
		"sort":           "desc",
	}); err != nil {
		t.Fatalf("handleSearchWorkitems() error = %v", err)
	}
}

func assertConditionFormat(t *testing.T, conditionGroups [][]map[string]any, fieldIdentifier, format string) {
	t.Helper()

	for _, group := range conditionGroups {
		for _, condition := range group {
			if condition["fieldIdentifier"] == fieldIdentifier {
				if condition["format"] != format {
					t.Fatalf("%s format = %v, want %s", fieldIdentifier, condition["format"], format)
				}
				return
			}
		}
	}
	t.Fatalf("condition %s not found in %#v", fieldIdentifier, conditionGroups)
}

func TestHandleGetWorkitemBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems/workitem-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"workitem-1"}`))
	})

	if _, err := handleGetWorkitem(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "workitem-1",
	}); err != nil {
		t.Fatalf("handleGetWorkitem() error = %v", err)
	}
}

func TestSearchWorkitemsRequiresCategoryAndSpace(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without required params")
	})

	if _, err := handleSearchWorkitems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"spaceId":        "project-1",
	}); err == nil {
		t.Fatal("handleSearchWorkitems() expected missing category error")
	}
	if _, err := handleSearchWorkitems(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"category":       "Task",
	}); err == nil {
		t.Fatal("handleSearchWorkitems() expected missing spaceId error")
	}
}
