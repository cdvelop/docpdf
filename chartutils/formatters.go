package chartutils

import (
	"github.com/cdvelop/tinystring"
)

// LabelFormatter es un tipo para las funciones que formatean etiquetas
// Recibe la etiqueta original y el ancho disponible para ella.
type LabelFormatter func(label string, availableWidth int) string

// ValueFormatter es un tipo para las funciones que formatean valores
type ValueFormatter func(v interface{}) string

// DefaultLabelFormatter es el formateador de etiquetas por defecto (no hace nada)
func DefaultLabelFormatter(label string, availableWidth int) string {
	return label
}

// TruncateNameLabelFormatter crea un formateador que trunca las etiquetas usando TruncateName
// maxCharsPerWord: máximo de caracteres por palabra
// El ancho máximo se pasará como segundo argumento a la función devuelta.
func TruncateNameLabelFormatter(maxCharsPerWord int) LabelFormatter {
	return func(label string, availableWidth int) string {
		// Usar availableWidth como el maxWidth para TruncateName
		if availableWidth <= 0 {
			// Si el ancho no es válido, devolver la etiqueta original
			return label
		}
		return tinystring.Convert(label).TruncateName(maxCharsPerWord, availableWidth).String()
	}
}

// FormatNumberValueFormatter formatea los valores numéricos con separadores de miles
func FormatNumberValueFormatter(v any) string {
	// tinystring.Convert puede manejar números directamente
	return tinystring.Convert(v).FormatNumber().String()
}
