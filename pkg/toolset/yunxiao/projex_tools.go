package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 39)
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
	return tools
}

func projexProjectTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("search_projects",
				mcp.WithDescription("Search Projex projects in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("name", mcp.Description("Project name contains filter.")),
				mcp.WithString("status", mcp.Description("Comma-separated project status IDs.")),
				mcp.WithString("creator", mcp.Description("Comma-separated creator user IDs.")),
				mcp.WithString("conditions", mcp.Description("Advanced conditions JSON string. Overrides simple filters.")),
				mcp.WithString("extraConditions", mcp.Description("Extra conditions JSON string.")),
				mcp.WithString("orderBy", mcp.Description("Sort field, such as gmtCreate or name.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleSearchProjects,
		},
		{
			Tool: mcp.NewTool("get_project",
				mcp.WithDescription("Get a Projex project by ID."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetProject,
		},
	}
}

func projexSprintTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_sprint",
				mcp.WithDescription("Get a Projex sprint by ID."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Sprint ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetSprint,
		},
		{
			Tool: mcp.NewTool("list_sprints",
				mcp.WithDescription("List Projex sprints in a project."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithString("status", mcp.Description("Comma-separated sprint statuses: TODO, DOING, ARCHIVED.")),
				mcp.WithString("name", mcp.Description("Sprint name filter.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
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
				mcp.WithDescription("Search work items in one Projex project space."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("category", mcp.Required(), mcp.Description("Work item category, such as Req, Task, Bug, or Risk.")),
				mcp.WithString("spaceId", mcp.Required(), mcp.Description("Project ID. This API searches one project at a time.")),
				mcp.WithString("subject", mcp.Description("Subject contains filter.")),
				mcp.WithString("status", mcp.Description("Comma-separated status IDs.")),
				mcp.WithString("assignedTo", mcp.Description("Comma-separated assignee user IDs.")),
				mcp.WithString("creator", mcp.Description("Comma-separated creator user IDs.")),
				mcp.WithString("tag", mcp.Description("Comma-separated tag IDs.")),
				mcp.WithString("conditions", mcp.Description("Advanced conditions JSON string. Overrides simple filters.")),
				mcp.WithString("orderBy", mcp.Description("Sort field. Defaults are controlled by Yunxiao.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleSearchWorkitems,
		},
		{
			Tool: mcp.NewTool("get_workitem",
				mcp.WithDescription("Get a Projex work item by ID."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Work item ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetWorkitem,
		},
		{
			Tool: mcp.NewTool("list_work_item_comments",
				mcp.WithDescription("List comments for a Projex work item."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("workItemId", mcp.Required(), mcp.Description("Work item ID.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
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
				mcp.WithDescription("List all work item types in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("categories", mcp.Description("Optional work item type categories, such as Req, Bug, or Task.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAllWorkItemTypes,
		},
		{
			Tool: mcp.NewTool("list_work_item_types",
				mcp.WithDescription("List work item types in one Projex project."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithString("category", mcp.Required(), mcp.Description("Work item category, such as Req, Bug, or Task.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListWorkItemTypes,
		},
		{
			Tool: mcp.NewTool("get_work_item_type",
				mcp.WithDescription("Get a Projex work item type by ID."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Work item type ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetWorkItemType,
		},
	}
}

func projexWorkitemTypeMetadataTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_work_item_relation_work_item_types",
				mcp.WithDescription("List work item types that can be related to a Projex work item type."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("workItemTypeId", mcp.Required(), mcp.Description("Work item type ID.")),
				mcp.WithString("relationType", mcp.Description("Relation type: PARENT, SUB, ASSOCIATED, DEPEND_ON, or DEPENDED_BY.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListWorkItemRelationWorkItemTypes,
		},
		{
			Tool: mcp.NewTool("get_work_item_type_field_config",
				mcp.WithDescription("Get field configuration for a Projex work item type."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithString("workItemTypeId", mcp.Required(), mcp.Description("Work item type ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetWorkItemTypeFieldConfig,
		},
		{
			Tool: mcp.NewTool("get_work_item_workflow",
				mcp.WithDescription("Get workflow information for a Projex work item type."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithString("workItemTypeId", mcp.Required(), mcp.Description("Work item type ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetWorkItemWorkflow,
		},
	}
}
