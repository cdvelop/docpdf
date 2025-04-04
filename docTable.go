package docpdf

import "strconv"

// tableWidthMode defines the table width calculation mode: automatic, fixed or percentage
type tableWidthMode int

const (
	widthModeAuto    tableWidthMode = iota // Automatic width based on content
	widthModeFixed                         // Fixed specified widths
	widthModePercent                       // Widths as percentage of available space
)

// docTable represents a table to be added to the document
type docTable struct {
	doc          *Document
	columns      []tableColumn
	rows         [][]tableCell
	width        float64
	rowHeight    float64
	cellPadding  float64
	headerStyle  CellStyle
	cellStyle    CellStyle
	alignment    position // Left, Center, or Right alignment
	currentWidth float64
}

// tableColumn represents a column in the table
type tableColumn struct {
	header      string   // Header text
	width       float64  // Width of the column
	headerAlign position // Text alignment for header (Left, Center, Right)
	align       position // Text alignment within cells (Left, Center, Right)
	prefix      string   // Prefix to add before each value in the column
	suffix      string   // Suffix to add after each value in the column
}

// tableCell represents a cell in the table
type tableCell struct {
	content      string    // Content of the cell
	useCellStyle bool      // If true, use custom style instead of the table's default
	cellStyle    CellStyle // Custom style for this cell
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
//   - "Price|HR,CR" - Right-aligned header, right-aligned column
//   - "Price|HR,CR,P:$" - Right-aligned header, right-aligned column with "$" prefix
//   - "Percentage|HC,CC,S:%" - Center-aligned header, center-aligned column with "%" suffix
//   - "Name|HL,CL,W:30%" - Left-aligned header, left-aligned column with 30% of available width
//   - "Age|HC,CR,W:20" - Center-aligned header, right-aligned column with fixed width of 20 units
func (doc *Document) NewTable(headers ...string) *docTable {
	// Almacenar los encabezados para verificación posterior
	doc.lastTableHeaders = headers

	// Crear una nueva tabla con configuración predeterminada
	table := &docTable{
		doc:         doc,
		rowHeight:   25, // Default row height
		cellPadding: 5,  // Default padding
		alignment:   Center,
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
		options := parseHeaderFormat(headerStr)
		if options.WidthMode == widthModePercent {
			return widthModePercent
		}
	}

	// Si todos los encabezados usan ancho fijo, usamos modo fijo
	allFixed := true
	for _, headerStr := range headers {
		options := parseHeaderFormat(headerStr)
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
func (t *docTable) SetHeaderAlignment(columnIndex int, alignment position) *docTable {
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
	t.alignment = Left
	return t
}

// AlignCenter centers the table horizontally (default)
func (t *docTable) AlignCenter() *docTable {
	t.alignment = Center
	return t
}

// AlignRight aligns the table to the right margin
func (t *docTable) AlignRight() *docTable {
	t.alignment = Right
	return t
}

// HeaderStyle sets the style for the header row
func (t *docTable) HeaderStyle(style CellStyle) *docTable {
	t.headerStyle = style
	return t
}

// CellStyle sets the default style for regular cells
func (t *docTable) CellStyle(style CellStyle) *docTable {
	t.cellStyle = style
	return t
}

// anyToString converts any value to a string without using fmt
// Uses TinyGo-compatible approach to convert numbers to strings
func (t *docTable) anyToString(v any) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.FormatInt(int64(val), 10)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		// For any other type, return empty string
		// In a full implementation, you might want to handle more types
		t.doc.log("docTable anyToString error Unsupported type for conversion to string:", val)
		return ""
	}
}

// AddRow adds a row of data to the table
// Accepts any value type which will be converted to strings
func (t *docTable) AddRow(cells ...any) *docTable {
	rowCells := make([]tableCell, len(cells))
	for i, content := range cells {
		formattedContent := t.anyToString(content)

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
	Style   CellStyle
}

// NewStyledCell creates a new cell with custom styling
func (doc *Document) NewStyledCell(content string, style CellStyle) styledCell {
	return styledCell{
		Content: content,
		Style:   style,
	}
}

// Draw renders the table on the document
func (t *docTable) Draw() error {
	// Calculate total height of table
	totalHeight := t.rowHeight * float64(len(t.rows)+1) // +1 for header row

	// Check if table fits on current page
	y := t.doc.ensureElementFits(totalHeight, t.doc.fontConfig.Normal.SpaceAfter)

	// Calculate starting X position - no recalculation of widths
	x := t.calculatePosition()

	// Colección para guardar información de los encabezados para dibujar sus bordes al final
	type headerInfo struct {
		x, y, width, height float64
		style               CellStyle
	}
	headers := []headerInfo{}

	// Draw header row (sin bordes por ahora, solo contenido y fondo)
	currentX := x
	for _, col := range t.columns {
		// Guardamos información del encabezado para dibujar los bordes después
		headers = append(headers, headerInfo{
			x:      currentX,
			y:      y,
			width:  col.width,
			height: t.rowHeight,
			style:  t.headerStyle,
		})

		// Dibujar contenido y fondo de la celda de encabezado usando la alineación especificada
		t.drawCellContent(
			currentX,
			y,
			col.width,
			t.rowHeight,
			col.header,
			col.headerAlign, // Use header-specific alignment
			t.headerStyle,
		)
		currentX += col.width
	}

	// Draw data rows
	for rowIndex, row := range t.rows {
		currentY := y + ((float64(rowIndex) + 1) * t.rowHeight)
		currentX = x

		// Check if this row fits on the current page
		if currentY+t.rowHeight > t.doc.curr.pageSize.H-t.doc.margins.Bottom {
			t.doc.AddPage()
			currentY = t.doc.margins.Top

			// Limpiar la lista de encabezados para la nueva página
			headers = []headerInfo{}

			// Redraw the header row on the new page (solo contenido y fondo)
			headerY := currentY
			headerX := x
			for _, col := range t.columns {
				// Guardamos información del encabezado para dibujar los bordes después
				headers = append(headers, headerInfo{
					x:      headerX,
					y:      headerY,
					width:  col.width,
					height: t.rowHeight,
					style:  t.headerStyle,
				})

				// Dibujar contenido y fondo de la celda de encabezado usando la alineación especificada
				t.drawCellContent(
					headerX,
					headerY,
					col.width,
					t.rowHeight,
					col.header,
					col.headerAlign, // Use header-specific alignment
					t.headerStyle,
				)
				headerX += col.width
			}

			// Adjust currentY to start below the header
			currentY += t.rowHeight
		}

		for colIndex, cell := range row {
			// Use column width, handle case where row has fewer cells than columns
			if colIndex < len(t.columns) {
				cellWidth := t.columns[colIndex].width
				cellAlign := t.columns[colIndex].align

				// Determine which style to use
				style := t.cellStyle
				if cell.useCellStyle {
					style = cell.cellStyle
				}

				// Dibujar celda completa (contenido, fondo y bordes)
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
	}

	// Ahora dibujamos los bordes de los encabezados al final
	for _, h := range headers {
		t.drawCellBorder(h.x, h.y, h.width, h.height, h.style.BorderStyle)
	}

	// Update document position to after the table
	t.doc.SetY(y + totalHeight + t.doc.fontConfig.Normal.SpaceAfter)

	return nil
}

// drawCellContent dibuja solo el contenido y fondo de una celda (sin bordes)
func (t *docTable) drawCellContent(
	x float64,
	y float64,
	width float64,
	height float64,
	content string,
	align position,
	style CellStyle,
) {
	// Fill the cell background if a fill color is specified
	if (style.FillColor != RGBColor{}) {
		t.doc.SetFillColor(style.FillColor.R, style.FillColor.G, style.FillColor.B)
		t.doc.RectFromUpperLeftWithStyle(x, y, width, height, "F")
	}

	// Set text properties
	if style.Font != "" {
		t.doc.SetFont(style.Font, "", style.FontSize)
	}
	t.doc.SetTextColor(style.TextColor.R, style.TextColor.G, style.TextColor.B)

	// Calculate text position with padding
	textX := x + t.cellPadding
	textY := y + t.cellPadding
	textWidth := width - (2 * t.cellPadding)
	textHeight := height - (2 * t.cellPadding)

	// Create cell options
	cellOpt := cellOption{
		Align:  align | Middle, // Combine horizontal alignment with vertical middle alignment
		Border: 0,              // No borders
		breakOption: &breakOption{
			Mode:           breakModeIndicatorSensitive,
			BreakIndicator: ' ',
		},
	}

	// Draw the cell content
	t.doc.SetXY(textX, textY)
	err := t.doc.MultiCellWithOption(&Rect{W: textWidth, H: textHeight}, content, cellOpt)
	if err != nil && err.Error() != "empty string" {
		t.doc.log("Error drawing table cell:", err)
	}
}

// drawCellBorder dibuja solo los bordes de una celda
func (t *docTable) drawCellBorder(
	x float64,
	y float64,
	width float64,
	height float64,
	borderStyle BorderStyle,
) {
	if borderStyle.Width > 0 {
		t.doc.SetLineWidth(borderStyle.Width)
		t.doc.SetStrokeColor(
			borderStyle.RGBColor.R,
			borderStyle.RGBColor.G,
			borderStyle.RGBColor.B,
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
	align position,
	style CellStyle,
) {
	// Primero dibujamos el contenido y el fondo
	t.drawCellContent(x, y, width, height, content, align, style)

	// Luego dibujamos los bordes
	t.drawCellBorder(x, y, width, height, style.BorderStyle)
}

// calculatePosition determina donde colocar la tabla
func (t *docTable) calculatePosition() float64 {
	// Posición X inicial (margen izquierdo)
	x := t.doc.margins.Left

	// Aplicar alineación
	switch t.alignment {
	case Center:
		// Centrar la tabla: margen izquierdo + (espacio disponible - ancho tabla) / 2
		x = t.doc.margins.Left + (t.doc.contentAreaWidth-t.width)/2
	case Right:
		// Alinear a la derecha: margen izquierdo + espacio disponible - ancho tabla
		x = t.doc.margins.Left + t.doc.contentAreaWidth - t.width
	case Left:
		// Alinear a la izquierda: simplemente el margen izquierdo
		x = t.doc.margins.Left
	}

	return x
}
