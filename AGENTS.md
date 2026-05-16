# Yunxiao MCP Server

Go implementation of an MCP server for the Alibaba Yunxiao (云效) DevOps platform. Defaults to read-only; write tools available with `read_only=false`.

## Quick Start

Build:
```bash
make build
```

Run tests:
```bash
make test
make coverage-check   # verify 98% threshold
```

Run locally (HTTP mode):
```bash
./bin/yunxiao-mcp-server --mode http --port 8080
```

Health check:
```bash
curl http://localhost:8080/healthz
```

## Environment

The server reads configuration from a `.env` file or environment variables:
- `YUNXiao_ACCESS_TOKEN` — default Yunxiao personal access token
- `YUNXiao_ORGANIZATION_ID` — default organization ID (auto-injected when omitted in tool calls)

To test tool calls interactively, use an MCP client (e.g., `mcpc`) pointed at the server.

## Project Structure

```
cmd/yunxiao-mcp-server/    # Entry point
internal/cmd/              # CLI commands (root, serve flags)
pkg/server/http/           # HTTP/SSE transport
pkg/server/mcp/            # MCP server construction, stdio/HTTP/SSE
pkg/toolset/yunxiao/       # All Yunxiao tool definitions and handlers
  *_tools.go               # Tool schema definitions (mcp.NewTool)
  *_handlers.go            # Handler implementations
  *_test.go                # Handler and helper tests
  client.go                # Yunxiao OpenAPI HTTP client
  conditions.go            # Conditions JSON builders
  helpers.go               # Common handler helpers
  toolset.go               # Toolset registration and filtering
```

## Key Files

| File | Purpose |
|------|---------|
| `pkg/toolset/yunxiao/*_tools.go` | 140 MCP tool schemas across 7 domains (130 read-only + 10 write-capable tools gated behind read_only=false) |
| `pkg/toolset/yunxiao/projex_enhanced_*.go` | 8 enhanced aggregation tools |
| `pkg/toolset/yunxiao/codeup_write_*.go` | 6 Codeup write tools (create_change_request, add_change_request_comment, create_merge_request, close_change_request, reopen_change_request, merge_change_request) |
| `pkg/toolset/yunxiao/client.go` | HTTP client with auth, pagination, and error handling |
| `docs/ga-readiness.md` | Release gate checklist and deferred endpoints |
| `docs/quickstart.md` | Common AI conversation patterns |
| `docs/conditions-cookbook.md` | Conditions JSON examples |
| `docs/pagination-guide.md` | Pagination mode reference |

## Adding a New Tool

1. Define the tool in the appropriate `*_tools.go` file using `mcp.NewTool`.
2. Implement the handler in the matching `*_handlers.go` file.
3. Add tests in the matching `*_test.go` file.
4. Run `make ci` and `make coverage-check`.
5. Run `make docs` to regenerate domain documentation.
6. Update the corresponding `docs/<domain>-enhanced-tools.md` and `docs/enhanced-tools-index.md` if the tool is an enhanced aggregation tool.

## Reference Projects

- `third-party-projects/alibabacloud-devops-mcp-server` — Node.js Yunxiao MCP server (useful for API/auth reference)
- `third-party-projects/flashduty-mcp-server` — Go MCP server reference for project structure
