package validate

import (
	"fmt"

	"github.com/johanbellander/prism/internal/types"
)

// DarkModeIssue represents a dark mode validation issue
type DarkModeIssue struct {
	ComponentID string `json:"component_id"`
	Message     string `json:"message"`
	Severity    string `json:"severity"` // "error", "warning", "info"
	Mode        string `json:"mode,omitempty"` // "light", "dark", "both"
}

// DarkModeResult contains the validation results
type DarkModeResult struct {
	Passed bool            `json:"passed"`
	Issues []DarkModeIssue `json:"issues"`
}

// DarkModeRule defines the dark mode validation rules
type DarkModeRule struct {
	RequireSemanticColors bool    // Whether semantic color tokens are required
	ValidateContrast      bool    // Whether to validate contrast in both modes
	MinContrastRatio      float64 // Minimum contrast ratio for both modes
	RecommendAdaptive     bool    // Whether to recommend adaptive colors
}

// DefaultDarkModeRule returns the default dark mode validation rules
func DefaultDarkModeRule() DarkModeRule {
	return DarkModeRule{
		RequireSemanticColors: true,
		ValidateContrast:      true,
		MinContrastRatio:      4.5, // WCAG AA standard
		RecommendAdaptive:     true,
	}
}

// ValidateDarkMode validates dark mode support in the design
func ValidateDarkMode(structure *types.Structure, rule DarkModeRule) DarkModeResult {
	result := DarkModeResult{
		Passed: true,
		Issues: []DarkModeIssue{},
	}

	// Check if semantic colors are defined
	if rule.RequireSemanticColors {
		// In Phase 1, we don't have semantic colors in the schema yet
		// This validator provides recommendations for dark mode support
		result.Issues = append(result.Issues, DarkModeIssue{
			ComponentID: "structure",
			Message:     "Consider defining semantic color tokens for dark mode support (e.g., 'text.primary', 'background.surface')",
			Severity:    "info",
			Mode:        "both",
		})
	}

	// Check components for hardcoded colors
	for _, component := range structure.Components {
		validateComponentDarkMode(&result, &component, rule)
	}

	// If no errors found, mark as passed
	if len(result.Issues) == 0 {
		result.Passed = true
	} else {
		// Check if there are any errors (not just warnings/info)
		hasErrors := false
		for _, issue := range result.Issues {
			if issue.Severity == "error" {
				hasErrors = true
				break
			}
		}
		result.Passed = !hasErrors
	}

	return result
}

func validateComponentDarkMode(result *DarkModeResult, component *types.Component, rule DarkModeRule) {
	// Check for hardcoded colors that might not work well in dark mode
	if component.Color != "" && rule.RecommendAdaptive {
		// Pure black or pure white text might not be ideal for both modes
		if component.Color == "#000000" || component.Color == "#FFFFFF" {
			result.Issues = append(result.Issues, DarkModeIssue{
				ComponentID: component.ID,
				Message:     fmt.Sprintf("Component '%s' uses absolute color '%s' which may not adapt well to dark mode. Consider using semantic color tokens.", component.ID, component.Color),
				Severity:    "info",
				Mode:        "both",
			})
		}
	}

	// Check background colors
	if component.Layout.Background != "" && rule.RecommendAdaptive {
		if component.Layout.Background == "#FFFFFF" || component.Layout.Background == "#000000" {
			result.Issues = append(result.Issues, DarkModeIssue{
				ComponentID: component.ID,
				Message:     fmt.Sprintf("Component '%s' uses absolute background color '%s' which may not adapt to dark mode. Consider semantic tokens like 'background.primary'.", component.ID, component.Layout.Background),
				Severity:    "info",
				Mode:        "both",
			})
		}
	}

	// Check children recursively
	for _, child := range component.Children {
		validateComponentDarkMode(result, &child, rule)
	}
}
