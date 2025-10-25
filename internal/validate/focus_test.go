package validate

import (
	"testing"

	"github.com/johanbellander/prism/internal/types"
)

func TestDefaultFocusRule(t *testing.T) {
	rule := DefaultFocusRule()

	if !rule.RequireFocusState {
		t.Error("Expected RequireFocusState to be true")
	}
	if rule.MinOutlineWidth != 2 {
		t.Errorf("Expected MinOutlineWidth 2, got %d", rule.MinOutlineWidth)
	}
	if rule.MinContrastRatio != 3.0 {
		t.Errorf("Expected MinContrastRatio 3.0, got %f", rule.MinContrastRatio)
	}
	if len(rule.InteractiveTypes) != 2 {
		t.Errorf("Expected 2 interactive types, got %d", len(rule.InteractiveTypes))
	}
	if !rule.RequireVisibleFocus {
		t.Error("Expected RequireVisibleFocus to be true")
	}
}

func TestValidateFocus_NoInteractiveElements(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "text-1",
				Type: "text",
			},
			{
				ID:   "image-1",
				Type: "image",
			},
		},
	}

	result := ValidateFocus(structure, DefaultFocusRule())

	if !result.Passed {
		t.Error("Expected validation to pass with no interactive elements")
	}
	if len(result.Issues) != 0 {
		t.Errorf("Expected no issues, got %d", len(result.Issues))
	}
}

func TestValidateFocus_ButtonElement(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "btn-1",
				Type: "button",
			},
		},
	}

	result := ValidateFocus(structure, DefaultFocusRule())

	if !result.Passed {
		t.Error("Expected validation to pass (info only)")
	}
	if len(result.Issues) == 0 {
		t.Error("Expected info message about focus state")
	}

	// Should have info severity
	if result.Issues[0].Severity != "info" {
		t.Errorf("Expected info severity, got %s", result.Issues[0].Severity)
	}
	if result.Issues[0].ComponentID != "btn-1" {
		t.Errorf("Expected component ID 'btn-1', got %s", result.Issues[0].ComponentID)
	}
}

func TestValidateFocus_InputElement(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "input-1",
				Type: "input",
			},
		},
	}

	result := ValidateFocus(structure, DefaultFocusRule())

	if !result.Passed {
		t.Error("Expected validation to pass (info only)")
	}
	if len(result.Issues) == 0 {
		t.Error("Expected info message about focus state")
	}

	// Should have info severity
	if result.Issues[0].Severity != "info" {
		t.Errorf("Expected info severity, got %s", result.Issues[0].Severity)
	}
	if result.Issues[0].ComponentID != "input-1" {
		t.Errorf("Expected component ID 'input-1', got %s", result.Issues[0].ComponentID)
	}
}

func TestValidateFocus_MultipleInteractiveElements(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "btn-1",
				Type: "button",
			},
			{
				ID:   "input-1",
				Type: "input",
			},
			{
				ID:   "text-1",
				Type: "text",
			},
			{
				ID:   "btn-2",
				Type: "button",
			},
		},
	}

	result := ValidateFocus(structure, DefaultFocusRule())

	if !result.Passed {
		t.Error("Expected validation to pass (info only)")
	}

	// Should have 3 info messages (2 buttons + 1 input)
	if len(result.Issues) != 3 {
		t.Errorf("Expected 3 issues, got %d", len(result.Issues))
	}

	// All should be info severity
	for _, issue := range result.Issues {
		if issue.Severity != "info" {
			t.Errorf("Expected info severity, got %s", issue.Severity)
		}
	}
}

func TestValidateFocus_NestedInteractiveElements(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "container",
				Type: "box",
				Children: []types.Component{
					{
						ID:   "nested-btn",
						Type: "button",
					},
					{
						ID:   "nested-input",
						Type: "input",
					},
				},
			},
		},
	}

	result := ValidateFocus(structure, DefaultFocusRule())

	if !result.Passed {
		t.Error("Expected validation to pass (info only)")
	}

	// Should have 2 info messages for nested interactive elements
	if len(result.Issues) != 2 {
		t.Errorf("Expected 2 issues, got %d", len(result.Issues))
	}

	// Check component IDs
	foundBtn := false
	foundInput := false
	for _, issue := range result.Issues {
		if issue.ComponentID == "nested-btn" {
			foundBtn = true
		}
		if issue.ComponentID == "nested-input" {
			foundInput = true
		}
	}
	if !foundBtn {
		t.Error("Expected issue for nested-btn")
	}
	if !foundInput {
		t.Error("Expected issue for nested-input")
	}
}

func TestValidateFocus_CustomRule_NoRequirement(t *testing.T) {
	customRule := FocusRule{
		RequireFocusState:   false, // Don't require focus states
		MinOutlineWidth:     2,
		MinContrastRatio:    3.0,
		InteractiveTypes:    []string{"button", "input"},
		RequireVisibleFocus: true,
	}

	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "btn-1",
				Type: "button",
			},
		},
	}

	result := ValidateFocus(structure, customRule)

	if !result.Passed {
		t.Error("Expected validation to pass")
	}
	// Should have no issues since we don't require focus states
	if len(result.Issues) != 0 {
		t.Errorf("Expected no issues, got %d", len(result.Issues))
	}
}

func TestValidateFocus_CustomRule_CustomInteractiveTypes(t *testing.T) {
	customRule := FocusRule{
		RequireFocusState:   true,
		MinOutlineWidth:     2,
		MinContrastRatio:    3.0,
		InteractiveTypes:    []string{"text"}, // Only text is interactive in this custom rule
		RequireVisibleFocus: true,
	}

	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "btn-1",
				Type: "button",
			},
			{
				ID:   "text-1",
				Type: "text",
			},
		},
	}

	result := ValidateFocus(structure, customRule)

	if !result.Passed {
		t.Error("Expected validation to pass (info only)")
	}

	// Should only have 1 issue for text-1
	if len(result.Issues) != 1 {
		t.Errorf("Expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].ComponentID != "text-1" {
		t.Errorf("Expected component ID 'text-1', got %s", result.Issues[0].ComponentID)
	}
}

func TestValidateFocus_DeepNesting(t *testing.T) {
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
								ID:   "deep-btn",
								Type: "button",
							},
						},
					},
				},
			},
		},
	}

	result := ValidateFocus(structure, DefaultFocusRule())

	if !result.Passed {
		t.Error("Expected validation to pass (info only)")
	}

	// Should find the deeply nested button
	if len(result.Issues) != 1 {
		t.Errorf("Expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].ComponentID != "deep-btn" {
		t.Errorf("Expected component ID 'deep-btn', got %s", result.Issues[0].ComponentID)
	}
}
