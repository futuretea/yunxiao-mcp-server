# Align Tool Surface Documentation

## Summary

The documentation has drifted from the current Yunxiao MCP tool surface. README and GA readiness text still frame the catalog as entirely read-only, while the implementation now supports four opt-in Projex write tools that are filtered out by the default `read_only=true` mode. The generated domain docs also describe every domain as read-only, even when Projex includes write tools.

This proposal updates the public docs and the documentation generator so the safety boundary is explicit and durable: default usage remains read-only, write tools are opt-in, and generated docs distinguish read-only from write-capable tools.

## Motivation

### Goals

- Replace stale all-read-only claims with accurate current-mode wording.
- Document that write tools require `read_only=false`.
- Document `minimal` and `project_focused` configuration in the example config.
- Make generated domain docs identify tool access type so `make docs` preserves the correction.
- Keep runtime behavior unchanged and prove that with existing build, lint, and test gates.

### Non-Goals

- No runtime behavior changes.
- No new Yunxiao API endpoints.
- No mutation confirmation flow.
- No package publishing or remote side effects.

## Proposal

### User Stories

- As a user installing the MCP server, I can tell that the default mode is read-only and safe for exploration.
- As an operator who wants Projex write actions, I can see that those tools require an explicit `read_only=false` choice.
- As a maintainer, I can run `make docs` without reintroducing false read-only wording.

### User Experience

`README.md` should describe the default read-only catalog and the opt-in write tools separately. `docs/ga-readiness.md` should describe the GA safety posture as default read-only plus opt-in mutation scope, rather than saying the entire project is read-only. `config.example.yaml` should show `project_focused`, `minimal`, `enabled_domains`, and `disabled_domains` keys.

Generated domain docs should include an access/type column and use wording that does not call write tools read-only.

### API Changes

None.

## Plan Readiness

### Open Questions and Assumptions

- blocking_unknown: none
- assumption: current code behavior is authoritative for tool counts and safety modes; generated access labels follow each tool's `ReadOnlyHint` annotation, with unannotated tools treated as write-capable.
- advisory: exact catalog counts may continue to change as tools are added; generated docs should carry per-tool access facts.

### Decision Log

| Decision | Reason | Evidence |
| --- | --- | --- |
| Fix the generator, not only checked-in docs | Prevents `make docs` from restoring stale wording. | `scripts/gen-tool-docs.go` |
| Keep write mode clearly opt-in | Matches current safety boundary. | `read_only` defaults in CLI/config |
| Avoid runtime code changes | The issue is documentation drift, not behavior. | `make lint` and `make test` already pass |

## Design

### Architecture

```text
Tool definitions + annotations
        |
        v
scripts/gen-tool-docs.go
        |
        v
docs/*-tools.md with access type

Runtime config and tool filtering
        |
        v
README.md / docs/ga-readiness.md / config.example.yaml
```

### Implementation Overview

- Extend `scripts/gen-tool-docs.go` to extract `mcp.WithReadOnlyHintAnnotation(true)` and mark tools without the annotation as write-capable.
- Update generated domain docs through `make docs`.
- Update README, GA readiness, and config example text.
- Add focused tests for generator output so read/write wording is covered. Tests should assert that generated Projex docs mark the four write tools as write-capable, keep read-only tools read-only, and avoid domain-level all-read-only wording.

### Test Plan

- Run `go test ./scripts`.
- Run `make docs` and inspect generated doc diffs.
- Run `make lint`.
- Run `make test`.

### Upgrade Strategy

No runtime migration is required. If the wording is wrong, revert this commit or rerun `make docs` after fixing the generator.
