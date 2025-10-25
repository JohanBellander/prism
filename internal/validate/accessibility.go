package validate

import (
	"fmt"
	"strings"

	"github.com/johanbellander/prism/internal/types"
)

// A11yRule defines validation rules for accessibility
type A11yRule struct {
	RequireLabels         bool // All interactive elements need labels
	RequireHeadingOrder   bool // h1 → h2 → h3 (no skipping)
	MaxNestingDepth       int  // 4 levels
	RequireFocusIndicator bool // All interactive elements
	CheckTabOrder         bool // Verify logical tab sequence
}

// DefaultA11yRule returns the default accessibility validation rules
func DefaultA11yRule() A11yRule {
	return A11yRule{
		RequireLabels:         true,
		RequireHeadingOrder:   true,
		MaxNestingDepth:       4,
		RequireFocusIndicator: true,
		CheckTabOrder:         true,
	}
}

// A11yIssue represents a single accessibility validation issue
type A11yIssue struct {
	Severity  string // "error", "warning", "info"
	Message   string
	Component string // Component ID if applicable
}

// A11yResult represents the result of accessibility validation
type A11yResult struct {
	Passed bool
	Issues []A11yIssue
}

// ComponentWithOrder represents a component with its tab order
type ComponentWithOrder struct {
	Component *types.Component
	Order     int
	Depth     int
}

// ValidateAccessibility validates accessibility requirements
func ValidateAccessibility(structure *types.Structure, rule A11yRule) A11yResult {
	result := A11yResult{
		Passed: true,
		Issues: []A11yIssue{},
	}

	// Collect all components with their order and depth
	orderedComponents := []ComponentWithOrder{}
	interactiveComponents := []*types.Component{}
	headings := []struct {
		component *types.Component
		level     int
	}{}

	var traverse func(comp *types.Component, order *int, depth int)
	traverse = func(comp *types.Component, order *int, depth int) {
		// Check nesting depth
		if depth > rule.MaxNestingDepth {
			result.Issues = append(result.Issues, A11yIssue{
				Severity:  "error",
				Message:   fmt.Sprintf("A11y: Component '%s' exceeds max nesting depth (%d levels)", comp.ID, rule.MaxNestingDepth),
				Component: comp.ID,
			})
			result.Passed = false
		}

		// Track component order
		orderedComponents = append(orderedComponents, ComponentWithOrder{
			Component: comp,
			Order:     *order,
			Depth:     depth,
		})
		*order++

		// Check if it's interactive
		if isInteractiveElement(comp) {
			interactiveComponents = append(interactiveComponents, comp)
		}

		// Check if it's a heading
		if comp.Type == "text" {
			level := getHeadingLevel(comp)
			if level > 0 {
				headings = append(headings, struct {
					component *types.Component
					level     int
				}{comp, level})
			}
		}

		// Recurse into children
		for i := range comp.Children {
			traverse(&comp.Children[i], order, depth+1)
		}
	}

	// Start traversal
	order := 0
	for i := range structure.Components {
		traverse(&structure.Components[i], &order, 0)
	}

	// Check for missing labels on interactive elements
	if rule.RequireLabels {
		for _, comp := range interactiveComponents {
			if !hasLabel(comp, structure) {
				result.Issues = append(result.Issues, A11yIssue{
					Severity:  "error",
					Message:   fmt.Sprintf("A11y: '%s' missing label", comp.ID),
					Component: comp.ID,
				})
				result.Passed = false
			}
		}
	}

	// Check heading order
	if rule.RequireHeadingOrder && len(headings) > 1 {
		for i := 1; i < len(headings); i++ {
			prevLevel := headings[i-1].level
			currLevel := headings[i].level

			// Check if skipping levels (e.g., h1 to h3)
			if currLevel > prevLevel+1 {
				result.Issues = append(result.Issues, A11yIssue{
					Severity:  "error",
					Message:   fmt.Sprintf("A11y: Heading structure jumps from h%d to h%d (missing h%d)", prevLevel, currLevel, prevLevel+1),
					Component: headings[i].component.ID,
				})
				result.Passed = false
			}
		}
	}

	// Check focus indicators
	if rule.RequireFocusIndicator {
		// In Phase 1 structure, we check that focus_indicators is defined in accessibility
		if structure.Accessibility.FocusIndicators == "" {
			result.Issues = append(result.Issues, A11yIssue{
				Severity:  "warning",
				Message:   "A11y: Focus indicators not defined in accessibility settings",
				Component: "",
			})
		} else if structure.Accessibility.FocusIndicators != "visible" {
			result.Issues = append(result.Issues, A11yIssue{
				Severity:  "warning",
				Message:   fmt.Sprintf("A11y: Focus indicators set to '%s' - recommend 'visible'", structure.Accessibility.FocusIndicators),
				Component: "",
			})
		}
	}

	// Check tab order
	if rule.CheckTabOrder {
		// Analyze if interactive elements appear in a logical order
		interactiveOrder := []ComponentWithOrder{}
		for _, ordered := range orderedComponents {
			if isInteractiveElement(ordered.Component) {
				interactiveOrder = append(interactiveOrder, ordered)
			}
		}

		// Check for potential confusing tab order (e.g., button before its related input)
		for i := 0; i < len(interactiveOrder)-1; i++ {
			curr := interactiveOrder[i]
			next := interactiveOrder[i+1]

			// If a button comes before an input that shares a name prefix, it might be confusing
			if curr.Component.Type == "button" && next.Component.Type == "input" {
				if sharesPrefix(curr.Component.ID, next.Component.ID) {
					result.Issues = append(result.Issues, A11yIssue{
						Severity:  "warning",
						Message:   fmt.Sprintf("A11y: Tab order may be confusing - '%s' comes before '%s' in layout", curr.Component.ID, next.Component.ID),
						Component: curr.Component.ID,
					})
				}
			}
		}
	}

	// Check semantic structure
	if structure.Accessibility.SemanticStructure {
		// Verify that roles are used appropriately
		roleCount := 0
		for _, comp := range orderedComponents {
			if comp.Component.Role != "" {
				roleCount++
			}
		}
		
		if roleCount == 0 {
			result.Issues = append(result.Issues, A11yIssue{
				Severity:  "info",
				Message:   "A11y: Semantic structure enabled but no roles defined - consider adding roles like 'header', 'navigation', 'main', 'footer'",
				Component: "",
			})
		}
	}

	// Add success messages if no major issues found
	if len(result.Issues) == 0 || (result.Passed && len(result.Issues) <= 2) {
		if rule.RequireLabels && len(interactiveComponents) > 0 {
			result.Issues = append(result.Issues, A11yIssue{
				Severity: "info",
				Message:  "✓ All interactive elements have labels",
			})
		}
		
		if rule.RequireHeadingOrder && len(headings) > 0 {
			result.Issues = append(result.Issues, A11yIssue{
				Severity: "info",
				Message:  "✓ Heading hierarchy is correct",
			})
		}
		
		if rule.RequireFocusIndicator && structure.Accessibility.FocusIndicators == "visible" {
			result.Issues = append(result.Issues, A11yIssue{
				Severity: "info",
				Message:  "✓ Focus indicators are properly defined",
			})
		}
		
		if len(orderedComponents) > 0 {
			maxDepth := 0
			for _, comp := range orderedComponents {
				if comp.Depth > maxDepth {
					maxDepth = comp.Depth
				}
			}
			result.Issues = append(result.Issues, A11yIssue{
				Severity: "info",
				Message:  fmt.Sprintf("✓ Nesting depth (%d) within acceptable limits (%d)", maxDepth, rule.MaxNestingDepth),
			})
		}
	}

	return result
}

// getHeadingLevel extracts heading level from component ID or size
func getHeadingLevel(comp *types.Component) int {
	// Check ID for explicit heading level (h1, h2, h3, etc.)
	idLower := strings.ToLower(comp.ID)
	if strings.HasPrefix(idLower, "h") && len(idLower) >= 2 {
		if idLower[1] >= '1' && idLower[1] <= '6' {
			return int(idLower[1] - '0')
		}
	}

	// Check for heading in ID or role
	if strings.Contains(idLower, "heading") || strings.Contains(idLower, "title") {
		// Infer level from size
		sizeMap := map[string]int{
			"4xl": 1,
			"3xl": 2,
			"2xl": 3,
			"xl":  4,
		}
		
		if level, ok := sizeMap[comp.Size]; ok {
			return level
		}
	}

	return 0
}

// hasLabel checks if an interactive component has an associated label
func hasLabel(comp *types.Component, structure *types.Structure) bool {
	// Check if there's a text component with a matching ID pattern
	// e.g., "username-input" should have "username-label"
	labelID := ""
	
	// Try to find label by removing common input suffixes
	baseID := strings.TrimSuffix(comp.ID, "-input")
	baseID = strings.TrimSuffix(baseID, "-field")
	baseID = strings.TrimSuffix(baseID, "-button")
	baseID = strings.TrimSuffix(baseID, "-btn")
	
	if baseID != comp.ID {
		labelID = baseID + "-label"
	}
	
	// Search for the label
	var findLabel func(components []types.Component) bool
	findLabel = func(components []types.Component) bool {
		for i := range components {
			if components[i].ID == labelID && components[i].Type == "text" {
				return true
			}
			if findLabel(components[i].Children) {
				return true
			}
		}
		return false
	}
	
	if labelID != "" && findLabel(structure.Components) {
		return true
	}
	
	// Check if the component itself has content (self-labeling button)
	if comp.Content != "" {
		return true
	}
	
	// Check accessibility labels field
	if structure.Accessibility.Labels == "all_interactive_elements" {
		// Assume labels are planned/will be added
		return true
	}
	
	return false
}

// sharesPrefix checks if two component IDs share a common prefix
func sharesPrefix(id1, id2 string) bool {
	parts1 := strings.Split(id1, "-")
	parts2 := strings.Split(id2, "-")
	
	if len(parts1) > 0 && len(parts2) > 0 {
		return parts1[0] == parts2[0]
	}
	
	return false
}
