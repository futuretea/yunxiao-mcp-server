package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func TestYunxiaoCLIToolsDescribeJSONIsOfflineAndFiltered(t *testing.T) {
	var out bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--enabled-tools", "search_projects", "tools", "describe", "search_projects", "--json"})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var description toolDescription
	if err := json.Unmarshal(out.Bytes(), &description); err != nil {
		t.Fatalf("unmarshal stdout: %v\n%s", err, out.String())
	}
	if description.Name != "search_projects" || description.Domain != "projex" || description.Access != "read-only" {
		t.Fatalf("description = %#v", description)
	}
	if description.InputSchema == nil {
		t.Fatal("InputSchema is nil")
	}

	params := parametersByName(description.Parameters)
	for name, wantType := range map[string]string{
		"organizationId": "string",
		"name":           "string",
		"page":           "number",
		"perPage":        "number",
	} {
		param, ok := params[name]
		if !ok {
			t.Fatalf("parameters missing %q: %#v", name, description.Parameters)
		}
		if param.Type != wantType {
			t.Fatalf("parameters[%q].Type = %q, want %q", name, param.Type, wantType)
		}
		if param.Required {
			t.Fatalf("parameters[%q].Required = true, want false", name)
		}
	}
	if !strings.Contains(params["organizationId"].Description, "default organization") {
		t.Fatalf("organizationId description = %q", params["organizationId"].Description)
	}
}

func TestYunxiaoCLIToolsDescribePrintsTable(t *testing.T) {
	var out bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--enabled-tools", "search_projects", "tools", "describe", "search_projects"})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"NAME", "search_projects", "DOMAIN", "projex", "REQUIRED", "-", "PARAMETERS", "organizationId", "page", "number"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
}

func TestYunxiaoCLIToolsDescribePrintsRequiredParams(t *testing.T) {
	var out bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--enabled-tools", "list_sprints", "tools", "describe", "list_sprints"})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"REQUIRED", "projectId", "projectId", "string", "yes"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
}

func TestYunxiaoCLIToolsSchemaAliasDescribesWriteToolWhenEnabled(t *testing.T) {
	var out bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--read-only=false",
		"--enabled-tools", "create_workitem",
		"tools", "schema", "create_workitem", "--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var description toolDescription
	if err := json.Unmarshal(out.Bytes(), &description); err != nil {
		t.Fatalf("unmarshal stdout: %v\n%s", err, out.String())
	}
	if description.Name != "create_workitem" || description.Access != "write" {
		t.Fatalf("description = %#v", description)
	}
	params := parametersByName(description.Parameters)
	for _, name := range []string{"projectId", "category", "workitemTypeId", "subject"} {
		if !params[name].Required {
			t.Fatalf("parameters[%q].Required = false, want true", name)
		}
	}
}

func TestYunxiaoCLIToolsDescribeRejectsReadOnlyFilteredWriteTool(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"tools", "describe", "create_workitem"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected filtered tool error")
	}
	if !strings.Contains(err.Error(), `yunxiao tool "create_workitem" is not enabled by current CLI filters`) {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIToolsDescribeCompactFiltering(t *testing.T) {
	t.Run("default compact hides raw tool", func(t *testing.T) {
		command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
		command.SetArgs([]string{"tools", "describe", "get_application"})

		err := command.Execute()
		if err == nil {
			t.Fatal("Execute() expected filtered tool error")
		}
		if !strings.Contains(err.Error(), `yunxiao tool "get_application" is not enabled by current CLI filters`) {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("compact false allows raw tool", func(t *testing.T) {
		var out bytes.Buffer
		command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
		command.SetArgs([]string{"--compact=false", "tools", "describe", "get_application", "--json"})

		if err := command.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		var description toolDescription
		if err := json.Unmarshal(out.Bytes(), &description); err != nil {
			t.Fatalf("unmarshal stdout: %v\n%s", err, out.String())
		}
		if description.Name != "get_application" {
			t.Fatalf("description.Name = %q, want get_application", description.Name)
		}
	})
}

func TestYunxiaoCLIToolsDescribeDomainFiltering(t *testing.T) {
	t.Run("enabled domain allows tool", func(t *testing.T) {
		var out bytes.Buffer
		command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
		command.SetArgs([]string{"--enable-domains", "projex", "tools", "describe", "search_projects", "--json"})

		if err := command.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if !strings.Contains(out.String(), `"domain": "projex"`) {
			t.Fatalf("stdout = %q", out.String())
		}
	})

	t.Run("enabled domain filters tool out", func(t *testing.T) {
		command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
		command.SetArgs([]string{"--enable-domains", "platform", "tools", "describe", "search_projects"})

		err := command.Execute()
		if err == nil {
			t.Fatal("Execute() expected filtered tool error")
		}
		if !strings.Contains(err.Error(), `yunxiao tool "search_projects" is not enabled by current CLI filters`) {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("disabled domain filters tool out", func(t *testing.T) {
		command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
		command.SetArgs([]string{"--disable-domains", "projex", "tools", "describe", "search_projects"})

		err := command.Execute()
		if err == nil {
			t.Fatal("Execute() expected filtered tool error")
		}
		if !strings.Contains(err.Error(), `yunxiao tool "search_projects" is not enabled by current CLI filters`) {
			t.Fatalf("error = %v", err)
		}
	})
}

func TestYunxiaoCLIToolsDescribeReturnsConfigError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--base-url", "://invalid-url", "tools", "describe", "search_projects"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected config error")
	}
	if !strings.Contains(err.Error(), "load config") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIToolsDescribeReturnsCatalogError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--enabled-tools", "not_a_tool", "tools", "describe", "not_a_tool"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected catalog error")
	}
	if !strings.Contains(err.Error(), `unknown MCP tool "not_a_tool"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIToolsDescribeRejectsFilteredTool(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--enabled-tools", "get_current_user", "tools", "describe", "search_projects"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected filtered tool error")
	}
	if !strings.Contains(err.Error(), `yunxiao tool "search_projects" is not enabled by current CLI filters`) {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIToolsDescribeRejectsUnknownTool(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"tools", "describe", "not_a_tool"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected unknown tool error")
	}
	if !strings.Contains(err.Error(), `unknown yunxiao tool "not_a_tool"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestPrintToolDescriptionJSONReturnsMarshalError(t *testing.T) {
	if err := printToolDescriptionJSON(&bytes.Buffer{}, toolDescription{InputSchema: func() {}}); err == nil {
		t.Fatal("printToolDescriptionJSON() expected marshal error")
	}
}

func TestNewToolDescriptionUsesRawInputSchema(t *testing.T) {
	description := newToolDescription(toolset.ServerTool{
		Tool: mcp.Tool{
			Name: "raw_tool",
			RawInputSchema: json.RawMessage(`{
				"type": "object",
				"required": ["path"],
				"properties": {
					"path": {"type": "string", "description": "Request path."},
					"tags": {"type": "array", "items": {"type": "string"}}
				}
			}`),
		},
		Domain: "api",
	})

	params := parametersByName(description.Parameters)
	if !params["path"].Required || params["path"].Type != "string" {
		t.Fatalf("path parameter = %#v", params["path"])
	}
	if params["tags"].Type != "array<string>" {
		t.Fatalf("tags type = %q, want array<string>", params["tags"].Type)
	}
	if len(description.Required) != 1 || description.Required[0] != "path" {
		t.Fatalf("Required = %#v, want [path]", description.Required)
	}
}

func TestNewToolDescriptionHandlesInvalidRawInputSchema(t *testing.T) {
	description := newToolDescription(toolset.ServerTool{
		Tool: mcp.Tool{
			Name:           "raw_tool",
			RawInputSchema: json.RawMessage(`{`),
		},
		Domain: "api",
	})

	if description.InputSchema != "{" {
		t.Fatalf("InputSchema = %#v, want raw string fallback", description.InputSchema)
	}
	if len(description.Parameters) != 0 || len(description.Required) != 0 {
		t.Fatalf("Parameters = %#v, Required = %#v; want empty", description.Parameters, description.Required)
	}
}

func TestToolSchemaHelpersHandleUnsupportedValues(t *testing.T) {
	if got := schemaType(map[string]any{"type": []any{"string", "null"}}); got != "string|null" {
		t.Fatalf("schemaType() = %q, want string|null", got)
	}
	if got := schemaType(map[string]any{"type": "array"}); got != "array" {
		t.Fatalf("schemaType() = %q, want array", got)
	}
	if got := schemaType("not-a-schema"); got != "" {
		t.Fatalf("schemaType() = %q, want empty", got)
	}
	if got := sortedSchemaStrings([]string{"path"}); got != nil {
		t.Fatalf("sortedSchemaStrings() = %#v, want nil for unsupported input", got)
	}
}

func parametersByName(parameters []toolParameterSummary) map[string]toolParameterSummary {
	byName := make(map[string]toolParameterSummary, len(parameters))
	for _, parameter := range parameters {
		byName[parameter.Name] = parameter
	}
	return byName
}
