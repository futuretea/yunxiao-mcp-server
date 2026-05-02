package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListCommitStatusesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/oapi/v1/codeup/organizations/org-1/repositories/group%2Frepo/commits/abc123/statuses") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		if r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"state":"success"}]`))
	})

	result, err := handleListCommitStatuses(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"sha":            "abc123",
		"page":           float64(2),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleListCommitStatuses() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleListCheckRunsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/oapi/v1/codeup/organizations/org-1/repositories/group%2Frepo/checkRuns") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		if r.URL.Query().Get("ref") != "refs/heads/main" ||
			r.URL.Query().Get("page") != "3" ||
			r.URL.Query().Get("perPage") != "30" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "2")
		_, _ = w.Write([]byte(`[{"id":1},{"id":2}]`))
	})

	result, err := handleListCheckRuns(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"ref":            "refs/heads/main",
		"page":           float64(3),
		"perPage":        float64(30),
	})
	if err != nil {
		t.Fatalf("handleListCheckRuns() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleListCheckRunsRequiresRef(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	_, err := handleListCheckRuns(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
	})
	if err == nil || !strings.Contains(err.Error(), "ref is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestHandleGetCheckRunBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/repositories/group%2Frepo/checkRuns/9223372036854775807" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"id":9223372036854775807}`))
	})

	if _, err := handleGetCheckRun(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"checkRunId":     "9223372036854775807",
	}); err != nil {
		t.Fatalf("handleGetCheckRun() error = %v", err)
	}
}

func TestCodeupCheckHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleListCommitStatuses(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListCommitStatuses(context.Background(), client, map[string]any{"organizationId": "org-1", "repositoryId": "repo-1"}); err == nil {
		t.Fatal("expected missing sha error")
	}
	if _, err := handleListCommitStatuses(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "repositoryId": "repo-1", "sha": "abc123"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListCheckRuns(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListCheckRuns(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "repositoryId": "repo-1", "ref": "main"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetCheckRun(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetCheckRun(context.Background(), client, map[string]any{"organizationId": "org-1", "repositoryId": "repo-1"}); err == nil {
		t.Fatal("expected missing checkRunId error")
	}
	if _, err := handleGetCheckRun(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "repositoryId": "repo-1", "checkRunId": "1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}
