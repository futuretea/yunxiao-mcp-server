# Lingma Tools

This document describes the 6 read-only MCP tools in the lingma domain.

## Pagination

Tools in this domain use the following pagination scheme(s):

- Offset (page/perPage)

## Tool Inventory

| Tool | Description |
|------|-------------|
| `list_knowledge_bases` | List Tongyi Lingma knowledge bases. |
| `list_kb_files` | List Tongyi Lingma knowledge base files. |
| `list_kb_members` | List Tongyi Lingma knowledge base members. |
| `get_department_usage` | Get Tongyi Lingma daily usage data for a department. |
| `list_developer_members` | List Tongyi Lingma developer members. |
| `get_developer_usage` | Get Tongyi Lingma daily usage data for developers. |

### list_knowledge_bases

**Description**: List Tongyi Lingma knowledge bases.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `query` | string | No | Knowledge base name fuzzy query. |
| `sceneType` | string | No | Scene type, such as chat or completion. |
| `orderBy` | string | No | Sort field. |
| `sort` | string | No | Sort order: desc or asc. |
| `userId` | string | No | User ID permission filter. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. Default is 20. |

### list_kb_files

**Description**: List Tongyi Lingma knowledge base files.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `kbId` | string | Yes | Knowledge base ID. |
| `query` | string | No | File name fuzzy query. |
| `orderBy` | string | No | Sort field. |
| `sort` | string | No | Sort order: desc or asc. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. Default is 20. |

### list_kb_members

**Description**: List Tongyi Lingma knowledge base members.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `kbId` | string | Yes | Knowledge base ID. |
| `query` | string | No | Member name fuzzy query. |
| `orderBy` | string | No | Sort field. |
| `sort` | string | No | Sort order: desc or asc. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. Default is 20. |

### get_department_usage

**Description**: Get Tongyi Lingma daily usage data for a department.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `departmentId` | string | Yes | Department ID. |
| `startTime` | string | Yes | Start date in YYYY-MM-DD format. |
| `endTime` | string | Yes | End date in YYYY-MM-DD format. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. Default is 100. |

### list_developer_members

**Description**: List Tongyi Lingma developer members.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `departmentId` | string | No | Department ID filter. |
| `userId` | string | No | User ID filter. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. Default is 100. |

### get_developer_usage

**Description**: Get Tongyi Lingma daily usage data for developers.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `userId` | string | No | User ID. Either userId or departmentId is required. |
| `departmentId` | string | No | Department ID. Either userId or departmentId is required. |
| `startTime` | string | Yes | Start date in YYYY-MM-DD format. |
| `endTime` | string | Yes | End date in YYYY-MM-DD format. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. Default is 100. |

