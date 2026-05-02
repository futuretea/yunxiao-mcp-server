package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetPodContainerLogBuildsPathAndDefaultQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/cluster%2F1/ns%2F1/pods/pod%2F1/containers/app%2F1:logs?tailingLines=1000" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte("line1\nline2"))
	})

	got, err := handleGetPodContainerLog(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"resourcePath":   "cluster/1",
		"namespace":      "ns/1",
		"name":           "pod/1",
		"container":      "app/1",
	})
	if err != nil {
		t.Fatalf("handleGetPodContainerLog() error = %v", err)
	}
	if got != "line1\nline2" {
		t.Fatalf("handleGetPodContainerLog() = %q", got)
	}
}

func TestHandleGetPodContainerLogOverridesTailingLines(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Query().Get("tailingLines") != "200" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`"ok"`))
	})

	if _, err := handleGetPodContainerLog(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"resourcePath":   "cluster-1",
		"namespace":      "default",
		"name":           "pod-1",
		"container":      "app",
		"tailingLines":   float64(200),
	}); err != nil {
		t.Fatalf("handleGetPodContainerLog() error = %v", err)
	}
}

func TestHandleGetPodInfoBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/appstack/organizations/org-1/cluster%2F1/ns%2F1/pods/pod%2F1/info?") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		if r.URL.Query().Get("taskSn") != "task/1" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"name":"pod-1"}`))
	})

	if _, err := handleGetPodInfo(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"resourcePath":   "cluster/1",
		"namespace":      "ns/1",
		"name":           "pod/1",
		"taskSn":         "task/1",
	}); err != nil {
		t.Fatalf("handleGetPodInfo() error = %v", err)
	}
}

func TestHandleGetKubernetesObjectInfoBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/cluster%2F1/ns%2F1/Deployment/deploy%2F1/info?taskSn=task%2F1" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"name":"deploy-1"}`))
	})

	if _, err := handleGetKubernetesObjectInfo(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"resourcePath":   "cluster/1",
		"namespace":      "ns/1",
		"kind":           "Deployment",
		"name":           "deploy/1",
		"taskSn":         "task/1",
	}); err != nil {
		t.Fatalf("handleGetKubernetesObjectInfo() error = %v", err)
	}
}

func TestHandleGetDeploymentRevisionInfoBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/apps/app%2F1/envs/env%2F1/ns/ns%2F1/deployments/deploy%2F1/revisions/rev%2F1?taskSn=task%2F1" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"name":"deploy-1"}`))
	})

	if _, err := handleGetDeploymentRevisionInfo(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"appName":        "app/1",
		"envName":        "env/1",
		"namespace":      "ns/1",
		"name":           "deploy/1",
		"revision":       "rev/1",
		"taskSn":         "task/1",
	}); err != nil {
		t.Fatalf("handleGetDeploymentRevisionInfo() error = %v", err)
	}
}

func TestRequiredResourceProxyObjectRequiresResourcePath(t *testing.T) {
	_, _, _, _, err := requiredResourceProxyObject(map[string]any{"organizationId": "org-1"})
	if err == nil {
		t.Fatal("expected missing resourcePath error")
	}
}

func TestRequiredResourceProxyObjectRequiresNamespace(t *testing.T) {
	_, _, _, _, err := requiredResourceProxyObject(map[string]any{"organizationId": "org-1", "resourcePath": "rp-1"})
	if err == nil {
		t.Fatal("expected missing namespace error")
	}
}

func TestRequiredResourceProxyObjectRequiresName(t *testing.T) {
	_, _, _, _, err := requiredResourceProxyObject(map[string]any{"organizationId": "org-1", "resourcePath": "rp-1", "namespace": "ns-1"})
	if err == nil {
		t.Fatal("expected missing name error")
	}
}

func TestRequiredDeploymentRevisionRequiresNamespace(t *testing.T) {
	_, _, _, _, _, _, err := requiredDeploymentRevision(map[string]any{"organizationId": "org-1", "appName": "app-1", "envName": "env-1"})
	if err == nil {
		t.Fatal("expected missing namespace error")
	}
}

func TestRequiredDeploymentRevisionRequiresName(t *testing.T) {
	_, _, _, _, _, _, err := requiredDeploymentRevision(map[string]any{"organizationId": "org-1", "appName": "app-1", "envName": "env-1", "namespace": "ns-1"})
	if err == nil {
		t.Fatal("expected missing name error")
	}
}

func TestRequiredDeploymentRevisionRequiresRevision(t *testing.T) {
	_, _, _, _, _, _, err := requiredDeploymentRevision(map[string]any{"organizationId": "org-1", "appName": "app-1", "envName": "env-1", "namespace": "ns-1", "name": "deploy-1"})
	if err == nil {
		t.Fatal("expected missing revision error")
	}
}

func TestAppstackResourceProxyHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetPodContainerLog(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleGetPodContainerLog(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "resourcePath": "rp-1", "namespace": "ns-1", "name": "pod-1", "container": "app"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetPodInfo(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleGetPodInfo(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "resourcePath": "rp-1", "namespace": "ns-1", "name": "pod-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetKubernetesObjectInfo(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleGetKubernetesObjectInfo(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "resourcePath": "rp-1", "namespace": "ns-1", "name": "obj-1", "kind": "Deployment"}); err == nil {
		t.Fatal("expected getClient error")
	}
	if _, err := handleGetDeploymentRevisionInfo(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleGetDeploymentRevisionInfo(context.Background(), "invalid-client", map[string]any{"organizationId": "org-1", "appName": "app-1", "envName": "env-1", "namespace": "ns-1", "name": "deploy-1", "revision": "rev-1"}); err == nil {
		t.Fatal("expected getClient error")
	}
}

func TestRequiredResourceProxyObjectRequiresOrganizationId(t *testing.T) {
	_, _, _, _, err := requiredResourceProxyObject(map[string]any{})
	if err == nil {
		t.Fatal("expected missing organizationId error")
	}
}

func TestRequiredResourceProxyObjectRequiresResourcePathOnly(t *testing.T) {
	_, _, _, _, err := requiredResourceProxyObject(map[string]any{"organizationId": "org-1"})
	if err == nil {
		t.Fatal("expected missing resourcePath error")
	}
}

func TestRequiredDeploymentRevisionRequiresAppEnvironment(t *testing.T) {
	_, _, _, _, _, _, err := requiredDeploymentRevision(map[string]any{})
	if err == nil {
		t.Fatal("expected missing organizationId error")
	}
}
