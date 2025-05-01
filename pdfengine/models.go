package pdfengine

import (
	"io"
	"time"

	"github.com/cdvelop/docpdf/canvas"
)

type anchorOption struct {
	page int
	y    float64
}

type linkOption struct {
	x, y, w, h float64
	url        string
	anchor     string
}

// PdfInfo Document Information Dictionary
type PdfInfo struct {
	Title        string    // The document’s title
	Author       string    // The name of the person who created the document
	Subject      string    // The subject of the document
	Creator      string    // If the document was converted to PDF from another format, the name of the application which created the original document
	Producer     string    // If the document was converted to PDF from another format, the name of the application that converted the original document to PDF
	CreationDate time.Time // The date and time the document was created, in human-readable form
}

// iObj inteface for all pdf object
type iObj interface {
	init(func() *PdfEngine)
	getType() string
	write(w io.Writer, objID int) error
}

type placeHolderTextInfo struct {
	indexOfContent   int
	indexInContent   int
	fontISubset      *subsetFontObj
	placeHolderWidth float64
	fontSize         float64
	charSpacing      float64
}

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

	FontISubset *subsetFontObj // FontType == CURRENT_FONT_TYPE_SUBSET

	//page
	IndexOfPageObj int

	//img
	CountOfImg int
	//cache of image in pdf file
	ImgCaches map[int]imageCache

	//text color mode
	txtColorMode string //color, gray

	//text color
	txtColor iCacheColorText

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

type iCacheContent interface {
	write(w io.Writer, protection *pdfProtection) error
}

type iCacheColorText interface {
	iCacheContent
	equal(obj iCacheColorText) bool
}
