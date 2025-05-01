package pdfengine

// point a point in a two-dimensional
type point struct {
	X float64
	Y float64
}

// SetX : set current alignment.Alignment X
func (gp *PdfEngine) SetX(x float64) {
	gp.UnitsToPointsVar(&x)
	gp.Curr.setXCount++
	gp.Curr.X = x
}

// GetX : get current alignment.Alignment X
func (gp *PdfEngine) GetX() float64 {
	return gp.pointsToUnits(gp.Curr.X)
}

// SetNewY : set current alignment.Alignment y, and modified y if add a new page.
// Example:
// For example, if the page height is set to 841px, MarginTop is 20px,
// MarginBottom is 10px, and the height of the element(such as text) to be inserted is 25px,
// because 10<25, you need to add another page and set y to 20px.
// Because of called AddPage(), X is set to MarginLeft, so you should specify X if needed,
// or make sure SetX() is after SetNewY(), or using SetNewXY().
// SetNewYIfNoOffset is more suitable for scenarios where the offset does not change, such as pdf.Image().
func (gp *PdfEngine) SetNewY(y float64, h float64) {
	gp.UnitsToPointsVar(&y)
	gp.UnitsToPointsVar(&h)
	if gp.Curr.Y+h > gp.Curr.pageSize.H-gp.MarginBottom() {
		gp.AddPage()
		y = gp.MarginTop() // reset to top of the page.
	}
	gp.Curr.Y = y
}

// SetNewYIfNoOffset : set current alignment.Alignment y, and modified y if add a new page.
// Example:
// For example, if the page height is set to 841px, MarginTop is 20px,
// MarginBottom is 10px, and the height of the element(such as image) to be inserted is 200px,
// because 10<200, you need to add another page and set y to 20px.
// Tips: gp.Curr.X and gp.Curr.Y do not change when pdf.Image() is called.
func (gp *PdfEngine) SetNewYIfNoOffset(y float64, h float64) {
	gp.UnitsToPointsVar(&y)
	gp.UnitsToPointsVar(&h)
	if y+h > gp.Curr.pageSize.H-gp.MarginBottom() { // using new y(*y) instead of gp.Curr.Y
		gp.AddPage()
		y = gp.MarginTop() // reset to top of the page.
	}
	gp.Curr.Y = y
}

// SetNewXY : set current alignment.Alignment x and y, and modified y if add a new page.
// Example:
// For example, if the page height is set to 841px, MarginTop is 20px,
// MarginBottom is 10px, and the height of the element to be inserted is 25px,
// because 10<25, you need to add another page and set y to 20px.
// Because of AddPage(), X is set to MarginLeft, so you should specify X if needed,
// or make sure SetX() is after SetNewY().
func (gp *PdfEngine) SetNewXY(y float64, x, h float64) {
	gp.UnitsToPointsVar(&y)
	gp.UnitsToPointsVar(&h)
	if gp.Curr.Y+h > gp.Curr.pageSize.H-gp.MarginBottom() {
		gp.AddPage()
		y = gp.MarginTop() // reset to top of the page.
	}
	gp.Curr.Y = y
	gp.SetX(x)
}

// SetY : set current alignment.Alignment y
func (gp *PdfEngine) SetY(y float64) {
	gp.UnitsToPointsVar(&y)
	gp.Curr.Y = y
}

// GetY : get current alignment.Alignment y
func (gp *PdfEngine) GetY() float64 {
	return gp.pointsToUnits(gp.Curr.Y)
}

// SetXY : set current alignment.Alignment x and y
func (gp *PdfEngine) SetXY(x, y float64) {
	gp.UnitsToPointsVar(&x)
	gp.Curr.setXCount++
	gp.Curr.X = x

	gp.UnitsToPointsVar(&y)
	gp.Curr.Y = y
}
