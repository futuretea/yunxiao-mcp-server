package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackEnhancedTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_application_overview",
				mcp.WithDescription("Get a comprehensive overview of an Appstack application including basic info, environments, and recent orchestrations in one read-only call."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application unique name.")),
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
				mcp.WithDescription("Get a comprehensive overview of an Appstack environment including basic info, variable groups, and latest orchestration in one read-only call."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application unique name.")),
				mcp.WithString("envName", mcp.Required(), mcp.Description("Environment name.")),
				mcp.WithBoolean("includeVariableGroups", mcp.Description("Whether to include environment variable groups. Defaults to true.")),
				mcp.WithBoolean("includeLatestOrchestration", mcp.Description("Whether to include the latest available orchestration. Defaults to true.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetEnvironmentOverview,
		},
		{
			Tool: mcp.NewTool("get_release_overview",
				mcp.WithDescription("Get a comprehensive overview of an Appstack system release including basic info, members, products, and attached change requests in one read-only call."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("systemName", mcp.Required(), mcp.Description("System name.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Release serial number.")),
				mcp.WithBoolean("includeMembers", mcp.Description("Whether to include release members. Defaults to true.")),
				mcp.WithBoolean("includeProducts", mcp.Description("Whether to include release products. Defaults to true.")),
				mcp.WithBoolean("includeChangeRequests", mcp.Description("Whether to include attached change requests. Defaults to true.")),
				mcp.WithNumber("changeRequestLimit", mcp.Description("Max change requests returned. Defaults to 5.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetReleaseOverview,
		},
	}
}
