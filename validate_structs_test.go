package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStructWithRules struct {
	Field1 string `json:"field_1" rules:"required"`
	Field2 int    `json:"field_2" rules:"required"`
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
			fields, err := GetStructFields(tt.input, nil)
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
