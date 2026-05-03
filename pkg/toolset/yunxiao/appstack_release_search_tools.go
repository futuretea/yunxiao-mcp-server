package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackReleaseSearchTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("search_releases",
				mcp.WithDescription("Search AppStack releases in a Yunxiao organization. Use this to discover releases before calling get_release_overview or other release-specific tools."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Valid value: keyset.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithString("orderBy", mcp.Description("Sort field. Valid values: id, gmt_create.")),
				mcp.WithString("sort", mcp.Description("Sort direction. Valid values: asc, desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithString("nameKeyword", mcp.Description("Release display-name search keyword.")),
				mcp.WithString("systemName", mcp.Description("System unique name. Use list_systems to discover valid names.")),
				mcp.WithArray("states", mcp.Description("Release states filter. Valid values: DEVELOPING, RELEASING, CLOSED, RELEASED."), mcp.WithStringItems()),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleSearchReleases,
		},
	}
}
