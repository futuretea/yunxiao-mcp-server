package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func codeupRepositoryMetadataTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 6)
	tools = append(tools, codeupRepositoryTagTools()...)
	tools = append(tools, codeupRepositoryMemberTools()...)
	tools = append(tools, codeupRepositoryPolicyTools()...)
	return tools
}

func codeupRepositoryTagTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_tags",
				mcp.WithDescription("List tags (version markers) in a CodeUp repository. Use this to discover release versions."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository ID (numeric ID or full path like org/repo). Use list_repositories to find the repository ID.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithString("search", mcp.Description("Tag search keyword.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: name or create.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListTags,
		},
	}
}

func codeupRepositoryMemberTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_repository_members",
				mcp.WithDescription("List members who have access to a CodeUp repository. Use this to discover user IDs for assignment or review."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository ID (numeric ID or full path like org/repo). Use list_repositories to find the repository ID.")),
				mcp.WithNumber("accessLevel", mcp.Description("Minimum access level: 20, 30, or 40.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListRepositoryMembers,
		},
	}
}

func codeupRepositoryPolicyTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_protected_branches",
				mcp.WithDescription("List protected branch rules in a CodeUp repository. Protected branches enforce review and CI requirements before merging."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository ID (numeric ID or full path like org/repo). Use list_repositories to find the repository ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListProtectedBranches,
		},
		{
			Tool: mcp.NewTool("list_push_rules",
				mcp.WithDescription("List push rules (commit restrictions) in a CodeUp repository. Push rules enforce commit message formats and file path restrictions."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository ID (numeric ID or full path like org/repo). Use list_repositories to find the repository ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListPushRules,
		},
	}
}
