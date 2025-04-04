package docpdf

import (
	"fmt"
	"io"
	"strings"
)

// pagesObj pdf pages object
type pagesObj struct { //impl iObj
	PageCount int
	Kids      string
	getRoot   func() *pdfEngine
}

// GetNumberOfPages gets the number of pages from the PDF.
func (gp *pdfEngine) GetNumberOfPages() int {
	return gp.numOfPagesObj
}

func (p *pagesObj) init(funcGetRoot func() *pdfEngine) {
	p.PageCount = 0
	p.getRoot = funcGetRoot
}

func (p *pagesObj) write(w writer, objID int) error {

	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "  /Type /%s\n", p.getType())

	rootConfig := p.getRoot().config
	fmt.Fprintf(w, "  /MediaBox [ 0 0 %0.2f %0.2f ]\n", rootConfig.PageSize.W, rootConfig.PageSize.H)
	fmt.Fprintf(w, "  /Count %d\n", p.PageCount)
	fmt.Fprintf(w, "  /Kids [ %s ]\n", p.Kids) //sample Kids [ 3 0 R ]
	io.WriteString(w, ">>\n")
	return nil
}

func (p *pagesObj) getType() string {
	return "Pages"
}

func (p *pagesObj) test() {
	fmt.Print(p.getType() + "\n")
}

// pageOption option of page
type pageOption struct {
	TrimBox  *box
	PageSize *Rect
}

func (p pageOption) isEmpty() bool {
	return p.PageSize == nil
}

func (p pageOption) isTrimBoxSet() bool {
	if p.TrimBox == nil {
		return false
	}
	if p.TrimBox.Top == 0 && p.TrimBox.Left == 0 && p.TrimBox.Bottom == 0 && p.TrimBox.Right == 0 {
		return false
	}

	return true
}

// pageObj pdf page object
type pageObj struct { //impl iObj
	Contents        string
	ResourcesRelate string
	pageOption      pageOption
	LinkObjIds      []int
	getRoot         func() *pdfEngine
}

func (p *pageObj) init(funcGetRoot func() *pdfEngine) {
	p.getRoot = funcGetRoot
	p.LinkObjIds = make([]int, 0)
}

func (p *pageObj) setOption(opt pageOption) {
	p.pageOption = opt
}

func (p *pageObj) write(w writer, objID int) error {
	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "  /Type /%s\n", p.getType())
	io.WriteString(w, "  /Parent 2 0 R\n")
	fmt.Fprintf(w, "  /Resources %s\n", p.ResourcesRelate)

	var err error
	if len(p.LinkObjIds) > 0 {
		io.WriteString(w, "  /Annots [")
		for _, l := range p.LinkObjIds {
			_, err = fmt.Fprintf(w, "%d 0 R ", l)
			if err != nil {
				return err
			}
		}
		io.WriteString(w, "]\n")
	}

	/*me.buffer.WriteString("    /Font <<\n")
	i := 0
	max := len(me.Realtes)
	for i < max {
		realte := me.Realtes[i]
		me.buffer.WriteString(fmt.Sprintf("      /F%d %d 0 R\n",realte.CountOfFont + 1, realte.IndexOfObj + 1))
		i++
	}
	me.buffer.WriteString("    >>\n")*/
	//me.buffer.WriteString("  >>\n")
	fmt.Fprintf(w, "  /Contents %s\n", p.Contents) //sample  Contents 8 0 R
	if !p.pageOption.isEmpty() {
		fmt.Fprintf(w, " /MediaBox [ 0 0 %0.2f %0.2f ]\n", p.pageOption.PageSize.W, p.pageOption.PageSize.H)
	}
	if p.pageOption.isTrimBoxSet() {
		trimBox := p.pageOption.TrimBox
		fmt.Fprintf(w, " /TrimBox [ %0.2f %0.2f %0.2f %0.2f ]\n", trimBox.Left, trimBox.Top, trimBox.Right, trimBox.Bottom)
	}
	io.WriteString(w, ">>\n")
	return nil
}

func (p *pageObj) writeExternalLink(w writer, l linkOption, objID int) error {
	protection := p.getRoot().protection()
	url := l.url
	if protection != nil {
		tmp, err := rc4Cip(protection.objectkey(objID), []byte(url))
		if err != nil {
			return err
		}
		url = string(tmp)
	}
	url = strings.Replace(url, "\\", "\\\\", -1)
	url = strings.Replace(url, "(", "\\(", -1)
	url = strings.Replace(url, ")", "\\)", -1)
	url = strings.Replace(url, "\r", "\\r", -1)

	_, err := fmt.Fprintf(w, "<</Type /Annot /Subtype /Link /Rect [%.2f %.2f %.2f %.2f] /Border [0 0 0] /A <</S /URI /URI (%s)>>>>",
		l.x, l.y, l.x+l.w, l.y-l.h, url)
	return err
}

func (p *pageObj) writeInternalLink(w writer, l linkOption, anchors map[string]anchorOption) error {
	a, ok := anchors[l.anchor]
	if !ok {
		return nil
	}
	_, err := fmt.Fprintf(w, "<</Type /Annot /Subtype /Link /Rect [%.2f %.2f %.2f %.2f] /Border [0 0 0] /Dest [%d 0 R /XYZ 0 %.2f null]>>",
		l.x, l.y, l.x+l.w, l.y-l.h, a.page+1, a.y)
	return err
}

func (p *pageObj) getType() string {
	return "Page"
}
