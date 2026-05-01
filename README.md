# yunxiao-mcp-server

Go 语言版本的 Yunxiao MCP Server。当前实现提供可构建的 stdio MCP 服务骨架，并接入云效 OpenAPI 的基础只读工具。

## 功能

- stdio MCP transport
- HTTP Streamable MCP transport 与 SSE transport
- HTTP health endpoint：`/healthz`
- Cobra CLI 与 YAML/env/flag 配置加载
- Yunxiao OpenAPI token 认证：`x-yunxiao-token`
- 工具启用/禁用过滤
- 基础只读工具：
  - `get_current_user`
  - `get_current_organization_info`
  - `get_user_organizations`
  - `list_organizations`
  - `get_organization`
  - `list_repositories`
  - `get_repository`
  - `list_branches`
  - `list_pipelines`
  - `get_pipeline`
  - `list_pipeline_runs`
  - `get_pipeline_run`
  - `get_latest_pipeline_run`

## 快速开始

```bash
make build
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao-mcp-server
```

默认使用 stdio 模式，适合 MCP client 直接拉起。也可以通过 `--config config.example.yaml` 使用配置文件。

HTTP 模式：

```bash
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao-mcp-server --port 3000
```

HTTP endpoints：

- `/mcp`：Streamable HTTP MCP endpoint
- `/sse`：SSE MCP endpoint
- `/message`：SSE message endpoint
- `/healthz`：readiness check；缺少 access token 或未注册工具时返回 `503`

## 配置

配置优先级为 programmatic explicit values > flag > environment > config file > defaults。普通运行只会用到 flag 及其之后的层级；programmatic explicit values 用于测试或内嵌调用。

常用环境变量：

- `YUNXIAO_MCP_ACCESS_TOKEN`：云效 access token
- `YUNXIAO_MCP_BASE_URL`：云效 OpenAPI host 或 API base URL，默认 `https://openapi-rdc.aliyuncs.com`
- `YUNXIAO_MCP_SSE_BASE_URL`：反向代理场景下 SSE message endpoint 的 public base URL
- `YUNXIAO_ACCESS_TOKEN`：兼容 Node 参考实现的 token 变量
- `YUNXIAO_API_BASE_URL`：兼容 Node 参考实现的 base URL 变量

若同时设置新旧环境变量，`YUNXIAO_MCP_*` 优先于兼容用的 legacy 变量。

`base_url` 可以是主域名，也可以已经包含 `/oapi/v1`。客户端会避免重复追加 `/oapi/v1`。

## 开发

```bash
make format
make tidy
make test
make build
```
