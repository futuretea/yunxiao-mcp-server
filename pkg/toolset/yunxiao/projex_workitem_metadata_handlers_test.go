package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListWorkitemAttachmentsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/workitems/workitem-1/attachments" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"id":"file-1"}]`))
	})

	if _, err := handleListWorkitemAttachments(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "workitem-1",
	}); err != nil {
		t.Fatalf("handleListWorkitemAttachments() error = %v", err)
	}
}

func TestHandleGetWorkitemFileBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/projex/organizations/org-1/workitems/workitem-1/files/file-1" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"id":"file-1"}`))
	})

	if _, err := handleGetWorkitemFile(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"workitemId":     "workitem-1",
		"id":             "file-1",
	}); err != nil {
		t.Fatalf("handleGetWorkitemFile() error = %v", err)
	}
}

func TestHandleListWorkitemRelationRecordsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/workitems/workitem-1/relationRecords" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("relationType") != "ASSOCIATED" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[{"id":"relation-1"}]`))
	})

	if _, err := handleListWorkitemRelationRecords(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "workitem-1",
		"relationType":   "ASSOCIATED",
	}); err != nil {
		t.Fatalf("handleListWorkitemRelationRecords() error = %v", err)
	}
}

func TestHandleListWorkitemRelationRecordsRequiresRelationType(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	_, err := handleListWorkitemRelationRecords(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "workitem-1",
	})
	if err == nil || !strings.Contains(err.Error(), "relationType is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestHandleListLabelsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/projex/organizations/org-1/projects/project-1/labels" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "20" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"label-1"}]`))
	})

	result, err := handleListLabels(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"id":             "project-1",
		"page":           float64(2),
		"perPage":        float64(20),
	})
	if err != nil {
		t.Fatalf("handleListLabels() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}
