// filepath: c:\Users\Cesar\Packages\Internal\docpdf\docChart.go
package docpdf

import (
	"bytes"

	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/chart"
	"github.com/cdvelop/docpdf/chartutils"
	"github.com/cdvelop/docpdf/config"
	"github.com/cdvelop/docpdf/fontbridge"
	"github.com/cdvelop/tinystring"
)

// docChart representa la configuración base común a todos los gráficos
type docChart struct {
	doc       *Document
	width     float64
	height    float64
	keepRatio bool
	alignment config.Alignment
	x, y      float64
	hasPos    bool
	inline    bool
	chartType chartType

	// Propiedades comunes a todos los gráficos
	title          string
	xAxisStyle     chart.Style
	yAxisStyle     chart.Style
	background     chart.Style
	canvas         chart.Style
	labelFormatter chartutils.LabelFormatter // Formateador para etiquetas
	valueFormatter chart.ValueFormatter      // Formateador para valores

	// Propiedades específicas para distintos tipos de gráficos
	// Para BarChart
	barWidth   int
	barSpacing int
	bars       []chart.Value

	// Para DonutChart
	values []chart.Value

	// Propiedades para control de calidad
	dpi         float64 // Resolución del gráfico en DPI
	strokeWidth float64 // Ancho de línea para los contornos
}

// WithStyle aplica un estilo personalizado al gráfico
func (c *docChart) WithStyle(backgroundColor, chartColor config.Color) *docChart {
	c.background = chart.Style{
		FillColor: backgroundColor,
	}

	// Usar el color proporcionado directamente en lugar de intentar crear una nueva paleta
	c.canvas = chart.Style{
		StrokeColor: chartColor,
		StrokeWidth: c.strokeWidth,
		FillColor:   chartColor.WithAlpha(100), // Versión semi-transparente para relleno
	}

	return c
}

// WithAxis configura la visibilidad del eje X e Y
func (c *docChart) WithAxis(showX, showY bool) *docChart {
	if !showX {
		c.xAxisStyle.Hidden = true
	} else {
		c.xAxisStyle = chart.Shown()
		// Aplicar configuración específica para eje X
		c.xAxisStyle.Padding = canvas.Box{
			Bottom: 20, // Espacio predeterminado para etiquetas del eje X
			Top:    5,
		}
	}

	if !showY {
		c.yAxisStyle.Hidden = true
	} else {
		c.yAxisStyle = chart.Shown()
		// Aplicar configuración específica para eje Y
		c.yAxisStyle.Padding = canvas.Box{
			Left:  10,
			Right: 5,
		}
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
	// Verificar si podemos usar renderizado directo a PDF para el gráfico donut
	// Solo habilitado para donut como prueba de integración
	if c.chartType == donutChartType {
		// Transferir los valores desde el docDonutChart.values
		// Esto asegura que tenemos los valores correctos al renderizar
		return c.drawWithPdfRenderer()
	}

	// Si no podemos usar renderizado directo, continuamos con el método tradicional
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
			c.doc.Log("FATAL: Could not load default chart font:", errDefault)
			return errDefault
		}
		fontbridge.SharedFontConfig.Font = defaultFont
	}

	// Inicializar el estilo del título (común para todos los tipos de gráfico)
	titleStyle := chart.Style{}
	fontbridge.ApplyToChartStyle(&titleStyle, "title")
	titleStyle.Padding = canvas.Box{
		Top:    int(10 * scaleFactor),
		Bottom: int(5 * scaleFactor),
	}

	// Preparar el estilo de fondo común
	backgroundStyle := chart.Style{
		Padding: canvas.Box{
			Top:    int(10 * scaleFactor),
			Left:   int(10 * scaleFactor),
			Right:  int(10 * scaleFactor),
			Bottom: int(15 * scaleFactor),
		},
	}

	// Si hay estilos personalizados, aplicarlos
	if c.background.FillColor.A > 0 {
		backgroundStyle.FillColor = c.background.FillColor
		backgroundStyle.StrokeColor = c.background.StrokeColor
		backgroundStyle.StrokeWidth = c.background.StrokeWidth
	}

	var err error

	// Renderizar según el tipo de gráfico
	switch c.chartType {
	case barChartType:
		err = c.drawBarChart(&buf, widthInPixels, heightInPixels, scaleFactor, titleStyle, backgroundStyle)
	case donutChartType:
		err = c.drawDonutChart(&buf, widthInPixels, heightInPixels, scaleFactor, titleStyle, backgroundStyle)
	default:
		// Si no se reconoce el tipo, usar el gráfico de barras como predeterminado
		err = c.drawBarChart(&buf, widthInPixels, heightInPixels, scaleFactor, titleStyle, backgroundStyle)
	}

	if err != nil {
		return err
	}

	// Usar los bytes del buffer directamente sin necesidad de archivo temporal
	docImage := c.doc.AddImage(buf.Bytes())
	//docImage := c.doc.AddChartSvg(buf.Bytes())
	// ajustar la alineación según la configuración
	docImage.alignment = c.alignment

	return docImage.Draw()
}

// drawBarChart renderiza un gráfico de barras específicamente
func (c *docChart) drawBarChart(buf *bytes.Buffer, widthInPixels, heightInPixels int, scaleFactor float64, titleStyle, backgroundStyle chart.Style) error {
	// Calcular automáticamente el ancho de barras y espaciado ANTES de formatear
	if len(c.bars) > 0 {
		c.calculateBarLayout()
	}

	// Aplicar formateadores a las etiquetas y valores antes de crear el gráfico
	formattedBars := make([]chart.Value, len(c.bars))
	for i, bar := range c.bars {
		// Copia la barra original
		formattedBars[i] = bar

		// Aplicar formateador de etiquetas si está definido
		if c.labelFormatter != nil { // Asegurarse de que el formateador existe
			// Pasar el ancho de barra calculado como segundo argumento
			formattedBars[i].Label = c.labelFormatter(bar.Label, c.barWidth-10) // Ajustar el ancho para evitar que se corte
		}
	}

	// Crear el gráfico de barras (usando las barras ya formateadas)
	barChart := chart.BarChart{
		Title:      c.title,
		Width:      widthInPixels,
		Height:     heightInPixels,
		BarWidth:   int(float64(c.barWidth) * scaleFactor),   // Usar el barWidth calculado
		BarSpacing: int(float64(c.barSpacing) * scaleFactor), // Usar el barSpacing calculado
		Bars:       formattedBars,                            // Usar las barras con etiquetas ya formateadas
		DPI:        c.dpi,
		Font:       fontbridge.SharedFontConfig.Font, // Usar la fuente compartida
		TitleStyle: titleStyle,
		Background: backgroundStyle,
	}

	// Configurar el formateador de valores en el eje Y
	if c.valueFormatter != nil {
		barChart.YAxis = chart.YAxis{
			ValueFormatter: c.valueFormatter,
		}
	}

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
		Padding: canvas.Box{
			Bottom: bottomCanvasPadding, // Espacio adicional para etiquetas
		},
	}

	// Aplicar estilos personalizados al canvas si están configurados
	if c.canvas.FillColor.A > 0 {
		// Mantener el padding adicional
		c.canvas.Padding = canvas.Box{
			Bottom: bottomCanvasPadding,
		}
		barChart.Canvas = c.canvas
	}

	// Configuración de ejes
	if !c.xAxisStyle.Hidden {
		xStyle := chart.Style{}
		fontbridge.ApplyToChartStyle(&xStyle, "axis")
		xStyle.Hidden = false
		xStyle.StrokeWidth = c.strokeWidth
		xStyle.StrokeColor = c.xAxisStyle.StrokeColor

		// Aplicar el padding personalizado si se ha configurado, o usar el valor predeterminado
		bottomPadding := int(20 * scaleFactor) // Valor predeterminado
		if c.xAxisStyle.Padding.Bottom > 0 {
			bottomPadding = c.xAxisStyle.Padding.Bottom
		}

		xStyle.Padding = canvas.Box{
			Top:    int(5 * scaleFactor),
			Bottom: bottomPadding, // Usar el padding configurado o el predeterminado
		}
		barChart.XAxis = xStyle
	}

	if !c.yAxisStyle.Hidden {
		yStyle := chart.Style{}
		fontbridge.ApplyToChartStyle(&yStyle, "axis")
		yStyle.Hidden = false
		yStyle.StrokeWidth = c.strokeWidth
		yStyle.Padding = canvas.Box{
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
		c.bars[i].Style = valueStyle
	}

	// Renderizar el gráfico directamente al buffer en memoria
	return barChart.Render(chart.PNG, buf)
}

// drawDonutChart renderiza un gráfico de tipo donut específicamente
func (c *docChart) drawDonutChart(buf *bytes.Buffer, widthInPixels, heightInPixels int, scaleFactor float64, titleStyle, backgroundStyle chart.Style) error {
	// Asegurarse de que tenemos al menos un valor
	if len(c.values) == 0 {
		return nil // No hay valores para mostrar
	}

	// Calcular ancho óptimo para etiquetas basado en el tamaño del gráfico
	// Para el donut, usaremos un ancho adaptado al espacio disponible
	labelWidth := widthInPixels / 20 // Un valor proporcional al ancho total del gráfico
	if labelWidth < 10 {
		labelWidth = 10 // Mínimo para evitar etiquetas muy truncadas
	} else if labelWidth > 50 {
		labelWidth = 50 // Máximo para evitar etiquetas demasiado largas
	}

	// Si tenemos un formateador de etiquetas, aplicarlo a cada valor
	formattedValues := make([]chart.Value, len(c.values))
	for i, val := range c.values {
		formattedValues[i] = val

		// Aplicar formateador de etiquetas si existe
		if c.labelFormatter != nil {
			// Usar ancho calculado para las etiquetas del donut
			formattedValues[i].Label = c.labelFormatter(val.Label, labelWidth)
		}

		// Aplicar formateador de valores si existe (para mostrar valores junto a etiquetas)
		if c.valueFormatter != nil {
			formattedValues[i].Label = formattedValues[i].Label + " (" + c.valueFormatter(val.Value) + ")"
		}
	}

	donutChart := chart.DonutChart{
		Title:        c.title,
		Width:        widthInPixels,
		Height:       heightInPixels,
		DPI:          c.dpi,
		Values:       formattedValues,
		TitleStyle:   titleStyle,
		Background:   backgroundStyle,
		Font:         fontbridge.SharedFontConfig.Font,
		ColorPalette: chart.AlternateColorPalette,
	}

	// Asegurarse de que todos los valores tengan el mismo estilo
	for i := range formattedValues {
		valueStyle := chart.Style{}
		fontbridge.ApplyToChartStyle(&valueStyle, "value")
		formattedValues[i].Style = valueStyle
	}

	// Aplicar estilos personalizados al canvas si están configurados
	if c.canvas.FillColor.A > 0 {
		donutChart.Canvas = c.canvas
	}

	// Renderizar el gráfico directamente al buffer en memoria
	return donutChart.Render(chart.PNG, buf)
}

// configureBaseChart establece la configuración común para cualquier tipo de gráfico
// Esta función centraliza la lógica de configuración para mantener consistencia
func configureBaseChart(doc *Document, chartType chartType) *docChart {
	// Crear instancia base con configuración común
	base := &docChart{
		doc:            doc,
		width:          500, // Ancho predeterminado
		height:         300, // Alto predeterminado
		keepRatio:      true,
		alignment:      config.Center,
		dpi:            150,                                   // DPI reducido a 150
		strokeWidth:    1.0,                                   // Ancho de línea por defecto
		valueFormatter: chartutils.FormatNumberValueFormatter, // Formateador de valores predeterminado con separadores de miles
		chartType:      chartType,
	}

	// Configuración automática del formateador de etiquetas
	base.labelFormatter = chartutils.TruncateNameLabelFormatter(3)

	return base
}
