package pdfengine

import ()

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
