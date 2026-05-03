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
		{
			Tool: mcp.NewTool("get_change_request_overview",
				mcp.WithDescription("Get a comprehensive overview of a CodeUp change request (merge request) including basic info, patch sets, and comments in one read-only call."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("localId", mcp.Required(), mcp.Description("Change request local ID.")),
				mcp.WithBoolean("includePatchSets", mcp.Description("Whether to include patch sets. Defaults to true.")),
				mcp.WithBoolean("includeComments", mcp.Description("Whether to include comments. Defaults to true.")),
				mcp.WithString("commentState", mcp.Description("Comment state filter: OPENED or RESOLVED. Defaults to OPENED.")),
				mcp.WithBoolean("commentResolved", mcp.Description("Whether to show resolved comments. Defaults to false.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetChangeRequestOverview,
		},
		{
			Tool: mcp.NewTool("get_commit_overview",
				mcp.WithDescription("Get a comprehensive overview of a CodeUp commit including commit details, commit statuses, and check runs in one read-only call."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("sha", mcp.Required(), mcp.Description("Commit SHA.")),
				mcp.WithBoolean("includeStatuses", mcp.Description("Whether to include commit statuses. Defaults to true.")),
				mcp.WithBoolean("includeCheckRuns", mcp.Description("Whether to include check runs. Defaults to true.")),
				mcp.WithNumber("statusLimit", mcp.Description("Max commit statuses returned. Defaults to 5.")),
				mcp.WithNumber("checkRunLimit", mcp.Description("Max check runs returned. Defaults to 5.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetCommitOverview,
		},
	}
}
