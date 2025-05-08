package pdfengine

import (
	"strings"
	"unicode"

	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/config"
	"github.com/cdvelop/docpdf/errs"
)

// JustifiedText representa un texto justificado con sus espacios ajustados
type JustifiedText struct {
	words       []string
	spaces      []float64
	originalStr string
	width       float64
}

// GetWords retorna las palabras del texto justificado
func (jt *JustifiedText) GetWords() []string {
	return jt.words
}

// GetSpaces retorna los espacios entre palabras del texto justificado
func (jt *JustifiedText) GetSpaces() []float64 {
	return jt.spaces
}

// GetOriginalString retorna el texto original antes de ser justificado
func (jt *JustifiedText) GetOriginalString() string {
	return jt.originalStr
}

// GetWidth retorna el ancho total del texto justificado
func (jt *JustifiedText) GetWidth() float64 {
	return jt.width
}

// WordCount retorna el número de palabras del texto justificado
func (jt *JustifiedText) WordCount() int {
	return len(jt.words)
}

// SpaceCount retorna el número de espacios del texto justificado
func (jt *JustifiedText) SpaceCount() int {
	return len(jt.spaces)
}

// ParseTextForJustification divide un texto en sus palabras y calcula los espacios necesarios
func (gp *PdfEngine) ParseTextForJustification(text string, width float64) (*JustifiedText, error) {
	// Si el texto está vacío o no tiene espacios, no hay nada que justificar
	if text == "" {
		return nil, errs.EmptyString
	}

	// Ignorar espacios iniciales y finales
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errs.EmptyString
	}

	// Dividir el texto en palabras
	words := strings.FieldsFunc(text, unicode.IsSpace)
	if len(words) <= 1 {
		// No hay suficientes palabras para justificar
		return &JustifiedText{
			words:       words,
			spaces:      []float64{0},
			originalStr: text,
			width:       width,
		}, nil
	}

	// Calcular el ancho de cada palabra y el ancho total de las palabras
	wordsWidth := 0.0
	for _, word := range words {
		w, err := gp.MeasureTextWidth(word)
		if err != nil {
			return nil, err
		}
		wordsWidth += w
	}

	// Calcular el ancho normal de un espacio
	normalSpaceWidth, err := gp.MeasureTextWidth(" ")
	if err != nil {
		return nil, err
	}

	// Calcular el espacio disponible para distribuir
	spaceCount := len(words) - 1
	availableSpace := width - wordsWidth

	// Si el espacio disponible es negativo, usar espacios normales
	if availableSpace < 0 {
		// Crear array de espacios normales
		spaces := make([]float64, spaceCount)
		for i := range spaces {
			spaces[i] = normalSpaceWidth
		}

		return &JustifiedText{
			words:       words,
			spaces:      spaces,
			originalStr: text,
			width:       width,
		}, nil
	}

	// Calcular el ancho de cada espacio
	spaceWidth := availableSpace / float64(spaceCount)

	// Si el espacio calculado es menor que un espacio normal, usar el espacio normal
	if spaceWidth < normalSpaceWidth {
		spaceWidth = normalSpaceWidth
	}

	// Crear array de espacios (todos iguales en este caso)
	spaces := make([]float64, spaceCount)
	for i := range spaces {
		spaces[i] = spaceWidth
	}

	return &JustifiedText{
		words:       words,
		spaces:      spaces,
		originalStr: text,
		width:       width,
	}, nil
}

// drawJustifiedLine dibuja una línea de texto justificado
func drawJustifiedLine(gp *PdfEngine, jText *JustifiedText, x, y float64) error {
	if len(jText.words) == 0 {
		return nil
	}

	currentX := x

	// Si solo hay una palabra, simplemente la dibujamos sin justificar
	if len(jText.words) == 1 {
		return gp.Cell(&canvas.Rect{W: jText.width, H: 0}, jText.words[0])
	}

	// Guardar el estado actual
	originalX := gp.GetX()
	originalY := gp.GetY()

	// Dibujar cada palabra con el espacio calculado
	for i, word := range jText.words {
		gp.SetX(currentX)
		gp.SetY(y)

		err := gp.Cell(&canvas.Rect{W: 0, H: 0}, word)
		if err != nil {
			return err
		}

		if i < len(jText.words)-1 {
			wordWidth, err := gp.MeasureTextWidth(word)
			if err != nil {
				return err
			}
			currentX += wordWidth + jText.spaces[i]
		}
	}

	// Restaurar la posición
	gp.SetX(originalX)
	gp.SetY(originalY)

	return nil
}

// Draw dibuja el texto justificado en las coordenadas especificadas
func (jt *JustifiedText) Draw(gp *PdfEngine, x, y float64) error {
	return drawJustifiedLine(gp, jt, x, y)
}

// MultiCellJustified dibuja texto justificado dentro de un rectángulo
func (gp *PdfEngine) MultiCellJustified(rectangle *canvas.Rect, text string) error {
	opt := CellOption{
		Align: config.Justify,
	}
	return gp.MultiCellWithOption(rectangle, text, opt)
}
