# Projex Tools

This document describes the 45 MCP tools in the projex domain.

Access summary: 41 read-only, 4 write-capable.

## Enhanced Tools

These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations. Prefer them when available.

| Tool | Description |
|------|-------------|
| `get_project_overview` | Get a compact Projex project overview with members, sprints, milestones, versions, and labels in one read-only call. This is the best starting point when exploring a new project. |
| `get_project_workitem_summary` | Summarize Projex work items by category for one project with samples and pagination totals. This is ideal for dashboards and overviews. For detailed searching with complex filters, use search_workitems instead. |
| `get_project_workitem_context` | Get project work item metadata context: types, labels, members, and optional field/workflow configuration for one type. Use this before creating or updating work items to discover valid IDs for status, priority, labels, and members. |
| `get_sprint_overview` | Get a sprint overview with work item samples and totals by category for one project sprint. Combines sprint metadata with categorized work item summaries in a single call. |
| `get_my_project_workitems` | Get work items assigned to or created by a specific user for one project, grouped by category. Useful for workload review and stand-up summaries. |
| `get_project_workitem_board` | Get a Kanban-style board view of work items grouped by status for one project. Returns work items organized by their current status columns. For flat search results, use search_workitems instead. |
| `get_project_workitem_detail` | Get comprehensive details for a single work item including basic info, activities, attachments, comments, and relation records in one read-only call. Use this when you need the full picture of one specific work item. |
| `get_work_item_type_overview` | Get a comprehensive overview of a Projex work item type including basic info, field configuration, and workflow in one read-only call. Use this to understand the schema and lifecycle of a specific work item type. |

## Pagination

Tools in this domain use the following pagination scheme(s):

- Offset (page/perPage)

## Tool Inventory

Tools marked in **bold** are enhanced aggregation tools.

| Tool | Access | Description |
|------|--------|-------------|
| `list_current_user_effort_records` | Read-only | List actual effort (time tracking) records for the current user within a date range. |
| `list_effort_records` | Read-only | List actual effort (time tracking) records logged against a specific work item. |
| `list_estimated_efforts` | Read-only | List estimated effort records for a specific work item. |
| **`get_project_overview`** | Read-only | Get a compact Projex project overview with members, sprints, milestones, versions, and labels in one read-only call. This is the best starting point when exploring a new project. |
| **`get_project_workitem_summary`** | Read-only | Summarize Projex work items by category for one project with samples and pagination totals. This is ideal for dashboards and overviews. For detailed searching with complex filters, use search_workitems instead. |
| **`get_project_workitem_context`** | Read-only | Get project work item metadata context: types, labels, members, and optional field/workflow configuration for one type. Use this before creating or updating work items to discover valid IDs for status, priority, labels, and members. |
| **`get_sprint_overview`** | Read-only | Get a sprint overview with work item samples and totals by category for one project sprint. Combines sprint metadata with categorized work item summaries in a single call. |
| **`get_my_project_workitems`** | Read-only | Get work items assigned to or created by a specific user for one project, grouped by category. Useful for workload review and stand-up summaries. |
| **`get_project_workitem_board`** | Read-only | Get a Kanban-style board view of work items grouped by status for one project. Returns work items organized by their current status columns. For flat search results, use search_workitems instead. |
| **`get_project_workitem_detail`** | Read-only | Get comprehensive details for a single work item including basic info, activities, attachments, comments, and relation records in one read-only call. Use this when you need the full picture of one specific work item. |
| **`get_work_item_type_overview`** | Read-only | Get a comprehensive overview of a Projex work item type including basic info, field configuration, and workflow in one read-only call. Use this to understand the schema and lifecycle of a specific work item type. |
| `get_project_risk_dashboard` | Read-only | Get a read-only project risk dashboard with category samples, overdue work items, and optional high-priority/stale sections. Best used for project health checks and sprint retrospectives. |
| `get_project_member_task_status` | Read-only | Get per-member task status for one project with assigned, overdue, and optional status-group sections. Useful for workload balancing and identifying blocked team members. |
| `get_sprint_velocity` | Read-only | Get historical sprint velocity metrics including completion rates, work item counts, and trend analysis per sprint. Combines sprint metadata with categorized work item completion statistics. Useful for sprint planning and delivery forecasting. |
| `get_workitem_status_timeline` | Read-only | Get detailed status change timeline for a work item including timestamps, duration in each status, and responsible users. Parses activity history to extract workflow bottlenecks. Useful for analyzing why a task is delayed. |
| `get_blocker_analysis` | Read-only | Get aggregated dependency blocker analysis for a project including work items that are blocked by unresolved dependencies and work items that are blocking others. Useful for identifying critical path risks and unblocking the team. |
| `get_member_workload_trend` | Read-only | Get workload trend analysis for project members including current task counts, status distribution, recent activity, and overdue items. Useful for capacity planning and identifying overloaded team members. |
| `get_team_workload_breakdown` | Read-only | Get a detailed per-member workload breakdown with current active tasks, task subjects, status, labels, and creation dates. Useful for stand-ups, sprint planning, and identifying what each team member is working on right now. |
| `list_milestones` | Read-only | List milestones (planning checkpoints) in a Projex project. Milestones track progress against goals, distinct from sprints (time-boxed iterations) and versions (releases). |
| `list_testcase_repositories` | Read-only | List Projex testcase repositories in a Yunxiao organization. A testcase repository is a container for organizing test cases. |
| `list_directories` | Read-only | List directories (folders) within a Projex testcase repository. |
| `search_testcases` | Read-only | Search test cases within a single Projex testcase repository. |
| `list_test_plans` | Read-only | List test plans in a Yunxiao organization. A test plan groups test cases for execution tracking. |
| `get_test_result_list` | Read-only | Get test execution result summaries for a specific directory within a test plan. |
| `list_project_members` | Read-only | List members in a Projex project. Use this to discover user IDs for filtering work items or assigning tasks. |
| `list_project_templates` | Read-only | List Projex project templates in a Yunxiao organization. Useful when setting up new projects. |
| `list_project_program` | Read-only | List Projex projects bound to a project program (project group). |
| `list_project_roles` | Read-only | List roles defined in a specific Projex project. |
| `list_all_project_roles` | Read-only | List all Projex project roles across a Yunxiao organization. |
| `search_projects` | Read-only | Search Projex projects in a Yunxiao organization. Use this to discover available projects and obtain their project IDs before calling other project-scoped tools. |
| `list_sprints` | Read-only | List sprints (iterations) in a Projex project. Sprints are time-boxed development cycles. Use this when you need iteration-level planning data. For a broader view that also includes milestones and versions, use get_project_overview instead. |
| `search_workitems` | Read-only | Search work items in a single Projex project with flexible filtering and pagination. This is the primary tool for finding work items when you know the project. For a high-level summary across categories without detailed filtering, use get_project_workitem_summary. For a Kanban board view grouped by status, use get_project_workitem_board. |
| `list_work_item_comments` | Read-only | List comments for a single Projex work item. For a comprehensive view that includes comments along with activities, attachments, and relations, use get_project_workitem_detail instead. |
| `list_all_work_item_types` | Read-only | List all work item types across a Yunxiao organization. This returns organization-level type definitions. For project-scoped types (including custom fields), use list_work_item_types instead. |
| `list_work_item_types` | Read-only | List work item types available in a specific Projex project. Use this before creating a work item to discover valid workitemTypeId values and categories for the target project. |
| `list_work_item_relation_work_item_types` | Read-only | List work item types that can be related to a given work item type. Use this when setting up work item relationships to discover which types are valid for the chosen relationType. |
| `list_versions` | Read-only | List versions (releases) in a Projex project. Versions represent release milestones, distinct from sprints (iterations) and milestones (planning checkpoints). |
| `list_workitem_activities` | Read-only | List activity events (history log) for a single Projex work item. For a comprehensive view including comments and attachments, use get_project_workitem_detail instead. |
| `list_workitem_attachments` | Read-only | List file attachments for a Projex work item. Use search_workitems to discover work item IDs. |
| `list_workitem_relation_records` | Read-only | List relation records (parent, subtask, dependency links) for a Projex work item. Use search_workitems to discover work item IDs. |
| `list_labels` | Read-only | List labels in a Projex project. Labels are used to categorize and filter work items. Use search_projects to discover project IDs. |
| `create_workitem` | Write-capable | Create a new work item in a Projex project. Before calling this, use list_work_item_types to discover the correct category and workitemTypeId for the target project. |
| `update_workitem` | Write-capable | Update an existing work item's fields. Only provided fields will be changed; omitted fields retain their current values. Use search_workitems to find the workitemId if unknown. |
| `update_workitem_status` | Write-capable | Change the status of a work item. Use get_project_workitem_context to discover valid statusId values for the work item type before updating. |
| `add_workitem_comment` | Write-capable | Add a comment to a work item. Use search_workitems to find the workitemId if unknown. |

### list_current_user_effort_records

**Description**: List actual effort (time tracking) records for the current user within a date range.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `startDate` | string | Yes | Start date in YYYY-MM-DD format (e.g., 2024-01-01). |
| `endDate` | string | Yes | End date in YYYY-MM-DD format (e.g., 2024-01-31). |

### list_effort_records

**Description**: List actual effort (time tracking) records logged against a specific work item.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `workitemId` | string | Yes | Work item ID (numeric string). Find it via search_workitems. |

### list_estimated_efforts

**Description**: List estimated effort records for a specific work item.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `workitemId` | string | Yes | Work item ID (numeric string). Find it via search_workitems. |

### get_project_overview

**Description**: Get a compact Projex project overview with members, sprints, milestones, versions, and labels in one read-only call. This is the best starting point when exploring a new project.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `includeMembers` | boolean | No | Whether to include project members. Defaults to true. |
| `includeSprints` | boolean | No | Whether to include sprints. Defaults to true. |
| `includeMilestones` | boolean | No | Whether to include milestones. Defaults to true. |
| `includeVersions` | boolean | No | Whether to include versions. Defaults to true. |
| `includeLabels` | boolean | No | Whether to include labels. Defaults to true. |
| `activeOnly` | boolean | No | Whether sprint, milestone, and version sections use the status filter. Defaults to true. |
| `status` | string | No | Comma-separated statuses for sprint, milestone, and version sections when activeOnly is true. Defaults to TODO,DOING. |
| `page` | number | No | Page number for paged list sections. Defaults to 1. |
| `perPage` | number | No | Page size for paged list sections. Defaults to 20. |

### get_project_workitem_summary

**Description**: Summarize Projex work items by category for one project with samples and pagination totals. This is ideal for dashboards and overviews. For detailed searching with complex filters, use search_workitems instead.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `categories` | string | No | Comma-separated work item categories. Defaults to Req,Task,Bug,Risk. |
| `subject` | string | No | Subject contains filter applied to every category. |
| `status` | string | No | Comma-separated status IDs applied to every category. |
| `assignedTo` | string | No | Comma-separated assignee user IDs applied to every category. |
| `creator` | string | No | Comma-separated creator user IDs applied to every category. |
| `tag` | string | No | Comma-separated tag IDs applied to every category. |
| `conditions` | string | No | Advanced conditions JSON string applied to every category. Overrides simple filters. |
| `orderBy` | string | No | Sort field. Defaults are controlled by Yunxiao. |
| `sort` | string | No | Sort direction: asc or desc. |
| `sampleLimit` | number | No | Samples returned per category. Defaults to 5, clamped to 0-200. |

### get_project_workitem_context

**Description**: Get project work item metadata context: types, labels, members, and optional field/workflow configuration for one type. Use this before creating or updating work items to discover valid IDs for status, priority, labels, and members.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `category` | string | Yes | Work item category, such as Req, Task, Bug, or Risk. |
| `workItemTypeId` | string | No | Optional work item type ID (numeric string) for field and workflow metadata. Use list_work_item_types to find available types. |
| `includeMembers` | boolean | No | Whether to include project members. Defaults to true. |
| `includeLabels` | boolean | No | Whether to include project labels. Defaults to true. |
| `includeFields` | boolean | No | Whether to include field configuration when workItemTypeId is set. Defaults to true. |
| `includeWorkflow` | boolean | No | Whether to include workflow metadata when workItemTypeId is set. Defaults to true. |
| `page` | number | No | Page number for labels. Defaults to 1. |
| `perPage` | number | No | Page size for labels. Defaults to 20. |

### get_sprint_overview

**Description**: Get a sprint overview with work item samples and totals by category for one project sprint. Combines sprint metadata with categorized work item summaries in a single call.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `sprintId` | string | Yes | Sprint ID (string). Use list_sprints to discover sprint IDs in the project. |
| `categories` | string | No | Comma-separated work item categories. Defaults to Task,Bug. |
| `subject` | string | No | Subject contains filter applied to every category. |
| `status` | string | No | Comma-separated status IDs applied to every category. |
| `assignedTo` | string | No | Comma-separated assignee user IDs applied to every category. |
| `creator` | string | No | Comma-separated creator user IDs applied to every category. |
| `sampleLimit` | number | No | Samples returned per category. Defaults to 5, clamped to 0-200. |

### get_my_project_workitems

**Description**: Get work items assigned to or created by a specific user for one project, grouped by category. Useful for workload review and stand-up summaries.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `userId` | string | Yes | User ID (string) to filter work items by. Use list_project_members to find user IDs in the project. |
| `relation` | string | No | Filter relation: assigned or created. Defaults to assigned. |
| `categories` | string | No | Comma-separated work item categories. Defaults to Task,Bug. |
| `status` | string | No | Comma-separated status IDs applied to every category. |
| `sampleLimit` | number | No | Samples returned per category. Defaults to 5, clamped to 0-200. |

### get_project_workitem_board

**Description**: Get a Kanban-style board view of work items grouped by status for one project. Returns work items organized by their current status columns. For flat search results, use search_workitems instead.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `category` | string | Yes | Work item category, such as Task or Bug. |
| `sprint` | string | No | Optional sprint ID to filter work items. |
| `subject` | string | No | Subject contains filter applied to the search. |
| `status` | string | No | Comma-separated status IDs applied to the search. |
| `assignedTo` | string | No | Comma-separated assignee user IDs applied to the search. |
| `creator` | string | No | Comma-separated creator user IDs applied to the search. |
| `sampleLimit` | number | No | Max work items returned. Defaults to 5, clamped to 0-200. |

### get_project_workitem_detail

**Description**: Get comprehensive details for a single work item including basic info, activities, attachments, comments, and relation records in one read-only call. Use this when you need the full picture of one specific work item.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `workitemId` | string | Yes | Work item ID (numeric string). Find it via search_workitems or get_project_workitem_summary. |
| `includeActivities` | boolean | No | Whether to include activity history. Defaults to true. |
| `includeRelations` | boolean | No | Whether to include relation records. Defaults to true. |
| `relationTypes` | string | No | Comma-separated relation types for relation records. Defaults to ASSOCIATED,SUB. |
| `includeAttachments` | boolean | No | Whether to include attachments. Defaults to true. |
| `includeComments` | boolean | No | Whether to include comments. Defaults to true. |
| `page` | number | No | Page number for comments. Defaults to 1. |
| `perPage` | number | No | Page size for comments. Defaults to 20. |

### get_work_item_type_overview

**Description**: Get a comprehensive overview of a Projex work item type including basic info, field configuration, and workflow in one read-only call. Use this to understand the schema and lifecycle of a specific work item type.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `workItemTypeId` | string | Yes | Work item type ID (numeric string). Use list_work_item_types to find available types. |
| `includeFieldConfig` | boolean | No | Whether to include field configuration. Defaults to true. |
| `includeWorkflow` | boolean | No | Whether to include workflow metadata. Defaults to true. |

### get_project_risk_dashboard

**Description**: Get a read-only project risk dashboard with category samples, overdue work items, and optional high-priority/stale sections. Best used for project health checks and sprint retrospectives.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `categories` | string | No | Comma-separated categories for category totals. Defaults to Risk,Bug,Task. |
| `subject` | string | No | Subject contains filter applied to every section. |
| `status` | string | No | Comma-separated active status IDs applied to every section. |
| `statusStage` | string | No | Comma-separated status stage IDs applied to every section. |
| `assignedTo` | string | No | Comma-separated assignee user IDs applied to every section. |
| `creator` | string | No | Comma-separated creator user IDs applied to every section. |
| `sprint` | string | No | Comma-separated sprint IDs applied to every section. |
| `workitemType` | string | No | Comma-separated work item type IDs applied to every section. |
| `tag` | string | No | Comma-separated tag IDs applied to every section. |
| `overdueBefore` | string | No | Planned finish date upper bound for overdue work items. Defaults to today. |
| `highPriority` | string | No | Comma-separated priority IDs for the high-priority section. |
| `staleBefore` | string | No | Status update date upper bound for the stale section. |
| `sampleLimit` | number | No | Samples returned per section. Defaults to 5, clamped to 0-200. |

### get_project_member_task_status

**Description**: Get per-member task status for one project with assigned, overdue, and optional status-group sections. Useful for workload balancing and identifying blocked team members.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `assigneeIds` | string | No | Comma-separated assignee user IDs. Defaults to project members up to memberLimit. |
| `categories` | string | No | Comma-separated work item categories. Defaults to Task,Bug. |
| `subject` | string | No | Subject contains filter applied to every section. |
| `status` | string | No | Comma-separated status IDs applied to assigned and overdue sections. |
| `statusStage` | string | No | Comma-separated status stage IDs applied to assigned and overdue sections. |
| `assignedTo` | string | No | Comma-separated assignee user IDs applied to every section. |
| `creator` | string | No | Comma-separated creator user IDs applied to every section. |
| `sprint` | string | No | Comma-separated sprint IDs applied to every section. |
| `workitemType` | string | No | Comma-separated work item type IDs applied to every section. |
| `tag` | string | No | Comma-separated tag IDs applied to every section. |
| `overdueBefore` | string | No | Planned finish date upper bound for overdue work items. Defaults to today. |
| `statusGroups` | string | No | Optional JSON object mapping group names to comma-separated status IDs. |
| `memberLimit` | number | No | Max project members to inspect when assigneeIds is omitted. Defaults to 20. |
| `sampleLimit` | number | No | Samples returned per member section. Defaults to 5, clamped to 0-200. |

### get_sprint_velocity

**Description**: Get historical sprint velocity metrics including completion rates, work item counts, and trend analysis per sprint. Combines sprint metadata with categorized work item completion statistics. Useful for sprint planning and delivery forecasting.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `categories` | string | No | Comma-separated work item categories for velocity calculation. Defaults to Task,Bug. |
| `sprintCount` | number | No | Number of recent sprints to analyze. Defaults to 5, max 20. |
| `sprintStatus` | string | No | Comma-separated sprint statuses to include. Defaults to ARCHIVED,DONE. |

### get_workitem_status_timeline

**Description**: Get detailed status change timeline for a work item including timestamps, duration in each status, and responsible users. Parses activity history to extract workflow bottlenecks. Useful for analyzing why a task is delayed.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `workitemId` | string | Yes | Work item ID (numeric string). Find it via search_workitems. |
| `includeWorkitem` | boolean | No | Whether to include basic work item info. Defaults to true. |

### get_blocker_analysis

**Description**: Get aggregated dependency blocker analysis for a project including work items that are blocked by unresolved dependencies and work items that are blocking others. Useful for identifying critical path risks and unblocking the team.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `categories` | string | No | Comma-separated work item categories to analyze. Defaults to Task,Bug. |
| `sampleLimit` | number | No | Max work items to analyze per category. Defaults to 20, clamped to 0-200. |

### get_member_workload_trend

**Description**: Get workload trend analysis for project members including current task counts, status distribution, recent activity, and overdue items. Useful for capacity planning and identifying overloaded team members.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `assigneeIds` | string | No | Comma-separated assignee user IDs. Defaults to project members up to memberLimit. |
| `categories` | string | No | Comma-separated work item categories. Defaults to Task,Bug. |
| `memberLimit` | number | No | Max project members to analyze when assigneeIds is omitted. Defaults to 20. |
| `daysBack` | number | No | Number of days to look back for recent activity. Defaults to 30. |

### get_team_workload_breakdown

**Description**: Get a detailed per-member workload breakdown with current active tasks, task subjects, status, labels, and creation dates. Useful for stand-ups, sprint planning, and identifying what each team member is working on right now.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `assigneeIds` | string | No | Comma-separated assignee user IDs. Defaults to project members up to memberLimit. |
| `categories` | string | No | Comma-separated work item categories. Defaults to Task,Bug. |
| `status` | string | No | Comma-separated status IDs to include. Defaults to all active statuses (backlog,待处理,处理中,测试中,In Review,pending). |
| `memberLimit` | number | No | Max project members to inspect when assigneeIds is omitted. Defaults to 20. |
| `taskLimit` | number | No | Max tasks returned per member. Defaults to 10, clamped to 1-50. |

### list_milestones

**Description**: List milestones (planning checkpoints) in a Projex project. Milestones track progress against goals, distinct from sprints (time-boxed iterations) and versions (releases).

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `status` | string | No | Comma-separated milestone status IDs. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### list_testcase_repositories

**Description**: List Projex testcase repositories in a Yunxiao organization. A testcase repository is a container for organizing test cases.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### list_directories

**Description**: List directories (folders) within a Projex testcase repository.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `testRepoId` | string | Yes | Testcase repository ID (string). Use list_testcase_repositories to discover available repositories. |

### search_testcases

**Description**: Search test cases within a single Projex testcase repository.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `testRepoId` | string | Yes | Testcase repository ID (string). Use list_testcase_repositories to discover available repositories. |
| `directoryId` | string | No | Directory ID (string) to filter by. Use list_directories to discover directory IDs. |
| `subject` | string | No | Filter by testcase subject/title (contains match). |
| `conditions` | string | No | Advanced conditions JSON string. Overrides simple filters. |
| `orderBy` | string | No | Sort field: gmtCreate or name. |
| `sort` | string | No | Sort direction: asc or desc. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### list_test_plans

**Description**: List test plans in a Yunxiao organization. A test plan groups test cases for execution tracking.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |

### get_test_result_list

**Description**: Get test execution result summaries for a specific directory within a test plan.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `testPlanIdentifier` | string | Yes | Test plan ID. |
| `directoryIdentifier` | string | Yes | Test plan directory ID. |

### list_project_members

**Description**: List members in a Projex project. Use this to discover user IDs for filtering work items or assigning tasks.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `name` | string | No | Filter by member name (contains match). |
| `roleId` | string | No | Filter by project role ID, such as project.admin. Use list_project_roles to discover available roles. |

### list_project_templates

**Description**: List Projex project templates in a Yunxiao organization. Useful when setting up new projects.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |

### list_project_program

**Description**: List Projex projects bound to a project program (project group).

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `programIdentifier` | string | Yes | Project program identifier (string). |

### list_project_roles

**Description**: List roles defined in a specific Projex project.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |

### list_all_project_roles

**Description**: List all Projex project roles across a Yunxiao organization.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |

### search_projects

**Description**: Search Projex projects in a Yunxiao organization. Use this to discover available projects and obtain their project IDs before calling other project-scoped tools.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `name` | string | No | Filter by project name (contains match). |
| `status` | string | No | Comma-separated project status IDs. Use without filters to list all projects and inspect their statuses. |
| `creator` | string | No | Comma-separated creator user IDs. Use list_project_members to find user IDs in the project. |
| `conditions` | string | No | Advanced filter conditions as a JSON string. Overrides simple filters when provided. |
| `extraConditions` | string | No | Extra filter conditions as a JSON string, combined with conditions. |
| `orderBy` | string | No | Sort field, such as gmtCreate (creation time) or name. |
| `sort` | string | No | Sort direction: asc (ascending) or desc (descending). |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Most endpoints default to 20. |

### list_sprints

**Description**: List sprints (iterations) in a Projex project. Sprints are time-boxed development cycles. Use this when you need iteration-level planning data. For a broader view that also includes milestones and versions, use get_project_overview instead.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID (string). Use search_projects if you do not know the project ID. |
| `status` | string | No | Comma-separated sprint statuses. Common values: TODO (not started), DOING (active), ARCHIVED (closed). |
| `name` | string | No | Filter by sprint name (contains match). |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Most endpoints default to 20. |

### search_workitems

**Description**: Search work items in a single Projex project with flexible filtering and pagination. This is the primary tool for finding work items when you know the project. For a high-level summary across categories without detailed filtering, use get_project_workitem_summary. For a Kanban board view grouped by status, use get_project_workitem_board.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `category` | string | Yes | Work item category. Common values: Req (需求/Requirement), Task (任务), Bug (缺陷), Risk (风险). Use list_work_item_types to discover available categories for a project. |
| `projectId` | string | Yes | Project ID (string). This API searches one project at a time. Use search_projects to find the project ID. |
| `subject` | string | No | Filter by work item subject/title (contains match). |
| `status` | string | No | Comma-separated status IDs. Use get_project_workitem_context to discover valid status IDs for the target category. |
| `assignedTo` | string | No | Comma-separated assignee user IDs. Use list_project_members to find user IDs in the project. |
| `creator` | string | No | Comma-separated creator user IDs. Use list_project_members to find user IDs in the project. |
| `tag` | string | No | Comma-separated tag IDs. Use get_project_workitem_context to discover available tags. |
| `sprint` | string | No | Comma-separated sprint IDs. Use list_sprints to discover sprint IDs in the project. |
| `workitemType` | string | No | Comma-separated work item type IDs. Use list_work_item_types to discover available types. |
| `statusStage` | string | No | Comma-separated status stage IDs (higher-level status grouping). Use get_project_workitem_context to discover valid stages. |
| `priority` | string | No | Comma-separated priority IDs. Use get_project_workitem_context to discover available priorities. |
| `subjectDescription` | string | No | Filter by subject or description content (contains match). |
| `createdAfter` | string | No | Created date lower bound. Accepts YYYY-MM-DD (e.g., 2024-01-01) or ISO datetime. |
| `createdBefore` | string | No | Created date upper bound. Accepts YYYY-MM-DD or ISO datetime. |
| `updatedAfter` | string | No | Updated date lower bound. Accepts YYYY-MM-DD or ISO datetime. |
| `updatedBefore` | string | No | Updated date upper bound. Accepts YYYY-MM-DD or ISO datetime. |
| `finishTimeAfter` | string | No | Planned finish date lower bound. Accepts YYYY-MM-DD or ISO datetime. |
| `finishTimeBefore` | string | No | Planned finish date upper bound. Accepts YYYY-MM-DD or ISO datetime. |
| `updateStatusAtAfter` | string | No | Status update date lower bound. Accepts YYYY-MM-DD or ISO datetime. |
| `updateStatusAtBefore` | string | No | Status update date upper bound. Accepts YYYY-MM-DD or ISO datetime. |
| `conditions` | string | No | Advanced filter conditions as a JSON string. Overrides simple filters when provided. |
| `orderBy` | string | No | Sort field. Defaults are controlled by Yunxiao. Common values: gmtCreate, gmtModified, priority, serialNumber. |
| `sort` | string | No | Sort direction: asc (ascending) or desc (descending). |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Most endpoints default to 20. |

### list_work_item_comments

**Description**: List comments for a single Projex work item. For a comprehensive view that includes comments along with activities, attachments, and relations, use get_project_workitem_detail instead.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `workitemId` | string | Yes | Work item ID (numeric string). Find it via search_workitems or get_project_workitem_summary. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Most endpoints default to 20. |

### list_all_work_item_types

**Description**: List all work item types across a Yunxiao organization. This returns organization-level type definitions. For project-scoped types (including custom fields), use list_work_item_types instead.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `categories` | string | No | Comma-separated work item type categories to filter by. Common values: Req, Bug, Task, Risk. |

### list_work_item_types

**Description**: List work item types available in a specific Projex project. Use this before creating a work item to discover valid workitemTypeId values and categories for the target project.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID (string). Use search_projects if you do not know the project ID. |
| `category` | string | Yes | Work item category to filter by. Common values: Req, Bug, Task, Risk. Use list_all_work_item_types to see available categories. |

### list_work_item_relation_work_item_types

**Description**: List work item types that can be related to a given work item type. Use this when setting up work item relationships to discover which types are valid for the chosen relationType.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `workItemTypeId` | string | Yes | Work item type ID (numeric string). Use list_work_item_types to find available types. |
| `relationType` | string | No | Relation type filter. Valid values: PARENT, SUB, ASSOCIATED, DEPEND_ON, DEPENDED_BY. Omit to list all compatible types. |

### list_versions

**Description**: List versions (releases) in a Projex project. Versions represent release milestones, distinct from sprints (iterations) and milestones (planning checkpoints).

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `status` | string | No | Comma-separated version statuses. Common values: TODO, DOING, ARCHIVED. |
| `name` | string | No | Version name filter. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### list_workitem_activities

**Description**: List activity events (history log) for a single Projex work item. For a comprehensive view including comments and attachments, use get_project_workitem_detail instead.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `workitemId` | string | Yes | Work item ID. Use search_workitems to discover valid IDs. |

### list_workitem_attachments

**Description**: List file attachments for a Projex work item. Use search_workitems to discover work item IDs.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `workitemId` | string | Yes | Work item ID. Use search_workitems to discover valid IDs. |

### list_workitem_relation_records

**Description**: List relation records (parent, subtask, dependency links) for a Projex work item. Use search_workitems to discover work item IDs.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `workitemId` | string | Yes | Work item ID. Use search_workitems to discover valid IDs. |
| `relationType` | string | Yes | Relation type: PARENT, SUB, ASSOCIATED, DEPEND_ON, or DEPENDED_BY. |

### list_labels

**Description**: List labels in a Projex project. Labels are used to categorize and filter work items. Use search_projects to discover project IDs.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID. Use search_projects to discover valid IDs. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### create_workitem

**Description**: Create a new work item in a Projex project. Before calling this, use list_work_item_types to discover the correct category and workitemTypeId for the target project.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `projectId` | string | Yes | Project ID where the work item will be created. |
| `category` | string | Yes | Work item category. Common values: Req, Task, Bug, Risk. Use list_work_item_types to discover valid categories for the project. |
| `workitemTypeId` | string | Yes | Work item type ID (numeric string). Use list_work_item_types to find available types for the project. |
| `subject` | string | Yes | Work item title/subject. |
| `description` | string | No | Work item description. |
| `assignedTo` | string | No | Assignee user ID (string). Use list_project_members or get_project_workitem_context to discover valid user IDs. |
| `priority` | string | No | Priority ID (string). Use get_project_workitem_context to discover available priorities for the work item type. |
| `parentId` | string | No | Parent work item ID (numeric string) for creating sub-items. Use search_workitems to find the parent ID. |
| `sprint` | string | No | Sprint ID (string) to associate the work item with. Use list_sprints to discover active sprints. |
| `labels` | array | No | Label IDs to attach to the work item. |

### update_workitem

**Description**: Update an existing work item's fields. Only provided fields will be changed; omitted fields retain their current values. Use search_workitems to find the workitemId if unknown.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `workitemId` | string | Yes | Work item ID to update. |
| `subject` | string | No | New title/subject. |
| `description` | string | No | New description. |
| `assignedTo` | string | No | New assignee user ID. |
| `priority` | string | No | New priority ID. |
| `sprint` | string | No | New sprint ID. |
| `labels` | array | No | New label IDs (replaces existing labels). |

### update_workitem_status

**Description**: Change the status of a work item. Use get_project_workitem_context to discover valid statusId values for the work item type before updating.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `workitemId` | string | Yes | Work item ID (numeric string). Find it via search_workitems or get_project_workitem_summary. |
| `statusId` | string | Yes | Target status ID (string). Use get_project_workitem_context to discover valid status IDs for the work item type. |
| `comment` | string | No | Optional comment explaining the status change. |

### add_workitem_comment

**Description**: Add a comment to a work item. Use search_workitems to find the workitemId if unknown.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `workitemId` | string | Yes | Work item ID (numeric string). Find it via search_workitems or get_project_workitem_summary. |
| `content` | string | Yes | Comment content (plain text or rich text format supported by Yunxiao). |

