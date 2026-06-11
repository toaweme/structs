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

// Option configures a Struct. See WithTags, WithEncodingTags, WithRules, WithValidationMessageTag.
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

// WithRules sets the rule set used by Validate, keyed by the name used in a
// `rules:` tag. Defaults to DefaultRules; pass your own map to extend or replace it.
func WithRules(rules map[string]RuleFunc) Option {
	return func(s *Struct) { s.ruleFuncs = rules }
}

// WithValidationMessageTag sets the struct tag whose value is used as the field
// key in the map returned by Validate. When a field carries this tag, its value
// replaces the tag-priority-resolved name in the validation messages, letting you
// report errors under a stable, caller-facing name. Defaults to "rules".
func WithValidationMessageTag(tag string) Option {
	return func(s *Struct) { s.validationMessageTag = tag }
}

// DefaultTags is the default tag priority order for input lookup and validation.
var DefaultTags = []string{"arg", "short", "json", "yaml"}

// DefaultEncodingTags are the struct tags whose values follow the stdlib
// convention of a name followed by comma-separated options (e.g.
// `json:"name,omitempty"`). Only these tags have their comma-options stripped
// when parsed; freeform tags (help, default, rules, ...) keep their value
// verbatim. Override to change which tags are treated as encoding tags.
var DefaultEncodingTags = []string{"json", "yaml", "toml", "xml"}

// New binds a pointer to a struct into a reusable *Struct. structure must be a
// pointer to a struct. Options override the defaults: no tag priority,
// DefaultRules, DefaultEncodingTags, and the "rules" validation message tag.
func New(structure any, opts ...Option) *Struct {
	s := &Struct{
		structure:            structure,
		validationMessageTag: validationTag,
		ruleFuncs:            DefaultRules,
		encodingTags:         DefaultEncodingTags,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Validate runs the configured rules over inputs and returns the validation
// errors as a map of field name (resolved by tag priority, or the validation
// tag when present) to messages. An empty map means everything passed.
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

// Set populates the bound struct from inputs, resolving keys by tag priority
// and applying `default:` tag values to fields left zero.
func (m *Struct) Set(inputs map[string]any) error {
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
