package pdfengine

import (
	"bytes"

	"github.com/cdvelop/docpdf/env"
)

// relateFonts is a slice of relateFont.
type relateFonts []relateFont

// relateFont is a metadata index for fonts?
type relateFont struct {
	Family string
	//etc /F1
	CountOfFont int
	//etc  5 0 R
	IndexOfObj int
	Style      int // Regular|Bold|Italic
}

// IsContainsFamily checks if font family exists.
func (re *relateFonts) IsContainsFamily(family string) bool {
	for _, rf := range *re {
		if rf.Family == family {
			return true
		}
	}
	return false
}

// IsContainsFamilyAndStyle checks if font with same name and style already exists .
func (re *relateFonts) IsContainsFamilyAndStyle(family string, style int) bool {
	for _, rf := range *re {
		if rf.Family == family && rf.Style == style {
			return true
		}
	}
	return false
}

// AddTTFFont : add font file
func (gp *PdfEngine) AddTTFFont(family string, ttfpath string) error {
	return gp.AddTTFFontWithOption(family, ttfpath, defaultTtfFontOption())
}

// AddTTFFontByReader adds font file by reader.
func (gp *PdfEngine) AddTTFFontByReader(family string, rd Reader) error {
	return gp.AddTTFFontByReaderWithOption(family, rd, defaultTtfFontOption())
}

// AddTTFFontWithOption : add font file
func (gp *PdfEngine) AddTTFFontWithOption(family string, ttfpath string, option TtfOption) error {

	data, err := env.FileExists(ttfpath)
	if err != nil {
		return err
	}
	rd := bytes.NewReader(data)
	return gp.AddTTFFontByReaderWithOption(family, rd, option)
}

// AddTTFFontByReaderWithOption adds font file by reader with option.
func (gp *PdfEngine) AddTTFFontByReaderWithOption(family string, rd Reader, option TtfOption) error {
	subsetFont := new(ttfSubsetObj)
	subsetFont.Init(func() *PdfEngine {
		return gp
	})
	subsetFont.SetTtfFontOption(option)
	subsetFont.SetFamily(family)
	err := subsetFont.SetTTFByReader(rd)
	if err != nil {
		return err
	}

	return gp.setSubsetFontObject(subsetFont, family, option)
}

// setSubsetFontObject sets ttfSubsetObj.
// The given ttfSubsetObj is expected to be configured in advance.
func (gp *PdfEngine) setSubsetFontObject(subsetFont *ttfSubsetObj, family string, option TtfOption) error {
	unicodemap := new(unicodeMap)
	unicodemap.Init(func() *PdfEngine {
		return gp
	})
	unicodemap.setProtection(gp.protection())
	unicodemap.SetPtrToSubsetFontObj(subsetFont)
	unicodeindex := gp.addObj(unicodemap)

	pdfdic := new(pdfDictionaryObj)
	pdfdic.Init(func() *PdfEngine {
		return gp
	})
	pdfdic.setProtection(gp.protection())
	pdfdic.SetPtrToSubsetFontObj(subsetFont)
	pdfdicindex := gp.addObj(pdfdic)

	subfontdesc := new(subfontDescriptorObj)
	subfontdesc.Init(func() *PdfEngine {
		return gp
	})
	subfontdesc.SetPtrToSubsetFontObj(subsetFont)
	subfontdesc.SetIndexObjPdfDictionary(pdfdicindex)
	subfontdescindex := gp.addObj(subfontdesc)

	cidfont := new(cidFontObj)
	cidfont.Init(func() *PdfEngine {
		return gp
	})
	cidfont.SetPtrToSubsetFontObj(subsetFont)
	cidfont.SetIndexObjSubfontDescriptor(subfontdescindex)
	cidindex := gp.addObj(cidfont)

	subsetFont.SetIndexObjCIDFont(cidindex)
	subsetFont.SetIndexObjUnicodeMap(unicodeindex)
	index := gp.addObj(subsetFont) //add หลังสุด

	if gp.indexOfProcSet != -1 {
		procset := gp.pdfObjs[gp.indexOfProcSet].(*procSetObj)
		if !procset.Relates.IsContainsFamilyAndStyle(family, option.Style&^Underline) {
			procset.Relates = append(procset.Relates, relateFont{Family: family, IndexOfObj: index, CountOfFont: gp.curr.CountOfFont, Style: option.Style &^ Underline})
			subsetFont.CountOfFont = gp.curr.CountOfFont
			gp.curr.CountOfFont++
		}
	}
	return nil
}
