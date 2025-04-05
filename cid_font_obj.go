package docpdf

import (
	"fmt"
	"io"
)

// cidFontObj is a CID-keyed font.
// cf. https://www.adobe.com/content/dam/acom/en/devnet/font/pdfs/5014.CIDFont_Spec.pdf
type cidFontObj struct {
	PtrToSubsetFontObj        *subsetFontObj
	indexObjSubfontDescriptor int
}

func (ci *cidFontObj) init(funcGetRoot func() *pdfEngine) {
}

// SetIndexObjSubfontDescriptor set  indexObjSubfontDescriptor
func (ci *cidFontObj) SetIndexObjSubfontDescriptor(index int) {
	ci.indexObjSubfontDescriptor = index
}

func (ci *cidFontObj) getType() string {
	return "CIDFont"
}

func (ci *cidFontObj) write(w writer, objID int) error {
	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "/BaseFont /%s\n", createEmbeddedFontSubsetName(ci.PtrToSubsetFontObj.GetFamily()))
	io.WriteString(w, "/CIDSystemInfo\n")
	io.WriteString(w, "<<\n")
	io.WriteString(w, "  /Ordering (Identity)\n")
	io.WriteString(w, "  /Registry (Adobe)\n")
	io.WriteString(w, "  /Supplement 0\n")
	io.WriteString(w, ">>\n")
	fmt.Fprintf(w, "/FontDescriptor %d 0 R\n", ci.indexObjSubfontDescriptor+1) //TODO fix
	io.WriteString(w, "/Subtype /CIDFontType2\n")
	io.WriteString(w, "/Type /Font\n")
	glyphIndexs := ci.PtrToSubsetFontObj.CharacterToGlyphIndex.AllVals()
	io.WriteString(w, "/W [")
	for _, v := range glyphIndexs {
		width := ci.PtrToSubsetFontObj.GlyphIndexToPdfWidth(v)
		fmt.Fprintf(w, "%d[%d]", v, width)
	}
	io.WriteString(w, "]\n")
	io.WriteString(w, ">>\n")
	return nil
}

// SetPtrToSubsetFontObj set PtrToSubsetFontObj
func (ci *cidFontObj) SetPtrToSubsetFontObj(ptr *subsetFontObj) {
	ci.PtrToSubsetFontObj = ptr
}
