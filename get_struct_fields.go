package structs

import (
	"reflect"
	"strings"
)

// GetStructFields reflects over structure (a pointer to a struct) and returns
// its fields as []Field, recursing into named nested structs and building each
// nested field's FQN. Embedded (anonymous) struct fields are promoted: their
// fields are returned inline at this level, with no wrapper field and no FQN,
// matching Go's own field promotion. encodingTags selects which tags get their
// comma options stripped (see DefaultEncodingTags). It returns ErrInputPointer
// or ErrInputPointerStruct when structure is not a pointer to a struct.
func GetStructFields(structure any, parent *Field, encodingTags []string) ([]Field, error) {
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

		tags := parseTags(string(field.Tag), encodingTags)

		// an untagged embedded (anonymous) struct promotes its fields to this
		// level, just as Go's own field promotion does: the fields appear inline,
		// with no wrapper field and no FQN prefix. a tagged embed (or any named
		// struct field) instead groups its fields under a dotted FQN, matching
		// encoding/json: a tag on an anonymous field names it rather than
		// promoting it.
		if field.Anonymous && field.Type.Kind() == reflect.Struct && len(tags) == 0 {
			promoted, err := GetStructFields(val.Field(i).Addr().Interface(), parent, encodingTags)
			if err != nil {
				return nil, err
			}
			fields = append(fields, promoted...)
			continue
		}

		f := NewField(field.Name, field.Type.Kind(), fieldValue, tags, parent)

		if field.Type.Kind() == reflect.Struct {
			nestedFields, err := GetStructFields(val.Field(i).Addr().Interface(), &f, encodingTags)
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
func parseTags(line string, encodingTags []string) map[string]string {
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
				// comma-suffixed options (",omitempty", ",flow", ...) are a
				// convention of encoding tags only. strip them for those tags so
				// the stored value is just the name (e.g. "filters" not
				// "filters,omitempty"); freeform tags (help, default, rules, ...)
				// keep their value verbatim.
				if isEncodingTag(lastTagName, encodingTags) {
					if idx := strings.IndexByte(lastTagValue, ','); idx >= 0 {
						lastTagValue = lastTagValue[:idx]
					}
				}
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

// isEncodingTag reports whether name is one of encodingTags, the set of tags
// whose values use comma-separated options. See DefaultEncodingTags.
func isEncodingTag(name string, encodingTags []string) bool {
	for _, t := range encodingTags {
		if t == name {
			return true
		}
	}
	return false
}
