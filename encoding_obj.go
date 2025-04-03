package docpdf

import (
	"io"
)

// encodingObj is a font object.
type encodingObj struct {
	font iFont
}

func (e *encodingObj) init(funcGetRoot func() *pdfEngine) {

}
func (e *encodingObj) getType() string {
	return "Encoding"
}
func (e *encodingObj) write(w io.Writer, objID int) error {
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
