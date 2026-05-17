package yunxiao

import "testing"

func TestToolsetIncludesBaseReadTools(t *testing.T) {
	tools := (&Toolset{ReadOnly: true}).GetTools(nil)
	if len(tools) < 100 {
		t.Fatalf("tool count = %d, expected at least 100 read-only tools", len(tools))
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
}

func TestGetCompactToolsHidesSupersededTools(t *testing.T) {
	all := (&Toolset{ReadOnly: true}).GetTools(nil)
	compact := (&Toolset{ReadOnly: true}).GetCompactTools(all)

	if len(compact) >= len(all) {
		t.Fatalf("compact tool count = %d, should be less than full count = %d", len(compact), len(all))
	}

	compactNames := make(map[string]bool, len(compact))
	for _, tool := range compact {
		if _, hidden := compactHiddenTools[tool.Tool.Name]; hidden {
			t.Fatalf("compact should hide %q", tool.Tool.Name)
		}
		compactNames[tool.Tool.Name] = true
	}

	for _, tool := range all {
		if _, hidden := compactHiddenTools[tool.Tool.Name]; hidden {
			if compactNames[tool.Tool.Name] {
				t.Fatalf("tool %q should be hidden in compact mode", tool.Tool.Name)
			}
		}
	}
}

func TestGetCompactToolsIncludesEnhancedAlternatives(t *testing.T) {
	all := (&Toolset{ReadOnly: true}).GetTools(nil)
	compact := (&Toolset{ReadOnly: true}).GetCompactTools(all)

	compactNames := make(map[string]bool, len(compact))
	for _, tool := range compact {
		compactNames[tool.Tool.Name] = true
	}

	want := []string{
		"get_application_overview",
		"get_change_order_overview",
		"get_environment_overview",
		"get_release_overview",
		"get_system_overview",
		"get_project_overview",
		"get_sprint_overview",
		"get_project_workitem_detail",
		"get_work_item_type_overview",
		"get_project_workitem_context",
		"get_organization_overview",
		"get_organization_department_overview",
		"get_organization_group_overview",
		"get_repository_overview",
		"get_change_request_overview",
		"get_commit_overview",
		"get_branch_overview",
		"get_pipeline_overview",
		"get_pipeline_run_overview",
	}
	for _, w := range want {
		if !compactNames[w] {
			t.Fatalf("compact tools should include %q", w)
		}
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
