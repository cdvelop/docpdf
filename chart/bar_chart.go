package chart

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/cdvelop/docpdf/freetype/truetype"
	"github.com/cdvelop/docpdf/mathutils"
)

// BarChart is a chart that draws bars on a range.
type BarChart struct {
	Title      string
	TitleStyle Style

	ColorPalette ColorPalette

	Width  int
	Height int
	DPI    float64

	BarWidth int

	Background Style
	Canvas     Style

	XAxis Style
	YAxis YAxis

	BarSpacing int

	UseBaseValue bool
	BaseValue    float64

	Font        *truetype.Font
	defaultFont *truetype.Font

	Bars     []Value
	Elements []Renderable
}

// GetDPI returns the dpi for the chart.
func (bc BarChart) GetDPI() float64 {

	if bc.DPI == 0 {
		return DefaultDPI
	}
	return bc.DPI
}

// GetFont returns the text font.
func (bc BarChart) GetFont() *truetype.Font {
	if bc.Font == nil {
		return bc.defaultFont
	}
	return bc.Font
}

// GetWidth returns the chart width or the default value.
func (bc BarChart) GetWidth() int {
	if bc.Width == 0 {
		return DefaultChartWidth
	}
	return bc.Width
}

// GetHeight returns the chart height or the default value.
func (bc BarChart) GetHeight() int {
	if bc.Height == 0 {
		return DefaultChartHeight
	}
	return bc.Height
}

// GetBarSpacing returns the spacing between bars.
func (bc BarChart) GetBarSpacing() int {
	if bc.BarSpacing == 0 {
		return DefaultBarSpacing
	}
	return bc.BarSpacing
}

// GetBarWidth returns the default bar width.
func (bc BarChart) GetBarWidth() int {
	if bc.BarWidth == 0 {
		return DefaultBarWidth
	}
	return bc.BarWidth
}

// Render renders the chart with the given renderer to the given io.Writer.
func (bc BarChart) Render(rp RendererProvider, w io.Writer) error {
	if len(bc.Bars) == 0 {
		return errors.New("please provide at least one bar")
	}

	r, err := rp(bc.GetWidth(), bc.GetHeight())
	if err != nil {
		return err
	}

	if bc.Font == nil {
		defaultFont, err := GetDefaultFont()
		if err != nil {
			return err
		}
		bc.defaultFont = defaultFont
	}
	r.SetDPI(bc.GetDPI())

	bc.drawBackground(r)

	var canvasBox Box
	var yt []Tick
	var yr Range
	var yf ValueFormatter

	canvasBox = bc.getDefaultCanvasBox()
	yr = bc.getRanges()
	if yr.GetMax()-yr.GetMin() == 0 {
		return fmt.Errorf("invalid data range; cannot be zero")
	}
	yr = bc.setRangeDomains(canvasBox, yr)
	yf = bc.getValueFormatters()

	if bc.hasAxes() {
		yt = bc.getAxesTicks(r, yr, yf)
		canvasBox = bc.getAdjustedCanvasBox(r, canvasBox, yr, yt)
		yr = bc.setRangeDomains(canvasBox, yr)
	}
	bc.drawCanvas(r, canvasBox)
	bc.drawBars(r, canvasBox, yr)
	bc.drawXAxis(r, canvasBox)
	bc.drawYAxis(r, canvasBox, yr, yt)

	bc.drawTitle(r)
	for _, a := range bc.Elements {
		a(r, canvasBox, bc.styleDefaultsElements())
	}

	return r.Save(w)
}

func (bc BarChart) drawCanvas(r Renderer, canvasBox Box) {
	Draw.Box(r, canvasBox, bc.getCanvasStyle())
}

func (bc BarChart) getRanges() Range {
	var yrange Range
	if bc.YAxis.Range != nil && !bc.YAxis.Range.IsZero() {
		yrange = bc.YAxis.Range
	} else {
		yrange = &ContinuousRange{}
	}

	if !yrange.IsZero() {
		return yrange
	}

	if len(bc.YAxis.Ticks) > 0 {
		tickMin, tickMax := math.MaxFloat64, -math.MaxFloat64
		for _, t := range bc.YAxis.Ticks {
			tickMin = math.Min(tickMin, t.Value)
			tickMax = math.Max(tickMax, t.Value)
		}
		yrange.SetMin(tickMin)
		yrange.SetMax(tickMax)
		return yrange
	}

	min, max := math.MaxFloat64, -math.MaxFloat64
	for _, b := range bc.Bars {
		min = math.Min(b.Value, min)
		max = math.Max(b.Value, max)
	}

	yrange.SetMin(min)
	yrange.SetMax(max)

	return yrange
}

func (bc BarChart) drawBackground(r Renderer) {
	Draw.Box(r, Box{
		Right:  bc.GetWidth(),
		Bottom: bc.GetHeight(),
	}, bc.getBackgroundStyle())
}

func (bc BarChart) drawBars(r Renderer, canvasBox Box, yr Range) {
	xoffset := canvasBox.Left

	width, spacing, _ := bc.calculateScaledTotalWidth(canvasBox)
	bs2 := spacing >> 1

	var barBox Box
	var bxl, bxr, by int
	for index, bar := range bc.Bars {
		bxl = xoffset + bs2
		bxr = bxl + width

		by = canvasBox.Bottom - yr.Translate(bar.Value)

		if bc.UseBaseValue {
			barBox = Box{
				Top:    by,
				Left:   bxl,
				Right:  bxr,
				Bottom: canvasBox.Bottom - yr.Translate(bc.BaseValue),
			}
		} else {
			barBox = Box{
				Top:    by,
				Left:   bxl,
				Right:  bxr,
				Bottom: canvasBox.Bottom,
			}
		}

		Draw.Box(r, barBox, bar.Style.InheritFrom(bc.styleDefaultsBar(index)))

		xoffset += width + spacing
	}
}

func (bc BarChart) drawXAxis(r Renderer, canvasBox Box) {
	if !bc.XAxis.Hidden {
		axisStyle := bc.XAxis.InheritFrom(bc.styleDefaultsAxes())
		axisStyle.WriteToRenderer(r)

		width, spacing, _ := bc.calculateScaledTotalWidth(canvasBox)

		r.MoveTo(canvasBox.Left, canvasBox.Bottom)
		r.LineTo(canvasBox.Right, canvasBox.Bottom)
		r.Stroke()

		r.MoveTo(canvasBox.Left, canvasBox.Bottom)
		r.LineTo(canvasBox.Left, canvasBox.Bottom+DefaultVerticalTickHeight)
		r.Stroke()

		cursor := canvasBox.Left
		for index, bar := range bc.Bars {
			barLabelBox := Box{
				Top:    canvasBox.Bottom + DefaultXAxisMargin,
				Left:   cursor,
				Right:  cursor + width + spacing,
				Bottom: bc.GetHeight(),
			}

			if len(bar.Label) > 0 {
				Draw.TextWithin(r, bar.Label, barLabelBox, axisStyle)
			}

			axisStyle.WriteToRenderer(r)
			if index < len(bc.Bars)-1 {
				r.MoveTo(barLabelBox.Right, canvasBox.Bottom)
				r.LineTo(barLabelBox.Right, canvasBox.Bottom+DefaultVerticalTickHeight)
				r.Stroke()
			}
			cursor += width + spacing
		}
	}
}

func (bc BarChart) drawYAxis(r Renderer, canvasBox Box, yr Range, ticks []Tick) {
	if !bc.YAxis.Style.Hidden {
		bc.YAxis.Render(r, canvasBox, yr, bc.styleDefaultsAxes(), ticks)
	}
}

func (bc BarChart) drawTitle(r Renderer) {
	if len(bc.Title) > 0 && !bc.TitleStyle.Hidden {
		r.SetFont(bc.TitleStyle.GetFont(bc.GetFont()))
		r.SetFontColor(bc.TitleStyle.GetFontColor(bc.GetColorPalette().TextColor()))
		titleFontSize := bc.TitleStyle.GetFontSize(bc.getTitleFontSize())
		r.SetFontSize(titleFontSize)

		textBox := r.MeasureText(bc.Title)

		textWidth := textBox.Width()
		textHeight := textBox.Height()

		titleX := (bc.GetWidth() >> 1) - (textWidth >> 1)
		titleY := bc.TitleStyle.Padding.GetTop(DefaultTitleTop) + textHeight

		r.Text(bc.Title, titleX, titleY)
	}
}

func (bc BarChart) getCanvasStyle() Style {
	return bc.Canvas.InheritFrom(bc.styleDefaultsCanvas())
}

func (bc BarChart) styleDefaultsCanvas() Style {
	return Style{
		FillColor:   bc.GetColorPalette().CanvasColor(),
		StrokeColor: bc.GetColorPalette().CanvasStrokeColor(),
		StrokeWidth: DefaultCanvasStrokeWidth,
	}
}

func (bc BarChart) hasAxes() bool {
	return !bc.YAxis.Style.Hidden
}

func (bc BarChart) setRangeDomains(canvasBox Box, yr Range) Range {
	yr.SetDomain(canvasBox.Height())
	return yr
}

func (bc BarChart) getDefaultCanvasBox() Box {
	return bc.box()
}

func (bc BarChart) getValueFormatters() ValueFormatter {
	if bc.YAxis.ValueFormatter != nil {
		return bc.YAxis.ValueFormatter
	}
	return FloatValueFormatter
}

func (bc BarChart) getAxesTicks(r Renderer, yr Range, yf ValueFormatter) (yticks []Tick) {
	if !bc.YAxis.Style.Hidden {
		yticks = bc.YAxis.GetTicks(r, yr, bc.styleDefaultsAxes(), yf)
	}
	return
}

func (bc BarChart) calculateEffectiveBarSpacing(canvasBox Box) int {
	totalWithBaseSpacing := bc.calculateTotalBarWidth(bc.GetBarWidth(), bc.GetBarSpacing())
	if totalWithBaseSpacing > canvasBox.Width() {
		lessBarWidths := canvasBox.Width() - (len(bc.Bars) * bc.GetBarWidth())
		if lessBarWidths > 0 {
			return int(math.Ceil(float64(lessBarWidths) / float64(len(bc.Bars))))
		}
		return 0
	}
	return bc.GetBarSpacing()
}

func (bc BarChart) calculateEffectiveBarWidth(canvasBox Box, spacing int) int {
	totalWithBaseWidth := bc.calculateTotalBarWidth(bc.GetBarWidth(), spacing)
	if totalWithBaseWidth > canvasBox.Width() {
		totalLessBarSpacings := canvasBox.Width() - (len(bc.Bars) * spacing)
		if totalLessBarSpacings > 0 {
			return int(math.Ceil(float64(totalLessBarSpacings) / float64(len(bc.Bars))))
		}
		return 0
	}
	return bc.GetBarWidth()
}

func (bc BarChart) calculateTotalBarWidth(barWidth, spacing int) int {
	return len(bc.Bars) * (barWidth + spacing)
}

func (bc BarChart) calculateScaledTotalWidth(canvasBox Box) (width, spacing, total int) {
	spacing = bc.calculateEffectiveBarSpacing(canvasBox)
	width = bc.calculateEffectiveBarWidth(canvasBox, spacing)
	total = bc.calculateTotalBarWidth(width, spacing)
	return
}

func (bc BarChart) getAdjustedCanvasBox(r Renderer, canvasBox Box, yrange Range, yticks []Tick) Box {
	// This section is just for calculating xaxisHeight for later use
	var xaxisHeight int
	if !bc.XAxis.Hidden {
		xaxisHeight = DefaultVerticalTickHeight

		axisStyle := bc.XAxis.InheritFrom(bc.styleDefaultsAxes())
		axisStyle.WriteToRenderer(r)

		cursor := canvasBox.Left
		for _, bar := range bc.Bars {
			if len(bar.Label) > 0 {
				barLabelBox := Box{
					Top:    canvasBox.Bottom + DefaultXAxisMargin,
					Left:   cursor,
					Right:  cursor + bc.GetBarWidth() + bc.GetBarSpacing(),
					Bottom: bc.GetHeight(),
				}
				lines := Text.WrapFit(r, bar.Label, barLabelBox.Width(), axisStyle)
				linesBox := Text.MeasureLines(r, lines, axisStyle)

				xaxisHeight = mathutils.MinInt(linesBox.Height()+(2*DefaultXAxisMargin), xaxisHeight)
			}
		}
	}
	initialBox := bc.box()               // Get the box based on background padding first
	finalCanvasBox := initialBox.Clone() // Start adjusting from the initial box

	// --- Adjust for Y Axis (Primary) ---
	if !bc.YAxis.Style.Hidden {
		yAxisStyleDefaults := bc.styleDefaultsAxes()
		yAxisStyleDefaults.WriteToRenderer(r) // Set style for measurement

		// Calculate max label width manually
		maxLabelWidth := 0
		for _, t := range yticks {
			// Ensure tick style is applied for measurement if different
			tickStyle := bc.YAxis.TickStyle.InheritFrom(yAxisStyleDefaults)
			tickStyle.WriteToRenderer(r)
			labelBox := r.MeasureText(t.Label)
			if labelBox.Width() > maxLabelWidth {
				maxLabelWidth = labelBox.Width()
			}
		}

		// Calculate total width needed:
		yAxisTotalWidth := maxLabelWidth + DefaultYAxisMargin + DefaultHorizontalTickWidth

		// Add space for axis name if present
		if !bc.YAxis.NameStyle.Hidden && len(bc.YAxis.Name) > 0 {
			// Use the actual name style defined on the axis, inheriting defaults
			nameStyle := bc.YAxis.NameStyle.InheritFrom(yAxisStyleDefaults)
			// Assume default rotation is 90 degrees if not specified otherwise
			if nameStyle.TextRotationDegrees == 0 {
				nameStyle.TextRotationDegrees = 90
			}
			nameStyle.WriteToRenderer(r)
			nameBox := r.MeasureText(bc.YAxis.Name)

			var nameWidth int
			// Use height for width if rotated vertically
			if nameStyle.TextRotationDegrees == 90 || nameStyle.TextRotationDegrees == 270 {
				nameWidth = nameBox.Height()
			} else {
				nameWidth = nameBox.Width()
			}
			// Add name width and another margin
			yAxisTotalWidth += nameWidth + DefaultYAxisMargin
		}

		// The final canvas needs to start at the initial position + the calculated width.
		requiredCanvasLeft := initialBox.Left + yAxisTotalWidth
		finalCanvasBox.Left = requiredCanvasLeft
	}

	// --- Adjust for X Axis ---
	if !bc.XAxis.Hidden {
		// Calculate X-axis height
		calculatedXAxisHeight := DefaultVerticalTickHeight // Minimum height
		axisStyle := bc.XAxis.InheritFrom(bc.styleDefaultsAxes())
		axisStyle.WriteToRenderer(r) // Ensure renderer has the correct style for measurement

		width, spacing, _ := bc.calculateScaledTotalWidth(canvasBox) // Use canvasBox for width calculation context

		for _, bar := range bc.Bars {
			if len(bar.Label) > 0 {
				// Estimate label box width more accurately based on scaled width/spacing
				labelBoxWidth := width + spacing
				// Ensure minimum width for very narrow bars
				if labelBoxWidth < 10 {
					labelBoxWidth = 10
				}

				// Use a temporary box for measurement, respecting potential canvas adjustments
				tempLabelBox := Box{
					Left:  0, // Use 0 for width measurement context
					Right: labelBoxWidth,
				}

				lines := Text.WrapFit(r, bar.Label, tempLabelBox.Width(), axisStyle)
				linesBox := Text.MeasureLines(r, lines, axisStyle)
				labelHeight := linesBox.Height() + DefaultXAxisMargin // Add margin

				if labelHeight > calculatedXAxisHeight {
					calculatedXAxisHeight = labelHeight
				}
			}
		}
		// Ensure bottom padding from style is considered
		calculatedXAxisHeight += axisStyle.Padding.Bottom

		finalCanvasBox.Bottom -= calculatedXAxisHeight
	}

	// Ensure top padding for title is considered
	if len(bc.Title) > 0 && !bc.TitleStyle.Hidden {
		titleHeight := bc.measureTitleHeight(r)
		// Use DefaultTitleTop for bottom margin as well for consistency, or define a specific constant if needed
		finalCanvasBox.Top += titleHeight + bc.TitleStyle.Padding.GetBottom(DefaultTitleTop)
	}

	// Ensure canvas box doesn't collapse
	if finalCanvasBox.Bottom <= finalCanvasBox.Top {
		finalCanvasBox.Bottom = finalCanvasBox.Top + 1 // Minimum 1 pixel height
	}
	if finalCanvasBox.Right <= finalCanvasBox.Left {
		finalCanvasBox.Right = finalCanvasBox.Left + 1 // Minimum 1 pixel width
	}

	return finalCanvasBox
}

// measureTitleHeight calculates the height needed for the title.
func (bc BarChart) measureTitleHeight(r Renderer) int {
	if len(bc.Title) == 0 || bc.TitleStyle.Hidden {
		return 0
	}
	style := bc.styleDefaultsTitle() // Get combined style
	r.SetFont(style.GetFont(bc.GetFont()))
	r.SetFontColor(style.GetFontColor(bc.GetColorPalette().TextColor()))
	r.SetFontSize(style.GetFontSize(bc.getTitleFontSize()))

	textBox := r.MeasureText(bc.Title)
	// Use DefaultTitleTop for bottom margin as well for consistency
	return textBox.Height() + style.Padding.GetTop(DefaultTitleTop) + style.Padding.GetBottom(DefaultTitleTop)
}

// box returns the chart bounds as a box, considering background padding.
func (bc BarChart) box() Box {
	return Box{
		Top:    bc.Background.Padding.GetTop(DefaultBackgroundPadding.Top),
		Left:   bc.Background.Padding.GetLeft(DefaultBackgroundPadding.Left),
		Right:  bc.GetWidth() - bc.Background.Padding.GetRight(DefaultBackgroundPadding.Right),
		Bottom: bc.GetHeight() - bc.Background.Padding.GetBottom(DefaultBackgroundPadding.Bottom),
	}
}

func (bc BarChart) getBackgroundStyle() Style {
	return bc.Background.InheritFrom(bc.styleDefaultsBackground())
}

func (bc BarChart) styleDefaultsBackground() Style {
	return Style{
		FillColor:   bc.GetColorPalette().BackgroundColor(),
		StrokeColor: bc.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: DefaultStrokeWidth,
	}
}

func (bc BarChart) styleDefaultsBar(index int) Style {
	return Style{
		StrokeColor: bc.GetColorPalette().GetSeriesColor(index),
		StrokeWidth: 3.0,
		FillColor:   bc.GetColorPalette().GetSeriesColor(index),
	}
}

func (bc BarChart) styleDefaultsTitle() Style {
	return bc.TitleStyle.InheritFrom(Style{
		FontColor:           bc.GetColorPalette().TextColor(),
		Font:                bc.GetFont(),
		FontSize:            bc.getTitleFontSize(),
		TextHorizontalAlign: TextHorizontalAlignCenter,
		TextVerticalAlign:   TextVerticalAlignTop,
		TextWrap:            TextWrapWord,
	})
}

func (bc BarChart) getTitleFontSize() float64 {
	effectiveDimension := mathutils.MinInt(bc.GetWidth(), bc.GetHeight())
	if effectiveDimension >= 2048 {
		return 48
	} else if effectiveDimension >= 1024 {
		return 24
	} else if effectiveDimension >= 512 {
		return 18
	} else if effectiveDimension >= 256 {
		return 12
	}
	return 10
}

func (bc BarChart) styleDefaultsAxes() Style {
	return Style{
		StrokeColor:         bc.GetColorPalette().AxisStrokeColor(),
		Font:                bc.GetFont(),
		FontSize:            DefaultAxisFontSize,
		FontColor:           bc.GetColorPalette().TextColor(),
		TextHorizontalAlign: TextHorizontalAlignCenter,
		TextVerticalAlign:   TextVerticalAlignTop,
		TextWrap:            TextWrapWord,
	}
}

func (bc BarChart) styleDefaultsElements() Style {
	return Style{
		Font: bc.GetFont(),
	}
}

// GetColorPalette returns the color palette for the chart.
func (bc BarChart) GetColorPalette() ColorPalette {
	if bc.ColorPalette != nil {
		return bc.ColorPalette
	}
	return AlternateColorPalette
}
