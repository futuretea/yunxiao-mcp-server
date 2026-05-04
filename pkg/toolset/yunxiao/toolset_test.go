package yunxiao

import "testing"

func TestToolsetIncludesBaseReadTools(t *testing.T) {
	tools := (&Toolset{ReadOnly: true}).GetTools(nil)
	wantTools := expectedToolNames()
	if len(tools) != len(wantTools) {
		t.Fatalf("tool count = %d, want %d", len(tools), len(wantTools))
	}

	names := make(map[string]bool, len(tools))
	for _, tool := range tools {
		if names[tool.Tool.Name] {
			t.Fatalf("duplicate tool %q", tool.Tool.Name)
		}
		names[tool.Tool.Name] = true
		if tool.Tool.Annotations.ReadOnlyHint == nil || !*tool.Tool.Annotations.ReadOnlyHint {
			t.Fatalf("tool %q should be marked read-only", tool.Tool.Name)
		}
	}

	for _, want := range wantTools {
		if !names[want] {
			t.Fatalf("expected tool %q", want)
		}
	}
}

func expectedToolNames() []string {
	var names []string
	names = append(names, expectedPlatformToolNames()...)
	names = append(names, expectedCodeUpToolNames()...)
	names = append(names, expectedFlowToolNames()...)
	names = append(names, expectedProjexToolNames()...)
	names = append(names, expectedPackageToolNames()...)
	names = append(names, expectedAppStackToolNames()...)
	names = append(names, expectedLingmaToolNames()...)
	names = append(names, "call_yunxiao_api")
	names = append(names, "describe_toolset")
	return names
}

func expectedPlatformToolNames() []string {
	return []string{
		"get_current_user",
		"get_current_organization_info",
		"get_user_organizations",
		"list_organizations",
		"list_organization_departments",
		"list_enterprise_departments",
		"list_organization_members",
		"search_organization_members",
		"list_organization_groups",
		"list_organization_group_members",
		"list_organization_roles",
		"list_users",
		"list_audit_logs",
		"get_user",
		"list_app_extension_features",
		"get_organization_overview",
		"get_organization_department_overview",
		"get_organization_group_overview",
	}
}

func expectedCodeUpToolNames() []string {
	return []string{
		"list_repositories",
		"list_branches",
		"list_tags",
		"list_repository_members",
		"list_protected_branches",
		"list_push_rules",
		"list_template_repositories",
		"list_namespaces",
		"list_group_members",
		"list_ssh_keys",
		"list_user_ssh_keys",
		"list_webhooks",
		"list_files",
		"list_commits",
		"list_commit_statuses",
		"list_check_runs",
		"list_merge_requests",
		"list_change_requests",
		"list_change_request_patch_sets",
		"list_change_request_comments",
		"get_repository_overview",
		"get_change_request_overview",
		"get_commit_overview",
		"get_branch_overview",
	}
}

func expectedFlowToolNames() []string {
	return []string{
		"list_pipelines",
		"list_pipeline_runs",
		"list_pipeline_jobs_by_category",
		"list_pipeline_job_historys",
		"list_pipeline_relations",
		"list_resource_members",
		"get_pipeline_overview",
		"get_pipeline_run_overview",
	}
}

func expectedProjexToolNames() []string {
	return []string{
		"search_projects",
		"get_project_overview",
		"get_project_workitem_summary",
		"get_project_workitem_context",
		"get_sprint_overview",
		"get_my_project_workitems",
		"get_project_workitem_board",
		"get_project_workitem_detail",
		"get_work_item_type_overview",
		"get_project_risk_dashboard",
		"get_project_member_task_status",
		"get_sprint_velocity",
		"get_workitem_status_timeline",
		"get_blocker_analysis",
		"get_member_workload_trend",
		"get_team_workload_breakdown",
		"list_project_members",
		"list_project_templates",
		"list_project_program",
		"list_project_roles",
		"list_all_project_roles",
		"search_workitems",
		"list_sprints",
		"list_all_work_item_types",
		"list_work_item_types",
		"list_work_item_relation_work_item_types",
		"list_versions",
		"list_workitem_activities",
		"list_current_user_effort_records",
		"list_effort_records",
		"list_estimated_efforts",
		"list_workitem_attachments",
		"list_workitem_relation_records",
		"list_labels",
		"list_milestones",
		"list_testcase_repositories",
		"list_directories",
		"search_testcases",
		"list_test_plans",
		"get_test_result_list",
		"list_work_item_comments",
	}
}

func expectedPackageToolNames() []string {
	return []string{
		"list_package_repositories",
		"list_artifacts",
	}
}

func expectedAppStackToolNames() []string {
	return []string{
		"list_applications",
		"search_app_templates",
		"list_environments",
		"list_application_members",
		"list_application_sources",
		"list_resource_instances",
		"list_global_vars",
		"list_app_orchestration",
		"list_app_release_workflows",
		"list_app_release_workflow_briefs",
		"list_app_release_stage_briefs",
		"list_app_release_stage_runs",
		"list_app_release_stage_exec_metadata",
		"list_systems",
		"list_attached_apps",
		"list_system_members",
		"list_system_release_workflows",
		"list_release_members",
		"list_release_products",
		"list_attached_change_requests",
		"list_release_executions",
		"search_releases",
		"list_appstack_change_request_executions",
		"list_appstack_change_request_work_items",
		"list_change_order_versions",
		"list_change_order_job_logs",
		"find_task_operation_log",
		"list_change_orders_by_origin",
		"get_application_overview",
		"get_environment_overview",
		"get_release_overview",
	}
}

func expectedLingmaToolNames() []string {
	return []string{
		"list_developer_members",
		"list_knowledge_bases",
		"list_kb_files",
		"list_kb_members",
	}
}

func TestToolsetGetNameAndDescription(t *testing.T) {
	ts := &Toolset{ReadOnly: true}
	if got := ts.GetName(); got != "yunxiao" {
		t.Fatalf("GetName() = %q, want yunxiao", got)
	}
	if got := ts.GetDescription(); got != "Yunxiao organization and DevOps OpenAPI tools" {
		t.Fatalf("GetDescription() = %q", got)
	}
}

func TestGetMinimalToolsReturnsExactSet(t *testing.T) {
	want := []string{
		"get_current_user",
		"get_current_organization_info",
		"search_projects",
		"get_project_overview",
		"get_project_risk_dashboard",
		"search_workitems",
		"get_project_workitem_summary",
		"get_project_workitem_detail",
		"get_my_project_workitems",
		"get_project_workitem_board",
		"list_sprints",
		"get_sprint_overview",
		"list_project_members",
	}

	tools := (&Toolset{ReadOnly: true}).GetMinimalTools(nil)
	if len(tools) != len(want) {
		t.Fatalf("tool count = %d, want %d", len(tools), len(want))
	}

	names := make(map[string]bool, len(tools))
	for _, tool := range tools {
		names[tool.Tool.Name] = true
		if tool.Domain != "platform" && tool.Domain != "projex" {
			t.Fatalf("tool %q has unexpected domain %q", tool.Tool.Name, tool.Domain)
		}
	}
	for _, w := range want {
		if !names[w] {
			t.Fatalf("expected minimal tool %q", w)
		}
	}
}

func TestGetProjectFocusedToolsIncludesPlatformAndProjex(t *testing.T) {
	tools := (&Toolset{ReadOnly: true}).GetProjectFocusedTools(nil)

	hasPlatform := false
	hasProjex := false
	for _, tool := range tools {
		switch tool.Domain {
		case "platform":
			hasPlatform = true
		case "projex":
			hasProjex = true
		default:
			t.Fatalf("unexpected domain %q for tool %q", tool.Domain, tool.Tool.Name)
		}
	}
	if !hasPlatform {
		t.Fatal("project-focused tools should include platform tools")
	}
	if !hasProjex {
		t.Fatal("project-focused tools should include projex tools")
	}
}

func TestGetProjectFocusedToolsHidesSupersededTools(t *testing.T) {
	tools := (&Toolset{ReadOnly: true}).GetProjectFocusedTools(nil)
	for _, tool := range tools {
		if _, ok := projectFocusedHiddenTools[tool.Tool.Name]; ok {
			t.Fatalf("project-focused tools should hide %q", tool.Tool.Name)
		}
	}
}

func TestGetProjectFocusedToolsIncludesEnhancedAlternatives(t *testing.T) {
	tools := (&Toolset{ReadOnly: true}).GetProjectFocusedTools(nil)

	want := []string{
		"get_project_overview",
		"get_sprint_overview",
		"get_project_workitem_detail",
		"get_work_item_type_overview",
		"get_project_workitem_context",
		"get_organization_overview",
		"get_organization_department_overview",
		"get_organization_group_overview",
	}
	for _, w := range want {
		found := false
		for _, tool := range tools {
			if tool.Tool.Name == w {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("project-focused tools should include %q", w)
		}
	}
}

func TestReadOnlyModeExcludesWriteTools(t *testing.T) {
	ts := &Toolset{ReadOnly: true}
	all := ts.GetTools(nil)
	for _, tool := range all {
		if _, ok := writeToolNames[tool.Tool.Name]; ok {
			t.Fatalf("read-only mode should exclude write tool %q", tool.Tool.Name)
		}
	}
}

func TestWriteModeIncludesWriteTools(t *testing.T) {
	ts := &Toolset{ReadOnly: false}
	all := ts.GetTools(nil)
	names := make(map[string]bool, len(all))
	for _, tool := range all {
		names[tool.Tool.Name] = true
	}
	for want := range writeToolNames {
		if !names[want] {
			t.Fatalf("write mode should include write tool %q", want)
		}
	}
}
