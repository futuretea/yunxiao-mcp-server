package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_applications",
				mcp.WithDescription("List AppStack applications in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Yunxiao currently supports keyset.")),
				mcp.WithNumber("perPage", mcp.Description("Page size, up to 100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: id or gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithNumber("page", mcp.Description("Page number when using page pagination.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListApplications,
		},
		{
			Tool: mcp.NewTool("get_application",
				mcp.WithDescription("Get an AppStack application by name."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetApplication,
		},
	}
}
