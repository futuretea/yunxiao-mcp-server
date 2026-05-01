package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetUserBuildsEscapedPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/platform/users/user%2F1" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"id":"user-1"}`))
	})

	if _, err := handleGetUser(context.Background(), client, map[string]any{
		"idOrUsername": "user/1",
	}); err != nil {
		t.Fatalf("handleGetUser() error = %v", err)
	}
}

func TestHandleGetUserPreservesEncodedPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/oapi/v1/platform/users/user%2F1" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"id":"user-1"}`))
	})

	if _, err := handleGetUser(context.Background(), client, map[string]any{
		"idOrUsername": "user%2F1",
	}); err != nil {
		t.Fatalf("handleGetUser() error = %v", err)
	}
}

func TestHandleListAppExtensionFeaturesBuildsEscapedPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/platform/organizations/org-1/appExtensions/ext%2F1/features" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"features":[]}`))
	})

	if _, err := handleListAppExtensionFeatures(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"type":           "ext/1",
	}); err != nil {
		t.Fatalf("handleListAppExtensionFeatures() error = %v", err)
	}
}

func TestPlatformMetadataHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetUser(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("handleGetUser() expected missing idOrUsername error")
	}

	_, err := handleListAppExtensionFeatures(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil {
		t.Fatal("handleListAppExtensionFeatures() expected missing type error")
	}
}
