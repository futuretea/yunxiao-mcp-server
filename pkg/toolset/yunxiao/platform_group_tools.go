package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func platformGroupTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_organization_groups",
				mcp.WithDescription("List groups in a Yunxiao organization. Groups are permission-bound collections of users and resources. Use list_organization_members to discover users who can be added to groups."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListOrganizationGroups,
		},
		{
			Tool: mcp.NewTool("list_organization_group_members",
				mcp.WithDescription("List members in a Yunxiao organization group. Use this to check who belongs to a specific group and their roles."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("groupId", mcp.Required(), mcp.Description("Group ID. Use list_organization_groups to discover valid IDs.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListOrganizationGroupMembers,
		},
	}
}
