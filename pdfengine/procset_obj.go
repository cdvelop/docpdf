package pdfengine

import (
	"fmt"
	"io"
)

// procSetObj is a PDF procSet object.
type procSetObj struct {
	//Font
	Relates             relateFonts
	RelateXobjs         relateXobjects
	ExtGStates          []extGS
	ImportedTemplateIds map[string]int
	getRoot             func() *PdfEngine
}

func (pr *procSetObj) Init(funcGetRoot func() *PdfEngine) {
	pr.getRoot = funcGetRoot
	pr.ImportedTemplateIds = make(map[string]int, 0)
	pr.ExtGStates = make([]extGS, 0)
}

func (pr *procSetObj) Write(w Writer, objID int) error {
	content := "<<\n"
	content += "\t/ProcSet [/PDF /Text /ImageB /ImageC /ImageI]\n"

	fonts := "\t/Font <<\n"
	for _, relate := range pr.Relates {
		fonts += fmt.Sprintf("\t\t/F%d %d 0 R\n", relate.CountOfFont+1, relate.IndexOfObj+1)
	}
	fonts += "\t>>\n"

	content += fonts

	xobjects := "\t/XObject <<\n"
	for _, XObject := range pr.RelateXobjs {
		xobjects += fmt.Sprintf("\t\t/I%d %d 0 R\n", XObject.IndexOfObj+1, XObject.IndexOfObj+1)
	}
	// Write imported template name and their ids
	for tplName, objID := range pr.ImportedTemplateIds {
		xobjects += fmt.Sprintf("\t\t%s %d 0 R\n", tplName, objID)
	}
	xobjects += "\t>>\n"

	content += xobjects

	extGStates := "\t/extGState <<\n"
	for _, extGState := range pr.ExtGStates {
		extGStates += fmt.Sprintf("\t\t/GS%d %d 0 R\n", extGState.Index+1, extGState.Index+1)
	}
	extGStates += "\t>>\n"

	content += extGStates

	content += ">>\n"

	if _, err := io.WriteString(w, content); err != nil {
		return err
	}

	return nil
}

func (pr *procSetObj) GetType() string {
	return "ProcSet"
}

// relateFonts is a slice of relateFont.
type relateFonts []relateFont

// IsContainsFamily checks if font family exists.
func (re *relateFonts) IsContainsFamily(family string) bool {
	for _, rf := range *re {
		if rf.Family == family {
			return true
		}
	}
	return false
}

// IsContainsFamilyAndStyle checks if font with same name and style already exists .
func (re *relateFonts) IsContainsFamilyAndStyle(family string, style int) bool {
	for _, rf := range *re {
		if rf.Family == family && rf.Style == style {
			return true
		}
	}
	return false
}

// relateFont is a metadata index for fonts?
type relateFont struct {
	Family string
	//etc /F1
	CountOfFont int
	//etc  5 0 R
	IndexOfObj int
	Style      int // Regular|Bold|Italic
}

// relateXobjects is a slice of relateXobject.
type relateXobjects []relateXobject

// relateXobject is an index for ???
type relateXobject struct {
	IndexOfObj int
}

// extGS represents an External Graphics State object in PDF.
// It stores an index reference to a graphics state parameter dictionary
// that defines visual characteristics like transparency, line width, etc.
// These objects are referenced in the extGState dictionary of the resource dictionary.
type extGS struct {
	Index int // Index identifier used to reference this graphics state in the PDF
}
