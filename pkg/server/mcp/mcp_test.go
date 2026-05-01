package mcp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
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
