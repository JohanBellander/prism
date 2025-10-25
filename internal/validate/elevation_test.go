package validate

import (
	"testing"

	"github.com/johanbellander/prism/internal/types"
)

func TestDefaultElevationRule(t *testing.T) {
	rule := DefaultElevationRule()
	
	expectedLevels := map[string]string{
		"0": "none",
		"1": "0 1px 2px 0 rgba(0,0,0,0.05)",
		"2": "0 2px 4px 0 rgba(0,0,0,0.1)",
		"3": "0 4px 8px 0 rgba(0,0,0,0.12)",
		"4": "0 8px 16px 0 rgba(0,0,0,0.15)",
		"5": "0 16px 32px 0 rgba(0,0,0,0.2)",
	}
	
	for level, expectedShadow := range expectedLevels {
		actualShadow, exists := rule.Levels[level]
		if !exists {
			t.Errorf("Missing elevation level '%s'", level)
			continue
		}
		if actualShadow != expectedShadow {
			t.Errorf("Level '%s': expected '%s', got '%s'", level, expectedShadow, actualShadow)
		}
	}
}

func TestValidateElevation_EmptyStructure(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{},
	}
	
	rule := DefaultElevationRule()
	result := ValidateElevation(structure, rule)
	
	if !result.Passed {
		t.Errorf("Expected validation to pass for empty structure")
	}
	
	if len(result.Issues) > 0 {
		t.Errorf("Expected no issues for empty structure, got %d", len(result.Issues))
	}
}

func TestValidateElevation_ComponentRecommendations(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "card1",
				Type: "box",
				Role: "card",
			},
			{
				ID:   "button1",
				Type: "button",
			},
			{
				ID:   "modal1",
				Type: "box",
				Role: "modal",
			},
		},
	}
	
	rule := DefaultElevationRule()
	result := ValidateElevation(structure, rule)
	
	// Should have info messages for recommendations
	if len(result.Issues) == 0 {
		t.Errorf("Expected recommendation info messages")
	}
	
	// Check for specific recommendations
	foundCardRec := false
	foundButtonRec := false
	foundModalRec := false
	
	for _, issue := range result.Issues {
		if issue.ComponentID == "card1" && issue.Severity == "info" {
			foundCardRec = true
		}
		if issue.ComponentID == "button1" && issue.Severity == "info" {
			foundButtonRec = true
		}
		if issue.ComponentID == "modal1" && issue.Severity == "info" {
			foundModalRec = true
		}
	}
	
	if !foundCardRec {
		t.Errorf("Expected recommendation for card component")
	}
	if !foundButtonRec {
		t.Errorf("Expected recommendation for button component")
	}
	if !foundModalRec {
		t.Errorf("Expected recommendation for modal component")
	}
}

func TestGetRecommendedElevationLevel(t *testing.T) {
	tests := []struct {
		componentType string
		role          string
		expected      string
		name          string
	}{
		{"box", "card", "1", "card"},
		{"button", "", "2", "button"},
		{"box", "modal", "4", "modal"},
		{"box", "dialog", "4", "dialog"},
		{"box", "dropdown", "3", "dropdown"},
		{"box", "menu", "3", "menu"},
		{"text", "", "", "text (no recommendation)"},
		{"box", "content", "", "box content (no recommendation)"},
	}
	
	for _, tt := range tests {
		result := getRecommendedElevationLevel(tt.componentType, tt.role)
		if result != tt.expected {
			t.Errorf("%s: expected '%s', got '%s'", tt.name, tt.expected, result)
		}
	}
}

func TestParseShadowValue(t *testing.T) {
	rule := DefaultElevationRule()
	
	tests := []struct {
		shadow        string
		expectedLevel string
		shouldMatch   bool
		name          string
	}{
		{"none", "0", true, "none matches level 0"},
		{"0 1px 2px 0 rgba(0,0,0,0.05)", "1", true, "exact match level 1"},
		{"0 2px 4px 0 rgba(0,0,0,0.1)", "2", true, "exact match level 2"},
		{"0 4px 8px 0 rgba(0,0,0,0.12)", "3", true, "exact match level 3"},
		{"  0 4px 8px 0 rgba(0,0,0,0.12)  ", "3", true, "whitespace trimmed"},
		{"0 4PX 8PX 0 RGBA(0,0,0,0.12)", "3", true, "case insensitive"},
		{"0 3px 6px 0 rgba(0,0,0,0.15)", "", false, "custom shadow doesn't match"},
		{"", "", false, "empty shadow"},
	}
	
	for _, tt := range tests {
		level, matches := ParseShadowValue(tt.shadow, rule)
		if matches != tt.shouldMatch {
			t.Errorf("%s: expected match=%v, got %v", tt.name, tt.shouldMatch, matches)
		}
		if matches && level != tt.expectedLevel {
			t.Errorf("%s: expected level '%s', got '%s'", tt.name, tt.expectedLevel, level)
		}
	}
}

func TestNormalizeShadow(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		name     string
	}{
		{"0 1px 2px 0 rgba(0,0,0,0.05)", "0 1px 2px 0 rgba(0,0,0,0.05)", "no change needed"},
		{"  0 1px 2px 0 rgba(0,0,0,0.05)  ", "0 1px 2px 0 rgba(0,0,0,0.05)", "trim whitespace"},
		{"0  1px  2px  0  rgba(0,0,0,0.05)", "0 1px 2px 0 rgba(0,0,0,0.05)", "normalize spaces"},
		{"0 1PX 2PX 0 RGBA(0,0,0,0.05)", "0 1px 2px 0 rgba(0,0,0,0.05)", "lowercase"},
		{"NONE", "none", "none lowercase"},
	}
	
	for _, tt := range tests {
		result := normalizeShadow(tt.input)
		if result != tt.expected {
			t.Errorf("%s: expected '%s', got '%s'", tt.name, tt.expected, result)
		}
	}
}

func TestValidateShadowValue(t *testing.T) {
	rule := DefaultElevationRule()
	
	tests := []struct {
		shadow          string
		expectedValid   bool
		expectedLevel   string
		hasSuggestion   bool
		name            string
	}{
		{"none", true, "0", false, "none is valid"},
		{"", true, "0", false, "empty is valid (no shadow)"},
		{"0 1px 2px 0 rgba(0,0,0,0.05)", true, "1", false, "level 1 is valid"},
		{"0 8px 16px 0 rgba(0,0,0,0.15)", true, "4", false, "level 4 is valid"},
		{"0 3px 6px 0 rgba(0,0,0,0.08)", false, "", true, "custom shadow gets suggestion"},
	}
	
	for _, tt := range tests {
		valid, level, suggestion := ValidateShadowValue(tt.shadow, rule)
		if valid != tt.expectedValid {
			t.Errorf("%s: expected valid=%v, got %v", tt.name, tt.expectedValid, valid)
		}
		if valid && level != tt.expectedLevel {
			t.Errorf("%s: expected level '%s', got '%s'", tt.name, tt.expectedLevel, level)
		}
		if tt.hasSuggestion && suggestion == "" {
			t.Errorf("%s: expected suggestion but got none", tt.name)
		}
		if !tt.hasSuggestion && suggestion != "" {
			t.Errorf("%s: expected no suggestion but got '%s'", tt.name, suggestion)
		}
	}
}

func TestExtractBlurRadius(t *testing.T) {
	tests := []struct {
		shadow   string
		expected int
		name     string
	}{
		{"0 1px 2px 0 rgba(0,0,0,0.05)", 2, "level 1"},
		{"0 2px 4px 0 rgba(0,0,0,0.1)", 4, "level 2"},
		{"0 4px 8px 0 rgba(0,0,0,0.12)", 8, "level 3"},
		{"0 8px 16px 0 rgba(0,0,0,0.15)", 16, "level 4"},
		{"0 16px 32px 0 rgba(0,0,0,0.2)", 32, "level 5"},
		{"invalid", 0, "invalid format"},
		{"none", 0, "none"},
	}
	
	for _, tt := range tests {
		result := extractBlurRadius(tt.shadow)
		if result != tt.expected {
			t.Errorf("%s: expected %d, got %d", tt.name, tt.expected, result)
		}
	}
}

func TestFindClosestElevationLevel(t *testing.T) {
	rule := DefaultElevationRule()
	
	tests := []struct {
		shadow   string
		expected string
		name     string
	}{
		{"0 1px 1px 0 rgba(0,0,0,0.05)", "1", "blur 1px -> level 1"},
		{"0 2px 3px 0 rgba(0,0,0,0.1)", "2", "blur 3px -> level 2"},
		{"0 3px 6px 0 rgba(0,0,0,0.12)", "3", "blur 6px -> level 3"},
		{"0 5px 12px 0 rgba(0,0,0,0.15)", "4", "blur 12px -> level 4"},
		{"0 10px 40px 0 rgba(0,0,0,0.2)", "5", "blur 40px -> level 5"},
	}
	
	for _, tt := range tests {
		result := findClosestElevationLevel(tt.shadow, rule)
		if result != tt.expected {
			t.Errorf("%s: expected '%s', got '%s'", tt.name, tt.expected, result)
		}
	}
}

func TestValidateElevation_NestedComponents(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "header",
				Type: "box",
				Role: "header",
				Children: []types.Component{
					{
						ID:   "logo-card",
						Type: "box",
						Role: "card",
					},
				},
			},
			{
				ID:   "content",
				Type: "box",
				Role: "content",
				Children: []types.Component{
					{
						ID:   "submit-button",
						Type: "button",
					},
				},
			},
		},
	}
	
	rule := DefaultElevationRule()
	result := ValidateElevation(structure, rule)
	
	// Should have recommendations for nested components
	foundCardRec := false
	foundButtonRec := false
	
	for _, issue := range result.Issues {
		if issue.ComponentID == "logo-card" {
			foundCardRec = true
		}
		if issue.ComponentID == "submit-button" {
			foundButtonRec = true
		}
	}
	
	if !foundCardRec {
		t.Errorf("Expected recommendation for nested card component")
	}
	if !foundButtonRec {
		t.Errorf("Expected recommendation for nested button component")
	}
}
