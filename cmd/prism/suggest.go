package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/johanbellander/prism/internal/types"
	"github.com/johanbellander/prism/internal/validate"
	"github.com/spf13/cobra"
)

var suggestCmd = &cobra.Command{
	Use:   "suggest [project-path]",
	Short: "Get proactive design improvement suggestions",
	Long: `Analyze design patterns and provide actionable best practice suggestions.

Goes beyond validation to suggest proactive improvements and optimizations.

Suggestion Categories:

  hierarchy          Visual hierarchy improvements
                     • Heading emphasis, call-to-action prominence
                     • Size ratio recommendations, visual flow

  accessibility      Accessibility enhancements beyond WCAG minimums
                     • Alt text improvements, ARIA attributes
                     • Keyboard navigation, screen reader optimization

  consistency        Pattern consistency across components
                     • Repeated patterns that should be shared components
                     • Spacing and sizing inconsistencies

  performance        Loading and performance optimizations
                     • Lazy loading opportunities, image optimization
                     • Skeleton screen recommendations

  responsiveness     Responsive design improvements
                     • Breakpoint recommendations, flexible layouts
                     • Mobile-first considerations

  microinteractions  Animation and feedback suggestions
                     • Hover states, focus transitions
                     • Loading indicators, success confirmations

  errorprevention    Error prevention and validation
                     • Confirming destructive actions
                     • Inline validation, helpful error messages

Suggestion Output Format:
  {
    "category": "hierarchy",
    "suggestions": [
      {
        "component": "pricing-cards",
        "priority": "medium",                    // low, medium, high
        "suggestion": "Highlight recommended tier",
        "rationale": "Users need guidance...",
        "implementation": "Add recommended: true flag"
      }
    ]
  }

Examples:
  # Get all suggestions
  prism suggest ./my-dashboard

  # Focus on specific category
  prism suggest ./my-dashboard --category hierarchy
  prism suggest ./my-dashboard --category accessibility
  prism suggest ./my-dashboard --category performance

  # Show all categories (same as no flag)
  prism suggest ./my-dashboard --all

  # Get JSON output for tooling integration
  prism suggest ./my-dashboard --json

  # Combine with validation
  prism audit ./my-dashboard && prism suggest ./my-dashboard

Available categories: hierarchy, accessibility, consistency, performance,
                      responsiveness, microinteractions, errorprevention

For validation rules, see: VALIDATION_RULES.md
For design tokens, see: DESIGN_TOKENS.md`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSuggest,
}

func init() {
	suggestCmd.Flags().String("category", "", "Specific category (forms, navigation, layouts, buttons, cards, tables, modals)")
	suggestCmd.Flags().Bool("all", false, "Show suggestions for all categories")
}

func runSuggest(cmd *cobra.Command, args []string) error {
	// Get flags
	projectPath := "./"
	if len(args) > 0 {
		projectPath = args[0]
	}

	categoryFlag, _ := cmd.Flags().GetString("category")
	showAll, _ := cmd.Flags().GetBool("all")
	outputJSON, _ := cmd.Parent().PersistentFlags().GetBool("json")

	// Determine category
	var category validate.SuggestionCategory
	if showAll {
		category = validate.CategoryAll
	} else if categoryFlag != "" {
		category = validate.SuggestionCategory(categoryFlag)
	} else {
		category = validate.CategoryAll
	}

	// Find the structure file
	structurePath := filepath.Join(projectPath, "phase1-structure")
	
	var structureFile string
	if _, err := os.Stat(filepath.Join(structurePath, "approved.json")); err == nil {
		structureFile = filepath.Join(structurePath, "approved.json")
	} else {
		// Find latest version
		files, err := filepath.Glob(filepath.Join(structurePath, "v*.json"))
		if err != nil || len(files) == 0 {
			if outputJSON {
				result := map[string]interface{}{
					"status": "error",
					"error":  "No structure files found in " + structurePath,
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(result)
			}
			return fmt.Errorf("no structure files found in %s", structurePath)
		}
		structureFile = files[len(files)-1]
	}

	// Load and parse the structure
	data, err := os.ReadFile(structureFile)
	if err != nil {
		if outputJSON {
			result := map[string]interface{}{
				"status": "error",
				"error":  fmt.Sprintf("Failed to read file: %v", err),
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		return fmt.Errorf("failed to read file: %w", err)
	}

	var structure types.Structure
	if err := json.Unmarshal(data, &structure); err != nil {
		if outputJSON {
			result := map[string]interface{}{
				"status": "error",
				"error":  fmt.Sprintf("Failed to parse JSON: %v", err),
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Generate suggestions
	result := validate.GenerateSuggestions(&structure, category)

	// Output results
	if outputJSON {
		output := map[string]interface{}{
			"file":        structureFile,
			"version":     structure.Version,
			"phase":       structure.Phase,
			"suggestions": result,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(output)
	}

	// Console output
	fmt.Printf("💡 Design Suggestions for %s\n", structureFile)
	fmt.Printf("   Version: %s\n", structure.Version)
	fmt.Printf("   Phase: %s\n", structure.Phase)
	fmt.Printf("   Components: %d\n\n", len(structure.Components))

	if result.Total == 0 {
		fmt.Println("✨ No suggestions found - design looks good!")
		return nil
	}

	fmt.Printf("Found %d suggestion(s) across %d categor(ies)\n\n", result.Total, len(result.Categories))
	fmt.Println("═══════════════════════════════════════════════════════")

	// Print suggestions by category
	categories := []string{"forms", "navigation", "layouts", "buttons", "cards", "tables", "modals"}
	
	for _, cat := range categories {
		suggestions, exists := result.Categories[cat]
		if !exists || len(suggestions) == 0 {
			continue
		}

		// Print category header
		icon := getCategoryIcon(cat)
		fmt.Printf("\n%s %s Best Practices:\n", icon, formatCategoryName(cat))

		// Group by type
		goods := []validate.Suggestion{}
		considers := []validate.Suggestion{}
		suggestionList := []validate.Suggestion{}

		for _, s := range suggestions {
			switch s.Type {
			case "good":
				goods = append(goods, s)
			case "consider":
				considers = append(considers, s)
			case "suggestion":
				suggestionList = append(suggestionList, s)
			}
		}

		// Print goods
		for _, s := range goods {
			fmt.Printf("   ✅ %s\n", s.Message)
		}

		// Print considers
		for _, s := range considers {
			fmt.Printf("   💭 %s\n", s.Message)
		}

		// Print suggestions
		for _, s := range suggestionList {
			fmt.Printf("   💡 Suggestion: %s\n", s.Message)
		}
	}

	fmt.Println("\n═══════════════════════════════════════════════════════")
	fmt.Println("\nRun with --category to focus on specific areas:")
	fmt.Println("  prism suggest --category forms")
	fmt.Println("  prism suggest --category navigation")
	fmt.Println("  prism suggest --category layouts")

	return nil
}

func getCategoryIcon(category string) string {
	icons := map[string]string{
		"forms":      "📝",
		"navigation": "🧭",
		"layouts":    "📐",
		"buttons":    "🔘",
		"cards":      "🃏",
		"tables":     "📊",
		"modals":     "🗨️",
	}
	if icon, ok := icons[category]; ok {
		return icon
	}
	return "📋"
}

func formatCategoryName(category string) string {
	names := map[string]string{
		"forms":      "Forms",
		"navigation": "Navigation",
		"layouts":    "Layouts",
		"buttons":    "Buttons",
		"cards":      "Cards",
		"tables":     "Tables",
		"modals":     "Modals",
	}
	if name, ok := names[category]; ok {
		return name
	}
	return category
}
