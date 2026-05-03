package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func codeupGroupMemberTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_group_members",
				mcp.WithDescription("List CodeUp group members."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("groupId", mcp.Required(), mcp.Description("Group ID or URL-encoded full path.")),
				mcp.WithNumber("accessLevel", mcp.Description("Minimum access level: 20 viewer, 30 developer, 40 admin. Defaults to no filter.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListGroupMembers,
		},
		{
			Tool: mcp.NewTool("get_member_https_clone_username",
				mcp.WithDescription("Get a CodeUp user's HTTPS clone username."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("userId", mcp.Required(), mcp.Description("Yunxiao user ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetMemberHTTPSCloneUsername,
		},
	}
}
