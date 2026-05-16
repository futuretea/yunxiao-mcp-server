# Platform Tools

This document describes the 18 MCP tools in the platform domain.

Access summary: 18 read-only, 0 write-capable.

## Enhanced Tools

These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations. Prefer them when available.

| Tool | Description |
|------|-------------|
| `get_organization_overview` | Get a comprehensive overview of a Yunxiao organization including basic info, departments, members, groups, and roles in one read-only call. |
| `get_organization_department_overview` | Get a comprehensive overview of a Yunxiao organization department including basic info and ancestor chain in one read-only call. |
| `get_organization_group_overview` | Get a comprehensive overview of a Yunxiao organization group including basic info and members in one read-only call. |

## Pagination

Tools in this domain use the following pagination scheme(s):

- Keyset (nextToken)
- Offset (page/perPage)

## Tool Inventory

Tools marked in **bold** are enhanced aggregation tools.

| Tool | Access | Description |
|------|--------|-------------|
| `list_audit_logs` | Read-only | List audit logs in a Yunxiao organization. Use this to track user actions, resource changes, and security events within a time range. |
| **`get_organization_overview`** | Read-only | Get a comprehensive overview of a Yunxiao organization including basic info, departments, members, groups, and roles in one read-only call. |
| **`get_organization_department_overview`** | Read-only | Get a comprehensive overview of a Yunxiao organization department including basic info and ancestor chain in one read-only call. |
| **`get_organization_group_overview`** | Read-only | Get a comprehensive overview of a Yunxiao organization group including basic info and members in one read-only call. |
| `list_enterprise_departments` | Read-only | List enterprise departments visible to the current Yunxiao user. |
| `list_organization_groups` | Read-only | List groups in a Yunxiao organization. Groups are permission-bound collections of users and resources. Use list_organization_members to discover users who can be added to groups. |
| `list_organization_group_members` | Read-only | List members in a Yunxiao organization group. Use this to check who belongs to a specific group and their roles. |
| `get_user` | Read-only | Get a Yunxiao user by ID or username. |
| `list_app_extension_features` | Read-only | List app extension feature implementations for a Yunxiao organization. |
| `get_current_user` | Read-only | Get the current Yunxiao user profile for the configured access token. Use this to verify authentication and discover the user's identity, account ID, and default organization. |
| `get_current_organization_info` | Read-only | Get the current user's default Yunxiao organization context, including organization ID and name. Use this to discover the default organizationId before calling organization-scoped tools. |
| `get_user_organizations` | Read-only | Get Yunxiao organizations visible to the current user. Use this to discover organization IDs and names when the default organization is not the desired one. |
| `list_organizations` | Read-only | List Yunxiao organizations visible to the current user. Use this to discover organization IDs and names when the default organization is not the desired one. |
| `list_organization_departments` | Read-only | List departments in a Yunxiao organization. Use this to discover department IDs for filtering members or assigning work items. |
| `list_organization_members` | Read-only | List members in a Yunxiao organization. Use this to discover user IDs, names, and roles for assigning work items or mentioning in comments. |
| `search_organization_members` | Read-only | Search members in a Yunxiao organization with filters. Use this to find specific users by name, department, or role for assignment or review purposes. |
| `list_organization_roles` | Read-only | List roles defined in a Yunxiao organization. Use this to discover role IDs for filtering members or checking permissions. |
| `list_users` | Read-only | List Yunxiao users across organizations. Use this to discover user IDs and account information for mentions, assignments, or cross-org collaboration. |

### list_audit_logs

**Description**: List audit logs in a Yunxiao organization. Use this to track user actions, resource changes, and security events within a time range.

**Access**: Read-only

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `actionTimeStart` | string | Yes | Inclusive action-time lower bound. Format: RFC3339 timestamp (e.g., 2024-01-01T00:00:00+08:00). |
| `actionTimeEnd` | string | No | Action-time upper bound. Format: RFC3339 timestamp (e.g., 2024-01-31T23:59:59+08:00). Defaults to current time when omitted. |
| `userIds` | string | No | Filter by user IDs. Format: comma-separated numeric user IDs. Use list_users or list_organization_members to discover valid IDs. |
| `apps` | string | No | Filter by application identities. Format: comma-separated application names or identifiers. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `nextToken` | string | No | Pagination token from the previous response x-next-token header. |

### get_organization_overview

**Description**: Get a comprehensive overview of a Yunxiao organization including basic info, departments, members, groups, and roles in one read-only call.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `includeDepartments` | boolean | No | Whether to include departments list. Defaults to true. |
| `includeMembers` | boolean | No | Whether to include members list. Defaults to true. |
| `includeGroups` | boolean | No | Whether to include groups list. Defaults to true. |
| `includeRoles` | boolean | No | Whether to include roles list. Defaults to true. |
| `departmentLimit` | number | No | Maximum departments to include in the overview. Defaults to 5. Set to 0 to exclude. |
| `memberLimit` | number | No | Maximum members to include in the overview. Defaults to 5. Set to 0 to exclude. |
| `groupLimit` | number | No | Maximum groups to include in the overview. Defaults to 5. Set to 0 to exclude. |

### get_organization_department_overview

**Description**: Get a comprehensive overview of a Yunxiao organization department including basic info and ancestor chain in one read-only call.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `departmentId` | string | Yes | Department ID. Use list_organization_departments or list_enterprise_departments to discover valid IDs. |
| `includeAncestors` | boolean | No | Whether to include the ancestor chain. Defaults to true. |

### get_organization_group_overview

**Description**: Get a comprehensive overview of a Yunxiao organization group including basic info and members in one read-only call.

**Access**: Read-only

**Type**: Enhanced aggregation tool

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `groupId` | string | Yes | Group ID. Use list_organization_groups to discover valid group IDs. |
| `includeMembers` | boolean | No | Whether to include group members. Defaults to true. |
| `memberLimit` | number | No | Maximum members to include in the overview. Defaults to 5. Set to 0 to exclude. |

### list_enterprise_departments

**Description**: List enterprise departments visible to the current Yunxiao user.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `parentId` | string | No | Parent department ID. Use list_organization_departments with an empty parentId to discover top-level departments, then drill down by setting this to a discovered department ID. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### list_organization_groups

**Description**: List groups in a Yunxiao organization. Groups are permission-bound collections of users and resources. Use list_organization_members to discover users who can be added to groups.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### list_organization_group_members

**Description**: List members in a Yunxiao organization group. Use this to check who belongs to a specific group and their roles.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `groupId` | string | Yes | Group ID. Use list_organization_groups to discover valid IDs. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### get_user

**Description**: Get a Yunxiao user by ID or username.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `idOrUsername` | string | Yes | User ID (numeric) or login username. Use list_users to discover valid values. |

### list_app_extension_features

**Description**: List app extension feature implementations for a Yunxiao organization.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `type` | string | Yes | App extension type identifier. Contact your organization admin for valid type values. |

### get_current_user

**Description**: Get the current Yunxiao user profile for the configured access token. Use this to verify authentication and discover the user's identity, account ID, and default organization.

**Access**: Read-only

**Parameters**: None

### get_current_organization_info

**Description**: Get the current user's default Yunxiao organization context, including organization ID and name. Use this to discover the default organizationId before calling organization-scoped tools.

**Access**: Read-only

**Parameters**: None

### get_user_organizations

**Description**: Get Yunxiao organizations visible to the current user. Use this to discover organization IDs and names when the default organization is not the desired one.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `page` | number | No | Page number. Defaults to 1 when omitted by Yunxiao. |
| `perPage` | number | No | Page size from 1 to 100. Defaults to 100 when omitted by Yunxiao. |

### list_organizations

**Description**: List Yunxiao organizations visible to the current user. Use this to discover organization IDs and names when the default organization is not the desired one.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `page` | number | No | Page number. Defaults to 1 when omitted by Yunxiao. |
| `perPage` | number | No | Page size from 1 to 100. Defaults to 100 when omitted by Yunxiao. |

### list_organization_departments

**Description**: List departments in a Yunxiao organization. Use this to discover department IDs for filtering members or assigning work items.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `parentId` | string | No | Parent department ID. Use list_organization_departments with an empty parentId to discover top-level departments, then drill down by setting this to a discovered department ID. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. |

### list_organization_members

**Description**: List members in a Yunxiao organization. Use this to discover user IDs, names, and roles for assigning work items or mentioning in comments.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. |

### search_organization_members

**Description**: Search members in a Yunxiao organization with filters. Use this to find specific users by name, department, or role for assignment or review purposes.

**Access**: Read-only

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `deptIds` | array | No | Department IDs to filter by. Use list_organization_departments to discover valid department IDs. |
| `query` | string | No | Member search query. Matches username, display name, or email. |
| `includeChildren` | boolean | No | Whether to include members from child departments. Set to true for broader searches across the org tree. |
| `nextToken` | string | No | Pagination next token from a previous response. |
| `roleIds` | array | No | Role IDs to filter by. Use list_organization_roles to discover valid role IDs. |
| `statuses` | array | No | Member statuses to filter by. Common values: enabled, disabled. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. |

### list_organization_roles

**Description**: List roles defined in a Yunxiao organization. Use this to discover role IDs for filtering members or checking permissions.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |

### list_users

**Description**: List Yunxiao users across organizations. Use this to discover user IDs and account information for mentions, assignments, or cross-org collaboration.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `filter` | string | No | Fuzzy filter for username, login, email, or phone. |
| `status` | string | No | User status filter. Common values: enabled, deleted. |
| `deptId` | string | No | Department ID. Use list_organization_departments to discover valid IDs. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

