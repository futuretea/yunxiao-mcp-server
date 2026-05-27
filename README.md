# yunxiao-mcp-server

MCP (Model Context Protocol) server and CLI for [Alibaba Cloud Yunxiao](https://www.aliyun.com/product/yunxiao) DevOps platform. Lets AI coding assistants (Claude Code, Cursor, etc.) browse projects, track work items, review code, and monitor pipelines directly from your IDE.

## Quick start

```bash
# One-shot via npx (no install)
npx -y @futuretea/yunxiao-mcp-server

# Or install locally
go install github.com/futuretea/yunxiao-mcp-server/cmd/yunxiao@latest
```

You need a Yunxiao access token. Set it via environment variable:

```bash
export YUNXIAO_MCP_ACCESS_TOKEN=<your-token>
```

## Features

- **193 MCP tools** across 9 domains (177 read-only, 16 write)
- **Dual mode**: MCP server for AI assistants + standalone CLI for humans (`yunxiao`)
- **Zero-install**: available via `npx` with automatic platform binary download
- **Multi-transport**: stdio, SSE, and Streamable HTTP

### Tool domains

| Domain | Tools | What you can do |
|--------|-------|-----------------|
| **Projex** | 47 | Projects, work items, sprints, milestones, test cases |
| **Codeup** | 37 | Repositories, branches, commits, MRs, CRs, code review |
| **Flow** | 18 | Pipelines, pipeline runs, build jobs, validation gates |
| **Appstack** | 62 | App stacks, environments, releases, change orders |
| **Platform** | 18 | Organizations, departments, members, roles |
| **Packages** | 3 | Artifact repositories and versions |
| **Lingma** | 6 | Knowledge base and AI usage statistics |
| **API** | 1 | Generic API fallback |
| **Meta** | 1 | Tool discovery |

### Enhanced tools

Some tools aggregate multiple API calls into a single operation to reduce AI round-trips:

- `get_project_overview` — project info + members + sprints + milestones + versions + labels
- `get_project_workitem_detail` — work item + activities + comments + attachments + relations
- `get_repository_overview` — repo + default branch + recent commits + recent MRs
- `get_change_request_overview` — CR details + patch sets + comments
- `get_pipeline_overview` — pipeline info + recent runs + history
- `get_pipeline_run_overview` — pipeline run info + categorized jobs

### Insight tools

- `get_project_risk_dashboard` — risk dashboard with category samples and overdue items
- `get_member_workload_trend` — member workload trends with recent activity
- `get_team_workload_breakdown` — per-member workload breakdown with task details
- `get_sprint_velocity` — historical sprint velocity metrics
- `get_blocker_analysis` — dependency blocker analysis
- `get_workitem_status_timeline` — status change history with dwell time

## MCP client configuration

Add to your MCP client config (e.g., Claude Code's `mcp.json`):

```json
{
  "mcpServers": {
    "yunxiao": {
      "command": "npx",
      "args": ["-y", "@futuretea/yunxiao-mcp-server"],
      "env": {
        "YUNXIAO_MCP_ACCESS_TOKEN": "<your-token>"
      }
    }
  }
}
```

For HTTP mode:

```json
{
  "mcpServers": {
    "yunxiao": {
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

## Configuration

Priority: **flags > environment variables > config file > defaults**.

### Environment variables

| Variable | Description | Default |
|----------|-------------|---------|
| `YUNXIAO_MCP_ACCESS_TOKEN` | Yunxiao access token | (required) |
| `YUNXIAO_MCP_BASE_URL` | API base URL | `https://openapi-rdc.aliyuncs.com` |
| `YUNXIAO_MCP_SSE_BASE_URL` | Public SSE base URL (for reverse proxies) | `""` |
| `YUNXIAO_MCP_INSECURE_SKIP_TLS_VERIFY` | Skip TLS verification | `false` |

### Config file

Create `config.yaml`:

```yaml
port: 0                    # 0 = stdio mode, >0 = HTTP server
sse_base_url: ""           # public URL when behind a reverse proxy
log_level: info            # debug, info, warn, error
base_url: https://openapi-rdc.aliyuncs.com
access_token: ""           # set via env var instead
insecure_skip_tls_verify: false
read_only: true            # deny write tools when true
compact: true              # compact tool descriptions
enabled_tools: []          # whitelist specific tools
disabled_tools: []         # blacklist specific tools
enabled_domains: []        # whitelist specific domains
disabled_domains: []       # blacklist specific domains
request_timeout_seconds: 30
```

### Per-request tokens

For multi-tenant scenarios, pass the token via HTTP header or query parameter:

```bash
# Header
curl -H "x-yunxiao-token: <token>" http://localhost:8080/mcp

# Query parameter
curl "http://localhost:8080/sse?yunxiao_access_token=<token>"
```

## CLI usage

The `yunxiao` binary doubles as a standalone CLI for humans:

```bash
# List projects
yunxiao project list

# Search by name
yunxiao project list --name demo

# Output as JSON
yunxiao project list --json

# View project overview
yunxiao project view 123

# List tasks in a sprint
yunxiao task list --project-id 123 --sprint sprint-456

# View pipeline status
yunxiao pipeline view pipeline-789

# Show your own tasks
yunxiao task my 123
```

### Global flags

```
--config                    path to config file
--base-url                  API base URL
--access-token              Yunxiao access token
--insecure-skip-tls-verify  skip TLS certificate verification
--read-only                 deny write operations (default true)
--output                    output format: table, json, csv
--no-color                  disable ANSI color output
--compact                   compact output
--enabled-tools             comma-separated tool whitelist
--disabled-tools            comma-separated tool blacklist
--enable-domains            comma-separated domain whitelist
--disable-domains           comma-separated domain blacklist
--request-timeout-seconds   API request timeout (default 30)
```

### Command tree

```
yunxiao
├── mcp                 start MCP server
├── version             print version
├── completion          shell completions (bash, zsh, fish, powershell, install)
├── organization        list, view, info
├── member              list, search
├── group               list, view
├── department          list, view
├── pipeline            list, view, run, job, resource-member
├── project             list, view, summary, context, risk, board,
│                       labels, milestones, member, role, member-tasks,
│                       member-trend, templates
├── role                list, view
├── repo                list, view, branch, change-request, mr, commit,
│                       compare, file
├── sprint              list, view, velocity
├── task                list, view, type-list, timeline, my,
│                       type-view, type-all, relation-types
├── testcase            repo-list, search, view, field-config,
│                       directories, plan-list
├── tools               list, describe, call
└── user                whoami, list, get, orgs
```

## Build from source

```bash
git clone https://github.com/futuretea/yunxiao-mcp-server.git
cd yunxiao-mcp-server
make build
```

Requires Go 1.25+.

### Make targets

| Target | Description |
|--------|-------------|
| `build` | Build `bin/yunxiao` binary |
| `test` | Run all tests |
| `lint` | Run go vet + gofmt + golangci-lint + gocyclo |
| `ci` | Full CI pipeline (lint + mod verify + race tests + build) |
| `coverage` | Generate test coverage report |
| `smoke` | Start server and verify `/healthz` |
| `docs` | Generate tool documentation |
| `build-all-platforms` | Cross-compile for darwin/linux/windows × amd64/arm64 |
| `clean` | Remove build artifacts |

## Docker

```bash
# Build
docker build -t yunxiao-mcp-server .

# Run in stdio mode (for MCP clients)
docker run -i --rm -e YUNXIAO_MCP_ACCESS_TOKEN=<token> yunxiao-mcp-server

# Run in HTTP mode
docker run -p 8080:8080 -e YUNXIAO_MCP_ACCESS_TOKEN=<token> \
  yunxiao-mcp-server yunxiao mcp --port 8080
```

Pre-built images: `ghcr.io/futuretea/yunxiao-mcp-server`

## Architecture

```
AI Assistant (IDE) <--MCP JSON-RPC--> yunxiao binary (stdio or HTTP)
                                           |
                                           v
                                    Yunxiao OpenAPI
                              https://openapi-rdc.aliyuncs.com
```

The binary has two entry points that share the same core SDK:

- **MCP server** (`yunxiao mcp`): speaks MCP JSON-RPC over stdio, SSE, or Streamable HTTP. Designed for AI assistants.
- **Standalone CLI** (`yunxiao`): human-friendly output with table/JSON/CSV formatting. Same API surface.

**Security**: write tools are disabled by default (`read_only: true`). Per-request token override supports multi-tenant deployments.

## Documentation

- [Quickstart guide](docs/quickstart.md)
- [MCP client configuration](docs/mcp-client-config.md)
- [Enhanced tools index](docs/enhanced-tools-index.md)
- [Conditions cookbook](docs/conditions-cookbook.md)
- [Pagination guide](docs/pagination-guide.md)
- Per-domain tool references in `docs/*-tools.md`

## License

MIT
