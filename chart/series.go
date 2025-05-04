package chart

import "github.com/cdvelop/docpdf/canvas"

// Series is an alias to Renderable.
type Series interface {
	GetName() string
	GetYAxis() YAxisType
	GetStyle() Style
	Validate() error
	Render(r Renderer, canvasBox canvas.Box, xrange, yrange Range, s Style)
}
