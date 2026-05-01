package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexEnhancedTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_project_overview",
				mcp.WithDescription("Get a compact Projex project overview with common project-management lists in one read-only call."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithBoolean("includeMembers", mcp.Description("Whether to include project members. Defaults to true.")),
				mcp.WithBoolean("includeSprints", mcp.Description("Whether to include sprints. Defaults to true.")),
				mcp.WithBoolean("includeMilestones", mcp.Description("Whether to include milestones. Defaults to true.")),
				mcp.WithBoolean("includeVersions", mcp.Description("Whether to include versions. Defaults to true.")),
				mcp.WithBoolean("includeLabels", mcp.Description("Whether to include labels. Defaults to true.")),
				mcp.WithBoolean("activeOnly", mcp.Description("Whether sprint, milestone, and version sections use the status filter. Defaults to true.")),
				mcp.WithString("status", mcp.Description("Comma-separated statuses for sprint, milestone, and version sections when activeOnly is true. Defaults to TODO,DOING.")),
				mcp.WithNumber("page", mcp.Description("Page number for paged list sections. Defaults to 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for paged list sections. Defaults to 20.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetProjectOverview,
		},
	}
}
