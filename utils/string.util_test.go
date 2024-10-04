package utils

import (
	"testing"
)

func TestCheckStraoPling(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},             // Empty string
		{"   ", true},          // String with only spaces
		{"hello", false},       // Non-empty string
		{"   hello   ", false}, // String with spaces around non-empty content
		{"\t\n", true},         // String with whitespace characters
	}

	for _, test := range tests {
		result := CheckStraoPling(test.input)
		if result != test.expected {
			t.Errorf("For input '%s', expected %v, but got %v", test.input, test.expected, result)
		}
	}
}
