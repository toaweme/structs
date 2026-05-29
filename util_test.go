package structs

import (
	"reflect"
	"testing"

	"github.com/toaweme/structs/utils"
)

// Test_ParseTags_CommaHandling pins the intended contract for comma handling in
// struct tags: comma-suffixed options (",omitempty", ",flow", ...) are a
// convention of encoding tags only (json/yaml/toml/xml), so only those tags get
// truncated at the first comma. Freeform tags (help, default, rules, arg, ...)
// must keep their value verbatim.
//
// the default/rules cases currently FAIL: the implementation strips commas from
// every tag except "help", so values like default:"a,b,c" and rules:"oneof:a,b,c"
// are silently truncated. these assertions describe the target behavior.
func Test_ParseTags_CommaHandling(t *testing.T) {
	tests := []struct {
		name string
		// input is the raw struct tag line
		input string
		// tag is the key whose value we assert
		tag  string
		want string
	}{
		{
			name:  "json strips omitempty",
			input: `json:"filters,omitempty"`,
			tag:   "json",
			want:  "filters",
		},
		{
			name:  "yaml strips multiple options",
			input: `yaml:"filters,omitempty,flow"`,
			tag:   "yaml",
			want:  "filters",
		},
		{
			name:  "toml strips omitempty",
			input: `toml:"name,omitempty"`,
			tag:   "toml",
			want:  "name",
		},
		{
			name:  "xml strips options",
			input: `xml:"name,attr"`,
			tag:   "xml",
			want:  "name",
		},
		{
			name:  "help keeps commas verbatim",
			input: `help:"do x, then y, then z"`,
			tag:   "help",
			want:  "do x, then y, then z",
		},
		{
			name:  "default keeps commas verbatim",
			input: `default:"a,b,c"`,
			tag:   "default",
			want:  "a,b,c",
		},
		{
			name:  "rules keeps comma-separated args verbatim",
			input: `rules:"oneof:a,b,c"`,
			tag:   "rules",
			want:  "oneof:a,b,c",
		},
		{
			name:  "arg without comma is unaffected",
			input: `arg:"cwd"`,
			tag:   "arg",
			want:  "cwd",
		},
		{
			name:  "short without comma is unaffected",
			input: `short:"c"`,
			tag:   "short",
			want:  "c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTags(tt.input, DefaultEncodingTags)
			if got[tt.tag] != tt.want {
				t.Fatalf("parseTags(%q)[%q] = %q, want %q", tt.input, tt.tag, got[tt.tag], tt.want)
			}
		})
	}
}

// Test_ParseTags_MixedLine asserts that stripping one tag's option does not
// affect sibling tags on the same line.
func Test_ParseTags_MixedLine(t *testing.T) {
	got := parseTags(`json:"filters,omitempty" default:"x,y,z" help:"a, b" arg:"cwd"`, DefaultEncodingTags)
	want := map[string]string{
		"json":    "filters",
		"default": "x,y,z",
		"help":    "a, b",
		"arg":     "cwd",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseTags mixed line = %#v, want %#v", got, want)
	}
}

func Test_ToInt(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    int
		wantErr bool
	}{
		{name: "int", input: 5, want: 5},
		{name: "int8", input: int8(5), want: 5},
		{name: "int16", input: int16(5), want: 5},
		{name: "int32", input: int32(5), want: 5},
		{name: "int64", input: int64(5), want: 5},
		{name: "uint", input: uint(7), want: 7},
		{name: "uint8", input: uint8(7), want: 7},
		{name: "uint64", input: uint64(7), want: 7},
		{name: "whole float32", input: float32(2), want: 2},
		{name: "whole float64", input: float64(3), want: 3},
		{name: "fractional float32 errors", input: float32(2.5), wantErr: true},
		{name: "fractional float64 errors", input: float64(3.5), wantErr: true},
		{name: "numeric string", input: "42", want: 42},
		{name: "non-numeric string errors", input: "x", wantErr: true},
		{name: "bool errors", input: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ToInt(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ToInt(%v) = %d, want error", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ToInt(%v) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("ToInt(%v) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func Test_ToFloat(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    float64
		wantErr bool
	}{
		{name: "float64", input: float64(1.5), want: 1.5},
		{name: "float32", input: float32(2.5), want: 2.5},
		{name: "int", input: 3, want: 3},
		{name: "int64", input: int64(4), want: 4},
		{name: "uint8", input: uint8(5), want: 5},
		{name: "uint64", input: uint64(6), want: 6},
		{name: "numeric string", input: "2.5", want: 2.5},
		{name: "non-numeric string errors", input: "x", wantErr: true},
		{name: "bool errors", input: false, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ToFloat(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ToFloat(%v) = %v, want error", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ToFloat(%v) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("ToFloat(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func Test_ToString(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    string
		wantErr bool
	}{
		{name: "string", input: "hello", want: "hello"},
		{name: "empty string", input: "", want: ""},
		{name: "int", input: 42, want: "42"},
		{name: "float64", input: float64(1.5), want: "1.500000"},
		{name: "float32", input: float32(1.5), want: "1.500000"},
		{name: "bool errors", input: true, wantErr: true},
		{name: "int64 errors", input: int64(5), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ToString(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ToString(%v) = %q, want error", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ToString(%v) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("ToString(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func Test_ParseBool(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "true", input: "true", want: true},
		{name: "true uppercase", input: "TRUE", want: true},
		{name: "yes", input: "yes", want: true},
		{name: "one", input: "1", want: true},
		{name: "false", input: "false", want: false},
		{name: "no", input: "no", want: false},
		{name: "zero", input: "0", want: false},
		{name: "empty", input: "", want: false},
		{name: "garbage falls back to false", input: "maybe", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.ParseBool(tt.input); got != tt.want {
				t.Fatalf("ParseBool(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func Test_ToAnySlice(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    []any
		wantErr bool
	}{
		{name: "any slice passthrough", input: []any{1, "a", true}, want: []any{1, "a", true}},
		{name: "empty string yields empty slice", input: "", want: []any{}},
		{name: "non-empty string yields single element", input: "a", want: []any{"a"}},
		{name: "string slice", input: []string{"a", "b"}, want: []any{"a", "b"}},
		{name: "int slice via reflection", input: []int{1, 2, 3}, want: []any{1, 2, 3}},
		{name: "non-slice errors", input: 5, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ToAnySlice(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ToAnySlice(%v) = %v, want error", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ToAnySlice(%v) unexpected error: %v", tt.input, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ToAnySlice(%v) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}
