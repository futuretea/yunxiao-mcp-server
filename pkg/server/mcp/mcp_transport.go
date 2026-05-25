package mcp

import (
	"context"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
	"github.com/futuretea/yunxiao-mcp-server/pkg/core/version"
	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
	yunxiaoToolset "github.com/futuretea/yunxiao-mcp-server/pkg/toolset/yunxiao"
	yunxiaoSDK "github.com/futuretea/yunxiao-mcp-server/pkg/yunxiao"
)

func (s *Server) registerTool(tool toolset.ServerTool) {
	handler := server.ToolHandlerFunc(func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, _ := request.Params.Arguments.(map[string]any)
		if params == nil {
			params = map[string]any{}
		}

		result, err := yunxiaoToolset.InvokeTool(ctx, s.client, tool, params)
		return NewTextResult(result, err), nil
	})

	s.server.AddTool(tool.Tool, handler)
	s.enabledTools = append(s.enabledTools, tool.Tool.Name)
}

// ServeStdio starts the MCP server over stdin/stdout.
func (s *Server) ServeStdio() error {
	return server.ServeStdio(s.server)
}

// ServeSSE creates an SSE MCP HTTP handler.
func (s *Server) ServeSSE(baseURL string, httpServer *http.Server) *server.SSEServer {
	options := []server.SSEOption{
		server.WithHTTPServer(httpServer),
		server.WithAppendQueryToMessageEndpoint(),
		server.WithSSEContextFunc(withRequestAccessToken),
	}
	if baseURL != "" {
		options = append(options, server.WithBaseURL(baseURL))
	}
	return server.NewSSEServer(s.server, options...)
}

// ServeStreamableHTTP creates a streamable HTTP MCP handler.
func (s *Server) ServeStreamableHTTP(httpServer *http.Server) *server.StreamableHTTPServer {
	return server.NewStreamableHTTPServer(
		s.server,
		server.WithStreamableHTTPServer(httpServer),
		server.WithHTTPContextFunc(withRequestAccessToken),
		server.WithStateLess(true),
	)
}

// NewTestServer creates a server with explicit fields for testing.
func NewTestServer(client *yunxiaoSDK.Client, enabledTools []string) *Server {
	return &Server{
		configuration: &Configuration{StaticConfig: &config.StaticConfig{}},
		client:        client,
		server:        server.NewMCPServer(version.BinaryName, version.Version),
		enabledTools:  enabledTools,
	}
}

func withRequestAccessToken(ctx context.Context, r *http.Request) context.Context {
	return yunxiaoSDK.WithAccessToken(ctx, requestAccessToken(r))
}

func requestAccessToken(r *http.Request) string {
	if r == nil {
		return ""
	}
	if accessToken := strings.TrimSpace(r.Header.Get(yunxiaoSDK.AccessTokenHeader)); accessToken != "" {
		return accessToken
	}
	return strings.TrimSpace(r.URL.Query().Get(yunxiaoSDK.AccessTokenQueryParam))
}

// NewTextResult creates a standard MCP text result.
func NewTextResult(content string, err error) *mcp.CallToolResult {
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{Type: "text", Text: err.Error()},
			},
		}
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: content},
		},
	}
}
