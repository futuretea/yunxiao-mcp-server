package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func platformTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_current_user",
				mcp.WithDescription("Get the current Yunxiao user for the configured access token."),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetCurrentUser,
		},
		{
			Tool: mcp.NewTool("get_current_organization_info",
				mcp.WithDescription("Get current user context, including the last organization returned by Yunxiao."),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetCurrentOrganizationInfo,
		},
		{
			Tool: mcp.NewTool("get_user_organizations",
				mcp.WithDescription("Get Yunxiao organizations visible to the current user."),
				mcp.WithNumber("page",
					mcp.Description("Page number. Defaults to 1 when omitted by Yunxiao."),
				),
				mcp.WithNumber("perPage",
					mcp.Description("Page size from 1 to 100. Defaults to 100 when omitted by Yunxiao."),
				),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetUserOrganizations,
		},
		{
			Tool: mcp.NewTool("list_organizations",
				mcp.WithDescription("List Yunxiao organizations visible to the current user."),
				mcp.WithNumber("page",
					mcp.Description("Page number. Defaults to 1 when omitted by Yunxiao."),
				),
				mcp.WithNumber("perPage",
					mcp.Description("Page size from 1 to 100. Defaults to 100 when omitted by Yunxiao."),
				),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListOrganizations,
		},
		{
			Tool: mcp.NewTool("get_organization",
				mcp.WithDescription("Get a Yunxiao organization by ID."),
				mcp.WithString("id",
					mcp.Required(),
					mcp.Description("Organization ID."),
				),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetOrganization,
		},
	}
}
