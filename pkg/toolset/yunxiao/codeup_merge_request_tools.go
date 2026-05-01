package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func codeupMergeRequestTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_merge_request",
				mcp.WithDescription("Get legacy CodeUp merge request details."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("iid", mcp.Required(), mcp.Description("Legacy merge request IID within the repository.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetMergeRequest,
		},
	}
}
