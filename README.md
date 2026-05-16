# Yunxiao MCP Server

Go implementation of an MCP server for Alibaba [Yunxiao](https://www.aliyun.com/product/yunxiao) (云效). Exposes a default read-only MCP catalog via stdio, HTTP Streamable, and SSE transports, with ten write-capable tools (4 Projex work item operations, 6 Codeup change request/merge request operations) available only when `read_only=false`.

---

## Quick Start

### npx (Recommended)

No installation required — `npx` downloads the correct platform binary automatically:

```bash
npx -y @futuretea/yunxiao-mcp-server
```

With environment variables:

```bash
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> npx -y @futuretea/yunxiao-mcp-server
```

### Docker

Pre-built images are published to `ghcr.io/futuretea/yunxiao-mcp-server`:

**Stdio mode:**

```bash
docker run -i --rm -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ghcr.io/futuretea/yunxiao-mcp-server:latest
```

**HTTP mode:**

```bash
docker run --rm -p 3000:3000 -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ghcr.io/futuretea/yunxiao-mcp-server:latest --port 3000
```

### Build from source

```bash
make build
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao-mcp-server
```

See [MCP Client Config](docs/mcp-client-config.md) for Claude, Cursor, and other IDE setup examples.

---

## Tools Overview

The default `read_only=true` catalog exposes 130 read-only tools. The full catalog has 140 tools: 130 read-only tools plus ten write-capable tools — four Projex work item operations (`create_workitem`, `update_workitem`, `update_workitem_status`, `add_workitem_comment`) and six Codeup change request/merge request operations (`create_change_request`, `add_change_request_comment`, `create_merge_request`, `close_change_request`, `reopen_change_request`, `merge_change_request`) — when `read_only=false`.

| Domain | Tools | Access | Description |
|--------|-------|--------|-------------|
| **Projex** | 45 | 41 read-only, 4 write-capable | Projects, iterations, work items, milestones, test cases |
| **Codeup** | 30 | 24 read-only, 6 write-capable | Repositories, branches, commits, merge requests, change requests, code review |
| **Flow** | 8 | read-only | Pipelines, runs, build tasks |
| **Appstack** | 31 | read-only | Applications, environments, releases, change orders |
| **Platform** | 18 | read-only | Organizations, departments, members, roles |
| **Packages** | 2 | read-only | Artifact repositories and versions |
| **Lingma** | 4 | read-only | Knowledge bases and usage |
| **API** | 1 | read-only | `call_yunxiao_api` — fallback for supported read-only OpenAPI calls |
| **Meta** | 1 | read-only | `describe_toolset` — discover all available tools |

**Enhanced tools** aggregate multiple API calls into single user-friendly operations. For example, `get_project_overview` returns project info, members, iterations, milestones, and tags in one call.

---

## Configuration

Priority: explicit values > flags > environment > config file > defaults.

### Environment Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `YUNXIAO_MCP_ACCESS_TOKEN` | Yunxiao access token | — |
| `YUNXIAO_MCP_BASE_URL` | API base URL or host | `https://openapi-rdc.aliyuncs.com` |
| `YUNXIAO_MCP_SSE_BASE_URL` | Public SSE base URL (reverse proxy) | — |
| `YUNXIAO_MCP_INSECURE_SKIP_TLS_VERIFY` | Skip Yunxiao server TLS certificate verification for private/self-signed endpoints | `false` |

Legacy aliases `YUNXIAO_ACCESS_TOKEN` and `YUNXIAO_API_BASE_URL` are also supported.

Use `YUNXIAO_MCP_INSECURE_SKIP_TLS_VERIFY=true` or `--insecure-skip-tls-verify` only for trusted internal endpoints where certificate validation cannot be fixed.

### Tool Modes

| Option | Default | Purpose |
|--------|---------|---------|
| `read_only` / `--read-only` | `true` | Excludes write-capable tools. Set `read_only=false` only when Projex work item mutations are intended. |
| `project_focused` / `--project-focused` | `false` | Registers a focused platform + Projex catalog, hiding low-value raw tools with enhanced alternatives. |
| `minimal` / `--minimal` | `false` | Registers the smallest project-centric catalog. Write tools still require `read_only=false`. |
| `enabled_tools` / `--enabled-tools` | `[]` | Explicit tool allow-list by name. |
| `disabled_tools` / `--disabled-tools` | `[]` | Explicit tool deny-list by name. |
| `enabled_domains` / `--enable-domains` | `[]` | Explicit domain allow-list, overriding `project_focused`. |
| `disabled_domains` / `--disable-domains` | `[]` | Explicit domain deny-list. |

### Per-Request Tokens (HTTP/SSE)

Clients can override the default token per request:

```bash
curl -H "x-yunxiao-token: <token>" http://localhost:3000/mcp
# or
http://localhost:3000/sse?yunxiao_access_token=<token>
```

### Config File

```bash
./bin/yunxiao-mcp-server --config config.example.yaml
```

See [config.example.yaml](config.example.yaml) for all options.

---

## HTTP Endpoints

| Endpoint | Purpose |
|----------|---------|
| `/mcp` | Streamable HTTP MCP |
| `/sse` | SSE MCP |
| `/message` | SSE message endpoint |
| `/healthz` | Health check (returns 503 if tools not registered) |

---

## Development

```bash
make format   # gofmt
make tidy     # go mod tidy
make lint     # golangci-lint
make test     # go test ./...
make build    # build binary to bin/
make smoke    # local smoke test (default port 39393)
```

Coverage threshold: 98%. Run `make coverage-check` to verify.

---

## Documentation

- [MCP Client Config](docs/mcp-client-config.md) — IDE integration examples
- [Quick Start Guide](docs/quickstart.md) — Common AI conversation patterns
- [Enhanced Tools Index](docs/enhanced-tools-index.md) — Aggregated tool reference
- [Conditions Cookbook](docs/conditions-cookbook.md) — Query filter examples
- [Pagination Guide](docs/pagination-guide.md) — Pagination mode reference
- [GA Readiness](docs/ga-readiness.md) — Release checklist and deferred endpoints

---

## License

MIT
