package pdfengine

import "github.com/cdvelop/docpdf/canvas"

// SetLeftMargin sets left margin.
func (gp *PdfEngine) SetLeftMargin(margin float64) {
	gp.UnitsToPointsVar(&margin)
	gp.margins.Left = margin
}

// SetTopMargin sets top margin.
func (gp *PdfEngine) SetTopMargin(margin float64) {
	gp.UnitsToPointsVar(&margin)
	gp.margins.Top = margin
}

// SetMargins defines the left, top, right and bottom canvas.Margins. By default, they equal 1 cm. Call this method to change them.
func (gp *PdfEngine) SetMargins(left, top, right, bottom float64) {
	gp.UnitsToPointsVar(&left, &top, &right, &bottom)
	gp.margins = canvas.Margins{left, top, right, bottom}
}

// SetMarginLeft sets the left margin
func (gp *PdfEngine) SetMarginLeft(margin float64) {
	gp.margins.Left = gp.UnitsToPoints(margin)
}

// SetMarginTop sets the top margin
func (gp *PdfEngine) SetMarginTop(margin float64) {
	gp.margins.Top = gp.UnitsToPoints(margin)
}

// SetMarginRight sets the right margin
func (gp *PdfEngine) SetMarginRight(margin float64) {
	gp.margins.Right = gp.UnitsToPoints(margin)
}

// SetMarginBottom set the bottom margin
func (gp *PdfEngine) SetMarginBottom(margin float64) {
	gp.margins.Bottom = gp.UnitsToPoints(margin)
}

func (gp *PdfEngine) Margins() canvas.Margins {
	return gp.margins
}

// AllMargins gets the current Margins, The Margins will be converted back to the documents units. Returned values will be in the following order Left, Top, Right, Bottom
func (gp *PdfEngine) AllMargins() (float64, float64, float64, float64) {
	return gp.pointsToUnits(gp.margins.Left),
		gp.pointsToUnits(gp.margins.Top),
		gp.pointsToUnits(gp.margins.Right),
		gp.pointsToUnits(gp.margins.Bottom)
}

// MarginLeft returns the left margin.
func (gp *PdfEngine) MarginLeft() float64 {
	return gp.pointsToUnits(gp.margins.Left)
}

// MarginTop returns the top margin.
func (gp *PdfEngine) MarginTop() float64 {
	return gp.pointsToUnits(gp.margins.Top)
}

// MarginRight returns the right margin.
func (gp *PdfEngine) MarginRight() float64 {
	return gp.pointsToUnits(gp.margins.Right)
}

// MarginBottom returns the bottom margin.
func (gp *PdfEngine) MarginBottom() float64 {
	return gp.pointsToUnits(gp.margins.Bottom)
}
