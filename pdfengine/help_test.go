package pdfengine_test

import (
	"os"
	"testing"

	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/config"
	"github.com/cdvelop/docpdf/pdfengine"
)

func initTesting() error {
	err := os.MkdirAll("./test/out", 0777)
	if err != nil {
		return err
	}
	return nil
}

// setupDefaultA4PDF creates an A4 sized pdf with a plain configuration adding and setting the required fonts for
// further processing. Tests will fail in case adding or setting the font fails.
func setupDefaultA4PDF(t *testing.T) *pdfengine.PdfEngine {
	pdf := pdfengine.PdfEngine{}
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	err := pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}

	err = pdf.SetFont(config.NewFontStyle("LiberationSerif-Regular", 14))
	if err != nil {
		t.Fatal(err)
	}
	return &pdf
}
