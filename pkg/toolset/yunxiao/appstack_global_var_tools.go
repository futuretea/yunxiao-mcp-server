package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackGlobalVarTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_global_vars",
				mcp.WithDescription("Search AppStack global variable groups. Use this to discover variable group IDs before reading or updating specific groups."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithNumber("current", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithString("search", mcp.Description("Optional search keyword for filtering variable groups by name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListGlobalVars,
		},
		{
			Tool: mcp.NewTool("get_global_var",
				mcp.WithDescription("Get an AppStack global variable by name. Use list_global_vars to discover valid global variable names."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("name", mcp.Required(), mcp.Description("Global variable name. Use list_global_vars to discover valid names.")),
				mcp.WithString("revisionSha", mcp.Description("Optional revision SHA to retrieve a specific version of the variable.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetGlobalVar,
		},
	}
}
