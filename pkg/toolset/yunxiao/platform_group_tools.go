package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func platformGroupTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_organization_groups",
				mcp.WithDescription("List groups in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListOrganizationGroups,
		},
		{
			Tool: mcp.NewTool("list_organization_group_members",
				mcp.WithDescription("List members in a Yunxiao organization group."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("groupId", mcp.Required(), mcp.Description("Group ID.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size. Default is 100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListOrganizationGroupMembers,
		},
	}
}
