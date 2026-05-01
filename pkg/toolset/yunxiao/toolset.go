package yunxiao

import "github.com/futuretea/yunxiao-mcp-server/pkg/toolset"

// Toolset exposes Yunxiao OpenAPI tools.
type Toolset struct {
	ReadOnly bool
}

func (t *Toolset) GetName() string {
	return "yunxiao"
}

func (t *Toolset) GetDescription() string {
	return "Yunxiao organization and DevOps OpenAPI tools"
}

func (t *Toolset) GetTools(_ any) []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 37)
	tools = append(tools, platformTools()...)
	tools = append(tools, codeupTools()...)
	tools = append(tools, flowTools()...)
	tools = append(tools, projexTools()...)
	tools = append(tools, packageTools()...)
	return tools
}
