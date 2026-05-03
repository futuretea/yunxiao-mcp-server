package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexVersionActivityTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_versions",
				mcp.WithDescription("List versions (releases) in a Projex project. Versions represent release milestones, distinct from sprints (iterations) and milestones (planning checkpoints)."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithString("status", mcp.Description("Comma-separated version statuses. Common values: TODO, DOING, ARCHIVED.")),
				mcp.WithString("name", mcp.Description("Version name filter.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListVersions,
		},
		{
			Tool: mcp.NewTool("list_workitem_activities",
				mcp.WithDescription("List activity events (history log) for a single Projex work item. For a comprehensive view including comments and attachments, use get_project_workitem_detail instead."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("workitemId", mcp.Required(), mcp.Description("Work item ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListWorkitemActivities,
		},
	}
}
