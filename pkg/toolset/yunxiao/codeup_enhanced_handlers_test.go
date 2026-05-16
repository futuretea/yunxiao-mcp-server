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
