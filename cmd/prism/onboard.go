package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Set up PRISM in a new project",
	Long: `Initialize a new project with PRISM documentation and examples.

This command creates:
- Project directory structure (phase1-structure/, mockups/)
- DESIGNPROCESS.md with two-phase workflow guide
- Example Phase 1 structure file
- .gitignore for mockup outputs

Run this once when starting a new UI design project.`,
	RunE: runOnboard,
}

func init() {
	onboardCmd.Flags().BoolP("force", "f", false, "Overwrite existing files")
}

func runOnboard(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")
	projectPath, _ := cmd.Flags().GetString("project")

	fmt.Println("🎨 Setting up PRISM in your project...")
	fmt.Println()

	// Create directory structure
	dirs := []string{
		"phase1-structure",
		"phase2-design",
		"mockups",
		"history",
	}

	for _, dir := range dirs {
		path := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		fmt.Printf("✅ Created %s/\n", dir)
	}

	fmt.Println()

	// Create DESIGNPROCESS.md
	designProcessPath := filepath.Join(projectPath, "DESIGNPROCESS.md")
	if _, err := os.Stat(designProcessPath); err == nil && !force {
		fmt.Printf("⚠️  DESIGNPROCESS.md already exists. Use --force to overwrite.\n")
	} else {
		if err := createDesignProcessFile(designProcessPath); err != nil {
			return err
		}
		fmt.Println("✅ Created DESIGNPROCESS.md")
	}

	// Create example structure
	examplePath := filepath.Join(projectPath, "phase1-structure", "example.json")
	if _, err := os.Stat(examplePath); err == nil && !force {
		fmt.Printf("⚠️  example.json already exists. Use --force to overwrite.\n")
	} else {
		if err := createExampleStructure(examplePath); err != nil {
			return err
		}
		fmt.Println("✅ Created phase1-structure/example.json")
	}

	// Create .gitignore
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil && !force {
		fmt.Printf("⚠️  .gitignore already exists. Use --force to overwrite.\n")
	} else {
		if err := createGitignore(gitignorePath); err != nil {
			return err
		}
		fmt.Println("✅ Created .gitignore")
	}

	// Create README.md if it doesn't exist
	readmePath := filepath.Join(projectPath, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		if err := createProjectReadme(readmePath); err != nil {
			return err
		}
		fmt.Println("✅ Created README.md")
	}

	fmt.Println()
	fmt.Println("🎉 PRISM setup complete!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Read DESIGNPROCESS.md to understand the two-phase workflow")
	fmt.Println("  2. Review phase1-structure/example.json")
	fmt.Println("  3. Render the example: prism render . --version example")
	fmt.Println("  4. Start creating your Phase 1 structures in phase1-structure/")
	fmt.Println()
	fmt.Println("Quick commands:")
	fmt.Println("  prism render .              # Render latest version")
	fmt.Println("  prism validate .            # Validate structure")
	fmt.Println("  prism list                  # List all versions")
	fmt.Println("  prism compare --from v1 --to v2  # Compare versions")
	fmt.Println()

	return nil
}

func createDesignProcessFile(path string) error {
	content := `# Two-Phase Design Process with PRISM

## Overview

This project follows a strict two-phase design process:

**Phase 1**: Structure-only (black & white wireframes)  
**Phase 2**: Design system styling (applied after Phase 1 approval)

PRISM renders your Phase 1 JSON structures into visual PNG mockups for easy review.

## Phase 1: Structure

Focus: Layout, hierarchy, spacing, component grouping

**Constraints:**
- Colors: Only black (#000000), white (#FFFFFF), grays (#E5E5E5, #737373, #525252)
- No styling: No shadows, gradients, rounded corners
- No decorative elements
- Typography: Size variations only for hierarchy

**Workflow:**
1. Create ` + "`phase1-structure/v1.json`" + `
2. Validate: ` + "`prism validate .`" + `
3. Render: ` + "`prism render . --version v1`" + `
4. Review mockup in ` + "`mockups/v1.png`" + `
5. Iterate: Create v2.json, v3.json, etc.
6. Compare: ` + "`prism compare --from v1 --to v2`" + `
7. When approved, copy to ` + "`approved.json`" + `

## Phase 2: Design

Only after Phase 1 is approved and locked.

Focus: Apply design system tokens (colors, typography, spacing, shadows)

**Rules:**
- NO structural changes
- Use design tokens only
- Create ` + "`phase2-design/v1.json`" + `

## Quick Reference

### Rendering
` + "```bash" + `
prism render .                    # Latest version
prism render . --version v2       # Specific version
prism render . --viewport mobile  # Different viewport
prism render . --all              # All versions
` + "```" + `

### Validation
` + "```bash" + `
prism validate .                  # Validate latest
prism validate . --phase 1        # Phase 1 constraints
` + "```" + `

### Comparison
` + "```bash" + `
prism compare --from v1 --to v2
` + "```" + `

### Version Management
` + "```bash" + `
prism list                        # List all versions
prism show v1                     # Show version details
` + "```" + `

## File Structure

` + "```" + `
your-project/
├── phase1-structure/
│   ├── v1.json
│   ├── v2.json
│   └── approved.json
├── phase2-design/
│   └── v1.json
├── mockups/          # Generated by PRISM
│   ├── v1.png
│   └── v2.png
├── history/
│   └── decisions.json
└── DESIGNPROCESS.md  # This file
` + "```" + `

## Why Two Phases?

Separating structure from styling prevents:
- ❌ Structural changes sneaking in during styling
- ❌ Design decisions affecting information architecture
- ❌ Scope creep and endless revisions
- ❌ Loss of approved structure during iteration

Benefits:
- ✅ Clear approval checkpoints
- ✅ Faster iteration on structure
- ✅ Design system consistency
- ✅ Version control and rollback
- ✅ Better collaboration between designers and developers

## Learn More

- PRISM Documentation: https://github.com/JohanBellander/prism
- Full Design Process Guide: See repository DESIGNPROCESS.md
`

	return os.WriteFile(path, []byte(content), 0644)
}

func createExampleStructure(path string) error {
	content := `{
  "version": "example",
  "phase": "structure",
  "created_at": "2025-10-25T12:00:00Z",
  "locked": false,
  "intent": {
    "purpose": "Example dashboard showing Phase 1 structure",
    "primary_action": "View key metrics and recent activity",
    "user_context": "Admin user checking system health",
    "key_interactions": ["view_metrics", "filter_data", "drill_down"]
  },
  "layout": {
    "type": "stack",
    "direction": "vertical",
    "spacing": 16,
    "max_width": 1200,
    "padding": 24
  },
  "components": [
    {
      "id": "header",
      "type": "box",
      "role": "header",
      "layout": {
        "display": "flex",
        "direction": "horizontal",
        "padding": 16,
        "background": "#FFFFFF",
        "border": "1px solid #E5E5E5",
        "gap": 16
      },
      "children": [
        {
          "id": "title",
          "type": "text",
          "content": "Dashboard",
          "size": "2xl",
          "weight": "bold",
          "color": "#000000"
        }
      ]
    },
    {
      "id": "metrics",
      "type": "box",
      "role": "content",
      "layout": {
        "display": "grid",
        "padding": 16,
        "background": "#FFFFFF",
        "border": "1px solid #E5E5E5",
        "gap": 16
      },
      "children": [
        {
          "id": "metric-1",
          "type": "box",
          "layout": {
            "padding": 16,
            "background": "#E5E5E5"
          },
          "children": [
            {
              "id": "metric-1-label",
              "type": "text",
              "content": "Active Users",
              "size": "sm",
              "color": "#737373"
            },
            {
              "id": "metric-1-value",
              "type": "text",
              "content": "1,234",
              "size": "3xl",
              "weight": "bold",
              "color": "#000000"
            }
          ]
        }
      ]
    }
  ],
  "validation": {
    "visual_hierarchy": "passed",
    "touch_targets": "passed",
    "max_nesting_depth": 3,
    "responsive_tested": true
  }
}
`

	return os.WriteFile(path, []byte(content), 0644)
}

func createGitignore(path string) error {
	content := `# PRISM generated mockups
mockups/*.png
!mockups/.gitkeep

# Temporary files
*.tmp
.DS_Store
Thumbs.db
`

	return os.WriteFile(path, []byte(content), 0644)
}

func createProjectReadme(path string) error {
	content := `# UI Design Project

This project uses the two-phase design process with PRISM.

## Quick Start

1. Create Phase 1 structure in ` + "`phase1-structure/v1.json`" + `
2. Render mockup: ` + "`prism render .`" + `
3. Review mockup in ` + "`mockups/v1.png`" + `
4. Iterate and improve
5. Get approval on structure
6. Move to Phase 2 styling

## Commands

` + "```bash" + `
# Render
prism render .

# Validate
prism validate .

# Compare versions
prism compare --from v1 --to v2

# List versions
prism list
` + "```" + `

## Documentation

See ` + "`DESIGNPROCESS.md`" + ` for the full two-phase workflow guide.
`

	return os.WriteFile(path, []byte(content), 0644)
}
