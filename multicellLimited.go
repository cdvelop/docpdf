package docpdf

import "github.com/cdvelop/tinystring"

// MultiCellWithOptionAndMaxLines dibuja texto multi-línea con un límite máximo de líneas
// El texto se truncará automáticamente si excede el número máximo de líneas
func (gp *pdfEngine) MultiCellWithOptionAndMaxLines(rectangle *Rect, text string, opt cellOption, maxLines int) error {
	// Validar entrada
	if maxLines <= 0 {
		maxLines = 1 // Al menos una línea
	}

	// Asegurarse de que las opciones de ruptura estén configuradas para no cortar palabras
	if opt.breakOption == nil {
		opt.breakOption = &defaultBreakOption
	}

	// Asegurarse que siempre use el modo sensible a espacios para evitar cortar palabras
	originalMode := opt.breakOption.Mode
	originalIndicator := opt.breakOption.BreakIndicator

	opt.breakOption.Mode = breakModeIndicatorSensitive
	opt.breakOption.BreakIndicator = ' '

	// Dividir el texto en líneas según el ancho
	textSplits, err := gp.SplitTextWithOption(text, rectangle.W, opt.breakOption)
	if err != nil {
		return err
	}

	// Restaurar opciones originales
	opt.breakOption.Mode = originalMode
	opt.breakOption.BreakIndicator = originalIndicator

	// Si el número de líneas excede el máximo, truncar el texto
	if len(textSplits) > maxLines {
		// Tomamos solo las primeras líneas según maxLines
		limitedLines := textSplits[:maxLines]

		// Si estamos truncando y tenemos al menos una línea para mostrar
		if maxLines >= 1 && len(limitedLines) > 0 {
			// Si tenemos más de una línea, trabajamos con la última para añadir puntos suspensivos
			if maxLines > 1 {
				lastLineIdx := maxLines - 1
				lastLine := limitedLines[lastLineIdx]

				// Calcular el ancho disponible y cuántos caracteres caben con los puntos suspensivos
				ellipsis := "..."

				// Medir ancho de puntos suspensivos
				ellipsisWidth, _ := gp.MeasureTextWidth(ellipsis)

				// Calcular ancho disponible para el texto sin los puntos suspensivos
				availableWidth := rectangle.W - ellipsisWidth

				// Estimar caracteres basado en tamaño de fuente actual
				// Usamos un factor conservador (0.55) para chars por punto de ancho
				charWidthEstimate := gp.curr.FontSize * 0.55
				if charWidthEstimate <= 0 {
					charWidthEstimate = 1
				}

				// Usar la biblioteca tinystring para truncar la línea correctamente
				// Reservamos 3 caracteres para "..."
				truncatedLastLine := tinystring.Convert(lastLine).
					Truncate(int(availableWidth/charWidthEstimate), 3).
					String()

				limitedLines[lastLineIdx] = truncatedLastLine
			}
		}

		// Unir las líneas truncadas
		truncatedText := ""
		for i, line := range limitedLines {
			if i > 0 {
				truncatedText += "\n"
			}
			truncatedText += line
		}

		// Usar el texto truncado
		text = truncatedText
	}

	// Dibujar el texto utilizando la función MultiCellWithOption estándar
	return gp.MultiCellWithOption(rectangle, text, opt)
}
