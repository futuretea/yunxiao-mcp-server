package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func flowArtifactRelationTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_pipeline_scan_report_url",
				mcp.WithDescription("Get a temporary download URL for a Flow pipeline scan report."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("reportPath", mcp.Required(), mcp.Description("Scan report path returned by Flow pipeline APIs.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPipelineScanReportURL,
		},
		{
			Tool: mcp.NewTool("get_pipeline_artifact_url",
				mcp.WithDescription("Get a temporary download URL for a Flow pipeline artifact."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("filePath", mcp.Required(), mcp.Description("Artifact file path returned by Flow pipeline APIs.")),
				mcp.WithString("fileName", mcp.Required(), mcp.Description("Artifact file name.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPipelineArtifactURL,
		},
		{
			Tool: mcp.NewTool("get_pipeline_emas_artifact_url",
				mcp.WithDescription("Get a temporary download URL for a Flow EMAS artifact."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("emasJobInstanceId", mcp.Required(), mcp.Description("EMAS job instance ID.")),
				mcp.WithString("md5", mcp.Required(), mcp.Description("EMAS artifact MD5.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID.")),
				mcp.WithString("pipelineRunId", mcp.Required(), mcp.Description("Pipeline run ID.")),
				mcp.WithString("serviceConnectionId", mcp.Required(), mcp.Description("Service connection ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPipelineEmasArtifactURL,
		},
		{
			Tool: mcp.NewTool("list_pipeline_relations",
				mcp.WithDescription("List Flow pipeline related objects."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID.")),
				mcp.WithString("relObjectType", mcp.Required(), mcp.Description("Related object type, such as VARIABLE_GROUP.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListPipelineRelations,
		},
		{
			Tool: mcp.NewTool("get_last_instance",
				mcp.WithDescription("Get the latest Flow pipeline run instance detail."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("pipelineId", mcp.Required(), mcp.Description("Pipeline ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetLastInstance,
		},
	}
}
