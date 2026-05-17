package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func codeupMergeRequestTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_merge_requests",
				mcp.WithDescription("List legacy CodeUp merge requests across repositories in a Yunxiao organization. For change requests (new merge request format), use list_change_requests instead."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100.")),
				mcp.WithArray("repositoryIds", mcp.Description("Repository IDs as strings to preserve int64 precision."), mcp.WithStringItems()),
				mcp.WithArray("authorUserIds", mcp.Description("Author user IDs."), mcp.WithStringItems()),
				mcp.WithArray("assigneeUserIds", mcp.Description("Assignee user IDs."), mcp.WithStringItems()),
				mcp.WithArray("subscriberUserIds", mcp.Description("Subscriber user IDs."), mcp.WithStringItems()),
				mcp.WithString("state", mcp.Description("Merge request state: merged, opened, closed, reopened, accepted, canceled, or all.")),
				mcp.WithString("search", mcp.Description("Title search keyword.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: id or updated_at.")),
				mcp.WithString("createdAfter", mcp.Description("Created-after date in yyyy-MM-dd format.")),
				mcp.WithString("createdBefore", mcp.Description("Created-before date in yyyy-MM-dd format.")),
				mcp.WithString("targetBranch", mcp.Description("Target branch filter.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListMergeRequests,
		},
		{
			Tool: mcp.NewTool("get_merge_request",
				mcp.WithDescription("Get a single legacy CodeUp merge request by ID. Use list_merge_requests to discover valid merge request IDs. For the new format, use get_change_request_overview instead."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("mergeRequestId", mcp.Required(), mcp.Description("Merge request local ID. Use list_merge_requests to discover valid IDs.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetMergeRequest,
		},
	}
}
