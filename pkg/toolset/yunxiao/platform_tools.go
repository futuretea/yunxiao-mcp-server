package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func platformTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 25)
	tools = append(tools, platformBasicTools()...)
	tools = append(tools, platformDepartmentTools()...)
	tools = append(tools, platformEnterpriseDepartmentTools()...)
	tools = append(tools, platformMemberTools()...)
	tools = append(tools, platformGroupTools()...)
	tools = append(tools, platformRoleAndUserTools()...)
	tools = append(tools, platformAuditTools()...)
	tools = append(tools, platformMetadataTools()...)
	tools = append(tools, platformEnhancedTools()...)
	return tools
}

func platformBasicTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_current_user",
				mcp.WithDescription("Get the current Yunxiao user profile for the configured access token. Use this to verify authentication and discover the user's identity, account ID, and default organization."),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetCurrentUser,
		},
		{
			Tool: mcp.NewTool("get_current_organization_info",
				mcp.WithDescription("Get the current user's default Yunxiao organization context, including organization ID and name. Use this to discover the default organizationId before calling organization-scoped tools."),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetCurrentOrganizationInfo,
		},
		{
			Tool: mcp.NewTool("get_user_organizations",
				mcp.WithDescription("Get Yunxiao organizations visible to the current user. Use this to discover organization IDs and names when the default organization is not the desired one."),
				mcp.WithNumber("page",
					mcp.Description("Page number. Defaults to 1 when omitted by Yunxiao."),
				),
				mcp.WithNumber("perPage",
					mcp.Description("Page size from 1 to 100. Defaults to 100 when omitted by Yunxiao."),
				),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetUserOrganizations,
		},
		{
			Tool: mcp.NewTool("list_organizations",
				mcp.WithDescription("List Yunxiao organizations visible to the current user. Use this to discover organization IDs and names when the default organization is not the desired one."),
				mcp.WithNumber("page",
					mcp.Description("Page number. Defaults to 1 when omitted by Yunxiao."),
				),
				mcp.WithNumber("perPage",
					mcp.Description("Page size from 1 to 100. Defaults to 100 when omitted by Yunxiao."),
				),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListOrganizations,
		},
	}
}

func platformDepartmentTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_organization_departments",
				mcp.WithDescription("List departments in a Yunxiao organization. Use this to discover department IDs for filtering members or assigning work items."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("parentId", mcp.Description("Parent department ID. Use list_organization_departments with an empty parentId to discover top-level departments, then drill down by setting this to a discovered department ID.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListOrganizationDepartments,
		},
	}
}

func platformMemberTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_organization_members",
				mcp.WithDescription("List members in a Yunxiao organization. Use this to discover user IDs, names, and roles for assigning work items or mentioning in comments."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListOrganizationMembers,
		},
		{
			Tool: mcp.NewTool("search_organization_members",
				mcp.WithDescription("Search members in a Yunxiao organization with filters. Use this to find specific users by name, department, or role for assignment or review purposes."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithArray("deptIds", mcp.Description("Department IDs to filter by. Use list_organization_departments to discover valid department IDs."), mcp.WithStringItems()),
				mcp.WithString("query", mcp.Description("Member search query. Matches username, display name, or email.")),
				mcp.WithBoolean("includeChildren", mcp.Description("Whether to include members from child departments. Set to true for broader searches across the org tree.")),
				mcp.WithString("nextToken", mcp.Description("Pagination next token from a previous response.")),
				mcp.WithArray("roleIds", mcp.Description("Role IDs to filter by. Use list_organization_roles to discover valid role IDs."), mcp.WithStringItems()),
				mcp.WithArray("statuses", mcp.Description("Member statuses to filter by. Common values: enabled, disabled."), mcp.WithStringItems()),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleSearchOrganizationMembers,
		},
	}
}

func platformRoleAndUserTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_organization_roles",
				mcp.WithDescription("List roles defined in a Yunxiao organization. Use this to discover role IDs for filtering members or checking permissions."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListOrganizationRoles,
		},
		{
			Tool: mcp.NewTool("list_users",
				mcp.WithDescription("List Yunxiao users across organizations. Use this to discover user IDs and account information for mentions, assignments, or cross-org collaboration."),
				mcp.WithString("filter", mcp.Description("Fuzzy filter for username, login, email, or phone.")),
				mcp.WithString("status", mcp.Description("User status filter. Common values: enabled, deleted.")),
				mcp.WithString("deptId", mcp.Description("Department ID. Use list_organization_departments to discover valid IDs.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListUsers,
		},
	}
}
