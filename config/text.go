package config

import (
	"github.com/cdvelop/docpdf/style"
)

// TextStyle defines the style configuration for different text sections
type TextStyle struct {
	fontStyle   FontStyle // Contains font name, style, size, and color
	lineSpacing float64
	alignment   Alignment
	spaceBefore float64
	spaceAfter  float64
}

// NewTextStyle creates a new TextStyle with all properties
func NewTextStyle(fontStyle FontStyle, lineSpacing float64, alignment Alignment, spaceBefore, spaceAfter float64) TextStyle {
	return TextStyle{
		fontStyle:   fontStyle,
		lineSpacing: lineSpacing,
		alignment:   alignment,
		spaceBefore: spaceBefore,
		spaceAfter:  spaceAfter,
	}
}

// GetFontStyle returns the FontStyle
func (ts TextStyle) GetFontStyle() FontStyle {
	return ts.fontStyle
}

// SetFontStyle sets the FontStyle
func (ts *TextStyle) SetFontStyle(fontStyle FontStyle) {
	ts.fontStyle = fontStyle
}

// GetLineSpacing returns the line spacing
func (ts TextStyle) GetLineSpacing() float64 {
	return ts.lineSpacing
}

// SetLineSpacing sets the line spacing
func (ts *TextStyle) SetLineSpacing(lineSpacing float64) {
	ts.lineSpacing = lineSpacing
}

// GetAlignment returns the text alignment
func (ts TextStyle) GetAlignment() Alignment {
	return ts.alignment
}

// SetAlignment sets the text alignment
func (ts *TextStyle) SetAlignment(alignment Alignment) {
	ts.alignment = alignment
}

// GetSpaceBefore returns the space before the text
func (ts TextStyle) GetSpaceBefore() float64 {
	return ts.spaceBefore
}

// SetSpaceBefore sets the space before the text
func (ts *TextStyle) SetSpaceBefore(spaceBefore float64) {
	ts.spaceBefore = spaceBefore
}

// GetSpaceAfter returns the space after the text
func (ts TextStyle) GetSpaceAfter() float64 {
	return ts.spaceAfter
}

// SetSpaceAfter sets the space after the text
func (ts *TextStyle) SetSpaceAfter(spaceAfter float64) {
	ts.spaceAfter = spaceAfter
}

// TextStyles represents different font configurations for document sections.
// This was previously named FontConfig in docFont.go.
type TextStyles struct {
	fontFamily     FontFamily
	normal         TextStyle
	header1        TextStyle
	header2        TextStyle
	header3        TextStyle
	footnote       TextStyle
	pageHeader     TextStyle
	pageFooter     TextStyle
	chartLabel     TextStyle // For chart bar labels/values
	chartAxisLabel TextStyle // For chart axis labels (X/Y axes)
}

type pdfEngine interface {
	AddFontFamilyConfig(fontFamily FontFamily) error
	SetFont(FontStyle) error
	SetTextColor(r uint8, g uint8, b uint8)
	SetStrokeColor(r uint8, g uint8, b uint8)
	SetLineWidth(width float64)
}

// LoadFonts prepares and loads the font family using the provided engine.
// This method centralizes font loading logic within TextStyles.
func (ts *TextStyles) LoadFonts(pdf pdfEngine) error {
	fontFamily := ts.GetFontFamily()
	if fontFamily.Path == "" {
		fontFamily.Path = "fonts/"
	}
	// Set default values if not provided for Bold and Italic
	if fontFamily.Bold == "" {
		fontFamily.Bold = fontFamily.Regular
	}
	if fontFamily.Italic == "" {
		fontFamily.Italic = fontFamily.Regular
	}
	ts.SetFontFamily(fontFamily)

	return pdf.AddFontFamilyConfig(fontFamily)
}

// SetDefaultTextConfig applies the normal text style
func (ts *TextStyles) SetDefaultTextConfig(pdf pdfEngine) {
	fontStyle := ts.normal.GetFontStyle()
	pdf.SetFont(fontStyle)
	color := fontStyle.GetColor()
	pdf.SetTextColor(color.R, color.G, color.B)
	pdf.SetLineWidth(0.7)
	pdf.SetStrokeColor(0, 0, 0)
}

// DefaultTextStyles returns word-processor like defaults
func DefaultTextConfig() TextStyles {
	textStyles := TextStyles{}

	// Configure FontFamily
	textStyles.SetFontFamily(FontFamily{
		Regular:   "regular.ttf",
		Bold:      "bold.ttf",
		Italic:    "italic.ttf",
		Underline: "", // No default underline font, underlining is often a text decoration
		Path:      "fonts/",
	})
	// Configure Normal
	normalFont := FontStyleRegular.WithSize(11).WithColor(style.Color{R: 0, G: 0, B: 0, A: 255})
	textStyles.SetNormal(NewTextStyle(
		normalFont,
		1.15,
		Left|Top,
		0,
		8,
	))
	// Configure Header1
	header1Font := FontStyleBold.WithSize(16).WithColor(style.Color{R: 0, G: 0, B: 0, A: 255})
	textStyles.SetHeader1(NewTextStyle(
		header1Font,
		1.5,
		Left|Top,
		12,
		8,
	))

	// Configure Header2
	header2Font := FontStyleBold.WithSize(14).WithColor(style.Color{R: 0, G: 0, B: 0, A: 255})
	textStyles.SetHeader2(NewTextStyle(
		header2Font,
		1.3,
		Left|Top,
		10,
		6,
	))

	// Configure Header3
	header3Font := FontStyleBold.WithSize(12).WithColor(style.Color{R: 0, G: 0, B: 0, A: 255})
	textStyles.SetHeader3(NewTextStyle(
		header3Font,
		1.2,
		Left|Top,
		8,
		4,
	))

	// Configure Footnote
	footnoteFont := FontStyleRegular.WithSize(8).WithColor(style.Color{R: 128, G: 128, B: 128, A: 255})
	textStyles.SetFootnote(NewTextStyle(
		footnoteFont,
		1.0,
		Left|Top,
		2,
		2,
	))
	// Configure PageHeader
	pageHeaderFont := FontStyleRegular.WithSize(9).WithColor(style.Color{R: 128, G: 128, B: 128, A: 255})
	textStyles.SetPageHeader(NewTextStyle(
		pageHeaderFont,
		1.0,
		Center|Top,
		0,
		12,
	))

	// Configure PageFooter
	pageFooterFont := FontStyleRegular.WithSize(9).WithColor(style.Color{R: 128, G: 128, B: 128, A: 255})
	textStyles.SetPageFooter(NewTextStyle(
		pageFooterFont,
		1.0,
		Right|Top,
		2,
		0,
	))

	// Configure ChartLabel
	chartLabelFont := FontStyleRegular.WithSize(9).WithColor(style.Color{R: 50, G: 50, B: 50, A: 255})
	textStyles.SetChartLabel(NewTextStyle(
		chartLabelFont,
		1.0,
		Left|Top,
		0,
		0,
	))

	// Configure ChartAxisLabel
	chartAxisLabelFont := FontStyleRegular.WithSize(8).WithColor(style.Color{R: 80, G: 80, B: 80, A: 255})
	textStyles.SetChartAxisLabel(NewTextStyle(
		chartAxisLabelFont,
		1.0,
		Left|Top,
		0,
		0,
	))

	return textStyles
}

// GetFontFamily returns the font family configuration
func (ts *TextStyles) GetFontFamily() FontFamily {
	return ts.fontFamily
}

// SetFontFamily sets the font family configuration
func (ts *TextStyles) SetFontFamily(fontFamily FontFamily) {
	ts.fontFamily = fontFamily
}

// GetNormal returns the normal text style
func (ts *TextStyles) GetNormal() TextStyle {
	return ts.normal
}

// SetNormal sets the normal text style
func (ts *TextStyles) SetNormal(style TextStyle) {
	ts.normal = style
}

// GetHeader1 returns the header1 text style
func (ts *TextStyles) GetHeader1() TextStyle {
	return ts.header1
}

// SetHeader1 sets the header1 text style
func (ts *TextStyles) SetHeader1(style TextStyle) {
	ts.header1 = style
}

// GetHeader2 returns the header2 text style
func (ts *TextStyles) GetHeader2() TextStyle {
	return ts.header2
}

// SetHeader2 sets the header2 text style
func (ts *TextStyles) SetHeader2(style TextStyle) {
	ts.header2 = style
}

// GetHeader3 returns the header3 text style
func (ts *TextStyles) GetHeader3() TextStyle {
	return ts.header3
}

// SetHeader3 sets the header3 text style
func (ts *TextStyles) SetHeader3(style TextStyle) {
	ts.header3 = style
}

// GetFootnote returns the footnote text style
func (ts *TextStyles) GetFootnote() TextStyle {
	return ts.footnote
}

// SetFootnote sets the footnote text style
func (ts *TextStyles) SetFootnote(style TextStyle) {
	ts.footnote = style
}

// GetPageHeader returns the page header text style
func (ts *TextStyles) GetPageHeader() TextStyle {
	return ts.pageHeader
}

// SetPageHeader sets the page header text style
func (ts *TextStyles) SetPageHeader(style TextStyle) {
	ts.pageHeader = style
}

// GetPageFooter returns the page footer text style
func (ts *TextStyles) GetPageFooter() TextStyle {
	return ts.pageFooter
}

// SetPageFooter sets the page footer text style
func (ts *TextStyles) SetPageFooter(style TextStyle) {
	ts.pageFooter = style
}

// GetChartLabel returns the chart label text style
func (ts *TextStyles) GetChartLabel() TextStyle {
	return ts.chartLabel
}

// SetChartLabel sets the chart label text style
func (ts *TextStyles) SetChartLabel(style TextStyle) {
	ts.chartLabel = style
}

// GetChartAxisLabel returns the chart axis label text style
func (ts *TextStyles) GetChartAxisLabel() TextStyle {
	return ts.chartAxisLabel
}

// SetChartAxisLabel sets the chart axis label text style
func (ts *TextStyles) SetChartAxisLabel(style TextStyle) {
	ts.chartAxisLabel = style
}
