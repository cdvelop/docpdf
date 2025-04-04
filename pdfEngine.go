package docpdf

import (
	// NO tinygo supported
	"bytes"
	"compress/zlib" // for constants tinygo OK

	// "fmt" // NO tinygo supported

	"io"

	"math"    // tinygo OK
	"os"      // tinygo OK
	"strconv" // tinygo OK
	"time"    // NO tinygo supported
)

const subsetFont = "SubsetFont"

// the default margin if no margins are set
const defaultMargin = 10.0 //for backward compatible

type drawableRectOptions struct {
	Rect
	X            float64
	Y            float64
	paintStyle   paintStyle
	transparency *transparency

	extGStateIndexes []int
}

type cropOptions struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type imageOptions struct {
	DegreeAngle    float64
	VerticalFlip   bool
	HorizontalFlip bool
	X              float64
	Y              float64
	Rect           *Rect
	Mask           *maskOptions
	Crop           *cropOptions
	transparency   *transparency

	extGStateIndexes []int
}
type imageFromOption struct {
	Format string //jpeg,png
	X      float64
	Y      float64
	Rect   *Rect
}

type maskOptions struct {
	imageOptions
	BBox   *[4]float64
	Holder imageHolder
}

type lineOptions struct {
	extGStateIndexes []int
}

type polygonOptions struct {
	extGStateIndexes []int
}

// SetLineWidth : set line width
func (gp *pdfEngine) SetLineWidth(width float64) {
	gp.curr.lineWidth = gp.unitsToPoints(width)
	gp.getContent().AppendStreamSetLineWidth(gp.unitsToPoints(width))
}

// SetCompressLevel : set compress Level for content streams
// Possible values for level:
//
//	-2 HuffmanOnly, -1 DefaultCompression (which is level 6)
//	 0 No compression,
//	 1 fastest compression, but not very good ratio
//	 9 best compression, but slowest
func (gp *pdfEngine) SetCompressLevel(level int) {
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
func (gp *pdfEngine) SetNoCompression() {
	gp.compressLevel = zlib.NoCompression
}

// SetLineType : set line type  ("dashed" ,"dotted")
//
//	Usage:
//	pdf.SetLineType("dashed")
//	pdf.Line(50, 200, 550, 200)
//	pdf.SetLineType("dotted")
//	pdf.Line(50, 400, 550, 400)
func (gp *pdfEngine) SetLineType(linetype string) {
	gp.getContent().AppendStreamSetLineType(linetype)
}

// SetCustomLineType : set custom line type
//
//	Usage:
//	pdf.SetCustomLineType([]float64{0.8, 0.8}, 0)
//	pdf.Line(50, 200, 550, 200)
func (gp *pdfEngine) SetCustomLineType(dashArray []float64, dashPhase float64) {
	for i := range dashArray {
		gp.unitsToPointsVar(&dashArray[i])
	}
	gp.unitsToPointsVar(&dashPhase)
	gp.getContent().AppendStreamSetCustomLineType(dashArray, dashPhase)
}

// Line : draw line
//
//	Usage:
//	pdf.SetTransparency(docpdf.transparency{Alpha: 0.5,blendModeType: docpdf.colorBurn})
//	pdf.SetLineType("dotted")
//	pdf.SetStrokeColor(255, 0, 0)
//	pdf.SetLineWidth(2)
//	pdf.Line(10, 30, 585, 30)
//	pdf.ClearTransparency()
func (gp *pdfEngine) Line(x1 float64, y1 float64, x2 float64, y2 float64) {
	gp.unitsToPointsVar(&x1, &y1, &x2, &y2)
	transparency, err := gp.getCachedTransparency(nil)
	if err != nil {
		transparency = nil
	}
	var opts = lineOptions{}
	if transparency != nil {
		opts.extGStateIndexes = append(opts.extGStateIndexes, transparency.extGStateIndex)
	}
	gp.getContent().AppendStreamLine(x1, y1, x2, y2, opts)
}

// RectFromLowerLeft : draw rectangle from lower-left corner (x, y)
func (gp *pdfEngine) RectFromLowerLeft(x float64, y float64, wdth float64, hght float64) {
	gp.unitsToPointsVar(&x, &y, &wdth, &hght)

	opts := drawableRectOptions{
		X:          x,
		Y:          y,
		paintStyle: drawPaintStyle,
		Rect:       Rect{W: wdth, H: hght},
	}

	gp.getContent().AppendStreamRectangle(opts)
}

// RectFromUpperLeft : draw rectangle from upper-left corner (x, y)
func (gp *pdfEngine) RectFromUpperLeft(x float64, y float64, wdth float64, hght float64) {
	gp.unitsToPointsVar(&x, &y, &wdth, &hght)

	opts := drawableRectOptions{
		X:          x,
		Y:          y + hght,
		paintStyle: drawPaintStyle,
		Rect:       Rect{W: wdth, H: hght},
	}

	gp.getContent().AppendStreamRectangle(opts)
}

// RectFromLowerLeftWithStyle : draw rectangle from lower-left corner (x, y)
//   - style: Style of rectangule (draw and/or fill: D, F, DF, FD)
//     D or empty string: draw. This is the default value.
//     F: fill
//     DF or FD: draw and fill
func (gp *pdfEngine) RectFromLowerLeftWithStyle(x float64, y float64, wdth float64, hght float64, style string) {
	opts := drawableRectOptions{
		X: x,
		Y: y,
		Rect: Rect{
			H: hght,
			W: wdth,
		},
		paintStyle: parseStyle(style),
	}
	gp.RectFromLowerLeftWithOpts(opts)
}

func (gp *pdfEngine) RectFromLowerLeftWithOpts(opts drawableRectOptions) error {
	gp.unitsToPointsVar(&opts.X, &opts.Y, &opts.W, &opts.H)

	imageTransparency, err := gp.getCachedTransparency(opts.transparency)
	if err != nil {
		return err
	}

	if imageTransparency != nil {
		opts.extGStateIndexes = append(opts.extGStateIndexes, imageTransparency.extGStateIndex)
	}

	gp.getContent().AppendStreamRectangle(opts)

	return nil
}

// RectFromUpperLeftWithStyle : draw rectangle from upper-left corner (x, y)
//   - style: Style of rectangule (draw and/or fill: D, F, DF, FD)
//     D or empty string: draw. This is the default value.
//     F: fill
//     DF or FD: draw and fill
func (gp *pdfEngine) RectFromUpperLeftWithStyle(x float64, y float64, wdth float64, hght float64, style string) {
	opts := drawableRectOptions{
		X: x,
		Y: y,
		Rect: Rect{
			H: hght,
			W: wdth,
		},
		paintStyle: parseStyle(style),
	}
	gp.RectFromUpperLeftWithOpts(opts)
}

func (gp *pdfEngine) RectFromUpperLeftWithOpts(opts drawableRectOptions) error {
	gp.unitsToPointsVar(&opts.X, &opts.Y, &opts.W, &opts.H)

	opts.Y += opts.H

	imageTransparency, err := gp.getCachedTransparency(opts.transparency)
	if err != nil {
		return err
	}

	if imageTransparency != nil {
		opts.extGStateIndexes = append(opts.extGStateIndexes, imageTransparency.extGStateIndex)
	}

	gp.getContent().AppendStreamRectangle(opts)

	return nil
}

// Oval : draw oval
func (gp *pdfEngine) Oval(x1 float64, y1 float64, x2 float64, y2 float64) {
	gp.unitsToPointsVar(&x1, &y1, &x2, &y2)
	gp.getContent().AppendStreamOval(x1, y1, x2, y2)
}

// Br : new line
func (gp *pdfEngine) Br(h float64) {
	gp.unitsToPointsVar(&h)
	gp.curr.Y += h
	gp.curr.X = gp.margins.Left
}

// SetGrayFill set the grayscale for the fill, takes a float64 between 0.0 and 1.0
func (gp *pdfEngine) SetGrayFill(grayScale float64) {
	gp.curr.txtColorMode = "gray"
	gp.curr.grayFill = grayScale
	gp.getContent().AppendStreamSetGrayFill(grayScale)
}

// SetGrayStroke set the grayscale for the stroke, takes a float64 between 0.0 and 1.0
func (gp *pdfEngine) SetGrayStroke(grayScale float64) {
	gp.curr.grayStroke = grayScale
	gp.getContent().AppendStreamSetGrayStroke(grayScale)
}

func (gp *pdfEngine) AddOutline(title string) {
	gp.outlines.AddOutline(gp.curr.IndexOfPageObj+1, title)
}

// AddOutlineWithPosition add an outline with position
func (gp *pdfEngine) AddOutlineWithPosition(title string) *outlineObj {
	return gp.outlines.AddOutlinesWithPosition(gp.curr.IndexOfPageObj+1, title, gp.config.PageSize.H-gp.curr.Y+20)
}

// Start : init gopdf
func (gp *pdfEngine) Start(config config) {

	gp.start(config)

}

func (gp *pdfEngine) StartWithImporter(config config, importer *importer) {

	gp.start(config, importer)

}

func (gp *pdfEngine) start(config config, importer ...*importer) {

	gp.config = config
	gp.init(importer...)
	//init all basic obj
	catalog := new(catalogObj)
	catalog.init(func() *pdfEngine {
		return gp
	})
	pages := new(pagesObj)
	pages.init(func() *pdfEngine {
		return gp
	})
	gp.outlines = new(outlinesObj)
	gp.outlines.init(func() *pdfEngine {
		return gp
	})
	gp.indexOfCatalogObj = gp.addObj(catalog)
	gp.indexOfPagesObj = gp.addObj(pages)
	gp.indexOfOutlinesObj = gp.addObj(gp.outlines)
	gp.outlines.SetIndexObjOutlines(gp.indexOfOutlinesObj)

	//indexOfProcSet
	procset := new(procSetObj)
	procset.init(func() *pdfEngine {
		return gp
	})
	gp.indexOfProcSet = gp.addObj(procset)

	if gp.isUseProtection() {
		gp.pdfProtection = gp.createProtection()
	}

	gp.placeHolderTexts = make(map[string][]placeHolderTextInfo)

}

// ImportPage imports a page and return template id.
// gofpdi code
func (gp *pdfEngine) ImportPage(sourceFile string, pageno int, box string) int {
	// Set source file for fpdi
	gp.fpdi.SetSourceFile(sourceFile)

	// gofpdi needs to know where to start the object id at.
	// By default, it starts at 1, but gopdf adds a few objects initially.
	startObjID := gp.GetNextObjectID()

	// Set gofpdi next object ID to  whatever the value of startObjID is
	gp.fpdi.SetNextObjectID(startObjID)

	// Import page
	tpl := gp.fpdi.ImportPage(pageno, box)

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
func (gp *pdfEngine) ImportPageStream(sourceStream *io.ReadSeeker, pageno int, box string) int {
	// Set source file for fpdi
	gp.fpdi.SetSourceStream(sourceStream)

	// gofpdi needs to know where to start the object id at.
	// By default, it starts at 1, but gopdf adds a few objects initially.
	startObjID := gp.GetNextObjectID()

	// Set gofpdi next object ID to  whatever the value of startObjID is
	gp.fpdi.SetNextObjectID(startObjID)

	// Import page
	tpl := gp.fpdi.ImportPage(pageno, box)

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
func (gp *pdfEngine) UseImportedTemplate(tplid int, x float64, y float64, w float64, h float64) {
	gp.unitsToPointsVar(&x, &y, &w, &h)
	// Get template values to draw
	tplName, scaleX, scaleY, tX, tY := gp.fpdi.UseTemplate(tplid, x, y, w, h)
	gp.getContent().AppendStreamImportedTemplate(tplName, scaleX, scaleY, tX, tY)
}

// ImportPagesFromSource imports pages from a source pdf.
// The source can be a file path, byte slice, or (*)io.ReadSeeker.
func (gp *pdfEngine) ImportPagesFromSource(source interface{}, box string) error {
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
		return newErr("source type not supported")
	}

	// Get number of pages from source file
	pages := gp.fpdi.GetNumPages()

	// Get page sizes from source file
	sizes := gp.fpdi.GetPageSizes()

	for i := 0; i < pages; i++ {
		pageno := i + 1

		// Get the size of the page
		size, ok := sizes[pageno][box]
		if !ok {
			return newErr("can not get page size")
		}

		// Add a new page to the document
		gp.AddPage()

		// gofpdi needs to know where to start the object id at.
		// By default, it starts at 1, but gopdf adds a few objects initially.
		startObjID := gp.GetNextObjectID()

		// Set gofpdi next object ID to  whatever the value of startObjID is
		gp.fpdi.SetNextObjectID(startObjID)

		// Import page
		tpl := gp.fpdi.ImportPage(pageno, box)

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
func (gp *pdfEngine) GetNextObjectID() int {
	return len(gp.pdfObjs) + 1
}

// ImportObjects imports objects from gofpdi into current document.
func (gp *pdfEngine) ImportObjects(objs map[int]string, startObjID int) {
	for i := startObjID; i < len(objs)+startObjID; i++ {
		if objs[i] != "" {
			gp.addObj(&importedObj{Data: objs[i]})
		}
	}
}

// ImportTemplates names into procset dictionary.
func (gp *pdfEngine) ImportTemplates(tpls map[string]int) {
	procset := gp.pdfObjs[gp.indexOfProcSet].(*procSetObj)
	for tplName, tplID := range tpls {
		procset.ImportedTemplateIds[tplName] = tplID
	}
}

// AddExternalLink adds a new external link.
func (gp *pdfEngine) AddExternalLink(url string, x, y, w, h float64) {
	gp.unitsToPointsVar(&x, &y, &w, &h)

	linkOpt := linkOption{x, gp.config.PageSize.H - y, w, h, url, ""}
	gp.addLink(linkOpt)
}

// AddInternalLink adds a new internal link.
func (gp *pdfEngine) AddInternalLink(anchor string, x, y, w, h float64) {
	gp.unitsToPointsVar(&x, &y, &w, &h)

	linkOpt := linkOption{x, gp.config.PageSize.H - y, w, h, "", anchor}
	gp.addLink(linkOpt)
}

func (gp *pdfEngine) addLink(option linkOption) {
	page := gp.pdfObjs[gp.curr.IndexOfPageObj].(*pageObj)
	linkObj := gp.addObj(annotObj{option, func() *pdfEngine {
		return gp
	}})
	page.LinkObjIds = append(page.LinkObjIds, linkObj+1)
}

// SetAnchor creates a new anchor.
func (gp *pdfEngine) SetAnchor(name string) {
	y := gp.config.PageSize.H - gp.curr.Y + float64(gp.curr.FontSize)
	gp.anchors[name] = anchorOption{gp.curr.IndexOfPageObj, y}
}

// AddTTFFontByReader adds font data by reader.
func (gp *pdfEngine) AddTTFFontData(family string, fontData []byte) error {
	return gp.AddTTFFontDataWithOption(family, fontData, defaultTtfFontOption())
}

// AddTTFFontDataWithOption adds font data with option.
func (gp *pdfEngine) AddTTFFontDataWithOption(family string, fontData []byte, option ttfOption) error {
	subsetFont := new(subsetFontObj)
	subsetFont.init(func() *pdfEngine {
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

// AddTTFFontByReader adds font file by reader.
func (gp *pdfEngine) AddTTFFontByReader(family string, rd io.Reader) error {
	return gp.AddTTFFontByReaderWithOption(family, rd, defaultTtfFontOption())
}

// AddTTFFontByReaderWithOption adds font file by reader with option.
func (gp *pdfEngine) AddTTFFontByReaderWithOption(family string, rd io.Reader, option ttfOption) error {
	subsetFont := new(subsetFontObj)
	subsetFont.init(func() *pdfEngine {
		return gp
	})
	subsetFont.SetTtfFontOption(option)
	subsetFont.SetFamily(family)
	err := subsetFont.SetTTFByReader(rd)
	if err != nil {
		return err
	}

	return gp.setSubsetFontObject(subsetFont, family, option)
}

// setSubsetFontObject sets subsetFontObj.
// The given subsetFontObj is expected to be configured in advance.
func (gp *pdfEngine) setSubsetFontObject(subsetFont *subsetFontObj, family string, option ttfOption) error {
	unicodemap := new(unicodeMap)
	unicodemap.init(func() *pdfEngine {
		return gp
	})
	unicodemap.setProtection(gp.protection())
	unicodemap.SetPtrToSubsetFontObj(subsetFont)
	unicodeindex := gp.addObj(unicodemap)

	pdfdic := new(pdfDictionaryObj)
	pdfdic.init(func() *pdfEngine {
		return gp
	})
	pdfdic.setProtection(gp.protection())
	pdfdic.SetPtrToSubsetFontObj(subsetFont)
	pdfdicindex := gp.addObj(pdfdic)

	subfontdesc := new(subfontDescriptorObj)
	subfontdesc.init(func() *pdfEngine {
		return gp
	})
	subfontdesc.SetPtrToSubsetFontObj(subsetFont)
	subfontdesc.SetIndexObjPdfDictionary(pdfdicindex)
	subfontdescindex := gp.addObj(subfontdesc)

	cidfont := new(cidFontObj)
	cidfont.init(func() *pdfEngine {
		return gp
	})
	cidfont.SetPtrToSubsetFontObj(subsetFont)
	cidfont.SetIndexObjSubfontDescriptor(subfontdescindex)
	cidindex := gp.addObj(cidfont)

	subsetFont.SetIndexObjCIDFont(cidindex)
	subsetFont.SetIndexObjUnicodeMap(unicodeindex)
	index := gp.addObj(subsetFont) //add หลังสุด

	if gp.indexOfProcSet != -1 {
		procset := gp.pdfObjs[gp.indexOfProcSet].(*procSetObj)
		if !procset.Relates.IsContainsFamilyAndStyle(family, option.Style&^Underline) {
			procset.Relates = append(procset.Relates, relateFont{Family: family, IndexOfObj: index, CountOfFont: gp.curr.CountOfFont, Style: option.Style &^ Underline})
			subsetFont.CountOfFont = gp.curr.CountOfFont
			gp.curr.CountOfFont++
		}
	}
	return nil
}

// AddTTFFontWithOption : add font file
func (gp *pdfEngine) AddTTFFontWithOption(family string, ttfpath string, option ttfOption) error {

	if _, err := os.Stat(ttfpath); os.IsNotExist(err) {
		return err
	}
	data, err := os.ReadFile(ttfpath)
	if err != nil {
		return err
	}
	rd := bytes.NewReader(data)
	return gp.AddTTFFontByReaderWithOption(family, rd, option)
}

// AddTTFFont : add font file
func (gp *pdfEngine) AddTTFFont(family string, ttfpath string) error {
	return gp.AddTTFFontWithOption(family, ttfpath, defaultTtfFontOption())
}

// KernOverride override kern value
func (gp *pdfEngine) KernOverride(family string, fn funcKernOverride) error {
	i := 0
	max := len(gp.pdfObjs)
	for i < max {
		if gp.pdfObjs[i].getType() == subsetFont {
			obj := gp.pdfObjs[i]
			sub, ok := obj.(*subsetFontObj)
			if ok {
				if sub.GetFamily() == family {
					sub.funcKernOverride = fn
					return nil
				}
			}
		}
		i++
	}
	return errMissingFontFamily
}

func (c *currentPdf) setTextColor(color iCacheColorText) {
	c.txtColor = color
}

func (c *currentPdf) textColor() iCacheColorText {
	return c.txtColor
}

// SetTextColor :  function sets the text color
func (gp *pdfEngine) SetTextColor(r uint8, g uint8, b uint8) {
	gp.curr.txtColorMode = "color"
	rgb := cacheContentTextColorRGB{
		r: r,
		g: g,
		b: b,
	}
	gp.curr.setTextColor(rgb)
}

func (gp *pdfEngine) SetTextColorCMYK(c, m, y, k uint8) {
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
func (gp *pdfEngine) SetStrokeColor(r uint8, g uint8, b uint8) {
	gp.getContent().AppendStreamSetColorStroke(r, g, b)
}

// SetFillColor set the color for the stroke
func (gp *pdfEngine) SetFillColor(r uint8, g uint8, b uint8) {
	gp.getContent().AppendStreamSetColorFill(r, g, b)
}

// SetStrokeColorCMYK set the color for the stroke in CMYK color mode
func (gp *pdfEngine) SetStrokeColorCMYK(c, m, y, k uint8) {
	gp.getContent().AppendStreamSetColorStrokeCMYK(c, m, y, k)
}

// SetFillColorCMYK set the color for the fill in CMYK color mode
func (gp *pdfEngine) SetFillColorCMYK(c, m, y, k uint8) {
	gp.getContent().AppendStreamSetColorFillCMYK(c, m, y, k)
}

// MeasureTextWidth : measure Width of text (use current font)
func (gp *pdfEngine) MeasureTextWidth(text string) (float64, error) {

	text, err := gp.curr.FontISubset.AddChars(text) //AddChars for create CharacterToGlyphIndex
	if err != nil {
		return 0, err
	}

	_, _, textWidthPdfUnit, err := createContent(gp.curr.FontISubset, text, gp.curr.FontSize, gp.curr.CharSpacing, nil)
	if err != nil {
		return 0, err
	}
	return pointsToUnitsCfg(gp.config, textWidthPdfUnit), nil
}

// MeasureCellHeightByText : measure Height of cell by text (use current font)
func (gp *pdfEngine) MeasureCellHeightByText(text string) (float64, error) {

	text, err := gp.curr.FontISubset.AddChars(text) //AddChars for create CharacterToGlyphIndex
	if err != nil {
		return 0, err
	}

	_, cellHeightPdfUnit, _, err := createContent(gp.curr.FontISubset, text, gp.curr.FontSize, gp.curr.CharSpacing, nil)
	if err != nil {
		return 0, err
	}
	return pointsToUnitsCfg(gp.config, cellHeightPdfUnit), nil
}

// Curve Draws a Bézier curve (the Bézier curve is tangent to the line between the control points at either end of the curve)
// Parameters:
// - x0, y0: Start point
// - x1, y1: Control point 1
// - x2, y2: Control point 2
// - x3, y3: End point
// - style: Style of rectangule (draw and/or fill: D, F, DF, FD)
func (gp *pdfEngine) Curve(x0 float64, y0 float64, x1 float64, y1 float64, x2 float64, y2 float64, x3 float64, y3 float64, style string) {
	gp.unitsToPointsVar(&x0, &y0, &x1, &y1, &x2, &y2, &x3, &y3)
	gp.getContent().AppendStreamCurve(x0, y0, x1, y1, x2, y2, x3, y3, style)
}

/*
//SetProtection set permissions as well as user and owner passwords
func (gp *pdfEngine) SetProtection(permissions int, userPass []byte, ownerPass []byte) {
	gp.pdfProtection = new(pdfProtection)
	gp.pdfProtection.setProtection(permissions, userPass, ownerPass)
}*/

// SetInfo set Document Information Dictionary
func (gp *pdfEngine) SetInfo(info PdfInfo) {
	gp.info = &info
	gp.isUseInfo = true
}

// GetInfo get Document Information Dictionary
func (gp *pdfEngine) GetInfo() PdfInfo {
	return *gp.info
}

// Rotate rotate text or image
// angle is angle in degrees.
// x, y is rotation center
func (gp *pdfEngine) Rotate(angle, x, y float64) {
	gp.unitsToPointsVar(&x, &y)
	gp.getContent().appendRotate(angle, x, y)
}

// RotateReset reset rotate
func (gp *pdfEngine) RotateReset() {
	gp.getContent().appendRotateReset()
}

// Polygon : draw polygon
//   - style: Style of polygon (draw and/or fill: D, F, DF, FD)
//     D or empty string: draw. This is the default value.
//     F: fill
//     DF or FD: draw and fill
//
// Usage:
//
//	 pdf.SetStrokeColor(255, 0, 0)
//		pdf.SetLineWidth(2)
//		pdf.SetFillColor(0, 255, 0)
//		pdf.Polygon([]docpdf.point{{X: 10, Y: 30}, {X: 585, Y: 200}, {X: 585, Y: 250}}, "DF")
func (gp *pdfEngine) Polygon(points []point, style string) {

	transparency, err := gp.getCachedTransparency(nil)
	if err != nil {
		transparency = nil
	}

	var opts = polygonOptions{}
	if transparency != nil {
		opts.extGStateIndexes = append(opts.extGStateIndexes, transparency.extGStateIndex)
	}

	var pointReals []point
	for _, p := range points {
		x := p.X
		y := p.Y
		gp.unitsToPointsVar(&x, &y)
		pointReals = append(pointReals, point{X: x, Y: y})
	}
	gp.getContent().AppendStreamPolygon(pointReals, style, opts)
}

// Rectangle : draw rectangle, and add radius input to make a round corner, it helps to calculate the round corner coordinates and use Polygon functions to draw rectangle
//   - style: Style of Rectangle (draw and/or fill: D, F, DF, FD)
//     D or empty string: draw. This is the default value.
//     F: fill
//     DF or FD: draw and fill
//
// Usage:
//
//	 pdf.SetStrokeColor(255, 0, 0)
//		pdf.SetLineWidth(2)
//		pdf.SetFillColor(0, 255, 0)
//		pdf.Rectangle(196.6, 336.8, 398.3, 379.3, "DF", 3, 10)
func (gp *pdfEngine) Rectangle(x0 float64, y0 float64, x1 float64, y1 float64, style string, radius float64, radiusPointNum int) error {
	if x1 <= x0 || y1 <= y0 {
		return errInvalidRectangleCoordinates
	}
	if radiusPointNum <= 0 || radius <= 0 {
		//draw rectangle without round corner
		points := []point{}
		points = append(points, point{X: x0, Y: y0})
		points = append(points, point{X: x1, Y: y0})
		points = append(points, point{X: x1, Y: y1})
		points = append(points, point{X: x0, Y: y1})
		gp.Polygon(points, style)

	} else {

		if radius > (x1-x0) || radius > (y1-y0) {
			return errInvalidRectangleCoordinates
		}

		degrees := []float64{}
		angle := float64(90) / float64(radiusPointNum+1)
		accAngle := angle
		for accAngle < float64(90) {
			degrees = append(degrees, accAngle)
			accAngle += angle
		}

		radians := []float64{}
		for _, v := range degrees {
			radians = append(radians, v*math.Pi/180)
		}

		points := []point{}
		points = append(points, point{X: x0, Y: (y0 + radius)})
		for _, v := range radians {
			offsetX := radius * math.Cos(v)
			offsetY := radius * math.Sin(v)
			x := x0 + radius - offsetX
			y := y0 + radius - offsetY
			points = append(points, point{X: x, Y: y})
		}
		points = append(points, point{X: (x0 + radius), Y: y0})

		points = append(points, point{X: (x1 - radius), Y: y0})
		for i := range radians {
			v := radians[len(radians)-1-i]
			offsetX := radius * math.Cos(v)
			offsetY := radius * math.Sin(v)
			x := x1 - radius + offsetX
			y := y0 + radius - offsetY
			points = append(points, point{X: x, Y: y})
		}
		points = append(points, point{X: x1, Y: (y0 + radius)})

		points = append(points, point{X: x1, Y: (y1 - radius)})
		for _, v := range radians {
			offsetX := radius * math.Cos(v)
			offsetY := radius * math.Sin(v)
			x := x1 - radius + offsetX
			y := y1 - radius + offsetY
			points = append(points, point{X: x, Y: y})
		}
		points = append(points, point{X: (x1 - radius), Y: y1})

		points = append(points, point{X: (x0 + radius), Y: y1})
		for i := range radians {
			v := radians[len(radians)-1-i]
			offsetX := radius * math.Cos(v)
			offsetY := radius * math.Sin(v)
			x := x0 + radius - offsetX
			y := y1 - radius + offsetY
			points = append(points, point{X: x, Y: y})
		}
		points = append(points, point{X: x0, Y: y1 - radius})

		gp.Polygon(points, style)
	}
	return nil
}

/*---private---*/

// init
func (gp *pdfEngine) init(importer ...*importer) {
	gp.pdfObjs = []iObj{}
	gp.buf = bytes.Buffer{}
	gp.indexEncodingObjFonts = []int{}
	gp.pdfProtection = nil
	gp.encryptionObjID = 0
	gp.isUseInfo = false
	gp.info = nil

	//default
	gp.margins = Margins{
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
	gp.config.PageSize = *gp.config.PageSize.unitsToPoints(gp.config)
	gp.config.TrimBox = *gp.config.TrimBox.unitsToPoints(gp.config)

	// init gofpdi free pdf document importer
	gp.fpdi = gp.importerOrDefault(importer...)

}

func (gp *pdfEngine) importerOrDefault(importer ...*importer) *importer {
	if len(importer) != 0 {
		return importer[len(importer)-1]
	}

	return newImporter(gp.log)
}

func (gp *pdfEngine) resetCurrXY() {
	gp.curr.X = gp.margins.Left
	gp.curr.Y = gp.margins.Top
}

// unitsToPoints converts the units to the documents unit type
func (gp *pdfEngine) unitsToPoints(u float64) float64 {
	return unitsToPoints(gp.config, u)
}

// unitsToPointsVar converts the units to the documents unit type for all variables passed in
func (gp *pdfEngine) unitsToPointsVar(u ...*float64) {
	unitsToPointsVar(gp.config, u...)
}

// pointsToUnits converts the points to the documents unit type
func (gp *pdfEngine) pointsToUnits(u float64) float64 {
	return pointsToUnits(gp.config, u)
}

// pointsToUnitsVar converts the points to the documents unit type for all variables passed in
func (gp *pdfEngine) pointsToUnitsVar(u ...*float64) {
	pointsToUnitsVarCfg(gp.config, u...)
}

func (gp *pdfEngine) isUseProtection() bool {
	return gp.config.Protection.UseProtection
}

func (gp *pdfEngine) createProtection() *pdfProtection {
	var prot pdfProtection
	prot.setProtection(
		gp.config.Protection.Permissions,
		gp.config.Protection.UserPass,
		gp.config.Protection.OwnerPass,
	)
	return &prot
}

func (gp *pdfEngine) protection() *pdfProtection {
	return gp.pdfProtection
}

func (gp *pdfEngine) prepare() {

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
			objtype := gp.pdfObjs[i].getType()
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

func (gp *pdfEngine) xref(w writer, xrefbyteoffset int64, linelens []int64, i int) error {
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

func (gp *pdfEngine) writeInfo(w writer) {
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
func (gp *pdfEngine) formatXrefline(n int64) string {
	str := strconv.FormatInt(n, 10)
	for len(str) < 10 {
		str = "0" + str
	}
	return str
}

func (gp *pdfEngine) addObj(iobj iObj) int {
	index := len(gp.pdfObjs)
	gp.pdfObjs = append(gp.pdfObjs, iobj)
	return index
}

func (gp *pdfEngine) getContent() *contentObj {
	var content *contentObj
	if gp.indexOfContent <= -1 {
		content = new(contentObj)
		content.init(func() *pdfEngine {
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
func (gp *pdfEngine) SetTransparency(transparency transparency) error {
	t, err := gp.saveTransparency(&transparency)
	if err != nil {
		return err
	}

	gp.curr.transparency = t

	return nil
}

func (gp *pdfEngine) ClearTransparency() {
	gp.curr.transparency = nil
}

func (gp *pdfEngine) getCachedTransparency(transparency *transparency) (*transparency, error) {
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

func (gp *pdfEngine) saveTransparency(transparency *transparency) (*transparency, error) {
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
func (gp *pdfEngine) IsCurrFontContainGlyph(r rune) (bool, error) {
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

//tool for validate pdf https://www.pdf-online.com/osa/validate.aspx
