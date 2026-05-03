package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 53)
	tools = append(tools, appstackApplicationTools()...)
	tools = append(tools, appstackApplicationMetadataTools()...)
	tools = append(tools, appstackDeploymentResourceTools()...)
	tools = append(tools, appstackResourceProxyTools()...)
	tools = append(tools, appstackGlobalVarTools()...)
	tools = append(tools, appstackVariableGroupTools()...)
	tools = append(tools, appstackOrchestrationTools()...)
	tools = append(tools, appstackAppReleaseWorkflowTools()...)
	tools = append(tools, appstackSystemTools()...)
	tools = append(tools, appstackSystemReleaseTools()...)
	tools = append(tools, appstackReleaseSearchTools()...)
	tools = append(tools, appstackChangeRequestTools()...)
	tools = append(tools, appstackChangeOrderTools()...)
	tools = append(tools, appstackEnhancedTools()...)
	return tools
}

func appstackApplicationTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_applications",
				mcp.WithDescription("List AppStack applications in a Yunxiao organization. AppStack manages deployment environments and release pipelines."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Valid value: keyset.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithString("orderBy", mcp.Description("Sort field. Valid values: id, gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction. Valid values: asc, desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListApplications,
		},
	}
}

func appstackVariableGroupTools() []toolset.ServerTool {
	return nil
}

func appstackOrchestrationTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_app_orchestration",
				mcp.WithDescription("List AppStack orchestrations for an application. Orchestrations define deployment workflows and environment configurations."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppOrchestration,
		},
	}
}

func appstackAppReleaseWorkflowTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 8)
	tools = append(tools, appstackAppReleaseWorkflowOverviewTools()...)
	tools = append(tools, appstackAppReleaseStageExecutionTools()...)
	return tools
}

func appstackAppReleaseWorkflowOverviewTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_app_release_workflows",
				mcp.WithDescription("List AppStack release workflows for an application. Release workflows model multi-stage deployment pipelines."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppReleaseWorkflows,
		},
		{
			Tool: mcp.NewTool("list_app_release_workflow_briefs",
				mcp.WithDescription("List AppStack release workflow briefs for an application. Briefs provide a condensed view of workflow definitions."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppReleaseWorkflowBriefs,
		},
		{
			Tool: mcp.NewTool("list_app_release_stage_briefs",
				mcp.WithDescription("List AppStack release stage briefs for an application release workflow. Stages represent individual deployment phases within a workflow."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number. Use list_app_release_workflows to discover valid values.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppReleaseStageBriefs,
		},
	}
}

func appstackAppReleaseStageExecutionTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_app_release_stage_runs",
				mcp.WithDescription("List AppStack release stage execution records. Each record represents a single run of a deployment stage."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number. Use list_app_release_workflows to discover valid values.")),
				mcp.WithString("releaseStageSn", mcp.Required(), mcp.Description("Release stage serial number. Use list_app_release_stage_briefs to discover valid values.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Valid value: keyset.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithString("orderBy", mcp.Description("Sort field. Valid values: id, gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction. Valid values: asc, desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppReleaseStageRuns,
		},
		{
			Tool: mcp.NewTool("list_app_release_stage_exec_metadata",
				mcp.WithDescription("List integrated change metadata for an AppStack release stage execution. Metadata includes linked work items and commits."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number. Use list_app_release_workflows to discover valid values.")),
				mcp.WithString("releaseStageSn", mcp.Required(), mcp.Description("Release stage serial number. Use list_app_release_stage_briefs to discover valid values.")),
				mcp.WithString("executionNumber", mcp.Required(), mcp.Description("Release stage execution number. Use list_app_release_stage_runs to discover valid values.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppReleaseStageExecMetadata,
		},
	}
}

func appstackChangeRequestTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_appstack_change_request_executions",
				mcp.WithDescription("List execution records for an AppStack change request. Change requests track planned deployments."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Change request serial number. Typically discovered via list_attached_change_requests.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number. Use list_app_release_workflows to discover valid values.")),
				mcp.WithString("releaseStageSn", mcp.Required(), mcp.Description("Release stage serial number. Use list_app_release_stage_briefs to discover valid values.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithString("orderBy", mcp.Description("Sort field. Valid values: id, gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction. Valid values: asc, desc.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppStackChangeRequestExecutions,
		},
		{
			Tool: mcp.NewTool("list_appstack_change_request_work_items",
				mcp.WithDescription("List work items for an AppStack change request. Work items represent linked Projex tasks or requirements."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Change request serial number. Typically discovered via list_attached_change_requests.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppStackChangeRequestWorkItems,
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
				mcp.WithDescription("List AppStack change order versions. Change orders track actual deployments and their versions."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("envNames", mcp.Description("Comma-separated environment names. Use list_environments to discover valid environment names for an application.")),
				mcp.WithString("creators", mcp.Description("Comma-separated creator account IDs.")),
				mcp.WithNumber("current", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListChangeOrderVersions,
		},
		{
			Tool: mcp.NewTool("list_change_orders_by_origin",
				mcp.WithDescription("List AppStack change orders by creation origin. Use this to trace deployments back to their source (e.g., a Flow pipeline)."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("originType", mcp.Required(), mcp.Description("Origin type. Valid values: FLOW.")),
				mcp.WithString("originId", mcp.Required(), mcp.Description("Origin identifier.")),
				mcp.WithString("appName", mcp.Description("Application name filter. Use list_applications to discover valid app names.")),
				mcp.WithString("envName", mcp.Description("Environment name filter. Use list_environments to discover valid environment names for an application.")),
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
				mcp.WithDescription("List logs for an AppStack change order job. Job logs capture deployment script output."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("changeOrderSn", mcp.Required(), mcp.Description("Change order serial number. Use list_change_order_versions to discover valid values.")),
				mcp.WithString("jobSn", mcp.Required(), mcp.Description("Change order job serial number. Typically returned in change order details.")),
				mcp.WithNumber("current", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("pageSize", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListChangeOrderJobLogs,
		},
		{
			Tool: mcp.NewTool("find_task_operation_log",
				mcp.WithDescription("Get an AppStack deployment task operation log. Operation logs record manual or automated actions taken during deployment."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("changeOrderSn", mcp.Required(), mcp.Description("Change order serial number. Use list_change_order_versions to discover valid values.")),
				mcp.WithString("jobSn", mcp.Required(), mcp.Description("Change order job serial number. Typically returned in change order details.")),
				mcp.WithString("stageSn", mcp.Required(), mcp.Description("Deployment stage serial number. Typically returned in change order job details.")),
				mcp.WithString("taskSn", mcp.Required(), mcp.Description("Deployment task serial number. Typically returned in change order stage details.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleFindTaskOperationLog,
		},
	}
}
