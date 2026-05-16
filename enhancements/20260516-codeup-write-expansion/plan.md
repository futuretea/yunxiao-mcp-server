# Plan: Codeup Write Tool Expansion — First Safe Slice

## Scope

3 POST-based write tools: `create_change_request`, `add_change_request_comment`, `create_merge_request`. All gated behind `read_only=false`.

## Task DAG

1. Create `codeup_write_tools.go` — tool schemas for 3 new tools
   - Dependencies: none
   - Validation: `go build ./pkg/toolset/yunxiao`
   - Pattern: follow `projex_write_tools.go`

2. Create `codeup_write_handlers.go` — handler implementations
   - Dependencies: task 1
   - Validation: `go build ./pkg/toolset/yunxiao`
   - Pattern: follow `projex_write_handlers.go` (getClient, requiredString, PostJSONWithMetadata)

3. Create `codeup_write_handlers_test.go` — handler tests
   - Dependencies: task 2
   - Validation: `go test ./pkg/toolset/yunxiao -run CodeupWrite`
   - Pattern: follow `projex_write_handlers_test.go` with `newHandlerTestClient`

4. Register in toolset and minimal mode
   - Dependencies: tasks 1-3
   - Updates: `toolset.go` (register in projectFocusedHiddenTools or minimalToolNames), `toolset_test.go`
   - Validation: `make ci`

5. Full CI verification
   - Dependencies: tasks 1-4
   - Validation: `make ci`, verify `read_only=true` mode excludes new tools

## Implementation Details

### Tool schemas (codeup_write_tools.go)

```go
func codeupWriteTools() []toolset.ServerTool {
    return []toolset.ServerTool{
        {
            Tool: mcp.NewTool("create_change_request",
                mcp.WithDescription("Create a new Codeup change request..."),
                mcp.WithString("organizationId", mcp.Required(), ...),
                mcp.WithString("repositoryId", mcp.Required(), ...),
                mcp.WithString("title", mcp.Required(), ...),
                ...
            ),
            Handler: handleCreateChangeRequest,
        },
        // ... add_change_request_comment, create_merge_request
    }
}
```

### Handler pattern

Each handler: getClient → validate required params → build body → PostJSONWithMetadata. Mirror `handleCreateWorkitem` exactly.

### API paths

- `create_change_request`: POST `/codeup/organizations/{orgId}/repositories/{repoId}/changeRequests`
- `add_change_request_comment`: POST `/codeup/.../changeRequests/{localId}/comments`
- `create_merge_request`: POST `/codeup/organizations/{orgId}/repositories/{repoId}/mergeRequests`

Verify exact paths against reference project (`third-party-projects/alibabacloud-devops-mcp-server`) before implementing.
