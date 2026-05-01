package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexMilestoneTestcaseTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 4)
	tools = append(tools, projexMilestoneTools()...)
	tools = append(tools, projexTestcaseReadTools()...)
	return tools
}

func projexMilestoneTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_milestones",
				mcp.WithDescription("List milestones in a Projex project."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithString("status", mcp.Description("Comma-separated milestone statuses.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListMilestones,
		},
	}
}

func projexTestcaseReadTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_directories",
				mcp.WithDescription("List testcase directories in a Projex testcase repository."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Testcase repository ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListDirectories,
		},
		{
			Tool: mcp.NewTool("get_testcase_field_config",
				mcp.WithDescription("Get testcase field configuration in a Projex testcase repository."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Testcase repository ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetTestcaseFieldConfig,
		},
		{
			Tool: mcp.NewTool("get_testcase",
				mcp.WithDescription("Get a Projex testcase by ID."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("testRepoId", mcp.Required(), mcp.Description("Testcase repository ID.")),
				mcp.WithString("id", mcp.Required(), mcp.Description("Testcase ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetTestcase,
		},
	}
}
