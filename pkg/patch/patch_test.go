package patch

import (
	"reflect"
	"testing"
)

func TestPatch(t *testing.T) {
	tests := []struct {
		name     string
		base     any
		patch    any
		expected any
	}{
		{
			name:     "Simple merge",
			base:     map[string]any{"a": 1},
			patch:    map[string]any{"b": 2},
			expected: map[string]any{"a": 1, "b": 2},
		},
		{
			name:     "Recursive merge",
			base:     map[string]any{"nested": map[string]any{"a": 1}},
			patch:    map[string]any{"nested": map[string]any{"b": 2}},
			expected: map[string]any{"nested": map[string]any{"a": 1, "b": 2}},
		},
		{
			name:     "Force overwrite with !",
			base:     map[string]any{"a": map[string]any{"old": 1}},
			patch:    map[string]any{"a!": "new"},
			expected: map[string]any{"a": "new"},
		},
		{
			name:     "Append to slice with + suffix",
			base:     map[string]any{"list": []any{1, 2}},
			patch:    map[string]any{"list+": []any{3, 4}},
			expected: map[string]any{"list": []any{1, 2, 3, 4}},
		},
		{
			name:     "Prepend to slice with + prefix",
			base:     map[string]any{"list": []any{1, 2}},
			patch:    map[string]any{"+list": []any{3, 4}},
			expected: map[string]any{"list": []any{3, 4, 1, 2}},
		},
		{
			name:     "Key extraction with <>",
			base:     map[string]any{"realKey": 1},
			patch:    map[string]any{"<realKey>": 2},
			expected: map[string]any{"realKey": 2},
		},
		{
			name:     "Non-map patch replaces base",
			base:     map[string]any{"a": 1},
			patch:    "just a string",
			expected: "just a string",
		},
		{
			name:     "Type mismatch on append (falls through to overwrite)",
			base:     map[string]any{"list": "not a list"},
			patch:    map[string]any{"list+": []any{1}},
			expected: map[string]any{"list": []any{1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Patch(tt.base, tt.patch)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Patch() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPatch_ComplexCombinations(t *testing.T) {
	tests := []struct {
		name     string
		base     any
		patch    any
		expected any
	}{
		{
			name:     "Escaped key with force overwrite",
			base:     map[string]any{"tags": []any{"a", "b"}},
			patch:    map[string]any{"<tags>!": []any{"c"}},
			expected: map[string]any{"tags": []any{"c"}},
		},
		{
			name:     "Escaped key with prepend",
			base:     map[string]any{"list": []any{2, 3}},
			patch:    map[string]any{"+<list>": []any{1}},
			expected: map[string]any{"list": []any{1, 2, 3}},
		},
		{
			name:     "Escaped key with append",
			base:     map[string]any{"list": []any{1, 2}},
			patch:    map[string]any{"<list>+": []any{3}},
			expected: map[string]any{"list": []any{1, 2, 3}},
		},
		{
			name:     "Ambiguous prefix and suffix (Prepend priority)",
			base:     map[string]any{"list": []any{"middle"}},
			patch:    map[string]any{"+<list>+": []any{"start"}},
			expected: map[string]any{"list": []any{"start", "middle"}},
		},
		{
			name:     "Symbols inside brackets",
			base:     map[string]any{"my+key!": 10},
			patch:    map[string]any{"<my+key!>!": 20},
			expected: map[string]any{"my+key!": 20},
		},
		{
			name:     "Escaped overwrite on existing map",
			base:     map[string]any{"nested": map[string]any{"keep": "me"}},
			patch:    map[string]any{"<nested>!": "replaced"},
			expected: map[string]any{"nested": "replaced"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Patch(tt.base, tt.patch)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("%s: Patch() = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestRealKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"standard", "standard"},
		{"<wrapped>", "wrapped"},
		{"prefix<inside>suffix", "inside"},
		{"no-closing-<bracket", "no-closing-<bracket"},
		{"no-opening-bracket>", "no-opening-bracket>"},
		{"<>", ""},
		{"<reverse><order>", "reverse><order"},
		{"", ""},
		{"+!", "+"},
		{"++", "+"},
		{"!!", "!"},
		{"+prepend", "prepend"},
		{"append+", "append"},
		{"overwrite!", "overwrite"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := realKey(tt.input)
			if got != tt.expected {
				t.Errorf("realKey(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
