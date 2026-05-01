package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func codeupMergeRequestTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_merge_requests",
				mcp.WithDescription("List legacy CodeUp merge requests in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
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
