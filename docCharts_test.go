package docpdf

import (
	"sort"
	"testing"
)

func TestCharts(t *testing.T) {
	// Crear un documento con configuración predeterminada
	doc := NewDocument()

	// Añadir un título
	doc.AddHeader1("Prueba de integración de PdfRenderer con gráfico Donut").AlignCenter().Draw()
	doc.AddText("Esta prueba enfoca exclusivamente en el renderizado directo a PDF usando PdfRenderer.").Draw()

	// Datos para el gráfico
	data := []struct {
		val   float64
		label string
	}{
		{1234567, "Desarrollo"},
		{2345678, "Marketing"},
		{5678901, "Investigación"},
		{3456789, "RR.HH."},
		{4567890, "Ventas"},
	}

	// Ordenar de mayor a menor
	sort.Slice(data, func(i, j int) bool {
		return data[i].val > data[j].val
	})

	// Crear un gráfico de tipo donut - Prueba de renderizado directo a PDF
	doc.AddHeader2("Gráfico tipo Donut con renderizado directo a PDF").Draw()
	doc.AddText("Este gráfico utiliza PdfRenderer para dibujar directamente en el PDF:").Draw()

	// Crear el gráfico de donut con la nueva API unificada
	donutChart := doc.Chart().Donut().
		Title("Distribución de Ventas").
		WithTruncateNameFormatter(3, 30) // Aplicamos formato consistente para etiquetas

	// Añadir datos al gráfico donut
	for _, item := range data {
		donutChart.AddValue(item.val, item.label)
	}

	// Dibujar el gráfico de donut usando el renderizador directo a PDF
	donutChart.Draw()

	// Agregar texto explicativo
	doc.AddText("El gráfico anterior fue renderizado directamente en el PDF sin generar imágenes intermedias,").Draw()
	doc.AddText("utilizando el nuevo PdfRenderer que implementa la interfaz chart.Renderer.").Draw()

	// Guardar el documento
	err := doc.WritePdf("docCharts_test.pdf")
	if err != nil {
		t.Fatalf("Error al escribir PDF: %v", err)
	}
}
