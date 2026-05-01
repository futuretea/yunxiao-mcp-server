package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 7)
	tools = append(tools, appstackApplicationTools()...)
	tools = append(tools, appstackChangeOrderTools()...)
	return tools
}

func appstackApplicationTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_applications",
				mcp.WithDescription("List AppStack applications in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Yunxiao currently supports keyset.")),
				mcp.WithNumber("perPage", mcp.Description("Page size, up to 100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: id or gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithNumber("page", mcp.Description("Page number when using page pagination.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListApplications,
		},
		{
			Tool: mcp.NewTool("get_application",
				mcp.WithDescription("Get an AppStack application by name."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetApplication,
		},
	}
}

func appstackChangeOrderTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 5)
	tools = append(tools, appstackChangeOrderSummaryTools()...)
	tools = append(tools, appstackChangeOrderLogTools()...)
	return tools
}

func appstackChangeOrderSummaryTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_change_order_versions",
				mcp.WithDescription("List AppStack change order versions."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("envNames", mcp.Description("Comma-separated environment names.")),
				mcp.WithString("creators", mcp.Description("Comma-separated creator account IDs.")),
				mcp.WithNumber("current", mcp.Description("Current page number.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListChangeOrderVersions,
		},
		{
			Tool: mcp.NewTool("get_change_order",
				mcp.WithDescription("Get an AppStack change order by serial number."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("changeOrderSn", mcp.Required(), mcp.Description("Change order serial number.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetChangeOrder,
		},
		{
			Tool: mcp.NewTool("list_change_orders_by_origin",
				mcp.WithDescription("List AppStack change orders by creation origin."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("originType", mcp.Required(), mcp.Description("Origin type, such as FLOW.")),
				mcp.WithString("originId", mcp.Required(), mcp.Description("Origin identifier.")),
				mcp.WithString("appName", mcp.Description("Application name filter.")),
				mcp.WithString("envName", mcp.Description("Environment name filter.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListChangeOrdersByOrigin,
		},
	}
}

func appstackChangeOrderLogTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_change_order_job_logs",
				mcp.WithDescription("List logs for an AppStack change order job."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("changeOrderSn", mcp.Required(), mcp.Description("Change order serial number.")),
				mcp.WithString("jobSn", mcp.Required(), mcp.Description("Change order job serial number.")),
				mcp.WithNumber("current", mcp.Description("Current page number.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListChangeOrderJobLogs,
		},
		{
			Tool: mcp.NewTool("find_task_operation_log",
				mcp.WithDescription("Get an AppStack deployment task operation log."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("changeOrderSn", mcp.Required(), mcp.Description("Change order serial number.")),
				mcp.WithString("jobSn", mcp.Required(), mcp.Description("Change order job serial number.")),
				mcp.WithString("stageSn", mcp.Required(), mcp.Description("Deployment stage serial number.")),
				mcp.WithString("taskSn", mcp.Required(), mcp.Description("Deployment task serial number.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleFindTaskOperationLog,
		},
	}
}
