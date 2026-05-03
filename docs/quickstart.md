# Quickstart Guide

This guide provides common MCP conversation patterns for AI assistants using the Yunxiao toolset. Each pattern shows the recommended sequence of tool calls for a typical user question.

## Contents

- [Understand a Project](#pattern-understand-a-project)
- [Track Sprint Progress](#pattern-track-sprint-progress)
- [Find and Inspect a Work Item](#pattern-find-and-inspect-a-work-item)
- [My Tasks in a Project](#pattern-my-tasks-in-a-project)
- [Check Work Item Type Configuration](#pattern-check-work-item-type-configuration)
- [Member Workload](#pattern-member-workload)
- [Browse Code Repositories](#pattern-browse-code-repositories)
- [Review a Change Request](#pattern-review-a-change-request)
- [Check Branch Status](#pattern-check-branch-status)
- [Check Commit Status](#pattern-check-commit-status)
- [Check Pipeline Status](#pattern-check-pipeline-status)
- [Check Pipeline Run](#pattern-check-pipeline-run)
- [Check Application Status](#pattern-check-application-status)
- [Organization Overview](#pattern-organization-overview)
- [Department Overview](#pattern-department-overview)
- [Group Overview](#pattern-group-overview)
- [Check Release Status](#pattern-check-release-status)
- [Check Recent Deployments](#pattern-check-recent-deployments)
- [Check Test Case Status](#pattern-check-test-case-status)
- [Browse Artifacts](#pattern-browse-artifacts)
- [Browse Knowledge Bases](#pattern-browse-knowledge-bases)
- [Audit and Compliance](#pattern-audit-and-compliance)

## Pattern: Understand a Project

**User asks**: "Tell me about project X."

**Recommended flow**:
1. `get_project_overview` — compact dashboard with members, sprints, milestones, versions, labels.
2. If the user wants deeper work item shape: `get_project_workitem_summary` — totals and samples by category.
3. If the user wants risk visibility: `get_project_risk_dashboard` — overdue, high-priority, and stale items.

**Example**:
```json
{
  "projectId": "project-1",
  "includeVersions": false,
  "perPage": 10
}
```

## Pattern: Track Sprint Progress

**User asks**: "How is sprint Y progressing?"

**Recommended flow**:
1. `get_sprint_overview` — sprint metadata plus work items grouped by category.
2. If the user wants a Kanban view: `get_project_workitem_board` with `sprint` filter.

**Example**:
```json
{
  "projectId": "project-1",
  "sprintId": "sprint-2",
  "categories": "Task,Bug"
}
```

## Pattern: Find and Inspect a Work Item

**User asks**: "Show me details for the bug about login failure."

**Recommended flow**:
1. `search_workitems` with `subject=login` and `category=Bug` to find the work item ID.
2. `get_project_workitem_detail` with the `workitemId` to get full context (activities, comments, attachments, relations).

**Example step 1**:
```json
{
  "projectId": "project-1",
  "category": "Bug",
  "subject": "login"
}
```

**Example step 2**:
```json
{
  "workitemId": "wi-123"
}
```

## Pattern: My Tasks in a Project

**User asks**: "What tasks do I have in project X?"

**Recommended flow**:
1. `get_current_user` to obtain the user's ID.
2. `get_my_project_workitems` with `relation=assigned`.

**Example step 2**:
```json
{
  "projectId": "project-1",
  "userId": "user-123",
  "relation": "assigned",
  "categories": "Task,Bug"
}
```

## Pattern: Check Work Item Type Configuration

**User asks**: "What fields does the Task type have in project X?"

**Recommended flow**:
1. `get_work_item_type_overview` — work item type info plus field configuration and workflow in one call.

**Example**:
```json
{
  "projectId": "project-1",
  "workItemTypeId": "type-1",
  "includeFieldConfig": true,
  "includeWorkflow": true
}
```

## Pattern: Member Workload

**User asks**: "Who is overloaded in project X?"

**Recommended flow**:
1. `get_project_member_task_status` — per-member task counts and overdue items.
2. If custom status groups are needed, pass `statusGroups` as a JSON object.

**Example**:
```json
{
  "projectId": "project-1",
  "memberLimit": 20,
  "sampleLimit": 5
}
```

## Pattern: Browse Code Repositories

**User asks**: "What repositories are in namespace Y?"

**Recommended flow**:
1. `list_repositories` with the namespace ID.

**User asks**: "Tell me about repository Z."

**Recommended flow**:
1. `get_repository_overview` — repository info, branches, recent commits, and open merge requests in one call.
2. If deeper commit history is needed: `list_commits` with a larger `perPage`.

## Pattern: Review a Change Request

**User asks**: "Show me the change request about feature X."

**Recommended flow**:
1. `list_change_requests` with `subject=feature-x` to find the change request local ID.
2. `get_change_request_overview` with the `localId` to get full context (patch sets and comments).

**Example step 1**:
```json
{
  "repositoryId": "org/repo",
  "state": "opened"
}
```

**Example step 2**:
```json
{
  "repositoryId": "org/repo",
  "localId": "1"
}
```

## Pattern: Check Branch Status

**User asks**: "What's happening on the feature branch?"

**Recommended flow**:
1. `get_branch_overview` — branch details, recent commits, and merge requests targeting the branch in one call.

**Example**:
```json
{
  "repositoryId": "org/repo",
  "branchName": "feature/demo",
  "includeCommits": true,
  "includeMergeRequests": true,
  "commitLimit": 10,
  "mrLimit": 5
}
```

## Pattern: Check Commit Status

**User asks**: "What is the CI status of commit abc123?"

**Recommended flow**:
1. `get_commit_overview` — commit details, commit statuses, and check runs in one call.

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

## Pattern: Check Pipeline Status

**User asks**: "How is pipeline X doing?"

**Recommended flow**:
1. `get_pipeline_overview` — pipeline info, latest run, and recent run history in one call.

**Example**:
```json
{
  "pipelineId": "pipeline-1",
  "runLimit": 5
}
```

## Pattern: Check Pipeline Run

**User asks**: "How did pipeline run #123 go?"

**Recommended flow**:
1. `get_pipeline_run_overview` — run details plus pipeline jobs by category in one call.

**Example**:
```json
{
  "pipelineId": "pipeline-1",
  "pipelineRunId": "run-1",
  "includeJobs": true,
  "category": "DEPLOY"
}
```

## Pattern: Check Application Status

**User asks**: "Tell me about application Z."

**Recommended flow**:
1. `get_application_overview` — application info, environments, and recent orchestrations in one call.
2. If specific environment details are needed: `get_environment_overview` with the app and environment name.

**Example step 1**:
```json
{
  "appName": "my-app",
  "envLimit": 5,
  "orchestrationLimit": 5
}
```

**Example step 2**:
```json
{
  "appName": "my-app",
  "envName": "dev"
}
```

## Pattern: Organization Overview

**User asks**: "Tell me about my organization."

**Recommended flow**:
1. `get_organization_overview` — organization info, departments, members, groups, and roles in one call.

**Example**:
```json
{
  "departmentLimit": 5,
  "memberLimit": 5,
  "groupLimit": 5
}
```

## Pattern: Department Overview

**User asks**: "Tell me about department X in my organization."

**Recommended flow**:
1. `get_organization_department_overview` — department info and ancestor chain in one call.

**Example**:
```json
{
  "departmentId": "dept-1",
  "includeAncestors": true
}
```

## Pattern: Group Overview

**User asks**: "Who is in group X?"

**Recommended flow**:
1. `get_organization_group_overview` — group info and members in one call.

**Example**:
```json
{
  "groupId": "group-1",
  "includeMembers": true,
  "memberLimit": 10
}
```

## Pattern: Check Release Status

**User asks**: "Tell me about release X in system Y."

**Recommended flow**:
1. `get_release_overview` — release info, members, products, and attached change requests in one call.

**Example**:
```json
{
  "systemName": "my-system",
  "sn": "rel-1",
  "changeRequestLimit": 5
}
```

## Pattern: Check Recent Deployments

**User asks**: "What was the last deployment for application Z?"

**Recommended flow**:
1. `list_change_order_versions` with the application name, sorted by creation time descending.

**Tip**: Appstack change orders represent deployment records. Use `current` and `pageSize` for pagination on this endpoint.

## Pattern: Check Test Case Status

**User asks**: "Show me test cases for repository X."

**Recommended flow**:
1. `search_testcases` with `testRepoId` and optional `subject` filter.

**Example step 1**:
```json
{
  "testRepoId": "repo-1",
  "subject": "login",
  "page": 1,
  "perPage": 10
}
```

**Example step 2**:
```json
{
  "testRepoId": "repo-1",
  "testcaseId": "tc-123"
}
```

## Pattern: Browse Artifacts

**User asks**: "What artifacts are in package repository Y?"

**Recommended flow**:
1. `list_artifacts` with the package repository ID.

**Example step 1**:
```json
{
  "packageRepositoryId": "repo-1",
  "page": 1,
  "perPage": 10
}
```

## Pattern: Browse Knowledge Bases

**User asks**: "Show me files in knowledge base Z."

**Recommended flow**:
1. `list_knowledge_bases` to find the knowledge base ID.
2. `list_kb_files` with the knowledge base ID.

**Example step 2**:
```json
{
  "knowledgeBaseId": "kb-1",
  "page": 1,
  "perPage": 10
}
```

## Pattern: Audit and Compliance

**User asks**: "Show me recent audit logs."

**Recommended flow**:
1. `list_audit_logs` — organization-scoped audit events.
2. Use `nextToken` for pagination; this endpoint also uses keyset pagination.

## Tips for AI Assistants

- **Prefer enhanced tools**: When both a base tool and an enhanced tool exist for the same concept, use the enhanced tool. Enhanced tools return aggregated, filtered responses that are easier to summarize.
- **organizationId is optional**: If the user does not provide an organization ID, the server auto-injects the user's sole organization. Only ask for it when the user explicitly mentions multiple organizations.
- **Pagination defaults**: Always specify `perPage` or `pageSize` for predictable behavior. Default sizes vary by endpoint (usually 10 or 20).
- **Conditions JSON**: For advanced queries in `search_workitems` or `search_projects`, see `conditions-cookbook.md`. In most cases, simple filters (`subject`, `status`, `assignedTo`) are sufficient.
- **Status IDs**: Work item status values are IDs, not names. Use `get_project_workitem_context` or `get_work_item_workflow` to discover valid status IDs for a project.
