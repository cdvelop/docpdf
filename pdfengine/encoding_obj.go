package pdfengine

import (
	"io"
)

// encodingObj is a font object.
type encodingObj struct {
	font iFont
}

func (e *encodingObj) Init(funcGetRoot func() *PdfEngine) {

}
func (e *encodingObj) GetType() string {
	return "Encoding"
}
func (e *encodingObj) Write(w Writer, objID int) error {
	io.WriteString(w, "<</Type /Encoding /BaseEncoding /WinAnsiEncoding /Differences [")
	io.WriteString(w, e.font.GetDiff())
	io.WriteString(w, "]>>\n")
	return nil
}

// SetFont sets the font of an encoding object.
func (e *encodingObj) SetFont(font iFont) {
	e.font = font
}

// GetFont gets the font from an encoding object.
func (e *encodingObj) GetFont() iFont {
	return e.font
}
