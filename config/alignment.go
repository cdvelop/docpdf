package config

// position representa una posición o alineación en el documento
type Alignment int

const (
	// Left representa alineación a la izquierda
	Left Alignment = 8 //001000
	// Top representa alineación superior
	Top Alignment = 4 //000100
	// Right representa alineación a la derecha
	Right Alignment = 2 //000010
	// Bottom representa alineación inferior
	Bottom Alignment = 1 //000001
	// Center representa alineación al centro
	Center Alignment = 16 //010000
	// Middle representa alineación al medio
	Middle Alignment = 32 //100000
	// Justify representa texto justificado
	Justify Alignment = 64 //1000000
	// AllBorders representa todos los bordes
	AllBorders Alignment = 15 //001111
)
