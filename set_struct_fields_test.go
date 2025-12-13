package structs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetStructFields(t *testing.T) {
	var pointerInt *int
	tests := []struct {
		name      string
		structure any
		inputs    map[string]any
		settings  Settings
		expected  any
		wantErr   error
	}{
		{
			name:      "non-pointer struct",
			structure: struct{}{},
			inputs:    map[string]any{},
			wantErr:   ErrInputPointer,
		},
		{
			name:      "non-struct value",
			structure: pointerInt,
			inputs:    map[string]any{},
			wantErr:   ErrInputPointerStruct,
		},
		{
			name: "inputs as field names when AllowTagOverride is true",
			structure: &struct {
				Field1 string
				Field2 int
			}{},
			settings: Settings{
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"Field1": "value101",
				"Field2": 10,
			},
			expected: &struct {
				Field1 string
				Field2 int
			}{Field1: "value101", Field2: 10},
		},
		{
			name: "use last found tag when AllowTagOverride is true",
			structure: &struct {
				Field1 string `name:"field_1" tag2:"field_12" tag1:"field_11"`
				Field2 int
			}{},
			settings: Settings{
				TagOrder:         []string{"json", "tag1", "tag2"},
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"field_1":  "value100",
				"field_11": "value1000",
				"field_12": "value10000",
				"Field2":   100,
			},
			expected: &struct {
				Field1 string `name:"field_1" tag1:"field_11" tag2:"field_12"`
				Field2 int
			}{
				Field1: "value10000",
				Field2: 100,
			},
		},
		{
			name: "use first found tag when AllowTagOverride is false",
			structure: &struct {
				Field1 string `name:"field_1" tag1:"field_12" tag2:"field_11"`
				Field2 int    `name:"field_2" tag2:"field_22" tag1:"field_21"`
			}{},
			settings: Settings{
				TagOrder:         []string{"json", "tag1", "tag2"},
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"field_1":  "value1000",
				"field_11": "value10000",
				"field_12": "value100000",
				"field_2":  1000,
				"field_21": 10000,
				"field_22": 100000,
			},
			expected: &struct {
				Field1 string `name:"field_1" tag1:"field_1" tag2:"field_1"`
				Field2 int    `name:"field_2" tag2:"field_2" tag2:"field_2"`
			}{
				Field1: "value10000",
				Field2: 100000,
			},
		},
		{
			name: "nested struct",
			structure: &struct {
				Field1 string `name:"field_1"`
				Field2 int    `name:"field_2"`
				Field3 struct {
					Field1 string `name:"field_1"`
					Field2 int    `name:"field_2"`
				} `name:"field_3"`
				Field4 struct {
					Field1 string `name:"field_1"`
					Field2 int    `name:"field_2"`
				} `name:"field_4"`
			}{},
			settings: Settings{
				TagOrder:         []string{"name"},
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"field_1":         "------",
				"field_2":         11111,
				"field_3.field_1": "field3::field_1",
				"field_3.field_2": 10000,
				"field_4.field_1": "field4::field_1",
				"field_4.field_2": 100000,
			},
			expected: &struct {
				Field1 string
				Field2 int
				Field3 struct {
					Field1 string
					Field2 int
				}
				Field4 struct {
					Field1 string
					Field2 int
				}
			}{
				Field1: "------",
				Field2: 11111,
				Field3: struct {
					Field1 string
					Field2 int
				}{
					Field1: "field3::field_1",
					Field2: 10000,
				},
				Field4: struct {
					Field1 string
					Field2 int
				}{
					Field1: "field4::field_1",
					Field2: 100000,
				},
			},
		},
		{
			name: "cli app struct by env only",
			structure: &struct {
				Inner struct {
					ComplexFlag string `arg:"complex-flag" short:"c" env:"COMPLEX_FLAG"`
					SimpleFlag  bool   `arg:"simple-flag" short:"s" env:"SIMPLE_FLAG"`
				} `arg:"inner" short:"i" env:"INNER"`
				Outer int `arg:"outer" short:"o" env:"OUTER"`
			}{},
			settings: Settings{
				TagOrder:         []string{"arg", "short"},
				AllowEnvOverride: false,
			},
			inputs: map[string]any{
				"INNER_COMPLEX_FLAG": "complex-flag-value-env",
				"INNER_SIMPLE_FLAG":  "yes",
				"OUTER":              123,
				"i.c":                "i.c",
				"i.s":                "i.s",
				"o":                  -1,
			},
			expected: &struct {
				Inner struct {
					ComplexFlag string `arg:"complex-flag" short:"c" env:"COMPLEX_FLAG"`
					SimpleFlag  bool   `arg:"simple-flag" short:"s" env:"SIMPLE_FLAG"`
				} `arg:"inner" short:"i" env:"INNER"`
				Outer int `arg:"outer" short:"o" env:"OUTER"`
			}{
				Inner: struct {
					ComplexFlag string `arg:"complex-flag" short:"c" env:"COMPLEX_FLAG"`
					SimpleFlag  bool   `arg:"simple-flag" short:"s" env:"SIMPLE_FLAG"`
				}{
					ComplexFlag: "complex-flag-value-env",
					SimpleFlag:  true,
				},
				Outer: 123,
			},
		},
		{
			name: "cli app struct with env override by tag",
			structure: &struct {
				Inner struct {
					ComplexFlag string `arg:"complex-flag" short:"c" env:"COMPLEX_FLAG" help:"Complex flag"`
					SimpleFlag  bool   `arg:"simple-flag" short:"s" env:"SIMPLE_FLAG" help:"Simple flag"`
				} `arg:"inner" short:"i" env:"INNER" help:"Inner app"`
				Outer int `arg:"outer" short:"o" env:"OUTER" help:"Outer app"`
			}{},
			settings: Settings{
				TagOrder:         []string{"arg", "short"},
				AllowEnvOverride: true,
				// needs to be true to test from `arg` to `short` tag override
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"INNER_COMPLEX_FLAG": "envValue",
				"INNER_SIMPLE_FLAG":  "yes",
				"OUTER":              123,
				"i.c":                "i.c",
				"i.s":                false,
				"o":                  -1,
			},
			expected: &struct {
				Inner struct {
					ComplexFlag string `arg:"complex-flag" short:"c" env:"COMPLEX_FLAG" help:"Complex flag"`
					SimpleFlag  bool   `arg:"simple-flag" short:"s" env:"SIMPLE_FLAG" help:"Simple flag"`
				} `arg:"inner" short:"i" env:"INNER" help:"Inner app"`
				Outer int `arg:"outer" short:"o" env:"OUTER" help:"Outer app"`
			}{
				Inner: struct {
					ComplexFlag string `arg:"complex-flag" short:"c" env:"COMPLEX_FLAG" help:"Complex flag"`
					SimpleFlag  bool   `arg:"simple-flag" short:"s" env:"SIMPLE_FLAG" help:"Simple flag"`
				}{
					ComplexFlag: "i.c",
					SimpleFlag:  false,
				},
				Outer: -1,
			},
		},
		{
			name: "slices of strings",
			structure: &struct {
				Strings []string `name:"strings"`
			}{},
			inputs: map[string]any{},
			expected: &struct {
				Strings []string `name:"strings"`
			}{
				Strings: nil,
			},
		},
		{
			name: "ensure inner struct is not nil",
			structure: &struct {
				Inner struct {
					A string
				}
			}{},
			inputs: map[string]any{},
			expected: &struct {
				Inner struct {
					A string
				}
			}{
				Inner: struct {
					A string
				}{},
			},
		},
		{
			name: "default values are set",
			structure: &struct {
				Field1 string `name:"field_1" default:"default1"`
				Field2 int    `name:"field_2" default:"100"`
			}{
				Field1: "",
				Field2: 0,
			},
			expected: &struct {
				Field1 string `name:"field_1" default:"default1"`
				Field2 int    `name:"field_2" default:"100"`
			}{
				Field1: "default1",
				Field2: 100,
			},
		},
		{
			name: "default values don't override existing values",
			structure: &struct {
				Field1 string `name:"field_1" default:"default1"`
				Field2 int    `name:"field_2" default:"100"`
				Field3 bool   `name:"field_3" default:"yes"`
			}{
				Field1: "pem",
				Field2: 0,
				Field3: false,
			},
			expected: &struct {
				Field1 string `name:"field_1" default:"default1"`
				Field2 int    `name:"field_2" default:"100"`
				Field3 bool   `name:"field_3" default:"yes"`
			}{
				Field1: "pem",
				Field2: 100,
				Field3: true,
			},
		},
		{
			name: "values can be set via nested map[string]any",
			settings: Settings{
				TagOrder: []string{"name", "json"},
			},
			structure: &struct {
				Outer struct {
					Inner struct {
						Value string `name:"value" default:"default1"`
					} `name:"inner"`
				} `name:"outer"`
			}{
				Outer: struct {
					Inner struct {
						Value string `name:"value" default:"default1"`
					} `name:"inner"`
				}{},
			},
			expected: &struct {
				Outer struct {
					Inner struct {
						Value string `name:"value" default:"default1"`
					} `name:"inner"`
				} `name:"outer"`
			}{
				Outer: struct {
					Inner struct {
						Value string `name:"value" default:"default1"`
					} `name:"inner"`
				}{
					Inner: struct {
						Value string `name:"value" default:"default1"`
					}{
						Value: "woo!",
					},
				},
			},
			inputs: map[string]any{
				"outer": map[string]any{
					"inner": map[string]any{
						"value": "woo!",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetStructFields(tt.structure, tt.settings, tt.inputs)
			if tt.wantErr != nil {
				assert.ErrorIs(t, tt.wantErr, err)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.expected, tt.structure)
		})
	}
}

type metadata struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

func Test_SetStructFieldsWithStructSlice(t *testing.T) {
	tests := []struct {
		name      string
		structure any
		inputs    map[string]any
		settings  Settings
		expected  any
		wantErr   error
	}{
		{
			name: "inputs are slices",
			structure: &struct {
				Field1 []string
				Field2 []int
			}{},
			settings: Settings{
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"Field1": []string{"value101"},
				"Field2": []int{10},
			},
			expected: &struct {
				Field1 []string
				Field2 []int
			}{Field1: []string{"value101"}, Field2: []int{10}},
		},
		{
			name: "inputs are slices of structs",
			structure: &struct {
				Field1 []metadata
			}{},
			settings: Settings{
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"Field1": []metadata{
					{
						Title:       "value101",
						URL:         "value102",
						Description: "value103",
					},
				},
			},
			expected: &struct {
				Field1 []metadata
			}{Field1: []metadata{
				{
					Title:       "value101",
					URL:         "value102",
					Description: "value103",
				},
			}},
		},
		{
			name: "inputs are slices of map[string]any",
			structure: &struct {
				Field1 []metadata
			}{},
			settings: Settings{
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"Field1": []map[string]any{
					{
						"title":       "value101",
						"url":         "value102",
						"description": "value103",
					},
				},
			},
			expected: &struct {
				Field1 []metadata
			}{Field1: []metadata{
				{
					Title:       "value101",
					URL:         "value102",
					Description: "value103",
				},
			}},
		},
		// Add these test cases to your existing Test_SetStructFieldsWithStructSlice function

		{
			name: "slice of bools from interface slice",
			structure: &struct {
				Flags []bool
			}{},
			settings: Settings{
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"Flags": []any{true, false, true, "true", "false"},
			},
			expected: &struct {
				Flags []bool
			}{Flags: []bool{true, false, true, true, false}},
			wantErr: errors.New("failed to set field[Flags]: cannot assign or convert string to bool"),
		},
		{
			name: "slice of floats from mixed numeric types",
			structure: &struct {
				Scores []float64
			}{},
			settings: Settings{
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"Scores": []any{1.5, 2, 3.7, "4.2"},
			},
			expected: &struct {
				Scores []float64
			}{Scores: []float64{1.5, 2.0, 3.7, 4.2}},
			wantErr: errors.New("failed to set field[Scores]: cannot assign or convert string to float64"),
		},
		{
			name: "nested struct slice with maps of different key types",
			structure: &struct {
				Items []struct {
					ID    int
					Name  string
					Price float64
				}
			}{},
			settings: Settings{
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"Items": []map[string]any{
					{"id": 1, "name": "Item1", "price": 10.5},
					{"id": "2", "name": "Item2", "price": 20.75},
				},
			},
			expected: &struct {
				Items []struct {
					ID    int
					Name  string
					Price float64
				}
			}{Items: []struct {
				ID    int
				Name  string
				Price float64
			}{
				{ID: 1, Name: "Item1", Price: 10.5},
				{ID: 2, Name: "Item2", Price: 20.75},
			}},
		},
		{
			name: "slice of structs with nested slices",
			structure: &struct {
				Users []struct {
					Name  string
					Roles []string
				}
			}{},
			settings: Settings{
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"Users": []map[string]any{
					{"name": "Alice", "roles": []string{"admin", "user"}},
					{"name": "Bob", "roles": []any{"user"}},
				},
			},
			expected: &struct {
				Users []struct {
					Name  string
					Roles []string
				}
			}{Users: []struct {
				Name  string
				Roles []string
			}{
				{Name: "Alice", Roles: []string{"admin", "user"}},
				{Name: "Bob", Roles: []string{"user"}},
			}},
		},
		{
			name: "empty slice initialization",
			structure: &struct {
				Tags []string
			}{},
			settings: Settings{
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"Tags": []string{},
			},
			expected: &struct {
				Tags []string
			}{Tags: []string{}},
		},
		{
			name: "slice of interfaces with mixed types",
			structure: &struct {
				Mixed []interface{}
			}{},
			settings: Settings{
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"Mixed": []any{"string", 123, 45.67, true, nil},
			},
			expected: &struct {
				Mixed []interface{}
			}{Mixed: []interface{}{"string", 123, 45.67, true, nil}},
		},
		{
			name: "complex nested structure with multiple slice levels",
			structure: &struct {
				Departments []struct {
					Name  string
					Teams []struct {
						TeamName string
						Members  []string
					}
				}
			}{},
			settings: Settings{
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"Departments": []map[string]any{
					{
						"name": "Engineering",
						"teams": []map[string]any{
							{"teamname": "Backend", "members": []string{"Alice", "Bob"}},
							{"teamname": "Frontend", "members": []string{"Charlie"}},
						},
					},
				},
			},
			expected: &struct {
				Departments []struct {
					Name  string
					Teams []struct {
						TeamName string
						Members  []string
					}
				}
			}{Departments: []struct {
				Name  string
				Teams []struct {
					TeamName string
					Members  []string
				}
			}{
				{
					Name: "Engineering",
					Teams: []struct {
						TeamName string
						Members  []string
					}{
						{TeamName: "Backend", Members: []string{"Alice", "Bob"}},
						{TeamName: "Frontend", Members: []string{"Charlie"}},
					},
				},
			}},
		},
		{
			name: "slice of maps (not structs) - map type fields",
			structure: &struct {
				Data []map[string]string
			}{},
			settings: Settings{
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"Data": []map[string]string{
					{"key1": "value1", "key2": "value2"},
					{"key3": "value3"},
				},
			},
			expected: &struct {
				Data []map[string]string
			}{Data: []map[string]string{
				{"key1": "value1", "key2": "value2"},
				{"key3": "value3"},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetStructFields(tt.structure, tt.settings, tt.inputs)
			if tt.wantErr != nil {
				assert.ErrorContains(t, tt.wantErr, err.Error())
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.expected, tt.structure)
		})
	}
}
