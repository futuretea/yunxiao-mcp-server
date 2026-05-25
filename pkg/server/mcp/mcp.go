package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
	"github.com/futuretea/yunxiao-mcp-server/pkg/core/version"
	yunxiaoToolset "github.com/futuretea/yunxiao-mcp-server/pkg/toolset/yunxiao"
	yunxiaoSDK "github.com/futuretea/yunxiao-mcp-server/pkg/yunxiao"
)

// Configuration holds server startup settings.
type Configuration struct {
	*config.StaticConfig
}

// Server owns the MCP server and Yunxiao client.
type Server struct {
	configuration *Configuration
	server        *server.MCPServer
	client        *yunxiaoSDK.Client
	enabledTools  []string
}

// NewServer creates and configures an MCP server.
func NewServer(configuration Configuration) (*Server, error) {
	if configuration.StaticConfig == nil {
		return nil, fmt.Errorf("static config is required")
	}

	clientOptions := []yunxiaoSDK.ClientOption{}
	if configuration.InsecureSkipTLSVerify {
		clientOptions = append(clientOptions, yunxiaoSDK.WithInsecureSkipTLSVerify(true))
	}

	client, err := yunxiaoSDK.NewClient(
		configuration.BaseURL,
		configuration.AccessToken,
		time.Duration(configuration.RequestTimeoutSeconds)*time.Second,
		clientOptions...,
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
	if err := client.ResolveDefaultOrgID(context.Background()); err != nil {
		log.Warn().Err(err).Msg("failed to resolve default organization")
	}
	return s, nil
}

func (s *Server) registerTools() error {
	yunxiaoTools, err := yunxiaoToolset.BuildToolCatalog(s.client, yunxiaoToolset.ToolCatalogOptions{
		ReadOnly:        s.configuration.ReadOnly,
		CompactMode:     s.configuration.CompactMode,
		EnabledTools:    s.configuration.EnabledTools,
		DisabledTools:   s.configuration.DisabledTools,
		EnabledDomains:  s.configuration.EnabledDomains,
		DisabledDomains: s.configuration.DisabledDomains,
	})
	if err != nil {
		return err
	}

	for _, tool := range yunxiaoTools {
		s.registerTool(tool)
	}

	log.Info().Int("count", len(s.enabledTools)).Msg("registered MCP tools")
	return nil
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
