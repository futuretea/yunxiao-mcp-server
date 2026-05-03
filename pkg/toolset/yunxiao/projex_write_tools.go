package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexWriteTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("create_workitem",
				mcp.WithDescription("Create a new work item in a Projex project."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID where the work item will be created.")),
				mcp.WithString("category", mcp.Required(), mcp.Description("Work item category, such as Task, Bug, Req, or Risk.")),
				mcp.WithString("workitemTypeId", mcp.Required(), mcp.Description("Work item type ID. Use list_work_item_types to find available types.")),
				mcp.WithString("subject", mcp.Required(), mcp.Description("Work item title/subject.")),
				mcp.WithString("description", mcp.Description("Work item description.")),
				mcp.WithString("assignedTo", mcp.Description("Assignee user ID.")),
				mcp.WithString("priority", mcp.Description("Priority ID.")),
				mcp.WithString("parentId", mcp.Description("Parent work item ID for creating sub-items.")),
				mcp.WithString("sprint", mcp.Description("Sprint ID to associate the work item with.")),
				mcp.WithArray("labels", mcp.Description("Label IDs to attach to the work item.")),
			),
			Handler: handleCreateWorkitem,
		},
		{
			Tool: mcp.NewTool("update_workitem",
				mcp.WithDescription("Update an existing work item's fields. Only provided fields will be changed."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("workitemId", mcp.Required(), mcp.Description("Work item ID to update.")),
				mcp.WithString("subject", mcp.Description("New title/subject.")),
				mcp.WithString("description", mcp.Description("New description.")),
				mcp.WithString("assignedTo", mcp.Description("New assignee user ID.")),
				mcp.WithString("priority", mcp.Description("New priority ID.")),
				mcp.WithString("sprint", mcp.Description("New sprint ID.")),
				mcp.WithArray("labels", mcp.Description("New label IDs (replaces existing labels).")),
			),
			Handler: handleUpdateWorkitem,
		},
		{
			Tool: mcp.NewTool("update_workitem_status",
				mcp.WithDescription("Change the status of a work item."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("workitemId", mcp.Required(), mcp.Description("Work item ID.")),
				mcp.WithString("statusId", mcp.Required(), mcp.Description("Target status ID. Use get_project_workitem_context to find valid status IDs for the work item type.")),
				mcp.WithString("comment", mcp.Description("Optional comment explaining the status change.")),
			),
			Handler: handleUpdateWorkitemStatus,
		},
		{
			Tool: mcp.NewTool("add_workitem_comment",
				mcp.WithDescription("Add a comment to a work item."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("workitemId", mcp.Required(), mcp.Description("Work item ID.")),
				mcp.WithString("content", mcp.Required(), mcp.Description("Comment content.")),
			),
			Handler: handleAddWorkitemComment,
		},
	}
}
