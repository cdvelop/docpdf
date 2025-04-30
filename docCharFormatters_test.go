package docpdf

import (
	"sort"
	"testing"

	"github.com/cdvelop/docpdf/chart"
	"github.com/cdvelop/docpdf/chartutils"
)

func TestChartFormatters(t *testing.T) {
	// Crear un documento con configuración predeterminada
	doc := NewDocument()

	// Añadir un título
	doc.AddHeader1("Ejemplo de Formateadores para Gráficos").AlignCenter().Draw()
	// Añadir una explicación
	doc.AddText("Este ejemplo muestra cómo usar los nuevos formateadores de etiquetas y valores en gráficos de barras.").Draw()

	// Crear un gráfico sin formateo de miles (para comparación)
	doc.AddHeader2("Gráfico sin formateo de miles").Draw()
	doc.AddText("Las etiquetas y valores se muestran sin separadores de miles:").Draw()

	// Datos originales
	bars := []struct {
		val   float64
		label string
	}{
		{1234567, "Departamento de Desarrollo de Software"},
		{2345678, "Departamento de Marketing y Publicidad"},
		{5678901, "Departamento de Investigación"},
		{3456789, "Departamento de Recursos Humanos"},
		{4567890, "Departamento de Ventas y Atención al Cliente"},
		{6789012, "Departamento de Finanzas y Contabilidad"},
		{7890123, "Departamento de Soporte Técnico"},
		{8901234, "Departamento de Operaciones"},
		{9012345, "Departamento de Logística"},
		{1345678, "Departamento de Calidad"},
	}

	chartNoFormat := doc.AddBarChart().
		Title("Ventas por Departamento").
		// BarWidth(40).
		// BarSpacing(0).
		WithoutThousandsSeparator() // Explícitamente desactivar el separador de miles
	// Quality(150)

	// Ahora agrega las barras en orden
	for _, b := range bars {
		chartNoFormat.AddBar(b.val, b.label)
	}

	// Dibujar el gráfico sin formateo
	chartNoFormat.Draw()

	// Agregar algo de espacio
	doc.Br(1)
	// Crear un gráfico con formateo
	doc.AddHeader2("Gráfico con formateo").Draw()
	doc.AddText("Las etiquetas se truncan con TruncateName y los valores se muestran con separadores de miles (por defecto):").Draw()

	// Usar los mismos valores de barWidth y barSpacing para el segundo gráfico
	chartWithFormat := doc.AddBarChart().
		Title("Ventas por Departamento").
		// Height(250).
		WithTruncateNameFormatter(3, 15) // Máximo 3 caracteres por palabra, 15 en total

	// Ordenar de mayor a menor
	sort.Slice(bars, func(i, j int) bool {
		return bars[i].val > bars[j].val
	})

	// Añadir los mismos datos ordenados
	for _, b := range bars {
		chartWithFormat.AddBar(b.val, b.label)
	}

	// Dibujar el gráfico con formateo
	chartWithFormat.Draw()

	// Guardar el documento
	err := doc.WritePdf("docCharFormatters_test.pdf")
	if err != nil {
		t.Fatalf("Error al escribir PDF: %v", err)
	}
}

func TestChartFormattersLabels(t *testing.T) {
	// Crear un documento de prueba
	doc := NewDocument()

	// Prueba 1: Formateador de etiquetas
	t.Run("LabelFormatter", func(t *testing.T) {
		// Crear un gráfico con un formateador de etiquetas personalizado
		barChart := doc.AddBarChart()

		// Crear un formateador simple que añade un prefijo
		testFormatter := func(label string) string {
			return "Test-" + label
		}

		// Aplicar el formateador
		barChart.WithLabelFormatter(testFormatter)

		// Añadir una barra
		barChart.AddBar(100, "Label")

		// Verificar que el formateador se guardó
		if barChart.labelFormatter == nil {
			t.Error("El formateador de etiquetas no se guardó correctamente")
		}

		// Verificar que el formateador funciona como se espera
		result := barChart.labelFormatter("Label")
		expected := "Test-Label"
		if result != expected {
			t.Errorf("El formateador de etiquetas no funcionó correctamente. Esperado: %s, Obtenido: %s", expected, result)
		}
	})

	// Prueba 2: Formateador de valores
	t.Run("ValueFormatter", func(t *testing.T) {
		// Crear un gráfico con un formateador de valores personalizado
		barChart := doc.AddBarChart()

		// Crear un formateador simple que añade un sufijo
		testFormatter := func(v any) string {
			baseResult := chart.FloatValueFormatter(v)
			return baseResult + " units"
		}

		// Aplicar el formateador
		barChart.WithValueFormatter(testFormatter)

		// Verificar que el formateador se guardó
		if barChart.valueFormatter == nil {
			t.Error("El formateador de valores no se guardó correctamente")
		}

		// Verificar que el formateador funciona como se espera
		result := barChart.valueFormatter(100.0)
		expected := "100.00 units"
		if result != expected {
			t.Errorf("El formateador de valores no funcionó correctamente. Esperado: %s, Obtenido: %s", expected, result)
		}
	})

	// Prueba 3: Método de conveniencia WithTruncateNameFormatter
	t.Run("TruncateNameFormatter", func(t *testing.T) {
		// Crear un gráfico con el formateador TruncateName
		barChart := doc.AddBarChart()

		// Aplicar el formateador
		barChart.WithTruncateNameFormatter(3, 10)

		// Verificar que el formateador se guardó
		if barChart.labelFormatter == nil {
			t.Error("El formateador TruncateName no se guardó correctamente")
		}

		// Verificar que el formateador funciona como se espera
		result := barChart.labelFormatter("Departamento de Ventas")
		// Debería truncar a "Dep. de Ven..." o similar
		if len(result) > 10 {
			t.Errorf("El formateador TruncateName no limitó la longitud correctamente: %s", result)
		}
	})

	// Prueba 4: Método de conveniencia WithThousandsSeparator
	t.Run("ThousandsSeparator", func(t *testing.T) {
		// Crear un gráfico con el formateador FormatNumber
		barChart := doc.AddBarChart()

		// Aplicar el formateador
		barChart.WithThousandsSeparator()

		// Verificar que el formateador se guardó
		if barChart.valueFormatter == nil {
			t.Error("El formateador FormatNumber no se guardó correctamente")
		}

		// Verificar que el formateador funciona como se espera
		result := barChart.valueFormatter(1234567.0)
		expected := "1.234.567" // Con separadores de miles
		if result != expected {
			t.Errorf("El formateador FormatNumber no funcionó correctamente. Esperado: %s, Obtenido: %s", expected, result)
		}
	})

	// Prueba 5: Integración con chartutils.FormatNumberValueFormatter
	t.Run("FormatNumberValueFormatter", func(t *testing.T) {
		result := chartutils.FormatNumberValueFormatter(1234567.0)
		expected := "1.234.567"
		if result != expected {
			t.Errorf("FormatNumberValueFormatter no funcionó correctamente. Esperado: %s, Obtenido: %s", expected, result)
		}
	})
}
