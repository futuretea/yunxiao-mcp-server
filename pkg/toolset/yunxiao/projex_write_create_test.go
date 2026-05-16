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
