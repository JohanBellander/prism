# PRISM 🎨

**Phase Render & Inspection for Structural Mockups**

> Transform JSON mockups into visual wireframes instantly

PRISM is a CLI tool that generates visual PNG representations of Phase 1 structural mockups from your AI Design Agent workflow. It bridges the gap between JSON structure files and visual wireframes, enabling fast review and approval cycles.

Stop squinting at JSON. See your Phase 1 structures as actual wireframes in seconds.

## Features

✨ **Zero-config rendering** - Point at a project directory and get PNGs  
🖼️ **Multi-viewport support** - Mobile, tablet, and desktop layouts  
🔍 **Side-by-side comparison** - Visual diff between versions  
⚡ **Batch rendering** - Render all versions with `--all` flag  
✅ **Phase 1 validation** - Strict constraint checking (black/white only, no styling)  
🤖 **Agent-friendly** - `--json` output for programmatic integration  
📋 **Version tracking** - List, show, and compare structure versions  
🎯 **Layout engine** - Flexbox, grid, and stack layouts  
🧩 **Component library** - Box, text, button, input, image rendering

## Installation

**Quick install (Unix/macOS/Linux):**

```bash
curl -fsSL https://raw.githubusercontent.com/JohanBellander/prism/master/install.sh | bash
```

**Quick install (Windows PowerShell):**

```powershell
iwr -useb https://raw.githubusercontent.com/JohanBellander/prism/master/install.ps1 | iex
```

**From source:**

Requires Go 1.19 or later:

```bash
git clone https://github.com/johanbellander/prism.git
cd prism
make build
```

The binary will be in `bin/prism.exe` (Windows) or `bin/prism` (Unix).

Or install directly to your Go bin:

```bash
make install
```

## Quick Start

### For Humans

PRISM works with the [AI Design Agent two-phase process](DESIGNPROCESS.md). After your agent creates Phase 1 structure files, visualize them instantly:

```bash
# Render the latest version
prism render ./my-dashboard

# See all versions
prism list --project ./my-dashboard

# Compare two versions
prism compare ./my-dashboard --from v1 --to v2
```

Your agent creates the JSON. PRISM makes it visible.

### For AI Agents

PRISM integrates seamlessly into agentic workflows:

```bash
# After creating phase1-structure/v1.json, render it
prism render ./project --version v1 --json

# Validate before rendering
prism validate ./project --json

# Batch render all versions for review
prism render ./project --all --json

# Compare versions programmatically
prism compare ./project --from v1 --to v2 --json
```

All commands support `--json` for structured output your agent can parse.

All commands support `--json` for structured output your agent can parse.

## Usage

### Rendering Mockups

```bash
# Render latest version
prism render ./my-dashboard

# Render specific version
prism render ./my-dashboard --version v2

# Render for different viewports
prism render ./my-dashboard --viewport mobile
prism render ./my-dashboard --viewport tablet
prism render ./my-dashboard --viewport desktop

# Custom dimensions
prism render ./my-dashboard --width 1920 --height 1080

# Render all versions at once
prism render ./my-dashboard --all

# Save to specific location
prism render ./my-dashboard --output mockups/v1.png

# JSON output for agents
prism render ./my-dashboard --json
```

### Comparing Versions

Side-by-side visual comparison of structure changes:

```bash
# Compare v1 and v2
prism compare ./my-dashboard --from v1 --to v2

# Compare with custom output path
prism compare ./my-dashboard --from v1 --to approved --output diff.png

# JSON output
prism compare ./my-dashboard --from v1 --to v2 --json
```

### Validating Structures

Ensure Phase 1 constraints are met (black/white only, no styling):

```bash
# Validate project
prism validate ./my-dashboard

# Validate specific phase
prism validate ./my-dashboard --phase 1

# JSON output for CI/CD
prism validate ./my-dashboard --json
```

### Listing Versions

```bash
# List all versions in project
prism list --project ./my-dashboard

# JSON output
prism list --project ./my-dashboard --json
```

### Showing Version Details

```bash
# Show version metadata
prism show v1

# JSON output
prism show v1 --json
```

prism show v1 --json
```

## Global Flags

All commands support these flags:

- `--json` - Output in JSON format for programmatic use
- `--project`, `-p` - Project directory path (default: `./`)
- `--quiet`, `-q` - Suppress non-essential output
- `--config` - Config file path (default: `~/.prism`)

## The Design Process

PRISM is built for the [two-phase AI design workflow](DESIGNPROCESS.md):

**Phase 1**: Structure-only mockups (black & white, no styling)  
**Phase 2**: Design system styling (applied after Phase 1 approval)

PRISM renders Phase 1 structures, making them visual and reviewable before moving to Phase 2. This prevents the common problem of structural changes sneaking in during the styling phase.

Your AI agent creates `phase1-structure/v1.json`, `v2.json`, etc. PRISM turns them into PNGs. You approve. The agent moves to Phase 2.

Simple. Fast. No surprises.

## Project Structure

```
your-project/
├── phase1-structure/       # Created by your AI agent
│   ├── v1.json            # First structure iteration
│   ├── v2.json            # Revised structure
│   └── approved.json      # Approved for Phase 2
└── mockups/               # Created by PRISM
    ├── v1.png
    ├── v2.png
    └── approved.png
```

## Integration Examples

### CI/CD Pipeline

```bash
# Validate all structures in CI
prism validate ./project --json | jq '.valid'

# Generate mockups for pull request reviews
prism render ./project --all
```

### Git Hooks

```bash
# pre-commit hook: validate before committing
#!/bin/bash
if ! prism validate . --quiet; then
    echo "Phase 1 validation failed"
    exit 1
fi
```

### Agent Workflow

Your AI agent can use PRISM programmatically:

```python
import subprocess
import json

# After creating structure file
result = subprocess.run(
    ["prism", "render", "./project", "--version", "v1", "--json"],
    capture_output=True,
    text=True
)
output = json.loads(result.stdout)
print(f"Rendered: {output['output_path']}")
```

## Development

### Prerequisites

- Go 1.19 or later
- Make

### Building

```bash
# Build
make build

# Run tests
make test

# Format and lint
make fmt
make lint

# See all targets
make help
```

### Code Structure

```
prism/
├── cmd/prism/             # CLI commands (render, validate, list, show, compare)
├── internal/
│   ├── render/           # Rendering engine (layout calculation, PNG generation)
│   └── types/            # Data structures (Phase 1 schema)
├── test/
│   └── fixtures/         # Test structures and expected outputs
├── Makefile
└── README.md
```

## Documentation

- [README.md](README.md) - You are here!
- [DESIGNPROCESS.md](DESIGNPROCESS.md) - Two-phase design workflow guide
- [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) - Development roadmap and architecture
- [prism-spec.md](prism-spec.md) - Complete technical specification
- [AGENTS.md](AGENTS.md) - AI agent integration guide

## Status

✅ **Production Ready** - All Phase 1 features complete and tested

**Implemented:**
- Render command with PNG output
- Validate command with Phase 1 constraints
- List command for version discovery
- Show command for version details
- Compare command for side-by-side diffs
- Batch rendering with `--all`
- Multi-viewport support (mobile/tablet/desktop)
- JSON output for all commands
- Layout engine (flexbox, grid, stack)
- Component rendering (box, text, button, input, image)
- 18/18 tests passing

**Roadmap:**
See [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) for Phase 2 and Phase 3 plans (annotations, CI/CD integration, interactive viewer).

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## About

PRISM - Making AI-designed mockups visible

Built for developers working with AI design agents who need fast visual feedback on structural iterations.
