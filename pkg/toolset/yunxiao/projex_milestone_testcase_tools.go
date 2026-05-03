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
				mcp.WithDescription("List milestones (planning checkpoints) in a Projex project. Milestones track progress against goals, distinct from sprints (time-boxed iterations) and versions (releases)."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithString("status", mcp.Description("Comma-separated milestone status IDs.")),
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
				mcp.WithDescription("List Projex testcase repositories in a Yunxiao organization. A testcase repository is a container for organizing test cases."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListTestcaseRepositories,
		},
		{
			Tool: mcp.NewTool("list_directories",
				mcp.WithDescription("List directories (folders) within a Projex testcase repository."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("testRepoId", mcp.Required(), mcp.Description("Testcase repository ID (string). Use list_testcase_repositories to discover available repositories.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListDirectories,
		},
		{
			Tool: mcp.NewTool("search_testcases",
				mcp.WithDescription("Search test cases within a single Projex testcase repository."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("testRepoId", mcp.Required(), mcp.Description("Testcase repository ID (string). Use list_testcase_repositories to discover available repositories.")),
				mcp.WithString("directoryId", mcp.Description("Directory ID (string) to filter by. Use list_directories to discover directory IDs.")),
				mcp.WithString("subject", mcp.Description("Filter by testcase subject/title (contains match).")),
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
				mcp.WithDescription("List test plans in a Yunxiao organization. A test plan groups test cases for execution tracking."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListTestPlans,
		},
		{
			Tool: mcp.NewTool("get_test_result_list",
				mcp.WithDescription("Get test execution result summaries for a specific directory within a test plan."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("testPlanIdentifier", mcp.Required(), mcp.Description("Test plan ID.")),
				mcp.WithString("directoryIdentifier", mcp.Required(), mcp.Description("Test plan directory ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetTestResultList,
		},
	}
}
