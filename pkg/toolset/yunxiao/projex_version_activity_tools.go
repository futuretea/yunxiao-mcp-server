package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexVersionActivityTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_versions",
				mcp.WithDescription("List versions in a Projex project."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithString("status", mcp.Description("Comma-separated version statuses: TODO, DOING, ARCHIVED.")),
				mcp.WithString("name", mcp.Description("Version name filter.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListVersions,
		},
		{
			Tool: mcp.NewTool("list_workitem_activities",
				mcp.WithDescription("List activity events for a Projex work item."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Work item ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListWorkitemActivities,
		},
	}
}
