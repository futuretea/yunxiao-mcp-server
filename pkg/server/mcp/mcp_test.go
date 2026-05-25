package mcp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
	yunxiaoToolset "github.com/futuretea/yunxiao-mcp-server/pkg/toolset/yunxiao"
	yunxiaoSDK "github.com/futuretea/yunxiao-mcp-server/pkg/yunxiao"
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
	tools, err := yunxiaoToolset.BuildToolCatalog(nil, yunxiaoToolset.ToolCatalogOptions{ReadOnly: true})
	if err != nil {
		t.Fatalf("BuildToolCatalog() error = %v", err)
	}
	if _, ok := yunxiaoToolset.FindTool(tools, "get_current_user"); !ok {
		t.Fatal("tool should be enabled by default")
	}
}

func TestShouldEnableToolUsesAllowList(t *testing.T) {
	tools, err := yunxiaoToolset.BuildToolCatalog(nil, yunxiaoToolset.ToolCatalogOptions{
		ReadOnly:     true,
		EnabledTools: []string{"get_current_user"},
	})
	if err != nil {
		t.Fatalf("BuildToolCatalog() error = %v", err)
	}
	if _, ok := yunxiaoToolset.FindTool(tools, "get_current_user"); !ok {
		t.Fatal("get_current_user should be enabled")
	}
	if _, ok := yunxiaoToolset.FindTool(tools, "list_organizations"); ok {
		t.Fatal("list_organizations should not be enabled")
	}
}

func TestShouldEnableToolDisabledTakesPriority(t *testing.T) {
	tools, err := yunxiaoToolset.BuildToolCatalog(nil, yunxiaoToolset.ToolCatalogOptions{
		ReadOnly:      true,
		EnabledTools:  []string{"get_current_user", "list_organizations"},
		DisabledTools: []string{"get_current_user"},
	})
	if err != nil {
		t.Fatalf("BuildToolCatalog() error = %v", err)
	}
	if _, ok := yunxiaoToolset.FindTool(tools, "get_current_user"); ok {
		t.Fatal("disabled tool should not be present")
	}
	if _, ok := yunxiaoToolset.FindTool(tools, "list_organizations"); !ok {
		t.Fatal("non-disabled allow-listed tool should be present")
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

func TestNewServerRejectsUnknownEnabledToolBeforeDefaultOrgRequest(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		http.Error(w, "unexpected request", http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := NewServer(Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               server.URL,
		AccessToken:           "token-1",
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
		EnabledTools:          []string{"get_user_organizations_typo"},
	}})
	if err == nil {
		t.Fatal("NewServer() expected unknown enabled tool error")
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0", requests)
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
	req.Header.Set(yunxiaoSDK.AccessTokenHeader, "header-token")

	ctx := withRequestAccessToken(t.Context(), req)

	if got := yunxiaoSDK.AccessTokenFromContext(ctx); got != "header-token" {
		t.Fatalf("access token = %q", got)
	}
}

func TestRequestAccessTokenUsesQueryParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/message?yunxiao_access_token=query-token", nil)

	ctx := withRequestAccessToken(t.Context(), req)

	if got := yunxiaoSDK.AccessTokenFromContext(ctx); got != "query-token" {
		t.Fatalf("access token = %q", got)
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
	s.Close()
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
