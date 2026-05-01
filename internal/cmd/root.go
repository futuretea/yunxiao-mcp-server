package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
	"github.com/futuretea/yunxiao-mcp-server/pkg/core/logging"
	"github.com/futuretea/yunxiao-mcp-server/pkg/core/version"
	internalhttp "github.com/futuretea/yunxiao-mcp-server/pkg/server/http"
	mcpserver "github.com/futuretea/yunxiao-mcp-server/pkg/server/mcp"
)

// IOStreams groups the streams used by the CLI.
type IOStreams struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
}

// NewMCPServer creates the root command for the Yunxiao MCP server.
func NewMCPServer(streams IOStreams) *cobra.Command {
	var cfgFile string
	v := viper.New()

	command := &cobra.Command{
		Use:           version.BinaryName,
		Short:         "Yunxiao MCP Server - Model Context Protocol server for Alibaba Cloud DevOps",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := bindFlags(v, cmd); err != nil {
				return err
			}
			return runServer(cmd.Context(), cfgFile, streams, v)
		},
	}

	command.SetIn(streams.In)
	command.SetOut(streams.Out)
	command.SetErr(streams.ErrOut)

	command.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (YAML)")
	command.Flags().Int("port", 0, "port for HTTP mode; 0 runs stdio mode")
	command.Flags().String("sse-base-url", "", "public base URL for SSE message endpoints")
	command.Flags().String("log-level", "info", "log level: trace, debug, info, warn, error, fatal, panic, disabled")
	command.Flags().String("base-url", config.DefaultBaseURL, "Yunxiao OpenAPI host or API base URL")
	command.Flags().String("access-token", "", "Yunxiao access token; also read from YUNXIAO_MCP_ACCESS_TOKEN or YUNXIAO_ACCESS_TOKEN")
	command.Flags().Bool("read-only", true, "run in read-only mode")
	command.Flags().StringSlice("enabled-tools", []string{}, "comma-separated list of tool names to enable")
	command.Flags().StringSlice("disabled-tools", []string{}, "comma-separated list of tool names to disable")
	command.Flags().Int("request-timeout-seconds", 30, "Yunxiao API request timeout in seconds")

	command.AddCommand(newVersionCommand(streams))
	return command
}

func bindFlags(v *viper.Viper, cmd *cobra.Command) error {
	bindings := map[string]string{
		"port":                    "port",
		"sse_base_url":            "sse-base-url",
		"log_level":               "log-level",
		"base_url":                "base-url",
		"access_token":            "access-token",
		"read_only":               "read-only",
		"enabled_tools":           "enabled-tools",
		"disabled_tools":          "disabled-tools",
		"request_timeout_seconds": "request-timeout-seconds",
	}

	for key, flag := range bindings {
		if err := v.BindPFlag(key, cmd.Flags().Lookup(flag)); err != nil {
			return fmt.Errorf("bind flag %s: %w", flag, err)
		}
	}
	return nil
}

func runServer(ctx context.Context, cfgFile string, streams IOStreams, v *viper.Viper) error {
	cfg, err := config.LoadConfig(cfgFile, v)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if cfg.Port == 0 {
		logging.Disable()
	} else if err := logging.Initialize(cfg.LogLevel, streams.ErrOut); err != nil {
		return fmt.Errorf("initialize logging: %w", err)
	}

	server, err := mcpserver.NewServer(mcpserver.Configuration{StaticConfig: cfg})
	if err != nil {
		return fmt.Errorf("create MCP server: %w", err)
	}
	defer server.Close()

	if cfg.Port != 0 {
		return internalhttp.Serve(ctx, server, cfg)
	}

	return server.ServeStdio()
}

func newVersionCommand(streams IOStreams) *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: "print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(streams.Out, version.GetVersionInfo())
		},
	}
	command.SetOut(streams.Out)
	command.SetErr(streams.ErrOut)
	return command
}
