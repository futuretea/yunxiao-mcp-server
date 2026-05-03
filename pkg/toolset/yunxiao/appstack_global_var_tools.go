package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackGlobalVarTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_global_vars",
				mcp.WithDescription("Search AppStack global variable groups."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithNumber("current", mcp.Description("Current page number. Defaults to 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size. Defaults to 10.")),
				mcp.WithString("search", mcp.Description("Optional search keyword.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListGlobalVars,
		},
	}
}
