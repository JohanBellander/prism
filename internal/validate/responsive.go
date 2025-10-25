package validate

import (
	"fmt"

	"github.com/johanbellander/prism/internal/types"
)

// ResponsiveIssue represents a responsive design issue
type ResponsiveIssue struct {
	ComponentID string `json:"component_id"`
	Message     string `json:"message"`
	Severity    string `json:"severity"` // "error", "warning", "info"
	Viewport    string `json:"viewport,omitempty"`
}

// ResponsiveResult contains the validation results
type ResponsiveResult struct {
	Passed bool               `json:"passed"`
	Issues []ResponsiveIssue  `json:"issues"`
}

// ResponsiveRule defines the responsive validation rules
type ResponsiveRule struct {
	Breakpoints       map[string]int // Viewport name -> width in pixels
	MinTouchTarget    int            // Minimum touch target size for mobile
	CheckOverflow     bool           // Whether to check for content overflow
	CheckTouchTargets bool           // Whether to validate touch targets at each breakpoint
}

// DefaultResponsiveRule returns the default responsive validation rules
func DefaultResponsiveRule() ResponsiveRule {
	return ResponsiveRule{
		Breakpoints: map[string]int{
			"mobile":  375,
			"tablet":  768,
			"desktop": 1440,
		},
		MinTouchTarget:    44,
		CheckOverflow:     true,
		CheckTouchTargets: true,
	}
}

// ValidateResponsive validates responsive design at different breakpoints
func ValidateResponsive(structure *types.Structure, rule ResponsiveRule) ResponsiveResult {
	result := ResponsiveResult{
		Passed: true,
		Issues: []ResponsiveIssue{},
	}

	// Get the layout max width if defined
	layoutMaxWidth := structure.Layout.MaxWidth

	// Check each breakpoint
	for viewport, viewportWidth := range rule.Breakpoints {
		// Check if layout exceeds viewport
		if layoutMaxWidth > 0 && layoutMaxWidth > viewportWidth {
			result.Issues = append(result.Issues, ResponsiveIssue{
				ComponentID: "layout",
				Message:     fmt.Sprintf("Layout max-width (%dpx) exceeds %s viewport (%dpx)", layoutMaxWidth, viewport, viewportWidth),
				Severity:    "warning",
				Viewport:    viewport,
			})
		}

		// Validate components at this breakpoint
		for _, component := range structure.Components {
			validateComponentAtViewport(&result, &component, viewport, viewportWidth, rule, 0, 0)
		}
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

func validateComponentAtViewport(result *ResponsiveResult, component *types.Component, viewport string, viewportWidth int, rule ResponsiveRule, parentX, parentY int) {
	// Get component dimensions from layout
	width := component.Layout.Width
	height := component.Layout.Height

	// For absolute positioning, we would need X/Y coordinates
	// Since Phase 1 uses layout properties, we check max width
	if component.Layout.MaxWidth > 0 {
		if rule.CheckOverflow && component.Layout.MaxWidth > viewportWidth {
			result.Issues = append(result.Issues, ResponsiveIssue{
				ComponentID: component.ID,
				Message:     fmt.Sprintf("Component '%s' max-width (%dpx) exceeds %s viewport (%dpx)", component.ID, component.Layout.MaxWidth, viewport, viewportWidth),
				Severity:    "warning",
				Viewport:    viewport,
			})
		}
	}

	// Check width if defined
	if width > 0 && rule.CheckOverflow && width > viewportWidth {
		result.Issues = append(result.Issues, ResponsiveIssue{
			ComponentID: component.ID,
			Message:     fmt.Sprintf("Component '%s' width (%dpx) exceeds %s viewport (%dpx)", component.ID, width, viewport, viewportWidth),
			Severity:    "warning",
			Viewport:    viewport,
		})
	}

	// For mobile viewport, check touch targets
	if viewport == "mobile" && rule.CheckTouchTargets {
		if component.Type == "button" || component.Type == "input" {
			if width > 0 && height > 0 {
				if width < rule.MinTouchTarget || height < rule.MinTouchTarget {
					result.Issues = append(result.Issues, ResponsiveIssue{
						ComponentID: component.ID,
						Message:     fmt.Sprintf("Interactive element '%s' (%dx%dpx) is too small for mobile (minimum %dx%dpx recommended)", component.ID, width, height, rule.MinTouchTarget, rule.MinTouchTarget),
						Severity:    "warning",
						Viewport:    viewport,
					})
				}
			}
		}
	}

	// Check children recursively
	for _, child := range component.Children {
		validateComponentAtViewport(result, &child, viewport, viewportWidth, rule, parentX+width, parentY+height)
	}
}
