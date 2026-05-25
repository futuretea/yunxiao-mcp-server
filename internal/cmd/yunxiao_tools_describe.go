package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
	yunxiaoTools "github.com/futuretea/yunxiao-mcp-server/pkg/toolset/yunxiao"
)

type toolDescription struct {
	Name        string                 `json:"name"`
	Domain      string                 `json:"domain"`
	Access      string                 `json:"access"`
	Description string                 `json:"description,omitempty"`
	Required    []string               `json:"required"`
	Parameters  []toolParameterSummary `json:"parameters"`
	InputSchema any                    `json:"inputSchema"`
}

type toolParameterSummary struct {
	Name        string `json:"name"`
	Type        string `json:"type,omitempty"`
	Required    bool   `json:"required"`
	Description string `json:"description,omitempty"`
}

func newYunxiaoToolsDescribeCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var jsonOutput bool
	command := &cobra.Command{
		Use:     "describe <tool-name>",
		Aliases: []string{"schema"},
		Short:   "describe an enabled Yunxiao tool schema",
		Example: `  # Describe a tool
  yunxiao tools describe search_projects

  # Describe with tool name and sub-command
  yunxiao tools schema search_projects`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			tools, err := yunxiaoTools.BuildToolCatalog(nil, toolCatalogOptionsFromConfig(cfg))
			if err != nil {
				return err
			}
			tool, ok := yunxiaoTools.FindTool(tools, args[0])
			if !ok {
				return describeUnknownToolError(args[0])
			}

			description := newToolDescription(tool)
			if jsonOutput {
				return printToolDescriptionJSON(streams.Out, description)
			}
			return printToolDescription(streams.Out, description)
		},
	}
	command.Flags().BoolVar(&jsonOutput, "json", false, "print tool schema as JSON")
	return command
}

func describeUnknownToolError(name string) error {
	tools, err := yunxiaoTools.BuildToolCatalog(nil, yunxiaoTools.ToolCatalogOptions{
		ReadOnly:    false,
		CompactMode: false,
	})
	if err == nil {
		if _, ok := yunxiaoTools.FindTool(tools, name); ok {
			return fmt.Errorf("yunxiao tool %q is not enabled by current CLI filters; adjust --read-only, --compact, --enabled-tools, --disabled-tools, --enable-domains, or --disable-domains", name)
		}
	}
	return fmt.Errorf("unknown yunxiao tool %q", name)
}

func newToolDescription(tool toolset.ServerTool) toolDescription {
	summary := newToolSummary(tool)
	inputSchema, properties, required := toolInputSchema(tool)
	requiredSet := requiredParamSet(required)

	parameters := make([]toolParameterSummary, 0, len(properties))
	propertyNames := make([]string, 0, len(properties))
	for name := range properties {
		propertyNames = append(propertyNames, name)
	}
	sort.Strings(propertyNames)
	for _, name := range propertyNames {
		propertySchema := properties[name]
		_, isRequired := requiredSet[name]
		property := schemaObjectMap(propertySchema)
		description := ""
		if property != nil {
			description = schemaString(property["description"])
		}
		parameters = append(parameters, toolParameterSummary{
			Name:        name,
			Type:        schemaType(propertySchema),
			Required:    isRequired,
			Description: description,
		})
	}

	return toolDescription{
		Name:        summary.Name,
		Domain:      summary.Domain,
		Access:      summary.Access,
		Description: summary.Description,
		Required:    required,
		Parameters:  parameters,
		InputSchema: inputSchema,
	}
}

func toolInputSchema(tool toolset.ServerTool) (any, map[string]any, []string) {
	if len(tool.Tool.RawInputSchema) > 0 {
		var schema map[string]any
		if err := json.Unmarshal(tool.Tool.RawInputSchema, &schema); err == nil {
			return schema, schemaObjectMap(schema["properties"]), sortedSchemaStrings(schema["required"])
		}
		return string(tool.Tool.RawInputSchema), nil, nil
	}

	schema := tool.Tool.InputSchema
	return schema, schema.Properties, sortedStrings(schema.Required)
}

func requiredParamSet(required []string) map[string]struct{} {
	set := make(map[string]struct{}, len(required))
	for _, name := range required {
		set[name] = struct{}{}
	}
	return set
}

func schemaObjectMap(value any) map[string]any {
	if value == nil {
		return nil
	}
	if result, ok := value.(map[string]any); ok {
		return result
	}
	return nil
}

func sortedSchemaStrings(value any) []string {
	values, ok := value.([]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if s, ok := value.(string); ok {
			result = append(result, s)
		}
	}
	sort.Strings(result)
	return result
}

func sortedStrings(values []string) []string {
	result := append([]string{}, values...)
	sort.Strings(result)
	return result
}

func schemaType(value any) string {
	schema := schemaObjectMap(value)
	if schema == nil {
		return ""
	}

	typ := schemaString(schema["type"])
	if typ == "" {
		return ""
	}
	if typ == "array" {
		itemType := schemaType(schema["items"])
		if itemType != "" {
			return "array<" + itemType + ">"
		}
	}
	return typ
}

func schemaString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []any:
		values := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				values = append(values, s)
			}
		}
		return strings.Join(values, "|")
	default:
		return ""
	}
}

func printToolDescriptionJSON(out anyWriter, description toolDescription) error {
	encoded, err := json.MarshalIndent(description, "", "  ")
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(out, string(encoded))
	return nil
}

func printToolDescription(out anyWriter, description toolDescription) error {
	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintf(writer, "NAME\t%s\n", description.Name)
	_, _ = fmt.Fprintf(writer, "DOMAIN\t%s\n", description.Domain)
	_, _ = fmt.Fprintf(writer, "ACCESS\t%s\n", description.Access)
	if description.Description != "" {
		_, _ = fmt.Fprintf(writer, "DESCRIPTION\t%s\n", strings.ReplaceAll(description.Description, "\n", " "))
	}
	requiredText := "-"
	if len(description.Required) > 0 {
		requiredText = strings.Join(description.Required, ", ")
	}
	_, _ = fmt.Fprintf(writer, "REQUIRED\t%s\n", requiredText)
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, boldTableHeader("PARAMETERS"))
	_, _ = fmt.Fprintln(writer, boldTableHeader("NAME\tTYPE\tREQUIRED\tDESCRIPTION"))
	for _, parameter := range description.Parameters {
		_, _ = fmt.Fprintf(
			writer,
			"%s\t%s\t%s\t%s\n",
			parameter.Name,
			parameter.Type,
			yesNo(parameter.Required),
			strings.ReplaceAll(parameter.Description, "\n", " "),
		)
	}
	return writer.Flush()
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
