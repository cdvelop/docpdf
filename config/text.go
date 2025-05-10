package config

import (
	"github.com/cdvelop/docpdf/style"
)

const (
	FontRegular   = "regular"
	FontBold      = "bold"
	FontItalic    = "italic"
	FontUnderline = "underline"
)

// TextStyle defines the style configuration for different text sections
type TextStyle struct {
	Size        float64
	Color       style.Color
	LineSpacing float64
	Alignment   Alignment
	SpaceBefore float64
	SpaceAfter  float64
	FontStyle   FontIntStyle // Renamed from Style for clarity
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
	SetFont(family string, style string, size any) error
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
	style := ts.normal
	pdf.SetFont(FontRegular, "", style.Size)
	pdf.SetTextColor(style.Color.R, style.Color.G, style.Color.B)
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
	textStyles.SetNormal(TextStyle{
		Size:        11,
		Color:       style.Color{R: 0, G: 0, B: 0, A: 255},
		LineSpacing: 1.15,
		Alignment:   Left | Top,
		SpaceBefore: 0,
		SpaceAfter:  8,
		FontStyle:   FontStyleRegular,
	})

	// Configure Header1
	textStyles.SetHeader1(TextStyle{
		Size:        16,
		Color:       style.Color{R: 0, G: 0, B: 0, A: 255},
		LineSpacing: 1.5,
		Alignment:   Left | Top,
		SpaceBefore: 12,
		SpaceAfter:  8,
		FontStyle:   FontStyleBold,
	})

	// Configure Header2
	textStyles.SetHeader2(TextStyle{
		Size:        14,
		Color:       style.Color{R: 0, G: 0, B: 0, A: 255},
		LineSpacing: 1.3,
		Alignment:   Left | Top,
		SpaceBefore: 10,
		SpaceAfter:  6,
		FontStyle:   FontStyleBold,
	})

	// Configure Header3
	textStyles.SetHeader3(TextStyle{
		Size:        12,
		Color:       style.Color{R: 0, G: 0, B: 0, A: 255},
		LineSpacing: 1.2,
		Alignment:   Left | Top,
		SpaceBefore: 8,
		SpaceAfter:  4,
		FontStyle:   FontStyleBold,
	})

	// Configure Footnote
	textStyles.SetFootnote(TextStyle{
		Size:        8,
		Color:       style.Color{R: 128, G: 128, B: 128, A: 255},
		LineSpacing: 1.0,
		Alignment:   Left | Top,
		SpaceBefore: 2,
		SpaceAfter:  2,
		FontStyle:   FontStyleRegular,
	})

	// Configure PageHeader
	textStyles.SetPageHeader(TextStyle{
		Size:        9,
		Color:       style.Color{R: 128, G: 128, B: 128, A: 255},
		LineSpacing: 1.0,
		Alignment:   Center | Top,
		SpaceBefore: 0,
		SpaceAfter:  12,
		FontStyle:   FontStyleRegular,
	})

	// Configure PageFooter
	textStyles.SetPageFooter(TextStyle{
		Size:        9,
		Color:       style.Color{R: 128, G: 128, B: 128, A: 255},
		LineSpacing: 1.0,
		Alignment:   Right | Top,
		SpaceBefore: 2,
		SpaceAfter:  0,
		FontStyle:   FontStyleRegular,
	})

	// Configure ChartLabel
	textStyles.SetChartLabel(TextStyle{
		Size:        9,
		Color:       style.Color{R: 50, G: 50, B: 50, A: 255},
		LineSpacing: 1.0,
		Alignment:   Left | Top,
		SpaceBefore: 0,
		SpaceAfter:  0,
		FontStyle:   FontStyleRegular,
	})

	// Configure ChartAxisLabel
	textStyles.SetChartAxisLabel(TextStyle{
		Size:        8,
		Color:       style.Color{R: 80, G: 80, B: 80, A: 255},
		LineSpacing: 1.0,
		Alignment:   Left | Top,
		SpaceBefore: 0,
		SpaceAfter:  0,
		FontStyle:   FontStyleRegular,
	})

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
