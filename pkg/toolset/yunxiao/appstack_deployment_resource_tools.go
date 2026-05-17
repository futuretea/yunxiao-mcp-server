package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackDeploymentResourceTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_machine_deploy_log",
				mcp.WithDescription("Get deployment log for a specific machine in an AppStack deployment. Machine logs capture the agent-side output of a deployment."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("tunnelId", mcp.Required(), mcp.Description("Deployment tunnel ID. Typically discovered from change order or deployment details.")),
				mcp.WithString("machineSn", mcp.Required(), mcp.Description("Machine serial number. Use deployment details to discover valid machine identifiers.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetMachineDeployLog,
		},
		{
			Tool: mcp.NewTool("get_deploy_group",
				mcp.WithDescription("Get an AppStack deploy group by name within a resource pool. Deploy groups define subsets of machines for targeted deployments."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("poolName", mcp.Required(), mcp.Description("Resource pool name. Typically discovered from application environment details.")),
				mcp.WithString("deployGroupName", mcp.Required(), mcp.Description("Deploy group name. Typically discovered from deployment configuration.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetDeployGroup,
		},
		{
			Tool: mcp.NewTool("list_resource_instances",
				mcp.WithDescription("List AppStack resource instances in a resource pool. Pool names are typically found in application environment resource configurations."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("poolName", mcp.Required(), mcp.Description("Resource pool name. Typically found in application environment resource configurations.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Valid value: keyset.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithString("orderBy", mcp.Description("Sort field. Valid values: id, gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction. Valid values: asc, desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListResourceInstances,
		},
		{
			Tool: mcp.NewTool("get_resource_instance",
				mcp.WithDescription("Get an AppStack resource instance by name within a resource pool. Resource instances represent individual machines or hosts."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("poolName", mcp.Required(), mcp.Description("Resource pool name. Use list_resource_instances to discover valid pool names.")),
				mcp.WithString("instanceName", mcp.Required(), mcp.Description("Resource instance name. Use list_resource_instances to discover valid instance names.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetResourceInstance,
		},
	}
}
