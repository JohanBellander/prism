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

var validateCmd = &cobra.Command{
	Use:   "validate [project-path]",
	Short: "Validate design structure against UX principles and accessibility standards",
	Long: `Validate a Phase 1 structure JSON file against design principles and WCAG standards.

Run specific validators or use 'audit' command to run all validators at once.

Validation Categories:

  Phase 1 (Structural):
    --hierarchy          Visual hierarchy (heading scale, nesting depth)
    --touch-targets      Touch target sizing (44x44px minimum, Fitts's Law)
    --gestalt            Gestalt principles (proximity, similarity, continuity)
    --accessibility      WCAG compliance (labels, heading order, focus states)
    --choice-overload    Choice overload (Hick's Law, max 7 nav items)

  Phase 2 (Visual Design):
    --contrast           Color contrast (WCAG AA: 4.5:1 text, 3:1 UI)
    --typography         Typography scale (consistent ratios, 8-10 sizes)
    --spacing            8pt grid compliance (multiples of 4 or 8)
    --elevation          Shadow/elevation system (3-4 levels)
    --loading-states     Loading indicators and skeleton screens
    --responsive         Responsive breakpoints (mobile, tablet, desktop)
    --focus              Focus indicator visibility (2px outline, 3:1 contrast)
    --dark-mode          Dark mode support (separate palette, contrast)

Severity Levels:
  🔴 CRITICAL  - Must fix (accessibility violations, WCAG failures)
  🟠 WARNING   - Should fix (UX principle violations, degraded experience)
  🟡 INFO      - Consider fixing (style guide violations, minor issues)

Scoring:
  90-100  Excellent
  70-89   Good (passing)
  50-69   Fair (needs work)
  0-49    Poor (failing)

Examples:
  # Validate all Phase 1 rules
  prism validate ./my-dashboard

  # Run specific validator
  prism validate ./my-dashboard --hierarchy
  prism validate ./my-dashboard --touch-targets
  prism validate ./my-dashboard --accessibility

  # Get JSON output for CI/CD
  prism validate ./my-dashboard --json

  # Validate Phase 2 design (contrast, typography, etc.)
  prism validate ./my-dashboard --phase 2 --contrast

  # Run multiple validators
  prism validate ./my-dashboard --hierarchy --touch-targets --gestalt

For comprehensive audits, use: prism audit ./my-dashboard
For documentation, see: VALIDATION_RULES.md`,
	Args: cobra.MaximumNArgs(1),
	RunE: runValidate,
}

func init() {
	// Validate-specific flags
	validateCmd.Flags().Int("phase", 1, "Phase to validate against (1 or 2)")
	validateCmd.Flags().Bool("hierarchy", false, "Run visual hierarchy validation")
	validateCmd.Flags().Bool("touch-targets", false, "Run touch target and spacing validation")
	validateCmd.Flags().Bool("gestalt", false, "Run Gestalt principles validation (proximity and similarity)")
	validateCmd.Flags().Bool("accessibility", false, "Run accessibility (WCAG) validation")
	validateCmd.Flags().Bool("choice-overload", false, "Run choice overload (Hick's Law) validation")
	validateCmd.Flags().Bool("contrast", false, "Run color contrast (WCAG) validation")
	validateCmd.Flags().Bool("spacing", false, "Run spacing scale (8pt grid) validation")
	validateCmd.Flags().Bool("typography", false, "Run typography scale validation")
	validateCmd.Flags().Bool("elevation", false, "Run shadow/elevation system validation")
	validateCmd.Flags().Bool("loading-states", false, "Run loading states and skeleton screen validation")
	validateCmd.Flags().Bool("responsive", false, "Run responsive breakpoint validation (mobile, tablet, desktop)")
	validateCmd.Flags().Bool("focus", false, "Run focus indicator validation for interactive elements")
	validateCmd.Flags().Bool("dark-mode", false, "Run dark mode support validation")
}

func runValidate(cmd *cobra.Command, args []string) error {
	// Get flags
	projectPath := "./"
	if len(args) > 0 {
		projectPath = args[0]
	}

	phase, _ := cmd.Flags().GetInt("phase")
	outputJSON, _ := cmd.Parent().PersistentFlags().GetBool("json")
	hierarchyCheck, _ := cmd.Flags().GetBool("hierarchy")
	touchTargetsCheck, _ := cmd.Flags().GetBool("touch-targets")
	gestaltCheck, _ := cmd.Flags().GetBool("gestalt")
	a11yCheck, _ := cmd.Flags().GetBool("accessibility")
	choiceCheck, _ := cmd.Flags().GetBool("choice-overload")
	contrastCheck, _ := cmd.Flags().GetBool("contrast")
	spacingCheck, _ := cmd.Flags().GetBool("spacing")
	typographyCheck, _ := cmd.Flags().GetBool("typography")
	elevationCheck, _ := cmd.Flags().GetBool("elevation")
	loadingStatesCheck, _ := cmd.Flags().GetBool("loading-states")
	responsiveCheck, _ := cmd.Flags().GetBool("responsive")
	focusCheck, _ := cmd.Flags().GetBool("focus")
	darkModeCheck, _ := cmd.Flags().GetBool("dark-mode")

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
	
	// Try to find the latest version or approved.json
	var structureFile string
	if _, err := os.Stat(filepath.Join(structurePath, "approved.json")); err == nil {
		structureFile = filepath.Join(structurePath, "approved.json")
	} else if _, err := os.Stat(filepath.Join(structurePath, "v1.json")); err == nil {
		// Find the highest version number
		entries, err := os.ReadDir(structurePath)
		if err != nil {
			if outputJSON {
				result := map[string]interface{}{
					"status": "error",
					"error":  fmt.Sprintf("Failed to read directory: %v", err),
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(result)
			}
			return fmt.Errorf("failed to read directory %s: %w", structurePath, err)
		}

		// Find latest version
		latestVersion := 0
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
				var version int
				if _, err := fmt.Sscanf(entry.Name(), "v%d.json", &version); err == nil {
					if version > latestVersion {
						latestVersion = version
						structureFile = filepath.Join(structurePath, entry.Name())
					}
				}
			}
		}
	}

	if structureFile == "" {
		if outputJSON {
			result := map[string]interface{}{
				"status": "error",
				"error":  "No structure file found in " + structurePath,
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		return fmt.Errorf("no structure file found in %s", structurePath)
	}

	// Read the file
	data, err := os.ReadFile(structureFile)
	if err != nil {
		if outputJSON {
			result := map[string]interface{}{
				"status": "error",
				"file":   structureFile,
				"error":  fmt.Sprintf("Failed to read file: %v", err),
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		return fmt.Errorf("failed to read %s: %w", structureFile, err)
	}

	// Parse and validate
	structure, err := types.ParseAndValidateStructure(data)
	if err != nil {
		if outputJSON {
			result := map[string]interface{}{
				"status":     "failed",
				"file":       structureFile,
				"validation": "failed",
				"error":      err.Error(),
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		fmt.Printf("❌ Validation failed for %s\n", structureFile)
		return fmt.Errorf("validation error: %w", err)
	}

	// Success
	if outputJSON {
		result := map[string]interface{}{
			"status":     "success",
			"file":       structureFile,
			"validation": "passed",
			"version":    structure.Version,
			"phase":      structure.Phase,
			"components": len(structure.Components),
		}
		
		// Run hierarchy validation if requested
		if hierarchyCheck {
			hierarchyResult := validate.ValidateHierarchy(structure, validate.DefaultHierarchyRule())
			result["hierarchy"] = map[string]interface{}{
				"status": func() string {
					if hierarchyResult.Passed {
						return "passed"
					}
					return "failed"
				}(),
				"issues": hierarchyResult.Issues,
			}
		}
		
		// Run touch target validation if requested
		if touchTargetsCheck {
			touchResult := validate.ValidateTouchTargets(structure, validate.DefaultTouchTargetRule())
			result["touch_targets"] = map[string]interface{}{
				"status": func() string {
					if touchResult.Passed {
						return "passed"
					}
					return "failed"
				}(),
				"issues": touchResult.Issues,
			}
		}
		
		// Run Gestalt principles validation if requested
		if gestaltCheck {
			gestaltResult := validate.ValidateGestalt(structure, validate.DefaultGestaltRule())
			result["gestalt"] = map[string]interface{}{
				"status": func() string {
					if gestaltResult.Passed {
						return "passed"
					}
					return "failed"
				}(),
				"issues": gestaltResult.Issues,
			}
		}
		
		// Run accessibility validation if requested
		if a11yCheck {
			a11yResult := validate.ValidateAccessibility(structure, validate.DefaultA11yRule())
			result["accessibility"] = map[string]interface{}{
				"status": func() string {
					if a11yResult.Passed {
						return "passed"
					}
					return "failed"
				}(),
				"issues": a11yResult.Issues,
			}
		}
		
		// Run choice overload validation if requested
		if choiceCheck {
			choiceResult := validate.ValidateChoiceOverload(structure, validate.DefaultChoiceRule())
			result["choice_overload"] = map[string]interface{}{
				"status": func() string {
					if choiceResult.Passed {
						return "passed"
					}
					return "failed"
				}(),
				"issues": choiceResult.Issues,
			}
		}
		
		// Run contrast validation if requested
		if contrastCheck {
			contrastResult := validate.ValidateContrast(structure, validate.DefaultContrastRule())
			result["contrast"] = map[string]interface{}{
				"status": func() string {
					if contrastResult.Passed {
						return "passed"
					}
					return "failed"
				}(),
				"issues": contrastResult.Issues,
			}
		}
		
		// Run spacing validation if requested
		if spacingCheck {
			spacingResult := validate.ValidateSpacing(structure, validate.DefaultSpacingRule())
			result["spacing"] = map[string]interface{}{
				"status": func() string {
					if spacingResult.Passed {
						return "passed"
					}
					return "failed"
				}(),
				"issues": spacingResult.Issues,
			}
		}
		
		// Run typography validation if requested
		if typographyCheck {
			typographyResult := validate.ValidateTypography(structure, validate.DefaultTypographyRule())
			result["typography"] = map[string]interface{}{
				"status": func() string {
					if typographyResult.Passed {
						return "passed"
					}
					return "failed"
				}(),
				"issues": typographyResult.Issues,
			}
		}
		
		// Run elevation validation if requested
		if elevationCheck {
			elevationResult := validate.ValidateElevation(structure, validate.DefaultElevationRule())
			result["elevation"] = map[string]interface{}{
				"status": func() string {
					if elevationResult.Passed {
						return "passed"
					}
					return "failed"
				}(),
				"issues": elevationResult.Issues,
			}
		}
		
		// Run loading states validation if requested
		if loadingStatesCheck {
			loadingStatesResult := validate.ValidateLoadingStates(structure, validate.DefaultLoadingStateRule())
			result["loading_states"] = map[string]interface{}{
				"status": func() string {
					if loadingStatesResult.Passed {
						return "passed"
					}
					return "failed"
				}(),
				"issues": loadingStatesResult.Issues,
			}
		}
		
		// Run responsive breakpoint validation if requested
		if responsiveCheck {
			responsiveResult := validate.ValidateResponsive(structure, validate.DefaultResponsiveRule())
			result["responsive"] = map[string]interface{}{
				"status": func() string {
					if responsiveResult.Passed {
						return "passed"
					}
					return "failed"
				}(),
				"issues": responsiveResult.Issues,
			}
		}
		
		// Run focus indicator validation if requested
		if focusCheck {
			focusResult := validate.ValidateFocus(structure, validate.DefaultFocusRule())
			result["focus"] = map[string]interface{}{
				"status": func() string {
					if focusResult.Passed {
						return "passed"
					}
					return "failed"
				}(),
				"issues": focusResult.Issues,
			}
		}
		
		// Run dark mode support validation if requested
		if darkModeCheck {
			darkModeResult := validate.ValidateDarkMode(structure, validate.DefaultDarkModeRule())
			result["dark_mode"] = map[string]interface{}{
				"status": func() string {
					if darkModeResult.Passed {
						return "passed"
					}
					return "failed"
				}(),
				"issues": darkModeResult.Issues,
			}
		}
		
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	fmt.Printf("✅ Validation passed for %s\n", structureFile)
	fmt.Printf("   Version: %s\n", structure.Version)
	fmt.Printf("   Phase: %s\n", structure.Phase)
	fmt.Printf("   Components: %d\n", len(structure.Components))
	if structure.Locked {
		fmt.Println("   Status: Locked (approved)")
	} else {
		fmt.Println("   Status: Draft")
	}

	// Run hierarchy validation if requested
	if hierarchyCheck {
		fmt.Println("\n📊 Visual Hierarchy Validation:")
		hierarchyResult := validate.ValidateHierarchy(structure, validate.DefaultHierarchyRule())
		
		if hierarchyResult.Passed {
			fmt.Println("   Status: ✅ Passed")
		} else {
			fmt.Println("   Status: ⚠️  Issues Found")
		}
		
		// Group issues by severity
		errors := []validate.HierarchyIssue{}
		warnings := []validate.HierarchyIssue{}
		infos := []validate.HierarchyIssue{}
		
		for _, issue := range hierarchyResult.Issues {
			switch issue.Severity {
			case "error":
				errors = append(errors, issue)
			case "warning":
				warnings = append(warnings, issue)
			case "info":
				infos = append(infos, issue)
			}
		}
		
		// Print errors
		if len(errors) > 0 {
			fmt.Println("\n   Errors:")
			for _, issue := range errors {
				fmt.Printf("     ❌ %s\n", issue.Message)
			}
		}
		
		// Print warnings
		if len(warnings) > 0 {
			fmt.Println("\n   Warnings:")
			for _, issue := range warnings {
				fmt.Printf("     ⚠️  %s\n", issue.Message)
			}
		}
		
		// Print info
		if len(infos) > 0 {
			fmt.Println("\n   Info:")
			for _, issue := range infos {
				fmt.Printf("     ℹ️  %s\n", issue.Message)
			}
		}
	}

	// Run touch target validation if requested
	if touchTargetsCheck {
		fmt.Println("\n👆 Touch Target & Spacing Validation:")
		touchResult := validate.ValidateTouchTargets(structure, validate.DefaultTouchTargetRule())
		
		if touchResult.Passed {
			fmt.Println("   Status: ✅ Passed")
		} else {
			fmt.Println("   Status: ⚠️  Issues Found")
		}
		
		// Group issues by severity
		errors := []validate.TouchTargetIssue{}
		warnings := []validate.TouchTargetIssue{}
		infos := []validate.TouchTargetIssue{}
		
		for _, issue := range touchResult.Issues {
			switch issue.Severity {
			case "error":
				errors = append(errors, issue)
			case "warning":
				warnings = append(warnings, issue)
			case "info":
				infos = append(infos, issue)
			}
		}
		
		// Print errors
		if len(errors) > 0 {
			fmt.Println("\n   Errors:")
			for _, issue := range errors {
				fmt.Printf("     ❌ %s\n", issue.Message)
			}
		}
		
		// Print warnings
		if len(warnings) > 0 {
			fmt.Println("\n   Warnings:")
			for _, issue := range warnings {
				fmt.Printf("     ⚠️  %s\n", issue.Message)
			}
		}
		
		// Print info
		if len(infos) > 0 {
			fmt.Println("\n   Info:")
			for _, issue := range infos {
				fmt.Printf("     ℹ️  %s\n", issue.Message)
			}
		}
	}

	// Run Gestalt principles validation if requested
	if gestaltCheck {
		fmt.Println("\n🎨 Gestalt Principles Validation:")
		gestaltResult := validate.ValidateGestalt(structure, validate.DefaultGestaltRule())
		
		if gestaltResult.Passed {
			fmt.Println("   Status: ✅ Passed")
		} else {
			fmt.Println("   Status: ⚠️  Issues Found")
		}
		
		// Group issues by severity
		errors := []validate.GestaltIssue{}
		warnings := []validate.GestaltIssue{}
		infos := []validate.GestaltIssue{}
		
		for _, issue := range gestaltResult.Issues {
			switch issue.Severity {
			case "error":
				errors = append(errors, issue)
			case "warning":
				warnings = append(warnings, issue)
			case "info":
				infos = append(infos, issue)
			}
		}
		
		// Print errors
		if len(errors) > 0 {
			fmt.Println("\n   Errors:")
			for _, issue := range errors {
				fmt.Printf("     ❌ %s\n", issue.Message)
			}
		}
		
		// Print warnings
		if len(warnings) > 0 {
			fmt.Println("\n   Warnings:")
			for _, issue := range warnings {
				fmt.Printf("     ⚠️  %s\n", issue.Message)
			}
		}
		
		// Print info
		if len(infos) > 0 {
			fmt.Println("\n   Info:")
			for _, issue := range infos {
				fmt.Printf("     ℹ️  %s\n", issue.Message)
			}
		}
	}

	// Run accessibility validation if requested
	if a11yCheck {
		fmt.Println("\n♿ Accessibility (WCAG) Validation:")
		a11yResult := validate.ValidateAccessibility(structure, validate.DefaultA11yRule())
		
		if a11yResult.Passed {
			fmt.Println("   Status: ✅ Passed")
		} else {
			fmt.Println("   Status: ⚠️  Issues Found")
		}
		
		// Group issues by severity
		errors := []validate.A11yIssue{}
		warnings := []validate.A11yIssue{}
		infos := []validate.A11yIssue{}
		
		for _, issue := range a11yResult.Issues {
			switch issue.Severity {
			case "error":
				errors = append(errors, issue)
			case "warning":
				warnings = append(warnings, issue)
			case "info":
				infos = append(infos, issue)
			}
		}
		
		// Print errors
		if len(errors) > 0 {
			fmt.Println("\n   Errors:")
			for _, issue := range errors {
				fmt.Printf("     ❌ %s\n", issue.Message)
			}
		}
		
		// Print warnings
		if len(warnings) > 0 {
			fmt.Println("\n   Warnings:")
			for _, issue := range warnings {
				fmt.Printf("     ⚠️  %s\n", issue.Message)
			}
		}
		
		// Print info
		if len(infos) > 0 {
			fmt.Println("\n   Info:")
			for _, issue := range infos {
				fmt.Printf("     ℹ️  %s\n", issue.Message)
			}
		}
	}

	// Run choice overload validation if requested
	if choiceCheck {
		fmt.Println("\n🎯 Choice Overload (Hick's Law) Validation:")
		choiceResult := validate.ValidateChoiceOverload(structure, validate.DefaultChoiceRule())
		
		if choiceResult.Passed {
			fmt.Println("   Status: ✅ Passed")
		} else {
			fmt.Println("   Status: ⚠️  Issues Found")
		}
		
		// Group issues by severity
		errors := []validate.ChoiceIssue{}
		warnings := []validate.ChoiceIssue{}
		infos := []validate.ChoiceIssue{}
		
		for _, issue := range choiceResult.Issues {
			switch issue.Severity {
			case "error":
				errors = append(errors, issue)
			case "warning":
				warnings = append(warnings, issue)
			case "info":
				infos = append(infos, issue)
			}
		}
		
		// Print errors
		if len(errors) > 0 {
			fmt.Println("\n   Errors:")
			for _, issue := range errors {
				fmt.Printf("     ❌ %s\n", issue.Message)
			}
		}
		
		// Print warnings
		if len(warnings) > 0 {
			fmt.Println("\n   Warnings:")
			for _, issue := range warnings {
				fmt.Printf("     ⚠️  %s\n", issue.Message)
			}
		}
		
		// Print info
		if len(infos) > 0 {
			fmt.Println("\n   Info:")
			for _, issue := range infos {
				fmt.Printf("     ℹ️  %s\n", issue.Message)
			}
		}
	}

	// Run contrast validation if requested
	if contrastCheck {
		fmt.Println("\n🎨 Color Contrast (WCAG) Validation:")
		contrastResult := validate.ValidateContrast(structure, validate.DefaultContrastRule())
		
		if contrastResult.Passed {
			fmt.Println("   Status: ✅ Passed")
		} else {
			fmt.Println("   Status: ⚠️  Issues Found")
		}
		
		// Group issues by severity
		errors := []validate.ContrastIssue{}
		warnings := []validate.ContrastIssue{}
		infos := []validate.ContrastIssue{}
		
		for _, issue := range contrastResult.Issues {
			switch issue.Severity {
			case "error":
				errors = append(errors, issue)
			case "warning":
				warnings = append(warnings, issue)
			case "info":
				infos = append(infos, issue)
			}
		}
		
		// Print errors
		if len(errors) > 0 {
			fmt.Println("\n   Errors:")
			for _, issue := range errors {
				fmt.Printf("     ❌ %s\n", issue.Message)
			}
		}
		
		// Print warnings
		if len(warnings) > 0 {
			fmt.Println("\n   Warnings:")
			for _, issue := range warnings {
				fmt.Printf("     ⚠️  %s\n", issue.Message)
			}
		}
		
		// Print info
		if len(infos) > 0 {
			fmt.Println("\n   Info:")
			for _, issue := range infos {
				fmt.Printf("     ℹ️  %s\n", issue.Message)
			}
		}
	}

	// Run spacing validation if requested
	if spacingCheck {
		fmt.Println("\n📏 Spacing Scale (8pt Grid) Validation:")
		spacingResult := validate.ValidateSpacing(structure, validate.DefaultSpacingRule())
		
		if spacingResult.Passed {
			fmt.Println("   Status: ✅ Passed")
		} else {
			fmt.Println("   Status: ⚠️  Issues Found")
		}
		
		// Group issues by severity
		errors := []validate.SpacingIssue{}
		warnings := []validate.SpacingIssue{}
		infos := []validate.SpacingIssue{}
		
		for _, issue := range spacingResult.Issues {
			switch issue.Severity {
			case "error":
				errors = append(errors, issue)
			case "warning":
				warnings = append(warnings, issue)
			case "info":
				infos = append(infos, issue)
			}
		}
		
		// Print errors
		if len(errors) > 0 {
			fmt.Println("\n   Errors:")
			for _, issue := range errors {
				fmt.Printf("     ❌ %s\n", issue.Message)
			}
		}
		
		// Print warnings
		if len(warnings) > 0 {
			fmt.Println("\n   Warnings:")
			for _, issue := range warnings {
				fmt.Printf("     ⚠️  %s\n", issue.Message)
			}
		}
		
		// Print info
		if len(infos) > 0 {
			fmt.Println("\n   Info:")
			for _, issue := range infos {
				fmt.Printf("     ℹ️  %s\n", issue.Message)
			}
		}
	}

	// Run typography validation if requested
	if typographyCheck {
		fmt.Println("\n🔤 Typography Scale Validation:")
		typographyResult := validate.ValidateTypography(structure, validate.DefaultTypographyRule())
		
		if typographyResult.Passed {
			fmt.Println("   Status: ✅ Passed")
		} else {
			fmt.Println("   Status: ⚠️  Issues Found")
		}
		
		// Group issues by severity
		errors := []validate.TypographyIssue{}
		warnings := []validate.TypographyIssue{}
		infos := []validate.TypographyIssue{}
		
		for _, issue := range typographyResult.Issues {
			switch issue.Severity {
			case "error":
				errors = append(errors, issue)
			case "warning":
				warnings = append(warnings, issue)
			case "info":
				infos = append(infos, issue)
			}
		}
		
		// Print errors
		if len(errors) > 0 {
			fmt.Println("\n   Errors:")
			for _, issue := range errors {
				fmt.Printf("     ❌ %s\n", issue.Message)
			}
		}
		
		// Print warnings
		if len(warnings) > 0 {
			fmt.Println("\n   Warnings:")
			for _, issue := range warnings {
				fmt.Printf("     ⚠️  %s\n", issue.Message)
			}
		}
		
		// Print info
		if len(infos) > 0 {
			fmt.Println("\n   Info:")
			for _, issue := range infos {
				fmt.Printf("     ℹ️  %s\n", issue.Message)
			}
		}
	}

	// Run elevation validation if requested
	if elevationCheck {
		fmt.Println("\n⬆️  Shadow & Elevation Validation:")
		elevationResult := validate.ValidateElevation(structure, validate.DefaultElevationRule())
		
		if elevationResult.Passed {
			fmt.Println("   Status: ✅ Passed")
		} else {
			fmt.Println("   Status: ⚠️  Issues Found")
		}
		
		// Group issues by severity
		errors := []validate.ElevationIssue{}
		warnings := []validate.ElevationIssue{}
		infos := []validate.ElevationIssue{}
		
		for _, issue := range elevationResult.Issues {
			switch issue.Severity {
			case "error":
				errors = append(errors, issue)
			case "warning":
				warnings = append(warnings, issue)
			case "info":
				infos = append(infos, issue)
			}
		}
		
		// Print errors
		if len(errors) > 0 {
			fmt.Println("\n   Errors:")
			for _, issue := range errors {
				fmt.Printf("     ❌ %s\n", issue.Message)
			}
		}
		
		// Print warnings
		if len(warnings) > 0 {
			fmt.Println("\n   Warnings:")
			for _, issue := range warnings {
				fmt.Printf("     ⚠️  %s\n", issue.Message)
			}
		}
		
		// Print info
		if len(infos) > 0 {
			fmt.Println("\n   Info:")
			for _, issue := range infos {
				fmt.Printf("     ℹ️  %s\n", issue.Message)
			}
		}
	}

	// Run loading states validation if requested
	if loadingStatesCheck {
		fmt.Println("\n⏳ Loading States Validation:")
		loadingStatesResult := validate.ValidateLoadingStates(structure, validate.DefaultLoadingStateRule())
		
		if loadingStatesResult.Passed {
			fmt.Println("   Status: ✅ Passed")
		} else {
			fmt.Println("   Status: ⚠️  Issues Found")
		}
		
		// Group issues by severity
		errors := []validate.LoadingStateIssue{}
		warnings := []validate.LoadingStateIssue{}
		infos := []validate.LoadingStateIssue{}
		
		for _, issue := range loadingStatesResult.Issues {
			switch issue.Severity {
			case "error":
				errors = append(errors, issue)
			case "warning":
				warnings = append(warnings, issue)
			case "info":
				infos = append(infos, issue)
			}
		}
		
		// Print errors
		if len(errors) > 0 {
			fmt.Println("\n   Errors:")
			for _, issue := range errors {
				fmt.Printf("     ❌ %s\n", issue.Message)
			}
		}
		
		// Print warnings
		if len(warnings) > 0 {
			fmt.Println("\n   Warnings:")
			for _, issue := range warnings {
				fmt.Printf("     ⚠️  %s\n", issue.Message)
			}
		}
		
		// Print info
		if len(infos) > 0 {
			fmt.Println("\n   Info:")
			for _, issue := range infos {
				fmt.Printf("     ℹ️  %s\n", issue.Message)
			}
		}
	}

	// Run responsive breakpoint validation if requested
	if responsiveCheck {
		fmt.Println("\n📱 Responsive Breakpoint Validation:")
		responsiveResult := validate.ValidateResponsive(structure, validate.DefaultResponsiveRule())
		
		if responsiveResult.Passed {
			fmt.Println("   Status: ✅ Passed")
		} else {
			fmt.Println("   Status: ⚠️  Issues Found")
		}
		
		// Group issues by severity
		errors := []validate.ResponsiveIssue{}
		warnings := []validate.ResponsiveIssue{}
		infos := []validate.ResponsiveIssue{}
		
		for _, issue := range responsiveResult.Issues {
			switch issue.Severity {
			case "error":
				errors = append(errors, issue)
			case "warning":
				warnings = append(warnings, issue)
			case "info":
				infos = append(infos, issue)
			}
		}
		
		// Print errors
		if len(errors) > 0 {
			fmt.Println("\n   Errors:")
			for _, issue := range errors {
				fmt.Printf("     ❌ [%s] %s\n", issue.Viewport, issue.Message)
			}
		}
		
		// Print warnings
		if len(warnings) > 0 {
			fmt.Println("\n   Warnings:")
			for _, issue := range warnings {
				fmt.Printf("     ⚠️  [%s] %s\n", issue.Viewport, issue.Message)
			}
		}
		
		// Print info
		if len(infos) > 0 {
			fmt.Println("\n   Info:")
			for _, issue := range infos {
				fmt.Printf("     ℹ️  [%s] %s\n", issue.Viewport, issue.Message)
			}
		}
	}

	// Run focus indicator validation if requested
	if focusCheck {
		fmt.Println("\n🎯 Focus Indicator Validation:")
		focusResult := validate.ValidateFocus(structure, validate.DefaultFocusRule())
		
		if focusResult.Passed {
			fmt.Println("   Status: ✅ Passed")
		} else {
			fmt.Println("   Status: ⚠️  Issues Found")
		}
		
		// Group issues by severity
		errors := []validate.FocusIssue{}
		warnings := []validate.FocusIssue{}
		infos := []validate.FocusIssue{}
		
		for _, issue := range focusResult.Issues {
			switch issue.Severity {
			case "error":
				errors = append(errors, issue)
			case "warning":
				warnings = append(warnings, issue)
			case "info":
				infos = append(infos, issue)
			}
		}
		
		// Print errors
		if len(errors) > 0 {
			fmt.Println("\n   Errors:")
			for _, issue := range errors {
				fmt.Printf("     ❌ %s\n", issue.Message)
			}
		}
		
		// Print warnings
		if len(warnings) > 0 {
			fmt.Println("\n   Warnings:")
			for _, issue := range warnings {
				fmt.Printf("     ⚠️  %s\n", issue.Message)
			}
		}
		
		// Print info
		if len(infos) > 0 {
			fmt.Println("\n   Info:")
			for _, issue := range infos {
				fmt.Printf("     ℹ️  %s\n", issue.Message)
			}
		}
	}

	// Run dark mode support validation if requested
	if darkModeCheck {
		fmt.Println("\n🌓 Dark Mode Support Validation:")
		darkModeResult := validate.ValidateDarkMode(structure, validate.DefaultDarkModeRule())
		
		if darkModeResult.Passed {
			fmt.Println("   Status: ✅ Passed")
		} else {
			fmt.Println("   Status: ⚠️  Issues Found")
		}
		
		// Group issues by severity
		errors := []validate.DarkModeIssue{}
		warnings := []validate.DarkModeIssue{}
		infos := []validate.DarkModeIssue{}
		
		for _, issue := range darkModeResult.Issues {
			switch issue.Severity {
			case "error":
				errors = append(errors, issue)
			case "warning":
				warnings = append(warnings, issue)
			case "info":
				infos = append(infos, issue)
			}
		}
		
		// Print errors
		if len(errors) > 0 {
			fmt.Println("\n   Errors:")
			for _, issue := range errors {
				fmt.Printf("     ❌ %s\n", issue.Message)
			}
		}
		
		// Print warnings
		if len(warnings) > 0 {
			fmt.Println("\n   Warnings:")
			for _, issue := range warnings {
				fmt.Printf("     ⚠️  %s\n", issue.Message)
			}
		}
		
		// Print info
		if len(infos) > 0 {
			fmt.Println("\n   Info:")
			for _, issue := range infos {
				fmt.Printf("     ℹ️  %s\n", issue.Message)
			}
		}
	}

	return nil
}
