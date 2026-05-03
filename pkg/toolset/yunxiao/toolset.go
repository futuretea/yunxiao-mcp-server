package yunxiao

import "github.com/futuretea/yunxiao-mcp-server/pkg/toolset"

// minimalToolNames are the core query tools kept in --minimal mode.
var minimalToolNames = map[string]struct{}{
	// Platform
	"get_current_user":              {},
	"get_current_organization_info": {},

	// Projex — Project overview
	"search_projects":            {},
	"get_project_overview":       {},
	"get_project_risk_dashboard": {},

	// Projex — Workitem queries
	"search_workitems":             {},
	"get_project_workitem_summary": {},
	"get_project_workitem_detail":  {},
	"get_my_project_workitems":     {},
	"get_project_workitem_board":   {},

	// Projex — Sprint
	"list_sprints":        {},
	"get_sprint_overview": {},

	// Projex — Members
	"list_project_members": {},
}

// projectFocusedDomains are the tool domains enabled in project-focused mode.
var projectFocusedDomains = map[string]struct{}{
	"platform": {},
	"projex":   {},
}

// projectFocusedHiddenTools are raw tools superseded by enhanced alternatives.
var projectFocusedHiddenTools = map[string]struct{}{
	// Projex
	"get_project":                     {},
	"get_sprint":                      {},
	"get_workitem":                    {},
	"get_work_item_type":              {},
	"list_work_item_comments":         {},
	"list_workitem_activities":        {},
	"list_workitem_attachments":       {},
	"list_workitem_relation_records":  {},
	"list_versions":                   {},
	"list_milestones":                 {},
	"list_labels":                     {},
	"get_work_item_type_field_config": {},
	"get_work_item_workflow":          {},

	// Platform
	"get_organization":                      {},
	"get_organization_department_info":      {},
	"get_organization_department_ancestors": {},
	"get_organization_group":                {},
}

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
	tools := make([]toolset.ServerTool, 0, 200)
	tools = append(tools, withDomain(platformTools(), "platform")...)
	tools = append(tools, withDomain(codeupTools(), "codeup")...)
	tools = append(tools, withDomain(flowTools(), "flow")...)
	tools = append(tools, withDomain(projexTools(), "projex")...)
	tools = append(tools, withDomain(packageTools(), "packages")...)
	tools = append(tools, withDomain(appstackTools(), "appstack")...)
	tools = append(tools, withDomain(lingmaTools(), "lingma")...)
	tools = append(tools, withDomain(apiCallTools(), "api")...)
	return tools
}

// GetMinimalTools returns only the most essential query tools.
func (t *Toolset) GetMinimalTools(client any) []toolset.ServerTool {
	all := t.GetTools(client)
	filtered := make([]toolset.ServerTool, 0, len(minimalToolNames))
	for _, tool := range all {
		if _, ok := minimalToolNames[tool.Tool.Name]; ok {
			filtered = append(filtered, tool)
		}
	}
	return filtered
}

// GetProjectFocusedTools returns only platform + projex tools, hiding
// low-value raw tools that have enhanced alternatives.
func (t *Toolset) GetProjectFocusedTools(client any) []toolset.ServerTool {
	all := t.GetTools(client)
	filtered := make([]toolset.ServerTool, 0, len(all))
	for _, tool := range all {
		if _, ok := projectFocusedDomains[tool.Domain]; !ok {
			continue
		}
		if _, ok := projectFocusedHiddenTools[tool.Tool.Name]; ok {
			continue
		}
		filtered = append(filtered, tool)
	}
	return filtered
}

func withDomain(tools []toolset.ServerTool, domain string) []toolset.ServerTool {
	for i := range tools {
		tools[i].Domain = domain
	}
	return tools
}
