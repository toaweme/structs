package structs

import (
	"fmt"
	"reflect"
	"strings"
)

// DefaultRules is the built-in rule set, keyed by the name used in a `rules:`
// tag. Pass it to New/ValidateStructFields, or build your own map (optionally
// extending this one) to register custom RuleFuncs.
var DefaultRules = map[string]RuleFunc{
	"required": Required,
	"oneof":    OneOf,
}

// RuleFunc validates one field against one rule. It receives the lookup key
// (fieldName, already resolved by tag priority), the full input values, the
// field's default, its reflect.Value, and the rule's args. It returns a map of
// field name to error messages (empty when valid); the returned error is for
// internal failures, not validation failures.
type RuleFunc func(fieldName string, values map[string]any, defaultValue string, fieldValue reflect.Value, args []string) (map[string][]string, error)

// Required fails when a field has no non-empty input and no default and the
// field's current value is zero. Whitespace-only string inputs count as empty.
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
