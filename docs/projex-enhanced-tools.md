# Projex Enhanced Tools

This document describes the enhanced aggregation tools in the Projex (project management) module. These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations.

## Tool Inventory

| Tool | Purpose | API Calls |
|------|---------|-----------|
| `get_project_overview` | Compact project dashboard with members, sprints, milestones, versions, labels | 1 + up to 5 |
| `get_project_workitem_summary` | Work item samples and totals by category | 1 per category |
| `get_project_workitem_context` | Metadata context for a work item category | 1 + up to 4 |
| `get_sprint_overview` | Sprint-centric view of work items by category | 1 + 1 per category |
| `get_my_project_workitems` | Personal work item dashboard (assigned or created) | 1 per category |
| `get_project_risk_dashboard` | Risk-focused view with overdue, high-priority, stale items | 3 + statusGroups |
| `get_project_member_task_status` | Per-member workload and overdue tracking | 1 + 2N + statusGroups |

## Common Behaviors

### Pagination

List sections use `page`/`perPage` parameters. The response includes a `pagination` object when the upstream API returns `x-total` headers:

```json
{
  "data": [...],
  "pagination": {
    "total": 100,
    "page": 1,
    "perPage": 20
  }
}
```

Some endpoints (e.g., `list_project_members`) do not return pagination metadata; in these cases the section is returned as a plain array.

### Category Search

Tools that search by category (`get_project_workitem_summary`, `get_sprint_overview`, `get_my_project_workitems`) default to `Task,Bug` and issue one `workitems:search` request per category. The `sampleLimit` parameter (default 5, clamped 0-200) controls `perPage` for each category search.

### Conditions

Simple filters (`status`, `assignedTo`, `sprint`, etc.) are translated into Yunxiao `conditions` JSON. Advanced users can override with a raw `conditions` JSON string.

## Tool Details

### get_project_overview

**When to use**: You need a quick snapshot of a project — its basic info plus active sprints, milestones, versions, members, and labels.

**Parameters**:
- `organizationId`, `id` (project ID): required
- `includeMembers`, `includeSprints`, `includeMilestones`, `includeVersions`, `includeLabels`: toggle sections, default true
- `activeOnly`: when true, sprints/milestones/versions use `status=TODO,DOING`
- `page`, `perPage`: control list sections

**Example**:
```json
{
  "organizationId": "org-1",
  "id": "project-1",
  "includeVersions": false,
  "page": 1,
  "perPage": 10
}
```

### get_project_workitem_summary

**When to use**: You want to see the shape of work in a project by category — how many Tasks, Bugs, Requirements exist, with a few samples.

**Parameters**:
- `categories`: defaults to `Req,Task,Bug,Risk`
- `subject`, `status`, `assignedTo`, `creator`, `tag`: simple filters applied to every category
- `sampleLimit`: samples per category

### get_project_workitem_context

**When to use**: You are about to create or edit a work item and need to know the valid types, labels, members, and optionally the field config and workflow for a specific type.

**Parameters**:
- `category`: required (e.g., `Task`, `Bug`)
- `workItemTypeId`: optional, triggers `fields` and `workflow` sections

### get_sprint_overview

**When to use**: You want to understand what work is inside a specific sprint, grouped by category.

**Parameters**:
- `sprintId`: required
- `categories`: defaults to `Task,Bug`
- `status`: applied to every category search

**Note**: The tool first fetches sprint metadata, then searches workitems filtered by that sprint ID.

### get_my_project_workitems

**When to use**: You want to see work items that belong to a specific user in a project — either items assigned to them or items they created.

**Parameters**:
- `userId`: required
- `relation`: `assigned` (default) or `created`
- `categories`: defaults to `Task,Bug`

**Tip**: In an MCP conversation, first call `get_current_user` to obtain the user's ID, then call this tool.

### get_project_risk_dashboard

**When to use**: You need to identify at-risk work items — overdue, high-priority, or stale.

**Parameters**:
- `categories`: defaults to `Risk,Bug,Task`
- `overdueBefore`: defaults to today
- `highPriority`: priority IDs for the high-priority section (optional)
- `staleBefore`: status-update date upper bound for the stale section (optional)

### get_project_member_task_status

**When to use**: You want to see the workload distribution across project members, including overdue items and custom status groups.

**Parameters**:
- `assigneeIds`: comma-separated user IDs; defaults to project members up to `memberLimit`
- `statusGroups`: JSON object mapping group names to comma-separated status IDs (optional)
- `memberLimit`: max members when `assigneeIds` is omitted, default 20

**Note**: This tool issues 2 + (2 * statusGroups) requests per member. Use `memberLimit` and `sampleLimit` to control API load.
