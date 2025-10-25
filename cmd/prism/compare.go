package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"github.com/johanbellander/prism/internal/render"
	"github.com/johanbellander/prism/internal/types"
	"github.com/spf13/cobra"
)

var compareCmd = &cobra.Command{
	Use:   "compare [project-path]",
	Short: "Compare two versions side-by-side",
	Long: `Compare two versions of a Phase 1 structure by rendering them side-by-side.

This command renders two versions and places them next to each other in a single PNG
for easy visual comparison of changes between versions.

Examples:
  prism compare ./my-dashboard --from v1 --to v2
  prism compare ./my-dashboard --from v1 --to v2 --json
  prism compare ./my-dashboard --from v1 --to v2 --output comparison.png`,
	RunE: runCompare,
}

var (
	compareFrom   string
	compareTo     string
	compareOutput string
)

func init() {
	compareCmd.Flags().StringVar(&compareFrom, "from", "v1", "Source version to compare from")
	compareCmd.Flags().StringVar(&compareTo, "to", "v2", "Target version to compare to")
	compareCmd.Flags().StringVarP(&compareOutput, "output", "o", "", "Output file path (default: {project}-compare-{from}-{to}.png)")
}

func runCompare(cmd *cobra.Command, args []string) error {
	// Get project path - check parent flag first, then args
	projectPath, _ := cmd.Parent().PersistentFlags().GetString("project")
	if len(args) > 0 {
		projectPath = args[0]
	}

	outputJSON, _ := cmd.Parent().PersistentFlags().GetBool("json")

	// Get absolute project path
	absProjectPath, err := filepath.Abs(projectPath)
	if err != nil {
		return fmt.Errorf("invalid project path: %w", err)
	}

	projectName := filepath.Base(absProjectPath)

	// Find structure files
	fromFile := filepath.Join(absProjectPath, "phase1-structure", compareFrom+".json")
	toFile := filepath.Join(absProjectPath, "phase1-structure", compareTo+".json")

	// Check if files exist
	if _, err := os.Stat(fromFile); os.IsNotExist(err) {
		return fmt.Errorf("source version %s not found", compareFrom)
	}
	if _, err := os.Stat(toFile); os.IsNotExist(err) {
		return fmt.Errorf("target version %s not found", compareTo)
	}

	// Load both structures
	fromData, err := os.ReadFile(fromFile)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", compareFrom, err)
	}

	fromStructure, err := types.ParseAndValidateStructure(fromData)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", compareFrom, err)
	}

	toData, err := os.ReadFile(toFile)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", compareTo, err)
	}

	toStructure, err := types.ParseAndValidateStructure(toData)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", compareTo, err)
	}

	// Render both versions
	width := 1200
	height := 800

	opts := render.RenderOptions{
		Width:    width,
		Height:   height,
		Scale:    1,
		Viewport: "desktop",
	}
	renderer := render.NewRenderer(opts)

	fromResult, err := renderer.Render(fromStructure)
	if err != nil {
		return fmt.Errorf("failed to render %s: %w", compareFrom, err)
	}

	toResult, err := renderer.Render(toStructure)
	if err != nil {
		return fmt.Errorf("failed to render %s: %w", compareTo, err)
	}

	// Get images from results
	fromImg := fromResult.Image
	toImg := toResult.Image

	// Create side-by-side comparison image
	gap := 20 // pixels between images
	compWidth := fromImg.Bounds().Dx() + gap + toImg.Bounds().Dx()
	compHeight := fromImg.Bounds().Dy()
	if toImg.Bounds().Dy() > compHeight {
		compHeight = toImg.Bounds().Dy()
	}

	compImg := image.NewRGBA(image.Rect(0, 0, compWidth, compHeight))
	
	// Fill with white background
	draw.Draw(compImg, compImg.Bounds(), image.White, image.Point{}, draw.Src)

	// Draw both images
	draw.Draw(compImg, fromImg.Bounds(), fromImg, image.Point{}, draw.Src)
	toOffset := image.Pt(fromImg.Bounds().Dx()+gap, 0)
	draw.Draw(compImg, toImg.Bounds().Add(toOffset), toImg, image.Point{}, draw.Src)

	// Determine output filename
	outputFile := compareOutput
	if outputFile == "" {
		outputFile = fmt.Sprintf("%s-compare-%s-%s.png", projectName, compareFrom, compareTo)
	}

	// Save comparison image
	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	if err := png.Encode(out, compImg); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	// Output result
	if outputJSON {
		result := map[string]interface{}{
			"status":  "success",
			"command": "compare",
			"project": map[string]interface{}{
				"name": projectName,
				"path": absProjectPath,
			},
			"from": map[string]interface{}{
				"version": compareFrom,
				"file":    fromFile,
				"width":   fromImg.Bounds().Dx(),
				"height":  fromImg.Bounds().Dy(),
			},
			"to": map[string]interface{}{
				"version": compareTo,
				"file":    toFile,
				"width":   toImg.Bounds().Dx(),
				"height":  toImg.Bounds().Dy(),
			},
			"output": map[string]interface{}{
				"file":   outputFile,
				"format": "png",
				"dimensions": map[string]interface{}{
					"width":  compWidth,
					"height": compHeight,
				},
			},
			"summary": map[string]interface{}{
				"viewport":     "desktop",
				"gap_pixels":   gap,
				"layout":       "side-by-side",
				"from_purpose": fromStructure.Intent.Purpose,
				"to_purpose":   toStructure.Intent.Purpose,
				"from_locked":  fromStructure.Locked,
				"to_locked":    toStructure.Locked,
				"same_phase":   fromStructure.Phase == toStructure.Phase,
			},
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	fmt.Printf("✅ Compared %s vs %s\n", compareFrom, compareTo)
	fmt.Printf("   From: %s (%dx%d)\n", compareFrom, fromImg.Bounds().Dx(), fromImg.Bounds().Dy())
	fmt.Printf("   To: %s (%dx%d)\n", compareTo, toImg.Bounds().Dx(), toImg.Bounds().Dy())
	fmt.Printf("   Output: %s (%dx%d)\n", outputFile, compWidth, compHeight)
	fmt.Printf("   Layout: Side-by-side with %dpx gap\n", gap)
	if toStructure.ChangeSummary != "" {
		fmt.Printf("   Changes: %s\n", toStructure.ChangeSummary)
	}

	return nil
}
