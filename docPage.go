package docpdf

// ensureElementFits checks if an element with the specified height will fit on the current page.
// If it doesn't fit, it adds a new page and returns the new Y alignment.Alignment.
// Parameters:
//   - height: height of the element in document units
//   - minBottomMargin: optional minimum margin to leave at bottom of page
// Returns:
//   - positionY: the Y alignment.Alignment where the element should be drawn
//   - newPageAdded: true if a new page was added
func (doc *Document) ensureElementFits(height float64, minBottomMargin ...float64) float64 {
	// Convert height to points (internal PDF unit)
	doc.UnitsToPointsVar(&height)

	// Default minimum bottom margin
	bottomMargin := doc.Margins.alignment.Bottom
	if len(minBottomMargin) > 0 && minBottomMargin[0] > 0 {
		bottomMargin = minBottomMargin[0]
		doc.UnitsToPointsVar(&bottomMargin)
	}

	// Get current Y alignment.Alignment
	currentY := doc.Curr.Y

	// Calculate header/footer space if they exist
	headerSpace := 0.0
	footerSpace := 0.0

	if doc.header != nil && doc.header.initialized && (!doc.header.hideOnFirstPage || doc.NumOfPagesObj > 1) {
		// Considerar tanto el tamaño de la fuente como los espaciados
		headerSpace = doc.fontConfig.PageHeader.Size +
			doc.fontConfig.PageHeader.SpaceBefore +
			doc.fontConfig.PageHeader.SpaceAfter
	}

	if doc.footer != nil && doc.footer.initialized && (!doc.footer.hideOnFirstPage || doc.NumOfPagesObj > 1) {
		// Considerar tanto el tamaño de la fuente como los espaciados
		footerSpace = doc.fontConfig.PageFooter.Size +
			doc.fontConfig.PageFooter.SpaceBefore +
			doc.fontConfig.PageFooter.SpaceAfter
	}

	// Calculate available space considering header/footer
	availableSpace := doc.GetCurrentPageSize().H - currentY - bottomMargin - (headerSpace + footerSpace)

	// Check if we need to add a page
	if height > availableSpace {
		// Guardar el estado actual de la fuente antes de añadir la página
		currentFont := doc.Curr.FontFontCount
		currentFontSize := doc.Curr.FontSize
		currentFontStyle := doc.Curr.FontStyle
		currentFontType := doc.Curr.FontType
		currentIndexOfFontObj := doc.Curr.IndexOfFontObj
		currentCharSpacing := doc.Curr.CharSpacing

		// Guardar el actual modo de color del texto y valor de grayFill
		currentTextMode := doc.Curr.txtColorMode
		currentGrayFill := doc.Curr.grayFill

		// Si hay una estructura para el color de texto actual, guardarla
		var currentTxtColor iCacheColorText
		if doc.Curr.txtColor != nil {
			currentTxtColor = doc.Curr.txtColor
		}

		// Añadir nueva página
		doc.AddPage()

		// Restaurar el estado de la fuente después de añadir la página
		doc.Curr.FontFontCount = currentFont
		doc.Curr.FontSize = currentFontSize
		doc.Curr.FontStyle = currentFontStyle
		doc.Curr.FontType = currentFontType
		doc.Curr.IndexOfFontObj = currentIndexOfFontObj
		doc.Curr.CharSpacing = currentCharSpacing

		// Restaurar el modo de color y el grayFill
		doc.Curr.txtColorMode = currentTextMode
		doc.Curr.grayFill = currentGrayFill

		// Restaurar el color del texto si existía
		if currentTxtColor != nil {
			doc.Curr.txtColor = currentTxtColor
		}

		return doc.Curr.Y // Return the top margin alignment.Alignment of the new page
	}

	// The element fits on the current page
	return currentY
}
