# Design Principles Integration Plan for PRISM

## Executive Summary

This plan outlines how to incorporate industry-standard design principles into PRISM to ensure it produces accessible, usable, and visually effective UI mockups while maintaining the two-phase design process.

---

## Phase 1: Validation & Enforcement (Structural Principles)

These principles should be **automatically validated** during Phase 1 to catch issues early.

### 1.1 Visual Hierarchy Validation

**Goal**: Ensure clear visual hierarchy through size, spacing, and positioning.

**Implementation**:
- Add `prism validate --hierarchy` command
- Check that heading sizes follow a consistent scale (e.g., h1 > h2 > h3)
- Warn if primary CTA is smaller than secondary actions
- Validate that spacing increases with hierarchy level (more padding for important elements)

**Validation Rules**:
```go
type HierarchyRule struct {
    HeadingScaleRatio  float64  // e.g., 1.25 (each level 25% larger)
    MinPrimaryCTASize  int      // e.g., 120px width minimum
    SpacingScaleRatio  float64  // e.g., 1.5 (parent spacing > child spacing)
}
```

**Example Output**:
```
❌ Hierarchy Issue: Secondary button (140px) larger than primary (120px)
⚠️  Warning: h2 (18px) not sufficiently smaller than h1 (20px) - recommend 1.25x scale
✅ Spacing hierarchy is consistent
```

---

### 1.2 Touch Target & Fitts's Law Validation

**Goal**: Ensure interactive elements meet minimum touch target sizes.

**Implementation**:
- Extend existing touch_targets_min validation (currently 44px)
- Add proximity checks for frequently-used actions
- Validate dangerous actions (delete, submit) have adequate spacing from other buttons

**Validation Rules**:
```go
type TouchTargetRule struct {
    MinSize            int      // 44px (iOS) or 48px (Android)
    MinSpacing         int      // 8px between interactive elements
    DangerousSpacing   int      // 16px for destructive actions
    FrequentActions    []string // IDs of common actions to check proximity
}
```

**Example Output**:
```
❌ Touch Target: "close_button" is 32x32px (requires 44x44px minimum)
❌ Spacing: "delete_btn" only 4px from "cancel_btn" (requires 16px for destructive actions)
✅ All primary actions meet touch target requirements
```

---

### 1.3 Gestalt Principles Validation

**Goal**: Validate that related elements are grouped using proximity and similarity.

**Implementation**:
- Analyze spacing patterns to detect groupings
- Check that related components (e.g., form fields in a group) have consistent spacing
- Warn if unrelated items are closer than related items

**Validation Rules**:
```go
type GestaltRule struct {
    IntraGroupSpacing  int      // e.g., 8px within groups
    InterGroupSpacing  int      // e.g., 24px between groups
    MinGroupSize       int      // e.g., 2 items to form a group
    SimilarityCheck    bool     // Check if similar items look alike
}
```

**Example Output**:
```
⚠️  Proximity: "username_label" is 16px from "username_input" but 12px from "password_label"
   Suggestion: Increase spacing to 24px between form groups
✅ All card components use consistent spacing (16px)
```

---

### 1.4 Accessibility (WCAG) Validation

**Goal**: Catch accessibility issues in structure before they reach Phase 2.

**Implementation**:
- Validate semantic structure (headings, labels, hierarchy)
- Check for missing labels on interactive elements
- Ensure logical keyboard navigation order (based on layout order)
- Validate max nesting depth (already exists: 4 levels)

**Validation Rules**:
```go
type A11yRule struct {
    RequireLabels         bool     // All interactive elements need labels
    RequireHeadingOrder   bool     // h1 → h2 → h3 (no skipping)
    MaxNestingDepth       int      // 4 levels
    RequireFocusIndicator bool     // All interactive elements
    CheckTabOrder         bool     // Verify logical tab sequence
}
```

**Example Output**:
```
❌ A11y: "search_input" missing label
❌ A11y: Heading structure jumps from h1 to h3 (missing h2)
⚠️  A11y: Tab order may be confusing - "submit_btn" comes before "email_input" in layout
✅ All interactive elements have focus indicators defined
```

---

### 1.5 Hick's Law Validation (Choice Overload)

**Goal**: Warn when too many choices are presented at once.

**Implementation**:
- Count interactive elements in a single container
- Warn if navigation has too many top-level items
- Suggest progressive disclosure for complex forms

**Validation Rules**:
```go
type ChoiceRule struct {
    MaxNavItems        int      // e.g., 7 ± 2 (Miller's Law)
    MaxFormFields      int      // e.g., 5-7 per section
    MaxButtonGroup     int      // e.g., 3 buttons max in a group
    MaxCardGrid        int      // e.g., 12 cards before pagination
}
```

**Example Output**:
```
⚠️  Choice Overload: Navigation has 12 items - consider grouping or secondary menu
⚠️  Choice Overload: Form section has 8 fields - consider splitting into steps
✅ Button groups contain 3 or fewer options
```

---

## Phase 2: Design System Enhancements (Visual Principles)

These principles guide **design token creation** and **rendering** in Phase 2.

### 2.1 Color Contrast Validation (WCAG)

**Goal**: Ensure all text meets WCAG AA contrast ratios (4.5:1 for normal text, 3:1 for large text).

**Implementation**:
- Add contrast calculation to Phase 2 validation
- Check text color against background color
- Validate interactive element states (hover, focus, disabled)
- Provide suggestions for compliant alternatives

**Design Token Enhancement**:
```json
{
  "colors": {
    "primary": {
      "600": "#3B82F6",
      "700": "#2563EB",
      "contrast_ratio_white": 4.5,  // Auto-calculated
      "wcag_aa_compliant": true
    }
  }
}
```

**Validation Output**:
```
❌ Contrast: "secondary_text" (#999999) on white fails WCAG AA (2.8:1, requires 4.5:1)
   Suggestion: Use #767676 or darker for compliance
✅ "primary_button" text (#FFFFFF) on primary.600 passes (4.6:1)
```

---

### 2.2 Typography Scale (Visual Hierarchy)

**Goal**: Enforce consistent, harmonious type scales.

**Implementation**:
- Provide predefined type scale ratios (1.125, 1.200, 1.250, 1.333, 1.414, 1.618)
- Validate that custom scales follow mathematical progression
- Check line-height and letter-spacing for readability

**Design Token Enhancement**:
```json
{
  "typography": {
    "scale_ratio": 1.250,  // Major Third
    "base_size": 16,
    "sizes": {
      "xs": 12,   // 16 / 1.25^1
      "sm": 14,   // 16 / 1.25^0.5
      "base": 16,
      "lg": 20,   // 16 * 1.25
      "xl": 25,   // 16 * 1.25^2
      "2xl": 31,  // 16 * 1.25^3
      "3xl": 39   // 16 * 1.25^4
    },
    "line_height": {
      "tight": 1.25,
      "normal": 1.5,
      "relaxed": 1.75
    }
  }
}
```

---

### 2.3 Spacing Scale (8pt Grid System)

**Goal**: Use consistent, predictable spacing based on 8pt grid.

**Implementation**:
- Default spacing scale: 4, 8, 12, 16, 24, 32, 48, 64, 96, 128
- Validate that all margins/padding use scale values
- Allow 4px for fine-tuning, but warn if used excessively

**Design Token Enhancement**:
```json
{
  "spacing": {
    "base_unit": 8,
    "scale": {
      "0": 0,
      "1": 4,   // 0.5 * base (half step for fine-tuning)
      "2": 8,   // 1 * base
      "3": 12,  // 1.5 * base
      "4": 16,  // 2 * base
      "6": 24,  // 3 * base
      "8": 32,  // 4 * base
      "12": 48, // 6 * base
      "16": 64, // 8 * base
      "24": 96, // 12 * base
      "32": 128 // 16 * base
    }
  }
}
```

**Validation Output**:
```
⚠️  Spacing: "card_padding" uses 15px (not on 8pt grid)
   Suggestion: Use 16px (spacing.4) for consistency
✅ All container spacing follows 8pt grid system
```

---

### 2.4 Shadow & Elevation System

**Goal**: Create consistent depth hierarchy using shadows.

**Implementation**:
- Predefined elevation levels (0-5)
- Each level has specific shadow definition
- Validate that interactive elements have appropriate elevation changes on hover

**Design Token Enhancement**:
```json
{
  "elevation": {
    "0": "none",
    "1": "0 1px 2px 0 rgba(0,0,0,0.05)",      // Subtle (cards)
    "2": "0 2px 4px 0 rgba(0,0,0,0.1)",       // Raised (buttons)
    "3": "0 4px 8px 0 rgba(0,0,0,0.12)",      // Floating (dropdowns)
    "4": "0 8px 16px 0 rgba(0,0,0,0.15)",     // Overlays (modals)
    "5": "0 16px 32px 0 rgba(0,0,0,0.2)"      // Maximum (important dialogs)
  }
}
```

---

## Phase 3: Rendering Enhancements

### 3.1 Perceived Performance (Loading States)

**Goal**: Support skeleton screens and loading states in mockups.

**Implementation**:
- Add `state: "loading"` property to components
- Render skeleton versions (gray boxes with subtle animation hints)
- Support placeholder content

**Component Schema Addition**:
```json
{
  "id": "user_card",
  "type": "box",
  "state": "loading",  // NEW: loading | error | empty | default
  "layout": {
    "background": "#F3F4F6",
    "padding": 16
  },
  "skeleton": {
    "avatar": { "type": "circle", "size": 48 },
    "title": { "type": "text", "width": "60%" },
    "subtitle": { "type": "text", "width": "80%" }
  }
}
```

---

### 3.2 Responsive Breakpoints

**Goal**: Validate designs work at standard breakpoints.

**Implementation**:
- Add `--viewport` flag to render command
- Validate at: mobile (375px), tablet (768px), desktop (1440px)
- Check that touch targets scale appropriately
- Warn if content overflows

**Command Enhancement**:
```bash
# Render for all breakpoints
prism render ./project --all-viewports

# Specific viewport
prism render ./project --viewport mobile

# Custom width
prism render ./project --width 320
```

**Validation Output**:
```
✅ Mobile (375px): All touch targets ≥ 44px
⚠️  Tablet (768px): "sidebar" width fixed at 320px (42% of viewport)
❌ Mobile (375px): "nav_items" overflow container by 50px
   Suggestion: Use responsive layout or hamburger menu
```

---

### 3.3 Dark Mode Support

**Goal**: Enable automatic dark mode variant generation.

**Implementation**:
- Define semantic color tokens (not just hex values)
- Support `--mode dark` flag for rendering
- Validate contrast in both light and dark modes

**Design Token Enhancement**:
```json
{
  "semantic_colors": {
    "background": {
      "light": "#FFFFFF",
      "dark": "#1F2937"
    },
    "text": {
      "primary": {
        "light": "#111827",
        "dark": "#F9FAFB"
      },
      "secondary": {
        "light": "#6B7280",
        "dark": "#9CA3AF"
      }
    },
    "surface": {
      "light": "#F9FAFB",
      "dark": "#111827"
    }
  }
}
```

---

### 3.4 Focus Indicators (Accessibility)

**Goal**: Render visible focus states on interactive elements.

**Implementation**:
- Add `--show-focus` flag to render command
- Display 2px outline around focusable elements
- Use high-contrast focus color (minimum 3:1 ratio)

**Render Enhancement**:
```bash
# Show focus indicators on all interactive elements
prism render ./project --show-focus

# Show specific interaction state
prism render ./project --state hover
prism render ./project --state active
prism render ./project --state disabled
```

---

## Phase 4: New Validation Commands

### 4.1 Comprehensive Design Audit

**Goal**: Single command to validate all design principles.

**Implementation**:
```bash
prism audit ./project --report audit-report.json

# Specific categories
prism audit ./project --category accessibility
prism audit ./project --category hierarchy
prism audit ./project --category performance
```

**Audit Report Structure**:
```json
{
  "project": "my-project",
  "version": "v1",
  "phase": "structure",
  "audit_date": "2025-10-25T12:00:00Z",
  "results": {
    "accessibility": {
      "score": 85,
      "issues": [
        {
          "severity": "error",
          "rule": "missing_label",
          "component": "search_input",
          "message": "Interactive element missing accessible label"
        }
      ]
    },
    "visual_hierarchy": {
      "score": 92,
      "issues": []
    },
    "touch_targets": {
      "score": 100,
      "issues": []
    },
    "contrast": {
      "score": 78,
      "issues": [
        {
          "severity": "warning",
          "rule": "wcag_aa_contrast",
          "component": "secondary_text",
          "message": "Contrast ratio 2.8:1 fails WCAG AA (requires 4.5:1)"
        }
      ]
    }
  },
  "overall_score": 86,
  "recommendations": [
    "Add label to search_input for accessibility",
    "Increase contrast of secondary_text to #767676 or darker"
  ]
}
```

---

### 4.2 Best Practice Suggestions

**Goal**: Provide actionable suggestions based on design patterns.

**Implementation**:
```bash
prism suggest ./project --category forms
prism suggest ./project --category navigation
prism suggest ./project --category layouts
```

**Example Output**:
```
📋 Form Best Practices:
✅ Labels are above inputs (good for mobile)
⚠️  Consider adding field descriptions for complex inputs
💡 Suggestion: Group related fields with spacing.6 (24px) between groups

🧭 Navigation Best Practices:
✅ Primary navigation is in expected location (top)
⚠️  7 navigation items detected - consider dropdown for less common items
💡 Suggestion: Add visual indicator for current page
```

---

## Phase 5: Implementation Roadmap

### Milestone 1: Core Validation (v0.3.0)
**Timeline**: 2-3 weeks

- [ ] Implement touch target validation (extends existing)
- [ ] Add visual hierarchy checks (heading scale, CTA sizing)
- [ ] Implement basic accessibility validation (labels, heading order)
- [ ] Add spacing consistency checks (Gestalt proximity)
- [ ] Create validation report output (JSON + human-readable)

**Deliverables**:
- `prism validate --strict` (fails on any errors)
- `prism validate --report validation.json`
- Updated error messages with specific suggestions

---

### Milestone 2: Design System Enhancements (v0.4.0)
**Timeline**: 3-4 weeks

- [ ] Implement WCAG contrast calculation
- [ ] Add typography scale validation
- [ ] Enforce 8pt grid system
- [ ] Create elevation/shadow token system
- [ ] Add dark mode token support

**Deliverables**:
- `prism validate --phase2` (design-specific checks)
- Contrast validation with auto-suggestions
- Design token templates for common systems

---

### Milestone 3: Responsive & States (v0.5.0)
**Timeline**: 3-4 weeks

- [ ] Implement responsive breakpoint validation
- [ ] Add viewport-specific rendering (`--viewport mobile`)
- [ ] Support component states (loading, hover, focus, disabled)
- [ ] Add `--show-focus` rendering mode
- [ ] Implement skeleton/loading state rendering

**Deliverables**:
- `prism render --all-viewports`
- `prism render --state focus`
- `prism validate --responsive`

---

### Milestone 4: Comprehensive Auditing (v0.6.0)
**Timeline**: 2-3 weeks

- [ ] Implement `prism audit` command
- [ ] Create scoring system for all categories
- [ ] Generate detailed audit reports (JSON + PDF)
- [ ] Add best practice suggestions
- [ ] Create comparison reports (v1 vs v2)

**Deliverables**:
- `prism audit ./project --report audit.json`
- `prism suggest ./project`
- Automated improvement recommendations

---

### Milestone 5: Advanced Features (v0.7.0)
**Timeline**: 3-4 weeks

- [ ] Implement Gestalt principle analysis
- [ ] Add Hick's Law validation (choice overload)
- [ ] Create pattern library detection (recognize common UI patterns)
- [ ] Add Jakob's Law suggestions (compare to common patterns)
- [ ] Implement progressive disclosure analysis

**Deliverables**:
- `prism analyze --patterns`
- `prism compare --against industry-standard`
- Smart suggestions based on detected UI patterns

---

## Testing Strategy

### Unit Tests
- Test each validation rule independently
- Mock component structures for predictable testing
- Cover edge cases (nested layouts, complex grids)

### Integration Tests
- Test full validation pipeline
- Test rendering with various flags
- Test audit report generation

### Validation Dataset
Create test mockups that deliberately violate principles:
- `test-hierarchy-fail.json` - Poor visual hierarchy
- `test-touch-targets-fail.json` - Too-small interactive elements
- `test-contrast-fail.json` - WCAG failures
- `test-accessibility-fail.json` - Missing labels, bad heading order
- `test-all-pass.json` - Passes all validations (reference implementation)

---

## Documentation Updates

### 1. Update DESIGNPROCESS.md
Add section: "Design Principles Validation" explaining:
- What gets validated in Phase 1 vs Phase 2
- How to interpret validation errors
- How to use audit reports

### 2. Create VALIDATION_RULES.md
Document all validation rules:
- Rule name and category
- Why it matters (UX principle)
- How it's checked
- How to fix violations
- Examples of pass/fail

### 3. Create DESIGN_TOKENS.md
Comprehensive guide to:
- Typography scales
- Color systems and contrast
- Spacing (8pt grid)
- Elevation/shadows
- Responsive breakpoints
- Dark mode tokens

### 4. Update CLI Help
Add detailed help for all new commands:
```bash
prism validate --help
prism audit --help
prism suggest --help
prism render --help  # Update with new flags
```

---

## Success Metrics

### Quantitative
- **Validation Coverage**: 90%+ of WCAG AA criteria automated
- **Accuracy**: <5% false positives in validation
- **Performance**: Validation completes in <2 seconds for typical mockup
- **Adoption**: Users run `prism audit` before finalizing Phase 1

### Qualitative
- Users understand WHY validation fails (clear error messages)
- Suggestions are actionable (specific fixes, not vague)
- Validation catches real usability issues
- Design system tokens are intuitive to use

---

## Open Questions

1. **Configurability**: Should users be able to disable specific rules or customize thresholds?
   - Proposal: Add `.prismrc.json` config file with rule overrides

2. **Plugin System**: Should validation rules be pluggable for custom company standards?
   - Proposal: v1.0 feature - allow custom rule modules

3. **AI Suggestions**: Should we use LLM to suggest improvements based on context?
   - Proposal: Future exploration - may require external API

4. **Platform-Specific Rules**: Should iOS vs Android vs Web have different validation?
   - Proposal: Add `--platform ios|android|web` flag with platform-specific rules

5. **Performance Budget**: Should we validate asset sizes, image dimensions?
   - Proposal: v0.8.0 feature - add performance category to audit

---

## Conclusion

This plan systematically incorporates industry-standard design principles into PRISM across 5 major milestones. The phased approach ensures:

1. **Phase 1 validation** catches structural and accessibility issues early
2. **Phase 2 validation** ensures design system compliance and visual quality
3. **Rendering enhancements** support modern UI patterns (dark mode, responsive, states)
4. **Audit system** provides comprehensive quality scoring
5. **Extensibility** allows future additions without breaking existing workflows

**Next Steps**:
1. Review and approve this plan
2. Prioritize Milestone 1 tasks
3. Create GitHub issues for each task
4. Begin implementation with touch target validation (extends existing code)

**Estimated Total Timeline**: 13-17 weeks for full implementation (v0.3.0 → v0.7.0)
