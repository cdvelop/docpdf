package docpdf

import (
	"os"

	"github.com/cdvelop/docpdf/chart"
	"github.com/cdvelop/docpdf/drawing"
	"github.com/cdvelop/docpdf/fontbridge"
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
	title      string
	barWidth   int
	barSpacing int
	bars       []chart.Value
	xAxisStyle chart.Style
	yAxisStyle chart.Style
	background chart.Style
	canvas     chart.Style

	// Propiedades para control de calidad
	dpi         float64 // Resolución del gráfico en DPI (dots per inch)
	strokeWidth float64 // Ancho de línea para los contornos
}

// AddBarChart crea un nuevo elemento de gráfico de barras
func (doc *Document) AddBarChart() *docChart {
	// Inicializar fontbridge con la configuración de fuentes actual del documento
	// si aún no se ha inicializado
	if fontbridge.SharedFontConfig.Font == nil && doc != nil {
		// Convertir RGBColor a drawing.Color
		titleColor := fontbridge.GetDrawingColor(
			doc.fontConfig.Header1.Color.R,
			doc.fontConfig.Header1.Color.G,
			doc.fontConfig.Header1.Color.B,
		)
		normalColor := fontbridge.GetDrawingColor(
			doc.fontConfig.Normal.Color.R,
			doc.fontConfig.Normal.Color.G,
			doc.fontConfig.Normal.Color.B,
		)
		footnoteColor := fontbridge.GetDrawingColor(
			doc.fontConfig.Footnote.Color.R,
			doc.fontConfig.Footnote.Color.G,
			doc.fontConfig.Footnote.Color.B,
		)

		// Inicializar la configuración compartida
		fontbridge.InitFromDocConfig(
			doc.fontConfig.Family.Path,
			doc.fontConfig.Family.Regular,
			float64(doc.fontConfig.Header2.Size),
			float64(doc.fontConfig.Normal.Size),
			float64(doc.fontConfig.Footnote.Size),
			titleColor,
			normalColor,
			footnoteColor,
			doc.fontConfig.Normal.LineSpacing,
		)
	}

	return &docChart{
		doc:         doc,
		width:       500, // Ancho predeterminado
		height:      300, // Alto predeterminado
		keepRatio:   true,
		alignment:   Left,
		barWidth:    30,  // Ancho de barra predeterminado (ajustado)
		barSpacing:  15,  // Espacio entre barras predeterminado (ajustado)
		dpi:         150, // DPI reducido a 150
		strokeWidth: 1.0, // Ancho de línea por defecto
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

// Quality configura la calidad de la imagen del gráfico
// dpi - Resolución en puntos por pulgada (dots per inch)
// Valores recomendados:
// - 72: Calidad para pantalla
// - 150: Calidad media
// - 300: Alta calidad (por defecto)
// - 600: Calidad profesional (archivos más grandes)
func (c *docChart) Quality(dpi float64) *docChart {
	if dpi > 0 {
		c.dpi = dpi
	}
	return c
}

// Draw renderiza el gráfico en el documento con manejo de saltos de página
func (c *docChart) Draw() error {
	// Crear un archivo temporal para almacenar la imagen del gráfico
	tmpFile, err := os.CreateTemp("", "chart-*.png")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name()) // Eliminar el archivo temporal al finalizar

	// Ajustar las dimensiones del gráfico para la calidad deseada
	widthInPixels := int(c.width * c.dpi / 72.0)
	heightInPixels := int(c.height * c.dpi / 72.0)

	// Calcular factor de escala para ajustar elementos con el DPI
	// NO aplicamos el factor de escala a los tamaños de fuente
	// porque queremos que sean exactamente los definidos en FontConfig
	scaleFactor := c.dpi / 72.0

	// Asegurarnos de que fontbridge está inicializado correctamente
	if fontbridge.SharedFontConfig.Font == nil {
		// Intentar cargar la fuente predeterminada como última opción
		defaultFont, errDefault := chart.GetDefaultFont()
		if errDefault != nil {
			c.doc.log("FATAL: Could not load default chart font:", errDefault)
			return errDefault
		}
		fontbridge.SharedFontConfig.Font = defaultFont
	}

	// Crear el gráfico de barras
	barChart := chart.BarChart{
		Title:      c.title,
		Width:      widthInPixels,
		Height:     heightInPixels,
		BarWidth:   int(float64(c.barWidth) * scaleFactor),
		BarSpacing: int(float64(c.barSpacing) * scaleFactor),
		Bars:       c.bars,
		DPI:        c.dpi,
		Font:       fontbridge.SharedFontConfig.Font, // Usar la fuente compartida
	}

	// Aplicar estilos desde fontbridge - Sin escalar los tamaños de fuente
	titleStyle := chart.Style{}
	fontbridge.ApplyToChartStyle(&titleStyle, "title")
	// NO multiplicamos el tamaño de fuente por scaleFactor
	titleStyle.Padding = chart.Box{
		Top:    int(10 * scaleFactor),
		Bottom: int(5 * scaleFactor),
	}
	barChart.TitleStyle = titleStyle

	// Reservar más espacio para las etiquetas del eje X
	barChart.Height = int(float64(barChart.Height) * 0.9) // Reducir altura para dar más espacio a etiquetas

	// Configurar el canvas con más espacio en la parte inferior
	barChart.Canvas = chart.Style{
		Padding: chart.Box{
			Bottom: int(40 * scaleFactor), // Espacio adicional para etiquetas
		},
	}

	// Aplicar estilos personalizados si están configurados
	if c.background.FillColor.A > 0 {
		barChart.Background = c.background
	}

	if c.canvas.FillColor.A > 0 {
		// Mantener el padding adicional
		c.canvas.Padding = chart.Box{
			Bottom: int(40 * scaleFactor),
		}
		barChart.Canvas = c.canvas
	}

	// Configuración de ejes
	if !c.xAxisStyle.Hidden {
		xStyle := chart.Style{}
		fontbridge.ApplyToChartStyle(&xStyle, "axis")
		// NO multiplicamos el tamaño de fuente por scaleFactor
		xStyle.Hidden = false
		xStyle.StrokeWidth = c.strokeWidth
		xStyle.StrokeColor = c.xAxisStyle.StrokeColor
		xStyle.Padding = chart.Box{
			Top:    int(5 * scaleFactor),
			Bottom: int(20 * scaleFactor), // Más espacio para las etiquetas
		}
		barChart.XAxis = xStyle
	}

	if !c.yAxisStyle.Hidden {
		yStyle := chart.Style{}
		fontbridge.ApplyToChartStyle(&yStyle, "axis")
		// NO multiplicamos el tamaño de fuente por scaleFactor
		yStyle.Hidden = false
		yStyle.StrokeWidth = c.strokeWidth
		yStyle.Padding = chart.Box{
			Left:  int(5 * scaleFactor),
			Right: int(5 * scaleFactor),
		}
		barChart.YAxis = chart.YAxis{
			Style: yStyle,
		}
	}

	// Aplicar estilo para las etiquetas de las barras
	for i := range c.bars {
		valueStyle := chart.Style{}
		fontbridge.ApplyToChartStyle(&valueStyle, "value")
		// NO multiplicamos el tamaño de fuente por scaleFactor
		c.bars[i].Style = valueStyle
	}

	// Ajustar el espacio total del gráfico
	barChart.Background.Padding = chart.Box{
		Top:    int(10 * scaleFactor),
		Left:   int(10 * scaleFactor),
		Right:  int(10 * scaleFactor),
		Bottom: int(15 * scaleFactor),
	}

	// Renderizar el gráfico
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
		StrokeWidth: c.strokeWidth,
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
