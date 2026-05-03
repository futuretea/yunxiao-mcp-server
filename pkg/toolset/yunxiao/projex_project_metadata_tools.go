package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexProjectMetadataTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 5)
	tools = append(tools, projexProjectMemberTools()...)
	tools = append(tools, projexProjectTemplateTools()...)
	tools = append(tools, projexProjectProgramTools()...)
	tools = append(tools, projexProjectRoleTools()...)
	return tools
}

func projexProjectMemberTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_project_members",
				mcp.WithDescription("List members in a Projex project. Use this to discover user IDs for filtering work items or assigning tasks."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID. Use search_projects to discover valid IDs.")),
				mcp.WithString("name", mcp.Description("Filter by member name (contains match).")),
				mcp.WithString("roleId", mcp.Description("Filter by project role ID, such as project.admin. Use list_project_roles to discover available roles.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListProjectMembers,
		},
	}
}

func projexProjectTemplateTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_project_templates",
				mcp.WithDescription("List Projex project templates in a Yunxiao organization. Useful when setting up new projects."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListProjectTemplates,
		},
	}
}

func projexProjectProgramTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_project_program",
				mcp.WithDescription("List Projex projects bound to a project program (project group)."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("programIdentifier", mcp.Required(), mcp.Description("Project program identifier (string).")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListProjectProgram,
		},
	}
}

func projexProjectRoleTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_project_roles",
				mcp.WithDescription("List roles defined in a specific Projex project."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID. Use search_projects to discover valid IDs.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListProjectRoles,
		},
		{
			Tool: mcp.NewTool("list_all_project_roles",
				mcp.WithDescription("List all Projex project roles across a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAllProjectRoles,
		},
	}
}
