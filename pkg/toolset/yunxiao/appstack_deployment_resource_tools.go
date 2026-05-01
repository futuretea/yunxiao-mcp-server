package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackDeploymentResourceTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 4)
	tools = append(tools, appstackDeploymentLogTools()...)
	tools = append(tools, appstackResourcePoolTools()...)
	return tools
}

func appstackDeploymentLogTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_machine_deploy_log",
				mcp.WithDescription("Get an AppStack machine deployment log."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithNumber("tunnelId", mcp.Required(), mcp.Description("Deployment tunnel ID.")),
				mcp.WithString("machineSn", mcp.Required(), mcp.Description("Machine serial number.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetMachineDeployLog,
		},
	}
}

func appstackResourcePoolTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_deploy_group",
				mcp.WithDescription("Get an AppStack deploy group by pool and group name."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("poolName", mcp.Required(), mcp.Description("Resource pool name.")),
				mcp.WithString("deployGroupName", mcp.Required(), mcp.Description("Deploy group name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetDeployGroup,
		},
		{
			Tool: mcp.NewTool("list_resource_instances",
				mcp.WithDescription("List AppStack resource instances in a resource pool."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("poolName", mcp.Required(), mcp.Description("Resource pool name.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Yunxiao currently supports keyset.")),
				mcp.WithNumber("perPage", mcp.Description("Page size, up to 100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: id or gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithNumber("page", mcp.Description("Page number when using page pagination.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListResourceInstances,
		},
		{
			Tool: mcp.NewTool("get_resource_instance",
				mcp.WithDescription("Get an AppStack resource instance by pool and instance name."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("poolName", mcp.Required(), mcp.Description("Resource pool name.")),
				mcp.WithString("instanceName", mcp.Required(), mcp.Description("Resource instance name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetResourceInstance,
		},
	}
}
