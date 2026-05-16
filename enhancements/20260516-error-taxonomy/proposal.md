# Structured API Error Taxonomy

## Summary

Current error handling is ad-hoc: handlers use raw `fmt.Errorf`, `client_request.go` has `friendlyAPIError` for HTTP status mapping, and there's no way for callers or MCP clients (LLMs) to programmatically distinguish error categories. This leads to repeated error-propagation bugs (e.g., recent fixes for status comment failures, parameter validation, read-only boundary).

Define a lightweight error taxonomy with typed sentinel errors and category classification, then migrate the client layer and critical handler paths.

## Motivation

### Goals

- Define error categories: auth, permission, validation, rate-limit, server, network
- Each category maps to a distinct MCP-facing error message pattern that LLMs can recognize
- Migrate `friendlyAPIError` and `Request()` to use typed errors
- `make lint` + `make test` pass, no coverage regression
- Keep changes under 200 LOC for first safe slice

### Non-Goals

- No breaking changes to handler function signatures
- No runtime behavior changes — errors are still strings to MCP clients
- No new external dependencies
- Not migrating all handlers in first slice — only client layer

## Proposal

### Error Categories

```go
type ErrorCategory string

const (
    ErrAuth        ErrorCategory = "auth"        // 401 — bad/missing token
    ErrPermission  ErrorCategory = "permission"  // 403 — insufficient scope
    ErrValidation  ErrorCategory = "validation"  // 400 — bad params
    ErrRateLimit   ErrorCategory = "rate_limit"  // 429 — backoff
    ErrServer      ErrorCategory = "server"      // 5xx — upstream down
    ErrNetwork     ErrorCategory = "network"     // timeout, DNS, TLS
)
```

### Category Detection

- `ClassifyError(err)` — inspects `APIError.StatusCode` or `net.Error` to return category
- `friendlyAPIError` already has the status→suggestion mapping; extend it to also set category

### MCP Consumer Guidance

Each category gets a standard prefix in the error message so LLMs can pattern-match:

| Category | Prefix |
|----------|--------|
| auth | `[auth] Authentication failed:` |
| permission | `[permission] Access denied:` |
| validation | `[validation] Invalid request:` |
| rate_limit | `[rate_limit] Rate limited:` |
| server | `[server] Service unavailable:` |
| network | `[network] Connection failed:` |

### First Slice Scope

1. New file `pkg/toolset/yunxiao/errors.go` — category type, `ClassifyError()`, and `WrapError()` helper
2. Update `client_request.go` — `Request()` and `friendlyAPIError()` to use categories
3. Tests for `ClassifyError` with each HTTP status category

### Risk Assessment

- **Risk**: Low. Types and classification are additive; existing error paths unchanged.
- **Mitigation**: 98.2% test coverage catches regressions; `make lint` + `make test` as gate.

## Review Status

self-reviewed: zero P0/P1. Additive change, no existing behavior modified.
