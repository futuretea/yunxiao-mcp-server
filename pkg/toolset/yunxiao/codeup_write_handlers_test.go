package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestHandleAddChangeRequestComment(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/oapi/v1/codeup/organizations/org-1/repositories/repo-1/changeRequests/1/comments" {
			t.Fatalf("path = %q", r.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["content"] != "test comment" {
			t.Fatalf("content = %v", body["content"])
		}

		_, _ = w.Write([]byte(`{"id":"c-1","content":"test comment"}`))
	})

	result, err := handleAddChangeRequestComment(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
		"content":        "test comment",
	})
	if err != nil {
		t.Fatalf("handleAddChangeRequestComment() error = %v", err)
	}
	if !strings.Contains(result, "c-1") {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleAddChangeRequestCommentMissingContent(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleAddChangeRequestComment(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
	})
	if err == nil {
		t.Fatal("expected error for missing content")
	}
}

func TestHandleAddChangeRequestCommentNilClient(t *testing.T) {
	_, err := handleAddChangeRequestComment(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
		"content":        "test",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestHandleAddChangeRequestCommentAPIError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleAddChangeRequestComment(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
		"content":        "test",
	})
	if err == nil {
		t.Fatal("expected error for API failure")
	}
}
