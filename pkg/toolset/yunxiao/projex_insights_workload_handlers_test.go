package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetMemberWorkloadTrendBuildsCorrectRequests(t *testing.T) {
	requestCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/members" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			_, _ = w.Write([]byte(`[{"userId":"user-1","userName":"Alice"}]`))
		case http.MethodPost:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			conditions, _ := body["conditions"].(string)
			if !strings.Contains(conditions, `"fieldIdentifier":"assignedTo"`) {
				t.Fatalf("missing assignedTo condition: %q", conditions)
			}
			w.Header().Set("x-total", "1")
			_, _ = w.Write([]byte(`[{"id":"wi-1","status":{"name":"doing"},"finishTime":"2099-01-01"}]`))
		default:
			t.Fatalf("method = %s", r.Method)
		}
	})

	result, err := handleGetMemberWorkloadTrend(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
		"memberLimit":    float64(1),
	})
	if err != nil {
		t.Fatalf("handleGetMemberWorkloadTrend() error = %v", err)
	}
	if requestCount != 3 {
		t.Fatalf("requests = %d, want 3", requestCount)
	}
	if !strings.Contains(result, `"members"`) {
		t.Fatalf("result missing members: %q", result)
	}
	if !strings.Contains(result, `"summary"`) {
		t.Fatalf("result missing summary: %q", result)
	}
}

func TestHandleGetMemberWorkloadTrendRejectsEmptyCategories(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without categories")
	})
	if _, err := handleGetMemberWorkloadTrend(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     ", ,",
	}); err == nil {
		t.Fatal("expected missing categories error")
	}
}

func TestHandleGetMemberWorkloadTrendEmptyAssignees(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, _ = w.Write([]byte(`[]`))
		} else {
			_, _ = w.Write([]byte(`[]`))
		}
	})
	_, err := handleGetMemberWorkloadTrend(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	})
	if err == nil {
		t.Fatal("expected error for empty assignee IDs")
	}
}

func TestHandleGetMemberWorkloadTrendMissingProjectId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleGetMemberWorkloadTrend(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil {
		t.Fatal("expected error for missing projectId")
	}
}

func TestHandleGetMemberWorkloadTrendNilClient(t *testing.T) {
	_, err := handleGetMemberWorkloadTrend(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestHandleGetMemberWorkloadTrendReturnsMembersError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(`[]`))
	})
	if _, err := handleGetMemberWorkloadTrend(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	}); err == nil {
		t.Fatal("expected members fetch error")
	}
}

func TestHandleGetTeamWorkloadBreakdownBuildsCorrectRequests(t *testing.T) {
	requestCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/members" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			_, _ = w.Write([]byte(`[{"userId":"user-1","userName":"Alice"}]`))
		case http.MethodPost:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["perPage"].(float64) != 5 {
				t.Fatalf("taskLimit not applied: perPage = %v", body["perPage"])
			}
			w.Header().Set("x-total", "1")
			_, _ = w.Write([]byte(`[{"id":"wi-1","serialNumber":"ZGPQ-1","subject":"Test task","status":{"name":"doing"},"labels":[{"name":"area/test"}],"gmtCreate":1234567890}]`))
		default:
			t.Fatalf("method = %s", r.Method)
		}
	})

	result, err := handleGetTeamWorkloadBreakdown(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
		"memberLimit":    float64(1),
		"taskLimit":      float64(5),
	})
	if err != nil {
		t.Fatalf("handleGetTeamWorkloadBreakdown() error = %v", err)
	}
	if requestCount != 2 {
		t.Fatalf("requests = %d, want 2", requestCount)
	}
	if !strings.Contains(result, `"tasks"`) {
		t.Fatalf("result missing tasks: %q", result)
	}
	if !strings.Contains(result, `"ZGPQ-1"`) {
		t.Fatalf("result missing task serialNumber: %q", result)
	}
	if !strings.Contains(result, `"area/test"`) {
		t.Fatalf("result missing label: %q", result)
	}
}

func TestClampTaskLimit(t *testing.T) {
	tests := []struct {
		limit int
		want  int
	}{
		{-5, 1},
		{-1, 1},
		{0, 1},
		{1, 1},
		{25, 25},
		{50, 50},
		{51, 50},
		{100, 50},
	}
	for _, tt := range tests {
		got := clampTaskLimit(tt.limit)
		if got != tt.want {
			t.Errorf("clampTaskLimit(%d) = %d, want %d", tt.limit, got, tt.want)
		}
	}
}

func TestExtractWorkitemStatusName(t *testing.T) {
	tests := []struct {
		name    string
		itemMap map[string]any
		want    string
	}{
		{
			"valid status",
			map[string]any{"status": map[string]any{"name": "DOING"}},
			"DOING",
		},
		{
			"missing status key",
			map[string]any{"id": "wi-1"},
			"Unknown",
		},
		{
			"status is not a map",
			map[string]any{"status": "DOING"},
			"Unknown",
		},
		{
			"status map without name",
			map[string]any{"status": map[string]any{"stage": "dev"}},
			"Unknown",
		},
		{
			"empty map",
			map[string]any{},
			"Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractWorkitemStatusName(tt.itemMap)
			if got != tt.want {
				t.Errorf("extractWorkitemStatusName(%v) = %q, want %q", tt.itemMap, got, tt.want)
			}
		})
	}
}

func TestHandleGetTeamWorkloadBreakdownRejectsEmptyCategories(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without categories")
	})
	if _, err := handleGetTeamWorkloadBreakdown(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     ", ,",
	}); err == nil {
		t.Fatal("expected missing categories error")
	}
}

func TestHandleGetTeamWorkloadBreakdownReturnsMembersError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(`[]`))
	})
	if _, err := handleGetTeamWorkloadBreakdown(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	}); err == nil {
		t.Fatal("expected members fetch error")
	}
}
