package render

import (
	"strconv"
	"strings"

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

	// Check for explicit width/height in layout
	if comp.Layout.Width > 0 {
		box.Width = comp.Layout.Width * e.scale
	} else if comp.Layout.Flex > 0 {
		// Flex items take available width
		box.Width = availWidth
	} else {
		// Fallback to type-based sizing
		switch comp.Type {
		case "text":
			box.Width = availWidth
		case "button":
			box.Width = 120 * e.scale
		case "input", "box":
			box.Width = availWidth
		case "image":
			box.Width = availWidth
		default:
			box.Width = availWidth
		}
	}

	if comp.Layout.Height > 0 {
		box.Height = comp.Layout.Height * e.scale
	} else {
		// Calculate height based on component type
		switch comp.Type {
		case "text":
			box.Height = e.estimateTextHeight(comp)
		case "button":
			box.Height = 44 * e.scale
		case "input":
			box.Height = 40 * e.scale
		case "image":
			box.Height = 150 * e.scale
		case "box":
			if len(comp.Children) > 0 {
				box.Height = e.calculateContainerHeight(comp, box.Width)
			} else {
				box.Height = 100 * e.scale
			}
		default:
			box.Height = e.estimateContentHeight(comp)
		}
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
	
	// Add small default gap for vertical layouts if not specified
	if gap == 0 && direction == "vertical" {
		gap = 8 * e.scale
	}
	
	// For horizontal layouts with justify_content: space-between, we need to calculate positions differently
	if direction == "horizontal" && comp.Layout.JustifyContent == "space-between" && len(comp.Children) > 0 {
		// First pass: calculate all child boxes to get their widths
		childBoxes := make([]LayoutBox, len(comp.Children))
		totalChildWidth := 0
		
		for i, child := range comp.Children {
			// For text components, use intrinsic width instead of available width
			childWidth := width
			if child.Type == "text" {
				childWidth = e.estimateTextWidth(&child)
			}
			
			childBox, err := e.calculateComponentLayout(&child, 0, 0, childWidth, height)
			if err != nil {
				return err
			}
			childBoxes[i] = childBox
			totalChildWidth += childBox.Width
		}
		
		// Calculate spacing between items
		var spacing int
		if len(comp.Children) > 1 {
			availableSpace := width - totalChildWidth
			spacing = availableSpace / (len(comp.Children) - 1)
		}
		
		// Second pass: position children with calculated spacing
		currentX := x
		for i, child := range comp.Children {
			childBoxes[i].X = currentX
			childBoxes[i].Y = y
			boxes[child.ID] = childBoxes[i]
			
			// Recurse for grandchildren
			if err := e.calculateChildrenLayout(&child, childBoxes[i], boxes); err != nil {
				return err
			}
			
			currentX += childBoxes[i].Width + spacing
		}
		
		return nil
	}
	
	// Standard flex layout
	currentX := x
	currentY := y
	
	if direction == "horizontal" {
		// Two-pass layout for horizontal flex to handle flex-grow correctly
		// First pass: calculate fixed-width children and total flex
		fixedWidth := 0
		totalFlex := 0
		
		for _, child := range comp.Children {
			if child.Layout.Width > 0 {
				fixedWidth += child.Layout.Width * e.scale
			} else if child.Layout.Flex > 0 {
				totalFlex += child.Layout.Flex
			}
		}
		
		// Calculate available width for flex items
		availableForFlex := width - fixedWidth - (gap * (len(comp.Children) - 1))
		if availableForFlex < 0 {
			availableForFlex = 0
		}
		
		// Second pass: layout children with calculated widths
		currentX = x
		for _, child := range comp.Children {
			childWidth := width
			if child.Layout.Width > 0 {
				childWidth = child.Layout.Width * e.scale
			} else if child.Layout.Flex > 0 && totalFlex > 0 {
				childWidth = (availableForFlex * child.Layout.Flex) / totalFlex
			}
			
			childBox, err := e.calculateComponentLayout(&child, currentX, currentY, childWidth, height)
			if err != nil {
				return err
			}

			boxes[child.ID] = childBox

			// Recurse for grandchildren
			if err := e.calculateChildrenLayout(&child, childBox, boxes); err != nil {
				return err
			}

			currentX += childBox.Width + gap
		}
		
		return nil
	}
	
	// Vertical flex layout
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
		currentY += childBox.Height + gap
	}

	return nil
}

// layoutGridChildren layouts children using grid rules
func (e *LayoutEngine) layoutGridChildren(comp *types.Component, x, y, width, height int, boxes map[string]LayoutBox) error {
	gap := comp.Layout.Gap * e.scale
	
	// Parse grid_template_columns to get column widths
	columnWidths := e.parseGridColumnWidths(comp.Layout.GridTemplateColumns, width, gap)
	if len(columnWidths) == 0 {
		// Fallback to 2-column grid with equal widths
		cellWidth := (width - gap) / 2
		columnWidths = []int{cellWidth, cellWidth}
	}

	columns := len(columnWidths)
	currentX := x
	currentY := y
	col := 0
	maxRowHeight := 0

	for _, child := range comp.Children {
		cellWidth := columnWidths[col]
		
		childBox, err := e.calculateComponentLayout(&child, currentX, currentY, cellWidth, 0)
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
	// Use consistent 16px line height to match rendering
	lineHeight := 16
	
	// Count lines in content (split by newline)
	// This includes empty lines for spacing
	lines := 1
	if comp.Content != "" {
		lines = len(strings.Split(comp.Content, "\n"))
	}
	
	// Add 14px for first line baseline + (lines * lineHeight) + 8px bottom padding
	return (14 + (lines * lineHeight) + 8) * e.scale
}

// estimateTextWidth returns approximate width needed for text
func (e *LayoutEngine) estimateTextWidth(comp *types.Component) int {
	if comp.Content == "" {
		return 0
	}
	
	// Find longest line
	lines := strings.Split(comp.Content, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}
	
	// Rough approximation: 7 pixels per character (monospace-ish)
	// Adjust based on font size
	baseWidth := 7
	switch comp.Size {
	case "xs":
		baseWidth = 5
	case "sm":
		baseWidth = 6
	case "base":
		baseWidth = 7
	case "lg":
		baseWidth = 9
	case "xl":
		baseWidth = 11
	case "2xl":
		baseWidth = 14
	case "3xl":
		baseWidth = 18
	}
	
	return (maxLen * baseWidth) * e.scale
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
	
	// Add small default gap for vertical layouts if not specified
	if gap == 0 && direction == "vertical" {
		gap = 8 * e.scale
	}
	
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

// parseGridColumns parses CSS grid-template-columns value to determine number of columns
// Supports: "repeat(4, 1fr)", "1fr 1fr 1fr", "200px 1fr 1fr", etc.
func (e *LayoutEngine) parseGridColumns(gridTemplate string) int {
	if gridTemplate == "" {
		return 0
	}

	// Handle repeat() syntax: repeat(4, 1fr) -> 4 columns
	if strings.HasPrefix(gridTemplate, "repeat(") {
		// Extract the number from repeat(N, ...)
		parts := strings.TrimPrefix(gridTemplate, "repeat(")
		parts = strings.TrimSuffix(parts, ")")
		values := strings.Split(parts, ",")
		if len(values) >= 2 {
			countStr := strings.TrimSpace(values[0])
			if count, err := strconv.Atoi(countStr); err == nil && count > 0 {
				return count
			}
		}
		// Invalid repeat() syntax - don't fall through
		return 0
	}

	// Handle space-separated values: "1fr 1fr 1fr 1fr" -> 4 columns
	parts := strings.Fields(gridTemplate)
	if len(parts) > 0 {
		return len(parts)
	}

	return 0
}

// parseGridColumnWidths parses CSS grid-template-columns and returns actual pixel widths
// Supports: "repeat(4, 1fr)", "1fr 1fr 1fr", "300px 1fr 300px", etc.
func (e *LayoutEngine) parseGridColumnWidths(gridTemplate string, totalWidth, gap int) []int {
	if gridTemplate == "" {
		return nil
	}

	var columnDefs []string

	// Handle repeat() syntax: repeat(4, 1fr) -> ["1fr", "1fr", "1fr", "1fr"]
	if strings.HasPrefix(gridTemplate, "repeat(") {
		parts := strings.TrimPrefix(gridTemplate, "repeat(")
		parts = strings.TrimSuffix(parts, ")")
		values := strings.Split(parts, ",")
		if len(values) >= 2 {
			countStr := strings.TrimSpace(values[0])
			templateStr := strings.TrimSpace(values[1])
			if count, err := strconv.Atoi(countStr); err == nil && count > 0 {
				for i := 0; i < count; i++ {
					columnDefs = append(columnDefs, templateStr)
				}
			}
		}
	} else {
		// Handle space-separated values: "300px 1fr 300px" -> ["300px", "1fr", "300px"]
		columnDefs = strings.Fields(gridTemplate)
	}

	if len(columnDefs) == 0 {
		return nil
	}

	// Calculate widths
	widths := make([]int, len(columnDefs))
	fixedWidth := 0
	frCount := 0

	// First pass: calculate fixed widths and count fractional units
	for i, def := range columnDefs {
		if strings.HasSuffix(def, "px") {
			// Fixed pixel width
			pxStr := strings.TrimSuffix(def, "px")
			if px, err := strconv.Atoi(pxStr); err == nil {
				widths[i] = px * e.scale
				fixedWidth += widths[i]
			}
		} else if strings.HasSuffix(def, "fr") {
			// Fractional unit
			frStr := strings.TrimSuffix(def, "fr")
			if fr, err := strconv.Atoi(frStr); err == nil {
				frCount += fr
				widths[i] = -fr // Store negative to indicate fr unit
			} else {
				frCount++ // Default to 1fr
				widths[i] = -1
			}
		} else {
			// Unknown unit, treat as 1fr
			frCount++
			widths[i] = -1
		}
	}

	// Calculate available width for fractional units
	totalGap := gap * (len(columnDefs) - 1)
	availableForFr := totalWidth - fixedWidth - totalGap
	if availableForFr < 0 {
		availableForFr = 0
	}

	// Second pass: calculate fractional widths
	if frCount > 0 {
		for i, w := range widths {
			if w < 0 {
				// Negative value indicates fr unit
				fr := -w
				widths[i] = (availableForFr * fr) / frCount
			}
		}
	}

	return widths
}
