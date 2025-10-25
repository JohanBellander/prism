package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/johanbellander/prism/internal/render"
	"github.com/johanbellander/prism/internal/types"
	"github.com/spf13/cobra"
)

var renderCmd = &cobra.Command{
	Use:   "render [project-path]",
	Short: "Render mockup to PNG",
	Long: `Render a Phase 1 structure JSON file to a visual PNG mockup.

Examples:
  prism render ./my-dashboard
  prism render ./my-dashboard --version v2 --json
  prism render ./my-dashboard --viewport mobile --annotations`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRender,
}

func init() {
	// Render-specific flags
	renderCmd.Flags().StringP("version", "v", "latest", "Version to render (v1, v2, approved, latest)")
	renderCmd.Flags().StringP("output", "o", "", "Output file path (default: {project}-phase1-{version}.png)")
	renderCmd.Flags().IntP("width", "w", 1200, "Canvas width in pixels")
	renderCmd.Flags().Int("height", 0, "Canvas height in pixels (0 for auto)")
	renderCmd.Flags().IntP("scale", "s", 1, "Scale factor for high-DPI displays")
	renderCmd.Flags().String("viewport", "desktop", "Target viewport (mobile, tablet, desktop)")
	renderCmd.Flags().BoolP("annotations", "a", false, "Include annotations (IDs, dimensions)")
	renderCmd.Flags().BoolP("grid", "g", false, "Show layout grid overlay")
	renderCmd.Flags().StringP("format", "f", "png", "Output format (png, svg, pdf)")
	renderCmd.Flags().String("theme", "bw", "Color theme (bw, wireframe, blueprint)")
	renderCmd.Flags().Bool("all", false, "Render all versions found in phase1-structure directory")
}

func runRender(cmd *cobra.Command, args []string) error {
	// Get flags
	projectPath := "./"
	if len(args) > 0 {
		projectPath = args[0]
	}

	versionFlag, _ := cmd.Flags().GetString("version")
	outputPath, _ := cmd.Flags().GetString("output")
	width, _ := cmd.Flags().GetInt("width")
	height, _ := cmd.Flags().GetInt("height")
	scale, _ := cmd.Flags().GetInt("scale")
	viewport, _ := cmd.Flags().GetString("viewport")
	annotations, _ := cmd.Flags().GetBool("annotations")
	grid, _ := cmd.Flags().GetBool("grid")
	renderAll, _ := cmd.Flags().GetBool("all")
	outputJSON, _ := cmd.Parent().PersistentFlags().GetBool("json")

	// If --all flag is set, render all versions
	if renderAll {
		return renderAllVersions(cmd, projectPath, width, height, scale, viewport, annotations, grid, outputJSON)
	}

	// Find the structure file
	structurePath := filepath.Join(projectPath, "phase1-structure")
	
	var structureFile string
	if versionFlag == "approved" {
		structureFile = filepath.Join(structurePath, "approved.json")
	} else if versionFlag == "latest" {
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
	} else {
		// Specific version
		structureFile = filepath.Join(structurePath, versionFlag+".json")
	}

	if structureFile == "" {
		if outputJSON {
			result := map[string]interface{}{
				"status": "error",
				"error":  "No structure file found",
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		return fmt.Errorf("no structure file found in %s", structurePath)
	}

	// Read and parse the structure
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

	structure, err := types.ParseAndValidateStructure(data)
	if err != nil {
		if outputJSON {
			result := map[string]interface{}{
				"status": "error",
				"file":   structureFile,
				"error":  fmt.Sprintf("Failed to parse structure: %v", err),
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		return fmt.Errorf("failed to parse structure: %w", err)
	}

	// Adjust width based on viewport
	if viewport == "mobile" {
		width = 375
	} else if viewport == "tablet" {
		width = 768
	} else if viewport == "desktop" && width == 1200 {
		// Keep default
	}

	// Create renderer
	opts := render.RenderOptions{
		Width:       width,
		Height:      height,
		Scale:       scale,
		Viewport:    viewport,
		Annotations: annotations,
		Grid:        grid,
	}
	renderer := render.NewRenderer(opts)

	// Render the structure
	result, err := renderer.Render(structure)
	if err != nil {
		if outputJSON {
			errResult := map[string]interface{}{
				"status": "error",
				"error":  fmt.Sprintf("Rendering failed: %v", err),
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(errResult)
		}
		return fmt.Errorf("rendering failed: %w", err)
	}

	// Determine output path
	if outputPath == "" {
		baseName := filepath.Base(projectPath)
		if baseName == "." || baseName == "/" {
			baseName = "mockup"
		}
		outputPath = fmt.Sprintf("%s-phase1-%s.png", baseName, structure.Version)
	}

	// Save the result
	if err := result.SavePNG(outputPath); err != nil {
		if outputJSON {
			errResult := map[string]interface{}{
				"status": "error",
				"error":  fmt.Sprintf("Failed to save PNG: %v", err),
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(errResult)
		}
		return fmt.Errorf("failed to save PNG: %w", err)
	}

	// Success
	if outputJSON {
		successResult := map[string]interface{}{
			"status":  "success",
			"file":    structureFile,
			"output":  outputPath,
			"version": structure.Version,
			"width":   result.Width,
			"height":  result.Height,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(successResult)
	}

	fmt.Printf("✅ Rendered %s\n", structureFile)
	fmt.Printf("   Output: %s\n", outputPath)
	fmt.Printf("   Dimensions: %dx%d\n", result.Width, result.Height)
	fmt.Printf("   Viewport: %s\n", viewport)

	return nil
}

// renderAllVersions renders all JSON files found in the phase1-structure directory
func renderAllVersions(cmd *cobra.Command, projectPath string, width, height, scale int, viewport string, annotations, grid, outputJSON bool) error {
	structurePath := filepath.Join(projectPath, "phase1-structure")
	
	// Read all files in the directory
	entries, err := os.ReadDir(structurePath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", structurePath, err)
	}

	// Collect all JSON files
	var jsonFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			jsonFiles = append(jsonFiles, entry.Name())
		}
	}

	if len(jsonFiles) == 0 {
		if outputJSON {
			result := map[string]interface{}{
				"status": "error",
				"error":  "No JSON files found in phase1-structure directory",
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		return fmt.Errorf("no JSON files found in %s", structurePath)
	}

	projectName := filepath.Base(projectPath)
	results := []map[string]interface{}{}
	successCount := 0
	failCount := 0

	// Render each file
	for _, jsonFile := range jsonFiles {
		structureFile := filepath.Join(structurePath, jsonFile)
		versionName := jsonFile[:len(jsonFile)-5] // Remove .json extension

		// Read and parse the structure
		data, err := os.ReadFile(structureFile)
		if err != nil {
			if outputJSON {
				results = append(results, map[string]interface{}{
					"version": versionName,
					"status":  "error",
					"error":   fmt.Sprintf("Failed to read file: %v", err),
				})
			} else {
				fmt.Printf("❌ Failed to render %s: %v\n", versionName, err)
			}
			failCount++
			continue
		}

		structure, err := types.ParseAndValidateStructure(data)
		if err != nil {
			if outputJSON {
				results = append(results, map[string]interface{}{
					"version": versionName,
					"status":  "error",
					"error":   fmt.Sprintf("Failed to parse structure: %v", err),
				})
			} else {
				fmt.Printf("❌ Failed to render %s: %v\n", versionName, err)
			}
			failCount++
			continue
		}

		// Adjust width based on viewport
		renderWidth := width
		if viewport == "mobile" {
			renderWidth = 375
		} else if viewport == "tablet" {
			renderWidth = 768
		}

		// Create renderer
		opts := render.RenderOptions{
			Width:       renderWidth,
			Height:      height,
			Scale:       scale,
			Viewport:    viewport,
			Annotations: annotations,
			Grid:        grid,
		}
		renderer := render.NewRenderer(opts)

		// Render to PNG
		result, err := renderer.Render(structure)
		if err != nil {
			if outputJSON {
				results = append(results, map[string]interface{}{
					"version": versionName,
					"status":  "error",
					"error":   fmt.Sprintf("Render failed: %v", err),
				})
			} else {
				fmt.Printf("❌ Failed to render %s: %v\n", versionName, err)
			}
			failCount++
			continue
		}

		// Save the file
		outputPath := fmt.Sprintf("%s-phase1-%s.png", projectName, versionName)
		if err := result.SavePNG(outputPath); err != nil {
			if outputJSON {
				results = append(results, map[string]interface{}{
					"version": versionName,
					"status":  "error",
					"error":   fmt.Sprintf("Failed to save file: %v", err),
				})
			} else {
				fmt.Printf("❌ Failed to save %s: %v\n", versionName, err)
			}
			failCount++
			continue
		}

		// Success
		if outputJSON {
			results = append(results, map[string]interface{}{
				"version": versionName,
				"status":  "success",
				"file":    structureFile,
				"output":  outputPath,
				"width":   result.Width,
				"height":  result.Height,
			})
		} else {
			fmt.Printf("✅ Rendered %s\n", versionName)
			fmt.Printf("   Output: %s\n", outputPath)
			fmt.Printf("   Dimensions: %dx%d\n", result.Width, result.Height)
		}
		successCount++
	}

	// Output summary
	if outputJSON {
		summary := map[string]interface{}{
			"status":        "batch_complete",
			"command":       "render",
			"project":       projectName,
			"total":         len(jsonFiles),
			"success":       successCount,
			"failed":        failCount,
			"viewport":      viewport,
			"render_width":  width,
			"render_height": height,
			"results":       results,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(summary)
	}

	fmt.Printf("\n📊 Batch rendering complete:\n")
	fmt.Printf("   Total: %d versions\n", len(jsonFiles))
	fmt.Printf("   Success: %d\n", successCount)
	fmt.Printf("   Failed: %d\n", failCount)

	if failCount > 0 && successCount == 0 {
		return fmt.Errorf("all batch renders failed")
	}
	return nil
}
