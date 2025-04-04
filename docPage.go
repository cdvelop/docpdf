package docpdf

// ensureElementFits checks if an element with the specified height will fit on the current page.
// If it doesn't fit, it adds a new page and returns the new Y position.
// Parameters:
//   - height: height of the element in document units
//   - minBottomMargin: optional minimum margin to leave at bottom of page
// Returns:
//   - positionY: the Y position where the element should be drawn
//   - newPageAdded: true if a new page was added
func (doc *Document) ensureElementFits(height float64, minBottomMargin ...float64) float64 {
	// Convert height to points (internal PDF unit)
	doc.unitsToPointsVar(&height)

	// Default minimum bottom margin
	bottomMargin := doc.margins.Bottom
	if len(minBottomMargin) > 0 && minBottomMargin[0] > 0 {
		bottomMargin = minBottomMargin[0]
		doc.unitsToPointsVar(&bottomMargin)
	}

	// Get current Y position
	currentY := doc.curr.Y

	// Calculate header/footer space if they exist
	headerSpace := 0.0
	footerSpace := 0.0

	if doc.header != nil && doc.header.initialized && (!doc.header.hideOnFirstPage || doc.numOfPagesObj > 1) {
		// Considerar tanto el tamaño de la fuente como los espaciados
		headerSpace = doc.fontConfig.PageHeader.Size +
			doc.fontConfig.PageHeader.SpaceBefore +
			doc.fontConfig.PageHeader.SpaceAfter
	}

	if doc.footer != nil && doc.footer.initialized && (!doc.footer.hideOnFirstPage || doc.numOfPagesObj > 1) {
		// Considerar tanto el tamaño de la fuente como los espaciados
		footerSpace = doc.fontConfig.PageFooter.Size +
			doc.fontConfig.PageFooter.SpaceBefore +
			doc.fontConfig.PageFooter.SpaceAfter
	}

	// Calculate available space considering header/footer
	availableSpace := doc.curr.pageSize.H - currentY - bottomMargin - (headerSpace + footerSpace)

	// Check if we need to add a page
	if height > availableSpace {
		// Guardar el estado actual de la fuente antes de añadir la página
		currentFont := doc.curr.FontFontCount
		currentFontSize := doc.curr.FontSize
		currentFontStyle := doc.curr.FontStyle
		currentFontType := doc.curr.FontType
		currentIndexOfFontObj := doc.curr.IndexOfFontObj
		currentCharSpacing := doc.curr.CharSpacing

		// Guardar el actual modo de color del texto y valor de grayFill
		currentTextMode := doc.curr.txtColorMode
		currentGrayFill := doc.curr.grayFill

		// Si hay una estructura para el color de texto actual, guardarla
		var currentTxtColor iCacheColorText
		if doc.curr.txtColor != nil {
			currentTxtColor = doc.curr.txtColor
		}

		// Añadir nueva página
		doc.AddPage()

		// Restaurar el estado de la fuente después de añadir la página
		doc.curr.FontFontCount = currentFont
		doc.curr.FontSize = currentFontSize
		doc.curr.FontStyle = currentFontStyle
		doc.curr.FontType = currentFontType
		doc.curr.IndexOfFontObj = currentIndexOfFontObj
		doc.curr.CharSpacing = currentCharSpacing

		// Restaurar el modo de color y el grayFill
		doc.curr.txtColorMode = currentTextMode
		doc.curr.grayFill = currentGrayFill

		// Restaurar el color del texto si existía
		if currentTxtColor != nil {
			doc.curr.txtColor = currentTxtColor
		}

		return doc.curr.Y // Return the top margin position of the new page
	}

	// The element fits on the current page
	return currentY
}
