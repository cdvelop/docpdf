package pdfengine

import (
	"math"

	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/errs"
)

type lineOptions struct {
	extGStateIndexes []int
}

type polygonOptions struct {
	extGStateIndexes []int
}

type drawableRectOptions struct {
	canvas.Rect
	X            float64
	Y            float64
	paintStyle   paintStyle
	transparency *transparency

	extGStateIndexes []int
}

// SetLineWidth : set line width
func (gp *PdfEngine) SetLineWidth(width float64) {
	gp.curr.lineWidth = gp.UnitsToPoints(width)
	gp.getContent().AppendStreamSetLineWidth(gp.UnitsToPoints(width))
}

// SetLineType : set line type  ("dashed" ,"dotted")
//
//	Usage:
//	pdf.SetLineType("dashed")
//	pdf.Line(50, 200, 550, 200)
//	pdf.SetLineType("dotted")
//	pdf.Line(50, 400, 550, 400)
func (gp *PdfEngine) SetLineType(linetype string) {
	gp.getContent().AppendStreamSetLineType(linetype)
}

// SetCustomLineType : set custom line type
//
//	Usage:
//	pdf.SetCustomLineType([]float64{0.8, 0.8}, 0)
//	pdf.Line(50, 200, 550, 200)
func (gp *PdfEngine) SetCustomLineType(dashArray []float64, dashPhase float64) {
	for i := range dashArray {
		gp.UnitsToPointsVar(&dashArray[i])
	}
	gp.UnitsToPointsVar(&dashPhase)
	gp.getContent().AppendStreamSetCustomLineType(dashArray, dashPhase)
}

// Line : draw line
//
//	Usage:
//	pdf.SetTransparency(docpdf.transparency{Alpha: 0.5,blendModeType: docpdf.colorBurn})
//	pdf.SetLineType("dotted")
//	pdf.SetStrokeColor(255, 0, 0)
//	pdf.SetLineWidth(2)
//	pdf.Line(10, 30, 585, 30)
//	pdf.ClearTransparency()
func (gp *PdfEngine) Line(x1 float64, y1 float64, x2 float64, y2 float64) {
	gp.UnitsToPointsVar(&x1, &y1, &x2, &y2)
	transparency, err := gp.getCachedTransparency(nil)
	if err != nil {
		transparency = nil
	}
	var opts = lineOptions{}
	if transparency != nil {
		opts.extGStateIndexes = append(opts.extGStateIndexes, transparency.extGStateIndex)
	}
	gp.getContent().AppendStreamLine(x1, y1, x2, y2, opts)
}

// RectFromLowerLeft : draw rectangle from lower-left corner (x, y)
func (gp *PdfEngine) RectFromLowerLeft(x float64, y float64, wdth float64, hght float64) {
	gp.UnitsToPointsVar(&x, &y, &wdth, &hght)

	opts := drawableRectOptions{
		X:          x,
		Y:          y,
		paintStyle: drawPaintStyle,
		Rect:       canvas.Rect{W: wdth, H: hght},
	}

	gp.getContent().AppendStreamRectangle(opts)
}

// RectFromUpperLeft : draw rectangle from upper-left corner (x, y)
func (gp *PdfEngine) RectFromUpperLeft(x float64, y float64, wdth float64, hght float64) {
	gp.UnitsToPointsVar(&x, &y, &wdth, &hght)

	opts := drawableRectOptions{
		X:          x,
		Y:          y + hght,
		paintStyle: drawPaintStyle,
		Rect:       canvas.Rect{W: wdth, H: hght},
	}

	gp.getContent().AppendStreamRectangle(opts)
}

// RectFromLowerLeftWithStyle : draw rectangle from lower-left corner (x, y)
//   - style: Style of rectangule (draw and/or fill: D, F, DF, FD)
//     D or empty string: draw. This is the default value.
//     F: fill
//     DF or FD: draw and fill
func (gp *PdfEngine) RectFromLowerLeftWithStyle(x float64, y float64, wdth float64, hght float64, style string) {
	opts := drawableRectOptions{
		X: x,
		Y: y,
		Rect: canvas.Rect{
			H: hght,
			W: wdth,
		},
		paintStyle: parseStyle(style),
	}
	gp.RectFromLowerLeftWithOpts(opts)
}

func (gp *PdfEngine) RectFromLowerLeftWithOpts(opts drawableRectOptions) error {
	gp.UnitsToPointsVar(&opts.X, &opts.Y, &opts.W, &opts.H)

	imageTransparency, err := gp.getCachedTransparency(opts.transparency)
	if err != nil {
		return err
	}

	if imageTransparency != nil {
		opts.extGStateIndexes = append(opts.extGStateIndexes, imageTransparency.extGStateIndex)
	}

	gp.getContent().AppendStreamRectangle(opts)

	return nil
}

// RectFromUpperLeftWithStyle : draw rectangle from upper-left corner (x, y)
//   - style: Style of rectangule (draw and/or fill: D, F, DF, FD)
//     D or empty string: draw. This is the default value.
//     F: fill
//     DF or FD: draw and fill
func (gp *PdfEngine) RectFromUpperLeftWithStyle(x float64, y float64, wdth float64, hght float64, style string) {
	opts := drawableRectOptions{
		X: x,
		Y: y,
		Rect: canvas.Rect{
			H: hght,
			W: wdth,
		},
		paintStyle: parseStyle(style),
	}
	gp.RectFromUpperLeftWithOpts(opts)
}

func (gp *PdfEngine) RectFromUpperLeftWithOpts(opts drawableRectOptions) error {
	gp.UnitsToPointsVar(&opts.X, &opts.Y, &opts.W, &opts.H)

	opts.Y += opts.H

	imageTransparency, err := gp.getCachedTransparency(opts.transparency)
	if err != nil {
		return err
	}

	if imageTransparency != nil {
		opts.extGStateIndexes = append(opts.extGStateIndexes, imageTransparency.extGStateIndex)
	}

	gp.getContent().AppendStreamRectangle(opts)

	return nil
}

// Oval : draw oval
func (gp *PdfEngine) Oval(x1 float64, y1 float64, x2 float64, y2 float64) {
	gp.UnitsToPointsVar(&x1, &y1, &x2, &y2)
	gp.getContent().AppendStreamOval(x1, y1, x2, y2)
}

// Curve Draws a Bézier curve (the Bézier curve is tangent to the line between the control points at either end of the curve)
// Parameters:
// - x0, y0: Start point
// - x1, y1: Control point 1
// - x2, y2: Control point 2
// - x3, y3: End point
// - style: Style of rectangule (draw and/or fill: D, F, DF, FD)
func (gp *PdfEngine) Curve(x0 float64, y0 float64, x1 float64, y1 float64, x2 float64, y2 float64, x3 float64, y3 float64, style string) {
	gp.UnitsToPointsVar(&x0, &y0, &x1, &y1, &x2, &y2, &x3, &y3)
	gp.getContent().AppendStreamCurve(x0, y0, x1, y1, x2, y2, x3, y3, style)
}

// Polygon : draw polygon
//   - style: Style of polygon (draw and/or fill: D, F, DF, FD)
//     D or empty string: draw. This is the default value.
//     F: fill
//     DF or FD: draw and fill
//
// Usage:
//
//	 pdf.SetStrokeColor(255, 0, 0)
//		pdf.SetLineWidth(2)
//		pdf.SetFillColor(0, 255, 0)
//		pdf.Polygon([]docpdf.point{{X: 10, Y: 30}, {X: 585, Y: 200}, {X: 585, Y: 250}}, "DF")
func (gp *PdfEngine) Polygon(points []point, style string) {

	transparency, err := gp.getCachedTransparency(nil)
	if err != nil {
		transparency = nil
	}

	var opts = polygonOptions{}
	if transparency != nil {
		opts.extGStateIndexes = append(opts.extGStateIndexes, transparency.extGStateIndex)
	}

	var pointReals []point
	for _, p := range points {
		x := p.X
		y := p.Y
		gp.UnitsToPointsVar(&x, &y)
		pointReals = append(pointReals, point{X: x, Y: y})
	}
	gp.getContent().AppendStreamPolygon(pointReals, style, opts)
}

// Rectangle : draw rectangle, and add radius input to make a round corner, it helps to calculate the round corner coordinates and use Polygon functions to draw rectangle
//   - style: Style of Rectangle (draw and/or fill: D, F, DF, FD)
//     D or empty string: draw. This is the default value.
//     F: fill
//     DF or FD: draw and fill
//
// Usage:
//
//	 pdf.SetStrokeColor(255, 0, 0)
//		pdf.SetLineWidth(2)
//		pdf.SetFillColor(0, 255, 0)
//		pdf.Rectangle(196.6, 336.8, 398.3, 379.3, "DF", 3, 10)
func (gp *PdfEngine) Rectangle(x0 float64, y0 float64, x1 float64, y1 float64, style string, radius float64, radiusPointNum int) error {
	if x1 <= x0 || y1 <= y0 {
		return errs.InvalidRectangleCoordinates
	}
	if radiusPointNum <= 0 || radius <= 0 {
		//draw rectangle without round corner
		points := []point{}
		points = append(points, point{X: x0, Y: y0})
		points = append(points, point{X: x1, Y: y0})
		points = append(points, point{X: x1, Y: y1})
		points = append(points, point{X: x0, Y: y1})
		gp.Polygon(points, style)

	} else {

		if radius > (x1-x0) || radius > (y1-y0) {
			return errs.InvalidRectangleCoordinates
		}

		degrees := []float64{}
		angle := float64(90) / float64(radiusPointNum+1)
		accAngle := angle
		for accAngle < float64(90) {
			degrees = append(degrees, accAngle)
			accAngle += angle
		}

		radians := []float64{}
		for _, v := range degrees {
			radians = append(radians, v*math.Pi/180)
		}

		points := []point{}
		points = append(points, point{X: x0, Y: (y0 + radius)})
		for _, v := range radians {
			offsetX := radius * math.Cos(v)
			offsetY := radius * math.Sin(v)
			x := x0 + radius - offsetX
			y := y0 + radius - offsetY
			points = append(points, point{X: x, Y: y})
		}
		points = append(points, point{X: (x0 + radius), Y: y0})

		points = append(points, point{X: (x1 - radius), Y: y0})
		for i := range radians {
			v := radians[len(radians)-1-i]
			offsetX := radius * math.Cos(v)
			offsetY := radius * math.Sin(v)
			x := x1 - radius + offsetX
			y := y0 + radius - offsetY
			points = append(points, point{X: x, Y: y})
		}
		points = append(points, point{X: x1, Y: (y0 + radius)})

		points = append(points, point{X: x1, Y: (y1 - radius)})
		for _, v := range radians {
			offsetX := radius * math.Cos(v)
			offsetY := radius * math.Sin(v)
			x := x1 - radius + offsetX
			y := y1 - radius + offsetY
			points = append(points, point{X: x, Y: y})
		}
		points = append(points, point{X: (x1 - radius), Y: y1})

		points = append(points, point{X: (x0 + radius), Y: y1})
		for i := range radians {
			v := radians[len(radians)-1-i]
			offsetX := radius * math.Cos(v)
			offsetY := radius * math.Sin(v)
			x := x0 + radius - offsetX
			y := y1 - radius + offsetY
			points = append(points, point{X: x, Y: y})
		}
		points = append(points, point{X: x0, Y: y1 - radius})

		gp.Polygon(points, style)
	}
	return nil
}
