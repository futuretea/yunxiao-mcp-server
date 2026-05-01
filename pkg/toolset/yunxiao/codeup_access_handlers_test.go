package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListSSHKeysBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/codeup/organizations/org-1/keys" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "20" ||
			r.URL.Query().Get("orderBy") != "created_at" ||
			r.URL.Query().Get("sort") != "desc" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":1}]`))
	})

	result, err := handleListSSHKeys(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"page":           float64(2),
		"perPage":        float64(20),
		"orderBy":        "created_at",
		"sort":           "desc",
	})
	if err != nil {
		t.Fatalf("handleListSSHKeys() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetSSHKeyBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/keys/9223372036854775807" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"id":9223372036854775807}`))
	})

	if _, err := handleGetSSHKey(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"keyId":          "9223372036854775807",
	}); err != nil {
		t.Fatalf("handleGetSSHKey() error = %v", err)
	}
}

func TestHandleListUserSSHKeysBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/codeup/organizations/org-1/users/user-1/keys" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "3" ||
			r.URL.Query().Get("perPage") != "30" ||
			r.URL.Query().Get("orderBy") != "updated_at" ||
			r.URL.Query().Get("sort") != "asc" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "2")
		_, _ = w.Write([]byte(`[{"id":1},{"id":2}]`))
	})

	result, err := handleListUserSSHKeys(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"userId":         "user-1",
		"page":           float64(3),
		"perPage":        float64(30),
		"orderBy":        "updated_at",
		"sort":           "asc",
	})
	if err != nil {
		t.Fatalf("handleListUserSSHKeys() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleListWebHooksBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/repositories/group%2Frepo/webhooks?page=2&perPage=20" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":7}]`))
	})

	result, err := handleListWebHooks(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"page":           float64(2),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleListWebHooks() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetWebHookBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/repositories/group%2Frepo/webhooks/9223372036854775807" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"id":9223372036854775807}`))
	})

	if _, err := handleGetWebHook(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"hookId":         "9223372036854775807",
	}); err != nil {
		t.Fatalf("handleGetWebHook() error = %v", err)
	}
}
