package validate

import (
	"fmt"
	"strings"

	"github.com/johanbellander/prism/internal/types"
)

// GestaltRule defines validation rules for Gestalt principles
type GestaltRule struct {
	IntraGroupSpacing int  // e.g., 8px within groups
	InterGroupSpacing int  // e.g., 24px between groups
	MinGroupSize      int  // e.g., 2 items to form a group
	SimilarityCheck   bool // Check if similar items look alike
}

// DefaultGestaltRule returns the default Gestalt validation rules
func DefaultGestaltRule() GestaltRule {
	return GestaltRule{
		IntraGroupSpacing: 8,
		InterGroupSpacing: 24,
		MinGroupSize:      2,
		SimilarityCheck:   true,
	}
}

// GestaltIssue represents a single Gestalt validation issue
type GestaltIssue struct {
	Severity  string // "error", "warning", "info"
	Message   string
	Component string // Component ID if applicable
}

// GestaltResult represents the result of Gestalt validation
type GestaltResult struct {
	Passed bool
	Issues []GestaltIssue
}

// ComponentRelationship represents the relationship between components
type ComponentRelationship struct {
	ID1      string
	ID2      string
	Spacing  int
	Related  bool // Are they likely related based on naming/type
}

// ValidateGestalt validates Gestalt principles (proximity and similarity)
func ValidateGestalt(structure *types.Structure, rule GestaltRule) GestaltResult {
	result := GestaltResult{
		Passed: true,
		Issues: []GestaltIssue{},
	}

	// Collect all sibling relationships (components at the same level)
	var collectSiblings func(parent *types.Component, siblings []types.Component) []ComponentRelationship
	collectSiblings = func(parent *types.Component, siblings []types.Component) []ComponentRelationship {
		relationships := []ComponentRelationship{}
		
		// Analyze spacing between siblings
		for i := 0; i < len(siblings); i++ {
			for j := i + 1; j < len(siblings); j++ {
				comp1 := &siblings[i]
				comp2 := &siblings[j]
				
				// Calculate spacing between adjacent components
				var spacing int
				if parent != nil {
					if parent.Layout.Direction == "vertical" {
						spacing = parent.Layout.Gap
					} else if parent.Layout.Direction == "horizontal" {
						spacing = parent.Layout.Gap
					} else {
						spacing = parent.Layout.Gap
					}
				} else {
					spacing = structure.Layout.Spacing
				}
				
				// Determine if components are likely related
				related := areComponentsRelated(comp1, comp2)
				
				relationships = append(relationships, ComponentRelationship{
					ID1:     comp1.ID,
					ID2:     comp2.ID,
					Spacing: spacing,
					Related: related,
				})
			}
			
			// Recurse into children
			if len(siblings[i].Children) > 0 {
				childRels := collectSiblings(&siblings[i], siblings[i].Children)
				relationships = append(relationships, childRels...)
			}
		}
		
		return relationships
	}
	
	// Collect all relationships
	relationships := collectSiblings(nil, structure.Components)
	
	// Add relationships from children of top-level components
	for i := range structure.Components {
		if len(structure.Components[i].Children) > 0 {
			childRels := collectSiblings(&structure.Components[i], structure.Components[i].Children)
			relationships = append(relationships, childRels...)
		}
	}
	
	// Analyze spacing patterns
	relatedPairs := []ComponentRelationship{}
	unrelatedPairs := []ComponentRelationship{}
	
	for _, rel := range relationships {
		if rel.Related {
			relatedPairs = append(relatedPairs, rel)
		} else {
			unrelatedPairs = append(unrelatedPairs, rel)
		}
	}
	
	// Check that related items have consistent, close spacing
	spacingCounts := make(map[int]int)
	for _, rel := range relatedPairs {
		spacingCounts[rel.Spacing]++
		
		if rel.Spacing > rule.IntraGroupSpacing*2 {
			result.Issues = append(result.Issues, GestaltIssue{
				Severity:  "warning",
				Message:   fmt.Sprintf("Proximity: Related components '%s' and '%s' have large spacing (%dpx) - consider reducing to %dpx for better grouping", rel.ID1, rel.ID2, rel.Spacing, rule.IntraGroupSpacing),
				Component: rel.ID1,
			})
			result.Passed = false
		}
	}
	
	// Check that unrelated items have adequate spacing
	for _, rel := range unrelatedPairs {
		if rel.Spacing < rule.InterGroupSpacing {
			result.Issues = append(result.Issues, GestaltIssue{
				Severity:  "info",
				Message:   fmt.Sprintf("Suggestion: Increase spacing to %dpx between unrelated components '%s' and '%s' (currently %dpx)", rule.InterGroupSpacing, rel.ID1, rel.ID2, rel.Spacing),
				Component: rel.ID1,
			})
		}
	}
	
	// Check for similarity in related components
	if rule.SimilarityCheck {
		groups := findComponentGroups(structure)
		
		for groupName, components := range groups {
			if len(components) >= rule.MinGroupSize {
				// Check that similar components have consistent styling
				inconsistencies := checkSimilarity(components)
				
				if len(inconsistencies) > 0 {
					for _, inconsistency := range inconsistencies {
						result.Issues = append(result.Issues, GestaltIssue{
							Severity:  "warning",
							Message:   fmt.Sprintf("Similarity: %s in group '%s' - consider using consistent styling", inconsistency, groupName),
							Component: groupName,
						})
					}
				}
			}
		}
	}
	
	// Detect potential groupings by proximity
	detectedGroups := detectGroupsByProximity(structure, rule)
	for groupID, group := range detectedGroups {
		if len(group) >= rule.MinGroupSize {
			// Find the dominant spacing within the group
			if len(group) > 0 {
				result.Issues = append(result.Issues, GestaltIssue{
					Severity:  "info",
					Message:   fmt.Sprintf("✓ Detected well-formed group '%s' with %d components using consistent spacing", groupID, len(group)),
					Component: groupID,
				})
			}
		}
	}
	
	// Add success messages if no major issues found
	if len(result.Issues) == 0 {
		result.Issues = append(result.Issues, GestaltIssue{
			Severity: "info",
			Message:  "✓ Component grouping follows Gestalt proximity principles",
		})
		
		if rule.SimilarityCheck {
			result.Issues = append(result.Issues, GestaltIssue{
				Severity: "info",
				Message:  "✓ Similar components use consistent styling",
			})
		}
	}
	
	return result
}

// areComponentsRelated determines if two components are likely related
func areComponentsRelated(comp1, comp2 *types.Component) bool {
	// Check if they share a common prefix (e.g., "username-label" and "username-input")
	id1Parts := strings.Split(comp1.ID, "-")
	id2Parts := strings.Split(comp2.ID, "-")
	
	if len(id1Parts) > 1 && len(id2Parts) > 1 {
		if id1Parts[0] == id2Parts[0] {
			return true
		}
	}
	
	// Check if they're the same type and role
	if comp1.Type == comp2.Type && comp1.Role == comp2.Role && comp1.Role != "" {
		return true
	}
	
	// Check for label-input patterns
	if (comp1.Type == "text" && comp2.Type == "input") || (comp1.Type == "input" && comp2.Type == "text") {
		// If one contains "label" and they share a prefix, they're related
		if strings.Contains(comp1.ID, "label") || strings.Contains(comp2.ID, "label") {
			return true
		}
	}
	
	return false
}

// findComponentGroups groups components by their type and role
func findComponentGroups(structure *types.Structure) map[string][]*types.Component {
	groups := make(map[string][]*types.Component)
	
	var traverse func(comp *types.Component)
	traverse = func(comp *types.Component) {
		// Group by type-role combination
		groupKey := comp.Type
		if comp.Role != "" {
			groupKey = comp.Type + "-" + comp.Role
		}
		
		groups[groupKey] = append(groups[groupKey], comp)
		
		// Recurse into children
		for i := range comp.Children {
			traverse(&comp.Children[i])
		}
	}
	
	for i := range structure.Components {
		traverse(&structure.Components[i])
	}
	
	return groups
}

// checkSimilarity checks if similar components have consistent styling
func checkSimilarity(components []*types.Component) []string {
	inconsistencies := []string{}
	
	if len(components) < 2 {
		return inconsistencies
	}
	
	// Check for consistent sizes
	sizes := make(map[string]int)
	for _, comp := range components {
		if comp.Size != "" {
			sizes[comp.Size]++
		}
	}
	
	if len(sizes) > 1 {
		inconsistencies = append(inconsistencies, "inconsistent text sizes")
	}
	
	// Check for consistent colors
	colors := make(map[string]int)
	for _, comp := range components {
		if comp.Color != "" {
			colors[comp.Color]++
		}
	}
	
	if len(colors) > 1 {
		inconsistencies = append(inconsistencies, "inconsistent colors")
	}
	
	// Check for consistent padding
	paddings := make(map[int]int)
	for _, comp := range components {
		if comp.Layout.Padding > 0 {
			paddings[comp.Layout.Padding]++
		}
	}
	
	if len(paddings) > 2 { // Allow some variation
		inconsistencies = append(inconsistencies, "inconsistent padding")
	}
	
	return inconsistencies
}

// detectGroupsByProximity detects component groups based on spacing patterns
func detectGroupsByProximity(structure *types.Structure, rule GestaltRule) map[string][]*types.Component {
	groups := make(map[string][]*types.Component)
	
	var traverse func(parent *types.Component, parentID string)
	traverse = func(parent *types.Component, parentID string) {
		if parent != nil && len(parent.Children) > 0 {
			groupID := parent.ID
			if parent.Role != "" {
				groupID = parent.Role
			}
			
			// If spacing is tight (close to intra-group spacing), consider it a group
			if parent.Layout.Gap <= rule.IntraGroupSpacing*2 {
				groups[groupID] = make([]*types.Component, len(parent.Children))
				for i := range parent.Children {
					groups[groupID][i] = &parent.Children[i]
				}
			}
			
			// Recurse into children
			for i := range parent.Children {
				traverse(&parent.Children[i], groupID)
			}
		}
	}
	
	for i := range structure.Components {
		traverse(&structure.Components[i], structure.Components[i].ID)
	}
	
	return groups
}
