package pdfengine

import (
	"fmt"
	"io"
	"math"
)

const colorTypeStrokeRGB = "RG"

const colorTypeFillRGB = "rg"

type cacheContentColorRGB struct {
	colorType string
	r, g, b   uint8
}

func (c *cacheContentColorRGB) Write(w Writer, protection *pdfProtection) error {
	fmt.Fprintf(w, "%.3f %.3f %.3f %s\n", float64(c.r)/255, float64(c.g)/255, float64(c.b)/255, c.colorType)
	return nil
}

const colorTypeStrokeCMYK = "K"

const colorTypeFillCMYK = "k"

type cacheContentColorCMYK struct {
	colorType  string
	c, m, y, k uint8
}

func (c *cacheContentColorCMYK) Write(w Writer, protection *pdfProtection) error {
	fmt.Fprintf(w, "%.2f %.2f %.2f %.2f %s\n", float64(c.c)/100, float64(c.m)/100, float64(c.y)/100, float64(c.k)/100, c.colorType)
	return nil
}

type cacheContentCustomLineType struct {
	dashArray []float64
	dashPhase float64
}

func (c *cacheContentCustomLineType) Write(w Writer, protection *pdfProtection) error {
	fmt.Fprintf(w, "%0.2f %0.2f d\n", c.dashArray, c.dashPhase)
	return nil
}

const grayTypeFill = "g"
const grayTypeStroke = "G"

type cacheContentGray struct {
	grayType string
	scale    float64
}

func (c *cacheContentGray) Write(w Writer, protection *pdfProtection) error {
	fmt.Fprintf(w, "%.2f %s\n", c.scale, c.grayType)
	return nil
}

type cacheContentLineType struct {
	lineType string
}

func (c *cacheContentLineType) Write(w Writer, protection *pdfProtection) error {
	switch c.lineType {
	case "dashed":
		fmt.Fprint(w, "[5] 2 d\n")
	case "dotted":
		fmt.Fprint(w, "[2 3] 11 d\n")
	default:
		fmt.Fprint(w, "[] 0 d\n")
	}
	return nil
}

type cacheContentPolygon struct {
	pageHeight float64
	style      string
	points     []point
	opts       polygonOptions
}

func (c *cacheContentPolygon) Write(w Writer, protection *pdfProtection) error {

	fmt.Fprintf(w, "q\n")
	for _, extGStateIndex := range c.opts.extGStateIndexes {
		fmt.Fprintf(w, "/GS%d gs\n", extGStateIndex)
	}

	for i, point := range c.points {
		fmt.Fprintf(w, "%.2f %.2f", point.X, c.pageHeight-point.Y)
		if i == 0 {
			fmt.Fprintf(w, " m ")
		} else {
			fmt.Fprintf(w, " l ")
		}

	}

	if c.style == "F" {
		fmt.Fprintf(w, " f\n")
	} else if c.style == "FD" || c.style == "DF" {
		fmt.Fprintf(w, " b\n")
	} else {
		fmt.Fprintf(w, " s\n")
	}

	fmt.Fprintf(w, "Q\n")
	return nil
}

type cacheContentTextColorRGB struct {
	r, g, b uint8
}

func (c cacheContentTextColorRGB) Write(w Writer, protection *pdfProtection) error {
	fmt.Fprintf(w, "%.3f %.3f %.3f %s\n", float64(c.r)/255, float64(c.g)/255, float64(c.b)/255, colorTypeFillRGB)
	return nil
}

func (c cacheContentTextColorRGB) Equal(obj ICacheColorText) bool {
	rgb, ok := obj.(cacheContentTextColorRGB)
	if !ok {
		return false
	}

	return c.r == rgb.r && c.g == rgb.g && c.b == rgb.b
}

type cacheContentTextColorCMYK struct {
	c, m, y, k uint8
}

func (c cacheContentTextColorCMYK) Write(w Writer, protection *pdfProtection) error {
	fmt.Fprintf(w, "%.2f %.2f %.2f %.2f %s\n", float64(c.c)/100, float64(c.m)/100, float64(c.y)/100, float64(c.k)/100, colorTypeFillCMYK)
	return nil
}

func (c cacheContentTextColorCMYK) Equal(obj ICacheColorText) bool {
	cmyk, ok := obj.(cacheContentTextColorCMYK)
	if !ok {
		return false
	}

	return c.c == cmyk.c && c.m == cmyk.m && c.y == cmyk.y && c.k == cmyk.k
}

type cacheContentRotate struct {
	isReset     bool
	pageHeight  float64
	angle, x, y float64
}

func (cc *cacheContentRotate) Write(w Writer, protection *pdfProtection) error {
	if cc.isReset == true {
		if _, err := io.WriteString(w, "Q\n"); err != nil {
			return err
		}

		return nil
	}

	matrix := computeRotateTransformationMatrix(cc.x, cc.y, cc.angle, cc.pageHeight)
	contentStream := fmt.Sprintf("q\n %s", matrix)

	if _, err := io.WriteString(w, contentStream); err != nil {
		return err
	}

	return nil
}

func computeRotateTransformationMatrix(x, y, degreeAngle, pageHeight float64) string {
	radianAngle := degreeAngle * (math.Pi / 180)

	c := math.Cos(radianAngle)
	s := math.Sin(radianAngle)
	cy := pageHeight - y

	return fmt.Sprintf("%.5f %.5f %.5f\n %.5f %.2f %.2f cm\n 1 0 0\n 1 %.2f %.2f cm\n", c, s, -s, c, x, cy, -x, -cy)
}

type cacheContentLine struct {
	pageHeight float64
	x1         float64
	y1         float64
	x2         float64
	y2         float64
	opts       lineOptions
}

func (c *cacheContentLine) Write(w Writer, protection *pdfProtection) error {
	fmt.Fprintf(w, "q\n")
	for _, extGStateIndex := range c.opts.extGStateIndexes {
		fmt.Fprintf(w, "/GS%d gs\n", extGStateIndex)
	}
	fmt.Fprintf(w, "%0.2f %0.2f m %0.2f %0.2f l S\n", c.x1, c.pageHeight-c.y1, c.x2, c.pageHeight-c.y2)
	fmt.Fprintf(w, "Q\n")
	return nil
}

type cacheContentImportedTemplate struct {
	pageHeight float64
	tplName    string
	scaleX     float64
	scaleY     float64
	tX         float64
	tY         float64
}

func (c *cacheContentImportedTemplate) Write(w Writer, protection *pdfProtection) error {
	c.tY += c.pageHeight
	fmt.Fprintf(w, "q 0 J 1 w 0 j 0 G 0 g q %.4F 0 0 %.4F %.4F %.4F cm %s Do Q Q\n", c.scaleX, c.scaleY, c.tX, c.tY, c.tplName)
	return nil
}

type cacheContentLineWidth struct {
	width float64
}

func (c *cacheContentLineWidth) Write(w Writer, protection *pdfProtection) error {
	fmt.Fprintf(w, "%.2f w\n", c.width)
	return nil
}

type cacheContentOval struct {
	pageHeight float64
	x1         float64
	y1         float64
	x2         float64
	y2         float64
}

func (c *cacheContentOval) Write(w Writer, protection *pdfProtection) error {

	h := c.pageHeight
	x1 := c.x1
	y1 := c.y1
	x2 := c.x2
	y2 := c.y2

	cp := 0.55228                              // Magnification of the control point
	v1 := [2]float64{x1 + (x2-x1)/2, h - y2}   // Vertex of the lower
	v2 := [2]float64{x2, h - (y1 + (y2-y1)/2)} // .. config.Right
	v3 := [2]float64{x1 + (x2-x1)/2, h - y1}   // .. Upper
	v4 := [2]float64{x1, h - (y1 + (y2-y1)/2)} // .. config.Left

	fmt.Fprintf(w, "%0.2f %0.2f m\n", v1[0], v1[1])
	fmt.Fprintf(w,
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c\n",
		v1[0]+(x2-x1)/2*cp, v1[1], v2[0], v2[1]-(y2-y1)/2*cp, v2[0], v2[1],
	)
	fmt.Fprintf(w,
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c\n",
		v2[0], v2[1]+(y2-y1)/2*cp, v3[0]+(x2-x1)/2*cp, v3[1], v3[0], v3[1],
	)
	fmt.Fprintf(w,
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c\n",
		v3[0]-(x2-x1)/2*cp, v3[1], v4[0], v4[1]+(y2-y1)/2*cp, v4[0], v4[1],
	)
	fmt.Fprintf(w,
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c S\n",
		v4[0], v4[1]-(y2-y1)/2*cp, v1[0]-(x2-x1)/2*cp, v1[1], v1[0], v1[1],
	)

	return nil
}

type cacheContentRectangle struct {
	pageHeight       float64
	x                float64
	y                float64
	width            float64
	height           float64
	style            paintStyle
	extGStateIndexes []int
}

func newCacheContentRectangle(pageHeight float64, rectOpts drawableRectOptions) ICacheContent {
	if rectOpts.paintStyle == "" {
		rectOpts.paintStyle = drawPaintStyle
	}

	return cacheContentRectangle{
		x:                rectOpts.X,
		y:                rectOpts.Y,
		width:            rectOpts.W,
		height:           rectOpts.H,
		pageHeight:       pageHeight,
		style:            rectOpts.paintStyle,
		extGStateIndexes: rectOpts.extGStateIndexes,
	}
}

func (c cacheContentRectangle) Write(w Writer, protection *pdfProtection) error {
	stream := "q\n"

	for _, extGStateIndex := range c.extGStateIndexes {
		stream += fmt.Sprintf("/GS%d gs\n", extGStateIndex)
	}

	stream += fmt.Sprintf("%0.2f %0.2f %0.2f %0.2f re %s\n", c.x, c.pageHeight-c.y, c.width, c.height, c.style)

	stream += "Q\n"

	if _, err := io.WriteString(w, stream); err != nil {
		return err
	}

	return nil
}

type cacheContentCurve struct {
	pageHeight float64
	x0         float64
	y0         float64
	x1         float64
	y1         float64
	x2         float64
	y2         float64
	x3         float64
	y3         float64
	style      string
}

func (c *cacheContentCurve) Write(w Writer, protection *pdfProtection) error {

	h := c.pageHeight
	x0 := c.x0
	y0 := c.y0
	x1 := c.x1
	y1 := c.y1
	x2 := c.x2
	y2 := c.y2
	x3 := c.x3
	y3 := c.y3
	style := c.style

	//cp := 0.55228
	fmt.Fprintf(w, "%0.2f %0.2f m\n", x0, h-y0)
	fmt.Fprintf(w,
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c",
		x1, h-y1, x2, h-y2, x3, h-y3,
	)
	op := "S"
	if style == "F" {
		op = "f"
	} else if style == "FD" || style == "DF" {
		op = "B"
	}
	fmt.Fprintf(w, " %s\n", op)
	return nil
}
