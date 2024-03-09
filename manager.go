package structs

import (
	"fmt"
)

type Manager struct {
	structure any
	tags      []string
}

var DefaultTags = []string{"arg", "short", "env", "json", "yaml"}

func NewManager(structure any, tags ...string) *Manager {
	return &Manager{
		structure: structure,
		tags:      tags,
	}
}

func (m *Manager) Validate(inputs map[string]any) (map[string][]string, error) {
	structFields, err := GetStructFields(m.structure)
	if err != nil {
		return nil, fmt.Errorf("error getting struct fields for validation: %w", err)
	}
	errors, err := ValidateStructFields(structFields, inputs, "json", m.tags...)
	if err != nil {
		return nil, fmt.Errorf("error validating translator inputs: %w", err)
	}

	return errors, nil
}

func (m *Manager) SetFields(inputs map[string]any) error {
	err := SetStructFields(m.structure, m.tags, inputs)
	if err != nil {
		return err
	}

	return nil
}
