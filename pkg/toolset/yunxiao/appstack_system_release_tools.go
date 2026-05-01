package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackSystemReleaseTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 6)
	tools = append(tools, appstackSystemReleaseOverviewTools()...)
	tools = append(tools, appstackSystemReleaseExecutionTools()...)
	return tools
}

func appstackSystemReleaseOverviewTools() []toolset.ServerTool {
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

func appstackSystemReleaseExecutionTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_attached_change_requests",
				mcp.WithDescription("List change requests attached to an AppStack system release."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name.")),
				mcp.WithString("releaseSn", mcp.Required(), mcp.Description("Release serial number.")),
				mcp.WithNumber("current", mcp.Description("Current page number. Defaults to 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size. Defaults to 10.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAttachedChangeRequests,
		},
		{
			Tool: mcp.NewTool("list_release_executions",
				mcp.WithDescription("List execution records for an AppStack system release."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Release serial number.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number.")),
				mcp.WithString("releaseStageSn", mcp.Required(), mcp.Description("Release stage serial number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size, up to 100.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: id or gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListReleaseExecutions,
		},
	}
}
