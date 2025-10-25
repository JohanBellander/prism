package validate

import (
	"testing"
	"time"

	"github.com/johanbellander/prism/internal/types"
)

func TestValidateTouchTargets_MinimumSize(t *testing.T) {
	// Create a structure with touch target size issues
	structure := &types.Structure{
		Version:   "v1",
		Phase:     "structure",
		CreatedAt: time.Now(),
		Intent: types.Intent{
			Purpose:       "Test",
			PrimaryAction: "test",
		},
		Layout: types.Layout{
			Type:      "stack",
			Direction: "vertical",
			Spacing:   16,
		},
		Components: []types.Component{
			{
				ID:   "small-button",
				Type: "button",
				Layout: types.ComponentLayout{
					Display: "block",
					Width:   32,  // Too small
					Height:  32,  // Too small
				},
			},
			{
				ID:   "good-button",
				Type: "button",
				Layout: types.ComponentLayout{
					Display: "block",
					Width:   48,
					Height:  48,
				},
			},
		},
	}

	rule := DefaultTouchTargetRule()
	result := ValidateTouchTargets(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail due to small touch target")
	}

	// Should have at least one error about the small button
	errorCount := 0
	for _, issue := range result.Issues {
		if issue.Severity == "error" && issue.Component == "small-button" {
			errorCount++
		}
	}

	if errorCount == 0 {
		t.Error("Expected at least one error about small-button touch target size")
	}
}

func TestValidateTouchTargets_DangerousActionSpacing(t *testing.T) {
	// Create a structure with dangerous actions too close to other buttons
	structure := &types.Structure{
		Version:   "v1",
		Phase:     "structure",
		CreatedAt: time.Now(),
		Intent: types.Intent{
			Purpose:       "Test",
			PrimaryAction: "test",
		},
		Layout: types.Layout{
			Type:      "stack",
			Direction: "vertical",
			Spacing:   4, // Very small spacing
		},
		Components: []types.Component{
			{
				ID:   "delete-btn",
				Type: "button",
				Layout: types.ComponentLayout{
					Display: "block",
					Width:   100,
					Height:  44,
				},
			},
			{
				ID:   "cancel-btn",
				Type: "button",
				Layout: types.ComponentLayout{
					Display: "block",
					Width:   100,
					Height:  44,
				},
			},
		},
	}

	rule := DefaultTouchTargetRule()
	result := ValidateTouchTargets(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail due to dangerous action spacing")
	}

	// Should have an error about spacing for destructive actions
	foundSpacingError := false
	for _, issue := range result.Issues {
		if issue.Severity == "error" && issue.Component == "delete-btn" {
			foundSpacingError = true
			break
		}
	}

	if !foundSpacingError {
		t.Error("Expected an error about inadequate spacing for dangerous action")
	}
}

func TestValidateTouchTargets_ValidStructure(t *testing.T) {
	// Create a valid structure
	structure := &types.Structure{
		Version:   "v1",
		Phase:     "structure",
		CreatedAt: time.Now(),
		Intent: types.Intent{
			Purpose:       "Test",
			PrimaryAction: "submit",
		},
		Layout: types.Layout{
			Type:      "stack",
			Direction: "vertical",
			Spacing:   24,
		},
		Components: []types.Component{
			{
				ID:   "submit",
				Type: "button",
				Layout: types.ComponentLayout{
					Display: "block",
					Width:   150,
					Height:  48,
				},
			},
			{
				ID:   "cancel",
				Type: "button",
				Layout: types.ComponentLayout{
					Display: "block",
					Width:   120,
					Height:  48,
				},
			},
		},
	}

	rule := DefaultTouchTargetRule()
	result := ValidateTouchTargets(structure, rule)

	if !result.Passed {
		t.Error("Expected validation to pass for valid structure")
	}

	// Should only have info messages
	for _, issue := range result.Issues {
		if issue.Severity == "error" || issue.Severity == "warning" {
			t.Errorf("Unexpected issue in valid structure: %s - %s", issue.Severity, issue.Message)
		}
	}
}

func TestIsDangerousAction(t *testing.T) {
	tests := []struct {
		id       string
		role     string
		expected bool
	}{
		{"delete-btn", "", true},
		{"remove-item", "", true},
		{"cancel-action", "", true},
		{"submit-form", "", false},
		{"save-btn", "", false},
		{"", "delete", true},
		{"", "primary", false},
	}

	for _, test := range tests {
		comp := &types.Component{
			ID:   test.id,
			Role: test.role,
			Type: "button",
		}

		result := isDangerousAction(comp)
		if result != test.expected {
			t.Errorf("isDangerousAction(%s, %s) = %v, expected %v", test.id, test.role, result, test.expected)
		}
	}
}

func TestIsInteractiveElement(t *testing.T) {
	tests := []struct {
		compType string
		expected bool
	}{
		{"button", true},
		{"input", true},
		{"text", false},
		{"box", false},
		{"image", false},
	}

	for _, test := range tests {
		comp := &types.Component{
			ID:   "test",
			Type: test.compType,
		}

		result := isInteractiveElement(comp)
		if result != test.expected {
			t.Errorf("isInteractiveElement(%s) = %v, expected %v", test.compType, result, test.expected)
		}
	}
}
