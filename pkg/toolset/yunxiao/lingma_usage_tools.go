package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func lingmaUsageTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_developer_members",
				mcp.WithDescription("List Tongyi Lingma developer members in a Yunxiao organization. Use this to analyze AI coding assistant adoption across teams."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("departmentId", mcp.Description("Department ID filter. Use list_organization_departments to discover valid department IDs.")),
				mcp.WithString("userId", mcp.Description("User ID filter. Use list_organization_members to discover valid user IDs.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListDeveloperMembers,
		},
		{
			Tool: mcp.NewTool("get_department_usage",
				mcp.WithDescription("Get Tongyi Lingma usage metrics for a specific department over a time range. Use this to analyze AI coding assistant adoption within a department."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("departmentId", mcp.Required(), mcp.Description("Department ID. Use list_organization_departments to discover valid department IDs.")),
				mcp.WithString("startTime", mcp.Required(), mcp.Description("Start time for the usage query range (inclusive). Format: yyyy-MM-ddTHH:mm:ss+08:00 (e.g. 2024-01-01T00:00:00+08:00).")),
				mcp.WithString("endTime", mcp.Required(), mcp.Description("End time for the usage query range (inclusive). Format: yyyy-MM-ddTHH:mm:ss+08:00.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetDepartmentUsage,
		},
		{
			Tool: mcp.NewTool("get_developer_usage",
				mcp.WithDescription("Get Tongyi Lingma usage metrics for a specific developer or department over a time range. Use this to analyze individual or team AI coding assistant usage."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("startTime", mcp.Required(), mcp.Description("Start time for the usage query range (inclusive). Format: yyyy-MM-ddTHH:mm:ss+08:00 (e.g. 2024-01-01T00:00:00+08:00).")),
				mcp.WithString("endTime", mcp.Required(), mcp.Description("End time for the usage query range (inclusive). Format: yyyy-MM-ddTHH:mm:ss+08:00.")),
				mcp.WithString("userId", mcp.Description("User ID to query usage for a specific developer. Use list_organization_members to discover valid user IDs. Either userId or departmentId must be provided.")),
				mcp.WithString("departmentId", mcp.Description("Department ID to query usage for all developers in a department. Use list_organization_departments to discover valid department IDs. Either userId or departmentId must be provided.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetDeveloperUsage,
		},
	}
}
