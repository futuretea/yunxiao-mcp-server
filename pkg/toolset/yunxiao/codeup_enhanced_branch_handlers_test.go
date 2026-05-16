package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetBranchOverviewRequiresOrganizationId(t *testing.T) {
	_, err := handleGetBranchOverview(context.Background(), nil, map[string]any{
		"repositoryId": "repo-1",
		"branchName":   "main",
	})
	if err == nil || !strings.Contains(err.Error(), "organizationId is required") {
		t.Fatalf("expected organizationId required error, got %v", err)
	}
}

func TestHandleGetBranchOverviewRequiresRepositoryId(t *testing.T) {
	_, err := handleGetBranchOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"branchName":     "main",
	})
	if err == nil || !strings.Contains(err.Error(), "repositoryId is required") {
		t.Fatalf("expected repositoryId required error, got %v", err)
	}
}

func TestHandleGetBranchOverviewRequiresBranchName(t *testing.T) {
	_, err := handleGetBranchOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
	})
	if err == nil || !strings.Contains(err.Error(), "branchName is required") {
		t.Fatalf("expected branchName required error, got %v", err)
	}
}

func TestHandleGetBranchOverviewRequiresClient(t *testing.T) {
	_, err := handleGetBranchOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"branchName":     "main",
	})
	if err == nil || !strings.Contains(err.Error(), "yunxiao client is not configured") {
		t.Fatalf("expected client error, got %v", err)
	}
}

func TestHandleGetBranchOverviewReturnsErrorOnBranchFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleGetBranchOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"branchName":     "main",
	})
	if err == nil || !strings.Contains(err.Error(), "branch:") {
		t.Fatalf("expected branch error, got %v", err)
	}
}

func TestHandleGetBranchOverviewReturnsErrorOnCommitsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/commits") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"commits boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"name":"main"}`))
	})
	_, err := handleGetBranchOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"branchName":     "main",
	})
	if err == nil || !strings.Contains(err.Error(), "commits:") {
		t.Fatalf("expected commits error, got %v", err)
	}
}

func TestHandleGetBranchOverviewReturnsErrorOnMergeRequestsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/mergeRequests") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"mr boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"name":"main"}`))
	})
	_, err := handleGetBranchOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"branchName":     "main",
		"includeCommits": false,
	})
	if err == nil || !strings.Contains(err.Error(), "mergeRequests:") {
		t.Fatalf("expected mergeRequests error, got %v", err)
	}
}

func TestHandleGetBranchOverviewSuccessAllSections(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oapi/v1/codeup/organizations/org-1/repositories/repo-1/branches/main":
			_, _ = w.Write([]byte(`{"name":"main"}`))
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
			if r.URL.Query().Get("targetBranch") != "main" {
				t.Fatalf("mr targetBranch = %q", r.URL.Query().Get("targetBranch"))
			}
			if r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("mr perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[{"id":"mr-1"}]`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})

	result, err := handleGetBranchOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"branchName":     "main",
	})
	if err != nil {
		t.Fatalf("handleGetBranchOverview() error = %v", err)
	}
	if !strings.Contains(result, `"branch"`) {
		t.Fatalf("result missing branch: %q", result)
	}
	if !strings.Contains(result, `"commits"`) {
		t.Fatalf("result missing commits: %q", result)
	}
	if !strings.Contains(result, `"mergeRequests"`) {
		t.Fatalf("result missing mergeRequests: %q", result)
	}
}

func TestHandleGetBranchOverviewSkipsSectionsWhenDisabled(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/codeup/organizations/org-1/repositories/repo-1/branches/main" {
			_, _ = w.Write([]byte(`{"name":"main"}`))
			return
		}
		t.Fatalf("unexpected request to %q", r.URL.Path)
	})

	result, err := handleGetBranchOverview(context.Background(), client, map[string]any{
		"organizationId":       "org-1",
		"repositoryId":         "repo-1",
		"branchName":           "main",
		"includeCommits":       false,
		"includeMergeRequests": false,
	})
	if err != nil {
		t.Fatalf("handleGetBranchOverview() error = %v", err)
	}
	if strings.Contains(result, `"commits"`) {
		t.Fatalf("result should not contain commits: %q", result)
	}
	if strings.Contains(result, `"mergeRequests"`) {
		t.Fatalf("result should not contain mergeRequests: %q", result)
	}
}

func TestHandleGetBranchOverviewUsesCustomLimits(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/commits"):
			if r.URL.Query().Get("perPage") != "3" {
				t.Fatalf("commits perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[]`))
		case strings.HasSuffix(r.URL.Path, "/mergeRequests"):
			if r.URL.Query().Get("perPage") != "4" {
				t.Fatalf("mr perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[]`))
		default:
			_, _ = w.Write([]byte(`{"name":"main"}`))
		}
	})

	_, err := handleGetBranchOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"branchName":     "main",
		"commitLimit":    float64(3),
		"mrLimit":        float64(4),
	})
	if err != nil {
		t.Fatalf("handleGetBranchOverview() error = %v", err)
	}
}

func TestHandleGetBranchOverviewUsesCustomMrState(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/mergeRequests") {
			if r.URL.Query().Get("state") != "merged" {
				t.Fatalf("mr state = %q", r.URL.Query().Get("state"))
			}
			_, _ = w.Write([]byte(`[]`))
			return
		}
		_, _ = w.Write([]byte(`{"name":"main"}`))
	})

	_, err := handleGetBranchOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"branchName":     "main",
		"includeCommits": false,
		"mrState":        "merged",
	})
	if err != nil {
		t.Fatalf("handleGetBranchOverview() error = %v", err)
	}
}

func TestBranchOverviewFilters(t *testing.T) {
	params := map[string]any{
		"includeCommits":       false,
		"includeMergeRequests": false,
		"commitLimit":          float64(10),
		"mrLimit":              float64(20),
		"mrState":              "closed",
	}
	filters := branchOverviewFilters(params)
	if filters["includeCommits"].(bool) != false {
		t.Fatalf("includeCommits = %v", filters["includeCommits"])
	}
	if filters["includeMergeRequests"].(bool) != false {
		t.Fatalf("includeMergeRequests = %v", filters["includeMergeRequests"])
	}
	if filters["commitLimit"].(int) != 10 {
		t.Fatalf("commitLimit = %v", filters["commitLimit"])
	}
	if filters["mrLimit"].(int) != 20 {
		t.Fatalf("mrLimit = %v", filters["mrLimit"])
	}
	if filters["mrState"].(string) != "closed" {
		t.Fatalf("mrState = %q", filters["mrState"])
	}
}
