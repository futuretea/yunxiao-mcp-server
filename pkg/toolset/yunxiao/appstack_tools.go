package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 48)
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
	tools = append(tools, appstackChangeRequestTools()...)
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

func appstackVariableGroupTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_env_variable_groups",
				mcp.WithDescription("Get AppStack variable groups bound to an environment."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("envName", mcp.Required(), mcp.Description("Environment name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetEnvVariableGroups,
		},
		{
			Tool: mcp.NewTool("get_variable_group",
				mcp.WithDescription("Get an AppStack application variable group by name."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("variableGroupName", mcp.Required(), mcp.Description("Variable group name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetVariableGroup,
		},
		{
			Tool: mcp.NewTool("get_app_variable_groups",
				mcp.WithDescription("Get AppStack variable groups for an application."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetAppVariableGroups,
		},
		{
			Tool: mcp.NewTool("get_app_variable_groups_revision",
				mcp.WithDescription("Get the revision of AppStack variable groups for an application."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetAppVariableGroupsRevision,
		},
	}
}

func appstackOrchestrationTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_latest_orchestration",
				mcp.WithDescription("Get the latest AppStack orchestration available for an environment."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("envName", mcp.Required(), mcp.Description("Environment name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetLatestOrchestration,
		},
		{
			Tool: mcp.NewTool("list_app_orchestration",
				mcp.WithDescription("List AppStack orchestrations for an application."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppOrchestration,
		},
		{
			Tool: mcp.NewTool("get_app_orchestration",
				mcp.WithDescription("Get an AppStack application orchestration by serial number."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Orchestration serial number.")),
				mcp.WithString("tagName", mcp.Description("Optional orchestration tag.")),
				mcp.WithString("sha", mcp.Description("Optional orchestration commit SHA.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetAppOrchestration,
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
				mcp.WithDescription("List AppStack release workflows for an application."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppReleaseWorkflows,
		},
		{
			Tool: mcp.NewTool("list_app_release_workflow_briefs",
				mcp.WithDescription("List AppStack release workflow briefs for an application."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppReleaseWorkflowBriefs,
		},
		{
			Tool: mcp.NewTool("get_app_release_workflow_stage",
				mcp.WithDescription("Get an AppStack application release workflow stage."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number.")),
				mcp.WithString("releaseStageSn", mcp.Required(), mcp.Description("Release stage serial number.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetAppReleaseWorkflowStage,
		},
		{
			Tool: mcp.NewTool("list_app_release_stage_briefs",
				mcp.WithDescription("List AppStack release stage briefs for an application release workflow."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number.")),
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
				mcp.WithDescription("List AppStack release stage execution records."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number.")),
				mcp.WithString("releaseStageSn", mcp.Required(), mcp.Description("Release stage serial number.")),
				mcp.WithString("pagination", mcp.Description("Pagination mode. Use keyset for keyset pagination.")),
				mcp.WithNumber("perPage", mcp.Description("Page size, up to 100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: id or gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithString("nextToken", mcp.Description("Keyset pagination token from the previous response.")),
				mcp.WithNumber("page", mcp.Description("Page number when using page pagination.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppReleaseStageRuns,
		},
		{
			Tool: mcp.NewTool("list_app_release_stage_exec_metadata",
				mcp.WithDescription("List integrated change metadata for an AppStack release stage execution."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number.")),
				mcp.WithString("releaseStageSn", mcp.Required(), mcp.Description("Release stage serial number.")),
				mcp.WithString("executionNumber", mcp.Required(), mcp.Description("Release stage execution number.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppReleaseStageExecMetadata,
		},
		{
			Tool: mcp.NewTool("get_app_release_stage_pipeline_run",
				mcp.WithDescription("Get the pipeline run for an AppStack release stage execution."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number.")),
				mcp.WithString("releaseStageSn", mcp.Required(), mcp.Description("Release stage serial number.")),
				mcp.WithString("executionNumber", mcp.Required(), mcp.Description("Release stage execution number.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetAppReleaseStagePipelineRun,
		},
		{
			Tool: mcp.NewTool("get_app_release_stage_pipeline_job_log",
				mcp.WithDescription("Get a pipeline job log for an AppStack release stage execution."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number.")),
				mcp.WithString("releaseStageSn", mcp.Required(), mcp.Description("Release stage serial number.")),
				mcp.WithString("executionNumber", mcp.Required(), mcp.Description("Release stage execution number.")),
				mcp.WithString("jobId", mcp.Required(), mcp.Description("Pipeline job ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetAppReleaseStagePipelineJobLog,
		},
	}
}

func appstackChangeRequestTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_appstack_change_request_audit_items",
				mcp.WithDescription("Get audit items for an AppStack change request."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Change request serial number.")),
				mcp.WithString("refType", mcp.Required(), mcp.Description("Reference type, such as CR.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetAppStackChangeRequestAuditItems,
		},
		{
			Tool: mcp.NewTool("list_appstack_change_request_executions",
				mcp.WithDescription("List execution records for an AppStack change request."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Change request serial number.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number.")),
				mcp.WithString("releaseStageSn", mcp.Required(), mcp.Description("Release stage serial number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size, up to 100.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: id or gmtCreate.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppStackChangeRequestExecutions,
		},
		{
			Tool: mcp.NewTool("list_appstack_change_request_work_items",
				mcp.WithDescription("List work items for an AppStack change request."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Change request serial number.")),
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
