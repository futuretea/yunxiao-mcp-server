# MCP Client Configuration

## Stdio

```json
{
  "mcpServers": {
    "yunxiao": {
      "command": "/path/to/yunxiao-mcp-server",
      "env": {
        "YUNXIAO_MCP_ACCESS_TOKEN": "<your-token>"
      }
    }
  }
}
```

## HTTP

```bash
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> yunxiao-mcp-server --port 3000
```

Use these endpoints:

- Streamable HTTP: `http://127.0.0.1:3000/mcp`
- SSE: `http://127.0.0.1:3000/sse`
- Readiness: `http://127.0.0.1:3000/healthz`

When running behind a reverse proxy, set `YUNXIAO_MCP_SSE_BASE_URL` or `--sse-base-url` so SSE clients receive the public message URL.

For a shared HTTP service, clients can pass their own token per request:

```text
x-yunxiao-token: <your-token>
```

SSE clients can also put the token on the connection URL:

```text
http://127.0.0.1:3000/sse?yunxiao_access_token=<your-token>
```

The SSE message endpoint sent back to the client keeps that query token, so later JSON-RPC messages use the same token. Request tokens take precedence over `YUNXIAO_MCP_ACCESS_TOKEN`.

## Docker

For stdio clients, keep stdin open:

```bash
docker run -i --rm -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> yunxiao-mcp-server:local
```

For HTTP clients, publish the port and pass `--port`:

```bash
docker run --rm -p 3000:3000 -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> yunxiao-mcp-server:local --port 3000
```
