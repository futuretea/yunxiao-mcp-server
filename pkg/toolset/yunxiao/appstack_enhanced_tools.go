package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackEnhancedTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_application_overview",
				mcp.WithDescription("Get a comprehensive overview of an Appstack application including basic info, environments, and recent orchestrations in one read-only call. Use this after discovering applications via list_applications or list_attached_apps."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application unique name. Use list_applications or list_attached_apps to discover valid names.")),
				mcp.WithBoolean("includeEnvironments", mcp.Description("Whether to include environment list. Defaults to true.")),
				mcp.WithBoolean("includeOrchestrations", mcp.Description("Whether to include recent orchestrations. Defaults to true.")),
				mcp.WithNumber("envLimit", mcp.Description("Max environments returned. Defaults to 5.")),
				mcp.WithNumber("orchestrationLimit", mcp.Description("Max orchestrations returned. Defaults to 5.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetApplicationOverview,
		},
		{
			Tool: mcp.NewTool("get_environment_overview",
				mcp.WithDescription("Get a comprehensive overview of an Appstack environment including basic info, variable groups, and latest orchestration in one read-only call. Use this after identifying an application and environment via get_application_overview."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application unique name. Use list_applications or list_attached_apps to discover valid names.")),
				mcp.WithString("envName", mcp.Required(), mcp.Description("Environment name. Use get_application_overview to discover valid environment names.")),
				mcp.WithBoolean("includeVariableGroups", mcp.Description("Whether to include environment variable groups. Defaults to true.")),
				mcp.WithBoolean("includeLatestOrchestration", mcp.Description("Whether to include the latest available orchestration. Defaults to true.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetEnvironmentOverview,
		},
		{
			Tool: mcp.NewTool("get_release_overview",
				mcp.WithDescription("Get a comprehensive overview of an Appstack system release including basic info, members, products, and attached change requests in one read-only call. Use this after discovering a release via search_releases or list_system_release_workflows."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name. Use list_systems to discover valid names.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Release serial number. Use search_releases or list_system_release_workflows to discover valid serial numbers.")),
				mcp.WithBoolean("includeMembers", mcp.Description("Whether to include release members. Defaults to true.")),
				mcp.WithBoolean("includeProducts", mcp.Description("Whether to include release products. Defaults to true.")),
				mcp.WithBoolean("includeChangeRequests", mcp.Description("Whether to include attached change requests. Defaults to true.")),
				mcp.WithNumber("changeRequestLimit", mcp.Description("Max change requests returned. Defaults to 5.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetReleaseOverview,
		},
		{
			Tool: mcp.NewTool("get_system_overview",
				mcp.WithDescription("Get a comprehensive overview of an Appstack system including basic info, attached applications, and members in one read-only call. Use this after discovering a system via list_systems."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name. Use list_systems to discover valid names.")),
				mcp.WithBoolean("includeApps", mcp.Description("Whether to include attached applications. Defaults to true.")),
				mcp.WithBoolean("includeMembers", mcp.Description("Whether to include system members. Defaults to true.")),
				mcp.WithNumber("appLimit", mcp.Description("Max attached applications returned. Defaults to 10.")),
				mcp.WithNumber("memberLimit", mcp.Description("Max members returned. Defaults to 10.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetSystemOverview,
		},
		{
			Tool: mcp.NewTool("get_change_order_overview",
				mcp.WithDescription("Get a comprehensive overview of an Appstack change order including basic info and job list in one read-only call. Use this after discovering a change order via list_change_order_versions."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("changeOrderSn", mcp.Required(), mcp.Description("Change order serial number. Use list_change_order_versions to discover valid values.")),
				mcp.WithBoolean("includeJobLogs", mcp.Description("Whether to include job list. Defaults to true.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetChangeOrderOverview,
		},
	}
}
