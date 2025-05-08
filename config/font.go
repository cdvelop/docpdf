// Package config contains configuration structures for docpdf.
package config

import (
	"github.com/cdvelop/tinystring"
)

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

// FontIntStyle represents font styles as integer constants
// This was previously named FontStyle in docFont.go.
// eg: FontStyleBold = 2 (000010), FontStyleItalic = 1 (000001), etc.
type FontIntStyle int

// regular - font style regular
const FontStyleRegular FontIntStyle = 0 //000000
// Italic - font style italic
const FontStyleItalic FontIntStyle = 1 //000001
// Bold - font style bold
const FontStyleBold FontIntStyle = 2 //000010
// Underline - font style underline
const FontStyleUnderline FontIntStyle = 3 //000011

// GetFontStyleInIntFormat converts a string font style to its integer representation
// eg: "B" => FontStyleBold, "I" => FontStyleItalic, "U" => FontStyleUnderline
// defaults to FontStyleRegular if no valid style is found
func GetFontStyleInIntFormat(fontStyle string) (style FontIntStyle) {

	fontStyle = tinystring.Convert(fontStyle).ToUpper().String()

	style = FontStyleRegular // Default to regular style

	if tinystring.Contains(fontStyle, "B") {
		style = style | FontStyleBold
	}
	if tinystring.Contains(fontStyle, "I") {
		style = style | FontStyleItalic
	}
	if tinystring.Contains(fontStyle, "U") {
		style = style | FontStyleUnderline
	}

	return style
}
