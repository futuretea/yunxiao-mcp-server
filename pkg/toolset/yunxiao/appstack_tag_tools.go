package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackTagTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("search_app_tags",
				mcp.WithDescription("Search AppStack application tags in a Yunxiao organization. Application tags are used to categorize and label applications for organizational purposes."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithNumber("current", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 10 when omitted.")),
				mcp.WithString("orderBy", mcp.Description("Sort field. Valid values: tagName, id. Defaults to id.")),
				mcp.WithString("sort", mcp.Description("Sort direction. Valid values: asc, desc. Defaults to desc.")),
				mcp.WithString("search", mcp.Description("Optional search keyword for fuzzy matching against tag names.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleSearchAppTags,
		},
	}
}
