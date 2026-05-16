package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestHandleCreateChangeRequest(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/oapi/v1/codeup/organizations/org-1/repositories/repo-1/changeRequests" {
			t.Fatalf("path = %q", r.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["title"] != "Test CR" {
			t.Fatalf("title = %v", body["title"])
		}
		if body["sourceBranch"] != "feature/x" {
			t.Fatalf("sourceBranch = %v", body["sourceBranch"])
		}

		_, _ = w.Write([]byte(`{"id":"1","title":"Test CR"}`))
	})

	result, err := handleCreateChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"title":          "Test CR",
		"sourceBranch":   "feature/x",
		"targetBranch":   "main",
	})
	if err != nil {
		t.Fatalf("handleCreateChangeRequest() error = %v", err)
	}
	if !strings.Contains(result, "Test CR") {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleCreateChangeRequestMissingTitle(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleCreateChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"sourceBranch":   "feature/x",
		"targetBranch":   "main",
	})
	if err == nil {
		t.Fatal("expected error for missing title")
	}
}

func TestHandleCreateChangeRequestNilClient(t *testing.T) {
	_, err := handleCreateChangeRequest(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"title":          "Test",
		"sourceBranch":   "feature/x",
		"targetBranch":   "main",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestHandleCreateChangeRequestMissingSourceBranch(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleCreateChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"title":          "Test",
		"targetBranch":   "main",
	})
	if err == nil {
		t.Fatal("expected error for missing sourceBranch")
	}
}

func TestHandleCreateChangeRequestMissingTargetBranch(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleCreateChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"title":          "Test",
		"sourceBranch":   "feature/x",
	})
	if err == nil {
		t.Fatal("expected error for missing targetBranch")
	}
}

func TestHandleCreateChangeRequestAPIError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleCreateChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"title":          "Test",
		"sourceBranch":   "feature/x",
		"targetBranch":   "main",
	})
	if err == nil {
		t.Fatal("expected error for API failure")
	}
}

func TestHandleCreateChangeRequestMissingOrganizationId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleCreateChangeRequest(context.Background(), client, map[string]any{
		"repositoryId": "repo-1",
		"title":        "Test",
		"sourceBranch": "feature/x",
		"targetBranch": "main",
	})
	if err == nil {
		t.Fatal("expected error for missing organizationId")
	}
}

func TestHandleCreateChangeRequestMissingRepositoryId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleCreateChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"title":          "Test",
		"sourceBranch":   "feature/x",
		"targetBranch":   "main",
	})
	if err == nil {
		t.Fatal("expected error for missing repositoryId")
	}
}

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

func TestHandleAddChangeRequestCommentMissingOrganizationId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleAddChangeRequestComment(context.Background(), client, map[string]any{
		"repositoryId": "repo-1",
		"localId":      "1",
		"content":      "test",
	})
	if err == nil {
		t.Fatal("expected error for missing organizationId")
	}
}

func TestHandleAddChangeRequestCommentMissingRepositoryId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleAddChangeRequestComment(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"localId":        "1",
		"content":        "test",
	})
	if err == nil {
		t.Fatal("expected error for missing repositoryId")
	}
}

func TestHandleAddChangeRequestCommentMissingLocalId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleAddChangeRequestComment(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"content":        "test",
	})
	if err == nil {
		t.Fatal("expected error for missing localId")
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

func TestHandleCloseChangeRequest(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/oapi/v1/codeup/organizations/org-1/repositories/repo-1/changeRequests/1/close" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"result":true}`))
	})

	result, err := handleCloseChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
	})
	if err != nil {
		t.Fatalf("handleCloseChangeRequest() error = %v", err)
	}
	if !strings.Contains(result, "true") {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleCloseChangeRequestMissingOrganizationId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleCloseChangeRequest(context.Background(), client, map[string]any{
		"repositoryId": "repo-1",
		"localId":      "1",
	})
	if err == nil {
		t.Fatal("expected error for missing organizationId")
	}
}

func TestHandleCloseChangeRequestMissingRepositoryId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleCloseChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"localId":        "1",
	})
	if err == nil {
		t.Fatal("expected error for missing repositoryId")
	}
}

func TestHandleCloseChangeRequestMissingLocalId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleCloseChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
	})
	if err == nil {
		t.Fatal("expected error for missing localId")
	}
}

func TestHandleCloseChangeRequestNilClient(t *testing.T) {
	_, err := handleCloseChangeRequest(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestHandleCloseChangeRequestAPIError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleCloseChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
	})
	if err == nil {
		t.Fatal("expected error for API failure")
	}
}

func TestHandleReopenChangeRequest(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/oapi/v1/codeup/organizations/org-1/repositories/repo-1/changeRequests/1/reopen" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"result":true}`))
	})

	result, err := handleReopenChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
	})
	if err != nil {
		t.Fatalf("handleReopenChangeRequest() error = %v", err)
	}
	if !strings.Contains(result, "true") {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleReopenChangeRequestMissingOrganizationId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleReopenChangeRequest(context.Background(), client, map[string]any{
		"repositoryId": "repo-1",
		"localId":      "1",
	})
	if err == nil {
		t.Fatal("expected error for missing organizationId")
	}
}

func TestHandleReopenChangeRequestMissingRepositoryId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleReopenChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"localId":        "1",
	})
	if err == nil {
		t.Fatal("expected error for missing repositoryId")
	}
}

func TestHandleReopenChangeRequestMissingLocalId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleReopenChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
	})
	if err == nil {
		t.Fatal("expected error for missing localId")
	}
}

func TestHandleReopenChangeRequestNilClient(t *testing.T) {
	_, err := handleReopenChangeRequest(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestHandleReopenChangeRequestAPIError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleReopenChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
	})
	if err == nil {
		t.Fatal("expected error for API failure")
	}
}

func TestHandleMergeChangeRequest(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/oapi/v1/codeup/organizations/org-1/repositories/repo-1/changeRequests/1/merge" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"result":true}`))
	})

	result, err := handleMergeChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
	})
	if err != nil {
		t.Fatalf("handleMergeChangeRequest() error = %v", err)
	}
	if !strings.Contains(result, "true") {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleMergeChangeRequestMissingOrganizationId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleMergeChangeRequest(context.Background(), client, map[string]any{
		"repositoryId": "repo-1",
		"localId":      "1",
	})
	if err == nil {
		t.Fatal("expected error for missing organizationId")
	}
}

func TestHandleMergeChangeRequestMissingRepositoryId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleMergeChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"localId":        "1",
	})
	if err == nil {
		t.Fatal("expected error for missing repositoryId")
	}
}

func TestHandleMergeChangeRequestMissingLocalId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleMergeChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
	})
	if err == nil {
		t.Fatal("expected error for missing localId")
	}
}

func TestHandleMergeChangeRequestNilClient(t *testing.T) {
	_, err := handleMergeChangeRequest(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestHandleMergeChangeRequestAPIError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleMergeChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
		"localId":        "1",
	})
	if err == nil {
		t.Fatal("expected error for API failure")
	}
}
