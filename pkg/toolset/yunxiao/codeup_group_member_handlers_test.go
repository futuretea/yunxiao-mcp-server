package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListGroupMembersBuildsEscapedPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/groups/group%2Fsub/members?accessLevel=30" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"userId":"user-1"}]`))
	})

	if _, err := handleListGroupMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"groupId":        "group/sub",
		"accessLevel":    float64(30),
	}); err != nil {
		t.Fatalf("handleListGroupMembers() error = %v", err)
	}
}

func TestHandleListGroupMembersPreservesEncodedGroupID(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/groups/group%2Fsub/members" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"userId":"user-1"}]`))
	})

	if _, err := handleListGroupMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"groupId":        "group%2Fsub",
	}); err != nil {
		t.Fatalf("handleListGroupMembers() error = %v", err)
	}
}

func TestHandleGetMemberHTTPSCloneUsernameBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/codeup/organizations/org-1/users/user%2F1/httpsCloneUsername" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"username":"alice"}`))
	})

	if _, err := handleGetMemberHTTPSCloneUsername(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"userId":         "user/1",
	}); err != nil {
		t.Fatalf("handleGetMemberHTTPSCloneUsername() error = %v", err)
	}
}

func TestCodeUpGroupMemberHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	_, err := handleListGroupMembers(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil {
		t.Fatal("handleListGroupMembers() expected missing groupId error")
	}

	_, err = handleGetMemberHTTPSCloneUsername(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil {
		t.Fatal("handleGetMemberHTTPSCloneUsername() expected missing userId error")
	}
	if _, err := handleListGroupMembers(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "groupId": "g-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetMemberHTTPSCloneUsername(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "userId": "u-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}
