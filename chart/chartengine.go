package chart

import (
	"github.com/cdvelop/docpdf/freetype/truetype"
)

// ChartEngine centraliza la configuración y recursos para todos los gráficos
type ChartEngine struct {
	// Fuente predeterminada inicializada una sola vez deprecada posteriormente
	defaultFont *truetype.Font

	// Configuraciones comunes para todos los gráficos
	widthDefault  int     // ancho predeterminado
	heightDefault int     // altura predeterminada
	dpiDefault    float64 // DPI predeterminado para gráficos

	// Estilos predeterminados
	Background Style
	Canvas     Style

	// Paleta de colores
	ColorPalette ColorPalette
}

// NewEngine crea una nueva instancia de ChartEngine con la fuente proporcionada.
// Este método debe ser llamado una sola vez para centralizar la inicialización.
func NewEngine(fontBytes []byte) (*ChartEngine, error) {
	// Creamos una nueva instancia
	engine := &ChartEngine{
		// Valores predeterminados
		widthDefault:  800,
		heightDefault: 600,
		dpiDefault:    DefaultDPI,
		Background:    Style{},
		Canvas:        Style{},
		ColorPalette:  AlternateColorPalette,
	}

	// Inicializamos la fuente desde los bytes proporcionados
	if fontBytes != nil {
		font, err := truetype.Parse(fontBytes)
		if err != nil {
			return nil, err
		}
		engine.defaultFont = font
	}

	return engine, nil
}

// GetFont devuelve la fuente configurada en el ChartEngine
func (ce *ChartEngine) GetFont() *truetype.Font {
	return ce.defaultFont
}

// DonutChart crea un nuevo DonutChart usando la configuración del motor
func (ce *ChartEngine) DonutChart(values []Value) *DonutChart {
	return &DonutChart{
		Width:        ce.widthDefault,
		Height:       ce.heightDefault,
		DPI:          ce.dpiDefault,
		ColorPalette: ce.ColorPalette,
		Background:   ce.Background,
		Canvas:       ce.Canvas,
		Font:         ce.defaultFont,
		Values:       values,
	}
}

// PieChart crea un nuevo PieChart usando la configuración del motor
func (ce *ChartEngine) PieChart(values []Value) *PieChart {
	return &PieChart{
		Width:        ce.widthDefault,
		Height:       ce.heightDefault,
		DPI:          ce.dpiDefault,
		ColorPalette: ce.ColorPalette,
		Background:   ce.Background,
		Canvas:       ce.Canvas,
		Font:         ce.defaultFont,
		Values:       values,
	}
}

// BarChart crea un nuevo BarChart usando la configuración del motor
func (ce *ChartEngine) BarChart(barWidth int) *BarChart {
	bc := &BarChart{
		Width:        ce.widthDefault,
		Height:       ce.heightDefault,
		DPI:          ce.dpiDefault,
		BarWidth:     barWidth,
		ColorPalette: ce.ColorPalette,
		Background:   ce.Background,
		Canvas:       ce.Canvas,
		Font:         ce.defaultFont,
	}
	return bc
}

// SetDefaultDimensions configura el tamaño de los gráficos por defecto generados por este motor
func (ce *ChartEngine) SetDefaultDimensions(width, height int) *ChartEngine {
	ce.widthDefault = width
	ce.heightDefault = height
	return ce
}

// WithDPI configura el DPI para los gráficos generados por este motor
func (ce *ChartEngine) SetDefaultDPI(dpi float64) *ChartEngine {
	ce.dpiDefault = dpi
	return ce
}

// WithColorPalette configura la paleta de colores para los gráficos
func (ce *ChartEngine) SetDefaultColorPalette(palette ColorPalette) *ChartEngine {
	ce.ColorPalette = palette
	return ce
}

// WithBackgroundStyle configura el estilo de fondo
func (ce *ChartEngine) SetDefaultBackgroundStyle(style Style) *ChartEngine {
	ce.Background = style
	return ce
}

// WithCanvasStyle configura el estilo del canvas
func (ce *ChartEngine) SetDefaultCanvasStyle(style Style) *ChartEngine {
	ce.Canvas = style
	return ce
}
