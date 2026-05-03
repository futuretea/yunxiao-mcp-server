# Projex Tools

This document describes the 45 read-only MCP tools in the projex domain.

## Tool Inventory

| Tool | Description |
|------|-------------|
| `list_current_user_effort_records` | List actual effort records for the current user. |
| `list_effort_records` | List actual effort records for a Projex work item. |
| `list_estimated_efforts` | List estimated effort records for a Projex work item. |
| `get_project_overview` | Get a compact Projex project overview with common project-management lists in one read-only call. |
| `get_project_workitem_summary` | Summarize Projex work items by category for one project with samples and pagination totals. |
| `get_project_workitem_context` | Get project work item metadata context: types, labels, members, and optional fields/workflow for one type. |
| `get_sprint_overview` | Get a sprint overview with work item samples and totals by category for one project sprint. |
| `get_my_project_workitems` | Get work items assigned to or created by a specific user for one project, grouped by category. |
| `get_project_workitem_board` | Get a Kanban-style board view of work items grouped by status for one project. |
| `get_project_workitem_detail` | Get comprehensive details for a single work item including basic info, activities, attachments, comments, and relation records in one read-only call. |
| `get_project_risk_dashboard` | Get a read-only project risk dashboard with category samples, overdue work items, and optional high-priority/stale sections. |
| `get_project_member_task_status` | Get per-member task status for one project with assigned, overdue, and optional status-group sections. |
| `list_milestones` | List milestones in a Projex project. |
| `list_testcase_repositories` | List Projex testcase repositories in a Yunxiao organization. |
| `list_directories` | List testcase directories in a Projex testcase repository. |
| `get_testcase_field_config` | Get testcase field configuration in a Projex testcase repository. |
| `get_testcase` | Get a Projex testcase by ID. |
| `search_testcases` | Search Projex testcases in one testcase repository. |
| `list_test_plans` | List Projex test plans in a Yunxiao organization. |
| `get_test_result_list` | Get testcase result summaries in a Projex test plan directory. |
| `list_project_members` | List members in a Projex project. |
| `list_project_templates` | List Projex project templates in a Yunxiao organization. |
| `get_project_template_field_config` | Get field configuration for a Projex project template. |
| `list_project_program` | List Projex projects bound to a project program. |
| `list_project_roles` | List roles in a Projex project. |
| `list_all_project_roles` | List all Projex project roles in a Yunxiao organization. |
| `search_projects` | Search Projex projects in a Yunxiao organization. |
| `get_project` | Get a Projex project by ID. |
| `get_sprint` | Get a Projex sprint by ID. |
| `list_sprints` | List Projex sprints in a project. |
| `search_workitems` | Search work items in one Projex project space. |
| `get_workitem` | Get a Projex work item by ID. |
| `list_work_item_comments` | List comments for a Projex work item. |
| `list_all_work_item_types` | List all work item types in a Yunxiao organization. |
| `list_work_item_types` | List work item types in one Projex project. |
| `get_work_item_type` | Get a Projex work item type by ID. |
| `list_work_item_relation_work_item_types` | List work item types that can be related to a Projex work item type. |
| `get_work_item_type_field_config` | Get field configuration for a Projex work item type. |
| `get_work_item_workflow` | Get workflow information for a Projex work item type. |
| `list_versions` | List versions in a Projex project. |
| `list_workitem_activities` | List activity events for a Projex work item. |
| `list_workitem_attachments` | List attachments for a Projex work item. |
| `get_workitem_file` | Get file metadata for a Projex work item file. |
| `list_workitem_relation_records` | List relation records for a Projex work item. |
| `list_labels` | List labels in a Projex project. |

### list_current_user_effort_records

**Description**: List actual effort records for the current user.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `startDate` | string | Yes | Start date in yyyy-MM-dd format. |
| `endDate` | string | Yes | End date in yyyy-MM-dd format. |

### list_effort_records

**Description**: List actual effort records for a Projex work item.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `workitemId` | string | Yes | Work item ID. |

### list_estimated_efforts

**Description**: List estimated effort records for a Projex work item.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `workitemId` | string | Yes | Work item ID. |

### get_project_overview

**Description**: Get a compact Projex project overview with common project-management lists in one read-only call.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
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

**Description**: Summarize Projex work items by category for one project with samples and pagination totals.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
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

**Description**: Get project work item metadata context: types, labels, members, and optional fields/workflow for one type.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
| `category` | string | Yes | Work item category, such as Req, Task, Bug, or Risk. |
| `workItemTypeId` | string | No | Optional work item type ID for field and workflow metadata. |
| `includeMembers` | boolean | No | Whether to include project members. Defaults to true. |
| `includeLabels` | boolean | No | Whether to include project labels. Defaults to true. |
| `includeFields` | boolean | No | Whether to include field configuration when workItemTypeId is set. Defaults to true. |
| `includeWorkflow` | boolean | No | Whether to include workflow metadata when workItemTypeId is set. Defaults to true. |
| `page` | number | No | Page number for labels. Defaults to 1. |
| `perPage` | number | No | Page size for labels. Defaults to 20. |

### get_sprint_overview

**Description**: Get a sprint overview with work item samples and totals by category for one project sprint.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
| `sprintId` | string | Yes | Sprint ID. |
| `categories` | string | No | Comma-separated work item categories. Defaults to Task,Bug. |
| `subject` | string | No | Subject contains filter applied to every category. |
| `status` | string | No | Comma-separated status IDs applied to every category. |
| `assignedTo` | string | No | Comma-separated assignee user IDs applied to every category. |
| `creator` | string | No | Comma-separated creator user IDs applied to every category. |
| `sampleLimit` | number | No | Samples returned per category. Defaults to 5, clamped to 0-200. |

### get_my_project_workitems

**Description**: Get work items assigned to or created by a specific user for one project, grouped by category.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
| `userId` | string | Yes | User ID to filter work items by. |
| `relation` | string | No | Filter relation: assigned or created. Defaults to assigned. |
| `categories` | string | No | Comma-separated work item categories. Defaults to Task,Bug. |
| `status` | string | No | Comma-separated status IDs applied to every category. |
| `sampleLimit` | number | No | Samples returned per category. Defaults to 5, clamped to 0-200. |

### get_project_workitem_board

**Description**: Get a Kanban-style board view of work items grouped by status for one project.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
| `category` | string | Yes | Work item category, such as Task or Bug. |
| `sprint` | string | No | Optional sprint ID to filter work items. |
| `subject` | string | No | Subject contains filter applied to the search. |
| `status` | string | No | Comma-separated status IDs applied to the search. |
| `assignedTo` | string | No | Comma-separated assignee user IDs applied to the search. |
| `creator` | string | No | Comma-separated creator user IDs applied to the search. |
| `sampleLimit` | number | No | Max work items returned. Defaults to 5, clamped to 0-200. |

### get_project_workitem_detail

**Description**: Get comprehensive details for a single work item including basic info, activities, attachments, comments, and relation records in one read-only call.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `workitemId` | string | Yes | Work item ID. |
| `includeActivities` | boolean | No | Whether to include activity history. Defaults to true. |
| `includeRelations` | boolean | No | Whether to include relation records. Defaults to true. |
| `relationTypes` | string | No | Comma-separated relation types for relation records. Defaults to ASSOCIATED,SUB. |
| `includeAttachments` | boolean | No | Whether to include attachments. Defaults to true. |
| `includeComments` | boolean | No | Whether to include comments. Defaults to true. |
| `page` | number | No | Page number for comments. Defaults to 1. |
| `perPage` | number | No | Page size for comments. Defaults to 20. |

### get_project_risk_dashboard

**Description**: Get a read-only project risk dashboard with category samples, overdue work items, and optional high-priority/stale sections.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
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

**Description**: Get per-member task status for one project with assigned, overdue, and optional status-group sections.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
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

### list_milestones

**Description**: List milestones in a Projex project.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
| `status` | string | No | Comma-separated milestone statuses. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |

### list_testcase_repositories

**Description**: List Projex testcase repositories in a Yunxiao organization.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |

### list_directories

**Description**: List testcase directories in a Projex testcase repository.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `testRepoId` | string | Yes | Testcase repository ID. |

### get_testcase_field_config

**Description**: Get testcase field configuration in a Projex testcase repository.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `testRepoId` | string | Yes | Testcase repository ID. |

### get_testcase

**Description**: Get a Projex testcase by ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `testRepoId` | string | Yes | Testcase repository ID. |
| `testcaseId` | string | Yes | Testcase ID. |

### search_testcases

**Description**: Search Projex testcases in one testcase repository.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `testRepoId` | string | Yes | Testcase repository ID. |
| `directoryId` | string | No | Directory ID filter. |
| `subject` | string | No | Testcase subject contains filter. |
| `conditions` | string | No | Advanced conditions JSON string. Overrides simple filters. |
| `orderBy` | string | No | Sort field: gmtCreate or name. |
| `sort` | string | No | Sort direction: asc or desc. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |

### list_test_plans

**Description**: List Projex test plans in a Yunxiao organization.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |

### get_test_result_list

**Description**: Get testcase result summaries in a Projex test plan directory.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `testPlanIdentifier` | string | Yes | Test plan ID. |
| `directoryIdentifier` | string | Yes | Test plan directory ID. |

### list_project_members

**Description**: List members in a Projex project.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
| `name` | string | No | Member name filter. |
| `roleId` | string | No | Project role ID filter, such as project.admin. |

### list_project_templates

**Description**: List Projex project templates in a Yunxiao organization.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |

### get_project_template_field_config

**Description**: Get field configuration for a Projex project template.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectTemplateId` | string | Yes | Project template ID. |

### list_project_program

**Description**: List Projex projects bound to a project program.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `programIdentifier` | string | Yes | Project program identifier. |

### list_project_roles

**Description**: List roles in a Projex project.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |

### list_all_project_roles

**Description**: List all Projex project roles in a Yunxiao organization.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |

### search_projects

**Description**: Search Projex projects in a Yunxiao organization.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `name` | string | No | Project name contains filter. |
| `status` | string | No | Comma-separated project status IDs. |
| `creator` | string | No | Comma-separated creator user IDs. |
| `conditions` | string | No | Advanced conditions JSON string. Overrides simple filters. |
| `extraConditions` | string | No | Extra conditions JSON string. |
| `orderBy` | string | No | Sort field, such as gmtCreate or name. |
| `sort` | string | No | Sort direction: asc or desc. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |

### get_project

**Description**: Get a Projex project by ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |

### get_sprint

**Description**: Get a Projex sprint by ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
| `sprintId` | string | Yes | Sprint ID. |

### list_sprints

**Description**: List Projex sprints in a project.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
| `status` | string | No | Comma-separated sprint statuses: TODO, DOING, ARCHIVED. |
| `name` | string | No | Sprint name filter. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |

### search_workitems

**Description**: Search work items in one Projex project space.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `category` | string | Yes | Work item category, such as Req, Task, Bug, or Risk. |
| `projectId` | string | Yes | Project ID. This API searches one project at a time. |
| `subject` | string | No | Subject contains filter. |
| `status` | string | No | Comma-separated status IDs. |
| `assignedTo` | string | No | Comma-separated assignee user IDs. |
| `creator` | string | No | Comma-separated creator user IDs. |
| `tag` | string | No | Comma-separated tag IDs. |
| `sprint` | string | No | Comma-separated sprint IDs. |
| `workitemType` | string | No | Comma-separated work item type IDs. |
| `statusStage` | string | No | Comma-separated status stage IDs. |
| `priority` | string | No | Comma-separated priority IDs. |
| `subjectDescription` | string | No | Subject or description contains filter. |
| `createdAfter` | string | No | Created date lower bound, YYYY-MM-DD or datetime. |
| `createdBefore` | string | No | Created date upper bound, YYYY-MM-DD or datetime. |
| `updatedAfter` | string | No | Updated date lower bound, YYYY-MM-DD or datetime. |
| `updatedBefore` | string | No | Updated date upper bound, YYYY-MM-DD or datetime. |
| `finishTimeAfter` | string | No | Planned finish date lower bound, YYYY-MM-DD or datetime. |
| `finishTimeBefore` | string | No | Planned finish date upper bound, YYYY-MM-DD or datetime. |
| `updateStatusAtAfter` | string | No | Status update date lower bound, YYYY-MM-DD or datetime. |
| `updateStatusAtBefore` | string | No | Status update date upper bound, YYYY-MM-DD or datetime. |
| `conditions` | string | No | Advanced conditions JSON string. Overrides simple filters. |
| `orderBy` | string | No | Sort field. Defaults are controlled by Yunxiao. |
| `sort` | string | No | Sort direction: asc or desc. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |

### get_workitem

**Description**: Get a Projex work item by ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `workitemId` | string | Yes | Work item ID. |

### list_work_item_comments

**Description**: List comments for a Projex work item.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `workItemId` | string | Yes | Work item ID. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |

### list_all_work_item_types

**Description**: List all work item types in a Yunxiao organization.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `categories` | string | No | Optional work item type categories, such as Req, Bug, or Task. |

### list_work_item_types

**Description**: List work item types in one Projex project.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
| `category` | string | Yes | Work item category, such as Req, Bug, or Task. |

### get_work_item_type

**Description**: Get a Projex work item type by ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `workItemTypeId` | string | Yes | Work item type ID. |

### list_work_item_relation_work_item_types

**Description**: List work item types that can be related to a Projex work item type.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `workItemTypeId` | string | Yes | Work item type ID. |
| `relationType` | string | No | Relation type: PARENT, SUB, ASSOCIATED, DEPEND_ON, or DEPENDED_BY. |

### get_work_item_type_field_config

**Description**: Get field configuration for a Projex work item type.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
| `workItemTypeId` | string | Yes | Work item type ID. |

### get_work_item_workflow

**Description**: Get workflow information for a Projex work item type.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
| `workItemTypeId` | string | Yes | Work item type ID. |

### list_versions

**Description**: List versions in a Projex project.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
| `status` | string | No | Comma-separated version statuses: TODO, DOING, ARCHIVED. |
| `name` | string | No | Version name filter. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |

### list_workitem_activities

**Description**: List activity events for a Projex work item.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `workitemId` | string | Yes | Work item ID. |

### list_workitem_attachments

**Description**: List attachments for a Projex work item.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `workitemId` | string | Yes | Work item ID. |

### get_workitem_file

**Description**: Get file metadata for a Projex work item file.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `workitemId` | string | Yes | Work item ID. |
| `fileId` | string | Yes | File ID. |

### list_workitem_relation_records

**Description**: List relation records for a Projex work item.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `workitemId` | string | Yes | Work item ID. |
| `relationType` | string | Yes | Relation type: PARENT, SUB, ASSOCIATED, DEPEND_ON, or DEPENDED_BY. |

### list_labels

**Description**: List labels in a Projex project.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `projectId` | string | Yes | Project ID. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |

