package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetApplicationOverviewRequiresOrganizationId(t *testing.T) {
	_, err := handleGetApplicationOverview(context.Background(), nil, map[string]any{
		"appName": "app-1",
	})
	if err == nil || !strings.Contains(err.Error(), "organizationId is required") {
		t.Fatalf("expected organizationId required error, got %v", err)
	}
}

func TestHandleGetApplicationOverviewRequiresAppName(t *testing.T) {
	_, err := handleGetApplicationOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil || !strings.Contains(err.Error(), "appName is required") {
		t.Fatalf("expected appName required error, got %v", err)
	}
}

func TestHandleGetApplicationOverviewRequiresClient(t *testing.T) {
	_, err := handleGetApplicationOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	})
	if err == nil || !strings.Contains(err.Error(), "yunxiao client is not configured") {
		t.Fatalf("expected client error, got %v", err)
	}
}

func TestHandleGetApplicationOverviewReturnsErrorOnApplicationFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleGetApplicationOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	})
	if err == nil || !strings.Contains(err.Error(), "application:") {
		t.Fatalf("expected application error, got %v", err)
	}
}

func TestHandleGetApplicationOverviewReturnsErrorOnEnvironmentsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/envs") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"envs boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"app-1"}`))
	})
	_, err := handleGetApplicationOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	})
	if err == nil || !strings.Contains(err.Error(), "environments:") {
		t.Fatalf("expected environments error, got %v", err)
	}
}

func TestHandleGetApplicationOverviewReturnsErrorOnOrchestrationsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/orchestrations") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"orch boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"app-1"}`))
	})
	_, err := handleGetApplicationOverview(context.Background(), client, map[string]any{
		"organizationId":      "org-1",
		"appName":             "app-1",
		"includeEnvironments": false,
	})
	if err == nil || !strings.Contains(err.Error(), "orchestrations:") {
		t.Fatalf("expected orchestrations error, got %v", err)
	}
}

func TestHandleGetApplicationOverviewSuccessAllSections(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oapi/v1/appstack/organizations/org-1/apps/app-1":
			_, _ = w.Write([]byte(`{"id":"app-1"}`))
		case strings.HasSuffix(r.URL.Path, "/envs"):
			if r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("envs perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`["dev","prod"]`))
		case strings.HasSuffix(r.URL.Path, "/orchestrations"):
			if r.URL.Query().Get("perPage") != "5" {
				t.Fatalf("orchs perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[{"sn":"1"}]`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})

	result, err := handleGetApplicationOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	})
	if err != nil {
		t.Fatalf("handleGetApplicationOverview() error = %v", err)
	}
	if !strings.Contains(result, `"application"`) {
		t.Fatalf("result missing application: %q", result)
	}
	if !strings.Contains(result, `"environments"`) {
		t.Fatalf("result missing environments: %q", result)
	}
	if !strings.Contains(result, `"orchestrations"`) {
		t.Fatalf("result missing orchestrations: %q", result)
	}
}

func TestHandleGetApplicationOverviewSkipsSectionsWhenDisabled(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/appstack/organizations/org-1/apps/app-1" {
			_, _ = w.Write([]byte(`{"id":"app-1"}`))
			return
		}
		t.Fatalf("unexpected request to %q", r.URL.Path)
	})

	result, err := handleGetApplicationOverview(context.Background(), client, map[string]any{
		"organizationId":        "org-1",
		"appName":               "app-1",
		"includeEnvironments":   false,
		"includeOrchestrations": false,
	})
	if err != nil {
		t.Fatalf("handleGetApplicationOverview() error = %v", err)
	}
	if strings.Contains(result, `"environments"`) {
		t.Fatalf("result should not contain environments: %q", result)
	}
	if strings.Contains(result, `"orchestrations"`) {
		t.Fatalf("result should not contain orchestrations: %q", result)
	}
}

func TestHandleGetApplicationOverviewUsesCustomLimits(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/envs"):
			if r.URL.Query().Get("perPage") != "3" {
				t.Fatalf("envs perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[]`))
		case strings.HasSuffix(r.URL.Path, "/orchestrations"):
			if r.URL.Query().Get("perPage") != "2" {
				t.Fatalf("orchs perPage = %q", r.URL.Query().Get("perPage"))
			}
			_, _ = w.Write([]byte(`[]`))
		default:
			_, _ = w.Write([]byte(`{"id":"app-1"}`))
		}
	})

	_, err := handleGetApplicationOverview(context.Background(), client, map[string]any{
		"organizationId":     "org-1",
		"appName":            "app-1",
		"envLimit":           float64(3),
		"orchestrationLimit": float64(2),
	})
	if err != nil {
		t.Fatalf("handleGetApplicationOverview() error = %v", err)
	}
}

func TestApplicationOverviewFilters(t *testing.T) {
	params := map[string]any{
		"includeEnvironments":   false,
		"includeOrchestrations": false,
		"envLimit":              float64(10),
		"orchestrationLimit":    float64(20),
	}
	filters := applicationOverviewFilters(params)
	if filters["includeEnvironments"].(bool) != false {
		t.Fatalf("includeEnvironments = %v", filters["includeEnvironments"])
	}
	if filters["envLimit"].(int) != 10 {
		t.Fatalf("envLimit = %v", filters["envLimit"])
	}
	if filters["orchestrationLimit"].(int) != 20 {
		t.Fatalf("orchestrationLimit = %v", filters["orchestrationLimit"])
	}
}
