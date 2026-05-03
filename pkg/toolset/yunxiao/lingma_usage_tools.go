package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func lingmaUsageTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_developer_members",
				mcp.WithDescription("List Tongyi Lingma developer members."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("departmentId", mcp.Description("Department ID filter.")),
				mcp.WithString("userId", mcp.Description("User ID filter.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size. Default is 100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListDeveloperMembers,
		},
	}
}
