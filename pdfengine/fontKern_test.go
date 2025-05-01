package pdfengine

import (
	"fmt"
	"testing"

	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/errs"
)

func TestKern01(t *testing.T) {
	Wo, err := kern01("test/res/LiberationSerif-Regular.ttf", "LiberationSerif-Regular", 'W', 'o')
	if err != nil {
		t.Error(err)
		return
	}

	if Wo != -80 {
		t.Errorf("Wo must be -80 (but %d)", Wo)
		//return
	}

	Wi, err := kern01("test/res/LiberationSerif-Regular.ttf", "LiberationSerif-Regular", 'W', 'i')
	if err != nil {
		t.Error(err)
		return
	}

	if Wi != -40 {
		t.Errorf("Wi must be -40 (but %d)", Wi)
		//return
	}

}

func kern01(font string, prefix string, leftRune rune, rightRune rune) (int, error) {
	pdf := PdfEngine{}
	pdf.Start(Config{Unit: canvas.UnitPT, PageSize: *canvas.PageSizeA4})
	pdf.AddPage()
	err := pdf.AddTTFFontWithOption(prefix, font, TtfOption{
		UseKerning: true,
	})
	if err != nil {
		fmt.Errorf(err.Error())
		return 0, err
	}

	err = pdf.SetFont(prefix, "", 50)
	if err != nil {
		fmt.Errorf(err.Error())
		return 0, err
	}

	gindexleftRune, err := pdf.Curr.FontISubset.CharCodeToGlyphIndex(leftRune)
	if err != nil {
		return 0, err
	}

	gindexrightRune, err := pdf.Curr.FontISubset.CharCodeToGlyphIndex(rightRune)
	if err != nil {
		return 0, err
	}
	//fmt.Printf("gindexleftRune = %d  gindexrightRune=%d \n", gindexleftRune, gindexrightRune)
	kernTb := pdf.Curr.FontISubset.ttfp.Kern()

	//fmt.Printf("UnitsPerEm = %d\n", pdf.Curr.FontISubset.ttfp.UnitsPerEm())

	//fmt.Printf("len =%d\n", len(kernTb.Kerning))
	for left, kval := range kernTb.Kerning {
		if left == gindexleftRune {
			for right, val := range kval {
				if right == gindexrightRune {
					//fmt.Printf("left=%d right= %d  val=%d\n", left, right, val)
					valPdfUnit := convertTTFUnit2PDFUnit(int(val), int(pdf.Curr.FontISubset.ttfp.UnitsPerEm()))
					return valPdfUnit, nil
				}
			}
			break
		}
	}

	return 0, errs.New("not found")
}
