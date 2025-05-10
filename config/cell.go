package config

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
	Border    Border    // Border style for the cell
	FillColor Color     // Background color of the cell
	TextStyle TextStyle // Style for the text in the cell
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
	if c.TextStyle != other.TextStyle {
		return false
	}
	return true
}

// funcion que retorna un nuevo estilo de celda
func NewCell() *Cell {
	return &Cell{
		Border:    Border{},
		FillColor: Color{},
		TextStyle: TextStyle{},
	}
}
