package validate

import (
	"fmt"

	"github.com/johanbellander/prism/internal/types"
)

// LoadingStateRule defines the rules for loading state validation
type LoadingStateRule struct {
	ValidStates        []string // Valid state values
	RequireSkeleton    bool     // Require skeleton config for loading state
	RequireEmptyMessage bool    // Require message for empty state
}

// LoadingStateIssue represents a loading state validation issue
type LoadingStateIssue struct {
	ComponentID string `json:"component_id"`
	Message     string `json:"message"`
	Severity    string `json:"severity"` // "error", "warning", "info"
}

// LoadingStateResult represents the result of loading state validation
type LoadingStateResult struct {
	Passed bool                `json:"passed"`
	Issues []LoadingStateIssue `json:"issues"`
}

// DefaultLoadingStateRule returns the default loading state rule
func DefaultLoadingStateRule() LoadingStateRule {
	return LoadingStateRule{
		ValidStates:        []string{"loading", "error", "empty", "default", ""},
		RequireSkeleton:    false, // Optional but recommended
		RequireEmptyMessage: false, // Optional but recommended
	}
}

// ValidateLoadingStates validates that components use proper loading states
func ValidateLoadingStates(structure *types.Structure, rule LoadingStateRule) LoadingStateResult {
	result := LoadingStateResult{
		Passed: true,
		Issues: []LoadingStateIssue{},
	}

	// Validate all components recursively
	validateComponentStates(structure.Components, rule, &result)

	return result
}

func validateComponentStates(components []types.Component, rule LoadingStateRule, result *LoadingStateResult) {
	for _, comp := range components {
		// Check if state is valid
		if comp.State != "" && !isValidState(comp.State, rule.ValidStates) {
			result.Passed = false
			result.Issues = append(result.Issues, LoadingStateIssue{
				ComponentID: comp.ID,
				Message:     fmt.Sprintf("Loading State: '%s' has invalid state '%s'", comp.ID, comp.State),
				Severity:    "error",
			})
			
			result.Issues = append(result.Issues, LoadingStateIssue{
				ComponentID: comp.ID,
				Message:     fmt.Sprintf("   Valid states: %v", rule.ValidStates),
				Severity:    "info",
			})
		}

		// Check for loading state recommendations
		if comp.State == "loading" {
			if comp.Skeleton == nil {
				result.Issues = append(result.Issues, LoadingStateIssue{
					ComponentID: comp.ID,
					Message:     fmt.Sprintf("Loading State: '%s' in loading state but missing skeleton configuration", comp.ID),
					Severity:    "info",
				})
			} else {
				// Validate skeleton configuration
				validateSkeleton(comp, result)
			}
		}

		// Check for empty state
		if comp.State == "empty" {
			if comp.Content == "" && len(comp.Children) == 0 {
				result.Issues = append(result.Issues, LoadingStateIssue{
					ComponentID: comp.ID,
					Message:     fmt.Sprintf("Loading State: '%s' in empty state - consider adding empty state message", comp.ID),
					Severity:    "info",
				})
			}
		}

		// Check for error state
		if comp.State == "error" {
			if comp.Content == "" && len(comp.Children) == 0 {
				result.Issues = append(result.Issues, LoadingStateIssue{
					ComponentID: comp.ID,
					Message:     fmt.Sprintf("Loading State: '%s' in error state - consider adding error message", comp.ID),
					Severity:    "info",
				})
			}
		}

		// Recursively validate children
		if len(comp.Children) > 0 {
			validateComponentStates(comp.Children, rule, result)
		}
	}
}

func isValidState(state string, validStates []string) bool {
	for _, valid := range validStates {
		if state == valid {
			return true
		}
	}
	return false
}

func validateSkeleton(comp types.Component, result *LoadingStateResult) {
	if comp.Skeleton == nil {
		return
	}

	if len(comp.Skeleton.Elements) == 0 {
		result.Issues = append(result.Issues, LoadingStateIssue{
			ComponentID: comp.ID,
			Message:     fmt.Sprintf("Loading State: '%s' has skeleton config but no elements defined", comp.ID),
			Severity:    "warning",
		})
		return
	}

	// Validate skeleton elements
	for i, elem := range comp.Skeleton.Elements {
		if elem.Type == "" {
			result.Passed = false
			result.Issues = append(result.Issues, LoadingStateIssue{
				ComponentID: comp.ID,
				Message:     fmt.Sprintf("Loading State: '%s' skeleton element %d missing type", comp.ID, i),
				Severity:    "error",
			})
		}

		if !isValidSkeletonType(elem.Type) {
			result.Passed = false
			result.Issues = append(result.Issues, LoadingStateIssue{
				ComponentID: comp.ID,
				Message:     fmt.Sprintf("Loading State: '%s' skeleton element %d has invalid type '%s'", comp.ID, i, elem.Type),
				Severity:    "error",
			})
			
			result.Issues = append(result.Issues, LoadingStateIssue{
				ComponentID: comp.ID,
				Message:     "   Valid skeleton types: circle, text, rect",
				Severity:    "info",
			})
		}

		// Check for required dimensions
		if elem.Type == "circle" && elem.Size == 0 {
			result.Issues = append(result.Issues, LoadingStateIssue{
				ComponentID: comp.ID,
				Message:     fmt.Sprintf("Loading State: '%s' skeleton circle element %d should specify size", comp.ID, i),
				Severity:    "warning",
			})
		}

		if (elem.Type == "text" || elem.Type == "rect") && elem.Width == "" {
			result.Issues = append(result.Issues, LoadingStateIssue{
				ComponentID: comp.ID,
				Message:     fmt.Sprintf("Loading State: '%s' skeleton %s element %d should specify width", comp.ID, elem.Type, i),
				Severity:    "warning",
			})
		}
	}
}

func isValidSkeletonType(skeletonType string) bool {
	validTypes := []string{"circle", "text", "rect"}
	for _, valid := range validTypes {
		if skeletonType == valid {
			return true
		}
	}
	return false
}

// CountComponentsByState counts components in each state
func CountComponentsByState(structure *types.Structure) map[string]int {
	counts := make(map[string]int)
	countStates(structure.Components, counts)
	return counts
}

func countStates(components []types.Component, counts map[string]int) {
	for _, comp := range components {
		state := comp.State
		if state == "" {
			state = "default"
		}
		counts[state]++
		
		if len(comp.Children) > 0 {
			countStates(comp.Children, counts)
		}
	}
}
