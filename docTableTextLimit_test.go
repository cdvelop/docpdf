package docpdf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTableCellTextLimitation(t *testing.T) {
	// Create a simple document with FileWriter function
	doc := NewDocument()

	// Add a title
	doc.AddHeader1("Prueba de Límite de Texto en Celdas").AlignCenter().Draw()
	doc.AddText("Este ejemplo prueba la limitación de texto a 2 líneas en las celdas de la tabla").Draw()
	doc.SpaceBefore(1)

	// Crear tabla para probar la limitación de texto
	table := doc.NewTable(
		"Columna 1|W:30%",
		"Columna 2|W:30%",
		"Columna 3|W:40%",
	)

	// Añadir filas con diferentes cantidades de texto
	table.AddRow(
		"Texto corto en una línea",
		"Texto de dos líneas\nque debe verse completo",
		"Este es un texto largo que debería ocupar varias líneas al mostrarse en la celda y debería truncarse con puntos suspensivos",
	)

	// Añadir otra fila con texto multilínea explícito
	table.AddRow(
		"Una línea",
		"Línea 1\nLínea 2\nLínea 3 (no debería verse)",
		"Línea A\nLínea B\nLínea C (no debería verse)\nLínea D (no debería verse)",
	)

	// Dibujar la tabla
	table.Draw()

	// Añadir texto explicativo
	doc.SpaceBefore(1)
	doc.AddText("Las celdas anteriores deberían mostrar como máximo 2 líneas de texto").
		Bold().Draw()
	doc.AddText("Si hay más de 2 líneas, se deben mostrar puntos suspensivos (...)").Draw()

	// Crear directorio de salida si no existe
	outDir := "test/out"
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		t.Fatalf("Error al crear directorio de salida: %v", err)
	}

	// Establecer ruta del archivo de salida
	outFilePath := filepath.Join(outDir, "table_text_limitation_test.pdf")

	// Guardar el documento
	err = doc.WritePdf(outFilePath)
	if err != nil {
		t.Fatalf("Error al escribir PDF: %v", err)
	}

	absPath, _ := filepath.Abs(outFilePath)
	t.Logf("PDF de prueba de limitación de texto creado correctamente en: %s", absPath)
}
