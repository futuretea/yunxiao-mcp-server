package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetRepositoryOverviewRequiresOrganizationId(t *testing.T) {
	_, err := handleGetRepositoryOverview(context.Background(), nil, map[string]any{
		"repositoryId": "repo-1",
	})
	if err == nil || !strings.Contains(err.Error(), "organizationId is required") {
		t.Fatalf("expected organizationId required error, got %v", err)
	}
}

func TestHandleGetRepositoryOverviewRequiresRepositoryId(t *testing.T) {
	_, err := handleGetRepositoryOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil || !strings.Contains(err.Error(), "repositoryId is required") {
		t.Fatalf("expected repositoryId required error, got %v", err)
	}
}

func TestHandleGetRepositoryOverviewRequiresClient(t *testing.T) {
	_, err := handleGetRepositoryOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
	})
	if err == nil || !strings.Contains(err.Error(), "yunxiao client is not configured") {
		t.Fatalf("expected client error, got %v", err)
	}
}

func TestHandleGetRepositoryOverviewReturnsErrorOnRepositoryFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleGetRepositoryOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
	})
	if err == nil || !strings.Contains(err.Error(), "repository:") {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestHandleGetRepositoryOverviewReturnsErrorOnBranchesFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/branches") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"branches boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"repo-1","defaultBranch":"main"}`))
	})
	_, err := handleGetRepositoryOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
	})
	if err == nil || !strings.Contains(err.Error(), "branches:") {
		t.Fatalf("expected branches error, got %v", err)
	}
}

func TestHandleGetRepositoryOverviewReturnsErrorOnCommitsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/commits") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"commits boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"repo-1","defaultBranch":"main"}`))
	})
	_, err := handleGetRepositoryOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
	})
	if err == nil || !strings.Contains(err.Error(), "commits:") {
		t.Fatalf("expected commits error, got %v", err)
	}
}

func TestHandleGetRepositoryOverviewReturnsErrorOnMergeRequestsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/mergeRequests") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"mr boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"repo-1","defaultBranch":"main"}`))
	})
	_, err := handleGetRepositoryOverview(context.Background(), client, map[string]any{
		"organizationId":  "org-1",
		"repositoryId":    "repo-1",
		"includeBranches": false,
		"includeCommits":  false,
	})
	if err == nil || !strings.Contains(err.Error(), "mergeRequests:") {
		t.Fatalf("expected mergeRequests error, got %v", err)
	}
}

func TestHandleGetRepositoryOverviewSuccessAllSections(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oapi/v1/codeup/organizations/org-1/repositories/repo-1":
			_, _ = w.Write([]byte(`{"id":"repo-1","defaultBranch":"main"}`))
		case strings.HasSuffix(r.URL.Path, "/branches"):
			if r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("branches perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`["main","dev"]`))
		case strings.HasSuffix(r.URL.Path, "/commits"):
			if r.URL.Query().Get("refName") != "main" {
				t.Fatalf("commits refName = %q", r.URL.Query().Get("refName"))
			}
			if r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("commits perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[{"sha":"abc123"}]`))
		case strings.HasSuffix(r.URL.Path, "/mergeRequests"):
			if r.URL.Query().Get("state") != "opened" {
				t.Fatalf("mr state = %q", r.URL.Query().Get("state"))
			}
			if r.URL.Query().Get("repositoryIds") != "repo-1" {
				t.Fatalf("mr repositoryIds = %q", r.URL.Query().Get("repositoryIds"))
			}
			if r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("mr perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[{"id":"mr-1"}]`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})

	result, err := handleGetRepositoryOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
	})
	if err != nil {
		t.Fatalf("handleGetRepositoryOverview() error = %v", err)
	}
	if !strings.Contains(result, `"repository"`) {
		t.Fatalf("result missing repository: %q", result)
	}
	if !strings.Contains(result, `"branches"`) {
		t.Fatalf("result missing branches: %q", result)
	}
	if !strings.Contains(result, `"commits"`) {
		t.Fatalf("result missing commits: %q", result)
	}
	if !strings.Contains(result, `"mergeRequests"`) {
		t.Fatalf("result missing mergeRequests: %q", result)
	}
}

func TestHandleGetRepositoryOverviewSkipsSectionsWhenDisabled(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/codeup/organizations/org-1/repositories/repo-1" {
			_, _ = w.Write([]byte(`{"id":"repo-1","defaultBranch":"main"}`))
			return
		}
		t.Fatalf("unexpected request to %q", r.URL.Path)
	})

	result, err := handleGetRepositoryOverview(context.Background(), client, map[string]any{
		"organizationId":       "org-1",
		"repositoryId":         "repo-1",
		"includeBranches":      false,
		"includeCommits":       false,
		"includeMergeRequests": false,
	})
	if err != nil {
		t.Fatalf("handleGetRepositoryOverview() error = %v", err)
	}
	if strings.Contains(result, `"branches"`) {
		t.Fatalf("result should not contain branches: %q", result)
	}
	if strings.Contains(result, `"commits"`) {
		t.Fatalf("result should not contain commits: %q", result)
	}
	if strings.Contains(result, `"mergeRequests"`) {
		t.Fatalf("result should not contain mergeRequests: %q", result)
	}
}

func TestHandleGetRepositoryOverviewUsesCustomRefName(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/commits") {
			if r.URL.Query().Get("refName") != "feature-x" {
				t.Fatalf("commits refName = %q", r.URL.Query().Get("refName"))
			}
			_, _ = w.Write([]byte(`[{"sha":"def456"}]`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"repo-1","defaultBranch":"main"}`))
	})

	result, err := handleGetRepositoryOverview(context.Background(), client, map[string]any{
		"organizationId":       "org-1",
		"repositoryId":         "repo-1",
		"includeBranches":      false,
		"refName":              "feature-x",
		"includeMergeRequests": false,
	})
	if err != nil {
		t.Fatalf("handleGetRepositoryOverview() error = %v", err)
	}
	if !strings.Contains(result, `"commits"`) {
		t.Fatalf("result missing commits: %q", result)
	}
}

func TestHandleGetRepositoryOverviewSkipsCommitsWhenNoDefaultBranch(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/codeup/organizations/org-1/repositories/repo-1" {
			_, _ = w.Write([]byte(`{"id":"repo-1"}`))
			return
		}
		t.Fatalf("unexpected request to %q", r.URL.Path)
	})

	result, err := handleGetRepositoryOverview(context.Background(), client, map[string]any{
		"organizationId":       "org-1",
		"repositoryId":         "repo-1",
		"includeBranches":      false,
		"includeMergeRequests": false,
	})
	if err != nil {
		t.Fatalf("handleGetRepositoryOverview() error = %v", err)
	}
	if strings.Contains(result, `"commits"`) {
		t.Fatalf("result should not contain commits when no defaultBranch: %q", result)
	}
}

func TestHandleGetRepositoryOverviewUsesCustomLimits(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/branches"):
			if r.URL.Query().Get("perPage") != "3" {
				t.Fatalf("branches perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[]`))
		case strings.HasSuffix(r.URL.Path, "/commits"):
			if r.URL.Query().Get("perPage") != "2" {
				t.Fatalf("commits perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[]`))
		case strings.HasSuffix(r.URL.Path, "/mergeRequests"):
			if r.URL.Query().Get("perPage") != "4" {
				t.Fatalf("mr perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[]`))
		default:
			_, _ = w.Write([]byte(`{"id":"repo-1","defaultBranch":"main"}`))
		}
	})

	_, err := handleGetRepositoryOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"branchLimit":    float64(3),
		"commitLimit":    float64(2),
		"mrLimit":        float64(4),
	})
	if err != nil {
		t.Fatalf("handleGetRepositoryOverview() error = %v", err)
	}
}

func TestHandleGetRepositoryOverviewUsesCustomMrState(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/mergeRequests") {
			if r.URL.Query().Get("state") != "merged" {
				t.Fatalf("mr state = %q", r.URL.Query().Get("state"))
			}
			_, _ = w.Write([]byte(`[]`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"repo-1","defaultBranch":"main"}`))
	})

	_, err := handleGetRepositoryOverview(context.Background(), client, map[string]any{
		"organizationId":  "org-1",
		"repositoryId":    "repo-1",
		"includeBranches": false,
		"includeCommits":  false,
		"mrState":         "merged",
	})
	if err != nil {
		t.Fatalf("handleGetRepositoryOverview() error = %v", err)
	}
}

func TestRepositoryOverviewFilters(t *testing.T) {
	params := map[string]any{
		"includeBranches":      false,
		"includeCommits":       false,
		"includeMergeRequests": false,
		"refName":              "dev",
		"branchLimit":          float64(10),
		"commitLimit":          float64(20),
		"mrLimit":              float64(30),
		"mrState":              "closed",
	}
	filters := repositoryOverviewFilters(params)
	if filters["includeBranches"].(bool) != false {
		t.Fatalf("includeBranches = %v", filters["includeBranches"])
	}
	if filters["refName"].(string) != "dev" {
		t.Fatalf("refName = %q", filters["refName"])
	}
	if filters["branchLimit"].(int) != 10 {
		t.Fatalf("branchLimit = %v", filters["branchLimit"])
	}
	if filters["mrState"].(string) != "closed" {
		t.Fatalf("mrState = %q", filters["mrState"])
	}
}

func TestRepositoryLimitQuery(t *testing.T) {
	q := repositoryLimitQuery(map[string]any{"branchLimit": float64(7)}, "branchLimit", 5)
	if q.Get("page") != "1" {
		t.Fatalf("page = %q", q.Get("page"))
	}
	if q.Get("perPage") != "7" {
		t.Fatalf("perPage = %q", q.Get("perPage"))
	}

	q2 := repositoryLimitQuery(map[string]any{}, "branchLimit", 5)
	if q2.Get("perPage") != "5" {
		t.Fatalf("default perPage = %q", q2.Get("perPage"))
	}
}

func TestHandleGetRepositoryOverviewDefaultBranchStringArray(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/codeup/organizations/org-1/repositories/repo-1" {
			_, _ = w.Write([]byte(`["not-a-map"]`))
			return
		}
		t.Fatalf("unexpected request to %q", r.URL.Path)
	})

	result, err := handleGetRepositoryOverview(context.Background(), client, map[string]any{
		"organizationId":       "org-1",
		"repositoryId":         "repo-1",
		"includeBranches":      false,
		"includeMergeRequests": false,
	})
	if err != nil {
		t.Fatalf("handleGetRepositoryOverview() error = %v", err)
	}
	if strings.Contains(result, `"commits"`) {
		t.Fatalf("result should not contain commits when repository is not a map: %q", result)
	}
}

func TestHandleGetRepositoryOverviewWithRepositoryPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/codeup/organizations/org-1/repositories/org/repo" {
			_, _ = w.Write([]byte(`{"id":"org/repo","defaultBranch":"main"}`))
			return
		}
		if strings.HasSuffix(r.URL.Path, "/branches") {
			_, _ = w.Write([]byte(`[]`))
			return
		}
		if strings.HasSuffix(r.URL.Path, "/commits") {
			_, _ = w.Write([]byte(`[]`))
			return
		}
		if strings.HasSuffix(r.URL.Path, "/mergeRequests") {
			_, _ = w.Write([]byte(`[]`))
			return
		}
		t.Fatalf("unexpected request to %q", r.URL.Path)
	})

	result, err := handleGetRepositoryOverview(context.Background(), client, map[string]any{
		"organizationId":       "org-1",
		"repositoryId":         "org/repo",
		"includeBranches":      true,
		"includeCommits":       true,
		"includeMergeRequests": true,
	})
	if err != nil {
		t.Fatalf("handleGetRepositoryOverview() error = %v", err)
	}
	if !strings.Contains(result, `"repository"`) {
		t.Fatalf("result missing repository: %q", result)
	}
}
