package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func codeupCheckTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_commit_statuses",
				mcp.WithDescription("List commit statuses (CI checks) for a specific commit in a CodeUp repository. Use this to verify whether a commit has passed automated checks."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("sha", mcp.Required(), mcp.Description("Commit SHA.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListCommitStatuses,
		},
		{
			Tool: mcp.NewTool("list_check_runs",
				mcp.WithDescription("List check runs (CI pipeline executions) for a branch, tag, or commit in a CodeUp repository. Use this to monitor CI/CD status."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("ref", mcp.Required(), mcp.Description("Commit SHA, branch name, or tag name.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListCheckRuns,
		},
	}
}
