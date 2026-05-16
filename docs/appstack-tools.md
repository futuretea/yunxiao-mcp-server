# Appstack Tools

This document describes the 31 MCP tools in the appstack domain.

Access summary: 31 read-only, 0 write-capable.

## Enhanced Tools

These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations. Prefer them when available.

| Tool | Description |
|------|-------------|
| `get_application_overview` | Get a comprehensive overview of an Appstack application including basic info, environments, and recent orchestrations in one read-only call. Use this after discovering applications via list_applications or list_attached_apps. |
| `get_environment_overview` | Get a comprehensive overview of an Appstack environment including basic info, variable groups, and latest orchestration in one read-only call. Use this after identifying an application and environment via get_application_overview. |
| `get_release_overview` | Get a comprehensive overview of an Appstack system release including basic info, members, products, and attached change requests in one read-only call. Use this after discovering a release via search_releases or list_system_release_workflows. |

## Pagination

Tools in this domain use the following pagination scheme(s):

- Keyset (nextToken)
- Offset (current/pageSize)
- Offset (page/perPage)

## Tool Inventory

Tools marked in **bold** are enhanced aggregation tools.

| Tool | Access | Description |
|------|--------|-------------|
| `search_app_templates` | Read-only | Search AppStack application templates. Use this to discover templates before creating or deploying applications. |
| `list_environments` | Read-only | List AppStack environments for an application. Use list_applications to discover valid application names. |
| `list_application_members` | Read-only | List members of an AppStack application. Use list_applications to discover valid application names. |
| `list_application_sources` | Read-only | List source repositories attached to an AppStack application. Use list_applications to discover valid application names. |
| `list_resource_instances` | Read-only | List AppStack resource instances in a resource pool. Use list_resource_pools to discover valid pool names. |
| **`get_application_overview`** | Read-only | Get a comprehensive overview of an Appstack application including basic info, environments, and recent orchestrations in one read-only call. Use this after discovering applications via list_applications or list_attached_apps. |
| **`get_environment_overview`** | Read-only | Get a comprehensive overview of an Appstack environment including basic info, variable groups, and latest orchestration in one read-only call. Use this after identifying an application and environment via get_application_overview. |
| **`get_release_overview`** | Read-only | Get a comprehensive overview of an Appstack system release including basic info, members, products, and attached change requests in one read-only call. Use this after discovering a release via search_releases or list_system_release_workflows. |
| `list_global_vars` | Read-only | Search AppStack global variable groups. Use this to discover variable group IDs before reading or updating specific groups. |
| `search_releases` | Read-only | Search AppStack releases in a Yunxiao organization. Use this to discover releases before calling get_release_overview or other release-specific tools. |
| `list_system_release_workflows` | Read-only | List AppStack release workflows for a system. Use this after discovering a system via list_systems to find releases and their workflows. |
| `list_release_members` | Read-only | List members of an AppStack system release. Use this after discovering a release via search_releases or list_system_release_workflows. |
| `list_release_products` | Read-only | List products attached to an AppStack system release. Use this after discovering a release via search_releases or list_system_release_workflows. |
| `list_attached_change_requests` | Read-only | List change requests attached to an AppStack system release. Use this after discovering a release via search_releases or list_system_release_workflows. |
| `list_release_executions` | Read-only | List execution records for an AppStack system release. Use this after discovering release workflow and stage details via list_system_release_workflows or get_release_overview. |
| `list_systems` | Read-only | List AppStack systems in a Yunxiao organization. Use this as the entry point to discover systems before calling other system-specific tools. |
| `list_attached_apps` | Read-only | List applications attached to an AppStack system. Use this to discover applications within a system before calling get_application_overview. |
| `list_system_members` | Read-only | List members of an AppStack system. Use this after discovering a system via list_systems. |
| `list_applications` | Read-only | List AppStack applications in a Yunxiao organization. AppStack manages deployment environments and release pipelines. |
| `list_app_orchestration` | Read-only | List AppStack orchestrations for an application. Orchestrations define deployment workflows and environment configurations. |
| `list_app_release_workflows` | Read-only | List AppStack release workflows for an application. Release workflows model multi-stage deployment pipelines. |
| `list_app_release_workflow_briefs` | Read-only | List AppStack release workflow briefs for an application. Briefs provide a condensed view of workflow definitions. |
| `list_app_release_stage_briefs` | Read-only | List AppStack release stage briefs for an application release workflow. Stages represent individual deployment phases within a workflow. |
| `list_app_release_stage_runs` | Read-only | List AppStack release stage execution records. Each record represents a single run of a deployment stage. |
| `list_app_release_stage_exec_metadata` | Read-only | List integrated change metadata for an AppStack release stage execution. Metadata includes linked work items and commits. |
| `list_appstack_change_request_executions` | Read-only | List execution records for an AppStack change request. Change requests track planned deployments. |
| `list_appstack_change_request_work_items` | Read-only | List work items for an AppStack change request. Work items represent linked Projex tasks or requirements. |
| `list_change_order_versions` | Read-only | List AppStack change order versions. Change orders track actual deployments and their versions. |
| `list_change_orders_by_origin` | Read-only | List AppStack change orders by creation origin. Use this to trace deployments back to their source (e.g., a Flow pipeline). |
| `list_change_order_job_logs` | Read-only | List logs for an AppStack change order job. Job logs capture deployment script output. |
| `find_task_operation_log` | Read-only | Get an AppStack deployment task operation log. Operation logs record manual or automated actions taken during deployment. |

### search_app_templates

**Description**: Search AppStack application templates. Use this to discover templates before creating or deploying applications.

**Access**: Read-only

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `pagination` | string | No | Pagination mode. Valid value: keyset. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `orderBy` | string | No | Sort field. Valid values: id, gmtCreate. |
| `sort` | string | No | Sort direction. Valid values: asc, desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `displayNameKeyword` | string | No | Template display name keyword for filtering results. |
| `page` | number | No | Page number for pagination. Starts at 1. |

### list_environments

**Description**: List AppStack environments for an application. Use list_applications to discover valid application names.

**Access**: Read-only

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid names. |
| `pagination` | string | No | Pagination mode. Valid value: keyset. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `orderBy` | string | No | Sort field. Valid values: id, gmtCreate. |
| `sort` | string | No | Sort direction. Valid values: asc, desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `page` | number | No | Page number for pagination. Starts at 1. |

### list_application_members

**Description**: List members of an AppStack application. Use list_applications to discover valid application names.

**Access**: Read-only

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid names. |
| `current` | number | No | Page number for pagination. Starts at 1. |
| `pageSize` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### list_application_sources

**Description**: List source repositories attached to an AppStack application. Use list_applications to discover valid application names.

**Access**: Read-only

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid names. |
| `pagination` | string | No | Pagination mode. Valid value: keyset. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `orderBy` | string | No | Sort field. Valid values: id, gmtCreate. |
| `sort` | string | No | Sort direction. Valid values: asc, desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `page` | number | No | Page number for pagination. Starts at 1. |

### list_resource_instances

**Description**: List AppStack resource instances in a resource pool. Use list_resource_pools to discover valid pool names.

**Access**: Read-only

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `poolName` | string | Yes | Resource pool name. Use list_resource_pools to discover valid names. |
| `pagination` | string | No | Pagination mode. Valid value: keyset. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `orderBy` | string | No | Sort field. Valid values: id, gmtCreate. |
| `sort` | string | No | Sort direction. Valid values: asc, desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `page` | number | No | Page number for pagination. Starts at 1. |

### get_application_overview

**Description**: Get a comprehensive overview of an Appstack application including basic info, environments, and recent orchestrations in one read-only call. Use this after discovering applications via list_applications or list_attached_apps.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application unique name. Use list_applications or list_attached_apps to discover valid names. |
| `includeEnvironments` | boolean | No | Whether to include environment list. Defaults to true. |
| `includeOrchestrations` | boolean | No | Whether to include recent orchestrations. Defaults to true. |
| `envLimit` | number | No | Max environments returned. Defaults to 5. |
| `orchestrationLimit` | number | No | Max orchestrations returned. Defaults to 5. |

### get_environment_overview

**Description**: Get a comprehensive overview of an Appstack environment including basic info, variable groups, and latest orchestration in one read-only call. Use this after identifying an application and environment via get_application_overview.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application unique name. Use list_applications or list_attached_apps to discover valid names. |
| `envName` | string | Yes | Environment name. Use get_application_overview to discover valid environment names. |
| `includeVariableGroups` | boolean | No | Whether to include environment variable groups. Defaults to true. |
| `includeLatestOrchestration` | boolean | No | Whether to include the latest available orchestration. Defaults to true. |

### get_release_overview

**Description**: Get a comprehensive overview of an Appstack system release including basic info, members, products, and attached change requests in one read-only call. Use this after discovering a release via search_releases or list_system_release_workflows.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `systemName` | string | Yes | System name. Use list_systems to discover valid names. |
| `sn` | string | Yes | Release serial number. Use search_releases or list_system_release_workflows to discover valid serial numbers. |
| `includeMembers` | boolean | No | Whether to include release members. Defaults to true. |
| `includeProducts` | boolean | No | Whether to include release products. Defaults to true. |
| `includeChangeRequests` | boolean | No | Whether to include attached change requests. Defaults to true. |
| `changeRequestLimit` | number | No | Max change requests returned. Defaults to 5. |

### list_global_vars

**Description**: Search AppStack global variable groups. Use this to discover variable group IDs before reading or updating specific groups.

**Access**: Read-only

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `current` | number | No | Page number for pagination. Starts at 1. |
| `pageSize` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `search` | string | No | Optional search keyword for filtering variable groups by name. |

### search_releases

**Description**: Search AppStack releases in a Yunxiao organization. Use this to discover releases before calling get_release_overview or other release-specific tools.

**Access**: Read-only

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `pagination` | string | No | Pagination mode. Valid value: keyset. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `orderBy` | string | No | Sort field. Valid values: id, gmt_create. |
| `sort` | string | No | Sort direction. Valid values: asc, desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `nameKeyword` | string | No | Release display-name search keyword. |
| `systemName` | string | No | System unique name. Use list_systems to discover valid names. |
| `states` | array | No | Release states filter. Valid values: DEVELOPING, RELEASING, CLOSED, RELEASED. |

### list_system_release_workflows

**Description**: List AppStack release workflows for a system. Use this after discovering a system via list_systems to find releases and their workflows.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `systemName` | string | Yes | System name. Use list_systems to discover valid names. |

### list_release_members

**Description**: List members of an AppStack system release. Use this after discovering a release via search_releases or list_system_release_workflows.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `systemName` | string | Yes | System name. Use list_systems to discover valid names. |
| `sn` | string | Yes | Release serial number. Use search_releases or list_system_release_workflows to discover valid serial numbers. |

### list_release_products

**Description**: List products attached to an AppStack system release. Use this after discovering a release via search_releases or list_system_release_workflows.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `systemName` | string | Yes | System name. Use list_systems to discover valid names. |
| `sn` | string | Yes | Release serial number. Use search_releases or list_system_release_workflows to discover valid serial numbers. |

### list_attached_change_requests

**Description**: List change requests attached to an AppStack system release. Use this after discovering a release via search_releases or list_system_release_workflows.

**Access**: Read-only

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `systemName` | string | Yes | System name. Use list_systems to discover valid names. |
| `releaseSn` | string | Yes | Release serial number. Use search_releases or list_system_release_workflows to discover valid serial numbers. |
| `current` | number | No | Page number for pagination. Starts at 1. |
| `pageSize` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### list_release_executions

**Description**: List execution records for an AppStack system release. Use this after discovering release workflow and stage details via list_system_release_workflows or get_release_overview.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `systemName` | string | Yes | System name. Use list_systems to discover valid names. |
| `sn` | string | Yes | Release serial number. Use search_releases or list_system_release_workflows to discover valid serial numbers. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. Use list_system_release_workflows to discover valid serial numbers. |
| `releaseStageSn` | string | Yes | Release stage serial number. Use list_system_release_workflows to discover valid serial numbers. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `orderBy` | string | No | Sort field. Valid values: id, gmtCreate. |
| `sort` | string | No | Sort direction. Valid values: asc, desc. |

### list_systems

**Description**: List AppStack systems in a Yunxiao organization. Use this as the entry point to discover systems before calling other system-specific tools.

**Access**: Read-only

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `current` | number | No | Page number for pagination. Starts at 1. |
| `pageSize` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### list_attached_apps

**Description**: List applications attached to an AppStack system. Use this to discover applications within a system before calling get_application_overview.

**Access**: Read-only

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `systemName` | string | Yes | System name. Use list_systems to discover valid names. |
| `current` | number | No | Page number for pagination. Starts at 1. |
| `pageSize` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### list_system_members

**Description**: List members of an AppStack system. Use this after discovering a system via list_systems.

**Access**: Read-only

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `systemName` | string | Yes | System name. Use list_systems to discover valid names. |
| `current` | number | No | Page number for pagination. Starts at 1. |
| `pageSize` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### list_applications

**Description**: List AppStack applications in a Yunxiao organization. AppStack manages deployment environments and release pipelines.

**Access**: Read-only

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `pagination` | string | No | Pagination mode. Valid value: keyset. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `orderBy` | string | No | Sort field. Valid values: id, gmtCreate. |
| `sort` | string | No | Sort direction. Valid values: asc, desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `page` | number | No | Page number for pagination. Starts at 1. |

### list_app_orchestration

**Description**: List AppStack orchestrations for an application. Orchestrations define deployment workflows and environment configurations.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |

### list_app_release_workflows

**Description**: List AppStack release workflows for an application. Release workflows model multi-stage deployment pipelines.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |

### list_app_release_workflow_briefs

**Description**: List AppStack release workflow briefs for an application. Briefs provide a condensed view of workflow definitions.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |

### list_app_release_stage_briefs

**Description**: List AppStack release stage briefs for an application release workflow. Stages represent individual deployment phases within a workflow.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. Use list_app_release_workflows to discover valid values. |

### list_app_release_stage_runs

**Description**: List AppStack release stage execution records. Each record represents a single run of a deployment stage.

**Access**: Read-only

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. Use list_app_release_workflows to discover valid values. |
| `releaseStageSn` | string | Yes | Release stage serial number. Use list_app_release_stage_briefs to discover valid values. |
| `pagination` | string | No | Pagination mode. Valid value: keyset. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `orderBy` | string | No | Sort field. Valid values: id, gmtCreate. |
| `sort` | string | No | Sort direction. Valid values: asc, desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `page` | number | No | Page number for pagination. Starts at 1. |

### list_app_release_stage_exec_metadata

**Description**: List integrated change metadata for an AppStack release stage execution. Metadata includes linked work items and commits.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. Use list_app_release_workflows to discover valid values. |
| `releaseStageSn` | string | Yes | Release stage serial number. Use list_app_release_stage_briefs to discover valid values. |
| `executionNumber` | string | Yes | Release stage execution number. Use list_app_release_stage_runs to discover valid values. |

### list_appstack_change_request_executions

**Description**: List execution records for an AppStack change request. Change requests track planned deployments.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `sn` | string | Yes | Change request serial number. Typically discovered via list_attached_change_requests. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. Use list_app_release_workflows to discover valid values. |
| `releaseStageSn` | string | Yes | Release stage serial number. Use list_app_release_stage_briefs to discover valid values. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `orderBy` | string | No | Sort field. Valid values: id, gmtCreate. |
| `sort` | string | No | Sort direction. Valid values: asc, desc. |

### list_appstack_change_request_work_items

**Description**: List work items for an AppStack change request. Work items represent linked Projex tasks or requirements.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `sn` | string | Yes | Change request serial number. Typically discovered via list_attached_change_requests. |

### list_change_order_versions

**Description**: List AppStack change order versions. Change orders track actual deployments and their versions.

**Access**: Read-only

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `envNames` | string | No | Comma-separated environment names. Use list_environments to discover valid environment names for an application. |
| `creators` | string | No | Comma-separated creator account IDs. |
| `current` | number | No | Page number for pagination. Starts at 1. |
| `pageSize` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### list_change_orders_by_origin

**Description**: List AppStack change orders by creation origin. Use this to trace deployments back to their source (e.g., a Flow pipeline).

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `originType` | string | Yes | Origin type indicating the source system. Valid value: FLOW (Flow pipeline). |
| `originId` | string | Yes | Origin identifier from the source system. For FLOW origin, use a pipeline run ID or pipeline ID. |
| `appName` | string | No | Application name filter. Use list_applications to discover valid app names. |
| `envName` | string | No | Environment name filter. Use list_environments to discover valid environment names for an application. |

### list_change_order_job_logs

**Description**: List logs for an AppStack change order job. Job logs capture deployment script output.

**Access**: Read-only

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `changeOrderSn` | string | Yes | Change order serial number. Use list_change_order_versions to discover valid values. |
| `jobSn` | string | Yes | Change order job serial number. Typically returned in change order details. |
| `current` | number | No | Page number for pagination. Starts at 1. |
| `pageSize` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### find_task_operation_log

**Description**: Get an AppStack deployment task operation log. Operation logs record manual or automated actions taken during deployment.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `changeOrderSn` | string | Yes | Change order serial number. Use list_change_order_versions to discover valid values. |
| `jobSn` | string | Yes | Change order job serial number. Typically returned in change order details. |
| `stageSn` | string | Yes | Deployment stage serial number. Typically returned in change order job details. |
| `taskSn` | string | Yes | Deployment task serial number. Typically returned in change order stage details. |

