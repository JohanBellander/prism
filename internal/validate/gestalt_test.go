package validate

import (
	"testing"
	"time"

	"github.com/johanbellander/prism/internal/types"
)

func TestValidateGestalt_RelatedComponents(t *testing.T) {
	// Create a structure with related components
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
				ID:   "form-container",
				Type: "box",
				Layout: types.ComponentLayout{
					Display:   "flex",
					Direction: "vertical",
					Gap:       8, // Tight spacing suggests grouping
				},
				Children: []types.Component{
					{
						ID:   "username-label",
						Type: "text",
					},
					{
						ID:   "username-input",
						Type: "input",
					},
					{
						ID:   "password-label",
						Type: "text",
					},
					{
						ID:   "password-input",
						Type: "input",
					},
				},
			},
		},
	}

	rule := DefaultGestaltRule()
	result := ValidateGestalt(structure, rule)

	// Should detect the form container as a well-formed group
	foundGroup := false
	for _, issue := range result.Issues {
		if issue.Severity == "info" && issue.Component == "form-container" {
			foundGroup = true
			break
		}
	}

	if !foundGroup {
		t.Error("Expected to detect form-container as a well-formed group")
	}
}

func TestValidateGestalt_LargeSpacingBetweenRelated(t *testing.T) {
	// Create a structure with related components that have too much spacing
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
				ID:   "form-container",
				Type: "box",
				Layout: types.ComponentLayout{
					Display:   "flex",
					Direction: "vertical",
					Gap:       48, // Too much spacing for related items
				},
				Children: []types.Component{
					{
						ID:   "username-label",
						Type: "text",
					},
					{
						ID:   "username-input",
						Type: "input",
					},
				},
			},
		},
	}

	rule := DefaultGestaltRule()
	result := ValidateGestalt(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail due to large spacing between related components")
	}

	// Should have a warning about spacing
	foundWarning := false
	for _, issue := range result.Issues {
		if issue.Severity == "warning" {
			foundWarning = true
			break
		}
	}

	if !foundWarning {
		t.Error("Expected a warning about large spacing between related components")
	}
}

func TestValidateGestalt_SimilarityCheck(t *testing.T) {
	// Create a structure with similar components that have inconsistent styling
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
				ID:    "button1",
				Type:  "button",
				Size:  "lg",
				Color: "#000000",
				Layout: types.ComponentLayout{
					Padding: 16,
				},
			},
			{
				ID:    "button2",
				Type:  "button",
				Size:  "sm", // Different size
				Color: "#737373", // Different color
				Layout: types.ComponentLayout{
					Padding: 8, // Different padding
				},
			},
			{
				ID:    "button3",
				Type:  "button",
				Size:  "xl", // Another different size
				Color: "#000000",
				Layout: types.ComponentLayout{
					Padding: 24, // Another different padding
				},
			},
		},
	}

	rule := DefaultGestaltRule()
	rule.SimilarityCheck = true
	result := ValidateGestalt(structure, rule)

	// Should detect inconsistencies in similar components
	foundSimilarityWarning := false
	for _, issue := range result.Issues {
		if issue.Severity == "warning" && issue.Component == "button" {
			foundSimilarityWarning = true
			break
		}
	}

	if !foundSimilarityWarning {
		t.Error("Expected a warning about inconsistent styling in similar components")
	}
}

func TestAreComponentsRelated(t *testing.T) {
	tests := []struct {
		comp1    types.Component
		comp2    types.Component
		expected bool
	}{
		{
			types.Component{ID: "username-label", Type: "text"},
			types.Component{ID: "username-input", Type: "input"},
			true, // Same prefix
		},
		{
			types.Component{ID: "email-label", Type: "text"},
			types.Component{ID: "password-label", Type: "text"},
			false, // Different prefixes
		},
		{
			types.Component{ID: "btn1", Type: "button", Role: "primary"},
			types.Component{ID: "btn2", Type: "button", Role: "primary"},
			true, // Same type and role
		},
		{
			types.Component{ID: "label1", Type: "text"},
			types.Component{ID: "input1", Type: "input"},
			true, // Label-input pattern
		},
	}

	for _, test := range tests {
		result := areComponentsRelated(&test.comp1, &test.comp2)
		if result != test.expected {
			t.Errorf("areComponentsRelated(%s, %s) = %v, expected %v",
				test.comp1.ID, test.comp2.ID, result, test.expected)
		}
	}
}

func TestValidateGestalt_ValidStructure(t *testing.T) {
	// Create a valid structure with good grouping
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
			Spacing:   24, // Good inter-group spacing
		},
		Components: []types.Component{
			{
				ID:   "group1",
				Type: "box",
				Layout: types.ComponentLayout{
					Display:   "flex",
					Direction: "vertical",
					Gap:       8, // Good intra-group spacing
				},
				Children: []types.Component{
					{
						ID:   "item1",
						Type: "text",
						Size: "base",
					},
					{
						ID:   "item2",
						Type: "text",
						Size: "base", // Consistent styling
					},
				},
			},
			{
				ID:   "group2",
				Type: "box",
				Layout: types.ComponentLayout{
					Display:   "flex",
					Direction: "vertical",
					Gap:       8,
				},
				Children: []types.Component{
					{
						ID:   "item3",
						Type: "text",
						Size: "base",
					},
					{
						ID:   "item4",
						Type: "text",
						Size: "base",
					},
				},
			},
		},
	}

	rule := DefaultGestaltRule()
	result := ValidateGestalt(structure, rule)

	if !result.Passed {
		t.Error("Expected validation to pass for well-structured groups")
	}

	// Should have info messages about well-formed groups
	foundGroupInfo := false
	for _, issue := range result.Issues {
		if issue.Severity == "info" {
			foundGroupInfo = true
			break
		}
	}

	if !foundGroupInfo {
		t.Error("Expected info messages about well-formed groups")
	}
}
