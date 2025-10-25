package validate

import (
	"testing"

	"github.com/johanbellander/prism/internal/types"
)

func TestDefaultLoadingStateRule(t *testing.T) {
	rule := DefaultLoadingStateRule()
	
	expectedStates := []string{"loading", "error", "empty", "default", ""}
	
	if len(rule.ValidStates) != len(expectedStates) {
		t.Errorf("Expected %d valid states, got %d", len(expectedStates), len(rule.ValidStates))
	}
	
	for _, state := range expectedStates {
		found := false
		for _, valid := range rule.ValidStates {
			if state == valid {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected valid state '%s' not found", state)
		}
	}
}

func TestValidateLoadingStates_ValidStates(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{ID: "comp1", Type: "box", State: "loading"},
			{ID: "comp2", Type: "box", State: "error"},
			{ID: "comp3", Type: "box", State: "empty"},
			{ID: "comp4", Type: "box", State: "default"},
			{ID: "comp5", Type: "box", State: ""},
		},
	}
	
	rule := DefaultLoadingStateRule()
	result := ValidateLoadingStates(structure, rule)
	
	// Should pass but might have info messages
	if !result.Passed {
		t.Errorf("Expected validation to pass for valid states")
	}
}

func TestValidateLoadingStates_InvalidState(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{ID: "comp1", Type: "box", State: "invalid"},
		},
	}
	
	rule := DefaultLoadingStateRule()
	result := ValidateLoadingStates(structure, rule)
	
	if result.Passed {
		t.Errorf("Expected validation to fail for invalid state")
	}
	
	// Should have error about invalid state
	foundError := false
	for _, issue := range result.Issues {
		if issue.Severity == "error" && issue.ComponentID == "comp1" {
			foundError = true
			break
		}
	}
	
	if !foundError {
		t.Errorf("Expected error about invalid state")
	}
}

func TestValidateLoadingStates_LoadingWithSkeleton(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:    "user-card",
				Type:  "box",
				State: "loading",
				Skeleton: &types.SkeletonConfig{
					Elements: []types.SkeletonElement{
						{Type: "circle", Size: 48},
						{Type: "text", Width: "60%"},
						{Type: "text", Width: "80%"},
					},
				},
			},
		},
	}
	
	rule := DefaultLoadingStateRule()
	result := ValidateLoadingStates(structure, rule)
	
	if !result.Passed {
		t.Errorf("Expected validation to pass for loading state with skeleton")
	}
}

func TestValidateLoadingStates_LoadingWithoutSkeleton(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{ID: "comp1", Type: "box", State: "loading"},
		},
	}
	
	rule := DefaultLoadingStateRule()
	result := ValidateLoadingStates(structure, rule)
	
	// Should pass but have info message
	if !result.Passed {
		t.Errorf("Expected validation to pass")
	}
	
	// Should have info about missing skeleton
	foundInfo := false
	for _, issue := range result.Issues {
		if issue.Severity == "info" && issue.ComponentID == "comp1" {
			foundInfo = true
			break
		}
	}
	
	if !foundInfo {
		t.Errorf("Expected info about missing skeleton configuration")
	}
}

func TestValidateLoadingStates_InvalidSkeletonType(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:    "comp1",
				Type:  "box",
				State: "loading",
				Skeleton: &types.SkeletonConfig{
					Elements: []types.SkeletonElement{
						{Type: "invalid-type", Width: "100%"},
					},
				},
			},
		},
	}
	
	rule := DefaultLoadingStateRule()
	result := ValidateLoadingStates(structure, rule)
	
	if result.Passed {
		t.Errorf("Expected validation to fail for invalid skeleton type")
	}
}

func TestValidateLoadingStates_EmptySkeletonElements(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:       "comp1",
				Type:     "box",
				State:    "loading",
				Skeleton: &types.SkeletonConfig{
					Elements: []types.SkeletonElement{},
				},
			},
		},
	}
	
	rule := DefaultLoadingStateRule()
	result := ValidateLoadingStates(structure, rule)
	
	// Should have warning about empty skeleton
	foundWarning := false
	for _, issue := range result.Issues {
		if issue.Severity == "warning" && issue.ComponentID == "comp1" {
			foundWarning = true
			break
		}
	}
	
	if !foundWarning {
		t.Errorf("Expected warning about empty skeleton elements")
	}
}

func TestValidateLoadingStates_EmptyState(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{ID: "comp1", Type: "box", State: "empty"},
		},
	}
	
	rule := DefaultLoadingStateRule()
	result := ValidateLoadingStates(structure, rule)
	
	// Should have info about adding empty state message
	foundInfo := false
	for _, issue := range result.Issues {
		if issue.Severity == "info" && issue.ComponentID == "comp1" {
			foundInfo = true
			break
		}
	}
	
	if !foundInfo {
		t.Errorf("Expected info about adding empty state message")
	}
}

func TestValidateLoadingStates_ErrorState(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{ID: "comp1", Type: "box", State: "error"},
		},
	}
	
	rule := DefaultLoadingStateRule()
	result := ValidateLoadingStates(structure, rule)
	
	// Should have info about adding error message
	foundInfo := false
	for _, issue := range result.Issues {
		if issue.Severity == "info" && issue.ComponentID == "comp1" {
			foundInfo = true
			break
		}
	}
	
	if !foundInfo {
		t.Errorf("Expected info about adding error message")
	}
}

func TestValidateLoadingStates_NestedComponents(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:    "parent",
				Type:  "box",
				State: "default",
				Children: []types.Component{
					{ID: "child1", Type: "box", State: "loading"},
					{ID: "child2", Type: "box", State: "error"},
				},
			},
		},
	}
	
	rule := DefaultLoadingStateRule()
	result := ValidateLoadingStates(structure, rule)
	
	// Should validate nested components
	if !result.Passed {
		t.Errorf("Expected validation to pass for nested components")
	}
	
	// Should have messages for both children
	if len(result.Issues) < 2 {
		t.Errorf("Expected info messages for nested components")
	}
}

func TestIsValidState(t *testing.T) {
	validStates := []string{"loading", "error", "empty", "default", ""}
	
	tests := []struct {
		state    string
		expected bool
	}{
		{"loading", true},
		{"error", true},
		{"empty", true},
		{"default", true},
		{"", true},
		{"invalid", false},
		{"Loading", false}, // case sensitive
	}
	
	for _, tt := range tests {
		result := isValidState(tt.state, validStates)
		if result != tt.expected {
			t.Errorf("State '%s': expected %v, got %v", tt.state, tt.expected, result)
		}
	}
}

func TestIsValidSkeletonType(t *testing.T) {
	tests := []struct {
		skeletonType string
		expected     bool
	}{
		{"circle", true},
		{"text", true},
		{"rect", true},
		{"invalid", false},
		{"Circle", false}, // case sensitive
	}
	
	for _, tt := range tests {
		result := isValidSkeletonType(tt.skeletonType)
		if result != tt.expected {
			t.Errorf("Type '%s': expected %v, got %v", tt.skeletonType, tt.expected, result)
		}
	}
}

func TestCountComponentsByState(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{ID: "comp1", Type: "box", State: "loading"},
			{ID: "comp2", Type: "box", State: "loading"},
			{ID: "comp3", Type: "box", State: "error"},
			{ID: "comp4", Type: "box", State: ""},
			{ID: "comp5", Type: "box"},
		},
	}
	
	counts := CountComponentsByState(structure)
	
	if counts["loading"] != 2 {
		t.Errorf("Expected 2 loading components, got %d", counts["loading"])
	}
	
	if counts["error"] != 1 {
		t.Errorf("Expected 1 error component, got %d", counts["error"])
	}
	
	if counts["default"] != 2 {
		t.Errorf("Expected 2 default components (empty state treated as default), got %d", counts["default"])
	}
}

func TestValidateLoadingStates_MissingSkeletonDimensions(t *testing.T) {
	structure := &types.Structure{
		Components: []types.Component{
			{
				ID:    "comp1",
				Type:  "box",
				State: "loading",
				Skeleton: &types.SkeletonConfig{
					Elements: []types.SkeletonElement{
						{Type: "circle"}, // missing size
						{Type: "text"},   // missing width
					},
				},
			},
		},
	}
	
	rule := DefaultLoadingStateRule()
	result := ValidateLoadingStates(structure, rule)
	
	// Should have warnings about missing dimensions
	warningCount := 0
	for _, issue := range result.Issues {
		if issue.Severity == "warning" {
			warningCount++
		}
	}
	
	if warningCount < 2 {
		t.Errorf("Expected at least 2 warnings about missing dimensions, got %d", warningCount)
	}
}
