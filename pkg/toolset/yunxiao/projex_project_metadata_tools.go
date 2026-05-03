package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func projexProjectMetadataTools() []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 6)
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
				mcp.WithDescription("List members in a Projex project."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithString("name", mcp.Description("Member name filter.")),
				mcp.WithString("roleId", mcp.Description("Project role ID filter, such as project.admin.")),
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
				mcp.WithDescription("List Projex project templates in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListProjectTemplates,
		},
		{
			Tool: mcp.NewTool("get_project_template_field_config",
				mcp.WithDescription("Get field configuration for a Projex project template."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("projectTemplateId", mcp.Required(), mcp.Description("Project template ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetProjectTemplateFieldConfig,
		},
	}
}

func projexProjectProgramTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_project_program",
				mcp.WithDescription("List Projex projects bound to a project program."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("programIdentifier", mcp.Required(), mcp.Description("Project program identifier.")),
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
				mcp.WithDescription("List roles in a Projex project."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListProjectRoles,
		},
		{
			Tool: mcp.NewTool("list_all_project_roles",
				mcp.WithDescription("List all Projex project roles in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. Defaults to the user's sole organization when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListAllProjectRoles,
		},
	}
}
