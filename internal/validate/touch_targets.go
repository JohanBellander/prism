package validate

import (
	"fmt"
	"strings"

	"github.com/johanbellander/prism/internal/types"
)

// TouchTargetRule defines validation rules for touch targets and spacing
type TouchTargetRule struct {
	MinSize          int      // 44px (iOS) or 48px (Android)
	MinSpacing       int      // 8px between interactive elements
	DangerousSpacing int      // 16px for destructive actions
	FrequentActions  []string // IDs of common actions to check proximity
}

// DefaultTouchTargetRule returns the default touch target validation rules
func DefaultTouchTargetRule() TouchTargetRule {
	return TouchTargetRule{
		MinSize:          44, // iOS standard
		MinSpacing:       8,
		DangerousSpacing: 16,
		FrequentActions:  []string{},
	}
}

// TouchTargetIssue represents a single touch target validation issue
type TouchTargetIssue struct {
	Severity  string // "error", "warning", "info"
	Message   string
	Component string // Component ID if applicable
}

// TouchTargetResult represents the result of touch target validation
type TouchTargetResult struct {
	Passed bool
	Issues []TouchTargetIssue
}

// ComponentPosition represents a component's position and size
type ComponentPosition struct {
	ID           string
	X            int
	Y            int
	Width        int
	Height       int
	IsDangerous  bool
	IsInteractive bool
	Component    *types.Component
}

// ValidateTouchTargets validates touch targets and spacing
func ValidateTouchTargets(structure *types.Structure, rule TouchTargetRule) TouchTargetResult {
	result := TouchTargetResult{
		Passed: true,
		Issues: []TouchTargetIssue{},
	}

	// Collect all interactive elements with their positions
	positions := []ComponentPosition{}
	
	var traverse func(comp *types.Component, offsetX, offsetY int)
	traverse = func(comp *types.Component, offsetX, offsetY int) {
		isInteractive := isInteractiveElement(comp)
		
		if isInteractive {
			width := comp.Layout.Width
			height := comp.Layout.Height
			
			// If no explicit size, use minimum defaults
			if width == 0 {
				width = 100
			}
			if height == 0 {
				height = 44 // Default to minimum touch target
			}
			
			isDangerous := isDangerousAction(comp)
			
			positions = append(positions, ComponentPosition{
				ID:           comp.ID,
				X:            offsetX,
				Y:            offsetY,
				Width:        width,
				Height:       height,
				IsDangerous:  isDangerous,
				IsInteractive: true,
				Component:    comp,
			})
			
			// Validate minimum size
			if width < rule.MinSize || height < rule.MinSize {
				result.Issues = append(result.Issues, TouchTargetIssue{
					Severity:  "error",
					Message:   fmt.Sprintf("Touch Target: '%s' is %dx%dpx (requires %dx%dpx minimum)", comp.ID, width, height, rule.MinSize, rule.MinSize),
					Component: comp.ID,
				})
				result.Passed = false
			}
		}
		
		// Recurse into children with updated offsets
		childOffsetY := offsetY
		childOffsetX := offsetX
		
		for i := range comp.Children {
			child := &comp.Children[i]
			
			// Update offsets based on layout direction
			if comp.Layout.Direction == "vertical" {
				childOffsetY += comp.Layout.Gap
			} else if comp.Layout.Direction == "horizontal" {
				childOffsetX += comp.Layout.Gap
			}
			
			traverse(child, childOffsetX, childOffsetY)
			
			// Update offset for next sibling
			if comp.Layout.Direction == "vertical" {
				childOffsetY += child.Layout.Height
			} else if comp.Layout.Direction == "horizontal" {
				childOffsetX += child.Layout.Width
			}
		}
	}
	
	// Start traversal
	startY := 0
	for i := range structure.Components {
		traverse(&structure.Components[i], 0, startY)
		startY += structure.Components[i].Layout.Height + structure.Layout.Spacing
	}
	
	// Check spacing between interactive elements
	for i := 0; i < len(positions); i++ {
		for j := i + 1; j < len(positions); j++ {
			pos1 := positions[i]
			pos2 := positions[j]
			
			// Calculate spacing between elements
			spacing := calculateSpacing(pos1, pos2)
			
			// Determine required spacing
			requiredSpacing := rule.MinSpacing
			if pos1.IsDangerous || pos2.IsDangerous {
				requiredSpacing = rule.DangerousSpacing
			}
			
			// Check if spacing is adequate
			if spacing >= 0 && spacing < requiredSpacing {
				severity := "warning"
				actionType := "interactive elements"
				
				if pos1.IsDangerous || pos2.IsDangerous {
					severity = "error"
					actionType = "destructive action"
					result.Passed = false
				}
				
				result.Issues = append(result.Issues, TouchTargetIssue{
					Severity:  severity,
					Message:   fmt.Sprintf("Spacing: '%s' only %dpx from '%s' (requires %dpx for %s)", pos1.ID, spacing, pos2.ID, requiredSpacing, actionType),
					Component: pos1.ID,
				})
			}
		}
	}
	
	// Check frequent actions proximity (if specified)
	for _, freqAction := range rule.FrequentActions {
		var freqPos *ComponentPosition
		for i := range positions {
			if positions[i].ID == freqAction {
				freqPos = &positions[i]
				break
			}
		}
		
		if freqPos != nil {
			// Check if it's easily accessible (not too far from common interaction areas)
			// This is a basic check - could be enhanced with more sophisticated heuristics
			if freqPos.Y > 600 { // More than 600px down might be hard to reach
				result.Issues = append(result.Issues, TouchTargetIssue{
					Severity:  "info",
					Message:   fmt.Sprintf("Frequent action '%s' may be hard to reach (positioned at Y=%dpx)", freqAction, freqPos.Y),
					Component: freqAction,
				})
			}
		}
	}
	
	// Add success messages if no issues found
	if len(result.Issues) == 0 {
		result.Issues = append(result.Issues, TouchTargetIssue{
			Severity: "info",
			Message:  "✓ All interactive elements meet touch target requirements",
		})
		if len(positions) > 1 {
			result.Issues = append(result.Issues, TouchTargetIssue{
				Severity: "info",
				Message:  "✓ Spacing between interactive elements is adequate",
			})
		}
	}
	
	return result
}

// isInteractiveElement checks if a component is interactive
func isInteractiveElement(comp *types.Component) bool {
	interactiveTypes := map[string]bool{
		"button": true,
		"input":  true,
	}
	
	return interactiveTypes[comp.Type]
}

// isDangerousAction checks if a component represents a dangerous/destructive action
func isDangerousAction(comp *types.Component) bool {
	idLower := strings.ToLower(comp.ID)
	roleLower := strings.ToLower(comp.Role)
	
	dangerousKeywords := []string{"delete", "remove", "destroy", "clear", "reset", "cancel"}
	
	for _, keyword := range dangerousKeywords {
		if strings.Contains(idLower, keyword) || strings.Contains(roleLower, keyword) {
			return true
		}
	}
	
	return false
}

// calculateSpacing calculates the minimum spacing between two components
func calculateSpacing(pos1, pos2 ComponentPosition) int {
	// Calculate edges
	left1 := pos1.X
	right1 := pos1.X + pos1.Width
	top1 := pos1.Y
	bottom1 := pos1.Y + pos1.Height
	
	left2 := pos2.X
	right2 := pos2.X + pos2.Width
	top2 := pos2.Y
	bottom2 := pos2.Y + pos2.Height
	
	// Check if boxes overlap
	if right1 <= left2 {
		// pos1 is to the left of pos2
		horizontalGap := left2 - right1
		// Check if they're on similar vertical positions
		if bottom1 > top2 && top1 < bottom2 {
			return horizontalGap
		}
	}
	
	if right2 <= left1 {
		// pos2 is to the left of pos1
		horizontalGap := left1 - right2
		if bottom2 > top1 && top2 < bottom1 {
			return horizontalGap
		}
	}
	
	if bottom1 <= top2 {
		// pos1 is above pos2
		verticalGap := top2 - bottom1
		// Check if they're on similar horizontal positions
		if right1 > left2 && left1 < right2 {
			return verticalGap
		}
	}
	
	if bottom2 <= top1 {
		// pos2 is above pos1
		verticalGap := top1 - bottom2
		if right2 > left1 && left2 < right1 {
			return verticalGap
		}
	}
	
	// If they don't align horizontally or vertically, calculate diagonal distance
	// For simplicity, return -1 to indicate they're not adjacent
	return -1
}
