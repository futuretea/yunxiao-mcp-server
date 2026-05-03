# Flow Enhanced Tools

This document describes the enhanced aggregation tools in the Flow (CI/CD pipeline) module. These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations.

## Tool Inventory

| Tool | Purpose | API Calls |
|------|---------|-----------|
| `get_pipeline_overview` | Pipeline dashboard with basic info, latest run, and recent run history | 2 + up to 1 |
| `get_pipeline_run_overview` | Pipeline run details with pipeline jobs by category | 1 + up to 1 |

## Common Behaviors

### Pagination

Run history uses `page`/`perPage` parameters with a default page size of 5 (controlled by `runLimit`). The raw upstream response is returned; pagination metadata varies by endpoint.

## Tool Details

### get_pipeline_overview

**When to use**: You want a quick snapshot of a pipeline — its configuration plus the latest run status and recent run history.

**Parameters**:
- `organizationId`, `pipelineId`: required
- `includeRuns`: toggle recent run history, default true
- `runLimit`: max recent runs returned, default 5

**Example**:
```json
{
  "pipelineId": "pipeline-1",
  "includeRuns": true,
  "runLimit": 10
}
```

**Note**: The `latestRun` section is always included because it uses a dedicated lightweight endpoint.

### get_pipeline_run_overview

**When to use**: You want a quick snapshot of a specific pipeline run — its details plus the pipeline jobs that provide context for understanding the run.

**Parameters**:
- `organizationId`, `pipelineId`, `pipelineRunId`: required
- `includeJobs`: toggle pipeline jobs list, default true
- `category`: task category for job listing, default `DEPLOY`

**Example**:
```json
{
  "pipelineId": "pipeline-1",
  "pipelineRunId": "run-1",
  "includeJobs": true,
  "category": "DEPLOY"
}
```

**Note**: Pipeline jobs are fetched from the pipeline scope (not the run scope) because Yunxiao does not provide a run-scoped job listing endpoint. The jobs section shows the pipeline's job structure for context.
