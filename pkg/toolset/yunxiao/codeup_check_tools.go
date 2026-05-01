package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func codeupCheckTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_commit_statuses",
				mcp.WithDescription("List CodeUp commit statuses for a repository commit."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("sha", mcp.Required(), mcp.Description("Commit SHA.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListCommitStatuses,
		},
		{
			Tool: mcp.NewTool("list_check_runs",
				mcp.WithDescription("List CodeUp check runs for a repository ref."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("ref", mcp.Required(), mcp.Description("Commit SHA, branch name, or tag name.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListCheckRuns,
		},
		{
			Tool: mcp.NewTool("get_check_run",
				mcp.WithDescription("Get a CodeUp check run by ID."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("checkRunId", mcp.Required(), mcp.Description("Check run ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetCheckRun,
		},
	}
}
