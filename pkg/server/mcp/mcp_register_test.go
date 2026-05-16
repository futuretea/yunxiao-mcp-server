package mcp

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
	yunxiaoToolset "github.com/futuretea/yunxiao-mcp-server/pkg/toolset/yunxiao"
)

func TestRegisterToolFillsDefaultOrganizationID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"default-org"}]`))
	}))
	defer ts.Close()

	client, err := yunxiaoToolset.NewClient(ts.URL, "token", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if err := client.ResolveDefaultOrgID(context.Background()); err != nil {
		t.Fatalf("ResolveDefaultOrgID() error = %v", err)
	}

	s := &Server{
		configuration: &Configuration{StaticConfig: &config.StaticConfig{}},
		server:        server.NewMCPServer("test", "1.0.0"),
		client:        client,
	}

	var gotParams map[string]any
	mockTool := toolset.ServerTool{
		Tool: mcp.NewTool("mock_tool"),
		Handler: func(ctx context.Context, c any, params map[string]any) (string, error) {
			gotParams = params
			return "ok", nil
		},
	}

	s.registerTool(mockTool)

	registered := s.server.GetTool("mock_tool")
	if registered == nil {
		t.Fatal("mock_tool should be registered")
	}

	_, err = registered.Handler(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "mock_tool",
			Arguments: map[string]any{},
		},
	})
	if err != nil {
		t.Fatalf("handler error = %v", err)
	}

	if gotParams["organizationId"] != "default-org" {
		t.Fatalf("organizationId = %q, want default-org", gotParams["organizationId"])
	}
}

func TestRegisterToolReplacesBlankOrganizationIDWithDefault(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"default-org"}]`))
	}))
	defer ts.Close()

	client, err := yunxiaoToolset.NewClient(ts.URL, "token", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if err := client.ResolveDefaultOrgID(context.Background()); err != nil {
		t.Fatalf("ResolveDefaultOrgID() error = %v", err)
	}

	s := &Server{
		configuration: &Configuration{StaticConfig: &config.StaticConfig{}},
		server:        server.NewMCPServer("test", "1.0.0"),
		client:        client,
	}

	var gotParams map[string]any
	mockTool := toolset.ServerTool{
		Tool: mcp.NewTool("mock_tool"),
		Handler: func(ctx context.Context, c any, params map[string]any) (string, error) {
			gotParams = params
			return "ok", nil
		},
	}

	s.registerTool(mockTool)

	registered := s.server.GetTool("mock_tool")
	if registered == nil {
		t.Fatal("mock_tool should be registered")
	}

	_, err = registered.Handler(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "mock_tool",
			Arguments: map[string]any{"organizationId": " \t "},
		},
	})
	if err != nil {
		t.Fatalf("handler error = %v", err)
	}

	if gotParams["organizationId"] != "default-org" {
		t.Fatalf("organizationId = %q, want default-org", gotParams["organizationId"])
	}
}

func TestRegisterToolPreservesExistingOrganizationID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"default-org"}]`))
	}))
	defer ts.Close()

	client, err := yunxiaoToolset.NewClient(ts.URL, "token", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if err := client.ResolveDefaultOrgID(context.Background()); err != nil {
		t.Fatalf("ResolveDefaultOrgID() error = %v", err)
	}

	s := &Server{
		configuration: &Configuration{StaticConfig: &config.StaticConfig{}},
		server:        server.NewMCPServer("test", "1.0.0"),
		client:        client,
	}

	var gotParams map[string]any
	mockTool := toolset.ServerTool{
		Tool: mcp.NewTool("mock_tool"),
		Handler: func(ctx context.Context, c any, params map[string]any) (string, error) {
			gotParams = params
			return "ok", nil
		},
	}

	s.registerTool(mockTool)

	registered := s.server.GetTool("mock_tool")
	if registered == nil {
		t.Fatal("mock_tool should be registered")
	}

	_, err = registered.Handler(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "mock_tool",
			Arguments: map[string]any{"organizationId": "provided-org"},
		},
	})
	if err != nil {
		t.Fatalf("handler error = %v", err)
	}

	if gotParams["organizationId"] != "provided-org" {
		t.Fatalf("organizationId = %q, want provided-org", gotParams["organizationId"])
	}
}

func TestRegisterToolHandlesNilParams(t *testing.T) {
	s := &Server{
		configuration: &Configuration{StaticConfig: &config.StaticConfig{}},
		server:        server.NewMCPServer("test", "1.0.0"),
		client:        &yunxiaoToolset.Client{},
	}

	var gotParams map[string]any
	mockTool := toolset.ServerTool{
		Tool: mcp.NewTool("mock_tool"),
		Handler: func(ctx context.Context, c any, params map[string]any) (string, error) {
			gotParams = params
			return "ok", nil
		},
	}

	s.registerTool(mockTool)

	registered := s.server.GetTool("mock_tool")
	if registered == nil {
		t.Fatal("mock_tool should be registered")
	}

	_, err := registered.Handler(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "mock_tool",
			Arguments: nil,
		},
	})
	if err != nil {
		t.Fatalf("handler error = %v", err)
	}

	if gotParams == nil {
		t.Fatal("params should not be nil")
	}
	if _, ok := gotParams["organizationId"]; ok {
		t.Fatal("organizationId should not be set when no default org")
	}
}

func TestRegisterToolWrapsHandlerError(t *testing.T) {
	s := &Server{
		configuration: &Configuration{StaticConfig: &config.StaticConfig{}},
		server:        server.NewMCPServer("test", "1.0.0"),
		client:        &yunxiaoToolset.Client{},
	}

	mockTool := toolset.ServerTool{
		Tool: mcp.NewTool("mock_tool"),
		Handler: func(ctx context.Context, c any, params map[string]any) (string, error) {
			return "", fmt.Errorf("handler failed")
		},
	}

	s.registerTool(mockTool)

	registered := s.server.GetTool("mock_tool")
	if registered == nil {
		t.Fatal("mock_tool should be registered")
	}

	result, err := registered.Handler(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "mock_tool",
			Arguments: map[string]any{},
		},
	})
	if err != nil {
		t.Fatalf("handler should not return error, got %v", err)
	}
	if !result.IsError {
		t.Fatal("expected IsError = true")
	}
	if len(result.Content) != 1 {
		t.Fatalf("content count = %d, want 1", len(result.Content))
	}
	text, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("content type = %T, want TextContent", result.Content[0])
	}
	if text.Text != "handler failed" {
		t.Fatalf("text = %q, want handler failed", text.Text)
	}
}
