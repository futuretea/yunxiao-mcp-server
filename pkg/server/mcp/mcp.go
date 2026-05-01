package mcp

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
	"github.com/futuretea/yunxiao-mcp-server/pkg/core/version"
	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
	yunxiaoToolset "github.com/futuretea/yunxiao-mcp-server/pkg/toolset/yunxiao"
)

// Configuration holds server startup settings.
type Configuration struct {
	*config.StaticConfig
}

// Server owns the MCP server and Yunxiao client.
type Server struct {
	configuration *Configuration
	server        *server.MCPServer
	client        *yunxiaoToolset.Client
	enabledTools  []string
}

// NewServer creates and configures an MCP server.
func NewServer(configuration Configuration) (*Server, error) {
	if configuration.StaticConfig == nil {
		return nil, fmt.Errorf("static config is required")
	}

	client, err := yunxiaoToolset.NewClient(
		configuration.BaseURL,
		configuration.AccessToken,
		time.Duration(configuration.RequestTimeoutSeconds)*time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("create Yunxiao client: %w", err)
	}

	s := &Server{
		configuration: &configuration,
		server: server.NewMCPServer(
			version.BinaryName,
			version.Version,
			server.WithToolCapabilities(true),
			server.WithLogging(),
		),
		client: client,
	}

	if err := s.registerTools(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) registerTools() error {
	yunxiaoTools := (&yunxiaoToolset.Toolset{ReadOnly: s.configuration.ReadOnly}).GetTools(s.client)
	if err := validateToolFilters(yunxiaoTools, s.configuration.EnabledTools, s.configuration.DisabledTools); err != nil {
		return err
	}

	for _, tool := range yunxiaoTools {
		if !s.shouldEnableTool(tool.Tool.Name) {
			continue
		}
		s.registerTool(tool)
	}
	if len(s.enabledTools) == 0 {
		return fmt.Errorf("no MCP tools enabled; check enabled_tools and disabled_tools")
	}

	log.Info().Int("count", len(s.enabledTools)).Msg("registered MCP tools")
	return nil
}

func validateToolFilters(tools []toolset.ServerTool, enabledTools, disabledTools []string) error {
	known := make(map[string]struct{}, len(tools))
	for _, tool := range tools {
		name := tool.Tool.Name
		if _, exists := known[name]; exists {
			return fmt.Errorf("duplicate MCP tool registered: %s", name)
		}
		known[name] = struct{}{}
	}

	for _, name := range append(append([]string{}, enabledTools...), disabledTools...) {
		if _, exists := known[name]; !exists {
			return fmt.Errorf("unknown MCP tool %q; known tools: %s", name, strings.Join(knownToolNames(known), ", "))
		}
	}
	return nil
}

func knownToolNames(known map[string]struct{}) []string {
	names := make([]string, 0, len(known))
	for name := range known {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func (s *Server) shouldEnableTool(toolName string) bool {
	if slices.Contains(s.configuration.DisabledTools, toolName) {
		return false
	}
	if len(s.configuration.EnabledTools) > 0 {
		return slices.Contains(s.configuration.EnabledTools, toolName)
	}
	return true
}

func (s *Server) registerTool(tool toolset.ServerTool) {
	handler := server.ToolHandlerFunc(func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, _ := request.Params.Arguments.(map[string]any)
		if params == nil {
			params = map[string]any{}
		}

		result, err := tool.Handler(ctx, s.client, params)
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
		server.WithStateLess(true),
	)
}

// GetEnabledTools returns registered tool names.
func (s *Server) GetEnabledTools() []string {
	return append([]string(nil), s.enabledTools...)
}

// IsHealthy reports whether the server has a configured API client and registered tools.
func (s *Server) IsHealthy() bool {
	return s != nil &&
		s.client != nil &&
		s.configuration != nil &&
		s.configuration.AccessToken != "" &&
		len(s.enabledTools) > 0
}

// Close releases server resources.
func (s *Server) Close() {}

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
