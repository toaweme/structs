package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStructWithRules struct {
	Field1 string `json:"field_1" rules:"required"`
	Field2 int    `json:"field_2" rules:"required"`
}

type TestStructWithOneOf struct {
	Format string `json:"format" rules:"oneof:json,yaml,toml" default:"json"`
	Mode   string `json:"mode" rules:"required|oneof:fast,slow"`
}

func Test_Validate_OneOf(t *testing.T) {
	tests := []struct {
		name           string
		values         map[string]any
		expectedErrors map[string][]string
	}{
		{
			name:           "allowed value passes",
			values:         map[string]any{"format": "yaml", "mode": "fast"},
			expectedErrors: map[string][]string{},
		},
		{
			name:           "disallowed value fails with the allowed set",
			values:         map[string]any{"format": "xml", "mode": "fast"},
			expectedErrors: map[string][]string{"format": {"must be one of: json, yaml, toml"}},
		},
		{
			name:           "omitted value falls back to default and passes",
			values:         map[string]any{"mode": "slow"},
			expectedErrors: map[string][]string{},
		},
		{
			name:           "non-empty but disallowed fails oneof, passes required",
			values:         map[string]any{"mode": "turbo"},
			expectedErrors: map[string][]string{"mode": {"must be one of: fast, slow"}},
		},
		{
			name:           "absent required field reports required, not oneof",
			values:         map[string]any{},
			expectedErrors: map[string][]string{"mode": {"required"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := GetStructFields(&TestStructWithOneOf{}, nil, DefaultEncodingTags)
			assert.NoError(t, err)
			errors, err := ValidateStructFields(DefaultRules, fields, tt.values, "json", "json")
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedErrors, errors)
		})
	}
}

func Test_Validate_StructFields(t *testing.T) {
	tests := []struct {
		name           string
		input          any
		tagPriority    []string
		values         map[string]any
		expectedErrors map[string][]string
		wantErr        error
	}{
		{
			name:        "all json fields valid",
			input:       &TestStructWithRules{},
			tagPriority: []string{"json"},
			values: map[string]any{
				"field_1": "field_1_value",
				"field_2": "field_2_value",
			},
			expectedErrors: map[string][]string{},
		},
		{
			name:        "one json field invalid",
			input:       &TestStructWithRules{},
			tagPriority: []string{"json"},
			values: map[string]any{
				"field_2": "field_2_value",
			},
			expectedErrors: map[string][]string{
				"field_1": {"required"},
			},
		},
		{
			name:        "all json fields invalid",
			input:       &TestStructWithRules{},
			tagPriority: []string{"json"},
			values:      map[string]any{},
			expectedErrors: map[string][]string{
				"field_1": {"required"},
				"field_2": {"required"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := GetStructFields(tt.input, nil, DefaultEncodingTags)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			errors, err := ValidateStructFields(DefaultRules, fields, tt.values, "json", tt.tagPriority...)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedErrors, errors)
		})
	}
}
