package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexEffortTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_current_user_effort_records",
				mcp.WithDescription("List actual effort (time tracking) records for the current user within a date range."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("startDate", mcp.Required(), mcp.Description("Start date in YYYY-MM-DD format (e.g., 2024-01-01).")),
				mcp.WithString("endDate", mcp.Required(), mcp.Description("End date in YYYY-MM-DD format (e.g., 2024-01-31).")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListCurrentUserEffortRecords,
		},
		{
			Tool: mcp.NewTool("list_effort_records",
				mcp.WithDescription("List actual effort (time tracking) records logged against a specific work item."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("workitemId", mcp.Required(), mcp.Description("Work item ID (numeric string). Find it via search_workitems.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListEffortRecords,
		},
		{
			Tool: mcp.NewTool("list_estimated_efforts",
				mcp.WithDescription("List estimated effort records for a specific work item."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("workitemId", mcp.Required(), mcp.Description("Work item ID (numeric string). Find it via search_workitems.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListEstimatedEfforts,
		},
	}
}
