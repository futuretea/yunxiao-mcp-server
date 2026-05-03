package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func codeupTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 32)
	tools = append(tools, codeupRepositoryTools()...)
	tools = append(tools, codeupRepositoryMetadataTools()...)
	tools = append(tools, codeupNamespaceTools()...)
	tools = append(tools, codeupGroupMemberTools()...)
	tools = append(tools, codeupAccessTools()...)
	tools = append(tools, codeupFileAndCommitTools()...)
	tools = append(tools, codeupCheckTools()...)
	tools = append(tools, codeupMergeRequestTools()...)
	tools = append(tools, codeupChangeRequestTools()...)
	tools = append(tools, codeupEnhancedTools()...)
	return tools
}

func codeupRepositoryTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_repositories",
				mcp.WithDescription("List CodeUp repositories in a Yunxiao organization."),
				mcp.WithString("organizationId",
					mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted."),
				),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: created_at, name, path, or last_activity_at.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithString("search", mcp.Description("Fuzzy repository path search keyword.")),
				mcp.WithBoolean("archived", mcp.Description("Filter archived repositories.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListRepositories,
		},
		{
			Tool: mcp.NewTool("list_branches",
				mcp.WithDescription("List branches in a CodeUp repository."),
				mcp.WithString("organizationId",
					mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted."),
				),
				mcp.WithString("repositoryId",
					mcp.Required(),
					mcp.Description("Repository numeric ID or full path such as org/repo."),
				),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithString("sort", mcp.Description("Sort mode: name_asc, name_desc, updated_asc, or updated_desc.")),
				mcp.WithString("search", mcp.Description("Branch search keyword.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListBranches,
		},
	}
}

func codeupFileAndCommitTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_files",
				mcp.WithDescription("List files in a CodeUp repository tree."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("path", mcp.Description("Directory path to query.")),
				mcp.WithString("ref", mcp.Description("Branch, tag, or commit SHA. Defaults to the repository default branch when omitted.")),
				mcp.WithString("type", mcp.Description("Tree mode: DIRECT, RECURSIVE, or FLATTEN.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListFiles,
		},
		{
			Tool: mcp.NewTool("list_commits",
				mcp.WithDescription("List commits in a CodeUp repository."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("refName", mcp.Required(), mcp.Description("Branch, tag, or commit SHA.")),
				mcp.WithString("since", mcp.Description("Start time in ISO 8601 format.")),
				mcp.WithString("until", mcp.Description("End time in ISO 8601 format.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithString("path", mcp.Description("Filter commits touching this path.")),
				mcp.WithString("search", mcp.Description("Commit search keyword.")),
				mcp.WithBoolean("showSignature", mcp.Description("Whether to include commit signatures.")),
				mcp.WithString("committerIds", mcp.Description("Comma-separated committer user IDs.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListCommits,
		},
	}
}

func codeupChangeRequestTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 6)
	tools = append(tools, codeupChangeRequestCoreTools()...)
	tools = append(tools, codeupChangeRequestDiffTools()...)
	tools = append(tools, codeupChangeRequestCommentTools()...)
	return tools
}

func codeupChangeRequestCoreTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_change_requests",
				mcp.WithDescription("List CodeUp merge requests in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithString("projectIds", mcp.Description("Comma-separated repository IDs or full paths.")),
				mcp.WithString("authorIds", mcp.Description("Comma-separated author user IDs.")),
				mcp.WithString("reviewerIds", mcp.Description("Comma-separated reviewer user IDs.")),
				mcp.WithString("state", mcp.Description("Merge request state: opened, merged, or closed.")),
				mcp.WithString("search", mcp.Description("Title search keyword.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: created_at or updated_at.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithString("createdBefore", mcp.Description("Created-before time in ISO 8601 format.")),
				mcp.WithString("createdAfter", mcp.Description("Created-after time in ISO 8601 format.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListChangeRequests,
		},
	}
}

func codeupChangeRequestDiffTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_change_request_patch_sets",
				mcp.WithDescription("List CodeUp merge request patch sets."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("localId", mcp.Required(), mcp.Description("Merge request local ID within the repository.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListChangeRequestPatchSets,
		},
	}
}

func codeupChangeRequestCommentTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_change_request_comments",
				mcp.WithDescription("List CodeUp merge request comments."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("localId", mcp.Required(), mcp.Description("Merge request local ID within the repository.")),
				mcp.WithString("patchSetBizIds", mcp.Description("Comma-separated patch set IDs.")),
				mcp.WithString("commentType", mcp.Description("Comment type: GLOBAL_COMMENT or INLINE_COMMENT. Defaults to GLOBAL_COMMENT.")),
				mcp.WithString("state", mcp.Description("Comment state: OPENED or DRAFT. Defaults to OPENED.")),
				mcp.WithBoolean("resolved", mcp.Description("Whether to list resolved comments. Defaults to false.")),
				mcp.WithString("filePath", mcp.Description("File path filter for inline comments.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListChangeRequestComments,
		},
	}
}
