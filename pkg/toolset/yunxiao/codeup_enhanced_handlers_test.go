package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestPageOneLimitQuery(t *testing.T) {
	q := pageOneLimitQuery(map[string]any{"branchLimit": float64(7)}, "branchLimit", 5)
	if q.Get("page") != "1" {
		t.Fatalf("page = %q", q.Get("page"))
	}
	if q.Get("perPage") != "7" {
		t.Fatalf("perPage = %q", q.Get("perPage"))
	}

	q2 := pageOneLimitQuery(map[string]any{}, "branchLimit", 5)
	if q2.Get("perPage") != "5" {
		t.Fatalf("default perPage = %q", q2.Get("perPage"))
	}
}

func TestHandleGetChangeRequestOverviewRequiresOrganizationId(t *testing.T) {
	_, err := handleGetChangeRequestOverview(context.Background(), nil, map[string]any{
		"repositoryId": "repo-1",
		"localId":      "1",
	})
	if err == nil || !strings.Contains(err.Error(), "organizationId is required") {
		t.Fatalf("expected organizationId required error, got %v", err)
	}
}

func TestHandleGetChangeRequestOverviewRequiresRepositoryId(t *testing.T) {
	_, err := handleGetChangeRequestOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"localId":        "1",
	})
	if err == nil || !strings.Contains(err.Error(), "repositoryId is required") {
		t.Fatalf("expected repositoryId required error, got %v", err)
	}
}

func TestHandleGetChangeRequestOverviewRequiresLocalId(t *testing.T) {
	_, err := handleGetChangeRequestOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
	})
	if err == nil || !strings.Contains(err.Error(), "localId is required") {
		t.Fatalf("expected localId required error, got %v", err)
	}
}

func TestHandleGetChangeRequestOverviewRequiresClient(t *testing.T) {
	_, err := handleGetChangeRequestOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
	})
	if err == nil || !strings.Contains(err.Error(), "yunxiao client is not configured") {
		t.Fatalf("expected client error, got %v", err)
	}
}

func TestHandleGetChangeRequestOverviewReturnsErrorOnCRFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleGetChangeRequestOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetChangeRequestOverviewReturnsErrorOnPatchSetsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/diffs/patches") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"patches boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"cr-1"}`))
	})
	_, err := handleGetChangeRequestOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetChangeRequestOverviewReturnsErrorOnCommentsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/comments/list") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"comments boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"cr-1"}`))
	})
	_, err := handleGetChangeRequestOverview(context.Background(), client, map[string]any{
		"organizationId":   "org-1",
		"repositoryId":     "repo-1",
		"localId":          "1",
		"includePatchSets": false,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetChangeRequestOverviewSuccessAllSections(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oapi/v1/codeup/organizations/org-1/repositories/repo-1/changeRequests/1":
			_, _ = w.Write([]byte(`{"id":"cr-1"}`))
		case strings.HasSuffix(r.URL.Path, "/diffs/patches"):
			_, _ = w.Write([]byte(`["patch-1"]`))
		case strings.HasSuffix(r.URL.Path, "/comments/list"):
			_, _ = w.Write([]byte(`["comment-1"]`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})

	result, err := handleGetChangeRequestOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
	})
	if err != nil {
		t.Fatalf("handleGetChangeRequestOverview() error = %v", err)
	}
	if !strings.Contains(result, `"changeRequest"`) {
		t.Fatalf("result missing changeRequest: %q", result)
	}
	if !strings.Contains(result, `"patchSets"`) {
		t.Fatalf("result missing patchSets: %q", result)
	}
	if !strings.Contains(result, `"comments"`) {
		t.Fatalf("result missing comments: %q", result)
	}
}

func TestHandleGetChangeRequestOverviewSkipsSectionsWhenDisabled(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/codeup/organizations/org-1/repositories/repo-1/changeRequests/1" {
			_, _ = w.Write([]byte(`{"id":"cr-1"}`))
			return
		}
		t.Fatalf("unexpected request to %q", r.URL.Path)
	})

	result, err := handleGetChangeRequestOverview(context.Background(), client, map[string]any{
		"organizationId":   "org-1",
		"repositoryId":     "repo-1",
		"localId":          "1",
		"includePatchSets": false,
		"includeComments":  false,
	})
	if err != nil {
		t.Fatalf("handleGetChangeRequestOverview() error = %v", err)
	}
	if strings.Contains(result, `"patchSets"`) {
		t.Fatalf("result should not contain patchSets: %q", result)
	}
	if strings.Contains(result, `"comments"`) {
		t.Fatalf("result should not contain comments: %q", result)
	}
}

func TestChangeRequestOverviewFilters(t *testing.T) {
	params := map[string]any{
		"includePatchSets": false,
		"includeComments":  false,
		"commentState":     "RESOLVED",
		"commentResolved":  true,
	}
	filters := changeRequestOverviewFilters(params)
	if filters["includePatchSets"].(bool) != false {
		t.Fatalf("includePatchSets = %v", filters["includePatchSets"])
	}
	if filters["includeComments"].(bool) != false {
		t.Fatalf("includeComments = %v", filters["includeComments"])
	}
	if filters["commentState"] != "RESOLVED" {
		t.Fatalf("commentState = %v", filters["commentState"])
	}
	if filters["commentResolved"].(bool) != true {
		t.Fatalf("commentResolved = %v", filters["commentResolved"])
	}
}

func TestHandleGetCommitOverviewRequiresOrganizationId(t *testing.T) {
	_, err := handleGetCommitOverview(context.Background(), nil, map[string]any{
		"repositoryId": "repo-1",
		"sha":          "abc123",
	})
	if err == nil || !strings.Contains(err.Error(), "organizationId is required") {
		t.Fatalf("expected organizationId required error, got %v", err)
	}
}

func TestHandleGetCommitOverviewRequiresRepositoryId(t *testing.T) {
	_, err := handleGetCommitOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"sha":            "abc123",
	})
	if err == nil || !strings.Contains(err.Error(), "repositoryId is required") {
		t.Fatalf("expected repositoryId required error, got %v", err)
	}
}

func TestHandleGetCommitOverviewRequiresSha(t *testing.T) {
	_, err := handleGetCommitOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
	})
	if err == nil || !strings.Contains(err.Error(), "sha is required") {
		t.Fatalf("expected sha required error, got %v", err)
	}
}

func TestHandleGetCommitOverviewRequiresClient(t *testing.T) {
	_, err := handleGetCommitOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"sha":            "abc123",
	})
	if err == nil || !strings.Contains(err.Error(), "yunxiao client is not configured") {
		t.Fatalf("expected client error, got %v", err)
	}
}

func TestHandleGetCommitOverviewReturnsErrorOnCommitFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleGetCommitOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"sha":            "abc123",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetCommitOverviewReturnsErrorOnStatusesFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/statuses") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"statuses boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"sha":"abc123"}`))
	})
	_, err := handleGetCommitOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"sha":            "abc123",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetCommitOverviewReturnsErrorOnCheckRunsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/checkRuns") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"checkRuns boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"sha":"abc123"}`))
	})
	_, err := handleGetCommitOverview(context.Background(), client, map[string]any{
		"organizationId":  "org-1",
		"repositoryId":    "repo-1",
		"sha":             "abc123",
		"includeStatuses": false,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetCommitOverviewSuccessAllSections(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oapi/v1/codeup/organizations/org-1/repositories/repo-1/commits/abc123":
			_, _ = w.Write([]byte(`{"sha":"abc123"}`))
		case strings.HasSuffix(r.URL.Path, "/statuses"):
			if r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("statuses perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`["success"]`))
		case strings.HasSuffix(r.URL.Path, "/checkRuns"):
			if r.URL.Query().Get("ref") != "abc123" {
				t.Fatalf("checkRuns ref = %q", r.URL.Query().Get("ref"))
			}
			if r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("checkRuns perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`["check-1"]`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})

	result, err := handleGetCommitOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"sha":            "abc123",
	})
	if err != nil {
		t.Fatalf("handleGetCommitOverview() error = %v", err)
	}
	if !strings.Contains(result, `"commit"`) {
		t.Fatalf("result missing commit: %q", result)
	}
	if !strings.Contains(result, `"statuses"`) {
		t.Fatalf("result missing statuses: %q", result)
	}
	if !strings.Contains(result, `"checkRuns"`) {
		t.Fatalf("result missing checkRuns: %q", result)
	}
}

func TestHandleGetCommitOverviewSkipsSectionsWhenDisabled(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/codeup/organizations/org-1/repositories/repo-1/commits/abc123" {
			_, _ = w.Write([]byte(`{"sha":"abc123"}`))
			return
		}
		t.Fatalf("unexpected request to %q", r.URL.Path)
	})

	result, err := handleGetCommitOverview(context.Background(), client, map[string]any{
		"organizationId":   "org-1",
		"repositoryId":     "repo-1",
		"sha":              "abc123",
		"includeStatuses":  false,
		"includeCheckRuns": false,
	})
	if err != nil {
		t.Fatalf("handleGetCommitOverview() error = %v", err)
	}
	if strings.Contains(result, `"statuses"`) {
		t.Fatalf("result should not contain statuses: %q", result)
	}
	if strings.Contains(result, `"checkRuns"`) {
		t.Fatalf("result should not contain checkRuns: %q", result)
	}
}

func TestHandleGetCommitOverviewUsesCustomLimits(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/statuses"):
			if r.URL.Query().Get("perPage") != "3" {
				t.Fatalf("statuses perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[]`))
		case strings.HasSuffix(r.URL.Path, "/checkRuns"):
			if r.URL.Query().Get("perPage") != "2" {
				t.Fatalf("checkRuns perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[]`))
		default:
			_, _ = w.Write([]byte(`{"sha":"abc123"}`))
		}
	})

	_, err := handleGetCommitOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"sha":            "abc123",
		"statusLimit":    float64(3),
		"checkRunLimit":  float64(2),
	})
	if err != nil {
		t.Fatalf("handleGetCommitOverview() error = %v", err)
	}
}

func TestCommitOverviewFilters(t *testing.T) {
	params := map[string]any{
		"includeStatuses":  false,
		"includeCheckRuns": false,
		"statusLimit":      float64(10),
		"checkRunLimit":    float64(20),
	}
	filters := commitOverviewFilters(params)
	if filters["includeStatuses"].(bool) != false {
		t.Fatalf("includeStatuses = %v", filters["includeStatuses"])
	}
	if filters["includeCheckRuns"].(bool) != false {
		t.Fatalf("includeCheckRuns = %v", filters["includeCheckRuns"])
	}
	if filters["statusLimit"].(int) != 10 {
		t.Fatalf("statusLimit = %v", filters["statusLimit"])
	}
	if filters["checkRunLimit"].(int) != 20 {
		t.Fatalf("checkRunLimit = %v", filters["checkRunLimit"])
	}
}
