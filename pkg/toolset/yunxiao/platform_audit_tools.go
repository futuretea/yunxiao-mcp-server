package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func platformAuditTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_audit_logs",
				mcp.WithDescription("List audit logs in a Yunxiao organization. Use this to track user actions, resource changes, and security events within a time range."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("actionTimeStart", mcp.Required(), mcp.Description("Inclusive action-time lower bound. Format: RFC3339 timestamp (e.g., 2024-01-01T00:00:00+08:00).")),
				mcp.WithString("actionTimeEnd", mcp.Description("Action-time upper bound. Format: RFC3339 timestamp (e.g., 2024-01-31T23:59:59+08:00). Defaults to current time when omitted.")),
				mcp.WithString("userIds", mcp.Description("Filter by user IDs. Format: comma-separated numeric user IDs. Use list_users or list_organization_members to discover valid IDs.")),
				mcp.WithString("apps", mcp.Description("Filter by application identities. Format: comma-separated application names or identifiers.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithString("nextToken", mcp.Description("Pagination token from the previous response x-next-token header.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAuditLogs,
		},
	}
}
