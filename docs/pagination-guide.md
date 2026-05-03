# Pagination Guide

Yunxiao APIs use three pagination schemes. This guide explains which scheme each tool family uses and how to paginate correctly.

## Pagination Modes

### Mode A: `page` / `perPage` (Offset Pagination)

The most common scheme. Pass a 1-based page number and a page size.

**Used by**: Projex, Codeup, Flow, Packages, Platform, Lingma, and most Appstack tools.

**Example tools**:
- `search_projects`
- `search_workitems`
- `list_sprints`
- `list_project_members`
- `list_artifacts`
- `list_audit_logs`

**Request**:
```json
{
  "page": 2,
  "perPage": 20
}
```

**Response metadata**: Some endpoints return `x-total` and `x-total-pages` headers. Enhanced tools surface this as a `pagination` object:

```json
{
  "data": [...],
  "pagination": {
    "total": 100,
    "page": 2,
    "perPage": 20,
    "totalPages": 5
  }
}
```

Base tools return the raw upstream response; pagination metadata varies by endpoint.

### Mode B: `current` / `pageSize` (Offset Pagination)

A variation used by a subset of Appstack endpoints. Behavior is the same as Mode A, only the parameter names differ.

**Used by**: Some Appstack system, global variable, and release endpoints.

**Example tools**:
- `list_appstack_systems`
- `list_appstack_global_variables`
- `list_appstack_system_releases`
- `list_appstack_environments`

**Request**:
```json
{
  "current": 2,
  "pageSize": 20
}
```

### Mode C: `nextToken` (+ `page` / `perPage`) (Keyset Pagination)

Keyset pagination uses a continuation token from the previous response. This is more efficient for large datasets than offset pagination because it avoids skipped or duplicated rows when the underlying data changes between requests.

Some tools support **both** keyset and offset pagination; the `page` parameter is a fallback when `nextToken` is omitted.

**Used by**: Appstack application metadata, deployment resources, release search, release workflows, and Platform audit logs.

**Example tools**:
- `list_appstack_applications`
- `list_appstack_application_templates`
- `list_appstack_deployments`
- `search_appstack_releases`
- `list_audit_logs`

**First request** (omit `nextToken`):
```json
{
  "perPage": 20
}
```

**Subsequent request**:
```json
{
  "nextToken": "eyJwYWdlIjoyfQ==",
  "perPage": 20
}
```

**Response**: The raw response includes a `nextToken` field when more pages are available. If `nextToken` is absent or empty, you have reached the last page.

## Quick Reference by Domain

| Domain | Mode | Parameters |
|--------|------|------------|
| Projex | A | `page` / `perPage` |
| Codeup | A | `page` / `perPage` |
| Flow | A | `page` / `perPage` |
| Packages | A | `page` / `perPage` |
| Platform | A / C | `page` / `perPage` (most); `nextToken` (audit) |
| Lingma | A | `page` / `perPage` |
| Appstack | A / B / C | `page` / `perPage` (most); `current` / `pageSize` (system/global var); `nextToken` (app metadata, deployments, releases) |

## Pagination Tips

- **Default sizes**: Most tools do not enforce a default page size in the MCP schema. If you omit `perPage` or `pageSize`, the Yunxiao API chooses its own default (usually 10 or 20). For predictable behavior, always specify the page size.
- **Maximum sizes**: Some Appstack endpoints document a maximum `perPage` of 100. Other endpoints may silently truncate larger values.
- **Enhanced tools**: Aggregation tools such as `get_project_overview` and `get_project_workitem_detail` use `page` / `perPage` for their list sub-sections and include `pagination` metadata when the upstream API provides it.
- **When to stop paginating**:
  - Mode A: stop when the returned `data` array is empty or when `page` exceeds `totalPages`.
  - Mode C: stop when `nextToken` is absent or empty.
- **Do not mix modes**: Passing both `page` and `nextToken` to the same request may produce undefined behavior depending on the endpoint.
