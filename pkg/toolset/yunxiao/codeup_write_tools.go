package yunxiao

import (
	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
	"github.com/mark3labs/mcp-go/mcp"
)

func codeupWriteTools() []toolset.ServerTool {
	return []toolset.ServerTool{
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
	}
}
