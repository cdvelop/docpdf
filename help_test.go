package docpdf

import (
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
func setupDefaultA4PDF(filePath string, t *testing.T) *pdfEngine {
	pdf := pdfEngine{
		fileWrite: fw(filePath),
		log: func(a ...any) {

			t.Log(a...)
		},
	}
	pdf.Start(config{PageSize: *PageSizeA4})
	err := pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}

	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Fatal(err)
	}
	return &pdf
}

func pdfBenchEngine(filePath string, b *testing.B) *pdfEngine {
	pdf := pdfEngine{
		fileWrite: fw(filePath),
		log: func(a ...any) {
			b.Log(a...)
		},
	}

	return &pdf
}
