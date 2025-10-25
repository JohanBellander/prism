package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/johanbellander/prism/internal/types"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [version]",
	Short: "Show version details",
	Long: `Show detailed information about a specific version.

Examples:
  prism show v1
  prism show v2 --json`,
	Args: cobra.ExactArgs(1),
	RunE: runShow,
}

func runShow(cmd *cobra.Command, args []string) error {
	// Get flags
	version := args[0]
	projectPath, _ := cmd.Parent().PersistentFlags().GetString("project")
	outputJSON, _ := cmd.Parent().PersistentFlags().GetBool("json")

	// Find the structure file
	structurePath := filepath.Join(projectPath, "phase1-structure")
	
	// Determine the file name
	var fileName string
	if version == "approved" || version == "latest" {
		fileName = version + ".json"
	} else {
		fileName = version + ".json"
	}

	filePath := filepath.Join(structurePath, fileName)

	// If "latest", find the highest version number
	if version == "latest" {
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

		latestVersion := 0
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
				var v int
				if _, err := fmt.Sscanf(entry.Name(), "v%d.json", &v); err == nil {
					if v > latestVersion {
						latestVersion = v
						filePath = filepath.Join(structurePath, entry.Name())
						fileName = entry.Name()
					}
				}
			}
		}

		if latestVersion == 0 {
			if outputJSON {
				result := map[string]interface{}{
					"status": "error",
					"error":  "No versions found",
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(result)
			}
			return fmt.Errorf("no versions found in %s", structurePath)
		}
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if outputJSON {
			result := map[string]interface{}{
				"status": "error",
				"error":  fmt.Sprintf("Version '%s' not found", version),
				"path":   filePath,
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		return fmt.Errorf("version '%s' not found at %s", version, filePath)
	}

	// Read and parse the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		if outputJSON {
			result := map[string]interface{}{
				"status": "error",
				"file":   filePath,
				"error":  fmt.Sprintf("Failed to read file: %v", err),
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		return fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	structure, err := types.ParseStructure(data)
	if err != nil {
		if outputJSON {
			result := map[string]interface{}{
				"status": "error",
				"file":   filePath,
				"error":  fmt.Sprintf("Failed to parse structure: %v", err),
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		return fmt.Errorf("failed to parse structure: %w", err)
	}

	// Output results
	if outputJSON {
		// For JSON output, include the full structure
		result := map[string]interface{}{
			"status":    "success",
			"file":      fileName,
			"path":      filePath,
			"structure": structure,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	// Human-readable output
	fmt.Printf("Version: %s\n", structure.Version)
	fmt.Printf("File: %s\n", fileName)
	fmt.Printf("Phase: %s\n", structure.Phase)
	fmt.Printf("Created: %s\n", structure.CreatedAt.Format("2006-01-02 15:04:05"))
	
	if structure.Locked {
		fmt.Printf("Status: Locked ⚡\n")
		if structure.LockedAt != nil {
			fmt.Printf("Locked At: %s\n", structure.LockedAt.Format("2006-01-02 15:04:05"))
		}
		if structure.ApprovedBy != "" {
			fmt.Printf("Approved By: %s\n", structure.ApprovedBy)
		}
	} else {
		fmt.Printf("Status: Draft\n")
	}

	if structure.ParentVersion != "" {
		fmt.Printf("Parent Version: %s\n", structure.ParentVersion)
	}

	fmt.Printf("\n--- Intent ---\n")
	fmt.Printf("Purpose: %s\n", structure.Intent.Purpose)
	fmt.Printf("Primary Action: %s\n", structure.Intent.PrimaryAction)
	fmt.Printf("User Context: %s\n", structure.Intent.UserContext)
	if len(structure.Intent.KeyInteractions) > 0 {
		fmt.Printf("Key Interactions:\n")
		for _, interaction := range structure.Intent.KeyInteractions {
			fmt.Printf("  - %s\n", interaction)
		}
	}

	fmt.Printf("\n--- Layout ---\n")
	fmt.Printf("Type: %s\n", structure.Layout.Type)
	fmt.Printf("Direction: %s\n", structure.Layout.Direction)
	fmt.Printf("Spacing: %dpx\n", structure.Layout.Spacing)
	fmt.Printf("Max Width: %dpx\n", structure.Layout.MaxWidth)
	fmt.Printf("Padding: %dpx\n", structure.Layout.Padding)

	fmt.Printf("\n--- Components ---\n")
	fmt.Printf("Total Components: %d\n", len(structure.Components))
	for i, comp := range structure.Components {
		fmt.Printf("\n%d. %s (%s)\n", i+1, comp.ID, comp.Type)
		if comp.Role != "" {
			fmt.Printf("   Role: %s\n", comp.Role)
		}
		if comp.Content != "" {
			fmt.Printf("   Content: %s\n", comp.Content)
		}
		if len(comp.Children) > 0 {
			fmt.Printf("   Children: %d\n", len(comp.Children))
		}
	}

	fmt.Printf("\n--- Responsive ---\n")
	fmt.Printf("Mobile Breakpoint: %dpx\n", structure.Responsive.Mobile.Breakpoint)
	fmt.Printf("Tablet Breakpoint: %dpx\n", structure.Responsive.Tablet.Breakpoint)

	fmt.Printf("\n--- Accessibility ---\n")
	fmt.Printf("Touch Targets Min: %dpx\n", structure.Accessibility.TouchTargetsMin)
	fmt.Printf("Focus Indicators: %s\n", structure.Accessibility.FocusIndicators)
	fmt.Printf("Labels: %s\n", structure.Accessibility.Labels)
	fmt.Printf("Semantic Structure: %v\n", structure.Accessibility.SemanticStructure)

	fmt.Printf("\n--- Validation ---\n")
	fmt.Printf("Visual Hierarchy: %s\n", structure.Validation.VisualHierarchy)
	fmt.Printf("Touch Targets: %s\n", structure.Validation.TouchTargets)
	fmt.Printf("Max Nesting Depth: %d\n", structure.Validation.MaxNestingDepth)
	fmt.Printf("Responsive Tested: %v\n", structure.Validation.ResponsiveTested)
	if structure.Validation.Notes != "" {
		fmt.Printf("Notes: %s\n", structure.Validation.Notes)
	}

	if structure.ChangeSummary != "" {
		fmt.Printf("\n--- Changes ---\n")
		fmt.Printf("Summary: %s\n", structure.ChangeSummary)
		if structure.Rationale != "" {
			fmt.Printf("Rationale: %s\n", structure.Rationale)
		}
	}

	return nil
}
