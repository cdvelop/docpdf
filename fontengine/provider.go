package fontengine

// FontProvider es una interfaz que abstrae las propiedades necesarias de una fuente
// para que el renderizador pueda trabajar con ella independientemente de su implementación
type FontProvider interface {
	// Identificación de la fuente
	Name() string   // Nombre de la fuente
	Family() string // Familia de la fuente

	// Propiedades de estilo
	Weight() string // Peso: regular, bold, etc.
	Style() string  // Estilo: normal, italic, etc.

	// Propiedades para renderizado SVG
	SVGFontID() string // ID para referenciar en SVG

	// Opcionalmente, para sistemas que necesiten la ruta al archivo
	Path() string // Ruta al archivo de la fuente
}
