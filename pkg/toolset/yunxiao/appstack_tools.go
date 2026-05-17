package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 54)
	tools = append(tools, appstackApplicationTools()...)
	tools = append(tools, appstackApplicationMetadataTools()...)
	tools = append(tools, appstackDeploymentResourceTools()...)
	tools = append(tools, appstackResourceProxyTools()...)
	tools = append(tools, appstackGlobalVarTools()...)
	tools = append(tools, appstackTagTools()...)
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
		{
			Tool: mcp.NewTool("get_application",
				mcp.WithDescription("Get a single AppStack application by name. Use list_applications to discover valid application names."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid names.")),
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
				mcp.WithDescription("Get AppStack variable groups for a specific environment. Variable groups define environment-specific configuration values."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("envName", mcp.Required(), mcp.Description("Environment name. Use list_environments to discover valid environment names.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetEnvVariableGroups,
		},
		{
			Tool: mcp.NewTool("get_variable_group",
				mcp.WithDescription("Get an AppStack variable group by name. Variable groups contain key-value configuration pairs for deployment environments."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("variableGroupName", mcp.Required(), mcp.Description("Variable group name. Use get_app_variable_groups to discover valid names.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetVariableGroup,
		},
		{
			Tool: mcp.NewTool("get_app_variable_groups",
				mcp.WithDescription("List AppStack variable groups for an application. Use this to discover variable group names before reading specific groups."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetAppVariableGroups,
		},
		{
			Tool: mcp.NewTool("get_app_variable_groups_revision",
				mcp.WithDescription("Get the revision metadata for AppStack variable groups of an application. Use this to check if variable groups have been updated."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
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
				mcp.WithDescription("Get the latest available AppStack orchestration for an application environment. Orchestrations define the deployment configuration and resource layout."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("envName", mcp.Required(), mcp.Description("Environment name. Use list_environments to discover valid environment names.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetLatestOrchestration,
		},
		{
			Tool: mcp.NewTool("list_app_orchestration",
				mcp.WithDescription("List AppStack orchestrations for an application. Orchestrations define deployment workflows and environment configurations."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAppOrchestration,
		},
		{
			Tool: mcp.NewTool("get_app_orchestration",
				mcp.WithDescription("Get an AppStack orchestration by serial number. Use list_app_orchestration to discover valid serial numbers."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Orchestration serial number. Use list_app_orchestration to discover valid serial numbers.")),
				mcp.WithString("tagName", mcp.Description("Optional tag name to retrieve a tagged version of the orchestration.")),
				mcp.WithString("sha", mcp.Description("Optional SHA to retrieve a specific version of the orchestration.")),
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
			Tool: mcp.NewTool("get_app_release_workflow_stage",
				mcp.WithDescription("Get an AppStack release workflow stage by serial number. Stages represent individual deployment phases with configuration details."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number. Use list_app_release_workflows to discover valid values.")),
				mcp.WithString("releaseStageSn", mcp.Required(), mcp.Description("Release stage serial number. Use list_app_release_stage_briefs to discover valid values.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetAppReleaseWorkflowStage,
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
			Tool: mcp.NewTool("get_app_release_stage_pipeline_run",
				mcp.WithDescription("Get the Flow pipeline run associated with an AppStack release stage execution. Pipeline runs show CI/CD pipeline execution details."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number. Use list_app_release_workflows to discover valid values.")),
				mcp.WithString("releaseStageSn", mcp.Required(), mcp.Description("Release stage serial number. Use list_app_release_stage_briefs to discover valid values.")),
				mcp.WithString("executionNumber", mcp.Required(), mcp.Description("Release stage execution number. Use list_app_release_stage_runs to discover valid values.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetAppReleaseStagePipelineRun,
		},
		{
			Tool: mcp.NewTool("get_app_release_stage_pipeline_job_log",
				mcp.WithDescription("Get the pipeline job log for an AppStack release stage execution. Job logs contain detailed CI/CD pipeline output for a specific job."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("releaseWorkflowSn", mcp.Required(), mcp.Description("Release workflow serial number. Use list_app_release_workflows to discover valid values.")),
				mcp.WithString("releaseStageSn", mcp.Required(), mcp.Description("Release stage serial number. Use list_app_release_stage_briefs to discover valid values.")),
				mcp.WithString("executionNumber", mcp.Required(), mcp.Description("Release stage execution number. Use list_app_release_stage_runs to discover valid values.")),
				mcp.WithString("jobId", mcp.Required(), mcp.Description("Pipeline job ID. Typically discovered from the pipeline run details.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetAppReleaseStagePipelineJobLog,
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
			Tool: mcp.NewTool("get_appstack_change_request_audit_items",
				mcp.WithDescription("Get audit items for an AppStack change request. Audit items represent approval checkpoints that must be satisfied before deployment."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("sn", mcp.Required(), mcp.Description("Change request serial number. Typically discovered via list_attached_change_requests.")),
				mcp.WithString("refType", mcp.Required(), mcp.Description("Reference type for audit items. Valid values: RELEASE, CHANGE_REQUEST.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetAppStackChangeRequestAuditItems,
		},
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
			Tool: mcp.NewTool("get_change_order",
				mcp.WithDescription("Get an AppStack change order by serial number. Change orders track actual deployments with full details."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("changeOrderSn", mcp.Required(), mcp.Description("Change order serial number. Use list_change_order_versions to discover valid values.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetChangeOrder,
		},
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
				mcp.WithString("originType", mcp.Required(), mcp.Description("Origin type indicating the source system. Valid value: FLOW (Flow pipeline).")),
				mcp.WithString("originId", mcp.Required(), mcp.Description("Origin identifier from the source system. For FLOW origin, use a pipeline run ID or pipeline ID.")),
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
