package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetProjectRiskDashboardBuildsRiskQueries(t *testing.T) {
	seenCategories := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		category, _ := body["category"].(string)
		seenCategories[category] = true
		if body["spaceId"] != "project-1" || body["perPage"].(float64) != 2 {
			t.Fatalf("body = %#v", body)
		}
		conditions, _ := body["conditions"].(string)
		if !strings.Contains(conditions, `"fieldIdentifier":"status"`) {
			t.Fatalf("conditions = %q", conditions)
		}
		if category == "Risk,Bug" && !strings.Contains(conditions, `"fieldIdentifier":"finishTime"`) {
			t.Fatalf("overdue conditions = %q", conditions)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"workitem-1"}]`))
	})

	result, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Risk,Bug",
		"status":         "open,doing",
		"overdueBefore":  "2026-05-01",
		"sampleLimit":    float64(2),
	})
	if err != nil {
		t.Fatalf("handleGetProjectRiskDashboard() error = %v", err)
	}
	for _, category := range []string{"Risk", "Bug", "Risk,Bug"} {
		if !seenCategories[category] {
			t.Fatalf("missing category query %q in %#v", category, seenCategories)
		}
	}
	if !strings.Contains(result, `"overdue"`) || !strings.Contains(result, `"byCategory"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetProjectRiskDashboardRejectsEmptyCategories(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without categories")
	})

	if _, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     ", ,",
	}); err == nil {
		t.Fatal("handleGetProjectRiskDashboard() expected missing categories error")
	}
}

func TestHandleGetProjectRiskDashboardReturnsOverdueSearchError(t *testing.T) {
	requestCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(`[]`))
	})

	if _, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Risk",
		"overdueBefore":  "2026-05-01",
	}); err == nil {
		t.Fatal("handleGetProjectRiskDashboard() expected overdue search error")
	}
}

func TestHandleGetProjectRiskDashboardWithHighPriorityAndStale(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		conditions, _ := body["conditions"].(string)
		if strings.Contains(conditions, `"fieldIdentifier":"priority"`) {
			seen["highPriority"] = true
		}
		if strings.Contains(conditions, `"fieldIdentifier":"updateStatusAt"`) {
			seen["stale"] = true
		}
		if strings.Contains(conditions, `"fieldIdentifier":"finishTime"`) {
			seen["overdue"] = true
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[]`))
	})

	_, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
		"highPriority":   "p0",
		"staleBefore":    "2024-01-01",
	})
	if err != nil {
		t.Fatalf("handleGetProjectRiskDashboard() error = %v", err)
	}
	if !seen["overdue"] {
		t.Fatal("expected overdue search")
	}
	if !seen["highPriority"] {
		t.Fatal("expected highPriority search")
	}
	if !seen["stale"] {
		t.Fatal("expected stale search")
	}
}

func TestHandleGetProjectRiskDashboardReturnsSearchError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	}); err == nil {
		t.Fatal("expected search error")
	}
}

func TestHandleGetProjectRiskDashboardReturnsFocusSearchError(t *testing.T) {
	requests := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		conditions, _ := body["conditions"].(string)
		if strings.Contains(conditions, `"fieldIdentifier":"priority"`) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[]`))
	})
	if _, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
		"highPriority":   "p0",
	}); err == nil {
		t.Fatal("expected highPriority search error")
	}
}

func TestHandleGetProjectRiskDashboardReturnsStaleSearchError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		conditions, _ := body["conditions"].(string)
		if strings.Contains(conditions, `"fieldIdentifier":"updateStatusAt"`) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[]`))
	})
	if _, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
		"staleBefore":    "2024-01-01",
	}); err == nil {
		t.Fatal("expected stale search error")
	}
}

func TestProjexInsightsRiskHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetProjectRiskDashboard(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleGetProjectRiskDashboard(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "projectId": "p-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}
