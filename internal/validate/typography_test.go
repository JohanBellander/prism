package validate

import (
	"testing"

	"github.com/johanbellander/prism/internal/types"
)

func TestDefaultTypographyRule(t *testing.T) {
	rule := DefaultTypographyRule()
	
	if rule.ScaleRatio != 1.250 {
		t.Errorf("Expected scale ratio 1.250, got %.3f", rule.ScaleRatio)
	}
	
	if rule.BaseSize != 16.0 {
		t.Errorf("Expected base size 16, got %.1f", rule.BaseSize)
	}
	
	// Check expected sizes
	expectedSizes := map[string]float64{
		"xs":   12,
		"sm":   14,
		"base": 16,
		"md":   18,
		"lg":   20,
		"xl":   25,
		"2xl":  31,
		"3xl":  39,
		"4xl":  49,
	}
	
	for token, expectedSize := range expectedSizes {
		actualSize, exists := rule.Sizes[token]
		if !exists {
			t.Errorf("Missing size token '%s'", token)
			continue
		}
		if actualSize != expectedSize {
			t.Errorf("Token '%s': expected %.0f, got %.0f", token, expectedSize, actualSize)
		}
	}
}

func TestPredefinedScales(t *testing.T) {
	scales := PredefinedScales()
	
	expectedScales := map[string]float64{
		"minor-second":     1.067,
		"major-second":     1.125,
		"minor-third":      1.200,
		"major-third":      1.250,
		"perfect-fourth":   1.333,
		"augmented-fourth": 1.414,
		"perfect-fifth":    1.500,
		"golden-ratio":     1.618,
	}
	
	for name, expectedRatio := range expectedScales {
		actualRatio, exists := scales[name]
		if !exists {
			t.Errorf("Missing scale '%s'", name)
			continue
		}
		if actualRatio != expectedRatio {
			t.Errorf("Scale '%s': expected %.3f, got %.3f", name, expectedRatio, actualRatio)
		}
	}
}

func TestValidateTypography_ValidTokens(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "heading",
				Type: "text",
				Size: "2xl",
			},
			{
				ID:   "subheading",
				Type: "text",
				Size: "xl",
			},
			{
				ID:   "body",
				Type: "text",
				Size: "base",
			},
		},
	}
	
	rule := DefaultTypographyRule()
	result := ValidateTypography(structure, rule)
	
	if !result.Passed {
		t.Errorf("Expected validation to pass for valid tokens")
	}
	
	if len(result.Issues) > 0 {
		t.Errorf("Expected no issues, got %d", len(result.Issues))
	}
}

func TestValidateTypography_InvalidToken(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "heading",
				Type: "text",
				Size: "huge", // Invalid token
			},
		},
	}
	
	rule := DefaultTypographyRule()
	result := ValidateTypography(structure, rule)
	
	if result.Passed {
		t.Errorf("Expected validation to fail for invalid token")
	}
	
	if len(result.Issues) == 0 {
		t.Errorf("Expected issues to be reported")
	}
	
	// Check for warning about unknown token
	foundWarning := false
	foundInfo := false
	for _, issue := range result.Issues {
		if issue.Severity == "warning" && issue.ComponentID == "heading" {
			foundWarning = true
		}
		if issue.Severity == "info" {
			foundInfo = true
		}
	}
	
	if !foundWarning {
		t.Errorf("Expected warning about unknown token")
	}
	
	if !foundInfo {
		t.Errorf("Expected info about valid tokens")
	}
}

func TestValidateTypography_NestedComponents(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "header",
				Type: "box",
				Children: []types.Component{
					{
						ID:   "title",
						Type: "text",
						Size: "3xl",
					},
					{
						ID:   "subtitle",
						Type: "text",
						Size: "lg",
					},
				},
			},
			{
				ID:   "content",
				Type: "box",
				Children: []types.Component{
					{
						ID:   "paragraph",
						Type: "text",
						Size: "base",
					},
				},
			},
		},
	}
	
	rule := DefaultTypographyRule()
	result := ValidateTypography(structure, rule)
	
	if !result.Passed {
		t.Errorf("Expected validation to pass for nested valid tokens")
	}
	
	if len(result.Issues) > 0 {
		t.Errorf("Expected no issues, got %d", len(result.Issues))
	}
}

func TestValidateTypography_NonTextComponents(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "button",
				Type: "button",
				Size: "lg", // Size on non-text component should be ignored
			},
			{
				ID:   "container",
				Type: "box",
			},
			{
				ID:   "input",
				Type: "input",
			},
		},
	}
	
	rule := DefaultTypographyRule()
	result := ValidateTypography(structure, rule)
	
	if !result.Passed {
		t.Errorf("Expected validation to pass when no text components present")
	}
	
	if len(result.Issues) > 0 {
		t.Errorf("Expected no issues for non-text components, got %d", len(result.Issues))
	}
}

func TestValidateTypography_EmptySize(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "text-no-size",
				Type: "text",
				Size: "", // Empty size should be ignored
			},
		},
	}
	
	rule := DefaultTypographyRule()
	result := ValidateTypography(structure, rule)
	
	if !result.Passed {
		t.Errorf("Expected validation to pass when text has no size specified")
	}
	
	if len(result.Issues) > 0 {
		t.Errorf("Expected no issues for text without size, got %d", len(result.Issues))
	}
}

func TestValidateTypography_MultipleInvalidTokens(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "text1",
				Type: "text",
				Size: "invalid1",
			},
			{
				ID:   "text2",
				Type: "text",
				Size: "invalid2",
			},
			{
				ID:   "text3",
				Type: "text",
				Size: "base", // Valid
			},
		},
	}
	
	rule := DefaultTypographyRule()
	result := ValidateTypography(structure, rule)
	
	if result.Passed {
		t.Errorf("Expected validation to fail for multiple invalid tokens")
	}
	
	// Should have warnings for both invalid tokens + info messages
	if len(result.Issues) < 2 {
		t.Errorf("Expected at least 2 issues (warnings + infos), got %d", len(result.Issues))
	}
	
	// Count warnings
	warningCount := 0
	for _, issue := range result.Issues {
		if issue.Severity == "warning" {
			warningCount++
		}
	}
	
	if warningCount != 2 {
		t.Errorf("Expected 2 warnings, got %d", warningCount)
	}
}

func TestValidateTypography_AllStandardTokens(t *testing.T) {
	// Test all standard tokens from the default rule
	structure := &types.Structure{
		Components: []types.Component{
			{ID: "xs-text", Type: "text", Size: "xs"},
			{ID: "sm-text", Type: "text", Size: "sm"},
			{ID: "base-text", Type: "text", Size: "base"},
			{ID: "md-text", Type: "text", Size: "md"},
			{ID: "lg-text", Type: "text", Size: "lg"},
			{ID: "xl-text", Type: "text", Size: "xl"},
			{ID: "2xl-text", Type: "text", Size: "2xl"},
			{ID: "3xl-text", Type: "text", Size: "3xl"},
			{ID: "4xl-text", Type: "text", Size: "4xl"},
		},
	}
	
	rule := DefaultTypographyRule()
	result := ValidateTypography(structure, rule)
	
	if !result.Passed {
		t.Errorf("Expected validation to pass for all standard tokens")
	}
	
	if len(result.Issues) > 0 {
		t.Errorf("Expected no issues for standard tokens, got %d: %v", len(result.Issues), result.Issues)
	}
}

func TestIsOnTypographyScale(t *testing.T) {
	rule := DefaultTypographyRule()
	
	tests := []struct {
		size     float64
		expected bool
		name     string
	}{
		{16, true, "base size"},
		{20, true, "16 * 1.25"},
		{25, true, "16 * 1.25^2"},
		{31, true, "16 * 1.25^3 (with rounding)"},
		{12.8, true, "16 / 1.25 (with tolerance)"},
		{18, true, "16 * 1.25^0.5 (md size)"},
		{13, true, "close to 12.8, within tolerance"},
		{22, true, "close to scale value with half-steps"},
	}
	
	for _, tt := range tests {
		result := isOnTypographyScale(tt.size, rule)
		if result != tt.expected {
			t.Errorf("%s: size %.1f, expected %v, got %v", tt.name, tt.size, tt.expected, result)
		}
	}
}

func TestGetScaleName(t *testing.T) {
	tests := []struct {
		ratio    float64
		expected string
	}{
		{1.250, "major-third"},
		{1.125, "major-second"},
		{1.618, "golden-ratio"},
		{1.500, "perfect-fifth"},
		{1.999, "custom"}, // Unknown ratio
	}
	
	for _, tt := range tests {
		result := getScaleName(tt.ratio)
		if result != tt.expected {
			t.Errorf("Ratio %.3f: expected '%s', got '%s'", tt.ratio, tt.expected, result)
		}
	}
}
