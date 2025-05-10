package pdfengine

import (
	// NO tinygo supported
	"bytes"
	"compress/zlib" // for constants tinygo OK

	// "fmt" // NO tinygo supported

	"io"

	"os"      // tinygo OK
	"strconv" // tinygo OK
	"time"    // NO tinygo supported

	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/env"
	"github.com/cdvelop/docpdf/errs"
)

// PdfEngine : core library for generating PDF
type PdfEngine struct {
	// FileWriter function for custom file writing logic
	FileWriter env.FileWriter

	// Page canvas.Margins
	margins canvas.Margins

	pdfObjs []iObj
	Config  Config
	anchors map[string]anchorOption

	indexOfCatalogObj int

	/*--- Important obj indexes stored to reduce search loops ---*/
	// Index of pages obj
	indexOfPagesObj int

	// Number of pages obj
	NumOfPagesObj int

	// Index of first page obj
	indexOfFirstPageObj int

	// currentPdf config.Alignment
	curr currentPdf

	indexEncodingObjFonts []int
	indexOfContent        int

	// Index of procset which should be unique
	indexOfProcSet int

	// Buffer for Reader compliance
	buf bytes.Buffer

	// PDF protection
	pdfProtection   *pdfProtection
	encryptionObjID int

	// Content streams only
	compressLevel int

	// Document info
	isUseInfo bool
	info      *PdfInfo

	// Outlines/bookmarks
	outlines           *outlinesObj
	indexOfOutlinesObj int

	// Header and footer functions
	headerFunc func()
	footerFunc func()

	// gofpdi free pdf document importer
	fpdi *importer

	// Placeholder text
	placeHolderTexts map[string]([]placeHolderTextInfo)

	// Log function for debugging
	Log func(...any)
}

// metodo que retorna currentPdf
func (gp *PdfEngine) CurrentPdf() *currentPdf {
	return &gp.curr
}

const subsetFont = "SubsetFont"

// the default margin if no canvas.Margins are set
const defaultMargin = 10.0 //for backward compatible

type CropOptions struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type ImageOptions struct {
	DegreeAngle    float64
	VerticalFlip   bool
	HorizontalFlip bool
	X              float64
	Y              float64
	Rect           *canvas.Rect
	Mask           *maskOptions
	Crop           *CropOptions
	transparency   *transparency

	extGStateIndexes []int
}
type imageFromOption struct {
	Format string //jpeg,png
	X      float64
	Y      float64
	Rect   *canvas.Rect
}

type maskOptions struct {
	ImageOptions
	BBox   *[4]float64
	Holder imageHolder
}

// SetCompressLevel : set compress Level for content streams
// Possible values for level:
//
//	-2 HuffmanOnly, -1 DefaultCompression (which is level 6)
//	 0 No compression,
//	 1 fastest compression, but not very good ratio
//	 9 best compression, but slowest
func (gp *PdfEngine) SetCompressLevel(level int) {
	if level < -2 { //-2 = zlib.HuffmanOnly
		io.WriteString(os.Stderr, "compress level too small, using DefaultCompression instead\n")
		level = zlib.DefaultCompression
	} else if level > zlib.BestCompression {
		io.WriteString(os.Stderr, "compress level too big, using BestCompression instead\n")
		level = zlib.BestCompression
		return
	}
	// sanity check complete
	gp.compressLevel = level
}

// SetNoCompression : compressLevel = 0
func (gp *PdfEngine) SetNoCompression() {
	gp.compressLevel = zlib.NoCompression
}

// metodo que retorna anchors
func (gp *PdfEngine) GetAnchors() map[string]anchorOption {
	return gp.anchors
}

// metodo que retorna iObjs
func (gp *PdfEngine) GetPdfObjs() []iObj {
	return gp.pdfObjs
}

// metodo que retorna el objeto de proteccion
func (gp *PdfEngine) GetPdfProtection() *pdfProtection {
	return gp.pdfProtection
}

// Br : new line
func (gp *PdfEngine) Br(h float64) {
	gp.UnitsToPointsVar(&h)
	gp.curr.Y += h
	gp.curr.X = gp.margins.Left
}

// SetGrayFill set the grayscale for the fill, takes a float64 between 0.0 and 1.0
func (gp *PdfEngine) SetGrayFill(grayScale float64) {
	gp.curr.txtColorMode = "gray"
	gp.curr.grayFill = grayScale
	gp.getContent().AppendStreamSetGrayFill(grayScale)
}

// SetGrayStroke set the grayscale for the stroke, takes a float64 between 0.0 and 1.0
func (gp *PdfEngine) SetGrayStroke(grayScale float64) {
	gp.curr.grayStroke = grayScale
	gp.getContent().AppendStreamSetGrayStroke(grayScale)
}

func (gp *PdfEngine) AddOutline(title string) {
	gp.outlines.AddOutline(gp.curr.IndexOfPageObj+1, title)
}

// AddOutlineWithPosition add an outline with config.Alignment
func (gp *PdfEngine) AddOutlineWithPosition(title string) *outlineObj {
	return gp.outlines.AddOutlinesWithPosition(gp.curr.IndexOfPageObj+1, title, gp.Config.PageSize.H-gp.curr.Y+20)
}

// Start : init gopdf
func (gp *PdfEngine) Start(Config Config) {

	gp.start(Config)

}

func (gp *PdfEngine) StartWithImporter(Config Config, importer *importer) {

	gp.start(Config, importer)

}

func (gp *PdfEngine) start(Config Config, importer ...*importer) {

	// setup Log function
	if gp.Log == nil {
		gp.Log = env.SetupDefaultLogger()
	}

	// setup file writer
	if gp.FileWriter == nil {
		gp.FileWriter = env.SetupDefaultFileWriter()
	}

	gp.Config = Config
	gp.Init(importer...)
	//init all basic obj
	catalog := new(catalogObj)
	catalog.Init(func() *PdfEngine {
		return gp
	})
	pages := new(pagesObj)
	pages.Init(func() *PdfEngine {
		return gp
	})
	gp.outlines = new(outlinesObj)
	gp.outlines.Init(func() *PdfEngine {
		return gp
	})
	gp.indexOfCatalogObj = gp.addObj(catalog)
	gp.indexOfPagesObj = gp.addObj(pages)
	gp.indexOfOutlinesObj = gp.addObj(gp.outlines)
	gp.outlines.SetIndexObjOutlines(gp.indexOfOutlinesObj)

	//IndexOfProcSet
	procset := new(procSetObj)
	procset.Init(func() *PdfEngine {
		return gp
	})
	gp.indexOfProcSet = gp.addObj(procset)

	if gp.isUseProtection() {
		gp.pdfProtection = gp.createProtection()
	}

	gp.placeHolderTexts = make(map[string][]placeHolderTextInfo)

}

// metodo para setar IndexOfProcSet
func (gp *PdfEngine) SetIndexOfProcSet(index int) {
	gp.indexOfProcSet = index
}

// ImportPage imports a page and return template id.
// gofpdi code
func (gp *PdfEngine) ImportPage(sourceFile string, pageno int, Box string) int {
	// Set source file for fpdi
	gp.fpdi.SetSourceFile(sourceFile)

	// gofpdi needs to know where to start the object id at.
	// By default, it starts at 1, but gopdf adds a few objects initially.
	startObjID := gp.GetNextObjectID()

	// Set gofpdi next object ID to  whatever the value of startObjID is
	gp.fpdi.SetNextObjectID(startObjID)

	// Import page
	tpl := gp.fpdi.ImportPage(pageno, Box)

	// Import objects into current pdf document
	tplObjIDs := gp.fpdi.PutFormXobjects()

	// Set template names and ids in gopdf
	gp.ImportTemplates(tplObjIDs)

	// Get a map[int]string of the imported objects.
	// The map keys will be the ID of each object.
	imported := gp.fpdi.GetImportedObjects()

	// Import gofpdi objects into gopdf, starting at whatever the value of startObjID is
	gp.ImportObjects(imported, startObjID)

	// Return template ID
	return tpl
}

// ImportPageStream imports page using a stream.
// Return template id after importing.
// gofpdi code
func (gp *PdfEngine) ImportPageStream(sourceStream *io.ReadSeeker, pageno int, Box string) int {
	// Set source file for fpdi
	gp.fpdi.SetSourceStream(sourceStream)

	// gofpdi needs to know where to start the object id at.
	// By default, it starts at 1, but gopdf adds a few objects initially.
	startObjID := gp.GetNextObjectID()

	// Set gofpdi next object ID to  whatever the value of startObjID is
	gp.fpdi.SetNextObjectID(startObjID)

	// Import page
	tpl := gp.fpdi.ImportPage(pageno, Box)

	// Import objects into current pdf document
	tplObjIDs := gp.fpdi.PutFormXobjects()

	// Set template names and ids in gopdf
	gp.ImportTemplates(tplObjIDs)

	// Get a map[int]string of the imported objects.
	// The map keys will be the ID of each object.
	imported := gp.fpdi.GetImportedObjects()

	// Import gofpdi objects into gopdf, starting at whatever the value of startObjID is
	gp.ImportObjects(imported, startObjID)

	// Return template ID
	return tpl
}

// UseImportedTemplate draws an imported PDF page.
func (gp *PdfEngine) UseImportedTemplate(tplid int, x float64, y float64, w float64, h float64) {
	gp.UnitsToPointsVar(&x, &y, &w, &h)
	// Get template values to draw
	tplName, scaleX, scaleY, tX, tY := gp.fpdi.UseTemplate(tplid, x, y, w, h)
	gp.getContent().AppendStreamImportedTemplate(tplName, scaleX, scaleY, tX, tY)
}

// ImportPagesFromSource imports pages from a source pdf.
// The source can be a file path, byte slice, or (*)io.ReadSeeker.
func (gp *PdfEngine) ImportPagesFromSource(source any, Box string) error {
	switch v := source.(type) {
	case string:
		// Set source file for fpdi
		gp.fpdi.SetSourceFile(v)
	case []byte:
		// Set source stream for fpdi
		rs := io.ReadSeeker(bytes.NewReader(v))
		gp.fpdi.SetSourceStream(&rs)
	case io.ReadSeeker:
		// Set source stream for fpdi
		gp.fpdi.SetSourceStream(&v)
	case *io.ReadSeeker:
		// Set source stream for fpdi
		gp.fpdi.SetSourceStream(v)
	default:
		return errs.New("source type not supported")
	}

	// Get number of pages from source file
	pages := gp.fpdi.GetNumPages()

	// Get page sizes from source file
	sizes := gp.fpdi.GetPageSizes()

	for i := 0; i < pages; i++ {
		pageno := i + 1

		// Get the size of the page
		size, ok := sizes[pageno][Box]
		if !ok {
			return errs.New("can not get page size")
		}

		// Add a new page to the document
		gp.AddPage()

		// gofpdi needs to know where to start the object id at.
		// By default, it starts at 1, but gopdf adds a few objects initially.
		startObjID := gp.GetNextObjectID()

		// Set gofpdi next object ID to  whatever the value of startObjID is
		gp.fpdi.SetNextObjectID(startObjID)

		// Import page
		tpl := gp.fpdi.ImportPage(pageno, Box)

		// Import objects into current pdf document
		tplObjIDs := gp.fpdi.PutFormXobjects()

		// Set template names and ids in gopdf
		gp.ImportTemplates(tplObjIDs)

		// Get a map[int]string of the imported objects.
		// The map keys will be the ID of each object.
		imported := gp.fpdi.GetImportedObjects()

		// Import gofpdi objects into gopdf, starting at whatever the value of startObjID is
		gp.ImportObjects(imported, startObjID)

		// Draws the imported template on the current page
		gp.UseImportedTemplate(tpl, 0, 0, size["w"], size["h"])
	}

	return nil
}

// GetNextObjectID gets the next object ID so that gofpdi knows where to start the object IDs.
func (gp *PdfEngine) GetNextObjectID() int {
	return len(gp.pdfObjs) + 1
}

// ImportObjects imports objects from gofpdi into current document.
func (gp *PdfEngine) ImportObjects(objs map[int]string, startObjID int) {
	for i := startObjID; i < len(objs)+startObjID; i++ {
		if objs[i] != "" {
			gp.addObj(&importedObj{Data: objs[i]})
		}
	}
}

// ImportTemplates names into procset dictionary.
func (gp *PdfEngine) ImportTemplates(tpls map[string]int) {
	procset := gp.pdfObjs[gp.indexOfProcSet].(*procSetObj)
	for tplName, tplID := range tpls {
		procset.ImportedTemplateIds[tplName] = tplID
	}
}

// AddExternalLink adds a new external link.
func (gp *PdfEngine) AddExternalLink(url string, x, y, w, h float64) {
	gp.UnitsToPointsVar(&x, &y, &w, &h)

	linkOpt := linkOption{x, gp.Config.PageSize.H - y, w, h, url, ""}
	gp.addLink(linkOpt)
}

// AddInternalLink adds a new internal link.
func (gp *PdfEngine) AddInternalLink(anchor string, x, y, w, h float64) {
	gp.UnitsToPointsVar(&x, &y, &w, &h)

	linkOpt := linkOption{x, gp.Config.PageSize.H - y, w, h, "", anchor}
	gp.addLink(linkOpt)
}

func (gp *PdfEngine) addLink(option linkOption) {
	page := gp.pdfObjs[gp.curr.IndexOfPageObj].(*pageObj)
	linkObj := gp.addObj(annotObj{option, func() *PdfEngine {
		return gp
	}})
	page.LinkObjIds = append(page.LinkObjIds, linkObj+1)
}

// SetAnchor creates a new anchor.
func (gp *PdfEngine) SetAnchor(name string) {
	y := gp.Config.PageSize.H - gp.curr.Y + float64(gp.curr.FontStyle.GetSize())/2
	gp.anchors[name] = anchorOption{gp.curr.IndexOfPageObj, y}
}

// AddTTFFontByReader adds font data by reader.
func (gp *PdfEngine) AddTTFFontData(family string, fontData []byte) error {
	return gp.AddTTFFontDataWithOption(family, fontData, defaultTtfFontOption())
}

// AddTTFFontDataWithOption adds font data with option.
func (gp *PdfEngine) AddTTFFontDataWithOption(family string, fontData []byte, option TtfOption) error {
	subsetFont := new(ttfSubsetObj)
	subsetFont.Init(func() *PdfEngine {
		return gp
	})
	subsetFont.SetTtfFontOption(option)
	subsetFont.SetFamily(family)
	err := subsetFont.SetTTFData(fontData)
	if err != nil {
		return err
	}

	return gp.setSubsetFontObject(subsetFont, family, option)
}

// KernOverride override kern value
func (gp *PdfEngine) KernOverride(family string, fn funcKernOverride) error {
	i := 0
	max := len(gp.pdfObjs)
	for i < max {
		if gp.pdfObjs[i].GetType() == subsetFont {
			obj := gp.pdfObjs[i]
			sub, ok := obj.(*ttfSubsetObj)
			if ok {
				if sub.GetFamily() == family {
					sub.funcKernOverride = fn
					return nil
				}
			}
		}
		i++
	}
	return errs.MissingFontFamily
}

func (c *currentPdf) setTextColor(color ICacheColorText) {
	c.txtColor = color
}

func (c *currentPdf) textColor() ICacheColorText {
	return c.txtColor
}

// SetTextColor :  function sets the text color
func (gp *PdfEngine) SetTextColor(r uint8, g uint8, b uint8) {
	gp.curr.txtColorMode = "color"
	rgb := cacheContentTextColorRGB{
		r: r,
		g: g,
		b: b,
	}
	gp.curr.setTextColor(rgb)
}

func (gp *PdfEngine) SetTextColorCMYK(c, m, y, k uint8) {
	gp.curr.txtColorMode = "color"
	cmyk := cacheContentTextColorCMYK{
		c: c,
		m: m,
		y: y,
		k: k,
	}
	gp.curr.setTextColor(cmyk)
}

// SetStrokeColor set the color for the stroke
func (gp *PdfEngine) SetStrokeColor(r uint8, g uint8, b uint8) {
	gp.getContent().AppendStreamSetColorStroke(r, g, b)
}

// SetFillColor set the color for the stroke
func (gp *PdfEngine) SetFillColor(r uint8, g uint8, b uint8) {
	gp.getContent().AppendStreamSetColorFill(r, g, b)
}

// SetStrokeColorCMYK set the color for the stroke in CMYK color mode
func (gp *PdfEngine) SetStrokeColorCMYK(c, m, y, k uint8) {
	gp.getContent().AppendStreamSetColorStrokeCMYK(c, m, y, k)
}

// SetFillColorCMYK set the color for the fill in CMYK color mode
func (gp *PdfEngine) SetFillColorCMYK(c, m, y, k uint8) {
	gp.getContent().AppendStreamSetColorFillCMYK(c, m, y, k)
}

// MeasureTextWidth : measure Width of text (use current font)
func (gp *PdfEngine) MeasureTextWidth(text string) (float64, error) {

	text, err := gp.curr.FontISubset.AddChars(text) //AddChars for create CharacterToGlyphIndex
	if err != nil {
		return 0, err
	}

	_, _, textWidthPdfUnit, err := CreateContent(gp.curr.FontISubset, text, gp.curr.FontStyle.GetSize(), gp.curr.CharSpacing, nil)
	if err != nil {
		return 0, err
	}
	return canvas.PointsToUnitsCfg(gp.Config, textWidthPdfUnit), nil
}

// MeasureCellHeightByText : measure Height of cell by text (use current font)
func (gp *PdfEngine) MeasureCellHeightByText(text string) (float64, error) {

	text, err := gp.curr.FontISubset.AddChars(text) //AddChars for create CharacterToGlyphIndex
	if err != nil {
		return 0, err
	}

	_, cellHeightPdfUnit, _, err := CreateContent(gp.curr.FontISubset, text, gp.curr.FontStyle.GetSize(), gp.curr.CharSpacing, nil)
	if err != nil {
		return 0, err
	}
	return canvas.PointsToUnitsCfg(gp.Config, cellHeightPdfUnit), nil
}

/*
//SetProtection set permissions as well as user and owner passwords
func (gp *PdfEngine) SetProtection(permissions int, userPass []byte, ownerPass []byte) {
	gp.pdfProtection = new(pdfProtection)
	gp.pdfProtection.setProtection(permissions, userPass, ownerPass)
}*/

// SetInfo set Document Information Dictionary
func (gp *PdfEngine) SetInfo(info PdfInfo) {
	gp.info = &info
	gp.isUseInfo = true
}

// GetInfo get Document Information Dictionary
func (gp *PdfEngine) GetInfo() *PdfInfo {
	return gp.info
}

// Rotate rotate text or image
// angle is angle in degrees.
// x, y is rotation center
func (gp *PdfEngine) Rotate(angle, x, y float64) {
	gp.UnitsToPointsVar(&x, &y)
	gp.getContent().appendRotate(angle, x, y)
}

// RotateReset reset rotate
func (gp *PdfEngine) RotateReset() {
	gp.getContent().appendRotateReset()
}

/*---private---*/

// init
func (gp *PdfEngine) Init(importer ...*importer) {
	gp.pdfObjs = []iObj{}
	gp.buf = bytes.Buffer{}
	gp.indexEncodingObjFonts = []int{}
	gp.pdfProtection = nil
	gp.encryptionObjID = 0
	gp.isUseInfo = false
	gp.info = nil

	//default
	gp.margins = canvas.Margins{
		Left:   defaultMargin,
		Top:    defaultMargin,
		Right:  defaultMargin,
		Bottom: defaultMargin,
	}

	//init curr
	gp.resetCurrXY()
	gp.curr = currentPdf{}
	gp.curr.IndexOfPageObj = -1
	gp.curr.CountOfFont = 0
	gp.curr.CountOfL = 0
	gp.curr.CountOfImg = 0                       //img
	gp.curr.ImgCaches = make(map[int]imageCache) //= *new([]imageCache)
	gp.curr.sMasksMap = newSMaskMap()
	gp.curr.extGStatesMap = newExtGStatesMap()
	gp.curr.transparencyMap = newTransparencyMap()
	gp.anchors = make(map[string]anchorOption)
	gp.curr.txtColorMode = "gray"

	//init index
	gp.indexOfPagesObj = -1
	gp.indexOfFirstPageObj = -1
	gp.indexOfContent = -1

	//No underline
	//gp.IsUnderline = false
	gp.curr.lineWidth = 1

	// default to zlib.DefaultCompression
	gp.compressLevel = zlib.DefaultCompression

	// change the unit type
	gp.Config.PageSize = *gp.Config.PageSize.UnitsToPoints(gp.Config)
	gp.Config.TrimBox = *gp.Config.TrimBox.UnitsToPoints(gp.Config)

	// init gofpdi free pdf document importer
	gp.fpdi = gp.importerOrDefault(importer...)

}

func (gp *PdfEngine) importerOrDefault(importer ...*importer) *importer {
	if len(importer) != 0 {
		return importer[len(importer)-1]
	}

	return newImporter(gp.Log)
}

func (gp *PdfEngine) resetCurrXY() {
	gp.curr.X = gp.margins.Left
	gp.curr.Y = gp.margins.Top
}

// UnitsToPoints converts the units to the documents unit type
func (gp *PdfEngine) UnitsToPoints(u float64) float64 {
	return canvas.UnitsToPoints(gp.Config, u)
}

// UnitsToPointsVar converts the units to the documents unit type for all variables passed in
func (gp *PdfEngine) UnitsToPointsVar(u ...*float64) {
	canvas.UnitsToPointsVar(gp.Config, u...)
}

// pointsToUnits converts the points to the documents unit type
func (gp *PdfEngine) pointsToUnits(u float64) float64 {
	return canvas.PointsToUnits(gp.Config, u)
}

// PointsToUnitsVar converts the points to the documents unit type for all variables passed in
func (gp *PdfEngine) PointsToUnitsVar(u ...*float64) {
	canvas.PointsToUnitsVarCfg(gp.Config, u...)
}

func (gp *PdfEngine) isUseProtection() bool {
	return gp.Config.Protection.UseProtection
}

func (gp *PdfEngine) createProtection() *pdfProtection {
	var prot pdfProtection
	prot.setProtection(
		gp.Config.Protection.Permissions,
		gp.Config.Protection.UserPass,
		gp.Config.Protection.OwnerPass,
	)
	return &prot
}

func (gp *PdfEngine) protection() *pdfProtection {
	return gp.pdfProtection
}

func (gp *PdfEngine) prepare() {

	if gp.isUseProtection() {
		encObj := gp.pdfProtection.encryptionObj()
		gp.addObj(encObj)
	}

	if gp.outlines.Count() > 0 {
		catalogObj := gp.pdfObjs[gp.indexOfCatalogObj].(*catalogObj)
		catalogObj.SetIndexObjOutlines(gp.indexOfOutlinesObj)
	}

	if gp.indexOfPagesObj != -1 {
		indexCurrPage := -1
		pagesObj := gp.pdfObjs[gp.indexOfPagesObj].(*pagesObj)
		i := 0 //gp.indexOfFirstPageObj
		max := len(gp.pdfObjs)
		for i < max {
			objtype := gp.pdfObjs[i].GetType()
			switch objtype {
			case "Page":
				pagesObj.Kids = pagesObj.Kids + strconv.Itoa(i+1) + " 0 R "
				pagesObj.PageCount++
				indexCurrPage = i
			case "Content":
				if indexCurrPage != -1 {
					gp.pdfObjs[indexCurrPage].(*pageObj).Contents = gp.pdfObjs[indexCurrPage].(*pageObj).Contents + strconv.Itoa(i+1) + " 0 R "
				}
			case "Font":
				tmpfont := gp.pdfObjs[i].(*fontObj)
				j := 0
				jmax := len(gp.indexEncodingObjFonts)
				for j < jmax {
					tmpencoding := gp.pdfObjs[gp.indexEncodingObjFonts[j]].(*encodingObj).GetFont()
					if tmpfont.Family == tmpencoding.GetFamily() { //ใส่ ข้อมูลของ embed font
						tmpfont.IsEmbedFont = true
						tmpfont.SetIndexObjEncoding(gp.indexEncodingObjFonts[j] + 1)
						tmpfont.SetIndexObjWidth(gp.indexEncodingObjFonts[j] + 2)
						tmpfont.SetIndexObjFontDescriptor(gp.indexEncodingObjFonts[j] + 3)
						break
					}
					j++
				}
			case "Encryption":
				gp.encryptionObjID = i + 1
			}
			i++
		}
	}
}

func (gp *PdfEngine) xref(w Writer, xrefbyteoffset int64, linelens []int64, i int) error {
	io.WriteString(w, "xref\n")
	io.WriteString(w, "0 "+strconv.Itoa(i+1)+"\n")
	io.WriteString(w, "0000000000 65535 f \n")
	j := 0
	max := len(linelens)
	for j < max {
		linelen := linelens[j]
		io.WriteString(w, gp.formatXrefline(linelen))
		io.WriteString(w, " 00000 n \n")
		j++
	}
	io.WriteString(w, "trailer\n")
	io.WriteString(w, "<<\n")
	io.WriteString(w, "/Size "+strconv.Itoa(max+1)+"\n")
	io.WriteString(w, "/Root 1 0 R\n")
	if gp.isUseProtection() {
		io.WriteString(w, "/Encrypt "+strconv.Itoa(gp.encryptionObjID)+" 0 R\n")
		io.WriteString(w, "/ID [()()]\n")
	}
	if gp.isUseInfo {
		gp.writeInfo(w)
	}
	io.WriteString(w, ">>\n")
	io.WriteString(w, "startxref\n")
	io.WriteString(w, strconv.FormatInt(xrefbyteoffset, 10))
	io.WriteString(w, "\n%%EOF\n")

	return nil
}

func (gp *PdfEngine) writeInfo(w Writer) {
	var zerotime time.Time
	io.WriteString(w, "/Info <<\n")

	if gp.info.Author != "" {
		io.WriteString(w, "/Author <FEFF"+encodeUtf8(gp.info.Author)+">\n")
	}

	if gp.info.Title != "" {
		io.WriteString(w, "/Title <FEFF"+encodeUtf8(gp.info.Title)+">\n")
	}

	if gp.info.Subject != "" {
		io.WriteString(w, "/Subject <FEFF"+encodeUtf8(gp.info.Subject)+">\n")
	}

	if gp.info.Creator != "" {
		io.WriteString(w, "/Creator <FEFF"+encodeUtf8(gp.info.Creator)+">\n")
	}

	if gp.info.Producer != "" {
		io.WriteString(w, "/Producer <FEFF"+encodeUtf8(gp.info.Producer)+">\n")
	}

	if !zerotime.Equal(gp.info.CreationDate) {
		io.WriteString(w, "/CreationDate(D:"+infodate(gp.info.CreationDate)+")\n")
	}

	io.WriteString(w, " >>\n")
}

// ปรับ xref ให้เป็น 10 หลัก
func (gp *PdfEngine) formatXrefline(n int64) string {
	str := strconv.FormatInt(n, 10)
	for len(str) < 10 {
		str = "0" + str
	}
	return str
}

func (gp *PdfEngine) addObj(iobj iObj) int {
	index := len(gp.pdfObjs)
	gp.pdfObjs = append(gp.pdfObjs, iobj)
	return index
}

func (gp *PdfEngine) getContent() *contentObj {
	var content *contentObj
	if gp.indexOfContent <= -1 {
		content = new(contentObj)
		content.Init(func() *PdfEngine {
			return gp
		})
		gp.indexOfContent = gp.addObj(content)
	} else {
		content = gp.pdfObjs[gp.indexOfContent].(*contentObj)
	}
	return content
}

func encodeUtf8(str string) string {
	var buff bytes.Buffer
	for _, r := range str {
		// Convertir runa a hexadecimal usando strconv en lugar de fmt
		hex := strconv.FormatInt(int64(r), 16)

		// Asegurar que tenga 4 caracteres (rellenando con ceros)
		for len(hex) < 4 {
			hex = "0" + hex
		}

		// Convertir a mayúsculas manualmente y añadir al buffer
		for i := 0; i < len(hex); i++ {
			c := hex[i]
			if c >= 'a' && c <= 'f' {
				buff.WriteByte(c - 32) // 'a'-'A' = 32 en ASCII
			} else {
				buff.WriteByte(c)
			}
		}
	}
	return buff.String()
}

func infodate(t time.Time) string {
	ft := t.Format("20060102150405-07'00'")
	return ft
}

// SetTransparency sets transparency.
// alpha: 		value from 0 (transparent) to 1 (opaque)
// blendMode:   blend mode, one of the following:
//
//	Normal, multiply, screen, overlay, darken, lighten, colorDodge, colorBurn,
//	hardLight, softLight, difference, exclusion, hue, saturation, Color, luminosity
func (gp *PdfEngine) SetTransparency(transparency transparency) error {
	t, err := gp.saveTransparency(&transparency)
	if err != nil {
		return err
	}

	gp.curr.transparency = t

	return nil
}

func (gp *PdfEngine) ClearTransparency() {
	gp.curr.transparency = nil
}

func (gp *PdfEngine) getCachedTransparency(transparency *transparency) (*transparency, error) {
	if transparency == nil {
		transparency = gp.curr.transparency
	} else {
		cached, err := gp.saveTransparency(transparency)
		if err != nil {
			return nil, err
		}

		transparency = cached
	}

	return transparency, nil
}

func (gp *PdfEngine) saveTransparency(transparency *transparency) (*transparency, error) {
	cached, ok := gp.curr.transparencyMap.Find(*transparency)
	if ok {
		return &cached, nil
	} else if transparency.Alpha != defaultAplhaValue {
		bm := transparency.blendModeType
		opts := extGStateOptions{
			BlendMode:     &bm,
			StrokingCA:    &transparency.Alpha,
			NonStrokingCa: &transparency.Alpha,
		}

		extGState, err := getCachedExtGState(opts, gp)
		if err != nil {
			return nil, err
		}

		transparency.extGStateIndex = extGState.Index + 1

		gp.curr.transparencyMap.Save(*transparency)

		return transparency, nil
	}

	return nil, nil
}

// IsCurrFontContainGlyph defines is current font contains to a glyph
// r:           any rune
func (gp *PdfEngine) IsCurrFontContainGlyph(r rune) (bool, error) {
	fontISubset := gp.curr.FontISubset
	if fontISubset == nil {
		return false, nil
	}

	glyphIndex, err := fontISubset.CharCodeToGlyphIndex(r)
	if err == errGlyphNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if glyphIndex == 0 {
		return false, nil
	}

	return true, nil
}
