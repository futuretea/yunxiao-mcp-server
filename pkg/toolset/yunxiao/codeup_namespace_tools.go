package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func codeupNamespaceTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 2)
	tools = append(tools, codeupTemplateRepositoryTools()...)
	tools = append(tools, codeupNamespaceListTools()...)
	return tools
}

func codeupTemplateRepositoryTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_template_repositories",
				mcp.WithDescription("List CodeUp template repositories in a Yunxiao organization. Templates are pre-configured repositories used as starting points for new projects."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithNumber("templateType", mcp.Required(), mcp.Description("Template type: 1 for custom templates, 2 for built-in templates.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListTemplateRepositories,
		},
	}
}

func codeupNamespaceListTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_namespaces",
				mcp.WithDescription("List CodeUp namespaces or groups in a Yunxiao organization. Namespaces organize repositories into hierarchical groups."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithNumber("parentId", mcp.Description("Parent namespace ID. Omit to list namespaces available to the current user.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100.")),
				mcp.WithString("search", mcp.Description("Namespace search keyword.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: created_at or updated_at.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListNamespaces,
		},
		{
			Tool: mcp.NewTool("get_org_namespace",
				mcp.WithDescription("Get a CodeUp organization namespace by ID with nested sub-namespaces. Use list_namespaces to discover valid namespace IDs."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("namespaceId", mcp.Required(), mcp.Description("Namespace ID (string). Use list_namespaces to discover valid IDs.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetOrgNamespace,
		},
	}
}
