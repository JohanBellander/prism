package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseStructure(t *testing.T) {
	validJSON := `{
		"version": "v1",
		"phase": "structure",
		"created_at": "2025-10-25T12:00:00Z",
		"locked": false,
		"intent": {
			"purpose": "Test dashboard",
			"primary_action": "View metrics",
			"user_context": "Admin users",
			"key_interactions": ["view", "filter"]
		},
		"layout": {
			"type": "stack",
			"direction": "vertical",
			"spacing": 8,
			"max_width": 1200,
			"padding": 24
		},
		"components": [
			{
				"id": "header",
				"type": "box",
				"role": "header",
				"layout": {
					"display": "flex",
					"direction": "horizontal",
					"padding": 16,
					"background": "#FFFFFF",
					"border": "1px solid #E5E5E5",
					"gap": 8
				},
				"children": [
					{
						"id": "title",
						"type": "text",
						"content": "Dashboard",
						"size": "2xl",
						"weight": "bold",
						"color": "#000000"
					}
				]
			}
		],
		"responsive": {
			"mobile": {
				"breakpoint": 640,
				"changes": {
					"layout.padding": 16
				}
			},
			"tablet": {
				"breakpoint": 1024,
				"changes": {}
			}
		},
		"accessibility": {
			"touch_targets_min": 44,
			"focus_indicators": "visible",
			"labels": "all_interactive_elements",
			"semantic_structure": true
		},
		"validation": {
			"visual_hierarchy": "passed",
			"touch_targets": "passed",
			"max_nesting_depth": 4,
			"responsive_tested": true,
			"notes": "All checks passed"
		}
	}`

	s, err := ParseStructure([]byte(validJSON))
	if err != nil {
		t.Fatalf("ParseStructure failed: %v", err)
	}

	// Verify basic fields
	if s.Version != "v1" {
		t.Errorf("Expected version 'v1', got '%s'", s.Version)
	}
	if s.Phase != "structure" {
		t.Errorf("Expected phase 'structure', got '%s'", s.Phase)
	}
	if s.Intent.Purpose != "Test dashboard" {
		t.Errorf("Expected purpose 'Test dashboard', got '%s'", s.Intent.Purpose)
	}
	if s.Layout.Type != "stack" {
		t.Errorf("Expected layout type 'stack', got '%s'", s.Layout.Type)
	}
	if len(s.Components) != 1 {
		t.Errorf("Expected 1 component, got %d", len(s.Components))
	}
	if s.Components[0].ID != "header" {
		t.Errorf("Expected component ID 'header', got '%s'", s.Components[0].ID)
	}
	if len(s.Components[0].Children) != 1 {
		t.Errorf("Expected 1 child component, got %d", len(s.Components[0].Children))
	}
}

func TestParseStructure_InvalidJSON(t *testing.T) {
	invalidJSON := `{"version": "v1", "phase": "structure"`

	_, err := ParseStructure([]byte(invalidJSON))
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

func TestValidatePhase1_Valid(t *testing.T) {
	s := &Structure{
		Version: "v1",
		Phase:   "structure",
		Intent: Intent{
			Purpose:       "Test",
			PrimaryAction: "Action",
			UserContext:   "User",
		},
		Layout: Layout{
			Type:      "stack",
			Direction: "vertical",
			Spacing:   8,
		},
		Components: []Component{
			{
				ID:   "comp1",
				Type: "box",
				Role: "content",
				Layout: ComponentLayout{
					Display: "flex",
				},
			},
		},
	}

	err := s.ValidatePhase1()
	if err != nil {
		t.Errorf("ValidatePhase1 failed for valid structure: %v", err)
	}
}

func TestValidatePhase1_InvalidPhase(t *testing.T) {
	s := &Structure{
		Phase: "design",
	}

	err := s.ValidatePhase1()
	if err == nil {
		t.Error("Expected error for invalid phase, got nil")
	}
}

func TestValidatePhase1_MissingVersion(t *testing.T) {
	s := &Structure{
		Phase: "structure",
	}

	err := s.ValidatePhase1()
	if err == nil {
		t.Error("Expected error for missing version, got nil")
	}
}

func TestValidatePhase1_MissingIntent(t *testing.T) {
	s := &Structure{
		Version: "v1",
		Phase:   "structure",
		Layout: Layout{
			Type: "stack",
		},
		Components: []Component{
			{ID: "comp1", Type: "box"},
		},
	}

	err := s.ValidatePhase1()
	if err == nil {
		t.Error("Expected error for missing intent.purpose, got nil")
	}
}

func TestValidatePhase1_InvalidLayoutType(t *testing.T) {
	s := &Structure{
		Version: "v1",
		Phase:   "structure",
		Intent: Intent{
			Purpose: "Test",
		},
		Layout: Layout{
			Type: "invalid",
		},
		Components: []Component{
			{ID: "comp1", Type: "box"},
		},
	}

	err := s.ValidatePhase1()
	if err == nil {
		t.Error("Expected error for invalid layout type, got nil")
	}
}

func TestValidatePhase1_NoComponents(t *testing.T) {
	s := &Structure{
		Version: "v1",
		Phase:   "structure",
		Intent: Intent{
			Purpose: "Test",
		},
		Layout: Layout{
			Type: "stack",
		},
		Components: []Component{},
	}

	err := s.ValidatePhase1()
	if err == nil {
		t.Error("Expected error for no components, got nil")
	}
}

func TestValidateComponent_InvalidType(t *testing.T) {
	c := &Component{
		ID:   "comp1",
		Type: "invalid",
	}

	err := validateComponent(c, 0)
	if err == nil {
		t.Error("Expected error for invalid component type, got nil")
	}
}

func TestValidateComponent_InvalidColor(t *testing.T) {
	c := &Component{
		ID:    "comp1",
		Type:  "text",
		Color: "#FF0000", // Red - not allowed in Phase 1
	}

	err := validateComponent(c, 0)
	if err == nil {
		t.Error("Expected error for invalid color in Phase 1, got nil")
	}
}

func TestValidateComponent_ValidColors(t *testing.T) {
	validColors := []string{"#FFFFFF", "#000000", "#E5E5E5", "#737373", "#525252"}

	for _, color := range validColors {
		c := &Component{
			ID:    "comp1",
			Type:  "text",
			Color: color,
		}

		err := validateComponent(c, 0)
		if err != nil {
			t.Errorf("Expected valid color %s to pass, got error: %v", color, err)
		}
	}
}

func TestValidateComponent_MaxNestingDepth(t *testing.T) {
	// Create a component with 5 levels of nesting (exceeds max of 4)
	c := &Component{
		ID:   "level0",
		Type: "box",
		Children: []Component{
			{
				ID:   "level1",
				Type: "box",
				Children: []Component{
					{
						ID:   "level2",
						Type: "box",
						Children: []Component{
							{
								ID:   "level3",
								Type: "box",
								Children: []Component{
									{
										ID:   "level4",
										Type: "box",
										Children: []Component{
											{
												ID:   "level5",
												Type: "text",
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
	}

	err := validateComponent(c, 0)
	if err == nil {
		t.Error("Expected error for exceeding max nesting depth, got nil")
	}
}

func TestValidateComponent_ValidNesting(t *testing.T) {
	// Create a component with 3 levels of nesting (within max of 4)
	c := &Component{
		ID:   "level0",
		Type: "box",
		Children: []Component{
			{
				ID:   "level1",
				Type: "box",
				Children: []Component{
					{
						ID:   "level2",
						Type: "text",
					},
				},
			},
		},
	}

	err := validateComponent(c, 0)
	if err != nil {
		t.Errorf("Expected valid nesting to pass, got error: %v", err)
	}
}

func TestParseAndValidateStructure_Valid(t *testing.T) {
	validJSON := `{
		"version": "v1",
		"phase": "structure",
		"created_at": "2025-10-25T12:00:00Z",
		"locked": false,
		"intent": {
			"purpose": "Test",
			"primary_action": "Action",
			"user_context": "User",
			"key_interactions": ["action1"]
		},
		"layout": {
			"type": "stack",
			"direction": "vertical",
			"spacing": 8,
			"max_width": 1200,
			"padding": 24
		},
		"components": [
			{
				"id": "comp1",
				"type": "box",
				"role": "content",
				"layout": {
					"display": "flex"
				}
			}
		],
		"responsive": {
			"mobile": {"breakpoint": 640, "changes": {}},
			"tablet": {"breakpoint": 1024, "changes": {}}
		},
		"accessibility": {
			"touch_targets_min": 44,
			"focus_indicators": "visible",
			"labels": "all",
			"semantic_structure": true
		},
		"validation": {
			"visual_hierarchy": "passed",
			"touch_targets": "passed",
			"max_nesting_depth": 4,
			"responsive_tested": true
		}
	}`

	s, err := ParseAndValidateStructure([]byte(validJSON))
	if err != nil {
		t.Fatalf("ParseAndValidateStructure failed: %v", err)
	}

	if s.Version != "v1" {
		t.Errorf("Expected version 'v1', got '%s'", s.Version)
	}
}

func TestParseAndValidateStructure_Invalid(t *testing.T) {
	invalidJSON := `{
		"version": "v1",
		"phase": "design",
		"created_at": "2025-10-25T12:00:00Z",
		"locked": false,
		"intent": {"purpose": "Test"},
		"layout": {"type": "stack"},
		"components": [{"id": "comp1", "type": "box"}],
		"responsive": {
			"mobile": {"breakpoint": 640, "changes": {}},
			"tablet": {"breakpoint": 1024, "changes": {}}
		},
		"accessibility": {
			"touch_targets_min": 44,
			"focus_indicators": "visible",
			"labels": "all",
			"semantic_structure": true
		},
		"validation": {
			"visual_hierarchy": "passed",
			"touch_targets": "passed",
			"max_nesting_depth": 4,
			"responsive_tested": true
		}
	}`

	_, err := ParseAndValidateStructure([]byte(invalidJSON))
	if err == nil {
		t.Error("Expected error for invalid phase, got nil")
	}
}

func TestStructure_JSONRoundTrip(t *testing.T) {
	original := &Structure{
		Version:   "v1",
		Phase:     "structure",
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		Locked:    false,
		Intent: Intent{
			Purpose:         "Test roundtrip",
			PrimaryAction:   "Test",
			UserContext:     "Testing",
			KeyInteractions: []string{"test1", "test2"},
		},
		Layout: Layout{
			Type:      "stack",
			Direction: "vertical",
			Spacing:   8,
			MaxWidth:  1200,
			Padding:   24,
		},
		Components: []Component{
			{
				ID:   "test-component",
				Type: "box",
				Role: "content",
				Layout: ComponentLayout{
					Display:    "flex",
					Background: "#FFFFFF",
				},
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal back
	var parsed Structure
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Compare key fields
	if parsed.Version != original.Version {
		t.Errorf("Version mismatch: expected '%s', got '%s'", original.Version, parsed.Version)
	}
	if parsed.Phase != original.Phase {
		t.Errorf("Phase mismatch: expected '%s', got '%s'", original.Phase, parsed.Phase)
	}
	if parsed.Intent.Purpose != original.Intent.Purpose {
		t.Errorf("Intent.Purpose mismatch: expected '%s', got '%s'", original.Intent.Purpose, parsed.Intent.Purpose)
	}
}
