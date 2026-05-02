package yunxiao

import (
	"context"
	"net/http"
	"testing"
)

func TestHandleGetMachineDeployLogBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/host/deployLog?machineSn=machine%2F1&tunnelId=123" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	if _, err := handleGetMachineDeployLog(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"tunnelId":       float64(123),
		"machineSn":      "machine/1",
	}); err != nil {
		t.Fatalf("handleGetMachineDeployLog() error = %v", err)
	}
}

func TestHandleGetDeployGroupBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/pools/pool%2F1/deployGroups/group%2F1" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"name":"group-1"}`))
	})

	if _, err := handleGetDeployGroup(context.Background(), client, map[string]any{
		"organizationId":  "org-1",
		"poolName":        "pool/1",
		"deployGroupName": "group/1",
	}); err != nil {
		t.Fatalf("handleGetDeployGroup() error = %v", err)
	}
}

func TestHandleListResourceInstancesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/appstack/organizations/org-1/pools/pool-1/instances" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("pagination") != "keyset" ||
			r.URL.Query().Get("perPage") != "20" ||
			r.URL.Query().Get("orderBy") != "id" ||
			r.URL.Query().Get("sort") != "asc" ||
			r.URL.Query().Get("nextToken") != "token-1" ||
			r.URL.Query().Get("page") != "2" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	if _, err := handleListResourceInstances(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"poolName":       "pool-1",
		"pagination":     "keyset",
		"perPage":        float64(20),
		"orderBy":        "id",
		"sort":           "asc",
		"nextToken":      "token-1",
		"page":           float64(2),
	}); err != nil {
		t.Fatalf("handleListResourceInstances() error = %v", err)
	}
}

func TestHandleGetResourceInstanceBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/appstack/organizations/org-1/pools/pool%2F1/instances/instance%2F1" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"name":"instance-1"}`))
	})

	if _, err := handleGetResourceInstance(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"poolName":       "pool/1",
		"instanceName":   "instance/1",
	}); err != nil {
		t.Fatalf("handleGetResourceInstance() error = %v", err)
	}
}

func TestHandleListResourceInstancesRequiresPoolName(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleListResourceInstances(context.Background(), client, map[string]any{"organizationId": "org-1"}); err == nil {
		t.Fatal("expected missing poolName error")
	}
}

func TestHandleGetDeployGroupRequiresDeployGroupName(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleGetDeployGroup(context.Background(), client, map[string]any{"organizationId": "org-1", "poolName": "pool-1"}); err == nil {
		t.Fatal("expected missing deployGroupName error")
	}
}

func TestHandleGetResourceInstanceRequiresInstanceName(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpected request")
	})
	if _, err := handleGetResourceInstance(context.Background(), client, map[string]any{"organizationId": "org-1", "poolName": "pool-1"}); err == nil {
		t.Fatal("expected missing instanceName error")
	}
}

func TestRequiredOrganizationAndPoolRequiresPoolName(t *testing.T) {
	_, _, err := requiredOrganizationAndPool(map[string]any{"organizationId": "org-1"})
	if err == nil {
		t.Fatal("expected missing poolName error")
	}
}
