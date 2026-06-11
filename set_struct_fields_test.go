package structs

import (
	"testing"
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
				Field2 int    `name:"field_2" tag1:"field_2" tag2:"field_2"`
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
				requireErrorIs(t, err, tt.wantErr)
				return
			}
			requireNoError(t, err)
			requireEqual(t, tt.expected, tt.structure)
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
				Mixed []any
			}{},
			settings: Settings{
				AllowTagOverride: true,
			},
			inputs: map[string]any{
				"Mixed": []any{"string", 123, 45.67, true, nil},
			},
			expected: &struct {
				Mixed []any
			}{Mixed: []any{"string", 123, 45.67, true, nil}},
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
				requireErrorContains(t, err, tt.wantErr.Error())
				return
			}
			requireNoError(t, err)
			requireEqual(t, tt.expected, tt.structure)
		})
	}
}

func Test_SetStructFields_NestedTagOmitempty(t *testing.T) {
	// reproduces the bug where nested struct fields whose json tags carry
	// ",omitempty" (or similar suffixes) silently fail to populate from a
	// nested map[string]any input. before the parseTags fix, fqn lookups
	// were keyed against "query.filters,omitempty" instead of "query.filters".

	type inner struct {
		Filters []map[string]any `json:"filters,omitempty"`
		Limit   int              `json:"limit,omitempty"`
		Offset  int              `json:"offset,omitempty"`
	}

	type outer struct {
		OrgID string `json:"org_id"`
		Query inner  `json:"query"`
	}

	got := &outer{}
	inputs := map[string]any{
		"org_id": "org-1",
		"query": map[string]any{
			"filters": []map[string]any{
				{"field": "bank_iban", "op": "eq", "value": "LT123"},
			},
			"limit":  5,
			"offset": 10,
		},
	}

	err := SetStructFields(got, Settings{TagOrder: DefaultTags, EncodingTags: DefaultEncodingTags}, inputs)
	requireNoError(t, err)
	requireEqual(t, "org-1", got.OrgID)
	requireEqual(t, 5, got.Query.Limit)
	requireEqual(t, 10, got.Query.Offset)
	requireLen(t, got.Query.Filters, 1)
}

func Test_SetField_CommaSeparatedSlice(t *testing.T) {
	type withDefaultSep struct {
		Tags []string `arg:"tags"`
	}
	type withCustomSep struct {
		Tags []string `arg:"tags" sep:"|"`
	}
	type withInts struct {
		Ports []int `arg:"ports"`
	}

	tests := []struct {
		name   string
		target any
		inputs map[string]any
		assert func(t *testing.T, target any)
	}{
		{
			name:   "default comma separator splits and trims",
			target: &withDefaultSep{},
			inputs: map[string]any{"tags": "a, b ,c"},
			assert: func(t *testing.T, target any) {
				t.Helper()
				requireEqual(t, []string{"a", "b", "c"}, target.(*withDefaultSep).Tags)
			},
		},
		{
			name:   "single value becomes one element",
			target: &withDefaultSep{},
			inputs: map[string]any{"tags": "solo"},
			assert: func(t *testing.T, target any) {
				t.Helper()
				requireEqual(t, []string{"solo"}, target.(*withDefaultSep).Tags)
			},
		},
		{
			name:   "empty string yields empty slice",
			target: &withDefaultSep{},
			inputs: map[string]any{"tags": ""},
			assert: func(t *testing.T, target any) {
				t.Helper()
				requireEqual(t, []string{}, target.(*withDefaultSep).Tags)
			},
		},
		{
			name:   "custom separator via sep tag",
			target: &withCustomSep{},
			inputs: map[string]any{"tags": "a|b|c"},
			assert: func(t *testing.T, target any) {
				t.Helper()
				requireEqual(t, []string{"a", "b", "c"}, target.(*withCustomSep).Tags)
			},
		},
		{
			name:   "already a slice passes through untouched",
			target: &withDefaultSep{},
			inputs: map[string]any{"tags": []string{"x,y", "z"}},
			assert: func(t *testing.T, target any) {
				t.Helper()
				requireEqual(t, []string{"x,y", "z"}, target.(*withDefaultSep).Tags)
			},
		},
		{
			name:   "int slice splits and converts",
			target: &withInts{},
			inputs: map[string]any{"ports": "8080,9090"},
			assert: func(t *testing.T, target any) {
				t.Helper()
				requireEqual(t, []int{8080, 9090}, target.(*withInts).Ports)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetStructFields(tt.target, Settings{TagOrder: DefaultTags, EncodingTags: DefaultEncodingTags}, tt.inputs)
			requireNoError(t, err)
			tt.assert(t, tt.target)
		})
	}
}

func Test_ParseTags_StripsOmitempty(t *testing.T) {
	tags := parseTags(`json:"filters,omitempty" yaml:"filters,omitempty,flow" rules:"required"`, DefaultEncodingTags)
	requireEqual(t, "filters", tags["json"])
	requireEqual(t, "filters", tags["yaml"])
	// non-stdlib tags without commas are unaffected
	requireEqual(t, "required", tags["rules"])
}
