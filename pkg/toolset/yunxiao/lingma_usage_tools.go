package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func lingmaUsageTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_developer_members",
				mcp.WithDescription("List Tongyi Lingma developer members in a Yunxiao organization. Use this to analyze AI coding assistant adoption across teams."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("departmentId", mcp.Description("Department ID filter.")),
				mcp.WithString("userId", mcp.Description("User ID filter.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListDeveloperMembers,
		},
	}
}
