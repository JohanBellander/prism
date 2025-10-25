package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// Structure represents the complete Phase 1 structure JSON
type Structure struct {
	Version       string        `json:"version"`
	Phase         string        `json:"phase"`
	CreatedAt     time.Time     `json:"created_at"`
	Locked        bool          `json:"locked"`
	ParentVersion string        `json:"parent_version,omitempty"`
	ChangeSummary string        `json:"change_summary,omitempty"`
	Rationale     string        `json:"rationale,omitempty"`
	LockedAt      *time.Time    `json:"locked_at,omitempty"`
	ApprovedBy    string        `json:"approved_by,omitempty"`
	Checksum      string        `json:"checksum,omitempty"`
	Note          string        `json:"note,omitempty"`
	Intent        Intent        `json:"intent"`
	Layout        Layout        `json:"layout"`
	Components    []Component   `json:"components"`
	Responsive    Responsive    `json:"responsive"`
	Accessibility Accessibility `json:"accessibility"`
	Validation    Validation    `json:"validation"`
}

// Intent describes the purpose and context of the UI
type Intent struct {
	Purpose         string   `json:"purpose"`
	PrimaryAction   string   `json:"primary_action"`
	UserContext     string   `json:"user_context"`
	KeyInteractions []string `json:"key_interactions"`
}

// Layout defines the top-level layout configuration
type Layout struct {
	Type      string `json:"type"`       // "stack", "grid", "sidebar"
	Direction string `json:"direction"`  // "vertical", "horizontal"
	Spacing   int    `json:"spacing"`    // spacing in pixels
	MaxWidth  int    `json:"max_width"`  // max width in pixels
	Padding   int    `json:"padding"`    // padding in pixels
}

// Component represents a UI component
type Component struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`     // "box", "text", "input", "button", "image"
	Role     string           `json:"role"`     // "header", "navigation", "content", "footer", etc
	State    string           `json:"state,omitempty"`    // "loading", "error", "empty", "default"
	Layout   ComponentLayout  `json:"layout"`
	Content  string           `json:"content,omitempty"`
	Size     string           `json:"size,omitempty"`     // "xs", "sm", "base", "lg", "xl", "2xl", "3xl", "4xl"
	Weight   string           `json:"weight,omitempty"`   // "normal", "bold"
	Color    string           `json:"color,omitempty"`    // hex color
	Children []Component      `json:"children,omitempty"`
	Skeleton *SkeletonConfig  `json:"skeleton,omitempty"` // Skeleton placeholder configuration
}

// SkeletonConfig defines the skeleton/placeholder structure for loading states
type SkeletonConfig struct {
	Elements []SkeletonElement `json:"elements,omitempty"`
}

// SkeletonElement represents a placeholder element in skeleton screen
type SkeletonElement struct {
	Type   string `json:"type"`            // "circle", "text", "rect"
	Width  string `json:"width,omitempty"` // e.g., "60%" or "120px"
	Height string `json:"height,omitempty"`
	Size   int    `json:"size,omitempty"`  // For circles
}

// ComponentLayout defines layout properties for a component
type ComponentLayout struct {
	Display             string `json:"display"`                        // "flex", "block", "grid"
	Direction           string `json:"direction,omitempty"`            // "horizontal", "vertical"
	Padding             int    `json:"padding,omitempty"`              // padding in pixels
	Background          string `json:"background,omitempty"`           // hex color
	Border              string `json:"border,omitempty"`               // e.g., "1px solid #E5E5E5"
	BorderBottom        string `json:"border_bottom,omitempty"`        // e.g., "1px solid #E5E5E5"
	BorderRight         string `json:"border_right,omitempty"`         // e.g., "1px solid #E5E5E5"
	Gap                 int    `json:"gap,omitempty"`                  // gap in pixels
	GridTemplateColumns string `json:"grid_template_columns,omitempty"` // e.g., "repeat(4, 1fr)"
	Width               int    `json:"width,omitempty"`                // width in pixels
	Height              int    `json:"height,omitempty"`               // height in pixels
	MinHeight           string `json:"min_height,omitempty"`           // e.g., "calc(100vh - 64px)"
	MaxWidth            int    `json:"max_width,omitempty"`            // max width in pixels
	Flex                int    `json:"flex,omitempty"`                 // flex grow factor
	JustifyContent      string `json:"justify_content,omitempty"`      // "flex-start", "center", "space-between"
	AlignItems          string `json:"align_items,omitempty"`          // "flex-start", "center", "flex-end"
	MarginBottom        int    `json:"margin_bottom,omitempty"`        // margin bottom in pixels
}

// Responsive defines responsive breakpoints and changes
type Responsive struct {
	Mobile ResponsiveBreakpoint `json:"mobile"`
	Tablet ResponsiveBreakpoint `json:"tablet"`
}

// ResponsiveBreakpoint defines a responsive breakpoint configuration
type ResponsiveBreakpoint struct {
	Breakpoint int                    `json:"breakpoint"`
	Changes    map[string]interface{} `json:"changes"`
}

// Accessibility defines accessibility requirements
type Accessibility struct {
	TouchTargetsMin    int    `json:"touch_targets_min"`
	FocusIndicators    string `json:"focus_indicators"`
	Labels             string `json:"labels"`
	SemanticStructure  bool   `json:"semantic_structure"`
}

// Validation defines validation results
type Validation struct {
	VisualHierarchy   string `json:"visual_hierarchy"`   // "passed", "failed"
	TouchTargets      string `json:"touch_targets"`      // "passed", "failed"
	MaxNestingDepth   int    `json:"max_nesting_depth"`
	ResponsiveTested  bool   `json:"responsive_tested"`
	Notes             string `json:"notes,omitempty"`
	AspectImproved    string `json:"aspect_improved,omitempty"`
	ChecksPassed      []string `json:"checks_passed,omitempty"`
}

// ValidatePhase1 validates that the structure conforms to Phase 1 constraints
func (s *Structure) ValidatePhase1() error {
	// Check phase
	if s.Phase != "structure" {
		return fmt.Errorf("invalid phase: expected 'structure', got '%s'", s.Phase)
	}

	// Validate required fields
	if s.Version == "" {
		return fmt.Errorf("version is required")
	}
	if s.Intent.Purpose == "" {
		return fmt.Errorf("intent.purpose is required")
	}
	if s.Layout.Type == "" {
		return fmt.Errorf("layout.type is required")
	}
	if len(s.Components) == 0 {
		return fmt.Errorf("at least one component is required")
	}

	// Validate layout type
	validLayoutTypes := map[string]bool{"stack": true, "grid": true, "sidebar": true}
	if !validLayoutTypes[s.Layout.Type] {
		return fmt.Errorf("invalid layout.type: %s (must be stack, grid, or sidebar)", s.Layout.Type)
	}

	// Validate components
	for i, comp := range s.Components {
		if err := validateComponent(&comp, 0); err != nil {
			return fmt.Errorf("component[%d]: %w", i, err)
		}
	}

	return nil
}

// validateComponent recursively validates a component and its children
func validateComponent(c *Component, depth int) error {
	// Check max nesting depth
	if depth > 4 {
		return fmt.Errorf("component '%s': max nesting depth (4) exceeded", c.ID)
	}

	// Validate required fields
	if c.ID == "" {
		return fmt.Errorf("component ID is required")
	}
	if c.Type == "" {
		return fmt.Errorf("component '%s': type is required", c.ID)
	}

	// Validate component type
	validTypes := map[string]bool{"box": true, "text": true, "input": true, "button": true, "image": true}
	if !validTypes[c.Type] {
		return fmt.Errorf("component '%s': invalid type '%s' (must be box, text, input, button, or image)", c.ID, c.Type)
	}

	// Validate colors (Phase 1 constraint: only black, white, and grays)
	validColors := map[string]bool{
		"#FFFFFF": true,
		"#000000": true,
		"#E5E5E5": true,
		"#737373": true,
		"#525252": true,
	}
	
	if c.Color != "" && !validColors[c.Color] {
		return fmt.Errorf("component '%s': invalid color '%s' (Phase 1 only allows #FFFFFF, #000000, #E5E5E5, #737373, #525252)", c.ID, c.Color)
	}
	
	if c.Layout.Background != "" && !validColors[c.Layout.Background] {
		return fmt.Errorf("component '%s': invalid background color '%s' (Phase 1 only allows #FFFFFF, #000000, #E5E5E5, #737373, #525252)", c.ID, c.Layout.Background)
	}

	// Validate children recursively
	for i, child := range c.Children {
		if err := validateComponent(&child, depth+1); err != nil {
			return fmt.Errorf("component '%s'.children[%d]: %w", c.ID, i, err)
		}
	}

	return nil
}

// ParseStructure parses a JSON byte array into a Structure
func ParseStructure(data []byte) (*Structure, error) {
	var s Structure
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return &s, nil
}

// ParseAndValidateStructure parses and validates a Phase 1 structure
func ParseAndValidateStructure(data []byte) (*Structure, error) {
	s, err := ParseStructure(data)
	if err != nil {
		return nil, err
	}

	if err := s.ValidatePhase1(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return s, nil
}
