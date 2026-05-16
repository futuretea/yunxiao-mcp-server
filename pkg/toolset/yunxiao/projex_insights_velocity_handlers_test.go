package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetSprintVelocityBuildsCorrectRequests(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/sprints" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			seen["sprints"] = true
			w.Header().Set("x-total", "2")
			_, _ = w.Write([]byte(`[{"id":"sp-1","name":"Sprint 1"},{"id":"sp-2","name":"Sprint 2"}]`))
		case http.MethodPost:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems:search" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			category, _ := body["category"].(string)
			seen["cat:"+category] = true
			conditions, _ := body["conditions"].(string)
			if !strings.Contains(conditions, `"fieldIdentifier":"sprint"`) {
				t.Fatalf("missing sprint condition: %q", conditions)
			}
			w.Header().Set("x-total", "3")
			_, _ = w.Write([]byte(`[{"id":"wi-1","status":{"stage":"DONE"}},{"id":"wi-2","status":{"stage":"DOING"}},{"id":"wi-3","status":{"stage":"DONE"}}]`))
		default:
			t.Fatalf("method = %s", r.Method)
		}
	})

	result, err := handleGetSprintVelocity(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
		"sprintCount":    float64(2),
	})
	if err != nil {
		t.Fatalf("handleGetSprintVelocity() error = %v", err)
	}
	if !seen["sprints"] {
		t.Fatal("missing sprints request")
	}
	if !seen["cat:Task"] {
		t.Fatal("missing Task search request")
	}
	if !strings.Contains(result, `"sprints"`) {
		t.Fatalf("result missing sprints: %q", result)
	}
	if !strings.Contains(result, `"rate"`) {
		t.Fatalf("result missing rate: %q", result)
	}
}

func TestHandleGetSprintVelocityRejectsEmptyCategories(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without categories")
	})
	if _, err := handleGetSprintVelocity(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     ", ,",
	}); err == nil {
		t.Fatal("expected missing categories error")
	}
}

func TestClampSprintCount(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{-5, 1},
		{-1, 1},
		{0, 1},
		{1, 1},
		{10, 10},
		{20, 20},
		{21, 20},
		{100, 20},
	}
	for _, tt := range tests {
		got := clampSprintCount(tt.n)
		if got != tt.want {
			t.Errorf("clampSprintCount(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestParseSprintList(t *testing.T) {
	tests := []struct {
		name string
		data any
		want int
	}{
		{"slice", []any{"a"}, 1},
		{"map with data", map[string]any{"data": []any{"x", "y"}}, 2},
		{"empty map", map[string]any{}, -1},
		{"map with non-list data", map[string]any{"data": 123}, -1},
		{"string", "nope", -1},
		{"nil", nil, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSprintList(tt.data)
			if tt.want < 0 {
				if got != nil {
					t.Fatalf("parseSprintList(%v) = %v, want nil", tt.data, got)
				}
			} else if len(got) != tt.want {
				t.Fatalf("parseSprintList(%v) len = %d, want %d", tt.data, len(got), tt.want)
			}
		})
	}
}

func TestHandleGetSprintVelocityReturnsSprintError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := handleGetSprintVelocity(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"categories":     "Task",
	}); err == nil {
		t.Fatal("expected sprint list error")
	}
}
