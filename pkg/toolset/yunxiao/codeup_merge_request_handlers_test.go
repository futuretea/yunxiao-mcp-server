package yunxiao

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestHandleListMergeRequestsBuildsExplodedArrayQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/codeup/organizations/org-1/mergeRequests" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		query := r.URL.Query()
		if query.Get("page") != "2" ||
			query.Get("perPage") != "20" ||
			query.Get("state") != "opened" ||
			query.Get("search") != "demo" ||
			query.Get("orderBy") != "updated_at" ||
			query.Get("createdAfter") != "2026-01-01" ||
			query.Get("createdBefore") != "2026-02-01" ||
			query.Get("targetBranch") != "main" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		if !reflect.DeepEqual(query["repositoryIds"], []string{"2813489", "9223372036854775807"}) {
			t.Fatalf("repositoryIds = %#v, raw query = %q", query["repositoryIds"], r.URL.RawQuery)
		}
		if !reflect.DeepEqual(query["authorUserIds"], []string{"author-1", "author-2"}) {
			t.Fatalf("authorUserIds = %#v, raw query = %q", query["authorUserIds"], r.URL.RawQuery)
		}
		if !reflect.DeepEqual(query["assigneeUserIds"], []string{"assignee-1"}) {
			t.Fatalf("assigneeUserIds = %#v, raw query = %q", query["assigneeUserIds"], r.URL.RawQuery)
		}
		if !reflect.DeepEqual(query["subscriberUserIds"], []string{"subscriber-1"}) {
			t.Fatalf("subscriberUserIds = %#v, raw query = %q", query["subscriberUserIds"], r.URL.RawQuery)
		}
		w.Header().Set("x-total", "2")
		_, _ = w.Write([]byte(`[{"iid":1}]`))
	})

	result, err := handleListMergeRequests(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"page":              float64(2),
		"perPage":           float64(20),
		"repositoryIds":     []any{"2813489", "9223372036854775807"},
		"authorUserIds":     []string{"author-1", "author-2"},
		"assigneeUserIds":   "assignee-1",
		"subscriberUserIds": "subscriber-1",
		"state":             "opened",
		"search":            "demo",
		"orderBy":           "updated_at",
		"createdAfter":      "2026-01-01",
		"createdBefore":     "2026-02-01",
		"targetBranch":      "main",
	})
	if err != nil {
		t.Fatalf("handleListMergeRequests() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetMergeRequestBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/repositories/group%2Frepo/mergeRequests/9223372036854775807" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"iid":9223372036854775807}`))
	})

	if _, err := handleGetMergeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"iid":            "9223372036854775807",
	}); err != nil {
		t.Fatalf("handleGetMergeRequest() error = %v", err)
	}
}

func TestHandleGetMergeRequestPreservesEncodedRepositoryID(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/repositories/group%2Frepo/mergeRequests/12" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"iid":12}`))
	})

	if _, err := handleGetMergeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group%2Frepo",
		"iid":            float64(12),
	}); err != nil {
		t.Fatalf("handleGetMergeRequest() error = %v", err)
	}
}

func TestCodeUpMergeRequestHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	_, err := handleListMergeRequests(context.Background(), client, map[string]any{})
	if err == nil {
		t.Fatal("handleListMergeRequests() expected missing organizationId error")
	}

	_, err = handleGetMergeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
	})
	if err == nil {
		t.Fatal("handleGetMergeRequest() expected missing iid error")
	}
}
