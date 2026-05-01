package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexWorkitemMetadataTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_workitem_attachments",
				mcp.WithDescription("List attachments for a Projex work item."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Work item ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListWorkitemAttachments,
		},
		{
			Tool: mcp.NewTool("get_workitem_file",
				mcp.WithDescription("Get file metadata for a Projex work item file."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("workitemId", mcp.Required(), mcp.Description("Work item ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("File ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetWorkitemFile,
		},
		{
			Tool: mcp.NewTool("list_workitem_relation_records",
				mcp.WithDescription("List relation records for a Projex work item."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Work item ID.")),
				mcp.WithString("relationType", mcp.Required(), mcp.Description("Relation type: PARENT, SUB, ASSOCIATED, DEPEND_ON, or DEPENDED_BY.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListWorkitemRelationRecords,
		},
		{
			Tool: mcp.NewTool("list_labels",
				mcp.WithDescription("List labels in a Projex project."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListLabels,
		},
	}
}
