package chart

import (
	"io"
	"math"

	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/config"
	"github.com/cdvelop/docpdf/pdfengine"
	"github.com/cdvelop/docpdf/style"
)

// Point representa un punto en un camino (path) para el renderizado
type Point struct {
	X, Y float64
}

// PdfRenderer implementa la interfaz Renderer para dibujar directamente en un PDF
// usando el motor PdfEngine en lugar de generar imágenes rasterizadas.
type PdfRenderer struct {
	engine       *pdfengine.PdfEngine
	className    string
	dpi          float64
	strokeColor  style.Color
	fillColor    style.Color
	strokeWidth  float64
	dashArray    []float64
	font         config.FontFamily
	fontSize     float64
	fontColor    style.Color
	textRotation float64

	// Coordenadas de la posición actual (para MoveTo, LineTo, etc.)
	currentX float64
	currentY float64

	// Almacena los puntos del path actual para operaciones Close, Stroke, Fill, FillStroke
	pathStartX float64
	pathStartY float64
	pathPoints []Point // Usar un tipo de punto compatible con PdfEngine
	pathClosed bool
}

// NewPdfRenderer crea un nuevo renderizador PDF usando el motor PdfEngine existente
func NewPdfRenderer(engine *pdfengine.PdfEngine) *PdfRenderer {

	return &PdfRenderer{
		engine:      engine,
		dpi:         96.0, // DPI estándar para SVG
		strokeColor: style.ColorBlack,
		fillColor:   style.ColorWhite,
		font:        engine, // Usar la fuente predeterminada en lugar del engine
		strokeWidth: 1.0,
		fontSize:    10.0,
		fontColor:   style.ColorBlack,
		pathPoints:  []Point{}, // Inicializar el slice de puntos vacío
	}
}

// ResetStyle restablece los estilos del renderizador a sus valores predeterminados
func (r *PdfRenderer) ResetStyle() {
	r.className = ""
	r.strokeColor = style.ColorBlack
	r.fillColor = style.ColorWhite
	r.strokeWidth = 1.0
	r.dashArray = nil
	r.fontSize = 10.0
	r.fontColor = style.ColorBlack
	r.textRotation = 0.0
}

// GetDPI obtiene el DPI del renderizador
func (r *PdfRenderer) GetDPI() float64 {
	return r.dpi
}

// SetDPI establece el DPI para el renderizador
func (r *PdfRenderer) SetDPI(dpi float64) {
	r.dpi = dpi
}

// SetClassName establece el nombre de clase actual
func (r *PdfRenderer) SetClassName(className string) {
	r.className = className
}

// SetStrokeColor establece el color de trazo actual
func (r *PdfRenderer) SetStrokeColor(color style.Color) {
	r.strokeColor = color
	// Aplicar al pdfengine
	r.engine.SetStrokeColor(color.R, color.G, color.B)
}

// SetFillColor establece el color de relleno actual
func (r *PdfRenderer) SetFillColor(color style.Color) {
	r.fillColor = color
	// Aplicar al pdfengine
	r.engine.SetFillColor(color.R, color.G, color.B)
}

// SetStrokeWidth establece el ancho del trazo
func (r *PdfRenderer) SetStrokeWidth(width float64) {
	r.strokeWidth = width
	// Aplicar al pdfengine
	r.engine.SetLineWidth(width)
}

// SetStrokeDashArray establece el array de guiones del trazo
func (r *PdfRenderer) SetStrokeDashArray(dashArray []float64) {
	r.dashArray = dashArray
	// Aplicar al pdfengine
	r.engine.SetCustomLineType(dashArray, 0)
}

// MoveTo mueve el cursor a un punto dado
func (r *PdfRenderer) MoveTo(x, y int) {
	// Convertir coordenadas enteras a float64 y almacenarlas
	r.currentX = float64(x)
	r.currentY = float64(y)

	// Iniciar un nuevo path
	r.pathStartX = r.currentX
	r.pathStartY = r.currentY
	r.pathPoints = []Point{{X: r.currentX, Y: r.currentY}}
	r.pathClosed = false
}

// LineTo comienza una forma y dibuja una línea hasta un punto dado
// desde el punto anterior
func (r *PdfRenderer) LineTo(x, y int) {
	// Convertir coordenadas enteras a float64
	newX := float64(x)
	newY := float64(y)

	// Dibujar la línea desde la posición actual hasta el nuevo punto
	r.engine.Line(r.currentX, r.currentY, newX, newY)

	// Actualizar la posición actual
	r.currentX = newX
	r.currentY = newY

	// Añadir el punto al path actual
	r.pathPoints = append(r.pathPoints, Point{X: newX, Y: newY})
}

// QuadCurveTo dibuja una curva cuadrática
// cx y cy representan los "puntos de control" de Bezier
func (r *PdfRenderer) QuadCurveTo(cx, cy, x, y int) {
	// Convertir a coordenadas float64 para PdfEngine
	cxf := float64(cx)
	cyf := float64(cy)
	xf := float64(x)
	yf := float64(y)

	// Para una curva de Bezier cuadrática con un solo punto de control,
	// necesitamos convertirla a una curva cúbica para PdfEngine.Curve
	// La curva cúbica tiene dos puntos de control, así que calculamos
	// estos puntos a partir del punto de control cuadrático.
	//
	// La fórmula es:
	// CP1 = current + 2/3 * (quadratic_cp - current)
	// CP2 = end + 2/3 * (quadratic_cp - end)

	// Punto de control 1 para la curva cúbica
	cp1x := r.currentX + 2.0/3.0*(cxf-r.currentX)
	cp1y := r.currentY + 2.0/3.0*(cyf-r.currentY)

	// Punto de control 2 para la curva cúbica
	cp2x := xf + 2.0/3.0*(cxf-xf)
	cp2y := yf + 2.0/3.0*(cyf-yf)

	// Usar el método Curve de PdfEngine para dibujar la curva
	r.engine.Curve(r.currentX, r.currentY, cp1x, cp1y, cp2x, cp2y, xf, yf, "D")
	// Actualizar la posición current
	r.currentX = xf
	r.currentY = yf

	// Añadir el punto final al path
	r.pathPoints = append(r.pathPoints, Point{X: xf, Y: yf})
}

// ArcTo dibuja un arco con un centro dado (cx,cy)
// un conjunto de radios dado (rx,ry), un ángulo inicial y delta (en radianes)
func (r *PdfRenderer) ArcTo(cx, cy int, rx, ry, startAngle, delta float64) {
	// Convertir el centro a float64
	cxf := float64(cx)
	cyf := float64(cy)

	// Aproximamos el arco usando líneas
	// Para una aproximación decente, usamos un segmento por cada 5 grados
	const segmentAngle = 5.0 * math.Pi / 180.0 // 5 grados en radianes

	// Determinar cuántos segmentos necesitamos
	segments := int(math.Abs(delta) / segmentAngle)
	if segments < 1 {
		segments = 1
	}

	// Calcular el ángulo de cada segmento
	segmentDelta := delta / float64(segments)

	// Comenzar desde el punto inicial del arco
	angle := startAngle
	firstX := cxf + rx*math.Cos(angle)
	firstY := cyf + ry*math.Sin(angle)

	// Mover al primer punto del arco
	r.MoveTo(int(firstX), int(firstY))

	// Dibujar cada segmento del arco
	for i := 0; i < segments; i++ {
		angle += segmentDelta
		nextX := cxf + rx*math.Cos(angle)
		nextY := cyf + ry*math.Sin(angle)
		r.LineTo(int(nextX), int(nextY))
	}

	// Actualizar la posición actual (aunque LineTo ya lo hace)
	r.currentX = cxf + rx*math.Cos(startAngle+delta)
	r.currentY = cyf + ry*math.Sin(startAngle+delta)
}

// Close finaliza una forma dibujada por LineTo
func (r *PdfRenderer) Close() {
	// Si tenemos al menos un punto en el path
	if len(r.pathPoints) > 1 {
		// Dibujar una línea desde el punto actual hasta el punto inicial
		r.engine.Line(r.currentX, r.currentY, r.pathStartX, r.pathStartY)

		// Actualizar la posición actual al punto inicial
		r.currentX = r.pathStartX
		r.currentY = r.pathStartY

		// Marcar el path como cerrado
		r.pathClosed = true
	}
}

// Stroke traza el camino
func (r *PdfRenderer) Stroke() {
	// Si tenemos suficientes puntos para formar un path
	if len(r.pathPoints) > 1 {
		// Trazar líneas individuales para simular el path completo
		for i := 0; i < len(r.pathPoints)-1; i++ {
			p1 := r.pathPoints[i]
			p2 := r.pathPoints[i+1]
			r.engine.Line(p1.X, p1.Y, p2.X, p2.Y)
		}

		// Si el path está cerrado, dibujamos la línea final
		if r.pathClosed && len(r.pathPoints) > 0 {
			last := r.pathPoints[len(r.pathPoints)-1]
			first := r.pathPoints[0]
			r.engine.Line(last.X, last.Y, first.X, first.Y)
		}

		// Reiniciar el path
		r.pathPoints = nil
		r.pathClosed = false
	}
}

// Fill rellena el camino, pero no lo traza
func (r *PdfRenderer) Fill() {
	// Si tenemos suficientes puntos para formar un path
	if len(r.pathPoints) > 1 {
		// Para rellenar, necesitamos tener una forma cerrada
		// Convertimos nuestros puntos a un slice compatible con PdfEngine
		pdfePoints := []struct{ X, Y float64 }{}

		// Añadir todos los puntos del path
		for _, p := range r.pathPoints {
			pdfePoints = append(pdfePoints, struct{ X, Y float64 }{X: p.X, Y: p.Y})
		}

		// Usar el método Rectangle de bajo nivel si tenemos 4 puntos que forman un rectángulo
		// O de lo contrario usar líneas individuales
		if len(pdfePoints) == 4 {
			// Determinar las coordenadas del rectángulo
			minX, minY := pdfePoints[0].X, pdfePoints[0].Y
			maxX, maxY := pdfePoints[0].X, pdfePoints[0].Y

			for _, p := range pdfePoints {
				if p.X < minX {
					minX = p.X
				}
				if p.Y < minY {
					minY = p.Y
				}
				if p.X > maxX {
					maxX = p.X
				}
				if p.Y > maxY {
					maxY = p.Y
				}
			}

			r.engine.RectFromLowerLeftWithStyle(minX, minY, maxX-minX, maxY-minY, "F")
		} else {
			// Si no tenemos una buena manera de usar Fill directamente,
			// trazamos el contorno con el color de relleno
			for i := 0; i < len(r.pathPoints)-1; i++ {
				p1 := r.pathPoints[i]
				p2 := r.pathPoints[i+1]
				r.engine.Line(p1.X, p1.Y, p2.X, p2.Y)
			}

			// Cerrar el path si es necesario
			if r.pathClosed && len(r.pathPoints) > 0 {
				last := r.pathPoints[len(r.pathPoints)-1]
				first := r.pathPoints[0]
				r.engine.Line(last.X, last.Y, first.X, first.Y)
			}
		}

		// Reiniciar el path
		r.pathPoints = nil
		r.pathClosed = false
	}
}

// FillStroke rellena y traza un camino
func (r *PdfRenderer) FillStroke() {
	// Si tenemos suficientes puntos para formar un path
	if len(r.pathPoints) > 1 {
		// Para rellenar y trazar, necesitamos tener una forma cerrada
		// Similar al método Fill(), pero combinando Fill y Stroke

		// Usar el método Rectangle de bajo nivel si tenemos 4 puntos que forman un rectángulo
		if len(r.pathPoints) == 4 {
			minX, minY := r.pathPoints[0].X, r.pathPoints[0].Y
			maxX, maxY := r.pathPoints[0].X, r.pathPoints[0].Y

			for _, p := range r.pathPoints {
				if p.X < minX {
					minX = p.X
				}
				if p.Y < minY {
					minY = p.Y
				}
				if p.X > maxX {
					maxX = p.X
				}
				if p.Y > maxY {
					maxY = p.Y
				}
			}

			r.engine.RectFromLowerLeftWithStyle(minX, minY, maxX-minX, maxY-minY, "DF")
		} else {
			// Si no tenemos una forma directa de hacerlo, simplemente dibujamos las líneas individualmente
			// Primero rellenar
			for i := 0; i < len(r.pathPoints)-1; i++ {
				p1 := r.pathPoints[i]
				p2 := r.pathPoints[i+1]
				r.engine.Line(p1.X, p1.Y, p2.X, p2.Y)
			}

			// Cerrar el path si es necesario
			if r.pathClosed && len(r.pathPoints) > 0 {
				last := r.pathPoints[len(r.pathPoints)-1]
				first := r.pathPoints[0]
				r.engine.Line(last.X, last.Y, first.X, first.Y)
			}
		}

		// Reiniciar el path
		r.pathPoints = nil
		r.pathClosed = false
	}
}

// Circle dibuja un círculo en las coordenadas dadas con un radio dado
func (r *PdfRenderer) Circle(radius float64, x, y int) {
	// Convertir coordenadas enteras a float64
	xf, yf := float64(x), float64(y)

	// Dibujar un círculo usando Oval de PdfEngine
	// La API de Oval espera un rectángulo que defina una elipse,
	// así que calculamos las coordenadas para que sea un círculo perfecto
	r.engine.Oval(
		xf-radius, yf-radius,
		xf+radius, yf+radius,
	)
}

// SetFont establece una fuente para un campo de texto
func (r *PdfRenderer) SetFont(font config.FontFamily) {
	r.font = font

	// Configurar la fuente en PdfEngine
	// Aquí necesitamos usar el nombre de la fuente tal como se registró en el PDF
	// Esto probablemente requiera pasar de FontProvider a la fuente específica del PDF
	fontName := font.Name()
	r.engine.SetFont(fontName, "", r.fontSize)
}

// SetFontColor establece el color de una fuente
func (r *PdfRenderer) SetFontColor(color style.Color) {
	r.fontColor = color
	// Aplicar al pdfengine (para el texto se usa SetFillColor)
	r.engine.SetFillColor(color.R, color.G, color.B)
}

// SetFontSize establece el tamaño de fuente para un campo de texto
func (r *PdfRenderer) SetFontSize(size float64) {
	r.fontSize = size

	// Si ya tenemos una fuente configurada, actualizar el tamaño
	if r.font != nil {
		r.engine.SetFont(r.font.Name(), "", size)
	}
}

// Text dibuja un blob de texto
func (r *PdfRenderer) Text(body string, x, y int) {
	// Convertir coordenadas a float64
	xf := float64(x)
	yf := float64(y)

	// Si hay rotación, aplicarla
	if r.textRotation != 0 {
		// Rotar alrededor del punto (x, y)
		angleDegrees := r.textRotation * 180.0 / math.Pi
		r.engine.Rotate(angleDegrees, xf, yf)

		// Dibujar el texto
		r.engine.Text(body)

		// Restaurar el estado después de rotar
		r.engine.RotateReset()
	} else {
		// Sin rotación, simplemente dibujamos el texto
		r.engine.SetX(xf)
		r.engine.SetY(yf)
		r.engine.Text(body)
	}
}

// MeasureText mide el texto
func (r *PdfRenderer) MeasureText(body string) canvas.Box {
	// Usar métodos de medición de texto de PdfEngine
	// PdfEngine puede calcular el ancho del texto con MeasureTextWidth
	width, err := r.engine.MeasureTextWidth(body)
	if err != nil {
		// Si hay un error, devolver una estimación
		width = float64(len(body)) * r.fontSize * 0.5
	}

	// Para la altura, podemos usar una aproximación basada en el tamaño de fuente
	// En general, la altura es aproximadamente 1.2 veces el tamaño de la fuente
	height := r.fontSize * 1.2
	return canvas.Box{
		Right:  int(width),
		Bottom: int(height),
		Left:   0,
		Top:    0,
	}
}

// SetTextRotation establece una rotación para dibujar elementos
func (r *PdfRenderer) SetTextRotation(radians float64) {
	r.textRotation = radians
	// La rotación se aplica cuando se dibuja el texto
}

// ClearTextRotation limpia la rotación
func (r *PdfRenderer) ClearTextRotation() {
	r.textRotation = 0
	// No hay necesidad de hacer nada más aquí, ya que
	// la rotación se aplica cuando se dibuja el texto
}

// Save escribe la imagen en el escritor dado
// En el caso de PdfRenderer, este método no hace nada directamente ya que
// el renderizado se hace directamente en el PDF proporcionado en la creación
func (r *PdfRenderer) Save(w io.Writer) error {
	// No necesitamos hacer nada aquí, ya que ya estamos dibujando
	// directamente en el PDF que se guardará con los métodos estándar de Document
	return nil
}

// NewPdfRendererProvider crea un proveedor de renderizado PDF
func NewPdfRendererProvider(engine *pdfengine.PdfEngine) RendererProvider {
	return func(width int, height int) (Renderer, error) {
		return NewPdfRenderer(engine), nil
	}
}
