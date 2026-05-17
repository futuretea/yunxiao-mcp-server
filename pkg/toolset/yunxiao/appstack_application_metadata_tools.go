package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackApplicationMetadataTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 5)
	tools = append(tools, appstackTemplateTools()...)
	tools = append(tools, appstackEnvironmentTools()...)
	tools = append(tools, appstackApplicationAssociationTools()...)
	return tools
}

func appstackTemplateTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("search_app_templates",
				mcp.WithDescription("Search AppStack application templates. Use this to discover templates before creating or deploying applications."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Valid value: keyset.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithString("orderBy", mcp.Description("Sort field. Valid values: id, gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction. Valid values: asc, desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithString("displayNameKeyword", mcp.Description("Template display name keyword for filtering results.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleSearchAppTemplates,
		},
	}
}

func appstackEnvironmentTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_environments",
				mcp.WithDescription("List AppStack environments for an application. Use list_applications to discover valid application names."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid names.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Valid value: keyset.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithString("orderBy", mcp.Description("Sort field. Valid values: id, gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction. Valid values: asc, desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListEnvironments,
		},
		{
			Tool: mcp.NewTool("get_environment",
				mcp.WithDescription("Get a single AppStack environment by name. Use list_environments or get_application_overview to discover valid environment names."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid names.")),
				mcp.WithString("envName", mcp.Required(), mcp.Description("Environment name. Use list_environments to discover valid names.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetEnvironment,
		},
	}
}

func appstackApplicationAssociationTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_application_members",
				mcp.WithDescription("List members of an AppStack application. Use list_applications to discover valid application names."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid names.")),
				mcp.WithNumber("current", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListApplicationMembers,
		},
		{
			Tool: mcp.NewTool("list_application_sources",
				mcp.WithDescription("List source repositories attached to an AppStack application. Use list_applications to discover valid application names."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid names.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Valid value: keyset.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithString("orderBy", mcp.Description("Sort field. Valid values: id, gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction. Valid values: asc, desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListApplicationSources,
		},
	}
}
