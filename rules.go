package structs

import "strings"

var DefaultRules = map[string]RuleFunc{
	"required": Required,
}

type RuleFunc func(fieldName string, values map[string]any, defaultValue string, args []string) (map[string][]string, error)

var Required RuleFunc = func(fieldName string, values map[string]any, defaultValue string, args []string) (map[string][]string, error) {
	value, ok := values[fieldName]
	if !ok && defaultValue == "" {
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
