package structs

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/toaweme/structs/utils"
)

var ErrInputPointer = errors.New("structure should be a pointer")
var ErrInputPointerStruct = errors.New("structure should be a pointer to a struct")

// TODO: file flags
// TODO: Fosdem stickers

type Settings struct {
	TagOrder []string
	// AllowEnvOverride toggles whether env vars can be overridden by tag inputs
	// env var handling takes priority over tag handling and is enabled by having a tag `env:"ENV_VAR"`
	// this toggles whether `env:"ENV_VAR"` can be overridden by `anything:"env_var"`
	AllowEnvOverride bool

	// AllowTagOverride toggles whether tag inputs can be overridden by other tags or exact FieldName matches
	// if true, only the first tag that matches will be used and if nothing is matched
	// we'll look for the field name as the structure key
	AllowTagOverride bool
}

func printFields(fields []Field) {
	for _, field := range fields {
		fmt.Println(field.Name, field.Kind, "fields", len(field.Fields))
		printFields(field.Fields)
	}
}

// SetStructFields sets the fields of a struct based on the inputs provided
func SetStructFields(structure any, settings Settings, inputs map[string]any) error {
	// log.Info("SetStructFields", "structure", reflect.TypeOf(structure), "settings", settings)
	fields, err := GetStructFields(structure, nil)
	if err != nil {
		return err
	}

	// for _, field := range fields {
	// 	log.Info("GetStructFields", "field", field.Name, "type", field.Kind)
	// }

	// print recursive fields with nested field count
	// printFields(fields)

	err = SetFields(fields, settings, inputs)
	if err != nil {
		return err
	}

	return nil
}

func SetFields(fields []Field, settings Settings, inputs map[string]any) error {
	// log.Info("SetFields", "fields", len(fields))
	for _, field := range fields {
		if field.Kind == reflect.Struct {
			// log.Info("SetFields.struct", "field", field.Name, "type", field.Kind, "tags", field.Tags)
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

func SetField(field Field, settings Settings, inputs map[string]any) error {
	// log.Info("SetField", "field", field.Name, "type", field.Kind, "tags", field.Tags)
	// set default value if it exists
	if field.Default != "" {
		// check if field has already a value set
		if !field.Value.IsValid() || field.Value.IsZero() {
			// log.Info("SetField.default", "field", field.Name, "value", field.Default)
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
			} else {
				// log.Info("SetField.tag", "field", field.Name, "tag", tag, "not found in inputs")
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
		// log.Info("SetField.tag1", "field", field.Name, "tag", tag, "fqnTag", fieldTag)
		if val, ok := inputs[fieldTag]; ok {
			// log.Info("SetField.tag2", "field", field.Name, "tag", tag, "value", val)
			err := setField(field, val)
			if err != nil {
				return err
			}

			if !settings.AllowTagOverride {
				return nil
			}
		} else {
			// log.Info("SetField.tag", "field", field.Name, "tag", tag, "not found in inputs", "fqnTag", fieldTag)
			split := strings.Split(fieldTag, ".")
			if len(split) == 1 {
				continue
			}

			found, value := findNestedValue(inputs, split)
			if found {
				// log.Info("SetField.nested", "field", field.Name, "tag", tag, "value", value)
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

	// log.Warn("SetField", "stage", "passed fqn tag checks", "field", field.Name)

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

	// Navigate through all but the last key
	for _, key := range path[:len(path)-1] {
		value, exists := current[key]
		if !exists {
			return false, nil
		}

		// Check if the value is a map[string]any
		if nestedMap, ok := value.(map[string]any); ok {
			current = nestedMap
		} else {
			// Value exists but isn't a nested map, can't continue traversal
			return false, nil
		}
	}

	// Check for the final key
	finalKey := path[len(path)-1]
	value, exists := current[finalKey]
	return exists, value
}

func setField(field Field, input any) error {
	err := setValue(field.Name, input, field.Kind, field.Value)
	if err != nil {
		return fmt.Errorf("failed to set field[%s]: %w", field.Name, err)
	}

	return nil
}

func setValue(fieldName string, value any, fieldType reflect.Kind, fieldValue reflect.Value) error {
	// log.Info("setValue", "field", fieldName, "value", value, "type", fieldType)
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
		elemVal := reflect.ValueOf(val)

		if !elemVal.Type().AssignableTo(elemType) {
			if elemVal.Type().ConvertibleTo(elemType) {
				elemVal = elemVal.Convert(elemType)
			} else {
				return fmt.Errorf("cannot assign or convert %T to %s", val, elemType)
			}
		}

		newSlice.Index(i).Set(elemVal)
	}

	fieldValue.Set(newSlice)

	return nil
}
