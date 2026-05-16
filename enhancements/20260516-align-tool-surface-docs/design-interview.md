# Design Interview

## Requirement

Align public documentation and generated tool docs with the current Yunxiao MCP tool surface:

- Default mode is `read_only: true`.
- Four Projex write tools exist but are only available when `read_only: false`.
- `minimal` and `project_focused` modes affect the registered tool set.
- Generated domain docs must not describe write-capable tools as read-only.

## Interview Mode

skipped

## Confirmed Answers

| Question | Answer | Evidence |
| --- | --- | --- |
| What is the default safety mode? | Read-only is enabled by default. | `internal/cmd/root.go`, `pkg/core/config/config.go`, `config.example.yaml` |
| Are write tools part of the full catalog? | Yes, four Projex write tools are present when `read_only=false`. | `pkg/toolset/yunxiao/projex_write_tools.go`, `pkg/toolset/yunxiao/toolset.go` |
| Should this change alter runtime behavior? | No. Documentation and generator output only. | Current tests already cover filtering behavior. |

## Blocking Questions

None.

## Assumptions

- The current code behavior is authoritative.
- Documentation should describe both default safety posture and opt-in write mode.
