package yunxiao

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestHandleSearchReleasesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/releases:search" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		query := r.URL.Query()
		if query.Get("pagination") != "keyset" ||
			query.Get("perPage") != "20" ||
			query.Get("orderBy") != "id" ||
			query.Get("sort") != "asc" ||
			query.Get("nextToken") != "token-1" ||
			query.Get("nameKeyword") != "demo-release" ||
			query.Get("systemName") != "system-1" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		if !reflect.DeepEqual(query["states"], []string{"DEVELOPING", "RELEASING"}) {
			t.Fatalf("states = %#v, raw query = %q", query["states"], r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"data":[{"name":"release"}],"nextToken":"token-2"}`))
	})

	if _, err := handleSearchReleases(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pagination":     "keyset",
		"perPage":        float64(20),
		"orderBy":        "id",
		"sort":           "asc",
		"nextToken":      "token-1",
		"nameKeyword":    "demo-release",
		"systemName":     "system-1",
		"states":         []any{"DEVELOPING", "RELEASING"},
	}); err != nil {
		t.Fatalf("handleSearchReleases() error = %v", err)
	}
}

func TestHandleSearchReleasesRequiresOrganizationID(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleSearchReleases(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("handleSearchReleases() expected missing organizationId error")
	}
}
