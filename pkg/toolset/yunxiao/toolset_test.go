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
		"list_organization_members",
		"get_organization_member_info",
		"get_organization_member_info_by_user_id",
		"search_organization_members",
		"list_organization_roles",
		"get_organization_role",
		"list_users",
	}
}

func expectedCodeUpToolNames() []string {
	return []string{
		"list_repositories",
		"get_repository",
		"list_branches",
		"get_branch",
		"list_files",
		"get_file_blobs",
		"list_commits",
		"get_commit",
		"compare",
		"list_change_requests",
		"get_change_request",
		"list_change_request_patch_sets",
		"get_change_request_tree",
		"list_change_request_comments",
		"get_change_request_comment",
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
	}
}

func expectedProjexToolNames() []string {
	return []string{
		"search_projects",
		"get_project",
		"search_workitems",
		"get_workitem",
		"get_sprint",
		"list_sprints",
		"list_all_work_item_types",
		"list_work_item_types",
		"get_work_item_type",
		"list_work_item_relation_work_item_types",
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
		"get_appstack_change_request_audit_items",
		"list_appstack_change_request_executions",
		"list_appstack_change_request_work_items",
		"list_change_order_versions",
		"get_change_order",
		"list_change_order_job_logs",
		"find_task_operation_log",
		"list_change_orders_by_origin",
	}
}
