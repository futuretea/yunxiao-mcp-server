package yunxiao

import "testing"

func TestToolsetIncludesBaseReadTools(t *testing.T) {
	tools := (&Toolset{ReadOnly: true}).GetTools(nil)
	names := make(map[string]bool, len(tools))
	for _, tool := range tools {
		names[tool.Tool.Name] = true
		if tool.Tool.Annotations.ReadOnlyHint == nil || !*tool.Tool.Annotations.ReadOnlyHint {
			t.Fatalf("tool %q should be marked read-only", tool.Tool.Name)
		}
	}

	for _, want := range []string{
		"get_current_user",
		"get_current_organization_info",
		"get_user_organizations",
		"list_organizations",
		"get_organization",
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
		"search_projects",
		"get_project",
		"search_workitems",
		"get_workitem",
	} {
		if !names[want] {
			t.Fatalf("expected tool %q", want)
		}
	}
}
