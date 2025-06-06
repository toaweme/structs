package structs

import (
	"reflect"
	"strings"
)

func GetStructFields(structure any, parent *Field) ([]Field, error) {
	val := reflect.ValueOf(structure)
	if val.Kind() != reflect.Ptr {
		return nil, ErrInputPointer
	}
	if val.Elem().Kind() != reflect.Struct {
		return nil, ErrInputPointerStruct
	}

	fields := make([]Field, 0)

	val = val.Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		tags := parseTags(string(field.Tag))
		f := NewField(field.Name, field.Type.Kind(), fieldValue, tags, parent)

		if field.Type.Kind() == reflect.Struct {
			nestedFields, err := GetStructFields(val.Field(i).Addr().Interface(), &f)
			if err != nil {
				return nil, err
			}
			for j := range nestedFields {
				nestedFields[j].Parent = &f
				nestedFields[j].FQN = nestedFields[j].buildFQN()
			}
			f.Fields = nestedFields
			fields = append(fields, f)
		} else {
			fields = append(fields, f)
		}
	}
	return fields, nil
}

// parseTags extracts inputs and their inputs from a given line of text
// arg:"cwd" short:"c" help:"Current working directory"
// The function returns a map: {"arg": "cwd", "short": "c", "help": "Current working directory"}
func parseTags(line string) map[string]string {
	line = strings.TrimSpace(line)
	result := make(map[string]string)

	inTag := true
	lastTagName := ""

	inValue := false
	lastTagValue := ""

	for i, char := range line {
		switch {
		case char == ':':
			if inValue {
				lastTagValue += string(char)
				continue
			}
			inTag = false
			continue
		case char == '"':
			// allow escaping quotes
			if i > 0 && line[i-1] == '\\' {
				lastTagValue = lastTagValue[:len(lastTagValue)-1] + string(char)
				continue
			}
			// entering or exiting a tag value
			inValue = !inValue
			if inValue {
				continue
			} else {
				// exiting a tag value
				lastTagName = strings.TrimSpace(lastTagName)
				lastTagValue = strings.TrimSpace(lastTagValue)
				result[lastTagName] = lastTagValue
				lastTagName = ""
				lastTagValue = ""
				inTag = true
			}
		default:
			// Collect characters for the tag name or value
			if inTag {
				lastTagName += string(char)
			}
			if inValue {
				lastTagValue += string(char)
			}
		}
	}

	return result
}
