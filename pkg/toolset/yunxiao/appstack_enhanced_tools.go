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
	}
}
