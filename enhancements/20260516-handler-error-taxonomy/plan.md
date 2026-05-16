# Plan: Apply Error Taxonomy to Handler Helpers

## Task DAG

1. Add `ValidationError` type + update `ClassifyError` in `errors.go`
   - Validation: `go build ./pkg/toolset/yunxiao`
2. Add tests for `ValidationError` classification in `errors_test.go`
   - Validation: `go test ./pkg/toolset/yunxiao -run ValidationError`
3. Update helpers (`requiredString`, `requiredNumberPathString`, `getClient`) in `helpers.go`
   - Validation: `go test ./pkg/toolset/yunxiao`
4. Full CI check
   - Validation: `make ci`

## Implementation

### errors.go changes

- Add `ValidationError` struct with `Msg string`
- Implement `Error() string` method
- Update `ClassifyError` to check `errors.As(err, &ValidationError)` and return `ErrValidation`

### helpers.go changes

- `requiredString`: return `&ValidationError{Msg: fmt.Sprintf("%s is required", key)}`
- `requiredNumberPathString`: return `&ValidationError{Msg: ...}`
- `getClient`: return `&ValidationError{Msg: "yunxiao client is not configured"}`
