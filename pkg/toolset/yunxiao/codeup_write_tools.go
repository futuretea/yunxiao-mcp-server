package yunxiao

import (
	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
	"github.com/mark3labs/mcp-go/mcp"
)

func codeupWriteTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("create_change_request",
				mcp.WithDescription("Create a new Codeup change request."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository ID or path (e.g., 'org%2Frepo').")),
				mcp.WithString("title", mcp.Required(), mcp.Description("Change request title.")),
				mcp.WithString("sourceBranch", mcp.Required(), mcp.Description("Source branch name.")),
				mcp.WithString("targetBranch", mcp.Required(), mcp.Description("Target branch name.")),
				mcp.WithString("sourceProjectId", mcp.Description("Source project numeric ID. Defaults to repository numeric ID if omitted.")),
				mcp.WithString("targetProjectId", mcp.Description("Target project numeric ID. Defaults to repository numeric ID if omitted.")),
				mcp.WithString("description", mcp.Description("Change request description.")),
			),
			Handler: handleCreateChangeRequest,
		},
		{
			Tool: mcp.NewTool("add_change_request_comment",
				mcp.WithDescription("Add a comment to a Codeup change request."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository ID or path (e.g., 'org%2Frepo').")),
				mcp.WithString("localId", mcp.Required(), mcp.Description("Change request local ID.")),
				mcp.WithString("content", mcp.Required(), mcp.Description("Comment content (plain text or Markdown).")),
			),
			Handler: handleAddChangeRequestComment,
		},
		{
			Tool: mcp.NewTool("create_merge_request",
				mcp.WithDescription("Create a new Codeup merge request."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository ID or path (e.g., 'org%2Frepo').")),
				mcp.WithString("title", mcp.Required(), mcp.Description("Merge request title.")),
				mcp.WithString("sourceBranch", mcp.Required(), mcp.Description("Source branch name.")),
				mcp.WithString("targetBranch", mcp.Required(), mcp.Description("Target branch name.")),
				mcp.WithString("description", mcp.Description("Merge request description.")),
				mcp.WithString("sourceProjectId", mcp.Description("Source project numeric ID. Defaults to repository numeric ID if omitted.")),
				mcp.WithString("targetProjectId", mcp.Description("Target project numeric ID. Defaults to repository numeric ID if omitted.")),
				mcp.WithArray("assigneeIds", mcp.Description("List of reviewer user IDs.")),
			),
			Handler: handleCreateMergeRequest,
		},
		{
			Tool: mcp.NewTool("close_change_request",
				mcp.WithDescription("Close a Codeup change request."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository ID or path (e.g., 'org%2Frepo').")),
				mcp.WithString("localId", mcp.Required(), mcp.Description("Change request local ID.")),
			),
			Handler: handleCloseChangeRequest,
		},
		{
			Tool: mcp.NewTool("reopen_change_request",
				mcp.WithDescription("Reopen a closed Codeup change request."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository ID or path (e.g., 'org%2Frepo').")),
				mcp.WithString("localId", mcp.Required(), mcp.Description("Change request local ID.")),
			),
			Handler: handleReopenChangeRequest,
		},
		{
			Tool: mcp.NewTool("merge_change_request",
				mcp.WithDescription("Merge a Codeup change request."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository ID or path (e.g., 'org%2Frepo').")),
				mcp.WithString("localId", mcp.Required(), mcp.Description("Change request local ID.")),
			),
			Handler: handleMergeChangeRequest,
		},
	}
}
