package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackWriteTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("create_change_order",
				mcp.WithDescription("Create an AppStack change order (deployment order). Change orders trigger application deployments to environments. This is a write operation and requires read_only=false."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("changeOrder", mcp.Required(), mcp.Description("JSON string with change order details: {changeOrderName, type (Deploy|Scale|Rollback|Destroy), envs (object), orchestrationRevisionSha, description}.")),
			),
			Handler: handleCreateChangeOrder,
		},
	}
}
