package structs

import "reflect"

func GetStructFields(structure any) ([]Field, error) {
	val := reflect.ValueOf(structure)
	if val.Kind() != reflect.Ptr {
		return nil, ErrInputPointerStruct
	}
	if val.Elem().Kind() != reflect.Struct {
		return nil, ErrInputPointerStruct
	}

	fields := make([]Field, 0)

	val = val.Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.Type.Kind() == reflect.Struct {
			nestedFields, err := GetStructFields(val.Field(i).Addr().Interface())
			if err != nil {
				return nil, err
			}
			for _, nestedField := range nestedFields {
				f := NewField(field.Name+"."+nestedField.Name, nestedField.Tags)
				fields = append(fields, f)
			}
		} else {
			tags := getFieldTags(field)
			f := NewField(field.Name, tags)
			fields = append(fields, f)
		}
	}

	return fields, nil
}
