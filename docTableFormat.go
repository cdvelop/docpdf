package docpdf

import (
	"strconv"
	"strings"

	"github.com/cdvelop/docpdf/config"
	"github.com/cdvelop/docpdf/style"
)

// tableFormat represents the formatting options for a table header
type tableFormat struct {
	HeaderTitle     string           // The displayed title text
	HeaderAlignment config.Alignment // Alignment of the header text (config.Left, config.Center, config.Right)
	ColumnAlignment config.Alignment // Alignment of the column values (config.Left, config.Center, config.Right)
	Prefix          string           // Prefix to add before each value in the column
	Suffix          string           // Suffix to add after each value in the column
	Width           float64          // Width of the column (0 = auto)
	WidthMode       tableWidthMode   // Table width mode: auto, fixed, or percent
}

// parseTableFormat parses a header, column string with format options
// Format: "headerTitle|option1,option2,option3,..."
//
// Examples:
//
//	"Name" -> {HeaderTitle: "Name", HeaderAlignment: config.Center, ColumnAlignment: config.Left}
//	"Price|HR,CR" -> {HeaderTitle: "Price", HeaderAlignment: config.Right, ColumnAlignment: config.Right}
//	"Amount|HR,CR,P:$" -> {HeaderTitle: "Amount", HeaderAlignment: config.Right, ColumnAlignment: config.Right, Prefix: "$"}
//	"Percent|HC,CC,S:%" -> {HeaderTitle: "Percent", HeaderAlignment: config.Center, ColumnAlignment: config.Center, Suffix: "%"}
//	"Name|HL,CL,W:30%" -> {HeaderTitle: "Name", HeaderAlignment: config.Left, ColumnAlignment: config.Left, Width: 30, WidthMode: percent}
func parseTableFormat(headerStr string) tableFormat {
	result := tableFormat{
		HeaderAlignment: config.Center, // Default header alignment is center
		ColumnAlignment: config.Left,   // Default column alignment is left
		WidthMode:       widthModeAuto, // Default width mode is auto (enum)
	}

	// Split by the separator character
	parts := strings.Split(headerStr, "|")
	result.HeaderTitle = parts[0]

	// If there are no format options, return with defaults
	if len(parts) < 2 || parts[1] == "" {
		return result
	}

	// Get the format options
	formatOptions := parts[1]

	// Process each option separately
	options := strings.Split(formatOptions, ",")
	for _, option := range options {
		option = strings.TrimSpace(option)

		// Process header alignment
		if strings.HasPrefix(option, "H") {
			if option == "HL" {
				result.HeaderAlignment = config.Left
			} else if option == "HR" {
				result.HeaderAlignment = config.Right
			} else if option == "HC" {
				result.HeaderAlignment = config.Center
			}
		}

		// Process column alignment
		if strings.HasPrefix(option, "C") {
			if option == "CL" {
				result.ColumnAlignment = config.Left
			} else if option == "CR" {
				result.ColumnAlignment = config.Right
			} else if option == "CC" {
				result.ColumnAlignment = config.Center
			}
		}

		// Process prefix
		if strings.HasPrefix(option, "P:") {
			result.Prefix = option[2:]
		}

		// Process suffix
		if strings.HasPrefix(option, "S:") {
			result.Suffix = option[2:]
		}

		// Process width
		if strings.HasPrefix(option, "W:") {
			widthStr := option[2:]
			if strings.HasSuffix(widthStr, "%") {
				result.WidthMode = widthModePercent
				// Parse percentage value (remove the % sign)
				percentStr := widthStr[:len(widthStr)-1]
				if val, err := strconv.ParseFloat(percentStr, 64); err == nil {
					result.Width = val
				}
			} else {
				// Si no tiene %, se asume un ancho fijo
				result.WidthMode = widthModeFixed
				if val, err := strconv.ParseFloat(widthStr, 64); err == nil {
					result.Width = val
				}
			}
		}
	}

	return result
}

// createDefaultTableStyles crea los estilos predeterminados para una tabla basados en la configuración del documento
// Devuelve los estilos para el encabezado y las celdas normales
func createDefaultTableStyles(doc *Document) (headerStyle, cellStyle style.Cell) {
	// Estilo para el encabezado
	headerStyle = style.Cell{
		Border: style.Border{
			Top:    true,
			Left:   true,
			Bottom: true,
			Right:  true,
			Width:  1.0,
			Color:  style.Color{R: 100, G: 100, B: 100},
		},
		FillColor: style.Color{R: 240, G: 240, B: 240},
		TextColor: style.Color{R: 0, G: 0, B: 0},
		Font:      FontBold,
		FontSize:  doc.textConfig.GetHeader3().Size(),
	}

	// Estilo para las celdas normales
	cellStyle = style.Cell{
		Border: style.Border{
			Top:    true,
			Left:   true,
			Bottom: true,
			Right:  true,
			Width:  0.5,
			Color:  style.Color{R: 200, G: 200, B: 200},
		},
		FillColor: style.Color{R: 255, G: 255, B: 255},
		TextColor: style.Color{R: 0, G: 0, B: 0},
		Font:      FontRegular,
		FontSize:  doc.textConfig.GetNormal().Size,
	}

	return headerStyle, cellStyle
}

// calculateAutoColumnWidth estima el ancho necesario para una columna basado en el texto del encabezado
func calculateAutoColumnWidth(headerTitle string, fontSize float64, cellPadding float64) float64 {
	textWidthFactor := 0.85
	estWidth := float64(len(headerTitle)) * fontSize * textWidthFactor

	// Añadir padding a ambos lados
	colWidth := estWidth + (cellPadding * 2)

	// Asegurar un ancho mínimo razonable
	if colWidth < 40 {
		colWidth = 40
	}

	return colWidth
}

// initializePercentageColumns inicializa columnas con anchos porcentuales
func initializePercentageColumns(doc *Document, headers []string, availableWidth float64, cellPadding float64) []tableColumn {
	columns := make([]tableColumn, len(headers))
	percentages := make([]float64, len(headers))
	totalPercentage := 0.0

	// Primera pasada: obtener porcentajes explícitos
	for i, headerStr := range headers {
		options := parseTableFormat(headerStr)

		// Crear columna con las opciones analizadas
		columns[i] = tableColumn{
			header:      options.HeaderTitle,
			headerAlign: options.HeaderAlignment,
			align:       options.ColumnAlignment,
			prefix:      options.Prefix,
			suffix:      options.Suffix,
		}

		// Guardar porcentaje si está especificado
		if options.WidthMode == widthModePercent {
			percentages[i] = options.Width
			totalPercentage += options.Width
		}
	}

	// Segunda pasada: distribuir porcentaje restante y aplicar
	remainingColumns := 0
	for i := range percentages {
		if percentages[i] == 0 {
			remainingColumns++
		}
	}

	// Si hay columnas sin porcentaje explícito, distribuir el resto equitativamente
	if remainingColumns > 0 && totalPercentage < 100 {
		remainingPercent := (100 - totalPercentage) / float64(remainingColumns)
		for i := range percentages {
			if percentages[i] == 0 {
				percentages[i] = remainingPercent
				totalPercentage += remainingPercent
			}
		}
	}

	// Normalizar porcentajes para que sumen exactamente 100%
	if totalPercentage != 100.0 && totalPercentage > 0 {
		scaleFactor := 100.0 / totalPercentage
		for i := range percentages {
			percentages[i] *= scaleFactor
		}
	}

	// Aplicar porcentajes al ancho disponible
	for i := range columns {
		columns[i].width = (availableWidth * percentages[i]) / 100.0
	}

	return columns
}

// initializeFixedWidthColumns inicializa columnas con anchos fijos
func initializeFixedWidthColumns(doc *Document, headers []string, cellPadding float64) []tableColumn {
	columns := make([]tableColumn, len(headers))

	// Procesar cada encabezado
	for i, headerStr := range headers {
		options := parseTableFormat(headerStr)

		// Crear columna con las opciones analizadas
		columns[i] = tableColumn{
			header:      options.HeaderTitle,
			width:       options.Width,
			headerAlign: options.HeaderAlignment,
			align:       options.ColumnAlignment,
			prefix:      options.Prefix,
			suffix:      options.Suffix,
		}
	}

	return columns
}

// initializeAutoWidthColumns inicializa columnas con anchos automáticos
func initializeAutoWidthColumns(doc *Document, headers []string, cellPadding float64) []tableColumn {
	columns := make([]tableColumn, len(headers))

	// Configurar la fuente para estimaciones de ancho
	doc.SetFont(FontBold, "", doc.textConfig.GetHeader3().Size)

	// Procesar cada encabezado
	for i, headerStr := range headers {
		options := parseTableFormat(headerStr)

		// Determinar ancho de columna
		var colWidth float64
		if options.WidthMode == widthModeFixed {
			colWidth = options.Width
		} else {
			colWidth = calculateAutoColumnWidth(options.HeaderTitle, doc.textConfig.GetHeader3().Size, cellPadding)
		}

		// Crear columna con las opciones analizadas
		columns[i] = tableColumn{
			header:      options.HeaderTitle,
			width:       colWidth,
			headerAlign: options.HeaderAlignment,
			align:       options.ColumnAlignment,
			prefix:      options.Prefix,
			suffix:      options.Suffix,
		}
	}

	return columns
}
