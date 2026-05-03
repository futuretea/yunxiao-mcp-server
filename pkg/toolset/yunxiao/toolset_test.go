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
	return names
}

func expectedPlatformToolNames() []string {
	return []string{
		"get_current_user",
		"get_current_organization_info",
		"get_user_organizations",
		"list_organizations",
		"get_organization",
		"list_organization_departments",
		"get_organization_department_info",
		"get_organization_department_ancestors",
		"list_enterprise_departments",
		"get_enterprise_department",
		"list_organization_members",
		"get_organization_member_info",
		"get_organization_member_info_by_user_id",
		"search_organization_members",
		"list_organization_groups",
		"get_organization_group",
		"list_organization_group_members",
		"list_organization_roles",
		"get_organization_role",
		"list_users",
		"list_audit_logs",
		"get_user",
		"list_app_extension_features",
		"get_organization_overview",
	}
}

func expectedCodeUpToolNames() []string {
	return []string{
		"list_repositories",
		"get_repository",
		"list_branches",
		"get_branch",
		"list_tags",
		"list_repository_members",
		"list_protected_branches",
		"get_protected_branch",
		"list_push_rules",
		"get_push_rule",
		"list_template_repositories",
		"list_namespaces",
		"get_namespace",
		"get_org_namespace",
		"list_group_members",
		"get_member_https_clone_username",
		"list_ssh_keys",
		"get_ssh_key",
		"list_user_ssh_keys",
		"list_webhooks",
		"get_webhook",
		"list_files",
		"get_file_blobs",
		"list_commits",
		"get_commit",
		"compare",
		"list_commit_statuses",
		"list_check_runs",
		"get_check_run",
		"list_merge_requests",
		"get_merge_request",
		"list_change_requests",
		"get_change_request",
		"list_change_request_patch_sets",
		"get_change_request_tree",
		"list_change_request_comments",
		"get_change_request_comment",
		"get_repository_overview",
		"get_change_request_overview",
	}
}

func expectedFlowToolNames() []string {
	return []string{
		"list_pipelines",
		"get_pipeline",
		"list_pipeline_runs",
		"get_pipeline_run",
		"get_latest_pipeline_run",
		"list_pipeline_jobs_by_category",
		"list_pipeline_job_historys",
		"get_pipeline_job_run_log",
		"get_pipeline_scan_report_url",
		"get_pipeline_artifact_url",
		"get_pipeline_emas_artifact_url",
		"list_pipeline_relations",
		"get_last_instance",
		"list_resource_members",
		"get_pipeline_overview",
	}
}

func expectedProjexToolNames() []string {
	return []string{
		"search_projects",
		"get_project",
		"get_project_overview",
		"get_project_workitem_summary",
		"get_project_workitem_context",
		"get_sprint_overview",
		"get_my_project_workitems",
		"get_project_workitem_board",
		"get_project_workitem_detail",
		"get_project_risk_dashboard",
		"get_project_member_task_status",
		"list_project_members",
		"list_project_templates",
		"get_project_template_field_config",
		"list_project_program",
		"list_project_roles",
		"list_all_project_roles",
		"search_workitems",
		"get_workitem",
		"get_sprint",
		"list_sprints",
		"list_all_work_item_types",
		"list_work_item_types",
		"get_work_item_type",
		"list_work_item_relation_work_item_types",
		"list_versions",
		"list_workitem_activities",
		"list_current_user_effort_records",
		"list_effort_records",
		"list_estimated_efforts",
		"list_workitem_attachments",
		"get_workitem_file",
		"list_workitem_relation_records",
		"list_labels",
		"list_milestones",
		"list_testcase_repositories",
		"list_directories",
		"get_testcase_field_config",
		"get_testcase",
		"search_testcases",
		"list_test_plans",
		"get_test_result_list",
		"get_work_item_type_field_config",
		"get_work_item_workflow",
		"list_work_item_comments",
	}
}

func expectedPackageToolNames() []string {
	return []string{
		"list_package_repositories",
		"list_artifacts",
		"get_artifact",
	}
}

func expectedAppStackToolNames() []string {
	return []string{
		"list_applications",
		"get_application",
		"search_app_templates",
		"list_environments",
		"get_environment",
		"list_application_members",
		"list_application_sources",
		"get_machine_deploy_log",
		"get_deploy_group",
		"list_resource_instances",
		"get_resource_instance",
		"get_pod_container_log",
		"get_pod_info",
		"get_kubernetes_object_info",
		"get_deployment_revision_info",
		"get_global_var",
		"list_global_vars",
		"get_env_variable_groups",
		"get_variable_group",
		"get_app_variable_groups",
		"get_app_variable_groups_revision",
		"get_latest_orchestration",
		"list_app_orchestration",
		"get_app_orchestration",
		"list_app_release_workflows",
		"list_app_release_workflow_briefs",
		"get_app_release_workflow_stage",
		"list_app_release_stage_briefs",
		"list_app_release_stage_runs",
		"list_app_release_stage_exec_metadata",
		"get_app_release_stage_pipeline_run",
		"get_app_release_stage_pipeline_job_log",
		"list_systems",
		"list_attached_apps",
		"list_system_members",
		"list_system_release_workflows",
		"get_release",
		"list_release_members",
		"list_release_products",
		"list_attached_change_requests",
		"list_release_executions",
		"search_releases",
		"get_appstack_change_request_audit_items",
		"list_appstack_change_request_executions",
		"list_appstack_change_request_work_items",
		"list_change_order_versions",
		"get_change_order",
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
		"get_department_usage",
		"list_developer_members",
		"get_developer_usage",
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
		"get_workitem",
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
		if tool.Tool.Name == "get_project" {
			t.Fatal("project-focused tools should hide get_project")
		}
		if tool.Tool.Name == "get_sprint" {
			t.Fatal("project-focused tools should hide get_sprint")
		}
		if tool.Tool.Name == "list_work_item_comments" {
			t.Fatal("project-focused tools should hide list_work_item_comments")
		}
	}
}

func TestGetProjectFocusedToolsIncludesEnhancedAlternatives(t *testing.T) {
	tools := (&Toolset{ReadOnly: true}).GetProjectFocusedTools(nil)

	want := []string{"get_project_overview", "get_sprint_overview"}
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
