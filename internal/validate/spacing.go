package validate

import (
	"fmt"
	"math"

	"github.com/johanbellander/prism/internal/types"
)

// SpacingRule defines validation rules for spacing (8pt grid system)
type SpacingRule struct {
	BaseUnit         int     // 8px base unit
	AllowedScale     []int   // Allowed spacing values: 0, 4, 8, 12, 16, 24, 32, 48, 64, 96, 128
	WarnOffGrid      bool    // Warn when values are off-grid
	AllowHalfStep    bool    // Allow 4px for fine-tuning
	MaxHalfStepUsage int     // Maximum number of 4px usages before warning
}

// DefaultSpacingRule returns the default 8pt grid validation rules
func DefaultSpacingRule() SpacingRule {
	return SpacingRule{
		BaseUnit:         8,
		AllowedScale:     []int{0, 4, 8, 12, 16, 24, 32, 48, 64, 96, 128},
		WarnOffGrid:      true,
		AllowHalfStep:    true,
		MaxHalfStepUsage: 5,
	}
}

// SpacingIssue represents a single spacing validation issue
type SpacingIssue struct {
	Severity    string // "error", "warning", "info"
	Category    string // e.g., "off_grid", "excessive_half_step"
	Message     string
	ComponentID string // Component ID if applicable
	Property    string // e.g., "padding", "gap", "margin_bottom"
	Value       int    // Actual value used
	Suggested   int    // Suggested value on grid
}

// SpacingResult represents the result of spacing validation
type SpacingResult struct {
	Passed bool
	Issues []SpacingIssue
}

// ValidateSpacing validates that spacing follows 8pt grid system
func ValidateSpacing(structure *types.Structure, rule SpacingRule) SpacingResult {
	result := SpacingResult{
		Passed: true,
		Issues: []SpacingIssue{},
	}

	halfStepCount := 0

	// Analyze all components for spacing values
	var analyzeComponent func(comp *types.Component, depth int)
	analyzeComponent = func(comp *types.Component, depth int) {
		// Check layout padding
		if comp.Layout.Padding > 0 {
			if !isOnGrid(comp.Layout.Padding, rule.AllowedScale) {
				suggested := findNearestGridValue(comp.Layout.Padding, rule.AllowedScale)
				result.Issues = append(result.Issues, SpacingIssue{
					Severity:    "warning",
					Category:    "off_grid",
					Message:     fmt.Sprintf("Spacing: '%s' padding uses %dpx (not on 8pt grid)", comp.ID, comp.Layout.Padding),
					ComponentID: comp.ID,
					Property:    "padding",
					Value:       comp.Layout.Padding,
					Suggested:   suggested,
				})
				result.Passed = false
				
				// Add suggestion
				result.Issues = append(result.Issues, SpacingIssue{
					Severity:    "info",
					Category:    "suggestion",
					Message:     fmt.Sprintf("   Suggestion: Use %dpx for consistency", suggested),
					ComponentID: comp.ID,
					Property:    "padding",
					Suggested:   suggested,
				})
				
				// Track half-step usage
				if comp.Layout.Padding%4 == 0 && comp.Layout.Padding%8 != 0 {
					halfStepCount++
				}
			}
		}

		// Check layout gap
		if comp.Layout.Gap > 0 {
			if !isOnGrid(comp.Layout.Gap, rule.AllowedScale) {
				suggested := findNearestGridValue(comp.Layout.Gap, rule.AllowedScale)
				result.Issues = append(result.Issues, SpacingIssue{
					Severity:    "warning",
					Category:    "off_grid",
					Message:     fmt.Sprintf("Spacing: '%s' gap uses %dpx (not on 8pt grid)", comp.ID, comp.Layout.Gap),
					ComponentID: comp.ID,
					Property:    "gap",
					Value:       comp.Layout.Gap,
					Suggested:   suggested,
				})
				result.Passed = false
				
				result.Issues = append(result.Issues, SpacingIssue{
					Severity:    "info",
					Category:    "suggestion",
					Message:     fmt.Sprintf("   Suggestion: Use %dpx for consistency", suggested),
					ComponentID: comp.ID,
					Property:    "gap",
					Suggested:   suggested,
				})
				
				if comp.Layout.Gap%4 == 0 && comp.Layout.Gap%8 != 0 {
					halfStepCount++
				}
			}
		}

		// Check margin bottom
		if comp.Layout.MarginBottom > 0 {
			if !isOnGrid(comp.Layout.MarginBottom, rule.AllowedScale) {
				suggested := findNearestGridValue(comp.Layout.MarginBottom, rule.AllowedScale)
				result.Issues = append(result.Issues, SpacingIssue{
					Severity:    "warning",
					Category:    "off_grid",
					Message:     fmt.Sprintf("Spacing: '%s' margin_bottom uses %dpx (not on 8pt grid)", comp.ID, comp.Layout.MarginBottom),
					ComponentID: comp.ID,
					Property:    "margin_bottom",
					Value:       comp.Layout.MarginBottom,
					Suggested:   suggested,
				})
				result.Passed = false
				
				result.Issues = append(result.Issues, SpacingIssue{
					Severity:    "info",
					Category:    "suggestion",
					Message:     fmt.Sprintf("   Suggestion: Use %dpx for consistency", suggested),
					ComponentID: comp.ID,
					Property:    "margin_bottom",
					Suggested:   suggested,
				})
				
				if comp.Layout.MarginBottom%4 == 0 && comp.Layout.MarginBottom%8 != 0 {
					halfStepCount++
				}
			}
		}

		// Recurse into children
		for i := range comp.Children {
			analyzeComponent(&comp.Children[i], depth+1)
		}
	}

	// Check top-level layout spacing
	if structure.Layout.Spacing > 0 {
		if !isOnGrid(structure.Layout.Spacing, rule.AllowedScale) {
			suggested := findNearestGridValue(structure.Layout.Spacing, rule.AllowedScale)
			result.Issues = append(result.Issues, SpacingIssue{
				Severity:    "warning",
				Category:    "off_grid",
				Message:     fmt.Sprintf("Spacing: Layout spacing uses %dpx (not on 8pt grid)", structure.Layout.Spacing),
				ComponentID: "layout",
				Property:    "spacing",
				Value:       structure.Layout.Spacing,
				Suggested:   suggested,
			})
			result.Passed = false
			
			result.Issues = append(result.Issues, SpacingIssue{
				Severity:    "info",
				Category:    "suggestion",
				Message:     fmt.Sprintf("   Suggestion: Use %dpx for consistency", suggested),
				ComponentID: "layout",
				Property:    "spacing",
				Suggested:   suggested,
			})
			
			if structure.Layout.Spacing%4 == 0 && structure.Layout.Spacing%8 != 0 {
				halfStepCount++
			}
		}
	}

	// Check top-level layout padding
	if structure.Layout.Padding > 0 {
		if !isOnGrid(structure.Layout.Padding, rule.AllowedScale) {
			suggested := findNearestGridValue(structure.Layout.Padding, rule.AllowedScale)
			result.Issues = append(result.Issues, SpacingIssue{
				Severity:    "warning",
				Category:    "off_grid",
				Message:     fmt.Sprintf("Spacing: Layout padding uses %dpx (not on 8pt grid)", structure.Layout.Padding),
				ComponentID: "layout",
				Property:    "padding",
				Value:       structure.Layout.Padding,
				Suggested:   suggested,
			})
			result.Passed = false
			
			result.Issues = append(result.Issues, SpacingIssue{
				Severity:    "info",
				Category:    "suggestion",
				Message:     fmt.Sprintf("   Suggestion: Use %dpx for consistency", suggested),
				ComponentID: "layout",
				Property:    "padding",
				Suggested:   suggested,
			})
			
			if structure.Layout.Padding%4 == 0 && structure.Layout.Padding%8 != 0 {
				halfStepCount++
			}
		}
	}

	// Analyze all top-level components
	for i := range structure.Components {
		analyzeComponent(&structure.Components[i], 0)
	}

	// Check for excessive half-step usage
	if rule.AllowHalfStep && halfStepCount > rule.MaxHalfStepUsage {
		result.Issues = append(result.Issues, SpacingIssue{
			Severity: "warning",
			Category: "excessive_half_step",
			Message:  fmt.Sprintf("Excessive use of 4px half-steps (%d occurrences) - consider using 8px base unit", halfStepCount),
		})
	}

	return result
}

// isOnGrid checks if a value is on the allowed spacing scale
func isOnGrid(value int, allowedScale []int) bool {
	for _, allowed := range allowedScale {
		if value == allowed {
			return true
		}
	}
	return false
}

// findNearestGridValue finds the nearest value on the spacing grid
func findNearestGridValue(value int, allowedScale []int) int {
	if len(allowedScale) == 0 {
		return value
	}

	nearest := allowedScale[0]
	minDiff := math.Abs(float64(value - allowedScale[0]))

	for _, allowed := range allowedScale {
		diff := math.Abs(float64(value - allowed))
		if diff < minDiff {
			minDiff = diff
			nearest = allowed
		}
	}

	return nearest
}
