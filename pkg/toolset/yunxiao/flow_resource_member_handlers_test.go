package yunxiao

import (
	"context"
	"net/http"
	"testing"
)

func TestHandleListResourceMembersBuildsEscapedPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/flow/organizations/org-1/resourceMembers/resourceTypes/hostGroup/resourceIds/group%2F1" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"username":"alice","userId":"u-1","role":"owner"}]`))
	})

	if _, err := handleListResourceMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"resourceType":   "hostGroup",
		"resourceId":     "group/1",
	}); err != nil {
		t.Fatalf("handleListResourceMembers() error = %v", err)
	}
}

func TestHandleListResourceMembersRequiresParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleListResourceMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"resourceType":   "pipeline",
	}); err == nil {
		t.Fatal("handleListResourceMembers() expected missing resourceId error")
	}
	if _, err := handleListResourceMembers(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing organizationId error")
	}
	if _, err := handleListResourceMembers(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing resourceType error")
	}
	if _, err := handleListResourceMembers(context.Background(), "invalid-client", map[string]any{}); err == nil {
		t.Fatal("expected getClient error")
	}
}
