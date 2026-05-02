package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListCurrentUserEffortRecordsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/effortRecords" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("startDate") != "2026-01-01" ||
			r.URL.Query().Get("endDate") != "2026-02-01" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[{"id":"effort-1"}]`))
	})

	if _, err := handleListCurrentUserEffortRecords(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"startDate":      "2026-01-01",
		"endDate":        "2026-02-01",
	}); err != nil {
		t.Fatalf("handleListCurrentUserEffortRecords() error = %v", err)
	}
}

func TestHandleListCurrentUserEffortRecordsRequiresDateRange(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	_, err := handleListCurrentUserEffortRecords(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"startDate":      "2026-01-01",
	})
	if err == nil || !strings.Contains(err.Error(), "endDate is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestHandleListEffortRecordsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/workitems/workitem-1/effortRecords" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"id":"effort-1"}]`))
	})

	if _, err := handleListEffortRecords(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "workitem-1",
	}); err != nil {
		t.Fatalf("handleListEffortRecords() error = %v", err)
	}
}

func TestHandleListEstimatedEffortsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/workitems/workitem-1/estimatedEfforts" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"id":"estimated-1"}]`))
	})

	if _, err := handleListEstimatedEfforts(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "workitem-1",
	}); err != nil {
		t.Fatalf("handleListEstimatedEfforts() error = %v", err)
	}
}

func TestProjexEffortHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleListCurrentUserEffortRecords(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListCurrentUserEffortRecords(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing startDate error")
	}
	if _, err := handleListCurrentUserEffortRecords(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "startDate": "2026-01-01", "endDate": "2026-02-01"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListEffortRecords(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListEffortRecords(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "id": "wi-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleListEstimatedEfforts(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListEstimatedEfforts(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "id": "wi-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}

func TestProjexWorkitemPath(t *testing.T) {
	if got := projexWorkitemPath("org-1", "wi/1"); got != "/projex/organizations/org-1/workitems/wi%2F1" {
		t.Fatalf("projexWorkitemPath() = %q", got)
	}
}
