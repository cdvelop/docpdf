package docpdf

import "strconv"

type Document struct {
	*pdfEngine
	fontConfig         FontConfig
	contentAreaWidth   float64       // Width of the content area (page width - margins)
	inlineMode         bool          // Add this field to track inline element state
	lastInlineWidth    float64       // Track the width of the last inline element
	header             *headerFooter // New field for document header
	footer             *headerFooter // New field for document footer
	inHeaderFooterDraw bool          // Flag to prevent recursion in header/footer drawing
	lastTableHeaders   []string      // Store the last table headers for width type verification
}

// NewDocument creates a new PDF document with configurable settings
// Accepts optional configurations including:
//   - FontConfig: Custom font configuration
//   - Font: Custom font family
//   - Margins: Custom margins in millimeters (more intuitive than points)
//   - PageSize: Custom page size with desired units
//   - *Rect: Predefined page size (like PageSizeLetter, PageSizeA4, etc.)
//
// Examples:
//   - NewDocument(fmt.Println, Margins{Left: 15, Top: 10, Right: 10, Bottom: 10})
//   - NewDocument(fmt.Println, PageSize{Width: 210, Height: 297, Unit: UnitMM}) // A4 size in mm
//   - NewDocument(fmt.Println, PageSizeA4) // Using predefined page size
func NewDocument(logPrint func(a ...any), configs ...any) *Document {

	doc := &Document{
		pdfEngine: &pdfEngine{
			log: logPrint,
		},
		fontConfig:      defaultFontConfig(),
		inlineMode:      false,
		lastInlineWidth: 0,
	}

	// Default margins: 1.5 cm left, 1 cm on other sides
	leftMargin := 42.52   // 1.5 cm in points
	otherMargins := 28.35 // 1 cm in points

	// Start with default page configuration (will be overridden if PageSize is provided)
	doc.Start(config{
		PageSize: *PageSizeLetter,
	})

	// Set default margins explicitly
	doc.SetMargins(leftMargin, otherMargins, otherMargins, otherMargins)

	// Process all configurations in one place
	for _, v := range configs {
		switch v := v.(type) {
		case FontConfig:
			doc.fontConfig = v
		case Font:
			doc.fontConfig.Family = v
		case Margins:
			// Convert millimeters to points (1 mm = 72.0/25.4 points)
			doc.SetMargins(
				v.Left*(72.0/25.4),
				v.Top*(72.0/25.4),
				v.Right*(72.0/25.4),
				v.Bottom*(72.0/25.4),
			)
		case PageSize:
			// User provided a custom page size with specific units
			doc.Start(config{
				PageSize: *v.ToRect(),
			})
		case *Rect:
			// User provided a predefined page size (like PageSizeLetter, PageSizeA4, etc.)
			doc.Start(config{
				PageSize: *v,
			})
		}
	}

	err := doc.loadFonts()
	if err != nil {
		doc.log("Error loading fonts: ", err)
	}

	doc.contentAreaWidth = doc.config.PageSize.W - (doc.margins.Left + doc.margins.Right)

	// Initialize header and footer
	doc.initHeaderFooter()

	// Importante: Agregar la primera página después de inicializar el header y footer
	doc.AddPage()
	doc.setDefaultFont()

	return doc
}

// GetLineHeight returns the current line height based on the active font and size
func (doc *Document) GetLineHeight() float64 {
	// Get current font size and add some padding
	fontSize := doc.curr.FontSize
	if fontSize <= 0 {
		fontSize = doc.fontConfig.Normal.Size // Default font size as fallback
	}

	// Line height is typically 1.2 to 1.5 times the font size
	// Using 1.2 as a conservative multiplier
	return fontSize * 1.2
}

// AddPage : add new page
func (gp *pdfEngine) AddPage() {
	emptyOpt := pageOption{}
	gp.AddPageWithOption(emptyOpt)
}

// AddPageWithOption  : add new page with option
func (gp *pdfEngine) AddPageWithOption(opt pageOption) {
	opt.TrimBox = opt.TrimBox.unitsToPoints(gp.config.Unit)
	opt.PageSize = opt.PageSize.unitsToPoints(gp.config.Unit)

	page := new(pageObj)
	page.init(func() *pdfEngine {
		return gp
	})

	if !opt.isEmpty() { //use page option
		page.setOption(opt)
		gp.curr.pageSize = opt.PageSize

		if opt.isTrimBoxSet() {
			gp.curr.trimBox = opt.TrimBox
		}
	} else { //use default
		gp.curr.pageSize = &gp.config.PageSize
		gp.curr.trimBox = &gp.config.TrimBox
	}

	page.ResourcesRelate = strconv.Itoa(gp.indexOfProcSet+1) + " 0 R"
	index := gp.addObj(page)
	if gp.indexOfFirstPageObj == -1 {
		gp.indexOfFirstPageObj = index
	}
	gp.curr.IndexOfPageObj = index

	gp.numOfPagesObj++

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
func (gp *pdfEngine) SetPage(pageno int) error {
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

	return newErr("invalid page number")
}

// AddPage añade una nueva página y actualiza el contador de páginas para el header y footer
func (doc *Document) AddPage() {
	// Llamar al método subyacente de pdfEngine
	doc.pdfEngine.AddPage()

	// Actualizar el contador de página actual para header y footer
	if doc.header != nil {
		doc.header.currentPage = doc.numOfPagesObj
	}
	if doc.footer != nil {
		doc.footer.currentPage = doc.numOfPagesObj
	}
}

// AddPageWithOption añade una nueva página con opciones y actualiza el contador de páginas para el header y footer
func (doc *Document) AddPageWithOption(opt pageOption) {
	// Llamar al método subyacente de pdfEngine
	doc.pdfEngine.AddPageWithOption(opt)

	// Actualizar el contador de página actual para header y footer
	if doc.header != nil {
		doc.header.currentPage = doc.numOfPagesObj
	}
	if doc.footer != nil {
		doc.footer.currentPage = doc.numOfPagesObj
	}
}

// RedrawHeaderFooter fuerza el redibujado del encabezado y pie de página en la página actual
func (doc *Document) RedrawHeaderFooter() {
	// Guardar la posición actual
	prevX, prevY := doc.GetX(), doc.GetY()

	// Si estamos en la primera página y se han modificado las opciones de visibilidad
	if doc.numOfPagesObj == 1 {
		// Forzar el redibujado del encabezado si está configurado
		if doc.header != nil && doc.header.initialized {
			doc.header.draw()
		}

		// Forzar el redibujado del pie de página si está configurado
		if doc.footer != nil && doc.footer.initialized {
			doc.footer.draw()
		}
	}

	// Restaurar la posición
	doc.SetXY(prevX, prevY)
}

// calculateElementPosition determina la posición X de un elemento basado en su alineación y ancho
func (doc *Document) calculateElementPosition(width float64, alignment position, withPadding bool) float64 {
	// Ancho total disponible en la página (incluyendo márgenes)
	// totalWidth := doc.contentAreaWidth

	// Ancho disponible para contenido
	contentWidth := doc.contentAreaWidth - (doc.margins.Left + doc.margins.Right)

	// Padding solo si se requiere
	padding := 0.0
	if withPadding {
		padding = 10.0
		// No restamos padding del ancho disponible, solo lo aplicaremos al posicionar
	}

	// Calcular posición X basada en la alineación
	var x float64
	switch alignment {
	case Center:
		// Para centrado: margen izquierdo + mitad del espacio disponible - mitad del ancho
		x = doc.margins.Left + (contentWidth / 2) - (width / 2)
	case Right:
		// Para alineado a la derecha: posición derecha - ancho
		x = doc.contentAreaWidth - doc.margins.Right - width
	default: // Left
		// Para alineado a la izquierda: simplemente el margen izquierdo
		x = doc.margins.Left
	}

	// Aplicar padding solo a la posición, no al cálculo del ancho
	if withPadding {
		if alignment == Left {
			x += padding
		} else if alignment == Right {
			x -= padding
		}
		// Para centrado, no aplicamos padding adicional
	}

	return x
}

// returns the current PDF engine instance
func (doc *Document) PdfEngine() *pdfEngine {
	return doc.pdfEngine
}
