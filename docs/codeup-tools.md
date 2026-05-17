# Codeup Tools

This document describes the 35 MCP tools in the codeup domain.

Access summary: 29 read-only, 6 write-capable.

## Enhanced Tools

These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations. Prefer them when available.

| Tool | Description |
|------|-------------|
| `get_repository_overview` | Get a comprehensive overview of a CodeUp repository including basic info, branches, recent commits, and merge requests in one read-only call. This is the best starting point when exploring a new repository. |
| `get_change_request_overview` | Get a comprehensive overview of a CodeUp change request (merge request) including basic info, patch sets, and comments in one read-only call. |
| `get_commit_overview` | Get a comprehensive overview of a CodeUp commit including commit details, commit statuses, and check runs in one read-only call. |
| `get_branch_overview` | Get a comprehensive overview of a CodeUp branch including branch details, recent commits, and merge requests targeting the branch in one read-only call. |

## Pagination

Tools in this domain use the following pagination scheme(s):

- Offset (page/perPage)

## Tool Inventory

Tools marked in **bold** are enhanced aggregation tools.

| Tool | Access | Description |
|------|--------|-------------|
| `list_ssh_keys` | Read-only | List SSH keys registered for CodeUp access in a Yunxiao organization. |
| `list_user_ssh_keys` | Read-only | List SSH keys registered for a specific Yunxiao user in CodeUp. |
| `list_webhooks` | Read-only | List webhooks configured for a CodeUp repository. Webhooks trigger external integrations on repository events. |
| `list_commit_statuses` | Read-only | List commit statuses (CI checks) for a specific commit in a CodeUp repository. Use this to verify whether a commit has passed automated checks. |
| `list_check_runs` | Read-only | List check runs (CI pipeline executions) for a branch, tag, or commit in a CodeUp repository. Use this to monitor CI/CD status. |
| **`get_repository_overview`** | Read-only | Get a comprehensive overview of a CodeUp repository including basic info, branches, recent commits, and merge requests in one read-only call. This is the best starting point when exploring a new repository. |
| **`get_change_request_overview`** | Read-only | Get a comprehensive overview of a CodeUp change request (merge request) including basic info, patch sets, and comments in one read-only call. |
| **`get_commit_overview`** | Read-only | Get a comprehensive overview of a CodeUp commit including commit details, commit statuses, and check runs in one read-only call. |
| **`get_branch_overview`** | Read-only | Get a comprehensive overview of a CodeUp branch including branch details, recent commits, and merge requests targeting the branch in one read-only call. |
| `list_group_members` | Read-only | List members of a CodeUp group (namespace). Use this to discover who has access to repositories within the group. |
| `list_merge_requests` | Read-only | List legacy CodeUp merge requests across repositories in a Yunxiao organization. For change requests (new merge request format), use list_change_requests instead. |
| `get_merge_request` | Read-only | Get a single legacy CodeUp merge request by ID. Use list_merge_requests to discover valid merge request IDs. For the new format, use get_change_request_overview instead. |
| `list_template_repositories` | Read-only | List CodeUp template repositories in a Yunxiao organization. Templates are pre-configured repositories used as starting points for new projects. |
| `list_namespaces` | Read-only | List CodeUp namespaces or groups in a Yunxiao organization. Namespaces organize repositories into hierarchical groups. |
| `get_org_namespace` | Read-only | Get a CodeUp organization namespace by ID with nested sub-namespaces. Use list_namespaces to discover valid namespace IDs. |
| `list_tags` | Read-only | List tags (version markers) in a CodeUp repository. Use this to discover release versions. |
| `list_repository_members` | Read-only | List members who have access to a CodeUp repository. Use this to discover user IDs for assignment or review. |
| `list_protected_branches` | Read-only | List protected branch rules in a CodeUp repository. Protected branches enforce review and CI requirements before merging. |
| `list_push_rules` | Read-only | List push rules (commit restrictions) in a CodeUp repository. Push rules enforce commit message formats and file path restrictions. |
| `list_repositories` | Read-only | List CodeUp (Git) repositories in a Yunxiao organization. Use this to discover repositories and obtain their IDs before calling repository-scoped tools. For a comprehensive view of a single repository, use get_repository_overview instead. |
| `list_branches` | Read-only | List branches in a CodeUp repository. Use this to discover available branches before checking out code or reviewing merge requests. |
| `list_files` | Read-only | List files and directories in a CodeUp repository tree. Use this to explore repository structure. |
| `list_commits` | Read-only | List commits in a CodeUp repository. Use this to review recent changes and commit history. |
| `get_commit` | Read-only | Get a single CodeUp commit by SHA. Use list_commits to discover valid commit SHAs. For a comprehensive view with statuses and check runs, use get_commit_overview instead. |
| `compare` | Read-only | Compare two commits, branches, or tags in a CodeUp repository. Returns the diff between the two refs. |
| `list_change_requests` | Read-only | List CodeUp change requests (merge requests) across repositories in a Yunxiao organization. Use this to find pending reviews or track merged changes. |
| `list_change_request_patch_sets` | Read-only | List patch sets (diff iterations) for a CodeUp merge request. Use this to review how a merge request evolved across multiple pushes. |
| `list_change_request_comments` | Read-only | List comments on a CodeUp merge request. Use this to review feedback, inline discussions, and approval threads. |
| `get_change_request_comment` | Read-only | Get a single comment on a CodeUp change request by ID. Use list_change_request_comments to discover valid comment IDs. |
| `create_change_request` | Write-capable | Create a new Codeup change request. |
| `add_change_request_comment` | Write-capable | Add a comment to a Codeup change request. |
| `create_merge_request` | Write-capable | Create a new Codeup merge request. |
| `close_change_request` | Write-capable | Close a Codeup change request. |
| `reopen_change_request` | Write-capable | Reopen a closed Codeup change request. |
| `merge_change_request` | Write-capable | Merge a Codeup change request. |

### list_ssh_keys

**Description**: List SSH keys registered for CodeUp access in a Yunxiao organization.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. |
| `orderBy` | string | No | Sort field: created_at or updated_at. |
| `sort` | string | No | Sort direction: asc or desc. |

### list_user_ssh_keys

**Description**: List SSH keys registered for a specific Yunxiao user in CodeUp.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `userId` | string | Yes | Yunxiao user ID. Use list_organization_members or search_organization_members to discover valid user IDs. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. |
| `orderBy` | string | No | Sort field: created_at or updated_at. |
| `sort` | string | No | Sort direction: asc or desc. |

### list_webhooks

**Description**: List webhooks configured for a CodeUp repository. Webhooks trigger external integrations on repository events.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. |

### list_commit_statuses

**Description**: List commit statuses (CI checks) for a specific commit in a CodeUp repository. Use this to verify whether a commit has passed automated checks.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `sha` | string | Yes | Commit SHA (full 40-character hash). Use list_commits to discover valid SHAs. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. |

### list_check_runs

**Description**: List check runs (CI pipeline executions) for a branch, tag, or commit in a CodeUp repository. Use this to monitor CI/CD status.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `ref` | string | Yes | Commit SHA, branch name, or tag name. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. |

### get_repository_overview

**Description**: Get a comprehensive overview of a CodeUp repository including basic info, branches, recent commits, and merge requests in one read-only call. This is the best starting point when exploring a new repository.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository ID (numeric ID or full path like org/repo). Use list_repositories to find the repository ID. |
| `includeBranches` | boolean | No | Whether to include branch list. Defaults to true. |
| `includeCommits` | boolean | No | Whether to include recent commits. Defaults to true. |
| `includeMergeRequests` | boolean | No | Whether to include merge requests. Defaults to true. |
| `refName` | string | No | Branch, tag, or commit SHA for commit listing. Defaults to the repository default branch when omitted. |
| `branchLimit` | number | No | Max branches returned. Defaults to 5. |
| `commitLimit` | number | No | Max commits returned. Defaults to 5. |
| `mrLimit` | number | No | Max merge requests returned. Defaults to 5. |
| `mrState` | string | No | Merge request state filter: opened, merged, or closed. Defaults to opened. |

### get_change_request_overview

**Description**: Get a comprehensive overview of a CodeUp change request (merge request) including basic info, patch sets, and comments in one read-only call.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository ID (numeric ID or full path like org/repo). Use list_repositories to find the repository ID. |
| `localId` | string | Yes | Change request local ID within the repository. Use list_change_requests to discover valid local IDs. |
| `includePatchSets` | boolean | No | Whether to include patch sets. Defaults to true. |
| `includeComments` | boolean | No | Whether to include comments. Defaults to true. |
| `commentState` | string | No | Comment state filter: OPENED or RESOLVED. Defaults to OPENED. |
| `commentResolved` | boolean | No | Whether to show resolved comments. Defaults to false. |

### get_commit_overview

**Description**: Get a comprehensive overview of a CodeUp commit including commit details, commit statuses, and check runs in one read-only call.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository ID (numeric ID or full path like org/repo). Use list_repositories to find the repository ID. |
| `sha` | string | Yes | Commit SHA (full 40-character hash). Use list_commits to discover valid SHAs. |
| `includeStatuses` | boolean | No | Whether to include commit statuses. Defaults to true. |
| `includeCheckRuns` | boolean | No | Whether to include check runs. Defaults to true. |
| `statusLimit` | number | No | Max commit statuses returned. Defaults to 5. |
| `checkRunLimit` | number | No | Max check runs returned. Defaults to 5. |

### get_branch_overview

**Description**: Get a comprehensive overview of a CodeUp branch including branch details, recent commits, and merge requests targeting the branch in one read-only call.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository ID (numeric ID or full path like org/repo). Use list_repositories to find the repository ID. |
| `branchName` | string | Yes | Branch name, such as main or feature/demo. |
| `includeCommits` | boolean | No | Whether to include recent commits on the branch. Defaults to true. |
| `includeMergeRequests` | boolean | No | Whether to include merge requests targeting the branch. Defaults to true. |
| `commitLimit` | number | No | Max commits returned. Defaults to 5. |
| `mrLimit` | number | No | Max merge requests returned. Defaults to 5. |
| `mrState` | string | No | Merge request state filter: opened, merged, or closed. Defaults to opened. |

### list_group_members

**Description**: List members of a CodeUp group (namespace). Use this to discover who has access to repositories within the group.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `groupId` | string | Yes | Group ID or URL-encoded full path. Use list_namespaces to discover valid group IDs. |
| `accessLevel` | number | No | Minimum access level: 20 viewer, 30 developer, 40 admin. Defaults to no filter. |

### list_merge_requests

**Description**: List legacy CodeUp merge requests across repositories in a Yunxiao organization. For change requests (new merge request format), use list_change_requests instead.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. |
| `repositoryIds` | array | No | Repository IDs as strings to preserve int64 precision. |
| `authorUserIds` | array | No | Author user IDs. |
| `assigneeUserIds` | array | No | Assignee user IDs. |
| `subscriberUserIds` | array | No | Subscriber user IDs. |
| `state` | string | No | Merge request state: merged, opened, closed, reopened, accepted, canceled, or all. |
| `search` | string | No | Title search keyword. |
| `orderBy` | string | No | Sort field: id or updated_at. |
| `createdAfter` | string | No | Created-after date in yyyy-MM-dd format. |
| `createdBefore` | string | No | Created-before date in yyyy-MM-dd format. |
| `targetBranch` | string | No | Target branch filter. |

### get_merge_request

**Description**: Get a single legacy CodeUp merge request by ID. Use list_merge_requests to discover valid merge request IDs. For the new format, use get_change_request_overview instead.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `mergeRequestId` | string | Yes | Merge request local ID. Use list_merge_requests to discover valid IDs. |

### list_template_repositories

**Description**: List CodeUp template repositories in a Yunxiao organization. Templates are pre-configured repositories used as starting points for new projects.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `templateType` | number | Yes | Template type: 1 for custom templates, 2 for built-in templates. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. |

### list_namespaces

**Description**: List CodeUp namespaces or groups in a Yunxiao organization. Namespaces organize repositories into hierarchical groups.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `parentId` | number | No | Parent namespace ID. Omit to list namespaces available to the current user. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. |
| `search` | string | No | Namespace search keyword. |
| `orderBy` | string | No | Sort field: created_at or updated_at. |
| `sort` | string | No | Sort direction: asc or desc. |

### get_org_namespace

**Description**: Get a CodeUp organization namespace by ID with nested sub-namespaces. Use list_namespaces to discover valid namespace IDs.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `namespaceId` | string | Yes | Namespace ID (string). Use list_namespaces to discover valid IDs. |

### list_tags

**Description**: List tags (version markers) in a CodeUp repository. Use this to discover release versions.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository ID (numeric ID or full path like org/repo). Use list_repositories to find the repository ID. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `search` | string | No | Tag search keyword. |
| `sort` | string | No | Sort direction: asc or desc. |
| `orderBy` | string | No | Sort field: name or create. |

### list_repository_members

**Description**: List members who have access to a CodeUp repository. Use this to discover user IDs for assignment or review.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository ID (numeric ID or full path like org/repo). Use list_repositories to find the repository ID. |
| `accessLevel` | number | No | Minimum access level: 20, 30, or 40. |

### list_protected_branches

**Description**: List protected branch rules in a CodeUp repository. Protected branches enforce review and CI requirements before merging.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository ID (numeric ID or full path like org/repo). Use list_repositories to find the repository ID. |

### list_push_rules

**Description**: List push rules (commit restrictions) in a CodeUp repository. Push rules enforce commit message formats and file path restrictions.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository ID (numeric ID or full path like org/repo). Use list_repositories to find the repository ID. |

### list_repositories

**Description**: List CodeUp (Git) repositories in a Yunxiao organization. Use this to discover repositories and obtain their IDs before calling repository-scoped tools. For a comprehensive view of a single repository, use get_repository_overview instead.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. |
| `orderBy` | string | No | Sort field. Common values: created_at, name, path, last_activity_at. |
| `sort` | string | No | Sort direction: asc (ascending) or desc (descending). |
| `search` | string | No | Fuzzy repository path search keyword. |
| `archived` | boolean | No | Filter archived repositories. |

### list_branches

**Description**: List branches in a CodeUp repository. Use this to discover available branches before checking out code or reviewing merge requests.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository ID (numeric ID or full path like org/repo). Use list_repositories to find the repository ID. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `sort` | string | No | Sort mode: name_asc, name_desc, updated_asc, or updated_desc. |
| `search` | string | No | Branch search keyword. |

### list_files

**Description**: List files and directories in a CodeUp repository tree. Use this to explore repository structure.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `path` | string | No | Directory path to query. |
| `ref` | string | No | Branch, tag, or commit SHA. Defaults to the repository default branch when omitted. |
| `type` | string | No | Tree mode: DIRECT, RECURSIVE, or FLATTEN. |

### list_commits

**Description**: List commits in a CodeUp repository. Use this to review recent changes and commit history.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `refName` | string | Yes | Branch, tag, or commit SHA. |
| `since` | string | No | Start time in ISO 8601 format. |
| `until` | string | No | End time in ISO 8601 format. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `path` | string | No | Filter commits touching this path. |
| `search` | string | No | Commit search keyword. |
| `showSignature` | boolean | No | Whether to include commit signatures. |
| `committerIds` | string | No | Comma-separated committer user IDs. |

### get_commit

**Description**: Get a single CodeUp commit by SHA. Use list_commits to discover valid commit SHAs. For a comprehensive view with statuses and check runs, use get_commit_overview instead.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `sha` | string | Yes | Commit SHA (full 40-character hash). Use list_commits to discover valid SHAs. |

### compare

**Description**: Compare two commits, branches, or tags in a CodeUp repository. Returns the diff between the two refs.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `from` | string | Yes | Source commit SHA, branch, or tag. |
| `to` | string | Yes | Target commit SHA, branch, or tag. |
| `sourceType` | string | No | Source ref type: branch, tag, or commit. |
| `targetType` | string | No | Target ref type: branch, tag, or commit. |
| `straight` | string | No | Whether to compare directly without merge base. |

### list_change_requests

**Description**: List CodeUp change requests (merge requests) across repositories in a Yunxiao organization. Use this to find pending reviews or track merged changes.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `projectIds` | string | No | Comma-separated repository IDs or full paths (e.g., org/repo). Use list_repositories to discover repositories. |
| `authorIds` | string | No | Comma-separated author user IDs. |
| `reviewerIds` | string | No | Comma-separated reviewer user IDs. |
| `state` | string | No | Merge request state: opened, merged, or closed. |
| `search` | string | No | Title search keyword. |
| `orderBy` | string | No | Sort field: created_at or updated_at. |
| `sort` | string | No | Sort direction: asc (ascending) or desc (descending). |
| `createdBefore` | string | No | Created-before time in ISO 8601 format. |
| `createdAfter` | string | No | Created-after time in ISO 8601 format. |

### list_change_request_patch_sets

**Description**: List patch sets (diff iterations) for a CodeUp merge request. Use this to review how a merge request evolved across multiple pushes.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. Use list_repositories to discover valid repositories. |
| `localId` | string | Yes | Merge request local ID within the repository. Use list_change_requests to discover valid local IDs. |

### list_change_request_comments

**Description**: List comments on a CodeUp merge request. Use this to review feedback, inline discussions, and approval threads.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. Use list_repositories to discover valid repositories. |
| `localId` | string | Yes | Merge request local ID within the repository. Use list_change_requests to discover valid local IDs. |
| `patchSetBizIds` | string | No | Comma-separated patch set IDs to filter comments by. Use list_change_request_patch_sets to discover valid patch set IDs. |
| `commentType` | string | No | Comment type: GLOBAL_COMMENT (general comments) or INLINE_COMMENT (code-level comments). Defaults to GLOBAL_COMMENT. |
| `state` | string | No | Comment state: OPENED or DRAFT. Defaults to OPENED. |
| `resolved` | boolean | No | Whether to list resolved comments. Defaults to false. Set to true to see resolved threads. |
| `filePath` | string | No | File path filter for inline comments. Use this to narrow comments to a specific file. |

### get_change_request_comment

**Description**: Get a single comment on a CodeUp change request by ID. Use list_change_request_comments to discover valid comment IDs.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `localId` | string | Yes | Change request local ID. Use list_change_requests to discover valid local IDs. |
| `commentId` | string | Yes | Comment ID. Use list_change_request_comments to discover valid comment IDs. |

### create_change_request

**Description**: Create a new Codeup change request.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | Yes | Yunxiao organization ID. |
| `repositoryId` | string | Yes | Repository ID or path (e.g., 'org%2Frepo'). |
| `title` | string | Yes | Change request title. |
| `sourceBranch` | string | Yes | Source branch name. |
| `targetBranch` | string | Yes | Target branch name. |
| `sourceProjectId` | string | No | Source project numeric ID. Defaults to repository numeric ID if omitted. |
| `targetProjectId` | string | No | Target project numeric ID. Defaults to repository numeric ID if omitted. |
| `description` | string | No | Change request description. |

### add_change_request_comment

**Description**: Add a comment to a Codeup change request.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | Yes | Yunxiao organization ID. |
| `repositoryId` | string | Yes | Repository ID or path (e.g., 'org%2Frepo'). |
| `localId` | string | Yes | Change request local ID. |
| `content` | string | Yes | Comment content (plain text or Markdown). |

### create_merge_request

**Description**: Create a new Codeup merge request.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | Yes | Yunxiao organization ID. |
| `repositoryId` | string | Yes | Repository ID or path (e.g., 'org%2Frepo'). |
| `title` | string | Yes | Merge request title. |
| `sourceBranch` | string | Yes | Source branch name. |
| `targetBranch` | string | Yes | Target branch name. |
| `description` | string | No | Merge request description. |
| `sourceProjectId` | string | No | Source project numeric ID. Defaults to repository numeric ID if omitted. |
| `targetProjectId` | string | No | Target project numeric ID. Defaults to repository numeric ID if omitted. |
| `assigneeIds` | array | No | List of reviewer user IDs. |

### close_change_request

**Description**: Close a Codeup change request.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | Yes | Yunxiao organization ID. |
| `repositoryId` | string | Yes | Repository ID or path (e.g., 'org%2Frepo'). |
| `localId` | string | Yes | Change request local ID. |

### reopen_change_request

**Description**: Reopen a closed Codeup change request.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | Yes | Yunxiao organization ID. |
| `repositoryId` | string | Yes | Repository ID or path (e.g., 'org%2Frepo'). |
| `localId` | string | Yes | Change request local ID. |

### merge_change_request

**Description**: Merge a Codeup change request.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | Yes | Yunxiao organization ID. |
| `repositoryId` | string | Yes | Repository ID or path (e.g., 'org%2Frepo'). |
| `localId` | string | Yes | Change request local ID. |

