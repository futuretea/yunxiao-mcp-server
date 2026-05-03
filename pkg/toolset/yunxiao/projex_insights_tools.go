package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexInsightTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_project_risk_dashboard",
				mcp.WithDescription("Get a read-only project risk dashboard with category samples, overdue work items, and optional high-priority/stale sections. Best used for project health checks and sprint retrospectives."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID. Use search_projects to discover valid IDs.")),
				mcp.WithString("categories", mcp.Description("Comma-separated categories for category totals. Defaults to Risk,Bug,Task.")),
				mcp.WithString("subject", mcp.Description("Subject contains filter applied to every section.")),
				mcp.WithString("status", mcp.Description("Comma-separated active status IDs applied to every section.")),
				mcp.WithString("statusStage", mcp.Description("Comma-separated status stage IDs applied to every section.")),
				mcp.WithString("assignedTo", mcp.Description("Comma-separated assignee user IDs applied to every section.")),
				mcp.WithString("creator", mcp.Description("Comma-separated creator user IDs applied to every section.")),
				mcp.WithString("sprint", mcp.Description("Comma-separated sprint IDs applied to every section.")),
				mcp.WithString("workitemType", mcp.Description("Comma-separated work item type IDs applied to every section.")),
				mcp.WithString("tag", mcp.Description("Comma-separated tag IDs applied to every section.")),
				mcp.WithString("overdueBefore", mcp.Description("Planned finish date upper bound for overdue work items. Defaults to today.")),
				mcp.WithString("highPriority", mcp.Description("Comma-separated priority IDs for the high-priority section.")),
				mcp.WithString("staleBefore", mcp.Description("Status update date upper bound for the stale section.")),
				mcp.WithNumber("sampleLimit", mcp.Description("Samples returned per section. Defaults to 5, clamped to 0-200.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetProjectRiskDashboard,
		},
		{
			Tool: mcp.NewTool("get_project_member_task_status",
				mcp.WithDescription("Get per-member task status for one project with assigned, overdue, and optional status-group sections. Useful for workload balancing and identifying blocked team members."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID. Use search_projects to discover valid IDs.")),
				mcp.WithString("assigneeIds", mcp.Description("Comma-separated assignee user IDs. Defaults to project members up to memberLimit.")),
				mcp.WithString("categories", mcp.Description("Comma-separated work item categories. Defaults to Task,Bug.")),
				mcp.WithString("subject", mcp.Description("Subject contains filter applied to every section.")),
				mcp.WithString("status", mcp.Description("Comma-separated status IDs applied to assigned and overdue sections.")),
				mcp.WithString("statusStage", mcp.Description("Comma-separated status stage IDs applied to assigned and overdue sections.")),
				mcp.WithString("assignedTo", mcp.Description("Comma-separated assignee user IDs applied to every section.")),
				mcp.WithString("creator", mcp.Description("Comma-separated creator user IDs applied to every section.")),
				mcp.WithString("sprint", mcp.Description("Comma-separated sprint IDs applied to every section.")),
				mcp.WithString("workitemType", mcp.Description("Comma-separated work item type IDs applied to every section.")),
				mcp.WithString("tag", mcp.Description("Comma-separated tag IDs applied to every section.")),
				mcp.WithString("overdueBefore", mcp.Description("Planned finish date upper bound for overdue work items. Defaults to today.")),
				mcp.WithString("statusGroups", mcp.Description("Optional JSON object mapping group names to comma-separated status IDs.")),
				mcp.WithNumber("memberLimit", mcp.Description("Max project members to inspect when assigneeIds is omitted. Defaults to 20.")),
				mcp.WithNumber("sampleLimit", mcp.Description("Samples returned per member section. Defaults to 5, clamped to 0-200.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetProjectMemberTaskStatus,
		},
	}
}
