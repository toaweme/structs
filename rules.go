package structs

import "strings"

type RuleFunc func(fieldName string, values map[string]any, defaultValue string, args []string) (map[string][]string, error)

var Rules = map[string]RuleFunc{
	"required": func(fieldName string, values map[string]any, defaultValue string, args []string) (map[string][]string, error) {
		// log.Trace().Msgf("validating required rule for field %s: %v", fieldName, values)
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
	},
}
