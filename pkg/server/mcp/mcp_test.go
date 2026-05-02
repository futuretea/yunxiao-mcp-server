package mcp

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
	yunxiaoToolset "github.com/futuretea/yunxiao-mcp-server/pkg/toolset/yunxiao"
)

func newTestServer(enabledTools, disabledTools []string) *Server {
	return &Server{
		configuration: &Configuration{
			StaticConfig: &config.StaticConfig{
				EnabledTools:  enabledTools,
				DisabledTools: disabledTools,
			},
		},
	}
}

func TestShouldEnableToolAllEnabledByDefault(t *testing.T) {
	s := newTestServer(nil, nil)
	if !s.shouldEnableTool("get_current_user") {
		t.Fatal("tool should be enabled by default")
	}
}

func TestShouldEnableToolUsesAllowList(t *testing.T) {
	s := newTestServer([]string{"get_current_user"}, nil)
	if !s.shouldEnableTool("get_current_user") {
		t.Fatal("get_current_user should be enabled")
	}
	if s.shouldEnableTool("list_organizations") {
		t.Fatal("list_organizations should not be enabled")
	}
}

func TestShouldEnableToolDisabledTakesPriority(t *testing.T) {
	s := newTestServer([]string{"get_current_user"}, []string{"get_current_user"})
	if s.shouldEnableTool("get_current_user") {
		t.Fatal("disabled tool should not be enabled")
	}
}

func TestNewServerRegistersTools(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		EnabledTools:          []string{"get_current_user"},
	}})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()
	if len(enabled) != 1 || enabled[0] != "get_current_user" {
		t.Fatalf("enabled tools = %#v", enabled)
	}
}

func TestNewServerRejectsUnknownEnabledTool(t *testing.T) {
	_, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		EnabledTools:          []string{"get_user_organizations_typo"},
	}})
	if err == nil {
		t.Fatal("NewServer() expected unknown enabled tool error")
	}
}

func TestNewServerRejectsUnknownDisabledTool(t *testing.T) {
	_, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		DisabledTools:         []string{"get_user_organizations_typo"},
	}})
	if err == nil {
		t.Fatal("NewServer() expected unknown disabled tool error")
	}
}

func TestNewServerRejectsZeroEnabledTools(t *testing.T) {
	_, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		EnabledTools:          []string{"get_current_user"},
		DisabledTools:         []string{"get_current_user"},
	}})
	if err == nil {
		t.Fatal("NewServer() expected zero enabled tools error")
	}
}

func TestRequestAccessTokenPrefersHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/message?yunxiao_access_token=query-token", nil)
	req.Header.Set(yunxiaoToolset.AccessTokenHeader, "header-token")

	ctx := withRequestAccessToken(t.Context(), req)

	if got := yunxiaoToolset.AccessTokenFromContext(ctx); got != "header-token" {
		t.Fatalf("access token = %q", got)
	}
}

func TestRequestAccessTokenUsesQueryParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/message?yunxiao_access_token=query-token", nil)

	ctx := withRequestAccessToken(t.Context(), req)

	if got := yunxiaoToolset.AccessTokenFromContext(ctx); got != "query-token" {
		t.Fatalf("access token = %q", got)
	}
}

func TestNewServerProjectFocusedMode(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		ProjectFocused:        true,
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
		t.Fatalf("project-focused mode should include platform tools, got %v", enabled)
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
		t.Fatalf("project-focused mode should include projex enhanced tools, got %v", enabled)
	}

	// Should NOT include codeup tools
	for _, name := range enabled {
		if name == "list_repositories" {
			t.Fatalf("project-focused mode should not include codeup tools, got %v", enabled)
		}
	}

	// Should NOT include superseded raw projex tools
	for _, name := range enabled {
		if name == "get_project" {
			t.Fatalf("project-focused mode should hide superseded raw tool get_project, got %v", enabled)
		}
	}
}

func TestNewServerProjectFocusedModeAllowsExplicitEnable(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		ProjectFocused:        true,
		EnabledTools:          []string{"get_current_user"},
	}})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()
	if len(enabled) != 1 || enabled[0] != "get_current_user" {
		t.Fatalf("explicit enabled tools should override project-focused defaults, got %v", enabled)
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
	// Should still include platform and projex
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

func TestNewServerEnableDomainsOverridesProjectFocused(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		ProjectFocused:        true,
		EnabledDomains:        []string{"platform", "projex", "codeup"},
	}})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()

	// Should include codeup because enabled_domains overrides project_focused
	hasCodeup := false
	for _, name := range enabled {
		if name == "list_repositories" {
			hasCodeup = true
		}
	}
	if !hasCodeup {
		t.Fatalf("enabled_domains should override project_focused, expected codeup tools, got %v", enabled)
	}

	// Should include superseded raw tools because enabled_domains does not hide them
	hasGetProject := false
	for _, name := range enabled {
		if name == "get_project" {
			hasGetProject = true
		}
	}
	if !hasGetProject {
		t.Fatalf("enabled_domains should include all projex tools including get_project, got %v", enabled)
	}
}

func TestNewServerMinimalMode(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		MinimalMode:           true,
	}})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()

	// Should include core platform tools
	hasCurrentUser := false
	for _, name := range enabled {
		if name == "get_current_user" {
			hasCurrentUser = true
			break
		}
	}
	if !hasCurrentUser {
		t.Fatalf("minimal mode should include get_current_user, got %v", enabled)
	}

	// Should include core projex tools
	hasProjectOverview := false
	for _, name := range enabled {
		if name == "get_project_overview" {
			hasProjectOverview = true
			break
		}
	}
	if !hasProjectOverview {
		t.Fatalf("minimal mode should include get_project_overview, got %v", enabled)
	}

	// Should NOT include non-core tools like codeup
	for _, name := range enabled {
		if name == "list_repositories" {
			t.Fatalf("minimal mode should not include codeup tools, got %v", enabled)
		}
	}

	// Should NOT include platform admin tools
	for _, name := range enabled {
		if name == "list_organization_departments" {
			t.Fatalf("minimal mode should not include platform admin tools, got %v", enabled)
		}
	}

	// Should NOT include projex metadata tools
	for _, name := range enabled {
		if name == "get_work_item_type" {
			t.Fatalf("minimal mode should not include metadata tools, got %v", enabled)
		}
	}
}

func TestNewServerMinimalModeOverridesProjectFocused(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		MinimalMode:           true,
		ProjectFocused:        true,
	}})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	enabled := s.GetEnabledTools()

	// Minimal should take priority over project-focused
	for _, name := range enabled {
		if name == "list_organization_departments" {
			t.Fatalf("minimal mode should override project_focused, got %v", enabled)
		}
	}
}

func TestIsHealthy(t *testing.T) {
	if newTestServer(nil, nil).IsHealthy() {
		t.Fatal("test server without client should not be healthy")
	}

	s, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		EnabledTools:          []string{"get_current_user"},
	}})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}
	if !s.IsHealthy() {
		t.Fatal("server with client and tools should be healthy")
	}
}

func TestNewTextResultReturnsContent(t *testing.T) {
	result := NewTextResult("hello", nil)
	if result.IsError {
		t.Fatal("expected IsError = false")
	}
	if len(result.Content) != 1 {
		t.Fatalf("content count = %d, want 1", len(result.Content))
	}
	text, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("content type = %T, want TextContent", result.Content[0])
	}
	if text.Text != "hello" {
		t.Fatalf("text = %q, want hello", text.Text)
	}
}

func TestNewTextResultReturnsError(t *testing.T) {
	result := NewTextResult("", fmt.Errorf("something went wrong"))
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
	if text.Text != "something went wrong" {
		t.Fatalf("text = %q, want error message", text.Text)
	}
}

func TestRequestAccessTokenHandlesNilRequest(t *testing.T) {
	if got := requestAccessToken(nil); got != "" {
		t.Fatalf("requestAccessToken(nil) = %q, want empty", got)
	}
}

func TestValidateToolFiltersDetectsDuplicate(t *testing.T) {
	tools := []toolset.ServerTool{
		{Tool: mcp.NewTool("dup_tool")},
		{Tool: mcp.NewTool("dup_tool")},
	}
	if err := validateToolFilters(tools, nil, nil); err == nil {
		t.Fatal("validateToolFilters expected duplicate error")
	}
}

func TestRegisterToolFillsDefaultOrganizationID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"default-org"}]`))
	}))
	defer ts.Close()

	client, _ := yunxiaoToolset.NewClient(ts.URL, "token", time.Second)
	_ = client.ResolveDefaultOrgID(context.Background())

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

	_, err := registered.Handler(context.Background(), mcp.CallToolRequest{
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

func TestRegisterToolPreservesExistingOrganizationID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"default-org"}]`))
	}))
	defer ts.Close()

	client, _ := yunxiaoToolset.NewClient(ts.URL, "token", time.Second)
	_ = client.ResolveDefaultOrgID(context.Background())

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

	_, err := registered.Handler(context.Background(), mcp.CallToolRequest{
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

func TestServeSSEReturnsWorkingSSEServer(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		AccessToken:           "token",
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
	}})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	sseServer := s.ServeSSE("", &http.Server{})
	if sseServer == nil {
		t.Fatal("ServeSSE() returned nil")
	}

	ts := httptest.NewServer(sseServer.SSEHandler())
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Fatalf("GET /sse: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "text/event-stream") {
		t.Fatalf("Content-Type = %q, want text/event-stream", ct)
	}
}

func TestServeStreamableHTTPReturnsWorkingHandler(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
		AccessToken:           "token",
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
	}})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	streamable := s.ServeStreamableHTTP(&http.Server{})
	if streamable == nil {
		t.Fatal("ServeStreamableHTTP() returned nil")
	}

	ts := httptest.NewServer(streamable)
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Fatalf("GET /mcp: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
}

func TestCloseIsNoOp(t *testing.T) {
	s := newTestServer(nil, nil)
	s.Close() // should not panic
}
