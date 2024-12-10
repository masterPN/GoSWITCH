package helpers

import (
	"testing"
)

func TestNormalizeDestinationNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "string starting with '0'",
			input:    "01234",
			expected: "661234",
		},
		{
			name:     "string not starting with '0'",
			input:    "12345",
			expected: "12345",
		},
		{
			name:     "string with single character '0'",
			input:    "0",
			expected: "66",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := NormalizeDestinationNumber(test.input)
			if actual != test.expected {
				t.Errorf("NormalizeDestinationNumber(%q) = %q, want %q", test.input, actual, test.expected)
			}
		})
	}
}
