package chart

import (
	"github.com/cdvelop/docpdf/chart/roboto"
	"github.com/cdvelop/docpdf/freetype/truetype"
)

// Variable global para almacenar el motor predeterminado
// cuando se llama a GetDefaultFont sin inicialización explícita
var defaultEngine *ChartEngine

// GetDefaultFont returns the default font (Roboto-Medium).
// Esta función ahora inicializa un ChartEngine si no existe uno
// para mantener compatibilidad con el código existente.
func GetDefaultFont() (*truetype.Font, error) {
	// Si ya tenemos un motor inicializado con una fuente, lo usamos
	if defaultEngine != nil && defaultEngine.defaultFont != nil {
		return defaultEngine.defaultFont, nil
	}

	// Si no hay motor, lo inicializamos con la fuente Roboto por defecto
	var err error
	defaultEngine, err = NewEngine(roboto.Roboto)
	if err != nil {
		return nil, err
	}

	return defaultEngine.defaultFont, nil
}
