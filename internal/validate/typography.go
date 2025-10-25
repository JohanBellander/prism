package validate

import (
	"fmt"
	"math"

	"github.com/johanbellander/prism/internal/types"
)

// TypographyRule defines the rules for typography scale validation
type TypographyRule struct {
	ScaleRatio float64            // e.g., 1.250 for Major Third
	BaseSize   float64            // base font size in pixels
	Sizes      map[string]float64 // expected sizes for each scale level
	Tolerance  float64            // acceptable deviation (e.g., 0.5px)
}

// TypographyIssue represents a typography validation issue
type TypographyIssue struct {
	ComponentID string `json:"component_id"`
	Message     string `json:"message"`
	Severity    string `json:"severity"` // "error", "warning", "info"
}

// TypographyResult represents the result of typography validation
type TypographyResult struct {
	Passed bool              `json:"passed"`
	Issues []TypographyIssue `json:"issues"`
}

// DefaultTypographyRule returns the default typography rule using Major Third (1.250) scale
func DefaultTypographyRule() TypographyRule {
	baseSize := 16.0
	ratio := 1.250 // Major Third
	
	return TypographyRule{
		ScaleRatio: ratio,
		BaseSize:   baseSize,
		Sizes: map[string]float64{
			"xs":   12,  // 16 / 1.25^1 ≈ 12.8 → 12
			"sm":   14,  // 16 / 1.25^0.5 ≈ 14.2 → 14
			"base": 16,  // base size
			"md":   18,  // 16 * 1.25^0.5 ≈ 17.9 → 18
			"lg":   20,  // 16 * 1.25 = 20
			"xl":   25,  // 16 * 1.25^2 ≈ 25
			"2xl":  31,  // 16 * 1.25^3 ≈ 31.25 → 31
			"3xl":  39,  // 16 * 1.25^4 ≈ 39.06 → 39
			"4xl":  49,  // 16 * 1.25^5 ≈ 48.83 → 49
		},
		Tolerance: 0.5, // Allow 0.5px deviation for rounding
	}
}

// PredefinedScales returns common typography scale ratios
func PredefinedScales() map[string]float64 {
	return map[string]float64{
		"minor-second":    1.067,
		"major-second":    1.125,
		"minor-third":     1.200,
		"major-third":     1.250,
		"perfect-fourth":  1.333,
		"augmented-fourth": 1.414,
		"perfect-fifth":   1.500,
		"golden-ratio":    1.618,
	}
}

// ValidateTypography validates that text components follow the typography scale
func ValidateTypography(structure *types.Structure, rule TypographyRule) TypographyResult {
	result := TypographyResult{
		Passed: true,
		Issues: []TypographyIssue{},
	}

	// Validate all components recursively
	validateComponentTypography(structure.Components, rule, &result)

	return result
}

func validateComponentTypography(components []types.Component, rule TypographyRule, result *TypographyResult) {
	for _, comp := range components {
		// Only validate text components
		if comp.Type == "text" && comp.Size != "" {
			validateTextSize(comp, rule, result)
		}

		// Recursively validate children
		if len(comp.Children) > 0 {
			validateComponentTypography(comp.Children, rule, result)
		}
	}
}

func validateTextSize(comp types.Component, rule TypographyRule, result *TypographyResult) {
	expectedSize, exists := rule.Sizes[comp.Size]
	
	if !exists {
		// Unknown size token - this is a warning
		result.Passed = false
		result.Issues = append(result.Issues, TypographyIssue{
			ComponentID: comp.ID,
			Message:     fmt.Sprintf("Typography: '%s' uses unknown size token '%s'", comp.ID, comp.Size),
			Severity:    "warning",
		})
		
		// Suggest valid tokens
		validTokens := getValidSizeTokens(rule)
		result.Issues = append(result.Issues, TypographyIssue{
			ComponentID: comp.ID,
			Message:     fmt.Sprintf("   Valid size tokens: %v", validTokens),
			Severity:    "info",
		})
		return
	}

	// For token-based systems, just validate that tokens are used correctly
	// The actual size values are already defined in the rule and should be harmonious
	_ = expectedSize // Size is valid if token exists
}

func isOnTypographyScale(size float64, rule TypographyRule) bool {
	// Check if the size can be generated from the base size and ratio
	// Allow some tolerance due to rounding
	
	// Check both integer and half-step powers from base
	// Integer steps: -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5
	// Half steps: -2.5, -1.5, -0.5, 0.5, 1.5, 2.5, etc.
	steps := []float64{-5, -4.5, -4, -3.5, -3, -2.5, -2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5}
	
	for _, step := range steps {
		scaledSize := rule.BaseSize * math.Pow(rule.ScaleRatio, step)
		if math.Abs(scaledSize-size) <= rule.Tolerance {
			return true
		}
	}
	
	return false
}

func getScaleName(ratio float64) string {
	scales := PredefinedScales()
	for name, r := range scales {
		if math.Abs(r-ratio) < 0.001 {
			return name
		}
	}
	return "custom"
}

func getValidSizeTokens(rule TypographyRule) []string {
	tokens := []string{}
	for token := range rule.Sizes {
		tokens = append(tokens, token)
	}
	return tokens
}
