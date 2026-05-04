# Yunxiao MCP Server

Go implementation of an MCP server for Alibaba [Yunxiao](https://www.aliyun.com/product/yunxiao) (云效). Exposes 148 read-only tools across 7 domains via stdio, HTTP Streamable, and SSE transports.

---

## Quick Start

**Build and run (stdio mode):**

```bash
make build
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao-mcp-server
```

**HTTP mode:**

```bash
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao-mcp-server --port 3000
```

**Docker:**

```bash
docker build -t yunxiao-mcp-server:local .
docker run -i --rm -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> yunxiao-mcp-server:local
```

See [docs/mcp-client-config.md](docs/mcp-client-config.md) for Claude/Cursor/IDE configuration examples.

---

## Tools Overview

| Domain | Basic Tools | Enhanced Tools | Description |
|--------|------------|----------------|-------------|
| **Projex** | 44 | 9 | Projects, iterations, work items, milestones, test cases |
| **Codeup** | 28 | 5 | Repositories, branches, commits, merge requests, code review |
| **Flow** | 22 | 3 | Pipelines, runs, build tasks |
| **Appstack** | 13 | 2 | Applications, environments, releases, change orders |
| **Platform** | 10 | 2 | Organizations, departments, members, roles |
| **Packages** | 4 | 0 | Artifact repositories and versions |
| **Lingma** | 4 | 1 | Knowledge bases and usage |
| **Meta** | 1 | 0 | `describe_toolset` — discover all available tools |

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

Legacy aliases `YUNXIAO_ACCESS_TOKEN` and `YUNXIAO_API_BASE_URL` are also supported.

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
