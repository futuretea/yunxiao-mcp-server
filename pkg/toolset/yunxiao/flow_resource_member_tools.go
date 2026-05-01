package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func flowResourceMemberTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_resource_members",
				mcp.WithDescription("List members for a Flow resource such as a pipeline or host group."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("resourceType", mcp.Required(), mcp.Description("Resource type, such as pipeline or hostGroup.")),
				mcp.WithString("resourceId", mcp.Required(), mcp.Description("Resource ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListResourceMembers,
		},
	}
}
