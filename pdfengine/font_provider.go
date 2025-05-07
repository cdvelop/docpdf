package pdfengine

// Name returns the unique name of the current font.
// It implements the fontengine.FontProvider interface.
func (gp *PdfEngine) Name() string {
	// The Family, Weight, and Style methods handle nil checks for FontISubset internally
	// and provide default values if necessary.
	return gp.Family() + "-" + gp.Weight() + "-" + gp.Style()
}

// Family returns the family of the current font.
// It implements the fontengine.FontProvider interface.
func (gp *PdfEngine) Family() string {
	if gp.curr.FontISubset == nil {
		return "UnknownFamily" // Default family if FontISubset is nil
	}
	return gp.curr.FontISubset.GetFamily()
}

// Weight returns the weight of the current font (e.g., "regular", "bold").
// It implements the fontengine.FontProvider interface.
func (gp *PdfEngine) Weight() string {
	// FontStyle is part of currentPdf, which is not a pointer, so gp.curr itself is not nil.
	if (gp.curr.FontStyle & Bold) != 0 {
		return "bold"
	}
	return "regular"
}

// Style returns the style of the current font (e.g., "normal", "italic").
// It implements the fontengine.FontProvider interface.
func (gp *PdfEngine) Style() string {
	// FontStyle is part of currentPdf.
	if (gp.curr.FontStyle & Italic) != 0 {
		return "italic"
	}
	return "normal"
}

// SVGFontID returns the ID for referencing the font in SVG.
// It implements the fontengine.FontProvider interface.
func (gp *PdfEngine) SVGFontID() string {
	return gp.Name()
}

// Path returns the file path of the font, if available.
// For PdfEngine, fonts are typically embedded, so a file path is not directly relevant.
// It implements the fontengine.FontProvider interface.
func (gp *PdfEngine) Path() string {
	return ""
}

//tool for validate pdf https://www.pdf-online.com/osa/validate.aspx
