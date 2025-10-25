package validate

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/johanbellander/prism/internal/types"
)

// ContrastRule defines validation rules for color contrast (WCAG)
type ContrastRule struct {
	RequireWCAG_AA    bool    // WCAG AA compliance (4.5:1 for normal, 3:1 for large)
	RequireWCAG_AAA   bool    // WCAG AAA compliance (7:1 for normal, 4.5:1 for large)
	NormalTextRatio   float64 // 4.5:1 for AA, 7:1 for AAA
	LargeTextRatio    float64 // 3:1 for AA, 4.5:1 for AAA
	LargeTextSizePx   int     // 18px bold or 24px normal
}

// DefaultContrastRule returns the default WCAG AA contrast validation rules
func DefaultContrastRule() ContrastRule {
	return ContrastRule{
		RequireWCAG_AA:  true,
		RequireWCAG_AAA: false,
		NormalTextRatio: 4.5,
		LargeTextRatio:  3.0,
		LargeTextSizePx: 18,
	}
}

// ContrastIssue represents a single contrast validation issue
type ContrastIssue struct {
	Severity       string  // "error", "warning", "info"
	Category       string  // e.g., "contrast_fail", "contrast_aaa"
	Message        string
	ComponentID    string  // Component ID if applicable
	ForegroundColor string  // Hex color
	BackgroundColor string  // Hex color
	ContrastRatio   float64 // Calculated ratio
	RequiredRatio   float64 // Required ratio for compliance
}

// ContrastResult represents the result of contrast validation
type ContrastResult struct {
	Passed bool
	Issues []ContrastIssue
}

// ValidateContrast validates WCAG color contrast ratios
func ValidateContrast(structure *types.Structure, rule ContrastRule) ContrastResult {
	result := ContrastResult{
		Passed: true,
		Issues: []ContrastIssue{},
	}

	// Analyze all components for text/background color combinations
	var analyzeComponent func(comp *types.Component, parentBg string, depth int)
	analyzeComponent = func(comp *types.Component, parentBg string, depth int) {
		// Determine the effective background color for this component
		effectiveBg := parentBg
		if comp.Layout.Background != "" {
			effectiveBg = comp.Layout.Background
		}

		// Check if this component has text with a color
		if comp.Type == "text" && comp.Color != "" && effectiveBg != "" {
			// Calculate contrast ratio
			ratio := calculateContrastRatio(comp.Color, effectiveBg)
			
			// Determine if this is large text
			isLargeText := isLargeTextSize(comp.Size, comp.Weight)
			
			// Determine required ratio
			requiredRatio := rule.NormalTextRatio
			if isLargeText {
				requiredRatio = rule.LargeTextRatio
			}
			
			// Check compliance
			if ratio < requiredRatio {
				result.Issues = append(result.Issues, ContrastIssue{
					Severity:        "error",
					Category:        "contrast_fail",
					Message:         fmt.Sprintf("Contrast: '%s' (%s) on %s fails WCAG AA (%.1f:1, requires %.1f:1)", comp.ID, comp.Color, effectiveBg, ratio, requiredRatio),
					ComponentID:     comp.ID,
					ForegroundColor: comp.Color,
					BackgroundColor: effectiveBg,
					ContrastRatio:   ratio,
					RequiredRatio:   requiredRatio,
				})
				result.Passed = false
				
				// Provide suggestion
				suggestion := suggestCompliantColor(comp.Color, effectiveBg, requiredRatio)
				if suggestion != "" {
					result.Issues = append(result.Issues, ContrastIssue{
						Severity:        "info",
						Category:        "contrast_suggestion",
						Message:         fmt.Sprintf("   Suggestion: Use %s or similar for compliance", suggestion),
						ComponentID:     comp.ID,
						ForegroundColor: suggestion,
						BackgroundColor: effectiveBg,
					})
				}
			} else if rule.RequireWCAG_AAA {
				// Check AAA compliance
				aaaRatio := 7.0
				if isLargeText {
					aaaRatio = 4.5
				}
				
				if ratio < aaaRatio {
					result.Issues = append(result.Issues, ContrastIssue{
						Severity:        "warning",
						Category:        "contrast_aaa",
						Message:         fmt.Sprintf("Contrast: '%s' passes AA but fails AAA (%.1f:1, requires %.1f:1 for AAA)", comp.ID, ratio, aaaRatio),
						ComponentID:     comp.ID,
						ForegroundColor: comp.Color,
						BackgroundColor: effectiveBg,
						ContrastRatio:   ratio,
						RequiredRatio:   aaaRatio,
					})
				}
			}
		}

		// Check button text contrast
		if comp.Type == "button" && comp.Content != "" {
			// Buttons typically have white text on colored background
			textColor := "#FFFFFF" // Default button text color
			buttonBg := effectiveBg
			if comp.Layout.Background != "" {
				buttonBg = comp.Layout.Background
			}
			
			if buttonBg != "" {
				ratio := calculateContrastRatio(textColor, buttonBg)
				requiredRatio := rule.NormalTextRatio
				
				if ratio < requiredRatio {
					result.Issues = append(result.Issues, ContrastIssue{
						Severity:        "error",
						Category:        "contrast_fail",
						Message:         fmt.Sprintf("Contrast: Button '%s' text (%s) on %s fails WCAG AA (%.1f:1, requires %.1f:1)", comp.ID, textColor, buttonBg, ratio, requiredRatio),
						ComponentID:     comp.ID,
						ForegroundColor: textColor,
						BackgroundColor: buttonBg,
						ContrastRatio:   ratio,
						RequiredRatio:   requiredRatio,
					})
					result.Passed = false
				}
			}
		}

		// Recurse into children
		for i := range comp.Children {
			analyzeComponent(&comp.Children[i], effectiveBg, depth+1)
		}
	}

	// Default background is white for Phase 1
	defaultBg := "#FFFFFF"
	
	// Analyze all top-level components
	for i := range structure.Components {
		analyzeComponent(&structure.Components[i], defaultBg, 0)
	}

	return result
}

// calculateContrastRatio calculates the WCAG contrast ratio between two colors
func calculateContrastRatio(fg, bg string) float64 {
	fgLum := relativeLuminance(fg)
	bgLum := relativeLuminance(bg)
	
	lighter := math.Max(fgLum, bgLum)
	darker := math.Min(fgLum, bgLum)
	
	return (lighter + 0.05) / (darker + 0.05)
}

// relativeLuminance calculates the relative luminance of a color
// Formula from WCAG 2.0: https://www.w3.org/TR/WCAG20/#relativeluminancedef
func relativeLuminance(hexColor string) float64 {
	r, g, b := hexToRGB(hexColor)
	
	// Convert to 0-1 range
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0
	
	// Apply gamma correction
	rf = gammaCorrect(rf)
	gf = gammaCorrect(gf)
	bf = gammaCorrect(bf)
	
	// Calculate luminance
	return 0.2126*rf + 0.7152*gf + 0.0722*bf
}

// gammaCorrect applies gamma correction to a color channel
func gammaCorrect(channel float64) float64 {
	if channel <= 0.03928 {
		return channel / 12.92
	}
	return math.Pow((channel+0.055)/1.055, 2.4)
}

// hexToRGB converts a hex color string to RGB values
func hexToRGB(hexColor string) (r, g, b int) {
	// Remove # if present
	hex := strings.TrimPrefix(hexColor, "#")
	
	// Parse hex values
	if len(hex) == 6 {
		val, _ := strconv.ParseInt(hex, 16, 64)
		r = int((val >> 16) & 0xFF)
		g = int((val >> 8) & 0xFF)
		b = int(val & 0xFF)
	} else if len(hex) == 3 {
		// Handle shorthand hex (#RGB)
		rh, _ := strconv.ParseInt(string(hex[0]), 16, 64)
		gh, _ := strconv.ParseInt(string(hex[1]), 16, 64)
		bh, _ := strconv.ParseInt(string(hex[2]), 16, 64)
		r = int(rh*17) // Convert F to FF
		g = int(gh*17)
		b = int(bh*17)
	}
	
	return r, g, b
}

// isLargeTextSize determines if text is considered "large" for WCAG purposes
// Large text is 18px bold or 24px normal
func isLargeTextSize(size, weight string) bool {
	// Map size names to approximate pixel values
	sizeMap := map[string]int{
		"xs":   12,
		"sm":   14,
		"base": 16,
		"lg":   20,
		"xl":   24,
		"2xl":  30,
		"3xl":  36,
		"4xl":  48,
	}
	
	sizePx := sizeMap[size]
	
	// 18px bold or 24px normal is considered large
	if weight == "bold" && sizePx >= 18 {
		return true
	}
	if sizePx >= 24 {
		return true
	}
	
	return false
}

// suggestCompliantColor suggests a darker/lighter version of the color for compliance
func suggestCompliantColor(fg, bg string, requiredRatio float64) string {
	// Simple approach: darken or lighten the foreground color
	r, g, b := hexToRGB(fg)
	
	// Check if we should darken or lighten
	bgLum := relativeLuminance(bg)
	
	// Try darkening
	for i := 0; i < 10; i++ {
		factor := 1.0 - float64(i)*0.1
		newR := int(float64(r) * factor)
		newG := int(float64(g) * factor)
		newB := int(float64(b) * factor)
		
		newHex := rgbToHex(newR, newG, newB)
		ratio := calculateContrastRatio(newHex, bg)
		
		if ratio >= requiredRatio {
			return newHex
		}
	}
	
	// If darkening doesn't work, try lightening
	if bgLum < 0.5 {
		for i := 1; i <= 10; i++ {
			factor := 1.0 + float64(i)*0.1
			newR := int(math.Min(255, float64(r)*factor))
			newG := int(math.Min(255, float64(g)*factor))
			newB := int(math.Min(255, float64(b)*factor))
			
			newHex := rgbToHex(newR, newG, newB)
			ratio := calculateContrastRatio(newHex, bg)
			
			if ratio >= requiredRatio {
				return newHex
			}
		}
	}
	
	return ""
}

// rgbToHex converts RGB values to hex color string
func rgbToHex(r, g, b int) string {
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}
