package docpdf

import (
	"fmt"
	"io"
)

type transparencyXObjectGroup struct {
	Index            int
	BBox             [4]float64
	Matrix           [6]float64
	ExtGStateIndexes []int
	XObjects         []cacheContentImage

	getRoot       func() *pdfEngine
	pdfProtection *pdfProtection
}

type transparencyXObjectGroupOptions struct {
	Protection       *pdfProtection
	ExtGStateIndexes []int
	BBox             [4]float64
	XObjects         []cacheContentImage
}

func getCachedTransparencyXObjectGroup(opts transparencyXObjectGroupOptions, gp *pdfEngine) (transparencyXObjectGroup, error) {
	group := transparencyXObjectGroup{
		BBox:             opts.BBox,
		XObjects:         opts.XObjects,
		pdfProtection:    opts.Protection,
		ExtGStateIndexes: opts.ExtGStateIndexes,
	}
	group.Index = gp.addObj(group)
	group.init(func() *pdfEngine {
		return gp
	})

	return group, nil
}

func (s transparencyXObjectGroup) init(funcGetRoot func() *pdfEngine) {
	s.getRoot = funcGetRoot
}

func (s *transparencyXObjectGroup) setProtection(p *pdfProtection) {
	s.pdfProtection = p
}

func (s transparencyXObjectGroup) protection() *pdfProtection {
	return s.pdfProtection
}

func (s transparencyXObjectGroup) getType() string {
	return "XObject"
}

func (s transparencyXObjectGroup) write(w writer, objId int) error {
	streamBuff := getBuffer()
	defer putBuffer(streamBuff)

	for _, XObject := range s.XObjects {
		if err := XObject.write(streamBuff, nil); err != nil {
			return err
		}
	}

	content := "<<\n"
	content += "\t/FormType 1\n"
	content += "\t/Subtype /Form\n"
	content += fmt.Sprintf("\t/Type /%s\n", s.getType())
	content += fmt.Sprintf("\t/Matrix [1 0 0 1 0 0]\n")
	content += fmt.Sprintf("\t/BBox [%.3F %.3F %.3F %.3F]\n", s.BBox[0], s.BBox[1], s.BBox[2], s.BBox[3])
	content += "\t/Group<</CS /deviceGray /S /transparency>>\n"

	content += fmt.Sprintf("\t/Length %d\n", len(streamBuff.Bytes()))
	content += ">>\n"
	content += "stream\n"

	if _, err := io.WriteString(w, content); err != nil {
		return err
	}

	if _, err := w.Write(streamBuff.Bytes()); err != nil {
		return err
	}

	if _, err := io.WriteString(w, "endstream\n"); err != nil {
		return err
	}

	return nil
}
