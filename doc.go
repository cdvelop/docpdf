package docpdf

import (
	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/config"
	"github.com/cdvelop/docpdf/env"
	"github.com/cdvelop/docpdf/pdfengine"
)

type Document struct {
	*pdfengine.PdfEngine
	textConfig         config.TextStyles
	contentAreaWidth   float64       // Width of the content area (page width - canvas.Margins)
	inlineMode         bool          // Add this field to track inline element state
	lastInlineWidth    float64       // Track the width of the last inline element
	header             *headerFooter // New field for document header
	footer             *headerFooter // New field for document footer
	inHeaderFooterDraw bool          // Flag to prevent recursion in header/footer drawing
	lastTableHeaders   []string      // Store the last table headers for width type verification
}

// NewDocument creates a new PDF document with configurable settings
// Accepts optional configurations:
//
// Optional configurations include:
//   - config.TextStyles: Custom text styles for different sections
//   - config.FontFamily: Custom font family
//   - canvas.Margins: Custom canvas.Margins in millimeters (more intuitive than points)
//   - canvas.PageSize: Custom page size with desired units
//   - *canvas.Rect: Predefined page size (like PageSizeLetter, canvas.PageSizeA4, etc.)
//   - func(string, []byte) error: Custom file writer function (if not provided, defaults to env.SetupDefaultFileWriter())
//
// Examples:
//   - NewDocument() // Uses default file writer
//   - NewDocument(os.WriteFile) // Custom file writer
//   - NewDocument(canvas.Margins{config.Left: 15, config.Top: 10, config.Right: 10, config.Bottom: 10})
//   - NewDocument(canvas.PageSize{Width: 210, Height: 297, Unit: canvas.UnitMM}) // A4 size in mm
//   - NewDocument(canvas.PageSizeA4) // Using predefined page size
//   - NewDocument(os.WriteFile, canvas.PageSizeA4, canvas.Margins{config.Left: 20, config.Top: 10, config.Right: 20, config.Bottom: 10})
//
// For web applications:
//   - NewDocument(func(filename string, data []byte) error {
//     response.Header().Set("Content-Type", "application/pdf")
//     _, err := response.Write(data)
//     return err
//     })
func NewDocument(configs ...any) *Document {

	doc := &Document{
		PdfEngine:       &pdfengine.PdfEngine{},
		textConfig:      config.DefaultTextConfig(),
		inlineMode:      false,
		lastInlineWidth: 0,
	}

	// Default canvas.Margins: 1.5 cm left, 1 cm on other sides
	leftMargin := 42.52   // 1.5 cm in points
	otherMargins := 28.35 // 1 cm in points

	// Default page size (will be used if no canvas.PageSize is provided)
	defaultPageSize := canvas.PageSizeLetter

	// Process all configurations in one place
	for _, v := range configs {
		switch v := v.(type) {
		case env.FileWriter:
			// Set custom file writer if provided
			doc.FileWriter = v
		case func(...any):
			// Custom logger
			doc.Log = v
		case config.TextStyles:
			doc.textConfig = v
		case config.FontFamily:
			doc.textConfig.SetFontFamily(v)
		case canvas.Margins:
			// Convert millimeters to points (1 mm = 72.0/25.4 points)
			doc.SetMargins(
				v.Left*(72.0/25.4),
				v.Top*(72.0/25.4),
				v.Right*(72.0/25.4),
				v.Bottom*(72.0/25.4),
			)
		case canvas.PageSize:
			// User provided a custom page size with specific units
			// User provided a custom page size with specific units
			defaultPageSize = v.ToRect()
		case *canvas.Rect:
			// User provided a predefined page size (like PageSizeLetter, canvas.PageSizeA4, etc.)
			defaultPageSize = v
		}
	}

	// Start with default page configuration (will be overridden if canvas.PageSize is provided)
	doc.Start(pdfengine.Config{
		PageSize: *defaultPageSize,
	})

	// Set default canvas.Margins explicitly
	doc.SetMargins(leftMargin, otherMargins, otherMargins, otherMargins)

	err := doc.textConfig.LoadFonts(doc.PdfEngine)
	if err != nil {
		doc.Log("Error loading fonts: ", err)
	}

	doc.contentAreaWidth = doc.Config.PageSize.W - (doc.Margins().Left + doc.Margins().Right)

	// Initialize header and footer
	doc.initHeaderFooter()

	// Importante: Agregar la primera página después de inicializar el header y footer
	doc.AddPage()
	doc.textConfig.SetDefaultTextConfig(doc.PdfEngine)

	return doc
}

// GetLineHeight returns the current line height based on the active font and size
func (doc *Document) GetLineHeight() float64 {
	// Get current font size and add some padding
	fontSize := doc.CurrentPdf().FontSize
	if fontSize <= 0 {
		fontSize = doc.textConfig.GetNormal().Size // Default font size as fallback
	}

	// Line height is typically 1.2 to 1.5 times the font size
	// Using 1.2 as a conservative multiplier
	return fontSize * 1.2
}

// AddPage añade una nueva página y actualiza el contador de páginas para el header y footer
func (doc *Document) AddPage() {
	// Llamar al método subyacente de PdfEngine
	doc.PdfEngine.AddPage()

	// Actualizar el contador de página actual para header y footer
	if doc.header != nil {
		doc.header.currentPage = doc.NumOfPagesObj
	}
	if doc.footer != nil {
		doc.footer.currentPage = doc.NumOfPagesObj
	}

	// Respetar el SpaceAfter del encabezado para el contenido inicial de la página
	if doc.header != nil && doc.header.initialized && (!doc.header.hideOnFirstPage || doc.NumOfPagesObj > 1) {
		// Ajustar la posición Y inicial para incluir el espacio después del encabezado
		doc.SetY(doc.Margins().Top + doc.textConfig.GetPageHeader().SpaceAfter)
	}
}

// AddPageWithOption añade una nueva página con opciones y actualiza el contador de páginas para el header y footer
func (doc *Document) AddPageWithOption(opt pdfengine.PageOption) {
	// Llamar al método subyacente de PdfEngine
	doc.PdfEngine.AddPageWithOption(opt)

	// Actualizar el contador de página actual para header y footer
	if doc.header != nil {
		doc.header.currentPage = doc.NumOfPagesObj
	}
	if doc.footer != nil {
		doc.footer.currentPage = doc.NumOfPagesObj
	}

	// Respetar el SpaceAfter del encabezado para el contenido inicial de la página
	if doc.header != nil && doc.header.initialized && (!doc.header.hideOnFirstPage || doc.NumOfPagesObj > 1) {
		// Ajustar la posición Y inicial para incluir el espacio después del encabezado
		doc.SetY(doc.Margins().Top + doc.textConfig.GetPageHeader().SpaceAfter)
	}
}

// RedrawHeaderFooter fuerza el redibujado del encabezado y pie de página en la página actual
func (doc *Document) RedrawHeaderFooter() {
	// Guardar la posición actual
	prevX, prevY := doc.GetX(), doc.GetY()

	// Si estamos en la primera página y se han modificado las opciones de visibilidad
	if doc.NumOfPagesObj == 1 {
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
func (doc *Document) calculateElementPosition(width float64, align config.Alignment, withPadding bool) float64 {
	// Ancho total disponible en la página (incluyendo márgenes)
	// totalWidth := doc.contentAreaWidth

	// Ancho disponible para contenido
	contentWidth := doc.contentAreaWidth - (doc.Margins().Left + doc.Margins().Right)

	// Padding solo si se requiere
	padding := 0.0
	if withPadding {
		padding = 10.0
		// No restamos padding del ancho disponible, solo lo aplicaremos al posicionar
	}

	// Calcular posición X basada en la alineación
	var x float64
	switch align {
	case config.Center:
		// Para centrado: margen izquierdo + mitad del espacio disponible - mitad del ancho
		x = doc.Margins().Left + (contentWidth / 2) - (width / 2)
	case config.Right:
		// Para alineado a la derecha: posición derecha - ancho
		x = doc.contentAreaWidth - doc.Margins().Right - width
	default: // config.Left
		// Para alineado a la izquierda: simplemente el margen izquierdo
		x = doc.Margins().Left
	}

	// Aplicar padding solo a la posición, no al cálculo del ancho
	if withPadding {
		if align == config.Left {
			x += padding
		} else if align == config.Right {
			x -= padding
		}
		// Para centrado, no aplicamos padding adicional
	}

	return x
}
