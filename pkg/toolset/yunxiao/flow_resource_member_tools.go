package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func flowResourceMemberTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_resource_members",
				mcp.WithDescription("List members who have access to a Flow resource (e.g., a pipeline or host group). Use this to discover who can manage or trigger a pipeline."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("resourceType", mcp.Required(), mcp.Description("Resource type. Examples: pipeline, hostGroup.")),
				mcp.WithString("resourceId", mcp.Required(), mcp.Description("Resource ID (string). Use list_pipelines or other list tools to find the resource ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListResourceMembers,
		},
	}
}
