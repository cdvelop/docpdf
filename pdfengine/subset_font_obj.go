package pdfengine

import (
	"fmt"
	"io"

	"github.com/cdvelop/docpdf/errs"
	"github.com/cdvelop/docpdf/fontmaker/core"
)

// errCharNotFound char not found
var errCharNotFound = errs.New("char not found")

// errGlyphNotFound font file not contain glyph
var errGlyphNotFound = errs.New("glyph not found")

// subsetFontObj pdf subsetFont object
type subsetFontObj struct {
	ttfp                  core.TTFParser
	Family                string
	CharacterToGlyphIndex *mapOfCharacterToGlyphIndex
	CountOfFont           int
	indexObjCIDFont       int
	indexObjUnicodeMap    int
	ttfFontOption         TtfOption
	funcKernOverride      funcKernOverride
	funcGetRoot           func() *PdfEngine
	addCharsBuff          []rune
}

func (s *subsetFontObj) init(funcGetRoot func() *PdfEngine) {
	s.CharacterToGlyphIndex = newMapOfCharacterToGlyphIndex() //make(map[rune]uint)
	s.funcKernOverride = nil
	s.funcGetRoot = funcGetRoot

}

func (s *subsetFontObj) write(w io.Writer, objID int) error {
	//me.AddChars("จ")
	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "/BaseFont /%s\n", createEmbeddedFontSubsetName(s.Family))
	fmt.Fprintf(w, "/DescendantFonts [%d 0 R]\n", s.indexObjCIDFont+1)
	io.WriteString(w, "/Encoding /Identity-H\n")
	io.WriteString(w, "/Subtype /Type0\n")
	fmt.Fprintf(w, "/ToUnicode %d 0 R\n", s.indexObjUnicodeMap+1)
	io.WriteString(w, "/Type /Font\n")
	io.WriteString(w, ">>\n")
	return nil
}

// SetIndexObjCIDFont set IndexObjCIDFont
func (s *subsetFontObj) SetIndexObjCIDFont(index int) {
	s.indexObjCIDFont = index
}

// SetIndexObjUnicodeMap set IndexObjUnicodeMap
func (s *subsetFontObj) SetIndexObjUnicodeMap(index int) {
	s.indexObjUnicodeMap = index
}

// SetFamily set font family name
func (s *subsetFontObj) SetFamily(familyname string) {
	s.Family = familyname
}

// GetFamily get font family name
func (s *subsetFontObj) GetFamily() string {
	return s.Family
}

// SetTtfFontOption set TtfOption must set before SetTTFByPath
func (s *subsetFontObj) SetTtfFontOption(option TtfOption) {
	if option.OnGlyphNotFoundSubstitute == nil {
		option.OnGlyphNotFoundSubstitute = defaultOnGlyphNotFoundSubstitute
	}
	s.ttfFontOption = option
}

// GetTtfFontOption get TtfOption must set before SetTTFByPath
func (s *subsetFontObj) GetTtfFontOption() TtfOption {
	return s.ttfFontOption
}

// KernValueByLeft find kern value from kern table by left
func (s *subsetFontObj) KernValueByLeft(left uint) (bool, *core.KernValue) {

	if !s.ttfFontOption.UseKerning {
		return false, nil
	}

	k := s.ttfp.Kern()
	if k == nil {
		return false, nil
	}

	if kval, ok := k.Kerning[left]; ok {
		return true, &kval
	}

	return false, nil
}

// SetTTFByPath set ttf
func (s *subsetFontObj) SetTTFByPath(ttfpath string) error {
	useKerning := s.ttfFontOption.UseKerning
	s.ttfp.SetUseKerning(useKerning)
	err := s.ttfp.Parse(ttfpath)
	if err != nil {
		return err
	}
	return nil
}

// SetTTFByReader set ttf
func (s *subsetFontObj) SetTTFByReader(rd io.Reader) error {
	useKerning := s.ttfFontOption.UseKerning
	s.ttfp.SetUseKerning(useKerning)
	err := s.ttfp.ParseByReader(rd)
	if err != nil {
		return err
	}
	return nil
}

// SetTTFData set ttf
func (s *subsetFontObj) SetTTFData(data []byte) error {
	useKerning := s.ttfFontOption.UseKerning
	s.ttfp.SetUseKerning(useKerning)
	err := s.ttfp.ParseFontData(data)
	if err != nil {
		return err
	}
	return nil
}

// AddChars add char to map CharacterToGlyphIndex
func (s *subsetFontObj) AddChars(txt string) (string, error) {
	s.addCharsBuff = s.addCharsBuff[:0]
	for _, runeValue := range txt {
		if s.CharacterToGlyphIndex.KeyExists(runeValue) {
			s.addCharsBuff = append(s.addCharsBuff, runeValue)
			continue
		}
		glyphIndex, err := s.CharCodeToGlyphIndex(runeValue)
		if err == errGlyphNotFound {
			//never return error on this, just call function OnGlyphNotFound
			if s.ttfFontOption.OnGlyphNotFound != nil {
				s.ttfFontOption.OnGlyphNotFound(runeValue)
			}
			//start: try to find rune for replace
			alreadyExists, runeValueReplace, glyphIndexReplace := s.replaceGlyphThatNotFound(runeValue)
			if !alreadyExists {
				s.CharacterToGlyphIndex.Set(runeValueReplace, glyphIndexReplace) // [runeValue] = glyphIndex
			}
			//end: try to find rune for replace
			s.addCharsBuff = append(s.addCharsBuff, runeValueReplace)
			continue
		} else if err != nil {
			return "", err
		}
		s.CharacterToGlyphIndex.Set(runeValue, glyphIndex) // [runeValue] = glyphIndex
		s.addCharsBuff = append(s.addCharsBuff, runeValue)
	}
	return string(s.addCharsBuff), nil
}

/*
//AddChars add char to map CharacterToGlyphIndex
func (s *subsetFontObj) AddChars(txt string) error {

	for _, runeValue := range txt {
		if s.CharacterToGlyphIndex.KeyExists(runeValue) {
			continue
		}
		glyphIndex, err := s.CharCodeToGlyphIndex(runeValue)
		if err == errGlyphNotFound {
			//never return error on this, just call function OnGlyphNotFound
			if s.ttfFontOption.OnGlyphNotFound != nil {
				s.ttfFontOption.OnGlyphNotFound(runeValue)
			}
			//start: try to find rune for replace
			runeValueReplace, glyphIndexReplace, ok := s.replaceGlyphThatNotFound(runeValue)
			if ok {
				s.CharacterToGlyphIndex.Set(runeValueReplace, glyphIndexReplace) // [runeValue] = glyphIndex
			}
			//end: try to find rune for replace
			continue
		} else if err != nil {
			return err
		}
		s.CharacterToGlyphIndex.Set(runeValue, glyphIndex) // [runeValue] = glyphIndex
	}
	return nil
}
*/

// replaceGlyphThatNotFound find glyph to replaced
// it returns
// - true if rune already add to CharacterToGlyphIndex
// - rune for replace
// - rune for replace is found or not
// - glyph index for replace
func (s *subsetFontObj) replaceGlyphThatNotFound(runeNotFound rune) (bool, rune, uint) {
	if s.ttfFontOption.OnGlyphNotFoundSubstitute != nil {
		runeForReplace := s.ttfFontOption.OnGlyphNotFoundSubstitute(runeNotFound)
		if s.CharacterToGlyphIndex.KeyExists(runeForReplace) {
			return true, runeForReplace, 0
		}
		glyphIndexForReplace, err := s.CharCodeToGlyphIndex(runeForReplace)
		if err != nil {
			return false, runeForReplace, 0
		}
		return false, runeForReplace, glyphIndexForReplace
	}
	return false, runeNotFound, 0
}

// CharIndex index of char in glyph table
func (s *subsetFontObj) CharIndex(r rune) (uint, error) {
	glyIndex, ok := s.CharacterToGlyphIndex.Val(r)
	if ok {
		return glyIndex, nil
	}
	return 0, errCharNotFound
}

// CharWidth with of char
func (s *subsetFontObj) CharWidth(r rune) (uint, error) {
	glyIndex, ok := s.CharacterToGlyphIndex.Val(r)
	if ok {
		return s.GlyphIndexToPdfWidth(glyIndex), nil
	}
	return 0, errCharNotFound
}

func (s *subsetFontObj) getType() string {
	return "SubsetFont"
}

func (s *subsetFontObj) charCodeToGlyphIndexFormat12(r rune) (uint, error) {

	value := uint(r)
	gTbs := s.ttfp.GroupingTables()
	for _, gTb := range gTbs {
		if value >= gTb.StartCharCode && value <= gTb.EndCharCode {
			gIndex := (value - gTb.StartCharCode) + gTb.GlyphID
			return gIndex, nil
		}
	}

	return uint(0), errGlyphNotFound
}

func (s *subsetFontObj) charCodeToGlyphIndexFormat4(r rune) (uint, error) {
	value := uint(r)
	seg := uint(0)
	segCount := s.ttfp.SegCount
	for seg < segCount {
		if value <= s.ttfp.EndCount[seg] {
			break
		}
		seg++
	}
	//fmt.Printf("\ncccc--->%#v\n", me.ttfp.Chars())
	if value < s.ttfp.StartCount[seg] {
		return 0, errGlyphNotFound
	}

	if s.ttfp.IdRangeOffset[seg] == 0 {

		return (value + s.ttfp.IdDelta[seg]) & 0xFFFF, nil
	}
	//fmt.Printf("IdRangeOffset=%d\n", me.ttfp.IdRangeOffset[seg])
	idx := s.ttfp.IdRangeOffset[seg]/2 + (value - s.ttfp.StartCount[seg]) - (segCount - seg)

	if s.ttfp.GlyphIdArray[int(idx)] == uint(0) {
		return 0, nil
	}

	return (s.ttfp.GlyphIdArray[int(idx)] + s.ttfp.IdDelta[seg]) & 0xFFFF, nil
}

// CharCodeToGlyphIndex gets glyph index from char code.
func (s *subsetFontObj) CharCodeToGlyphIndex(r rune) (uint, error) {
	value := uint64(r)
	if value <= 0xFFFF {
		gIndex, err := s.charCodeToGlyphIndexFormat4(r)
		if err != nil {
			return 0, err
		}
		return gIndex, nil
	}
	gIndex, err := s.charCodeToGlyphIndexFormat12(r)
	if err != nil {
		return 0, err
	}
	return gIndex, nil
}

// GlyphIndexToPdfWidth gets width from glyphIndex.
func (s *subsetFontObj) GlyphIndexToPdfWidth(glyphIndex uint) uint {

	numberOfHMetrics := s.ttfp.NumberOfHMetrics()
	unitsPerEm := s.ttfp.UnitsPerEm()
	if glyphIndex >= numberOfHMetrics {
		glyphIndex = numberOfHMetrics - 1
	}

	width := s.ttfp.Widths()[glyphIndex]
	if unitsPerEm == 1000 {
		return width
	}
	return width * 1000 / unitsPerEm
}

// GetTTFParser gets TTFParser.
func (s *subsetFontObj) GetTTFParser() *core.TTFParser {
	return &s.ttfp
}

// GetUnderlineThickness underlineThickness.
func (s *subsetFontObj) GetUnderlineThickness() int {
	return s.ttfp.UnderlineThickness()
}

func (s *subsetFontObj) GetUnderlineThicknessPx(fontSize float64) float64 {
	return (float64(s.ttfp.UnderlineThickness()) / float64(s.ttfp.UnitsPerEm())) * fontSize
}

// GetUnderlinePosition underline alignment.Alignment.
func (s *subsetFontObj) GetUnderlinePosition() int {
	return s.ttfp.UnderlinePosition()
}

func (s *subsetFontObj) GetUnderlinePositionPx(fontSize float64) float64 {
	return (float64(s.ttfp.UnderlinePosition()) / float64(s.ttfp.UnitsPerEm())) * fontSize
}

func (s *subsetFontObj) GetAscender() int {
	return s.ttfp.Ascender()
}

func (s *subsetFontObj) GetAscenderPx(fontSize float64) float64 {
	return (float64(s.ttfp.Ascender()) / float64(s.ttfp.UnitsPerEm())) * fontSize
}

func (s *subsetFontObj) GetDescender() int {
	return s.ttfp.Descender()
}

func (s *subsetFontObj) GetDescenderPx(fontSize float64) float64 {
	return (float64(s.ttfp.Descender()) / float64(s.ttfp.UnitsPerEm())) * fontSize
}
