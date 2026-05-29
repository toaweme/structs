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
	encodingTags         []string
}

// Option configures a Struct. See WithTags, WithEncodingTags, WithValidationTag.
type Option func(*Struct)

// WithTags sets the tag priority order used for input lookup and validation.
func WithTags(tags ...string) Option {
	return func(s *Struct) { s.tags = tags }
}

// WithEncodingTags sets the tags whose values use comma-separated options
// (see DefaultEncodingTags). Pass none to disable comma stripping entirely.
func WithEncodingTags(tags ...string) Option {
	return func(s *Struct) { s.encodingTags = tags }
}

// WithValidationTag sets the struct tag used to source validation rules.
func WithValidationTag(tag string) Option {
	return func(s *Struct) { s.validationMessageTag = tag }
}

var DefaultTags = []string{"arg", "short", "json", "yaml"}

// DefaultEncodingTags are the struct tags whose values follow the stdlib
// convention of a name followed by comma-separated options (e.g.
// `json:"name,omitempty"`). Only these tags have their comma-options stripped
// when parsed; freeform tags (help, default, rules, ...) keep their value
// verbatim. Override to change which tags are treated as encoding tags.
var DefaultEncodingTags = []string{"json", "yaml", "toml", "xml"}

func New(structure any, rules map[string]RuleFunc, opts ...Option) *Struct {
	s := &Struct{
		structure:            structure,
		validationMessageTag: validationTag,
		ruleFuncs:            rules,
		encodingTags:         DefaultEncodingTags,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (m *Struct) Validate(inputs map[string]any) (map[string][]string, error) {
	structFields, err := GetStructFields(m.structure, nil, m.encodingTags)
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
	// log.Info("Set", "structure", m.structure)
	err := SetStructFields(m.structure, Settings{
		TagOrder:         m.tags,
		AllowEnvOverride: false,
		AllowTagOverride: false,
		EncodingTags:     m.encodingTags,
	}, inputs)
	if err != nil {
		return fmt.Errorf("error setting struct fields: %w", err)
	}

	return nil
}
