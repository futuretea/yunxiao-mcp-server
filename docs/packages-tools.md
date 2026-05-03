# Packages Tools

This document describes the 3 read-only MCP tools in the packages domain.

## Pagination

Tools in this domain use the following pagination scheme(s):

- Offset (page/perPage)

## Tool Inventory

| Tool | Description |
|------|-------------|
| `list_package_repositories` | List Packages repositories in a Yunxiao organization. |
| `list_artifacts` | List artifacts in a Packages repository. |
| `get_artifact` | Get one artifact from a Packages repository. |

### list_package_repositories

**Description**: List Packages repositories in a Yunxiao organization.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repoTypes` | string | No | Comma-separated repository types: GENERIC, DOCKER, MAVEN, NPM, or NUGET. |
| `repoCategories` | string | No | Comma-separated repository modes: Hybrid, Local, Proxy, ProxyCache, or Group. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |

### list_artifacts

**Description**: List artifacts in a Packages repository.

**Pagination**: Offset (page/perPage)

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repoId` | string | Yes | Packages repository ID. |
| `repoType` | string | Yes | Repository type: GENERIC, DOCKER, MAVEN, NPM, NUGET, or PYPI. |
| `page` | number | No | Page number. |
| `perPage` | number | No | Page size. |
| `search` | string | No | Package name search text. |
| `orderBy` | string | No | Sort field: latestUpdate or gmtDownload. |
| `sort` | string | No | Sort direction: asc or desc. |

### get_artifact

**Description**: Get one artifact from a Packages repository.

**Parameters**:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `organizationId` | string | No | Yunxiao organization ID. Defaults to the user's sole organization when omitted. |
| `repoId` | string | Yes | Packages repository ID. |
| `artifactId` | number | Yes | Artifact ID. |
| `repoType` | string | Yes | Repository type: GENERIC, DOCKER, MAVEN, NPM, NUGET, or PYPI. |

