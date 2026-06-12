package structs

import "testing"

// exercises the fluent public SDK (New + options + Validate/Set), which the
// lower-level GetStructFields/SetStructFields/ValidateStructFields tests reach
// only indirectly.

func Test_New_Options(t *testing.T) {
	s := New(&struct{}{},
		WithTags("arg", "json"),
		WithEncodingTags("json"),
		WithValidationTag("validate"),
	)

	requireEqual(t, []string{"arg", "json"}, s.tags)
	requireEqual(t, []string{"json"}, s.encodingTags)
	requireEqual(t, "validate", s.validationTag)
}

func Test_New_Defaults(t *testing.T) {
	s := New(&struct{}{})

	// no WithTags -> tag priority defaults to DefaultTags
	requireEqual(t, DefaultTags, s.tags)
	requireEqual(t, DefaultEncodingTags, s.encodingTags)
	requireEqual(t, "rules", s.validationTag)
}

func Test_Struct_Validate(t *testing.T) {
	type target struct {
		Format string `json:"format" rules:"oneof:json,yaml"`
		Mode   string `json:"mode" rules:"required"`
	}

	s := New(&target{}, WithTags("json"))

	t.Run("valid inputs report no errors", func(t *testing.T) {
		errs, err := s.Validate(map[string]any{"format": "json", "mode": "fast"})
		requireNoError(t, err)
		requireEqual(t, map[string][]string{}, errs)
	})

	t.Run("invalid inputs surface per-field errors", func(t *testing.T) {
		errs, err := s.Validate(map[string]any{"format": "xml"})
		requireNoError(t, err)
		requireEqual(t, map[string][]string{
			"format": {"must be one of: json, yaml"},
			"mode":   {"required"},
		}, errs)
	})

	t.Run("non-pointer structure errors", func(t *testing.T) {
		bad := New(target{}, WithTags("json"))
		_, err := bad.Validate(map[string]any{})
		requireErrorIs(t, err, ErrInputPointer)
	})
}

func Test_Struct_Set(t *testing.T) {
	t.Run("sets fields and applies defaults", func(t *testing.T) {
		type target struct {
			Name string `json:"name"`
			Port int    `json:"port" default:"8080"`
		}
		got := &target{}
		s := New(got, WithTags("json"))

		err := s.Set(map[string]any{"name": "svc"})
		requireNoError(t, err)
		requireEqual(t, "svc", got.Name)
		requireEqual(t, 8080, got.Port)
	})

	t.Run("non-pointer structure errors", func(t *testing.T) {
		s := New(struct{}{}, WithTags("json"))
		err := s.Set(map[string]any{})
		requireErrorIs(t, err, ErrInputPointer)
	})
}

func Test_MapDefaultValues(t *testing.T) {
	fields := []Field{
		{Name: "Field1", Tags: map[string]string{"json": "field_1"}, Default: "d1"},
		{Name: "Field2", Tags: map[string]string{"json": "field_2"}, Default: "d2"},
		{Name: "Field3", Tags: map[string]string{"json": "field_3"}},
	}

	t.Run("fills defaults for absent fields, keeps provided tag-keyed values", func(t *testing.T) {
		// field_2 is supplied under its tag key, so its default must not clobber
		// it; field_1 is absent and gets its default; Field3 has no default, so
		// it gets no entry.
		got := MapDefaultValues(fields, map[string]any{"field_2": "set"}, "json")
		requireEqual(t, map[string]any{
			"field_1": "d1",
			"field_2": "set",
		}, got)
	})

	t.Run("does not override field present by struct name", func(t *testing.T) {
		got := MapDefaultValues(fields, map[string]any{"Field1": "keep"}, "json")
		requireEqual(t, "keep", got["Field1"])
		if _, ok := got["field_1"]; ok {
			t.Fatalf("expected no default for field present by name, got %#v", got)
		}
	})

	t.Run("no tag priority yields a copy of the inputs", func(t *testing.T) {
		in := map[string]any{"x": 1}
		got := MapDefaultValues(fields, in)
		requireEqual(t, in, got)
		got["y"] = 2
		if _, ok := in["y"]; ok {
			t.Fatalf("MapDefaultValues mutated the input map")
		}
	})
}
