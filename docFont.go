package docpdf

import (
	"strings"
)

// loadFonts loads the fonts from the Font struct
func (d *Document) loadFonts() error {
	// Set default path if not provided
	if d.fontConfig.FontFamily.Path == "" {
		d.fontConfig.FontFamily.Path = "fonts/"
	}

	// Set default values if not provided
	if d.fontConfig.FontFamily.Bold == "" {
		d.fontConfig.FontFamily.Bold = d.fontConfig.FontFamily.Regular
	}

	if d.fontConfig.FontFamily.Italic == "" {
		d.fontConfig.FontFamily.Italic = d.fontConfig.FontFamily.Regular
	}
	// Use AddFontFamilyConfig to load all fonts at once
	return d.PdfEngine.AddFontFamilyConfig(d.fontConfig.FontFamily)
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
