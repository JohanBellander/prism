package render

import (
	"testing"
)

func TestParseGridColumns(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "repeat syntax with 4 columns",
			input:    "repeat(4, 1fr)",
			expected: 4,
		},
		{
			name:     "repeat syntax with 3 columns",
			input:    "repeat(3, 1fr)",
			expected: 3,
		},
		{
			name:     "repeat syntax with 6 columns",
			input:    "repeat(6, 1fr)",
			expected: 6,
		},
		{
			name:     "space-separated 1fr values",
			input:    "1fr 1fr 1fr 1fr",
			expected: 4,
		},
		{
			name:     "space-separated mixed values",
			input:    "200px 1fr 1fr",
			expected: 3,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "invalid repeat syntax",
			input:    "repeat(abc, 1fr)",
			expected: 0,
		},
	}

	engine := NewLayoutEngine(1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.parseGridColumns(tt.input)
			if result != tt.expected {
				t.Errorf("parseGridColumns(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseGridColumns_EdgeCases(t *testing.T) {
	engine := NewLayoutEngine(1)

	// Test with extra whitespace
	result := engine.parseGridColumns("repeat( 4 , 1fr )")
	if result != 4 {
		t.Errorf("parseGridColumns with whitespace failed: got %d, expected 4", result)
	}

	// Test single column
	result = engine.parseGridColumns("1fr")
	if result != 1 {
		t.Errorf("parseGridColumns single column failed: got %d, expected 1", result)
	}

	// Test many columns
	result = engine.parseGridColumns("1fr 1fr 1fr 1fr 1fr 1fr 1fr 1fr")
	if result != 8 {
		t.Errorf("parseGridColumns 8 columns failed: got %d, expected 8", result)
	}
}
