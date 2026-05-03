package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func platformEnterpriseDepartmentTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_enterprise_departments",
				mcp.WithDescription("List enterprise departments visible to the current Yunxiao user."),
				mcp.WithString("parentId", mcp.Description("Parent department ID.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListEnterpriseDepartments,
		},
	}
}
