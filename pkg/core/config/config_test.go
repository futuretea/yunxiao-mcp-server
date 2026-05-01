package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestValidateValidConfig(t *testing.T) {
	cfg := &StaticConfig{
		Port:                  0,
		LogLevel:              "info",
		BaseURL:               DefaultBaseURL,
		ReadOnly:              true,
		RequestTimeoutSeconds: 30,
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestValidateRejectsInvalidPort(t *testing.T) {
	cfg := &StaticConfig{
		Port:                  70000,
		LogLevel:              "info",
		BaseURL:               DefaultBaseURL,
		RequestTimeoutSeconds: 30,
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() expected invalid port error")
	}
}

func TestValidateRejectsInvalidBaseURL(t *testing.T) {
	cfg := &StaticConfig{
		LogLevel:              "info",
		BaseURL:               "openapi-rdc.aliyuncs.com",
		RequestTimeoutSeconds: 30,
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() expected invalid base_url error")
	}
}

func TestGetPortString(t *testing.T) {
	tests := []struct {
		name string
		port int
		want string
	}{
		{name: "stdio", port: 0, want: ""},
		{name: "http", port: 8080, want: ":8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &StaticConfig{Port: tt.port}
			if got := cfg.GetPortString(); got != tt.want {
				t.Fatalf("GetPortString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLoadConfigReadsPrefixedEnvironment(t *testing.T) {
	t.Setenv("YUNXIAO_MCP_ACCESS_TOKEN", "prefixed-token")
	t.Setenv("YUNXIAO_MCP_BASE_URL", "https://example.com/yunxiao")

	cfg, err := LoadConfig("", viper.New())
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.AccessToken != "prefixed-token" {
		t.Fatalf("AccessToken = %q, want prefixed-token", cfg.AccessToken)
	}
	if cfg.BaseURL != "https://example.com/yunxiao" {
		t.Fatalf("BaseURL = %q", cfg.BaseURL)
	}
}

func TestLoadConfigReadsLegacyEnvironment(t *testing.T) {
	t.Setenv("YUNXIAO_ACCESS_TOKEN", "legacy-token")
	t.Setenv("YUNXIAO_API_BASE_URL", "https://legacy.example.com")

	cfg, err := LoadConfig("", viper.New())
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.AccessToken != "legacy-token" {
		t.Fatalf("AccessToken = %q, want legacy-token", cfg.AccessToken)
	}
	if cfg.BaseURL != "https://legacy.example.com" {
		t.Fatalf("BaseURL = %q", cfg.BaseURL)
	}
}

func TestLoadConfigPrefixedEnvironmentOverridesLegacyEnvironment(t *testing.T) {
	t.Setenv("YUNXIAO_MCP_ACCESS_TOKEN", "prefixed-token")
	t.Setenv("YUNXIAO_ACCESS_TOKEN", "legacy-token")
	t.Setenv("YUNXIAO_MCP_BASE_URL", "https://prefixed.example.com")
	t.Setenv("YUNXIAO_API_BASE_URL", "https://legacy.example.com")

	cfg, err := LoadConfig("", viper.New())
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.AccessToken != "prefixed-token" {
		t.Fatalf("AccessToken = %q, want prefixed-token", cfg.AccessToken)
	}
	if cfg.BaseURL != "https://prefixed.example.com" {
		t.Fatalf("BaseURL = %q", cfg.BaseURL)
	}
}

func TestLoadConfigLegacyEnvironmentOverridesConfigFile(t *testing.T) {
	t.Setenv("YUNXIAO_ACCESS_TOKEN", "legacy-token")
	t.Setenv("YUNXIAO_API_BASE_URL", "https://legacy.example.com")

	configPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(configPath, []byte(`
access_token: file-token
base_url: https://file.example.com
`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadConfig(configPath, viper.New())
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.AccessToken != "legacy-token" {
		t.Fatalf("AccessToken = %q, want legacy-token", cfg.AccessToken)
	}
	if cfg.BaseURL != "https://legacy.example.com" {
		t.Fatalf("BaseURL = %q", cfg.BaseURL)
	}
}

func TestLoadConfigExplicitSetOverridesEnvironment(t *testing.T) {
	t.Setenv("YUNXIAO_MCP_ACCESS_TOKEN", "env-token")

	v := viper.New()
	v.Set("access_token", "explicit-token")

	cfg, err := LoadConfig("", v)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.AccessToken != "explicit-token" {
		t.Fatalf("AccessToken = %q, want explicit-token", cfg.AccessToken)
	}
}

func TestLoadConfigNormalizesToolFilters(t *testing.T) {
	v := viper.New()
	v.Set("enabled_tools", []string{" get_current_user ", "", "list_organizations"})

	cfg, err := LoadConfig("", v)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if len(cfg.EnabledTools) != 2 {
		t.Fatalf("EnabledTools = %#v", cfg.EnabledTools)
	}
	if cfg.EnabledTools[0] != "get_current_user" || cfg.EnabledTools[1] != "list_organizations" {
		t.Fatalf("EnabledTools = %#v", cfg.EnabledTools)
	}
}
