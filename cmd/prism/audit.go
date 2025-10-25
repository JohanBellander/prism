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

var auditCmd = &cobra.Command{
	Use:   "audit [project-path]",
	Short: "Run comprehensive design audit with all validations",
	Long: `Run all available design principle validations and generate an audit report.

This command runs ALL validators appropriate for the current phase:

Phase 1 Validators (Structural):
  ✓ Visual Hierarchy       - Heading scale, nesting depth, primary action prominence
  ✓ Touch Targets          - 44x44px minimum, Fitts's Law compliance
  ✓ Gestalt Principles     - Proximity, similarity, continuity
  ✓ Accessibility (WCAG)   - Labels, heading order, semantic structure
  ✓ Choice Overload        - Hick's Law (max 7 nav items, 5 form fields)

Phase 2 Validators (Visual Design):
  ✓ Color Contrast         - WCAG AA (4.5:1 text, 3:1 large text/UI)
  ✓ Typography Scale       - Consistent ratios, 8-10 sizes
  ✓ Spacing (8pt Grid)     - Multiples of 4 or 8 pixels
  ✓ Shadow & Elevation     - 3-4 elevation levels, appropriate usage
  ✓ Loading States         - Indicators, skeleton screens, feedback
  ✓ Responsive Design      - Mobile, tablet, desktop breakpoints
  ✓ Focus Indicators       - 2px outline, 3:1 contrast minimum
  ✓ Dark Mode Support      - Separate palette, maintained contrast

Audit Report Structure:
  {
    "overall_score": 85,           // 0-100 aggregate score
    "validators": [
      {
        "name": "visual_hierarchy",
        "score": 90,                 // 0-100 per validator
        "status": "passed",          // passed/failed
        "issues": [...]              // Array of issues found
      }
    ],
    "summary": {
      "total_validators": 13,
      "passed": 12,
      "failed": 1,
      "critical_issues": 0,
      "warnings": 3
    }
  }

Scoring System:
  90-100  Excellent - Production ready, exceeds standards
  70-89   Good      - Passing, minor improvements possible
  50-69   Fair      - Needs attention, several issues
  0-49    Poor      - Failing, significant problems

Examples:
  # Run full audit on Phase 1 structure
  prism audit ./my-dashboard

  # Get JSON output for CI/CD pipeline
  prism audit ./my-dashboard --json

  # Audit Phase 2 design (includes all Phase 1 + Phase 2 validators)
  prism audit ./my-dashboard --phase 2

  # Audit specific version
  prism audit ./my-dashboard --version v2

For individual validators, use: prism validate ./my-dashboard --hierarchy
For documentation, see: VALIDATION_RULES.md, TESTING_STRATEGY.md`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAudit,
}

func init() {
	auditCmd.Flags().Int("phase", 1, "Phase to validate against (1 or 2)")
}

func runAudit(cmd *cobra.Command, args []string) error {
	// Get flags
	projectPath := "./"
	if len(args) > 0 {
		projectPath = args[0]
	}

	phase, _ := cmd.Flags().GetInt("phase")
	outputJSON, _ := cmd.Parent().PersistentFlags().GetBool("json")

	// Only Phase 1 validation is currently supported
	if phase != 1 {
		if outputJSON {
			result := map[string]interface{}{
				"status": "error",
				"error":  fmt.Sprintf("Phase %d validation not yet implemented", phase),
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		return fmt.Errorf("phase %d validation not yet implemented", phase)
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

	// Run all validations
	hierarchyResult := validate.ValidateHierarchy(&structure, validate.DefaultHierarchyRule())
	touchTargetsResult := validate.ValidateTouchTargets(&structure, validate.DefaultTouchTargetRule())
	gestaltResult := validate.ValidateGestalt(&structure, validate.DefaultGestaltRule())
	a11yResult := validate.ValidateAccessibility(&structure, validate.DefaultA11yRule())
	choiceResult := validate.ValidateChoiceOverload(&structure, validate.DefaultChoiceRule())
	contrastResult := validate.ValidateContrast(&structure, validate.DefaultContrastRule())
	spacingResult := validate.ValidateSpacing(&structure, validate.DefaultSpacingRule())
	typographyResult := validate.ValidateTypography(&structure, validate.DefaultTypographyRule())
	elevationResult := validate.ValidateElevation(&structure, validate.DefaultElevationRule())
	loadingStatesResult := validate.ValidateLoadingStates(&structure, validate.DefaultLoadingStateRule())
	responsiveResult := validate.ValidateResponsive(&structure, validate.DefaultResponsiveRule())
	focusResult := validate.ValidateFocus(&structure, validate.DefaultFocusRule())
	darkModeResult := validate.ValidateDarkMode(&structure, validate.DefaultDarkModeRule())

	// Calculate overall pass/fail
	allPassed := hierarchyResult.Passed && touchTargetsResult.Passed && gestaltResult.Passed &&
		a11yResult.Passed && choiceResult.Passed && contrastResult.Passed &&
		spacingResult.Passed && typographyResult.Passed && elevationResult.Passed &&
		loadingStatesResult.Passed && responsiveResult.Passed && focusResult.Passed &&
		darkModeResult.Passed

	if outputJSON {
		result := map[string]interface{}{
			"file":       structureFile,
			"version":    structure.Version,
			"phase":      structure.Phase,
			"status":     func() string { if allPassed { return "passed" } else { return "failed" } }(),
			"components": len(structure.Components),
			"audits": map[string]interface{}{
				"hierarchy": map[string]interface{}{
					"status": func() string { if hierarchyResult.Passed { return "passed" } else { return "failed" } }(),
					"issues": hierarchyResult.Issues,
				},
				"touch_targets": map[string]interface{}{
					"status": func() string { if touchTargetsResult.Passed { return "passed" } else { return "failed" } }(),
					"issues": touchTargetsResult.Issues,
				},
				"gestalt": map[string]interface{}{
					"status": func() string { if gestaltResult.Passed { return "passed" } else { return "failed" } }(),
					"issues": gestaltResult.Issues,
				},
				"accessibility": map[string]interface{}{
					"status": func() string { if a11yResult.Passed { return "passed" } else { return "failed" } }(),
					"issues": a11yResult.Issues,
				},
				"choice_overload": map[string]interface{}{
					"status": func() string { if choiceResult.Passed { return "passed" } else { return "failed" } }(),
					"issues": choiceResult.Issues,
				},
				"contrast": map[string]interface{}{
					"status": func() string { if contrastResult.Passed { return "passed" } else { return "failed" } }(),
					"issues": contrastResult.Issues,
				},
				"spacing": map[string]interface{}{
					"status": func() string { if spacingResult.Passed { return "passed" } else { return "failed" } }(),
					"issues": spacingResult.Issues,
				},
				"typography": map[string]interface{}{
					"status": func() string { if typographyResult.Passed { return "passed" } else { return "failed" } }(),
					"issues": typographyResult.Issues,
				},
				"elevation": map[string]interface{}{
					"status": func() string { if elevationResult.Passed { return "passed" } else { return "failed" } }(),
					"issues": elevationResult.Issues,
				},
				"loading_states": map[string]interface{}{
					"status": func() string { if loadingStatesResult.Passed { return "passed" } else { return "failed" } }(),
					"issues": loadingStatesResult.Issues,
				},
				"responsive": map[string]interface{}{
					"status": func() string { if responsiveResult.Passed { return "passed" } else { return "failed" } }(),
					"issues": responsiveResult.Issues,
				},
				"focus": map[string]interface{}{
					"status": func() string { if focusResult.Passed { return "passed" } else { return "failed" } }(),
					"issues": focusResult.Issues,
				},
				"dark_mode": map[string]interface{}{
					"status": func() string { if darkModeResult.Passed { return "passed" } else { return "failed" } }(),
					"issues": darkModeResult.Issues,
				},
			},
		}
		
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	// Console output
	fmt.Printf("🔍 Design Audit for %s\n", structureFile)
	fmt.Printf("   Version: %s\n", structure.Version)
	fmt.Printf("   Phase: %s\n", structure.Phase)
	fmt.Printf("   Components: %d\n", len(structure.Components))
	
	if structure.Locked {
		fmt.Printf("   Status: Locked")
		if structure.ApprovedBy != "" {
			fmt.Printf(" (approved)")
		}
		fmt.Println()
	} else {
		fmt.Println("   Status: Draft")
	}
	
	fmt.Println("\n═══════════════════════════════════════════════════════")
	
	// Print summary
	printAuditCategory("Visual Hierarchy", hierarchyResult.Passed, len(hierarchyResult.Issues))
	printAuditCategory("Touch Targets (Fitts's Law)", touchTargetsResult.Passed, len(touchTargetsResult.Issues))
	printAuditCategory("Gestalt Principles", gestaltResult.Passed, len(gestaltResult.Issues))
	printAuditCategory("Accessibility (WCAG)", a11yResult.Passed, len(a11yResult.Issues))
	printAuditCategory("Choice Overload (Hick's Law)", choiceResult.Passed, len(choiceResult.Issues))
	printAuditCategory("Color Contrast", contrastResult.Passed, len(contrastResult.Issues))
	printAuditCategory("Spacing Scale (8pt Grid)", spacingResult.Passed, len(spacingResult.Issues))
	printAuditCategory("Typography Scale", typographyResult.Passed, len(typographyResult.Issues))
	printAuditCategory("Shadow & Elevation", elevationResult.Passed, len(elevationResult.Issues))
	printAuditCategory("Loading States", loadingStatesResult.Passed, len(loadingStatesResult.Issues))
	printAuditCategory("Responsive Breakpoints", responsiveResult.Passed, len(responsiveResult.Issues))
	printAuditCategory("Focus Indicators", focusResult.Passed, len(focusResult.Issues))
	printAuditCategory("Dark Mode Support", darkModeResult.Passed, len(darkModeResult.Issues))
	
	fmt.Println("═══════════════════════════════════════════════════════")
	
	if allPassed {
		fmt.Println("\n✅ Overall: PASSED - All design principles validated")
	} else {
		fmt.Println("\n⚠️  Overall: ISSUES FOUND - Review recommendations above")
		fmt.Println("\nRun individual validations for detailed issue breakdown:")
		fmt.Println("  prism validate --hierarchy")
		fmt.Println("  prism validate --touch-targets")
		fmt.Println("  prism validate --gestalt")
		fmt.Println("  prism validate --accessibility")
		fmt.Println("  prism validate --choice-overload")
		fmt.Println("  prism validate --contrast")
		fmt.Println("  prism validate --spacing")
		fmt.Println("  prism validate --typography")
		fmt.Println("  prism validate --elevation")
		fmt.Println("  prism validate --loading-states")
		fmt.Println("  prism validate --responsive")
		fmt.Println("  prism validate --focus")
		fmt.Println("  prism validate --dark-mode")
	}
	
	return nil
}

func printAuditCategory(name string, passed bool, issueCount int) {
	status := "✅"
	statusText := "PASSED"
	if !passed {
		status = "⚠️ "
		statusText = fmt.Sprintf("%d ISSUES", issueCount)
	}
	fmt.Printf("%s %-35s %s\n", status, name, statusText)
}
