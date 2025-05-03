package pdfengine

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/errs"
)

// AddPage : add new page
func (gp *PdfEngine) AddPage() {
	emptyOpt := PageOption{}
	gp.AddPageWithOption(emptyOpt)
}

// AddHeader - add a header function, if present this will be automatically called by AddPage()
func (gp *PdfEngine) AddHeader(f func()) {
	gp.headerFunc = f
}

// AddFooter - add a footer function, if present this will be automatically called by AddPage()
func (gp *PdfEngine) AddFooter(f func()) {
	gp.footerFunc = f
}

// metodo que obteine la estrucutua  tamaño de la pagina
func (gp *PdfEngine) GetCurrentPageSize() *canvas.Rect {
	return gp.curr.pageSize
}

// AddPageWithOption  : add new page with option
func (gp *PdfEngine) AddPageWithOption(opt PageOption) {
	opt.TrimBox = opt.TrimBox.UnitsToPoints(gp.Config.Unit)
	opt.PageSize = opt.PageSize.UnitsToPoints(gp.Config.Unit)

	page := new(pageObj)
	page.Init(func() *PdfEngine {
		return gp
	})

	if !opt.isEmpty() { //use page option
		page.setOption(opt)
		gp.curr.pageSize = opt.PageSize

		if opt.isTrimBoxSet() {
			gp.curr.trimBox = opt.TrimBox
		}
	} else { //use default
		gp.curr.pageSize = &gp.Config.PageSize
		gp.curr.trimBox = &gp.Config.TrimBox
	}

	page.ResourcesRelate = strconv.Itoa(gp.indexOfProcSet+1) + " 0 R"
	index := gp.addObj(page)
	if gp.indexOfFirstPageObj == -1 {
		gp.indexOfFirstPageObj = index
	}
	gp.curr.IndexOfPageObj = index

	gp.NumOfPagesObj++

	//reset
	gp.indexOfContent = -1
	gp.resetCurrXY()

	if gp.headerFunc != nil {
		gp.headerFunc()
		gp.resetCurrXY()
	}

	if gp.footerFunc != nil {
		gp.footerFunc()
		gp.resetCurrXY()
	}
}

// SetPage set current page
func (gp *PdfEngine) SetPage(pageno int) error {
	var pageIndex int
	for i := 0; i < len(gp.pdfObjs); i++ {
		switch gp.pdfObjs[i].(type) {
		case *contentObj:
			pageIndex += 1
			if pageIndex == pageno {
				gp.indexOfContent = i
				return nil
			}
		}
	}

	return errs.New("invalid page number")
}

// pagesObj pdf pages object
type pagesObj struct { //impl iObj
	PageCount int
	Kids      string
	getRoot   func() *PdfEngine
}

// GetNumberOfPages gets the number of pages from the PDF.
func (gp *PdfEngine) GetNumberOfPages() int {
	return gp.NumOfPagesObj
}

func (p *pagesObj) Init(funcGetRoot func() *PdfEngine) {
	p.PageCount = 0
	p.getRoot = funcGetRoot
}

func (p *pagesObj) Write(w Writer, objID int) error {

	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "  /Type /%s\n", p.GetType())

	rootConfig := p.getRoot().Config
	fmt.Fprintf(w, "  /MediaBox [ 0 0 %0.2f %0.2f ]\n", rootConfig.PageSize.W, rootConfig.PageSize.H)
	fmt.Fprintf(w, "  /Count %d\n", p.PageCount)
	fmt.Fprintf(w, "  /Kids [ %s ]\n", p.Kids) //sample Kids [ 3 0 R ]
	io.WriteString(w, ">>\n")
	return nil
}

func (p *pagesObj) GetType() string {
	return "Pages"
}

func (p *pagesObj) test() {
	fmt.Print(p.GetType() + "\n")
}

// PageOption option of page
type PageOption struct {
	TrimBox  *canvas.Box
	PageSize *canvas.Rect
}

func (p PageOption) isEmpty() bool {
	return p.PageSize == nil
}

func (p PageOption) isTrimBoxSet() bool {
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
	PageOption      PageOption
	LinkObjIds      []int
	getRoot         func() *PdfEngine
}

func (p *pageObj) Init(funcGetRoot func() *PdfEngine) {
	p.getRoot = funcGetRoot
	p.LinkObjIds = make([]int, 0)
}

func (p *pageObj) setOption(opt PageOption) {
	p.PageOption = opt
}

func (p *pageObj) Write(w Writer, objID int) error {
	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "  /Type /%s\n", p.GetType())
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
	if !p.PageOption.isEmpty() {
		fmt.Fprintf(w, " /MediaBox [ 0 0 %0.2f %0.2f ]\n", p.PageOption.PageSize.W, p.PageOption.PageSize.H)
	}
	if p.PageOption.isTrimBoxSet() {
		trimBox := p.PageOption.TrimBox
		fmt.Fprintf(w, " /TrimBox [ %0.2f %0.2f %0.2f %0.2f ]\n", trimBox.Left, trimBox.Top, trimBox.Right, trimBox.Bottom)
	}
	io.WriteString(w, ">>\n")
	return nil
}

func (p *pageObj) writeExternalLink(w Writer, l linkOption, objID int) error {
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

	_, err := fmt.Fprintf(w, "<</Type /Annot /Subtype /Link /canvas.Rect [%.2f %.2f %.2f %.2f] /Border [0 0 0] /A <</S /URI /URI (%s)>>>>",
		l.x, l.y, l.x+l.w, l.y-l.h, url)
	return err
}

func (p *pageObj) writeInternalLink(w Writer, l linkOption, anchors map[string]anchorOption) error {
	a, ok := anchors[l.anchor]
	if !ok {
		return nil
	}
	_, err := fmt.Fprintf(w, "<</Type /Annot /Subtype /Link /canvas.Rect [%.2f %.2f %.2f %.2f] /Border [0 0 0] /Dest [%d 0 R /XYZ 0 %.2f null]>>",
		l.x, l.y, l.x+l.w, l.y-l.h, a.page+1, a.y)
	return err
}

func (p *pageObj) GetType() string {
	return "Page"
}
