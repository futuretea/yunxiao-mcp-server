# Lingma Tools

This document describes the 4 read-only MCP tools in the lingma domain.

## Pagination

Tools in this domain use the following pagination scheme(s):

- Offset (page/perPage)

## Tool Inventory

| Tool | Description |
|------|-------------|
| `list_knowledge_bases` | List Tongyi Lingma knowledge bases in a Yunxiao organization. Knowledge bases contain curated documents for AI-assisted code completion and chat. |
| `list_kb_files` | List files in a Tongyi Lingma knowledge base. Use list_knowledge_bases to discover valid kbId values. |
| `list_kb_members` | List members with access to a Tongyi Lingma knowledge base. Use list_knowledge_bases to discover valid kbId values. |
| `list_developer_members` | List Tongyi Lingma developer members in a Yunxiao organization. Use this to analyze AI coding assistant adoption across teams. |

### list_knowledge_bases

**Description**: List Tongyi Lingma knowledge bases in a Yunxiao organization. Knowledge bases contain curated documents for AI-assisted code completion and chat.

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

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `departmentId` | string | No | Department ID filter. Use list_organization_departments to discover valid department IDs. |
| `userId` | string | No | User ID filter. Use list_organization_members to discover valid user IDs. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

