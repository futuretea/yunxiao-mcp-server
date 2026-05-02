package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListKnowledgeBasesBuildsQueryWithMetadata(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/lingma/organizations/org-1/knowledgeBases" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("query") != "docs" ||
			r.URL.Query().Get("sceneType") != "chat" ||
			r.URL.Query().Get("orderBy") != "gmt_created" ||
			r.URL.Query().Get("sort") != "desc" ||
			r.URL.Query().Get("userId") != "user-1" ||
			r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"kbId":"kb-1"}]`))
	})

	result, err := handleListKnowledgeBases(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"query":          "docs",
		"sceneType":      "chat",
		"orderBy":        "gmt_created",
		"sort":           "desc",
		"userId":         "user-1",
		"page":           float64(2),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleListKnowledgeBases() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleListKbFilesBuildsEscapedPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/lingma/organizations/org-1/knowledgeBases/kb%2F1/files?") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		if r.URL.Query().Get("query") != "readme" ||
			r.URL.Query().Get("orderBy") != "gmt_added" ||
			r.URL.Query().Get("sort") != "asc" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"fileId":"file-1"}]`))
	})

	result, err := handleListKbFiles(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"kbId":           "kb/1",
		"query":          "readme",
		"orderBy":        "gmt_added",
		"sort":           "asc",
	})
	if err != nil {
		t.Fatalf("handleListKbFiles() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleListKbMembersBuildsEscapedPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/lingma/organizations/org-1/knowledgeBases/kb%2F1/members?") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if r.URL.Query().Get("query") != "alice" ||
			r.URL.Query().Get("page") != "1" ||
			r.URL.Query().Get("perPage") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"userId":"user-1"}]`))
	})

	result, err := handleListKbMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"kbId":           "kb/1",
		"query":          "alice",
		"page":           float64(1),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleListKbMembers() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleListKbMembersPreservesEncodedKBID(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/oapi/v1/lingma/organizations/org-1/knowledgeBases/kb%2F1/members" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"userId":"user-1"}]`))
	})

	if _, err := handleListKbMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"kbId":           "kb%2F1",
	}); err != nil {
		t.Fatalf("handleListKbMembers() error = %v", err)
	}
}

func TestLingmaKnowledgeBaseHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleListKnowledgeBases(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListKnowledgeBases(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListKbFiles(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("handleListKbFiles() expected missing kbId error")
	}
	if _, err := handleListKbFiles(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "kbId": "kb-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListKbMembers(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing kbId error")
	}
	if _, err := handleListKbMembers(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "kbId": "kb-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}

func TestLingmaKnowledgeBaseHelpers(t *testing.T) {
	q := lingmaKnowledgeBaseListQuery(map[string]any{
		"query":     "docs",
		"sceneType": "chat",
		"userId":    "user-1",
		"page":      float64(2),
		"perPage":   float64(20),
	})
	if q.Get("query") != "docs" || q.Get("sceneType") != "chat" || q.Get("userId") != "user-1" || q.Get("page") != "2" || q.Get("perPage") != "20" {
		t.Fatalf("lingmaKnowledgeBaseListQuery() = %v", q)
	}

	cq := lingmaKnowledgeBaseChildListQuery(map[string]any{"query": "readme", "orderBy": "id", "sort": "asc", "page": float64(1)})
	if cq.Get("query") != "readme" || cq.Get("orderBy") != "id" || cq.Get("sort") != "asc" || cq.Get("page") != "1" {
		t.Fatalf("lingmaKnowledgeBaseChildListQuery() = %v", cq)
	}

	if got := lingmaKnowledgeBasePath("org-1", "kb/1"); got != "/lingma/organizations/org-1/knowledgeBases/kb%2F1" {
		t.Fatalf("lingmaKnowledgeBasePath() = %q", got)
	}
}
