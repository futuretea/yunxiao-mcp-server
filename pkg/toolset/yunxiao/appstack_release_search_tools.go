package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackReleaseSearchTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("search_releases",
				mcp.WithDescription("Search AppStack releases in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Yunxiao currently supports keyset.")),
				mcp.WithNumber("perPage", mcp.Description("Page size, up to 100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: id or gmt_create.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithString("nameKeyword", mcp.Description("Release display-name search keyword.")),
				mcp.WithString("systemName", mcp.Description("System unique name.")),
				mcp.WithArray("states", mcp.Description("Release states such as DEVELOPING, RELEASING, CLOSED, or RELEASED."), mcp.WithStringItems()),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleSearchReleases,
		},
	}
}
