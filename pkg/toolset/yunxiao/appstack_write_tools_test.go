package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleCreateChangeOrderBuildsPathAndBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/apps/app-1/changeOrders" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"sn":"co-1","status":"RUNNING"}`))
	})

	result, err := handleCreateChangeOrder(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app-1",
		"changeOrder":    `{"changeOrderName":"release-1","type":"Deploy"}`,
	})
	if err != nil {
		t.Fatalf("handleCreateChangeOrder() error = %v", err)
	}
	if !strings.Contains(result, "co-1") {
		t.Fatalf("result = %q, want sn:co-1", result)
	}
}

func TestHandleCreateChangeOrderRequiresParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleCreateChangeOrder(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleCreateChangeOrder(context.Background(), client, map[string]any{
		"organizationId": "org-1", "appName": "app-1",
	}); err == nil {
		t.Fatal("expected missing changeOrder error")
	}
	if _, err := handleCreateChangeOrder(context.Background(), client, map[string]any{
		"organizationId": "org-1", "appName": "app-1", "changeOrder": "not-json",
	}); err == nil {
		t.Fatal("expected invalid JSON error")
	}
	if _, err := handleCreateChangeOrder(context.Background(), "invalid-client", map[string]any{
		"organizationId": "org-1", "appName": "app-1", "changeOrder": `{}`,
	}); err == nil {
		t.Fatal("expected getClient error")
	}
}
