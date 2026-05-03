# Platform Tools

This document describes the 23 read-only MCP tools in the platform domain.

## Pagination

Tools in this domain use the following pagination scheme(s):

- Keyset (nextToken)
- Offset (page/perPage)

## Tool Inventory

| Tool | Description |
|------|-------------|
| `list_audit_logs` | List audit logs in a Yunxiao organization. |
| `list_enterprise_departments` | List enterprise departments visible to the current Yunxiao user. |
| `get_enterprise_department` | Get an enterprise department by ID. |
| `list_organization_groups` | List groups in a Yunxiao organization. |
| `get_organization_group` | Get a Yunxiao organization group by ID. |
| `list_organization_group_members` | List members in a Yunxiao organization group. |
| `get_user` | Get a Yunxiao user by ID or username. |
| `list_app_extension_features` | List app extension feature implementations for a Yunxiao organization. |
| `get_current_user` | Get the current Yunxiao user for the configured access token. |
| `get_current_organization_info` | Get current user context, including the last organization returned by Yunxiao. |
| `get_user_organizations` | Get Yunxiao organizations visible to the current user. |
| `list_organizations` | List Yunxiao organizations visible to the current user. |
| `get_organization` | Get a Yunxiao organization by ID. |
| `list_organization_departments` | List departments in a Yunxiao organization. |
| `get_organization_department_info` | Get Yunxiao organization department details. |
| `get_organization_department_ancestors` | List ancestor departments for a Yunxiao organization department. |
| `list_organization_members` | List members in a Yunxiao organization. |
| `get_organization_member_info` | Get Yunxiao organization member details by member ID. |
| `get_organization_member_info_by_user_id` | Get Yunxiao organization member details by user ID. |
| `search_organization_members` | Search members in a Yunxiao organization. |
| `list_organization_roles` | List roles in a Yunxiao organization. |
| `get_organization_role` | Get a Yunxiao organization role by ID. |
| `list_users` | List Yunxiao users. |

### list_audit_logs

**Description**: List audit logs in a Yunxiao organization.

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `actionTimeStart` | string | Yes | Inclusive RFC3339 action-time lower bound. |
| `actionTimeEnd` | string | No | RFC3339 action-time upper bound. Defaults to current time when omitted by Yunxiao. |
| `userIds` | string | No | Comma-separated user IDs. |
| `apps` | string | No | Comma-separated application identities. |
| `perPage` | number | No | Page size from 1 to 100. |
| `nextToken` | string | No | Pagination token from the previous response x-next-token header. |

### list_enterprise_departments

**Description**: List enterprise departments visible to the current Yunxiao user.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `parentId` | string | No | Parent department ID. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |

### get_enterprise_department

**Description**: Get an enterprise department by ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `departmentId` | string | Yes | Enterprise department ID. |

### list_organization_groups

**Description**: List groups in a Yunxiao organization.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |

### get_organization_group

**Description**: Get a Yunxiao organization group by ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `groupId` | string | Yes | Group ID. |

### list_organization_group_members

**Description**: List members in a Yunxiao organization group.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `groupId` | string | Yes | Group ID. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. Default is 100. |

### get_user

**Description**: Get a Yunxiao user by ID or username.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `idOrUsername` | string | Yes | User ID or login username. |

### list_app_extension_features

**Description**: List app extension feature implementations for a Yunxiao organization.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `type` | string | Yes | App extension type. |

### get_current_user

**Description**: Get the current Yunxiao user for the configured access token.

**Parameters**: None

### get_current_organization_info

**Description**: Get current user context, including the last organization returned by Yunxiao.

**Parameters**: None

### get_user_organizations

**Description**: Get Yunxiao organizations visible to the current user.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `page` | number | No | Page number. Defaults to 1 when omitted by Yunxiao. |
| `perPage` | number | No | Page size from 1 to 100. Defaults to 100 when omitted by Yunxiao. |

### list_organizations

**Description**: List Yunxiao organizations visible to the current user.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `page` | number | No | Page number. Defaults to 1 when omitted by Yunxiao. |
| `perPage` | number | No | Page size from 1 to 100. Defaults to 100 when omitted by Yunxiao. |

### get_organization

**Description**: Get a Yunxiao organization by ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Organization ID. Defaults to the user's sole organization when omitted. |

### list_organization_departments

**Description**: List departments in a Yunxiao organization.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `parentId` | string | No | Parent department ID. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |

### get_organization_department_info

**Description**: Get Yunxiao organization department details.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `departmentId` | string | Yes | Department ID. |

### get_organization_department_ancestors

**Description**: List ancestor departments for a Yunxiao organization department.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `departmentId` | string | Yes | Department ID. |

### list_organization_members

**Description**: List members in a Yunxiao organization.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |

### get_organization_member_info

**Description**: Get Yunxiao organization member details by member ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `memberId` | string | Yes | Organization member ID. |

### get_organization_member_info_by_user_id

**Description**: Get Yunxiao organization member details by user ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `userId` | string | Yes | Yunxiao user ID. |

### search_organization_members

**Description**: Search members in a Yunxiao organization.

**Pagination**: Keyset (nextToken)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `deptIds` | array | No | Department IDs. |
| `query` | string | No | Member search query. |
| `includeChildren` | boolean | No | Whether to include child departments. |
| `nextToken` | string | No | Pagination next token. |
| `roleIds` | array | No | Role IDs. |
| `statuses` | array | No | Member statuses. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size from 1 to 100. |

### list_organization_roles

**Description**: List roles in a Yunxiao organization.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |

### get_organization_role

**Description**: Get a Yunxiao organization role by ID.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `roleId` | string | Yes | Organization role ID. |

### list_users

**Description**: List Yunxiao users.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `filter` | string | No | Fuzzy filter for username, login, email, or phone. |
| `status` | string | No | User status: enabled or deleted. |
| `deptId` | string | No | Department ID. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |

