package structs

import (
	"fmt"
	"reflect"
	"strings"
)

var DefaultRules = map[string]RuleFunc{
	"required": Required,
	"oneof":    OneOf,
}

type RuleFunc func(fieldName string, values map[string]any, defaultValue string, fieldValue reflect.Value, args []string) (map[string][]string, error)

var Required RuleFunc = func(fieldName string, values map[string]any, defaultValue string, fieldValue reflect.Value, args []string) (map[string][]string, error) {
	value, ok := values[fieldName]
	if !ok && defaultValue == "" {
		if fieldValue.IsValid() && !fieldValue.IsZero() {
			return nil, nil
		}
		errors := map[string][]string{
			fieldName: {"required"},
		}
		return errors, nil
	}

	if valStr, ok := value.(string); ok {
		value = strings.TrimSpace(valStr)
		if value == "" {
			value = defaultValue
		}
		if value == "" {
			errors := map[string][]string{
				fieldName: {"required"},
			}
			return errors, nil
		}
	}

	return nil, nil
}

// OneOf restricts a field to one of the values listed in the rule args, e.g.
// `rules:"oneof:json,yaml,toml"`. An empty/absent value passes (pair with
// `required` to force presence); the default value, when set, is what gets
// checked if no input was provided.
var OneOf RuleFunc = func(fieldName string, values map[string]any, defaultValue string, fieldValue reflect.Value, args []string) (map[string][]string, error) {
	value := ""
	if raw, ok := values[fieldName]; ok {
		if s, isStr := raw.(string); isStr {
			value = strings.TrimSpace(s)
		} else if raw != nil {
			value = fmt.Sprintf("%v", raw)
		}
	}
	if value == "" {
		value = defaultValue
	}
	if value == "" {
		return nil, nil
	}

	for _, allowed := range args {
		if value == allowed {
			return nil, nil
		}
	}

	return map[string][]string{
		fieldName: {fmt.Sprintf("must be one of: %s", strings.Join(args, ", "))},
	}, nil
}
