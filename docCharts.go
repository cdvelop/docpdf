package docpdf

import (
	"github.com/cdvelop/docpdf/chart"
	"github.com/cdvelop/docpdf/chartutils"
	"github.com/cdvelop/docpdf/drawing"
	"github.com/cdvelop/docpdf/fontbridge"
)

// ChartType define los tipos de gráficos soportados
type ChartType string

const (
	BarChartType   ChartType = "bar"
	DonutChartType ChartType = "donut"
)

// docCharts es el punto de entrada para la API de gráficos
// Se accede mediante doc.Chart()
type docCharts struct {
	doc *Document
}

// Chart es el punto de entrada para la API de gráficos
func (doc *Document) Chart() *docCharts {
	// Inicializar fontbridge con la configuración de fuentes del documento
	if fontbridge.SharedFontConfig.Font == nil && doc != nil {
		// Color para elementos del gráfico
		titleColor := fontbridge.GetDrawingColor(
			doc.fontConfig.Header1.Color.R,
			doc.fontConfig.Header1.Color.G,
			doc.fontConfig.Header1.Color.B,
		)
		normalColor := fontbridge.GetDrawingColor(
			doc.fontConfig.ChartLabel.Color.R,
			doc.fontConfig.ChartLabel.Color.G,
			doc.fontConfig.ChartLabel.Color.B,
		)
		footnoteColor := fontbridge.GetDrawingColor(
			doc.fontConfig.Footnote.Color.R,
			doc.fontConfig.Footnote.Color.G,
			doc.fontConfig.Footnote.Color.B,
		)

		// Color específico para ejes
		axisColor := fontbridge.GetDrawingColor(
			doc.fontConfig.ChartAxisLabel.Color.R,
			doc.fontConfig.ChartAxisLabel.Color.G,
			doc.fontConfig.ChartAxisLabel.Color.B,
		)

		// Inicializar la configuración compartida
		fontbridge.InitFromDocConfig(
			doc.fontConfig.Family.Path,
			doc.fontConfig.Family.Regular,
			float64(doc.fontConfig.Header2.Size),
			float64(doc.fontConfig.ChartLabel.Size),
			float64(doc.fontConfig.Footnote.Size),
			titleColor,
			normalColor,
			footnoteColor,
			doc.fontConfig.Normal.LineSpacing,
		)

		// Actualizar específicamente el tamaño y color de las etiquetas de ejes
		fontbridge.SharedFontConfig.AxisLabelSize = float64(doc.fontConfig.ChartAxisLabel.Size)
		fontbridge.SharedFontConfig.AxisLabelColor = axisColor
	}

	return &docCharts{
		doc: doc,
	}
}

// Bar crea un nuevo elemento de gráfico de barras
// default alignment: Center
func (c *docCharts) Bar() *docBarChart {
	// Crear una instancia base del gráfico
	base := &docChart{
		doc:            c.doc,
		width:          500, // Ancho predeterminado
		height:         300, // Alto predeterminado
		keepRatio:      true,
		alignment:      Center,
		dpi:            150,                                   // DPI reducido a 150
		strokeWidth:    1.0,                                   // Ancho de línea por defecto
		valueFormatter: chartutils.FormatNumberValueFormatter, // Formateador de valores predeterminado con separadores de miles
		chartType:      BarChartType,
	}

	// Crear el gráfico de barras específico
	barChart := &docBarChart{
		docChart:   base,
		barWidth:   40, // Ancho de barra inicial (será ajustado automáticamente)
		barSpacing: 15, // Espacio entre barras inicial (será ajustado automáticamente)
	}

	// Configuración automática del formateador de etiquetas
	barChart.docChart.labelFormatter = chartutils.TruncateNameLabelFormatter(3)

	return barChart
}

// Donut crea un nuevo elemento de gráfico de tipo donut
// default alignment: Center
func (c *docCharts) Donut() *docDonutChart {
	// Crear una instancia base del gráfico
	base := &docChart{
		doc:            c.doc,
		width:          400, // Ancho predeterminado para donut
		height:         400, // Alto predeterminado para donut (normalmente cuadrado)
		keepRatio:      true,
		alignment:      Center,
		dpi:            150,                                   // DPI reducido a 150
		strokeWidth:    1.0,                                   // Ancho de línea por defecto
		valueFormatter: chartutils.FormatNumberValueFormatter, // Formateador de valores predeterminado con separadores de miles
		chartType:      DonutChartType,
	}

	// Crear el gráfico específico de tipo donut
	donutChart := &docDonutChart{
		docChart: base,
	}

	return donutChart
}

// docBarChart representa un gráfico de barras específico
type docBarChart struct {
	docChart   *docChart
	barWidth   int
	barSpacing int
	bars       []chart.Value
}

// docDonutChart representa un gráfico de tipo donut específico
type docDonutChart struct {
	docChart *docChart
	values   []chart.Value
}

// Title establece el título del gráfico de barras
func (c *docBarChart) Title(title string) *docBarChart {
	c.docChart.title = title
	return c
}

// Width establece el ancho del gráfico de barras
func (c *docBarChart) Width(w float64) *docBarChart {
	c.docChart.width = w
	return c
}

// Height establece la altura del gráfico de barras
func (c *docBarChart) Height(h float64) *docBarChart {
	c.docChart.height = h
	return c
}

// Size establece tanto el ancho como la altura explícitamente
func (c *docBarChart) Size(w, h float64) *docBarChart {
	c.docChart.width = w
	c.docChart.height = h
	c.docChart.keepRatio = false
	return c
}

// FixedPosition coloca el gráfico de barras en coordenadas específicas
func (c *docBarChart) FixedPosition(x, y float64) *docBarChart {
	c.docChart.x = x
	c.docChart.y = y
	c.docChart.hasPos = true
	return c
}

// AlignLeft alinea el gráfico de barras al margen izquierdo
func (c *docBarChart) AlignLeft() *docBarChart {
	c.docChart.alignment = Left
	return c
}

// AlignCenter centra el gráfico de barras horizontalmente
func (c *docBarChart) AlignCenter() *docBarChart {
	c.docChart.alignment = Center
	return c
}

// AlignRight alinea el gráfico de barras al margen derecho
func (c *docBarChart) AlignRight() *docBarChart {
	c.docChart.alignment = Right
	return c
}

// Inline hace que el gráfico de barras se muestre en línea con el texto
func (c *docBarChart) Inline() *docBarChart {
	c.docChart.inline = true
	return c
}

// BarWidth establece un ancho de barra inicial antes de los cálculos automáticos
func (c *docBarChart) BarWidth(width int) *docBarChart {
	c.barWidth = width
	return c
}

// BarSpacing establece un espaciado inicial entre barras antes de los cálculos automáticos
func (c *docBarChart) BarSpacing(spacing int) *docBarChart {
	c.barSpacing = spacing
	return c
}

// AddBar añade una barra con un valor y etiqueta al gráfico
func (c *docBarChart) AddBar(value float64, label string) *docBarChart {
	c.bars = append(c.bars, chart.Value{
		Value: value,
		Label: label,
	})
	return c
}

// Quality configura la calidad de la imagen del gráfico
func (c *docBarChart) Quality(dpi float64) *docBarChart {
	if dpi > 0 {
		c.docChart.dpi = dpi
	}
	return c
}

// WithAxis configura la visibilidad del eje X e Y
func (c *docBarChart) WithAxis(showX, showY bool) *docBarChart {
	c.docChart.WithAxis(showX, showY)
	return c
}

// WithLabelFormatter configura un formateador personalizado para las etiquetas
func (c *docBarChart) WithLabelFormatter(formatter chartutils.LabelFormatter) *docBarChart {
	c.docChart.labelFormatter = formatter
	return c
}

// WithValueFormatter configura un formateador personalizado para los valores numéricos
func (c *docBarChart) WithValueFormatter(formatter chart.ValueFormatter) *docBarChart {
	c.docChart.valueFormatter = formatter
	return c
}

// WithTruncateNameFormatter configura un formateador para truncar las etiquetas
func (c *docBarChart) WithTruncateNameFormatter(maxCharsPerWord, maxWidth int) *docBarChart {
	c.docChart.WithTruncateNameFormatter(maxCharsPerWord, maxWidth)
	return c
}

// WithThousandsSeparator configura un formateador para mostrar los valores con separadores de miles
func (c *docBarChart) WithThousandsSeparator() *docBarChart {
	c.docChart.valueFormatter = chartutils.FormatNumberValueFormatter
	return c
}

// WithoutThousandsSeparator configura un formateador para mostrar los valores sin separadores de miles
func (c *docBarChart) WithoutThousandsSeparator() *docBarChart {
	c.docChart.valueFormatter = chart.FloatValueFormatter
	return c
}

// WithStyle aplica un estilo personalizado al gráfico de barras
func (c *docBarChart) WithStyle(backgroundColor, barColor drawing.Color) *docBarChart {
	c.docChart.WithStyle(backgroundColor, barColor)
	return c
}

// Draw renderiza el gráfico de barras en el documento
func (c *docBarChart) Draw() error {
	// Transferir las barras al docChart base
	c.docChart.bars = c.bars
	c.docChart.barWidth = c.barWidth
	c.docChart.barSpacing = c.barSpacing

	// Dibujar usando la lógica común
	return c.docChart.Draw()
}

// Title establece el título del gráfico de donut
func (c *docDonutChart) Title(title string) *docDonutChart {
	c.docChart.title = title
	return c
}

// Width establece el ancho del gráfico de donut
func (c *docDonutChart) Width(w float64) *docDonutChart {
	c.docChart.width = w
	return c
}

// Height establece la altura del gráfico de donut
func (c *docDonutChart) Height(h float64) *docDonutChart {
	c.docChart.height = h
	return c
}

// Size establece tanto el ancho como la altura explícitamente
func (c *docDonutChart) Size(w, h float64) *docDonutChart {
	c.docChart.width = w
	c.docChart.height = h
	c.docChart.keepRatio = false
	return c
}

// FixedPosition coloca el gráfico de donut en coordenadas específicas
func (c *docDonutChart) FixedPosition(x, y float64) *docDonutChart {
	c.docChart.x = x
	c.docChart.y = y
	c.docChart.hasPos = true
	return c
}

// AlignLeft alinea el gráfico de donut al margen izquierdo
func (c *docDonutChart) AlignLeft() *docDonutChart {
	c.docChart.alignment = Left
	return c
}

// AlignCenter centra el gráfico de donut horizontalmente
func (c *docDonutChart) AlignCenter() *docDonutChart {
	c.docChart.alignment = Center
	return c
}

// AlignRight alinea el gráfico de donut al margen derecho
func (c *docDonutChart) AlignRight() *docDonutChart {
	c.docChart.alignment = Right
	return c
}

// Inline hace que el gráfico de donut se muestre en línea con el texto
func (c *docDonutChart) Inline() *docDonutChart {
	c.docChart.inline = true
	return c
}

// Quality configura la calidad de la imagen del gráfico
func (c *docDonutChart) Quality(dpi float64) *docDonutChart {
	if dpi > 0 {
		c.docChart.dpi = dpi
	}
	return c
}

// AddValue añade un valor y etiqueta al gráfico de donut
func (c *docDonutChart) AddValue(value float64, label string) *docDonutChart {
	c.values = append(c.values, chart.Value{
		Value: value,
		Label: label,
	})
	return c
}

// WithValueFormatter configura un formateador personalizado para los valores numéricos
func (c *docDonutChart) WithValueFormatter(formatter chart.ValueFormatter) *docDonutChart {
	c.docChart.valueFormatter = formatter
	return c
}

// WithThousandsSeparator configura un formateador para mostrar los valores con separadores de miles
func (c *docDonutChart) WithThousandsSeparator() *docDonutChart {
	c.docChart.valueFormatter = chartutils.FormatNumberValueFormatter
	return c
}

// WithoutThousandsSeparator configura un formateador para mostrar los valores sin separadores de miles
func (c *docDonutChart) WithoutThousandsSeparator() *docDonutChart {
	c.docChart.valueFormatter = chart.FloatValueFormatter
	return c
}

// WithStyle aplica un estilo personalizado al gráfico de donut
func (c *docDonutChart) WithStyle(backgroundColor, sliceColor drawing.Color) *docDonutChart {
	c.docChart.WithStyle(backgroundColor, sliceColor)
	return c
}

// Draw renderiza el gráfico de donut en el documento
func (c *docDonutChart) Draw() error {
	// Transferir los valores al docChart base
	c.docChart.values = c.values

	// Dibujar usando la lógica común
	return c.docChart.Draw()
}
