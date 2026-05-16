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

func TestRequiredStringRejectsWhitespaceOnly(t *testing.T) {
	_, err := requiredString(map[string]any{"organizationId": " \t\n "}, "organizationId")
	if err == nil {
		t.Fatal("expected whitespace-only organizationId error")
	}
}

func TestRequiredStringPreservesNonBlankWhitespace(t *testing.T) {
	got, err := requiredString(map[string]any{"content": "  keep padding  "}, "content")
	if err != nil {
		t.Fatalf("requiredString() error = %v", err)
	}
	if got != "  keep padding  " {
		t.Fatalf("requiredString() = %q, want original string", got)
	}
}

func TestRequiredNumberPathStringAcceptsTypes(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want string
	}{
		{"float64", float64(123), "123"},
		{"int", int(456), "456"},
		{"int64", int64(789), "789"},
		{"string", "abc", "abc"},
		{"trimmed string", " 123 ", "123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := requiredNumberPathString(map[string]any{"k": tt.val}, "k")
			if err != nil {
				t.Fatalf("requiredNumberPathString() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("requiredNumberPathString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRequiredNumberPathStringRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name string
		val  any
	}{
		{"empty string", ""},
		{"whitespace string", " \t "},
		{"fractional float64", float64(1.9)},
		{"missing", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := map[string]any{}
			if tt.name != "missing" {
				params["k"] = tt.val
			}
			if _, err := requiredNumberPathString(params, "k"); err == nil {
				t.Fatal("expected error")
			}
		})
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
