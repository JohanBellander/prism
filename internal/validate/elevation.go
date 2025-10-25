package validate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/johanbellander/prism/internal/types"
)

// ElevationRule defines the rules for shadow/elevation validation
type ElevationRule struct {
	Levels map[string]string // elevation level -> shadow CSS value
}

// ElevationIssue represents an elevation validation issue
type ElevationIssue struct {
	ComponentID string `json:"component_id"`
	Message     string `json:"message"`
	Severity    string `json:"severity"` // "error", "warning", "info"
}

// ElevationResult represents the result of elevation validation
type ElevationResult struct {
	Passed bool              `json:"passed"`
	Issues []ElevationIssue  `json:"issues"`
}

// DefaultElevationRule returns the default elevation rule
func DefaultElevationRule() ElevationRule {
	return ElevationRule{
		Levels: map[string]string{
			"0": "none",
			"1": "0 1px 2px 0 rgba(0,0,0,0.05)",   // Subtle (cards)
			"2": "0 2px 4px 0 rgba(0,0,0,0.1)",    // Raised (buttons)
			"3": "0 4px 8px 0 rgba(0,0,0,0.12)",   // Floating (dropdowns)
			"4": "0 8px 16px 0 rgba(0,0,0,0.15)",  // Overlays (modals)
			"5": "0 16px 32px 0 rgba(0,0,0,0.2)",  // Maximum (important dialogs)
		},
	}
}

// ValidateElevation validates that components use consistent elevation/shadow values
func ValidateElevation(structure *types.Structure, rule ElevationRule) ElevationResult {
	result := ElevationResult{
		Passed: true,
		Issues: []ElevationIssue{},
	}

	// Validate all components recursively
	validateComponentElevation(structure.Components, rule, &result)

	return result
}

func validateComponentElevation(components []types.Component, rule ElevationRule, result *ElevationResult) {
	for _, comp := range components {
		// Check for shadow property in layout
		// Note: Shadow property doesn't exist in current Phase 1 schema
		// This validator is prepared for when it's added
		validateShadow(comp, rule, result)

		// Recursively validate children
		if len(comp.Children) > 0 {
			validateComponentElevation(comp.Children, rule, result)
		}
	}
}

func validateShadow(comp types.Component, rule ElevationRule, result *ElevationResult) {
	// For future implementation when shadow/elevation is added to schema
	// Currently this would check comp.Layout.Shadow or comp.Elevation
	// For now, we'll provide informational validation framework
	
	// Check component type for recommended elevation levels
	recommendedLevel := getRecommendedElevationLevel(comp.Type, comp.Role)
	if recommendedLevel != "" {
		result.Issues = append(result.Issues, ElevationIssue{
			ComponentID: comp.ID,
			Message:     fmt.Sprintf("Info: Component '%s' (%s) should use elevation %s: %s", 
				comp.ID, comp.Type, recommendedLevel, rule.Levels[recommendedLevel]),
			Severity:    "info",
		})
	}
}

func getRecommendedElevationLevel(componentType, role string) string {
	// Recommend elevation levels based on component type/role
	switch componentType {
	case "box":
		switch role {
		case "card":
			return "1"
		case "modal", "dialog":
			return "4"
		case "dropdown", "menu":
			return "3"
		}
	case "button":
		return "2"
	}
	return ""
}

// ParseShadowValue parses a CSS box-shadow value and returns elevation level if it matches
func ParseShadowValue(shadow string, rule ElevationRule) (string, bool) {
	shadow = strings.TrimSpace(shadow)
	
	// Check if it matches any predefined level
	for level, definedShadow := range rule.Levels {
		if normalizeShadow(shadow) == normalizeShadow(definedShadow) {
			return level, true
		}
	}
	
	return "", false
}

// normalizeShadow normalizes a shadow string for comparison
func normalizeShadow(shadow string) string {
	// Remove extra whitespace
	shadow = strings.TrimSpace(shadow)
	shadow = regexp.MustCompile(`\s+`).ReplaceAllString(shadow, " ")
	
	// Convert to lowercase for case-insensitive comparison
	shadow = strings.ToLower(shadow)
	
	return shadow
}

// ValidateShadowValue validates a shadow CSS value against the elevation system
func ValidateShadowValue(shadow string, rule ElevationRule) (bool, string, string) {
	if shadow == "" || shadow == "none" {
		return true, "0", ""
	}
	
	level, matches := ParseShadowValue(shadow, rule)
	if matches {
		return true, level, ""
	}
	
	// Find closest matching elevation level
	closestLevel := findClosestElevationLevel(shadow, rule)
	suggestion := fmt.Sprintf("Consider using elevation %s: %s", closestLevel, rule.Levels[closestLevel])
	
	return false, "", suggestion
}

func findClosestElevationLevel(shadow string, rule ElevationRule) string {
	// Extract blur radius from shadow as a simple heuristic
	blurRadius := extractBlurRadius(shadow)
	
	// Map blur radius to elevation level
	switch {
	case blurRadius <= 1:
		return "1"
	case blurRadius <= 3:
		return "2"
	case blurRadius <= 6:
		return "3"
	case blurRadius <= 12:
		return "4"
	default:
		return "5"
	}
}

func extractBlurRadius(shadow string) int {
	// Simple regex to extract blur radius (3rd number in box-shadow)
	// Format: offset-x offset-y blur-radius spread-radius color
	// Example: "0 4px 8px 0 rgba(0,0,0,0.12)"
	re := regexp.MustCompile(`(-?\d+)\s+(-?\d+)px\s+(-?\d+)px`)
	matches := re.FindStringSubmatch(shadow)
	
	if len(matches) >= 4 {
		blur, err := strconv.Atoi(matches[3])
		if err == nil {
			return blur
		}
	}
	
	return 0
}
