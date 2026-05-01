package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetProjectOverviewBuildsCommonRequests(t *testing.T) {
	seen := map[string]bool{}
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seen[r.URL.Path] = true
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}

		switch r.URL.Path {
		case "/oapi/v1/projex/organizations/org-1/projects/project-1":
			_, _ = w.Write([]byte(`{"id":"project-1"}`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/members":
			_, _ = w.Write([]byte(`[{"id":"user-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/sprints",
			"/oapi/v1/projex/organizations/org-1/projects/project-1/milestones",
			"/oapi/v1/projex/organizations/org-1/projects/project-1/versions":
			if r.URL.Query().Get("status") != "TODO,DOING" ||
				r.URL.Query().Get("page") != "2" ||
				r.URL.Query().Get("perPage") != "3" {
				t.Fatalf("query = %q", r.URL.RawQuery)
			}
			w.Header().Set("x-total", "1")
			_, _ = w.Write([]byte(`[{"id":"item-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/projects/project-1/labels":
			if r.URL.Query().Get("status") != "" ||
				r.URL.Query().Get("page") != "2" ||
				r.URL.Query().Get("perPage") != "3" {
				t.Fatalf("query = %q", r.URL.RawQuery)
			}
			w.Header().Set("x-total", "1")
			_, _ = w.Write([]byte(`[{"id":"label-1"}]`))
		default:
			t.Fatalf("unexpected path = %q", r.URL.Path)
		}
	})

	result, err := handleGetProjectOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"page":           float64(2),
		"perPage":        float64(3),
	})
	if err != nil {
		t.Fatalf("handleGetProjectOverview() error = %v", err)
	}
	for _, path := range []string{
		"/oapi/v1/projex/organizations/org-1/projects/project-1",
		"/oapi/v1/projex/organizations/org-1/projects/project-1/members",
		"/oapi/v1/projex/organizations/org-1/projects/project-1/sprints",
		"/oapi/v1/projex/organizations/org-1/projects/project-1/milestones",
		"/oapi/v1/projex/organizations/org-1/projects/project-1/versions",
		"/oapi/v1/projex/organizations/org-1/projects/project-1/labels",
	} {
		if !seen[path] {
			t.Fatalf("missing request to %s", path)
		}
	}
	if !strings.Contains(result, `"sprints"`) || !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetProjectOverviewHonorsIncludeFlags(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1" {
			t.Fatalf("unexpected path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"project-1"}`))
	})

	result, err := handleGetProjectOverview(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"id":                "project-1",
		"includeMembers":    false,
		"includeSprints":    false,
		"includeMilestones": false,
		"includeVersions":   false,
		"includeLabels":     false,
	})
	if err != nil {
		t.Fatalf("handleGetProjectOverview() error = %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	for _, absent := range []string{"members", "sprints", "milestones", "versions", "labels"} {
		if _, ok := payload[absent]; ok {
			t.Fatalf("section %q should be absent in %#v", absent, payload)
		}
	}
}
