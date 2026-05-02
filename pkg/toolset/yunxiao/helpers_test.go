package yunxiao

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetClientReturnsErrorForNil(t *testing.T) {
	if _, err := getClient(nil); err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestGetClientReturnsErrorForWrongType(t *testing.T) {
	if _, err := getClient("not-a-client"); err == nil {
		t.Fatal("expected error for wrong type")
	}
}

func TestGetClientReturnsClientForValidValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	client, err := NewClient(server.URL, "token", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	got, err := getClient(client)
	if err != nil {
		t.Fatalf("getClient() error = %v", err)
	}
	if got != client {
		t.Fatal("getClient() returned different client")
	}
}

func TestRequiredOrganizationAndRepositoryRequiresRepositoryId(t *testing.T) {
	_, _, err := requiredOrganizationAndRepository(map[string]any{"organizationId": "org-1"})
	if err == nil {
		t.Fatal("expected missing repositoryId error")
	}
}

func TestRequiredOrganizationRepositoryAndLocalIDRequiresLocalId(t *testing.T) {
	_, _, _, err := requiredOrganizationRepositoryAndLocalID(map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "repo-1",
	})
	if err == nil {
		t.Fatal("expected missing localId error")
	}
}

func TestRequiredOrganizationAndPipelineRequiresPipelineId(t *testing.T) {
	_, _, err := requiredOrganizationAndPipeline(map[string]any{"organizationId": "org-1"})
	if err == nil {
		t.Fatal("expected missing pipelineId error")
	}
}
