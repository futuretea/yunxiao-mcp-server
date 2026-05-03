package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func platformMetadataTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_user",
				mcp.WithDescription("Get a Yunxiao user by ID or username."),
				mcp.WithString("idOrUsername", mcp.Required(), mcp.Description("User ID (numeric) or login username. Use list_users to discover valid values.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetUser,
		},
		{
			Tool: mcp.NewTool("list_app_extension_features",
				mcp.WithDescription("List app extension feature implementations for a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("type", mcp.Required(), mcp.Description("App extension type identifier. Contact your organization admin for valid type values.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppExtensionFeatures,
		},
	}
}
