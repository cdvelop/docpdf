package docpdf

import (
	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/config"
	"github.com/cdvelop/docpdf/pdfengine"
	"github.com/cdvelop/docpdf/style"
	"github.com/cdvelop/tinystring"
)

// tableWidthMode defines the table width calculation mode: automatic, fixed or percentage
type tableWidthMode int

const (
	widthModeAuto    tableWidthMode = iota // Automatic width based on content
	widthModeFixed                         // Fixed specified widths
	widthModePercent                       // Widths as percentage of available space
)

// docTable represents a table to be added to the document
type docTable struct {
	doc                       *Document
	columns                   []tableColumn
	rows                      [][]tableCell
	width                     float64
	rowHeight                 float64
	maxLinesTextForRowInACell int // Maximum number of text lines per row in a cell
	cellPadding               float64
	headerStyle               style.Cell
	cellStyle                 style.Cell
	alignment                 config.Alignment // config.Left, config.Center, or config.Right alignment
	currentWidth              float64
}

// tableColumn represents a column in the table
type tableColumn struct {
	header      string           // Header text
	width       float64          // Width of the column
	headerAlign config.Alignment // Text alignment for header (config.Left, config.Center, config.Right)
	align       config.Alignment // Text alignment within cells (config.Left, config.Center, config.Right)
	prefix      string           // Prefix to add before each value in the column
	suffix      string           // Suffix to add after each value in the column
}

// tableCell represents a cell in the table
type tableCell struct {
	content      string     // Content of the cell
	useCellStyle bool       // If true, use custom style instead of the table's default
	cellStyle    style.Cell // Custom style for this cell
}

// NewTable creates a new table with the specified headers
// Headers can include formatting options in the following format.
//
// Format: "headerTitle|option1,option2,option3,..."
//
// Available options:
//   - Header alignment: "HL" (left), "HC" (center, default), "HR" (right)
//   - Column alignment: "CL" (left, default), "CC" (center), "CR" (right)
//   - Prefix/Suffix: "P:value" (adds prefix), "S:value" (adds suffix)
//   - Width: "W:number" (fixed width), "W:number%" (percentage width)
//     Note: If width is not specified, auto width is used by default
//
// Examples:
//   - "Name" - Normal header, left-aligned column, auto width
//   - "Price|HR,CR" - config.Right-aligned header, right-aligned column
//   - "Price|HR,CR,P:$" - config.Right-aligned header, right-aligned column with "$" prefix
//   - "Percentage|HC,CC,S:%" - config.Center-aligned header, center-aligned column with "%" suffix
//   - "Name|HL,CL,W:30%" - config.Left-aligned header, left-aligned column with 30% of available width
//   - "Age|HC,CR,W:20" - config.Center-aligned header, right-aligned column with fixed width of 20 units
func (doc *Document) NewTable(headers ...string) *docTable {
	// Almacenar los encabezados para verificación posterior
	doc.lastTableHeaders = headers

	// Crear una nueva tabla con configuración predeterminada
	table := &docTable{
		doc:                       doc,
		rowHeight:                 25, // Default row height
		cellPadding:               5,  // Default padding
		alignment:                 config.Center,
		maxLinesTextForRowInACell: 2, // Limitamos a máximo 2 líneas de texto por fila en cada celda
	}

	// Crear los estilos predeterminados para la tabla
	table.headerStyle, table.cellStyle = createDefaultTableStyles(doc)

	// Detectar el modo de ancho basado en los encabezados
	widthMode := detectWidthModeFromHeaders(headers)

	// Procesar los encabezados y calcular los anchos iniciales según su tipo
	var columns []tableColumn

	switch widthMode {
	case widthModePercent:
		// Si hay porcentajes, todas las columnas son tratadas como porcentaje
		// Y la tabla debe ocupar exactamente el ancho disponible
		columns = initializePercentageColumns(doc, headers, doc.contentAreaWidth, table.cellPadding)
		table.width = doc.contentAreaWidth

	case widthModeFixed:
		// Si todas son anchos fijos, respetamos esos anchos exactamente
		columns = initializeFixedWidthColumns(doc, headers, table.cellPadding)
		// Calcular el ancho total sumando los anchos de las columnas
		totalWidth := 0.0
		for _, col := range columns {
			totalWidth += col.width
		}
		table.width = totalWidth

	default: // widthModeAuto
		// Calculamos los anchos automáticamente basados en el contenido
		columns = initializeAutoWidthColumns(doc, headers, table.cellPadding)
		// Calcular el ancho total sumando los anchos de las columnas
		totalWidth := 0.0
		for _, col := range columns {
			totalWidth += col.width
		}

		// Si el ancho total excede el disponible, escalar proporcionalmente
		if totalWidth > doc.contentAreaWidth {
			scaleFactor := doc.contentAreaWidth / totalWidth
			for i := range columns {
				columns[i].width *= scaleFactor
			}
			table.width = doc.contentAreaWidth
		} else {
			table.width = totalWidth
		}
	}

	// Establecer las columnas
	table.columns = columns

	return table
}

// detectWidthModeFromHeaders analiza los encabezados para determinar el modo de ancho
func detectWidthModeFromHeaders(headers []string) tableWidthMode {
	// Si algún encabezado usa porcentaje, usamos modo porcentaje
	for _, headerStr := range headers {
		options := parseTableFormat(headerStr)
		if options.WidthMode == widthModePercent {
			return widthModePercent
		}
	}

	// Si todos los encabezados usan ancho fijo, usamos modo fijo
	allFixed := true
	for _, headerStr := range headers {
		options := parseTableFormat(headerStr)
		if options.WidthMode != widthModeFixed {
			allFixed = false
			break
		}
	}

	if allFixed {
		return widthModeFixed
	}

	// Por defecto, usamos modo automático
	return widthModeAuto
}

// Width sets the total width of the table
func (t *docTable) Width(width float64) *docTable {
	if width > 0 {
		scaleFactor := width / t.currentWidth
		for i := range t.columns {
			t.columns[i].width *= scaleFactor
		}
		t.width = width
		t.currentWidth = width
	}
	return t
}

// RowHeight sets the height of rows
func (t *docTable) RowHeight(height float64) *docTable {
	if height > 0 {
		t.rowHeight = height
	}
	return t
}

// SetColumnWidth sets the width of a specific column by index
func (t *docTable) SetColumnWidth(columnIndex int, width float64) *docTable {
	if columnIndex >= 0 && columnIndex < len(t.columns) && width > 0 {
		// Adjust total width
		t.currentWidth = t.currentWidth - t.columns[columnIndex].width + width
		t.columns[columnIndex].width = width
	}
	return t
}

// SetHeaderAlignment sets the text alignment for a specific header
func (t *docTable) SetHeaderAlignment(columnIndex int, alignment config.Alignment) *docTable {
	if columnIndex >= 0 && columnIndex < len(t.columns) {
		t.columns[columnIndex].headerAlign = alignment
	}
	return t
}

// SetColumnPrefix sets a prefix for all values in a column
func (t *docTable) SetColumnPrefix(columnIndex int, prefix string) *docTable {
	if columnIndex >= 0 && columnIndex < len(t.columns) {
		t.columns[columnIndex].prefix = prefix
	}
	return t
}

// SetColumnSuffix sets a suffix for all values in a column
func (t *docTable) SetColumnSuffix(columnIndex int, suffix string) *docTable {
	if columnIndex >= 0 && columnIndex < len(t.columns) {
		t.columns[columnIndex].suffix = suffix
	}
	return t
}

// AlignLeft aligns the table to the left margin
func (t *docTable) AlignLeft() *docTable {
	t.alignment = config.Left
	return t
}

// AlignCenter centers the table horizontally (default)
func (t *docTable) AlignCenter() *docTable {
	t.alignment = config.Center
	return t
}

// AlignRight aligns the table to the right margin
func (t *docTable) AlignRight() *docTable {
	t.alignment = config.Right
	return t
}

// HeaderStyle sets the style for the header row
func (t *docTable) HeaderStyle(s style.Cell) *docTable {
	t.headerStyle = s
	return t
}

// style.Cell sets the default style for regular cells
func (t *docTable) Cell(s style.Cell) *docTable {
	t.cellStyle = s
	return t
}

// AddRow adds a row of data to the table
// Accepts any value type which will be converted to strings
func (t *docTable) AddRow(cells ...any) *docTable {
	rowCells := make([]tableCell, len(cells))
	for i, content := range cells {
		formattedContent := tinystring.Convert(content).String()

		// Apply prefix and suffix if column exists
		if i < len(t.columns) {
			if t.columns[i].prefix != "" {
				formattedContent = t.columns[i].prefix + formattedContent
			}
			if t.columns[i].suffix != "" {
				formattedContent = formattedContent + t.columns[i].suffix
			}
		}

		rowCells[i] = tableCell{
			content:      formattedContent,
			useCellStyle: false,
		}
	}
	t.rows = append(t.rows, rowCells)
	return t
}

// AddStyledRow adds a row with individually styled cells
func (t *docTable) AddStyledRow(cells ...styledCell) *docTable {
	rowCells := make([]tableCell, len(cells))
	for i, cell := range cells {
		formattedContent := cell.Content

		// Apply prefix and suffix if column exists
		if i < len(t.columns) {
			if t.columns[i].prefix != "" {
				formattedContent = t.columns[i].prefix + formattedContent
			}
			if t.columns[i].suffix != "" {
				formattedContent = formattedContent + t.columns[i].suffix
			}
		}

		rowCells[i] = tableCell{
			content:      formattedContent,
			useCellStyle: true,
			cellStyle:    cell.Style,
		}
	}
	t.rows = append(t.rows, rowCells)
	return t
}

// styledCell represents a cell with custom styling
type styledCell struct {
	Content string
	Style   style.Cell
}

// NewStyledCell creates a new cell with custom styling
func (doc *Document) NewStyledCell(content string, style style.Cell) styledCell {
	return styledCell{
		Content: content,
		Style:   style,
	}
}

// Draw renders the table on the document
func (t *docTable) Draw() error {
	// Calculate starting X config.Alignment
	x := t.calculatePosition()

	// Colección para guardar información de los encabezados para dibujar sus bordes al final
	type headerInfo struct {
		x, y, width, height float64
		style               style.Cell
	}
	headers := []headerInfo{}

	// Obtenemos la posición Y inicial, considerando el espacio disponible en la página actual
	y := t.doc.CurrentPdf().Y

	// Verificar si al menos la cabecera cabe en la página actual
	headerHeight := t.rowHeight
	bottomMargin := t.doc.Margins().Bottom

	// Calcular espacio disponible en la página actual
	availableSpace := t.doc.CurrentPdf().PageSize().H - y - bottomMargin

	// Si ni siquiera la cabecera cabe, ir a una nueva página
	if headerHeight > availableSpace {
		t.doc.AddPage()
		y = t.doc.CurrentPdf().Y
	}

	// Dibujar encabezado
	currentX := x
	for _, col := range t.columns {
		headers = append(headers, headerInfo{
			x:      currentX,
			y:      y,
			width:  col.width,
			height: t.rowHeight,
			style:  t.headerStyle,
		})

		t.drawCellContent(
			currentX,
			y,
			col.width,
			t.rowHeight,
			col.header,
			col.headerAlign,
			t.headerStyle,
		)
		currentX += col.width
	}

	currentY := y + headerHeight

	// Dibujar los bordes de los encabezados
	for _, h := range headers {
		t.drawCellBorder(h.x, h.y, h.width, h.height, h.style.Border)
	}
	// Dibujar filas de datos
	for _, row := range t.rows {
		// Verificar si la fila actual cabe en la página actual
		availableSpace = t.doc.CurrentPdf().PageSize().H - currentY - bottomMargin

		// Si esta fila no cabe en la página actual, crear una nueva página
		if t.rowHeight > availableSpace {
			t.doc.AddPage()
			currentY = t.doc.CurrentPdf().Y

			// Limpiar la lista de encabezados para la nueva página
			headers = []headerInfo{}

			// Redibujar la fila de encabezados en la nueva página
			headerY := currentY
			headerX := x
			for _, col := range t.columns {
				headers = append(headers, headerInfo{
					x:      headerX,
					y:      headerY,
					width:  col.width,
					height: t.rowHeight,
					style:  t.headerStyle,
				})

				t.drawCellContent(
					headerX,
					headerY,
					col.width,
					t.rowHeight,
					col.header,
					col.headerAlign,
					t.headerStyle,
				)
				headerX += col.width
			}

			// Dibujar los bordes de los encabezados
			for _, h := range headers {
				t.drawCellBorder(h.x, h.y, h.width, h.height, h.style.Border)
			}

			// Ajustar currentY para empezar debajo del encabezado
			currentY += t.rowHeight
		}

		// Dibujar la fila actual
		currentX = x
		for colIndex, cell := range row {
			if colIndex < len(t.columns) {
				cellWidth := t.columns[colIndex].width
				cellAlign := t.columns[colIndex].align

				style := t.cellStyle
				if cell.useCellStyle {
					style = cell.cellStyle
				}

				t.drawCell(
					currentX,
					currentY,
					cellWidth,
					t.rowHeight,
					cell.content,
					cellAlign,
					style,
				)

				currentX += cellWidth
			}
		}

		// Avanzar a la siguiente fila
		currentY += t.rowHeight
	}

	// Actualizar la posición del documento para que sea después de la tabla
	t.doc.SetY(currentY + t.doc.fontConfig.Normal.SpaceAfter)

	return nil
}

// drawCellContent dibuja solo el contenido y fondo de una celda (sin bordes)
func (t *docTable) drawCellContent(
	x float64,
	y float64,
	width float64,
	height float64,
	content string,
	align config.Alignment,
	stCell style.Cell,
) {
	// Fill the cell background if a fill color is specified
	if (stCell.FillColor != style.Color{}) {
		t.doc.SetFillColor(stCell.FillColor.R, stCell.FillColor.G, stCell.FillColor.B)
		t.doc.RectFromUpperLeftWithStyle(x, y, width, height, "F")
	}

	// Set text properties
	if stCell.Font != "" {
		t.doc.SetFont(stCell.Font, "", stCell.FontSize)
	}
	t.doc.SetTextColor(stCell.TextColor.R, stCell.TextColor.G, stCell.TextColor.B)

	// Calculate available text area considering padding
	textX := x + t.cellPadding
	textWidth := width - (2 * t.cellPadding)

	// Para textos con más de una línea, usamos justificación para mejorar la apariencia
	cellAlign := align

	// Estimar el número de líneas que tendría el texto
	lines, _ := t.doc.SplitTextWithWordWrap(content, textWidth)
	numLines := len(lines)

	if numLines >= 2 {
		cellAlign = config.Justify // Justificar texto multilínea
	}

	// Calcular altura de línea y posicionamiento vertical
	lineHeight := t.doc.CurrentPdf().FontSize * 1.2

	// Altura estimada del texto (con límite de 2 líneas)
	estimatedLines := min(numLines, t.maxLinesTextForRowInACell)
	totalTextHeight := float64(estimatedLines) * lineHeight

	// Centrar verticalmente el texto en la celda
	verticalPadding := max((height-totalTextHeight)/2, 0)

	// Pequeño ajuste para textos de 2 líneas para compensar descendentes de fuente
	bottomSpacingAdjustment := 0.0
	if estimatedLines == 2 {
		bottomSpacingAdjustment = lineHeight * 0.1
	}
	drawY := max(y+verticalPadding-bottomSpacingAdjustment, y)

	// Configurar opciones para dibujar la celda
	cellOpt := pdfengine.CellOption{
		Align:         cellAlign,
		Border:        0, // No borders for content drawing
		TruncateLines: t.maxLinesTextForRowInACell,
		BreakOption: &pdfengine.BreakOption{
			Mode:           pdfengine.BreakModeIndicatorSensitive,
			BreakIndicator: ' ',
			Separator:      "",
		},
	}

	// Dibujar el texto utilizando el método MultiCellWithOption
	// que se encarga automáticamente de truncar el texto si es necesario
	t.doc.SetXY(textX, drawY)
	drawRectH := totalTextHeight + (lineHeight * 0.1)
	err := t.doc.MultiCellWithOption(
		&canvas.Rect{W: textWidth, H: drawRectH},
		content,
		cellOpt,
	)

	if err != nil && err.Error() != "empty string" {
		t.doc.Log("Error drawing table cell content:", err)
	}
}

// drawCellBorder dibuja solo los bordes de una celda
func (t *docTable) drawCellBorder(
	x float64,
	y float64,
	width float64,
	height float64,
	borderStyle style.Border,
) {
	if borderStyle.Width > 0 {
		t.doc.SetLineWidth(borderStyle.Width)
		t.doc.SetStrokeColor(
			borderStyle.Color.R,
			borderStyle.Color.G,
			borderStyle.Color.B,
		)

		if borderStyle.Top {
			t.doc.Line(x, y, x+width, y)
		}
		if borderStyle.Bottom {
			t.doc.Line(x, y+height, x+width, y+height)
		}
		if borderStyle.Left {
			t.doc.Line(x, y, x, y+height)
		}
		if borderStyle.Right {
			t.doc.Line(x+width, y, x+width, y+height)
		}
	}
}

// drawCell draws a single cell of the table
func (t *docTable) drawCell(
	x float64,
	y float64,
	width float64,
	height float64,
	content string,
	align config.Alignment,
	style style.Cell,
) {
	// Primero dibujamos el contenido y el fondo
	t.drawCellContent(x, y, width, height, content, align, style)

	// Luego dibujamos los bordes
	t.drawCellBorder(x, y, width, height, style.Border)
}

// calculatePosition determina donde colocar la tabla
func (t *docTable) calculatePosition() float64 {
	// Posición X inicial (margen izquierdo)
	x := t.doc.Margins().Left

	// Aplicar alineación
	switch t.alignment {
	case config.Center:
		// Centrar la tabla: margen izquierdo + (espacio disponible - ancho tabla) / 2
		x = t.doc.Margins().Left + (t.doc.contentAreaWidth-t.width)/2
	case config.Right:
		// Alinear a la derecha: margen izquierdo + espacio disponible - ancho tabla
		x = t.doc.Margins().Left + t.doc.contentAreaWidth - t.width
	case config.Left:
		// Alinear a la izquierda: simplemente el margen izquierdo
		x = t.doc.Margins().Left
	}

	return x
}
