# Platform Enhanced Tools

This document describes the enhanced aggregation tools in the Platform (organization and user management) module. These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations.

## Tool Inventory

| Tool | Purpose | API Calls |
|------|---------|-----------|
| `get_organization_overview` | Organization dashboard with basic info, departments, members, groups, and roles | 1 + up to 4 |

## Common Behaviors

### Pagination

Sub-sections use `page`/`perPage` parameters with a default page size of 5 (controlled by `departmentLimit`, `memberLimit`, and `groupLimit`). The `roles` section does not support pagination and returns all roles. The raw upstream response is returned for each section; pagination metadata varies by endpoint.

## Tool Details

### get_organization_overview

**When to use**: You want a quick snapshot of an organization — its basic info plus departments, members, groups, and roles.

**Parameters**:
- `organizationId`: optional, defaults to user's sole organization
- `includeDepartments`, `includeMembers`, `includeGroups`, `includeRoles`: toggle sections, default true
- `departmentLimit`, `memberLimit`, `groupLimit`: control section sizes, default 5

**Example**:
```json
{
  "organizationId": "org-1",
  "includeDepartments": true,
  "includeMembers": true,
  "includeGroups": true,
  "includeRoles": true,
  "departmentLimit": 5,
  "memberLimit": 5,
  "groupLimit": 5
}
```

**Note**: The `roles` section returns all organization roles without pagination limits.
