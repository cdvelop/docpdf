package docpdf

import "strings"

// FontConfig represents different font configurations for document sections
type FontConfig struct {
	Family     Font
	Normal     TextStyle
	Header1    TextStyle
	Header2    TextStyle
	Header3    TextStyle
	Footnote   TextStyle
	PageHeader TextStyle
	PageFooter TextStyle
}

// Font represents font files for different styles
type Font struct {
	Regular string
	Bold    string
	Italic  string
	Path    string // Base path for fonts
}

// loadFonts loads the fonts from the Font struct
func (d *Document) loadFonts() error {
	fontPath := d.fontConfig.Family.Path

	// add regular font
	if err := d.AddTTFFont(FontRegular, fontPath+d.fontConfig.Family.Regular); err != nil {
		return err
	}
	// add bold font
	if d.fontConfig.Family.Bold == "" {
		d.fontConfig.Family.Bold = d.fontConfig.Family.Regular
	} else {
		if err := d.AddTTFFont(FontBold, fontPath+d.fontConfig.Family.Bold); err != nil {
			return err
		}
	}
	// add italic font
	if d.fontConfig.Family.Italic == "" {
		d.fontConfig.Family.Italic = d.fontConfig.Family.Regular
	} else {
		if err := d.AddTTFFont(FontItalic, fontPath+d.fontConfig.Family.Italic); err != nil {
			return err
		}
	}

	if d.fontConfig.Family.Path == "" {
		d.fontConfig.Family.Path = "fonts/"
	}

	return nil
}

// extracts the font name from the font path eg: "fonts/Rubik-Regular.ttf" => "Rubik-Regular"
func extractNameFromPath(path string) string {
	if path == "" {
		return ""
	}
	// normalize path separators to forward slash
	path = strings.ReplaceAll(path, "\\", "/")

	// split the path by "/"
	parts := strings.Split(path, "/")

	// get the last part (filename)
	filename := parts[len(parts)-1]

	// split by dot to get all parts
	nameParts := strings.Split(filename, ".")

	// remove the last part if it's an extension (ttf, otf, etc)
	if len(nameParts) > 1 {
		nameParts = nameParts[:len(nameParts)-1]
	}

	// join all parts without dots
	return strings.Join(nameParts, "")
}

func (d *Document) setDefaultFont() {
	style := d.fontConfig.Normal
	d.SetFont(FontRegular, "", style.Size)
	d.SetTextColor(style.Color.R, style.Color.G, style.Color.B)
	d.SetLineWidth(0.7)
	d.SetStrokeColor(0, 0, 0)
}

// defaultFontConfig returns word-processor like defaults
func defaultFontConfig() FontConfig {
	return FontConfig{
		Family: Font{
			Regular: "Rubik-Regular.ttf",
			Bold:    "Rubik-Bold.ttf",
			Italic:  "Rubik-Italic.ttf",
			Path:    "fonts/",
		},

		Normal: TextStyle{
			Size:        11,
			Color:       RGBColor{0, 0, 0},
			LineSpacing: 1.15,
			Alignment:   Left | Top,
			SpaceBefore: 0,
			SpaceAfter:  8, // ~0.73x font size (Word default is similar)
		},
		Header1: TextStyle{
			Size:        16,
			Color:       RGBColor{0, 0, 0},
			LineSpacing: 1.5,
			Alignment:   Left | Top,
			SpaceBefore: 12,
			SpaceAfter:  8,
		},
		Header2: TextStyle{
			Size:        14,
			Color:       RGBColor{0, 0, 0},
			LineSpacing: 1.3,
			Alignment:   Left | Top,
			SpaceBefore: 10,
			SpaceAfter:  6,
		},
		Header3: TextStyle{
			Size:        12,
			Color:       RGBColor{0, 0, 0},
			LineSpacing: 1.2,
			Alignment:   Left | Top,
			SpaceBefore: 8,
			SpaceAfter:  4,
		},
		Footnote: TextStyle{
			Size:        8,
			Color:       RGBColor{128, 128, 128},
			LineSpacing: 1.0,
			Alignment:   Left | Top,
			SpaceBefore: 2,
			SpaceAfter:  2,
		},
		PageHeader: TextStyle{
			Size:        9,
			Color:       RGBColor{128, 128, 128},
			LineSpacing: 1.0,
			Alignment:   Center | Top,
			SpaceBefore: 0,
			SpaceAfter:  3,
		},
		PageFooter: TextStyle{
			Size:        9,
			Color:       RGBColor{128, 128, 128},
			LineSpacing: 1.0,
			Alignment:   Right | Top,
			SpaceBefore: 2,
			SpaceAfter:  0,
		},
	}
}
