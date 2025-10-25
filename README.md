# PRISM

A CLI tool that generates visual PNG representations of Phase 1 structural mockups from AI Design Agent processes.

PRISM (Phase Render & Inspection for Structural Mockups) transforms JSON structure files into visual wireframes, enabling quick review and approval workflows.

## Features

- **Visual Mockup Generation**: Convert Phase 1 JSON structures to PNG wireframes
- **Multiple Viewports**: Support for mobile, tablet, and desktop layouts
- **Version Comparison**: Side-by-side visual comparison of different versions
- **Batch Rendering**: Render all versions at once with `--all` flag
- **Validation**: Validate JSON structures against Phase 1 constraints
- **JSON-First Output**: All commands support `--json` for programmatic use
- **Version Management**: Track and list different structure versions

## Installation

### From Source

Requires Go 1.19 or later:

```bash
# Clone the repository
git clone https://github.com/johanbellander/prism.git
cd prism

# Build the binary
make build

# Or install to GOPATH/bin
make install
```

### Pre-built Binaries

Download the latest release from the [releases page](https://github.com/johanbellander/prismer/releases).

## Quick Start

```bash
# Render the latest version
prism render ./my-project

# Render a specific version
prism render ./my-project --version v2

# Render all versions at once
prism render ./my-project --all

# Compare two versions side-by-side
prism compare ./my-project --from v1 --to v2

# List available versions
prism list --project ./my-project

# Validate structure
prism validate ./my-project

# Show version details
prism show v1
```

## Commands

### render

Render a Phase 1 structure JSON file to a visual PNG mockup.

```bash
prism render [project-path] [flags]
```

**Flags:**
- `--version`, `-v`: Version to render (default: "latest")
- `--output`, `-o`: Output file path
- `--width`, `-w`: Canvas width in pixels (default: 1200)
- `--viewport`: Target viewport (mobile, tablet, desktop)
- `--annotations`, `-a`: Include annotations
- `--all`: Render all versions found in phase1-structure directory
- `--json`: Output in JSON format

**Examples:**

```bash
# Render latest version
prism render ./my-dashboard

# Render specific version for mobile
prism render ./my-dashboard --version v2 --viewport mobile

# Render all versions at once
prism render ./my-dashboard --all

# Get JSON output
prism render ./my-dashboard --all --json
```

### validate

Validate a Phase 1 structure JSON file.

```bash
prism validate [project-path] [flags]
```

**Flags:**
- `--phase`: Phase to validate against (default: 1)
- `--json`: Output in JSON format

### list

List all available versions in a project.

```bash
prism list [flags]
```

**Flags:**
- `--project`, `-p`: Project directory path
- `--json`: Output in JSON format

### show

Show detailed information about a specific version.

```bash
prism show [version] [flags]
```

**Flags:**
- `--json`: Output in JSON format

### compare

Compare two versions side-by-side in a single PNG.

```bash
prism compare [project-path] [flags]
```

**Flags:**
- `--from`: Source version to compare from (default: "v1")
- `--to`: Target version to compare to (default: "v2")
- `--output`, `-o`: Output file path (default: {project}-compare-{from}-{to}.png)
- `--json`: Output in JSON format

**Examples:**

```bash
# Compare v1 and v2
prism compare ./my-dashboard --from v1 --to v2

# Compare with custom output
prism compare ./my-dashboard --from v1 --to approved --output comparison.png

# Get JSON output
prism compare ./my-dashboard --from v1 --to v2 --json
```

## Global Flags

- `--json`: Output in JSON format
- `--project`, `-p`: Project directory path (default: "./")
- `--quiet`, `-q`: Suppress non-essential output
- `--config`: Config file path

## Development

### Prerequisites

- Go 1.19 or later
- Make (optional, but recommended)

### Building

```bash
# Build the binary
make build

# Run tests
make test

# Format code
make fmt

# Run linters
make lint

# Show all available targets
make help
```

### Project Structure

```
prism/
├── cmd/prism/             # CLI commands
├── internal/              # Internal packages
│   ├── config/           # Configuration
│   ├── render/           # Rendering engine
│   ├── types/            # Data structures
│   └── validate/         # Validation logic
├── Makefile              # Build automation
└── README.md             # This file
```

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Status

✅ **Core Functionality Complete** - All Phase 1 features implemented and tested.

**Completed Features:**
- ✅ Render command with PNG output
- ✅ Validate command with Phase 1 constraint checking
- ✅ List command for version discovery
- ✅ Show command for version details
- ✅ Compare command for side-by-side version comparison
- ✅ Batch rendering with `--all` flag
- ✅ Multi-viewport support (mobile, tablet, desktop)
- ✅ JSON output for all commands
- ✅ Comprehensive test suite (18/18 tests passing)
- ✅ Layout calculation engine (flex, grid, stack)
- ✅ Component rendering (box, text, button, input, image)

See [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) for the development roadmap and future enhancements.
