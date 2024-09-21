package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:  "spaces in value",
			input: `arg:"cwd" short:"c" help:"Current working directory"`,
			expected: map[string]string{
				"arg":   "cwd",
				"short": "c",
				"help":  "Current working directory",
			},
		},
		{
			// TODO: unsupported quotes in value
			name:  "TODO: unsupported quotes in value",
			input: `arg:"cwd" short:"c" help:"Current \"working directory\""`,
			expected: map[string]string{
				"arg":   "cwd",
				"short": "c",
				"help":  `Current "working directory"`,
			},
		},
		{
			// TODO: unsupported quotes in value
			name:  "TODO: unsupported quotes in value",
			input: `arg:"cwd" short:"c" default:"http://127.0.0.1:3888" help:"Current \"working directory\""`,
			expected: map[string]string{
				"arg":     "cwd",
				"short":   "c",
				"default": "http://127.0.0.1:3888",
				"help":    `Current "working directory"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTags(tt.input)

			assert.Equal(t, tt.expected, result)
		})
	}
}
