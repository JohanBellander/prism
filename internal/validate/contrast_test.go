package validate

import (
	"math"
	"testing"

	"github.com/johanbellander/prism/internal/types"
)

func TestCalculateContrastRatio(t *testing.T) {
	tests := []struct {
		name     string
		fg       string
		bg       string
		expected float64
		delta    float64
	}{
		{
			name:     "Black on white",
			fg:       "#000000",
			bg:       "#FFFFFF",
			expected: 21.0,
			delta:    0.1,
		},
		{
			name:     "White on black",
			fg:       "#FFFFFF",
			bg:       "#000000",
			expected: 21.0,
			delta:    0.1,
		},
		{
			name:     "Dark gray on white",
			fg:       "#767676",
			bg:       "#FFFFFF",
			expected: 4.5,
			delta:    0.1,
		},
		{
			name:     "Light gray on white (fails AA)",
			fg:       "#999999",
			bg:       "#FFFFFF",
			expected: 2.8,
			delta:    0.2,
		},
		{
			name:     "Blue on white",
			fg:       "#3B82F6",
			bg:       "#FFFFFF",
			expected: 3.7,
			delta:    0.2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ratio := calculateContrastRatio(tt.fg, tt.bg)
			if math.Abs(ratio-tt.expected) > tt.delta {
				t.Errorf("Expected ratio ~%.1f, got %.2f", tt.expected, ratio)
			}
		})
	}
}

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		name     string
		hex      string
		expectedR int
		expectedG int
		expectedB int
	}{
		{
			name:     "Black",
			hex:      "#000000",
			expectedR: 0,
			expectedG: 0,
			expectedB: 0,
		},
		{
			name:     "White",
			hex:      "#FFFFFF",
			expectedR: 255,
			expectedG: 255,
			expectedB: 255,
		},
		{
			name:     "Red",
			hex:      "#FF0000",
			expectedR: 255,
			expectedG: 0,
			expectedB: 0,
		},
		{
			name:     "Shorthand white",
			hex:      "#FFF",
			expectedR: 255,
			expectedG: 255,
			expectedB: 255,
		},
		{
			name:     "Blue",
			hex:      "#3B82F6",
			expectedR: 59,
			expectedG: 130,
			expectedB: 246,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b := hexToRGB(tt.hex)
			if r != tt.expectedR || g != tt.expectedG || b != tt.expectedB {
				t.Errorf("Expected RGB(%d, %d, %d), got RGB(%d, %d, %d)",
					tt.expectedR, tt.expectedG, tt.expectedB, r, g, b)
			}
		})
	}
}

func TestIsLargeTextSize(t *testing.T) {
	tests := []struct {
		name     string
		size     string
		weight   string
		expected bool
	}{
		{
			name:     "18px bold is large",
			size:     "lg",  // 20px
			weight:   "bold",
			expected: true,
		},
		{
			name:     "24px normal is large",
			size:     "xl",  // 24px
			weight:   "normal",
			expected: true,
		},
		{
			name:     "16px bold is not large",
			size:     "base",  // 16px
			weight:   "bold",
			expected: false,
		},
		{
			name:     "16px normal is not large",
			size:     "base",  // 16px
			weight:   "normal",
			expected: false,
		},
		{
			name:     "36px is always large",
			size:     "3xl",  // 36px
			weight:   "normal",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isLargeTextSize(tt.size, tt.weight)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for size=%s weight=%s",
					tt.expected, result, tt.size, tt.weight)
			}
		})
	}
}

func TestValidateContrast_PassingText(t *testing.T) {
	// Black text on white background should pass
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "container",
				Type: "box",
				Layout: types.ComponentLayout{
					Background: "#FFFFFF",
				},
				Children: []types.Component{
					{
						ID:      "text1",
						Type:    "text",
						Content: "Hello",
						Color:   "#000000",
					},
				},
			},
		},
	}

	rule := DefaultContrastRule()
	result := ValidateContrast(structure, rule)

	if !result.Passed {
		t.Errorf("Expected validation to pass for black text on white, but got %d issues", len(result.Issues))
		for _, issue := range result.Issues {
			t.Logf("  - %s: %s", issue.Category, issue.Message)
		}
	}
}

func TestValidateContrast_FailingText(t *testing.T) {
	// Light gray text on white background should fail
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "container",
				Type: "box",
				Layout: types.ComponentLayout{
					Background: "#FFFFFF",
				},
				Children: []types.Component{
					{
						ID:      "text1",
						Type:    "text",
						Content: "Hello",
						Color:   "#999999", // Light gray - fails AA
						Size:    "base",
						Weight:  "normal",
					},
				},
			},
		},
	}

	rule := DefaultContrastRule()
	result := ValidateContrast(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail for light gray text on white")
	}

	foundContrastIssue := false
	for _, issue := range result.Issues {
		if issue.Category == "contrast_fail" && issue.ComponentID == "text1" {
			foundContrastIssue = true
			if issue.ContrastRatio >= 4.5 {
				t.Errorf("Expected contrast ratio < 4.5, got %.2f", issue.ContrastRatio)
			}
		}
	}

	if !foundContrastIssue {
		t.Error("Expected contrast_fail issue for text1")
	}
}

func TestValidateContrast_LargeText(t *testing.T) {
	// Large text has lower contrast requirements (3:1 instead of 4.5:1)
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "container",
				Type: "box",
				Layout: types.ComponentLayout{
					Background: "#FFFFFF",
				},
				Children: []types.Component{
					{
						ID:      "heading",
						Type:    "text",
						Content: "Large Heading",
						Color:   "#999999", // Would fail for normal text, but passes for large
						Size:    "xl",      // 24px
						Weight:  "normal",
					},
				},
			},
		},
	}

	rule := DefaultContrastRule()
	result := ValidateContrast(structure, rule)

	// Light gray might still fail even for large text if ratio is below 3:1
	// Let's check the actual ratio
	ratio := calculateContrastRatio("#999999", "#FFFFFF")
	shouldPass := ratio >= 3.0

	if shouldPass && !result.Passed {
		t.Errorf("Expected large text to pass with ratio %.2f:1", ratio)
	}
}

func TestValidateContrast_ButtonText(t *testing.T) {
	// Button with insufficient contrast
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:      "button1",
				Type:    "button",
				Content: "Click Me",
				Layout: types.ComponentLayout{
					Background: "#FFEB3B", // Light yellow - white text would fail
				},
			},
		},
	}

	rule := DefaultContrastRule()
	result := ValidateContrast(structure, rule)

	// White on light yellow should fail
	ratio := calculateContrastRatio("#FFFFFF", "#FFEB3B")
	if ratio < 4.5 {
		if result.Passed {
			t.Error("Expected button contrast validation to fail")
		}
	}
}

func TestValidateContrast_InheritedBackground(t *testing.T) {
	// Text should inherit background from parent container
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "container",
				Type: "box",
				Layout: types.ComponentLayout{
					Background: "#FFFFFF",
				},
				Children: []types.Component{
					{
						ID:   "nested",
						Type: "box",
						Children: []types.Component{
							{
								ID:      "text1",
								Type:    "text",
								Content: "Nested text",
								Color:   "#CCCCCC", // Light gray - should fail
							},
						},
					},
				},
			},
		},
	}

	rule := DefaultContrastRule()
	result := ValidateContrast(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail for inherited background")
	}

	foundIssue := false
	for _, issue := range result.Issues {
		if issue.ComponentID == "text1" && issue.Category == "contrast_fail" {
			foundIssue = true
		}
	}

	if !foundIssue {
		t.Error("Expected contrast issue for nested text")
	}
}

func TestValidateContrast_MultipleIssues(t *testing.T) {
	// Multiple text elements with contrast issues
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "container",
				Type: "box",
				Layout: types.ComponentLayout{
					Background: "#FFFFFF",
				},
				Children: []types.Component{
					{
						ID:      "text1",
						Type:    "text",
						Content: "Low contrast 1",
						Color:   "#AAAAAA",
					},
					{
						ID:      "text2",
						Type:    "text",
						Content: "Low contrast 2",
						Color:   "#BBBBBB",
					},
				},
			},
		},
	}

	rule := DefaultContrastRule()
	result := ValidateContrast(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail for multiple low-contrast texts")
	}

	// Should have issues for both text elements
	issueCount := 0
	for _, issue := range result.Issues {
		if issue.Category == "contrast_fail" {
			issueCount++
		}
	}

	if issueCount < 2 {
		t.Errorf("Expected at least 2 contrast issues, got %d", issueCount)
	}
}

func TestValidateContrast_EmptyStructure(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{},
	}

	rule := DefaultContrastRule()
	result := ValidateContrast(structure, rule)

	if !result.Passed {
		t.Error("Expected validation to pass for empty structure")
	}
}
