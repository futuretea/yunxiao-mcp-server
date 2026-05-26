package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
	"github.com/futuretea/yunxiao-mcp-server/pkg/core/version"
	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
	yunxiaoTools "github.com/futuretea/yunxiao-mcp-server/pkg/toolset/yunxiao"
	yunxiaoSDK "github.com/futuretea/yunxiao-mcp-server/pkg/yunxiao"
)

const yunxiaoCLIBinaryName = "yunxiao"

// NewYunxiaoCLI creates the standalone Yunxiao CLI root command.
func NewYunxiaoCLI(streams IOStreams) *cobra.Command {
	var cfgFile string
	v := viper.New()

	command := &cobra.Command{
		Use:           yunxiaoCLIBinaryName,
		Short:         "Yunxiao CLI for Alibaba Cloud DevOps",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	command.SetIn(streams.In)
	command.SetOut(streams.Out)
	command.SetErr(streams.ErrOut)

	addYunxiaoCommonFlags(command, &cfgFile)
	command.AddCommand(newYunxiaoMCPCommand(streams, &cfgFile, v))
	command.AddCommand(newYunxiaoVersionCommand(streams))
	command.AddCommand(newYunxiaoOrganizationCommand(streams, &cfgFile, v))
	command.AddCommand(newYunxiaoMemberCommand(streams, &cfgFile, v))
	command.AddCommand(newYunxiaoDepartmentCommand(streams, &cfgFile, v))
	command.AddCommand(newYunxiaoPipelineCommand(streams, &cfgFile, v))
	command.AddCommand(newYunxiaoProjectCommand(streams, &cfgFile, v))
	command.AddCommand(newYunxiaoRepoCommand(streams, &cfgFile, v))
	command.AddCommand(newYunxiaoSprintCommand(streams, &cfgFile, v))
	command.AddCommand(newYunxiaoTaskCommand(streams, &cfgFile, v))
	command.AddCommand(newYunxiaoToolsCommand(streams, &cfgFile, v))
	command.AddCommand(newYunxiaoUserCommand(streams, &cfgFile, v))
	command.AddCommand(newYunxiaoCompletionCommand(streams))
	return command
}

func addYunxiaoCommonFlags(command *cobra.Command, cfgFile *string) {
	flags := command.PersistentFlags()
	flags.StringVar(cfgFile, "config", "", "config file path (YAML)")
	flags.String("base-url", config.DefaultBaseURL, "Yunxiao OpenAPI host or API base URL")
	flags.String("access-token", "", "Yunxiao access token; also read from YUNXIAO_MCP_ACCESS_TOKEN or YUNXIAO_ACCESS_TOKEN")
	flags.Bool("insecure-skip-tls-verify", false, "skip Yunxiao server TLS certificate verification")
	flags.Bool("read-only", true, "enable only read-only tools")
	flags.StringSlice("enabled-tools", []string{}, "comma-separated list of tool names to enable")
	flags.StringSlice("disabled-tools", []string{}, "comma-separated list of tool names to disable")
	flags.StringSlice("enable-domains", []string{}, "comma-separated list of tool domains to enable (e.g. platform,projex)")
	flags.StringSlice("disable-domains", []string{}, "comma-separated list of tool domains to disable (e.g. codeup,flow)")
	flags.Bool("compact", true, "hide raw API tools that have enhanced overview alternatives; set false to show all tools")
	flags.Int("request-timeout-seconds", 30, "Yunxiao API request timeout in seconds")
}

func bindYunxiaoCLIFlags(v *viper.Viper, cmd *cobra.Command) error {
	flags := cmd.Root().PersistentFlags()
	bindings := map[string]string{
		"base_url":                 "base-url",
		"access_token":             "access-token",
		"insecure_skip_tls_verify": "insecure-skip-tls-verify",
		"read_only":                "read-only",
		"enabled_tools":            "enabled-tools",
		"disabled_tools":           "disabled-tools",
		"enabled_domains":          "enable-domains",
		"disabled_domains":         "disable-domains",
		"compact":                  "compact",
		"request_timeout_seconds":  "request-timeout-seconds",
	}

	for key, flag := range bindings {
		if err := v.BindPFlag(key, flags.Lookup(flag)); err != nil {
			return fmt.Errorf("bind flag %s: %w", flag, err)
		}
	}
	return nil
}

func loadYunxiaoCLIConfig(cmd *cobra.Command, cfgFile string, v *viper.Viper) (*config.StaticConfig, error) {
	if err := bindYunxiaoCLIFlags(v, cmd); err != nil {
		return nil, err
	}
	cfg, err := config.LoadConfig(cfgFile, v)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return cfg, nil
}

func newYunxiaoVersionCommand(streams IOStreams) *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: "print version information",
		Run: func(cmd *cobra.Command, args []string) {
			_, _ = fmt.Fprintln(streams.Out, version.GetVersionInfoFor(yunxiaoCLIBinaryName))
		},
	}
	command.SetOut(streams.Out)
	command.SetErr(streams.ErrOut)
	return command
}

func newYunxiaoToolsCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "tools",
		Aliases: []string{"tool"},
		Short:   "list and call Yunxiao tools",
	}
	command.AddCommand(newYunxiaoToolsListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoToolsDescribeCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoToolsCallCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoTaskCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "task",
		Aliases: []string{"tasks"},
		Short:   "work with Projex tasks",
	}
	command.AddCommand(newYunxiaoTaskListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoTaskViewCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoTaskListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options taskListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Projex tasks in a project",
		Example: `  # List tasks in a project
  yunxiao task list --project-id 123

  # Filter by status and assignee
  yunxiao task list --project-id 123 --status "处理中" --assigned-to user1

  # Output as JSON
  yunxiao task list --project-id 123 --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "search_workitems", params)
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printTaskList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.ProjectID, "project-id", "", "Projex project ID")
	flags.StringVar(&options.Category, "category", "Task", "work item category, e.g. Task, Bug, Requirement")
	flags.StringVar(&options.Subject, "subject", "", "subject/title keyword")
	flags.StringVar(&options.Status, "status", "", "comma-separated status IDs")
	flags.StringVar(&options.AssignedTo, "assigned-to", "", "assignee user ID")
	flags.StringVar(&options.Creator, "creator", "", "creator user ID")
	flags.StringVar(&options.Sprint, "sprint", "", "sprint ID")
	flags.StringVar(&options.OrderBy, "order-by", "", "sort field")
	flags.StringVar(&options.Sort, "sort", "", "sort direction, e.g. asc or desc")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func newYunxiaoToolsListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var jsonOutput bool
	command := &cobra.Command{
		Use:   "list",
		Short: "list enabled Yunxiao tools",
		Example: `  # List all enabled tools
  yunxiao tools list

  # Output as JSON
  yunxiao tools list --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}

			tools, err := yunxiaoTools.BuildToolCatalog(nil, toolCatalogOptionsFromConfig(cfg))
			if err != nil {
				return err
			}
			if jsonOutput {
				return printToolSummariesJSON(streams.Out, tools)
			}
			return printToolSummariesTable(streams.Out, tools)
		},
	}
	command.Flags().BoolVar(&jsonOutput, "json", false, "print tools as JSON")
	return command
}

func newYunxiaoToolsCallCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var rawParams string
	var paramsFile string
	command := &cobra.Command{
		Use:   "call <tool-name>",
		Short: "call an enabled Yunxiao tool with JSON parameters",
		Example: `  # Call a tool with inline JSON
  yunxiao tools call search_projects --params '{"name":"demo"}'

  # Call a tool with params from stdin
  echo '{"name":"demo"}' | yunxiao tools call search_projects --params-file -

  # Call a tool with params from a file
  yunxiao tools call search_projects --params-file query.json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := parseToolParamsWithInput(rawParams, paramsFile, cmd.InOrStdin())
			if err != nil {
				return err
			}

			result, err := callYunxiaoTool(cmd, cfg, args[0], params)
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	command.Flags().StringVar(&rawParams, "params", "{}", "tool parameters as a JSON object")
	command.Flags().StringVar(&paramsFile, "params-file", "", "path to a JSON file containing tool parameters; use - to read stdin")
	return command
}

func callYunxiaoTool(cmd *cobra.Command, cfg *config.StaticConfig, toolName string, params map[string]any) (string, error) {
	tools, err := yunxiaoTools.BuildToolCatalog(nil, toolCatalogOptionsFromConfig(cfg))
	if err != nil {
		return "", err
	}
	tool, ok := yunxiaoTools.FindTool(tools, toolName)
	if !ok {
		return "", fmt.Errorf("unknown Yunxiao tool %q; run 'yunxiao tools list' to see available tools", toolName)
	}
	if err := yunxiaoTools.ValidateToolRequiredParams(tool, params, "organizationId"); err != nil {
		return "", yunxiaoSDK.WrapError(err)
	}

	client, err := newSDKClientFromConfig(cfg)
	if err != nil {
		return "", err
	}
	if toolNeedsDefaultOrganizationID(tool, params) {
		if err := client.ResolveDefaultOrgID(cmd.Context()); err != nil {
			return "", fmt.Errorf("resolve default organization: %w", yunxiaoSDK.WrapError(yunxiaoSDK.FriendlyAPIError(err)))
		}
	}

	return yunxiaoTools.InvokeTool(cmd.Context(), client, tool, params)
}

func toolCatalogOptionsFromConfig(cfg *config.StaticConfig) yunxiaoTools.ToolCatalogOptions {
	return yunxiaoTools.ToolCatalogOptions{
		ReadOnly:        cfg.ReadOnly,
		CompactMode:     cfg.CompactMode,
		EnabledTools:    cfg.EnabledTools,
		DisabledTools:   cfg.DisabledTools,
		EnabledDomains:  cfg.EnabledDomains,
		DisabledDomains: cfg.DisabledDomains,
	}
}

func newSDKClientFromConfig(cfg *config.StaticConfig) (*yunxiaoSDK.Client, error) {
	clientOptions := []yunxiaoSDK.ClientOption{}
	if cfg.InsecureSkipTLSVerify {
		clientOptions = append(clientOptions, yunxiaoSDK.WithInsecureSkipTLSVerify(true))
	}
	return yunxiaoSDK.NewClient(
		cfg.BaseURL,
		cfg.AccessToken,
		time.Duration(cfg.RequestTimeoutSeconds)*time.Second,
		clientOptions...,
	)
}

type toolSummary struct {
	Name        string `json:"name"`
	Domain      string `json:"domain"`
	Access      string `json:"access"`
	Description string `json:"description,omitempty"`
}

func newToolSummary(tool toolset.ServerTool) toolSummary {
	access := "read-only"
	if yunxiaoTools.IsWriteTool(tool.Tool.Name) {
		access = "write"
	}
	return toolSummary{
		Name:        tool.Tool.Name,
		Domain:      tool.Domain,
		Access:      access,
		Description: tool.Tool.Description,
	}
}

func toolNeedsDefaultOrganizationID(tool toolset.ServerTool, params map[string]any) bool {
	if _, ok := tool.Tool.InputSchema.Properties["organizationId"]; !ok {
		return false
	}
	orgID, _ := params["organizationId"].(string)
	return strings.TrimSpace(orgID) == ""
}

func printToolSummariesJSON(out anyWriter, tools []toolset.ServerTool) error {
	summaries := make([]toolSummary, 0, len(tools))
	for _, tool := range tools {
		summaries = append(summaries, newToolSummary(tool))
	}
	encoded, err := json.MarshalIndent(summaries, "", "  ")
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(out, string(encoded))
	return nil
}

func printToolSummariesTable(out anyWriter, tools []toolset.ServerTool) error {
	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("NAME\tDOMAIN\tACCESS\tDESCRIPTION"))
	for _, tool := range tools {
		summary := newToolSummary(tool)
		description := strings.ReplaceAll(summary.Description, "\n", " ")
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", summary.Name, summary.Domain, summary.Access, description)
	}
	return writer.Flush()
}


func parseToolParamsWithInput(rawParams, paramsFile string, in io.Reader) (map[string]any, error) {
	if strings.TrimSpace(paramsFile) != "" {
		if strings.TrimSpace(paramsFile) == "-" {
			if in == nil {
				return nil, fmt.Errorf("read params stdin: no input")
			}
			content, err := io.ReadAll(in)
			if err != nil {
				return nil, fmt.Errorf("read params stdin: %w", err)
			}
			rawParams = string(content)
		} else {
			content, err := os.ReadFile(paramsFile)
			if err != nil {
				return nil, fmt.Errorf("read params file: %w", err)
			}
			rawParams = string(content)
		}
	}
	if strings.TrimSpace(rawParams) == "" {
		return map[string]any{}, nil
	}

	var params map[string]any
	decoder := json.NewDecoder(strings.NewReader(rawParams))
	if err := decoder.Decode(&params); err != nil {
		return nil, fmt.Errorf("invalid params JSON: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("params JSON must contain a single object")
		}
		return nil, fmt.Errorf("invalid params JSON: %w", err)
	}
	if params == nil {
		return nil, fmt.Errorf("params JSON must be an object")
	}
	return params, nil
}

type anyWriter interface {
	Write([]byte) (int, error)
}

func newYunxiaoCompletionCommand(streams IOStreams) *cobra.Command {
	command := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "generate shell completion script",
		Long: `To load completions:

  Bash:
    source <(yunxiao completion bash)

  Zsh:
    source <(yunxiao completion zsh)

  Fish:
    yunxiao completion fish | source

  PowerShell:
    yunxiao completion powershell | Out-String | Invoke-Expression`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(streams.Out)
			case "zsh":
				cmd.Root().GenZshCompletion(streams.Out)
			case "fish":
				cmd.Root().GenFishCompletion(streams.Out, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(streams.Out)
			}
		},
	}
	command.SetOut(streams.Out)
	command.SetErr(streams.ErrOut)
	return command
}
