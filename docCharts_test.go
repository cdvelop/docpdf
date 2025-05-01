package docpdf

import (
	"sort"
	"testing"
)

func TestCharts(t *testing.T) {
	// Crear un documento con configuración predeterminada
	doc := NewDocument()

	// Añadir un título
	doc.AddHeader1("Ejemplo de gráficos con la nueva API").AlignCenter().Draw()
	doc.AddText("Este ejemplo muestra cómo usar la nueva API de gráficos con diferentes tipos.").Draw()

	// Crear un gráfico de barras usando la nueva API
	doc.AddHeader2("1. Gráfico de barras").Draw()
	doc.AddText("Se puede crear un gráfico de barras con Chart().Bar():").Draw()

	// Datos para el gráfico de barras
	bars := []struct {
		val   float64
		label string
	}{
		{1234567, "Desarrollo"},
		{2345678, "Marketing"},
		{5678901, "Investigación"},
		{3456789, "RR.HH."},
		{4567890, "Ventas"},
	}

	// Crear el gráfico de barras con la nueva API
	barChart := doc.Chart().Bar().
		Title("Ventas por Departamento")
		// BarWidth(60).
		// BarSpacing(20).
		// WithAxis(true, true).
		// AlignCenter().
		// Height(250)

	// Ordenar de mayor a menor
	sort.Slice(bars, func(i, j int) bool {
		return bars[i].val > bars[j].val
	})

	// Añadir las barras
	for _, b := range bars {
		barChart.AddBar(b.val, b.label)
	}

	// Dibujar el gráfico de barras
	barChart.Draw()

	// Agregar espacio
	doc.Br(1)

	// Crear un gráfico de tipo donut
	doc.AddHeader2("2. Gráfico tipo Donut").Draw()
	doc.AddText("Se puede crear un gráfico de tipo donut con Chart().Donut():").Draw()

	// Crear el gráfico de donut con la nueva API
	donutChart := doc.Chart().Donut().
		Title("Distribución de Ventas").
		Size(350, 350). // Los gráficos donut suelen ser cuadrados
		AlignCenter()

	// Usar los mismos datos que el gráfico de barras
	for _, b := range bars {
		donutChart.AddValue(b.val, b.label)
	}

	// Aplicar un estilo personalizado
	// bgColor := drawing.Color{R: 245, G: 245, B: 245, A: 255}   // Gris muy claro
	// sliceColor := drawing.Color{R: 30, G: 144, B: 255, A: 255} // Azul real

	// donutChart.WithStyle(bgColor, sliceColor)
	// Dibujar el gráfico de donut
	donutChart.Draw()

	// Guardar el documento
	err := doc.WritePdf("docCharts_test.pdf")
	if err != nil {
		t.Fatalf("Error al escribir PDF: %v", err)
	}
}
