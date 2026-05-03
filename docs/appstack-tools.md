# Appstack Tools

This document describes the 53 read-only MCP tools in the appstack domain.

## Enhanced Tools

These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations. Prefer them when available.

| Tool | Description |
|------|-------------|
| `get_application_overview` | Get a comprehensive overview of an Appstack application including basic info, environments, and recent orchestrations in one read-only call. |
| `get_environment_overview` | Get a comprehensive overview of an Appstack environment including basic info, variable groups, and latest orchestration in one read-only call. |
| `get_release_overview` | Get a comprehensive overview of an Appstack system release including basic info, members, products, and attached change requests in one read-only call. |

## Pagination

Tools in this domain use the following pagination scheme(s):

- Keyset (nextToken)
- Offset (current/pageSize)
- Offset (page/perPage)

## Tool Inventory

Tools marked in **bold** are enhanced aggregation tools.

| Tool | Description |
|------|-------------|
| `search_app_templates` | Search AppStack application templates. |
| `list_environments` | List AppStack environments for an application. |
| `get_environment` | Get an AppStack environment by application and environment name. |
| `list_application_members` | List members of an AppStack application. |
| `list_application_sources` | List source repositories attached to an AppStack application. |
| `get_machine_deploy_log` | Get an AppStack machine deployment log. |
| `get_deploy_group` | Get an AppStack deploy group by pool and group name. |
| `list_resource_instances` | List AppStack resource instances in a resource pool. |
| `get_resource_instance` | Get an AppStack resource instance by pool and instance name. |
| **`get_application_overview`** | Get a comprehensive overview of an Appstack application including basic info, environments, and recent orchestrations in one read-only call. |
| **`get_environment_overview`** | Get a comprehensive overview of an Appstack environment including basic info, variable groups, and latest orchestration in one read-only call. |
| **`get_release_overview`** | Get a comprehensive overview of an Appstack system release including basic info, members, products, and attached change requests in one read-only call. |
| `get_global_var` | Get an AppStack global variable group. |
| `list_global_vars` | Search AppStack global variable groups. |
| `search_releases` | Search AppStack releases in a Yunxiao organization. |
| `get_pod_container_log` | Get logs from a pod container through AppStack resource proxy. |
| `get_pod_info` | Get pod information through AppStack resource proxy. |
| `get_kubernetes_object_info` | Get Kubernetes object information through AppStack resource proxy. |
| `get_deployment_revision_info` | Get AppStack deployment workload revision information. |
| `list_system_release_workflows` | List AppStack release workflows for a system. |
| `get_release` | Get an AppStack system release by serial number. |
| `list_release_members` | List members of an AppStack system release. |
| `list_release_products` | List products attached to an AppStack system release. |
| `list_attached_change_requests` | List change requests attached to an AppStack system release. |
| `list_release_executions` | List execution records for an AppStack system release. |
| `list_systems` | List AppStack systems in a Yunxiao organization. |
| `list_attached_apps` | List applications attached to an AppStack system. |
| `list_system_members` | List members of an AppStack system. |
| `list_applications` | List AppStack applications in a Yunxiao organization. |
| `get_application` | Get an AppStack application by name. |
| `get_env_variable_groups` | Get AppStack variable groups bound to an environment. |
| `get_variable_group` | Get an AppStack application variable group by name. |
| `get_app_variable_groups` | Get AppStack variable groups for an application. |
| `get_app_variable_groups_revision` | Get the revision of AppStack variable groups for an application. |
| `get_latest_orchestration` | Get the latest AppStack orchestration available for an environment. |
| `list_app_orchestration` | List AppStack orchestrations for an application. |
| `get_app_orchestration` | Get an AppStack application orchestration by serial number. |
| `list_app_release_workflows` | List AppStack release workflows for an application. |
| `list_app_release_workflow_briefs` | List AppStack release workflow briefs for an application. |
| `get_app_release_workflow_stage` | Get an AppStack application release workflow stage. |
| `list_app_release_stage_briefs` | List AppStack release stage briefs for an application release workflow. |
| `list_app_release_stage_runs` | List AppStack release stage execution records. |
| `list_app_release_stage_exec_metadata` | List integrated change metadata for an AppStack release stage execution. |
| `get_app_release_stage_pipeline_run` | Get the pipeline run for an AppStack release stage execution. |
| `get_app_release_stage_pipeline_job_log` | Get a pipeline job log for an AppStack release stage execution. |
| `get_appstack_change_request_audit_items` | Get audit items for an AppStack change request. |
| `list_appstack_change_request_executions` | List execution records for an AppStack change request. |
| `list_appstack_change_request_work_items` | List work items for an AppStack change request. |
| `list_change_order_versions` | List AppStack change order versions. |
| `get_change_order` | Get an AppStack change order by serial number. |
| `list_change_orders_by_origin` | List AppStack change orders by creation origin. |
| `list_change_order_job_logs` | List logs for an AppStack change order job. |
| `find_task_operation_log` | Get an AppStack deployment task operation log. |

### search_app_templates

**Description**: Search AppStack application templates.

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `pagination` | string | No | Pagination mode. Yunxiao currently supports keyset. |
| `perPage` | number | No | Page size, up to 100. |
| `orderBy` | string | No | Sort field: id or gmtCreate. |
| `sort` | string | No | Sort direction: asc or desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `displayNameKeyword` | string | No | Template display name keyword. |
| `page` | number | No | Page number when using page pagination. |

### list_environments

**Description**: List AppStack environments for an application.

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `pagination` | string | No | Pagination mode. Yunxiao currently supports keyset. |
| `perPage` | number | No | Page size, up to 100. |
| `orderBy` | string | No | Sort field: id or gmtCreate. |
| `sort` | string | No | Sort direction: asc or desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `page` | number | No | Page number when using page pagination. |

### get_environment

**Description**: Get an AppStack environment by application and environment name.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `envName` | string | Yes | Environment name. |

### list_application_members

**Description**: List members of an AppStack application.

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `current` | number | No | Current page number. Defaults to 1. |
| `pageSize` | number | No | Page size. Defaults to 10. |

### list_application_sources

**Description**: List source repositories attached to an AppStack application.

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `pagination` | string | No | Pagination mode. Use keyset for keyset pagination. |
| `perPage` | number | No | Page size, up to 100. |
| `orderBy` | string | No | Sort field: id or gmtCreate. |
| `sort` | string | No | Sort direction: asc or desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `page` | number | No | Page number when using page pagination. |

### get_machine_deploy_log

**Description**: Get an AppStack machine deployment log.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `tunnelId` | number | Yes | Deployment tunnel ID. |
| `machineSn` | string | Yes | Machine serial number. |

### get_deploy_group

**Description**: Get an AppStack deploy group by pool and group name.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `poolName` | string | Yes | Resource pool name. |
| `deployGroupName` | string | Yes | Deploy group name. |

### list_resource_instances

**Description**: List AppStack resource instances in a resource pool.

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `poolName` | string | Yes | Resource pool name. |
| `pagination` | string | No | Pagination mode. Yunxiao currently supports keyset. |
| `perPage` | number | No | Page size, up to 100. |
| `orderBy` | string | No | Sort field: id or gmtCreate. |
| `sort` | string | No | Sort direction: asc or desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `page` | number | No | Page number when using page pagination. |

### get_resource_instance

**Description**: Get an AppStack resource instance by pool and instance name.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `poolName` | string | Yes | Resource pool name. |
| `instanceName` | string | Yes | Resource instance name. |

### get_application_overview

**Description**: Get a comprehensive overview of an Appstack application including basic info, environments, and recent orchestrations in one read-only call.

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application unique name. |
| `includeEnvironments` | boolean | No | Whether to include environment list. Defaults to true. |
| `includeOrchestrations` | boolean | No | Whether to include recent orchestrations. Defaults to true. |
| `envLimit` | number | No | Max environments returned. Defaults to 5. |
| `orchestrationLimit` | number | No | Max orchestrations returned. Defaults to 5. |

### get_environment_overview

**Description**: Get a comprehensive overview of an Appstack environment including basic info, variable groups, and latest orchestration in one read-only call.

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application unique name. |
| `envName` | string | Yes | Environment name. |
| `includeVariableGroups` | boolean | No | Whether to include environment variable groups. Defaults to true. |
| `includeLatestOrchestration` | boolean | No | Whether to include the latest available orchestration. Defaults to true. |

### get_release_overview

**Description**: Get a comprehensive overview of an Appstack system release including basic info, members, products, and attached change requests in one read-only call.

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `systemName` | string | Yes | System name. |
| `sn` | string | Yes | Release serial number. |
| `includeMembers` | boolean | No | Whether to include release members. Defaults to true. |
| `includeProducts` | boolean | No | Whether to include release products. Defaults to true. |
| `includeChangeRequests` | boolean | No | Whether to include attached change requests. Defaults to true. |
| `changeRequestLimit` | number | No | Max change requests returned. Defaults to 5. |

### get_global_var

**Description**: Get an AppStack global variable group.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `name` | string | Yes | Global variable group name. |
| `revisionSha` | string | No | Optional global variable group revision SHA. |

### list_global_vars

**Description**: Search AppStack global variable groups.

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `current` | number | No | Current page number. Defaults to 1. |
| `pageSize` | number | No | Page size. Defaults to 10. |
| `search` | string | No | Optional search keyword. |

### search_releases

**Description**: Search AppStack releases in a Yunxiao organization.

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `pagination` | string | No | Pagination mode. Yunxiao currently supports keyset. |
| `perPage` | number | No | Page size, up to 100. |
| `orderBy` | string | No | Sort field: id or gmt_create. |
| `sort` | string | No | Sort direction: asc or desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `nameKeyword` | string | No | Release display-name search keyword. |
| `systemName` | string | No | System unique name. |
| `states` | array | No | Release states such as DEVELOPING, RELEASING, CLOSED, or RELEASED. |

### get_pod_container_log

**Description**: Get logs from a pod container through AppStack resource proxy.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `resourcePath` | string | Yes | Resource path segment. |
| `namespace` | string | Yes | Kubernetes namespace. |
| `name` | string | Yes | Pod name. |
| `container` | string | Yes | Container name. |
| `tailingLines` | number | No | Number of log tail lines. Defaults to 1000. |

### get_pod_info

**Description**: Get pod information through AppStack resource proxy.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `resourcePath` | string | Yes | Resource path segment. |
| `namespace` | string | Yes | Kubernetes namespace. |
| `name` | string | Yes | Pod name. |
| `taskSn` | string | No | Optional task serial number. |

### get_kubernetes_object_info

**Description**: Get Kubernetes object information through AppStack resource proxy.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `resourcePath` | string | Yes | Resource path segment. |
| `namespace` | string | Yes | Kubernetes namespace. |
| `kind` | string | Yes | Kubernetes object kind. |
| `name` | string | Yes | Kubernetes object name. |
| `taskSn` | string | No | Optional task serial number. |

### get_deployment_revision_info

**Description**: Get AppStack deployment workload revision information.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `envName` | string | Yes | Environment name. |
| `namespace` | string | Yes | Kubernetes namespace. |
| `name` | string | Yes | Deployment name. |
| `revision` | string | Yes | Deployment revision. |
| `taskSn` | string | No | Optional task serial number. |

### list_system_release_workflows

**Description**: List AppStack release workflows for a system.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `systemName` | string | Yes | System name. |

### get_release

**Description**: Get an AppStack system release by serial number.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `systemName` | string | Yes | System name. |
| `sn` | string | Yes | Release serial number. |

### list_release_members

**Description**: List members of an AppStack system release.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `systemName` | string | Yes | System name. |
| `sn` | string | Yes | Release serial number. |

### list_release_products

**Description**: List products attached to an AppStack system release.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `systemName` | string | Yes | System name. |
| `sn` | string | Yes | Release serial number. |

### list_attached_change_requests

**Description**: List change requests attached to an AppStack system release.

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `systemName` | string | Yes | System name. |
| `releaseSn` | string | Yes | Release serial number. |
| `current` | number | No | Current page number. Defaults to 1. |
| `pageSize` | number | No | Page size. Defaults to 10. |

### list_release_executions

**Description**: List execution records for an AppStack system release.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `systemName` | string | Yes | System name. |
| `sn` | string | Yes | Release serial number. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. |
| `releaseStageSn` | string | Yes | Release stage serial number. |
| `perPage` | number | No | Page size, up to 100. |
| `page` | number | No | Page number. |
| `orderBy` | string | No | Sort field: id or gmtCreate. |
| `sort` | string | No | Sort direction: asc or desc. |

### list_systems

**Description**: List AppStack systems in a Yunxiao organization.

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `current` | number | No | Current page number. Defaults to 1. |
| `pageSize` | number | No | Page size. Defaults to 10. |

### list_attached_apps

**Description**: List applications attached to an AppStack system.

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `systemName` | string | Yes | System name. |
| `current` | number | No | Current page number. Defaults to 1. |
| `pageSize` | number | No | Page size. Defaults to 10. |

### list_system_members

**Description**: List members of an AppStack system.

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `systemName` | string | Yes | System name. |
| `current` | number | No | Current page number. Defaults to 1. |
| `pageSize` | number | No | Page size. Defaults to 10. |

### list_applications

**Description**: List AppStack applications in a Yunxiao organization.

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `pagination` | string | No | Pagination mode. Yunxiao currently supports keyset. |
| `perPage` | number | No | Page size, up to 100. |
| `orderBy` | string | No | Sort field: id or gmtCreate. |
| `sort` | string | No | Sort direction: asc or desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `page` | number | No | Page number when using page pagination. |

### get_application

**Description**: Get an AppStack application by name.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |

### get_env_variable_groups

**Description**: Get AppStack variable groups bound to an environment.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `envName` | string | Yes | Environment name. |

### get_variable_group

**Description**: Get an AppStack application variable group by name.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `variableGroupName` | string | Yes | Variable group name. |

### get_app_variable_groups

**Description**: Get AppStack variable groups for an application.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |

### get_app_variable_groups_revision

**Description**: Get the revision of AppStack variable groups for an application.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |

### get_latest_orchestration

**Description**: Get the latest AppStack orchestration available for an environment.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `envName` | string | Yes | Environment name. |

### list_app_orchestration

**Description**: List AppStack orchestrations for an application.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |

### get_app_orchestration

**Description**: Get an AppStack application orchestration by serial number.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `sn` | string | Yes | Orchestration serial number. |
| `tagName` | string | No | Optional orchestration tag. |
| `sha` | string | No | Optional orchestration commit SHA. |

### list_app_release_workflows

**Description**: List AppStack release workflows for an application.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |

### list_app_release_workflow_briefs

**Description**: List AppStack release workflow briefs for an application.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |

### get_app_release_workflow_stage

**Description**: Get an AppStack application release workflow stage.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. |
| `releaseStageSn` | string | Yes | Release stage serial number. |

### list_app_release_stage_briefs

**Description**: List AppStack release stage briefs for an application release workflow.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. |

### list_app_release_stage_runs

**Description**: List AppStack release stage execution records.

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. |
| `releaseStageSn` | string | Yes | Release stage serial number. |
| `pagination` | string | No | Pagination mode. Use keyset for keyset pagination. |
| `perPage` | number | No | Page size, up to 100. |
| `orderBy` | string | No | Sort field: id or gmtCreate. |
| `sort` | string | No | Sort direction: asc or desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `page` | number | No | Page number when using page pagination. |

### list_app_release_stage_exec_metadata

**Description**: List integrated change metadata for an AppStack release stage execution.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. |
| `releaseStageSn` | string | Yes | Release stage serial number. |
| `executionNumber` | string | Yes | Release stage execution number. |

### get_app_release_stage_pipeline_run

**Description**: Get the pipeline run for an AppStack release stage execution.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. |
| `releaseStageSn` | string | Yes | Release stage serial number. |
| `executionNumber` | string | Yes | Release stage execution number. |

### get_app_release_stage_pipeline_job_log

**Description**: Get a pipeline job log for an AppStack release stage execution.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. |
| `releaseStageSn` | string | Yes | Release stage serial number. |
| `executionNumber` | string | Yes | Release stage execution number. |
| `jobId` | string | Yes | Pipeline job ID. |

### get_appstack_change_request_audit_items

**Description**: Get audit items for an AppStack change request.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `sn` | string | Yes | Change request serial number. |
| `refType` | string | Yes | Reference type, such as CR. |

### list_appstack_change_request_executions

**Description**: List execution records for an AppStack change request.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `sn` | string | Yes | Change request serial number. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. |
| `releaseStageSn` | string | Yes | Release stage serial number. |
| `perPage` | number | No | Page size, up to 100. |
| `page` | number | No | Page number. |
| `orderBy` | string | No | Sort field: id or gmtCreate. |
| `sort` | string | No | Sort direction: asc or desc. |

### list_appstack_change_request_work_items

**Description**: List work items for an AppStack change request.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `sn` | string | Yes | Change request serial number. |

### list_change_order_versions

**Description**: List AppStack change order versions.

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `envNames` | string | No | Comma-separated environment names. |
| `creators` | string | No | Comma-separated creator account IDs. |
| `current` | number | No | Current page number. |
| `pageSize` | number | No | Page size. |

### get_change_order

**Description**: Get an AppStack change order by serial number.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `changeOrderSn` | string | Yes | Change order serial number. |

### list_change_orders_by_origin

**Description**: List AppStack change orders by creation origin.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `originType` | string | Yes | Origin type, such as FLOW. |
| `originId` | string | Yes | Origin identifier. |
| `appName` | string | No | Application name filter. |
| `envName` | string | No | Environment name filter. |

### list_change_order_job_logs

**Description**: List logs for an AppStack change order job.

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `changeOrderSn` | string | Yes | Change order serial number. |
| `jobSn` | string | Yes | Change order job serial number. |
| `current` | number | No | Current page number. |
| `pageSize` | number | No | Page size. |

### find_task_operation_log

**Description**: Get an AppStack deployment task operation log.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `appName` | string | Yes | Application name. |
| `changeOrderSn` | string | Yes | Change order serial number. |
| `jobSn` | string | Yes | Change order job serial number. |
| `stageSn` | string | Yes | Deployment stage serial number. |
| `taskSn` | string | Yes | Deployment task serial number. |

