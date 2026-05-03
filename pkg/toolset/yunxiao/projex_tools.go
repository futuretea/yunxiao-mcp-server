package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 41)
	tools = append(tools, projexProjectTools()...)
	tools = append(tools, projexEnhancedTools()...)
	tools = append(tools, projexProjectMetadataTools()...)
	tools = append(tools, projexSprintTools()...)
	tools = append(tools, projexWorkitemTools()...)
	tools = append(tools, projexVersionActivityTools()...)
	tools = append(tools, projexEffortTools()...)
	tools = append(tools, projexWorkitemMetadataTools()...)
	tools = append(tools, projexMilestoneTestcaseTools()...)
	tools = append(tools, projexWorkitemTypeTools()...)
	tools = append(tools, projexWriteTools()...)
	return tools
}

func projexProjectTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("search_projects",
				mcp.WithDescription("Search Projex projects in a Yunxiao organization. Use this to discover available projects and obtain their project IDs before calling other project-scoped tools."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("name", mcp.Description("Filter by project name (contains match).")),
				mcp.WithString("status", mcp.Description("Comma-separated project status IDs. Use without filters to list all projects and inspect their statuses.")),
				mcp.WithString("creator", mcp.Description("Comma-separated creator user IDs. Use list_project_members to find user IDs in the project.")),
				mcp.WithString("conditions", mcp.Description("Advanced filter conditions as a JSON string. Overrides simple filters when provided.")),
				mcp.WithString("extraConditions", mcp.Description("Extra filter conditions as a JSON string, combined with conditions.")),
				mcp.WithString("orderBy", mcp.Description("Sort field, such as gmtCreate (creation time) or name.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc (ascending) or desc (descending).")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Most endpoints default to 20.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleSearchProjects,
		},
	}
}

func projexSprintTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_sprints",
				mcp.WithDescription("List sprints (iterations) in a Projex project. Sprints are time-boxed development cycles. Use this when you need iteration-level planning data. For a broader view that also includes milestones and versions, use get_project_overview instead."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID (string). Use search_projects if you do not know the project ID.")),
				mcp.WithString("status", mcp.Description("Comma-separated sprint statuses. Common values: TODO (not started), DOING (active), ARCHIVED (closed).")),
				mcp.WithString("name", mcp.Description("Filter by sprint name (contains match).")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Most endpoints default to 20.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListSprints,
		},
	}
}

func projexWorkitemTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("search_workitems",
				mcp.WithDescription("Search work items in a single Projex project with flexible filtering and pagination. This is the primary tool for finding work items when you know the project. For a high-level summary across categories without detailed filtering, use get_project_workitem_summary. For a Kanban board view grouped by status, use get_project_workitem_board."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("category", mcp.Required(), mcp.Description("Work item category. Common values: Req (需求/Requirement), Task (任务), Bug (缺陷), Risk (风险). Use list_work_item_types to discover available categories for a project.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID (string). This API searches one project at a time. Use search_projects to find the project ID.")),
				mcp.WithString("subject", mcp.Description("Filter by work item subject/title (contains match).")),
				mcp.WithString("status", mcp.Description("Comma-separated status IDs. Use get_project_workitem_context to discover valid status IDs for the target category.")),
				mcp.WithString("assignedTo", mcp.Description("Comma-separated assignee user IDs. Use list_project_members to find user IDs in the project.")),
				mcp.WithString("creator", mcp.Description("Comma-separated creator user IDs. Use list_project_members to find user IDs in the project.")),
				mcp.WithString("tag", mcp.Description("Comma-separated tag IDs. Use get_project_workitem_context to discover available tags.")),
				mcp.WithString("sprint", mcp.Description("Comma-separated sprint IDs. Use list_sprints to discover sprint IDs in the project.")),
				mcp.WithString("workitemType", mcp.Description("Comma-separated work item type IDs. Use list_work_item_types to discover available types.")),
				mcp.WithString("statusStage", mcp.Description("Comma-separated status stage IDs (higher-level status grouping). Use get_project_workitem_context to discover valid stages.")),
				mcp.WithString("priority", mcp.Description("Comma-separated priority IDs. Use get_project_workitem_context to discover available priorities.")),
				mcp.WithString("subjectDescription", mcp.Description("Filter by subject or description content (contains match).")),
				mcp.WithString("createdAfter", mcp.Description("Created date lower bound. Accepts YYYY-MM-DD (e.g., 2024-01-01) or ISO datetime.")),
				mcp.WithString("createdBefore", mcp.Description("Created date upper bound. Accepts YYYY-MM-DD or ISO datetime.")),
				mcp.WithString("updatedAfter", mcp.Description("Updated date lower bound. Accepts YYYY-MM-DD or ISO datetime.")),
				mcp.WithString("updatedBefore", mcp.Description("Updated date upper bound. Accepts YYYY-MM-DD or ISO datetime.")),
				mcp.WithString("finishTimeAfter", mcp.Description("Planned finish date lower bound. Accepts YYYY-MM-DD or ISO datetime.")),
				mcp.WithString("finishTimeBefore", mcp.Description("Planned finish date upper bound. Accepts YYYY-MM-DD or ISO datetime.")),
				mcp.WithString("updateStatusAtAfter", mcp.Description("Status update date lower bound. Accepts YYYY-MM-DD or ISO datetime.")),
				mcp.WithString("updateStatusAtBefore", mcp.Description("Status update date upper bound. Accepts YYYY-MM-DD or ISO datetime.")),
				mcp.WithString("conditions", mcp.Description("Advanced filter conditions as a JSON string. Overrides simple filters when provided.")),
				mcp.WithString("orderBy", mcp.Description("Sort field. Defaults are controlled by Yunxiao. Common values: gmtCreate, gmtModified, priority, serialNumber.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc (ascending) or desc (descending).")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Most endpoints default to 20.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleSearchWorkitems,
		},
		{
			Tool: mcp.NewTool("list_work_item_comments",
				mcp.WithDescription("List comments for a single Projex work item. For a comprehensive view that includes comments along with activities, attachments, and relations, use get_project_workitem_detail instead."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("workitemId", mcp.Required(), mcp.Description("Work item ID (numeric string). Find it via search_workitems or get_project_workitem_summary.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Most endpoints default to 20.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListWorkItemComments,
		},
	}
}

func projexWorkitemTypeTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 6)
	tools = append(tools, projexWorkitemTypeListTools()...)
	tools = append(tools, projexWorkitemTypeMetadataTools()...)
	return tools
}

func projexWorkitemTypeListTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_all_work_item_types",
				mcp.WithDescription("List all work item types across a Yunxiao organization. This returns organization-level type definitions. For project-scoped types (including custom fields), use list_work_item_types instead."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("categories", mcp.Description("Comma-separated work item type categories to filter by. Common values: Req, Bug, Task, Risk.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAllWorkItemTypes,
		},
		{
			Tool: mcp.NewTool("list_work_item_types",
				mcp.WithDescription("List work item types available in a specific Projex project. Use this before creating a work item to discover valid workitemTypeId values and categories for the target project."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID (string). Use search_projects if you do not know the project ID.")),
				mcp.WithString("category", mcp.Required(), mcp.Description("Work item category to filter by. Common values: Req, Bug, Task, Risk. Use list_all_work_item_types to see available categories.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListWorkItemTypes,
		},
	}
}

func projexWorkitemTypeMetadataTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_work_item_relation_work_item_types",
				mcp.WithDescription("List work item types that can be related to a given work item type. Use this when setting up work item relationships to discover which types are valid for the chosen relationType."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("workItemTypeId", mcp.Required(), mcp.Description("Work item type ID (numeric string). Use list_work_item_types to find available types.")),
				mcp.WithString("relationType", mcp.Description("Relation type filter. Valid values: PARENT, SUB, ASSOCIATED, DEPEND_ON, DEPENDED_BY. Omit to list all compatible types.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListWorkItemRelationWorkItemTypes,
		},
	}
}
