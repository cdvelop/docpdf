package pdfengine_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cdvelop/docpdf"
	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/pdfengine"
)

func GetFont(doc *docpdf.Document, fontPath string) (err error) {

	b, err := os.Open(fontPath)
	if err != nil {
		return err
	}
	err = doc.AddTTFFontByReader("Ubuntu-L", b)
	if err != nil {
		return err
	}
	return err
}

func TestSetY(t *testing.T) {
	var err error

	doc := docpdf.NewDocument(func(filename string, data []byte) error {
		// Use proper FileWriter function that writes to file
		return os.WriteFile(filename, data, 0644)
	})

	pdf := doc.PdfEngine
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	err = GetFont(doc, "test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}
	err = pdf.SetFont("Ubuntu-L", "", 14)
	if err != nil {
		t.Fatal(err)
	}
	pdf.SetMargins(0, 20, 0, 10)
	pdf.AddPage()

	var x float64 = 100
	var y float64 = 10
	for i := 0; i < 200; i++ {
		text := fmt.Sprintf("---------line no: %d -----------", i)
		//var textH float64 = 25 // if text height is 25px.
		pdf.SetXY(x, y)
		err = pdf.Text(text)
		if err != nil {
			t.Fatal(err)
		}
		y += 20
	}

	err = pdf.WritePdf(fmt.Sprintf("test/pagination/page_sety-%s.pdf", time.Now().Format("01-02-15-04-05")))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetNewY(t *testing.T) {
	var err error

	doc := docpdf.NewDocument(func(filename string, data []byte) error {
		return os.WriteFile(filename, data, 0644)
	})

	pdf := doc.PdfEngine
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	err = GetFont(doc, "test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}
	err = pdf.SetFont("Ubuntu-L", "", 14)
	if err != nil {
		t.Fatal(err)
	}
	pdf.SetMargins(0, 20, 0, 10)
	pdf.AddPage()

	var x float64 = 100
	var y float64 = 10
	for i := 0; i < 200; i++ {
		text := fmt.Sprintf("---------line no: %d -----------", i)
		var textH float64 = 25 // if text height is 25px.
		pdf.SetX(x)
		pdf.SetNewY(y, textH)
		y = pdf.GetY()
		err = pdf.Text(text)
		if err != nil {
			t.Fatal(err)
		}
		y += 20
	}

	err = pdf.WritePdf(fmt.Sprintf("test/pagination/page_setnewy-%s.pdf", time.Now().Format("01-02-15-04-05")))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetNewXY(t *testing.T) {
	var err error

	doc := docpdf.NewDocument(func(filename string, data []byte) error {
		return os.WriteFile(filename, data, 0644)
	})

	pdf := doc.PdfEngine
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	err = GetFont(doc, "test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}
	err = pdf.SetFont("Ubuntu-L", "", 14)
	if err != nil {
		t.Fatal(err)
	}
	pdf.SetMargins(0, 20, 0, 10)
	pdf.AddPage()

	var x float64 = 100
	var y float64 = 10
	for i := 0; i < 200; i++ {
		text := fmt.Sprintf("---------line no: %d -----------", i)
		var textH float64 = 25 // if text height is 25px.
		//pdf.SetX(x)
		pdf.SetNewXY(y, x, textH)
		y = pdf.GetY()
		err = pdf.Text(text)
		if err != nil {
			t.Fatal(err)
		}
		y += 20
	}

	err = pdf.WritePdf(fmt.Sprintf("test/pagination/page_setnewxy-%s.pdf", time.Now().Format("01-02-15-04-05")))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetNewYX(t *testing.T) {
	var err error

	doc := docpdf.NewDocument(func(filename string, data []byte) error {
		return os.WriteFile(filename, data, 0644)
	})

	pdf := doc.PdfEngine
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	err = GetFont(doc, "test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}
	err = pdf.SetFont("Ubuntu-L", "", 14)
	if err != nil {
		t.Fatal(err)
	}
	pdf.SetMargins(0, 20, 0, 10)
	pdf.AddPage()

	var x float64 = 100
	var y float64 = 10
	for i := 0; i < 200; i++ {
		text := fmt.Sprintf("---------line no: %d -----------", i)
		var textH float64 = 25 // if text height is 25px.
		pdf.SetNewY(y, textH)
		y = pdf.GetY()
		pdf.SetX(x) // must after pdf.SetNewY() called.
		err = pdf.Text(text)
		if err != nil {
			t.Fatal(err)
		}
		y += 20
	}

	err = pdf.WritePdf(fmt.Sprintf("test/pagination/page_setnewyx-%s.pdf", time.Now().Format("01-02-15-04-05")))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetNewYCheckHeight(t *testing.T) {
	var err error

	doc := docpdf.NewDocument(func(filename string, data []byte) error {
		return os.WriteFile(filename, data, 0644)
	})

	pdf := doc.PdfEngine
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	err = GetFont(doc, "test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}
	err = pdf.SetFont("Ubuntu-L", "", 14)
	if err != nil {
		t.Fatal(err)
	}
	pdf.SetMargins(0, 20, 0, 10)
	pdf.AddPage()

	y := 10.0
	pdf.SetNewY(y, 0)
	if y != pdf.GetY() {
		t.Fatal(" y != pdf.GetY()")
	}

	y = 1000.0
	pdf.SetNewY(y, 0)
	if y != pdf.GetY() {
		t.Fatal(" y != pdf.GetY()")
	}
}

func TestLineBreak(t *testing.T) {
	var err error

	doc := docpdf.NewDocument(func(filename string, data []byte) error {
		return os.WriteFile(filename, data, 0644)
	})

	pdf := doc.PdfEngine
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	err = GetFont(doc, "test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}
	err = pdf.SetFont("Ubuntu-L", "", 28)
	if err != nil {
		t.Fatal(err)
	}
	pdf.SetMargins(0, 20, 0, 10)
	pdf.AddPage()

	w := 500.0

	var breakOptionTests = []*pdfengine.BreakOption{
		&pdfengine.DefaultBreakOption,
		{
			Mode:           pdfengine.BreakModeIndicatorSensitive,
			BreakIndicator: ' ',
		},
	}

	y := (canvas.PageSizeA4.H/2 + 100.0*float64(len(breakOptionTests))) / 2
	linebreakText := strings.Repeat("MultiCell* methods don't respect linebreaking rules.", 2)
	for i, opt := range breakOptionTests {
		pdf.SetXY(canvas.PageSizeA4.W/2-w/2, y+100.0*float64(i))
		err = pdf.MultiCellWithOption(&canvas.Rect{
			W: w,
			H: 1000,
		}, linebreakText, pdfengine.CellOption{
			BreakOption: opt,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	err = pdf.WritePdf("test/pagination/page_linebreak.pdf")
	if err != nil {
		t.Fatal(err)
	}
}
