package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestHandleSearchAppTagsBuildsPathQueryAndBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/appTags:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("current") != "2" || r.URL.Query().Get("pageSize") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["search"] != "demo" {
			t.Fatalf("body = %#v", body)
		}
		_, _ = w.Write([]byte(`{"total":1,"data":[]}`))
	})

	if _, err := handleSearchAppTags(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"current":        float64(2),
		"pageSize":       float64(20),
		"search":         "demo",
	}); err != nil {
		t.Fatalf("handleSearchAppTags() error = %v", err)
	}
}

func TestHandleSearchAppTagsUsesDefaultPagination(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Query().Get("current") != "1" || r.URL.Query().Get("pageSize") != "10" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"total":0,"data":[]}`))
	})

	if _, err := handleSearchAppTags(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	}); err != nil {
		t.Fatalf("handleSearchAppTags() error = %v", err)
	}
}

func TestHandleSearchAppTagsRequiresOrganizationID(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleSearchAppTags(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleSearchAppTags(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}

func TestHandleSearchAppTagsReturnsAPIError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := handleSearchAppTags(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	}); err == nil {
		t.Fatal("expected API error")
	}
}
