package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexEnhancedTools() []toolset.ServerTool {
	tools := []toolset.ServerTool{
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
		{
			Tool: mcp.NewTool("get_project_workitem_summary",
				mcp.WithDescription("Summarize Projex work items by category for one project with samples and pagination totals."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithString("categories", mcp.Description("Comma-separated work item categories. Defaults to Req,Task,Bug,Risk.")),
				mcp.WithString("subject", mcp.Description("Subject contains filter applied to every category.")),
				mcp.WithString("status", mcp.Description("Comma-separated status IDs applied to every category.")),
				mcp.WithString("assignedTo", mcp.Description("Comma-separated assignee user IDs applied to every category.")),
				mcp.WithString("creator", mcp.Description("Comma-separated creator user IDs applied to every category.")),
				mcp.WithString("tag", mcp.Description("Comma-separated tag IDs applied to every category.")),
				mcp.WithString("conditions", mcp.Description("Advanced conditions JSON string applied to every category. Overrides simple filters.")),
				mcp.WithString("orderBy", mcp.Description("Sort field. Defaults are controlled by Yunxiao.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithNumber("sampleLimit", mcp.Description("Samples returned per category. Defaults to 5, clamped to 0-200.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetProjectWorkitemSummary,
		},
		{
			Tool: mcp.NewTool("get_project_workitem_context",
				mcp.WithDescription("Get project work item metadata context: types, labels, members, and optional fields/workflow for one type."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithString("category", mcp.Required(), mcp.Description("Work item category, such as Req, Task, Bug, or Risk.")),
				mcp.WithString("workItemTypeId", mcp.Description("Optional work item type ID for field and workflow metadata.")),
				mcp.WithBoolean("includeMembers", mcp.Description("Whether to include project members. Defaults to true.")),
				mcp.WithBoolean("includeLabels", mcp.Description("Whether to include project labels. Defaults to true.")),
				mcp.WithBoolean("includeFields", mcp.Description("Whether to include field configuration when workItemTypeId is set. Defaults to true.")),
				mcp.WithBoolean("includeWorkflow", mcp.Description("Whether to include workflow metadata when workItemTypeId is set. Defaults to true.")),
				mcp.WithNumber("page", mcp.Description("Page number for labels. Defaults to 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for labels. Defaults to 20.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetProjectWorkitemContext,
		},
	}
	tools = append(tools, projexInsightTools()...)
	return tools
}
