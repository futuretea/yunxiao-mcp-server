package yunxiao

import "github.com/futuretea/yunxiao-mcp-server/pkg/toolset"

func lingmaTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 3)
	tools = append(tools, lingmaUsageTools()...)
	return tools
}
