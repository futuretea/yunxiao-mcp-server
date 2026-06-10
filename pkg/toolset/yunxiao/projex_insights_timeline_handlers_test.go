package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetWorkitemStatusTimelineBuildsCorrectRequests(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			switch r.URL.Path {
			case "/oapi/v1/projex/organizations/org-1/workitems/wi-1":
				seen["workitem"] = true
				_, _ = w.Write([]byte(`{"id":"wi-1","subject":"Test"}`))
			case "/oapi/v1/projex/organizations/org-1/workitems/wi-1/activities":
				seen["activities"] = true
				_, _ = w.Write([]byte(`[{"action":"UPDATE","field":"status","gmtCreate":1234567890,"oldValue":{"name":"backlog"},"newValue":{"name":"doing"}}]`))
			default:
				t.Fatalf("unexpected path = %q", r.URL.Path)
			}
		default:
			t.Fatalf("method = %s", r.Method)
		}
	})

	result, err := handleGetWorkitemStatusTimeline(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
	})
	if err != nil {
		t.Fatalf("handleGetWorkitemStatusTimeline() error = %v", err)
	}
	if !seen["workitem"] {
		t.Fatal("missing workitem request")
	}
	if !seen["activities"] {
		t.Fatal("missing activities request")
	}
	if !strings.Contains(result, `"timeline"`) {
		t.Fatalf("result missing timeline: %q", result)
	}
	if !strings.Contains(result, `"summary"`) {
		t.Fatalf("result missing summary: %q", result)
	}
}

func TestHandleGetWorkitemStatusTimelineSkipsWorkitemWhenDisabled(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/projex/organizations/org-1/workitems/wi-1" {
			seen["workitem"] = true
		}
		if r.URL.Path == "/oapi/v1/projex/organizations/org-1/workitems/wi-1/activities" {
			seen["activities"] = true
			_, _ = w.Write([]byte(`[]`))
		}
	})

	_, err := handleGetWorkitemStatusTimeline(context.Background(), client, map[string]any{
		"organizationId":  "org-1",
		"workitemId":      "wi-1",
		"includeWorkitem": false,
	})
	if err != nil {
		t.Fatalf("handleGetWorkitemStatusTimeline() error = %v", err)
	}
	if seen["workitem"] {
		t.Fatal("should not fetch workitem when includeWorkitem=false")
	}
	if !seen["activities"] {
		t.Fatal("missing activities request")
	}
}

func TestParseStatusTimeline(t *testing.T) {
	tests := []struct {
		name       string
		activities any
		wantLen    int
	}{
		{"nil", nil, 0},
		{"string input", "not-activities", 0},
		{
			"array with status update",
			[]any{
				map[string]any{"action": "UPDATE", "field": "status", "gmtCreate": float64(100), "operator": "user1", "oldValue": "backlog", "newValue": "doing"},
			},
			1,
		},
		{
			"array with non-status update",
			[]any{
				map[string]any{"action": "UPDATE", "field": "subject", "gmtCreate": float64(200)},
				map[string]any{"action": "CREATE", "field": "status", "gmtCreate": float64(300)},
			},
			0,
		},
		{
			"map with data key",
			map[string]any{"data": []any{
				map[string]any{"action": "UPDATE", "field": "status", "gmtCreate": float64(400), "operator": "u2", "oldValue": "doing", "newValue": "done"},
			}},
			1,
		},
		{"map without data key", map[string]any{"total": 5}, 0},
		{"map with non-list data", map[string]any{"data": "not-a-list"}, 0},
		{
			"skips non-map items",
			[]any{"not-a-map", map[string]any{"action": "UPDATE", "field": "status", "gmtCreate": float64(500)}},
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseStatusTimeline(tt.activities)
			if len(got) != tt.wantLen {
				t.Fatalf("parseStatusTimeline() len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestHandleGetWorkitemStatusTimelineNilClient(t *testing.T) {
	_, err := handleGetWorkitemStatusTimeline(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestHandleGetWorkitemStatusTimelineActivitiesError(t *testing.T) {
	callCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			_, _ = w.Write([]byte(`{"id":"wi-1","subject":"Test"}`))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	_, err := handleGetWorkitemStatusTimeline(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
	})
	if err == nil {
		t.Fatal("expected error for activities API failure")
	}
}

func TestHandleGetWorkitemStatusTimelineWorkitemFetchError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	_, err := handleGetWorkitemStatusTimeline(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
	})
	if err == nil {
		t.Fatal("expected error for workitem fetch failure")
	}
}

func TestCalculateTimelineStatsEmpty(t *testing.T) {
	got := calculateTimelineStats(nil)
	if len(got) != 1 {
		t.Fatalf("calculateTimelineStats(nil) len = %d, want 1", len(got))
	}
	if got["totalChanges"] != 0 {
		t.Fatalf("calculateTimelineStats(nil) totalChanges = %v, want 0", got["totalChanges"])
	}
	if _, ok := got["statusVisits"]; ok {
		t.Fatal("calculateTimelineStats(nil) should not include statusVisits")
	}
}

func TestHandleGetWorkitemStatusTimelineRequiresWorkitemId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request without workitemId")
	})
	if _, err := handleGetWorkitemStatusTimeline(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	}); err == nil {
		t.Fatal("expected missing workitemId error")
	}
}
