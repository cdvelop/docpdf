package docpdf

import (
	"strings"
	"unicode"

	"github.com/cdvelop/docpdf/errs"
)

// justifiedText representa un texto justificado con sus espacios ajustados
type justifiedText struct {
	words       []string
	spaces      []float64
	originalStr string
	width       float64
}

// parseTextForJustification divide un texto en sus palabras y calcula los espacios necesarios
func parseTextForJustification(gp *pdfEngine, text string, width float64) (*justifiedText, error) {
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
		return &justifiedText{
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

		return &justifiedText{
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

	return &justifiedText{
		words:       words,
		spaces:      spaces,
		originalStr: text,
		width:       width,
	}, nil
}

// drawJustifiedLine dibuja una línea de texto justificado
func drawJustifiedLine(gp *pdfEngine, jText *justifiedText, x, y float64) error {
	if len(jText.words) == 0 {
		return nil
	}

	currentX := x

	// Si solo hay una palabra, simplemente la dibujamos sin justificar
	if len(jText.words) == 1 {
		return gp.Cell(&Rect{W: jText.width, H: 0}, jText.words[0])
	}

	// Guardar el estado actual
	originalX := gp.GetX()
	originalY := gp.GetY()

	// Dibujar cada palabra con el espacio calculado
	for i, word := range jText.words {
		gp.SetX(currentX)
		gp.SetY(y)

		err := gp.Cell(&Rect{W: 0, H: 0}, word)
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

// MultiCellJustified dibuja texto justificado dentro de un rectángulo
func (gp *pdfEngine) MultiCellJustified(rectangle *Rect, text string) error {
	opt := cellOption{
		Align: Justify,
	}
	return gp.MultiCellWithOption(rectangle, text, opt)
}
