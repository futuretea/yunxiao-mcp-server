package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

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

	_, err := handleGetMergeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
	})
	if err == nil {
		t.Fatal("handleGetMergeRequest() expected missing iid error")
	}
}
