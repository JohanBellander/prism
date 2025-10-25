package validate

import (
	"testing"
	"time"

	"github.com/johanbellander/prism/internal/types"
)

func TestValidateHierarchy_HeadingSizes(t *testing.T) {
	// Create a structure with heading hierarchy issues
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
				Size: "2xl", // 24px
				Layout: types.ComponentLayout{
					Display: "block",
				},
			},
			{
				ID:   "h2",
				Type: "text",
				Size: "3xl", // 30px - should be smaller than h1
				Layout: types.ComponentLayout{
					Display: "block",
				},
			},
		},
	}

	rule := DefaultHierarchyRule()
	result := ValidateHierarchy(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail due to inverted heading sizes")
	}

	if len(result.Issues) == 0 {
		t.Error("Expected at least one issue to be reported")
	}

	// Check that we have a warning about heading sizes
	foundWarning := false
	for _, issue := range result.Issues {
		if issue.Severity == "warning" {
			foundWarning = true
			break
		}
	}

	if !foundWarning {
		t.Error("Expected a warning about heading size hierarchy")
	}
}

func TestValidateHierarchy_ButtonSizes(t *testing.T) {
	// Create a structure with button size issues
	structure := &types.Structure{
		Version:   "v1",
		Phase:     "structure",
		CreatedAt: time.Now(),
		Intent: types.Intent{
			Purpose:       "Test",
			PrimaryAction: "primary-btn",
		},
		Layout: types.Layout{
			Type:      "stack",
			Direction: "vertical",
			Spacing:   16,
		},
		Components: []types.Component{
			{
				ID:   "primary-btn",
				Type: "button",
				Layout: types.ComponentLayout{
					Display: "block",
					Width:   100, // Less than minimum 120px
				},
			},
			{
				ID:   "secondary-btn",
				Type: "button",
				Layout: types.ComponentLayout{
					Display: "block",
					Width:   150, // Larger than primary
				},
			},
		},
	}

	rule := DefaultHierarchyRule()
	result := ValidateHierarchy(structure, rule)

	if result.Passed {
		t.Error("Expected validation to fail due to button size issues")
	}

	// Should have at least two issues: primary too small, secondary larger than primary
	errorCount := 0
	warningCount := 0
	for _, issue := range result.Issues {
		if issue.Severity == "error" {
			errorCount++
		}
		if issue.Severity == "warning" {
			warningCount++
		}
	}

	if errorCount == 0 {
		t.Error("Expected at least one error about secondary button being larger than primary")
	}

	if warningCount == 0 {
		t.Error("Expected at least one warning about primary button size")
	}
}

func TestValidateHierarchy_ValidStructure(t *testing.T) {
	// Create a valid structure
	structure := &types.Structure{
		Version:   "v1",
		Phase:     "structure",
		CreatedAt: time.Now(),
		Intent: types.Intent{
			Purpose:       "Test",
			PrimaryAction: "primary-btn",
		},
		Layout: types.Layout{
			Type:      "stack",
			Direction: "vertical",
			Spacing:   24,
		},
		Components: []types.Component{
			{
				ID:   "title",
				Type: "text",
				Size: "3xl",
				Layout: types.ComponentLayout{
					Display: "block",
					Padding: 16,
				},
			},
			{
				ID:   "subtitle",
				Type: "text",
				Size: "xl",
				Layout: types.ComponentLayout{
					Display: "block",
					Padding: 8,
				},
			},
			{
				ID:   "primary-btn",
				Type: "button",
				Layout: types.ComponentLayout{
					Display: "block",
					Width:   150,
				},
			},
			{
				ID:   "secondary-btn",
				Type: "button",
				Layout: types.ComponentLayout{
					Display: "block",
					Width:   120,
				},
			},
		},
	}

	rule := DefaultHierarchyRule()
	result := ValidateHierarchy(structure, rule)

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
