package yunxiao

import (
	"context"
	"fmt"
	"net/http"
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

func TestBuildCategoryResultReturnsErrorOnSearchFailure(t *testing.T) {
	_, err := buildCategoryResult(context.Background(), []string{"Bug"}, map[string]any{}, func(string) (any, error) {
		return nil, fmt.Errorf("search failed")
	})
	if err == nil {
		t.Fatal("buildCategoryResult() expected error")
	}
}

func TestGroupWorkitemsByStatus(t *testing.T) {
	data := []any{
		map[string]any{"status": map[string]any{"name": "TODO"}},
		map[string]any{"status": map[string]any{"name": "TODO"}},
		map[string]any{"status": map[string]any{"name": "DOING"}},
		map[string]any{"id": "no-status"},
	}

	columns, counts := groupWorkitemsByStatus(data)

	if len(columns["TODO"].([]any)) != 2 {
		t.Fatalf("TODO column = %v", columns["TODO"])
	}
	if len(columns["DOING"].([]any)) != 1 {
		t.Fatalf("DOING column = %v", columns["DOING"])
	}
	if len(columns["Unknown"].([]any)) != 1 {
		t.Fatalf("Unknown column = %v", columns["Unknown"])
	}
	if counts["TODO"] != 2 || counts["DOING"] != 1 || counts["Unknown"] != 1 {
		t.Fatalf("counts = %v", counts)
	}
}

func TestParseListData(t *testing.T) {
	tests := []struct {
		name string
		data any
		want int // expected length, -1 for nil
	}{
		{"slice", []any{"a", "b"}, 2},
		{"map with data key", map[string]any{"data": []any{"x", "y", "z"}}, 3},
		{"empty map", map[string]any{}, -1},
		{"map with non-list data", map[string]any{"data": "string"}, -1},
		{"string input", "not a list", -1},
		{"int input", 42, -1},
		{"nil input", nil, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseListData(tt.data)
			if tt.want < 0 {
				if got != nil {
					t.Fatalf("parseListData(%v) = %v, want nil", tt.data, got)
				}
			} else if len(got) != tt.want {
				t.Fatalf("parseListData(%v) len = %d, want %d", tt.data, len(got), tt.want)
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
		{"map with int total", map[string]any{"data": []any{map[string]any{"id": "1"}}, "pagination": map[string]any{"total": int(5)}}, 1, 5, false},
		{"map with int64 total", map[string]any{"data": []any{map[string]any{"id": "1"}}, "pagination": map[string]any{"total": int64(7)}}, 1, 7, false},
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

func TestExtractWorkitemStatusName(t *testing.T) {
	tests := []struct {
		name    string
		itemMap map[string]any
		want    string
	}{
		{
			"valid status",
			map[string]any{"status": map[string]any{"name": "DOING"}},
			"DOING",
		},
		{
			"missing status key",
			map[string]any{"id": "wi-1"},
			"Unknown",
		},
		{
			"status is not a map",
			map[string]any{"status": "DOING"},
			"Unknown",
		},
		{
			"status map without name",
			map[string]any{"status": map[string]any{"stage": "dev"}},
			"Unknown",
		},
		{
			"empty map",
			map[string]any{},
			"Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractWorkitemStatusName(tt.itemMap)
			if got != tt.want {
				t.Errorf("extractWorkitemStatusName(%v) = %q, want %q", tt.itemMap, got, tt.want)
			}
		})
	}
}

func TestSearchProjectWorkitemsReturnsError(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	c, _ := getClient(client)
	if _, err := searchProjectWorkitems(context.Background(), c, "org-1", "project-1", "Task", map[string]any{}); err == nil {
		t.Fatal("expected search error")
	}
}
