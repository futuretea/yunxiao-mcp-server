package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const (
	// DefaultBaseURL is the public Yunxiao OpenAPI host used by the Node reference server.
	DefaultBaseURL = "https://openapi-rdc.aliyuncs.com"
)

// StaticConfig contains configuration that is fixed when the server starts.
type StaticConfig struct {
	Port                  int      `mapstructure:"port"`
	SSEBaseURL            string   `mapstructure:"sse_base_url"`
	LogLevel              string   `mapstructure:"log_level"`
	BaseURL               string   `mapstructure:"base_url"`
	AccessToken           string   `mapstructure:"access_token"`
	ReadOnly              bool     `mapstructure:"read_only"`
	EnabledTools          []string `mapstructure:"enabled_tools"`
	DisabledTools         []string `mapstructure:"disabled_tools"`
	EnabledDomains        []string `mapstructure:"enabled_domains"`
	DisabledDomains       []string `mapstructure:"disabled_domains"`
	ProjectFocused        bool     `mapstructure:"project_focused"`
	MinimalMode           bool     `mapstructure:"minimal"`
	RequestTimeoutSeconds int      `mapstructure:"request_timeout_seconds"`
}

// Validate checks whether the configuration can be used to start the server.
func (c *StaticConfig) Validate() error {
	if c.Port < 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 0 and 65535, got %d", c.Port)
	}
	if _, err := zerolog.ParseLevel(c.LogLevel); err != nil {
		return fmt.Errorf("invalid log_level %q: %w", c.LogLevel, err)
	}
	if c.RequestTimeoutSeconds <= 0 {
		return fmt.Errorf("request_timeout_seconds must be positive, got %d", c.RequestTimeoutSeconds)
	}
	parsed, err := url.Parse(c.BaseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("base_url must be an absolute URL, got %q", c.BaseURL)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("base_url scheme must be http or https, got %q", parsed.Scheme)
	}
	return nil
}

// GetPortString returns the configured HTTP listen address.
func (c *StaticConfig) GetPortString() string {
	if c.Port == 0 {
		return ""
	}
	return fmt.Sprintf(":%d", c.Port)
}

// LoadConfig loads configuration from defaults, optional YAML, environment, and flags.
func LoadConfig(configPath string, v *viper.Viper) (*StaticConfig, error) {
	if v == nil {
		v = viper.New()
	}

	v.SetDefault("port", 0)
	v.SetDefault("sse_base_url", "")
	v.SetDefault("log_level", "info")
	v.SetDefault("base_url", DefaultBaseURL)
	v.SetDefault("read_only", true)
	v.SetDefault("enabled_tools", []string{})
	v.SetDefault("disabled_tools", []string{})
	v.SetDefault("enabled_domains", []string{})
	v.SetDefault("disabled_domains", []string{})
	v.SetDefault("project_focused", false)
	v.SetDefault("minimal", false)
	v.SetDefault("request_timeout_seconds", 30)

	if configPath != "" {
		v.SetConfigFile(configPath)
		v.SetConfigType("yaml")
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config file: %w", err)
		}
	}

	v.SetEnvPrefix("YUNXIAO_MCP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
	if err := bindEnvironment(v); err != nil {
		return nil, err
	}

	cfg := &StaticConfig{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	cfg.BaseURL = strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	cfg.LogLevel = strings.ToLower(strings.TrimSpace(cfg.LogLevel))
	cfg.EnabledTools = normalizeStringSlice(cfg.EnabledTools)
	cfg.DisabledTools = normalizeStringSlice(cfg.DisabledTools)
	cfg.EnabledDomains = normalizeStringSlice(cfg.EnabledDomains)
	cfg.DisabledDomains = normalizeStringSlice(cfg.DisabledDomains)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func bindEnvironment(v *viper.Viper) error {
	envBindings := map[string][]string{
		"access_token": {"YUNXIAO_MCP_ACCESS_TOKEN", "YUNXIAO_ACCESS_TOKEN"},
		"base_url":     {"YUNXIAO_MCP_BASE_URL", "YUNXIAO_API_BASE_URL"},
	}

	for key, names := range envBindings {
		args := append([]string{key}, names...)
		if err := v.BindEnv(args...); err != nil {
			return fmt.Errorf("bind environment %s: %w", key, err)
		}
	}
	return nil
}

func normalizeStringSlice(values []string) []string {
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			normalized = append(normalized, value)
		}
	}
	return normalized
}
