# Flow Tools

This document describes the 10 MCP tools in the flow domain.

Access summary: 8 read-only, 2 write-capable.

## Enhanced Tools

These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations. Prefer them when available.

| Tool | Description |
|------|-------------|
| `get_pipeline_overview` | Get a comprehensive overview of a Flow pipeline including basic info, latest run, and recent run history in one read-only call. |
| `get_pipeline_run_overview` | Get a comprehensive overview of a Flow pipeline run including run details and pipeline jobs by category in one read-only call. |

## Pagination

Tools in this domain use the following pagination scheme(s):

- Offset (page/perPage)

## Tool Inventory

Tools marked in **bold** are enhanced aggregation tools.

| Tool | Access | Description |
|------|--------|-------------|
| `list_pipeline_relations` | Read-only | List objects related to a Flow pipeline, such as variable groups. Use this to discover pipeline dependencies and linked resources. |
| **`get_pipeline_overview`** | Read-only | Get a comprehensive overview of a Flow pipeline including basic info, latest run, and recent run history in one read-only call. |
| **`get_pipeline_run_overview`** | Read-only | Get a comprehensive overview of a Flow pipeline run including run details and pipeline jobs by category in one read-only call. |
| `list_resource_members` | Read-only | List members who have access to a Flow resource (e.g., a pipeline or host group). Use this to discover who can manage or trigger a pipeline. |
| `list_pipelines` | Read-only | List Flow CI/CD pipelines in a Yunxiao organization. Use this to discover pipelines and obtain their IDs before calling pipeline-scoped tools. For a comprehensive view of a single pipeline including latest run and history, use get_pipeline_overview instead. |
| `list_pipeline_runs` | Read-only | List execution runs for a Flow pipeline. Use this to review historical runs and their statuses. For the latest run only, use get_latest_pipeline_run. |
| `list_pipeline_jobs_by_category` | Read-only | List jobs (tasks) within a Flow pipeline grouped by category. Use this after identifying a pipeline to see its build, deploy, and test stages. |
| `list_pipeline_job_historys` | Read-only | List execution history for a specific Flow pipeline job. Use this to track how a particular job (e.g., a deploy step) has performed across multiple runs. |
| `pass_pipeline_validate` | Write-capable | Pass (approve) a pipeline validation job. |
| `refuse_pipeline_validate` | Write-capable | Refuse (reject) a pipeline validation job. |

### list_pipeline_relations

**Description**: List objects related to a Flow pipeline, such as variable groups. Use this to discover pipeline dependencies and linked resources.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `pipelineId` | string | Yes | Pipeline ID (string). Use list_pipelines to find the pipeline ID. |
| `relObjectType` | string | Yes | Related object type. Example: VARIABLE_GROUP. |

### get_pipeline_overview

**Description**: Get a comprehensive overview of a Flow pipeline including basic info, latest run, and recent run history in one read-only call.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `pipelineId` | string | Yes | Pipeline ID (string). Use list_pipelines to find the pipeline ID. |
| `includeRuns` | boolean | No | Whether to include recent run history. Defaults to true. |
| `runLimit` | number | No | Max recent runs returned. Defaults to 5. |

### get_pipeline_run_overview

**Description**: Get a comprehensive overview of a Flow pipeline run including run details and pipeline jobs by category in one read-only call.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `pipelineId` | string | Yes | Pipeline ID (string). Use list_pipelines to find the pipeline ID. |
| `pipelineRunId` | string | Yes | Pipeline run ID (string). Use list_pipeline_runs to find the run ID. |
| `includeJobs` | boolean | No | Whether to include pipeline jobs by category. Defaults to true. |
| `category` | string | No | Task category for job listing. Common value: DEPLOY. Use list_pipeline_jobs_by_category to discover available categories. |

### list_resource_members

**Description**: List members who have access to a Flow resource (e.g., a pipeline or host group). Use this to discover who can manage or trigger a pipeline.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `resourceType` | string | Yes | Resource type. Examples: pipeline, hostGroup. |
| `resourceId` | string | Yes | Resource ID (string). Use list_pipelines or other list tools to find the resource ID. |

### list_pipelines

**Description**: List Flow CI/CD pipelines in a Yunxiao organization. Use this to discover pipelines and obtain their IDs before calling pipeline-scoped tools. For a comprehensive view of a single pipeline including latest run and history, use get_pipeline_overview instead.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `createStartTime` | number | No | Pipeline creation start time as a Unix timestamp in milliseconds (e.g., 1704067200000). |
| `createEndTime` | number | No | Pipeline creation end time as a Unix timestamp in milliseconds. |
| `executeStartTime` | number | No | Pipeline execution start time as a Unix timestamp in milliseconds. |
| `executeEndTime` | number | No | Pipeline execution end time as a Unix timestamp in milliseconds. |
| `pipelineName` | string | No | Filter by pipeline name (contains match). |
| `statusList` | string | No | Comma-separated pipeline statuses. Common values: RUNNING, SUCCESS, FAIL, CANCELED, WAITING. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Yunxiao supports up to 30. |

### list_pipeline_runs

**Description**: List execution runs for a Flow pipeline. Use this to review historical runs and their statuses. For the latest run only, use get_latest_pipeline_run.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `pipelineId` | string | Yes | Pipeline ID (string). Use list_pipelines to find the pipeline ID. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Yunxiao supports up to 30. |
| `startTime` | number | No | Run start time as a Unix timestamp in milliseconds. |
| `endTime` | number | No | Run end time as a Unix timestamp in milliseconds. |
| `status` | string | No | Filter by run status. Common values: FAIL, SUCCESS, RUNNING. |
| `triggerMode` | number | No | Filter by trigger mode: 1 manual, 2 scheduled, 3 code push, 5 pipeline, 6 webhook. |

### list_pipeline_jobs_by_category

**Description**: List jobs (tasks) within a Flow pipeline grouped by category. Use this after identifying a pipeline to see its build, deploy, and test stages.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `pipelineId` | string | Yes | Pipeline ID. Use list_pipelines to discover valid IDs. |
| `category` | string | Yes | Task category. Common value: DEPLOY (for deployment tasks). |

### list_pipeline_job_historys

**Description**: List execution history for a specific Flow pipeline job. Use this to track how a particular job (e.g., a deploy step) has performed across multiple runs.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `pipelineId` | string | Yes | Pipeline ID. Use list_pipelines to discover valid IDs. |
| `category` | string | Yes | Task category. Common value: DEPLOY (for deployment tasks). |
| `identifier` | string | Yes | Pipeline job identifier (string). Use list_pipeline_jobs_by_category to discover job identifiers. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Yunxiao supports up to 30. |

### pass_pipeline_validate

**Description**: Pass (approve) a pipeline validation job.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | Yes | Yunxiao organization ID. |
| `pipelineId` | string | Yes | Pipeline ID. |
| `pipelineRunId` | string | Yes | Pipeline run ID. |
| `jobId` | string | Yes | Validation job ID. |

### refuse_pipeline_validate

**Description**: Refuse (reject) a pipeline validation job.

**Access**: Write-capable (requires `read_only=false`)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | Yes | Yunxiao organization ID. |
| `pipelineId` | string | Yes | Pipeline ID. |
| `pipelineRunId` | string | Yes | Pipeline run ID. |
| `jobId` | string | Yes | Validation job ID. |

