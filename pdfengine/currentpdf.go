package pdfengine

import "github.com/cdvelop/docpdf/canvas"

// currentPdf current state
type currentPdf struct {
	setXCount int //many times we go func SetX()
	X         float64
	Y         float64

	//font
	IndexOfFontObj int
	CountOfFont    int
	CountOfL       int

	FontSize      float64
	FontStyle     int // Regular|Bold|Italic|Underline
	FontFontCount int
	FontType      int // CURRENT_FONT_TYPE_IFONT or  CURRENT_FONT_TYPE_SUBSET

	CharSpacing float64

	FontISubset *ttfSubsetObj // FontType == CURRENT_FONT_TYPE_SUBSET

	//page
	IndexOfPageObj int

	//img
	CountOfImg int
	//cache of image in pdf file
	ImgCaches map[int]imageCache

	//text color mode
	txtColorMode string //color, gray

	//text color
	txtColor ICacheColorText

	//text grayscale
	grayFill float64
	//draw grayscale
	grayStroke float64

	lineWidth float64

	//current page size
	pageSize *canvas.Rect

	//current trim canvas.Box
	trimBox *canvas.Box

	sMasksMap       sMaskMap
	extGStatesMap   extGStatesMap
	transparency    *transparency
	transparencyMap transparencyMap
}

// metodo que retorna pageSize
func (c *currentPdf) PageSize() *canvas.Rect {
	return c.pageSize
}

// metodo que retorna txtColor o nil
func (c *currentPdf) TxtColor() ICacheColorText {
	return c.txtColor
}

// metodo que set txtColor
func (c *currentPdf) SetTxtColor(txtColor ICacheColorText) {
	c.txtColor = txtColor
}

// metodo que retorna txtColorMode
func (c *currentPdf) TxtColorMode() string {
	return c.txtColorMode
}

// metodo que set txtColorMode
func (c *currentPdf) SetTxtColorMode(txtColorMode string) {
	c.txtColorMode = txtColorMode
}

// metodo que retorna grayFill
func (c *currentPdf) GrayFill() float64 {
	return c.grayFill
}

// metodo que set grayFill
func (c *currentPdf) SetGrayFill(grayFill float64) {
	c.grayFill = grayFill
}
