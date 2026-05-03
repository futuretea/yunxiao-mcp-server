# Codeup Enhanced Tools

This document describes the enhanced aggregation tools in the Codeup (source code hosting) module. These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations.

## Tool Inventory

| Tool | Purpose | API Calls |
|------|---------|-----------|
| `get_repository_overview` | Repository dashboard with branches, recent commits, and merge requests | 1 + up to 3 |
| `get_change_request_overview` | Change request snapshot with patch sets and comments | 1 + up to 2 |
| `get_commit_overview` | Commit snapshot with details, commit statuses, and check runs | 1 + up to 2 |

## Common Behaviors

### Pagination

Sub-sections use `page`/`perPage` parameters with a default page size of 5 (controlled by `branchLimit`, `commitLimit`, and `mrLimit`). The raw upstream response is returned for each section; pagination metadata varies by endpoint.

### Default Branch Detection

When `refName` is omitted and `includeCommits` is true, the tool attempts to read `defaultBranch` from the repository response. If the repository response does not contain a `defaultBranch` field, the commits section is skipped.

## Tool Details

### get_repository_overview

**When to use**: You want a quick snapshot of a repository — its basic info plus branches, recent commits, and open merge requests.

**Parameters**:
- `organizationId`, `repositoryId`: required
- `includeBranches`, `includeCommits`, `includeMergeRequests`: toggle sections, default true
- `refName`: optional branch/tag/SHA for commits; defaults to repository default branch
- `branchLimit`, `commitLimit`, `mrLimit`: control section sizes, default 5
- `mrState`: merge request state filter, default `opened`

**Example**:
```json
{
  "repositoryId": "org/my-repo",
  "includeBranches": true,
  "includeCommits": true,
  "includeMergeRequests": true,
  "branchLimit": 5,
  "commitLimit": 10,
  "mrLimit": 5,
  "mrState": "opened"
}
```

### get_change_request_overview

**When to use**: You want a quick snapshot of a change request (merge request) — its basic info plus patch sets and comments.

**Parameters**:
- `organizationId`, `repositoryId`, `localId`: required
- `includePatchSets`: toggle patch sets list, default true
- `includeComments`: toggle comments list, default true
- `commentState`: comment state filter, default `OPENED`
- `commentResolved`: whether to include resolved comments, default false

**Example**:
```json
{
  "repositoryId": "org/repo",
  "localId": "1",
  "includePatchSets": true,
  "includeComments": true
}
```

**Note**: Comments are fetched via POST with `comment_type` set to `GLOBAL_COMMENT`. Merge requests are filtered to the specified repository using the `repositoryIds` query parameter.

### get_commit_overview

**When to use**: You want a quick snapshot of a commit — its details plus CI statuses and check runs.

**Parameters**:
- `organizationId`, `repositoryId`, `sha`: required
- `includeStatuses`: toggle commit statuses list, default true
- `includeCheckRuns`: toggle check runs list, default true
- `statusLimit`: max commit statuses returned, default 5
- `checkRunLimit`: max check runs returned, default 5

**Example**:
```json
{
  "repositoryId": "org/repo",
  "sha": "abc123def456",
  "includeStatuses": true,
  "includeCheckRuns": true,
  "statusLimit": 5,
  "checkRunLimit": 5
}
```

**Note**: Check runs are fetched using the commit SHA as the `ref` query parameter.
