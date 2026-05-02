package yunxiao

import (
	"testing"
)

func TestBuildProjectConditions(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]any
		want   string
	}{
		{"empty", map[string]any{}, ""},
		{"name only", map[string]any{"name": "demo"}, `{"conditionGroups":[[{"className":"string","fieldIdentifier":"name","format":"input","operator":"CONTAINS","toValue":null,"value":["demo"]}]]}`},
		{"status only", map[string]any{"status": "TODO,DOING"}, `{"conditionGroups":[[{"className":"status","fieldIdentifier":"status","format":"list","operator":"CONTAINS","toValue":null,"value":["TODO","DOING"]}]]}`},
		{"creator only", map[string]any{"creator": "alice"}, `{"conditionGroups":[[{"className":"user","fieldIdentifier":"creator","format":"list","operator":"CONTAINS","toValue":null,"value":["alice"]}]]}`},
		{"multiple", map[string]any{"name": "demo", "status": "TODO"}, `{"conditionGroups":[[{"className":"string","fieldIdentifier":"name","format":"input","operator":"CONTAINS","toValue":null,"value":["demo"]},{"className":"status","fieldIdentifier":"status","format":"list","operator":"CONTAINS","toValue":null,"value":["TODO"]}]]}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildProjectConditions(tt.params)
			if got != tt.want {
				t.Fatalf("buildProjectConditions() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStartOfDay(t *testing.T) {
	if got := startOfDay("2026-05-03"); got != "2026-05-03 00:00:00" {
		t.Fatalf("startOfDay(date) = %q", got)
	}
	if got := startOfDay("2026-05-03 10:30:00"); got != "2026-05-03 10:30:00" {
		t.Fatalf("startOfDay(datetime) = %q", got)
	}
}

func TestEndOfDay(t *testing.T) {
	if got := endOfDay("2026-05-03"); got != "2026-05-03 23:59:59" {
		t.Fatalf("endOfDay(date) = %q", got)
	}
	if got := endOfDay("2026-05-03 10:30:00"); got != "2026-05-03 10:30:00" {
		t.Fatalf("endOfDay(datetime) = %q", got)
	}
}
