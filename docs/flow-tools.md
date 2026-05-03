# Flow Tools

This document describes the 14 read-only MCP tools in the flow domain.

## Pagination

Tools in this domain use the following pagination scheme(s):

- Offset (page/perPage)

## Tool Inventory

| Tool | Description |
|------|-------------|
| `get_pipeline_scan_report_url` | Get a temporary download URL for a Flow pipeline scan report. |
| `get_pipeline_artifact_url` | Get a temporary download URL for a Flow pipeline artifact. |
| `get_pipeline_emas_artifact_url` | Get a temporary download URL for a Flow EMAS artifact. |
| `list_pipeline_relations` | List Flow pipeline related objects. |
| `get_last_instance` | Get the latest Flow pipeline run instance detail. |
| `list_resource_members` | List members for a Flow resource such as a pipeline or host group. |
| `list_pipelines` | List Flow pipelines in a Yunxiao organization. |
| `get_pipeline` | Get Flow pipeline details. |
| `list_pipeline_runs` | List Flow pipeline runs. |
| `get_pipeline_run` | Get a Flow pipeline run by ID. |
| `get_latest_pipeline_run` | Get the latest Flow pipeline run. |
| `list_pipeline_jobs_by_category` | List Flow pipeline jobs by task category. |
| `list_pipeline_job_historys` | List Flow pipeline job execution history. |
| `get_pipeline_job_run_log` | Get a Flow pipeline job run log. |

### get_pipeline_scan_report_url

**Description**: Get a temporary download URL for a Flow pipeline scan report.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `reportPath` | string | Yes | Scan report path returned by Flow pipeline APIs. |

### get_pipeline_artifact_url

**Description**: Get a temporary download URL for a Flow pipeline artifact.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `filePath` | string | Yes | Artifact file path returned by Flow pipeline APIs. |
| `fileName` | string | Yes | Artifact file name. |

### get_pipeline_emas_artifact_url

**Description**: Get a temporary download URL for a Flow EMAS artifact.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `emasJobInstanceId` | string | Yes | EMAS job instance ID. |
| `md5` | string | Yes | EMAS artifact MD5. |
| `pipelineId` | string | Yes | Pipeline ID. |
| `pipelineRunId` | string | Yes | Pipeline run ID. |
| `serviceConnectionId` | string | Yes | Service connection ID. |

### list_pipeline_relations

**Description**: List Flow pipeline related objects.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `pipelineId` | string | Yes | Pipeline ID. |
| `relObjectType` | string | Yes | Related object type, such as VARIABLE_GROUP. |

### get_last_instance

**Description**: Get the latest Flow pipeline run instance detail.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `pipelineId` | string | Yes | Pipeline ID. |

### list_resource_members

**Description**: List members for a Flow resource such as a pipeline or host group.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `resourceType` | string | Yes | Resource type, such as pipeline or hostGroup. |
| `resourceId` | string | Yes | Resource ID. |

### list_pipelines

**Description**: List Flow pipelines in a Yunxiao organization.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `createStartTime` | number | No | Pipeline creation start time in milliseconds. |
| `createEndTime` | number | No | Pipeline creation end time in milliseconds. |
| `executeStartTime` | number | No | Pipeline execution start time in milliseconds. |
| `executeEndTime` | number | No | Pipeline execution end time in milliseconds. |
| `pipelineName` | string | No | Pipeline name filter. |
| `statusList` | string | No | Comma-separated statuses such as RUNNING,SUCCESS,FAIL,CANCELED,WAITING. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. Yunxiao supports up to 30. |

### get_pipeline

**Description**: Get Flow pipeline details.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `pipelineId` | string | Yes | Pipeline ID. |

### list_pipeline_runs

**Description**: List Flow pipeline runs.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `pipelineId` | string | Yes | Pipeline ID. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. Yunxiao supports up to 30. |
| `startTime` | number | No | Run start time in milliseconds. |
| `endTime` | number | No | Run end time in milliseconds. |
| `status` | string | No | Run status: FAIL, SUCCESS, or RUNNING. |
| `triggerMode` | number | No | Trigger mode: 1 manual, 2 scheduled, 3 code, 5 pipeline, 6 webhook. |

### get_pipeline_run

**Description**: Get a Flow pipeline run by ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `pipelineId` | string | Yes | Pipeline ID. |
| `pipelineRunId` | string | Yes | Pipeline run ID. |

### get_latest_pipeline_run

**Description**: Get the latest Flow pipeline run.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `pipelineId` | string | Yes | Pipeline ID. |

### list_pipeline_jobs_by_category

**Description**: List Flow pipeline jobs by task category.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `pipelineId` | string | Yes | Pipeline ID. |
| `category` | string | Yes | Task category, currently DEPLOY. |

### list_pipeline_job_historys

**Description**: List Flow pipeline job execution history.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `pipelineId` | string | Yes | Pipeline ID. |
| `category` | string | Yes | Task category, currently DEPLOY. |
| `identifier` | string | Yes | Pipeline job identifier. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. Yunxiao supports up to 30. |

### get_pipeline_job_run_log

**Description**: Get a Flow pipeline job run log.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `pipelineId` | string | Yes | Pipeline ID. |
| `pipelineRunId` | string | Yes | Pipeline run ID. |
| `jobId` | string | Yes | Pipeline job ID. |

