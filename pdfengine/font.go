package pdfengine

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/cdvelop/docpdf/errs"
)

// Regular - font style regular
const Regular = 0 //000000
// Italic - font style italic
const Italic = 1 //000001
// Bold - font style bold
const Bold = 2 //000010
// Underline - font style underline
const Underline = 4 //000100

func getConvertedStyle(fontStyle string) (style int) {
	fontStyle = strings.ToUpper(fontStyle)
	if strings.Contains(fontStyle, "B") {
		style = style | Bold
	}
	if strings.Contains(fontStyle, "I") {
		style = style | Italic
	}
	if strings.Contains(fontStyle, "U") {
		style = style | Underline
	}
	return
}

// iFont represents a font interface.
type iFont interface {
	Init()
	GetType() string
	GetName() string
	GetDesc() []fontDescItem
	GetUp() int
	GetUt() int
	GetCw() fontCw
	GetEnc() string
	GetDiff() string
	GetOriginalsize() int

	SetFamily(family string)
	GetFamily() string
}

// TtfOption  font option
type TtfOption struct {
	UseKerning                bool
	Style                     int               //Regular|Bold|Italic
	OnGlyphNotFound           func(r rune)      //Called when a glyph cannot be found, just for debugging
	OnGlyphNotFoundSubstitute func(r rune) rune //Called when a glyph cannot be found, we can return a new rune to replace it.
}

func defaultTtfFontOption() TtfOption {
	var defa TtfOption
	defa.UseKerning = false
	defa.Style = Regular
	defa.OnGlyphNotFoundSubstitute = defaultOnGlyphNotFoundSubstitute
	return defa
}

func defaultOnGlyphNotFoundSubstitute(r rune) rune {
	return rune('\u0020')
}

// fontCw maps characters to integers.
type fontCw map[byte]int

// fontDescItem is a (key, value) pair.
type fontDescItem struct {
	Key string
	Val string
}

// fontObj font obj
type fontObj struct {
	Family string
	//Style string
	//Size int
	IsEmbedFont bool

	indexObjWidth          int
	indexObjFontDescriptor int
	indexObjEncoding       int

	Font        iFont
	CountOfFont int
}

func (f *fontObj) Init(funcGetRoot func() *PdfEngine) {
	f.IsEmbedFont = false
	//me.CountOfFont = -1
}

func (f *fontObj) Write(w Writer, objID int) error {
	baseFont := f.Family
	if f.Font != nil {
		baseFont = f.Font.GetName()
	}

	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "  /Type /%s\n", f.GetType())
	io.WriteString(w, "  /Subtype /TrueType\n")
	fmt.Fprintf(w, "  /BaseFont /%s\n", baseFont)
	if f.IsEmbedFont {
		io.WriteString(w, "  /FirstChar 32 /LastChar 255\n")
		fmt.Fprintf(w, "  /Widths %d 0 R\n", f.indexObjWidth)
		fmt.Fprintf(w, "  /FontDescriptor %d 0 R\n", f.indexObjFontDescriptor)
		fmt.Fprintf(w, "  /Encoding %d 0 R\n", f.indexObjEncoding)
	}
	io.WriteString(w, ">>\n")
	return nil
}

func (f *fontObj) GetType() string {
	return "Font"
}

// SetIndexObjWidth sets the width of a font object.
func (f *fontObj) SetIndexObjWidth(index int) {
	f.indexObjWidth = index
}

// SetIndexObjFontDescriptor sets the font descriptor.
func (f *fontObj) SetIndexObjFontDescriptor(index int) {
	f.indexObjFontDescriptor = index
}

// SetIndexObjEncoding sets the encoding.
func (f *fontObj) SetIndexObjEncoding(index int) {
	f.indexObjEncoding = index
}

// SetFontWithStyle : set font style support Regular or Underline
// for Bold|Italic should be loaded appropriate fonts with same styles defined
// size MUST be uint*, int* or float64*
func (gp *PdfEngine) SetFontWithStyle(family string, style int, size any) error {
	fontSize, err := convertNumericToFloat64(size)
	if err != nil {
		return err
	}
	found := false
	i := 0
	max := len(gp.pdfObjs)
	for i < max {
		if gp.pdfObjs[i].GetType() == subsetFont {
			obj := gp.pdfObjs[i]
			sub, ok := obj.(*ttfSubsetObj)
			if ok {
				if sub.GetFamily() == family && sub.GetTtfFontOption().Style == style&^Underline {
					gp.curr.FontSize = fontSize
					gp.curr.FontStyle = style
					gp.curr.FontFontCount = sub.CountOfFont
					gp.curr.FontISubset = sub
					found = true
					break
				}
			}
		}
		i++
	}

	if !found {
		return errs.MissingFontFamily
	}

	return nil
}

// SetFont : set font style support "" or "U"
// for "B" and "I" should be loaded appropriate fonts with same styles defined
// size MUST be uint*, int* or float64*
func (gp *PdfEngine) SetFont(family string, style string, size any) error {
	return gp.SetFontWithStyle(family, getConvertedStyle(style), size)
}

// SetFontSize : set the font size (and only the font size) of the currently
// active font
func (gp *PdfEngine) SetFontSize(fontSize float64) error {
	gp.curr.FontSize = fontSize
	return nil
}

// SetCharSpacing : set the character spacing of the currently active font
func (gp *PdfEngine) SetCharSpacing(charSpacing float64) error {
	gp.UnitsToPointsVar(&charSpacing)
	gp.curr.CharSpacing = charSpacing
	return nil
}

// fontConvertHelperCw2Str converts main ASCII characters of a FontCW to a string.
func fontConvertHelperCw2Str(cw fontCw) string {
	buff := new(bytes.Buffer)
	buff.WriteString(" ")
	i := 32
	for i <= 255 {
		buff.WriteString(strconv.Itoa(cw[byte(i)]) + " ")
		i++
	}
	return buff.String()
}

// fontConvertHelper_Cw2Str converts main ASCII characters of a FontCW to a string. (for backward compatibility)
// Deprecated: Use fontConvertHelperCw2Str(cw fontCw) instead
func fontConvertHelper_Cw2Str(cw fontCw) string {
	return fontConvertHelperCw2Str(cw)
}

// fontDescriptorObj is a font descriptor object.
type fontDescriptorObj struct {
	font              iFont
	fontFileObjRelate string
}

func (f *fontDescriptorObj) Init(funcGetRoot func() *PdfEngine) {

}

func (f *fontDescriptorObj) Write(w Writer, objID int) error {

	fmt.Fprintf(w, "<</Type /FontDescriptor /FontName /%s ", f.font.GetName())
	descs := f.font.GetDesc()
	i := 0
	max := len(descs)
	for i < max {
		fmt.Fprintf(w, "/%s %s ", descs[i].Key, descs[i].Val)
		i++
	}

	if f.GetType() == "Type1" {
		io.WriteString(w, "/FontFile ")
	} else {
		io.WriteString(w, "/FontFile2 ")
	}

	io.WriteString(w, f.fontFileObjRelate)
	io.WriteString(w, ">>\n")

	return nil
}

func (f *fontDescriptorObj) GetType() string {
	return "FontDescriptor"
}

// SetFont sets the font in descriptor.
func (f *fontDescriptorObj) SetFont(font iFont) {
	f.font = font
}

// GetFont gets font from descriptor.
func (f *fontDescriptorObj) GetFont() iFont {
	return f.font
}

// SetFontFileObjRelate ???
func (f *fontDescriptorObj) SetFontFileObjRelate(relate string) {
	f.fontFileObjRelate = relate
}
