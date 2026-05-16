package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestHandleCreateWorkitem(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems" {
			t.Fatalf("path = %q", r.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["spaceId"] != "proj-1" {
			t.Fatalf("spaceId = %v", body["spaceId"])
		}
		if body["subject"] != "Test Task" {
			t.Fatalf("subject = %v", body["subject"])
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"wi-1","subject":"Test Task"}`))
	})

	result, err := handleCreateWorkitem(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "proj-1",
		"category":       "Task",
		"workitemTypeId": "type-1",
		"subject":        "Test Task",
		"description":    "A test task",
	})
	if err != nil {
		t.Fatalf("handleCreateWorkitem() error = %v", err)
	}
	if !strings.Contains(result, "wi-1") {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleCreateWorkitemMissingRequiredFields(t *testing.T) {
	_, err := handleCreateWorkitem(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil {
		t.Fatal("expected error for missing required fields")
	}
}

func TestHandleCreateWorkitemAPIError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleCreateWorkitem(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "proj-1",
		"category":       "Task",
		"workitemTypeId": "type-1",
		"subject":        "Test",
	})
	if err == nil {
		t.Fatal("expected error for API failure")
	}
}

func TestHandleCreateWorkitemMissingOrgId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleCreateWorkitem(context.Background(), client, map[string]any{
		"projectId":      "proj-1",
		"category":       "Task",
		"workitemTypeId": "type-1",
		"subject":        "Test",
	})
	if err == nil {
		t.Fatal("expected error for missing organizationId")
	}
}

func TestHandleCreateWorkitemMissingCategory(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleCreateWorkitem(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "proj-1",
		"workitemTypeId": "type-1",
		"subject":        "Test",
	})
	if err == nil {
		t.Fatal("expected error for missing category")
	}
}

func TestHandleCreateWorkitemMissingWorkitemTypeId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleCreateWorkitem(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "proj-1",
		"category":       "Task",
		"subject":        "Test",
	})
	if err == nil {
		t.Fatal("expected error for missing workitemTypeId")
	}
}

func TestHandleCreateWorkitemMissingSubject(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleCreateWorkitem(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectId":      "proj-1",
		"category":       "Task",
		"workitemTypeId": "type-1",
	})
	if err == nil {
		t.Fatal("expected error for missing subject")
	}
}

func TestHandleCreateWorkitemMissingProjectId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleCreateWorkitem(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"category":       "Task",
		"workitemTypeId": "type-1",
		"subject":        "Test",
	})
	if err == nil {
		t.Fatal("expected error for missing projectId")
	}
}

func TestHandleUpdateWorkitem(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems/wi-1" {
			t.Fatalf("path = %q", r.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["subject"] != "Updated Title" {
			t.Fatalf("subject = %v", body["subject"])
		}

		_, _ = w.Write([]byte(`{"id":"wi-1","subject":"Updated Title"}`))
	})

	result, err := handleUpdateWorkitem(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
		"subject":        "Updated Title",
	})
	if err != nil {
		t.Fatalf("handleUpdateWorkitem() error = %v", err)
	}
	if !strings.Contains(result, "Updated Title") {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleUpdateWorkitemMissingOrgId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleUpdateWorkitem(context.Background(), client, map[string]any{
		"workitemId": "wi-1",
		"subject":    "Updated",
	})
	if err == nil {
		t.Fatal("expected error for missing organizationId")
	}
}

func TestHandleUpdateWorkitemMissingWorkitemId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleUpdateWorkitem(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"subject":        "Updated",
	})
	if err == nil {
		t.Fatal("expected error for missing workitemId")
	}
}

func TestHandleUpdateWorkitemNoFields(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleUpdateWorkitem(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
	})
	if err == nil {
		t.Fatal("expected error when no fields to update")
	}
}

func TestHandleUpdateWorkitemAPIError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	_, err := handleUpdateWorkitem(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
		"subject":        "Updated",
	})
	if err == nil {
		t.Fatal("expected error for API failure")
	}
}

func TestHandleUpdateWorkitemStatus(t *testing.T) {
	callCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}

		switch callCount {
		case 1:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems/wi-1/status" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["statusId"] != "status-2" {
				t.Fatalf("statusId = %v", body["statusId"])
			}
			_, _ = w.Write([]byte(`{"success":true}`))
		case 2:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems/wi-1/comments" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["content"] != "Moving to done" {
				t.Fatalf("content = %v", body["content"])
			}
			_, _ = w.Write([]byte(`{"id":"c-1"}`))
		default:
			t.Fatalf("unexpected call %d", callCount)
		}
	})

	result, err := handleUpdateWorkitemStatus(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
		"statusId":       "status-2",
		"comment":        "Moving to done",
	})
	if err != nil {
		t.Fatalf("handleUpdateWorkitemStatus() error = %v", err)
	}
	if !strings.Contains(result, "success") {
		t.Fatalf("result = %q", result)
	}
	if callCount != 2 {
		t.Fatalf("expected 2 API calls, got %d", callCount)
	}
}

func TestHandleUpdateWorkitemStatusWithoutComment(t *testing.T) {
	callCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if r.URL.Path == "/oapi/v1/projex/organizations/org-1/workitems/wi-1/status" {
			_, _ = w.Write([]byte(`{"success":true}`))
			return
		}
		t.Fatalf("unexpected path %q", r.URL.Path)
	})

	_, err := handleUpdateWorkitemStatus(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
		"statusId":       "status-2",
	})
	if err != nil {
		t.Fatalf("handleUpdateWorkitemStatus() error = %v", err)
	}
	if callCount != 1 {
		t.Fatalf("expected 1 API call, got %d", callCount)
	}
}

func TestHandleUpdateWorkitemStatusReturnsCommentError(t *testing.T) {
	callCount := 0
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		switch callCount {
		case 1:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems/wi-1/status" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			_, _ = w.Write([]byte(`{"success":true}`))
		case 2:
			if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems/wi-1/comments" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			http.Error(w, "comment failed", http.StatusInternalServerError)
		default:
			t.Fatalf("unexpected call %d", callCount)
		}
	})

	_, err := handleUpdateWorkitemStatus(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
		"statusId":       "status-2",
		"comment":        "Moving to done",
	})
	if err == nil || !strings.Contains(err.Error(), "add status comment") {
		t.Fatalf("expected comment error, got %v", err)
	}
	if callCount != 2 {
		t.Fatalf("expected 2 API calls, got %d", callCount)
	}
}

func TestHandleAddWorkitemComment(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems/wi-1/comments" {
			t.Fatalf("path = %q", r.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["content"] != "This is a comment" {
			t.Fatalf("content = %v", body["content"])
		}

		_, _ = w.Write([]byte(`{"id":"c-1","content":"This is a comment"}`))
	})

	result, err := handleAddWorkitemComment(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
		"content":        "This is a comment",
	})
	if err != nil {
		t.Fatalf("handleAddWorkitemComment() error = %v", err)
	}
	if !strings.Contains(result, "c-1") {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleAddWorkitemCommentMissingOrganizationId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleAddWorkitemComment(context.Background(), client, map[string]any{
		"workitemId": "wi-1",
		"content":    "test",
	})
	if err == nil {
		t.Fatal("expected error for missing organizationId")
	}
}

func TestHandleAddWorkitemCommentMissingWorkitemId(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleAddWorkitemComment(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"content":        "test",
	})
	if err == nil {
		t.Fatal("expected error for missing workitemId")
	}
}

func TestHandleAddWorkitemCommentMissingContent(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected API call")
	})
	_, err := handleAddWorkitemComment(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
	})
	if err == nil {
		t.Fatal("expected error for missing content")
	}
}

func TestHandleAddWorkitemCommentNilClient(t *testing.T) {
	_, err := handleAddWorkitemComment(context.Background(), nil, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
		"content":        "test",
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestHandleAddWorkitemCommentAPIError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"internal error"}`))
	})
	_, err := handleAddWorkitemComment(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "wi-1",
		"content":        "test",
	})
	if err == nil {
		t.Fatal("expected error for API error response")
	}
}
