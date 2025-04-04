package docpdf

import (
	"fmt"
	"io"
	"os"
)

// embedFontObj is an embedded font object.
type embedFontObj struct {
	Data      string
	zfontpath string
	font      iFont
	getRoot   func() *pdfEngine
}

func (e *embedFontObj) init(funcGetRoot func() *pdfEngine) {
	e.getRoot = funcGetRoot
}

func (e *embedFontObj) protection() *pdfProtection {
	return e.getRoot().protection()
}

func (e *embedFontObj) write(w writer, objID int) error {
	b, err := os.ReadFile(e.zfontpath)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "<</Length %d\n", len(b))
	io.WriteString(w, "/Filter /FlateDecode\n")
	fmt.Fprintf(w, "/Length1 %d\n", e.font.GetOriginalsize())
	io.WriteString(w, ">>\n")
	io.WriteString(w, "stream\n")
	if e.protection() != nil {
		tmp, err := rc4Cip(e.protection().objectkey(objID), b)
		if err != nil {
			return err
		}
		w.Write(tmp)
		io.WriteString(w, "\n")
	} else {
		w.Write(b)
	}
	io.WriteString(w, "\nendstream\n")
	return nil
}

func (e *embedFontObj) getType() string {
	return "EmbedFont"
}

// SetFont sets the font of an embedded font object.
func (e *embedFontObj) SetFont(font iFont, zfontpath string) {
	e.font = font
	e.zfontpath = zfontpath
}
