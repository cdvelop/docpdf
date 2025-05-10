// Package config contains configuration structures for docpdf.
package config

import (
	"github.com/cdvelop/docpdf/style"
	"github.com/cdvelop/tinystring"
)

// Regular font style (000000)
var FontStyleRegular = FontStyle{name: "regular", intStyle: 0, size: 12, color: style.Color{R: 0, G: 0, B: 0, A: 255}}

// Italic font style (000001)
var FontStyleItalic = FontStyle{name: "italic", intStyle: 1, size: 12, color: style.Color{R: 0, G: 0, B: 0, A: 255}}

// Bold font style (000010)
var FontStyleBold = FontStyle{name: "bold", intStyle: 2, size: 12, color: style.Color{R: 0, G: 0, B: 0, A: 255}}

// Underline font style (000011)
var FontStyleUnderline = FontStyle{name: "underline", intStyle: 3, size: 12, color: style.Color{R: 0, G: 0, B: 0, A: 255}}

// FontFamily represents font files for different styles
// It contains the regular, bold, italic, and other styles.
type FontFamily struct {
	// Regular specifies the filename for the regular font style.
	// It's recommended to name this file "regular.ttf".
	Regular string
	// Bold specifies the filename for the bold font style.
	// It's recommended to name this file "bold.ttf".
	Bold string
	// Italic specifies the filename for the italic font style.
	// It's recommended to name this file "italic.ttf".
	Italic string
	// Underline specifies the filename for the underline font style.
	// It's recommended to name this file "underline.ttf".
	Underline string
	// Path specifies the base directory where the font files are located.
	// Defaults to "fonts/".
	Path string // Base path for fonts defaults to "fonts/"
}

// FontStyle represents a complete font configuration with styling, size, color, and family
// This is the unified font structure for the entire docpdf project
type FontStyle struct {
	name     string
	intStyle int
	size     float64
	color    style.Color
}

// NewFontStyle creates a new FontStyle with the given properties
// name: The font style name (e.g., "regular", "bold", "italic", "underline")
// size: Font size in points
// color: Optional color parameter. If not provided, defaults to black (0,0,0,255)
func NewFontStyle(name string, size float64, color ...style.Color) FontStyle {
	// Default color is black
	defaultColor := style.Color{R: 0, G: 0, B: 0, A: 255}

	// Use provided color if available
	if len(color) > 0 {
		defaultColor = color[0]
	}

	// Determine the intStyle based on name
	var intStyle int
	switch name {
	case "regular":
		intStyle = 0
	case "italic":
		intStyle = 1
	case "bold":
		intStyle = 2
	case "underline":
		intStyle = 3
	default:
		// Default to regular if name is not recognized
		intStyle = 0
	}

	return FontStyle{
		name:     name,
		intStyle: intStyle,
		size:     size,
		color:    defaultColor,
	}
}

// WithSize returns a copy of this style with the specified size
func (fs FontStyle) WithSize(size float64) FontStyle {
	copy := fs
	copy.size = size
	return copy
}

// WithColor returns a copy of this style with the specified color
func (fs FontStyle) WithColor(color style.Color) FontStyle {
	copy := fs
	copy.color = color
	return copy
}

// Bitwise operations for FontStyle
func (fs FontStyle) AndNot(other FontStyle) FontStyle {
	result := fs
	result.intStyle = fs.intStyle &^ other.intStyle
	return result
}

// BitwiseAnd implements the & operator for FontStyle
func (fs FontStyle) BitwiseAnd(other FontStyle) FontStyle {
	result := fs
	result.intStyle = fs.intStyle & other.intStyle
	return result
}

// Equals checks if two FontStyle values are equal based on their intStyle
func (fs FontStyle) Equals(other FontStyle) bool {
	return fs.intStyle == other.intStyle
}

// GetName returns the name of the font style
func (fs FontStyle) GetName() string {
	return fs.name
}

// GetIntStyle returns the integer style value
func (fs FontStyle) GetIntStyle() int {
	return fs.intStyle
}

// GetSize returns the font size
func (fs FontStyle) GetSize() float64 {
	return fs.size
}

// GetColor returns the font color
func (fs FontStyle) GetColor() style.Color {
	return fs.color
}

// GetFamily returns the font family name
func (fs FontStyle) GetFamily() string {
	return fs.name
}

// SetName sets the name of the font style
func (fs *FontStyle) SetName(name string) {
	fs.name = name
}

// SetIntStyle sets the integer style value
func (fs *FontStyle) SetIntStyle(intStyle int) {
	fs.intStyle = intStyle
}

// SetSize sets the font size
func (fs *FontStyle) SetSize(size float64) {
	fs.size = size
}

// SetColor sets the font color
func (fs *FontStyle) SetColor(color style.Color) {
	fs.color = color
}

// SetFamily sets the font family name
func (fs *FontStyle) SetFamily(family string) {
	fs.name = family
}

// GetFontStyle converts a string font style to its FontStyle representation
// eg: "B" => FontStyleBold, "I" => FontStyleItalic, "U" => FontStyleUnderline
// defaults to FontStyleRegular if no valid style is found
func GetFontStyle(fontStyleStr string) FontStyle {
	fontStyleStr = tinystring.Convert(fontStyleStr).ToUpper().String()

	// Clone the regular style as a base
	result := FontStyleRegular

	if tinystring.Contains(fontStyleStr, "B") {
		result.intStyle = result.intStyle | FontStyleBold.intStyle

		// Update name to reflect the style combination
		if result.name == FontStyleRegular.name {
			result.name = FontStyleBold.name
		} else {
			result.name += "+" + FontStyleBold.name
		}
	}

	if tinystring.Contains(fontStyleStr, "I") {
		result.intStyle = result.intStyle | FontStyleItalic.intStyle

		// Update name to reflect the style combination
		if result.name == FontStyleRegular.name {
			result.name = FontStyleItalic.name
		} else {
			result.name += "+" + FontStyleItalic.name
		}
	}

	if tinystring.Contains(fontStyleStr, "U") {
		result.intStyle = result.intStyle | FontStyleUnderline.intStyle

		// Update name to reflect the style combination
		if result.name == FontStyleRegular.name {
			result.name = FontStyleUnderline.name
		} else {
			result.name += "+" + FontStyleUnderline.name
		}
	}

	return result
}

// GetCompleteFont creates a FontStyle with all properties specified
// Use this to create a complete font definition with style, size, color and family
func GetCompleteFont(styleStr string, size float64, color style.Color, family string) FontStyle {
	baseStyle := GetFontStyle(styleStr)
	baseStyle.SetSize(size)
	baseStyle.SetColor(color)
	baseStyle.SetFamily(family)
	return baseStyle
}

// GetFontStyleInIntFormat converts a string font style to its integer representation
// Kept for backward compatibility, but uses the new FontStyle struct internally
// eg: "B" => FontStyleBold, "I" => FontStyleItalic, "U" => FontStyleUnderline
// defaults to FontStyleRegular if no valid style is found
func GetFontStyleInIntFormat(fontStyle string) FontStyle {
	return GetFontStyle(fontStyle)
}
