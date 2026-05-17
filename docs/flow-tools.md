# Flow Tools

This document describes the 18 MCP tools in the flow domain.

Access summary: 16 read-only, 2 write-capable.

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
| `get_pipeline_scan_report_url` | Read-only | Get the scan report URL for a Flow pipeline. Use this to retrieve security scan or code quality report links. |
| `get_pipeline_artifact_url` | Read-only | Get the artifact download URL for a Flow pipeline build output. Use this to retrieve build artifact download links. |
| `get_pipeline_emas_artifact_url` | Read-only | Get the EMAS (Enterprise Mobile Application Studio) artifact download URL for a Flow pipeline. Use this to retrieve mobile app build artifact links. |
| `get_last_instance` | Read-only | Get the last service connection instance for a Flow pipeline. Use this to retrieve the most recent service connection configuration. |
| **`get_pipeline_overview`** | Read-only | Get a comprehensive overview of a Flow pipeline including basic info, latest run, and recent run history in one read-only call. |
| **`get_pipeline_run_overview`** | Read-only | Get a comprehensive overview of a Flow pipeline run including run details and pipeline jobs by category in one read-only call. |
| `list_resource_members` | Read-only | List members who have access to a Flow resource (e.g., a pipeline or host group). Use this to discover who can manage or trigger a pipeline. |
| `list_pipelines` | Read-only | List Flow CI/CD pipelines in a Yunxiao organization. Use this to discover pipelines and obtain their IDs before calling pipeline-scoped tools. For a comprehensive view of a single pipeline including latest run and history, use get_pipeline_overview instead. |
| `get_pipeline` | Read-only | Get a single Flow pipeline by ID. Use list_pipelines to discover valid pipeline IDs. For a comprehensive view with latest run info, use get_pipeline_overview instead. |
| `list_pipeline_runs` | Read-only | List execution runs for a Flow pipeline. Use this to review historical runs and their statuses. For the latest run only, use get_latest_pipeline_run. |
| `get_latest_pipeline_run` | Read-only | Get the latest execution run for a Flow pipeline. Use this for a quick status check without listing all historical runs. |
| `get_pipeline_run` | Read-only | Get a specific Flow pipeline run by ID. Use list_pipeline_runs to discover valid run IDs. For a comprehensive view with metadata, use get_pipeline_run_overview instead. |
| `list_pipeline_jobs_by_category` | Read-only | List jobs (tasks) within a Flow pipeline grouped by category. Use this after identifying a pipeline to see its build, deploy, and test stages. |
| `list_pipeline_job_historys` | Read-only | List execution history for a specific Flow pipeline job. Use this to track how a particular job (e.g., a deploy step) has performed across multiple runs. |
| `get_pipeline_job_run_log` | Read-only | Get the execution log for a specific job within a Flow pipeline run. Use this to debug pipeline failures by inspecting individual job output. |
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

### get_pipeline_scan_report_url

**Description**: Get the scan report URL for a Flow pipeline. Use this to retrieve security scan or code quality report links.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `reportPath` | string | Yes | Report path provided by the pipeline scan step output. |

### get_pipeline_artifact_url

**Description**: Get the artifact download URL for a Flow pipeline build output. Use this to retrieve build artifact download links.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `filePath` | string | Yes | Artifact file path. Typically returned by pipeline build step output. |
| `fileName` | string | Yes | Artifact file name. Typically returned by pipeline build step output. |

### get_pipeline_emas_artifact_url

**Description**: Get the EMAS (Enterprise Mobile Application Studio) artifact download URL for a Flow pipeline. Use this to retrieve mobile app build artifact links.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `emasJobInstanceId` | string | Yes | EMAS job instance ID. Discovered from pipeline EMAS step output. |
| `md5` | string | Yes | MD5 checksum of the EMAS artifact. Discovered from pipeline EMAS step output. |
| `pipelineId` | string | Yes | Pipeline ID (integer or string). Use list_pipelines to find the pipeline ID. |
| `pipelineRunId` | string | Yes | Pipeline run ID (integer or string). Use list_pipeline_runs to discover valid run IDs. |
| `serviceConnectionId` | string | Yes | Service connection ID (integer or string). Typically discovered from pipeline configuration. |

### get_last_instance

**Description**: Get the last service connection instance for a Flow pipeline. Use this to retrieve the most recent service connection configuration.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `pipelineId` | string | Yes | Pipeline ID (string). Use list_pipelines to find the pipeline ID. |

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

### get_pipeline

**Description**: Get a single Flow pipeline by ID. Use list_pipelines to discover valid pipeline IDs. For a comprehensive view with latest run info, use get_pipeline_overview instead.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `pipelineId` | string | Yes | Pipeline ID (string). Use list_pipelines to discover valid IDs. |

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

### get_latest_pipeline_run

**Description**: Get the latest execution run for a Flow pipeline. Use this for a quick status check without listing all historical runs.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `pipelineId` | string | Yes | Pipeline ID (string). Use list_pipelines to find the pipeline ID. |

### get_pipeline_run

**Description**: Get a specific Flow pipeline run by ID. Use list_pipeline_runs to discover valid run IDs. For a comprehensive view with metadata, use get_pipeline_run_overview instead.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `pipelineId` | string | Yes | Pipeline ID (string). Use list_pipelines to find the pipeline ID. |
| `pipelineRunId` | string | Yes | Pipeline run ID. Use list_pipeline_runs to discover valid run IDs. |

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

### get_pipeline_job_run_log

**Description**: Get the execution log for a specific job within a Flow pipeline run. Use this to debug pipeline failures by inspecting individual job output.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `pipelineId` | string | Yes | Pipeline ID (string). Use list_pipelines to find the pipeline ID. |
| `pipelineRunId` | string | Yes | Pipeline run ID. Use list_pipeline_runs to discover valid run IDs. |
| `jobId` | string | Yes | Job ID within the pipeline run. Use list_pipeline_jobs_by_category to discover valid job IDs. |

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

