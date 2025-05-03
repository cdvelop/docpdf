package pdfengine

import (
	"time"
)

type Writer interface {
	Write(p []byte) (n int, err error)
}

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
	Init(func() *PdfEngine)
	GetType() string
	Write(w Writer, objID int) error
}

type placeHolderTextInfo struct {
	indexOfContent   int
	indexInContent   int
	fontISubset      *subsetFontObj
	placeHolderWidth float64
	fontSize         float64
	charSpacing      float64
}

type ICacheContent interface {
	Write(w Writer, protection *pdfProtection) error
}

type ICacheColorText interface {
	ICacheContent
	Equal(obj ICacheColorText) bool
}
