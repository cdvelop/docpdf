package chartutils

import (
	"github.com/cdvelop/tinystring"
)

// LabelFormatter es un tipo para las funciones que formatean etiquetas
type LabelFormatter func(label string) string

// ValueFormatter es un tipo para las funciones que formatean valores
type ValueFormatter func(v interface{}) string

// DefaultLabelFormatter es el formateador de etiquetas por defecto (no hace nada)
func DefaultLabelFormatter(label string) string {
	return label
}

// TruncateNameLabelFormatter crea un formateador que trunca las etiquetas usando TruncateName
// maxCharsPerWord: máximo de caracteres por palabra
// maxWidth: máximo ancho total de la etiqueta
func TruncateNameLabelFormatter(maxCharsPerWord, maxWidth int) LabelFormatter {
	return func(label string) string {
		return tinystring.Convert(label).TruncateName(maxCharsPerWord, maxWidth).String()
	}
}

// FormatNumberValueFormatter formatea los valores numéricos con separadores de miles
func FormatNumberValueFormatter(v any) string {
	// tinystring.Convert puede manejar números directamente
	return tinystring.Convert(v).FormatNumber().String()
}
