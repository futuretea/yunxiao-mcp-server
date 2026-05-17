package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func flowArtifactRelationTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_pipeline_relations",
				mcp.WithDescription("List objects related to a Flow pipeline, such as variable groups. Use this to discover pipeline dependencies and linked resources."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID (string). Use list_pipelines to find the pipeline ID.")),
				mcp.WithString("relObjectType", mcp.Required(), mcp.Description("Related object type. Example: VARIABLE_GROUP.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListPipelineRelations,
		},
		{
			Tool: mcp.NewTool("get_pipeline_scan_report_url",
				mcp.WithDescription("Get the scan report URL for a Flow pipeline. Use this to retrieve security scan or code quality report links."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("reportPath", mcp.Required(), mcp.Description("Report path provided by the pipeline scan step output.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPipelineScanReportURL,
		},
		{
			Tool: mcp.NewTool("get_pipeline_artifact_url",
				mcp.WithDescription("Get the artifact download URL for a Flow pipeline build output. Use this to retrieve build artifact download links."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("filePath", mcp.Required(), mcp.Description("Artifact file path. Typically returned by pipeline build step output.")),
				mcp.WithString("fileName", mcp.Required(), mcp.Description("Artifact file name. Typically returned by pipeline build step output.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPipelineArtifactURL,
		},
		{
			Tool: mcp.NewTool("get_pipeline_emas_artifact_url",
				mcp.WithDescription("Get the EMAS (Enterprise Mobile Application Studio) artifact download URL for a Flow pipeline. Use this to retrieve mobile app build artifact links."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("emasJobInstanceId", mcp.Required(), mcp.Description("EMAS job instance ID. Discovered from pipeline EMAS step output.")),
				mcp.WithString("md5", mcp.Required(), mcp.Description("MD5 checksum of the EMAS artifact. Discovered from pipeline EMAS step output.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID (integer or string). Use list_pipelines to find the pipeline ID.")),
				mcp.WithString("pipelineRunId", mcp.Required(), mcp.Description("Pipeline run ID (integer or string). Use list_pipeline_runs to discover valid run IDs.")),
				mcp.WithString("serviceConnectionId", mcp.Required(), mcp.Description("Service connection ID (integer or string). Typically discovered from pipeline configuration.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPipelineEmasArtifactURL,
		},
		{
			Tool: mcp.NewTool("get_last_instance",
				mcp.WithDescription("Get the last service connection instance for a Flow pipeline. Use this to retrieve the most recent service connection configuration."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID (string). Use list_pipelines to find the pipeline ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetLastInstance,
		},
	}
}
