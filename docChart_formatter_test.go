package docpdf

import (
	"testing"

	"github.com/cdvelop/docpdf/chart"
	"github.com/cdvelop/docpdf/chartutils"
)

func TestChartFormatters(t *testing.T) {
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
