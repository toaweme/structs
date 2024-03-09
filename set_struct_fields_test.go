package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStructWithoutTags struct {
	Field1 string
	Field2 int
}

type TestStructWithOneFieldWithTags struct {
	Field1 string `json:"field_1" tag1:"field_1" tag2:"field_1"`
	Field2 int
}

type TestStructWithDefaultTags struct {
	Field1 string `json:"field_1" tag1:"field_1" tag2:"field_1"`
	Field2 int    `json:"field_2" tag2:"field_2" tag2:"field_2"`
}

type TestNestedStruct struct {
	Field1 string
	Field2 int
	Field3 TestStructWithDefaultTags `json:"field_3" tag1:"field_3" tag2:"field_3"`
	Field4 TestStructWithDefaultTags
}

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
			name:  "no tags",
			input: &TestStructWithoutTags{},
			expected: []Field{
				{Name: "Field1", Tags: map[string]string{}},
				{Name: "Field2", Tags: map[string]string{}},
			},
		},
		{
			name:  "one tag",
			input: &TestStructWithOneFieldWithTags{},
			expected: []Field{
				{Name: "Field1", Tags: map[string]string{"json": "field_1", "tag1": "field_1", "tag2": "field_1"}},
				{Name: "Field2", Tags: map[string]string{}},
			},
		},
		{
			name:  "multiple tags",
			input: &TestStructWithDefaultTags{},
			expected: []Field{
				{Name: "Field1", Tags: map[string]string{"json": "field_1", "tag1": "field_1", "tag2": "field_1"}},
				{Name: "Field2", Tags: map[string]string{"json": "field_2", "tag2": "field_2"}},
			},
		},
		{
			name:  "nested struct",
			input: &TestNestedStruct{},
			expected: []Field{
				{Name: "Field1", Tags: map[string]string{}},
				{Name: "Field2", Tags: map[string]string{}},
				{Name: "Field3.Field1", Tags: map[string]string{"json": "field_1", "tag1": "field_1", "tag2": "field_1"}},
				{Name: "Field3.Field2", Tags: map[string]string{"json": "field_2", "tag2": "field_2"}},
				{Name: "Field4.Field1", Tags: map[string]string{"json": "field_1", "tag1": "field_1", "tag2": "field_1"}},
				{Name: "Field4.Field2", Tags: map[string]string{"json": "field_2", "tag2": "field_2"}},
			},
		},
		{
			name:    "non-pointer struct",
			input:   TestStructWithoutTags{},
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

func Test_SetStructFields(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		tags     map[string]interface{}
		expected interface{}
		wantErr  error
		tagOrder []string
	}{
		{
			name:  "no tags",
			input: &TestStructWithoutTags{},
			tags: map[string]interface{}{
				"Field1": "value10",
				"Field2": 10,
			},
			expected: &TestStructWithoutTags{
				Field1: "value10",
				Field2: 10,
			},
		},
		{
			name:     "one tag",
			input:    &TestStructWithOneFieldWithTags{},
			tagOrder: []string{"json", "tag1", "tag2"},
			tags: map[string]interface{}{
				"field_1": "value100",
				"Field2":  100,
			},
			expected: &TestStructWithOneFieldWithTags{
				Field1: "value100",
				Field2: 100,
			},
		},
		{
			name:     "multiple tags",
			input:    &TestStructWithDefaultTags{},
			tagOrder: []string{"json", "tag1", "tag2"},
			tags: map[string]interface{}{
				"field_1": "value1000",
				"field_2": 1000,
			},
			expected: &TestStructWithDefaultTags{
				Field1: "value1000",
				Field2: 1000,
			},
		},
		{
			name: "nested struct",
			input: &TestNestedStruct{
				Field3: TestStructWithDefaultTags{
					Field1: "value10000",
					Field2: 10000,
				},
			},
			tagOrder: []string{"json", "tag1", "tag2"},
			tags: map[string]interface{}{
				"Field1":          "11111",
				"Field2":          11111,
				"field_3.field_1": "value10000",
				"field_3.field_2": 10000,
				"Field4.field_1":  "value11111",
				"Field4.field_2":  11111,
			},
			expected: &TestNestedStruct{
				Field1: "11111",
				Field2: 11111,
				Field3: TestStructWithDefaultTags{
					Field1: "value10000",
					Field2: 10000,
				},
				Field4: TestStructWithDefaultTags{
					Field1: "value11111",
					Field2: 11111,
				},
			},
		},
		{
			name:    "non-pointer struct",
			input:   TestStructWithoutTags{},
			tags:    map[string]interface{}{},
			wantErr: ErrInputPointerStruct,
		},
		{
			name:    "non-struct value",
			input:   123,
			tags:    map[string]interface{}{},
			wantErr: ErrInputPointerStruct,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetStructFields(tt.input, tt.tagOrder, tt.tags)
			if tt.wantErr != nil {
				assert.ErrorIs(t, tt.wantErr, err)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.expected, tt.input)
		})
	}
}
