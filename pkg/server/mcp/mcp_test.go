package mcp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

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

func TestNewServerRejectsNilStaticConfig(t *testing.T) {
	_, err := NewServer(Configuration{})
	if err == nil {
		t.Fatal("NewServer() expected error for nil StaticConfig")
	}
}

func TestNewServerRejectsInvalidBaseURL(t *testing.T) {
	_, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               "://invalid-url",
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
	}})
	if err == nil {
		t.Fatal("NewServer() expected error for invalid base URL")
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

func TestFilterToolsByDomainsNoFilterReturnsAll(t *testing.T) {
	tools := []toolset.ServerTool{
		{Tool: mcp.NewTool("a"), Domain: "platform"},
		{Tool: mcp.NewTool("b"), Domain: "projex"},
	}
	got := filterToolsByDomains(tools, nil, nil)
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
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
