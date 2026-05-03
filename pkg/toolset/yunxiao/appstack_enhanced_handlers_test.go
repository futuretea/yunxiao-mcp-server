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

func TestHandleGetEnvironmentOverviewRequiresOrganizationId(t *testing.T) {
	_, err := handleGetEnvironmentOverview(context.Background(), nil, map[string]any{
		"appName": "app-1",
		"envName": "dev",
	})
	if err == nil || !strings.Contains(err.Error(), "organizationId is required") {
		t.Fatalf("expected organizationId required error, got %v", err)
	}
}

func TestHandleGetEnvironmentOverviewRequiresAppName(t *testing.T) {
	_, err := handleGetEnvironmentOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"envName":        "dev",
	})
	if err == nil || !strings.Contains(err.Error(), "appName is required") {
		t.Fatalf("expected appName required error, got %v", err)
	}
}

func TestHandleGetEnvironmentOverviewRequiresEnvName(t *testing.T) {
	_, err := handleGetEnvironmentOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
	})
	if err == nil || !strings.Contains(err.Error(), "envName is required") {
		t.Fatalf("expected envName required error, got %v", err)
	}
}

func TestHandleGetEnvironmentOverviewRequiresClient(t *testing.T) {
	_, err := handleGetEnvironmentOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"envName":        "dev",
	})
	if err == nil || !strings.Contains(err.Error(), "yunxiao client is not configured") {
		t.Fatalf("expected client error, got %v", err)
	}
}

func TestHandleGetEnvironmentOverviewReturnsErrorOnEnvFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleGetEnvironmentOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"envName":        "dev",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetEnvironmentOverviewReturnsErrorOnVariableGroupsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/variableGroups") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"vg boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"env-1"}`))
	})
	_, err := handleGetEnvironmentOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"envName":        "dev",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetEnvironmentOverviewReturnsErrorOnOrchestrationFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ":latestAvailable") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"orch boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"env-1"}`))
	})
	_, err := handleGetEnvironmentOverview(context.Background(), client, map[string]any{
		"organizationId":        "org-1",
		"appName":               "app-1",
		"envName":               "dev",
		"includeVariableGroups": false,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetEnvironmentOverviewSuccessAllSections(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oapi/v1/appstack/organizations/org-1/apps/app-1/envs/dev":
			_, _ = w.Write([]byte(`{"id":"env-1"}`))
		case strings.HasSuffix(r.URL.Path, "/variableGroups"):
			_, _ = w.Write([]byte(`["vg-1"]`))
		case strings.HasSuffix(r.URL.Path, ":latestAvailable"):
			_, _ = w.Write([]byte(`{"sn":"1"}`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})

	result, err := handleGetEnvironmentOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"envName":        "dev",
	})
	if err != nil {
		t.Fatalf("handleGetEnvironmentOverview() error = %v", err)
	}
	if !strings.Contains(result, `"environment"`) {
		t.Fatalf("result missing environment: %q", result)
	}
	if !strings.Contains(result, `"variableGroups"`) {
		t.Fatalf("result missing variableGroups: %q", result)
	}
	if !strings.Contains(result, `"latestOrchestration"`) {
		t.Fatalf("result missing latestOrchestration: %q", result)
	}
}

func TestHandleGetEnvironmentOverviewSkipsSectionsWhenDisabled(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/appstack/organizations/org-1/apps/app-1/envs/dev" {
			_, _ = w.Write([]byte(`{"id":"env-1"}`))
			return
		}
		t.Fatalf("unexpected request to %q", r.URL.Path)
	})

	result, err := handleGetEnvironmentOverview(context.Background(), client, map[string]any{
		"organizationId":             "org-1",
		"appName":                    "app-1",
		"envName":                    "dev",
		"includeVariableGroups":      false,
		"includeLatestOrchestration": false,
	})
	if err != nil {
		t.Fatalf("handleGetEnvironmentOverview() error = %v", err)
	}
	if strings.Contains(result, `"variableGroups"`) {
		t.Fatalf("result should not contain variableGroups: %q", result)
	}
	if strings.Contains(result, `"latestOrchestration"`) {
		t.Fatalf("result should not contain latestOrchestration: %q", result)
	}
}

func TestEnvironmentOverviewFilters(t *testing.T) {
	params := map[string]any{
		"includeVariableGroups":      false,
		"includeLatestOrchestration": false,
	}
	filters := environmentOverviewFilters(params)
	if filters["includeVariableGroups"].(bool) != false {
		t.Fatalf("includeVariableGroups = %v", filters["includeVariableGroups"])
	}
	if filters["includeLatestOrchestration"].(bool) != false {
		t.Fatalf("includeLatestOrchestration = %v", filters["includeLatestOrchestration"])
	}
}

func TestHandleGetReleaseOverviewRequiresOrganizationId(t *testing.T) {
	_, err := handleGetReleaseOverview(context.Background(), nil, map[string]any{
		"systemName": "sys-1",
		"sn":         "rel-1",
	})
	if err == nil || !strings.Contains(err.Error(), "organizationId is required") {
		t.Fatalf("expected organizationId required error, got %v", err)
	}
}

func TestHandleGetReleaseOverviewRequiresSystemName(t *testing.T) {
	_, err := handleGetReleaseOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"sn":             "rel-1",
	})
	if err == nil || !strings.Contains(err.Error(), "systemName is required") {
		t.Fatalf("expected systemName required error, got %v", err)
	}
}

func TestHandleGetReleaseOverviewRequiresSn(t *testing.T) {
	_, err := handleGetReleaseOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"systemName":     "sys-1",
	})
	if err == nil || !strings.Contains(err.Error(), "sn is required") {
		t.Fatalf("expected sn required error, got %v", err)
	}
}

func TestHandleGetReleaseOverviewRequiresClient(t *testing.T) {
	_, err := handleGetReleaseOverview(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"systemName":     "sys-1",
		"sn":             "rel-1",
	})
	if err == nil || !strings.Contains(err.Error(), "yunxiao client is not configured") {
		t.Fatalf("expected client error, got %v", err)
	}
}

func TestHandleGetReleaseOverviewReturnsErrorOnReleaseFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleGetReleaseOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"systemName":     "sys-1",
		"sn":             "rel-1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetReleaseOverviewReturnsErrorOnMembersFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/members") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"members boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"rel-1"}`))
	})
	_, err := handleGetReleaseOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"systemName":     "sys-1",
		"sn":             "rel-1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetReleaseOverviewReturnsErrorOnProductsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/products") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"products boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"rel-1"}`))
	})
	_, err := handleGetReleaseOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"systemName":     "sys-1",
		"sn":             "rel-1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetReleaseOverviewReturnsErrorOnChangeRequestsFailure(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/changeRequests") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"cr boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"rel-1"}`))
	})
	_, err := handleGetReleaseOverview(context.Background(), client, map[string]any{
		"organizationId":  "org-1",
		"systemName":      "sys-1",
		"sn":              "rel-1",
		"includeMembers":  false,
		"includeProducts": false,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHandleGetReleaseOverviewSuccessAllSections(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oapi/v1/appstack/organizations/org-1/systems/sys-1/releases/rel-1":
			_, _ = w.Write([]byte(`{"id":"rel-1"}`))
		case strings.HasSuffix(r.URL.Path, "/members"):
			_, _ = w.Write([]byte(`["member-a"]`))
		case strings.HasSuffix(r.URL.Path, "/products"):
			_, _ = w.Write([]byte(`["product-a"]`))
		case strings.HasSuffix(r.URL.Path, "/changeRequests"):
			if r.URL.Query().Get("pageSize") != "5" {
				t.Fatalf("cr pageSize = %q", r.URL.Query().Get("pageSize"))
			}
			_, _ = w.Write([]byte(`["cr-a"]`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})

	result, err := handleGetReleaseOverview(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"systemName":     "sys-1",
		"sn":             "rel-1",
	})
	if err != nil {
		t.Fatalf("handleGetReleaseOverview() error = %v", err)
	}
	if !strings.Contains(result, `"release"`) {
		t.Fatalf("result missing release: %q", result)
	}
	if !strings.Contains(result, `"members"`) {
		t.Fatalf("result missing members: %q", result)
	}
	if !strings.Contains(result, `"products"`) {
		t.Fatalf("result missing products: %q", result)
	}
	if !strings.Contains(result, `"changeRequests"`) {
		t.Fatalf("result missing changeRequests: %q", result)
	}
}

func TestHandleGetReleaseOverviewSkipsSectionsWhenDisabled(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/appstack/organizations/org-1/systems/sys-1/releases/rel-1" {
			_, _ = w.Write([]byte(`{"id":"rel-1"}`))
			return
		}
		t.Fatalf("unexpected request to %q", r.URL.Path)
	})

	result, err := handleGetReleaseOverview(context.Background(), client, map[string]any{
		"organizationId":        "org-1",
		"systemName":            "sys-1",
		"sn":                    "rel-1",
		"includeMembers":        false,
		"includeProducts":       false,
		"includeChangeRequests": false,
	})
	if err != nil {
		t.Fatalf("handleGetReleaseOverview() error = %v", err)
	}
	if strings.Contains(result, `"members"`) {
		t.Fatalf("result should not contain members: %q", result)
	}
	if strings.Contains(result, `"products"`) {
		t.Fatalf("result should not contain products: %q", result)
	}
	if strings.Contains(result, `"changeRequests"`) {
		t.Fatalf("result should not contain changeRequests: %q", result)
	}
}

func TestHandleGetReleaseOverviewUsesCustomLimit(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/changeRequests") {
			if r.URL.Query().Get("pageSize") != "3" {
				t.Fatalf("cr pageSize = %q", r.URL.Query().Get("pageSize"))
			}
			_, _ = w.Write([]byte(`[]`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"rel-1"}`))
	})

	_, err := handleGetReleaseOverview(context.Background(), client, map[string]any{
		"organizationId":     "org-1",
		"systemName":         "sys-1",
		"sn":                 "rel-1",
		"includeMembers":     false,
		"includeProducts":    false,
		"changeRequestLimit": float64(3),
	})
	if err != nil {
		t.Fatalf("handleGetReleaseOverview() error = %v", err)
	}
}

func TestReleaseOverviewFilters(t *testing.T) {
	params := map[string]any{
		"includeMembers":        false,
		"includeProducts":       false,
		"includeChangeRequests": false,
		"changeRequestLimit":    float64(10),
	}
	filters := releaseOverviewFilters(params)
	if filters["includeMembers"].(bool) != false {
		t.Fatalf("includeMembers = %v", filters["includeMembers"])
	}
	if filters["includeProducts"].(bool) != false {
		t.Fatalf("includeProducts = %v", filters["includeProducts"])
	}
	if filters["includeChangeRequests"].(bool) != false {
		t.Fatalf("includeChangeRequests = %v", filters["includeChangeRequests"])
	}
	if filters["changeRequestLimit"].(int) != 10 {
		t.Fatalf("changeRequestLimit = %v", filters["changeRequestLimit"])
	}
}
