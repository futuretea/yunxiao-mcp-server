package yunxiao

import "github.com/futuretea/yunxiao-mcp-server/pkg/toolset"

func lingmaTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 6)
	tools = append(tools, lingmaUsageTools()...)
	tools = append(tools, lingmaKnowledgeBaseTools()...)
	return tools
}
