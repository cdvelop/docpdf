package docpdf

import (
	"log"
	"os"
	"testing"
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
func setupDefaultA4PDF(t *testing.T) *pdfEngine {
	pdf := pdfEngine{}
	pdf.Start(config{PageSize: *PageSizeA4})
	err := pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}

	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		log.Fatal(err)
	}
	return &pdf
}
