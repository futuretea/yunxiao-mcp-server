# Implementation Plan

## Scope

Current unit: align tool surface and mode documentation with existing runtime behavior.

In scope:

- `scripts/gen-tool-docs.go`
- `scripts/gen-tool-docs_test.go` (new)
- generated `docs/*-tools.md`
- `README.md`
- `docs/ga-readiness.md`
- `config.example.yaml`

Out of scope:

- Runtime tool filtering changes.
- Tool schema changes.
- New write operations.

## Baseline

```text
CLI/config defaults
  read_only: true
  project_focused: false
  minimal: false

Toolset filtering
  Toolset{ReadOnly:true} excludes writeToolNames
  Toolset{ReadOnly:false} includes the full catalog

Generated docs
  scripts/gen-tool-docs.go currently emits domain-level read-only wording
  Projex generated docs include four write tools
```

## Architecture

```text
pkg/toolset/yunxiao/*_tools.go
        |
        | parse mcp.NewTool + WithReadOnlyHintAnnotation
        v
scripts/gen-tool-docs.go
        |
        v
docs/<domain>-tools.md
  - domain summary with total/access counts
  - inventory Access column
  - per-tool Access line
```

README, GA readiness, and config example document the mode boundary that runtime config already enforces.

## Task DAG

1. Generator output tests
   - Source: proposal review findings and access-label requirement.
   - Dependencies: none.
   - Skill: tdd.
   - TDD: required; this is the red producer.
   - Validation: tests initially fail until the generator extracts access metadata.
   - Assertions:
     - Annotated tools using `mcp.WithReadOnlyHintAnnotation(true)` are read-only.
     - Unannotated tools are write-capable.
     - Generated Projex docs mark the four write tools as write-capable.
     - Generated docs keep read-only tools read-only.
     - Generated docs avoid domain-level all-read-only wording.

2. Generator model update
   - Source: proposal goals for durable generated docs.
   - Dependencies: task 1.
   - Skill: go-engineering.
   - TDD: required; this is the green implementation.
   - Validation: `go test ./scripts`.

3. Regenerate domain docs
   - Source: generator update.
   - Dependencies: tasks 1 and 2.
   - Skill: none.
   - TDD: not-applicable; generated artifact update.
   - Validation: `make docs`, then `git diff -- docs`.

4. Hand-written docs/config update
   - Source: proposal goals for README, GA readiness, config example.
   - Dependencies: task 3.
   - Skill: tech-doc-review.
   - TDD: not-applicable; documentation-only.
   - Validation:
     - content scan for stale all-read-only claims.
     - `config.example.yaml` contains `project_focused`, `minimal`, `enabled_domains`, and `disabled_domains`.
     - README states that write tools require explicit `read_only=false`.
     - GA readiness states that write tools require explicit `read_only=false`.

5. Gate
   - Source: auto-ready readiness contract.
   - Dependencies: tasks 1-4.
   - Skill: feature-ready.
   - TDD: not-applicable.
   - Validation: `go test ./scripts`, `make docs`, `make lint`, `make test`, independent feature-ready review.

## Test Strategy

- `go test ./scripts` covers generator read/write access labeling.
- `make docs` verifies generated docs can be produced from source.
- `gofmt -w scripts/gen-tool-docs.go scripts/gen-tool-docs_test.go` formats edited script files because `make lint` only checks formatting under `cmd`, `internal`, and `pkg`.
- `make lint` covers Go vet and gofmt for `cmd`, `internal`, and `pkg`.
- `make test` covers repository tests.

## Review Status

review-passed: revised after independent review-plan P1 findings; ready for implementation.
