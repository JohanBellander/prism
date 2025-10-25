package validate

import (
	"testing"

	"github.com/johanbellander/prism/internal/types"
)

func TestDefaultDarkModeRule(t *testing.T) {
	rule := DefaultDarkModeRule()

	if !rule.RequireSemanticColors {
		t.Error("Expected RequireSemanticColors to be true")
	}
	if !rule.ValidateContrast {
		t.Error("Expected ValidateContrast to be true")
	}
	if rule.MinContrastRatio != 4.5 {
		t.Errorf("Expected MinContrastRatio 4.5, got %f", rule.MinContrastRatio)
	}
	if !rule.RecommendAdaptive {
		t.Error("Expected RecommendAdaptive to be true")
	}
}

func TestValidateDarkMode_BasicRecommendation(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "text-1",
				Type: "text",
			},
		},
	}

	result := ValidateDarkMode(structure, DefaultDarkModeRule())

	if !result.Passed {
		t.Error("Expected validation to pass (info only)")
	}
	if len(result.Issues) == 0 {
		t.Error("Expected at least one info message about semantic colors")
	}

	// Should have info about semantic colors
	foundSemanticInfo := false
	for _, issue := range result.Issues {
		if issue.ComponentID == "structure" && issue.Severity == "info" {
			foundSemanticInfo = true
			break
		}
	}
	if !foundSemanticInfo {
		t.Error("Expected info message about semantic color tokens")
	}
}

func TestValidateDarkMode_HardcodedBlackText(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:    "text-1",
				Type:  "text",
				Color: "#000000",
			},
		},
	}

	result := ValidateDarkMode(structure, DefaultDarkModeRule())

	if !result.Passed {
		t.Error("Expected validation to pass (info only)")
	}

	// Should warn about hardcoded black color
	foundColorWarning := false
	for _, issue := range result.Issues {
		if issue.ComponentID == "text-1" && issue.Severity == "info" {
			foundColorWarning = true
			break
		}
	}
	if !foundColorWarning {
		t.Error("Expected info message about hardcoded black color")
	}
}

func TestValidateDarkMode_HardcodedWhiteText(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:    "text-1",
				Type:  "text",
				Color: "#FFFFFF",
			},
		},
	}

	result := ValidateDarkMode(structure, DefaultDarkModeRule())

	if !result.Passed {
		t.Error("Expected validation to pass (info only)")
	}

	// Should warn about hardcoded white color
	foundColorWarning := false
	for _, issue := range result.Issues {
		if issue.ComponentID == "text-1" && issue.Severity == "info" {
			foundColorWarning = true
			break
		}
	}
	if !foundColorWarning {
		t.Error("Expected info message about hardcoded white color")
	}
}

func TestValidateDarkMode_HardcodedWhiteBackground(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "box-1",
				Type: "box",
				Layout: types.ComponentLayout{
					Background: "#FFFFFF",
				},
			},
		},
	}

	result := ValidateDarkMode(structure, DefaultDarkModeRule())

	if !result.Passed {
		t.Error("Expected validation to pass (info only)")
	}

	// Should warn about hardcoded white background
	foundBgWarning := false
	for _, issue := range result.Issues {
		if issue.ComponentID == "box-1" && issue.Severity == "info" {
			foundBgWarning = true
			break
		}
	}
	if !foundBgWarning {
		t.Error("Expected info message about hardcoded white background")
	}
}

func TestValidateDarkMode_HardcodedBlackBackground(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "box-1",
				Type: "box",
				Layout: types.ComponentLayout{
					Background: "#000000",
				},
			},
		},
	}

	result := ValidateDarkMode(structure, DefaultDarkModeRule())

	if !result.Passed {
		t.Error("Expected validation to pass (info only)")
	}

	// Should warn about hardcoded black background
	foundBgWarning := false
	for _, issue := range result.Issues {
		if issue.ComponentID == "box-1" && issue.Severity == "info" {
			foundBgWarning = true
			break
		}
	}
	if !foundBgWarning {
		t.Error("Expected info message about hardcoded black background")
	}
}

func TestValidateDarkMode_AdaptiveColors(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:    "text-1",
				Type:  "text",
				Color: "#1F2937", // Dark gray, not pure black/white
			},
		},
	}

	result := ValidateDarkMode(structure, DefaultDarkModeRule())

	if !result.Passed {
		t.Error("Expected validation to pass")
	}

	// Should only have the general semantic colors recommendation
	// No specific warnings about this color
	specificWarnings := 0
	for _, issue := range result.Issues {
		if issue.ComponentID == "text-1" {
			specificWarnings++
		}
	}
	if specificWarnings > 0 {
		t.Error("Should not warn about adaptive color values")
	}
}

func TestValidateDarkMode_MultipleComponents(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:    "text-1",
				Type:  "text",
				Color: "#000000",
			},
			{
				ID:   "box-1",
				Type: "box",
				Layout: types.ComponentLayout{
					Background: "#FFFFFF",
				},
			},
			{
				ID:    "text-2",
				Type:  "text",
				Color: "#737373", // Adaptive color
			},
		},
	}

	result := ValidateDarkMode(structure, DefaultDarkModeRule())

	if !result.Passed {
		t.Error("Expected validation to pass (info only)")
	}

	// Should have warnings for text-1 and box-1, but not text-2
	text1Warning := false
	box1Warning := false
	text2Warning := false

	for _, issue := range result.Issues {
		if issue.ComponentID == "text-1" {
			text1Warning = true
		}
		if issue.ComponentID == "box-1" {
			box1Warning = true
		}
		if issue.ComponentID == "text-2" {
			text2Warning = true
		}
	}

	if !text1Warning {
		t.Error("Expected warning for text-1 with hardcoded black")
	}
	if !box1Warning {
		t.Error("Expected warning for box-1 with hardcoded white background")
	}
	if text2Warning {
		t.Error("Should not warn about text-2 with adaptive color")
	}
}

func TestValidateDarkMode_NestedComponents(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "container",
				Type: "box",
				Children: []types.Component{
					{
						ID:    "nested-text",
						Type:  "text",
						Color: "#FFFFFF",
					},
				},
			},
		},
	}

	result := ValidateDarkMode(structure, DefaultDarkModeRule())

	if !result.Passed {
		t.Error("Expected validation to pass (info only)")
	}

	// Should find the nested component with hardcoded color
	foundNestedWarning := false
	for _, issue := range result.Issues {
		if issue.ComponentID == "nested-text" {
			foundNestedWarning = true
			break
		}
	}
	if !foundNestedWarning {
		t.Error("Expected warning for nested component with hardcoded white")
	}
}

func TestValidateDarkMode_CustomRule_NoRecommendations(t *testing.T) {
	customRule := DarkModeRule{
		RequireSemanticColors: false,
		ValidateContrast:      false,
		MinContrastRatio:      4.5,
		RecommendAdaptive:     false,
	}

	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:    "text-1",
				Type:  "text",
				Color: "#000000",
			},
		},
	}

	result := ValidateDarkMode(structure, customRule)

	if !result.Passed {
		t.Error("Expected validation to pass")
	}

	// Should have no issues with custom rule disabling recommendations
	if len(result.Issues) != 0 {
		t.Errorf("Expected no issues with disabled recommendations, got %d", len(result.Issues))
	}
}

func TestValidateDarkMode_DeepNesting(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "level-1",
				Type: "box",
				Children: []types.Component{
					{
						ID:   "level-2",
						Type: "box",
						Children: []types.Component{
							{
								ID:    "deep-text",
								Type:  "text",
								Color: "#000000",
							},
						},
					},
				},
			},
		},
	}

	result := ValidateDarkMode(structure, DefaultDarkModeRule())

	if !result.Passed {
		t.Error("Expected validation to pass (info only)")
	}

	// Should find the deeply nested component
	foundDeepWarning := false
	for _, issue := range result.Issues {
		if issue.ComponentID == "deep-text" {
			foundDeepWarning = true
			break
		}
	}
	if !foundDeepWarning {
		t.Error("Expected warning for deeply nested component")
	}
}
