package docpdf

// SetLeftMargin sets left margin.
func (gp *pdfEngine) SetLeftMargin(margin float64) {
	gp.unitsToPointsVar(&margin)
	gp.margins.Left = margin
}

// SetTopMargin sets top margin.
func (gp *pdfEngine) SetTopMargin(margin float64) {
	gp.unitsToPointsVar(&margin)
	gp.margins.Top = margin
}

// SetMargins defines the left, top, right and bottom margins. By default, they equal 1 cm. Call this method to change them.
func (gp *pdfEngine) SetMargins(left, top, right, bottom float64) {
	gp.unitsToPointsVar(&left, &top, &right, &bottom)
	gp.margins = Margins{left, top, right, bottom}
}

// SetMarginLeft sets the left margin
func (gp *pdfEngine) SetMarginLeft(margin float64) {
	gp.margins.Left = gp.unitsToPoints(margin)
}

// SetMarginTop sets the top margin
func (gp *pdfEngine) SetMarginTop(margin float64) {
	gp.margins.Top = gp.unitsToPoints(margin)
}

// SetMarginRight sets the right margin
func (gp *pdfEngine) SetMarginRight(margin float64) {
	gp.margins.Right = gp.unitsToPoints(margin)
}

// SetMarginBottom set the bottom margin
func (gp *pdfEngine) SetMarginBottom(margin float64) {
	gp.margins.Bottom = gp.unitsToPoints(margin)
}

// Margins gets the current margins, The margins will be converted back to the documents units. Returned values will be in the following order Left, Top, Right, Bottom
func (gp *pdfEngine) Margins() (float64, float64, float64, float64) {
	return gp.pointsToUnits(gp.margins.Left),
		gp.pointsToUnits(gp.margins.Top),
		gp.pointsToUnits(gp.margins.Right),
		gp.pointsToUnits(gp.margins.Bottom)
}

// MarginLeft returns the left margin.
func (gp *pdfEngine) MarginLeft() float64 {
	return gp.pointsToUnits(gp.margins.Left)
}

// MarginTop returns the top margin.
func (gp *pdfEngine) MarginTop() float64 {
	return gp.pointsToUnits(gp.margins.Top)
}

// MarginRight returns the right margin.
func (gp *pdfEngine) MarginRight() float64 {
	return gp.pointsToUnits(gp.margins.Right)
}

// MarginBottom returns the bottom margin.
func (gp *pdfEngine) MarginBottom() float64 {
	return gp.pointsToUnits(gp.margins.Bottom)
}
