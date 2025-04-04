package docpdf

import (
	"strconv"
)

// headerFooterContent represents content that can be placed in a header or footer
type headerFooterContent struct {
	Text           string  // Text content
	Image          string  // Image path (if it's an image)
	Width          float64 // Image width if applicable
	Height         float64 // Image height if applicable
	IsImage        bool    // Whether this is an image
	WithPage       bool    // Whether to append page number
	WithTotalPages bool    // Whether to append total pages in format "X/Y"
	PageSeparator  string  // Custom separator for page numbering (default: "/")
}

// headerFooter represents a document header or footer with left, center, and right sections
type headerFooter struct {
	doc             *Document
	Left            headerFooterContent
	Center          headerFooterContent
	Right           headerFooterContent
	FontName        string
	isHeader        bool // true for header, false for footer
	initialized     bool
	currentPage     int  // Número de página actual para mostrar en el footer
	hideOnFirstPage bool // Controla si se oculta en la primera página
}

// AddHeader - add a header function, if present this will be automatically called by AddPage()
func (gp *pdfEngine) AddHeader(f func()) {
	gp.headerFunc = f
}

// AddFooter - add a footer function, if present this will be automatically called by AddPage()
func (gp *pdfEngine) AddFooter(f func()) {
	gp.footerFunc = f
}

// initHeaderFooter initializes the document's header and footer if not already done
func (d *Document) initHeaderFooter() {
	// Initialize header if not already done
	if d.header == nil {
		d.header = &headerFooter{
			doc:             d,
			FontName:        FontRegular,
			isHeader:        true,
			currentPage:     1,    // Inicializar en 1 para la primera página
			hideOnFirstPage: true, // El encabezado no se muestra por defecto en la primera página
		}

		// Set up header callback function
		d.AddHeader(func() {
			d.header.draw()
		})
	}

	// Initialize footer if not already done
	if d.footer == nil {
		d.footer = &headerFooter{
			doc:             d,
			FontName:        FontRegular,
			isHeader:        false,
			currentPage:     1,     // Inicializar en 1 para la primera página
			hideOnFirstPage: false, // El pie de página sí se muestra por defecto en la primera página
		}

		// Set up footer callback function
		d.AddFooter(func() {
			d.footer.draw()
		})
	}
}

// draw renders the header or footer on the current page
func (hf *headerFooter) draw() {
	if !hf.initialized {
		return // Nothing to draw if not initialized
	}

	// Asegurar que siempre tengamos un número de página válido
	// Sincronizar con el contador de páginas del documento
	hf.currentPage = hf.doc.numOfPagesObj
	if hf.currentPage <= 0 {
		hf.currentPage = 1 // Garantizar que la página mínima sea 1
	}

	// Verificar si debemos omitir el dibujo en la primera página
	if hf.hideOnFirstPage && hf.currentPage == 1 {
		return // No dibujar en la primera página si está configurado para ocultarse
	}

	// Save current position and drawing settings
	prevX, prevY := hf.doc.GetX(), hf.doc.GetY()

	// Determine Y position based on whether this is a header or footer
	var y float64
	if hf.isHeader {
		// position header in top margin area
		y = hf.doc.margins.Top / 2
	} else {
		// position footer in bottom margin area
		pageHeight := hf.doc.config.PageSize.H
		y = pageHeight - (hf.doc.margins.Bottom / 2) - hf.doc.fontConfig.PageFooter.Size
	}

	// Calculate column widths (3 equal sections)
	sectionWidth := hf.doc.contentAreaWidth / 3

	// Set font for header/footer
	var fontStyle TextStyle
	if hf.isHeader {
		fontStyle = hf.doc.fontConfig.PageHeader
	} else {
		fontStyle = hf.doc.fontConfig.PageFooter
	}
	hf.doc.SetFont(hf.FontName, "", fontStyle.Size)

	// Set a flag to prevent recursion during drawing the header/footer
	inHeaderFooterDraw := hf.doc.inHeaderFooterDraw
	hf.doc.inHeaderFooterDraw = true

	defer func() {
		// Restore original position and settings when done
		hf.doc.SetXY(prevX, prevY)
		// Reset inline mode
		hf.doc.inlineMode = false
		// Reset the flag
		hf.doc.inHeaderFooterDraw = inHeaderFooterDraw
	}()

	// Skip if we're already in header/footer drawing (prevents recursion)
	if inHeaderFooterDraw {
		return
	}

	// Draw left content - también dibujar si tiene paginación configurada
	if hf.Left.Text != "" || hf.Left.IsImage || hf.Left.WithPage || hf.Left.WithTotalPages {
		x := hf.doc.margins.Left
		hf.drawContent(hf.Left, x, y, sectionWidth, Left)
	}

	// Draw center content - también dibujar si tiene paginación configurada
	if hf.Center.Text != "" || hf.Center.IsImage || hf.Center.WithPage || hf.Center.WithTotalPages {
		x := hf.doc.margins.Left + sectionWidth
		hf.drawContent(hf.Center, x, y, sectionWidth, Center)
	}

	// Draw right content - también dibujar si tiene paginación configurada
	if hf.Right.Text != "" || hf.Right.IsImage || hf.Right.WithPage || hf.Right.WithTotalPages {
		x := hf.doc.margins.Left + 2*sectionWidth
		hf.drawContent(hf.Right, x, y, sectionWidth, Right)
	}
}

// drawContent draws a single content item (text or image) in the header/footer
func (hf *headerFooter) drawContent(content headerFooterContent, x, y, width float64, align position) {
	doc := hf.doc

	if content.IsImage {
		// Handle image content
		if content.Image != "" {
			img := doc.AddImage(content.Image)

			// Set fixed size if specified
			if content.Width > 0 && content.Height > 0 {
				img.Size(content.Width, content.Height)
			} else if content.Height > 0 {
				img.Height(content.Height)
			}

			// position based on alignment
			imgX := x
			if align == Center {
				imgX += width/2 - content.Width/2
			} else if align == Right {
				imgX += width - content.Width
			}

			// Place image at fixed position
			img.FixedPosition(imgX, y)
			img.Draw()
		}
	} else {
		// Handle text content
		text := content.Text

		// Add page number if requested
		if content.WithPage {
			// Para el encabezado usamos la pagina actual, para el pie de página incrementamos primero
			currentPage := hf.currentPage
			if text != "" {
				text += " "
			}
			text += strconv.Itoa(currentPage)
		}

		// Add total pages if requested
		if content.WithTotalPages {
			// Get current page and total pages
			currentPage := hf.currentPage
			totalPages := doc.GetNumberOfPages()

			// Get separator (default to "/" if not specified)
			separator := "/"
			if content.PageSeparator != "" {
				separator = content.PageSeparator
			}

			// Add space if there's existing text
			if text != "" {
				text += " "
			}

			// Format as "X/Y" or with custom separator
			text += strconv.Itoa(currentPage) + separator + strconv.Itoa(totalPages)
		}

		// Skip if there's no content to display
		if text == "" && !content.WithPage && !content.WithTotalPages {
			return
		}

		// Create text builder
		var fontStyle TextStyle
		if hf.isHeader {
			fontStyle = doc.fontConfig.PageHeader
		} else {
			fontStyle = doc.fontConfig.PageFooter
		}
		builder := doc.newTextBuilder(text, fontStyle, hf.FontName)
		builder.positioning = fixedPosition

		// Set position and width
		builder.rect.W = width

		// Save and set position
		prevX, prevY := doc.GetX(), doc.GetY()
		doc.SetXY(x, y)

		// Set alignment
		switch align {
		case Left:
			builder.AlignLeft()
		case Center:
			builder.AlignCenter()
		case Right:
			builder.AlignRight()
		}

		// Draw the text
		builder.Draw()

		// Si es encabezado, necesitamos ajustar la posición de inicio del contenido
		// pero sin modificar directamente los márgenes del documento
		if hf.isHeader && fontStyle.SpaceAfter > 0 {
			// En lugar de modificar doc.margins.Top, creamos un espacio adicional
			// Si estamos en la primera página, aplicamos un espacio adicional después del encabezado
			if doc.numOfPagesObj == 1 {
				// Añadir espacio después del encabezado (solo en la primera página para la inicialización)
				// Usamos SetY para establecer la posición correcta para el primer elemento
				topWithOffset := doc.margins.Top + fontStyle.SpaceAfter
				doc.SetY(topWithOffset)
			}
			// Para otras páginas, este ajuste se maneja en el método AddPage
		}

		// Restore position
		doc.SetXY(prevX, prevY)
	}
}

// SetPageHeader sets the document header
func (d *Document) SetPageHeader() *headerFooter {
	d.initHeaderFooter()
	d.header.initialized = true
	return d.header
}

// SetPageFooter sets the document footer
func (d *Document) SetPageFooter() *headerFooter {
	d.initHeaderFooter()
	d.footer.initialized = true
	return d.footer
}

// SetLeftText sets the left-aligned text in the header/footer
func (hf *headerFooter) SetLeftText(text string) *headerFooter {
	hf.Left = headerFooterContent{
		Text:     text,
		IsImage:  false,
		WithPage: false,
	}
	return hf
}

// SetCenterText sets the center-aligned text in the header/footer
func (hf *headerFooter) SetCenterText(text string) *headerFooter {
	hf.Center = headerFooterContent{
		Text:     text,
		IsImage:  false,
		WithPage: false,
	}
	return hf
}

// SetRightText sets the right-aligned text in the header/footer
func (hf *headerFooter) SetRightText(text string) *headerFooter {
	hf.Right = headerFooterContent{
		Text:     text,
		IsImage:  false,
		WithPage: false,
	}
	return hf
}

// SetLeftImage sets the left-aligned image in the header/footer
func (hf *headerFooter) SetLeftImage(imagePath string, width, height float64) *headerFooter {
	hf.Left = headerFooterContent{
		Image:   imagePath,
		Width:   width,
		Height:  height,
		IsImage: true,
	}
	return hf
}

// SetCenterImage sets the center-aligned image in the header/footer
func (hf *headerFooter) SetCenterImage(imagePath string, width, height float64) *headerFooter {
	hf.Center = headerFooterContent{
		Image:   imagePath,
		Width:   width,
		Height:  height,
		IsImage: true,
	}
	return hf
}

// SetRightImage sets the right-aligned image in the header/footer
func (hf *headerFooter) SetRightImage(imagePath string, width, height float64) *headerFooter {
	hf.Right = headerFooterContent{
		Image:   imagePath,
		Width:   width,
		Height:  height,
		IsImage: true,
	}
	return hf
}

// WithPageNumber adds the page number to specific section text
func (hf *headerFooter) WithPageNumber(position position) *headerFooter {
	switch position {
	case Left:
		hf.Left.WithPage = true
	case Center:
		hf.Center.WithPage = true
	case Right:
		hf.Right.WithPage = true
	default:
		// Default to center if position is invalid
		hf.Center.WithPage = true
	}
	return hf
}

// WithPageTotal adds the page number in format "X/Y" to specific section text
// You can specify a custom separator (default is "/")
// Example: WithPageTotal(Right, " de ") will display "1 de 3"
func (hf *headerFooter) WithPageTotal(position position, separator ...string) *headerFooter {
	// Set default separator to "/"
	pageSeparator := "/"

	// If custom separator provided, use it
	if len(separator) > 0 {
		pageSeparator = separator[0]
	}

	switch position {
	case Left:
		hf.Left.WithTotalPages = true
		hf.Left.WithPage = false // Disable simple page number if using total format
		hf.Left.PageSeparator = pageSeparator
	case Center:
		hf.Center.WithTotalPages = true
		hf.Center.WithPage = false // Disable simple page number if using total format
		hf.Center.PageSeparator = pageSeparator
	case Right:
		hf.Right.WithTotalPages = true
		hf.Right.WithPage = false // Disable simple page number if using total format
		hf.Right.PageSeparator = pageSeparator
	default:
		// Default to center if position is invalid
		hf.Center.WithTotalPages = true
		hf.Center.WithPage = false // Disable simple page number if using total format
		hf.Center.PageSeparator = pageSeparator
	}
	return hf
}

// SetFont sets the font for the header/footer
func (hf *headerFooter) SetFont(fontName string) *headerFooter {
	hf.FontName = fontName
	return hf
}

// AddPageHeader adds a header to the document (legacy method for backward compatibility)
func (d *Document) AddPageHeader(text string) *docText {
	// Create text builder with header style
	builder := d.newTextBuilder(text, d.fontConfig.PageHeader, FontRegular)

	// Mark as fixed position so it doesn't trigger page breaks
	builder.positioning = fixedPosition

	// Use the new header system
	d.SetPageHeader().SetCenterText(text)

	// Return the builder for method chaining (for backward compatibility)
	return builder
}

// AddPageFooter adds a footer to the document (legacy method for backward compatibility)
func (d *Document) AddPageFooter(text string) *docText {
	// Create text builder with footer style
	builder := d.newTextBuilder(text, d.fontConfig.PageFooter, FontRegular)

	// Mark as fixed position so it doesn't trigger page breaks
	builder.positioning = fixedPosition

	// Use the new footer system
	d.SetPageFooter().SetCenterText(text)

	// Return the builder for method chaining (for backward compatibility)
	return builder
}

// WithPageNumber adds page number to the text builder (legacy method for backward compatibility)
func (dt *docText) WithPageNumber() *docText {
	// Find if this is a header or footer by comparing styles
	isHeader := dt.style.Size == dt.doc.fontConfig.PageHeader.Size
	isFooter := dt.style.Size == dt.doc.fontConfig.PageFooter.Size

	// Determine if this is a header/footer text and update the appropriate structure
	if isHeader || isFooter {
		// Initialize header/footer if needed
		dt.doc.initHeaderFooter()

		// Update the appropriate header/footer
		if isHeader {
			dt.doc.header.WithPageNumber(Center)
		} else {
			dt.doc.footer.WithPageNumber(Center)
		}
	}

	// Original logic for backward compatibility
	currentText := dt.text

	// Add page number
	if currentText != "" {
		currentText += " "
	}

	// Update text in the builder
	dt.text = currentText

	return dt
}

// ShowOnFirstPage hace que el encabezado/pie de página se muestre en la primera página
func (hf *headerFooter) ShowOnFirstPage() *headerFooter {
	hf.hideOnFirstPage = false
	// Si estamos en la primera página, redibujamos para que el cambio tenga efecto inmediato
	if hf.doc.numOfPagesObj == 1 {
		hf.doc.RedrawHeaderFooter()
	}
	return hf
}

// HideOnFirstPage hace que el encabezado/pie de página se oculte en la primera página
func (hf *headerFooter) HideOnFirstPage() *headerFooter {
	hf.hideOnFirstPage = true
	// Si estamos en la primera página, redibujamos para que el cambio tenga efecto inmediato
	if hf.doc.numOfPagesObj == 1 {
		hf.doc.RedrawHeaderFooter()
	}
	return hf
}
