package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func lingmaUsageTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_department_usage",
				mcp.WithDescription("Get Tongyi Lingma daily usage data for a department."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("departmentId", mcp.Required(), mcp.Description("Department ID.")),
				mcp.WithString("startTime", mcp.Required(), mcp.Description("Start date in YYYY-MM-DD format.")),
				mcp.WithString("endTime", mcp.Required(), mcp.Description("End date in YYYY-MM-DD format.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size. Default is 100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetDepartmentUsage,
		},
		{
			Tool: mcp.NewTool("list_developer_members",
				mcp.WithDescription("List Tongyi Lingma developer members."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("departmentId", mcp.Description("Department ID filter.")),
				mcp.WithString("userId", mcp.Description("User ID filter.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size. Default is 100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListDeveloperMembers,
		},
		{
			Tool: mcp.NewTool("get_developer_usage",
				mcp.WithDescription("Get Tongyi Lingma daily usage data for developers."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("userId", mcp.Description("User ID. Either userId or departmentId is required.")),
				mcp.WithString("departmentId", mcp.Description("Department ID. Either userId or departmentId is required.")),
				mcp.WithString("startTime", mcp.Required(), mcp.Description("Start date in YYYY-MM-DD format.")),
				mcp.WithString("endTime", mcp.Required(), mcp.Description("End date in YYYY-MM-DD format.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size. Default is 100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetDeveloperUsage,
		},
	}
}
