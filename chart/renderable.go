package chart

import "github.com/cdvelop/docpdf/canvas"

// Renderable is a function that can be called to render custom elements on the chart.
type Renderable func(r Renderer, canvasBox canvas.Box, defaults Style)
