# MCP Client Configuration

## Stdio

### npx (Recommended)

No installation required — `npx` downloads the correct platform binary automatically.

**Claude Desktop** (`claude_desktop_config.json`):

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

**Cursor** (Settings → MCP):

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

### Local Binary

If you built from source or downloaded a release binary:

```json
{
  "mcpServers": {
    "yunxiao": {
      "command": "/path/to/yunxiao",
      "args": ["mcp"],
      "env": {
        "YUNXIAO_MCP_ACCESS_TOKEN": "<your-token>"
      }
    }
  }
}
```

### Docker

For stdio clients, keep stdin open with `-i`:

```json
{
  "mcpServers": {
    "yunxiao": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "YUNXIAO_MCP_ACCESS_TOKEN=<your-token>",
        "ghcr.io/futuretea/yunxiao-mcp-server:latest"
      ]
    }
  }
}
```

---

## HTTP

Start the server in HTTP mode:

```bash
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> yunxiao mcp --port 3000
```

Or with Docker:

```bash
docker run --rm -p 3000:3000 -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ghcr.io/futuretea/yunxiao-mcp-server:latest --port 3000
```

Use these endpoints:

- Streamable HTTP: `http://127.0.0.1:3000/mcp`
- SSE: `http://127.0.0.1:3000/sse`
- Readiness: `http://127.0.0.1:3000/healthz`

When running behind a reverse proxy, set `YUNXIAO_MCP_SSE_BASE_URL` or `--sse-base-url` so SSE clients receive the public message URL.

For trusted internal Yunxiao-compatible endpoints with private or self-signed certificates, set `YUNXIAO_MCP_INSECURE_SKIP_TLS_VERIFY=true` or pass `--insecure-skip-tls-verify`. Keep this disabled for public endpoints.

For a shared HTTP service, clients can pass their own token per request:

```text
x-yunxiao-token: <your-token>
```

SSE clients can also put the token on the connection URL:

```text
http://127.0.0.1:3000/sse?yunxiao_access_token=<your-token>
```

The SSE message endpoint sent back to the client keeps that query token, so later JSON-RPC messages use the same token. Request tokens take precedence over `YUNXIAO_MCP_ACCESS_TOKEN`.
