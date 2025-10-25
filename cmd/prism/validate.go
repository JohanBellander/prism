package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/johanbellander/prism/internal/types"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [project-path]",
	Short: "Validate JSON structure",
	Long: `Validate a Phase 1 structure JSON file against the schema.

Examples:
  prism validate ./my-dashboard
  prism validate ./my-dashboard --phase 1 --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: runValidate,
}

func init() {
	// Validate-specific flags
	validateCmd.Flags().Int("phase", 1, "Phase to validate against (1 or 2)")
}

func runValidate(cmd *cobra.Command, args []string) error {
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

	return nil
}
