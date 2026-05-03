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
	}
}
