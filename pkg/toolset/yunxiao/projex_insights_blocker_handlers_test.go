package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetBlockerAnalysisBuildsCorrectRequests(t *testing.T) {
	requestCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		switch r.Method {
		case http.MethodPost:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			w.Header().Set("x-total", "2")
			_, _ = w.Write([]byte(`[{"id":"wi-1","status":{"name":"doing"}},{"id":"wi-2","status":{"name":"done"}}]`))
		case http.MethodGet:
			if strings.HasSuffix(r.URL.Path, "/relationRecords") {
				if strings.Contains(r.URL.Path, "wi-1") {
					_, _ = w.Write([]byte(`[{"relationType":"DEPEND_ON","target":{"status":{"stage":"DOING"}}}]`))
				} else {
					_, _ = w.Write([]byte(`[]`))
				}
			} else {
				t.Fatalf("unexpected path = %q", r.URL.Path)
			}
		default:
			t.Fatalf("method = %s", r.Method)
		}
	})

	result, err := handleGetBlockerAnalysis(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
		"sampleLimit":    float64(2),
	})
	if err != nil {
		t.Fatalf("handleGetBlockerAnalysis() error = %v", err)
	}
	if requestCount != 3 {
		t.Fatalf("requests = %d, want 3", requestCount)
	}
	if !strings.Contains(result, `"blocked"`) {
		t.Fatalf("result missing blocked: %q", result)
	}
	if !strings.Contains(result, `"blocking"`) {
		t.Fatalf("result missing blocking: %q", result)
	}
	if !strings.Contains(result, `"summary"`) {
		t.Fatalf("result missing summary: %q", result)
	}
}

func TestHandleGetBlockerAnalysisNilClient(t *testing.T) {
	_, err := handleGetBlockerAnalysis(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestHandleGetBlockerAnalysisDetectsBlockingItems(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			_, _ = w.Write([]byte(`[{"id":"wi-1","status":{"name":"doing"}}]`))
		case http.MethodGet:
			_, _ = w.Write([]byte(`[{"relationType":"DEPENDED_BY","target":{"status":{"stage":"DOING"}}}]`))
		}
	})
	result, err := handleGetBlockerAnalysis(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	})
	if err != nil {
		t.Fatalf("handleGetBlockerAnalysis() error = %v", err)
	}
	if !strings.Contains(result, `"blocking"`) {
		t.Fatal("result missing blocking section")
	}
}

func TestHandleGetBlockerAnalysisSkipsDoneStage(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			_, _ = w.Write([]byte(`[{"id":"wi-1","status":{"name":"doing"}}]`))
		case http.MethodGet:
			_, _ = w.Write([]byte(`[{"relationType":"DEPEND_ON","target":{"status":{"stage":"DONE"}}}]`))
		}
	})
	result, err := handleGetBlockerAnalysis(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	})
	if err != nil {
		t.Fatalf("handleGetBlockerAnalysis() error = %v", err)
	}
	if !strings.Contains(result, `"totalBlocked"`) || !strings.Contains(result, `"totalBlocking"`) {
		t.Fatalf("expected summary in result, got %q", result)
	}
}

func TestHandleGetBlockerAnalysisNetworkError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	result, err := handleGetBlockerAnalysis(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	})
	if err != nil {
		t.Fatalf("handleGetBlockerAnalysis() error = %v", err)
	}
	if !strings.Contains(result, `"totalBlocked"`) || !strings.Contains(result, `"totalBlocking"`) {
		t.Fatalf("expected summary with zero counts, got %q", result)
	}
}

func TestHandleGetBlockerAnalysisSkipsItemsWithoutID(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			_, _ = w.Write([]byte(`[{"status":{"name":"doing"}},{"id":"","status":{"name":"doing"}}]`))
		} else {
			t.Fatal("unexpected request for item without id")
		}
	})
	result, err := handleGetBlockerAnalysis(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	})
	if err != nil {
		t.Fatalf("handleGetBlockerAnalysis() error = %v", err)
	}
	if !strings.Contains(result, `"totalBlocked"`) {
		t.Fatalf("expected summary, got %q", result)
	}
}

func TestHandleGetBlockerAnalysisSkipsNonMapItems(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			_, _ = w.Write([]byte(`["not-a-map",{"id":"wi-1","status":{"name":"doing"}}]`))
		} else {
			_, _ = w.Write([]byte(`[{"relationType":"DEPEND_ON","target":{"status":{"stage":"DOING"}}}]`))
		}
	})
	result, err := handleGetBlockerAnalysis(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	})
	if err != nil {
		t.Fatalf("handleGetBlockerAnalysis() error = %v", err)
	}
	if !strings.Contains(result, `"summary"`) {
		t.Fatalf("expected summary, got %q", result)
	}
}

func TestHandleGetBlockerAnalysisRejectsEmptyCategories(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without categories")
	})
	if _, err := handleGetBlockerAnalysis(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     ", ,",
	}); err == nil {
		t.Fatal("expected missing categories error")
	}
}
