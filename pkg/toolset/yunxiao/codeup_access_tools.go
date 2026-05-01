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
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: created_at or updated_at.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListSSHKeys,
		},
		{
			Tool: mcp.NewTool("get_ssh_key",
				mcp.WithDescription("Get a CodeUp SSH key by ID."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("keyId", mcp.Required(), mcp.Description("SSH key ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetSSHKey,
		},
		{
			Tool: mcp.NewTool("list_user_ssh_keys",
				mcp.WithDescription("List CodeUp SSH keys for a Yunxiao user."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
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
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListWebHooks,
		},
		{
			Tool: mcp.NewTool("get_webhook",
				mcp.WithDescription("Get a CodeUp webhook by ID."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repositoryId", mcp.Required(), mcp.Description("Repository numeric ID or full path such as org/repo.")),
				mcp.WithString("hookId", mcp.Required(), mcp.Description("Webhook ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetWebHook,
		},
	}
}
