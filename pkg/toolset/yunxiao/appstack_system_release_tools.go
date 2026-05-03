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
				mcp.WithDescription("List AppStack release workflows for a system. Use this after discovering a system via list_systems to find releases and their workflows."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name. Use list_systems to discover valid names.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListSystemReleaseWorkflows,
		},
		{
			Tool: mcp.NewTool("list_release_members",
				mcp.WithDescription("List members of an AppStack system release. Use this after discovering a release via search_releases or list_system_release_workflows."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name. Use list_systems to discover valid names.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Release serial number. Use search_releases or list_system_release_workflows to discover valid serial numbers.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListReleaseMembers,
		},
		{
			Tool: mcp.NewTool("list_release_products",
				mcp.WithDescription("List products attached to an AppStack system release. Use this after discovering a release via search_releases or list_system_release_workflows."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name. Use list_systems to discover valid names.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Release serial number. Use search_releases or list_system_release_workflows to discover valid serial numbers.")),
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
				mcp.WithDescription("List change requests attached to an AppStack system release. Use this after discovering a release via search_releases or list_system_release_workflows."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name. Use list_systems to discover valid names.")),
				mcp.WithString("releaseSn", mcp.Required(), mcp.Description("Release serial number. Use search_releases or list_system_release_workflows to discover valid serial numbers.")),
				mcp.WithNumber("current", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAttachedChangeRequests,
		},
		{
			Tool: mcp.NewTool("list_release_executions",
				mcp.WithDescription("List execution records for an AppStack system release. Use this after discovering release workflow and stage details via list_system_release_workflows or get_release_overview."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name. Use list_systems to discover valid names.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Release serial number. Use search_releases or list_system_release_workflows to discover valid serial numbers.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number. Use list_system_release_workflows to discover valid serial numbers.")),
				mcp.WithString("releaseStageSn", mcp.Required(), mcp.Description("Release stage serial number. Use list_system_release_workflows to discover valid serial numbers.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithString("orderBy", mcp.Description("Sort field. Valid values: id, gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction. Valid values: asc, desc.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListReleaseExecutions,
		},
	}
}
