package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetStructFields(t *testing.T) {
	var pointerInt *int
	// set pointerInt to 123
	pointerInt = new(int)
	*pointerInt = 123

	tests := []struct {
		name     string
		input    interface{}
		expected []Field
		wantErr  error
	}{
		{
			name: "non-pointer struct",
			input: struct {
				Field1 string
				Field2 int
			}{},
			wantErr: ErrInputPointerStruct,
		},
		{
			name:    "non-struct value",
			input:   123,
			wantErr: ErrInputPointerStruct,
		},
		{
			name:    "non-struct pointer",
			input:   pointerInt,
			wantErr: ErrInputPointerStruct,
		},
		{
			name: "no inputs",
			input: &struct {
				Field1 string
				Field2 int
			}{},
			expected: []Field{
				{Name: "Field1", Tags: map[string]string{}},
				{Name: "Field2", Tags: map[string]string{}},
			},
		},
		{
			name: "one tag",
			input: &struct {
				Field1 string `json:"field_1"`
				Field2 int
			}{},
			expected: []Field{
				{Name: "Field1", Tags: map[string]string{"json": "field_1"}},
				{Name: "Field2", Tags: map[string]string{}},
			},
		},
		{
			name: "multiple inputs",
			input: &struct {
				Field1 string `json:"field_1" tag1:"field_1" tag2:"field_1"`
				Field2 int    `json:"field_2" tag2:"field_2"`
			}{},
			expected: []Field{
				{Name: "Field1", Tags: map[string]string{"json": "field_1", "tag1": "field_1", "tag2": "field_1"}},
				{Name: "Field2", Tags: map[string]string{"json": "field_2", "tag2": "field_2"}},
			},
		},
		{
			name: "nested struct",
			input: &struct {
				Field1 string
				Field2 int
				Field3 struct {
					Field1 string `json:"field_1" tag1:"field_1" tag2:"field_1"`
					Field2 int    `json:"field_2" tag2:"field_2"`
				}
				Field4 struct {
					Field1 string `json:"field_1" tag1:"field_1" tag2:"field_1"`
					Field2 int    `json:"field_2" tag2:"field_2"`
				}
			}{},
			expected: []Field{
				{Name: "Field1", Tags: map[string]string{}},
				{Name: "Field2", Tags: map[string]string{}},
				{Name: "Field3.Field1", Tags: map[string]string{"json": "field_1", "tag1": "field_1", "tag2": "field_1"}},
				{Name: "Field3.Field2", Tags: map[string]string{"json": "field_2", "tag2": "field_2"}},
				{Name: "Field4.Field1", Tags: map[string]string{"json": "field_1", "tag1": "field_1", "tag2": "field_1"}},
				{Name: "Field4.Field2", Tags: map[string]string{"json": "field_2", "tag2": "field_2"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetStructFields(tt.input)
			if tt.wantErr != nil {
				assert.ErrorIs(t, tt.wantErr, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, result, len(tt.expected))
			for i, expectedField := range tt.expected {
				assert.Equal(t, expectedField.Name, result[i].Name)
				assert.Equal(t, expectedField.Tags, result[i].Tags)
			}
		})
	}
}

func Test_parseTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:  "spaces in value",
			input: `arg:"cwd" short:"c" help:"Current working directory"`,
			expected: map[string]string{
				"arg":   "cwd",
				"short": "c",
				"help":  "Current working directory",
			},
		},
		{
			name:  "quotes in value",
			input: `arg:"cwd" short:"c" help:"Current \"working directory\""`,
			expected: map[string]string{
				"arg":   "cwd",
				"short": "c",
				"help":  `Current "working directory"`,
			},
		},
		{
			name:  "quotes in value and url",
			input: `arg:"cwd" short:"c" default:"http://127.0.0.1:3888" help:"Current \"working directory\""`,
			expected: map[string]string{
				"arg":     "cwd",
				"short":   "c",
				"default": "http://127.0.0.1:3888",
				"help":    `Current "working directory"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTags(tt.input)

			assert.Equal(t, tt.expected, result)
		})
	}
}
