package docpdf

import (
	"fmt"
	"io"
)

// fontDescriptorObj is a font descriptor object.
type fontDescriptorObj struct {
	font              iFont
	fontFileObjRelate string
}

func (f *fontDescriptorObj) init(funcGetRoot func() *pdfEngine) {

}

func (f *fontDescriptorObj) write(w io.Writer, objID int) error {

	fmt.Fprintf(w, "<</Type /FontDescriptor /FontName /%s ", f.font.GetName())
	descs := f.font.GetDesc()
	i := 0
	max := len(descs)
	for i < max {
		fmt.Fprintf(w, "/%s %s ", descs[i].Key, descs[i].Val)
		i++
	}

	if f.getType() == "Type1" {
		io.WriteString(w, "/FontFile ")
	} else {
		io.WriteString(w, "/FontFile2 ")
	}

	io.WriteString(w, f.fontFileObjRelate)
	io.WriteString(w, ">>\n")

	return nil
}

func (f *fontDescriptorObj) getType() string {
	return "FontDescriptor"
}

// SetFont sets the font in descriptor.
func (f *fontDescriptorObj) SetFont(font iFont) {
	f.font = font
}

// GetFont gets font from descriptor.
func (f *fontDescriptorObj) GetFont() iFont {
	return f.font
}

// SetFontFileObjRelate ???
func (f *fontDescriptorObj) SetFontFileObjRelate(relate string) {
	f.fontFileObjRelate = relate
}
