package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexMilestoneTestcaseTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 6)
	tools = append(tools, projexMilestoneTools()...)
	tools = append(tools, projexTestcaseReadTools()...)
	tools = append(tools, projexTestPlanReadTools()...)
	return tools
}

func projexMilestoneTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_milestones",
				mcp.WithDescription("List milestones in a Projex project."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID.")),
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
			Tool: mcp.NewTool("list_testcase_repositories",
				mcp.WithDescription("List Projex testcase repositories in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListTestcaseRepositories,
		},
		{
			Tool: mcp.NewTool("list_directories",
				mcp.WithDescription("List testcase directories in a Projex testcase repository."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("testRepoId", mcp.Required(), mcp.Description("Testcase repository ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListDirectories,
		},
		{
			Tool: mcp.NewTool("search_testcases",
				mcp.WithDescription("Search Projex testcases in one testcase repository."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("testRepoId", mcp.Required(), mcp.Description("Testcase repository ID.")),
				mcp.WithString("directoryId", mcp.Description("Directory ID filter.")),
				mcp.WithString("subject", mcp.Description("Testcase subject contains filter.")),
				mcp.WithString("conditions", mcp.Description("Advanced conditions JSON string. Overrides simple filters.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: gmtCreate or name.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleSearchTestcases,
		},
	}
}

func projexTestPlanReadTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_test_plans",
				mcp.WithDescription("List Projex test plans in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListTestPlans,
		},
		{
			Tool: mcp.NewTool("get_test_result_list",
				mcp.WithDescription("Get testcase result summaries in a Projex test plan directory."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("testPlanIdentifier", mcp.Required(), mcp.Description("Test plan ID.")),
				mcp.WithString("directoryIdentifier", mcp.Required(), mcp.Description("Test plan directory ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetTestResultList,
		},
	}
}
