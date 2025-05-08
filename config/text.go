package config

import (
	"github.com/cdvelop/docpdf/style"
)

// TextStyle defines the style configuration for different text sections
type TextStyle struct {
	Size        float64
	Color       style.Color
	LineSpacing float64
	Alignment   Alignment
	SpaceBefore float64
	SpaceAfter  float64
}

// TextStyles represents different font configurations for document sections.
// This was previously named FontConfig in docFont.go.
type TextStyles struct {
	FontFamily     FontFamily
	Normal         TextStyle
	Header1        TextStyle
	Header2        TextStyle
	Header3        TextStyle
	Footnote       TextStyle
	PageHeader     TextStyle
	PageFooter     TextStyle
	ChartLabel     TextStyle // For chart bar labels/values
	ChartAxisLabel TextStyle // For chart axis labels (X/Y axes)
}

// DefaultTextStyles returns word-processor like defaults
func DefaultTextStyles() TextStyles {
	return TextStyles{
		FontFamily: FontFamily{
			// Use standardized filenames for default fonts
			Regular: "regular.ttf",
			Bold:    "bold.ttf",
			Italic:  "italic.ttf",
			Path:    "fonts/",
		},

		Normal: TextStyle{
			Size:        11,
			Color:       style.Color{R: 0, G: 0, B: 0, A: 255}, // Color negro opaco (A=255)
			LineSpacing: 1.15,
			Alignment:   Left | Top,
			SpaceBefore: 0,
			SpaceAfter:  8, // ~0.73x font size (Word default is similar)
		},
		Header1: TextStyle{
			Size:        16,
			Color:       style.Color{R: 0, G: 0, B: 0, A: 255}, // Color negro opaco (A=255)
			LineSpacing: 1.5,
			Alignment:   Left | Top,
			SpaceBefore: 12,
			SpaceAfter:  8,
		},
		Header2: TextStyle{
			Size:        14,
			Color:       style.Color{R: 0, G: 0, B: 0, A: 255}, // Color negro opaco (A=255)
			LineSpacing: 1.3,
			Alignment:   Left | Top,
			SpaceBefore: 10,
			SpaceAfter:  6,
		},
		Header3: TextStyle{
			Size:        12,
			Color:       style.Color{R: 0, G: 0, B: 0, A: 255}, // Color negro opaco (A=255)
			LineSpacing: 1.2,
			Alignment:   Left | Top,
			SpaceBefore: 8,
			SpaceAfter:  4,
		},
		Footnote: TextStyle{
			Size:        8,
			Color:       style.Color{R: 128, G: 128, B: 128, A: 255}, // Color gris opaco (A=255)
			LineSpacing: 1.0,
			Alignment:   Left | Top,
			SpaceBefore: 2,
			SpaceAfter:  2,
		},
		PageHeader: TextStyle{
			Size:        9,
			Color:       style.Color{R: 128, G: 128, B: 128, A: 255}, // Color gris opaco (A=255)
			LineSpacing: 1.0,
			Alignment:   Center | Top,
			SpaceBefore: 0,
			SpaceAfter:  12,
		},
		PageFooter: TextStyle{
			Size:        9,
			Color:       style.Color{R: 128, G: 128, B: 128, A: 255}, // Color gris opaco (A=255)
			LineSpacing: 1.0,
			Alignment:   Right | Top,
			SpaceBefore: 2,
			SpaceAfter:  0,
		},
		ChartLabel: TextStyle{ // Added default chart style
			Size:        9,                                        // Slightly smaller than Normal (11)
			Color:       style.Color{R: 50, G: 50, B: 50, A: 255}, // Dark Gray opaco (A=255)
			LineSpacing: 1.0,
			Alignment:   Left | Top, // Default alignment, chart might override
			SpaceBefore: 0,
			SpaceAfter:  0,
		},
		ChartAxisLabel: TextStyle{ // Style for X/Y axis labels
			Size:        8,                                        // Smaller than ChartLabel
			Color:       style.Color{R: 80, G: 80, B: 80, A: 255}, // Medium Gray opaco (A=255)
			LineSpacing: 1.0,
			Alignment:   Left | Top,
			SpaceBefore: 0,
			SpaceAfter:  0,
		},
	}
}
