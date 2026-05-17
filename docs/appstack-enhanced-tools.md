# Appstack Enhanced Tools

This document describes the enhanced aggregation tools in the Appstack (application deployment) module. These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations.

## Tool Inventory

| Tool | Purpose | API Calls |
|------|---------|-----------|
| `get_application_overview` | Application dashboard with basic info, environments, and recent orchestrations | 1 + up to 2 |
| `get_environment_overview` | Environment dashboard with basic info, variable groups, and latest orchestration | 1 + up to 2 |
| `get_release_overview` | System release dashboard with members, products, and attached change requests | 1 + up to 3 |
| `get_system_overview` | System dashboard with attached apps and members | 1 + up to 2 |
| `get_change_order_overview` | Change order dashboard with basic info and job list | 1 + 1 |
| `get_app_release_workflow_overview` | Release workflow dashboard with stage briefs | 1 + 1 |
| `get_app_release_stage_overview` | Stage execution dashboard with pipeline run and integrated metadata | 1 + up to 2 |

## Tool Details

### get_application_overview

**When to use**: Quick snapshot of an Appstack application — basic info plus environments and recent orchestrations.

**Parameters**: `organizationId`, `appName` (required), `includeEnvironments`, `includeOrchestrations`, `envLimit`, `orchestrationLimit`

### get_environment_overview

**When to use**: Quick snapshot of an Appstack environment — basic info plus variable groups and latest orchestration.

**Parameters**: `organizationId`, `appName`, `envName` (required), `includeVariableGroups`, `includeLatestOrchestration`

### get_release_overview

**When to use**: Quick snapshot of an Appstack system release — basic info plus members, products, and attached change requests.

**Parameters**: `organizationId`, `systemName`, `sn` (required), `includeMembers`, `includeProducts`, `includeChangeRequests`, `changeRequestLimit`

### get_system_overview

**When to use**: Quick snapshot of an Appstack system — basic info plus attached applications and members.

**Parameters**: `organizationId`, `systemName` (required), `includeApps`, `includeMembers`, `appLimit`, `memberLimit`

### get_change_order_overview

**When to use**: Quick snapshot of an Appstack change order — basic info plus job list.

**Parameters**: `organizationId`, `appName`, `changeOrderSn` (required), `includeJobLogs`

### get_app_release_workflow_overview

**When to use**: Quick snapshot of an AppStack application release workflow — basic info plus stage briefs for discovery.

**Parameters**: `organizationId`, `appName`, `releaseWorkflowSn` (required), `includeStageBriefs`

### get_app_release_stage_overview

**When to use**: Debug a specific stage execution — stage info, pipeline run, and integrated metadata in one call.

**Parameters**: `organizationId`, `appName`, `releaseWorkflowSn`, `releaseStageSn`, `executionNumber` (required), `includeStageInfo`, `includePipelineRun`, `includeMetadata`
