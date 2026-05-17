# Yunxiao MCP Server

[English](#english) | 中文

让 AI 编码助手直接对话 [阿里云·云效](https://www.aliyun.com/product/yunxiao) —— 查项目、管迭代、审代码、看流水线，无需离开 IDE。

Let your AI coding assistant talk directly to [Alibaba Yunxiao](https://www.aliyun.com/product/yunxiao) — browse projects, track iterations, review code, and monitor pipelines without leaving your IDE.

**默认只读，安全第一。** 142 个工具中 130 个为只读查询，12 个写操作需显式开启 `read_only=false`。

**Read-only by default, safety first.** 130 of 142 tools are read-only queries. 12 write operations require explicit `read_only=false`.

---

## 能做什么？ / What can you do?

| 场景 Scenario | 对应工具 Tools |
|------|----------|
| 📋 项目管理 Project Mgmt | 查项目、工作项、迭代、里程碑、成员，创建/更新工作项 |
| 🔍 代码审查 Code Review | 查仓库、分支、提交、合并请求、变更请求（CR），创建/关闭/合并 CR |
| 🚀 流水线 Pipelines | 查流水线、运行记录、构建任务，审批/拒绝人工卡点 |
| 📦 发布管控 Releases | 查应用、环境、发布单、变更单、资源 |
| 🧠 知识库 Knowledge Base | 查 Lingma 知识库、成员、文件 |
| 🤖 AI 使用分析 AI Adoption | 查团队成员 Lingma 采纳情况 |

> 写操作（创建/更新工作项、管理 CR/MR、流水线审批）仅在 `read_only=false` 时可用。
> Write operations (create/update work items, manage CR/MR, pipeline approvals) require `read_only=false`.

---

## 快速开始 / Quick Start

### npx（零安装 / zero install）

```bash
npx -y @futuretea/yunxiao-mcp-server
```

带上 token / With token：

```bash
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> npx -y @futuretea/yunxiao-mcp-server
```

### Docker

**Stdio 模式（默认） / Stdio mode (default)：**

```bash
docker run -i --rm -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> \
  ghcr.io/futuretea/yunxiao-mcp-server:latest
```

**HTTP 模式 / HTTP mode：**

```bash
docker run --rm -p 3000:3000 -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> \
  ghcr.io/futuretea/yunxiao-mcp-server:latest --port 3000
```

### 从源码构建 / Build from source

```bash
make build
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao-mcp-server
```

### 接入 IDE / IDE Setup

参考 [MCP Client Config](docs/mcp-client-config.md) 配置 Claude、Cursor 等客户端。

See [MCP Client Config](docs/mcp-client-config.md) for Claude, Cursor, and other IDE setup examples.

---

## 工具概览 / Tools Overview

| 领域 Domain | 工具数 Tools | 权限 Access | 说明 Description |
|------|--------|------|------|
| **Projex** | 45 | 41 read-only + 4 write | 项目、迭代、工作项、里程碑、测试用例 / Projects, iterations, work items, milestones, test cases |
| **Codeup** | 30 | 24 read-only + 6 write | 仓库、分支、提交、MR、CR、代码评审 / Repositories, branches, commits, MR, CR, code review |
| **Flow** | 10 | 8 read-only + 2 write | 流水线、运行记录、构建任务、人工审批 / Pipelines, runs, build tasks, validation |
| **Appstack** | 31 | read-only | 应用、环境、发布单、变更单 / Applications, environments, releases, change orders |
| **Platform** | 18 | read-only | 组织、部门、成员、角色 / Organizations, departments, members, roles |
| **Packages** | 2 | read-only | 制品仓库与版本 / Artifact repositories and versions |
| **Lingma** | 4 | read-only | 知识库与使用统计 / Knowledge bases and usage |
| **API** | 1 | read-only | 通用 API 调用（兜底） / Generic API fallback |
| **Meta** | 1 | read-only | 工具发现 / Tool discovery |

### 增强工具 / Enhanced Tools

增强工具将多次 API 调用聚合为一次查询，减少 AI 往返次数。
Enhanced tools aggregate multiple API calls into single operations, reducing AI round-trips.

| 工具 Tool | 聚合内容 What it combines |
|------|----------|
| `get_project_overview` | 项目信息 + 成员 + 迭代 + 里程碑 + 标签 / Project info + members + sprints + milestones + labels |
| `get_project_workitem_detail` | 工作项详情 + 活动 + 评论 + 附件 + 关联 / Work item + activities + comments + attachments + relations |
| `get_repository_overview` | 仓库信息 + 默认分支 + 最近提交 + 最近 MR / Repository + default branch + recent commits + recent MRs |
| `get_change_request_overview` | CR 详情 + Patch Sets + 评论 / CR detail + patch sets + comments |
| `get_pipeline_overview` | 流水线信息 + 最近运行 + 历史 / Pipeline info + latest run + history |

完整列表 / Full list: [Enhanced Tools Index](docs/enhanced-tools-index.md)

---

## 配置 / Configuration

优先级 / Priority：命令行参数 flags > 环境变量 env > 配置文件 config file > 默认值 defaults。

### 必需 / Required

| 变量 Variable | 说明 Description |
|------|------|
| `YUNXIAO_MCP_ACCESS_TOKEN` | 云效访问令牌 / Yunxiao access token |

### 可选 / Optional

| 变量 Variable | 说明 Description | 默认值 Default |
|------|------|--------|
| `YUNXIAO_MCP_BASE_URL` | API 地址 / API base URL | `https://openapi-rdc.aliyuncs.com` |
| `YUNXIAO_MCP_SSE_BASE_URL` | SSE 公网地址（反向代理场景） / Public SSE base URL | — |
| `YUNXIAO_MCP_INSECURE_SKIP_TLS_VERIFY` | 跳过 TLS 证书校验（仅内网） / Skip TLS verify (internal only) | `false` |

兼容旧环境变量 / Legacy aliases: `YUNXIAO_ACCESS_TOKEN`, `YUNXIAO_API_BASE_URL`.

### 模式切换 / Tool Modes

| 参数 Flag | 默认 Default | 作用 Purpose |
|------|------|------|
| `--read-only` | `true` | 只读模式。设为 `false` 开启写操作 / Set `false` to enable write tools |
| `--project-focused` | `false` | 仅加载 Platform + Projex 工具 / Platform + Projex only |
| `--minimal` | `false` | 最小工具集（约 14 个核心工具） / Minimal toolset (~14 tools) |
| `--enabled-tools` | — | 按名称白名单启用 / Explicit tool allow-list |
| `--disabled-tools` | — | 按名称黑名单禁用 / Explicit tool deny-list |
| `--enable-domains` | — | 按领域白名单启用 / Domain allow-list |
| `--disable-domains` | — | 按领域黑名单禁用 / Domain deny-list |

### 配置文件 / Config File

```bash
./bin/yunxiao-mcp-server --config config.example.yaml
```

### 按请求切换 Token（HTTP/SSE） / Per-Request Token

```bash
curl -H "x-yunxiao-token: <token>" http://localhost:3000/mcp
# 或 / or
http://localhost:3000/sse?yunxiao_access_token=<token>
```

---

## HTTP 端点 / HTTP Endpoints

| 端点 Endpoint | 用途 Purpose |
|------|------|
| `/mcp` | Streamable HTTP MCP |
| `/sse` | SSE MCP |
| `/message` | SSE 消息端点 / SSE message endpoint |
| `/healthz` | 健康检查（工具未注册时返回 503） / Health check (503 if no tools) |

---

## 安全 / Security

- **默认只读 / Read-only by default**：130 个工具无需写权限，可安全探索。130 tools safe for exploration.
- **显式开启写入 / Explicit write opt-in**：12 个写工具需手动设置 `read_only=false`。12 write tools require manual opt-in.
- **请求级 Token / Per-request token**：HTTP/SSE 模式下支持按请求覆盖 token，多用户场景互不干扰。Request-level token override for multi-tenant use.
- **不暴露敏感端点 / No sensitive endpoints**：管理员审计日志、个人令牌查询等高权限端点不在目录中。Admin audit logs, PAT queries excluded.

---

## 开发 / Development

```bash
make fmt      # 格式化 / gofmt
make tidy     # 整理依赖 / go mod tidy
make lint     # 静态检查 / go vet + gofmt
make test     # 运行测试 / go test
make build    # 构建 / build binary
make smoke    # 冒烟测试 / smoke test
make ci       # CI 全量检查 / full CI
```

覆盖率要求 / Coverage threshold: 98%. Run `make coverage-check`.

---

## 文档 / Documentation

- [MCP Client Config](docs/mcp-client-config.md) — IDE 集成示例 / IDE setup
- [Quick Start Guide](docs/quickstart.md) — 常见 AI 对话模式 / Conversation patterns
- [Enhanced Tools Index](docs/enhanced-tools-index.md) — 增强工具参考 / Enhanced tool reference
- [Conditions Cookbook](docs/conditions-cookbook.md) — 查询条件示例 / Query filter examples
- [Pagination Guide](docs/pagination-guide.md) — 分页模式参考 / Pagination reference
- [GA Readiness](docs/ga-readiness.md) — 发版检查清单 / Release checklist

---

## License

MIT
