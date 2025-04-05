package docpdf

import (
	"fmt"
	"io"

	"github.com/cdvelop/docpdf/fontmaker/core"
)

// subfontDescriptorObj pdf subfont descriptorObj object
type subfontDescriptorObj struct {
	PtrToSubsetFontObj    *subsetFontObj
	indexObjPdfDictionary int
}

func (s *subfontDescriptorObj) init(func() *pdfEngine) {}

func (s *subfontDescriptorObj) getType() string {
	return "SubFontDescriptor"
}

func (s *subfontDescriptorObj) write(w writer, objID int) error {
	ttfp := s.PtrToSubsetFontObj.GetTTFParser()
	//fmt.Printf("-->%d\n", ttfp.UnitsPerEm())
	io.WriteString(w, "<<\n")
	io.WriteString(w, "/Type /FontDescriptor\n")
	fmt.Fprintf(w, "/Ascent %d\n", designUnitsToPdf(ttfp.Ascender(), ttfp.UnitsPerEm()))
	fmt.Fprintf(w, "/CapHeight %d\n", designUnitsToPdf(ttfp.CapHeight(), ttfp.UnitsPerEm()))
	fmt.Fprintf(w, "/Descent %d\n", designUnitsToPdf(ttfp.Descender(), ttfp.UnitsPerEm()))
	fmt.Fprintf(w, "/Flags %d\n", ttfp.Flag())
	fmt.Fprintf(w, "/FontBBox [%d %d %d %d]\n",
		designUnitsToPdf(ttfp.XMin(), ttfp.UnitsPerEm()),
		designUnitsToPdf(ttfp.YMin(), ttfp.UnitsPerEm()),
		designUnitsToPdf(ttfp.XMax(), ttfp.UnitsPerEm()),
		designUnitsToPdf(ttfp.YMax(), ttfp.UnitsPerEm()),
	)
	fmt.Fprintf(w, "/FontFile2 %d 0 R\n", s.indexObjPdfDictionary+1)
	fmt.Fprintf(w, "/FontName /%s\n", createEmbeddedFontSubsetName(s.PtrToSubsetFontObj.GetFamily()))
	fmt.Fprintf(w, "/ItalicAngle %d\n", ttfp.ItalicAngle())
	io.WriteString(w, "/StemV 0\n")
	fmt.Fprintf(w, "/XHeight %d\n", designUnitsToPdf(ttfp.XHeight(), ttfp.UnitsPerEm()))
	io.WriteString(w, ">>\n")
	return nil
}

// SetIndexObjPdfDictionary set PdfDictionary pointer
func (s *subfontDescriptorObj) SetIndexObjPdfDictionary(index int) {
	s.indexObjPdfDictionary = index
}

// SetPtrToSubsetFontObj set SubsetFont pointer
func (s *subfontDescriptorObj) SetPtrToSubsetFontObj(ptr *subsetFontObj) {
	s.PtrToSubsetFontObj = ptr
}

// designUnitsToPdf convert unit
func designUnitsToPdf(val int, unitsPerEm uint) int {
	return core.Round(float64(float64(val) * 1000.00 / float64(unitsPerEm)))
}
