# Appstack Enhanced Tools

This document describes the enhanced aggregation tools in the Appstack (application deployment) module. These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations.

## Tool Inventory

| Tool | Purpose | API Calls |
|------|---------|-----------|
| `get_application_overview` | Application dashboard with basic info, environments, and recent orchestrations | 1 + up to 2 |

## Common Behaviors

### Pagination

Sub-sections use `page`/`perPage` parameters with a default page size of 5 (controlled by `envLimit` and `orchestrationLimit`). The raw upstream response is returned for each section; pagination metadata varies by endpoint.

## Tool Details

### get_application_overview

**When to use**: You want a quick snapshot of an Appstack application — its basic info plus environments and recent orchestrations.

**Parameters**:
- `organizationId`, `appName`: required
- `includeEnvironments`: toggle environment list, default true
- `includeOrchestrations`: toggle recent orchestrations, default true
- `envLimit`: max environments returned, default 5
- `orchestrationLimit`: max orchestrations returned, default 5

**Example**:
```json
{
  "appName": "my-app",
  "includeEnvironments": true,
  "includeOrchestrations": true,
  "envLimit": 5,
  "orchestrationLimit": 5
}
```
