# Quickstart Guide

This guide provides common MCP conversation patterns for AI assistants using the Yunxiao toolset. Each pattern shows the recommended sequence of tool calls for a typical user question.

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
1. `list_codeup_repositories` with the namespace ID.
2. `get_repository` for details on a specific repo.
3. `list_merge_requests` or `list_commits` for recent activity.

## Pattern: Check Recent Deployments

**User asks**: "What was the last deployment for application Z?"

**Recommended flow**:
1. `list_appstack_deployments` with the application name or ID, sorted by `gmtCreate` descending.
2. `get_appstack_deployment` for the specific deployment ID if needed.

**Tip**: Appstack tools use keyset pagination (`nextToken`). Omit `nextToken` on the first request; subsequent requests use the token from the previous response.

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
