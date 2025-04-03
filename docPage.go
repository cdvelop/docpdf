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

	// Calculate available space
	availableSpace := doc.curr.pageSize.H - currentY - bottomMargin

	// Check if we need to add a page
	if height > availableSpace {
		doc.AddPage()
		return doc.curr.Y // Return the top margin position of the new page
	}

	// The element fits on the current page
	return currentY
}
