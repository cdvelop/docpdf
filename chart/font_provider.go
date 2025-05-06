package chart

import (
	"github.com/cdvelop/docpdf/fontengine"
	"github.com/cdvelop/docpdf/freetype/truetype"
)

// TrueTypeFontAdapter adapta un *truetype.Font a la interfaz fontengine.FontProvider
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
func NewTrueTypeFontAdapter(font *truetype.Font, name, family, weight, style, path string) fontengine.FontProvider {
	return &TrueTypeFontAdapter{
		Font:       font,
		FontName:   name,
		FontFamily: family,
		FontWeight: weight,
		FontStyle:  style,
		FontPath:   path,
	}
}

// Implementación de la interfaz fontengine.FontProvider
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
