package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexEffortTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_current_user_effort_records",
				mcp.WithDescription("List actual effort records for the current user."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("startDate", mcp.Required(), mcp.Description("Start date in yyyy-MM-dd format.")),
				mcp.WithString("endDate", mcp.Required(), mcp.Description("End date in yyyy-MM-dd format.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListCurrentUserEffortRecords,
		},
		{
			Tool: mcp.NewTool("list_effort_records",
				mcp.WithDescription("List actual effort records for a Projex work item."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Work item ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListEffortRecords,
		},
		{
			Tool: mcp.NewTool("list_estimated_efforts",
				mcp.WithDescription("List estimated effort records for a Projex work item."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Work item ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListEstimatedEfforts,
		},
	}
}
