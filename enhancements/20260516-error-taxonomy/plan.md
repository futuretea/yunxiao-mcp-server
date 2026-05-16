# Implementation Plan

## Scope

First safe slice: define error category types, classification logic, and migrate client-layer error paths.

In scope:
- `pkg/toolset/yunxiao/errors.go` (new)
- `pkg/toolset/yunxiao/errors_test.go` (new)
- `pkg/toolset/yunxiao/client_request.go` (modify)
- `pkg/toolset/yunxiao/client.go` (modify — move friendlyAPIError to errors.go)

## Task DAG

1. Create errors.go with types and classification
   - Dependencies: none
   - Validation: `go build ./pkg/toolset/yunxiao`

2. Add tests for ClassifyError
   - Dependencies: task 1
   - Validation: `go test ./pkg/toolset/yunxiao -run ClassifyError`

3. Migrate friendlyAPIError and Request to use categories
   - Dependencies: tasks 1-2
   - Validation: `go test ./pkg/toolset/yunxiao`

4. Verify
   - Dependencies: tasks 1-3
   - Validation: `make lint`, `make test`

## Implementation Details

### errors.go

- `ErrorCategory` string type with 6 constants
- `ClassifyError(err error) ErrorCategory` — checks `APIError.StatusCode` and `net.Error`
- `WrapError(err error) error` — prepends category prefix to error message
- Move `friendlyAPIError` from client.go to errors.go (it's error classification logic)

### ClassifyError logic

```
401 → ErrAuth
403 → ErrPermission
400 → ErrValidation
429 → ErrRateLimit
5xx → ErrServer
net.Error (timeout/temporary) → ErrNetwork
other → "" (uncategorized)
```

### client_request.go changes

- `GetJSON`, `GetJSONWithMetadata`, `PostJSONWithMetadata`, `PutJSONWithMetadata` — already call `friendlyAPIError`, no change needed
- `Request()` — add category to APIError so classification works downstream

## Test Strategy

- `TestClassifyError` — table-driven: each HTTP status → expected category
- `TestClassifyErrorNetwork` — mock net.Error for timeout
- Existing 735 tests guard against regressions

## Review Status

self-reviewed: zero P0/P1.
