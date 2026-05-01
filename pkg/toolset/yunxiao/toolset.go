package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

// Toolset exposes Yunxiao OpenAPI tools.
type Toolset struct {
	ReadOnly bool
}

func (t *Toolset) GetName() string {
	return "yunxiao"
}

func (t *Toolset) GetDescription() string {
	return "Yunxiao organization and DevOps OpenAPI tools"
}

func (t *Toolset) GetTools(_ any) []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_current_user",
				mcp.WithDescription("Get the current Yunxiao user for the configured access token."),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetCurrentUser,
		},
		{
			Tool: mcp.NewTool("get_current_organization_info",
				mcp.WithDescription("Get current user context, including the last organization returned by Yunxiao."),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetCurrentOrganizationInfo,
		},
		{
			Tool: mcp.NewTool("get_user_organizations",
				mcp.WithDescription("Get Yunxiao organizations visible to the current user."),
				mcp.WithNumber("page",
					mcp.Description("Page number. Defaults to 1 when omitted by Yunxiao."),
				),
				mcp.WithNumber("perPage",
					mcp.Description("Page size from 1 to 100. Defaults to 100 when omitted by Yunxiao."),
				),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetUserOrganizations,
		},
		{
			Tool: mcp.NewTool("list_organizations",
				mcp.WithDescription("List Yunxiao organizations visible to the current user."),
				mcp.WithNumber("page",
					mcp.Description("Page number. Defaults to 1 when omitted by Yunxiao."),
				),
				mcp.WithNumber("perPage",
					mcp.Description("Page size from 1 to 100. Defaults to 100 when omitted by Yunxiao."),
				),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListOrganizations,
		},
		{
			Tool: mcp.NewTool("get_organization",
				mcp.WithDescription("Get a Yunxiao organization by ID."),
				mcp.WithString("id",
					mcp.Required(),
					mcp.Description("Organization ID."),
				),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetOrganization,
		},
		{
			Tool: mcp.NewTool("list_repositories",
				mcp.WithDescription("List CodeUp repositories in a Yunxiao organization."),
				mcp.WithString("organizationId",
					mcp.Required(),
					mcp.Description("Yunxiao organization ID."),
				),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size from 1 to 100.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: created_at, name, path, or last_activity_at.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithString("search", mcp.Description("Fuzzy repository path search keyword.")),
				mcp.WithBoolean("archived", mcp.Description("Filter archived repositories.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListRepositories,
		},
		{
			Tool: mcp.NewTool("get_repository",
				mcp.WithDescription("Get a CodeUp repository by numeric ID or full path."),
				mcp.WithString("organizationId",
					mcp.Required(),
					mcp.Description("Yunxiao organization ID."),
				),
				mcp.WithString("repositoryId",
					mcp.Required(),
					mcp.Description("Repository numeric ID or full path such as org/repo."),
				),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetRepository,
		},
		{
			Tool: mcp.NewTool("list_branches",
				mcp.WithDescription("List branches in a CodeUp repository."),
				mcp.WithString("organizationId",
					mcp.Required(),
					mcp.Description("Yunxiao organization ID."),
				),
				mcp.WithString("repositoryId",
					mcp.Required(),
					mcp.Description("Repository numeric ID or full path such as org/repo."),
				),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithString("sort", mcp.Description("Sort mode: name_asc, name_desc, updated_asc, or updated_desc.")),
				mcp.WithString("search", mcp.Description("Branch search keyword.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListBranches,
		},
	}
}
