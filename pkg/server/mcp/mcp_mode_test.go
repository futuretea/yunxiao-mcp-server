package mcp

import (
	"testing"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
)

func TestNewServerCompactMode(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		CompactMode:           true,
	}})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()

	// Should include enhanced tools
	hasProjectOverview := false
	for _, name := range enabled {
		if name == "get_project_overview" {
			hasProjectOverview = true
			break
		}
	}
	if !hasProjectOverview {
		t.Fatal("compact mode should include get_project_overview")
	}

	// Should NOT include superseded raw tools
	for _, name := range enabled {
		if name == "get_project" {
			t.Fatalf("compact mode should hide get_project, got %v", enabled)
		}
	}

	// Should NOT include superseded appstack tools
	for _, name := range enabled {
		if name == "get_application" {
			t.Fatalf("compact mode should hide get_application, got %v", enabled)
		}
	}

	// Should NOT include superseded flow tools
	for _, name := range enabled {
		if name == "get_pipeline" {
			t.Fatalf("compact mode should hide get_pipeline, got %v", enabled)
		}
	}

	// Should still include non-superseded raw tools
	hasPodLog := false
	for _, name := range enabled {
		if name == "get_pod_container_log" {
			hasPodLog = true
			break
		}
	}
	if !hasPodLog {
		t.Fatal("compact mode should include non-superseded tools like get_pod_container_log")
	}
}

func TestNewServerCompactModeAllowsExplicitEnable(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		CompactMode:           true,
		EnabledTools:          []string{"get_current_user"},
	}})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()
	if len(enabled) != 1 || enabled[0] != "get_current_user" {
		t.Fatalf("explicit enabled tools should override compact defaults, got %v", enabled)
	}
}

func TestNewServerEnableDomainsWhitelist(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		EnabledDomains:        []string{"platform"},
	}})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()
	for _, name := range enabled {
		if name == "search_projects" {
			t.Fatalf("enable-domains whitelist should exclude projex, got %v", enabled)
		}
	}
	if len(enabled) == 0 {
		t.Fatal("enable-domains whitelist should include platform tools")
	}
}

func TestNewServerDisableDomainsBlacklist(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		DisabledDomains:       []string{"codeup", "flow", "appstack", "lingma", "packages"},
	}})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()
	for _, name := range enabled {
		if name == "list_repositories" || name == "list_pipelines" {
			t.Fatalf("disable-domains blacklist should exclude codeup/flow, got %v", enabled)
		}
	}
	hasPlatform := false
	hasProjex := false
	for _, name := range enabled {
		if name == "get_current_user" {
			hasPlatform = true
		}
		if name == "get_project_overview" {
			hasProjex = true
		}
	}
	if !hasPlatform {
		t.Fatalf("disable-domains should keep platform, got %v", enabled)
	}
	if !hasProjex {
		t.Fatalf("disable-domains should keep projex, got %v", enabled)
	}
}

func TestNewServerEnableDomainsWithCompact(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		CompactMode:           true,
		EnabledDomains:        []string{"platform", "projex"},
	}})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()

	// Should include platform tools
	hasCurrentUser := false
	for _, name := range enabled {
		if name == "get_current_user" {
			hasCurrentUser = true
			break
		}
	}
	if !hasCurrentUser {
		t.Fatal("domain+compact should include platform tools")
	}

	// Should include projex tools
	hasProjectOverview := false
	for _, name := range enabled {
		if name == "get_project_overview" {
			hasProjectOverview = true
			break
		}
	}
	if !hasProjectOverview {
		t.Fatal("domain+compact should include projex enhanced tools")
	}

	// Should NOT include codeup tools
	for _, name := range enabled {
		if name == "list_repositories" {
			t.Fatalf("domain+compact should not include codeup tools, got %v", enabled)
		}
	}

	// Should NOT include superseded raw tools
	for _, name := range enabled {
		if name == "get_project" || name == "get_organization" {
			t.Fatalf("domain+compact should hide superseded raw tool %s, got %v", name, enabled)
		}
	}

	// Should still include non-superseded list tools
	hasSearchProjects := false
	for _, name := range enabled {
		if name == "search_projects" {
			hasSearchProjects = true
		}
	}
	if !hasSearchProjects {
		t.Fatal("domain+compact should include non-superseded tools like search_projects")
	}
}
