package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleCallYunxiaoAPIGetRequest(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/proj-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "1" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"id":"proj-1"}`))
	})

	result, err := handleCallYunxiaoAPI(context.Background(), client, map[string]any{
		"path":        "/projex/organizations/org-1/projects/proj-1",
		"queryParams": `{"page":"1"}`,
	})
	if err != nil {
		t.Fatalf("handleCallYunxiaoAPI() error = %v", err)
	}
	if !strings.Contains(result, `"id"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleCallYunxiaoAPIPostRequest(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/repo/list" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[]`))
	})

	_, err := handleCallYunxiaoAPI(context.Background(), client, map[string]any{
		"path":   "/projex/organizations/org-1/projects/repo/list",
		"method": "POST",
		"body":   `{}`,
	})
	if err != nil {
		t.Fatalf("handleCallYunxiaoAPI() error = %v", err)
	}
}

func TestHandleCallYunxiaoAPIRejectsInvalidMethod(t *testing.T) {
	_, err := handleCallYunxiaoAPI(context.Background(), nil, map[string]any{
		"path":   "/test",
		"method": "DELETE",
	})
	if err == nil || !strings.Contains(err.Error(), "method must be GET or POST") {
		t.Fatalf("expected method error, got %v", err)
	}
}

func TestHandleCallYunxiaoAPIRejectsBlockedPath(t *testing.T) {
	_, err := handleCallYunxiaoAPI(context.Background(), nil, map[string]any{
		"path": "/projex/organizations/org-1/workitems/wi-1/deletefile/1",
	})
	if err == nil || !strings.Contains(err.Error(), "blocked") {
		t.Fatalf("expected blocked error, got %v", err)
	}
}

func TestHandleCallYunxiaoAPIRejectsInvalidQueryParams(t *testing.T) {
	_, err := handleCallYunxiaoAPI(context.Background(), nil, map[string]any{
		"path":        "/test",
		"queryParams": "not-json",
	})
	if err == nil || !strings.Contains(err.Error(), "invalid queryParams JSON") {
		t.Fatalf("expected queryParams error, got %v", err)
	}
}

func TestHandleCallYunxiaoAPIRejectsInvalidBody(t *testing.T) {
	_, err := handleCallYunxiaoAPI(context.Background(), nil, map[string]any{
		"path": "/test",
		"body": "not-json",
	})
	if err == nil || !strings.Contains(err.Error(), "invalid body JSON") {
		t.Fatalf("expected body error, got %v", err)
	}
}

func TestHandleCallYunxiaoAPIDefaultsToGET(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		_, _ = w.Write([]byte(`{}`))
	})

	_, err := handleCallYunxiaoAPI(context.Background(), client, map[string]any{
		"path": "/platform/users:me",
	})
	if err != nil {
		t.Fatalf("handleCallYunxiaoAPI() error = %v", err)
	}
}
