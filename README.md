# yunxiao-mcp-server

Go 语言版本的 Yunxiao MCP Server。当前实现提供可构建的 stdio MCP 服务骨架，并接入云效 OpenAPI 的基础只读工具。

## 功能

- stdio MCP transport
- HTTP Streamable MCP transport 与 SSE transport
- HTTP health endpoint：`/healthz`
- Cobra CLI 与 YAML/env/flag 配置加载
- Yunxiao OpenAPI token 认证：启动时默认 token，HTTP/SSE 请求级 `x-yunxiao-token` 或 `yunxiao_access_token` 覆盖
- **191 个只读 MCP 工具**，覆盖 7 个领域：
  - **Projex**（项目管理）：项目、迭代、工作项、里程碑、测试用例等
  - **Codeup**（代码托管）：仓库、分支、提交、合并请求、代码评审等
  - **Flow**（CI/CD）：流水线、运行记录、构建任务等
  - **Appstack**（应用部署）：应用、环境、发布、变更单等
  - **Platform**（组织管理）：组织、部门、成员、角色等
  - **Packages**（制品管理）：制品仓库、制品版本等
  - **Lingma**（智能研发）：知识库、成员用量等
- **22 个增强聚合工具**：将多个 OpenAPI 调用合并为单次用户友好的操作。例如 `get_project_overview` 同时返回项目信息、成员、迭代、里程碑和标签；`get_repository_overview` 聚合仓库详情、分支、近期提交和合并请求。
- 工具启用/禁用过滤，支持 `--minimal` 精简模式和 `--project-focused` 项目聚焦模式。

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
- `/healthz`：readiness check；未注册工具时返回 `503`

Docker：

```bash
docker build -t yunxiao-mcp-server:local .
docker run -i --rm -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> yunxiao-mcp-server:local
```

Docker HTTP 模式：

```bash
docker run --rm -p 3000:3000 -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> yunxiao-mcp-server:local --port 3000
```

MCP client 配置示例见 [docs/mcp-client-config.md](docs/mcp-client-config.md)。

GA 边界、验证门禁和暂缓暴露的 OpenAPI endpoint 见 [docs/ga-readiness.md](docs/ga-readiness.md)。
本地二进制 smoke 可运行 `make smoke`；默认监听 `39393`，可用 `PORT=<port> make smoke` 覆盖。

## 配置

配置优先级为 programmatic explicit values > flag > environment > config file > defaults。普通运行只会用到 flag 及其之后的层级；programmatic explicit values 用于测试或内嵌调用。

常用环境变量：

- `YUNXIAO_MCP_ACCESS_TOKEN`：云效 access token
- `YUNXIAO_MCP_BASE_URL`：云效 OpenAPI host 或 API base URL，默认 `https://openapi-rdc.aliyuncs.com`
- `YUNXIAO_MCP_SSE_BASE_URL`：反向代理场景下 SSE message endpoint 的 public base URL
- `YUNXIAO_ACCESS_TOKEN`：兼容 Node 参考实现的 token 变量
- `YUNXIAO_API_BASE_URL`：兼容 Node 参考实现的 base URL 变量

若同时设置新旧环境变量，`YUNXIAO_MCP_*` 优先于兼容用的 legacy 变量。

HTTP/SSE 模式下，客户端也可以在请求 header `x-yunxiao-token` 或 query `yunxiao_access_token` 中传入 token；请求级 token 优先于启动时配置的默认 token，适合多人共享一个 HTTP MCP 服务。

`base_url` 可以是主域名，也可以已经包含 `/oapi/v1`。客户端会避免重复追加 `/oapi/v1`。

## 开发

```bash
make format
make tidy
make lint
make test
make build
```
