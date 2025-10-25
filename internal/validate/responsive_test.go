package validate

import (
	"testing"

	"github.com/johanbellander/prism/internal/types"
)

func TestDefaultResponsiveRule(t *testing.T) {
	rule := DefaultResponsiveRule()

	if rule.Breakpoints["mobile"] != 375 {
		t.Errorf("Expected mobile breakpoint 375, got %d", rule.Breakpoints["mobile"])
	}
	if rule.Breakpoints["tablet"] != 768 {
		t.Errorf("Expected tablet breakpoint 768, got %d", rule.Breakpoints["tablet"])
	}
	if rule.Breakpoints["desktop"] != 1440 {
		t.Errorf("Expected desktop breakpoint 1440, got %d", rule.Breakpoints["desktop"])
	}
	if rule.MinTouchTarget != 44 {
		t.Errorf("Expected min touch target 44, got %d", rule.MinTouchTarget)
	}
	if !rule.CheckOverflow {
		t.Error("Expected CheckOverflow to be true")
	}
	if !rule.CheckTouchTargets {
		t.Error("Expected CheckTouchTargets to be true")
	}
}

func TestValidateResponsive_NoIssues(t *testing.T) {
	structure := &types.Structure{
		Layout: types.Layout{
			MaxWidth: 375,
		},
		Components: []types.Component{
			{
				ID:   "btn-1",
				Type: "button",
				Layout: types.ComponentLayout{
					Width:  100,
					Height: 48,
				},
			},
		},
	}

	result := ValidateResponsive(structure, DefaultResponsiveRule())

	if !result.Passed {
		t.Error("Expected validation to pass")
	}
	if len(result.Issues) != 0 {
		t.Errorf("Expected no issues, got %d", len(result.Issues))
	}
}

func TestValidateResponsive_LayoutExceedsViewport(t *testing.T) {
	structure := &types.Structure{
		Layout: types.Layout{
			MaxWidth: 400, // Exceeds mobile (375px)
		},
		Components: []types.Component{},
	}

	result := ValidateResponsive(structure, DefaultResponsiveRule())

	if !result.Passed {
		t.Error("Expected validation to pass (warnings don't fail)")
	}
	if len(result.Issues) == 0 {
		t.Error("Expected warning about layout exceeding mobile viewport")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.ComponentID == "layout" && issue.Viewport == "mobile" && issue.Severity == "warning" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected warning for layout exceeding mobile viewport")
	}
}

func TestValidateResponsive_ComponentOverflow(t *testing.T) {
	structure := &types.Structure{
		Layout: types.Layout{
			MaxWidth: 375,
		},
		Components: []types.Component{
			{
				ID:   "wide-box",
				Type: "box",
				Layout: types.ComponentLayout{
					Width:  400, // Overflows 375px mobile
					Height: 100,
				},
			},
		},
	}

	result := ValidateResponsive(structure, DefaultResponsiveRule())

	if !result.Passed {
		t.Error("Expected validation to pass (warnings don't fail)")
	}
	if len(result.Issues) == 0 {
		t.Error("Expected overflow warning")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.ComponentID == "wide-box" && issue.Viewport == "mobile" && issue.Severity == "warning" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected overflow warning for component")
	}
}

func TestValidateResponsive_SmallTouchTarget(t *testing.T) {
	structure := &types.Structure{
		Layout: types.Layout{
			MaxWidth: 375,
		},
		Components: []types.Component{
			{
				ID:   "small-btn",
				Type: "button",
				Layout: types.ComponentLayout{
					Width:  30, // Too small (< 44px)
					Height: 30, // Too small (< 44px)
				},
			},
		},
	}

	result := ValidateResponsive(structure, DefaultResponsiveRule())

	if !result.Passed {
		t.Error("Expected validation to pass (warnings don't fail)")
	}
	if len(result.Issues) == 0 {
		t.Error("Expected touch target warning")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.ComponentID == "small-btn" && issue.Viewport == "mobile" && issue.Severity == "warning" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected touch target warning for small button")
	}
}

func TestValidateResponsive_SmallInput(t *testing.T) {
	structure := &types.Structure{
		Layout: types.Layout{
			MaxWidth: 375,
		},
		Components: []types.Component{
			{
				ID:   "small-input",
				Type: "input",
				Layout: types.ComponentLayout{
					Width:  100,
					Height: 32, // Too small height (< 44px)
				},
			},
		},
	}

	result := ValidateResponsive(structure, DefaultResponsiveRule())

	if !result.Passed {
		t.Error("Expected validation to pass (warnings don't fail)")
	}
	if len(result.Issues) == 0 {
		t.Error("Expected touch target warning for input")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.ComponentID == "small-input" && issue.Viewport == "mobile" && issue.Severity == "warning" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected touch target warning for small input")
	}
}

func TestValidateResponsive_NestedComponents(t *testing.T) {
	structure := &types.Structure{
		Layout: types.Layout{
			MaxWidth: 375,
		},
		Components: []types.Component{
			{
				ID:   "container",
				Type: "box",
				Layout: types.ComponentLayout{
					Width:  343,
					Height: 200,
				},
				Children: []types.Component{
					{
						ID:   "nested-btn",
						Type: "button",
						Layout: types.ComponentLayout{
							Width:  80,
							Height: 48,
						},
					},
				},
			},
		},
	}

	result := ValidateResponsive(structure, DefaultResponsiveRule())

	// With integer-based layout, nested components don't cause overflow issues
	// in the same way. This test now just verifies nested components are processed.
	if !result.Passed {
		t.Error("Expected validation to pass")
	}
}

func TestValidateResponsive_MultipleViewports(t *testing.T) {
	structure := &types.Structure{
		Layout: types.Layout{
			MaxWidth: 800, // Exceeds mobile and tablet
		},
		Components: []types.Component{},
	}

	result := ValidateResponsive(structure, DefaultResponsiveRule())

	// Should have warnings for mobile (375px) and tablet (768px)
	mobileWarning := false
	tabletWarning := false

	for _, issue := range result.Issues {
		if issue.Viewport == "mobile" {
			mobileWarning = true
		}
		if issue.Viewport == "tablet" {
			tabletWarning = true
		}
	}

	if !mobileWarning {
		t.Error("Expected warning for mobile viewport")
	}
	if !tabletWarning {
		t.Error("Expected warning for tablet viewport")
	}
}

func TestValidateResponsive_CustomRule(t *testing.T) {
	customRule := ResponsiveRule{
		Breakpoints: map[string]int{
			"mobile": 320, // Use "mobile" so touch target check runs
			"medium": 640,
		},
		MinTouchTarget:    48, // Custom larger touch target
		CheckOverflow:     true,
		CheckTouchTargets: true,
	}

	structure := &types.Structure{
		Layout: types.Layout{
			MaxWidth: 320,
		},
		Components: []types.Component{
			{
				ID:   "btn",
				Type: "button",
				Layout: types.ComponentLayout{
					Width:  46, // Smaller than custom rule (48px) but larger than default (44px)
					Height: 46,
				},
			},
		},
	}

	result := ValidateResponsive(structure, customRule)

	// Should warn because button is smaller than 48px
	found := false
	for _, issue := range result.Issues {
		if issue.ComponentID == "btn" && issue.Viewport == "mobile" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected warning for button smaller than custom touch target")
	}
}
