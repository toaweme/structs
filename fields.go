package structs

import (
	"reflect"
	"strings"
)

// Rule is a single validation rule parsed from a `rules:` tag entry. For
// `rules:"oneof:json,yaml"` Name is "oneof" and Args is ["json", "yaml"].
type Rule struct {
	// Name is the rule identifier, used to look up its RuleFunc in the rule set.
	Name string
	// Args are the colon-suffixed, comma-separated arguments, nil if none.
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

// Field is the reflected description of one struct field, produced by
// GetStructFields. Nested struct fields are described recursively through
// Fields, and each nested field also carries an FQN giving its dotted path and
// glued tags relative to the root struct.
type Field struct {
	// Name is the Go struct field name (e.g. "Field1").
	Name string
	// Type is the string form of Kind (e.g. "string", "struct").
	Type string
	// Tags holds the parsed struct tags as tag name to value (e.g. {"json": "field_1"}).
	Tags map[string]string
	// Default is the `default:` tag value, kept as a string because it comes
	// from the tag; it is converted to the field's type when the field is set.
	Default string
	// Kind is the field's reflect.Kind.
	Kind reflect.Kind
	// Value is the addressable reflect.Value backing the field, used to set it.
	Value reflect.Value
	// Rules are the validation rules parsed from the `rules:` tag.
	Rules []Rule
	// FQN is the fully-qualified view of a nested field: Name and Tags glued to
	// the parent's with "." (or "_" for the `env` tag). nil for top-level fields.
	FQN *Field
	// Parent points to the enclosing struct's field, nil for top-level fields.
	Parent *Field
	// Fields are the nested fields when Kind is reflect.Struct.
	Fields []Field
}

// NewField builds a Field from a struct field's name, kind, value, and parsed
// tags, extracting the `default:` and `rules:` tags into Default and Rules.
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
		for tag, value := range parent.Tags {
			if _, ok := newField.Tags[tag]; !ok {
				continue
			}
			if tag == envValueTag {
				newField.Tags[tag] = value + "_" + newField.Tags[tag]
			} else {
				newField.Tags[tag] = value + "." + newField.Tags[tag]
			}
		}
		parent = parent.Parent
	}

	return newField
}

// MapDefaultValues returns a copy of values with each field's `default:` tag
// value filled in under the field's tag name, for fields that have no non-empty
// value yet. Defaults are keyed by the first matching tag in tagPriority.
func MapDefaultValues(fields []Field, values map[string]any, tagPriority ...string) map[string]any {
	valuesCopy := make(map[string]any)
	for k, v := range values {
		valuesCopy[k] = v
	}

	for _, tag := range tagPriority {
		tag = strings.ToLower(tag)

		for _, field := range fields {
			fieldNameByTag, ok := field.Tags[tag]
			if !ok || field.Default == "" {
				continue
			}

			// don't override a value SetField could match, by tag name or field name
			if v, ok := values[fieldNameByTag]; ok && v != "" {
				continue
			}
			if v, ok := values[field.Name]; ok && v != "" {
				continue
			}

			valuesCopy[fieldNameByTag] = field.Default
		}
	}

	return valuesCopy
}
