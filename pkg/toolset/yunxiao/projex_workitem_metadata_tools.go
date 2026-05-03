package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexWorkitemMetadataTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_workitem_attachments",
				mcp.WithDescription("List file attachments for a Projex work item. Use search_workitems to discover work item IDs."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("workitemId", mcp.Required(), mcp.Description("Work item ID. Use search_workitems to discover valid IDs.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListWorkitemAttachments,
		},
		{
			Tool: mcp.NewTool("list_workitem_relation_records",
				mcp.WithDescription("List relation records (parent, subtask, dependency links) for a Projex work item. Use search_workitems to discover work item IDs."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("workitemId", mcp.Required(), mcp.Description("Work item ID. Use search_workitems to discover valid IDs.")),
				mcp.WithString("relationType", mcp.Required(), mcp.Description("Relation type: PARENT, SUB, ASSOCIATED, DEPEND_ON, or DEPENDED_BY.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListWorkitemRelationRecords,
		},
		{
			Tool: mcp.NewTool("list_labels",
				mcp.WithDescription("List labels in a Projex project. Labels are used to categorize and filter work items. Use search_projects to discover project IDs."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID. Use search_projects to discover valid IDs.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListLabels,
		},
	}
}
