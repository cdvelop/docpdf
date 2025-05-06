package chart

import (
	"github.com/cdvelop/docpdf/freetype/truetype"
)

// FontProvider es una interfaz que abstrae las propiedades necesarias de una fuente
// para que el renderizador pueda trabajar con ella independientemente de su implementación
type FontProvider interface {
	// Identificación de la fuente
	Name() string   // Nombre de la fuente
	Family() string // Familia de la fuente

	// Propiedades de estilo
	Weight() string // Peso: regular, bold, etc.
	Style() string  // Estilo: normal, italic, etc.

	// Propiedades para renderizado SVG
	SVGFontID() string // ID para referenciar en SVG

	// Opcionalmente, para sistemas que necesiten la ruta al archivo
	Path() string // Ruta al archivo de la fuente
}

// TrueTypeFontAdapter adapta un *truetype.Font a la interfaz FontProvider
// Esto se usará durante la transición mientras se refactoriza el código que usa freetype
type TrueTypeFontAdapter struct {
	Font       *truetype.Font
	FontName   string
	FontFamily string
	FontWeight string
	FontStyle  string
	FontPath   string
}

// NewTrueTypeFontAdapter crea un nuevo adaptador para truetype.Font
func NewTrueTypeFontAdapter(font *truetype.Font, name, family, weight, style, path string) FontProvider {
	return &TrueTypeFontAdapter{
		Font:       font,
		FontName:   name,
		FontFamily: family,
		FontWeight: weight,
		FontStyle:  style,
		FontPath:   path,
	}
}

// Implementación de la interfaz FontProvider
func (a *TrueTypeFontAdapter) Name() string   { return a.FontName }
func (a *TrueTypeFontAdapter) Family() string { return a.FontFamily }
func (a *TrueTypeFontAdapter) Weight() string { return a.FontWeight }
func (a *TrueTypeFontAdapter) Style() string  { return a.FontStyle }
func (a *TrueTypeFontAdapter) Path() string   { return a.FontPath }
func (a *TrueTypeFontAdapter) SVGFontID() string {
	id := a.FontFamily
	if a.FontWeight != "" && a.FontWeight != "regular" {
		id += "-" + a.FontWeight
	}
	if a.FontStyle != "" && a.FontStyle != "normal" {
		id += "-" + a.FontStyle
	}
	return id
}
