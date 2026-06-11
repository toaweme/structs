package structs

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/toaweme/structs/utils"
)

// ErrInputPointer is returned when the structure passed in is not a pointer.
var ErrInputPointer = errors.New("structure should be a pointer")

// ErrInputPointerStruct is returned when the structure passed in is a pointer
// but does not point to a struct.
var ErrInputPointerStruct = errors.New("structure should be a pointer to a struct")

// Settings controls how SetStructFields resolves inputs onto struct fields.
type Settings struct {
	// TagOrder is the tag priority used to match input keys to fields
	// the first tag in the list that a field carries wins.
	TagOrder []string
	// AllowEnvOverride toggles whether env vars can be overridden by tag inputs
	// env var handling takes priority over tag handling and is enabled by having a tag `env:"ENV_VAR"`
	// this toggles whether `env:"ENV_VAR"` can be overridden by `anything:"env_var"`
	AllowEnvOverride bool
	// AllowTagOverride toggles whether a later matching tag in TagOrder can
	// override a value already set by an earlier one.
	// if false, the first tag in TagOrder that matches wins and matching stops there.
	// if true, every matching tag is applied in order, so the last match wins.
	AllowTagOverride bool
	// EncodingTags are the tags whose values use comma-separated options (see
	// DefaultEncodingTags). Empty disables comma stripping.
	EncodingTags []string
}

// SetStructFields sets the fields of a struct based on the inputs provided
func SetStructFields(structure any, settings Settings, inputs map[string]any) error {
	fields, err := GetStructFields(structure, nil, settings.EncodingTags)
	if err != nil {
		return err
	}

	err = SetFields(fields, settings, inputs)
	if err != nil {
		return err
	}

	return nil
}

// SetFields sets each field in fields from inputs, recursing into nested
// structs. It is the recursive worker behind SetStructFields.
func SetFields(fields []Field, settings Settings, inputs map[string]any) error {
	for _, field := range fields {
		if field.Kind == reflect.Struct {
			err := SetFields(field.Fields, settings, inputs)
			if err != nil {
				return err
			}
			continue
		}

		err := SetField(field, settings, inputs)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetField sets a single field from inputs: it applies the field's default
// first, then looks up a value by env tag, exact field name, and tag priority
// (using the field's FQN for nested fields), honoring the override settings.
func SetField(field Field, settings Settings, inputs map[string]any) error {
	// set default value if it exists
	if field.Default != "" {
		// check if field has already a value set
		if !field.Value.IsValid() || field.Value.IsZero() {
			err := setField(field, field.Default)
			if err != nil {
				return fmt.Errorf("failed to set default value for field[%s]: %w", field.Name, err)
			}
		}
	}

	fqn := field.FQN
	// may be a top level field
	if fqn == nil {
		// check env var matches
		if _, ok := field.Tags[envValueTag]; ok {
			envKey := field.Tags[envValueTag]
			if _, ok := inputs[envKey]; ok {
				err := setField(field, inputs[envKey])
				if err != nil {
					return err
				}

				if !settings.AllowEnvOverride {
					return nil
				}
			}
		}

		// check exact field name match
		if val, ok := inputs[field.Name]; ok {
			err := setField(field, val)
			if err != nil {
				return err
			}
		}

		// check tag matches
		for _, tag := range settings.TagOrder {
			if val, ok := inputs[field.Tags[tag]]; ok {
				err := setField(field, val)
				if err != nil {
					return err
				}

				if !settings.AllowTagOverride {
					return nil
				}
			}
		}

		// check nested field matches
		if field.Fields != nil {
			err := SetFields(field.Fields, settings, inputs)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// check fqn env var matches
	if _, ok := fqn.Tags[envValueTag]; ok {
		envKey := fqn.Tags[envValueTag]
		if _, ok := inputs[envKey]; ok {
			err := setField(field, inputs[envKey])
			if err != nil {
				return err
			}

			if !settings.AllowEnvOverride {
				return nil
			}
		}
	}

	// check fqn exact field name match
	if val, ok := inputs[fqn.Name]; ok {
		err := setField(field, val)
		if err != nil {
			return err
		}
	}

	// check fqn tag matches
	for _, tag := range settings.TagOrder {
		fieldTag := fqn.Tags[tag]
		if val, ok := inputs[fieldTag]; ok {
			err := setField(field, val)
			if err != nil {
				return err
			}

			if !settings.AllowTagOverride {
				return nil
			}
		} else {
			split := strings.Split(fieldTag, ".")
			if len(split) == 1 {
				continue
			}

			found, value := findNestedValue(inputs, split)
			if found {
				err := setField(field, value)
				if err != nil {
					return err
				}

				if !settings.AllowTagOverride {
					return nil
				}
			}
		}
	}

	// check fqn nested field matches
	if field.Fields != nil {
		err := SetFields(field.Fields, settings, inputs)
		if err != nil {
			return err
		}
	}

	return nil
}

func findNestedValue(inputs map[string]any, path []string) (bool, any) {
	current := inputs

	// navigate through all but the last key
	for _, key := range path[:len(path)-1] {
		value, exists := current[key]
		if !exists {
			return false, nil
		}

		// can only descend into a nested map[string]any
		nestedMap, ok := value.(map[string]any)
		if !ok {
			return false, nil
		}
		current = nestedMap
	}

	// the final key holds the value
	finalKey := path[len(path)-1]
	value, exists := current[finalKey]
	return exists, value
}

const separatorTag = "sep"
const defaultSeparator = ","

func setField(field Field, input any) error {
	if field.Kind == reflect.Slice {
		input = splitSliceInput(field, input)
	}

	err := setValue(field.Name, input, field.Kind, field.Value)
	if err != nil {
		return fmt.Errorf("failed to set field[%s]: %w", field.Name, err)
	}

	return nil
}

// splitSliceInput turns a single string into a slice of trimmed elements,
// splitting on the field's `sep` tag (default ","). It only applies when the
// input is a string and the slice element is a scalar, so already-structured
// inputs ([]any from decoded config, struct-element slices) pass through
// untouched. This is what lets `--tags=a,b,c` and `TAGS=a,b,c` become
// ["a","b","c"] instead of a single-element ["a,b,c"].
func splitSliceInput(field Field, input any) any {
	s, ok := input.(string)
	if !ok {
		return input
	}
	if !field.Value.IsValid() || field.Value.Kind() != reflect.Slice {
		return input
	}
	if field.Value.Type().Elem().Kind() == reflect.Struct {
		return input
	}

	if s == "" {
		return []string{}
	}

	sep := field.Tags[separatorTag]
	if sep == "" {
		sep = defaultSeparator
	}

	parts := strings.Split(s, sep)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	return parts
}

func setValue(fieldName string, value any, fieldType reflect.Kind, fieldValue reflect.Value) error {
	switch fieldType {
	case reflect.String:
		s, err := utils.ToString(value)
		if err != nil {
			return err
		}
		fieldValue.SetString(s)
	case reflect.Float32, reflect.Float64:
		float, err := utils.ToFloat(value)
		if err != nil {
			return err
		}
		fieldValue.SetFloat(float)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		integer, err := utils.ToInt(value)
		if err != nil {
			return err
		}
		fieldValue.Set(reflect.ValueOf(integer))
	case reflect.Slice:
		err := setSliceValue(value, fieldValue)
		if err != nil {
			return err
		}
	case reflect.Bool:
		switch v := value.(type) {
		case bool:
			fieldValue.SetBool(v)
		case string:
			fieldValue.SetBool(utils.ParseBool(v))
		}
	case reflect.Interface:
		if value != nil {
			fieldValue.Set(reflect.ValueOf(value))
		}
	case reflect.Map:
		fieldValue.Set(reflect.ValueOf(value))
	default:
		return fmt.Errorf("unsupported field[%s] type: %s", fieldName, fieldType)
	}

	return nil
}

func setSliceValue(value any, fieldValue reflect.Value) error {
	fieldType := fieldValue.Type()
	elemType := fieldType.Elem()

	slice, err := utils.ToAnySlice(value)
	if err != nil {
		return err
	}

	newSlice := reflect.MakeSlice(fieldType, len(slice), len(slice))

	for i, val := range slice {
		// a struct element fed a map: populate its fields from the map entries
		if elemType.Kind() == reflect.Struct {
			valReflect := reflect.ValueOf(val)
			if valReflect.Kind() == reflect.Map {
				newStruct := reflect.New(elemType).Elem()

				for j := range elemType.NumField() {
					field := elemType.Field(j)
					structFieldValue := newStruct.Field(j)

					if !structFieldValue.CanSet() {
						continue
					}

					// match the map key to the field name, case-insensitively
					for _, key := range valReflect.MapKeys() {
						keyStr := fmt.Sprintf("%v", key.Interface())
						if strings.EqualFold(keyStr, field.Name) {
							mapValue := valReflect.MapIndex(key).Interface()
							err := setValue(field.Name, mapValue, field.Type.Kind(), structFieldValue)
							if err != nil {
								return fmt.Errorf("failed to set field %s: %w", field.Name, err)
							}
							break
						}
					}
				}

				newSlice.Index(i).Set(newStruct)
				continue
			}
		}

		if val == nil {
			newSlice.Index(i).Set(reflect.Zero(elemType))
			continue
		}

		// a string element targeting a scalar slice (e.g. []int from "8080,9090")
		// is coerced through the same converters used for top-level fields, since
		// reflect cannot convert "8080" to int directly.
		if s, ok := val.(string); ok && elemType.Kind() != reflect.String {
			elem := reflect.New(elemType).Elem()
			if err := setValue("", s, elemType.Kind(), elem); err != nil {
				return fmt.Errorf("failed to convert %q to %s: %w", s, elemType, err)
			}
			newSlice.Index(i).Set(elem)
			continue
		}

		elemVal := reflect.ValueOf(val)
		if !elemVal.Type().AssignableTo(elemType) {
			if !elemVal.Type().ConvertibleTo(elemType) {
				return fmt.Errorf("cannot assign or convert %T to %s", val, elemType)
			}
			elemVal = elemVal.Convert(elemType)
		}

		newSlice.Index(i).Set(elemVal)
	}

	fieldValue.Set(newSlice)

	return nil
}
