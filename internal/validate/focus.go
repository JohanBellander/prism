package validate

import (
	"fmt"

	"github.com/johanbellander/prism/internal/types"
)

// FocusIssue represents a focus indicator validation issue
type FocusIssue struct {
	ComponentID string `json:"component_id"`
	Message     string `json:"message"`
	Severity    string `json:"severity"` // "error", "warning", "info"
}

// FocusResult contains the validation results
type FocusResult struct {
	Passed bool         `json:"passed"`
	Issues []FocusIssue `json:"issues"`
}

// FocusRule defines the focus indicator validation rules
type FocusRule struct {
	RequireFocusState    bool     // Whether focus state is required for interactive elements
	MinOutlineWidth      int      // Minimum outline width in pixels (default: 2)
	MinContrastRatio     float64  // Minimum contrast ratio for focus indicator (default: 3.0)
	InteractiveTypes     []string // Component types that require focus indicators
	RequireVisibleFocus  bool     // Whether focus must be visibly different from default state
}

// DefaultFocusRule returns the default focus indicator validation rules
func DefaultFocusRule() FocusRule {
	return FocusRule{
		RequireFocusState:   true,
		MinOutlineWidth:     2,
		MinContrastRatio:    3.0,
		InteractiveTypes:    []string{"button", "input"},
		RequireVisibleFocus: true,
	}
}

// ValidateFocus validates focus indicators on interactive elements
func ValidateFocus(structure *types.Structure, rule FocusRule) FocusResult {
	result := FocusResult{
		Passed: true,
		Issues: []FocusIssue{},
	}

	// Check all components
	for _, component := range structure.Components {
		validateComponentFocus(&result, &component, rule)
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

func validateComponentFocus(result *FocusResult, component *types.Component, rule FocusRule) {
	// Check if this is an interactive component
	isInteractive := false
	for _, interactiveType := range rule.InteractiveTypes {
		if component.Type == interactiveType {
			isInteractive = true
			break
		}
	}

	if isInteractive && rule.RequireFocusState {
		// For Phase 1, we don't have explicit focus state in the schema yet
		// This is more of a documentation/reminder validator
		// In a real implementation, we'd check for focus state properties
		
		// Add informational message about focus states
		result.Issues = append(result.Issues, FocusIssue{
			ComponentID: component.ID,
			Message:     fmt.Sprintf("Interactive element '%s' of type '%s' should define a visible focus state for keyboard navigation (WCAG 2.4.7)", component.ID, component.Type),
			Severity:    "info",
		})
	}

	// Check children recursively
	for _, child := range component.Children {
		validateComponentFocus(result, &child, rule)
	}
}
