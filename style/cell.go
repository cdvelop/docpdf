package style

// Defines the border style for a cell or table
type Border struct {
	Top    bool    // Whether to draw the top border
	Left   bool    // Whether to draw the left border
	Right  bool    // Whether to draw the right border
	Bottom bool    // Whether to draw the bottom border
	Width  float64 // Width of the border line
	Color  Color   // Color of the border
}

// Defines the style for a cell, including border, fill, text, and font properties
type Cell struct {
	Border    Border  // Border style for the cell
	FillColor Color   // Background color of the cell
	TextColor Color   // Color of the text in the cell
	Font      string  // Font name for the cell text
	FontSize  float64 // Font size for the cell text
}

// metodo para comparar dos estilos de celda
func (c *Cell) Equals(other *Cell) bool {
	if c == nil && other == nil {
		return true
	}
	if c == nil || other == nil {
		return false
	}
	if c.Border != other.Border {
		return false
	}
	if c.FillColor != other.FillColor {
		return false
	}
	if c.TextColor != other.TextColor {
		return false
	}
	if c.Font != other.Font {
		return false
	}
	if c.FontSize != other.FontSize {
		return false
	}
	return true
}

// funcion que retorna un nuevo estilo de celda
func NewCell() *Cell {
	return &Cell{
		Border:    Border{},
		FillColor: Color{},
		TextColor: Color{},
		Font:      "",
		FontSize:  0,
	}
}
