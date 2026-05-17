# Packages Tools

This document describes the 3 MCP tools in the packages domain.

Access summary: 3 read-only, 0 write-capable.

## Pagination

Tools in this domain use the following pagination scheme(s):

- Offset (page/perPage)

## Tool Inventory

| Tool | Access | Description |
|------|--------|-------------|
| `list_package_repositories` | Read-only | List artifact repositories (Packages) in a Yunxiao organization. Use this to discover repository IDs for listing artifacts. |
| `list_artifacts` | Read-only | List artifacts in a Packages repository. Requires a repository ID from list_package_repositories. |
| `get_artifact` | Read-only | Get a specific artifact from a Packages repository by ID. Use list_artifacts to discover valid artifact IDs. |

### list_package_repositories

**Description**: List artifact repositories (Packages) in a Yunxiao organization. Use this to discover repository IDs for listing artifacts.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repoTypes` | string | No | Comma-separated repository types: GENERIC, DOCKER, MAVEN, NPM, or NUGET. |
| `repoCategories` | string | No | Comma-separated repository modes: Hybrid, Local, Proxy, ProxyCache, or Group. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |

### list_artifacts

**Description**: List artifacts in a Packages repository. Requires a repository ID from list_package_repositories.

**Access**: Read-only

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repoId` | string | Yes | Packages repository ID. Use list_package_repositories to discover valid IDs. |
| `repoType` | string | Yes | Repository type: GENERIC, DOCKER, MAVEN, NPM, NUGET, or PYPI. |
| `page` | number | No | Page number for pagination. Starts at 1. |
| `perPage` | number | No | Page size for pagination. Supports 1-100. Defaults to 100 when omitted. |
| `search` | string | No | Package name search text. |
| `orderBy` | string | No | Sort field: latestUpdate or gmtDownload. |
| `sort` | string | No | Sort direction: asc or desc. |

### get_artifact

**Description**: Get a specific artifact from a Packages repository by ID. Use list_artifacts to discover valid artifact IDs.

**Access**: Read-only

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. When omitted, the server uses the user's default organization. |
| `repoId` | string | Yes | Packages repository ID. Use list_package_repositories to discover valid IDs. |
| `artifactId` | string | Yes | Artifact ID (string or integer). Use list_artifacts to discover valid IDs. |
| `repoType` | string | Yes | Repository type: GENERIC, DOCKER, MAVEN, NPM, NUGET, or PYPI. |

