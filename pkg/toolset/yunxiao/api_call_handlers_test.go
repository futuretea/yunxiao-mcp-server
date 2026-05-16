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

func TestHandleCallYunxiaoAPIAllowsReadOnlyPostPaths(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "search suffix", path: "/projex/organizations/org-1/projects:search"},
		{name: "list suffix", path: "/projex/organizations/org-1/projects/repo/list"},
		{name: "result list segment", path: "/projex/organizations/org-1/plan/result/list/dir-1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatalf("method = %s", r.Method)
				}
				if r.URL.Path != "/oapi/v1"+tt.path {
					t.Fatalf("path = %q, want %q", r.URL.Path, "/oapi/v1"+tt.path)
				}
				_, _ = w.Write([]byte(`[]`))
			})

			_, err := handleCallYunxiaoAPI(context.Background(), client, map[string]any{
				"path":   tt.path,
				"method": "POST",
				"body":   `{}`,
			})
			if err != nil {
				t.Fatalf("handleCallYunxiaoAPI() error = %v", err)
			}
		})
	}
}

func TestHandleCallYunxiaoAPIRejectsInvalidMethod(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request with invalid method")
	})
	_, err := handleCallYunxiaoAPI(context.Background(), client, map[string]any{
		"path":   "/test",
		"method": "DELETE",
	})
	if err == nil || !strings.Contains(err.Error(), "method must be GET or POST") {
		t.Fatalf("expected method error, got %v", err)
	}
}

func TestHandleCallYunxiaoAPIRejectsBlockedPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request with blocked path")
	})
	_, err := handleCallYunxiaoAPI(context.Background(), client, map[string]any{
		"path": "/projex/organizations/org-1/workitems/wi-1/deletefile/1",
	})
	if err == nil || !strings.Contains(err.Error(), "blocked") {
		t.Fatalf("expected blocked error, got %v", err)
	}
}

func TestHandleCallYunxiaoAPIRejectsEncodedBlockedPath(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "encoded mutation term", path: "/projex/organizations/org-1/workitems/wi-1/%64eletefile/1"},
		{name: "encoded slash before mutation term", path: "/projex/organizations/org-1/workitems/wi-1%2Fdeletefile/1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("handler should not issue request with encoded blocked path")
			})
			_, err := handleCallYunxiaoAPI(context.Background(), client, map[string]any{
				"path": tt.path,
			})
			if err == nil || !strings.Contains(err.Error(), "blocked") {
				t.Fatalf("expected blocked error, got %v", err)
			}
		})
	}
}

func TestHandleCallYunxiaoAPIRejectsColonActionBlockedPath(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "raw colon action", path: "/appstack/organizations/org-1/apps/app-1/jobs/job-1:execute"},
		{name: "encoded colon action", path: "/appstack/organizations/org-1/apps/app-1/jobs/job-1%3Aexecute"},
		{name: "raw finish action", path: "/appstack/organizations/org-1/changeRequests/cr-1:finish"},
		{name: "encoded retry action", path: "/appstack/organizations/org-1/releaseWorkflows/rw-1%3Aretry"},
		{name: "raw skip action", path: "/appstack/organizations/org-1/releaseWorkflows/rw-1:skip"},
		{name: "unknown me action", path: "/platform/groups:me"},
		{name: "encoded unknown me action", path: "/platform/groups%3Ame"},
		{name: "mixed case users me action", path: "/Platform/Users:Me"},
		{name: "encoded mixed case users me action", path: "/platform/Users%3Ame"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("handler should not issue request with colon action path")
			})
			_, err := handleCallYunxiaoAPI(context.Background(), client, map[string]any{
				"path": tt.path,
			})
			if err == nil || !strings.Contains(err.Error(), "blocked") {
				t.Fatalf("expected blocked error, got %v", err)
			}
		})
	}
}

func TestHandleCallYunxiaoAPIRejectsNestedEscapedPath(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "double encoded slash", path: "/projex/organizations/org-1/workitems/wi-1%252Fdeletefile/1"},
		{name: "double encoded mutation term", path: "/projex/organizations/org-1/workitems/wi-1/%2564eletefile/1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("handler should not issue request with nested escaped path")
			})
			_, err := handleCallYunxiaoAPI(context.Background(), client, map[string]any{
				"path": tt.path,
			})
			if err == nil || !strings.Contains(err.Error(), "nested path escapes") {
				t.Fatalf("expected nested path escape error, got %v", err)
			}
		})
	}
}

func TestHandleCallYunxiaoAPIRejectsDotSegmentPath(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "raw dot segment", path: "/projex/organizations/org-1/result/list/../../workitems"},
		{name: "encoded dot segment", path: "/projex/organizations/org-1/result/list/%2e%2e/%2e%2e/workitems"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("handler should not issue request with dot-segment path")
			})
			_, err := handleCallYunxiaoAPI(context.Background(), client, map[string]any{
				"path":   tt.path,
				"method": "POST",
			})
			if err == nil || !strings.Contains(err.Error(), "dot segments") {
				t.Fatalf("expected dot segment error, got %v", err)
			}
		})
	}
}

func TestHandleCallYunxiaoAPIRejectsMutationPostPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request with mutation POST path")
	})
	_, err := handleCallYunxiaoAPI(context.Background(), client, map[string]any{
		"path":   "/projex/organizations/org-1/workitems",
		"method": "POST",
		"body":   `{}`,
	})
	if err == nil || !strings.Contains(err.Error(), "only allow read-only search/list endpoints") {
		t.Fatalf("expected read-only POST error, got %v", err)
	}
}

func TestHandleCallYunxiaoAPIRejectsInvalidEscapedPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request with invalid escaped path")
	})
	_, err := handleCallYunxiaoAPI(context.Background(), client, map[string]any{
		"path": "/projex/organizations/org-1/%zz",
	})
	if err == nil || !strings.Contains(err.Error(), "invalid path escape") {
		t.Fatalf("expected invalid path escape error, got %v", err)
	}
}

func TestHandleCallYunxiaoAPIRejectsInvalidQueryParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request with invalid queryParams")
	})
	_, err := handleCallYunxiaoAPI(context.Background(), client, map[string]any{
		"path":        "/test",
		"queryParams": "not-json",
	})
	if err == nil || !strings.Contains(err.Error(), "invalid queryParams JSON") {
		t.Fatalf("expected queryParams error, got %v", err)
	}
}

func TestHandleCallYunxiaoAPIRejectsInvalidBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not issue request with invalid body")
	})
	_, err := handleCallYunxiaoAPI(context.Background(), client, map[string]any{
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
