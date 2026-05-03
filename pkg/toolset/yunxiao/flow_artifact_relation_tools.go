package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func flowArtifactRelationTools() []toolset.ServerTool {
	return []toolset.ServerTool{
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
	}
}
