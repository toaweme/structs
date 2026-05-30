package structs

import (
	"testing"
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
			wantErr: ErrInputPointer,
		},
		{
			name:    "non-struct value",
			input:   123,
			wantErr: ErrInputPointer,
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
				{Name: "Field1", Tags: map[string]string{}, Type: "string"},
				{Name: "Field2", Tags: map[string]string{}, Type: "int"},
			},
		},
		{
			name: "one tag",
			input: &struct {
				Field1 string `json:"field_1"`
				Field2 int
			}{},
			expected: []Field{
				{Name: "Field1", Tags: map[string]string{"json": "field_1"}, Type: "string"},
				{Name: "Field2", Tags: map[string]string{}, Type: "int"},
			},
		},
		{
			name: "multiple inputs",
			input: &struct {
				Field1 string `json:"field_1" tag1:"field_1" tag2:"field_1"`
				Field2 int    `json:"field_2" tag2:"field_2"`
			}{},
			expected: []Field{
				{Name: "Field1", Tags: map[string]string{"json": "field_1", "tag1": "field_1", "tag2": "field_1"}, Type: "string"},
				{Name: "Field2", Tags: map[string]string{"json": "field_2", "tag2": "field_2"}, Type: "int"},
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
				} `json:"field_3" tag1:"field_3" tag2:"field_3"`
				Field4 struct {
					Field1 string `json:"field_1" tag1:"field_1" tag2:"field_1"`
					Field2 int    `json:"field_2" tag2:"field_2"`
				} `json:"field_4" tag1:"field_4" tag2:"field_4"`
			}{},
			expected: []Field{
				{Name: "Field1", Tags: map[string]string{}, Type: "string"},
				{Name: "Field2", Tags: map[string]string{}, Type: "int"},
				{
					Name: "Field3",
					Tags: map[string]string{"json": "field_3", "tag1": "field_3", "tag2": "field_3"},
					Type: "struct", Fields: []Field{
						{
							Name: "Field1",
							Type: "string",
							Tags: map[string]string{"json": "field_1", "tag1": "field_1", "tag2": "field_1"},
							FQN:  &Field{Name: "Field3.Field1", Tags: map[string]string{"json": "field_3.field_1", "tag1": "field_3.field_1", "tag2": "field_3.field_1"}},
						},
						{
							Name: "Field2",
							Type: "int",
							Tags: map[string]string{"json": "field_2", "tag2": "field_2"},
							FQN:  &Field{Name: "Field3.Field2", Tags: map[string]string{"json": "field_3.field_2", "tag2": "field_3.field_2"}},
						},
					}},
				{
					Name: "Field4",
					Tags: map[string]string{"json": "field_4", "tag1": "field_4", "tag2": "field_4"},
					Type: "struct",
					Fields: []Field{
						{
							Name: "Field1",
							Type: "string",
							Tags: map[string]string{"json": "field_1", "tag1": "field_1", "tag2": "field_1"},
							FQN:  &Field{Name: "Field4.Field1", Tags: map[string]string{"json": "field_4.field_1", "tag1": "field_4.field_1", "tag2": "field_4.field_1"}},
						},
						{
							Name: "Field2",
							Type: "int",
							Tags: map[string]string{"json": "field_2", "tag2": "field_2"},
							FQN:  &Field{Name: "Field4.Field2", Tags: map[string]string{"json": "field_4.field_2", "tag2": "field_4.field_2"}},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetStructFields(tt.input, nil, DefaultEncodingTags)
			if tt.wantErr != nil {
				requireErrorIs(t, err, tt.wantErr)
				return
			}
			requireNoError(t, err)
			requireLen(t, result, len(tt.expected))
			for i, res := range result {
				expectedField := tt.expected[i]
				requireEqual(t, expectedField.Name, res.Name, "Name")
				requireEqual(t, expectedField.Type, res.Type, "Type")
				requireEqual(t, expectedField.Tags, res.Tags, "Tags")
				requireEqual(t, expectedField.Default, res.Default, "Default")
				for j, resField := range result[i].Fields {
					expectedSubField := expectedField.Fields[j]
					requireEqual(t, expectedSubField.Name, resField.Name, "Sub.Name")
					requireEqual(t, expectedSubField.Tags, resField.Tags, "Sub.Tags")
					requireEqual(t, expectedSubField.Type, resField.Type, "Sub.Type")
					requireEqual(t, expectedSubField.Default, resField.Default, "Sub.Default")
					if expectedSubField.FQN != nil {
						requireNotNil(t, resField.FQN, "FQN")
						requireEqual(t, expectedSubField.FQN.Name, resField.FQN.Name, "FQN Name")
						requireEqual(t, expectedSubField.FQN.Tags, resField.FQN.Tags, "FQN Tags")
						requireEqual(t, expectedSubField.FQN.Type, resField.FQN.Type, "FQN Type")
						requireEqual(t, expectedSubField.FQN.Default, resField.FQN.Default, "FQN Default")
					}
				}
			}
		})
	}
}

//
// func assertField(t *testing.T, expectedField Field, result Field) {
// 	requireEqual(t, expectedField.Name, result.Name, "Name")
// 	requireEqual(t, expectedField.Type, result.Type, "Type")
// 	requireEqual(t, expectedField.Tags, result.Tags, "Tags")
// 	requireEqual(t, expectedField.Default, result.Default, "Default")
// 	if expectedField.FQN != nil {
// 		assertField(t, *expectedField.FQN, *result.FQN)
// 	}
// 	for i, expectedNestedField := range expectedField.Fields {
// 		assertField(t, expectedNestedField, result.Fields[i])
// 	}
// }

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
			result := parseTags(tt.input, DefaultEncodingTags)

			requireEqual(t, tt.expected, result)
		})
	}
}
