package docpdf

import (
	"github.com/cdvelop/docpdf/chart"
	"github.com/cdvelop/docpdf/config"
)

// drawWithPdfRenderer renderiza directamente al PDF usando PdfRenderer
// Implementación extremadamente simplificada para pruebas de concepto.
// No utiliza formateadores de etiquetas ni otros elementos complejos
// para enfocarse exclusivamente en probar el renderizado directo a PDF.
func (c *docChart) drawWithPdfRenderer() error {
	// Si no hay valores, no hay nada que hacer
	if len(c.values) == 0 {
		return nil
	}
	// Obtener el motor PDF
	pdfEngine := c.doc.PdfEngine // Assumes Document has a public PdfEngine field
	if pdfEngine == nil {
		return nil // Or handle error appropriately
	}

	// Calcular la posición donde dibujar
	x := c.x
	y := c.y

	if !c.hasPos {
		// Get current Y from the document, similar to docImage.go
		y = c.doc.GetY() // Assumes c.doc.GetY() exists

		// Calculate X based on alignment, using document properties
		// Assumes c.doc.Margins() and c.doc.contentAreaWidth (or means to calculate it) exist
		docMargins := c.doc.Margins() // Assumes c.doc.Margins() returns canvas.Margins
		// contentAreaWidth can be c.doc.contentAreaWidth or calculated:
		// c.doc.Config.PageSize.W - docMargins.Left - docMargins.Right
		// For simplicity, let's assume c.doc.contentAreaWidth is available as in docImage.go
		contentAreaWidth := c.doc.contentAreaWidth

		switch c.alignment {
		case config.Left:
			x = docMargins.Left
		case config.Center:
			x = docMargins.Left + (contentAreaWidth-c.width)/2
		case config.Right:
			x = docMargins.Left + contentAreaWidth - c.width
		default: // Default to left alignment
			x = docMargins.Left
		}
	}

	// Colocar el cursor del motor PDF en la posición correcta para dibujar este elemento
	pdfEngine.SetXY(x, y)

	// Asegurarse de que la fuente regular está activa en el motor PDF
	// (El tamaño de fuente podría venir de c.Style o c.doc.currentFont if applicable)
	pdfEngine.SetFont("regular", "", 10) // Using a default size for now

	// Crear el proveedor de renderizador
	rendererProvider := chart.NewPdfRendererProvider(pdfEngine)

	// Crear un gráfico de donut muy simple con los valores originales
	donut := chart.DonutChart{
		Width:  int(c.width),
		Height: int(c.height),
		Values: c.values,
		Title:  c.title,
		// Font is handled by the PdfRenderer using the pdfEngine's current font
	}

	// Renderizar directamente al PDF
	if err := donut.Render(rendererProvider, nil); err != nil {
		return err
	}

	// Actualizar la posición del cursor principal del DOCUMENTO después de dibujar el gráfico
	// (asumiendo el gráfico es un elemento de bloque)
	// Similar a docImage.go for block elements.
	if !c.hasPos { // Only update document cursor if not fixed position
		c.doc.SetY(y + c.height)         // Update Y position in the document
		c.doc.SetX(c.doc.Margins().Left) // Reset X to left margin for the next element
	}

	return nil
}
