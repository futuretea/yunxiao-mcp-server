package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func flowEnhancedTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_pipeline_overview",
				mcp.WithDescription("Get a comprehensive overview of a Flow pipeline including basic info, latest run, and recent run history in one read-only call."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID.")),
				mcp.WithBoolean("includeRuns", mcp.Description("Whether to include recent run history. Defaults to true.")),
				mcp.WithNumber("runLimit", mcp.Description("Max recent runs returned. Defaults to 5.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPipelineOverview,
		},
		{
			Tool: mcp.NewTool("get_pipeline_run_overview",
				mcp.WithDescription("Get a comprehensive overview of a Flow pipeline run including run details and pipeline jobs by category in one read-only call."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID.")),
				mcp.WithString("pipelineRunId", mcp.Required(), mcp.Description("Pipeline run ID.")),
				mcp.WithBoolean("includeJobs", mcp.Description("Whether to include pipeline jobs by category. Defaults to true.")),
				mcp.WithString("category", mcp.Description("Task category for job listing. Defaults to DEPLOY.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPipelineRunOverview,
		},
	}
}
