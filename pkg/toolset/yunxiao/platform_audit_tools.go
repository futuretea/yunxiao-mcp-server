package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func platformAuditTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_audit_logs",
				mcp.WithDescription("List audit logs in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("actionTimeStart", mcp.Required(), mcp.Description("Inclusive RFC3339 action-time lower bound.")),
				mcp.WithString("actionTimeEnd", mcp.Description("RFC3339 action-time upper bound. Defaults to current time when omitted by Yunxiao.")),
				mcp.WithString("userIds", mcp.Description("Comma-separated user IDs.")),
				mcp.WithString("apps", mcp.Description("Comma-separated application identities.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
				mcp.WithString("nextToken", mcp.Description("Pagination token from the previous response x-next-token header.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAuditLogs,
		},
	}
}
