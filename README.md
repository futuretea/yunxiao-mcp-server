# Yunxiao MCP Server

[中文文档](README.zh.md)

Let your AI coding assistant talk directly to [Alibaba Yunxiao](https://www.aliyun.com/product/yunxiao) — browse projects, track iterations, review code, and monitor pipelines without leaving your IDE.

**Read-only by default, safety first.** 177 of 193 tools are read-only queries. 16 write operations require explicit `read_only=false`.

---

## What can you do?

| Scenario | Tools |
|----------|-------|
| 📋 Project Management | Browse projects, work items, iterations, milestones, members; create/update work items |
| 🔍 Code Review | Browse repositories, branches, commits, merge requests, change requests; create/close/merge CRs |
| 🚀 Pipelines | List pipelines, runs, build tasks; approve/reject validation gates |
| 📦 Release Management | Browse applications, environments, releases, change orders, resources |
| 🧠 Knowledge Base | Browse Lingma knowledge bases, members, files |
| 🤖 AI Adoption | Analyze team Lingma usage |

> Write operations (create/update work items, manage CR/MR, pipeline approvals) require `read_only=false`.

---

## Quick Start

### npx (zero install)

```bash
npx -y @futuretea/yunxiao-mcp-server
```

With token:

```bash
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> npx -y @futuretea/yunxiao-mcp-server
```

### Docker

**Stdio MCP mode (default Docker entrypoint):**

```bash
docker run -i --rm -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> \
  ghcr.io/futuretea/yunxiao-mcp-server:latest
```

**HTTP MCP mode:**

```bash
docker run --rm -p 3000:3000 -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> \
  ghcr.io/futuretea/yunxiao-mcp-server:latest --port 3000
```

### Build from source

```bash
make build
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao mcp
```

The same `yunxiao` binary also provides human-facing CLI commands:

```bash
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao organization list
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao member list
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao project list
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao project member list --project-id <project-id>
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao project role list --project-id <project-id>
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao repo list
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao repo branch list --repository-id <repository-id>
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao repo commit list --repository-id <repository-id> --ref <branch>
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao sprint list --project-id <project-id>
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao sprint view <sprint-id> --project-id <project-id>
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao task list --project-id <project-id>
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao task view <workitem-id>
./bin/yunxiao tools describe search_projects
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao tools call get_current_user --params '{}'
printf '{"page":1}' | YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao tools call list_organization_members --params-file -
```

### IDE Setup

See [MCP Client Config](docs/mcp-client-config.md) for Claude, Cursor, and other IDE setup examples.

---

## Tools Overview

| Domain | Tools | Access | Description |
|--------|-------|--------|-------------|
| **Projex** | 47 | 43 read-only + 4 write | Projects, iterations, work items, milestones, test cases |
| **Codeup** | 37 | 31 read-only + 6 write | Repositories, branches, commits, MR, CR, code review |
| **Flow** | 18 | 16 read-only + 2 write | Pipelines, runs, build tasks, validation |
| **Appstack** | 62 | 58 read-only + 4 write | Applications, environments, releases, change orders |
| **Platform** | 18 | read-only | Organizations, departments, members, roles |
| **Packages** | 3 | read-only | Artifact repositories and versions |
| **Lingma** | 6 | read-only | Knowledge bases and usage |
| **API** | 1 | read-only | Generic API fallback |
| **Meta** | 1 | read-only | Tool discovery |

### Enhanced Tools

Enhanced tools aggregate multiple API calls into single operations, reducing AI round-trips.

| Tool | What it combines |
|------|------------------|
| `get_project_overview` | Project info + members + sprints + milestones + labels |
| `get_project_workitem_detail` | Work item + activities + comments + attachments + relations |
| `get_repository_overview` | Repository + default branch + recent commits + recent MRs |
| `get_change_request_overview` | CR detail + patch sets + comments |
| `get_pipeline_overview` | Pipeline info + latest run + history |

Full list: [Enhanced Tools Index](docs/enhanced-tools-index.md)

---

## Configuration

Priority: flags > environment variables > config file > defaults.

### Required

| Variable | Description |
|----------|-------------|
| `YUNXIAO_MCP_ACCESS_TOKEN` | Yunxiao access token |

### Optional

| Variable | Description | Default |
|----------|-------------|---------|
| `YUNXIAO_MCP_BASE_URL` | API base URL | `https://openapi-rdc.aliyuncs.com` |
| `YUNXIAO_MCP_SSE_BASE_URL` | Public SSE base URL (reverse proxy) | — |
| `YUNXIAO_MCP_INSECURE_SKIP_TLS_VERIFY` | Skip TLS verify (internal only) | `false` |

Legacy aliases: `YUNXIAO_ACCESS_TOKEN`, `YUNXIAO_API_BASE_URL`.

### Tool Modes

| Flag | Default | Purpose |
|------|---------|---------|
| `--read-only` | `true` | Set `false` to enable write tools |
| `--compact` | `true` | Hide raw tools with enhanced alternatives (set `false` to show all) |
| `--enabled-tools` | — | Explicit tool allow-list |
| `--disabled-tools` | — | Explicit tool deny-list |
| `--enable-domains` | — | Domain allow-list |
| `--disable-domains` | — | Domain deny-list |

### Config File

```bash
./bin/yunxiao mcp --config config.example.yaml
```

### Per-Request Token (HTTP/SSE)

```bash
curl -H "x-yunxiao-token: <token>" http://localhost:3000/mcp
# or
http://localhost:3000/sse?yunxiao_access_token=<token>
```

---

## HTTP Endpoints

| Endpoint | Purpose |
|----------|---------|
| `/mcp` | Streamable HTTP MCP |
| `/sse` | SSE MCP |
| `/message` | SSE message endpoint |
| `/healthz` | Health check (503 if no tools registered) |

---

## Security

- **Read-only by default**: 177 tools safe for exploration without write access.
- **Explicit write opt-in**: 16 write tools require manual `read_only=false`.
- **Per-request token**: HTTP/SSE support request-level token override for multi-tenant use.
- **No sensitive endpoints**: Admin audit logs, PAT queries, and other high-privilege endpoints are excluded.

---

## Development

```bash
make fmt      # gofmt
make tidy     # go mod tidy
make lint     # go vet + gofmt
make test     # go test ./...
make build    # build the yunxiao binary
make smoke    # smoke test
make ci       # full CI
```

The shared Yunxiao SDK lives in `pkg/yunxiao`; MCP mode and the CLI commands both use it for authenticated OpenAPI requests, path encoding, response metadata, and error classification.

Coverage threshold: 98%. Run `make coverage-check`.

---

## Documentation

- [MCP Client Config](docs/mcp-client-config.md) — IDE setup examples
- [Quick Start Guide](docs/quickstart.md) — AI conversation patterns
- [Enhanced Tools Index](docs/enhanced-tools-index.md) — Enhanced tool reference
- [Conditions Cookbook](docs/conditions-cookbook.md) — Query filter examples
- [Pagination Guide](docs/pagination-guide.md) — Pagination reference
- [GA Readiness](docs/ga-readiness.md) — Release checklist

---

## License

MIT
