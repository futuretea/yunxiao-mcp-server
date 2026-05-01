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
		"list_pipelines",
		"get_pipeline",
		"list_pipeline_runs",
		"get_pipeline_run",
		"get_latest_pipeline_run",
		"list_pipeline_jobs_by_category",
		"list_pipeline_job_historys",
		"get_pipeline_job_run_log",
		"search_projects",
		"get_project",
		"search_workitems",
		"get_workitem",
		"list_package_repositories",
		"list_artifacts",
		"get_artifact",
		"get_sprint",
		"list_sprints",
		"list_all_work_item_types",
		"list_work_item_types",
		"get_work_item_type",
		"list_work_item_relation_work_item_types",
		"get_work_item_type_field_config",
		"get_work_item_workflow",
		"list_work_item_comments",
		"list_applications",
		"get_application",
		"get_env_variable_groups",
		"get_variable_group",
		"get_app_variable_groups",
		"get_app_variable_groups_revision",
		"list_app_release_workflows",
		"list_app_release_workflow_briefs",
		"get_app_release_workflow_stage",
		"list_app_release_stage_briefs",
		"list_app_release_stage_runs",
		"list_app_release_stage_exec_metadata",
		"get_app_release_stage_pipeline_run",
		"get_app_release_stage_pipeline_job_log",
		"list_change_order_versions",
		"get_change_order",
		"list_change_order_job_logs",
		"find_task_operation_log",
		"list_change_orders_by_origin",
	}
}
