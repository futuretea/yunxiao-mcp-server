# Proposal: Codeup Write Tool Expansion

## Problem

The project currently exposes 4 write-capable tools, all in the Projex domain (create/update workitem, update status, add comment). The Codeup domain has 0 write tools, despite the Yunxiao API supporting Codeup mutation operations (create/update change requests, merge requests, comments, etc.). MCP consumers that need to perform Codeup mutations must fall back to the generic `call_yunxiao_api` tool, which lacks structured parameter validation and domain-specific guidance.

## Opportunity

Add 4-6 Codeup write tools following the established Projex write tool pattern:
- Tool schema definition in `codeup_write_tools.go`
- Handler implementation in `codeup_write_handlers.go`
- Tests in `codeup_write_handlers_test.go`
- All gated behind `read_only=false` configuration
- Each tool maps to a single, well-defined Yunxiao API endpoint

## Candidate Tools

| Tool | API Path | Risk |
|------|----------|------|
| `create_change_request` | POST `/codeup/.../changeRequests` | Low — follows Projex `create_workitem` pattern |
| `update_change_request` | PUT `/codeup/.../changeRequests/{id}` | Low — follows `update_workitem` pattern |
| `add_change_request_comment` | POST `/codeup/.../changeRequests/{id}/comments` | Low — follows `add_workitem_comment` pattern |
| `create_merge_request` | POST `/codeup/.../mergeRequests` | Low |
| `add_merge_request_comment` | POST `/codeup/.../mergeRequests/{id}/comments` | Low |
| `close_merge_request` | PUT `/codeup/.../mergeRequests/{id}` with state=closed | Medium — state transition |

## First Safe Slice

Start with 3 tools: `create_change_request`, `add_change_request_comment`, `create_merge_request`. These are the simplest POST operations, mirroring the existing `create_workitem` and `add_workitem_comment` patterns exactly.

## Out of Scope

- Delete/unlink operations (higher risk)
- Complex state transitions (close/merge/reopen)
- Batch operations

## Risks

- API contract changes could break new tools — mitigated by following existing request shapes validated in the reference project
- Write tool surface expansion increases the blast radius of auth token leaks — mitigated by existing `read_only=false` gate
- Enhanced tools that aggregate write + read operations (e.g., "create CR and add reviewers") — deferred to future slice
