# Codeup Tools

This document describes the 24 read-only MCP tools in the codeup domain.

## Enhanced Tools

These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations. Prefer them when available.

| Tool | Description |
|------|-------------|
| `get_repository_overview` | Get a comprehensive overview of a CodeUp repository including basic info, branches, recent commits, and merge requests in one read-only call. |
| `get_change_request_overview` | Get a comprehensive overview of a CodeUp change request (merge request) including basic info, patch sets, and comments in one read-only call. |
| `get_commit_overview` | Get a comprehensive overview of a CodeUp commit including commit details, commit statuses, and check runs in one read-only call. |
| `get_branch_overview` | Get a comprehensive overview of a CodeUp branch including branch details, recent commits, and merge requests targeting the branch in one read-only call. |

## Pagination

Tools in this domain use the following pagination scheme(s):

- Offset (page/perPage)

## Tool Inventory

Tools marked in **bold** are enhanced aggregation tools.

| Tool | Description |
|------|-------------|
| `list_ssh_keys` | List CodeUp SSH keys in a Yunxiao organization. |
| `list_user_ssh_keys` | List CodeUp SSH keys for a Yunxiao user. |
| `list_webhooks` | List CodeUp webhooks in a repository. |
| `list_commit_statuses` | List CodeUp commit statuses for a repository commit. |
| `list_check_runs` | List CodeUp check runs for a repository ref. |
| **`get_repository_overview`** | Get a comprehensive overview of a CodeUp repository including basic info, branches, recent commits, and merge requests in one read-only call. |
| **`get_change_request_overview`** | Get a comprehensive overview of a CodeUp change request (merge request) including basic info, patch sets, and comments in one read-only call. |
| **`get_commit_overview`** | Get a comprehensive overview of a CodeUp commit including commit details, commit statuses, and check runs in one read-only call. |
| **`get_branch_overview`** | Get a comprehensive overview of a CodeUp branch including branch details, recent commits, and merge requests targeting the branch in one read-only call. |
| `list_group_members` | List CodeUp group members. |
| `list_merge_requests` | List legacy CodeUp merge requests in a Yunxiao organization. |
| `list_template_repositories` | List CodeUp template repositories in a Yunxiao organization. |
| `list_namespaces` | List CodeUp namespaces or groups in a Yunxiao organization. |
| `list_tags` | List tags in a CodeUp repository. |
| `list_repository_members` | List members of a CodeUp repository. |
| `list_protected_branches` | List protected branch rules in a CodeUp repository. |
| `list_push_rules` | List push rules in a CodeUp repository. |
| `list_repositories` | List CodeUp repositories in a Yunxiao organization. |
| `list_branches` | List branches in a CodeUp repository. |
| `list_files` | List files in a CodeUp repository tree. |
| `list_commits` | List commits in a CodeUp repository. |
| `list_change_requests` | List CodeUp merge requests in a Yunxiao organization. |
| `list_change_request_patch_sets` | List CodeUp merge request patch sets. |
| `list_change_request_comments` | List CodeUp merge request comments. |

### list_ssh_keys

**Description**: List CodeUp SSH keys in a Yunxiao organization.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |
| `orderBy` | string | No | Sort field: created_at or updated_at. |
| `sort` | string | No | Sort direction: asc or desc. |

### list_user_ssh_keys

**Description**: List CodeUp SSH keys for a Yunxiao user.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `userId` | string | Yes | Yunxiao user ID. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |
| `orderBy` | string | No | Sort field: created_at or updated_at. |
| `sort` | string | No | Sort direction: asc or desc. |

### list_webhooks

**Description**: List CodeUp webhooks in a repository.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |

### list_commit_statuses

**Description**: List CodeUp commit statuses for a repository commit.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `sha` | string | Yes | Commit SHA. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |

### list_check_runs

**Description**: List CodeUp check runs for a repository ref.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `ref` | string | Yes | Commit SHA, branch name, or tag name. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |

### get_repository_overview

**Description**: Get a comprehensive overview of a CodeUp repository including basic info, branches, recent commits, and merge requests in one read-only call.

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
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

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `localId` | string | Yes | Change request local ID. |
| `includePatchSets` | boolean | No | Whether to include patch sets. Defaults to true. |
| `includeComments` | boolean | No | Whether to include comments. Defaults to true. |
| `commentState` | string | No | Comment state filter: OPENED or RESOLVED. Defaults to OPENED. |
| `commentResolved` | boolean | No | Whether to show resolved comments. Defaults to false. |

### get_commit_overview

**Description**: Get a comprehensive overview of a CodeUp commit including commit details, commit statuses, and check runs in one read-only call.

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `sha` | string | Yes | Commit SHA. |
| `includeStatuses` | boolean | No | Whether to include commit statuses. Defaults to true. |
| `includeCheckRuns` | boolean | No | Whether to include check runs. Defaults to true. |
| `statusLimit` | number | No | Max commit statuses returned. Defaults to 5. |
| `checkRunLimit` | number | No | Max check runs returned. Defaults to 5. |

### get_branch_overview

**Description**: Get a comprehensive overview of a CodeUp branch including branch details, recent commits, and merge requests targeting the branch in one read-only call.

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `branchName` | string | Yes | Branch name, such as main or feature/demo. |
| `includeCommits` | boolean | No | Whether to include recent commits on the branch. Defaults to true. |
| `includeMergeRequests` | boolean | No | Whether to include merge requests targeting the branch. Defaults to true. |
| `commitLimit` | number | No | Max commits returned. Defaults to 5. |
| `mrLimit` | number | No | Max merge requests returned. Defaults to 5. |
| `mrState` | string | No | Merge request state filter: opened, merged, or closed. Defaults to opened. |

### list_group_members

**Description**: List CodeUp group members.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `groupId` | string | Yes | Group ID or URL-encoded full path. |
| `accessLevel` | number | No | Minimum access level: 20 viewer, 30 developer, 40 admin. Defaults to no filter. |

### list_merge_requests

**Description**: List legacy CodeUp merge requests in a Yunxiao organization.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |
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

### list_template_repositories

**Description**: List CodeUp template repositories in a Yunxiao organization.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `templateType` | number | Yes | Template type: 1 for custom templates, 2 for built-in templates. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |

### list_namespaces

**Description**: List CodeUp namespaces or groups in a Yunxiao organization.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `parentId` | number | No | Parent namespace ID. Omit to list namespaces available to the current user. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |
| `search` | string | No | Namespace search keyword. |
| `orderBy` | string | No | Sort field: created_at or updated_at. |
| `sort` | string | No | Sort direction: asc or desc. |

### list_tags

**Description**: List tags in a CodeUp repository.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |
| `search` | string | No | Tag search keyword. |
| `sort` | string | No | Sort direction: asc or desc. |
| `orderBy` | string | No | Sort field: name or create. |

### list_repository_members

**Description**: List members of a CodeUp repository.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `accessLevel` | number | No | Minimum access level: 20, 30, or 40. |

### list_protected_branches

**Description**: List protected branch rules in a CodeUp repository.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |

### list_push_rules

**Description**: List push rules in a CodeUp repository.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |

### list_repositories

**Description**: List CodeUp repositories in a Yunxiao organization.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |
| `orderBy` | string | No | Sort field: created_at, name, path, or last_activity_at. |
| `sort` | string | No | Sort direction: asc or desc. |
| `search` | string | No | Fuzzy repository path search keyword. |
| `archived` | boolean | No | Filter archived repositories. |

### list_branches

**Description**: List branches in a CodeUp repository.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |
| `sort` | string | No | Sort mode: name_asc, name_desc, updated_asc, or updated_desc. |
| `search` | string | No | Branch search keyword. |

### list_files

**Description**: List files in a CodeUp repository tree.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `path` | string | No | Directory path to query. |
| `ref` | string | No | Branch, tag, or commit SHA. Defaults to the repository default branch when omitted. |
| `type` | string | No | Tree mode: DIRECT, RECURSIVE, or FLATTEN. |

### list_commits

**Description**: List commits in a CodeUp repository.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `refName` | string | Yes | Branch, tag, or commit SHA. |
| `since` | string | No | Start time in ISO 8601 format. |
| `until` | string | No | End time in ISO 8601 format. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |
| `path` | string | No | Filter commits touching this path. |
| `search` | string | No | Commit search keyword. |
| `showSignature` | boolean | No | Whether to include commit signatures. |
| `committerIds` | string | No | Comma-separated committer user IDs. |

### list_change_requests

**Description**: List CodeUp merge requests in a Yunxiao organization.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |
| `projectIds` | string | No | Comma-separated repository IDs or full paths. |
| `authorIds` | string | No | Comma-separated author user IDs. |
| `reviewerIds` | string | No | Comma-separated reviewer user IDs. |
| `state` | string | No | Merge request state: opened, merged, or closed. |
| `search` | string | No | Title search keyword. |
| `orderBy` | string | No | Sort field: created_at or updated_at. |
| `sort` | string | No | Sort direction: asc or desc. |
| `createdBefore` | string | No | Created-before time in ISO 8601 format. |
| `createdAfter` | string | No | Created-after time in ISO 8601 format. |

### list_change_request_patch_sets

**Description**: List CodeUp merge request patch sets.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `localId` | string | Yes | Merge request local ID within the repository. |

### list_change_request_comments

**Description**: List CodeUp merge request comments.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `localId` | string | Yes | Merge request local ID within the repository. |
| `patchSetBizIds` | string | No | Comma-separated patch set IDs. |
| `commentType` | string | No | Comment type: GLOBAL_COMMENT or INLINE_COMMENT. Defaults to GLOBAL_COMMENT. |
| `state` | string | No | Comment state: OPENED or DRAFT. Defaults to OPENED. |
| `resolved` | boolean | No | Whether to list resolved comments. Defaults to false. |
| `filePath` | string | No | File path filter for inline comments. |

