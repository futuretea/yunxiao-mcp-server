package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackSystemTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_systems",
				mcp.WithDescription("List AppStack systems in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithNumber("current", mcp.Description("Current page number. Defaults to 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size. Defaults to 10.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListSystems,
		},
		{
			Tool: mcp.NewTool("list_attached_apps",
				mcp.WithDescription("List applications attached to an AppStack system."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name.")),
				mcp.WithNumber("current", mcp.Description("Current page number. Defaults to 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size. Defaults to 10.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAttachedApps,
		},
	}
}
