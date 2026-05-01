package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackGlobalVarTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_global_var",
				mcp.WithDescription("Get an AppStack global variable group."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("name", mcp.Required(), mcp.Description("Global variable group name.")),
				mcp.WithString("revisionSha", mcp.Description("Optional global variable group revision SHA.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetGlobalVar,
		},
		{
			Tool: mcp.NewTool("list_global_vars",
				mcp.WithDescription("Search AppStack global variable groups."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithNumber("current", mcp.Description("Current page number. Defaults to 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size. Defaults to 10.")),
				mcp.WithString("search", mcp.Description("Optional search keyword.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListGlobalVars,
		},
	}
}
