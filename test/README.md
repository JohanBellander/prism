# Test Suite for PRISM

This directory contains comprehensive test fixtures for validating Phase 1 structure parsing, validation, and rendering.

## Test Fixtures

### 1. Simple Dashboard (`simple-dashboard`)
Basic dashboard layout demonstrating core Phase 1 features:
- Multiple versions (v1, v2, approved)
- Component nesting (boxes, text, buttons)
- Basic styling (colors, borders, padding)
- Version management
- Invalid fixture for testing validation

**Key tests:**
- Version listing and resolution
- Approved vs latest version handling
- Color validation (invalid-colors fixture)

### 2. Complex Layouts (`complex-layouts`)
Advanced layout patterns:
- Deep component nesting (header, sidebar, main content)
- Flexbox-style layouts with gaps and alignment
- Multi-level component hierarchies
- Spacer components for flexible layouts
- Grid-like stat cards

**Key tests:**
- Complex nested structures
- Horizontal and vertical layout combinations
- Component positioning and sizing

### 3. Form Components (`form-components`)
Form elements and input validation:
- Input fields with labels
- Buttons with various styles
- Form layout patterns
- Placeholder text
- Checkbox components
- Social login buttons

**Key tests:**
- Input rendering
- Label/input associations
- Form field styling
- Button variations

### 4. Edge Cases (`edge-cases`)
Boundary conditions and special cases:
- Zero padding (minimum values)
- All color formats (hex3, hex6, rgb, rgba)
- Border radius variations (0px to 50%)
- Empty inputs and buttons
- Transparency (rgba)
- Mixed component types

**Key tests:**
- Color parsing (all supported formats)
- Border radius handling (including percentage)
- Empty content handling
- Minimum/maximum values

### 5. Mobile Responsive (`mobile-responsive`)
Mobile-first design patterns:
- Mobile viewport (375x667)
- Bottom navigation
- Touch-friendly spacing
- Responsive card layouts
- Transaction list patterns
- Hero cards

**Key tests:**
- Small viewport rendering
- Viewport scaling to larger sizes
- Mobile UI patterns

## Running Tests

### Quick Test (PowerShell)
```powershell
.\test\test-all-fixtures.ps1
```

This script will:
1. Build the project
2. Validate all fixtures
3. List versions for each project
4. Render all fixtures to PNG
5. Test viewport variations
6. Run unit tests
7. Test invalid fixtures (should fail)
8. Display summary report

### Individual Tests

**Validate a specific fixture:**
```powershell
.\bin\prism.exe validate --project test\fixtures\complex-layouts
```

**List versions:**
```powershell
.\bin\prism.exe list --project test\fixtures\form-components
```

**Render specific fixture:**
```powershell
.\bin\prism.exe render --project test\fixtures\mobile-responsive --output test\output\mobile.png
```

**Test different viewports:**
```powershell
# Mobile
.\bin\prism.exe render --project test\fixtures\complex-layouts --output test\output\mobile.png --width 375 --height 667

# Tablet
.\bin\prism.exe render --project test\fixtures\complex-layouts --output test\output\tablet.png --width 768 --height 1024

# Desktop
.\bin\prism.exe render --project test\fixtures\complex-layouts --output test\output\desktop.png --width 1200 --height 800
```

**Test with scaling:**
```powershell
# 2x scale for retina displays
.\bin\prism.exe render --project test\fixtures\mobile-responsive --output test\output\mobile-2x.png --scale 2
```

**JSON output:**
```powershell
.\bin\prism.exe validate --project test\fixtures\edge-cases --json
.\bin\prism.exe list --project test\fixtures\form-components --json
```

### Unit Tests

Run Go unit tests for internal packages:
```powershell
go test ./internal/types -v
go test ./... -v
```

## Expected Behavior

### Valid Fixtures
All fixtures except `simple-dashboard/phase1-structure/invalid-colors.json` should:
- Pass validation
- Render successfully to PNG
- Support all viewport sizes
- Support 1x and 2x scaling

### Invalid Fixtures
The `invalid-colors.json` fixture should:
- Fail validation with color format error
- Not render

## Output

Test outputs are saved to:
```
test/output/
├── Simple-Dashboard/
│   └── output.png
├── Complex-Layouts/
│   └── output.png
├── Form-Components/
│   └── output.png
├── Edge-Cases/
│   └── output.png
└── Mobile-Responsive/
    ├── output.png
    ├── tablet.png
    └── desktop.png
```

## Coverage

These fixtures test:
- ✅ All Phase 1 component types (box, text, button, input, image)
- ✅ All layout properties (position, width, height, padding, gap, direction, align, flex)
- ✅ All style properties (background, color, border, borderRadius, fontSize, fontWeight, padding)
- ✅ Color formats (hex3, hex6, rgb, rgba)
- ✅ Nesting depth validation (max 10 levels)
- ✅ Required field validation
- ✅ Invalid color format detection
- ✅ Version management (v1, v2, approved)
- ✅ Viewport handling
- ✅ Scaling (1x, 2x)
- ✅ JSON and human-readable output

## Adding New Fixtures

To add a new test fixture:

1. Create directory structure:
   ```
   test/fixtures/your-fixture-name/
   └── phase1-structure/
       └── v1.json
   ```

2. Add valid Phase 1 structure JSON

3. Add to test script in `test-all-fixtures.ps1`

4. Document in this README

## Notes

- All fixtures use valid Phase 1 structure format
- Color values test all supported formats
- Layout calculations are basic (Phase 1)
- Font rendering uses Go's built-in font
- Images are placeholder boxes (image URLs not fetched in Phase 1)
