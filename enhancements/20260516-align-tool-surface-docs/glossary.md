# Glossary

| Term | Meaning |
| --- | --- |
| Full catalog | All tools registered by `Toolset{ReadOnly:false}` before mode and filter reductions. |
| Default read-only catalog | Tools available when `read_only=true`; write tools are excluded. |
| Write tool | A tool that performs a Yunxiao mutation and is excluded when `read_only=true`. |
| Minimal mode | A reduced project-centric catalog selected by `minimal=true`; write tools still require `read_only=false`. |
| Project-focused mode | A platform + Projex catalog selected by `project_focused=true`; superseded by explicit domain enables. |
| Generated domain docs | `docs/*-tools.md` files produced by `scripts/gen-tool-docs.go`. |
