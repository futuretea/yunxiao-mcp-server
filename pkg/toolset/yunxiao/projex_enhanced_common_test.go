package yunxiao

import (
	"testing"
)

func TestMergeConditions(t *testing.T) {
	tests := []struct {
		name     string
		existing string
		extra    string
		want     string
	}{
		{"both empty", "", "", ""},
		{"existing empty", "", `[{"field":"a"}]`, `[{"field":"a"}]`},
		{"extra empty", `[{"field":"a"}]`, "", `[{"field":"a"}]`},
		{"both arrays", `[{"field":"a"}]`, `[{"field":"b"}]`, `[{"field":"a"},{"field":"b"}]`},
		{"both objects", `{"conditionGroups":[[{"field":"a"}]]}`, `{"conditionGroups":[[{"field":"b"}]]}`, `{"conditionGroups":[[{"field":"a"},{"field":"b"}]]}`},
		{"existing object extra array", `{"conditionGroups":[[{"field":"a"}]]}`, `[{"field":"b"}]`, `{"conditionGroups":[[{"field":"a"}]]}`},
		{"existing array extra object", `[{"field":"a"}]`, `{"conditionGroups":[[{"field":"b"}]]}`, `[{"field":"a"}]`},
		{"existing invalid json", `not-json`, `[{"field":"b"}]`, `not-json`},
		{"empty existing groups extra groups", `{"conditionGroups":[]}`, `{"conditionGroups":[[{"field":"b"}]]}`, `{"conditionGroups":[[{"field":"b"}]]}`},
		{"both empty groups", `{"conditionGroups":[]}`, `{"conditionGroups":[]}`, `{"conditionGroups":[]}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeConditions(tt.existing, tt.extra)
			if got != tt.want {
				t.Fatalf("mergeConditions(%q, %q) = %q, want %q", tt.existing, tt.extra, got, tt.want)
			}
		})
	}
}

func TestOptionalIntDefault(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want int
	}{
		{"float64", float64(42), 42},
		{"int", int(7), 7},
		{"int64", int64(99), 99},
		{"string valid", "  123  ", 123},
		{"string invalid", "abc", 10},
		{"nil", nil, 10},
		{"bool", true, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := optionalIntDefault(map[string]any{"k": tt.val}, "k", 10)
			if got != tt.want {
				t.Fatalf("optionalIntDefault(%v) = %d, want %d", tt.val, got, tt.want)
			}
		})
	}
}

func TestNormalizedSampleLimit(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want int
	}{
		{"default", nil, 5},
		{"within range", float64(50), 50},
		{"negative", float64(-3), 0},
		{"overflow", float64(500), 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizedSampleLimit(map[string]any{"sampleLimit": tt.val})
			if got != tt.want {
				t.Fatalf("normalizedSampleLimit(%v) = %d, want %d", tt.val, got, tt.want)
			}
		})
	}
}

func TestExtractWorkitemData(t *testing.T) {
	tests := []struct {
		name      string
		payload   any
		wantLen   int
		wantTotal int
		wantErr   bool
	}{
		{"array", []any{map[string]any{"id": "1"}}, 1, 1, false},
		{"map with pagination", map[string]any{"data": []any{map[string]any{"id": "1"}}, "pagination": map[string]any{"total": float64(10)}}, 1, 10, false},
		{"map without pagination", map[string]any{"data": []any{map[string]any{"id": "1"}}}, 1, 0, false},
		{"invalid type", "string", 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, total, err := extractWorkitemData(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Fatalf("extractWorkitemData() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if len(data) != tt.wantLen {
					t.Fatalf("len(data) = %d, want %d", len(data), tt.wantLen)
				}
				if total != tt.wantTotal {
					t.Fatalf("total = %d, want %d", total, tt.wantTotal)
				}
			}
		})
	}
}
