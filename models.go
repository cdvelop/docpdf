package docpdf

import (
	"bytes"
	"io"
	"time"
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
	init(func() *pdfEngine)
	getType() string
	write(w io.Writer, objID int) error
}

// Margins type.
type Margins struct {
	Left, Top, Right, Bottom float64
}

// pdfEngine : core library for generating PDF
type pdfEngine struct {

	// Page margins
	margins Margins

	pdfObjs []iObj
	config  config
	anchors map[string]anchorOption

	indexOfCatalogObj int

	/*--- Important obj indexes stored to reduce search loops ---*/
	// Index of pages obj
	indexOfPagesObj int

	// Number of pages obj
	numOfPagesObj int

	// Index of first page obj
	indexOfFirstPageObj int

	// currentPdf position
	curr currentPdf

	indexEncodingObjFonts []int
	indexOfContent        int

	// Index of procset which should be unique
	indexOfProcSet int

	// Buffer for io.Reader compliance
	buf bytes.Buffer

	// PDF protection
	pdfProtection   *pdfProtection
	encryptionObjID int

	// Content streams only
	compressLevel int

	// Document info
	isUseInfo bool
	info      *PdfInfo

	// Outlines/bookmarks
	outlines           *outlinesObj
	indexOfOutlinesObj int

	// Header and footer functions
	headerFunc func()
	footerFunc func()

	// gofpdi free pdf document importer
	fpdi *importer

	// Placeholder text
	placeHolderTexts map[string]([]placeHolderTextInfo)

	// Log function for debugging
	log func(...any)
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
	pageSize *Rect

	//current trim box
	trimBox *box

	sMasksMap       sMaskMap
	extGStatesMap   extGStatesMap
	transparency    *transparency
	transparencyMap transparencyMap
}

// box represents a rectangular area with explicit coordinates for all four sides.
// It is used for defining boundaries in PDF documents, such as margins, trim boxes, etc.
// The coordinates are stored in the current unit system (points by default, but can be mm, cm, inches, or pixels).
type box struct {
	Left, Top, Right, Bottom float64
	unitOverride             defaultUnitConfig
}

// Rect defines a rectangle by its width and height.
// This is used for defining page sizes, content areas, and other rectangular regions in PDF documents.
// The dimensions are stored in the current unit system (points by default, but can be mm, cm, inches, or pixels).
type Rect struct {
	W            float64 // Width of the rectangle
	H            float64 // Height of the rectangle
	unitOverride defaultUnitConfig
}

// defaultUnitConfig is the standard implementation of the unitConfigurator interface.
// It stores the unit type and an optional custom conversion factor.
type defaultUnitConfig struct {
	// Unit specifies the unit type (UnitPT, UnitMM, UnitCM, UnitIN, UnitPX)
	Unit int

	// ConversionForUnit is an optional custom conversion factor
	ConversionForUnit float64
}

type iCacheContent interface {
	write(w io.Writer, protection *pdfProtection) error
}

type iCacheColorText interface {
	iCacheContent
	equal(obj iCacheColorText) bool
}
