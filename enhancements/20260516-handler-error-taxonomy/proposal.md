# Proposal: Apply Error Taxonomy to Handler Helpers

## Problem

The error taxonomy (`ErrorCategory`, `ClassifyError`, `WrapError`) is defined in `errors.go` and applied in `client_request.go`, but handler helper functions (`requiredString`, `requiredNumberPathString`, `getClient`) return plain `fmt.Errorf` errors that `ClassifyError` cannot categorize. These pass through `WrapError` unchanged, so MCP consumers get raw untagged errors for validation/config failures.

## Solution

Add a `ValidationError` type to `errors.go`, register it in `ClassifyError` to return `ErrValidation`, and update `requiredString`, `requiredNumberPathString`, and `getClient` in `helpers.go` to return `ValidationError` instead of plain `fmt.Errorf`.

## Scope

- `pkg/toolset/yunxiao/errors.go` — add `ValidationError` type, update `ClassifyError`
- `pkg/toolset/yunxiao/errors_test.go` — add tests for `ValidationError` classification
- `pkg/toolset/yunxiao/helpers.go` — use `ValidationError` in `requiredString`, `requiredNumberPathString`, `getClient`

## Out of scope

- Individual handler file modifications (helpers wrap all call sites automatically)
- New error categories beyond `ErrValidation`
- API contract changes
