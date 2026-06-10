package yunxiao

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

// ToolCatalogOptions controls which Yunxiao tools are exposed to a caller.
type ToolCatalogOptions struct {
	ReadOnly        bool
	CompactMode     bool
	EnabledTools    []string
	DisabledTools   []string
	EnabledDomains  []string
	DisabledDomains []string
}

// BuildToolCatalog returns the filtered Yunxiao tool catalog used by MCP and CLI callers.
func BuildToolCatalog(client any, options ToolCatalogOptions) ([]toolset.ServerTool, error) {
	toolsetBuilder := &Toolset{ReadOnly: options.ReadOnly}
	tools := toolsetBuilder.GetTools(client)

	if len(options.EnabledDomains) > 0 {
		tools = filterToolsByDomains(tools, options.EnabledDomains, nil)
	} else if len(options.DisabledDomains) > 0 {
		tools = filterToolsByDomains(tools, nil, options.DisabledDomains)
	}

	if options.CompactMode {
		tools = toolsetBuilder.GetCompactTools(tools)
	}

	if err := validateToolFilters(tools, options.EnabledTools, options.DisabledTools); err != nil {
		return nil, err
	}

	filtered := make([]toolset.ServerTool, 0, len(tools))
	for _, tool := range tools {
		if shouldEnableTool(tool.Tool.Name, options.EnabledTools, options.DisabledTools) {
			filtered = append(filtered, tool)
		}
	}
	if len(filtered) == 0 {
		return nil, fmt.Errorf("no Yunxiao tools enabled; check enabled_tools, disabled_tools, enable_domains, disable_domains, compact")
	}
	return filtered, nil
}

// FindTool returns the named tool from a filtered catalog.
func FindTool(tools []toolset.ServerTool, name string) (toolset.ServerTool, bool) {
	for _, tool := range tools {
		if tool.Tool.Name == name {
			return tool, true
		}
	}
	return toolset.ServerTool{}, false
}

// IsWriteTool reports whether a tool performs a mutation and requires read_only=false.
func IsWriteTool(name string) bool {
	_, ok := writeToolNames[name]
	return ok
}

// ValidateToolRequiredParams checks a tool schema's required parameters before a
// caller performs side effects such as resolving a default organization.
func ValidateToolRequiredParams(tool toolset.ServerTool, params map[string]any, defaultedKeys ...string) error {
	if params == nil {
		params = map[string]any{}
	}
	defaulted := make(map[string]struct{}, len(defaultedKeys))
	for _, key := range defaultedKeys {
		defaulted[key] = struct{}{}
	}
	for _, key := range toolRequiredParams(tool) {
		if _, ok := defaulted[key]; ok {
			continue
		}
		if isMissingToolParam(params[key]) {
			return &ValidationError{Msg: fmt.Sprintf("%s is required", key)}
		}
	}
	return nil
}

// InvokeTool applies shared Yunxiao invocation behavior before calling a tool handler.
func InvokeTool(ctx context.Context, client *Client, tool toolset.ServerTool, params map[string]any) (string, error) {
	if params == nil {
		params = map[string]any{}
	}
	fillDefaultOrganizationID(client, params)

	result, err := tool.Handler(ctx, client, params)
	return result, WrapError(err)
}

func toolRequiredParams(tool toolset.ServerTool) []string {
	if len(tool.Tool.InputSchema.Required) > 0 {
		return tool.Tool.InputSchema.Required
	}
	if len(tool.Tool.RawInputSchema) == 0 {
		return nil
	}
	var schema struct {
		Required []string `json:"required"`
	}
	if err := json.Unmarshal(tool.Tool.RawInputSchema, &schema); err != nil {
		return nil
	}
	return schema.Required
}

func isMissingToolParam(value any) bool {
	if value == nil {
		return true
	}
	if s, ok := value.(string); ok {
		return strings.TrimSpace(s) == ""
	}
	return false
}

func fillDefaultOrganizationID(client *Client, params map[string]any) {
	if client == nil || client.DefaultOrgID == "" {
		return
	}
	if orgID, ok := params["organizationId"].(string); !ok || strings.TrimSpace(orgID) == "" {
		params["organizationId"] = client.DefaultOrgID
	}
}

func filterToolsByDomains(tools []toolset.ServerTool, enabled, disabled []string) []toolset.ServerTool {
	if len(enabled) > 0 {
		allowed := make(map[string]struct{}, len(enabled))
		for _, d := range enabled {
			allowed[d] = struct{}{}
		}
		filtered := make([]toolset.ServerTool, 0, len(tools))
		for _, tool := range tools {
			if _, ok := allowed[tool.Domain]; ok {
				filtered = append(filtered, tool)
			}
		}
		return filtered
	}

	if len(disabled) > 0 {
		blocked := make(map[string]struct{}, len(disabled))
		for _, d := range disabled {
			blocked[d] = struct{}{}
		}
		filtered := make([]toolset.ServerTool, 0, len(tools))
		for _, tool := range tools {
			if _, ok := blocked[tool.Domain]; !ok {
				filtered = append(filtered, tool)
			}
		}
		return filtered
	}

	return tools
}

func validateToolFilters(tools []toolset.ServerTool, enabledTools, disabledTools []string) error {
	known := make(map[string]struct{}, len(tools))
	for _, tool := range tools {
		name := tool.Tool.Name
		if _, exists := known[name]; exists {
			return fmt.Errorf("duplicate MCP tool registered: %s", name)
		}
		known[name] = struct{}{}
	}

	for _, name := range enabledTools {
		if _, exists := known[name]; !exists {
			return fmt.Errorf("unknown MCP tool %q; known tools: %s", name, strings.Join(knownToolNames(known), ", "))
		}
	}
	for _, name := range disabledTools {
		if _, exists := known[name]; !exists {
			return fmt.Errorf("unknown MCP tool %q; known tools: %s", name, strings.Join(knownToolNames(known), ", "))
		}
	}
	return nil
}

func knownToolNames(known map[string]struct{}) []string {
	names := make([]string, 0, len(known))
	for name := range known {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func shouldEnableTool(toolName string, enabledTools, disabledTools []string) bool {
	if slices.Contains(disabledTools, toolName) {
		return false
	}
	if len(enabledTools) > 0 {
		return slices.Contains(enabledTools, toolName)
	}
	return true
}
