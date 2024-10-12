package structs

import (
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
	Rules   []Rule
}

func NewField(name, dataType string, tags map[string]string) Field {
	f := Field{Name: name, Type: dataType}
	if defaultVal, ok := tags[defaultValueTag]; ok {
		f.Default = defaultVal
		delete(tags, defaultValueTag)
	}
	if rules, ok := tags[validationTag]; ok {
		f.Rules = parseRules(strings.Split(rules, "|"))
		delete(tags, validationTag)
	}
	f.Tags = tags
	return f
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
