package validate

import (
	"testing"

	"github.com/johanbellander/prism/internal/types"
)

func TestValidateChoiceOverload_NavigationOverload(t *testing.T) {
	// Too many navigation items (>7)
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "nav1",
				Type: "container",
				Role: "navigation",
				Children: []types.Component{
					{ID: "item1", Type: "button"},
					{ID: "item2", Type: "button"},
					{ID: "item3", Type: "button"},
					{ID: "item4", Type: "button"},
					{ID: "item5", Type: "button"},
					{ID: "item6", Type: "button"},
					{ID: "item7", Type: "button"},
					{ID: "item8", Type: "button"}, // 8th item - exceeds limit
				},
			},
		},
	}

	rule := DefaultChoiceRule()
	result := ValidateChoiceOverload(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail for navigation with >7 items")
	}

	if len(result.Issues) == 0 {
		t.Error("Expected issues to be reported")
	}

	foundNavIssue := false
	for _, issue := range result.Issues {
		if issue.Category == "navigation_overload" {
			foundNavIssue = true
			if issue.ComponentID != "nav1" {
				t.Errorf("Expected component ID 'nav1', got '%s'", issue.ComponentID)
			}
		}
	}

	if !foundNavIssue {
		t.Error("Expected navigation_overload issue")
	}
}

func TestValidateChoiceOverload_FormFieldOverload(t *testing.T) {
	// Too many form fields (>7)
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "form1",
				Type: "container",
				Role: "form",
				Children: []types.Component{
					{ID: "field1", Type: "input", Role: "textbox"},
					{ID: "field2", Type: "input", Role: "textbox"},
					{ID: "field3", Type: "input", Role: "textbox"},
					{ID: "field4", Type: "input", Role: "textbox"},
					{ID: "field5", Type: "input", Role: "textbox"},
					{ID: "field6", Type: "input", Role: "textbox"},
					{ID: "field7", Type: "input", Role: "textbox"},
					{ID: "field8", Type: "input", Role: "textbox"}, // 8th field - exceeds limit
				},
			},
		},
	}

	rule := DefaultChoiceRule()
	result := ValidateChoiceOverload(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail for form with >7 fields")
	}

	foundFormIssue := false
	for _, issue := range result.Issues {
		if issue.Category == "form_overload" {
			foundFormIssue = true
			if issue.ComponentID != "form1" {
				t.Errorf("Expected component ID 'form1', got '%s'", issue.ComponentID)
			}
		}
	}

	if !foundFormIssue {
		t.Error("Expected form_overload issue")
	}
}

func TestValidateChoiceOverload_ButtonGroupOverload(t *testing.T) {
	// Too many buttons in container with buttons (>3)
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "btngroup1",
				Type: "container",
				Children: []types.Component{
					{ID: "btn1", Type: "button"},
					{ID: "btn2", Type: "button"},
					{ID: "btn3", Type: "button"},
					{ID: "btn4", Type: "button"}, // 4th button - exceeds limit
				},
			},
		},
	}

	rule := DefaultChoiceRule()
	result := ValidateChoiceOverload(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail for button group with >3 buttons")
	}

	foundGroupIssue := false
	for _, issue := range result.Issues {
		if issue.Category == "button_group_overload" {
			foundGroupIssue = true
			if issue.ComponentID != "btngroup1" {
				t.Errorf("Expected component ID 'btngroup1', got '%s'", issue.ComponentID)
			}
		}
	}

	if !foundGroupIssue {
		t.Error("Expected button_group_overload issue")
	}
}

func TestValidateChoiceOverload_CardGridOverload(t *testing.T) {
	// Too many cards in grid (>12)
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "grid1",
				Type: "container",
				Layout: types.ComponentLayout{
					Display: "grid",
				},
				Children: []types.Component{
					{ID: "card1", Type: "card"},
					{ID: "card2", Type: "card"},
					{ID: "card3", Type: "card"},
					{ID: "card4", Type: "card"},
					{ID: "card5", Type: "card"},
					{ID: "card6", Type: "card"},
					{ID: "card7", Type: "card"},
					{ID: "card8", Type: "card"},
					{ID: "card9", Type: "card"},
					{ID: "card10", Type: "card"},
					{ID: "card11", Type: "card"},
					{ID: "card12", Type: "card"},
					{ID: "card13", Type: "card"}, // 13th card - exceeds limit
				},
			},
		},
	}

	rule := DefaultChoiceRule()
	result := ValidateChoiceOverload(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail for card grid with >12 cards")
	}

	foundGridIssue := false
	for _, issue := range result.Issues {
		if issue.Category == "card_grid_overload" {
			foundGridIssue = true
			if issue.ComponentID != "grid1" {
				t.Errorf("Expected component ID 'grid1', got '%s'", issue.ComponentID)
			}
		}
	}

	if !foundGridIssue {
		t.Error("Expected card_grid_overload issue")
	}
}

func TestValidateChoiceOverload_ValidStructure(t *testing.T) {
	// All containers within limits
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "nav1",
				Type: "container",
				Role: "navigation",
				Children: []types.Component{
					{ID: "item1", Type: "button"},
					{ID: "item2", Type: "button"},
					{ID: "item3", Type: "button"},
					{ID: "item4", Type: "button"},
					{ID: "item5", Type: "button"},
				},
			},
			{
				ID:   "form1",
				Type: "container",
				Role: "form",
				Children: []types.Component{
					{ID: "field1", Type: "input", Role: "textbox"},
					{ID: "field2", Type: "input", Role: "textbox"},
					{ID: "field3", Type: "input", Role: "textbox"},
				},
			},
			{
				ID:   "btngroup1",
				Type: "container",
				Children: []types.Component{
					{ID: "btn1", Type: "button"},
					{ID: "btn2", Type: "button"},
					{ID: "btn3", Type: "button"},
				},
			},
			{
				ID:   "grid1",
				Type: "container",
				Layout: types.ComponentLayout{
					Display: "grid",
				},
				Children: []types.Component{
					{ID: "card1", Type: "card"},
					{ID: "card2", Type: "card"},
					{ID: "card3", Type: "card"},
					{ID: "card4", Type: "card"},
					{ID: "card5", Type: "card"},
					{ID: "card6", Type: "card"},
				},
			},
		},
	}

	rule := DefaultChoiceRule()
	result := ValidateChoiceOverload(structure, rule)

	if !result.Passed {
		t.Errorf("Expected validation to pass for valid structure, but got %d issues", len(result.Issues))
		for _, issue := range result.Issues {
			t.Logf("  - %s: %s", issue.Category, issue.Message)
		}
	}

	if len(result.Issues) != 0 {
		t.Errorf("Expected no issues for valid structure, got %d", len(result.Issues))
	}
}

func TestValidateChoiceOverload_MultipleIssues(t *testing.T) {
	// Multiple violations
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "nav1",
				Type: "container",
				Role: "navigation",
				Children: []types.Component{
					{ID: "item1", Type: "button"},
					{ID: "item2", Type: "button"},
					{ID: "item3", Type: "button"},
					{ID: "item4", Type: "button"},
					{ID: "item5", Type: "button"},
					{ID: "item6", Type: "button"},
					{ID: "item7", Type: "button"},
					{ID: "item8", Type: "button"}, // Navigation overload
				},
			},
			{
				ID:   "form1",
				Type: "container",
				Role: "form",
				Children: []types.Component{
					{ID: "field1", Type: "input", Role: "textbox"},
					{ID: "field2", Type: "input", Role: "textbox"},
					{ID: "field3", Type: "input", Role: "textbox"},
					{ID: "field4", Type: "input", Role: "textbox"},
					{ID: "field5", Type: "input", Role: "textbox"},
					{ID: "field6", Type: "input", Role: "textbox"},
					{ID: "field7", Type: "input", Role: "textbox"},
					{ID: "field8", Type: "input", Role: "textbox"}, // Form overload
				},
			},
		},
	}

	rule := DefaultChoiceRule()
	result := ValidateChoiceOverload(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail for multiple violations")
	}

	if len(result.Issues) < 2 {
		t.Errorf("Expected at least 2 issues, got %d", len(result.Issues))
	}

	hasNavIssue := false
	hasFormIssue := false
	for _, issue := range result.Issues {
		if issue.Category == "navigation_overload" {
			hasNavIssue = true
		}
		if issue.Category == "form_overload" {
			hasFormIssue = true
		}
	}

	if !hasNavIssue {
		t.Error("Expected navigation_overload issue")
	}
	if !hasFormIssue {
		t.Error("Expected form_overload issue")
	}
}

func TestValidateChoiceOverload_EmptyStructure(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{},
	}

	rule := DefaultChoiceRule()
	result := ValidateChoiceOverload(structure, rule)

	if !result.Passed {
		t.Error("Expected validation to pass for empty structure")
	}

	if len(result.Issues) != 0 {
		t.Errorf("Expected no issues for empty structure, got %d", len(result.Issues))
	}
}
