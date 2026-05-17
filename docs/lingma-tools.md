# Lingma Tools

This document describes the 6 MCP tools in the lingma domain.

Access summary: 6 read-only, 0 write-capable.

## Pagination

Tools in this domain use the following pagination scheme(s):

- Offset (page/perPage)

## Tool Inventory

| Tool | Access | Description |
|------|--------|-------------|
| `list_knowledge_bases` | Read-only | List Tongyi Lingma knowledge bases in a Yunxiao organization. Knowledge bases contain curated documents for AI-assisted code completion and chat. |
| `list_kb_files` | Read-only | List files in a Tongyi Lingma knowledge base. Use list_knowledge_bases to discover valid kbId values. |
| `list_kb_members` | Read-only | List members with access to a Tongyi Lingma knowledge base. Use list_knowledge_bases to discover valid kbId values. |
| `list_developer_members` | Read-only | List Tongyi Lingma developer members in a Yunxiao organization. Use this to analyze AI coding assistant adoption across teams. |
| `get_department_usage` | Read-only | Get Tongyi Lingma usage metrics for a specific department over a time range. Use this to analyze AI coding assistant adoption within a department. |
| `get_developer_usage` | Read-only | Get Tongyi Lingma usage metrics for a specific developer or department over a time range. Use this to analyze individual or team AI coding assistant usage. |

### list_knowledge_bases

**Description**: List Tongyi Lingma knowledge bases in a Yunxiao organization. Knowledge bases contain curated documents for AI-assisted code completion and chat.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `query` | string | No | Knowledge base name fuzzy query. |
| `sceneType` | string | No | Scene type, such as chat or completion. |
| `orderBy` | string | No | Sort field. |
| `sort` | string | No | Sort order: desc or asc. |
| `userId` | string | No | User ID permission filter. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 20 when omitted. |

### list_kb_files

**Description**: List files in a Tongyi Lingma knowledge base. Use list_knowledge_bases to discover valid kbId values.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `kbId` | string | Yes | Knowledge base ID. Use list_knowledge_bases to discover valid IDs. |
| `query` | string | No | File name fuzzy query. |
| `orderBy` | string | No | Sort field. |
| `sort` | string | No | Sort order: desc or asc. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 20 when omitted. |

### list_kb_members

**Description**: List members with access to a Tongyi Lingma knowledge base. Use list_knowledge_bases to discover valid kbId values.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `kbId` | string | Yes | Knowledge base ID. Use list_knowledge_bases to discover valid IDs. |
| `query` | string | No | Member name fuzzy query. |
| `orderBy` | string | No | Sort field. |
| `sort` | string | No | Sort order: desc or asc. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 20 when omitted. |

### list_developer_members

**Description**: List Tongyi Lingma developer members in a Yunxiao organization. Use this to analyze AI coding assistant adoption across teams.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `departmentId` | string | No | Department ID filter. Use list_organization_departments to discover valid department IDs. |
| `userId` | string | No | User ID filter. Use list_organization_members to discover valid user IDs. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### get_department_usage

**Description**: Get Tongyi Lingma usage metrics for a specific department over a time range. Use this to analyze AI coding assistant adoption within a department.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `departmentId` | string | Yes | Department ID. Use list_organization_departments to discover valid department IDs. |
| `startTime` | string | Yes | Start time for the usage query range (inclusive). Format: yyyy-MM-ddTHH:mm:ss+08:00 (e.g. 2024-01-01T00:00:00+08:00). |
| `endTime` | string | Yes | End time for the usage query range (inclusive). Format: yyyy-MM-ddTHH:mm:ss+08:00. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### get_developer_usage

**Description**: Get Tongyi Lingma usage metrics for a specific developer or department over a time range. Use this to analyze individual or team AI coding assistant usage.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `startTime` | string | Yes | Start time for the usage query range (inclusive). Format: yyyy-MM-ddTHH:mm:ss+08:00 (e.g. 2024-01-01T00:00:00+08:00). |
| `endTime` | string | Yes | End time for the usage query range (inclusive). Format: yyyy-MM-ddTHH:mm:ss+08:00. |
| `userId` | string | No | User ID to query usage for a specific developer. Use list_organization_members to discover valid user IDs. Either userId or departmentId must be provided. |
| `departmentId` | string | No | Department ID to query usage for all developers in a department. Use list_organization_departments to discover valid department IDs. Either userId or departmentId must be provided. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

