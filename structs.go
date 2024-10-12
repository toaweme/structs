package structs

import (
	"fmt"
)

// Struct holds the structure to be validated and the rules to validate it with
type Struct struct {
	structure            any
	ruleFuncs            map[string]RuleFunc
	validationMessageTag string
	tags                 []string
}

var DefaultTags = []string{"arg", "short", "json", "yaml"}

func New(structure any, rules map[string]RuleFunc, tags ...string) *Struct {
	return &Struct{
		structure:            structure,
		validationMessageTag: validationTag,
		ruleFuncs:            rules,
		tags:                 tags,
	}
}

func NewWithValidation(structure any, rules map[string]RuleFunc, validationTag string, tags ...string) *Struct {
	return &Struct{
		structure:            structure,
		validationMessageTag: validationTag,
		ruleFuncs:            rules,
		tags:                 tags,
	}
}

func (m *Struct) Validate(inputs map[string]any) (map[string][]string, error) {
	structFields, err := GetStructFields(m.structure)
	if err != nil {
		return nil, fmt.Errorf("error getting struct fields for validation: %w", err)
	}

	errors, err := ValidateStructFields(m.ruleFuncs, structFields, inputs, m.validationMessageTag, m.tags...)
	if err != nil {
		return nil, fmt.Errorf("error validating struct with inputs: %w", err)
	}

	return errors, nil
}

func (m *Struct) Set(inputs map[string]any) error {
	err := SetStructFields(m.structure, Settings{
		TagOrder:         m.tags,
		AllowEnvOverride: false,
		AllowTagOverride: false,
	}, inputs)
	if err != nil {
		return fmt.Errorf("error setting struct fields: %w", err)
	}

	return nil
}
