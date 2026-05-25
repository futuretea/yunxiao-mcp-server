package yunxiao

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func TestBuildToolCatalogFiltersDomains(t *testing.T) {
	tools, err := BuildToolCatalog(nil, ToolCatalogOptions{
		ReadOnly:       true,
		EnabledDomains: []string{"platform"},
	})
	if err != nil {
		t.Fatalf("BuildToolCatalog() error = %v", err)
	}
	for _, tool := range tools {
		if tool.Domain != "platform" {
			t.Fatalf("domain = %q, want platform", tool.Domain)
		}
	}
}

func TestBuildToolCatalogDisablesDomains(t *testing.T) {
	tools, err := BuildToolCatalog(nil, ToolCatalogOptions{
		ReadOnly:        true,
		DisabledDomains: []string{"platform"},
	})
	if err != nil {
		t.Fatalf("BuildToolCatalog() error = %v", err)
	}
	for _, tool := range tools {
		if tool.Domain == "platform" {
			t.Fatalf("platform tool %q should be disabled", tool.Tool.Name)
		}
	}
}

func TestBuildToolCatalogAppliesCompactBeforeToolValidation(t *testing.T) {
	_, err := BuildToolCatalog(nil, ToolCatalogOptions{
		ReadOnly:     true,
		CompactMode:  true,
		EnabledTools: []string{"get_project"},
	})
	if err == nil {
		t.Fatal("BuildToolCatalog() expected compact-hidden enabled tool error")
	}
}

func TestBuildToolCatalogAllowsCompactDisabledHiddenTool(t *testing.T) {
	if _, err := BuildToolCatalog(nil, ToolCatalogOptions{
		ReadOnly:      true,
		CompactMode:   true,
		DisabledTools: []string{"get_project"},
	}); err == nil {
		t.Fatal("BuildToolCatalog() expected compact-hidden disabled tool error")
	}
}

func TestBuildToolCatalogDisablesWriteToolsByDefault(t *testing.T) {
	_, err := BuildToolCatalog(nil, ToolCatalogOptions{
		ReadOnly:     true,
		EnabledTools: []string{"create_workitem"},
	})
	if err == nil {
		t.Fatal("BuildToolCatalog() expected read-only write tool error")
	}
}

func TestFindTool(t *testing.T) {
	tools, err := BuildToolCatalog(nil, ToolCatalogOptions{
		ReadOnly:     true,
		EnabledTools: []string{"get_current_user"},
	})
	if err != nil {
		t.Fatalf("BuildToolCatalog() error = %v", err)
	}

	if tool, ok := FindTool(tools, "get_current_user"); !ok || tool.Tool.Name != "get_current_user" {
		t.Fatalf("FindTool(get_current_user) = %#v, %v", tool, ok)
	}
	if _, ok := FindTool(tools, "missing_tool"); ok {
		t.Fatal("FindTool(missing_tool) should return false")
	}
}

func TestIsWriteTool(t *testing.T) {
	if !IsWriteTool("create_workitem") {
		t.Fatal("create_workitem should be a write tool")
	}
	if IsWriteTool("get_current_user") {
		t.Fatal("get_current_user should not be a write tool")
	}
}

func TestValidateToolRequiredParams(t *testing.T) {
	mockTool := toolset.ServerTool{
		Tool: mcp.NewTool("mock_tool",
			mcp.WithString("organizationId", mcp.Required()),
			mcp.WithString("projectId", mcp.Required()),
		),
	}

	err := ValidateToolRequiredParams(mockTool, map[string]any{"organizationId": "org-1"}, "organizationId")
	if err == nil {
		t.Fatal("ValidateToolRequiredParams() expected projectId error")
	}
	if !strings.Contains(err.Error(), "projectId is required") {
		t.Fatalf("error = %v", err)
	}

	if err := ValidateToolRequiredParams(mockTool, map[string]any{"projectId": "project-1"}, "organizationId"); err != nil {
		t.Fatalf("ValidateToolRequiredParams() error = %v", err)
	}
}

func TestValidateToolRequiredParamsUsesRawInputSchema(t *testing.T) {
	mockTool := toolset.ServerTool{
		Tool: mcp.Tool{
			Name:           "raw_schema_tool",
			RawInputSchema: json.RawMessage(`{"type":"object","required":["path"]}`),
		},
	}

	err := ValidateToolRequiredParams(mockTool, map[string]any{"path": " "})
	if err == nil {
		t.Fatal("ValidateToolRequiredParams() expected raw schema required param error")
	}
	if !strings.Contains(err.Error(), "path is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestValidateToolRequiredParamsHandlesNoRequiredInvalidRawAndNonStringValues(t *testing.T) {
	if err := ValidateToolRequiredParams(toolset.ServerTool{Tool: mcp.NewTool("no_required")}, nil); err != nil {
		t.Fatalf("ValidateToolRequiredParams() error = %v", err)
	}

	invalidRawTool := toolset.ServerTool{
		Tool: mcp.Tool{
			Name:           "invalid_raw_schema_tool",
			RawInputSchema: json.RawMessage(`{`),
		},
	}
	if err := ValidateToolRequiredParams(invalidRawTool, nil); err != nil {
		t.Fatalf("ValidateToolRequiredParams() invalid raw schema error = %v", err)
	}

	numberTool := toolset.ServerTool{
		Tool: mcp.NewTool("number_tool", mcp.WithNumber("page", mcp.Required())),
	}
	if err := ValidateToolRequiredParams(numberTool, map[string]any{"page": 1}); err != nil {
		t.Fatalf("ValidateToolRequiredParams() number param error = %v", err)
	}
	if err := ValidateToolRequiredParams(numberTool, nil); err == nil {
		t.Fatal("ValidateToolRequiredParams() expected nil params required error")
	}
}

func TestInvokeToolFillsDefaultOrganizationID(t *testing.T) {
	client := &Client{DefaultOrgID: "default-org"}
	var got map[string]any
	mockTool := toolset.ServerTool{
		Tool: mcp.NewTool("mock_tool"),
		Handler: func(ctx context.Context, client any, params map[string]any) (string, error) {
			got = params
			return "ok", nil
		},
	}

	result, err := InvokeTool(context.Background(), client, mockTool, map[string]any{})
	if err != nil {
		t.Fatalf("InvokeTool() error = %v", err)
	}
	if result != "ok" {
		t.Fatalf("result = %q, want ok", result)
	}
	if got["organizationId"] != "default-org" {
		t.Fatalf("organizationId = %q, want default-org", got["organizationId"])
	}
}

func TestInvokeToolPreservesProvidedOrganizationID(t *testing.T) {
	client := &Client{DefaultOrgID: "default-org"}
	var got map[string]any
	mockTool := toolset.ServerTool{
		Tool: mcp.NewTool("mock_tool"),
		Handler: func(ctx context.Context, client any, params map[string]any) (string, error) {
			got = params
			return "ok", nil
		},
	}

	_, err := InvokeTool(context.Background(), client, mockTool, map[string]any{"organizationId": "provided-org"})
	if err != nil {
		t.Fatalf("InvokeTool() error = %v", err)
	}
	if got["organizationId"] != "provided-org" {
		t.Fatalf("organizationId = %q, want provided-org", got["organizationId"])
	}
}

func TestInvokeToolHandlesNilParams(t *testing.T) {
	var got map[string]any
	mockTool := toolset.ServerTool{
		Tool: mcp.NewTool("mock_tool"),
		Handler: func(ctx context.Context, client any, params map[string]any) (string, error) {
			got = params
			return "ok", nil
		},
	}

	_, err := InvokeTool(context.Background(), &Client{}, mockTool, nil)
	if err != nil {
		t.Fatalf("InvokeTool() error = %v", err)
	}
	if got == nil {
		t.Fatal("params should be initialized")
	}
}
