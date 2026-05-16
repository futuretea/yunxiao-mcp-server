package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

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
