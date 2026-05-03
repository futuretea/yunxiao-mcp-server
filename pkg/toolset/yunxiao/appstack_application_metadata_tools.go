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
				mcp.WithDescription("Search AppStack application templates."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Yunxiao currently supports keyset.")),
				mcp.WithNumber("perPage", mcp.Description("Page size, up to 100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: id or gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithString("displayNameKeyword", mcp.Description("Template display name keyword.")),
				mcp.WithNumber("page", mcp.Description("Page number when using page pagination.")),
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
				mcp.WithDescription("List AppStack environments for an application."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Yunxiao currently supports keyset.")),
				mcp.WithNumber("perPage", mcp.Description("Page size, up to 100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: id or gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithNumber("page", mcp.Description("Page number when using page pagination.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListEnvironments,
		},
	}
}

func appstackApplicationAssociationTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_application_members",
				mcp.WithDescription("List members of an AppStack application."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithNumber("current", mcp.Description("Current page number. Defaults to 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size. Defaults to 10.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListApplicationMembers,
		},
		{
			Tool: mcp.NewTool("list_application_sources",
				mcp.WithDescription("List source repositories attached to an AppStack application."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Use keyset for keyset pagination.")),
				mcp.WithNumber("perPage", mcp.Description("Page size, up to 100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: id or gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithNumber("page", mcp.Description("Page number when using page pagination.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListApplicationSources,
		},
	}
}
