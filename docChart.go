package docpdf

import (
	"io/ioutil"
	"os"

	"github.com/cdvelop/docpdf/chart"
	"github.com/cdvelop/docpdf/drawing"
)

// docChart representa un gráfico que se añadirá al documento
type docChart struct {
	doc       *Document
	width     float64
	height    float64
	keepRatio bool
	alignment position
	x, y      float64
	hasPos    bool
	inline    bool
	valign    int

	// Propiedades específicas para BarChart
	title        string
	barWidth     int
	barSpacing   int
	bars         []chart.Value
	xAxisStyle   chart.Style
	yAxisStyle   chart.Style
	background   chart.Style
	canvas       chart.Style
	colorPalette chart.ColorPalette
}

// AddBarChart crea un nuevo elemento de gráfico de barras
func (doc *Document) AddBarChart() *docChart {
	return &docChart{
		doc:        doc,
		width:      500, // Ancho predeterminado
		height:     200, // Alto predeterminado
		keepRatio:  true,
		alignment:  Left,
		barWidth:   60, // Ancho de barra predeterminado
		barSpacing: 20, // Espacio entre barras predeterminado
	}
}

// Title establece el título del gráfico
func (c *docChart) Title(title string) *docChart {
	c.title = title
	return c
}

// Width establece el ancho del gráfico y mantiene la relación de aspecto si keepRatio es true
func (c *docChart) Width(w float64) *docChart {
	c.width = w
	return c
}

// Height establece la altura del gráfico y mantiene la relación de aspecto si keepRatio es true
func (c *docChart) Height(h float64) *docChart {
	c.height = h
	return c
}

// Size establece tanto el ancho como la altura explícitamente (sin preservar la relación de aspecto)
func (c *docChart) Size(w, h float64) *docChart {
	c.width = w
	c.height = h
	c.keepRatio = false
	return c
}

// FixedPosition coloca el gráfico en coordenadas específicas
func (c *docChart) FixedPosition(x, y float64) *docChart {
	c.x = x
	c.y = y
	c.hasPos = true
	return c
}

// AlignLeft alinea el gráfico al margen izquierdo
func (c *docChart) AlignLeft() *docChart {
	c.alignment = Left
	return c
}

// AlignCenter centra el gráfico horizontalmente
func (c *docChart) AlignCenter() *docChart {
	c.alignment = Center
	return c
}

// AlignRight alinea el gráfico al margen derecho
func (c *docChart) AlignRight() *docChart {
	c.alignment = Right
	return c
}

// Inline hace que el gráfico se muestre en línea con el texto en lugar de romper a una nueva línea
func (c *docChart) Inline() *docChart {
	c.inline = true
	return c
}

// VerticalAlignTop alinea el gráfico con la parte superior de la línea de texto cuando está en línea
func (c *docChart) VerticalAlignTop() *docChart {
	c.valign = 0
	return c
}

// VerticalAlignMiddle alinea el gráfico con el medio de la línea de texto cuando está en línea
func (c *docChart) VerticalAlignMiddle() *docChart {
	c.valign = 1
	return c
}

// VerticalAlignBottom alinea el gráfico con la parte inferior de la línea de texto cuando está en línea
func (c *docChart) VerticalAlignBottom() *docChart {
	c.valign = 2
	return c
}

// BarWidth establece el ancho de las barras en el gráfico de barras
func (c *docChart) BarWidth(width int) *docChart {
	c.barWidth = width
	return c
}

// BarSpacing establece el espaciado entre barras en el gráfico de barras
func (c *docChart) BarSpacing(spacing int) *docChart {
	c.barSpacing = spacing
	return c
}

// AddBar añade una barra con un valor y etiqueta al gráfico
func (c *docChart) AddBar(value float64, label string) *docChart {
	c.bars = append(c.bars, chart.Value{
		Value: value,
		Label: label,
	})
	return c
}

// Draw renderiza el gráfico en el documento con manejo de saltos de página
func (c *docChart) Draw() error {
	// Crear un archivo temporal para almacenar la imagen del gráfico
	tmpFile, err := ioutil.TempFile("", "chart-*.png")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name()) // Eliminar el archivo temporal al finalizar

	// Crear el gráfico de barras usando go-chart
	barChart := chart.BarChart{
		Title:      c.title,
		Width:      int(c.width),
		Height:     int(c.height),
		BarWidth:   c.barWidth,
		BarSpacing: c.barSpacing,
		Bars:       c.bars,
	}

	// Aplicar estilos si están configurados
	if c.background.FillColor.A > 0 {
		barChart.Background = c.background
	}

	if c.canvas.FillColor.A > 0 {
		barChart.Canvas = c.canvas
	}

	if !c.xAxisStyle.Hidden {
		barChart.XAxis = c.xAxisStyle
	}

	if !c.yAxisStyle.Hidden {
		barChart.YAxis = chart.YAxis{
			Style: c.yAxisStyle,
		}
	}

	// Renderizar el gráfico al archivo temporal
	err = barChart.Render(chart.PNG, tmpFile)
	if err != nil {
		return err
	}
	tmpFile.Close()

	// Verificar si el gráfico cabe en la página actual
	if !c.hasPos && !c.doc.inHeaderFooterDraw {
		newY := c.doc.ensureElementFits(c.height)
		if !c.inline {
			c.doc.SetY(newY)
		}
	}

	// Determinar la posición (después de un posible salto de página)
	x := c.doc.margins.Left
	y := c.doc.GetY()

	// Aplicar alineación
	switch c.alignment {
	case Center:
		x = c.doc.margins.Left + (c.doc.contentAreaWidth-c.width)/2
	case Right:
		x = c.doc.margins.Left + c.doc.contentAreaWidth - c.width
	}

	// Si se especificó una posición fija, usarla
	if c.hasPos {
		x, y = c.x, c.y
	}

	// Ajustar la posición vertical para gráficos en línea según la alineación
	if c.inline {
		lineHeight := c.doc.GetLineHeight()
		switch c.valign {
		case 0: // Alineación superior
			// No se necesita ajuste
		case 1: // Alineación media
			y = y + (lineHeight-c.height)/2
		case 2: // Alineación inferior
			y = y + lineHeight - c.height
		default:
			// Por defecto, alineación media
			y = y + (lineHeight-c.height)/2
		}
	}

	// Crear rectángulo para la imagen
	rect := &Rect{
		W: c.width,
		H: c.height,
	}

	// Dibujar la imagen usando la instancia pdfEngine subyacente
	err = c.doc.Image(tmpFile.Name(), x, y, rect)
	if err != nil {
		return err
	}

	// Manejar actualizaciones de posición según la configuración de línea
	if c.inline {
		// Para gráficos en línea, avanzar la posición X pero mantener Y sin cambios
		c.doc.SetX(x + c.width)
		c.doc.inlineMode = true
	} else {
		// Para gráficos de bloque, avanzar la posición Y para evitar que el texto se superponga con el gráfico
		if !c.hasPos {
			c.doc.newLineBreakBasedOnDefaultFont(y + c.height)
		}
		// Restablecer la posición X al margen izquierdo ya que este es un elemento de bloque
		c.doc.SetX(c.doc.margins.Left)
		c.doc.inlineMode = false
	}

	return nil
}

// WithStyle aplica un estilo personalizado al gráfico
func (c *docChart) WithStyle(backgroundColor, barColor drawing.Color) *docChart {
	c.background = chart.Style{
		FillColor: backgroundColor,
	}

	// Usar el color proporcionado directamente en lugar de intentar crear una nueva paleta
	c.canvas = chart.Style{
		StrokeColor: barColor,
		FillColor:   barColor.WithAlpha(100), // Versión semi-transparente para relleno
	}

	return c
}

// WithAxis configura la visibilidad del eje X e Y
func (c *docChart) WithAxis(showX, showY bool) *docChart {
	if !showX {
		c.xAxisStyle.Hidden = true
	} else {
		c.xAxisStyle = chart.Shown()
	}

	if !showY {
		c.yAxisStyle.Hidden = true
	} else {
		c.yAxisStyle = chart.Shown()
	}

	return c
}
