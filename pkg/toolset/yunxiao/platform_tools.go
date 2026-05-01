package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func platformTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 23)
	tools = append(tools, platformBasicTools()...)
	tools = append(tools, platformDepartmentTools()...)
	tools = append(tools, platformEnterpriseDepartmentTools()...)
	tools = append(tools, platformMemberTools()...)
	tools = append(tools, platformGroupTools()...)
	tools = append(tools, platformRoleAndUserTools()...)
	tools = append(tools, platformAuditTools()...)
	tools = append(tools, platformMetadataTools()...)
	return tools
}

func platformBasicTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_current_user",
				mcp.WithDescription("Get the current Yunxiao user for the configured access token."),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetCurrentUser,
		},
		{
			Tool: mcp.NewTool("get_current_organization_info",
				mcp.WithDescription("Get current user context, including the last organization returned by Yunxiao."),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetCurrentOrganizationInfo,
		},
		{
			Tool: mcp.NewTool("get_user_organizations",
				mcp.WithDescription("Get Yunxiao organizations visible to the current user."),
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
				mcp.WithDescription("List Yunxiao organizations visible to the current user."),
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
		{
			Tool: mcp.NewTool("get_organization",
				mcp.WithDescription("Get a Yunxiao organization by ID."),
				mcp.WithString("id",
					mcp.Required(),
					mcp.Description("Organization ID."),
				),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetOrganization,
		},
	}
}

func platformDepartmentTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_organization_departments",
				mcp.WithDescription("List departments in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("parentId", mcp.Description("Parent department ID.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListOrganizationDepartments,
		},
		{
			Tool: mcp.NewTool("get_organization_department_info",
				mcp.WithDescription("Get Yunxiao organization department details."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Department ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetOrganizationDepartmentInfo,
		},
		{
			Tool: mcp.NewTool("get_organization_department_ancestors",
				mcp.WithDescription("List ancestor departments for a Yunxiao organization department."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Department ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetOrganizationDepartmentAncestors,
		},
	}
}

func platformMemberTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_organization_members",
				mcp.WithDescription("List members in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListOrganizationMembers,
		},
		{
			Tool: mcp.NewTool("get_organization_member_info",
				mcp.WithDescription("Get Yunxiao organization member details by member ID."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("memberId", mcp.Required(), mcp.Description("Organization member ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetOrganizationMemberInfo,
		},
		{
			Tool: mcp.NewTool("get_organization_member_info_by_user_id",
				mcp.WithDescription("Get Yunxiao organization member details by user ID."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("userId", mcp.Required(), mcp.Description("Yunxiao user ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetOrganizationMemberInfoByUserID,
		},
		{
			Tool: mcp.NewTool("search_organization_members",
				mcp.WithDescription("Search members in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithArray("deptIds", mcp.Description("Department IDs."), mcp.WithStringItems()),
				mcp.WithString("query", mcp.Description("Member search query.")),
				mcp.WithBoolean("includeChildren", mcp.Description("Whether to include child departments.")),
				mcp.WithString("nextToken", mcp.Description("Pagination next token.")),
				mcp.WithArray("roleIds", mcp.Description("Role IDs."), mcp.WithStringItems()),
				mcp.WithArray("statuses", mcp.Description("Member statuses."), mcp.WithStringItems()),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
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
				mcp.WithDescription("List roles in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListOrganizationRoles,
		},
		{
			Tool: mcp.NewTool("get_organization_role",
				mcp.WithDescription("Get a Yunxiao organization role by ID."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("roleId", mcp.Required(), mcp.Description("Organization role ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetOrganizationRole,
		},
		{
			Tool: mcp.NewTool("list_users",
				mcp.WithDescription("List Yunxiao users."),
				mcp.WithString("filter", mcp.Description("Fuzzy filter for username, login, email, or phone.")),
				mcp.WithString("status", mcp.Description("User status: enabled or deleted.")),
				mcp.WithString("deptId", mcp.Description("Department ID.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListUsers,
		},
	}
}
