package docpdf

import (
	"bytes"

	"github.com/cdvelop/docpdf/chart"
	"github.com/cdvelop/docpdf/chartutils"
	"github.com/cdvelop/docpdf/drawing"
	"github.com/cdvelop/docpdf/fontbridge"
	"github.com/cdvelop/tinystring" // <--- Añadir import
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
	title          string
	barWidth       int
	barSpacing     int
	bars           []chart.Value
	xAxisStyle     chart.Style
	yAxisStyle     chart.Style
	background     chart.Style
	canvas         chart.Style
	labelFormatter chartutils.LabelFormatter // Formateador para etiquetas (recibe label y ancho disponible)
	valueFormatter chart.ValueFormatter      // Formateador para valores

	// Propiedades para control de calidad
	dpi         float64 // Resolución del gráfico en DPI (dots per inch)
	strokeWidth float64 // Ancho de línea para los contornos
}

// AddBarChart crea un nuevo elemento de gráfico de barras
// default alignment: Center
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
	chart := &docChart{
		doc:            doc,
		width:          500, // Ancho predeterminado
		height:         300, // Alto predeterminado
		keepRatio:      true,
		alignment:      Center,
		barWidth:       40,                                    // Ancho de barra inicial (será ajustado automáticamente)
		barSpacing:     15,                                    // Espacio entre barras inicial (será ajustado automáticamente)
		dpi:            150,                                   // DPI reducido a 150
		strokeWidth:    1.0,                                   // Ancho de línea por defecto
		valueFormatter: chartutils.FormatNumberValueFormatter, // Formateador de valores predeterminado con separadores de miles
	}

	// Configuración automática del formateador de etiquetas
	// Ahora TruncateNameLabelFormatter solo necesita maxCharsPerWord
	chart.labelFormatter = chartutils.TruncateNameLabelFormatter(3)

	return chart
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

// BarWidth establece un ancho de barra inicial antes de los cálculos automáticos
// Nota: El ancho final puede ser ajustado automáticamente para aprovechar el espacio disponible
func (c *docChart) BarWidth(width int) *docChart {
	c.barWidth = width
	return c
}

// BarSpacing establece un espaciado inicial entre barras antes de los cálculos automáticos
// Nota: El espaciado final puede ser ajustado automáticamente para aprovechar el espacio disponible
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

// calculateBarLayout calcula automáticamente el ancho de barras y espaciado
// para que ocupen todo el ancho de contenido disponible
func (c *docChart) calculateBarLayout() {
	if len(c.bars) == 0 {
		return // No hay barras para calcular
	}

	// Usar contentAreaWidth del documento como ancho total disponible
	chartWidth := c.doc.contentAreaWidth

	n := len(c.bars)
	if n > 1 {
		// Usar el ancho de barra configurado como base
		barWidth := float64(c.barWidth) // Usamos el ancho ya configurado como punto de partida

		// Calcular espaciado basado en el espacio disponible
		barSpacing := (chartWidth - float64(n)*barWidth) / float64(n-1)

		// Si el espaciado es negativo o muy pequeño, recalcular
		if barSpacing < 5 {
			// Asignar un espaciado mínimo
			barSpacing = 5
			// Recalcular el ancho de barras con este espaciado mínimo
			barWidth = (chartWidth - barSpacing*float64(n-1)) / float64(n)

			// Si el ancho es demasiado pequeño, usar un valor mínimo
			if barWidth < 20 {
				barWidth = 20
				// Ya no intentamos ajustar más, puede que el gráfico sea más ancho que contentAreaWidth
			}
		}

		// Actualizar las propiedades del objeto docChart
		c.barWidth = int(barWidth)
		c.barSpacing = int(barSpacing)
	} else {
		// Si solo hay una barra, usar todo el ancho disponible
		c.barWidth = int(chartWidth)
		c.barSpacing = 0
	}
}

// Draw renderiza el gráfico en el documento con manejo de saltos de página
func (c *docChart) Draw() error {
	// Calcular automáticamente el ancho de barras y espaciado ANTES de formatear
	if len(c.bars) > 0 {
		c.calculateBarLayout()
	}

	// Crear un buffer en memoria para almacenar la imagen del gráfico
	var buf bytes.Buffer

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

	// Aplicar formateadores a las etiquetas y valores antes de crear el gráfico
	formattedBars := make([]chart.Value, len(c.bars))
	for i, bar := range c.bars {
		// Copia la barra original
		formattedBars[i] = bar

		// Aplicar formateador de etiquetas si está definido
		if c.labelFormatter != nil { // Asegurarse de que el formateador existe

			// Pasar el ancho de barra calculado como segundo argumento
			formattedBars[i].Label = c.labelFormatter(bar.Label, c.barWidth-25) // Ajustar el ancho para evitar que se corte
			// formattedBars[i].Label = c.labelFormatter(bar.Label, c.barWidth)
		}
	} // Crear el gráfico de barras (usando las barras ya formateadas)
	barChart := chart.BarChart{
		Title:      c.title,
		Width:      widthInPixels,
		Height:     heightInPixels,
		BarWidth:   int(float64(c.barWidth) * scaleFactor),   // Usar el barWidth calculado
		BarSpacing: int(float64(c.barSpacing) * scaleFactor), // Usar el barSpacing calculado
		Bars:       formattedBars,                            // Usar las barras con etiquetas ya formateadas
		DPI:        c.dpi,
		Font:       fontbridge.SharedFontConfig.Font, // Usar la fuente compartida
	}
	// Configurar el formateador de valores en el eje Y
	if c.valueFormatter != nil {
		barChart.YAxis = chart.YAxis{
			ValueFormatter: c.valueFormatter,
		}
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
	// Factor de reducción predeterminado
	heightReductionFactor := 0.9

	// Si se ha configurado un padding específico para el eje X, ajustar la reducción proporcionalmente
	if c.xAxisStyle.Padding.Bottom > 20 {
		// Ajustamos el factor de reducción en función del padding
		// Mayor padding requiere más reducción de altura
		heightReductionFactor = 0.9 - float64(c.xAxisStyle.Padding.Bottom-20)/200
		// Limitamos el factor entre 0.75 y 0.9 para evitar reducciones extremas
		if heightReductionFactor < 0.75 {
			heightReductionFactor = 0.75
		}
	}

	barChart.Height = int(float64(barChart.Height) * heightReductionFactor) // Reducir altura para dar más espacio a etiquetas
	// Configurar el canvas con más espacio en la parte inferior
	// Valor predeterminado para el espacio inferior
	bottomCanvasPadding := int(40 * scaleFactor)

	// Si se ha configurado un padding específico para el eje X, aumentamos también el canvas
	if c.xAxisStyle.Padding.Bottom > 0 {
		bottomCanvasPadding = int(float64(c.xAxisStyle.Padding.Bottom) * 1.5 * scaleFactor)
	}

	barChart.Canvas = chart.Style{
		Padding: chart.Box{
			Bottom: bottomCanvasPadding, // Espacio adicional para etiquetas
		},
	}

	// Aplicar estilos personalizados si están configurados
	if c.background.FillColor.A > 0 {
		barChart.Background = c.background
	}
	if c.canvas.FillColor.A > 0 {
		// Mantener el padding adicional
		// Usar el mismo valor de bottomCanvasPadding que se calculó anteriormente
		c.canvas.Padding = chart.Box{
			Bottom: bottomCanvasPadding,
		}
		barChart.Canvas = c.canvas
	} // Configuración de ejes
	if !c.xAxisStyle.Hidden {
		xStyle := chart.Style{}
		fontbridge.ApplyToChartStyle(&xStyle, "axis")
		// NO multiplicamos el tamaño de fuente por scaleFactor
		xStyle.Hidden = false
		xStyle.StrokeWidth = c.strokeWidth
		xStyle.StrokeColor = c.xAxisStyle.StrokeColor

		// Aplicar el padding personalizado si se ha configurado, o usar el valor predeterminado
		bottomPadding := int(20 * scaleFactor) // Valor predeterminado
		if c.xAxisStyle.Padding.Bottom > 0 {
			bottomPadding = c.xAxisStyle.Padding.Bottom
		}

		xStyle.Padding = chart.Box{
			Top:    int(5 * scaleFactor),
			Bottom: bottomPadding, // Usar el padding configurado o el predeterminado
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

		// Configurar YAxis con estilo y formateador en un paso
		yAxis := chart.YAxis{
			Style: yStyle,
		}

		// Aplicar formateador de valores si existe
		if c.valueFormatter != nil {
			yAxis.ValueFormatter = c.valueFormatter
		}

		// Asignar YAxis configurado
		barChart.YAxis = yAxis
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
	// Renderizar el gráfico directamente al buffer en memoria
	err := barChart.Render(chart.PNG, &buf)
	if err != nil {
		return err
	}

	// Usar los bytes del buffer directamente sin necesidad de archivo temporal
	docImage := c.doc.AddImage(buf.Bytes())

	err = docImage.Draw()
	if err != nil {
		return err
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
		// No establecemos ValueFormatter aquí, se hace en Draw()
	}

	return c
}

// WithLabelFormatter configura un formateador personalizado para las etiquetas de las barras
// El formateador debe tener la firma: func(label string, availableWidth int) string
func (c *docChart) WithLabelFormatter(formatter chartutils.LabelFormatter) *docChart {
	c.labelFormatter = formatter
	return c
}

// WithValueFormatter configura un formateador personalizado para los valores numéricos
func (c *docChart) WithValueFormatter(formatter chart.ValueFormatter) *docChart {
	c.valueFormatter = formatter
	return c
}

// WithTruncateNameFormatter configura un formateador para truncar las etiquetas usando TruncateName
// con un ancho máximo FIJO, ignorando el ancho de barra calculado.
// maxCharsPerWord: máximo de caracteres por palabra
// maxWidth: máximo ancho total FIJO de la etiqueta
func (c *docChart) WithTruncateNameFormatter(maxCharsPerWord, maxWidth int) *docChart {
	// Crear una clausura que capture maxCharsPerWord y maxWidth fijos
	c.labelFormatter = func(label string, availableWidth int) string {
		// Ignorar availableWidth y usar el maxWidth fijo capturado
		if maxWidth <= 0 {
			return label // Si el ancho fijo no es válido, devolver original
		}
		return tinystring.Convert(label).TruncateName(maxCharsPerWord, maxWidth).String()
	}
	return c
}

// WithThousandsSeparator configura un formateador para mostrar los valores con separadores de miles
func (c *docChart) WithThousandsSeparator() *docChart {
	c.valueFormatter = chartutils.FormatNumberValueFormatter
	return c
}

// WithoutThousandsSeparator configura un formateador para mostrar los valores sin separadores de miles
func (c *docChart) WithoutThousandsSeparator() *docChart {
	c.valueFormatter = chart.FloatValueFormatter
	return c
}

// WithCustomLabelFormatter es un alias para WithLabelFormatter para mantener consistencia
func (c *docChart) WithCustomLabelFormatter(formatter chartutils.LabelFormatter) *docChart {
	return c.WithLabelFormatter(formatter)
}

// WithCustomValueFormatter es un alias para WithValueFormatter para mantener consistencia
func (c *docChart) WithCustomValueFormatter(formatter chart.ValueFormatter) *docChart {
	return c.WithValueFormatter(formatter)
}
