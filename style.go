package docpdf

// Represents an RGB color with red, green, and blue components
type RGBColor struct {
	R uint8 // Red component (0-255)
	G uint8 // Green component (0-255)
	B uint8 // Blue component (0-255)
}

// Defines the border style for a cell or table
type BorderStyle struct {
	Top      bool     // Whether to draw the top border
	Left     bool     // Whether to draw the left border
	Right    bool     // Whether to draw the right border
	Bottom   bool     // Whether to draw the bottom border
	Width    float64  // Width of the border line
	RGBColor RGBColor // Color of the border
}

// Defines the style for a cell, including border, fill, text, and font properties
type CellStyle struct {
	BorderStyle BorderStyle // Border style for the cell
	FillColor   RGBColor    // Background color of the cell
	TextColor   RGBColor    // Color of the text in the cell
	Font        string      // Font name for the cell text
	FontSize    float64     // Font size for the cell text
}

// paintStyle represents the painting style for graphics
type paintStyle string

const (
	// drawPaintStyle is for drawing only
	drawPaintStyle paintStyle = "S"

	// fillPaintStyle is for filling only
	fillPaintStyle paintStyle = "F"

	// drawAndFillPaintStyle is for drawing and filling
	drawAndFillPaintStyle paintStyle = "B"
)

// parseStyle converts style strings to paintStyle constants
func parseStyle(style string) paintStyle {
	switch style {
	case "F":
		return fillPaintStyle
	case "DF", "FD":
		return drawAndFillPaintStyle
	default: // "D" or any other string
		return drawPaintStyle
	}
}
