package validate

import (
	"testing"

	"github.com/johanbellander/prism/internal/types"
)

func TestIsOnGrid(t *testing.T) {
	allowedScale := []int{0, 4, 8, 12, 16, 24, 32, 48, 64, 96, 128}

	tests := []struct {
		name     string
		value    int
		expected bool
	}{
		{"Zero is on grid", 0, true},
		{"4px is on grid", 4, true},
		{"8px is on grid", 8, true},
		{"16px is on grid", 16, true},
		{"24px is on grid", 24, true},
		{"15px is off grid", 15, false},
		{"10px is off grid", 10, false},
		{"20px is off grid", 20, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOnGrid(tt.value, allowedScale)
			if result != tt.expected {
				t.Errorf("Expected %v for %dpx, got %v", tt.expected, tt.value, result)
			}
		})
	}
}

func TestFindNearestGridValue(t *testing.T) {
	allowedScale := []int{0, 4, 8, 12, 16, 24, 32, 48, 64, 96, 128}

	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{"15px nearest is 16px", 15, 16},
		{"10px nearest is 8px", 10, 8},
		{"11px nearest is 12px", 11, 12},
		{"20px nearest is 16px", 20, 16}, // Equidistant from 16 and 24, rounds to first (lower)
		{"22px nearest is 24px", 22, 24},
		{"6px nearest is 4px", 6, 4},
		{"7px nearest is 8px", 7, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findNearestGridValue(tt.value, allowedScale)
			if result != tt.expected {
				t.Errorf("Expected nearest grid value %dpx for %dpx, got %dpx", tt.expected, tt.value, result)
			}
		})
	}
}

func TestValidateSpacing_OnGrid(t *testing.T) {
	// All spacing values on the grid
	structure := &types.Structure{
		Layout: types.Layout{
			Spacing: 16,
			Padding: 24,
		},
		Components: []types.Component{
			{
				ID:   "container",
				Type: "box",
				Layout: types.ComponentLayout{
					Padding: 16,
					Gap:     8,
				},
				Children: []types.Component{
					{
						ID:   "child",
						Type: "box",
						Layout: types.ComponentLayout{
							MarginBottom: 12,
						},
					},
				},
			},
		},
	}

	rule := DefaultSpacingRule()
	result := ValidateSpacing(structure, rule)

	if !result.Passed {
		t.Errorf("Expected validation to pass for on-grid spacing, but got %d issues", len(result.Issues))
		for _, issue := range result.Issues {
			t.Logf("  - %s: %s", issue.Category, issue.Message)
		}
	}
}

func TestValidateSpacing_OffGrid(t *testing.T) {
	// Component with off-grid padding
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "container",
				Type: "box",
				Layout: types.ComponentLayout{
					Padding: 15, // Off grid
				},
			},
		},
	}

	rule := DefaultSpacingRule()
	result := ValidateSpacing(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail for off-grid spacing")
	}

	foundOffGridIssue := false
	foundSuggestion := false
	for _, issue := range result.Issues {
		if issue.Category == "off_grid" && issue.ComponentID == "container" {
			foundOffGridIssue = true
			if issue.Value != 15 {
				t.Errorf("Expected value 15, got %d", issue.Value)
			}
			if issue.Suggested != 16 {
				t.Errorf("Expected suggested value 16, got %d", issue.Suggested)
			}
		}
		if issue.Category == "suggestion" && issue.Suggested == 16 {
			foundSuggestion = true
		}
	}

	if !foundOffGridIssue {
		t.Error("Expected off_grid issue for container")
	}
	if !foundSuggestion {
		t.Error("Expected suggestion for compliant value")
	}
}

func TestValidateSpacing_LayoutSpacing(t *testing.T) {
	// Layout with off-grid spacing
	structure := &types.Structure{
		Layout: types.Layout{
			Spacing: 20, // Off grid - nearest is 16 (distance 4) vs 24 (distance 4), should pick lower
			Padding: 8,  // On grid
		},
		Components: []types.Component{},
	}

	rule := DefaultSpacingRule()
	result := ValidateSpacing(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail for off-grid layout spacing")
	}

	foundLayoutIssue := false
	for _, issue := range result.Issues {
		if issue.Category == "off_grid" && issue.ComponentID == "layout" && issue.Property == "spacing" {
			foundLayoutIssue = true
			if issue.Value != 20 {
				t.Errorf("Expected value 20, got %d", issue.Value)
			}
			// 20 is equidistant from 16 and 24, algorithm picks first match which is 16
			if issue.Suggested != 16 && issue.Suggested != 24 {
				t.Errorf("Expected suggested value 16 or 24, got %d", issue.Suggested)
			}
		}
	}

	if !foundLayoutIssue {
		t.Error("Expected off_grid issue for layout spacing")
	}
}

func TestValidateSpacing_Gap(t *testing.T) {
	// Component with off-grid gap
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "flex_container",
				Type: "box",
				Layout: types.ComponentLayout{
					Gap: 10, // Off grid
				},
			},
		},
	}

	rule := DefaultSpacingRule()
	result := ValidateSpacing(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail for off-grid gap")
	}

	foundGapIssue := false
	for _, issue := range result.Issues {
		if issue.Category == "off_grid" && issue.Property == "gap" {
			foundGapIssue = true
			if issue.Suggested != 8 {
				t.Errorf("Expected nearest value 8px for gap 10px, got %d", issue.Suggested)
			}
		}
	}

	if !foundGapIssue {
		t.Error("Expected off_grid issue for gap")
	}
}

func TestValidateSpacing_MarginBottom(t *testing.T) {
	// Component with off-grid margin bottom
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "spaced_box",
				Type: "box",
				Layout: types.ComponentLayout{
					MarginBottom: 18, // Off grid
				},
			},
		},
	}

	rule := DefaultSpacingRule()
	result := ValidateSpacing(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail for off-grid margin_bottom")
	}

	foundMarginIssue := false
	for _, issue := range result.Issues {
		if issue.Category == "off_grid" && issue.Property == "margin_bottom" {
			foundMarginIssue = true
			if issue.Suggested != 16 {
				t.Errorf("Expected nearest value 16px for margin 18px, got %d", issue.Suggested)
			}
		}
	}

	if !foundMarginIssue {
		t.Error("Expected off_grid issue for margin_bottom")
	}
}

func TestValidateSpacing_MultipleIssues(t *testing.T) {
	// Multiple components with spacing issues
	structure := &types.Structure{
		Layout: types.Layout{
			Spacing: 15,
		},
		Components: []types.Component{
			{
				ID:   "container1",
				Type: "box",
				Layout: types.ComponentLayout{
					Padding: 10,
				},
			},
			{
				ID:   "container2",
				Type: "box",
				Layout: types.ComponentLayout{
					Gap: 20,
				},
			},
		},
	}

	rule := DefaultSpacingRule()
	result := ValidateSpacing(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail for multiple spacing issues")
	}

	// Count off_grid issues (not suggestions)
	offGridCount := 0
	for _, issue := range result.Issues {
		if issue.Category == "off_grid" {
			offGridCount++
		}
	}

	if offGridCount < 3 {
		t.Errorf("Expected at least 3 off_grid issues, got %d", offGridCount)
	}
}

func TestValidateSpacing_NestedComponents(t *testing.T) {
	// Nested components with spacing
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:   "parent",
				Type: "box",
				Layout: types.ComponentLayout{
					Padding: 16, // On grid
				},
				Children: []types.Component{
					{
						ID:   "child",
						Type: "box",
						Layout: types.ComponentLayout{
							Padding: 13, // Off grid
						},
					},
				},
			},
		},
	}

	rule := DefaultSpacingRule()
	result := ValidateSpacing(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail for nested off-grid spacing")
	}

	foundChildIssue := false
	for _, issue := range result.Issues {
		if issue.Category == "off_grid" && issue.ComponentID == "child" {
			foundChildIssue = true
		}
	}

	if !foundChildIssue {
		t.Error("Expected off_grid issue for nested child component")
	}
}

func TestValidateSpacing_EmptyStructure(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{},
	}

	rule := DefaultSpacingRule()
	result := ValidateSpacing(structure, rule)

	if !result.Passed {
		t.Error("Expected validation to pass for empty structure")
	}
}

func TestValidateSpacing_ZeroValues(t *testing.T) {
	// Zero is a valid spacing value
	structure := &types.Structure{
		Layout: types.Layout{
			Spacing: 0,
			Padding: 0,
		},
		Components: []types.Component{
			{
				ID:   "container",
				Type: "box",
				Layout: types.ComponentLayout{
					Padding: 0,
					Gap:     0,
				},
			},
		},
	}

	rule := DefaultSpacingRule()
	result := ValidateSpacing(structure, rule)

	if !result.Passed {
		t.Errorf("Expected validation to pass for zero spacing values, but got %d issues", len(result.Issues))
	}
}
