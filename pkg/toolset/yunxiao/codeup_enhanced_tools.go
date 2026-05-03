package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func codeupEnhancedTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_repository_overview",
				mcp.WithDescription("Get a comprehensive overview of a CodeUp repository including basic info, branches, recent commits, and merge requests in one read-only call."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithBoolean("includeBranches", mcp.Description("Whether to include branch list. Defaults to true.")),
				mcp.WithBoolean("includeCommits", mcp.Description("Whether to include recent commits. Defaults to true.")),
				mcp.WithBoolean("includeMergeRequests", mcp.Description("Whether to include merge requests. Defaults to true.")),
				mcp.WithString("refName", mcp.Description("Branch, tag, or commit SHA for commit listing. Defaults to the repository default branch when omitted.")),
				mcp.WithNumber("branchLimit", mcp.Description("Max branches returned. Defaults to 5.")),
				mcp.WithNumber("commitLimit", mcp.Description("Max commits returned. Defaults to 5.")),
				mcp.WithNumber("mrLimit", mcp.Description("Max merge requests returned. Defaults to 5.")),
				mcp.WithString("mrState", mcp.Description("Merge request state filter: opened, merged, or closed. Defaults to opened.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetRepositoryOverview,
		},
	}
}
