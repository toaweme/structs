package structs

import (
	"reflect"
	"strings"
)

type Rule struct {
	Name string
	Args []string
}

func parseRules(rules []string) []Rule {
	parsedRules := make([]Rule, 0)
	for _, rule := range rules {
		parts := strings.Split(rule, ":")
		r := Rule{Name: parts[0]}
		if len(parts) > 1 {
			r.Args = strings.Split(parts[1], ",")
		}
		parsedRules = append(parsedRules, r)
	}
	return parsedRules

}

const defaultValueTag = "default"
const envValueTag = "env"
const validationTag = "rules"

type Field struct {
	Name    string
	Type    string
	Tags    map[string]string
	Default string
	Kind    reflect.Kind
	Value   reflect.Value
	Rules   []Rule
	FQN     *Field
	Parent  *Field
	Fields  []Field
}

func NewField(name string, dataType reflect.Kind, value reflect.Value, tags map[string]string, parentField *Field) Field {
	f := Field{Name: name, Kind: dataType, Type: dataType.String()}
	if defaultVal, ok := tags[defaultValueTag]; ok {
		f.Default = defaultVal
		delete(tags, defaultValueTag)
	}
	if rules, ok := tags[validationTag]; ok {
		f.Rules = parseRules(strings.Split(rules, "|"))
		delete(tags, validationTag)
	}
	f.Value = value
	f.Tags = tags
	f.Parent = parentField
	return f
}

func (f Field) buildFQN() *Field {
	if f.Parent == nil {
		return nil
	}
	parent := f.Parent

	newField := &Field{
		Name: f.Name,
		// don't modify the original tags
		Tags: make(map[string]string),
	}
	for tag, value := range f.Tags {
		newField.Tags[tag] = value
	}
	// recursively build the FQN for Name and Tags by gluing the parent's data with "."
	// `env` tag should be glued with "_"
	for parent != nil {
		newField.Name = parent.Name + "." + newField.Name
		// slog.Info("newField", "name", newField.Name)
		for tag, value := range parent.Tags {
			if _, ok := newField.Tags[tag]; !ok {
				continue
			}
			if tag == envValueTag {
				newField.Tags[tag] = value + "_" + newField.Tags[tag]
				// slog.Info("env", "name", newField.Tags[tag])
			} else {
				newField.Tags[tag] = value + "." + newField.Tags[tag]
			}

		}
		parent = parent.Parent
	}

	return newField
}

func MapDefaultValues(fields []Field, values map[string]any, tagPriority ...string) map[string]any {
	valuesCopy := make(map[string]any)
	for k, v := range values {
		valuesCopy[k] = v
	}

	for _, tag := range tagPriority {
		tag = strings.ToLower(tag)

		for _, field := range fields {
			foundValue, ok := values[field.Name]
			if ok && foundValue != "" {
				continue
			}

			fieldNameByTag, ok := field.Tags[tag]
			if ok {
				if field.Default != "" {
					valuesCopy[fieldNameByTag] = field.Default
					// log.Trace().Str("field", field.Name).Str("tag", tag).Str("value", field.Default).Msg("setting default value")
				}
			}
		}
	}

	return valuesCopy
}
