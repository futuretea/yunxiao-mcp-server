# yunxiao-mcp-server

Go 语言版本的 Yunxiao MCP Server。当前实现提供可构建的 stdio MCP 服务骨架，并接入云效 OpenAPI 的基础只读工具。

## 功能

- stdio MCP transport
- HTTP Streamable MCP transport 与 SSE transport
- HTTP health endpoint：`/healthz`
- Cobra CLI 与 YAML/env/flag 配置加载
- Yunxiao OpenAPI token 认证：启动时默认 token，HTTP/SSE 请求级 `x-yunxiao-token` 或 `yunxiao_access_token` 覆盖
- 工具启用/禁用过滤
- 基础只读工具：
  - `get_current_user`
  - `get_current_organization_info`
  - `get_user_organizations`
  - `list_organizations`
  - `get_organization`
  - `list_organization_departments`
  - `get_organization_department_info`
  - `get_organization_department_ancestors`
  - `list_organization_members`
  - `get_organization_member_info`
  - `get_organization_member_info_by_user_id`
  - `search_organization_members`
  - `list_organization_roles`
  - `get_organization_role`
  - `list_users`
  - `list_repositories`
  - `get_repository`
  - `list_branches`
  - `get_branch`
  - `list_files`
  - `get_file_blobs`
  - `list_commits`
  - `get_commit`
  - `compare`
  - `list_change_requests`
  - `get_change_request`
  - `list_change_request_patch_sets`
  - `get_change_request_tree`
  - `list_change_request_comments`
  - `get_change_request_comment`
  - `list_pipelines`
  - `get_pipeline`
  - `list_pipeline_runs`
  - `get_pipeline_run`
  - `get_latest_pipeline_run`
  - `list_pipeline_jobs_by_category`
  - `list_pipeline_job_historys`
  - `get_pipeline_job_run_log`
  - `search_projects`
  - `get_project`
  - `search_workitems`
  - `get_workitem`
  - `list_package_repositories`
  - `list_artifacts`
  - `get_artifact`
  - `get_sprint`
  - `list_sprints`
  - `list_all_work_item_types`
  - `list_work_item_types`
  - `get_work_item_type`
  - `list_work_item_relation_work_item_types`
  - `get_work_item_type_field_config`
  - `get_work_item_workflow`
  - `list_work_item_comments`
  - `list_applications`
  - `get_application`
  - `get_env_variable_groups`
  - `get_variable_group`
  - `get_app_variable_groups`
  - `get_app_variable_groups_revision`
  - `get_latest_orchestration`
  - `list_app_orchestration`
  - `get_app_orchestration`
  - `list_app_release_workflows`
  - `list_app_release_workflow_briefs`
  - `get_app_release_workflow_stage`
  - `list_app_release_stage_briefs`
  - `list_app_release_stage_runs`
  - `list_app_release_stage_exec_metadata`
  - `get_app_release_stage_pipeline_run`
  - `get_app_release_stage_pipeline_job_log`
  - `get_appstack_change_request_audit_items`
  - `list_appstack_change_request_executions`
  - `list_appstack_change_request_work_items`
  - `list_change_order_versions`
  - `get_change_order`
  - `list_change_order_job_logs`
  - `find_task_operation_log`
  - `list_change_orders_by_origin`

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
