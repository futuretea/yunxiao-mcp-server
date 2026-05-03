package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexWriteTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("create_workitem",
				mcp.WithDescription("Create a new work item in a Projex project. Before calling this, use list_work_item_types to discover the correct category and workitemTypeId for the target project."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID where the work item will be created.")),
				mcp.WithString("category", mcp.Required(), mcp.Description("Work item category. Common values: Req, Task, Bug, Risk. Use list_work_item_types to discover valid categories for the project.")),
				mcp.WithString("workitemTypeId", mcp.Required(), mcp.Description("Work item type ID (numeric string). Use list_work_item_types to find available types for the project.")),
				mcp.WithString("subject", mcp.Required(), mcp.Description("Work item title/subject.")),
				mcp.WithString("description", mcp.Description("Work item description.")),
				mcp.WithString("assignedTo", mcp.Description("Assignee user ID (string). Use list_project_members or get_project_workitem_context to discover valid user IDs.")),
				mcp.WithString("priority", mcp.Description("Priority ID (string). Use get_project_workitem_context to discover available priorities for the work item type.")),
				mcp.WithString("parentId", mcp.Description("Parent work item ID (numeric string) for creating sub-items. Use search_workitems to find the parent ID.")),
				mcp.WithString("sprint", mcp.Description("Sprint ID (string) to associate the work item with. Use list_sprints to discover active sprints.")),
				mcp.WithArray("labels", mcp.Description("Label IDs to attach to the work item.")),
			),
			Handler: handleCreateWorkitem,
		},
		{
			Tool: mcp.NewTool("update_workitem",
				mcp.WithDescription("Update an existing work item's fields. Only provided fields will be changed; omitted fields retain their current values. Use search_workitems to find the workitemId if unknown."),
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
				mcp.WithDescription("Change the status of a work item. Use get_project_workitem_context to discover valid statusId values for the work item type before updating."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("workitemId", mcp.Required(), mcp.Description("Work item ID (numeric string). Find it via search_workitems or get_project_workitem_summary.")),
				mcp.WithString("statusId", mcp.Required(), mcp.Description("Target status ID (string). Use get_project_workitem_context to discover valid status IDs for the work item type.")),
				mcp.WithString("comment", mcp.Description("Optional comment explaining the status change.")),
			),
			Handler: handleUpdateWorkitemStatus,
		},
		{
			Tool: mcp.NewTool("add_workitem_comment",
				mcp.WithDescription("Add a comment to a work item. Use search_workitems to find the workitemId if unknown."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("workitemId", mcp.Required(), mcp.Description("Work item ID (numeric string). Find it via search_workitems or get_project_workitem_summary.")),
				mcp.WithString("content", mcp.Required(), mcp.Description("Comment content (plain text or rich text format supported by Yunxiao).")),

			),
			Handler: handleAddWorkitemComment,
		},
	}
}
