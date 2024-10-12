package structs

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/contentforward/structs/utils"
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

func SetStructFieldsWithEnv(structure any, settings Settings, values map[string]any, env map[string]any) error {
	all := addInputs(values, env)

	return SetStructFields(structure, settings, all)
}

// SetStructFields sets the fields of a struct based on the inputs provided
func SetStructFields(structure any, settings Settings, inputs map[string]any) error {
	// slog.Info("setting struct fields", "structure", structure, "settings", settings, "inputs", inputs)
	// val is a pointer to the struct, not the struct itself
	val := reflect.ValueOf(structure)
	if val.Kind() != reflect.Ptr {
		return ErrInputPointer
	}
	if val.Elem().Kind() != reflect.Struct {
		return ErrInputPointerStruct
	}

	// val.Elem() is the struct
	val = val.Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		fieldType := field.Type.Kind()
		fieldTags := parseTags(string(field.Tag))

		err := setField(field, fieldValue, fieldType, fieldTags, settings, inputs)
		if err != nil {
			return err
		}
	}

	return nil
}

func setField(field reflect.StructField, fieldValue reflect.Value, fieldType reflect.Kind, fieldTags map[string]string, settings Settings, inputs map[string]any) error {
	// slog.Info("-------------")
	// slog.Info("inputs", "inputs", inputs)
	// slog.Info("field", "name", field.Name, "type", fieldType.String(), "tags", fieldTags)
	// structs are a special case, and we don't support `default:"*"` for them
	// maybe it's a good TODO to support JSON strings in tag/default inputs
	if fieldType == reflect.Struct {
		// handle env var based ["PARENT_CHILD_FIELD"] = value`
		if envVarName, ok := fieldTags[envValueTag]; ok {
			// env file:
			// PARENT_CHILD_FIELD=value
			// struct field:
			// Parent struct {
			//   ChildField canBeStruct `env:"CHILD_FIELD"`
			// }
			// convert PARENT_CHILD_FIELD => CHILD_FIELD
			childInputs := getChildInputs(inputs, envVarName, "_")
			withChildInputs := addInputs(inputs, childInputs)
			err := SetStructFields(fieldValue.Addr().Interface(), settings, withChildInputs)
			if err != nil {
				return fmt.Errorf("failed to set field-struct.env[%s] value: %w", envVarName, err)
			}
			if !settings.AllowEnvOverride {
				return nil
			}
		}

		// handle tag based ["parent.child.*.child.field"] = value`
		for _, tag := range settings.TagOrder {
			fieldNameInTag, ok := fieldTags[tag]
			if !ok {
				continue
			}

			childInputs := getChildInputs(inputs, fieldNameInTag, ".")
			withChildInputs := addInputs(inputs, childInputs)
			err := SetStructFields(fieldValue.Addr().Interface(), settings, withChildInputs)
			if err != nil {
				return fmt.Errorf("failed to set field-struct.tag[%s] value: %w", fieldNameInTag, err)
			}

			// skips the rest of the tags
			if !settings.AllowTagOverride {
				return nil
			}
		}

		// handle field name based ["ParentField.ChildField.*.ChildField.Field"] = value`
		childInputs := getChildInputs(inputs, field.Name, ".")
		withChildInputs := addInputs(inputs, childInputs)
		err := SetStructFields(fieldValue.Addr().Interface(), settings, withChildInputs)
		if err != nil {
			return fmt.Errorf("failed to set field-struct.name[%s] value: %w", field.Name, err)
		}

		return nil
	}

	// first set the default if possible
	if defaultValue, ok := fieldTags[defaultValueTag]; ok {
		// slog.Info("setting default value", "field", field.Name, "value", defaultValue)
		err := setValue(field.Name, defaultValue, fieldType, fieldValue)
		if err != nil {
			return fmt.Errorf("failed to set default value for field[%s]: %w", field.Name, err)
		}

		// do nothing here to allow the default value to be overridden
	}

	// then set the env var if possible
	if envVarName, ok := fieldTags[envValueTag]; ok {
		if value, ok := inputs[envVarName]; ok {
			// slog.Info("setting env value", "field", field.Name, "env", envVarName, "value", value)
			err := setValue(envVarName, value, fieldType, fieldValue)
			if err != nil {
				return fmt.Errorf("failed to set field[%s].env[%s] value: %w", field.Name, envVarName, err)
			}
			if !settings.AllowEnvOverride {
				return nil
			}
		}
	}

	// then iterate over tags
	for _, tag := range settings.TagOrder {
		fieldNameInTag, ok := fieldTags[tag]
		if !ok {
			continue
		}

		// found field name in tag e.g. `arg:"cwd"`
		// it's noteworthy that we allow emptying the field when a key is present
		// this ensures we can keep the ability to toggle the default value
		if value, ok := inputs[fieldNameInTag]; ok {
			// slog.Info("setting tag value", "field", field.Name, "tag", tag, "value", value)
			err := setValue(fieldNameInTag, value, fieldType, fieldValue)
			if err != nil {
				return fmt.Errorf("failed to set field[%s] value: %w", field.Name, err)
			}
			// skips the rest of the tags
			if !settings.AllowTagOverride {
				return nil
			}
		}
	}

	// if the field is not set, try to set it using the field name e.g. IsBeep vs is_beep
	if value, ok := inputs[field.Name]; ok {
		// slog.Info("setting field value", "field", field.Name, "value", value)
		err := setValue(field.Name, value, fieldType, fieldValue)
		if err != nil {
			return fmt.Errorf("failed to set field[%s] value: %w", field.Name, err)
		}
	}

	return nil
}

func addInputs(values map[string]any, nestedValues map[string]any) map[string]any {
	for key, value := range nestedValues {
		values[key] = value
	}

	return values
}

func getChildInputs(values map[string]any, fieldName, sep string) map[string]any {
	nestedValues := make(map[string]any)

	for key, value := range values {
		if strings.HasPrefix(key, fieldName+sep) {
			nestedKey := strings.TrimPrefix(key, fieldName+sep)
			nestedValues[nestedKey] = value
		}
	}

	return nestedValues
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
		fieldValue.Set(reflect.ValueOf(value))
	case reflect.Map:
		fieldValue.Set(reflect.ValueOf(value))
	default:
		return fmt.Errorf("unsupported field[%s] type: %s", fieldName, fieldType)
	}

	return nil
}

func setSliceValue(value any, fieldValue reflect.Value) error {
	// log.Trace().Str("field", fieldName).Any("val", value).Msg("setting slice")
	if fieldValue.Type().Elem().Kind() == reflect.String {
		// Handle the case where the field is a slice of strings
		if slice, ok := value.([]any); ok {
			stringSlice := make([]string, len(slice))
			for i, val := range slice {
				if str, ok := val.(string); ok {
					stringSlice[i] = str
				} else {
					return fmt.Errorf("cannot convert %T to string", val)
				}
			}
			fieldValue.Set(reflect.ValueOf(stringSlice))
		} else if slice, ok := value.([]string); ok {
			fieldValue.Set(reflect.ValueOf(slice))
			return nil
		}

		return nil
	}

	slice, err := utils.ToAnySlice(value)
	if err != nil {
		return err
	}
	fieldValue.Set(reflect.ValueOf(slice))

	return nil
}
