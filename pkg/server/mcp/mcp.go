package mcp

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

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

	clientOptions := []yunxiaoToolset.ClientOption{}
	if configuration.InsecureSkipTLSVerify {
		clientOptions = append(clientOptions, yunxiaoToolset.WithInsecureSkipTLSVerify(true))
	}

	client, err := yunxiaoToolset.NewClient(
		configuration.BaseURL,
		configuration.AccessToken,
		time.Duration(configuration.RequestTimeoutSeconds)*time.Second,
		clientOptions...,
	)
	if err != nil {
		return nil, fmt.Errorf("create Yunxiao client: %w", err)
	}

	if err := client.ResolveDefaultOrgID(context.Background()); err != nil {
		log.Warn().Err(err).Msg("failed to resolve default organization")
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
	toolsetBuilder := &yunxiaoToolset.Toolset{ReadOnly: s.configuration.ReadOnly}
	yunxiaoTools := toolsetBuilder.GetTools(s.client)

	// Stage 1: domain filter
	if len(s.configuration.EnabledDomains) > 0 {
		yunxiaoTools = filterToolsByDomains(yunxiaoTools, s.configuration.EnabledDomains, nil)
	} else if len(s.configuration.DisabledDomains) > 0 {
		yunxiaoTools = filterToolsByDomains(yunxiaoTools, nil, s.configuration.DisabledDomains)
	}

	// Stage 2: compact mode — hide raw tools with enhanced alternatives
	if s.configuration.CompactMode {
		yunxiaoTools = toolsetBuilder.GetCompactTools(yunxiaoTools)
	}

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
		return fmt.Errorf("no MCP tools enabled; check enabled_tools, disabled_tools, enable_domains, disable_domains, compact")
	}

	log.Info().Int("count", len(s.enabledTools)).Msg("registered MCP tools")
	return nil
}

func filterToolsByDomains(tools []toolset.ServerTool, enabled, disabled []string) []toolset.ServerTool {
	if len(enabled) > 0 {
		allowed := make(map[string]struct{}, len(enabled))
		for _, d := range enabled {
			allowed[d] = struct{}{}
		}
		filtered := make([]toolset.ServerTool, 0, len(tools))
		for _, tool := range tools {
			if _, ok := allowed[tool.Domain]; ok {
				filtered = append(filtered, tool)
			}
		}
		return filtered
	}

	if len(disabled) > 0 {
		blocked := make(map[string]struct{}, len(disabled))
		for _, d := range disabled {
			blocked[d] = struct{}{}
		}
		filtered := make([]toolset.ServerTool, 0, len(tools))
		for _, tool := range tools {
			if _, ok := blocked[tool.Domain]; !ok {
				filtered = append(filtered, tool)
			}
		}
		return filtered
	}

	return tools
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

// GetEnabledTools returns registered tool names.
func (s *Server) GetEnabledTools() []string {
	return append([]string(nil), s.enabledTools...)
}

// IsHealthy reports whether the server has a configured API client and registered tools.
func (s *Server) IsHealthy() bool {
	return s != nil &&
		s.client != nil &&
		s.configuration != nil &&
		len(s.enabledTools) > 0
}

// Close releases server resources.
func (s *Server) Close() {}
