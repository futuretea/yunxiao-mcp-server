package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackWriteTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("create_change_order",
				mcp.WithDescription("Create an AppStack change order (deployment order). Change orders trigger application deployments to environments. This is a write operation and requires read_only=false."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("changeOrder", mcp.Required(), mcp.Description("JSON string with change order details: {changeOrderName, type (Deploy|Scale|Rollback|Destroy), envs (object), orchestrationRevisionSha, description}.")),
			),
			Handler: handleCreateChangeOrder,
		},
		{
			Tool: mcp.NewTool("execute_job_action",
				mcp.WithDescription("Execute an action on an AppStack change order job. Use this to suspend, resume, rollback, or stop a deployment job. This is a write operation and requires read_only=false."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("changeOrderSn", mcp.Required(), mcp.Description("Change order serial number. Use list_change_order_versions to discover valid values.")),
				mcp.WithString("jobSn", mcp.Required(), mcp.Description("Job serial number. Typically returned in change order job details.")),
				mcp.WithString("action", mcp.Required(), mcp.Description("JSON string with job action: {actionType: SUSPEND|RESUME|ROLLBACK|STOP}.")),
			),
			Handler: handleExecuteJobAction,
		},
	}
}
