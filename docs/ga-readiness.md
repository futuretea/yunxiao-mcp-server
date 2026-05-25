# GA Readiness

This project is GA-ready for local MCP use as a default-read-only Yunxiao OpenAPI server.

## Scope

- Transports: stdio, Streamable HTTP, SSE, and `/healthz`.
- Authentication: a startup default Yunxiao token, with request-level HTTP/SSE token override through `x-yunxiao-token` or `yunxiao_access_token`.
- Tool surface: 177 read-only MCP tools in the default `read_only=true` catalog (63 hidden by default `--compact`). The full catalog contains 193 tools: 177 read-only tools plus 16 write-capable tools (4 Projex: `create_workitem`, `update_workitem`, `update_workitem_status`, `add_workitem_comment`; 6 Codeup: `create_change_request`, `add_change_request_comment`, `create_merge_request`, `close_change_request`, `reopen_change_request`, `merge_change_request`; 2 Flow: `pass_pipeline_validate`, `refuse_pipeline_validate`; 4 AppStack: `create_change_order`, `execute_job_action`, `execute_system_release_stage`, `execute_app_release_stage`) that are exposed only when `read_only=false`.
- Safety boundary: read-only API access by default. Write-capable tools require an explicit `read_only=false` configuration. Other endpoints with create, update, delete, execute, approve, refuse, or state-changing semantics are not exposed, even when Yunxiao models them as `GET`.

## Release Gate

Before treating a build as GA, run:

```bash
make ci
make smoke
git diff --check
```

`make ci` runs `go vet`, verifies module checksums, checks `gofmt`, runs `go test -race ./...`, and builds `bin/yunxiao`.
`make smoke` requires `curl` and `nc` on `PATH`; it runs the built binary's `version` and `--help` commands, starts local HTTP mode, and checks `/healthz`; it does not call Yunxiao OpenAPI.

## Deferred OpenAPI Endpoints

These endpoints are intentionally not exposed in the GA read-only surface.

| Operation | Path | Reason |
| --- | --- | --- |
| `ListAuditLogsForAdmin` | `/platform/auditLogs` | Enterprise-admin sensitive read. Organization-scoped audit logs are exposed as `list_audit_logs`; admin-wide audit access needs an explicit admin-mode and security review. |
| `ListUserPersonalAccessTokens` | `/platform/users/admin/personalAccessTokens` | Enterprise-admin sensitive token metadata. Expose only with an explicit admin-mode, redaction policy, and audit guidance. |
| `DeleteWorkitemFile` | `/projex/organizations/{organizationId}/workitems/{workitemId}/deletefile/{id}` | Uses `GET` but deletes a work item file. This is mutation behavior and is outside the read-only GA scope. |
| `GetPipelineRunV2` | `/flow/organizations/{organizationId}/pipelines/{pipelineId}/runs/{pipelineRunId}/v2` | The swagger path includes `pipelineRunId`, but the operation parameter list omits it. Defer until the contract is fixed or a trusted reference confirms the request shape. |
| `getMachineGroupDefault` | `/flow/organizations/{organizationId}/build/machinegroup/default` | Swagger declares a path parameter that is not present in the path. Defer until the contract is fixed or a trusted reference confirms the request shape. |

## Adding Deferred Endpoints Later

When adding any deferred endpoint:

- Keep it in a separate commit with focused handler and tool tests.
- Preserve large integer IDs as strings unless the API contract proves they fit safely in JSON numbers.
- Add path and query tests that assert exact request shape, including array query encoding.
- For mutation endpoints, keep them outside the default `read_only=true` catalog and require an explicit non-read-only mode.
