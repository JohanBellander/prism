# PRISM Validation Rules Reference

This document provides comprehensive details on all validation rules enforced by PRISM. Each rule is designed to ensure your UI follows established UX principles and accessibility standards.

## Table of Contents

- [Phase 1: Structural Validation](#phase-1-structural-validation)
  - [Visual Hierarchy](#visual-hierarchy)
  - [Touch Targets & Fitts's Law](#touch-targets--fittss-law)
  - [Gestalt Principles](#gestalt-principles)
  - [Accessibility (WCAG)](#accessibility-wcag)
  - [Choice Overload (Hick's Law)](#choice-overload-hicks-law)
- [Phase 2: Visual Design Validation](#phase-2-visual-design-validation)
  - [Color Contrast](#color-contrast)
  - [Typography Scale](#typography-scale)
  - [Spacing (8pt Grid)](#spacing-8pt-grid)
  - [Elevation & Shadows](#elevation--shadows)
  - [Loading States](#loading-states)
  - [Responsive Design](#responsive-design)
  - [Focus Indicators](#focus-indicators)
  - [Dark Mode Support](#dark-mode-support)
- [Severity Levels](#severity-levels)
- [Quick Reference](#quick-reference)

---

# Phase 1: Structural Validation

These rules validate the fundamental structure and usability of your interface before visual design is applied.

## Visual Hierarchy

**Category**: Structure  
**Command**: `prism validate --hierarchy`  
**Why it matters**: Users need to instantly understand what's most important on the page. Poor hierarchy causes confusion and increases cognitive load.

### Rules

#### 1. Heading Size Ratio

**Requirement**: Headings must be at least 1.5x larger than body text

**Why**: Clear size differentiation helps users scan content and understand structure

**How it's checked**:
```
heading_size / body_text_size >= 1.5
```

**Examples**:

✅ **PASS**:
```json
{
  "components": [
    {"type": "text", "size": "2xl", "role": "heading"},  // 24px
    {"type": "text", "size": "base", "role": "body"}     // 16px
  ]
}
// Ratio: 24/16 = 1.5 ✓
```

❌ **FAIL**:
```json
{
  "components": [
    {"type": "text", "size": "lg", "role": "heading"},   // 18px
    {"type": "text", "size": "base", "role": "body"}     // 16px
  ]
}
// Ratio: 18/16 = 1.125 ✗ (too small)
```

**How to fix**:
- Increase heading size to at least 24px (1.5x of 16px)
- Or use size tokens: `2xl` for headings, `base` for body
- Consider larger ratios (2x, 2.5x) for more dramatic hierarchy

---

#### 2. Nesting Depth Limit

**Requirement**: Component nesting must not exceed 4 levels deep

**Why**: Deep nesting increases cognitive complexity and makes layouts harder to understand

**How it's checked**:
```
max_depth(component_tree) <= 4
```

**Examples**:

✅ **PASS**:
```json
{
  "components": [                          // Level 1
    {
      "id": "container",
      "children": [                        // Level 2
        {
          "id": "card",
          "children": [                    // Level 3
            {
              "id": "content",
              "children": [                // Level 4
                {"id": "text", "content": "Hello"}
              ]
            }
          ]
        }
      ]
    }
  ]
}
// Depth: 4 ✓
```

❌ **FAIL**:
```json
{
  "components": [                          // Level 1
    {"children": [                         // Level 2
      {"children": [                       // Level 3
        {"children": [                     // Level 4
          {"children": [                   // Level 5 ✗
            {"id": "text"}
          ]}
        ]}
      ]}
    ]}
  ]
}
// Depth: 5 ✗ (too deep)
```

**How to fix**:
- Flatten the component structure
- Combine multiple wrapper components into one
- Use layout properties instead of nested containers
- Refactor deeply nested lists into flat structures

---

#### 3. Primary Action Prominence

**Requirement**: Primary CTA must be visually larger or more prominent than secondary actions

**Why**: Users need clear guidance on the most important action to take

**How it's checked**:
```
primary_button_size > secondary_button_size
OR
primary_button_weight > secondary_button_weight
```

**Examples**:

✅ **PASS**:
```json
{
  "components": [
    {
      "id": "save-button",
      "type": "button",
      "role": "primary",
      "size": "lg",           // 18px
      "weight": "bold"
    },
    {
      "id": "cancel-button",
      "type": "button",
      "role": "secondary",
      "size": "base",         // 16px
      "weight": "normal"
    }
  ]
}
// Primary is larger and bolder ✓
```

❌ **FAIL**:
```json
{
  "components": [
    {
      "id": "save-button",
      "type": "button",
      "role": "primary",
      "size": "base"          // Same size
    },
    {
      "id": "cancel-button",
      "type": "button",
      "role": "secondary",
      "size": "base"          // Same size
    }
  ]
}
// No visual differentiation ✗
```

**How to fix**:
- Increase primary button font size by at least one scale step
- Make primary button text bold (`weight: "bold"`)
- Increase primary button padding
- In Phase 2: Use stronger background color for primary

---

## Touch Targets & Fitts's Law

**Category**: Interaction Design  
**Command**: `prism validate --touch-targets`  
**Why it matters**: Small or closely-packed interactive elements are difficult to tap accurately, especially on mobile devices. This causes frustration and accessibility issues.

### Rules

#### 1. Minimum Touch Target Size

**Requirement**: All interactive elements must be at least 44x44px

**Why**: Based on iOS Human Interface Guidelines and WCAG 2.1 Level AA (Success Criterion 2.5.5)

**Standards**:
- iOS HIG: 44x44pt minimum
- Android Material: 48x48dp minimum
- WCAG 2.1 AA: 44x44px minimum
- PRISM uses: **44x44px** (most conservative)

**How it's checked**:
```
interactive_element.width >= 44 AND interactive_element.height >= 44
```

**Examples**:

✅ **PASS**:
```json
{
  "id": "submit-button",
  "type": "button",
  "layout": {
    "width": 120,
    "height": 44,
    "padding": 12
  }
}
// Total touch area: 120x44 ✓
```

✅ **PASS** (with padding):
```json
{
  "id": "icon-button",
  "type": "button",
  "layout": {
    "width": 24,      // Icon size
    "height": 24,
    "padding": 10     // Creates 44x44 total area
  }
}
// Total touch area: 24 + (10*2) = 44x44 ✓
```

❌ **FAIL**:
```json
{
  "id": "close-icon",
  "type": "button",
  "layout": {
    "width": 32,
    "height": 32,
    "padding": 0
  }
}
// Total touch area: 32x32 ✗ (too small)
```

**How to fix**:
- Add padding: `(44 - current_size) / 2`
- For 32px icon: add 6px padding → 32 + (6*2) = 44px
- For 24px icon: add 10px padding → 24 + (10*2) = 44px
- Increase button dimensions directly

---

#### 2. Interactive Element Spacing

**Requirement**: Interactive elements must have at least 8px spacing between them

**Why**: Prevents accidental taps on adjacent buttons (Fitts's Law)

**How it's checked**:
```
distance_between(element1, element2) >= 8
```

**Examples**:

✅ **PASS**:
```json
{
  "layout": {
    "direction": "horizontal",
    "gap": 8                    // ✓ Minimum spacing
  },
  "children": [
    {"type": "button", "content": "Cancel"},
    {"type": "button", "content": "Save"}
  ]
}
```

❌ **FAIL**:
```json
{
  "layout": {
    "direction": "horizontal",
    "gap": 4                    // ✗ Too close
  },
  "children": [
    {"type": "button", "content": "Cancel"},
    {"type": "button", "content": "Save"}
  ]
}
```

**How to fix**:
- Increase `gap` to at least 8px
- Add margin to individual buttons
- Use 12px or 16px spacing for better safety margin

---

#### 3. Dangerous Action Spacing

**Requirement**: Destructive actions must be at least 16px away from other interactive elements

**Why**: Prevents accidental deletion/destructive operations

**How it's checked**:
```
if button.role == "destructive":
    distance_to_nearest_interactive >= 16
```

**Examples**:

✅ **PASS**:
```json
{
  "layout": {
    "direction": "horizontal",
    "gap": 16
  },
  "children": [
    {"type": "button", "content": "Cancel", "role": "secondary"},
    {"type": "button", "content": "Delete", "role": "destructive"}
  ]
}
// Destructive action has 16px spacing ✓
```

❌ **FAIL**:
```json
{
  "layout": {
    "direction": "horizontal",
    "gap": 8
  },
  "children": [
    {"type": "button", "content": "Cancel"},
    {"type": "button", "content": "Delete", "role": "destructive"}
  ]
}
// Only 8px spacing for dangerous action ✗
```

**How to fix**:
- Increase spacing to 16px minimum
- Place destructive actions on opposite side of dialog
- Consider requiring confirmation for destructive actions
- Use 24px spacing for extra safety

---

## Gestalt Principles

**Category**: Visual Psychology  
**Command**: `prism validate --gestalt`  
**Why it matters**: Human brains naturally group visual elements based on proximity, similarity, and continuity. Following Gestalt principles makes interfaces more intuitive.

### Rules

#### 1. Proximity (Intra-group Spacing)

**Requirement**: Elements within a group must be closer together than elements between groups

**Formula**: `intra_group_spacing < inter_group_spacing`

**Why**: The Law of Proximity states that objects near each other are perceived as a group

**How it's checked**:
```
max(spacing_within_group) < min(spacing_between_groups)
```

**Examples**:

✅ **PASS**:
```json
{
  "components": [
    {
      "id": "form-group-1",
      "layout": {"gap": 8},        // Within group
      "children": [
        {"type": "label", "content": "Email"},
        {"type": "input", "name": "email"}
      ]
    },
    {
      "id": "form-group-2",
      "layout": {"margin_top": 24}, // Between groups
      "children": [
        {"type": "label", "content": "Password"},
        {"type": "input", "name": "password"}
      ]
    }
  ]
}
// Within: 8px, Between: 24px → 8 < 24 ✓
```

❌ **FAIL**:
```json
{
  "components": [
    {
      "id": "form-group-1",
      "layout": {"gap": 16},       // Within group
      "children": [
        {"type": "label"},
        {"type": "input"}
      ]
    },
    {
      "id": "form-group-2",
      "layout": {"margin_top": 16}, // Same spacing!
      "children": [
        {"type": "label"},
        {"type": "input"}
      ]
    }
  ]
}
// Within: 16px, Between: 16px → ambiguous grouping ✗
```

**How to fix**:
- Use 2-3x larger spacing between groups
- Common ratios: 8px within, 24px between
- Or: 12px within, 32px between
- Ensure clear visual separation

---

#### 2. Similarity (Consistent Styling)

**Requirement**: Similar elements must have similar visual properties

**Why**: The Law of Similarity states that similar-looking objects are perceived as related

**How it's checked**:
```
foreach element_type:
    variance(size, color, weight) <= threshold
```

**Examples**:

✅ **PASS**:
```json
{
  "components": [
    {
      "type": "button",
      "role": "primary",
      "size": "base",
      "weight": "medium"
    },
    {
      "type": "button",
      "role": "primary",
      "size": "base",      // Consistent
      "weight": "medium"    // Consistent
    }
  ]
}
// All primary buttons styled identically ✓
```

❌ **FAIL**:
```json
{
  "components": [
    {
      "type": "button",
      "role": "primary",
      "size": "base",
      "weight": "medium"
    },
    {
      "type": "button",
      "role": "primary",
      "size": "lg",        // Different size
      "weight": "bold"     // Different weight
    }
  ]
}
// Inconsistent styling for same role ✗
```

**How to fix**:
- Standardize all buttons of the same type
- Create component variants (primary, secondary, tertiary)
- Document when to use each variant
- Use design tokens for consistency

---

#### 3. Continuity (Visual Flow)

**Requirement**: Elements should follow predictable visual patterns (left-to-right, top-to-bottom)

**Why**: The Law of Continuity states that elements arranged in a line or curve are perceived as related

**How it's checked**:
```
alignment_check(elements) == consistent
flow_direction(elements) == predictable
```

**Examples**:

✅ **PASS**:
```json
{
  "layout": {
    "direction": "vertical",
    "alignment": "left"
  },
  "children": [
    {"type": "text", "content": "Title", "align": "left"},
    {"type": "text", "content": "Body", "align": "left"},
    {"type": "button", "content": "Action", "align": "left"}
  ]
}
// Consistent left alignment creates visual flow ✓
```

❌ **FAIL**:
```json
{
  "children": [
    {"type": "text", "content": "Title", "align": "left"},
    {"type": "text", "content": "Body", "align": "center"},
    {"type": "button", "content": "Action", "align": "right"}
  ]
}
// Inconsistent alignment breaks visual flow ✗
```

**How to fix**:
- Choose one primary alignment (usually left for LTR languages)
- Use center alignment sparingly (headers, empty states)
- Maintain consistent alignment within sections
- Follow reading direction (left-to-right for English)

---

## Accessibility (WCAG)

**Category**: Accessibility  
**Command**: `prism validate --accessibility`  
**Why it matters**: 15% of the world's population has some form of disability. Accessible design ensures everyone can use your interface.

### Rules

#### 1. Form Input Labels

**Requirement**: All form inputs must have associated labels

**Standard**: WCAG 2.1 Level A (Success Criterion 1.3.1, 3.3.2)

**Why**: Screen readers need labels to announce input purpose

**How it's checked**:
```
foreach input:
    has_label(input) OR has_aria_label(input)
```

**Examples**:

✅ **PASS**:
```json
{
  "components": [
    {
      "type": "label",
      "for": "email-input",
      "content": "Email Address"
    },
    {
      "type": "input",
      "id": "email-input",
      "name": "email"
    }
  ]
}
// Input has associated label ✓
```

✅ **PASS** (with aria-label):
```json
{
  "type": "input",
  "name": "search",
  "accessibility": {
    "aria_label": "Search products"
  }
}
// Input has aria-label ✓
```

❌ **FAIL**:
```json
{
  "type": "input",
  "name": "email",
  "placeholder": "Enter email"  // Placeholder is NOT a label
}
// No label or aria-label ✗
```

**How to fix**:
- Add a `<label>` element with `for` attribute
- Or add `aria-label` attribute to input
- Or use `aria-labelledby` to reference existing text
- Don't rely on placeholder text alone

---

#### 2. Heading Hierarchy

**Requirement**: Headings must follow logical order (h1 → h2 → h3, no skipping)

**Standard**: WCAG 2.1 Level A (Success Criterion 1.3.1)

**Why**: Screen readers use heading structure for navigation

**How it's checked**:
```
heading_sequence = extract_heading_levels(page)
is_sequential(heading_sequence) == true
```

**Examples**:

✅ **PASS**:
```json
{
  "components": [
    {"type": "text", "size": "4xl", "role": "h1", "content": "Page Title"},
    {"type": "text", "size": "2xl", "role": "h2", "content": "Section"},
    {"type": "text", "size": "xl", "role": "h3", "content": "Subsection"}
  ]
}
// Sequence: h1 → h2 → h3 ✓
```

❌ **FAIL**:
```json
{
  "components": [
    {"type": "text", "size": "4xl", "role": "h1", "content": "Page Title"},
    {"type": "text", "size": "xl", "role": "h3", "content": "Section"}
  ]
}
// Sequence: h1 → h3 (skipped h2) ✗
```

**How to fix**:
- Start with h1 for page title
- Use h2 for main sections
- Use h3 for subsections within h2
- Never skip heading levels
- Only one h1 per page

---

#### 3. Focus Indicators

**Requirement**: All interactive elements must have visible focus states

**Standard**: WCAG 2.1 Level AA (Success Criterion 2.4.7)

**Why**: Keyboard users need to see which element has focus

**How it's checked**:
```
foreach interactive_element:
    has_focus_state(element) == true
    focus_indicator_visible == true
```

**Examples**:

✅ **PASS**:
```json
{
  "id": "submit-button",
  "type": "button",
  "states": {
    "focus": {
      "outline": "2px solid",
      "outline_color": "#2563eb",
      "outline_offset": 2
    }
  }
}
// Focus state defined with visible outline ✓
```

❌ **FAIL**:
```json
{
  "id": "submit-button",
  "type": "button",
  "states": {
    "focus": {
      "outline": "none"  // Removes focus indicator!
    }
  }
}
// Focus indicator removed ✗
```

**How to fix**:
- Always define explicit focus states
- Use `outline: 2px solid` with high-contrast color
- Add `outline_offset: 2px` for better visibility
- Never set `outline: none` without a replacement
- Ensure 3:1 contrast ratio with background

---

#### 4. Semantic Structure

**Requirement**: Use semantic roles for components (header, navigation, main, footer, button, etc.)

**Standard**: WCAG 2.1 Level A (Success Criterion 4.1.2)

**Why**: Screen readers use semantic roles to understand page structure

**How it's checked**:
```
foreach component:
    has_semantic_role(component) == true
```

**Examples**:

✅ **PASS**:
```json
{
  "components": [
    {"id": "top", "role": "header"},
    {"id": "nav", "role": "navigation"},
    {"id": "content", "role": "main"},
    {"id": "bottom", "role": "footer"}
  ]
}
// All major sections have semantic roles ✓
```

❌ **FAIL**:
```json
{
  "components": [
    {"id": "top", "type": "box"},          // No role
    {"id": "nav", "type": "box"},          // No role
    {"id": "content", "type": "box"},      // No role
  ]
}
// Generic boxes without semantic meaning ✗
```

**How to fix**:
- Use `role: "header"` for page header
- Use `role: "navigation"` for nav menus
- Use `role: "main"` for primary content
- Use `role: "footer"` for page footer
- Use `role: "button"` for clickable actions
- Use `role: "form"` for form containers

---

## Choice Overload (Hick's Law)

**Category**: Cognitive Psychology  
**Command**: `prism validate --choice-overload`  
**Why it matters**: Decision time increases logarithmically with the number of choices (Hick's Law). Too many options overwhelm users.

### Rules

#### 1. Navigation Item Limit

**Requirement**: Primary navigation must have ≤ 7 items

**Why**: Based on Miller's Law (7±2 items in working memory) and Hick's Law

**How it's checked**:
```
count(navigation_items) <= 7
```

**Examples**:

✅ **PASS**:
```json
{
  "role": "navigation",
  "children": [
    {"type": "link", "content": "Home"},
    {"type": "link", "content": "Products"},
    {"type": "link", "content": "Pricing"},
    {"type": "link", "content": "About"},
    {"type": "link", "content": "Contact"}
  ]
}
// 5 navigation items ✓
```

❌ **FAIL**:
```json
{
  "role": "navigation",
  "children": [
    {"type": "link", "content": "Home"},
    {"type": "link", "content": "Products"},
    {"type": "link", "content": "Services"},
    {"type": "link", "content": "Pricing"},
    {"type": "link", "content": "Resources"},
    {"type": "link", "content": "Blog"},
    {"type": "link", "content": "About"},
    {"type": "link", "content": "Careers"},
    {"type": "link", "content": "Contact"}
  ]
}
// 9 navigation items ✗ (too many)
```

**How to fix**:
- Combine similar items (Resources → includes Blog, Docs, Help)
- Use dropdown menus for sub-items
- Move secondary items to footer
- Prioritize most important 5-7 items
- Use mega-menu for complex navigation

---

#### 2. Form Field Limit

**Requirement**: Form sections must have ≤ 5 fields visible at once

**Why**: Long forms are abandoned more frequently. Progressive disclosure reduces cognitive load.

**How it's checked**:
```
foreach form_section:
    count(visible_fields) <= 5
```

**Examples**:

✅ **PASS**:
```json
{
  "role": "form",
  "children": [
    {
      "id": "section-1",
      "children": [
        {"type": "input", "name": "first_name"},
        {"type": "input", "name": "last_name"},
        {"type": "input", "name": "email"}
      ]
    }
  ]
}
// 3 fields in section ✓
```

✅ **PASS** (multi-step):
```json
{
  "role": "form",
  "steps": [
    {
      "title": "Step 1: Personal Info",
      "fields": [
        {"type": "input", "name": "first_name"},
        {"type": "input", "name": "last_name"},
        {"type": "input", "name": "email"}
      ]
    },
    {
      "title": "Step 2: Address",
      "fields": [
        {"type": "input", "name": "street"},
        {"type": "input", "name": "city"},
        {"type": "input", "name": "zip"}
      ]
    }
  ]
}
// Each step has ≤5 fields ✓
```

❌ **FAIL**:
```json
{
  "role": "form",
  "children": [
    {"type": "input", "name": "first_name"},
    {"type": "input", "name": "last_name"},
    {"type": "input", "name": "email"},
    {"type": "input", "name": "phone"},
    {"type": "input", "name": "street"},
    {"type": "input", "name": "city"},
    {"type": "input", "name": "state"},
    {"type": "input", "name": "zip"}
  ]
}
// 8 fields visible at once ✗ (too many)
```

**How to fix**:
- Split into multiple steps/pages
- Group related fields into collapsible sections
- Show only required fields initially
- Use progressive disclosure for optional fields
- Aim for 3-5 fields per section

---

#### 3. Button Group Limit

**Requirement**: Action groups must have ≤ 3 primary actions

**Why**: Multiple primary actions create decision paralysis

**How it's checked**:
```
count(buttons.role == "primary") <= 3
```

**Examples**:

✅ **PASS**:
```json
{
  "components": [
    {"type": "button", "role": "primary", "content": "Save Changes"},
    {"type": "button", "role": "secondary", "content": "Cancel"},
    {"type": "button", "role": "tertiary", "content": "Reset"}
  ]
}
// 1 primary action ✓
```

❌ **FAIL**:
```json
{
  "components": [
    {"type": "button", "role": "primary", "content": "Save Draft"},
    {"type": "button", "role": "primary", "content": "Save & Continue"},
    {"type": "button", "role": "primary", "content": "Save & Exit"},
    {"type": "button", "role": "primary", "content": "Publish"}
  ]
}
// 4 primary actions ✗ (confusing)
```

**How to fix**:
- Choose ONE primary action (most common/important)
- Make other actions secondary or tertiary
- Hide advanced options in a menu
- Use smart defaults to reduce choices
- Consider user's primary goal

---

# Phase 2: Visual Design Validation

These rules validate visual polish and design system compliance after structure is approved.

## Color Contrast

**Category**: Accessibility  
**Command**: `prism validate --contrast`  
**Why it matters**: Low contrast text is difficult to read, especially for users with visual impairments. WCAG compliance is often legally required.

### Rules

#### 1. Normal Text Contrast (WCAG AA)

**Requirement**: Normal text must have ≥ 4.5:1 contrast ratio with background

**Standard**: WCAG 2.1 Level AA (Success Criterion 1.4.3)

**Definition**: Normal text is < 24px or < 19px bold

**How it's checked**:
```
if text_size < 24 OR (text_size < 19 AND weight < bold):
    contrast_ratio(foreground, background) >= 4.5
```

**Examples**:

✅ **PASS**:
```json
{
  "type": "text",
  "size": "base",           // 16px
  "color": "#171717",       // neutral.900
  "background": "#FFFFFF"   // white
}
// Contrast: 16.1:1 ✓ (exceeds 4.5:1)
```

❌ **FAIL**:
```json
{
  "type": "text",
  "size": "base",           // 16px
  "color": "#A3A3A3",       // neutral.400 (light gray)
  "background": "#FFFFFF"   // white
}
// Contrast: 2.9:1 ✗ (below 4.5:1)
```

**How to fix**:
- Use darker text colors: `neutral.900` (#171717) on light backgrounds
- Use lighter text colors: `neutral.50` on dark backgrounds
- Check contrast with online tools (WebAIM, Coolors)
- Use design tokens that guarantee WCAG compliance

**WCAG Levels**:
- **Level AA** (minimum): 4.5:1 for normal text, 3:1 for large text
- **Level AAA** (enhanced): 7:1 for normal text, 4.5:1 for large text

---

#### 2. Large Text Contrast (WCAG AA)

**Requirement**: Large text must have ≥ 3:1 contrast ratio with background

**Standard**: WCAG 2.1 Level AA (Success Criterion 1.4.3)

**Definition**: Large text is ≥ 24px or ≥ 19px bold

**How it's checked**:
```
if text_size >= 24 OR (text_size >= 19 AND weight >= bold):
    contrast_ratio(foreground, background) >= 3.0
```

**Examples**:

✅ **PASS**:
```json
{
  "type": "text",
  "size": "3xl",            // 30px (large text)
  "color": "#737373",       // neutral.500
  "background": "#FFFFFF"   // white
}
// Contrast: 4.6:1 ✓ (exceeds 3:1)
```

❌ **FAIL**:
```json
{
  "type": "text",
  "size": "2xl",            // 24px (large text)
  "color": "#D4D4D4",       // neutral.300 (very light)
  "background": "#FFFFFF"   // white
}
// Contrast: 1.8:1 ✗ (below 3:1)
```

**How to fix**:
- For large headings, use at least `neutral.500` on light backgrounds
- For dark backgrounds, use at least `neutral.400`
- Check if text truly qualifies as "large" (≥24px)
- Consider using AAA standard (4.5:1) for better readability

---

#### 3. Interactive Element Contrast

**Requirement**: Interactive element text must have ≥ 4.5:1 contrast (or 3:1 if large)

**Standard**: WCAG 2.1 Level AA (Success Criterion 1.4.3)

**Why**: Buttons, links, and interactive elements must be readable

**How it's checked**:
```
foreach interactive_element:
    contrast_ratio(text_color, button_background) >= 4.5
```

**Examples**:

✅ **PASS**:
```json
{
  "type": "button",
  "role": "primary",
  "style": {
    "background": "#2563eb",  // primary.600
    "color": "#FFFFFF"         // white
  }
}
// Contrast: 8.6:1 ✓
```

❌ **FAIL**:
```json
{
  "type": "button",
  "role": "secondary",
  "style": {
    "background": "#E5E5E5",  // neutral.200 (light gray)
    "color": "#FFFFFF"         // white
  }
}
// Contrast: 1.2:1 ✗ (white on light gray)
```

**How to fix**:
- Use dark text on light buttons: `neutral.900` on `neutral.200`
- Use white text on dark buttons: `white` on `primary.600`
- Check hover/active states too
- Disabled states can be lower contrast (but should still be visible)

---

#### 4. Non-Text Contrast

**Requirement**: UI components and graphical objects must have ≥ 3:1 contrast

**Standard**: WCAG 2.1 Level AA (Success Criterion 1.4.11)

**Applies to**: Icons, borders, focus indicators, chart elements

**How it's checked**:
```
contrast_ratio(component_color, background) >= 3.0
```

**Examples**:

✅ **PASS**:
```json
{
  "type": "input",
  "style": {
    "border": "1px solid",
    "border_color": "#737373",  // neutral.500
    "background": "#FFFFFF"      // white
  }
}
// Border contrast: 4.6:1 ✓
```

❌ **FAIL**:
```json
{
  "type": "input",
  "style": {
    "border": "1px solid",
    "border_color": "#E5E5E5",  // neutral.200
    "background": "#FFFFFF"      // white
  }
}
// Border contrast: 1.2:1 ✗ (too low)
```

**How to fix**:
- Use at least `neutral.400` for borders on white backgrounds
- Use at least `neutral.600` for better visibility
- Check focus indicators meet 3:1 contrast
- Icons should have 3:1 contrast or be accompanied by text

---

## Typography Scale

**Category**: Design System  
**Command**: `prism validate --typography`  
**Why it matters**: Consistent typography creates visual harmony and makes designs easier to maintain.

### Rules

#### 1. Font Size Scale Compliance

**Requirement**: All font sizes must be from the defined scale

**Scale**: 12, 14, 16, 18, 20, 24, 30, 36 (or custom scale with consistent ratio)

**Why**: Arbitrary font sizes create visual inconsistency

**How it's checked**:
```
is_on_scale(font_size, allowed_sizes, tolerance=0.5)
```

**Examples**:

✅ **PASS**:
```json
{
  "components": [
    {"type": "text", "size": "xs", "content": "Label"},      // 12px
    {"type": "text", "size": "base", "content": "Body"},     // 16px
    {"type": "text", "size": "2xl", "content": "Heading"}    // 24px
  ]
}
// All sizes on scale ✓
```

❌ **FAIL**:
```json
{
  "components": [
    {"type": "text", "size": 15, "content": "Body"},     // Not on scale
    {"type": "text", "size": 22, "content": "Heading"}   // Not on scale
  ]
}
// 15px and 22px are not in scale ✗
```

**How to fix**:
- Round 15px → 16px (base)
- Round 22px → 24px (2xl)
- Use size tokens: `xs`, `sm`, `base`, `lg`, `xl`, `2xl`, `3xl`, `4xl`
- Avoid arbitrary pixel values

**Typography Tokens**:
```json
{
  "xs": 12,
  "sm": 14,
  "base": 16,
  "lg": 18,
  "xl": 20,
  "2xl": 24,
  "3xl": 30,
  "4xl": 36
}
```

---

#### 2. Type Scale Ratio

**Requirement**: Custom scales must maintain consistent ratio (1.125 - 1.618)

**Common Ratios**:
- 1.125 (Major Second) - tight scale
- 1.200 (Minor Third) - default
- 1.250 (Major Third) - relaxed
- 1.333 (Perfect Fourth) - spacious
- 1.414 (Augmented Fourth) - dramatic
- 1.618 (Golden Ratio) - very dramatic

**How it's checked**:
```
ratio = size[n+1] / size[n]
1.125 <= ratio <= 1.618
```

**Examples**:

✅ **PASS** (1.200 ratio):
```json
{
  "typography_scale": {
    "base": 16,
    "lg": 19.2,    // 16 * 1.2
    "xl": 23.04,   // 19.2 * 1.2
    "2xl": 27.65   // 23.04 * 1.2
  }
}
// Consistent 1.2 ratio ✓
```

❌ **FAIL**:
```json
{
  "typography_scale": {
    "base": 16,
    "lg": 18,      // 1.125 ratio
    "xl": 24,      // 1.333 ratio (inconsistent!)
    "2xl": 28      // 1.167 ratio (inconsistent!)
  }
}
// Ratio varies between steps ✗
```

**How to fix**:
- Choose one ratio and stick to it
- Use a type scale generator (type-scale.com)
- Round to whole pixels but maintain approximate ratio
- PRISM default: 1.200 ratio starting from 16px base

---

## Spacing (8pt Grid)

**Category**: Design System  
**Command**: `prism validate --spacing`  
**Why it matters**: Consistent spacing creates visual rhythm and makes designs scalable across devices.

### Rules

#### 1. 8pt Grid Compliance

**Requirement**: All spacing values must be multiples of 8 (or 4 for fine-tuning)

**Scale**: 0, 4, 8, 12, 16, 24, 32, 48, 64, 96, 128

**Why**: Ensures consistency, scalability, and pixel-perfect rendering on all screens

**How it's checked**:
```
spacing % 4 == 0
```

**Examples**:

✅ **PASS**:
```json
{
  "layout": {
    "padding": 16,      // 16 = 8*2 ✓
    "margin": 24,       // 24 = 8*3 ✓
    "gap": 8            // 8 = 8*1 ✓
  }
}
```

❌ **FAIL**:
```json
{
  "layout": {
    "padding": 15,      // Not divisible by 4 ✗
    "margin": 22,       // Not divisible by 4 ✗
    "gap": 10           // Not divisible by 4 ✗
  }
}
```

**How to fix**:
- Round 15px → 16px
- Round 22px → 24px
- Round 10px → 8px or 12px
- Use spacing tokens: 4, 8, 12, 16, 24, 32, 48, 64, 96, 128

**When to use 4px (half-step)**:
- Fine-tuning optical alignment
- Icon spacing in buttons
- Compensating for font metrics
- **Not** for general layout spacing

---

#### 2. Spacing Scale Ratio

**Requirement**: Spacing should follow a consistent pattern (usually exponential)

**Why**: Predictable spacing scale makes design decisions easier

**PRISM Scale**: 0, 4, 8, 12, 16, 24, 32, 48, 64, 96, 128

**Pattern**:
- 0-16: increments of 4
- 16-32: increments of 8
- 32+: increments of 16 or more

**Examples**:

✅ **PASS**:
```json
{
  "card": {
    "padding": 16,        // Standard padding
    "gap": 12,            // Internal spacing
    "margin_bottom": 24   // Card separation
  }
}
// Uses scale values ✓
```

❌ **FAIL**:
```json
{
  "card": {
    "padding": 18,        // Not on scale
    "gap": 14,            // Not on scale
    "margin_bottom": 26   // Not on scale
  }
}
// Random values off scale ✗
```

**How to fix**:
- Consult spacing scale: 4, 8, 12, 16, 24, 32, 48, 64, 96, 128
- Round to nearest scale value
- Use spacing tokens in code
- 18px → 16px, 14px → 12px or 16px, 26px → 24px

---

## Elevation & Shadows

**Category**: Design System  
**Command**: `prism validate --elevation`  
**Why it matters**: Shadows indicate layering and interactivity. Consistent elevation creates depth hierarchy.

### Rules

#### 1. Shadow Definition Limit

**Requirement**: Use ≤ 4 unique shadow definitions

**Why**: Too many shadow variants create visual inconsistency

**PRISM Elevation Levels**:
- **Level 0**: No shadow (flat elements)
- **Level 1**: Subtle shadow (cards, default state)
- **Level 2**: Raised shadow (hover state, dropdowns)
- **Level 3**: Floating shadow (modals, popovers)

**How it's checked**:
```
unique_shadows(design) <= 4
```

**Examples**:

✅ **PASS**:
```json
{
  "shadows": {
    "sm": "0 1px 2px 0 rgb(0 0 0 / 0.05)",
    "md": "0 4px 6px -1px rgb(0 0 0 / 0.1)",
    "lg": "0 10px 15px -3px rgb(0 0 0 / 0.1)",
    "xl": "0 20px 25px -5px rgb(0 0 0 / 0.1)"
  }
}
// 4 elevation levels ✓
```

❌ **FAIL**:
```json
{
  "shadows": {
    "shadow1": "0 1px 3px rgba(0,0,0,0.1)",
    "shadow2": "0 2px 4px rgba(0,0,0,0.12)",
    "shadow3": "0 3px 5px rgba(0,0,0,0.15)",
    "shadow4": "0 4px 6px rgba(0,0,0,0.18)",
    "shadow5": "0 5px 7px rgba(0,0,0,0.2)",
    "shadow6": "0 10px 20px rgba(0,0,0,0.25)"
  }
}
// 6 unique shadows ✗ (too many)
```

**How to fix**:
- Consolidate similar shadows
- Stick to 3-4 elevation levels
- Use elevation tokens: `sm`, `md`, `lg`, `xl`
- Remove shadows that differ by tiny amounts

---

#### 2. Elevation Appropriateness

**Requirement**: Component elevation must match its importance and interactivity

**Guidelines**:
- **Level 0** (no shadow): Flat UI, backgrounds, contained content
- **Level 1** (sm): Default cards, panels, tiles
- **Level 2** (md): Hovered buttons, dropdowns, tooltips
- **Level 3** (lg/xl): Modals, dialogs, important floating content

**How it's checked**:
```
component.elevation <= component.importance_level
```

**Examples**:

✅ **PASS**:
```json
{
  "components": [
    {"type": "card", "shadow": "sm"},           // Resting state
    {"type": "dropdown", "shadow": "md"},       // Popup element
    {"type": "modal", "shadow": "xl"}           // Top layer
  ]
}
// Appropriate elevation for each component ✓
```

❌ **FAIL**:
```json
{
  "components": [
    {"type": "card", "shadow": "xl"},           // Too dramatic
    {"type": "modal", "shadow": "sm"}           // Too subtle
  ]
}
// Mismatched elevation and component type ✗
```

**How to fix**:
- Cards/panels: `sm` shadow
- Dropdowns/tooltips: `md` shadow
- Modals/dialogs: `lg` or `xl` shadow
- Don't over-elevate decorative elements
- Use higher elevation for interactive/temporary elements

---

## Loading States

**Category**: Feedback  
**Command**: `prism validate --loading-states`  
**Why it matters**: Users need feedback during async operations. Missing loading states cause confusion.

### Rules

#### 1. Loading Indicator Presence

**Requirement**: Async actions must have loading indicators

**Why**: Users need to know the system is working

**How it's checked**:
```
foreach async_action:
    has_loading_state(action) == true
```

**Examples**:

✅ **PASS**:
```json
{
  "id": "submit-form",
  "type": "button",
  "states": {
    "default": {"content": "Submit"},
    "loading": {
      "content": "Submitting...",
      "icon": "spinner",
      "disabled": true
    },
    "success": {"content": "Submitted!"},
    "error": {"content": "Failed - Retry"}
  }
}
// Loading state defined ✓
```

❌ **FAIL**:
```json
{
  "id": "submit-form",
  "type": "button",
  "states": {
    "default": {"content": "Submit"}
  }
}
// No loading state ✗
```

**How to fix**:
- Add `loading` state to buttons that trigger async actions
- Show spinner icon during loading
- Disable button to prevent double-submission
- Provide success/error feedback after completion

---

#### 2. Skeleton Screens

**Requirement**: Data-heavy components should have skeleton screens

**Why**: Skeleton screens feel faster than spinners and reduce perceived load time

**How it's checked**:
```
if component.loads_data:
    has_skeleton_state(component) == true
```

**Examples**:

✅ **PASS**:
```json
{
  "id": "user-profile",
  "type": "card",
  "states": {
    "loading": {
      "type": "skeleton",
      "elements": [
        {"type": "skeleton_circle", "size": 64},
        {"type": "skeleton_line", "width": "60%"},
        {"type": "skeleton_line", "width": "40%"}
      ]
    },
    "loaded": {
      "children": [
        {"type": "image", "src": "avatar.jpg"},
        {"type": "text", "content": "John Doe"},
        {"type": "text", "content": "@johndoe"}
      ]
    }
  }
}
// Skeleton state defined ✓
```

❌ **FAIL**:
```json
{
  "id": "user-profile",
  "type": "card",
  "children": [
    {"type": "image", "src": "avatar.jpg"},
    {"type": "text", "content": "John Doe"}
  ]
}
// No loading state, no skeleton ✗
```

**How to fix**:
- Add skeleton shapes matching final content layout
- Use gray rectangles for text lines
- Use gray circles for avatars
- Animate skeleton with subtle pulse or shimmer

---

## Responsive Design

**Category**: Multi-Device Support  
**Command**: `prism validate --responsive`  
**Why it matters**: 60%+ of web traffic is mobile. Non-responsive designs fail on small screens.

### Rules

#### 1. Mobile Breakpoint Defined

**Requirement**: Must define layout changes for ≤ 640px width

**Why**: Mobile is the most common viewport size

**How it's checked**:
```
has_breakpoint(design, max_width <= 640) == true
```

**Examples**:

✅ **PASS**:
```json
{
  "responsive": {
    "mobile": {
      "breakpoint": 640,
      "changes": {
        "layout.direction": "vertical",
        "layout.padding": 16,
        "sidebar.display": "none"
      }
    }
  }
}
// Mobile breakpoint defined ✓
```

❌ **FAIL**:
```json
{
  "layout": {
    "width": 1200,
    "padding": 24
  }
}
// No responsive breakpoints ✗
```

**How to fix**:
- Add mobile breakpoint (640px)
- Stack horizontal layouts vertically
- Reduce padding (24px → 16px)
- Hide or collapse secondary content
- Ensure touch targets remain ≥44px

---

#### 2. Flexible Layouts

**Requirement**: Use relative units (%, fr, flex) instead of fixed widths

**Why**: Fixed widths cause horizontal scrolling on small screens

**How it's checked**:
```
foreach component:
    uses_relative_units(component.width) == true
```

**Examples**:

✅ **PASS**:
```json
{
  "layout": {
    "display": "flex",
    "max_width": 1200,       // Max, not fixed
    "width": "100%"          // Relative
  }
}
// Flexible layout ✓
```

❌ **FAIL**:
```json
{
  "layout": {
    "width": 1200            // Fixed width
  }
}
// Fixed layout causes overflow ✗
```

**How to fix**:
- Use `max-width` instead of `width`
- Use percentages: `width: "100%"`
- Use flexbox or grid
- Ensure minimum width doesn't exceed 320px

---

## Focus Indicators

**Category**: Accessibility  
**Command**: `prism validate --focus`  
**Why it matters**: Keyboard users rely on focus indicators for navigation.

### Rules

#### 1. Focus State Defined

**Requirement**: All interactive elements must have explicit focus states

**Standard**: WCAG 2.1 Level AA (Success Criterion 2.4.7)

**How it's checked**:
```
foreach interactive_element:
    has_focus_state(element) == true
```

**Examples**:

✅ **PASS**:
```json
{
  "type": "button",
  "states": {
    "focus": {
      "outline": "2px solid",
      "outline_color": "#2563eb",
      "outline_offset": 2
    }
  }
}
// Focus state defined ✓
```

❌ **FAIL**:
```json
{
  "type": "button",
  "style": {
    "outline": "none"   // Removes default focus!
  }
}
// No focus state ✗
```

**How to fix**:
- Never remove outlines without replacement
- Use `outline: 2px solid [color]`
- Add `outline_offset: 2px` for spacing
- Ensure high contrast (3:1 minimum)

---

#### 2. Focus Indicator Contrast

**Requirement**: Focus indicators must have ≥ 3:1 contrast with background

**Standard**: WCAG 2.1 Level AA (Success Criterion 1.4.11)

**How it's checked**:
```
contrast_ratio(focus_color, background_color) >= 3.0
```

**Examples**:

✅ **PASS**:
```json
{
  "states": {
    "focus": {
      "outline_color": "#2563eb",  // primary.600
      "background": "#FFFFFF"       // white
    }
  }
}
// Contrast: 8.6:1 ✓
```

❌ **FAIL**:
```json
{
  "states": {
    "focus": {
      "outline_color": "#BFDBFE",  // primary.200 (light)
      "background": "#FFFFFF"       // white
    }
  }
}
// Contrast: 1.4:1 ✗ (too low)
```

**How to fix**:
- Use dark focus colors on light backgrounds (e.g., `primary.600`)
- Use light focus colors on dark backgrounds (e.g., `primary.300`)
- Check contrast ratio with online tools
- Avoid light blue on white or dark blue on black

---

## Dark Mode Support

**Category**: Theme Support  
**Command**: `prism validate --dark-mode`  
**Why it matters**: Dark mode reduces eye strain and is expected by users. Simple color inversion often breaks contrast.

### Rules

#### 1. Dark Mode Palette Defined

**Requirement**: Must define dark mode color variants

**Why**: Inverting light mode colors rarely works well

**How it's checked**:
```
has_dark_mode_palette(design) == true
```

**Examples**:

✅ **PASS**:
```json
{
  "theme": {
    "light": {
      "background": "#FFFFFF",
      "surface": "#F5F5F5",
      "text_primary": "#171717",
      "text_secondary": "#737373"
    },
    "dark": {
      "background": "#0A0A0A",
      "surface": "#171717",
      "text_primary": "#FAFAFA",
      "text_secondary": "#A3A3A3"
    }
  }
}
// Dark mode palette defined ✓
```

❌ **FAIL**:
```json
{
  "theme": {
    "background": "#FFFFFF",
    "text": "#000000"
  }
}
// Only light mode colors ✗
```

**How to fix**:
- Define separate dark mode palette
- Don't just invert colors
- Test contrast ratios in both modes
- Adjust semantic colors (success, error, warning)

---

#### 2. Dark Mode Contrast Compliance

**Requirement**: Dark mode must maintain WCAG AA contrast ratios

**Why**: Dark mode can have worse contrast than light mode if not carefully designed

**How it's checked**:
```
foreach color_pair in dark_mode:
    contrast_ratio(foreground, background) >= 4.5
```

**Examples**:

✅ **PASS**:
```json
{
  "dark_mode": {
    "text": "#FAFAFA",        // neutral.50
    "background": "#0A0A0A"   // Very dark
  }
}
// Contrast: 17.8:1 ✓
```

❌ **FAIL**:
```json
{
  "dark_mode": {
    "text": "#737373",        // neutral.500
    "background": "#262626"   // neutral.800
  }
}
// Contrast: 2.8:1 ✗ (too low)
```

**How to fix**:
- Use very light text on very dark backgrounds
- `neutral.50` (#FAFAFA) on `neutral.900` (#171717) or darker
- Check all color pairs, not just text
- Don't use pure black (#000000) for backgrounds (use #0A0A0A instead)

---

# Severity Levels

Validation issues are classified by severity:

| Level | Icon | Description | Action Required |
|-------|------|-------------|-----------------|
| **Critical** | 🔴 | Accessibility violations, WCAG failures, broken UX | Must fix before approval |
| **Warning** | 🟠 | Design principle violations, degraded UX | Should fix before Phase 2 |
| **Info** | 🟡 | Style guide violations, minor inconsistencies | Consider fixing for polish |

**Scoring**:
- Each validator returns 0-100 score
- **90-100**: Excellent
- **70-89**: Good (passing)
- **50-69**: Fair (needs work)
- **0-49**: Poor (failing)

**Passing Threshold**: 70+ for all validators

---

# Quick Reference

## Common Validation Failures & Fixes

| Issue | Quick Fix |
|-------|-----------|
| Touch target too small | Add padding: `(44 - size) / 2` |
| Text contrast too low | Use `neutral.900` on light, `neutral.50` on dark |
| Font size off-scale | Round to nearest: 12, 14, 16, 18, 20, 24, 30, 36 |
| Spacing off-grid | Round to nearest: 4, 8, 12, 16, 24, 32, 48, 64 |
| Too many nav items | Combine or use dropdowns (max 7 items) |
| Missing label | Add `aria-label` or `<label for="...">` |
| No focus indicator | Add `outline: 2px solid [color]` |
| Deep nesting | Flatten structure (max 4 levels) |
| Too many shadows | Consolidate to 3-4 levels |
| Missing loading state | Add spinner/skeleton screen |

## Validation Commands

```bash
# Run all validations for current phase
prism audit ./project

# Run specific validator
prism validate ./project --hierarchy
prism validate ./project --touch-targets
prism validate ./project --accessibility
prism validate ./project --contrast

# Get JSON output
prism audit ./project --json

# Get improvement suggestions
prism suggest ./project

# Compare validation scores
prism compare ./project --from v1 --to v2 --show-validation
```

## Design Token Quick Reference

**Typography**: 12, 14, 16, 18, 20, 24, 30, 36  
**Spacing**: 0, 4, 8, 12, 16, 24, 32, 48, 64, 96, 128  
**Touch Targets**: 44x44px minimum  
**Contrast**: 4.5:1 (text), 3:1 (large text, UI elements)  
**Nesting**: 4 levels maximum  
**Navigation**: 7 items maximum  
**Form Fields**: 5 per section maximum  
**Shadows**: 3-4 levels maximum

---

**Last Updated**: 2025-10-25  
**PRISM Version**: 0.2.0+  
**WCAG Version**: 2.1 Level AA
