package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackDeploymentResourceTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_resource_instances",
				mcp.WithDescription("List AppStack resource instances in a resource pool."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("poolName", mcp.Required(), mcp.Description("Resource pool name.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Yunxiao currently supports keyset.")),
				mcp.WithNumber("perPage", mcp.Description("Page size, up to 100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: id or gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithNumber("page", mcp.Description("Page number when using page pagination.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListResourceInstances,
		},
	}
}
