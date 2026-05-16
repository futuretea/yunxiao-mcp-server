# Clarification

## Goal

Make the project documentation accurately describe the current tool catalog and mode behavior without changing runtime behavior.

## Known Facts

- `Toolset{ReadOnly: true}` filters `create_workitem`, `update_workitem`, `update_workitem_status`, and `add_workitem_comment`.
- `Toolset{ReadOnly: false}` exposes the full catalog, including the four write tools.
- `minimal` mode keeps essential project tools and includes write tools only when read-only filtering is disabled.
- `project_focused` mode restricts the catalog to platform and Projex tools, hiding low-value raw tools with enhanced alternatives.
- `scripts/gen-tool-docs.go` currently labels each generated domain document as read-only even when that domain includes write tools.

## Constraints

- Do not change runtime tool registration or public API behavior.
- Keep write-mode documentation clearly opt-in.
- Regenerate generated docs using the existing `make docs` workflow.

## Out of Scope

- Adding new tools.
- Changing read-only filtering.
- Introducing mutation confirmation UX.
- Publishing packages or images.
