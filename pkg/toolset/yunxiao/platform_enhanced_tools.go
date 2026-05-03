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
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithBoolean("includeDepartments", mcp.Description("Whether to include departments list. Defaults to true.")),
				mcp.WithBoolean("includeMembers", mcp.Description("Whether to include members list. Defaults to true.")),
				mcp.WithBoolean("includeGroups", mcp.Description("Whether to include groups list. Defaults to true.")),
				mcp.WithBoolean("includeRoles", mcp.Description("Whether to include roles list. Defaults to true.")),
				mcp.WithNumber("departmentLimit", mcp.Description("Max departments returned. Defaults to 5.")),
				mcp.WithNumber("memberLimit", mcp.Description("Max members returned. Defaults to 5.")),
				mcp.WithNumber("groupLimit", mcp.Description("Max groups returned. Defaults to 5.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetOrganizationOverview,
		},
		{
			Tool: mcp.NewTool("get_organization_department_overview",
				mcp.WithDescription("Get a comprehensive overview of a Yunxiao organization department including basic info and ancestor chain in one read-only call."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("departmentId", mcp.Required(), mcp.Description("Department ID.")),
				mcp.WithBoolean("includeAncestors", mcp.Description("Whether to include the ancestor chain. Defaults to true.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetOrganizationDepartmentOverview,
		},
	}
}
