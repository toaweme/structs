package structs

import (
	"fmt"
)

func ValidateStructFields(structFields []Field, values map[string]any, validationMessageTag string, tagPriority ...string) (map[string][]string, error) {
	validationErrors := make(map[string][]string)
	for _, structField := range structFields {
		for _, rule := range structField.Rules {
			fieldName := structField.Name
			tags := structField.Tags
			fieldNameByTagPriority := getTagByPriority(tags, tagPriority)
			if fieldNameByTagPriority != "" {
				fieldName = fieldNameByTagPriority
			}

			fieldValidationRules, err := validateRule(rule, fieldName, values, structField.Default)
			if err != nil {
				return nil, fmt.Errorf("error running validator function for rule '%s' field '%s': %w", rule.Name, fieldName, err)
			}

			for errorFieldName, errorMessages := range fieldValidationRules {
				if fieldNameForValidationMessage, ok := tags[validationMessageTag]; ok {
					errorFieldName = fieldNameForValidationMessage
				}
				if _, ok := validationErrors[errorFieldName]; !ok {
					validationErrors[errorFieldName] = make([]string, 0)
				}
				validationErrors[errorFieldName] = append(validationErrors[errorFieldName], errorMessages...)
			}
		}
	}

	return validationErrors, nil
}

func getTagByPriority(tags map[string]string, priority []string) string {
	for _, p := range priority {
		if tag, ok := tags[p]; ok {
			return tag
		}
	}

	return ""
}

func validateRule(rule Rule, structFieldName string, values map[string]any, defaultValue string) (map[string][]string, error) {
	validationFunc, ok := Rules[rule.Name]
	if !ok {
		return nil, fmt.Errorf("struct field %s rule %s not found", structFieldName, rule.Name)
	}

	errors, err := validationFunc(structFieldName, values, defaultValue, rule.Args)
	if err != nil {
		return nil, fmt.Errorf("validation error for field %s rule %s: %w", structFieldName, rule.Name, err)
	}

	if len(errors) == 0 {
		return nil, nil
	}

	// log.Trace().Msgf("validation errors for field %s rule %s: %v", structFieldName, rule.Name, errors)

	return errors, nil
}
