# PRISM Design Tokens Guide

Design tokens are the fundamental design decisions expressed as data. They replace magic numbers and arbitrary values with a systematic, maintainable design language.

## Table of Contents

- [What Are Design Tokens?](#what-are-design-tokens)
- [Typography Scales](#typography-scales)
- [Color Systems](#color-systems)
- [Spacing (8pt Grid)](#spacing-8pt-grid)
- [Elevation & Shadows](#elevation--shadows)
- [Responsive Breakpoints](#responsive-breakpoints)
- [Dark Mode Tokens](#dark-mode-tokens)
- [Component States](#component-states)
- [Migration Guide](#migration-guide)
- [Best Practices](#best-practices)

---

## What Are Design Tokens?

Design tokens are **named entities that store visual design attributes**. They provide a single source of truth for design decisions and enable consistency across your entire UI.

### Why Use Design Tokens?

**Without tokens** (arbitrary values):
```json
{
  "button": {
    "padding": "13px 27px",
    "font_size": 15,
    "color": "#3498db",
    "border_radius": 6
  }
}
```
❌ Hard to maintain  
❌ Inconsistent across components  
❌ No semantic meaning  
❌ Difficult to theme  

**With tokens**:
```json
{
  "button": {
    "padding": "spacing.3 spacing.6",    // 12px 24px
    "font_size": "typography.base",      // 16px
    "color": "color.primary.600",
    "border_radius": "border.radius.md"  // 8px
  }
}
```
✅ Maintainable (change once, update everywhere)  
✅ Consistent (same values across all components)  
✅ Semantic (meaningful names)  
✅ Themable (swap token values for different themes)  

### Token Categories in PRISM

| Category | Purpose | Example Tokens |
|----------|---------|----------------|
| **Typography** | Font sizes, weights, line heights | `text.base`, `text.2xl`, `weight.bold` |
| **Color** | All color values | `primary.600`, `neutral.50`, `semantic.error` |
| **Spacing** | Margins, padding, gaps | `spacing.2`, `spacing.4`, `spacing.8` |
| **Elevation** | Shadows and layering | `shadow.sm`, `shadow.lg` |
| **Border** | Border radius, widths | `radius.md`, `border.width.2` |
| **Breakpoint** | Responsive viewports | `breakpoint.mobile`, `breakpoint.desktop` |

---

## Typography Scales

Typography tokens ensure consistent, harmonious text sizing across your UI.

### Predefined Scale Ratios

PRISM supports multiple type scale ratios based on music theory and the golden ratio:

| Ratio | Name | Character | Best For |
|-------|------|-----------|----------|
| **1.125** | Major Second | Tight, compact | Dense UIs, data tables |
| **1.200** | Minor Third | Balanced | General purpose (default) |
| **1.250** | Major Third | Relaxed | Marketing sites, blogs |
| **1.333** | Perfect Fourth | Spacious | Editorial content |
| **1.414** | Augmented Fourth | Dramatic | Landing pages |
| **1.618** | Golden Ratio | Very dramatic | Hero sections, presentations |

### Default Type Scale (1.200 ratio)

```json
{
  "typography": {
    "sizes": {
      "xs": 12,      // 16 / 1.2 / 1.2
      "sm": 14,      // 16 / 1.2
      "base": 16,    // Base size
      "lg": 18,      // 16 * 1.125 (adjusted for even number)
      "xl": 20,      // 16 * 1.2
      "2xl": 24,     // 16 * 1.2 * 1.2
      "3xl": 30,     // 16 * 1.2^3 (rounded)
      "4xl": 36      // 16 * 1.2^4 (rounded)
    }
  }
}
```

**Usage in JSON**:
```json
{
  "components": [
    {"type": "text", "size": "4xl", "content": "Page Title"},
    {"type": "text", "size": "2xl", "content": "Section Heading"},
    {"type": "text", "size": "base", "content": "Body text"}
  ]
}
```

### Custom Type Scale

To create a custom scale with a different ratio:

```json
{
  "typography": {
    "scale_ratio": 1.250,
    "base_size": 16,
    "sizes": {
      "xs": 10,      // 16 / 1.25 / 1.25 (rounded)
      "sm": 13,      // 16 / 1.25 (rounded)
      "base": 16,    // Base
      "lg": 20,      // 16 * 1.25
      "xl": 25,      // 16 * 1.25^2
      "2xl": 31,     // 16 * 1.25^3 (rounded)
      "3xl": 39,     // 16 * 1.25^4 (rounded)
      "4xl": 49      // 16 * 1.25^5 (rounded)
    }
  }
}
```

### Line Height Tokens

Line height ensures readable text with proper vertical rhythm:

```json
{
  "typography": {
    "line_height": {
      "tight": 1.25,      // Dense layouts, headings
      "normal": 1.5,      // Body text (default)
      "relaxed": 1.75,    // Long-form content
      "loose": 2.0        // Poetry, code blocks
    }
  }
}
```

**Usage**:
```json
{
  "type": "text",
  "size": "base",
  "line_height": "normal",    // 1.5 → 24px (16px * 1.5)
  "content": "Readable paragraph text with proper vertical spacing."
}
```

### Letter Spacing Tokens

Letter spacing (tracking) adjusts horizontal spacing between characters:

```json
{
  "typography": {
    "letter_spacing": {
      "tighter": "-0.05em",   // Tight headings
      "tight": "-0.025em",    // Large headings
      "normal": "0",          // Body text (default)
      "wide": "0.025em",      // Uppercase text
      "wider": "0.05em",      // Labels, buttons
      "widest": "0.1em"       // All-caps headings
    }
  }
}
```

**Usage**:
```json
{
  "type": "text",
  "size": "xs",
  "letter_spacing": "wider",
  "transform": "uppercase",
  "content": "LABEL"
}
```

### Font Weight Tokens

```json
{
  "typography": {
    "weights": {
      "light": 300,
      "normal": 400,
      "medium": 500,
      "semibold": 600,
      "bold": 700,
      "extrabold": 800,
      "black": 900
    }
  }
}
```

**Usage**:
```json
{
  "type": "text",
  "size": "2xl",
  "weight": "bold",
  "content": "Important Heading"
}
```

### Font Family Tokens

```json
{
  "typography": {
    "families": {
      "sans": "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
      "serif": "Georgia, Cambria, 'Times New Roman', Times, serif",
      "mono": "Menlo, Monaco, Consolas, 'Courier New', monospace"
    }
  }
}
```

**When to use**:
- **Sans**: UI elements, body text, headings (default)
- **Serif**: Editorial content, long-form reading
- **Mono**: Code blocks, terminal output, tabular data

---

## Color Systems

Color tokens create consistent, accessible, and themable color schemes.

### Neutral Color Palette

Grayscale colors for text, backgrounds, and borders:

```json
{
  "color": {
    "neutral": {
      "50": "#FAFAFA",    // Lightest - backgrounds
      "100": "#F5F5F5",   // Very light - surfaces
      "200": "#E5E5E5",   // Light - borders
      "300": "#D4D4D4",   // Medium light
      "400": "#A3A3A3",   // Medium - disabled text
      "500": "#737373",   // Medium dark - secondary text
      "600": "#525252",   // Dark - tertiary text
      "700": "#404040",   // Very dark
      "800": "#262626",   // Almost black
      "900": "#171717"    // Darkest - primary text
    }
  }
}
```

**Usage guidelines**:
- **50-100**: Light backgrounds, surfaces
- **200-300**: Borders, dividers, disabled backgrounds
- **400-500**: Secondary text, placeholders
- **600-700**: Tertiary text, icons
- **800-900**: Primary text, headings

### Primary Color Palette

Brand colors for interactive elements:

```json
{
  "color": {
    "primary": {
      "50": "#EFF6FF",    // Lightest - backgrounds, badges
      "100": "#DBEAFE",   // Very light - hover backgrounds
      "200": "#BFDBFE",   // Light - borders
      "300": "#93C5FD",   // Medium light - disabled states
      "400": "#60A5FA",   // Medium - hover states
      "500": "#3B82F6",   // Base - default interactive
      "600": "#2563EB",   // Dark - primary buttons
      "700": "#1D4ED8",   // Very dark - active states
      "800": "#1E40AF",   // Almost black - pressed states
      "900": "#1E3A8A"    // Darkest - text on light backgrounds
    }
  }
}
```

**Usage guidelines**:
- **50-200**: Subtle backgrounds, light accents
- **300-400**: Hover states, secondary buttons
- **500-600**: Primary actions, links, active states
- **700-900**: Dark mode, text, pressed states

### Semantic Color Tokens

Communicate meaning through color:

```json
{
  "color": {
    "semantic": {
      "success": {
        "light": "#D1FAE5",   // Light background
        "base": "#10B981",    // Icon, border
        "dark": "#065F46"     // Text
      },
      "warning": {
        "light": "#FEF3C7",
        "base": "#F59E0B",
        "dark": "#92400E"
      },
      "error": {
        "light": "#FEE2E2",
        "base": "#EF4444",
        "dark": "#991B1B"
      },
      "info": {
        "light": "#DBEAFE",
        "base": "#3B82F6",
        "dark": "#1E40AF"
      }
    }
  }
}
```

**Usage**:
```json
{
  "type": "alert",
  "variant": "error",
  "style": {
    "background": "semantic.error.light",
    "border_color": "semantic.error.base",
    "text_color": "semantic.error.dark"
  }
}
```

### WCAG Contrast Requirements

All color combinations must meet accessibility standards:

| Content Type | WCAG AA | WCAG AAA |
|--------------|---------|----------|
| **Normal text** (< 24px) | 4.5:1 | 7:1 |
| **Large text** (≥ 24px) | 3:1 | 4.5:1 |
| **UI components** | 3:1 | - |

**Pre-validated combinations**:

```json
{
  "color_combinations": {
    "light_mode": {
      "text_on_white": {
        "primary": "neutral.900",     // 16.1:1 ✓ AAA
        "secondary": "neutral.600",   // 7.0:1 ✓ AAA
        "tertiary": "neutral.500"     // 4.6:1 ✓ AA
      },
      "text_on_primary": {
        "color": "white",             // 8.6:1 ✓ AAA
        "on": "primary.600"
      }
    },
    "dark_mode": {
      "text_on_black": {
        "primary": "neutral.50",      // 17.8:1 ✓ AAA
        "secondary": "neutral.400",   // 5.8:1 ✓ AAA
        "tertiary": "neutral.500"     // 4.6:1 ✓ AA
      }
    }
  }
}
```

### Color Token Naming Convention

Use semantic, descriptive names:

✅ **Good**:
```json
{
  "background_primary": "neutral.50",
  "text_primary": "neutral.900",
  "border_default": "neutral.200",
  "button_primary_bg": "primary.600"
}
```

❌ **Bad**:
```json
{
  "bg1": "#FAFAFA",
  "text1": "#171717",
  "blue": "#2563EB"
}
```

---

## Spacing (8pt Grid)

The 8pt grid system creates consistent, scalable spacing throughout your UI.

### Why 8pt Grid?

- **Scalability**: 8 is divisible by 2, 4, making it easy to scale
- **Screen density**: Works well across 1x, 2x, 3x displays
- **Consistency**: Fewer arbitrary decisions
- **Rhythm**: Creates visual harmony

### Spacing Scale

```json
{
  "spacing": {
    "0": 0,       // No spacing
    "1": 4,       // 0.5 * base (half-step)
    "2": 8,       // 1 * base
    "3": 12,      // 1.5 * base
    "4": 16,      // 2 * base
    "5": 20,      // 2.5 * base
    "6": 24,      // 3 * base
    "8": 32,      // 4 * base
    "10": 40,     // 5 * base
    "12": 48,     // 6 * base
    "16": 64,     // 8 * base
    "20": 80,     // 10 * base
    "24": 96,     // 12 * base
    "32": 128,    // 16 * base
    "40": 160,    // 20 * base
    "48": 192     // 24 * base
  }
}
```

**Token naming**: Number represents the multiple of base unit (4px)
- `spacing.2` = 8px (2 * 4px)
- `spacing.4` = 16px (4 * 4px)
- `spacing.8` = 32px (8 * 4px)

### When to Use Each Spacing Value

| Value | Pixels | Use Cases |
|-------|--------|-----------|
| `0` | 0px | Remove spacing, flush layouts |
| `1` | 4px | Icon gaps, tight spacing, optical adjustments |
| `2` | 8px | Minimum spacing, compact layouts |
| `3` | 12px | Form element spacing, list items |
| `4` | 16px | Default padding, comfortable spacing |
| `6` | 24px | Section spacing, card padding |
| `8` | 32px | Large padding, component separation |
| `12` | 48px | Major section gaps |
| `16` | 64px | Page sections, hero spacing |
| `24+` | 96px+ | Large marketing layouts |

### Usage Examples

**Button padding**:
```json
{
  "type": "button",
  "layout": {
    "padding_y": "spacing.3",   // 12px vertical
    "padding_x": "spacing.6"    // 24px horizontal
  }
}
```

**Card layout**:
```json
{
  "type": "card",
  "layout": {
    "padding": "spacing.6",     // 24px all sides
    "gap": "spacing.4"          // 16px between children
  }
}
```

**Form spacing**:
```json
{
  "type": "form",
  "layout": {
    "gap": "spacing.4"          // 16px between fields
  },
  "sections": {
    "gap": "spacing.8"          // 32px between sections
  }
}
```

### Half-Step (4px) Guidelines

Use `spacing.1` (4px) only for:
- ✅ Optical alignment corrections
- ✅ Icon-to-text spacing in buttons
- ✅ Fine-tuning visual balance
- ❌ **NOT** for general layout spacing

**Example** (icon button):
```json
{
  "type": "button",
  "layout": {
    "padding": "spacing.3",     // 12px
    "gap": "spacing.1"          // 4px between icon and text
  },
  "children": [
    {"type": "icon", "name": "check"},
    {"type": "text", "content": "Submit"}
  ]
}
```

### Margin vs Padding Convention

- **Padding**: Internal spacing (within component)
- **Margin**: External spacing (between components)

```json
{
  "card": {
    "padding": "spacing.6",        // Internal: 24px
    "margin_bottom": "spacing.4"   // External: 16px to next card
  }
}
```

---

## Elevation & Shadows

Elevation creates depth hierarchy and indicates interactivity through shadows.

### Elevation Levels (0-5)

```json
{
  "elevation": {
    "0": "none",
    "1": "0 1px 2px 0 rgb(0 0 0 / 0.05)",
    "2": "0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)",
    "3": "0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)",
    "4": "0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1)",
    "5": "0 25px 50px -12px rgb(0 0 0 / 0.25)"
  }
}
```

**Token aliases** (semantic names):
```json
{
  "shadow": {
    "sm": "elevation.1",       // Subtle depth
    "md": "elevation.2",       // Default raised
    "lg": "elevation.3",       // Floating elements
    "xl": "elevation.4",       // Modals, overlays
    "2xl": "elevation.5"       // Maximum depth
  }
}
```

### When to Use Each Elevation

| Level | Shadow | Use Cases |
|-------|--------|-----------|
| **0** | None | Flat UI, backgrounds, contained elements |
| **1** | `sm` | Cards (default), panels, tiles |
| **2** | `md` | Raised buttons (hover), dropdowns, tooltips |
| **3** | `lg` | Floating action buttons, popovers, date pickers |
| **4** | `xl` | Modals, dialogs, drawer panels |
| **5** | `2xl` | Maximum emphasis (use sparingly) |

### Component Elevation Guidelines

```json
{
  "components": {
    "card": {
      "default": "shadow.sm",      // Resting state
      "hover": "shadow.md"         // Elevated on hover
    },
    "button": {
      "default": "elevation.0",    // Flat by default
      "hover": "shadow.sm",        // Subtle lift
      "active": "elevation.0"      // Pressed flat
    },
    "dropdown": {
      "default": "shadow.md"       // Clearly above content
    },
    "modal": {
      "default": "shadow.xl"       // Top layer
    },
    "tooltip": {
      "default": "shadow.md"       // Floating
    },
    "fab": {
      "default": "shadow.lg",      // Always elevated
      "hover": "shadow.xl"         // Extra lift on hover
    }
  }
}
```

### Interactive State Elevations

Elevation should respond to user interaction:

```json
{
  "type": "card",
  "states": {
    "default": {
      "shadow": "elevation.1"
    },
    "hover": {
      "shadow": "elevation.2",
      "transition": "box-shadow 200ms ease"
    },
    "active": {
      "shadow": "elevation.1"      // Lower on press
    }
  }
}
```

### Dark Mode Shadows

Shadows are less visible on dark backgrounds. Adjust for dark mode:

```json
{
  "shadow": {
    "light_mode": {
      "md": "0 4px 6px -1px rgb(0 0 0 / 0.1)"
    },
    "dark_mode": {
      "md": "0 4px 6px -1px rgb(0 0 0 / 0.3)"   // Higher opacity
    }
  }
}
```

---

## Responsive Breakpoints

Breakpoint tokens define where layouts adapt to different screen sizes.

### Standard Breakpoints

```json
{
  "breakpoint": {
    "mobile": 640,      // Phones
    "tablet": 768,      // Tablets
    "desktop": 1024,    // Laptops
    "wide": 1280,       // Desktops
    "ultrawide": 1536   // Large monitors
  }
}
```

**Mobile-first approach**:
```json
{
  "responsive": {
    "mobile": {
      "breakpoint": 640,
      "layout": {
        "direction": "vertical",
        "padding": "spacing.4"
      }
    },
    "tablet": {
      "breakpoint": 768,
      "layout": {
        "direction": "horizontal",
        "padding": "spacing.6"
      }
    },
    "desktop": {
      "breakpoint": 1024,
      "layout": {
        "max_width": 1200,
        "padding": "spacing.8"
      }
    }
  }
}
```

### Custom Breakpoints

For specific component needs:

```json
{
  "breakpoint": {
    "xs": 375,          // Small phones
    "sm": 640,          // Standard phones
    "md": 768,          // Tablets
    "lg": 1024,         // Laptops
    "xl": 1280,         // Desktops
    "2xl": 1536         // Large screens
  }
}
```

### Viewport-Specific Changes

```json
{
  "component": "navigation",
  "responsive": {
    "mobile": {
      "max_width": 640,
      "changes": {
        "layout.direction": "vertical",
        "menu.display": "none",
        "hamburger.display": "block"
      }
    },
    "desktop": {
      "min_width": 1024,
      "changes": {
        "layout.direction": "horizontal",
        "menu.display": "flex",
        "hamburger.display": "none"
      }
    }
  }
}
```

---

## Dark Mode Tokens

Dark mode requires carefully crafted color tokens that maintain contrast and usability.

### Semantic Color Structure

Instead of inverting colors, define semantic tokens:

```json
{
  "theme": {
    "light": {
      "background": {
        "primary": "neutral.50",
        "secondary": "white",
        "tertiary": "neutral.100"
      },
      "text": {
        "primary": "neutral.900",
        "secondary": "neutral.600",
        "tertiary": "neutral.500",
        "disabled": "neutral.400"
      },
      "border": {
        "default": "neutral.200",
        "hover": "neutral.300",
        "focus": "primary.500"
      }
    },
    "dark": {
      "background": {
        "primary": "#0A0A0A",      // Near black (not pure black)
        "secondary": "neutral.900",
        "tertiary": "neutral.800"
      },
      "text": {
        "primary": "neutral.50",
        "secondary": "neutral.400",
        "tertiary": "neutral.500",
        "disabled": "neutral.600"
      },
      "border": {
        "default": "neutral.700",
        "hover": "neutral.600",
        "focus": "primary.400"
      }
    }
  }
}
```

### Why Not Pure Black?

Avoid `#000000` for dark mode backgrounds:
- ❌ Pure black creates harsh contrast
- ❌ OLED smearing on pure black
- ❌ Harder to show elevation with shadows

✅ Use `#0A0A0A` or `neutral.900` (#171717) instead

### Surface Tokens

Surfaces are elevated above the background:

```json
{
  "surface": {
    "light": {
      "level_0": "white",          // Base surface
      "level_1": "neutral.50",     // Subtle elevation
      "level_2": "neutral.100"     // Higher elevation
    },
    "dark": {
      "level_0": "neutral.900",    // Base surface
      "level_1": "neutral.800",    // Subtle elevation
      "level_2": "neutral.700"     // Higher elevation
    }
  }
}
```

**Usage**:
```json
{
  "type": "card",
  "style": {
    "background": "surface.level_1",
    "border": "1px solid",
    "border_color": "border.default"
  }
}
```

### Contrast Validation in Both Modes

Ensure WCAG compliance in light AND dark mode:

```json
{
  "contrast_validation": {
    "light_mode": {
      "text_on_background": {
        "foreground": "text.primary",      // neutral.900
        "background": "background.primary", // neutral.50
        "ratio": 16.1,                     // ✓ AAA
        "passes": ["AA", "AAA"]
      }
    },
    "dark_mode": {
      "text_on_background": {
        "foreground": "text.primary",      // neutral.50
        "background": "background.primary", // #0A0A0A
        "ratio": 17.8,                     // ✓ AAA
        "passes": ["AA", "AAA"]
      }
    }
  }
}
```

### Brand Colors in Dark Mode

Primary colors often need adjustment for dark backgrounds:

```json
{
  "color": {
    "primary": {
      "light_mode": "primary.600",   // Darker for light bg
      "dark_mode": "primary.400"     // Lighter for dark bg
    },
    "success": {
      "light_mode": "semantic.success.dark",
      "dark_mode": "semantic.success.base"
    }
  }
}
```

---

## Component States

Define consistent states for all interactive components.

### Standard Interactive States

```json
{
  "states": {
    "default": {
      "description": "Resting state, no interaction"
    },
    "hover": {
      "description": "Mouse over the element",
      "triggers": ["mouse_enter"]
    },
    "active": {
      "description": "Element is being pressed",
      "triggers": ["mouse_down", "touch_start"]
    },
    "focus": {
      "description": "Element has keyboard focus",
      "triggers": ["tab", "click"]
    },
    "disabled": {
      "description": "Element is not interactive",
      "style": {
        "opacity": 0.5,
        "cursor": "not-allowed"
      }
    }
  }
}
```

### Button State Tokens

```json
{
  "button": {
    "primary": {
      "default": {
        "background": "primary.600",
        "color": "white",
        "border": "none"
      },
      "hover": {
        "background": "primary.700",
        "shadow": "elevation.1"
      },
      "active": {
        "background": "primary.800",
        "shadow": "none"
      },
      "focus": {
        "background": "primary.600",
        "outline": "2px solid",
        "outline_color": "primary.400",
        "outline_offset": "2px"
      },
      "disabled": {
        "background": "neutral.300",
        "color": "neutral.500",
        "cursor": "not-allowed"
      }
    },
    "secondary": {
      "default": {
        "background": "transparent",
        "color": "primary.600",
        "border": "1px solid",
        "border_color": "primary.600"
      },
      "hover": {
        "background": "primary.50",
        "border_color": "primary.700"
      },
      "active": {
        "background": "primary.100",
        "border_color": "primary.800"
      },
      "focus": {
        "outline": "2px solid",
        "outline_color": "primary.400",
        "outline_offset": "2px"
      },
      "disabled": {
        "border_color": "neutral.300",
        "color": "neutral.400"
      }
    }
  }
}
```

### Input Field State Tokens

```json
{
  "input": {
    "default": {
      "background": "white",
      "border": "1px solid",
      "border_color": "neutral.300",
      "text_color": "text.primary"
    },
    "hover": {
      "border_color": "neutral.400"
    },
    "focus": {
      "border_color": "primary.500",
      "outline": "2px solid",
      "outline_color": "primary.200",
      "outline_offset": "0"
    },
    "error": {
      "border_color": "semantic.error.base",
      "outline_color": "semantic.error.light"
    },
    "disabled": {
      "background": "neutral.100",
      "border_color": "neutral.200",
      "text_color": "neutral.400",
      "cursor": "not-allowed"
    }
  }
}
```

### Loading State Tokens

```json
{
  "loading": {
    "spinner": {
      "color": "primary.600",
      "size": "spacing.6",           // 24px
      "stroke_width": "2px"
    },
    "skeleton": {
      "background": "neutral.200",
      "highlight": "neutral.100",
      "animation": "pulse 2s ease-in-out infinite"
    },
    "overlay": {
      "background": "rgb(255 255 255 / 0.8)",
      "backdrop_blur": "4px"
    }
  }
}
```

**Usage**:
```json
{
  "type": "button",
  "states": {
    "loading": {
      "disabled": true,
      "content": "",
      "children": [
        {
          "type": "spinner",
          "color": "loading.spinner.color",
          "size": "loading.spinner.size"
        },
        {
          "type": "text",
          "content": "Loading...",
          "color": "text.secondary"
        }
      ]
    }
  }
}
```

### Empty State Tokens

```json
{
  "empty_state": {
    "icon": {
      "size": "spacing.16",          // 64px
      "color": "neutral.400"
    },
    "title": {
      "size": "typography.xl",
      "color": "text.primary",
      "weight": "weights.semibold"
    },
    "description": {
      "size": "typography.base",
      "color": "text.secondary"
    }
  }
}
```

### Error State Tokens

```json
{
  "error_state": {
    "message": {
      "background": "semantic.error.light",
      "border_color": "semantic.error.base",
      "text_color": "semantic.error.dark",
      "icon_color": "semantic.error.base"
    },
    "input_error": {
      "border_color": "semantic.error.base",
      "text_color": "semantic.error.dark"
    }
  }
}
```

---

## Migration Guide

Transitioning from hardcoded values to design tokens.

### Step 1: Audit Current Values

Find all unique values in your design:

```bash
# Find all color values
grep -r "#[0-9A-Fa-f]\{6\}" .

# Find all spacing values
grep -r "padding\|margin" . | grep -o "[0-9]\+px"

# Find all font sizes
grep -r "font-size\|size" . | grep -o "[0-9]\+px"
```

### Step 2: Map to Nearest Token

Convert hardcoded values to tokens:

| Category | Before | After |
|----------|--------|-------|
| **Color** | `#2563eb` | `color.primary.600` |
| **Spacing** | `padding: 13px` | `padding: spacing.3` (12px) |
| **Font** | `font-size: 15px` | `size: typography.base` (16px) |
| **Shadow** | `box-shadow: 0 2px 4px rgba(0,0,0,0.1)` | `shadow: elevation.1` |

### Step 3: Update Component Definitions

**Before** (hardcoded):
```json
{
  "type": "button",
  "style": {
    "padding": "13px 27px",
    "font_size": 15,
    "color": "#FFFFFF",
    "background": "#2563EB",
    "border_radius": 6,
    "box_shadow": "0 2px 4px rgba(0,0,0,0.1)"
  }
}
```

**After** (tokens):
```json
{
  "type": "button",
  "style": {
    "padding_y": "spacing.3",         // 12px
    "padding_x": "spacing.6",         // 24px
    "font_size": "typography.base",   // 16px
    "color": "white",
    "background": "primary.600",
    "border_radius": "border.radius.md",  // 8px
    "shadow": "elevation.1"
  }
}
```

### Step 4: Define Token File

Create a central token definition:

```json
// tokens.json
{
  "typography": { /* ... */ },
  "color": { /* ... */ },
  "spacing": { /* ... */ },
  "elevation": { /* ... */ },
  "border": { /* ... */ }
}
```

### Step 5: Validate Token Usage

Run PRISM validation to ensure compliance:

```bash
prism validate ./project --typography --spacing --contrast
```

---

## Best Practices

### DO ✅

1. **Use semantic names**: `background.primary` instead of `gray.50`
2. **Start with base tokens**: Build on `spacing.4`, `typography.base`
3. **Validate contrast**: Check WCAG compliance for all color pairs
4. **Document exceptions**: If you break the grid, document why
5. **Use design tools**: Leverage Figma variables, Tokens Studio
6. **Version your tokens**: Track changes like code
7. **Test in both modes**: Validate light and dark mode separately
8. **Round to scale**: 15px → 16px, 13px → 12px

### DON'T ❌

1. **Don't use arbitrary values**: Avoid `padding: 13px` or `color: #3498db`
2. **Don't invert for dark mode**: Create separate dark mode tokens
3. **Don't skip validation**: Always check contrast ratios
4. **Don't use too many tokens**: Limit to 3-4 shadows, 8-10 spacing values
5. **Don't name by appearance**: Avoid `blue`, `lightGray` (use semantic names)
6. **Don't break the 8pt grid**: Stick to multiples of 4 or 8
7. **Don't forget focus states**: Always define visible focus indicators
8. **Don't use pure black**: Use `#0A0A0A` for dark backgrounds

### Token Naming Convention

**Format**: `{category}.{subcategory}.{variant}`

✅ **Good**:
- `color.primary.600`
- `spacing.4`
- `typography.2xl`
- `shadow.lg`
- `border.radius.md`

❌ **Bad**:
- `blue`
- `space3`
- `bigText`
- `shadow2`
- `rounded`

### Performance Considerations

- **Bundle size**: Tokens add minimal overhead (< 5KB gzipped)
- **Runtime**: Token lookup is O(1) with proper structure
- **Caching**: Cache resolved token values in production
- **Tree-shaking**: Export tokens as ES modules for tree-shaking

---

## Summary

Design tokens transform your UI from a collection of arbitrary values into a systematic, maintainable design language.

**Key Takeaways**:
1. **Typography**: Use 1.200 ratio scale, 8-10 sizes maximum
2. **Color**: 9-step palettes, semantic tokens, WCAG validated
3. **Spacing**: 8pt grid (multiples of 4 or 8)
4. **Elevation**: 3-4 shadow levels
5. **Breakpoints**: Mobile-first, 640/768/1024px
6. **Dark Mode**: Separate palette, maintain contrast
7. **States**: Define all 5 states (default, hover, active, focus, disabled)

**Quick Reference**:
- Typography: `12, 14, 16, 18, 20, 24, 30, 36`
- Spacing: `0, 4, 8, 12, 16, 24, 32, 48, 64, 96, 128`
- Contrast: `4.5:1` (text), `3:1` (large text, UI)
- Touch Targets: `44x44px` minimum
- Shadows: `sm, md, lg, xl` (4 levels)

---

**Last Updated**: 2025-10-25  
**PRISM Version**: 0.2.0+  
**Related Docs**: VALIDATION_RULES.md, DESIGNPROCESS.md
