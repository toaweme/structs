package structs

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/contentforward/structs/utils"
)

// this code might be confusing to read
// we use reflection on the struct AND on the values in order to set the fields correctly

var ErrInputPointerStruct = errors.New("input should be a pointer to a struct")

func SetStructFields(structure any, tagOrder []string, values map[string]any) error {
	val := reflect.ValueOf(structure)
	if val.Kind() != reflect.Ptr {
		return ErrInputPointerStruct
	}
	if val.Elem().Kind() != reflect.Struct {
		return ErrInputPointerStruct
	}
	//name := val.Elem().Type().Name()
	val = val.Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		fieldType := field.Type.Kind()

		// structs are a special case
		if fieldType == reflect.Struct {
			nestedValues := make(map[string]any)
			for key, value := range values {
				if strings.HasPrefix(key, field.Name+".") {
					nestedKey := strings.TrimPrefix(key, field.Name+".")
					nestedValues[nestedKey] = value
				}
			}
			err := SetStructFields(fieldValue.Addr().Interface(), tagOrder, nestedValues)
			if err != nil {
				return err
			}
			continue
		}

		// other types
		fieldTags := getFieldTags(field)
		for _, tag := range tagOrder {
			fieldNameInTag, ok := fieldTags[tag]
			if !ok {
				continue
			}

			// found value name in tag
			if value, ok := values[fieldNameInTag]; ok {
				//log.Trace().Str("tag", tag).
				//	Str("tag-value", fieldNameInTag).
				//	Str("struct-field", field.Name).
				//	Str("struct", name).
				//	Str("struct-field-type", fieldType.String()).
				//	Any("value", value).
				//	Msg("setting struct field value via tag")

				err := setValue(fieldNameInTag, value, fieldType, fieldValue)
				if err != nil {
					return err
				}
				break
			}

			// found default value name in tag
			if defaultValue, ok := fieldTags[defaultValueTag]; ok {
				//log.Trace().Str("tag", tag).
				//	Str("tag-field", fieldNameInTag).
				//	Str("struct-field", field.Name).
				//	Str("struct", name).
				//	Str("struct-field-type", fieldType.String()).
				//	Any("val", fieldValue).
				//	Msg("setting struct field value via 'default:*` tag")
				err := setValue(fieldNameInTag, defaultValue, fieldType, fieldValue)
				if err != nil {
					return err
				}

				break
			}
		}

		// if the field is not set, try to set it using the field name e.g. IsBeep vs is_beep
		if value, ok := values[field.Name]; ok {
			//log.Trace().Str("field", field.Name).Any("val", value).Msg("setting field directly")
			fieldValue.Set(reflect.ValueOf(value))
		}
	}

	return nil
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
		// log.Trace().Str("field", fieldName).Any("val", value).Msg("setting slice")
		err := setSliceValue(fieldName, value, fieldValue)
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
		// spew.Dump(value)
		fieldValue.Set(reflect.ValueOf(value))
	case reflect.Map:
		// log.Trace().Str("field", fieldName).Any("val", value).Msg("setting map")
		fieldValue.Set(reflect.ValueOf(value))
	default:
		return fmt.Errorf("unsupported field[%s] type: %s", fieldName, fieldType)
	}

	return nil
}

func setSliceValue(fieldName string, value any, fieldValue reflect.Value) error {
	// log.Trace().Str("field", fieldName).Any("val", value).Msg("setting slice")
	if fieldValue.Type().Elem().Kind() == reflect.String {
		// Handle the case where the field is a slice of strings
		if slice, ok := value.([]interface{}); ok {
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

func getFieldTags(field reflect.StructField) map[string]string {
	tags := make(map[string]string)
	splitTags := strings.Split(string(field.Tag), " ")

	for _, tag := range splitTags {
		// spew.Dump(tag)
		if strings.Contains(tag, ":") {
			pair := strings.Split(tag, ":")
			tagName := pair[0]
			tagValue := strings.Trim(strings.Join(pair[1:], ":"), "\"")

			// log.Trace().Str("tag", tagName).Str("val", tagValue).Msg("found tag")

			tags[tagName] = tagValue
		}
	}
	return tags
}
