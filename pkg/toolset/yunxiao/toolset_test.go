package yunxiao

import "testing"

func TestToolsetIncludesBaseReadTools(t *testing.T) {
	tools := (&Toolset{ReadOnly: true}).GetTools(nil)
	names := make(map[string]bool, len(tools))
	for _, tool := range tools {
		names[tool.Tool.Name] = true
	}

	for _, want := range []string{"get_current_user", "get_current_organization_info", "get_user_organizations", "list_organizations", "get_organization"} {
		if !names[want] {
			t.Fatalf("expected tool %q", want)
		}
	}
}
