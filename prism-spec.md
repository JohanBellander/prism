# PRISM CLI Tool Specification

## Overview

The PRISM is a command-line tool that generates visual PNG representations of Phase 1 structural mockups from your AI Design Agent process. It takes the JSON structure files created in Phase 1 and renders them as black-and-white wireframe images for easy review and approval.

## Purpose

- **Enable visual review**: Convert JSON structure files into visual mockups
- **Facilitate approval process**: Provide clear visual representation of Phase 1 structure before moving to Phase 2
- **Support iteration**: Allow quick visual comparison between structure versions
- **Streamline workflow**: Reduce friction in the design process by making structures immediately visible

## Command-Line Interface

Inspired by Steve Yegge's Beads project architecture, the CLI follows these design principles:
- **JSON-first output**: All commands support `--json` for programmatic use
- **Composable commands**: Single-purpose commands that work well together
- **Helpful defaults**: Sensible defaults with easy overrides
- **Progressive disclosure**: Simple usage with advanced options available

### Basic Usage

```bash
prism [command] [options]
```

### Core Commands

| Command | Purpose | Key Flags |
|---------|---------|-----------|
| `prism render` | Render mockup to PNG | `--version`, `--output`, `--viewport`, `--json` |
| `prism list` | List available versions | `--project`, `--json` |
| `prism show` | Show version details | `--version`, `--json` |
| `prism validate` | Validate JSON structure | `--phase`, `--json` |
| `prism compare` | Compare versions side-by-side | `--from`, `--to`, `--json` |

### Global Flags

| Flag | Short | Description | Environment | Default |
|------|-------|-------------|-------------|---------|
| `--json` | | JSON output format | `MOCKUP_JSON` | `false` |
| `--project` | `-p` | Project directory path | `MOCKUP_PROJECT` | `./` |
| `--quiet` | `-q` | Suppress non-essential output | `MOCKUP_QUIET` | `false` |
| `--config` | | Config file path | `MOCKUP_CONFIG` | `~/.prism` |

### Render Command

The primary command for generating visual mockups:

```bash
prism render [project-path] [options]
```

#### Render Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--version` | `-v` | Version to render (v1, v2, approved, latest) | `latest` |
| `--output` | `-o` | Output file path | `{project}-phase1-{version}.png` |
| `--width` | `-w` | Canvas width in pixels | `1200` |
| `--height` | `-h` | Canvas height in pixels | `auto` |
| `--scale` | `-s` | Scale factor for high-DPI displays | `1` |
| `--viewport` | | Target viewport (mobile, tablet, desktop) | `desktop` |
| `--annotations` | `-a` | Include annotations (IDs, dimensions) | `false` |
| `--grid` | `-g` | Show layout grid overlay | `false` |
| `--format` | `-f` | Output format (png, svg, pdf) | `png` |
| `--theme` | | Color theme (bw, wireframe, blueprint) | `bw` |

### Examples

```bash
# Quick start - render latest version
prism render ./my-dashboard

# Render specific version with JSON output
prism render ./my-dashboard --version v2 --json

# Compare two versions side-by-side
prism compare ./my-dashboard --from v1 --to v2

# List all available versions
prism list --project ./my-dashboard --json

# Validate structure before rendering
prism validate ./my-dashboard --phase 1

# Mobile viewport with annotations
prism render ./my-dashboard --viewport mobile --annotations

# High-DPI rendering for presentations
prism render ./my-dashboard --scale 2 --format pdf
```

### Batch Operations

Following Beads' pattern of bulk operations:

```bash
# Render all versions in project
prism render ./my-dashboard --all-versions

# Render all viewports for a version
prism render ./my-dashboard --all-viewports

# Batch validation
prism validate ./projects/* --recursive
```

## Input Format

### Expected Project Structure

```
project-name/
├── phase1-structure/
│   ├── v1.json
│   ├── v2.json
│   └── approved.json (optional)
└── README.md (optional)
```

### JSON Structure Requirements

The tool expects JSON files following the format specified in DESIGNPROCESS.md:

- `version`: Version identifier
- `phase`: Must be "structure"
- `layout`: Root layout configuration
- `components`: Array of component definitions with layout and styling
- `responsive`: Responsive breakpoint definitions
- `accessibility`: Accessibility metadata

## Rendering Specifications

### Color Palette

Following Phase 1 constraints, only these colors are used:

- **Background**: `#FFFFFF` (white)
- **Primary elements**: `#000000` (black)
- **Secondary elements**: `#525252` (dark gray)
- **Tertiary elements**: `#737373` (medium gray)
- **Borders/dividers**: `#E5E5E5` (light gray)

### Typography

- **Font family**: System monospace font (consistent, readable)
- **Size mapping**: 
  - `xs`: 10px
  - `sm`: 12px
  - `base`: 14px
  - `lg`: 16px
  - `xl`: 18px
  - `2xl`: 20px
  - `3xl`: 24px
  - `4xl`: 28px
- **Weight mapping**:
  - `normal`: Regular weight
  - `bold`: Bold weight

### Visual Elements

#### Components
- **Boxes**: Rendered as rectangles with borders
- **Text**: Rendered with appropriate sizing and weight
- **Inputs**: Rendered as outlined rectangles with placeholder text
- **Buttons**: Rendered as filled rectangles with centered text
- **Images**: Rendered as gray rectangles with "IMAGE" placeholder

#### Layout
- **Flexbox**: Proper spacing and alignment
- **Grid**: CSS Grid layout simulation
- **Padding/Margins**: Accurate spacing representation
- **Borders**: 1px solid borders where specified

#### Annotations (when enabled)
- **Component IDs**: Small labels above each component
- **Dimensions**: Width x Height labels
- **Touch targets**: 44x44px minimum size indicators
- **Spacing**: Gap measurements between elements

### Responsive Rendering

When `--viewport` is specified:

- **Mobile**: 375px width (common mobile size)
- **Tablet**: 768px width
- **Desktop**: 1200px width (default)

The tool applies responsive overrides from the JSON structure.

## Output Specifications

### JSON-First Design

Following Beads' pattern, all commands support structured JSON output:

```bash
# Human-readable output
prism render ./my-dashboard

# JSON output for scripting
prism render ./my-dashboard --json
```

#### JSON Response Format

```json
{
  "status": "success",
  "command": "render",
  "project": {
    "name": "my-dashboard",
    "path": "./my-dashboard",
    "phase": 1
  },
  "input": {
    "version": "v2",
    "file": "phase1-structure/v2.json",
    "checksum": "sha256:abc123..."
  },
  "output": {
    "file": "my-dashboard-phase1-v2.png",
    "format": "png",
    "dimensions": {
      "width": 1200,
      "height": 800
    },
    "viewport": "desktop",
    "size_bytes": 245760
  },
  "validation": {
    "structure_valid": true,
    "phase1_constraints": true,
    "accessibility_passed": true,
    "touch_targets_ok": true
  },
  "render_time_ms": 1247,
  "metadata": {
    "tool_version": "1.0.0",
    "rendered_at": "2025-10-25T14:30:00Z",
    "git_commit": "abc123def"
  }
}
```

### PNG Properties
- **Format**: PNG with transparency support
- **Color depth**: 8-bit grayscale or 24-bit RGB  
- **Compression**: Optimized for file size
- **Metadata**: Include rendering timestamp and source file

### File Naming Convention
- Default: `{project-name}-phase1-{version}.png`
- With viewport: `{project-name}-phase1-{version}-{viewport}.png` 
- Custom: User-specified via `--output`

### Exit Codes

Following Unix conventions and Beads' error handling patterns:

- `0`: Success
- `1`: Invalid arguments or general error
- `2`: Project structure not found
- `3`: JSON validation failed
- `4`: Rendering failed
- `5`: File I/O error
- `6`: Phase constraint violation

## Error Handling

### Beads-Inspired Error Design

Following the robust error handling patterns from Steve Yegge's Beads:

#### Validation-First Approach
```bash
# Always validate before rendering
prism validate ./my-dashboard --phase 1
prism render ./my-dashboard  # Only proceeds if valid
```

#### Structured Error Messages

**CLI Output (Human-readable)**:
```
❌ Error: Phase 1 constraint violation

Project: my-dashboard
File: phase1-structure/v2.json
Issue: Color value '#3498db' found in component 'header-button'

Phase 1 allows only:
  - Black: #000000
  - White: #FFFFFF  
  - Grays: #E5E5E5, #737373, #525252

Suggestion: Remove color or move to Phase 2
```

**JSON Output (Machine-readable)**:
```json
{
  "status": "error",
  "error_code": 6,
  "error_type": "phase_constraint_violation", 
  "message": "Color value '#3498db' violates Phase 1 constraints",
  "details": {
    "project": "my-dashboard",
    "file": "phase1-structure/v2.json",
    "component_id": "header-button",
    "constraint": "phase1_colors_only",
    "allowed_colors": ["#000000", "#FFFFFF", "#E5E5E5", "#737373", "#525252"],
    "found_color": "#3498db"
  },
  "suggestion": "Remove color or move to Phase 2",
  "docs_url": "https://docs.prism.com/phase1-constraints"
}
```

### Progressive Error Recovery

Like Beads' graceful degradation:

1. **Validation errors**: Stop with helpful guidance
2. **Missing components**: Render with placeholders
3. **Partial failures**: Render what's possible, warn about issues
4. **Network timeouts**: Retry with exponential backoff

### Error Categories

| Category | Exit Code | Example | Recovery |
|----------|-----------|---------|----------|
| **User Error** | 1 | Invalid flags | Show help |
| **Project Error** | 2 | Missing structure | Guide setup |
| **Validation Error** | 3 | Invalid JSON | Show validation details |
| **Rendering Error** | 4 | Canvas failure | Fallback rendering |
| **I/O Error** | 5 | Permission denied | Suggest fixes |
| **Constraint Error** | 6 | Phase violation | Explain constraints |

## Dependencies and Requirements

### Architecture Inspired by Beads

Following Beads' clean dependency model:

#### Core Dependencies (Minimal)
- **Go 1.19+**: Following Beads' Go choice for CLI performance
- **Cobra**: CLI framework (same as Beads) for consistent UX
- **Color**: Terminal colors (github.com/fatih/color, same as Beads)
- **Image**: Go's built-in image libraries for PNG generation

#### Optional Dependencies
- **Canvas**: For advanced rendering features
- **Git**: For version metadata integration
- **Config**: For user preferences (Viper, like Beads)

### Installation Methods

```bash
# Go install (recommended)
go install github.com/yourusername/prism@latest

# Binary release
curl -L https://github.com/yourusername/prism/releases/latest/download/prism-linux-amd64 -o prism

# Build from source
git clone https://github.com/yourusername/prism
cd prism
make build
```

### System Requirements
- **Go**: 1.19+ for building
- **Memory**: 256MB minimum, 512MB recommended
- **Storage**: Space for output PNG files (typically 50-500KB each)
- **OS**: Windows, macOS, Linux (cross-platform like Beads)

## Performance Considerations

### Optimization Targets
- **Rendering time**: < 5 seconds for typical mockups
- **Memory usage**: < 256MB peak during rendering
- **Output size**: Optimized PNG compression

### Scaling Considerations
- **Large layouts**: Handle layouts up to 5000px height
- **Complex nesting**: Support up to 6 levels of component nesting
- **Multiple versions**: Batch processing capability

## Future Enhancements

### Phase 1: Core Tool (MVP)
Following Beads' incremental development approach:

- ✅ Basic `render` command with PNG output
- ✅ Phase 1 constraint validation
- ✅ JSON output for all commands  
- ✅ Multiple viewport rendering
- ✅ Annotation support

### Phase 2: Workflow Integration
- **Batch operations**: `--all-versions`, `--all-viewports`
- **Comparison mode**: Side-by-side version diffs
- **Watch mode**: Auto-render on file changes
- **Config system**: User preferences and themes

### Phase 3: Advanced Features
- **Multiple formats**: SVG, PDF export
- **Interactive mode**: CLI wizard like Beads
- **Plugin system**: Custom component renderers
- **Integration**: Git hooks, VS Code extension

### Phase 4: Ecosystem
- **Web interface**: Browser-based viewer
- **API server**: HTTP API for integrations
- **Cloud rendering**: Remote rendering service
- **CI/CD plugins**: GitHub Actions, Jenkins

### Beads-Style Quality Gates

Each phase requires:
- ✅ Comprehensive test coverage
- ✅ Performance benchmarks
- ✅ Cross-platform validation
- ✅ Documentation completeness
- ✅ User feedback integration

### Integration Patterns

Following Beads' integration philosophy:

```bash
# Git integration
git add phase1-structure/v2.json
git commit -m "Update dashboard structure"
prism render . --git-commit-info

# CI/CD pipeline
prism validate ./projects --recursive --json > validation.json
prism render ./projects --all-versions --output-dir ./artifacts

# Development workflow  
prism watch ./my-project --auto-render --open-browser
```

## Testing Strategy

### Unit Tests
- JSON parsing and validation
- Component rendering functions
- Layout calculation algorithms
- Color and typography mapping

### Integration Tests
- End-to-end rendering pipeline
- Multiple viewport rendering
- Error handling scenarios
- File I/O operations

### Visual Regression Tests
- Compare rendered output against reference images
- Test responsive breakpoint rendering
- Validate annotation positioning
- Check color constraint compliance

## Documentation Requirements

### User Documentation
- **Installation guide**: Setup instructions for different platforms
- **Usage examples**: Common use cases and workflows
- **Troubleshooting**: Common issues and solutions
- **API reference**: Complete option documentation

### Developer Documentation
- **Architecture overview**: Code organization and design patterns
- **Rendering pipeline**: Step-by-step rendering process
- **Extension guide**: How to add new component types
- **Contributing guide**: Development setup and contribution process

## Success Criteria

Drawing from Beads' principles of quality and user experience:

### Core Success Metrics

1. **Accuracy**: Visual output faithfully represents JSON structure intent
   - Pixel-perfect layout matching JSON specifications
   - Correct spacing, typography, and element positioning
   - Proper responsive breakpoint handling

2. **Performance**: Fast execution like Beads CLI commands
   - < 3 seconds for typical mockups (Beads-level responsiveness)
   - < 256MB memory usage during rendering
   - Efficient caching and incremental rendering

3. **Reliability**: Handles edge cases gracefully
   - Validates input before processing
   - Provides clear error messages with actionable suggestions
   - Graceful degradation for unsupported features

4. **Usability**: Intuitive CLI following Unix conventions
   - Helpful defaults requiring minimal flags
   - Consistent with Beads' command patterns
   - Progressive disclosure of advanced features

5. **Quality**: Professional output suitable for approval workflows
   - Clean, readable wireframes
   - Proper annotation placement when enabled
   - Consistent visual style across renders

### User Experience Goals

**Developer Workflow Integration**:
```bash
# Should feel natural in development flow
cd my-project
prism render .          # Quick preview
prism compare v1 v2     # Review changes  
git commit -m "Updated layout"   # Commit with confidence
```

**Stakeholder Review Process**:
```bash
# Generate review materials easily
prism render . --all-viewports --annotations
# Output: my-project-desktop.png, my-project-mobile.png, etc.
```

**Quality Validation**:
```bash
# Built-in quality checks
prism validate . --phase 1 --strict
# Ensures Phase 1 compliance before render
```

### Technical Success Indicators

- **Zero magic**: All behavior is explicit and documented
- **Composable**: Commands work well together
- **Debuggable**: Clear error messages and JSON output
- **Maintainable**: Clean code following Go/Beads patterns
- **Extensible**: Plugin architecture for custom needs

The tool succeeds when teams can confidently approve Phase 1 structures based on the visual output, leading to better Phase 2 designs and fewer iteration cycles.

## Risk Assessment

### Technical Risks
- **Canvas/rendering library limitations**: Mitigation through library evaluation
- **Memory usage with large layouts**: Implement streaming rendering
- **Cross-platform compatibility**: Test on Windows, macOS, Linux

### User Experience Risks
- **Learning curve**: Provide clear documentation and examples
- **Integration friction**: Design for easy integration with existing workflows
- **Output quality**: Validate visual fidelity with real users

---

This specification provides a comprehensive foundation for implementing the PRISM CLI tool. The tool should integrate seamlessly with your AI Design Agent workflow, making it easy to visualize and approve Phase 1 structures before proceeding to Phase 2 design application.
