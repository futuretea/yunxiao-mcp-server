package yunxiao

import (
	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
	"github.com/mark3labs/mcp-go/mcp"
)

func flowWriteTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("pass_pipeline_validate",
				mcp.WithDescription("Pass (approve) a pipeline validation job."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID.")),
				mcp.WithString("pipelineRunId", mcp.Required(), mcp.Description("Pipeline run ID.")),
				mcp.WithString("jobId", mcp.Required(), mcp.Description("Validation job ID.")),
			),
			Handler: handlePassPipelineValidate,
		},
		{
			Tool: mcp.NewTool("refuse_pipeline_validate",
				mcp.WithDescription("Refuse (reject) a pipeline validation job."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID.")),
				mcp.WithString("pipelineRunId", mcp.Required(), mcp.Description("Pipeline run ID.")),
				mcp.WithString("jobId", mcp.Required(), mcp.Description("Validation job ID.")),
			),
			Handler: handleRefusePipelineValidate,
		},
	}
}
