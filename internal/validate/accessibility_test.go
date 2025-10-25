package validate

import (
	"testing"
	"time"

	"github.com/johanbellander/prism/internal/types"
)

func TestValidateAccessibility_MissingLabels(t *testing.T) {
	// Create a structure with unlabeled interactive elements
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
				ID:   "search-input",
				Type: "input",
			},
		},
		Accessibility: types.Accessibility{
			TouchTargetsMin:   44,
			FocusIndicators:   "visible",
			Labels:            "", // Not all_interactive_elements
			SemanticStructure: true,
		},
	}

	rule := DefaultA11yRule()
	result := ValidateAccessibility(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail due to missing label")
	}

	foundMissingLabel := false
	for _, issue := range result.Issues {
		if issue.Severity == "error" && issue.Component == "search-input" {
			foundMissingLabel = true
			break
		}
	}

	if !foundMissingLabel {
		t.Error("Expected error about missing label for search-input")
	}
}

func TestValidateAccessibility_HeadingOrder(t *testing.T) {
	// Create a structure with skipped heading levels
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
				ID:   "h1",
				Type: "text",
				Size: "4xl",
			},
			{
				ID:   "h3", // Skipping h2
				Type: "text",
				Size: "2xl",
			},
		},
		Accessibility: types.Accessibility{
			TouchTargetsMin:   44,
			FocusIndicators:   "visible",
			Labels:            "all_interactive_elements",
			SemanticStructure: true,
		},
	}

	rule := DefaultA11yRule()
	result := ValidateAccessibility(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail due to skipped heading level")
	}

	foundHeadingError := false
	for _, issue := range result.Issues {
		if issue.Severity == "error" && issue.Component == "h3" {
			foundHeadingError = true
			break
		}
	}

	if !foundHeadingError {
		t.Error("Expected error about skipped heading level")
	}
}

func TestValidateAccessibility_NestingDepth(t *testing.T) {
	// Create a structure that exceeds max nesting depth
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
				ID:   "level0",
				Type: "box",
				Children: []types.Component{
					{
						ID:   "level1",
						Type: "box",
						Children: []types.Component{
							{
								ID:   "level2",
								Type: "box",
								Children: []types.Component{
									{
										ID:   "level3",
										Type: "box",
										Children: []types.Component{
											{
												ID:   "level4",
												Type: "box",
												Children: []types.Component{
													{
														ID:   "level5", // Too deep
														Type: "box",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Accessibility: types.Accessibility{
			TouchTargetsMin:   44,
			FocusIndicators:   "visible",
			Labels:            "all_interactive_elements",
			SemanticStructure: true,
		},
	}

	rule := DefaultA11yRule()
	result := ValidateAccessibility(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail due to excessive nesting depth")
	}

	foundDepthError := false
	for _, issue := range result.Issues {
		if issue.Severity == "error" && issue.Component == "level5" {
			foundDepthError = true
			break
		}
	}

	if !foundDepthError {
		t.Error("Expected error about excessive nesting depth")
	}
}

func TestValidateAccessibility_ValidStructure(t *testing.T) {
	// Create a valid accessible structure
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
				ID:   "h1",
				Type: "text",
				Size: "4xl",
				Role: "header",
			},
			{
				ID:   "h2",
				Type: "text",
				Size: "3xl",
			},
			{
				ID:   "username-label",
				Type: "text",
			},
			{
				ID:   "username-input",
				Type: "input",
			},
		},
		Accessibility: types.Accessibility{
			TouchTargetsMin:   44,
			FocusIndicators:   "visible",
			Labels:            "all_interactive_elements",
			SemanticStructure: true,
		},
	}

	rule := DefaultA11yRule()
	result := ValidateAccessibility(structure, rule)

	if !result.Passed {
		t.Errorf("Expected validation to pass for accessible structure, but got: %v", result.Issues)
	}

	// Should have info messages
	foundInfo := false
	for _, issue := range result.Issues {
		if issue.Severity == "info" {
			foundInfo = true
			break
		}
	}

	if !foundInfo {
		t.Error("Expected info messages for valid structure")
	}
}

func TestGetHeadingLevel(t *testing.T) {
	tests := []struct {
		comp     types.Component
		expected int
	}{
		{types.Component{ID: "h1", Type: "text"}, 1},
		{types.Component{ID: "h2", Type: "text"}, 2},
		{types.Component{ID: "h3", Type: "text"}, 3},
		{types.Component{ID: "title", Type: "text", Size: "4xl"}, 1},
		{types.Component{ID: "heading", Type: "text", Size: "3xl"}, 2},
		{types.Component{ID: "normal-text", Type: "text", Size: "base"}, 0},
	}

	for _, test := range tests {
		result := getHeadingLevel(&test.comp)
		if result != test.expected {
			t.Errorf("getHeadingLevel(%s, %s) = %d, expected %d",
				test.comp.ID, test.comp.Size, result, test.expected)
		}
	}
}

func TestHasLabel(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "username-label",
				Type: "text",
			},
			{
				ID:   "email-label",
				Type: "text",
			},
		},
		Accessibility: types.Accessibility{
			Labels: "",
		},
	}

	tests := []struct {
		comp     types.Component
		expected bool
	}{
		{types.Component{ID: "username-input", Type: "input"}, true},  // Has label
		{types.Component{ID: "email-input", Type: "input"}, true},     // Has label
		{types.Component{ID: "password-input", Type: "input"}, false}, // No label
		{types.Component{ID: "submit", Type: "button", Content: "Submit"}, true}, // Self-labeled
	}

	for _, test := range tests {
		result := hasLabel(&test.comp, structure)
		if result != test.expected {
			t.Errorf("hasLabel(%s) = %v, expected %v", test.comp.ID, result, test.expected)
		}
	}
}

func TestSharesPrefix(t *testing.T) {
	tests := []struct {
		id1      string
		id2      string
		expected bool
	}{
		{"username-label", "username-input", true},
		{"email-field", "password-field", false},
		{"submit-btn", "submit-form", true},
		{"cancel", "submit", false},
	}

	for _, test := range tests {
		result := sharesPrefix(test.id1, test.id2)
		if result != test.expected {
			t.Errorf("sharesPrefix(%s, %s) = %v, expected %v",
				test.id1, test.id2, result, test.expected)
		}
	}
}
