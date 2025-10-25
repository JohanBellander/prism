package validate

import (
	"fmt"
	"strings"

	"github.com/johanbellander/prism/internal/types"
)

// SuggestionCategory represents a category of design suggestions
type SuggestionCategory string

const (
	CategoryForms      SuggestionCategory = "forms"
	CategoryNavigation SuggestionCategory = "navigation"
	CategoryLayouts    SuggestionCategory = "layouts"
	CategoryButtons    SuggestionCategory = "buttons"
	CategoryCards      SuggestionCategory = "cards"
	CategoryTables     SuggestionCategory = "tables"
	CategoryModals     SuggestionCategory = "modals"
	CategoryAll        SuggestionCategory = "all"
)

// Suggestion represents a design best practice recommendation
type Suggestion struct {
	Category    string `json:"category"`
	Type        string `json:"type"` // "good", "consider", "suggestion"
	Message     string `json:"message"`
	ComponentID string `json:"component_id,omitempty"`
}

// SuggestionResult contains all suggestions for a structure
type SuggestionResult struct {
	Categories map[string][]Suggestion `json:"categories"`
	Total      int                     `json:"total"`
}

// GenerateSuggestions analyzes a structure and provides best practice suggestions
func GenerateSuggestions(structure *types.Structure, category SuggestionCategory) *SuggestionResult {
	result := &SuggestionResult{
		Categories: make(map[string][]Suggestion),
	}

	if category == CategoryAll || category == CategoryForms {
		formSuggestions := analyzeFormPatterns(structure)
		if len(formSuggestions) > 0 {
			result.Categories["forms"] = formSuggestions
			result.Total += len(formSuggestions)
		}
	}

	if category == CategoryAll || category == CategoryNavigation {
		navSuggestions := analyzeNavigationPatterns(structure)
		if len(navSuggestions) > 0 {
			result.Categories["navigation"] = navSuggestions
			result.Total += len(navSuggestions)
		}
	}

	if category == CategoryAll || category == CategoryLayouts {
		layoutSuggestions := analyzeLayoutPatterns(structure)
		if len(layoutSuggestions) > 0 {
			result.Categories["layouts"] = layoutSuggestions
			result.Total += len(layoutSuggestions)
		}
	}

	if category == CategoryAll || category == CategoryButtons {
		buttonSuggestions := analyzeButtonPatterns(structure)
		if len(buttonSuggestions) > 0 {
			result.Categories["buttons"] = buttonSuggestions
			result.Total += len(buttonSuggestions)
		}
	}

	if category == CategoryAll || category == CategoryCards {
		cardSuggestions := analyzeCardPatterns(structure)
		if len(cardSuggestions) > 0 {
			result.Categories["cards"] = cardSuggestions
			result.Total += len(cardSuggestions)
		}
	}

	if category == CategoryAll || category == CategoryTables {
		tableSuggestions := analyzeTablePatterns(structure)
		if len(tableSuggestions) > 0 {
			result.Categories["tables"] = tableSuggestions
			result.Total += len(tableSuggestions)
		}
	}

	if category == CategoryAll || category == CategoryModals {
		modalSuggestions := analyzeModalPatterns(structure)
		if len(modalSuggestions) > 0 {
			result.Categories["modals"] = modalSuggestions
			result.Total += len(modalSuggestions)
		}
	}

	return result
}

// analyzeFormPatterns provides suggestions for form components
func analyzeFormPatterns(structure *types.Structure) []Suggestion {
	var suggestions []Suggestion
	
	formComponents := findComponentsByType(structure, "form", "input", "text_input", "select", "checkbox", "radio")
	
	if len(formComponents) == 0 {
		return suggestions
	}

	// Check for label placement
	labelsAbove := 0
	labelsLeft := 0
	inputsWithoutLabels := []string{}

	for _, comp := range formComponents {
		if isInputField(comp.Type) {
			// In this structure, we'll check for text components that might be labels
			// by looking for text elements near inputs (in parent-child or sibling relationships)
			labelFound := false
			
			// Check if there are text children
			for _, child := range comp.Children {
				if child.Type == "text" || child.Type == "label" {
					labelFound = true
					labelsAbove++
					break
				}
			}
			
			if !labelFound {
				// Check siblings in parent containers
				for _, other := range structure.Components {
					if (other.Type == "text" || other.Type == "label") && other.ID != comp.ID {
						// Heuristic: if it's a text component, consider it a potential label
						labelFound = true
						labelsAbove++
						break
					}
				}
			}
			
			if !labelFound {
				inputsWithoutLabels = append(inputsWithoutLabels, comp.ID)
			}
		}
	}

	// Report label placement pattern
	if labelsAbove > labelsLeft {
		suggestions = append(suggestions, Suggestion{
			Category: "forms",
			Type:     "good",
			Message:  "Labels are above inputs (good for mobile and scanning)",
		})
	} else if labelsLeft > 0 {
		suggestions = append(suggestions, Suggestion{
			Category: "forms",
			Type:     "consider",
			Message:  "Labels are beside inputs. Consider placing above for better mobile experience",
		})
	}

	// Check for missing labels
	if len(inputsWithoutLabels) > 0 {
		suggestions = append(suggestions, Suggestion{
			Category:    "forms",
			Type:        "suggestion",
			Message:     fmt.Sprintf("Add labels for inputs: %s", strings.Join(inputsWithoutLabels, ", ")),
			ComponentID: inputsWithoutLabels[0],
		})
	}

	// Check field grouping
	if len(formComponents) > 5 {
		suggestions = append(suggestions, Suggestion{
			Category: "forms",
			Type:     "suggestion",
			Message:  fmt.Sprintf("%d form fields detected. Consider grouping related fields with spacing (24-32px between groups)", len(formComponents)),
		})
	}

	// Check for help text
	hasHelpText := false
	for _, comp := range structure.Components {
		if comp.Type == "text" && comp.Size != "" {
			// Small text sizes like "xs" or "sm" typically indicate help text
			if comp.Size == "xs" || comp.Size == "sm" {
				hasHelpText = true
				break
			}
		}
	}
	
	if !hasHelpText && len(formComponents) > 3 {
		suggestions = append(suggestions, Suggestion{
			Category: "forms",
			Type:     "consider",
			Message:  "Add field descriptions or help text for complex inputs (font-size: 12-13px, color: text.secondary)",
		})
	}

	return suggestions
}

// analyzeNavigationPatterns provides suggestions for navigation components
func analyzeNavigationPatterns(structure *types.Structure) []Suggestion {
	var suggestions []Suggestion
	
	navComponents := findComponentsByType(structure, "nav", "navbar", "menu", "navigation", "header")
	
	if len(navComponents) == 0 {
		return suggestions
	}

	// Check navigation location (header components are typically at top)
	for _, nav := range navComponents {
		if nav.Role == "header" || nav.Role == "navigation" {
			suggestions = append(suggestions, Suggestion{
				Category:    "navigation",
				Type:        "good",
				Message:     "Primary navigation is in expected location (header/top)",
				ComponentID: nav.ID,
			})
			break
		}
	}

	// Count navigation items by checking children
	navItemCount := 0
	for _, nav := range navComponents {
		navItemCount += countNavigationItems(nav)
	}
	
	if navItemCount > 7 {
		suggestions = append(suggestions, Suggestion{
			Category: "navigation",
			Type:     "consider",
			Message:  fmt.Sprintf("%d navigation items detected. Consider dropdown menus or grouping for less common items (optimal: 5-7 items)", navItemCount),
		})
	}

	// Check for active state indicators
	hasActiveState := false
	for _, nav := range navComponents {
		for _, child := range nav.Children {
			if child.Layout.Background != "" || child.Weight == "bold" {
				hasActiveState = true
				break
			}
		}
		if hasActiveState {
			break
		}
	}

	if !hasActiveState {
		suggestions = append(suggestions, Suggestion{
			Category: "navigation",
			Type:     "suggestion",
			Message:  "Add visual indicator for current/active page (background color, underline, or bold text)",
		})
	}

	return suggestions
}

// Helper to count navigation items
func countNavigationItems(comp types.Component) int {
	count := 0
	for _, child := range comp.Children {
		if child.Type == "text" || child.Type == "link" || child.Type == "button" {
			count++
		}
		count += countNavigationItems(child)
	}
	return count
}

// analyzeLayoutPatterns provides suggestions for layout and grid systems
func analyzeLayoutPatterns(structure *types.Structure) []Suggestion {
	var suggestions []Suggestion

	// Check for grid usage
	gridComponents := 0
	for _, comp := range structure.Components {
		if comp.Layout.Display == "grid" {
			gridComponents++
		}
	}

	if gridComponents > 0 {
		suggestions = append(suggestions, Suggestion{
			Category: "layouts",
			Type:     "good",
			Message:  "Layout uses CSS Grid for consistent structure",
		})
	} else if len(structure.Components) > 5 {
		suggestions = append(suggestions, Suggestion{
			Category: "layouts",
			Type:     "suggestion",
			Message:  "Consider using CSS Grid (display: grid) for consistent alignment",
		})
	}

	// Check for container max widths
	if structure.Layout.MaxWidth > 1440 {
		suggestions = append(suggestions, Suggestion{
			Category: "layouts",
			Type:     "consider",
			Message:  fmt.Sprintf("Max width is %dpx. Consider constraining to 1280-1440px for better readability", structure.Layout.MaxWidth),
		})
	} else if structure.Layout.MaxWidth > 0 {
		suggestions = append(suggestions, Suggestion{
			Category: "layouts",
			Type:     "good",
			Message:  fmt.Sprintf("Layout uses appropriate max-width (%dpx)", structure.Layout.MaxWidth),
		})
	}

	return suggestions
}

// analyzeButtonPatterns provides suggestions for button components
func analyzeButtonPatterns(structure *types.Structure) []Suggestion {
	var suggestions []Suggestion
	
	buttons := findComponentsByType(structure, "button", "cta", "action")
	
	if len(buttons) == 0 {
		return suggestions
	}

	// Check button sizing (min 44x44 for touch)
	smallButtons := []string{}
	for _, btn := range buttons {
		if btn.Layout.Width < 44 || btn.Layout.Height < 44 {
			smallButtons = append(smallButtons, btn.ID)
		}
	}

	if len(smallButtons) == 0 {
		suggestions = append(suggestions, Suggestion{
			Category: "buttons",
			Type:     "good",
			Message:  "All buttons meet minimum touch target size (44x44px)",
		})
	} else {
		suggestions = append(suggestions, Suggestion{
			Category:    "buttons",
			Type:        "suggestion",
			Message:     fmt.Sprintf("Increase size of buttons to minimum 44x44px: %s", strings.Join(smallButtons, ", ")),
			ComponentID: smallButtons[0],
		})
	}

	// Check for primary vs secondary buttons
	primaryButtons := 0
	for _, btn := range buttons {
		if strings.Contains(strings.ToLower(btn.ID), "primary") ||
		   strings.Contains(strings.ToLower(btn.Type), "primary") {
			primaryButtons++
		}
	}

	if primaryButtons > 1 {
		suggestions = append(suggestions, Suggestion{
			Category: "buttons",
			Type:     "consider",
			Message:  fmt.Sprintf("%d primary buttons detected. Use only 1 primary button per section for clear CTA hierarchy", primaryButtons),
		})
	}

	return suggestions
}

// analyzeCardPatterns provides suggestions for card components
func analyzeCardPatterns(structure *types.Structure) []Suggestion {
	var suggestions []Suggestion
	
	cards := findComponentsByType(structure, "card", "panel", "box")
	
	if len(cards) == 0 {
		return suggestions
	}

	// Check for consistent gap/spacing
	if len(cards) > 1 {
		// Check parent containers for gap property
		for _, comp := range structure.Components {
			if comp.Layout.Display == "grid" || comp.Layout.Display == "flex" {
				if comp.Layout.Gap > 0 {
					suggestions = append(suggestions, Suggestion{
						Category: "cards",
						Type:     "good",
						Message:  fmt.Sprintf("Cards use consistent spacing (gap: %dpx)", comp.Layout.Gap),
					})
					break
				}
			}
		}
	}

	// Check for shadows/borders
	cardsWithElevation := 0
	for _, card := range cards {
		if card.Layout.Border != "" {
			cardsWithElevation++
		}
	}

	if cardsWithElevation == len(cards) {
		suggestions = append(suggestions, Suggestion{
			Category: "cards",
			Type:     "good",
			Message:  "Cards use borders for visual separation",
		})
	} else if cardsWithElevation == 0 {
		suggestions = append(suggestions, Suggestion{
			Category: "cards",
			Type:     "consider",
			Message:  "Add subtle border to cards for visual separation (e.g., border: 1px solid #E5E5E5)",
		})
	}

	return suggestions
}

// analyzeTablePatterns provides suggestions for table components
func analyzeTablePatterns(structure *types.Structure) []Suggestion {
	var suggestions []Suggestion
	
	tables := findComponentsByType(structure, "table", "datagrid", "list")
	
	if len(tables) == 0 {
		return suggestions
	}

	// Check for headers with bold weight
	hasHeaders := false
	for _, comp := range structure.Components {
		if strings.Contains(strings.ToLower(comp.Role), "header") ||
		   strings.Contains(strings.ToLower(comp.ID), "header") {
			if comp.Weight == "bold" {
				hasHeaders = true
				break
			}
		}
	}

	if hasHeaders {
		suggestions = append(suggestions, Suggestion{
			Category: "tables",
			Type:     "good",
			Message:  "Table includes clear headers with appropriate weight",
		})
	} else {
		suggestions = append(suggestions, Suggestion{
			Category: "tables",
			Type:     "suggestion",
			Message:  "Add table headers with bold text (weight: bold) for better scannability",
		})
	}

	// Suggest sorting indicators
	suggestions = append(suggestions, Suggestion{
		Category: "tables",
		Type:     "consider",
		Message:  "Add sorting indicators (arrows) to sortable columns",
	})

	return suggestions
}

// analyzeModalPatterns provides suggestions for modal/dialog components
func analyzeModalPatterns(structure *types.Structure) []Suggestion {
	var suggestions []Suggestion
	
	modals := findComponentsByType(structure, "modal", "dialog", "popup", "overlay")
	
	if len(modals) == 0 {
		return suggestions
	}

	// Check for backdrop/overlay
	hasBackdrop := false
	for _, comp := range structure.Components {
		if strings.Contains(strings.ToLower(comp.Type), "overlay") ||
		   strings.Contains(strings.ToLower(comp.Role), "backdrop") {
			hasBackdrop = true
			break
		}
	}

	if hasBackdrop {
		suggestions = append(suggestions, Suggestion{
			Category: "modals",
			Type:     "good",
			Message:  "Modal includes backdrop/overlay for focus",
		})
	} else {
		suggestions = append(suggestions, Suggestion{
			Category: "modals",
			Type:     "suggestion",
			Message:  "Add semi-transparent backdrop (e.g., background: rgba(0,0,0,0.5)) to focus attention on modal",
		})
	}

	// Check for close button
	hasCloseButton := false
	for _, modal := range modals {
		for _, child := range modal.Children {
			if strings.Contains(strings.ToLower(child.ID), "close") ||
			   strings.Contains(strings.ToLower(child.Type), "close") {
				hasCloseButton = true
				break
			}
		}
		if hasCloseButton {
			break
		}
	}

	if !hasCloseButton {
		suggestions = append(suggestions, Suggestion{
			Category: "modals",
			Type:     "suggestion",
			Message:  "Add close button (X) in top-right corner for easy dismissal",
		})
	}

	return suggestions
}

// Helper functions

func findComponentsByType(structure *types.Structure, compTypes ...string) []types.Component {
	var result []types.Component
	for _, comp := range structure.Components {
		for _, t := range compTypes {
			if strings.Contains(strings.ToLower(comp.Type), strings.ToLower(t)) ||
			   strings.Contains(strings.ToLower(comp.ID), strings.ToLower(t)) {
				result = append(result, comp)
				break
			}
		}
	}
	return result
}

func isInputField(compType string) bool {
	inputTypes := []string{"input", "text_input", "select", "textarea", "checkbox", "radio"}
	lowerType := strings.ToLower(compType)
	for _, t := range inputTypes {
		if strings.Contains(lowerType, t) {
			return true
		}
	}
	return false
}

func findNestedComponents(structure *types.Structure, parents []types.Component) []types.Component {
	var result []types.Component
	for _, parent := range parents {
		result = append(result, getAllChildren(parent)...)
	}
	return result
}

func getAllChildren(comp types.Component) []types.Component {
	var result []types.Component
	result = append(result, comp.Children...)
	for _, child := range comp.Children {
		result = append(result, getAllChildren(child)...)
	}
	return result
}
