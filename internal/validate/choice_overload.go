package validate

import (
	"fmt"
	"strings"

	"github.com/johanbellander/prism/internal/types"
)

// ChoiceRule defines validation rules for choice overload (Hick's Law)
type ChoiceRule struct {
	MaxNavItems    int // e.g., 7 ± 2 (Miller's Law)
	MaxFormFields  int // e.g., 5-7 per section
	MaxButtonGroup int // e.g., 3 buttons max in a group
	MaxCardGrid    int // e.g., 12 cards before pagination
}

// DefaultChoiceRule returns the default choice overload validation rules
func DefaultChoiceRule() ChoiceRule {
	return ChoiceRule{
		MaxNavItems:    7,
		MaxFormFields:  7,
		MaxButtonGroup: 3,
		MaxCardGrid:    12,
	}
}

// ChoiceIssue represents a single choice overload validation issue
type ChoiceIssue struct {
	Severity    string // "error", "warning", "info"
	Category    string // e.g., "navigation_overload", "form_overload"
	Message     string
	ComponentID string // Component ID if applicable
}

// ChoiceResult represents the result of choice overload validation
type ChoiceResult struct {
	Passed bool
	Issues []ChoiceIssue
}

// ValidateChoiceOverload validates Hick's Law (choice overload)
func ValidateChoiceOverload(structure *types.Structure, rule ChoiceRule) ChoiceResult {
	result := ChoiceResult{
		Passed: true,
		Issues: []ChoiceIssue{},
	}

	// Track containers and their interactive element counts
	var analyzeContainer func(comp *types.Component, depth int)
	analyzeContainer = func(comp *types.Component, depth int) {
		// Check if this is a navigation container
		if isNavigationContainer(comp) {
			navItemCount := countInteractiveChildren(comp)
			if navItemCount > rule.MaxNavItems {
				result.Issues = append(result.Issues, ChoiceIssue{
					Severity:    "warning",
					Category:    "navigation_overload",
					Message:     fmt.Sprintf("Choice Overload: Navigation '%s' has %d items - consider grouping or secondary menu (recommended max: %d)", comp.ID, navItemCount, rule.MaxNavItems),
					ComponentID: comp.ID,
				})
				result.Passed = false
			}
		}

		// Check if this is a form container
		if isFormContainer(comp) {
			formFieldCount := countFormFields(comp)
			if formFieldCount > rule.MaxFormFields {
				result.Issues = append(result.Issues, ChoiceIssue{
					Severity:    "warning",
					Category:    "form_overload",
					Message:     fmt.Sprintf("Choice Overload: Form section '%s' has %d fields - consider splitting into steps (recommended max: %d)", comp.ID, formFieldCount, rule.MaxFormFields),
					ComponentID: comp.ID,
				})
				result.Passed = false
			}
		}

		// Check if this is a button group
		// Skip if already checked as navigation or form
		if !isNavigationContainer(comp) && !isFormContainer(comp) && isButtonGroup(comp) {
			buttonCount := countButtons(comp)
			if buttonCount > rule.MaxButtonGroup {
				result.Issues = append(result.Issues, ChoiceIssue{
					Severity:    "warning",
					Category:    "button_group_overload",
					Message:     fmt.Sprintf("Choice Overload: Button group '%s' has %d buttons - consider reducing options (recommended max: %d)", comp.ID, buttonCount, rule.MaxButtonGroup),
					ComponentID: comp.ID,
				})
				result.Passed = false
			}
		}

		// Check if this is a card/grid container
		if isCardGrid(comp) {
			cardCount := countCards(comp)
			if cardCount > rule.MaxCardGrid {
				result.Issues = append(result.Issues, ChoiceIssue{
					Severity:    "warning",
					Category:    "card_grid_overload",
					Message:     fmt.Sprintf("Choice Overload: Grid '%s' has %d items - consider pagination or filtering (recommended max: %d)", comp.ID, cardCount, rule.MaxCardGrid),
					ComponentID: comp.ID,
				})
				result.Passed = false
			}
		}

		// Recurse into children
		for i := range comp.Children {
			analyzeContainer(&comp.Children[i], depth+1)
		}
	}

	// Analyze all top-level components
	for i := range structure.Components {
		analyzeContainer(&structure.Components[i], 0)
	}

	return result
}

// isNavigationContainer checks if a component is a navigation container
func isNavigationContainer(comp *types.Component) bool {
	idLower := strings.ToLower(comp.ID)
	roleLower := strings.ToLower(comp.Role)
	
	return strings.Contains(idLower, "nav") ||
		strings.Contains(idLower, "menu") ||
		roleLower == "navigation" ||
		roleLower == "menu"
}

// isFormContainer checks if a component is a form container
func isFormContainer(comp *types.Component) bool {
	idLower := strings.ToLower(comp.ID)
	roleLower := strings.ToLower(comp.Role)
	
	return strings.Contains(idLower, "form") ||
		strings.Contains(idLower, "signup") ||
		strings.Contains(idLower, "login") ||
		strings.Contains(idLower, "register") ||
		roleLower == "form"
}

// isButtonGroup checks if a component is a button group
func isButtonGroup(comp *types.Component) bool {
	// Count buttons in direct children only
	buttonCount := 0
	for i := range comp.Children {
		if comp.Children[i].Type == "button" {
			buttonCount++
		}
	}
	
	// A button group must have at least 2 buttons
	return buttonCount >= 2
}

// isCardGrid checks if a component is a card/grid container
func isCardGrid(comp *types.Component) bool {
	idLower := strings.ToLower(comp.ID)
	
	return (strings.Contains(idLower, "grid") ||
		strings.Contains(idLower, "card") ||
		strings.Contains(idLower, "list")) &&
		comp.Layout.Display == "grid" &&
		len(comp.Children) > 0
}

// countInteractiveChildren counts interactive elements in direct children
func countInteractiveChildren(comp *types.Component) int {
	count := 0
	for i := range comp.Children {
		if isInteractiveElement(&comp.Children[i]) {
			count++
		}
		// Also count nested interactive elements (e.g., nav items with links)
		count += countInteractiveChildren(&comp.Children[i])
	}
	return count
}

// countFormFields counts input fields in a form
func countFormFields(comp *types.Component) int {
	count := 0
	
	var traverse func(c *types.Component)
	traverse = func(c *types.Component) {
		if c.Type == "input" {
			count++
		}
		for i := range c.Children {
			traverse(&c.Children[i])
		}
	}
	
	traverse(comp)
	return count
}

// countButtons counts buttons in a container
func countButtons(comp *types.Component) int {
	count := 0
	
	var traverse func(c *types.Component)
	traverse = func(c *types.Component) {
		if c.Type == "button" {
			count++
		}
		for i := range c.Children {
			traverse(&c.Children[i])
		}
	}
	
	traverse(comp)
	return count
}

// countCards counts card-like children in a grid
func countCards(comp *types.Component) int {
	// Count direct children that look like cards
	count := 0
	for i := range comp.Children {
		child := &comp.Children[i]
		if child.Type == "box" && child.Layout.Border != "" {
			count++
		} else {
			count++
		}
	}
	return count
}
