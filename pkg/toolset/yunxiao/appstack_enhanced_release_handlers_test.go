package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

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
