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
				mcp.WithDescription("List tags in a CodeUp repository."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
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
				mcp.WithDescription("List members of a CodeUp repository."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
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
				mcp.WithDescription("List protected branch rules in a CodeUp repository."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListProtectedBranches,
		},
		{
			Tool: mcp.NewTool("get_protected_branch",
				mcp.WithDescription("Get a protected branch rule in a CodeUp repository."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Protected branch rule ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetProtectedBranch,
		},
		{
			Tool: mcp.NewTool("list_push_rules",
				mcp.WithDescription("List push rules in a CodeUp repository."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListPushRules,
		},
		{
			Tool: mcp.NewTool("get_push_rule",
				mcp.WithDescription("Get a push rule in a CodeUp repository."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("pushRuleId", mcp.Required(), mcp.Description("Push rule ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPushRule,
		},
	}
}
