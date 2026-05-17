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
	tools = append(tools, codeupWriteTools()...)
	return tools
}

func codeupRepositoryTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_repositories",
				mcp.WithDescription("List CodeUp (Git) repositories in a Yunxiao organization. Use this to discover repositories and obtain their IDs before calling repository-scoped tools. For a comprehensive view of a single repository, use get_repository_overview instead."),
				mcp.WithString("organizationId",
					mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization."),
				),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field. Common values: created_at, name, path, last_activity_at.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc (ascending) or desc (descending).")),
				mcp.WithString("search", mcp.Description("Fuzzy repository path search keyword.")),
				mcp.WithBoolean("archived", mcp.Description("Filter archived repositories.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListRepositories,
		},
		{
			Tool: mcp.NewTool("list_branches",
				mcp.WithDescription("List branches in a CodeUp repository. Use this to discover available branches before checking out code or reviewing merge requests."),
				mcp.WithString("organizationId",
					mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization."),
				),
				mcp.WithString("repositoryId",
					mcp.Required(),
					mcp.Description("Repository ID (numeric ID or full path like org/repo). Use list_repositories to find the repository ID."),
				),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
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
				mcp.WithDescription("List files and directories in a CodeUp repository tree. Use this to explore repository structure."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
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
				mcp.WithDescription("List commits in a CodeUp repository. Use this to review recent changes and commit history."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("refName", mcp.Required(), mcp.Description("Branch, tag, or commit SHA.")),
				mcp.WithString("since", mcp.Description("Start time in ISO 8601 format.")),
				mcp.WithString("until", mcp.Description("End time in ISO 8601 format.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithString("path", mcp.Description("Filter commits touching this path.")),
				mcp.WithString("search", mcp.Description("Commit search keyword.")),
				mcp.WithBoolean("showSignature", mcp.Description("Whether to include commit signatures.")),
				mcp.WithString("committerIds", mcp.Description("Comma-separated committer user IDs.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListCommits,
		},
		{
			Tool: mcp.NewTool("get_commit",
				mcp.WithDescription("Get a single CodeUp commit by SHA. Use list_commits to discover valid commit SHAs. For a comprehensive view with statuses and check runs, use get_commit_overview instead."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("sha", mcp.Required(), mcp.Description("Commit SHA (full 40-character hash). Use list_commits to discover valid SHAs.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetCommit,
		},
		{
			Tool: mcp.NewTool("compare",
				mcp.WithDescription("Compare two commits, branches, or tags in a CodeUp repository. Returns the diff between the two refs."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("from", mcp.Required(), mcp.Description("Source commit SHA, branch, or tag.")),
				mcp.WithString("to", mcp.Required(), mcp.Description("Target commit SHA, branch, or tag.")),
				mcp.WithString("sourceType", mcp.Description("Source ref type: branch, tag, or commit.")),
				mcp.WithString("targetType", mcp.Description("Target ref type: branch, tag, or commit.")),
				mcp.WithString("straight", mcp.Description("Whether to compare directly without merge base.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleCompare,
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
				mcp.WithDescription("List CodeUp change requests (merge requests) across repositories in a Yunxiao organization. Use this to find pending reviews or track merged changes."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithString("projectIds", mcp.Description("Comma-separated repository IDs or full paths (e.g., org/repo). Use list_repositories to discover repositories.")),
				mcp.WithString("authorIds", mcp.Description("Comma-separated author user IDs.")),
				mcp.WithString("reviewerIds", mcp.Description("Comma-separated reviewer user IDs.")),
				mcp.WithString("state", mcp.Description("Merge request state: opened, merged, or closed.")),
				mcp.WithString("search", mcp.Description("Title search keyword.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: created_at or updated_at.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc (ascending) or desc (descending).")),
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
				mcp.WithDescription("List patch sets (diff iterations) for a CodeUp merge request. Use this to review how a merge request evolved across multiple pushes."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo. Use list_repositories to discover valid repositories.")),
				mcp.WithString("localId", mcp.Required(), mcp.Description("Merge request local ID within the repository. Use list_change_requests to discover valid local IDs.")),
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
				mcp.WithDescription("List comments on a CodeUp merge request. Use this to review feedback, inline discussions, and approval threads."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo. Use list_repositories to discover valid repositories.")),
				mcp.WithString("localId", mcp.Required(), mcp.Description("Merge request local ID within the repository. Use list_change_requests to discover valid local IDs.")),
				mcp.WithString("patchSetBizIds", mcp.Description("Comma-separated patch set IDs to filter comments by. Use list_change_request_patch_sets to discover valid patch set IDs.")),
				mcp.WithString("commentType", mcp.Description("Comment type: GLOBAL_COMMENT (general comments) or INLINE_COMMENT (code-level comments). Defaults to GLOBAL_COMMENT.")),
				mcp.WithString("state", mcp.Description("Comment state: OPENED or DRAFT. Defaults to OPENED.")),
				mcp.WithBoolean("resolved", mcp.Description("Whether to list resolved comments. Defaults to false. Set to true to see resolved threads.")),
				mcp.WithString("filePath", mcp.Description("File path filter for inline comments. Use this to narrow comments to a specific file.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListChangeRequestComments,
		},
	}
}
