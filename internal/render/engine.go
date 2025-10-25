package render

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"github.com/johanbellander/prism/internal/types"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// RenderOptions configures the rendering process
type RenderOptions struct {
	Width       int
	Height      int
	Scale       int
	Viewport    string // "mobile", "tablet", "desktop"
	Annotations bool
	Grid        bool
}

// RenderResult contains the result of a rendering operation
type RenderResult struct {
	Image      *image.RGBA
	Width      int
	Height     int
	OutputPath string
}

// Renderer handles rendering Phase 1 structures to images
type Renderer struct {
	opts RenderOptions
}

// NewRenderer creates a new renderer with the given options
func NewRenderer(opts RenderOptions) *Renderer {
	// Set defaults
	if opts.Width == 0 {
		opts.Width = 1200
	}
	if opts.Scale == 0 {
		opts.Scale = 1
	}
	if opts.Viewport == "" {
		opts.Viewport = "desktop"
	}

	return &Renderer{opts: opts}
}

// Render renders a structure to an image
func (r *Renderer) Render(structure *types.Structure) (*RenderResult, error) {
	// Calculate canvas dimensions
	width := r.opts.Width * r.opts.Scale
	height := r.opts.Height * r.opts.Scale
	
	// If height is 0 (auto), calculate based on content
	if height == 0 {
		height = r.calculateHeight(structure) * r.opts.Scale
	}

	// Create the image
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Fill with white background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Create layout engine
	layoutEngine := NewLayoutEngine(r.opts.Scale)
	
	// Calculate layout for all components
	boxes, err := layoutEngine.CalculateLayout(structure, width, height)
	if err != nil {
		return nil, fmt.Errorf("layout calculation failed: %w", err)
	}

	// Create render context
	ctx := &renderContext{
		img:    img,
		scale:  r.opts.Scale,
		boxes:  boxes,
	}

	// Render components using calculated layout
	for _, comp := range structure.Components {
		if err := r.renderComponent(ctx, &comp); err != nil {
			return nil, fmt.Errorf("failed to render component %s: %w", comp.ID, err)
		}
	}

	return &RenderResult{
		Image:  img,
		Width:  width,
		Height: height,
	}, nil
}

// SavePNG saves the rendered result to a PNG file
func (r *RenderResult) SavePNG(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, r.Image); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}

// renderContext holds the current rendering state
type renderContext struct {
	img   *image.RGBA
	scale int
	boxes map[string]LayoutBox // calculated layout boxes for all components
}

// calculateHeight estimates the height needed for the content
func (r *Renderer) calculateHeight(structure *types.Structure) int {
	// Simple heuristic: count components and estimate height
	baseHeight := structure.Layout.Padding * 2
	componentHeight := 0

	for _, comp := range structure.Components {
		componentHeight += r.estimateComponentHeight(&comp)
		componentHeight += structure.Layout.Spacing
	}

	totalHeight := baseHeight + componentHeight
	
	// Ensure minimum height
	if totalHeight < 400 {
		totalHeight = 400
	}

	return totalHeight
}

// estimateComponentHeight estimates the height of a component
func (r *Renderer) estimateComponentHeight(comp *types.Component) int {
	baseHeight := comp.Layout.Padding * 2

	// Estimate based on component type
	switch comp.Type {
	case "text":
		return baseHeight + r.getTextHeight(comp.Size)
	case "button":
		return baseHeight + 44 // minimum touch target
	case "input":
		return baseHeight + 40
	case "image":
		return baseHeight + 200 // placeholder size
	case "box":
		childHeight := 0
		for _, child := range comp.Children {
			childHeight += r.estimateComponentHeight(&child)
			if comp.Layout.Gap > 0 {
				childHeight += comp.Layout.Gap
			}
		}
		return baseHeight + childHeight
	}

	return baseHeight + 20
}

// getTextHeight returns the height for a given text size
func (r *Renderer) getTextHeight(size string) int {
	sizes := map[string]int{
		"xs":   12,
		"sm":   14,
		"base": 16,
		"lg":   18,
		"xl":   20,
		"2xl":  24,
		"3xl":  30,
		"4xl":  36,
	}
	
	if h, ok := sizes[size]; ok {
		return h
	}
	return 16 // default
}

// renderComponent renders a single component using pre-calculated layout
func (r *Renderer) renderComponent(ctx *renderContext, comp *types.Component) error {
	// Get the calculated layout box for this component
	box, ok := ctx.boxes[comp.ID]
	if !ok {
		return fmt.Errorf("no layout box found for component %s", comp.ID)
	}

	// Render based on component type
	switch comp.Type {
	case "box":
		return r.renderBox(ctx, comp, box)
	case "text":
		return r.renderText(ctx, comp, box)
	case "button":
		return r.renderButton(ctx, comp, box)
	case "input":
		return r.renderInput(ctx, comp, box)
	case "image":
		return r.renderImage(ctx, comp, box)
	default:
		return fmt.Errorf("unsupported component type: %s", comp.Type)
	}
}

// renderBox renders a box component
func (r *Renderer) renderBox(ctx *renderContext, comp *types.Component, box LayoutBox) error {
	// Draw background if specified
	if comp.Layout.Background != "" {
		bgColor := parseColor(comp.Layout.Background)
		rect := image.Rect(box.X, box.Y, box.X+box.Width, box.Y+box.Height)
		draw.Draw(ctx.img, rect, &image.Uniform{bgColor}, image.Point{}, draw.Src)
	}

	// Draw borders if specified
	borderColor := color.RGBA{229, 229, 229, 255} // #E5E5E5
	if comp.Layout.Border != "" {
		r.drawRect(ctx.img, box.X, box.Y, box.Width, box.Height, borderColor)
	}
	if comp.Layout.BorderBottom != "" {
		r.drawHorizontalLine(ctx.img, box.X, box.Y+box.Height-1, box.Width, borderColor)
	}
	if comp.Layout.BorderRight != "" {
		r.drawVerticalLine(ctx.img, box.X+box.Width-1, box.Y, box.Height, borderColor)
	}

	// Render children using their pre-calculated layouts
	for _, child := range comp.Children {
		if err := r.renderComponent(ctx, &child); err != nil {
			return err
		}
	}

	return nil
}

// renderText renders a text component
func (r *Renderer) renderText(ctx *renderContext, comp *types.Component, box LayoutBox) error {
	if comp.Content == "" {
		return nil
	}

	textColor := parseColor(comp.Color)
	if comp.Color == "" {
		textColor = color.Black
	}

	// Split content by newlines for multi-line text
	lines := strings.Split(comp.Content, "\n")
	lineHeight := 16 // pixels between lines
	
	d := &font.Drawer{
		Dst:  ctx.img,
		Src:  image.NewUniform(textColor),
		Face: basicfont.Face7x13,
	}

	// Draw each line separately
	currentLine := 0
	for _, line := range lines {
		if line == "" {
			currentLine++ // Skip empty lines but still count for spacing
			continue
		}
		
		point := fixed.Point26_6{
			X: fixed.Int26_6(box.X * 64),
			Y: fixed.Int26_6((box.Y + 14 + (currentLine * lineHeight)) * 64),
		}
		d.Dot = point
		d.DrawString(line)
		currentLine++
	}

	return nil
}

// renderButton renders a button component
func (r *Renderer) renderButton(ctx *renderContext, comp *types.Component, box LayoutBox) error {
	// Draw button background
	bgColor := parseColor(comp.Layout.Background)
	if comp.Layout.Background == "" {
		bgColor = color.Black
	}

	rect := image.Rect(box.X, box.Y, box.X+box.Width, box.Y+box.Height)
	draw.Draw(ctx.img, rect, &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Draw button text (centered)
	if comp.Content != "" {
		textColor := parseColor(comp.Color)
		if comp.Color == "" {
			textColor = color.White
		}

		point := fixed.Point26_6{
			X: fixed.Int26_6((box.X + 10) * 64),
			Y: fixed.Int26_6((box.Y + 25) * 64),
		}

		d := &font.Drawer{
			Dst:  ctx.img,
			Src:  image.NewUniform(textColor),
			Face: basicfont.Face7x13,
			Dot:  point,
		}

		d.DrawString(comp.Content)
	}

	return nil
}

// renderInput renders an input component
func (r *Renderer) renderInput(ctx *renderContext, comp *types.Component, box LayoutBox) error {
	// Draw input border
	borderColor := color.RGBA{229, 229, 229, 255} // #E5E5E5
	r.drawRect(ctx.img, box.X, box.Y, box.Width, box.Height, borderColor)

	// Draw placeholder text if present
	if comp.Content != "" {
		textColor := color.RGBA{115, 115, 115, 255} // #737373 (gray)
		point := fixed.Point26_6{
			X: fixed.Int26_6((box.X + 8) * 64),
			Y: fixed.Int26_6((box.Y + 22) * 64),
		}

		d := &font.Drawer{
			Dst:  ctx.img,
			Src:  image.NewUniform(textColor),
			Face: basicfont.Face7x13,
			Dot:  point,
		}

		d.DrawString(comp.Content)
	}

	return nil
}

// renderImage renders an image placeholder
func (r *Renderer) renderImage(ctx *renderContext, comp *types.Component, box LayoutBox) error {
	// Draw gray rectangle as placeholder
	bgColor := color.RGBA{229, 229, 229, 255} // #E5E5E5
	rect := image.Rect(box.X, box.Y, box.X+box.Width, box.Y+box.Height)
	draw.Draw(ctx.img, rect, &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Draw "IMAGE" text in center
	textColor := color.RGBA{115, 115, 115, 255} // #737373
	point := fixed.Point26_6{
		X: fixed.Int26_6((box.X + box.Width/2 - 20) * 64),
		Y: fixed.Int26_6((box.Y + box.Height/2) * 64),
	}

	d := &font.Drawer{
		Dst:  ctx.img,
		Src:  image.NewUniform(textColor),
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	d.DrawString("IMAGE")

	return nil
}

// drawRect draws a rectangle outline
func (r *Renderer) drawRect(img *image.RGBA, x, y, width, height int, col color.Color) {
	// Top
	for i := 0; i < width; i++ {
		img.Set(x+i, y, col)
	}
	// Bottom
	for i := 0; i < width; i++ {
		img.Set(x+i, y+height-1, col)
	}
	// Left
	for i := 0; i < height; i++ {
		img.Set(x, y+i, col)
	}
	// Right
	for i := 0; i < height; i++ {
		img.Set(x+width-1, y+i, col)
	}
}

// drawHorizontalLine draws a horizontal line
func (r *Renderer) drawHorizontalLine(img *image.RGBA, x, y, width int, col color.Color) {
	for i := 0; i < width; i++ {
		img.Set(x+i, y, col)
	}
}

// drawVerticalLine draws a vertical line
func (r *Renderer) drawVerticalLine(img *image.RGBA, x, y, height int, col color.Color) {
	for i := 0; i < height; i++ {
		img.Set(x, y+i, col)
	}
}

// parseColor converts a hex color string to color.Color
func parseColor(hex string) color.Color {
	if hex == "" || hex[0] != '#' || len(hex) != 7 {
		return color.Black
	}

	var r, g, b uint8
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return color.RGBA{r, g, b, 255}
}
