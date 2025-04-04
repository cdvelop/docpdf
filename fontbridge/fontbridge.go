package fontbridge

import (
	"os"

	"github.com/cdvelop/docpdf/chart"
	"github.com/cdvelop/docpdf/drawing"
	"github.com/cdvelop/docpdf/freetype/truetype"
)

// FontConfig contiene la configuración de fuente compartida entre docpdf y chart
type FontConfig struct {
	// Font es la fuente cargada desde TrueType
	Font *truetype.Font

	// Tamaños de fuente para diferentes elementos
	TitleSize      float64 // Para títulos de gráficos (Header1)
	AxisLabelSize  float64 // Para etiquetas de ejes (Normal)
	ValueLabelSize float64 // Para etiquetas de valores (Normal)
	LegendSize     float64 // Para leyendas (Normal o Footnote)

	// Colores
	TitleColor      drawing.Color
	AxisLabelColor  drawing.Color
	ValueLabelColor drawing.Color
	LegendColor     drawing.Color

	// Espaciado
	LineSpacing float64
}

// Instancia global de configuración
var SharedFontConfig *FontConfig

// Inicializar con valores por defecto
func init() {
	// Estos valores se actualizarán con InitFromDocConfig
	SharedFontConfig = &FontConfig{
		Font:            nil,
		TitleSize:       14,
		AxisLabelSize:   11,
		ValueLabelSize:  11,
		LegendSize:      9,
		TitleColor:      drawing.Color{R: 0, G: 0, B: 0, A: 255},       // Negro
		AxisLabelColor:  drawing.Color{R: 0, G: 0, B: 0, A: 255},       // Negro
		ValueLabelColor: drawing.Color{R: 0, G: 0, B: 0, A: 255},       // Negro
		LegendColor:     drawing.Color{R: 128, G: 128, B: 128, A: 255}, // Gris
		LineSpacing:     1.2,
	}
}

// InitFromDocConfig actualiza la configuración compartida desde docpdf.FontConfig
func InitFromDocConfig(fontPath string, fontFile string, titleSize, normalSize, footnoteSize float64,
	titleColor, normalColor, footnoteColor drawing.Color, lineSpacing float64) error {

	// Si ya tenemos una fuente cargada, no la cargamos de nuevo
	if SharedFontConfig.Font == nil {
		// Cargar la fuente desde el sistema de archivos
		font, err := LoadFont(fontPath, fontFile)
		if err != nil {
			// Si no podemos cargar la fuente configurada, intentamos usar la fuente por defecto de chart
			defaultFont, defaultErr := chart.GetDefaultFont()
			if defaultErr != nil {
				return defaultErr // Si ni siquiera podemos cargar la fuente por defecto, devolvemos el error
			}
			SharedFontConfig.Font = defaultFont
		} else {
			SharedFontConfig.Font = font
		}
	}

	// Actualizar tamaños
	SharedFontConfig.TitleSize = titleSize
	SharedFontConfig.AxisLabelSize = normalSize
	SharedFontConfig.ValueLabelSize = normalSize
	SharedFontConfig.LegendSize = footnoteSize

	// Actualizar colores
	SharedFontConfig.TitleColor = titleColor
	SharedFontConfig.AxisLabelColor = normalColor
	SharedFontConfig.ValueLabelColor = normalColor
	SharedFontConfig.LegendColor = footnoteColor

	// Actualizar espaciado
	SharedFontConfig.LineSpacing = lineSpacing

	return nil
}

// LoadFont carga una fuente desde el sistema de archivos
func LoadFont(fontPath string, fontFile string) (*truetype.Font, error) {
	fullPath := fontPath + fontFile
	fontBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	return font, nil
}

// ApplyToChartStyle aplica la configuración de fontbridge a los estilos de chart
func ApplyToChartStyle(style *chart.Style, fontType string) {
	style.Font = SharedFontConfig.Font

	switch fontType {
	case "title":
		style.FontSize = SharedFontConfig.TitleSize
		style.FontColor = SharedFontConfig.TitleColor
	case "axis":
		style.FontSize = SharedFontConfig.AxisLabelSize
		style.FontColor = SharedFontConfig.AxisLabelColor
	case "value":
		style.FontSize = SharedFontConfig.ValueLabelSize
		style.FontColor = SharedFontConfig.ValueLabelColor
	case "legend":
		style.FontSize = SharedFontConfig.LegendSize
		style.FontColor = SharedFontConfig.LegendColor
	default:
		// Usar configuración Normal por defecto
		style.FontSize = SharedFontConfig.AxisLabelSize
		style.FontColor = SharedFontConfig.AxisLabelColor
	}
}

// GetDrawingColor convierte RGBColor de docpdf a drawing.Color para chart
func GetDrawingColor(r, g, b uint8) drawing.Color {
	return drawing.Color{
		R: r,
		G: g,
		B: b,
		A: 255, // Opaco por defecto
	}
}
