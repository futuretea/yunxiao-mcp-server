package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func codeupAccessTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 5)
	tools = append(tools, codeupSSHKeyTools()...)
	tools = append(tools, codeupWebhookTools()...)
	return tools
}

func codeupSSHKeyTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_ssh_keys",
				mcp.WithDescription("List CodeUp SSH keys in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: created_at or updated_at.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListSSHKeys,
		},
		{
			Tool: mcp.NewTool("list_user_ssh_keys",
				mcp.WithDescription("List CodeUp SSH keys for a Yunxiao user."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("userId", mcp.Required(), mcp.Description("Yunxiao user ID.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: created_at or updated_at.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListUserSSHKeys,
		},
	}
}

func codeupWebhookTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_webhooks",
				mcp.WithDescription("List CodeUp webhooks in a repository."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListWebHooks,
		},
	}
}
