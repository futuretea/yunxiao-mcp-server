package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func flowTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 8)
	tools = append(tools, flowPipelineTools()...)
	tools = append(tools, flowPipelineRunTools()...)
	tools = append(tools, flowPipelineJobTools()...)
	tools = append(tools, flowArtifactRelationTools()...)
	tools = append(tools, flowResourceMemberTools()...)
	tools = append(tools, flowEnhancedTools()...)
	return tools
}

func flowPipelineTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_pipelines",
				mcp.WithDescription("List Flow CI/CD pipelines in a Yunxiao organization. Use this to discover pipelines and obtain their IDs before calling pipeline-scoped tools. For a comprehensive view of a single pipeline including latest run and history, use get_pipeline_overview instead."),
				mcp.WithString("organizationId",
					mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization."),
				),
				mcp.WithNumber("createStartTime", mcp.Description("Pipeline creation start time as a Unix timestamp in milliseconds (e.g., 1704067200000).")),
				mcp.WithNumber("createEndTime", mcp.Description("Pipeline creation end time as a Unix timestamp in milliseconds.")),
				mcp.WithNumber("executeStartTime", mcp.Description("Pipeline execution start time as a Unix timestamp in milliseconds.")),
				mcp.WithNumber("executeEndTime", mcp.Description("Pipeline execution end time as a Unix timestamp in milliseconds.")),
				mcp.WithString("pipelineName", mcp.Description("Filter by pipeline name (contains match).")),
				mcp.WithString("statusList", mcp.Description("Comma-separated pipeline statuses. Common values: RUNNING, SUCCESS, FAIL, CANCELED, WAITING.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Yunxiao supports up to 30.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListPipelines,
		},
		{
			Tool: mcp.NewTool("get_pipeline",
				mcp.WithDescription("Get a single Flow pipeline by ID. Use list_pipelines to discover valid pipeline IDs. For a comprehensive view with latest run info, use get_pipeline_overview instead."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID (string). Use list_pipelines to discover valid IDs.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPipeline,
		},
	}
}

func flowPipelineRunTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_pipeline_runs",
				mcp.WithDescription("List execution runs for a Flow pipeline. Use this to review historical runs and their statuses. For the latest run only, use get_latest_pipeline_run."),
				mcp.WithString("organizationId",
					mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization."),
				),
				mcp.WithString("pipelineId",
					mcp.Required(),
					mcp.Description("Pipeline ID (string). Use list_pipelines to find the pipeline ID."),
				),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Yunxiao supports up to 30.")),
				mcp.WithNumber("startTime", mcp.Description("Run start time as a Unix timestamp in milliseconds.")),
				mcp.WithNumber("endTime", mcp.Description("Run end time as a Unix timestamp in milliseconds.")),
				mcp.WithString("status", mcp.Description("Filter by run status. Common values: FAIL, SUCCESS, RUNNING.")),
				mcp.WithNumber("triggerMode", mcp.Description("Filter by trigger mode: 1 manual, 2 scheduled, 3 code push, 5 pipeline, 6 webhook.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListPipelineRuns,
		},
		{
			Tool: mcp.NewTool("get_latest_pipeline_run",
				mcp.WithDescription("Get the latest execution run for a Flow pipeline. Use this for a quick status check without listing all historical runs."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID (string). Use list_pipelines to find the pipeline ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetLatestPipelineRun,
		},
		{
			Tool: mcp.NewTool("get_pipeline_run",
				mcp.WithDescription("Get a specific Flow pipeline run by ID. Use list_pipeline_runs to discover valid run IDs. For a comprehensive view with metadata, use get_pipeline_run_overview instead."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID (string). Use list_pipelines to find the pipeline ID.")),
				mcp.WithString("pipelineRunId", mcp.Required(), mcp.Description("Pipeline run ID. Use list_pipeline_runs to discover valid run IDs.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPipelineRun,
		},
	}
}

func flowPipelineJobTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_pipeline_jobs_by_category",
				mcp.WithDescription("List jobs (tasks) within a Flow pipeline grouped by category. Use this after identifying a pipeline to see its build, deploy, and test stages."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID. Use list_pipelines to discover valid IDs.")),
				mcp.WithString("category", mcp.Required(), mcp.Description("Task category. Common value: DEPLOY (for deployment tasks).")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListPipelineJobsByCategory,
		},
		{
			Tool: mcp.NewTool("list_pipeline_job_historys",
				mcp.WithDescription("List execution history for a specific Flow pipeline job. Use this to track how a particular job (e.g., a deploy step) has performed across multiple runs."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID. Use list_pipelines to discover valid IDs.")),
				mcp.WithString("category", mcp.Required(), mcp.Description("Task category. Common value: DEPLOY (for deployment tasks).")),
				mcp.WithString("identifier", mcp.Required(), mcp.Description("Pipeline job identifier (string). Use list_pipeline_jobs_by_category to discover job identifiers.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Yunxiao supports up to 30.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListPipelineJobHistorys,
		},
		{
			Tool: mcp.NewTool("get_pipeline_job_run_log",
				mcp.WithDescription("Get the execution log for a specific job within a Flow pipeline run. Use this to debug pipeline failures by inspecting individual job output."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID (string). Use list_pipelines to find the pipeline ID.")),
				mcp.WithString("pipelineRunId", mcp.Required(), mcp.Description("Pipeline run ID. Use list_pipeline_runs to discover valid run IDs.")),
				mcp.WithString("jobId", mcp.Required(), mcp.Description("Job ID within the pipeline run. Use list_pipeline_jobs_by_category to discover valid job IDs.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPipelineJobRunLog,
		},
	}
}
