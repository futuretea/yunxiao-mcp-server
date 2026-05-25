# Yunxiao MCP Server

[English](README.md)

让 AI 编码助手直接对话 [阿里云·云效](https://www.aliyun.com/product/yunxiao) —— 查项目、管迭代、审代码、看流水线，无需离开 IDE。

**默认只读，安全第一。** 193 个工具中 177 个为只读查询，16 个写操作需显式开启 `read_only=false`。

---

## 能做什么？

| 场景 | 对应工具 |
|------|----------|
| 📋 项目管理 | 查项目、工作项、迭代、里程碑、成员，创建/更新工作项 |
| 🔍 代码审查 | 查仓库、分支、提交、合并请求、变更请求（CR），创建/关闭/合并 CR |
| 🚀 流水线 | 查流水线、运行记录、构建任务，审批/拒绝人工卡点 |
| 📦 发布管控 | 查应用、环境、发布单、变更单、资源 |
| 🧠 知识库 | 查 Lingma 知识库、成员、文件 |
| 🤖 AI 使用分析 | 查团队成员 Lingma 采纳情况 |

> 写操作（创建/更新工作项、管理 CR/MR、流水线审批）仅在 `read_only=false` 时可用。

---

## 快速开始

### npx（零安装）

```bash
npx -y @futuretea/yunxiao-mcp-server
```

带上 token：

```bash
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> npx -y @futuretea/yunxiao-mcp-server
```

### Docker

**Stdio MCP 模式（Docker 默认入口）：**

```bash
docker run -i --rm -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> \
  ghcr.io/futuretea/yunxiao-mcp-server:latest
```

**HTTP MCP 模式：**

```bash
docker run --rm -p 3000:3000 -e YUNXIAO_MCP_ACCESS_TOKEN=<your-token> \
  ghcr.io/futuretea/yunxiao-mcp-server:latest --port 3000
```

### 从源码构建

```bash
make build
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao mcp
```

同一个 `yunxiao` 二进制也提供面向人的 CLI 命令。

```bash
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao organization list
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao member list
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao project list
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao project member list --project-id <project-id>
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao sprint list --project-id <project-id>
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao task list --project-id <project-id>
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao task view <workitem-id>
./bin/yunxiao tools describe search_projects
YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao tools call get_current_user --params '{}'
printf '{"page":1}' | YUNXIAO_MCP_ACCESS_TOKEN=<your-token> ./bin/yunxiao tools call list_organization_members --params-file -
```

### 接入 IDE

参考 [MCP Client Config](docs/mcp-client-config.md) 配置 Claude、Cursor 等客户端。

---

## 工具概览

| 领域 | 工具数 | 权限 | 说明 |
|------|--------|------|------|
| **Projex** | 47 | 43 只读 + 4 可写 | 项目、迭代、工作项、里程碑、测试用例 |
| **Codeup** | 37 | 31 只读 + 6 可写 | 仓库、分支、提交、MR、CR、代码评审 |
| **Flow** | 18 | 16 只读 + 2 可写 | 流水线、运行记录、构建任务、人工审批 |
| **Appstack** | 62 | 58 只读 + 4 可写 | 应用、环境、发布单、变更单 |
| **Platform** | 18 | 只读 | 组织、部门、成员、角色 |
| **Packages** | 3 | 只读 | 制品仓库与版本 |
| **Lingma** | 6 | 只读 | 知识库与使用统计 |
| **API** | 1 | 只读 | 通用 API 调用（兜底） |
| **Meta** | 1 | 只读 | 工具发现 |

### 增强工具

增强工具将多次 API 调用聚合为一次查询，减少 AI 往返次数。

| 工具 | 聚合内容 |
|------|----------|
| `get_project_overview` | 项目信息 + 成员 + 迭代 + 里程碑 + 标签 |
| `get_project_workitem_detail` | 工作项详情 + 活动 + 评论 + 附件 + 关联 |
| `get_repository_overview` | 仓库信息 + 默认分支 + 最近提交 + 最近 MR |
| `get_change_request_overview` | CR 详情 + Patch Sets + 评论 |
| `get_pipeline_overview` | 流水线信息 + 最近运行 + 历史 |

完整列表见 [Enhanced Tools Index](docs/enhanced-tools-index.md)。

---

## 配置

优先级：命令行参数 > 环境变量 > 配置文件 > 默认值。

### 必需

| 变量 | 说明 |
|------|------|
| `YUNXIAO_MCP_ACCESS_TOKEN` | 云效访问令牌 |

### 可选

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `YUNXIAO_MCP_BASE_URL` | API 地址 | `https://openapi-rdc.aliyuncs.com` |
| `YUNXIAO_MCP_SSE_BASE_URL` | SSE 公网地址（反向代理场景） | — |
| `YUNXIAO_MCP_INSECURE_SKIP_TLS_VERIFY` | 跳过 TLS 证书校验（仅内网） | `false` |

兼容旧环境变量 `YUNXIAO_ACCESS_TOKEN` 和 `YUNXIAO_API_BASE_URL`。

### 模式切换

| 参数 | 默认 | 作用 |
|------|------|------|
| `--read-only` | `true` | 只读模式。设为 `false` 开启写操作 |
| `--compact` | `true` | 隐藏有增强替代的基础 API 工具（设为 `false` 显示全部） |
| `--enabled-tools` | — | 按名称白名单启用 |
| `--disabled-tools` | — | 按名称黑名单禁用 |
| `--enable-domains` | — | 按领域白名单启用 |
| `--disable-domains` | — | 按领域黑名单禁用 |

### 配置文件

```bash
./bin/yunxiao mcp --config config.example.yaml
```

### 按请求切换 Token（HTTP/SSE）

```bash
curl -H "x-yunxiao-token: <token>" http://localhost:3000/mcp
# 或
http://localhost:3000/sse?yunxiao_access_token=<token>
```

---

## HTTP 端点

| 端点 | 用途 |
|------|------|
| `/mcp` | Streamable HTTP MCP |
| `/sse` | SSE MCP |
| `/message` | SSE 消息端点 |
| `/healthz` | 健康检查（工具未注册时返回 503） |

---

## 安全

- **默认只读**：177 个工具无需写权限，可安全探索。
- **显式开启写入**：16 个写工具需手动设置 `read_only=false`。
- **请求级 Token**：HTTP/SSE 模式下支持按请求覆盖 token，多用户场景互不干扰。
- **不暴露敏感端点**：管理员审计日志、个人令牌查询等高权限端点不在目录中。

---

## 开发

```bash
make fmt      # 格式化
make tidy     # 整理依赖
make lint     # 静态检查
make test     # 运行测试
make build    # 构建 yunxiao 二进制
make smoke    # 冒烟测试
make ci       # CI 全量检查
```

公共 Yunxiao SDK 位于 `pkg/yunxiao`；MCP 模式和 CLI 命令都复用它发起认证 OpenAPI 请求、处理路径编码、响应元数据和错误分类。

覆盖率要求 98%，`make coverage-check` 验证。

---

## 文档

- [MCP Client Config](docs/mcp-client-config.md) — IDE 集成示例
- [Quick Start Guide](docs/quickstart.md) — 常见 AI 对话模式
- [Enhanced Tools Index](docs/enhanced-tools-index.md) — 增强工具参考
- [Conditions Cookbook](docs/conditions-cookbook.md) — 查询条件示例
- [Pagination Guide](docs/pagination-guide.md) — 分页模式参考
- [GA Readiness](docs/ga-readiness.md) — 发版检查清单

---

## License

MIT
