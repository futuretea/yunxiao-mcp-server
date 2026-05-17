package yunxiao

import "github.com/futuretea/yunxiao-mcp-server/pkg/toolset"

// compactHiddenTools are tools hidden in compact mode (default on).
// Categories: (A) raw getters with enhanced overview alternatives,
// (B) list tools whose data is included in enhanced overviews,
// (C) briefs that have full-detail alternatives,
// (D) very specialized tools rarely needed by LLMs.
var compactHiddenTools = map[string]struct{}{
	// Platform (A) — superseded by get_organization_overview, get_organization_department_overview, get_organization_group_overview
	"get_organization":                      {},
	"get_organization_department_info":      {},
	"get_organization_department_ancestors": {},
	"get_organization_group":                {},
	// Platform (B) — included in enhanced overviews
	"list_organization_group_members": {},
	// Platform (D) — admin / rarely needed
	"list_enterprise_departments": {},
	"list_audit_logs":             {},

	// Projex (A) — superseded by get_project_overview, get_sprint_overview, get_project_workitem_detail, etc.
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
	// Projex (D) — specialized / rarely needed
	"list_project_templates":           {},
	"list_project_program":             {},
	"list_all_project_roles":           {},
	"list_current_user_effort_records": {},
	"list_estimated_efforts":           {},
	"list_effort_records":              {},
	"list_directories":                 {},
	"list_test_plans":                  {},
	"get_test_result_list":             {},

	// AppStack (A) — superseded by get_application_overview, get_environment_overview, etc.
	"get_application":  {},
	"get_environment":  {},
	"get_change_order": {},
	"get_release":      {},
	// AppStack (B) — included in enhanced overviews
	"list_environments":             {},
	"list_application_members":      {},
	"list_release_members":          {},
	"list_release_products":         {},
	"list_attached_change_requests": {},
	"list_attached_apps":            {},
	"list_system_members":           {},
	// AppStack (C) — briefs with full-detail alternatives
	"list_app_release_workflow_briefs": {},
	"list_app_release_stage_briefs":    {},
	// AppStack (D) — specialized metadata
	"get_app_variable_groups_revision": {},
	"list_application_sources":         {},

	// CodeUp (A) — superseded by get_repository_overview, get_branch_overview, get_commit_overview, get_change_request_overview
	"get_repository":     {},
	"get_branch":         {},
	"get_commit":         {},
	"get_change_request": {},
	// CodeUp (B) — included in enhanced overviews
	"list_commit_statuses":           {},
	"list_check_runs":                {},
	"list_change_request_patch_sets": {},
	"list_change_request_comments":   {},
	// CodeUp (D) — specialized
	"list_template_repositories": {},
	"list_ssh_keys":              {},
	"list_user_ssh_keys":         {},
	"list_webhooks":              {},
	"list_protected_branches":    {},
	"list_push_rules":            {},
	"list_tags":                  {},
	"list_namespaces":            {},
	"list_group_members":         {},

	// Flow (A) — superseded by get_pipeline_overview, get_pipeline_run_overview
	"get_pipeline":     {},
	"get_pipeline_run": {},
	// Flow (D) — specialized
	"get_last_instance":              {},
	"get_pipeline_emas_artifact_url": {},
	"get_pipeline_artifact_url":      {},
	"get_pipeline_scan_report_url":   {},
	"list_pipeline_job_historys":     {},
}

// writeToolNames are tools that perform mutations and are excluded in read-only mode.
var writeToolNames = map[string]struct{}{
	"create_workitem":            {},
	"update_workitem":            {},
	"update_workitem_status":     {},
	"add_workitem_comment":       {},
	"create_change_request":      {},
	"add_change_request_comment": {},
	"create_merge_request":       {},
	"close_change_request":       {},
	"reopen_change_request":      {},
	"merge_change_request":       {},
	"pass_pipeline_validate":     {},
	"refuse_pipeline_validate":   {},
	"create_change_order":        {},
	"execute_job_action":            {},
	"execute_system_release_stage":  {},
	"execute_app_release_stage":     {},
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
	tools = append(tools, withDomain(flowWriteTools(), "flow")...)
	tools = append(tools, withDomain(capabilityTools(), "meta")...)
	return t.filterReadOnly(tools)
}

// GetCompactTools hides raw tools that have enhanced overview alternatives.
func (t *Toolset) GetCompactTools(all []toolset.ServerTool) []toolset.ServerTool {
	filtered := make([]toolset.ServerTool, 0, len(all))
	for _, tool := range all {
		if _, ok := compactHiddenTools[tool.Tool.Name]; ok {
			continue
		}
		filtered = append(filtered, tool)
	}
	return filtered
}

func (t *Toolset) filterReadOnly(tools []toolset.ServerTool) []toolset.ServerTool {
	if !t.ReadOnly {
		return tools
	}
	filtered := make([]toolset.ServerTool, 0, len(tools))
	for _, tool := range tools {
		if _, ok := writeToolNames[tool.Tool.Name]; !ok {
			filtered = append(filtered, tool)
		}
	}
	return filtered
}

func withDomain(tools []toolset.ServerTool, domain string) []toolset.ServerTool {
	for i := range tools {
		tools[i].Domain = domain
	}
	return tools
}
