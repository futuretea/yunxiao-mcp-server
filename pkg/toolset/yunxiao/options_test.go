package yunxiao

import (
	"net/url"
	"reflect"
	"slices"
	"testing"
)

func TestSetOptionalIntBody(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]any
		key    string
		want   any
		wantOk bool
	}{
		{"float64", map[string]any{"k": float64(42)}, "k", 42, true},
		{"int", map[string]any{"k": int(7)}, "k", 7, true},
		{"int64", map[string]any{"k": int64(99)}, "k", int64(99), true},
		{"string non-empty", map[string]any{"k": "123"}, "k", "123", true},
		{"string empty", map[string]any{"k": ""}, "k", nil, false},
		{"nil", map[string]any{}, "k", nil, false},
		{"bool", map[string]any{"k": true}, "k", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]any{}
			setOptionalIntBody(body, tt.params, tt.key)
			if tt.wantOk {
				if !reflect.DeepEqual(body[tt.key], tt.want) {
					t.Fatalf("body[%q] = %v (%T), want %v (%T)", tt.key, body[tt.key], body[tt.key], tt.want, tt.want)
				}
			} else {
				if _, ok := body[tt.key]; ok {
					t.Fatalf("body[%q] should not be set", tt.key)
				}
			}
		})
	}
}

func TestSetOptionalIntAs(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]any
		from   string
		to     string
		want   string
		wantOk bool
	}{
		{"float64", map[string]any{"k": float64(42)}, "k", "k", "42", true},
		{"int", map[string]any{"k": int(7)}, "k", "k", "7", true},
		{"int64", map[string]any{"k": int64(99)}, "k", "k", "99", true},
		{"string non-empty", map[string]any{"k": "123"}, "k", "k", "123", true},
		{"string empty", map[string]any{"k": ""}, "k", "k", "", false},
		{"nil", map[string]any{}, "k", "k", "", false},
		{"bool", map[string]any{"k": true}, "k", "k", "", false},
		{"rename", map[string]any{"from": "5"}, "from", "to", "5", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := url.Values{}
			setOptionalIntAs(query, tt.params, tt.from, tt.to)
			if tt.wantOk {
				if got := query.Get(tt.to); got != tt.want {
					t.Fatalf("query[%q] = %q, want %q", tt.to, got, tt.want)
				}
			} else {
				if query.Get(tt.to) != "" {
					t.Fatalf("query[%q] should not be set", tt.to)
				}
			}
		})
	}
}

func TestSetOptionalStringArrayBody(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]any
		key    string
		want   []string
		wantOk bool
	}{
		{"[]any with values", map[string]any{"k": []any{"a", "b"}}, "k", []string{"a", "b"}, true},
		{"[]any with whitespace", map[string]any{"k": []any{" a ", "  b  "}}, "k", []string{"a", "b"}, true},
		{"[]any mixed types", map[string]any{"k": []any{"a", 1, "b"}}, "k", []string{"a", "b"}, true},
		{"[]any all non-string", map[string]any{"k": []any{1, 2}}, "k", nil, false},
		{"[]any empty", map[string]any{"k": []any{}}, "k", nil, false},
		{"[]string with values", map[string]any{"k": []string{"a", "b"}}, "k", []string{"a", "b"}, true},
		{"[]string with whitespace", map[string]any{"k": []string{" a ", "b"}}, "k", []string{"a", "b"}, true},
		{"[]string all empty", map[string]any{"k": []string{"", "  "}}, "k", nil, false},
		{"[]string empty", map[string]any{"k": []string{}}, "k", nil, false},
		{"nil", map[string]any{}, "k", nil, false},
		{"string not array", map[string]any{"k": "not-array"}, "k", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]any{}
			setOptionalStringArrayBody(body, tt.params, tt.key)
			if tt.wantOk {
				got, ok := body[tt.key].([]string)
				if !ok {
					t.Fatalf("body[%q] type = %T, want []string", tt.key, body[tt.key])
				}
				if !slices.Equal(got, tt.want) {
					t.Fatalf("body[%q] = %v, want %v", tt.key, got, tt.want)
				}
			} else {
				if _, ok := body[tt.key]; ok {
					t.Fatalf("body[%q] should not be set", tt.key)
				}
			}
		})
	}
}
