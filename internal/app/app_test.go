package app

import (
	"testing"
)

func TestGenerateExportCustom(t *testing.T) {
	testCases := []struct {
		name        string
		format      string
		matches     [][]string
		expected    string
		expectError bool
	}{
		{
			name:   "Simple replacement",
			format: "$1-$2",
			matches: [][]string{
				{"a b", "a", "b"},
				{"c d", "c", "d"},
			},
			expected:    "a-b\nc-d",
			expectError: false,
		},
		{
			name:   "User case with newline and tab",
			format: `$2:$1\n\t$1`,
			matches: [][]string{
				{"- a b", "a", "b"},
			},
			expected:    "b:a\n\ta",
			expectError: false,
		},
		{
			name:   "Multiple matches with newlines",
			format: `$1\n`,
			matches: [][]string{
				{"a", "a"},
				{"b", "b"},
			},
			// The function adds its own newline between matches, so we expect a blank line
			expected:    "a\n\nb\n",
			expectError: false,
		},
		{
			name:   "Empty format string",
			format: "",
			matches: [][]string{
				{"a b", "a", "b"},
			},
			expected:    "",
			expectError: true,
		},
		{
			name:   "Reversed order",
			format: "$2 $1 $3",
			matches: [][]string{
				{"a b", "a", "b"},
				{"c d", "c", "d"},
			},
			expected:    "b a $3\nd c $3",
			expectError: false,
		},
		{
			name:   "nested $ signs",
			format: "$1-$2",
			matches: [][]string{
				{"a b", "a$2", "b"},
			},
			expected:    "a$2-b",
			expectError: false,
		},
		{
			name:        "No matches",
			format:      "$1",
			matches:     [][]string{},
			expected:    "",
			expectError: false,
		},
		{
			name:   "Format string with percentage sign",
			format: "$1 is 100%",
			matches: [][]string{
				{"a b", "a"},
			},
			expected:    "a is 100%",
			expectError: false,
		},
		{
			name:   "Format string with fmt verb",
			format: "$1 gives %s",
			matches: [][]string{
				{"a b", "a"},
			},
			expected:    "a gives %s",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a minimal App instance to call the method
			a := &App{
				matches: tc.matches,
			}

			result, err := a.generateExportCustom(tc.format)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tc.expected {
					t.Errorf("Expected result:\n---\n%s\n---\nGot:\n---\n%s\n---", tc.expected, result)
				}
			}
		})
	}
}
