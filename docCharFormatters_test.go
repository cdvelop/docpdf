package docpdf_test

import (
	"testing"

	"github.com/cdvelop/docpdf"
)

func TestChartFormatters(t *testing.T) {
	// Crear un documento con configuración predeterminada
	doc := docpdf.NewDocument()

	// Añadir un título
	doc.AddHeader1("Ejemplo de Formateadores para Gráficos").AlignCenter().Draw()

	// Añadir una explicación
	doc.AddText("Este ejemplo muestra cómo usar los nuevos formateadores de etiquetas y valores en gráficos de barras.").Draw()

	// Crear un gráfico sin formateo (para comparación)
	doc.AddHeader2("Gráfico sin formateo").Draw()
	doc.AddText("Las etiquetas y valores se muestran sin formatear:").Draw()

	chartNoFormat := doc.AddBarChart().
		Title("Ventas por Departamento").
		Height(250).
		Width(500).
		AlignCenter().
		BarWidth(40).
		BarSpacing(20).
		WithAxis(true, true).
		Quality(150)

	// Añadir datos con nombres largos y valores grandes
	chartNoFormat.AddBar(1234567, "Departamento de Desarrollo de Software")
	chartNoFormat.AddBar(2345678, "Departamento de Marketing y Publicidad")
	chartNoFormat.AddBar(3456789, "Departamento de Recursos Humanos")
	chartNoFormat.AddBar(4567890, "Departamento de Ventas y Atención al Cliente")
	chartNoFormat.AddBar(5678901, "Departamento de Investigación")

	// Dibujar el gráfico sin formateo
	chartNoFormat.Draw()

	// Agregar algo de espacio
	doc.AddText("").Draw()
	doc.AddText("").Draw()

	// Crear un gráfico con formateo
	doc.AddHeader2("Gráfico con formateo").Draw()
	doc.AddText("Las etiquetas se truncan con TruncateName y los valores se muestran con separadores de miles:").Draw()

	chartWithFormat := doc.AddBarChart().
		Title("Ventas por Departamento").
		Height(250).
		Width(500).
		AlignCenter().
		BarWidth(40).
		BarSpacing(20).
		WithAxis(true, true).
		Quality(150).
		// Usar los nuevos métodos de formateo
		WithTruncateNameFormatter(3, 15). // Máximo 3 caracteres por palabra, 15 en total
		WithThousandsSeparator()          // Añadir separadores de miles

	// Añadir los mismos datos
	chartWithFormat.AddBar(1234567, "Departamento de Desarrollo de Software")
	chartWithFormat.AddBar(2345678, "Departamento de Marketing y Publicidad")
	chartWithFormat.AddBar(3456789, "Departamento de Recursos Humanos")
	chartWithFormat.AddBar(4567890, "Departamento de Ventas y Atención al Cliente")
	chartWithFormat.AddBar(5678901, "Departamento de Investigación")

	// Dibujar el gráfico con formateo
	chartWithFormat.Draw()

	// Guardar el documento
	err := doc.WritePdf("docCharFormatters_test.pdf")
	if err != nil {
		t.Fatalf("Error al escribir PDF: %v", err)
	}
}
