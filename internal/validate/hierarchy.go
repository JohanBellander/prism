package validate

import (
	"fmt"
	"strings"

	"github.com/johanbellander/prism/internal/types"
)

// HierarchyRule defines validation rules for visual hierarchy
type HierarchyRule struct {
	HeadingScaleRatio float64 // e.g., 1.25 (each level 25% larger)
	MinPrimaryCTASize int     // e.g., 120px width minimum
	SpacingScaleRatio float64 // e.g., 1.5 (parent spacing > child spacing)
}

// DefaultHierarchyRule returns the default hierarchy validation rules
func DefaultHierarchyRule() HierarchyRule {
	return HierarchyRule{
		HeadingScaleRatio: 1.25,
		MinPrimaryCTASize: 120,
		SpacingScaleRatio: 1.5,
	}
}

// HierarchyIssue represents a single hierarchy validation issue
type HierarchyIssue struct {
	Severity string // "error", "warning", "info"
	Message  string
	Component string // Component ID if applicable
}

// HierarchyResult represents the result of hierarchy validation
type HierarchyResult struct {
	Passed bool
	Issues []HierarchyIssue
}

// ValidateHierarchy validates the visual hierarchy of a structure
func ValidateHierarchy(structure *types.Structure, rule HierarchyRule) HierarchyResult {
	result := HierarchyResult{
		Passed: true,
		Issues: []HierarchyIssue{},
	}

	// Map text size names to relative pixel sizes for comparison
	sizeMap := map[string]float64{
		"xs":   12,
		"sm":   14,
		"base": 16,
		"lg":   18,
		"xl":   20,
		"2xl":  24,
		"3xl":  30,
		"4xl":  36,
	}

	// Collect all text elements and buttons
	textElements := []struct {
		component *types.Component
		size      float64
		isHeading bool
		level     int // 1 for h1, 2 for h2, etc.
	}{}
	
	buttons := []struct {
		component *types.Component
		isPrimary bool
		width     int
	}{}

	// Traverse components to collect text elements and buttons
	var traverse func(comp *types.Component, parentSpacing int)
	traverse = func(comp *types.Component, parentSpacing int) {
		// Check if it's a text element
		if comp.Type == "text" {
			size := sizeMap["base"] // default
			if comp.Size != "" {
				if s, ok := sizeMap[comp.Size]; ok {
					size = s
				}
			}

			// Determine if it's a heading based on size and role
			isHeading := false
			level := 0
			
			// Check ID for explicit heading level (h1, h2, h3, etc.)
			idLower := strings.ToLower(comp.ID)
			if strings.HasPrefix(idLower, "h") && len(idLower) >= 2 {
				if idLower[1] >= '1' && idLower[1] <= '6' {
					isHeading = true
					level = int(idLower[1] - '0')
				}
			}
			
			// Infer heading level from size if not already determined (larger = higher level heading)
			if level == 0 && size >= sizeMap["2xl"] {
				isHeading = true
				if size >= sizeMap["4xl"] {
					level = 1
				} else if size >= sizeMap["3xl"] {
					level = 2
				} else if size >= sizeMap["2xl"] {
					level = 3
				}
			}

			// Also check role for explicit heading indication
			if strings.Contains(strings.ToLower(comp.Role), "heading") || 
			   strings.Contains(idLower, "title") ||
			   strings.Contains(idLower, "heading") {
				isHeading = true
				if level == 0 {
					// Assign level based on size if not already assigned
					if size >= sizeMap["3xl"] {
						level = 1
					} else if size >= sizeMap["2xl"] {
						level = 2
					} else {
						level = 3
					}
				}
			}

			textElements = append(textElements, struct {
				component *types.Component
				size      float64
				isHeading bool
				level     int
			}{comp, size, isHeading, level})
		}

		// Check if it's a button
		if comp.Type == "button" {
			isPrimary := strings.Contains(strings.ToLower(comp.ID), "primary") ||
			             strings.Contains(strings.ToLower(comp.Role), "primary") ||
			             comp.ID == structure.Intent.PrimaryAction
			
			width := comp.Layout.Width
			if width == 0 {
				width = 100 // default minimum
			}

			buttons = append(buttons, struct {
				component *types.Component
				isPrimary bool
				width     int
			}{comp, isPrimary, width})
		}

		// Check spacing hierarchy
		if comp.Layout.Padding > 0 && parentSpacing > 0 {
			// Parent should have equal or greater spacing
			if parentSpacing > comp.Layout.Padding {
				expectedChildSpacing := float64(parentSpacing) / rule.SpacingScaleRatio
				if float64(comp.Layout.Padding) < expectedChildSpacing*0.8 { // 20% tolerance
					result.Issues = append(result.Issues, HierarchyIssue{
						Severity:  "info",
						Message:   fmt.Sprintf("Spacing hierarchy: '%s' has padding %dpx (parent has %dpx) - consider using %.0fpx for consistent hierarchy", comp.ID, comp.Layout.Padding, parentSpacing, expectedChildSpacing),
						Component: comp.ID,
					})
				}
			}
		}

		// Recurse into children
		currentSpacing := comp.Layout.Padding
		if currentSpacing == 0 {
			currentSpacing = parentSpacing
		}
		for i := range comp.Children {
			traverse(&comp.Children[i], currentSpacing)
		}
	}

	// Start traversal with top-level spacing
	for i := range structure.Components {
		traverse(&structure.Components[i], structure.Layout.Spacing)
	}

	// Validate heading size hierarchy
	headings := []struct {
		component *types.Component
		size      float64
		level     int
	}{}
	for _, te := range textElements {
		if te.isHeading {
			headings = append(headings, struct {
				component *types.Component
				size      float64
				level     int
			}{te.component, te.size, te.level})
		}
	}

	// Check that higher-level headings are larger
	for i := 0; i < len(headings); i++ {
		for j := i + 1; j < len(headings); j++ {
			h1, h2 := headings[i], headings[j]
			
			// If h1 is a higher level (smaller number) than h2, it should be larger
			if h1.level < h2.level && h1.size <= h2.size {
				expectedRatio := 1.0
				for k := h2.level; k < h1.level; k++ {
					expectedRatio *= rule.HeadingScaleRatio
				}
				expectedSize := h2.size * expectedRatio
				
				result.Issues = append(result.Issues, HierarchyIssue{
					Severity:  "warning",
					Message:   fmt.Sprintf("h%d ('%s': %.0fpx) not sufficiently larger than h%d ('%s': %.0fpx) - recommend %.0fpx (%.2fx scale)", h1.level, h1.component.ID, h1.size, h2.level, h2.component.ID, h2.size, expectedSize, rule.HeadingScaleRatio),
					Component: h1.component.ID,
				})
				result.Passed = false
			}
		}
	}

	// Validate button sizes (primary CTA should meet minimum size)
	var primaryButtons []struct {
		component *types.Component
		isPrimary bool
		width     int
	}
	var secondaryButtons []struct {
		component *types.Component
		isPrimary bool
		width     int
	}

	for _, btn := range buttons {
		if btn.isPrimary {
			primaryButtons = append(primaryButtons, btn)
			if btn.width < rule.MinPrimaryCTASize {
				result.Issues = append(result.Issues, HierarchyIssue{
					Severity:  "warning",
					Message:   fmt.Sprintf("Primary button '%s' is %dpx wide (recommend minimum %dpx)", btn.component.ID, btn.width, rule.MinPrimaryCTASize),
					Component: btn.component.ID,
				})
				result.Passed = false
			}
		} else {
			secondaryButtons = append(secondaryButtons, btn)
		}
	}

	// Check that primary buttons are not smaller than secondary buttons
	for _, primary := range primaryButtons {
		for _, secondary := range secondaryButtons {
			if primary.width < secondary.width {
				result.Issues = append(result.Issues, HierarchyIssue{
					Severity:  "error",
					Message:   fmt.Sprintf("Secondary button '%s' (%dpx) larger than primary button '%s' (%dpx)", secondary.component.ID, secondary.width, primary.component.ID, primary.width),
					Component: primary.component.ID,
				})
				result.Passed = false
			}
		}
	}

	// If no issues found, add success message
	if len(result.Issues) == 0 {
		result.Issues = append(result.Issues, HierarchyIssue{
			Severity: "info",
			Message:  "✓ Spacing hierarchy is consistent",
		})
		if len(headings) > 0 {
			result.Issues = append(result.Issues, HierarchyIssue{
				Severity: "info",
				Message:  "✓ Heading sizes follow consistent scale",
			})
		}
		if len(primaryButtons) > 0 {
			result.Issues = append(result.Issues, HierarchyIssue{
				Severity: "info",
				Message:  "✓ Primary CTA buttons meet minimum size requirements",
			})
		}
	}

	return result
}
