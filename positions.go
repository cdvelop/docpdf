package docpdf

// position representa una posición o alineación en el documento
type position int

const (
	// Left representa alineación a la izquierda
	Left position = 8 //001000
	// Top representa alineación superior
	Top position = 4 //000100
	// Right representa alineación a la derecha
	Right position = 2 //000010
	// Bottom representa alineación inferior
	Bottom position = 1 //000001
	// Center representa alineación al centro
	Center position = 16 //010000
	// Middle representa alineación al medio
	Middle position = 32 //100000
	// Justify representa texto justificado
	Justify position = 64 //1000000
	// AllBorders representa todos los bordes
	AllBorders position = 15 //001111
)

// point a point in a two-dimensional
type point struct {
	X float64
	Y float64
}

// SetX : set current position X
func (gp *pdfEngine) SetX(x float64) {
	gp.unitsToPointsVar(&x)
	gp.curr.setXCount++
	gp.curr.X = x
}

// GetX : get current position X
func (gp *pdfEngine) GetX() float64 {
	return gp.pointsToUnits(gp.curr.X)
}

// SetNewY : set current position y, and modified y if add a new page.
// Example:
// For example, if the page height is set to 841px, MarginTop is 20px,
// MarginBottom is 10px, and the height of the element(such as text) to be inserted is 25px,
// because 10<25, you need to add another page and set y to 20px.
// Because of called AddPage(), X is set to MarginLeft, so you should specify X if needed,
// or make sure SetX() is after SetNewY(), or using SetNewXY().
// SetNewYIfNoOffset is more suitable for scenarios where the offset does not change, such as pdf.Image().
func (gp *pdfEngine) SetNewY(y float64, h float64) {
	gp.unitsToPointsVar(&y)
	gp.unitsToPointsVar(&h)
	if gp.curr.Y+h > gp.curr.pageSize.H-gp.MarginBottom() {
		gp.AddPage()
		y = gp.MarginTop() // reset to top of the page.
	}
	gp.curr.Y = y
}

// SetNewYIfNoOffset : set current position y, and modified y if add a new page.
// Example:
// For example, if the page height is set to 841px, MarginTop is 20px,
// MarginBottom is 10px, and the height of the element(such as image) to be inserted is 200px,
// because 10<200, you need to add another page and set y to 20px.
// Tips: gp.curr.X and gp.curr.Y do not change when pdf.Image() is called.
func (gp *pdfEngine) SetNewYIfNoOffset(y float64, h float64) {
	gp.unitsToPointsVar(&y)
	gp.unitsToPointsVar(&h)
	if y+h > gp.curr.pageSize.H-gp.MarginBottom() { // using new y(*y) instead of gp.curr.Y
		gp.AddPage()
		y = gp.MarginTop() // reset to top of the page.
	}
	gp.curr.Y = y
}

// SetNewXY : set current position x and y, and modified y if add a new page.
// Example:
// For example, if the page height is set to 841px, MarginTop is 20px,
// MarginBottom is 10px, and the height of the element to be inserted is 25px,
// because 10<25, you need to add another page and set y to 20px.
// Because of AddPage(), X is set to MarginLeft, so you should specify X if needed,
// or make sure SetX() is after SetNewY().
func (gp *pdfEngine) SetNewXY(y float64, x, h float64) {
	gp.unitsToPointsVar(&y)
	gp.unitsToPointsVar(&h)
	if gp.curr.Y+h > gp.curr.pageSize.H-gp.MarginBottom() {
		gp.AddPage()
		y = gp.MarginTop() // reset to top of the page.
	}
	gp.curr.Y = y
	gp.SetX(x)
}

// SetY : set current position y
func (gp *pdfEngine) SetY(y float64) {
	gp.unitsToPointsVar(&y)
	gp.curr.Y = y
}

// GetY : get current position y
func (gp *pdfEngine) GetY() float64 {
	return gp.pointsToUnits(gp.curr.Y)
}

// SetXY : set current position x and y
func (gp *pdfEngine) SetXY(x, y float64) {
	gp.unitsToPointsVar(&x)
	gp.curr.setXCount++
	gp.curr.X = x

	gp.unitsToPointsVar(&y)
	gp.curr.Y = y
}
