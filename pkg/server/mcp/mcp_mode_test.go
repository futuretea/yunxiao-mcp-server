package mcp

import (
	"testing"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
)

func contains(names []string, target string) bool {
	for _, name := range names {
		if name == target {
			return true
		}
	}
	return false
}

func assertContains(t *testing.T, names []string, target string) {
	t.Helper()
	if !contains(names, target) {
		t.Fatalf("expected %v to contain %q", names, target)
	}
}

func assertNotContains(t *testing.T, names []string, target string) {
	t.Helper()
	if contains(names, target) {
		t.Fatalf("expected %v not to contain %q", names, target)
	}
}

func TestNewServerCompactMode(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: newTestConfig(func(c *config.StaticConfig) {
		c.CompactMode = true
	})})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()

	assertContains(t, enabled, "get_project_overview")
	assertNotContains(t, enabled, "get_project")
	assertNotContains(t, enabled, "get_application")
	assertNotContains(t, enabled, "get_pipeline")
	assertContains(t, enabled, "get_pod_container_log")
}

func TestNewServerCompactModeAllowsExplicitEnable(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: newTestConfig(func(c *config.StaticConfig) {
		c.CompactMode = true
		c.EnabledTools = []string{"get_current_user"}
	})})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()
	if len(enabled) != 1 || enabled[0] != "get_current_user" {
		t.Fatalf("explicit enabled tools should override compact defaults, got %v", enabled)
	}
}

func TestNewServerEnableDomainsWhitelist(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: newTestConfig(func(c *config.StaticConfig) {
		c.EnabledDomains = []string{"platform"}
	})})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()
	assertNotContains(t, enabled, "search_projects")
	assertContains(t, enabled, "get_current_user")
}

func TestNewServerDisableDomainsBlacklist(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: newTestConfig(func(c *config.StaticConfig) {
		c.DisabledDomains = []string{"codeup", "flow", "appstack", "lingma", "packages"}
	})})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()
	assertNotContains(t, enabled, "list_repositories")
	assertNotContains(t, enabled, "list_pipelines")
	assertContains(t, enabled, "get_current_user")
	assertContains(t, enabled, "get_project_overview")
}

func TestNewServerDefaultCompactHidesRawTools(t *testing.T) {
	// Default CompactMode=true means raw tools with enhanced alternatives are hidden
	s, err := NewServer(Configuration{StaticConfig: newTestConfig(func(c *config.StaticConfig) {
		c.CompactMode = true
	})})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()
	assertNotContains(t, enabled, "get_application")
	assertNotContains(t, enabled, "get_pipeline")
}

func TestNewServerNoCompactShowsAllTools(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: newTestConfig(func(c *config.StaticConfig) {
		c.CompactMode = false
	})})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()
	assertContains(t, enabled, "get_application")
	assertContains(t, enabled, "get_pipeline")
}

func TestNewServerEnableDomainsWithCompact(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: newTestConfig(func(c *config.StaticConfig) {
		c.CompactMode = true
		c.EnabledDomains = []string{"platform", "projex"}
	})})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()

	assertContains(t, enabled, "get_current_user")
	assertContains(t, enabled, "get_project_overview")
	assertNotContains(t, enabled, "list_repositories")
	assertNotContains(t, enabled, "get_project")
	assertNotContains(t, enabled, "get_organization")
	assertContains(t, enabled, "search_projects")
}
