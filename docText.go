package docpdf

import (
	"github.com/cdvelop/docpdf/alignment"
	"github.com/cdvelop/docpdf/canvas"
)

// FontStyle defines the available font styles
const (
	FontRegular = "regular"
	FontBold    = "bold"
	FontItalic  = "italic"
)

type elementPosition string

const (
	inlinePosition  elementPosition = "inline"  // posicionamiento en línea (inline)
	newlinePosition elementPosition = "newline" // posicionamiento por defecto (salto de línea)
	fixedPosition   elementPosition = "fixed"   //posicionamiento fijo (no se mueve con el texto)
)

// TextStyle defines text appearance properties
type TextStyle struct {
	Size        float64
	Color       RGBColor
	LineSpacing float64
	Alignment   alignment.Alignment // Uses same alignment constants as cellOption (alignment.Left, alignment.Center, alignment.Right, etc)
	SpaceBefore float64             // Space before paragraph (in points)
	SpaceAfter  float64             // Space after paragraph (in points)
}

// docText is a helper struct to build text cells
type docText struct {
	doc         *Document
	text        string
	opts        cellOption
	rect        *canvas.Rect
	style       TextStyle
	fontName    string
	fullWidth   bool            // Por defecto es false (solo usa el ancho necesario)
	positioning elementPosition // "inline", "newline" (por defecto newline)
	wordWrap    bool            // Whether to use word-wrap (true) or allow mid-word breaks (false)
}

// newTextBuilder creates a new docText with the given style
func (d *Document) newTextBuilder(text string, style TextStyle, fontName string) *docText {
	builder := &docText{
		doc:         d,
		text:        text,
		style:       style, // Store the style
		fontName:    fontName,
		fullWidth:   true,            // Por defecto usar ancho completo para mantener compatibilidad
		positioning: newlinePosition, // Por defecto es newline
		wordWrap:    true,            // Por defecto usar word wrap (no cortar palabras)
		rect: &canvas.Rect{
			W: 0, // se calcula en Draw()
			H: 0,
		},
		opts: cellOption{
			Align: style.Alignment,
			// Border:         AllBorders,
			Float:          alignment.Bottom,
			CoefLineHeight: style.LineSpacing,
		},
	}

	// Apply style
	d.SetFont(fontName, "", style.Size)
	d.SetTextColor(style.Color.R, style.Color.G, style.Color.B)

	return builder
}

// AddText crea texto normal
func (d *Document) AddText(text string) *docText {
	dt := d.newTextBuilder(text, d.fontConfig.Normal, FontRegular)
	dt.fullWidth = false // Solo para texto normal, usar ancho automático
	return dt
}

// AddHeader1 crea un encabezado nivel 1
func (d *Document) AddHeader1(text string) *docText {
	return d.newTextBuilder(text, d.fontConfig.Header1, FontBold)
}

// AddHeader2 crea un encabezado nivel 2
func (d *Document) AddHeader2(text string) *docText {
	return d.newTextBuilder(text, d.fontConfig.Header2, FontBold)
}

// AddHeader3 crea un encabezado nivel 3
func (d *Document) AddHeader3(text string) *docText {
	return d.newTextBuilder(text, d.fontConfig.Header3, FontBold)
}

// AddFootnote crea una nota al pie
func (d *Document) AddFootnote(text string) *docText {
	return d.newTextBuilder(text, d.fontConfig.Footnote, FontItalic)
}

// AddJustifiedText crea texto justificado directamente
func (d *Document) AddJustifiedText(text string) *docText {
	dt := d.AddText(text)
	return dt.Justify()
}

func (dt *docText) AlignCenter() *docText {
	dt.opts.Align = alignment.Center | alignment.Top
	return dt
}

func (dt *docText) AlignRight() *docText {
	dt.opts.Align = alignment.Right | alignment.Top
	dt.fullWidth = true
	return dt
}

func (dt *docText) AlignLeft() *docText {
	dt.opts.Align = alignment.Left | alignment.Top
	return dt
}

func (dt *docText) Justify() *docText {
	dt.opts.Align = alignment.Justify | alignment.Top
	return dt
}

func (dt *docText) WithBorder() *docText {
	dt.opts.Border = AllBorders
	return dt
}

// Métodos para cambiar el estilo de fuente
func (dt *docText) Bold() *docText {
	dt.fontName = FontBold
	dt.doc.SetFont(FontBold, "", dt.style.Size)
	return dt
}

func (dt *docText) Italic() *docText {
	dt.fontName = FontItalic
	dt.doc.SetFont(FontItalic, "", dt.style.Size)
	return dt
}

func (dt *docText) Regular() *docText {
	dt.fontName = FontRegular
	dt.doc.SetFont(FontRegular, "", dt.style.Size)
	return dt
}

// SpaceBefore adds vertical space (in font spaces)
// example: SpaceBefore(2) adds 2 spaces before the text
func (d *Document) SpaceBefore(spaces ...float64) {
	space := 1.0 // Default value is 1 space if no parameter provided
	if len(spaces) > 0 && spaces[0] > 0 {
		space = spaces[0]
	}

	// Get the current font size
	fontSize := d.Curr.FontSize
	if fontSize <= 0 {
		fontSize = d.fontConfig.Normal.Size // Default font size if none is set
	}

	// Add vertical space based on font size
	d.SetY(d.GetY() + fontSize*space)
}

// FullWidth hace que el texto ocupe todo el ancho disponible
func (dt *docText) FullWidth() *docText {
	dt.fullWidth = true
	return dt
}

// WidthPercent establece el ancho como porcentaje del ancho de página
func (dt *docText) WidthPercent(percent float64) *docText {
	if percent > 0 && percent <= 100 {
		dt.rect.W = dt.doc.contentAreaWidth * (percent / 100)
	}
	return dt
}

// Inline intenta posicionar este elemento en la misma línea que el anterior
func (dt *docText) Inline() *docText {
	dt.positioning = inlinePosition
	return dt
}

// retorna el factor de ancho para una fuente específica
func (d *Document) measureTextWidthFactor(fontName string) float64 {
	// Factor de escala para el ancho de caracteres, varía según el estilo
	var widthFactor float64 = 0.6 // Factor para fuente regular

	// Ajustar factor según el estilo de fuente
	switch fontName {
	case FontBold:
		widthFactor = 0.65 // La negrita es un poco más ancha
	case FontItalic:
		widthFactor = 0.55 // La itálica suele ser ligeramente más estrecha
	}

	return widthFactor
}

// estima el ancho mínimo necesario para el texto
func (dt *docText) minimumWidthRequiredForText() {
	// Calcular ancho necesario para el texto si no se especificó un ancho fijo
	// o si se solicitó ancho completo
	if dt.fullWidth {
		// Usar ancho completo de la página
		dt.rect.W = dt.doc.contentAreaWidth
	} else {
		// Si no se especificó un ancho, calcular el ancho mínimo necesario
		if dt.rect.W <= 0 {
			// Obtener el factor de ancho según el tipo de fuente
			widthFactor := dt.doc.measureTextWidthFactor(dt.fontName)

			// Calcular ancho considerando un factor de reducción más realista
			charWidth := dt.style.Size * widthFactor

			// Considerar longitud efectiva (algunos caracteres son más estrechos)
			// Aumentar el factor para evitar problemas con signos de puntuación
			effectiveLength := float64(len(dt.text)) * 0.85 // Reducir solo 15% en lugar de 20%

			// Add extra space for text ending with punctuation marks
			if len(dt.text) > 0 {
				lastChar := dt.text[len(dt.text)-1]
				if lastChar == ':' || lastChar == '.' || lastChar == '!' || lastChar == '?' {
					effectiveLength += 0.5 // Add extra space for punctuation marks
				}
			}

			// Calcular ancho estimado
			width := effectiveLength * charWidth

			// Añadir un pequeño margen
			width += dt.style.Size * 1.1 // Aumentar el margen un 10%

			// Si el texto es largo usar el ancho de página
			if width >= dt.doc.contentAreaWidth {
				width = dt.doc.contentAreaWidth
			}

			// Asegurar un ancho mínimo razonable
			minWidth := dt.style.Size
			if width < minWidth {
				width = minWidth
			}

			dt.rect.W = width
		}
	}
}

// Draw renders the text on the document to include page break handling
func (dt *docText) Draw() error {
	// Apply space before the paragraph
	if dt.style.SpaceBefore > 0 {
		dt.doc.SetY(dt.doc.GetY() + dt.style.SpaceBefore)
	}

	// Special handling for right-aligned inline text
	isRightAligned := (dt.opts.Align == alignment.Right || dt.opts.Align == (alignment.Right|alignment.Top))

	// Configure word wrap to prevent cutting words
	if dt.wordWrap {
		// Use breakModeIndicatorSensitive with space as break indicator to avoid mid-word breaks
		dt.opts.breakOption = &breakOption{
			Mode:           breakModeIndicatorSensitive,
			BreakIndicator: ' ',
			Separator:      "",
		}
	} else {
		// Use default break option (which may allow mid-word breaks)
		dt.opts.breakOption = &defaultBreakOption
	}

	// Handle positioning
	if dt.positioning == inlinePosition {
		// For right-aligned inline text, calculate alignment.Alignment differently
		if isRightAligned {
			// First, calculate the width needed for the text
			dt.minimumWidthRequiredForText()

			// Save current Y alignment.Alignment
			currentY := dt.doc.GetY()

			// Set X to maintain right alignment while considering page canvas.Margins
			textWidth := dt.rect.W
			dt.doc.SetX(dt.doc.canvas.Margins.alignment.Left + dt.doc.contentAreaWidth - textWidth)

			// Ensure we're at the same Y alignment.Alignment
			dt.doc.SetY(currentY)
		} else {
			// Keep current X alignment.Alignment for regular inline elements
			// If we're in inline mode, adjust available width
			if dt.doc.inlineMode && dt.doc.lastInlineWidth > 0 {
				// Calculate remaining width on the current line
				currentX := dt.doc.GetX()
				availableWidth := dt.doc.contentAreaWidth - (currentX - dt.doc.canvas.Margins.alignment.Left)

				if !dt.fullWidth {
					// For auto-width text, adjust rectangle width
					dt.minimumWidthRequiredForText()
					// Check if there's enough space
					if dt.rect.W > availableWidth {
						// Not enough space, force to next line
						dt.doc.SetX(dt.doc.canvas.Margins.alignment.Left)
						dt.doc.inlineMode = false
						dt.doc.lastInlineWidth = 0
					}
				} else {
					// For full width text, adjust the width to available space
					dt.rect.W = availableWidth
				}
			}
		}
	} else {
		// Si no es inline, siempre restauramos la posición X al margen izquierdo
		// independientemente de si el elemento anterior era inline o no
		dt.doc.SetX(dt.doc.canvas.Margins.alignment.Left)
		dt.doc.inlineMode = false
		dt.doc.lastInlineWidth = 0
	}

	// If not inline or no previous inline width, and not right-aligned inline
	if (dt.positioning != inlinePosition || dt.doc.lastInlineWidth == 0) &&
		!(dt.positioning == inlinePosition && isRightAligned) {
		dt.minimumWidthRequiredForText()
	}

	// Detect if this is a single line of text
	textSplits, err := dt.doc.SplitTextWithOption(dt.text, dt.rect.W, dt.opts.breakOption)
	if err != nil {
		return err
	}

	// Get line height in current font and size
	_, lineHeight, _, err := createContent(dt.doc.Curr.FontISubset, dt.text,
		dt.doc.Curr.FontSize, dt.doc.Curr.CharSpacing, nil)
	if err != nil {
		return err
	}

	dt.doc.pointsToUnitsVar(&lineHeight)

	// Optimization for single-line text in inline mode - use Cell instead of MultiCell for better positioning
	isSingleLine := len(textSplits) == 1

	// Calculate total height needed for all lines
	totalHeight := float64(len(textSplits)) * lineHeight

	// Set the rectangle height to accommodate all text
	dt.rect.H = totalHeight

	// Skip page break check if we're in header/footer drawing mode
	if !dt.doc.inHeaderFooterDraw {
		// Check if the text fits on current page
		newY := dt.doc.ensureElementFits(totalHeight, dt.style.SpaceAfter)
		dt.doc.SetY(newY)
	}

	// Store current X alignment.Alignment to calculate width after drawing
	startX := dt.doc.GetX()

	// Choose the appropriate drawing method based on text characteristics
	if isSingleLine && dt.positioning == inlinePosition {
		// For single-line text in inline mode, use Cell for better positioning
		err = dt.doc.CellWithOption(dt.rect, dt.text, dt.opts)
	} else {
		// For multi-line text or non-inline text, use MultiCell
		err = dt.doc.MultiCellWithOption(dt.rect, dt.text, dt.opts)
	}

	if err != nil {
		return err
	}

	// Update the last inline width if in inline mode
	if dt.positioning == inlinePosition {
		// If not right-aligned, calculate actual width used
		if !isRightAligned {
			dt.doc.lastInlineWidth = dt.doc.GetX() - startX
		} else {
			// For right-aligned text, use the text width
			dt.doc.lastInlineWidth = dt.rect.W
		}
	} else {
		dt.doc.lastInlineWidth = 0
	}

	// Update inline mode based on current element's positioning
	dt.doc.inlineMode = (dt.positioning == inlinePosition)

	// If not inline, ensure we do a proper line break
	if dt.positioning == newlinePosition {
		dt.doc.newLineBreakBasedOnDefaultFont(dt.doc.GetY())
	}

	return nil
}

func (doc *Document) newLineBreakBasedOnDefaultFont(originY float64) {
	// Reset font to regular for next text (prevents style bleed)
	doc.setDefaultFont()

	// Apply space after the paragraph based on the current active text style
	// This ensures headers have their proper spacing
	var spaceAfter float64

	// Determine which style was used based on font size
	fontSize := doc.Curr.FontSize
	if fontSize >= doc.fontConfig.Header1.Size {
		spaceAfter = doc.fontConfig.Header1.SpaceAfter
	} else if fontSize >= doc.fontConfig.Header2.Size {
		spaceAfter = doc.fontConfig.Header2.SpaceAfter
	} else if fontSize >= doc.fontConfig.Header3.Size {
		spaceAfter = doc.fontConfig.Header3.SpaceAfter
	} else if fontSize <= doc.fontConfig.Footnote.Size {
		spaceAfter = doc.fontConfig.Footnote.SpaceAfter
	} else {
		spaceAfter = doc.fontConfig.Normal.SpaceAfter
	}

	// Apply the appropriate spacing
	doc.SetY(originY + fontSize + spaceAfter)
}
