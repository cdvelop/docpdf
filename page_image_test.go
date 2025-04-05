package docpdf

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestPageWithImage(t *testing.T) {
	var err error

	doc := NewDocument(fw(fmt.Sprintf("test/pagination/page_image-%s.pdf", time.Now().Format("01-02-15-04-05"))), func(a ...any) {
		t.Log(a...)
	})

	pdf := doc.PdfEngine()
	pdf.Start(config{PageSize: *PageSizeA4})
	pdf.SetMargins(0, 20, 0, 10)
	pdf.AddPage()

	var x float64 = 100
	var y float64 = 20
	imgRect := &Rect{
		W: 354 * 72 / 120,
		H: 241 * 72 / 120,
	}
	for i := 0; i < 10; i++ {
		var imgHeight float64 = 241 * 72 / 120
		pdf.SetNewYIfNoOffset(y, imgHeight)
		y = pdf.GetY()
		err = pdf.Image("test/res/gopher01.jpg", x, y, imgRect)
		if err != nil {
			log.Fatal(err)
		}
		y += imgHeight
	}

	pdf.WritePdfFile()

}
