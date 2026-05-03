# Codeup Tools

This document describes the 37 read-only MCP tools in the codeup domain.

## Tool Inventory

| Tool | Description |
|------|-------------|
| `list_ssh_keys` | List CodeUp SSH keys in a Yunxiao organization. |
| `get_ssh_key` | Get a CodeUp SSH key by ID. |
| `list_user_ssh_keys` | List CodeUp SSH keys for a Yunxiao user. |
| `list_webhooks` | List CodeUp webhooks in a repository. |
| `get_webhook` | Get a CodeUp webhook by ID. |
| `list_commit_statuses` | List CodeUp commit statuses for a repository commit. |
| `list_check_runs` | List CodeUp check runs for a repository ref. |
| `get_check_run` | Get a CodeUp check run by ID. |
| `list_group_members` | List CodeUp group members. |
| `get_member_https_clone_username` | Get a CodeUp user's HTTPS clone username. |
| `list_merge_requests` | List legacy CodeUp merge requests in a Yunxiao organization. |
| `get_merge_request` | Get legacy CodeUp merge request details. |
| `list_template_repositories` | List CodeUp template repositories in a Yunxiao organization. |
| `list_namespaces` | List CodeUp namespaces or groups in a Yunxiao organization. |
| `get_namespace` | Get a CodeUp namespace or group by ID or full path. |
| `get_org_namespace` | Get the organization-level CodeUp namespace. |
| `list_tags` | List tags in a CodeUp repository. |
| `list_repository_members` | List members of a CodeUp repository. |
| `list_protected_branches` | List protected branch rules in a CodeUp repository. |
| `get_protected_branch` | Get a protected branch rule in a CodeUp repository. |
| `list_push_rules` | List push rules in a CodeUp repository. |
| `get_push_rule` | Get a push rule in a CodeUp repository. |
| `list_repositories` | List CodeUp repositories in a Yunxiao organization. |
| `get_repository` | Get a CodeUp repository by numeric ID or full path. |
| `list_branches` | List branches in a CodeUp repository. |
| `get_branch` | Get CodeUp branch details. |
| `list_files` | List files in a CodeUp repository tree. |
| `get_file_blobs` | Get CodeUp file content. |
| `list_commits` | List commits in a CodeUp repository. |
| `get_commit` | Get CodeUp commit details. |
| `compare` | Compare two CodeUp refs or commits. |
| `list_change_requests` | List CodeUp merge requests in a Yunxiao organization. |
| `get_change_request` | Get CodeUp merge request details. |
| `list_change_request_patch_sets` | List CodeUp merge request patch sets. |
| `get_change_request_tree` | Get CodeUp merge request changed file tree. |
| `list_change_request_comments` | List CodeUp merge request comments. |
| `get_change_request_comment` | Get CodeUp merge request comment details. |

### list_ssh_keys

**Description**: List CodeUp SSH keys in a Yunxiao organization.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |
| `orderBy` | string | No | Sort field: created_at or updated_at. |
| `sort` | string | No | Sort direction: asc or desc. |

### get_ssh_key

**Description**: Get a CodeUp SSH key by ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `keyId` | string | Yes | SSH key ID. |

### list_user_ssh_keys

**Description**: List CodeUp SSH keys for a Yunxiao user.

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

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |

### get_webhook

**Description**: Get a CodeUp webhook by ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `hookId` | string | Yes | Webhook ID. |

### list_commit_statuses

**Description**: List CodeUp commit statuses for a repository commit.

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

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `ref` | string | Yes | Commit SHA, branch name, or tag name. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |

### get_check_run

**Description**: Get a CodeUp check run by ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `checkRunId` | string | Yes | Check run ID. |

### list_group_members

**Description**: List CodeUp group members.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `groupId` | string | Yes | Group ID or URL-encoded full path. |
| `accessLevel` | number | No | Minimum access level: 20 viewer, 30 developer, 40 admin. Defaults to no filter. |

### get_member_https_clone_username

**Description**: Get a CodeUp user's HTTPS clone username.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `userId` | string | Yes | Yunxiao user ID. |

### list_merge_requests

**Description**: List legacy CodeUp merge requests in a Yunxiao organization.

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

### get_merge_request

**Description**: Get legacy CodeUp merge request details.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `iid` | string | Yes | Legacy merge request IID within the repository. |

### list_template_repositories

**Description**: List CodeUp template repositories in a Yunxiao organization.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `templateType` | number | Yes | Template type: 1 for custom templates, 2 for built-in templates. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |

### list_namespaces

**Description**: List CodeUp namespaces or groups in a Yunxiao organization.

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

### get_namespace

**Description**: Get a CodeUp namespace or group by ID or full path.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `namespaceId` | string | Yes | Namespace ID or full path such as group/subgroup. |

### get_org_namespace

**Description**: Get the organization-level CodeUp namespace.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |

### list_tags

**Description**: List tags in a CodeUp repository.

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

### get_protected_branch

**Description**: Get a protected branch rule in a CodeUp repository.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `protectedBranchRuleId` | string | Yes | Protected branch rule ID. |

### list_push_rules

**Description**: List push rules in a CodeUp repository.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |

### get_push_rule

**Description**: Get a push rule in a CodeUp repository.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `pushRuleId` | string | Yes | Push rule ID. |

### list_repositories

**Description**: List CodeUp repositories in a Yunxiao organization.

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

### get_repository

**Description**: Get a CodeUp repository by numeric ID or full path.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |

### list_branches

**Description**: List branches in a CodeUp repository.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |
| `sort` | string | No | Sort mode: name_asc, name_desc, updated_asc, or updated_desc. |
| `search` | string | No | Branch search keyword. |

### get_branch

**Description**: Get CodeUp branch details.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `branchName` | string | Yes | Branch name, such as main or feature/demo. |

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

### get_file_blobs

**Description**: Get CodeUp file content.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `filePath` | string | Yes | File path, such as src/main.go. |
| `ref` | string | Yes | Branch, tag, or commit SHA. |

### list_commits

**Description**: List commits in a CodeUp repository.

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

### get_commit

**Description**: Get CodeUp commit details.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `sha` | string | Yes | Commit SHA. |

### compare

**Description**: Compare two CodeUp refs or commits.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `from` | string | Yes | Source commit SHA, branch, or tag. |
| `to` | string | Yes | Target commit SHA, branch, or tag. |
| `sourceType` | string | No | Source ref type: branch or tag. |
| `targetType` | string | No | Target ref type: branch or tag. |
| `straight` | string | No | Whether to compare directly without merge-base: true or false. |

### list_change_requests

**Description**: List CodeUp merge requests in a Yunxiao organization.

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

### get_change_request

**Description**: Get CodeUp merge request details.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `localId` | string | Yes | Merge request local ID within the repository. |

### list_change_request_patch_sets

**Description**: List CodeUp merge request patch sets.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `localId` | string | Yes | Merge request local ID within the repository. |

### get_change_request_tree

**Description**: Get CodeUp merge request changed file tree.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `localId` | string | Yes | Merge request local ID within the repository. |
| `fromPatchSetId` | string | Yes | Target-side patch set ID. |
| `toPatchSetId` | string | Yes | Source-side patch set ID. |

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

### get_change_request_comment

**Description**: Get CodeUp merge request comment details.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repositoryId` | string | Yes | Repository numeric ID or full path such as org/repo. |
| `localId` | string | Yes | Merge request local ID within the repository. |
| `commentBizId` | string | Yes | Comment business ID. |

