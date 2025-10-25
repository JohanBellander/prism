package render

import (
	"github.com/johanbellander/prism/internal/types"
)

// LayoutBox represents a calculated position and size for a component
type LayoutBox struct {
	X      int
	Y      int
	Width  int
	Height int
}

// LayoutEngine calculates layout positions for all components
type LayoutEngine struct {
	scale int
}

// NewLayoutEngine creates a new layout engine with given scale
func NewLayoutEngine(scale int) *LayoutEngine {
	return &LayoutEngine{scale: scale}
}

// CalculateLayout calculates positions and sizes for all components
func (e *LayoutEngine) CalculateLayout(structure *types.Structure, width, height int) (map[string]LayoutBox, error) {
	boxes := make(map[string]LayoutBox)

	// Calculate layout for top-level components
	currentY := 0
	for _, comp := range structure.Components {
		box, err := e.calculateComponentLayout(&comp, 0, currentY, width, height)
		if err != nil {
			return nil, err
		}

		boxes[comp.ID] = box

		// Recursively calculate children
		if err := e.calculateChildrenLayout(&comp, box, boxes); err != nil {
			return nil, err
		}

		currentY += box.Height + (structure.Layout.Spacing * e.scale)
	}

	return boxes, nil
}

// calculateComponentLayout calculates layout for a single component
func (e *LayoutEngine) calculateComponentLayout(comp *types.Component, x, y, availWidth, availHeight int) (LayoutBox, error) {
	box := LayoutBox{X: x, Y: y}

	// In Phase 1, we use simple layout rules based on component type
	switch comp.Type {
	case "text":
		box.Width = availWidth
		box.Height = e.estimateTextHeight(comp)
	case "button":
		// Buttons have minimum touch target size
		box.Width = 120 * e.scale
		box.Height = 44 * e.scale
	case "input":
		box.Width = availWidth
		box.Height = 40 * e.scale
	case "image":
		box.Width = availWidth
		box.Height = 150 * e.scale
	case "box":
		// Boxes fill available width
		box.Width = availWidth
		// Calculate height based on children
		if len(comp.Children) > 0 {
			box.Height = e.calculateContainerHeight(comp, availWidth)
		} else {
			box.Height = 100 * e.scale
		}
	default:
		box.Width = availWidth
		box.Height = e.estimateContentHeight(comp)
	}

	return box, nil
}

// calculateChildrenLayout recursively calculates layout for children
func (e *LayoutEngine) calculateChildrenLayout(comp *types.Component, parentBox LayoutBox, boxes map[string]LayoutBox) error {
	if len(comp.Children) == 0 {
		return nil
	}

	// Calculate content area (inside padding)
	padding := comp.Layout.Padding * e.scale
	contentX := parentBox.X + padding
	contentY := parentBox.Y + padding
	contentWidth := parentBox.Width - (padding * 2)
	contentHeight := parentBox.Height - (padding * 2)

	// Determine layout strategy based on display property
	display := comp.Layout.Display
	if display == "" {
		display = "flex" // default
	}

	switch display {
	case "flex":
		return e.layoutFlexChildren(comp, contentX, contentY, contentWidth, contentHeight, boxes)
	case "grid":
		return e.layoutGridChildren(comp, contentX, contentY, contentWidth, contentHeight, boxes)
	default:
		// Default to stack (vertical)
		return e.layoutStackChildren(comp, contentX, contentY, contentWidth, contentHeight, boxes)
	}
}

// layoutFlexChildren positions children using flexbox rules
func (e *LayoutEngine) layoutFlexChildren(comp *types.Component, x, y, width, height int, boxes map[string]LayoutBox) error {
	direction := comp.Layout.Direction
	if direction == "" {
		direction = "vertical"
	}

	gap := comp.Layout.Gap * e.scale
	currentX := x
	currentY := y

	// Calculate flex layout
	for _, child := range comp.Children {
		childBox, err := e.calculateComponentLayout(&child, currentX, currentY, width, height)
		if err != nil {
			return err
		}

		boxes[child.ID] = childBox

		// Recurse for grandchildren
		if err := e.calculateChildrenLayout(&child, childBox, boxes); err != nil {
			return err
		}

		// Advance position
		if direction == "vertical" {
			currentY += childBox.Height + gap
		} else {
			currentX += childBox.Width + gap
		}
	}

	return nil
}

// layoutGridChildren layouts children using grid rules
func (e *LayoutEngine) layoutGridChildren(comp *types.Component, x, y, width, height int, boxes map[string]LayoutBox) error {
	gap := comp.Layout.Gap * e.scale
	// In Phase 1, we use a simple 2-column grid by default
	columns := 2

	// Calculate cell dimensions
	cellWidth := (width - gap*(columns-1)) / columns
	
	currentX := x
	currentY := y
	col := 0
	maxRowHeight := 0

	for _, child := range comp.Children {
		childBox, err := e.calculateComponentLayout(&child, currentX, currentY, cellWidth, height)
		if err != nil {
			return err
		}

		boxes[child.ID] = childBox

		// Recurse for grandchildren
		if err := e.calculateChildrenLayout(&child, childBox, boxes); err != nil {
			return err
		}

		// Track max height in this row
		if childBox.Height > maxRowHeight {
			maxRowHeight = childBox.Height
		}

		col++
		if col >= columns {
			// Move to next row
			col = 0
			currentX = x
			currentY += maxRowHeight + gap
			maxRowHeight = 0
		} else {
			// Move to next column
			currentX += cellWidth + gap
		}
	}

	return nil
}

// layoutStackChildren layouts children in a vertical stack (default)
func (e *LayoutEngine) layoutStackChildren(comp *types.Component, x, y, width, height int, boxes map[string]LayoutBox) error {
	gap := comp.Layout.Gap * e.scale
	currentY := y

	for _, child := range comp.Children {
		childBox, err := e.calculateComponentLayout(&child, x, currentY, width, height)
		if err != nil {
			return err
		}

		boxes[child.ID] = childBox

		// Recurse for grandchildren
		if err := e.calculateChildrenLayout(&child, childBox, boxes); err != nil {
			return err
		}

		currentY += childBox.Height + gap
	}

	return nil
}

// estimateContentHeight estimates the intrinsic height of a component
func (e *LayoutEngine) estimateContentHeight(comp *types.Component) int {
	padding := comp.Layout.Padding * e.scale
	baseHeight := padding * 2

	switch comp.Type {
	case "text":
		return baseHeight + e.estimateTextHeight(comp)
	case "button":
		return baseHeight + 44*e.scale // minimum touch target
	case "input":
		return baseHeight + 40*e.scale
	case "image":
		return baseHeight + 150*e.scale
	case "box":
		return baseHeight + e.calculateContainerHeight(comp, 0)
	default:
		return baseHeight + 20*e.scale
	}
}

// estimateTextHeight returns height needed for text
func (e *LayoutEngine) estimateTextHeight(comp *types.Component) int {
	sizes := map[string]int{
		"xs":   12,
		"sm":   14,
		"base": 16,
		"md":   16,
		"lg":   18,
		"xl":   20,
		"2xl":  24,
		"3xl":  30,
		"4xl":  36,
	}

	if h, ok := sizes[comp.Size]; ok {
		return h * e.scale
	}
	return 16 * e.scale
}

// calculateContainerHeight calculates height for a container with children
func (e *LayoutEngine) calculateContainerHeight(comp *types.Component, width int) int {
	if len(comp.Children) == 0 {
		return 0
	}

	direction := comp.Layout.Direction
	if direction == "" {
		direction = "vertical"
	}

	gap := comp.Layout.Gap * e.scale
	totalHeight := 0

	if direction == "vertical" {
		// Stack children vertically
		for _, child := range comp.Children {
			totalHeight += e.estimateContentHeight(&child)
		}
		if len(comp.Children) > 1 {
			totalHeight += gap * (len(comp.Children) - 1)
		}
	} else {
		// Horizontal layout - use max child height
		maxHeight := 0
		for _, child := range comp.Children {
			h := e.estimateContentHeight(&child)
			if h > maxHeight {
				maxHeight = h
			}
		}
		totalHeight = maxHeight
	}

	return totalHeight
}
