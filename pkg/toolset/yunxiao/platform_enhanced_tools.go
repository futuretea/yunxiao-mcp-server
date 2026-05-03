package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func platformEnhancedTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_organization_overview",
				mcp.WithDescription("Get a comprehensive overview of a Yunxiao organization including basic info, departments, members, groups, and roles in one read-only call."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithBoolean("includeDepartments", mcp.Description("Whether to include departments list. Defaults to true.")),
				mcp.WithBoolean("includeMembers", mcp.Description("Whether to include members list. Defaults to true.")),
				mcp.WithBoolean("includeGroups", mcp.Description("Whether to include groups list. Defaults to true.")),
				mcp.WithBoolean("includeRoles", mcp.Description("Whether to include roles list. Defaults to true.")),
				mcp.WithNumber("departmentLimit", mcp.Description("Maximum departments to include in the overview. Defaults to 5. Set to 0 to exclude.")),
				mcp.WithNumber("memberLimit", mcp.Description("Maximum members to include in the overview. Defaults to 5. Set to 0 to exclude.")),
				mcp.WithNumber("groupLimit", mcp.Description("Maximum groups to include in the overview. Defaults to 5. Set to 0 to exclude.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetOrganizationOverview,
		},
		{
			Tool: mcp.NewTool("get_organization_department_overview",
				mcp.WithDescription("Get a comprehensive overview of a Yunxiao organization department including basic info and ancestor chain in one read-only call."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("departmentId", mcp.Required(), mcp.Description("Department ID. Use list_organization_departments or list_enterprise_departments to discover valid IDs.")),
				mcp.WithBoolean("includeAncestors", mcp.Description("Whether to include the ancestor chain. Defaults to true.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetOrganizationDepartmentOverview,
		},
		{
			Tool: mcp.NewTool("get_organization_group_overview",
				mcp.WithDescription("Get a comprehensive overview of a Yunxiao organization group including basic info and members in one read-only call."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("groupId", mcp.Required(), mcp.Description("Group ID. Use list_organization_groups to discover valid group IDs.")),
				mcp.WithBoolean("includeMembers", mcp.Description("Whether to include group members. Defaults to true.")),
				mcp.WithNumber("memberLimit", mcp.Description("Maximum members to include in the overview. Defaults to 5. Set to 0 to exclude.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetOrganizationGroupOverview,
		},
	}
}
