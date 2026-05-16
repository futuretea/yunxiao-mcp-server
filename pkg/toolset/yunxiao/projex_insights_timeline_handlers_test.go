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
