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

## Docker

For stdio clients, keep stdin open:

```bash
docker run -i --rm -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> yunxiao-mcp-server:local
```

For HTTP clients, publish the port and pass `--port`:

```bash
docker run --rm -p 3000:3000 -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> yunxiao-mcp-server:local --port 3000
```
