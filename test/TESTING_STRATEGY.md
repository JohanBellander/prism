# Testing Strategy

## Overview

This document outlines the comprehensive testing strategy for the Prism mockup renderer and validation engine.

## Test Coverage

Current test coverage:
- **internal/validate**: 70.6% coverage (all tests passing)
- **internal/types**: 86.7% coverage
- **internal/render**: 3.6% coverage (basic rendering tests)
- **cmd/prism**: 0% coverage (CLI integration - tested manually)

## Unit Tests

### Validation Rules (`internal/validate/*_test.go`)

Each validation rule has comprehensive unit tests covering:

1. **Visual Hierarchy** (`hierarchy_test.go`)
   - Heading size validation
   - Button hierarchy validation
   - Valid structure tests

2. **Touch Targets & Fitts's Law** (`touch_targets_test.go`)
   - Minimum touch target size (44x44px)
   - Dangerous action spacing
   - Interactive element detection

3. **Gestalt Principles** (`gestalt_test.go`)
   - Proximity detection
   - Similarity checks
   - Related component grouping

4. **Accessibility (WCAG)** (`accessibility_test.go`)
   - Missing labels detection
   - Heading order validation
   - Nesting depth checks
   - Tab order validation

5. **Choice Overload (Hick's Law)** (`choice_overload_test.go`)
   - Navigation item limits
   - Form field limits
   - Button group limits
   - Card grid pagination

6. **Color Contrast** (`contrast_test.go`)
   - WCAG AA ratio calculation (4.5:1 for normal, 3:1 for large text)
   - Hex to RGB conversion
   - Large text detection
   - Inherited background handling

7. **Spacing Scale (8pt Grid)** (`spacing_test.go`)
   - Grid value detection
   - Nearest grid value calculation
   - Layout spacing validation
   - Component padding/margin validation

8. **Typography Scale** (`typography_test.go`)
   - Scale ratio validation
   - Token validation
   - Predefined scales
   - Scale adherence with tolerance

9. **Shadow & Elevation** (`elevation_test.go`)
   - Elevation level recommendations
   - Shadow value parsing
   - Component-specific elevation

10. **Loading States** (`loading_states_test.go`)
    - State validation (loading, error, empty, default)
    - Skeleton configuration validation
    - Skeleton element types

11. **Responsive Breakpoints** (`responsive_test.go`)
    - Viewport overflow detection
    - Touch target scaling
    - Component sizing at breakpoints

12. **Focus Indicators** (`focus_test.go`)
    - Interactive element detection
    - Focus indicator requirements

13. **Dark Mode Support** (`darkmode_test.go`)
    - Hardcoded color detection
    - Semantic color token recommendations

## Integration Tests

### Test Fixtures

Located in `test/fixtures/`:

#### Passing Fixtures (`test/fixtures/*/phase1-structure/approved.json`)
- **simple-dashboard**: Reference implementation passing all validations
- **complex-layouts**: Multi-column grid layouts
- **form-components**: Form validation examples
- **mobile-responsive**: Responsive breakpoint handling
- **edge-cases**: Edge case coverage

#### Failing Fixtures (`test/fixtures/validation-fails/phase1-structure/`)
- **v1-hierarchy-fail.json**: Poor visual hierarchy
  - Heading sizes out of order
  - Secondary button larger than primary
  
- **v2-touch-targets-fail.json**: Touch target violations
  - Buttons smaller than 44x44px
  - Dangerous actions too close together
  
- **v3-contrast-fail.json**: WCAG contrast failures
  - Low contrast text (#CCCCCC on #FFFFFF)
  - Gray on gray combinations
  
- **v4-accessibility-fail.json**: Accessibility violations
  - Missing form labels
  - Heading order jumping (h1 to h3)

### Running Integration Tests

```bash
# Test all fixtures
powershell -File test/test-all-fixtures.ps1

# Test specific fixture
.\bin\prism.exe audit ./test/fixtures/simple-dashboard

# Test validation-fails (should report failures)
.\bin\prism.exe audit ./test/fixtures/validation-fails

# JSON output for automated testing
.\bin\prism.exe audit ./test/fixtures/simple-dashboard --json
```

## Manual Testing

### CLI Commands

Test all CLI commands with various flags:

```bash
# Rendering
prism render ./test/fixtures/simple-dashboard
prism render ./test/fixtures/simple-dashboard --viewport mobile
prism render ./test/fixtures/simple-dashboard --all-viewports

# Validation
prism validate ./test/fixtures/simple-dashboard --hierarchy
prism validate ./test/fixtures/simple-dashboard --touch-targets
prism validate ./test/fixtures/simple-dashboard --json

# Audit (all validations)
prism audit ./test/fixtures/simple-dashboard
prism audit ./test/fixtures/simple-dashboard --json

# Suggestions
prism suggest ./test/fixtures/simple-dashboard
prism suggest ./test/fixtures/simple-dashboard --category forms
prism suggest ./test/fixtures/simple-dashboard --json

# List & Compare
prism list ./test/fixtures/simple-dashboard
prism show ./test/fixtures/simple-dashboard
prism compare ./test/fixtures/simple-dashboard v1 v2
```

## CI/CD Integration

### GitHub Actions (Recommended)

```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run tests
        run: go test ./... -v -cover
      - name: Test all fixtures
        run: |
          go build -o bin/prism ./cmd/prism
          ./test/test-all-fixtures.sh
```

### Pre-commit Hooks

```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running tests..."
go test ./... || exit 1

echo "Testing fixtures..."
go build -o bin/prism ./cmd/prism || exit 1
./bin/prism audit ./test/fixtures/simple-dashboard --json > /dev/null || exit 1

echo "All checks passed!"
```

## Coverage Goals

- **Phase 1 (Current)**: 70%+ coverage on validation logic
- **Phase 2 (Target)**: 80%+ coverage overall
- **Phase 3 (Stretch)**: 90%+ coverage with edge cases

## Adding New Tests

When adding new validation rules:

1. Create `<feature>_test.go` in `internal/validate/`
2. Write tests covering:
   - Default rule configuration
   - Passing cases
   - Failing cases
   - Edge cases (empty, nested, complex)
   - Custom rule configurations
3. Add test fixture in `test/fixtures/` if needed
4. Update this document
5. Run full test suite: `go test ./... -cover`

## Test Maintenance

- Run tests before committing: `go test ./...`
- Check coverage: `go test ./... -cover`
- Update test fixtures when schema changes
- Keep test data realistic and representative
- Document expected failures in fixture comments

## Known Limitations

- CLI commands (`cmd/prism`) have 0% automated test coverage - rely on manual testing
- Rendering engine has minimal test coverage (image generation hard to unit test)
- Integration tests don't cover all flag combinations

## Future Improvements

1. Add CLI integration tests using golden files
2. Visual regression testing for rendering
3. Performance benchmarking tests
4. Fuzz testing for validation rules
5. Property-based testing for edge cases
