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
	return tools
}

func flowPipelineTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_pipelines",
				mcp.WithDescription("List Flow pipelines in a Yunxiao organization."),
				mcp.WithString("organizationId",
					mcp.Required(),
					mcp.Description("Yunxiao organization ID."),
				),
				mcp.WithNumber("createStartTime", mcp.Description("Pipeline creation start time in milliseconds.")),
				mcp.WithNumber("createEndTime", mcp.Description("Pipeline creation end time in milliseconds.")),
				mcp.WithNumber("executeStartTime", mcp.Description("Pipeline execution start time in milliseconds.")),
				mcp.WithNumber("executeEndTime", mcp.Description("Pipeline execution end time in milliseconds.")),
				mcp.WithString("pipelineName", mcp.Description("Pipeline name filter.")),
				mcp.WithString("statusList", mcp.Description("Comma-separated statuses such as RUNNING,SUCCESS,FAIL,CANCELED,WAITING.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size. Yunxiao supports up to 30.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListPipelines,
		},
		{
			Tool: mcp.NewTool("get_pipeline",
				mcp.WithDescription("Get Flow pipeline details."),
				mcp.WithString("organizationId",
					mcp.Required(),
					mcp.Description("Yunxiao organization ID."),
				),
				mcp.WithString("pipelineId",
					mcp.Required(),
					mcp.Description("Pipeline ID."),
				),
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
				mcp.WithDescription("List Flow pipeline runs."),
				mcp.WithString("organizationId",
					mcp.Required(),
					mcp.Description("Yunxiao organization ID."),
				),
				mcp.WithString("pipelineId",
					mcp.Required(),
					mcp.Description("Pipeline ID."),
				),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size. Yunxiao supports up to 30.")),
				mcp.WithNumber("startTime", mcp.Description("Run start time in milliseconds.")),
				mcp.WithNumber("endTime", mcp.Description("Run end time in milliseconds.")),
				mcp.WithString("status", mcp.Description("Run status: FAIL, SUCCESS, or RUNNING.")),
				mcp.WithNumber("triggerMode", mcp.Description("Trigger mode: 1 manual, 2 scheduled, 3 code, 5 pipeline, 6 webhook.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListPipelineRuns,
		},
		{
			Tool: mcp.NewTool("get_pipeline_run",
				mcp.WithDescription("Get a Flow pipeline run by ID."),
				mcp.WithString("organizationId",
					mcp.Required(),
					mcp.Description("Yunxiao organization ID."),
				),
				mcp.WithString("pipelineId",
					mcp.Required(),
					mcp.Description("Pipeline ID."),
				),
				mcp.WithString("pipelineRunId",
					mcp.Required(),
					mcp.Description("Pipeline run ID."),
				),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPipelineRun,
		},
		{
			Tool: mcp.NewTool("get_latest_pipeline_run",
				mcp.WithDescription("Get the latest Flow pipeline run."),
				mcp.WithString("organizationId",
					mcp.Required(),
					mcp.Description("Yunxiao organization ID."),
				),
				mcp.WithString("pipelineId",
					mcp.Required(),
					mcp.Description("Pipeline ID."),
				),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetLatestPipelineRun,
		},
	}
}

func flowPipelineJobTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_pipeline_jobs_by_category",
				mcp.WithDescription("List Flow pipeline jobs by task category."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID.")),
				mcp.WithString("category", mcp.Required(), mcp.Description("Task category, currently DEPLOY.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListPipelineJobsByCategory,
		},
		{
			Tool: mcp.NewTool("list_pipeline_job_historys",
				mcp.WithDescription("List Flow pipeline job execution history."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID.")),
				mcp.WithString("category", mcp.Required(), mcp.Description("Task category, currently DEPLOY.")),
				mcp.WithString("identifier", mcp.Required(), mcp.Description("Pipeline job identifier.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size. Yunxiao supports up to 30.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListPipelineJobHistorys,
		},
		{
			Tool: mcp.NewTool("get_pipeline_job_run_log",
				mcp.WithDescription("Get a Flow pipeline job run log."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID.")),
				mcp.WithString("pipelineRunId", mcp.Required(), mcp.Description("Pipeline run ID.")),
				mcp.WithString("jobId", mcp.Required(), mcp.Description("Pipeline job ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPipelineJobRunLog,
		},
	}
}
