package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/johanbellander/prism/internal/types"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available versions",
	Long: `List all available versions in the project's phase1-structure directory.

Examples:
  prism list
  prism list --project ./my-dashboard --json`,
	RunE: runList,
}

// VersionInfo holds information about a structure version
type VersionInfo struct {
	Version   string    `json:"version"`
	File      string    `json:"file"`
	Phase     string    `json:"phase"`
	Locked    bool      `json:"locked"`
	CreatedAt time.Time `json:"created_at"`
	Purpose   string    `json:"purpose,omitempty"`
}

func runList(cmd *cobra.Command, args []string) error {
	// Get flags
	projectPath, _ := cmd.Parent().PersistentFlags().GetString("project")
	outputJSON, _ := cmd.Parent().PersistentFlags().GetBool("json")

	// Find the phase1-structure directory
	structurePath := filepath.Join(projectPath, "phase1-structure")
	
	// Check if directory exists
	if _, err := os.Stat(structurePath); os.IsNotExist(err) {
		if outputJSON {
			result := map[string]interface{}{
				"status":   "error",
				"error":    "No phase1-structure directory found",
				"path":     structurePath,
				"versions": []VersionInfo{},
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		return fmt.Errorf("no phase1-structure directory found in %s", projectPath)
	}

	// Read directory
	entries, err := os.ReadDir(structurePath)
	if err != nil {
		if outputJSON {
			result := map[string]interface{}{
				"status": "error",
				"error":  fmt.Sprintf("Failed to read directory: %v", err),
				"path":   structurePath,
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		return fmt.Errorf("failed to read directory %s: %w", structurePath, err)
	}

	// Collect version information
	var versions []VersionInfo
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(structurePath, entry.Name())
		
		// Read and parse the file
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip files we can't read
		}

		structure, err := types.ParseStructure(data)
		if err != nil {
			continue // Skip files we can't parse
		}

		// Extract version name from filename
		versionName := strings.TrimSuffix(entry.Name(), ".json")
		
		versions = append(versions, VersionInfo{
			Version:   versionName,
			File:      entry.Name(),
			Phase:     structure.Phase,
			Locked:    structure.Locked,
			CreatedAt: structure.CreatedAt,
			Purpose:   structure.Intent.Purpose,
		})
	}

	// Sort versions (approved first, then by version number)
	sort.Slice(versions, func(i, j int) bool {
		// approved.json always comes first
		if versions[i].Version == "approved" {
			return true
		}
		if versions[j].Version == "approved" {
			return false
		}
		
		// Extract version numbers for sorting
		var vi, vj int
		fmt.Sscanf(versions[i].Version, "v%d", &vi)
		fmt.Sscanf(versions[j].Version, "v%d", &vj)
		return vi < vj
	})

	// Output results
	if outputJSON {
		result := map[string]interface{}{
			"status":   "success",
			"project":  projectPath,
			"path":     structurePath,
			"count":    len(versions),
			"versions": versions,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	// Human-readable output
	if len(versions) == 0 {
		fmt.Printf("No versions found in %s\n", projectPath)
		return nil
	}

	fmt.Printf("Versions in %s:\n\n", projectPath)
	for _, v := range versions {
		status := "draft"
		if v.Locked {
			status = "locked"
		}
		
		fmt.Printf("  %s", v.Version)
		if v.Locked {
			fmt.Printf(" ⚡")
		}
		fmt.Printf("\n")
		fmt.Printf("    File: %s\n", v.File)
		fmt.Printf("    Status: %s\n", status)
		fmt.Printf("    Created: %s\n", v.CreatedAt.Format("2006-01-02 15:04:05"))
		if v.Purpose != "" {
			fmt.Printf("    Purpose: %s\n", v.Purpose)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("Total: %d version(s)\n", len(versions))

	return nil
}
