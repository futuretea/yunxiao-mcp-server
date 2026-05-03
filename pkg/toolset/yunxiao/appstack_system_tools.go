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
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithNumber("current", mcp.Description("Current page number. Defaults to 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size. Defaults to 10.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListSystems,
		},
		{
			Tool: mcp.NewTool("list_attached_apps",
				mcp.WithDescription("List applications attached to an AppStack system."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name.")),
				mcp.WithNumber("current", mcp.Description("Current page number. Defaults to 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size. Defaults to 10.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAttachedApps,
		},
		{
			Tool: mcp.NewTool("list_system_members",
				mcp.WithDescription("List members of an AppStack system."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name.")),
				mcp.WithNumber("current", mcp.Description("Current page number. Defaults to 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size. Defaults to 10.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListSystemMembers,
		},
	}
}
