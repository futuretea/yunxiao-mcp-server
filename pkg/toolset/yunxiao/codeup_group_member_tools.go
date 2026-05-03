package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func codeupGroupMemberTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_group_members",
				mcp.WithDescription("List members of a CodeUp group (namespace). Use this to discover who has access to repositories within the group."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("groupId", mcp.Required(), mcp.Description("Group ID or URL-encoded full path. Use list_namespaces to discover valid group IDs.")),
				mcp.WithNumber("accessLevel", mcp.Description("Minimum access level: 20 viewer, 30 developer, 40 admin. Defaults to no filter.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListGroupMembers,
		},
	}
}
