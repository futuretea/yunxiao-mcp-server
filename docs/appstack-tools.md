# Appstack Tools

This document describes the 62 MCP tools in the appstack domain.

Access summary: 58 read-only, 4 write-capable.

## Enhanced Tools

These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations. Prefer them when available.

| Tool | Description |
|------|-------------|
| `get_application_overview` | Get a comprehensive overview of an Appstack application including basic info, environments, and recent orchestrations in one read-only call. Use this after discovering applications via list_applications or list_attached_apps. |
| `get_environment_overview` | Get a comprehensive overview of an Appstack environment including basic info, variable groups, and latest orchestration in one read-only call. Use this after identifying an application and environment via get_application_overview. |
| `get_release_overview` | Get a comprehensive overview of an Appstack system release including basic info, members, products, and attached change requests in one read-only call. Use this after discovering a release via search_releases or list_system_release_workflows. |
| `get_system_overview` | Get a comprehensive overview of an Appstack system including basic info, attached applications, and members in one read-only call. Use this after discovering a system via list_systems. |
| `get_change_order_overview` | Get a comprehensive overview of an Appstack change order including basic info and job list in one read-only call. Use this after discovering a change order via list_change_order_versions. |
| `get_app_release_workflow_overview` | Get a comprehensive overview of an AppStack application release workflow including workflow info and stage briefs in one read-only call. Use this after discovering a workflow via list_app_release_workflows. |
| `get_app_release_stage_overview` | Get a comprehensive overview of an AppStack application release stage execution including stage info, pipeline run, and integrated metadata in one read-only call. Use this after discovering executions via list_app_release_stage_runs. |

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
| `get_environment` | Read-only | Get a single AppStack environment by name. Use list_environments or get_application_overview to discover valid environment names. |
| `list_application_members` | Read-only | List members of an AppStack application. Use list_applications to discover valid application names. |
| `list_application_sources` | Read-only | List source repositories attached to an AppStack application. Use list_applications to discover valid application names. |
| `get_machine_deploy_log` | Read-only | Get deployment log for a specific machine in an AppStack deployment. Machine logs capture the agent-side output of a deployment. |
| `get_deploy_group` | Read-only | Get an AppStack deploy group by name within a resource pool. Deploy groups define subsets of machines for targeted deployments. |
| `list_resource_instances` | Read-only | List AppStack resource instances in a resource pool. Pool names are typically found in application environment resource configurations. |
| `get_resource_instance` | Read-only | Get an AppStack resource instance by name within a resource pool. Resource instances represent individual machines or hosts. |
| **`get_application_overview`** | Read-only | Get a comprehensive overview of an Appstack application including basic info, environments, and recent orchestrations in one read-only call. Use this after discovering applications via list_applications or list_attached_apps. |
| **`get_environment_overview`** | Read-only | Get a comprehensive overview of an Appstack environment including basic info, variable groups, and latest orchestration in one read-only call. Use this after identifying an application and environment via get_application_overview. |
| **`get_release_overview`** | Read-only | Get a comprehensive overview of an Appstack system release including basic info, members, products, and attached change requests in one read-only call. Use this after discovering a release via search_releases or list_system_release_workflows. |
| **`get_system_overview`** | Read-only | Get a comprehensive overview of an Appstack system including basic info, attached applications, and members in one read-only call. Use this after discovering a system via list_systems. |
| **`get_change_order_overview`** | Read-only | Get a comprehensive overview of an Appstack change order including basic info and job list in one read-only call. Use this after discovering a change order via list_change_order_versions. |
| **`get_app_release_workflow_overview`** | Read-only | Get a comprehensive overview of an AppStack application release workflow including workflow info and stage briefs in one read-only call. Use this after discovering a workflow via list_app_release_workflows. |
| **`get_app_release_stage_overview`** | Read-only | Get a comprehensive overview of an AppStack application release stage execution including stage info, pipeline run, and integrated metadata in one read-only call. Use this after discovering executions via list_app_release_stage_runs. |
| `list_global_vars` | Read-only | Search AppStack global variable groups. Use this to discover variable group IDs before reading or updating specific groups. |
| `get_global_var` | Read-only | Get an AppStack global variable by name. Use list_global_vars to discover valid global variable names. |
| `search_releases` | Read-only | Search AppStack releases in a Yunxiao organization. Use this to discover releases before calling get_release_overview or other release-specific tools. |
| `get_pod_container_log` | Read-only | Get container logs from a pod in an AppStack resource proxy. Use this to retrieve recent logs for debugging deployment issues. |
| `get_pod_info` | Read-only | Get detailed information about a pod in an AppStack resource proxy. Use this to inspect pod status, containers, and events. |
| `get_kubernetes_object_info` | Read-only | Get detailed information about a Kubernetes object in an AppStack resource proxy. Supports pods, deployments, services, and other Kubernetes resources. |
| `get_deployment_revision_info` | Read-only | Get revision information for an AppStack deployment. Use this to inspect the rollout history and revision details of a deployment. |
| `list_system_release_workflows` | Read-only | List AppStack release workflows for a system. Use this after discovering a system via list_systems to find releases and their workflows. |
| `get_release` | Read-only | Get an AppStack system release by serial number. Use list_system_release_workflows or search_releases to discover valid release serial numbers. |
| `list_release_members` | Read-only | List members of an AppStack system release. Use this after discovering a release via search_releases or list_system_release_workflows. |
| `list_release_products` | Read-only | List products attached to an AppStack system release. Use this after discovering a release via search_releases or list_system_release_workflows. |
| `list_attached_change_requests` | Read-only | List change requests attached to an AppStack system release. Use this after discovering a release via search_releases or list_system_release_workflows. |
| `list_release_executions` | Read-only | List execution records for an AppStack system release. Use this after discovering release workflow and stage details via list_system_release_workflows or get_release_overview. |
| `list_systems` | Read-only | List AppStack systems in a Yunxiao organization. Use this as the entry point to discover systems before calling other system-specific tools. |
| `list_attached_apps` | Read-only | List applications attached to an AppStack system. Use this to discover applications within a system before calling get_application_overview. |
| `list_system_members` | Read-only | List members of an AppStack system. Use this after discovering a system via list_systems. |
| `search_app_tags` | Read-only | Search AppStack application tags in a Yunxiao organization. Application tags are used to categorize and label applications for organizational purposes. |
| `list_applications` | Read-only | List AppStack applications in a Yunxiao organization. AppStack manages deployment environments and release pipelines. |
| `get_application` | Read-only | Get a single AppStack application by name. Use list_applications to discover valid application names. |
| `get_env_variable_groups` | Read-only | Get AppStack variable groups for a specific environment. Variable groups define environment-specific configuration values. |
| `get_variable_group` | Read-only | Get an AppStack variable group by name. Variable groups contain key-value configuration pairs for deployment environments. |
| `get_app_variable_groups` | Read-only | List AppStack variable groups for an application. Use this to discover variable group names before reading specific groups. |
| `get_app_variable_groups_revision` | Read-only | Get the revision metadata for AppStack variable groups of an application. Use this to check if variable groups have been updated. |
| `get_latest_orchestration` | Read-only | Get the latest available AppStack orchestration for an application environment. Orchestrations define the deployment configuration and resource layout. |
| `list_app_orchestration` | Read-only | List AppStack orchestrations for an application. Orchestrations define deployment workflows and environment configurations. |
| `get_app_orchestration` | Read-only | Get an AppStack orchestration by serial number. Use list_app_orchestration to discover valid serial numbers. |
| `list_app_release_workflows` | Read-only | List AppStack release workflows for an application. Release workflows model multi-stage deployment pipelines. |
| `list_app_release_workflow_briefs` | Read-only | List AppStack release workflow briefs for an application. Briefs provide a condensed view of workflow definitions. |
| `get_app_release_workflow_stage` | Read-only | Get an AppStack release workflow stage by serial number. Stages represent individual deployment phases with configuration details. |
| `list_app_release_stage_briefs` | Read-only | List AppStack release stage briefs for an application release workflow. Stages represent individual deployment phases within a workflow. |
| `list_app_release_stage_runs` | Read-only | List AppStack release stage execution records. Each record represents a single run of a deployment stage. |
| `get_app_release_stage_pipeline_run` | Read-only | Get the Flow pipeline run associated with an AppStack release stage execution. Pipeline runs show CI/CD pipeline execution details. |
| `get_app_release_stage_pipeline_job_log` | Read-only | Get the pipeline job log for an AppStack release stage execution. Job logs contain detailed CI/CD pipeline output for a specific job. |
| `list_app_release_stage_exec_metadata` | Read-only | List integrated change metadata for an AppStack release stage execution. Metadata includes linked work items and commits. |
| `get_appstack_change_request_audit_items` | Read-only | Get audit items for an AppStack change request. Audit items represent approval checkpoints that must be satisfied before deployment. |
| `list_appstack_change_request_executions` | Read-only | List execution records for an AppStack change request. Change requests track planned deployments. |
| `list_appstack_change_request_work_items` | Read-only | List work items for an AppStack change request. Work items represent linked Projex tasks or requirements. |
| `get_change_order` | Read-only | Get an AppStack change order by serial number. Change orders track actual deployments with full details. |
| `list_change_order_versions` | Read-only | List AppStack change order versions. Change orders track actual deployments and their versions. |
| `list_change_orders_by_origin` | Read-only | List AppStack change orders by creation origin. Use this to trace deployments back to their source (e.g., a Flow pipeline). |
| `list_change_order_job_logs` | Read-only | List logs for an AppStack change order job. Job logs capture deployment script output. |
| `find_task_operation_log` | Read-only | Get an AppStack deployment task operation log. Operation logs record manual or automated actions taken during deployment. |
| `create_change_order` | Write-capable | Create an AppStack change order (deployment order). Change orders trigger application deployments to environments. This is a write operation and requires read_only=false. |
| `execute_job_action` | Write-capable | Execute an action on an AppStack change order job. Use this to suspend, resume, rollback, or stop a deployment job. This is a write operation and requires read_only=false. |
| `execute_system_release_stage` | Write-capable | Execute a system release stage. Triggers the deployment pipeline for a release stage. This is a write operation and requires read_only=false. |
| `execute_app_release_stage` | Write-capable | Execute an application release stage. Triggers the deployment pipeline for an app-level release stage. This is a write operation and requires read_only=false. |

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

### get_environment

**Description**: Get a single AppStack environment by name. Use list_environments or get_application_overview to discover valid environment names.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid names. |
| `envName` | string | Yes | Environment name. Use list_environments to discover valid names. |

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

### get_machine_deploy_log

**Description**: Get deployment log for a specific machine in an AppStack deployment. Machine logs capture the agent-side output of a deployment.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `tunnelId` | string | Yes | Deployment tunnel ID. Typically discovered from change order or deployment details. |
| `machineSn` | string | Yes | Machine serial number. Use deployment details to discover valid machine identifiers. |

### get_deploy_group

**Description**: Get an AppStack deploy group by name within a resource pool. Deploy groups define subsets of machines for targeted deployments.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `poolName` | string | Yes | Resource pool name. Typically discovered from application environment details. |
| `deployGroupName` | string | Yes | Deploy group name. Typically discovered from deployment configuration. |

### list_resource_instances

**Description**: List AppStack resource instances in a resource pool. Pool names are typically found in application environment resource configurations.

**Access**: Read-only

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `poolName` | string | Yes | Resource pool name. Typically found in application environment resource configurations. |
| `pagination` | string | No | Pagination mode. Valid value: keyset. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `orderBy` | string | No | Sort field. Valid values: id, gmtCreate. |
| `sort` | string | No | Sort direction. Valid values: asc, desc. |
| `nextToken` | string | No | Keyset pagination token from the previous response. |
| `page` | number | No | Page number for pagination. Starts at 1. |

### get_resource_instance

**Description**: Get an AppStack resource instance by name within a resource pool. Resource instances represent individual machines or hosts.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `poolName` | string | Yes | Resource pool name. Use list_resource_instances to discover valid pool names. |
| `instanceName` | string | Yes | Resource instance name. Use list_resource_instances to discover valid instance names. |

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

### get_system_overview

**Description**: Get a comprehensive overview of an Appstack system including basic info, attached applications, and members in one read-only call. Use this after discovering a system via list_systems.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `systemName` | string | Yes | System name. Use list_systems to discover valid names. |
| `includeApps` | boolean | No | Whether to include attached applications. Defaults to true. |
| `includeMembers` | boolean | No | Whether to include system members. Defaults to true. |
| `appLimit` | number | No | Max attached applications returned. Defaults to 10. |
| `memberLimit` | number | No | Max members returned. Defaults to 10. |

### get_change_order_overview

**Description**: Get a comprehensive overview of an Appstack change order including basic info and job list in one read-only call. Use this after discovering a change order via list_change_order_versions.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `changeOrderSn` | string | Yes | Change order serial number. Use list_change_order_versions to discover valid values. |
| `includeJobLogs` | boolean | No | Whether to include job list. Defaults to true. |

### get_app_release_workflow_overview

**Description**: Get a comprehensive overview of an AppStack application release workflow including workflow info and stage briefs in one read-only call. Use this after discovering a workflow via list_app_release_workflows.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. Use list_app_release_workflows to discover valid values. |
| `includeStageBriefs` | boolean | No | Whether to include stage briefs. Defaults to true. |

### get_app_release_stage_overview

**Description**: Get a comprehensive overview of an AppStack application release stage execution including stage info, pipeline run, and integrated metadata in one read-only call. Use this after discovering executions via list_app_release_stage_runs.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. Use list_app_release_workflows to discover valid values. |
| `releaseStageSn` | string | Yes | Release stage serial number. Use list_app_release_stage_runs to discover valid values. |
| `executionNumber` | string | Yes | Release stage execution number. Use list_app_release_stage_runs to discover valid values. |
| `includeStageInfo` | boolean | No | Whether to include stage details. Defaults to true. |
| `includePipelineRun` | boolean | No | Whether to include pipeline run info. Defaults to true. |
| `includeMetadata` | boolean | No | Whether to include integrated metadata (linked work items, commits). Defaults to true. |

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

### get_global_var

**Description**: Get an AppStack global variable by name. Use list_global_vars to discover valid global variable names.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `name` | string | Yes | Global variable name. Use list_global_vars to discover valid names. |
| `revisionSha` | string | No | Optional revision SHA to retrieve a specific version of the variable. |

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

### get_pod_container_log

**Description**: Get container logs from a pod in an AppStack resource proxy. Use this to retrieve recent logs for debugging deployment issues.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `resourcePath` | string | Yes | Resource proxy path. Typically discovered from application environment resource configurations. |
| `namespace` | string | Yes | Kubernetes namespace. Use get_pod_info or get_kubernetes_object_info to discover valid namespaces. |
| `name` | string | Yes | Pod name. Use get_kubernetes_object_info to discover valid pod names. |
| `container` | string | Yes | Container name within the pod. Typically discovered from pod info. |
| `tailingLines` | number | No | Number of recent log lines to return. Defaults to 1000. |

### get_pod_info

**Description**: Get detailed information about a pod in an AppStack resource proxy. Use this to inspect pod status, containers, and events.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `resourcePath` | string | Yes | Resource proxy path. Typically discovered from application environment resource configurations. |
| `namespace` | string | Yes | Kubernetes namespace. Use get_kubernetes_object_info to discover valid namespaces. |
| `name` | string | Yes | Pod name. Use get_kubernetes_object_info to discover valid pod names. |
| `taskSn` | string | No | Optional deployment task serial number for filtering pod info by task. |

### get_kubernetes_object_info

**Description**: Get detailed information about a Kubernetes object in an AppStack resource proxy. Supports pods, deployments, services, and other Kubernetes resources.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `resourcePath` | string | Yes | Resource proxy path. Typically discovered from application environment resource configurations. |
| `namespace` | string | Yes | Kubernetes namespace. Use get_application_overview or get_environment_overview to discover valid namespaces. |
| `name` | string | Yes | Kubernetes object name. Use Kubernetes conventions to identify objects by name. |
| `kind` | string | Yes | Kubernetes object kind. Valid values: Pod, Deployment, Service, ConfigMap, Secret, Ingress, StatefulSet, DaemonSet, etc. |
| `taskSn` | string | No | Optional deployment task serial number for filtering object info by task. |

### get_deployment_revision_info

**Description**: Get revision information for an AppStack deployment. Use this to inspect the rollout history and revision details of a deployment.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `envName` | string | Yes | Environment name. Use list_environments to discover valid environment names. |
| `namespace` | string | Yes | Kubernetes namespace. Use get_application_overview or get_environment_overview to discover valid namespaces. |
| `name` | string | Yes | Deployment name. Use get_kubernetes_object_info with kind=Deployment to discover valid deployment names. |
| `revision` | string | Yes | Deployment revision number. Use Kubernetes rollout history to discover valid revision numbers. |
| `taskSn` | string | No | Optional deployment task serial number for filtering revision info by task. |

### list_system_release_workflows

**Description**: List AppStack release workflows for a system. Use this after discovering a system via list_systems to find releases and their workflows.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `systemName` | string | Yes | System name. Use list_systems to discover valid names. |

### get_release

**Description**: Get an AppStack system release by serial number. Use list_system_release_workflows or search_releases to discover valid release serial numbers.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `systemName` | string | Yes | System name. Use list_systems to discover valid names. |
| `sn` | string | Yes | Release serial number. Use search_releases or list_system_release_workflows to discover valid serial numbers. |

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

### search_app_tags

**Description**: Search AppStack application tags in a Yunxiao organization. Application tags are used to categorize and label applications for organizational purposes.

**Access**: Read-only

**Pagination**: Offset (current/pageSize)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `current` | number | No | Page number for pagination. Starts at 1. |
| `pageSize` | number | No | Page size for pagination. Supports 1-100. Defaults to 10 when omitted. |
| `orderBy` | string | No | Sort field. Valid values: tagName, id. Defaults to id. |
| `sort` | string | No | Sort direction. Valid values: asc, desc. Defaults to desc. |
| `search` | string | No | Optional search keyword for fuzzy matching against tag names. |

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

### get_application

**Description**: Get a single AppStack application by name. Use list_applications to discover valid application names.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid names. |

### get_env_variable_groups

**Description**: Get AppStack variable groups for a specific environment. Variable groups define environment-specific configuration values.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `envName` | string | Yes | Environment name. Use list_environments to discover valid environment names. |

### get_variable_group

**Description**: Get an AppStack variable group by name. Variable groups contain key-value configuration pairs for deployment environments.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `variableGroupName` | string | Yes | Variable group name. Use get_app_variable_groups to discover valid names. |

### get_app_variable_groups

**Description**: List AppStack variable groups for an application. Use this to discover variable group names before reading specific groups.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |

### get_app_variable_groups_revision

**Description**: Get the revision metadata for AppStack variable groups of an application. Use this to check if variable groups have been updated.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |

### get_latest_orchestration

**Description**: Get the latest available AppStack orchestration for an application environment. Orchestrations define the deployment configuration and resource layout.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `envName` | string | Yes | Environment name. Use list_environments to discover valid environment names. |

### list_app_orchestration

**Description**: List AppStack orchestrations for an application. Orchestrations define deployment workflows and environment configurations.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |

### get_app_orchestration

**Description**: Get an AppStack orchestration by serial number. Use list_app_orchestration to discover valid serial numbers.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `sn` | string | Yes | Orchestration serial number. Use list_app_orchestration to discover valid serial numbers. |
| `tagName` | string | No | Optional tag name to retrieve a tagged version of the orchestration. |
| `sha` | string | No | Optional SHA to retrieve a specific version of the orchestration. |

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

### get_app_release_workflow_stage

**Description**: Get an AppStack release workflow stage by serial number. Stages represent individual deployment phases with configuration details.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. Use list_app_release_workflows to discover valid values. |
| `releaseStageSn` | string | Yes | Release stage serial number. Use list_app_release_stage_briefs to discover valid values. |

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

### get_app_release_stage_pipeline_run

**Description**: Get the Flow pipeline run associated with an AppStack release stage execution. Pipeline runs show CI/CD pipeline execution details.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. Use list_app_release_workflows to discover valid values. |
| `releaseStageSn` | string | Yes | Release stage serial number. Use list_app_release_stage_briefs to discover valid values. |
| `executionNumber` | string | Yes | Release stage execution number. Use list_app_release_stage_runs to discover valid values. |

### get_app_release_stage_pipeline_job_log

**Description**: Get the pipeline job log for an AppStack release stage execution. Job logs contain detailed CI/CD pipeline output for a specific job.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. Use list_app_release_workflows to discover valid values. |
| `releaseStageSn` | string | Yes | Release stage serial number. Use list_app_release_stage_briefs to discover valid values. |
| `executionNumber` | string | Yes | Release stage execution number. Use list_app_release_stage_runs to discover valid values. |
| `jobId` | string | Yes | Pipeline job ID. Typically discovered from the pipeline run details. |

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

### get_appstack_change_request_audit_items

**Description**: Get audit items for an AppStack change request. Audit items represent approval checkpoints that must be satisfied before deployment.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `sn` | string | Yes | Change request serial number. Typically discovered via list_attached_change_requests. |
| `refType` | string | Yes | Reference type for audit items. Valid values: RELEASE, CHANGE_REQUEST. |

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

### get_change_order

**Description**: Get an AppStack change order by serial number. Change orders track actual deployments with full details.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `changeOrderSn` | string | Yes | Change order serial number. Use list_change_order_versions to discover valid values. |

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

### create_change_order

**Description**: Create an AppStack change order (deployment order). Change orders trigger application deployments to environments. This is a write operation and requires read_only=false.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `changeOrder` | string | Yes | JSON string with change order details: {changeOrderName, type (Deploy|Scale|Rollback|Destroy), envs (object), orchestrationRevisionSha, description}. |

### execute_job_action

**Description**: Execute an action on an AppStack change order job. Use this to suspend, resume, rollback, or stop a deployment job. This is a write operation and requires read_only=false.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `changeOrderSn` | string | Yes | Change order serial number. Use list_change_order_versions to discover valid values. |
| `jobSn` | string | Yes | Job serial number. Typically returned in change order job details. |
| `action` | string | Yes | JSON string with job action: {actionType: SUSPEND|RESUME|ROLLBACK|STOP}. |

### execute_system_release_stage

**Description**: Execute a system release stage. Triggers the deployment pipeline for a release stage. This is a write operation and requires read_only=false.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `systemName` | string | Yes | System name. Use list_systems to discover valid names. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. Use list_system_release_workflows to discover valid values. |
| `releaseStageSn` | string | Yes | Release stage serial number. Use list_system_release_workflows to discover valid values. |
| `execution` | string | Yes | JSON string with execution parameters: {appReleaseSn, params (object of key-value pairs)}. |

### execute_app_release_stage

**Description**: Execute an application release stage. Triggers the deployment pipeline for an app-level release stage. This is a write operation and requires read_only=false.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `appName` | string | Yes | Application name. Use list_applications to discover valid app names. |
| `releaseWorkflowSn` | string | Yes | Release workflow serial number. Use list_app_release_workflows to discover valid values. |
| `releaseStageSn` | string | Yes | Release stage serial number. Use list_app_release_stage_briefs to discover valid values. |
| `execution` | string | Yes | JSON string with execution parameters: {appReleaseSn, params (object of key-value pairs)}. |

