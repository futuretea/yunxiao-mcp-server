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
				mcp.WithString("idOrUsername", mcp.Required(), mcp.Description("User ID or login username.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetUser,
		},
		{
			Tool: mcp.NewTool("list_app_extension_features",
				mcp.WithDescription("List app extension feature implementations for a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("type", mcp.Required(), mcp.Description("App extension type.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppExtensionFeatures,
		},
	}
}
