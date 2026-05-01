package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackSystemReleaseTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_system_release_workflows",
				mcp.WithDescription("List AppStack release workflows for a system."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListSystemReleaseWorkflows,
		},
		{
			Tool: mcp.NewTool("get_release",
				mcp.WithDescription("Get an AppStack system release by serial number."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Release serial number.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetRelease,
		},
		{
			Tool: mcp.NewTool("list_release_members",
				mcp.WithDescription("List members of an AppStack system release."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Release serial number.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListReleaseMembers,
		},
		{
			Tool: mcp.NewTool("list_release_products",
				mcp.WithDescription("List products attached to an AppStack system release."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Release serial number.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListReleaseProducts,
		},
	}
}
