package yunxiao

import (
	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
	"github.com/mark3labs/mcp-go/mcp"
)

func apiCallTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("call_yunxiao_api",
				mcp.WithDescription("Execute a read-only Yunxiao OpenAPI call. Use this as a fallback for endpoints not covered by dedicated tools. GET is supported; POST is limited to read-only search/list-style endpoints. Prefer dedicated tools (e.g., search_workitems, list_pipelines) when available, as they provide better parameter validation and error guidance."),
				mcp.WithString("path", mcp.Required(), mcp.Description("API path relative to the base URL, e.g. /projex/organizations/{orgId}/projects/{projectId}")),
				mcp.WithString("method", mcp.Description("HTTP method: GET or POST. Defaults to GET.")),
				mcp.WithString("queryParams", mcp.Description("JSON string of query parameters to append to the URL.")),
				mcp.WithString("body", mcp.Description("JSON string of request body for POST requests.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleCallYunxiaoAPI,
		},
	}
}
