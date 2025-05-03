package pdfengine_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/cdvelop/docpdf/alignment"
	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/errs"
	"github.com/cdvelop/docpdf/pdfengine"
)

func BenchmarkPdfWithImageHolder(b *testing.B) {

	err := initTesting()
	if err != nil {
		b.Error(err)
		return
	}

	pdf := pdfengine.PdfEngine{}
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	pdf.AddPage()
	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		b.Error(err)
		return
	}

	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		b.Fatal(err)
		return
	}

	bytesOfImg, err := os.ReadFile("./test/res/chilli.jpg")
	if err != nil {
		b.Error(err)
		return
	}

	imgH, err := pdfengine.ImageHolderByBytes(bytesOfImg)
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		pdf.DrawImageByHolder(imgH, 20.0, float64(i)*2.0, nil)
	}

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher")

	pdf.WritePdf("./test/out/image_bench.pdf")
}

func TestPdfWithImageHolder(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := setupDefaultA4PDF(t)
	pdf.AddPage()

	bytesOfImg, err := os.ReadFile("./test/res/PNG_transparency_demonstration_1.png")
	if err != nil {
		t.Error(err)
		return
	}

	imgH, err := pdfengine.ImageHolderByBytes(bytesOfImg)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.DrawImageByHolder(imgH, 20.0, 20, nil)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.DrawImageByHolder(imgH, 20.0, 200, nil)
	if err != nil {
		t.Error(err)
		return
	}

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher")

	pdf.WritePdf("./test/out/image_test.pdf")
}

func TestPdfWithImageHolderGif(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := setupDefaultA4PDF(t)
	pdf.AddPage()

	bytesOfImg, err := os.ReadFile("./test/res/gopher03.gif")
	if err != nil {
		t.Error(err)
		return
	}

	imgH, err := pdfengine.ImageHolderByBytes(bytesOfImg)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.DrawImageByHolder(imgH, 20.0, 20, nil)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.DrawImageByHolder(imgH, 20.0, 200, nil)
	if err != nil {
		t.Error(err)
		return
	}

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher")

	pdf.WritePdf("./test/out/image_test_gif.pdf")
}

func TestRetrievingNumberOfPdfPage(t *testing.T) {
	pdf := setupDefaultA4PDF(t)
	if pdf.GetNumberOfPages() != 0 {
		t.Error("Invalid starting number of pages, should be 0")
		return
	}
	pdf.AddPage()

	bytesOfImg, err := os.ReadFile("./test/res/gopher01.jpg")
	if err != nil {
		t.Error(err)
		return
	}

	imgH, err := pdfengine.ImageHolderByBytes(bytesOfImg)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.DrawImageByHolder(imgH, 20.0, 20, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if pdf.GetNumberOfPages() != 1 {
		t.Error(err)
		return
	}

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher")

	pdf.AddPage()

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher again")

	if pdf.GetNumberOfPages() != 2 {
		t.Error(err)
		return
	}

	pdf.WritePdf("./test/out/number_of_pages_test.pdf")
}

func TestImageCrop(t *testing.T) {
	pdf := setupDefaultA4PDF(t)
	if pdf.GetNumberOfPages() != 0 {
		t.Error("Invalid starting number of pages, should be 0")
		return
	}

	pdf.AddPage()

	bytesOfImg, err := os.ReadFile("./test/res/gopher01.jpg")
	if err != nil {
		t.Error(err)
		return
	}

	imgH, err := pdfengine.ImageHolderByBytes(bytesOfImg)
	if err != nil {
		t.Error(err)
		return
	}

	//err = pdf.DrawImageByHolder(imgH, 20.0, 20, nil)
	err = pdf.ImageByHolderWithOptions(imgH, pdfengine.ImageOptions{
		//VerticalFlip: true,
		//HorizontalFlip: true,
		Rect: &canvas.Rect{
			W: 100,
			H: 100,
		},
		Crop: &pdfengine.CropOptions{
			X:      0,
			Y:      0,
			Width:  10,
			Height: 100,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	if pdf.GetNumberOfPages() != 1 {
		t.Error(err)
		return
	}

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher")

	pdf.AddPage()

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher again")

	if pdf.GetNumberOfPages() != 2 {
		t.Error(err)
		return
	}

	pdf.WritePdf("./test/out/image_crop.pdf")
}

func BenchmarkAddTTFFontByReader(b *testing.B) {
	ttf, err := os.Open("test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		b.Error(err)
		return
	}
	defer ttf.Close()

	fontData, err := io.ReadAll(ttf)
	if err != nil {
		b.Error(err)
		return
	}

	for n := 0; n < b.N; n++ {
		pdf := &pdfengine.PdfEngine{}
		pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
		if err := pdf.AddTTFFontByReader("LiberationSerif-Regular", bytes.NewReader(fontData)); err != nil {
			return
		}
	}
}

/*
func TestBuffer(t *testing.T) {
	b := bytes.NewReader([]byte("ssssssss"))

	b1, err := ioutil.ReadAll(b)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("->%s\n", string(b1))
	b.Seek(0, 0)
	b2, err := ioutil.ReadAll(b)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("+>%s\n", string(b2))
}*/

func BenchmarkAddTTFFontData(b *testing.B) {
	ttf, err := os.Open("test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		b.Error(err)
		return
	}
	defer ttf.Close()

	fontData, err := io.ReadAll(ttf)
	if err != nil {
		b.Error(err)
		return
	}

	for n := 0; n < b.N; n++ {
		pdf := &pdfengine.PdfEngine{}
		pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
		if err := pdf.AddTTFFontData("LiberationSerif-Regular", fontData); err != nil {
			return
		}
	}
}

func TestReuseFontData(t *testing.T) {
	ttf, err := os.Open("test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}
	defer ttf.Close()

	fontData, err := io.ReadAll(ttf)
	if err != nil {
		t.Error(err)
		return
	}

	pdf1 := &pdfengine.PdfEngine{}
	rst1, err := generatePDFBytesByAddTTFFontData(pdf1, fontData)
	if err != nil {
		t.Error(err)
		return
	}

	// Reuse the parsed font data.
	pdf2 := &pdfengine.PdfEngine{}
	rst2, err := generatePDFBytesByAddTTFFontData(pdf2, fontData)
	if err != nil {
		t.Error(err)
		return
	}

	if bytes.Compare(rst1, rst2) != 0 {
		t.Fatal("The generated files must be exactly the same.")
		return
	}

	if err := writeFile("./test/out/result1_by_parsed_ttf_font.pdf", rst1, 0644); err != nil {
		t.Error(err)
		return
	}
	if err := writeFile("./test/out/result2_by_parsed_ttf_font.pdf", rst1, 0644); err != nil {
		t.Error(err)
		return
	}
}

func writeFile(name string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}

func generatePDFBytesByAddTTFFontData(pdf *pdfengine.PdfEngine, fontData []byte) ([]byte, error) {
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	if pdf.GetNumberOfPages() != 0 {
		return nil, errs.New("Invalid starting number of pages, should be 0")
	}

	if err := pdf.AddTTFFontData("LiberationSerif-Regular", fontData); err != nil {
		return nil, err
	}

	if err := pdf.SetFont("LiberationSerif-Regular", "", 14); err != nil {
		return nil, err
	}

	pdf.AddPage()
	if err := pdf.Text("Test PDF content."); err != nil {
		return nil, err
	}

	return pdf.GetBytesPdfReturnErr()
}

func TestWhiteTransparent(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}
	// create pdf.
	pdf := &pdfengine.PdfEngine{}
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	pdf.AddPage()

	var glyphNotFoundOfLiberationSerif []rune
	err = pdf.AddTTFFontWithOption("LiberationSerif-Regular", "test/res/LiberationSerif-Regular.ttf", pdfengine.TtfOption{
		OnGlyphNotFound: func(r rune) { //call when can not find glyph inside ttf file.
			glyphNotFoundOfLiberationSerif = append(glyphNotFoundOfLiberationSerif, r)
			//Log.Printf("glyph not found %c", r)
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Error(err)
		return
	}
	// write text.
	op := pdfengine.CellOption{Align: alignment.Left | alignment.Middle}
	rect := canvas.Rect{W: 20, H: 30}
	pdf.SetXY(350, 50)
	err = pdf.Cell(&rect, "あい")
	//err = pdf.CellWithOption(&rect, "あい", op)
	//err = pdf.CellWithOption(&rect, "あ", op)
	//err = pdf.CellWithOption(&rect, "a", op)
	if err != nil {
		t.Error(err)
		return
	}
	pdf.SetY(100)
	err = pdf.CellWithOption(&rect, "abcdef.", op)
	if err != nil {
		t.Error(err)
		return
	}

	//coz あ and い  not contain in "test/res/LiberationSerif-Regular.ttf"
	if len(glyphNotFoundOfLiberationSerif) != 2 {
		t.Error(err)
		return
	}

	//pdf.SetNoCompression()
	err = pdf.WritePdf("./test/out/white_transparent.pdf")
	if err != nil {
		t.Error(err)
		return
	}

}

func TestRectangle(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}
	// create pdf.
	pdf := pdfengine.PdfEngine{}
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	pdf.AddPage()

	pdf.SetStrokeColor(240, 98, 146)
	pdf.SetLineWidth(1)
	pdf.SetFillColor(255, 255, 255)
	// draw rectangle with round radius
	err = pdf.Rectangle(100.6, 150.8, 150.3, 379.3, "DF", 20, 10)
	if err != nil {
		t.Error(err)
		return
	}

	// draw rectangle with round radius but less point number
	err = pdf.Rectangle(200.6, 150.8, 250.3, 379.3, "DF", 20, 2)
	if err != nil {
		t.Error(err)
		return
	}

	pdf.SetStrokeColor(240, 98, 146)
	pdf.SetLineWidth(1)
	pdf.SetFillColor(255, 255, 255)
	// draw rectangle directly
	err = pdf.Rectangle(100.6, 50.8, 130, 150, "DF", 0, 0)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.WritePdf("./test/out/rectangle_with_round_corner.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestWhiteTransparent195(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}
	// create pdf.
	pdf := &pdfengine.PdfEngine{}
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	pdf.AddPage()

	var glyphNotFoundOfLiberationSerif []rune
	//err = pdf.AddTTFFontWithOption("LiberationSerif-Regular", "/Users/oneplus/Code/Work/gopdf_old/test/res/Meera-Regular.ttf", TtfOption{
	err = pdf.AddTTFFontWithOption("LiberationSerif-Regular", "test/res/LiberationSerif-Regular.ttf", pdfengine.TtfOption{
		OnGlyphNotFound: func(r rune) { //call when can not find glyph inside ttf file.
			glyphNotFoundOfLiberationSerif = append(glyphNotFoundOfLiberationSerif, r)
		},
		OnGlyphNotFoundSubstitute: func(r rune) rune {
			return r
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Error(err)
		return
	}
	// write text.
	op := pdfengine.CellOption{Align: alignment.Left | alignment.Middle}
	rect := canvas.Rect{W: 20, H: 30}
	pdf.SetXY(350, 50)
	//err = pdf.Cell(&rect, "あいうえ") // OK.
	//err = pdf.Cell(&rect, "あうう") // OK.
	err = pdf.CellWithOption(&rect, "あいうえ", op) // NG. "abcdef." is White/Transparent.
	//err = pdf.Cell(&rect, " あいうえ") // NG. "abcdef." is White/Transparent.
	// err = pdf.Cell(&rect, "あいうえ ") // NG. "abcdef." is White/Transparent.
	if err != nil {
		t.Error(err)
		return
	}
	pdf.SetY(100)
	err = pdf.CellWithOption(&rect, "abcกdef.", op)
	if err != nil {
		t.Error(err)
		return
	}

	//coz あ い う え  not contain in "test/res/LiberationSerif-Regular.ttf"
	// if len(glyphNotFoundOfLiberationSerif) != 4 {
	// 	t.Error(err)
	// 	return
	// }

	pdf.SetNoCompression()
	err = pdf.WritePdf("./test/out/white_transparent195.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestClearValue(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := pdfengine.PdfEngine{}
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4, Protection: pdfengine.PdfProtectionConfig{
		UseProtection: true,
		OwnerPass:     []byte("123456"),
		UserPass:      []byte("123456"),
	}})
	pdf.AddPage()
	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Fatal(err)
		return
	}

	bytesOfImg, err := os.ReadFile("./test/res/PNG_transparency_demonstration_1.png")
	if err != nil {
		t.Error(err)
		return
	}

	imgH, err := pdfengine.ImageHolderByBytes(bytesOfImg)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.DrawImageByHolder(imgH, 20.0, 20, nil)
	if err != nil {
		t.Error(err)
		return
	}

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher")
	pdf.SetInfo(pdfengine.PdfInfo{
		Title: "xx",
	})
	pdf.WritePdf("./test/out/test_clear_value.pdf")

	//reset
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})

	pdf2 := pdfengine.PdfEngine{}
	pdf2.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})

	//check
	if pdf.Margins() != pdf2.Margins() {
		t.Fatal("pdf.Margins != pdf2.Margins")
	}

	if len(pdf2.GetPdfObjs()) != len(pdf.GetPdfObjs()) {
		t.Fatalf("len(pdf2.GetPdfObjs()) != len(pdf.GetPdfObjs())")
	}

	if len(pdf.GetAnchors()) > 0 {
		t.Fatalf("len( pdf.GetAnchors()) = %d", len(pdf.GetAnchors()))
	}

	if len(pdf.GetIndexEncodingObjFonts()) != len(pdf2.GetIndexEncodingObjFonts()) {
		t.Fatalf("len(pdf.GetIndexEncodingObjFonts()) != len(pdf2.GetIndexEncodingObjFonts())")
	}

	if pdf.GetIndexOfContent() != pdf2.GetIndexOfContent() {
		t.Fatalf("pdf.GetIndexOfContent() != pdf2.GetIndexOfContent()")
	}

	if pdf.GetBuf().Len() > 0 {
		t.Fatalf("pdf.GetBuf().Len() > 0")
	}

	if pdf.GetPdfProtection() != nil {
		t.Fatalf("pdf.pdfProtection is not nil")
	}
	if pdf.GetEncryptionObjID() != 0 {
		t.Fatalf("encryptionObjID %d", pdf.GetEncryptionObjID())
	}

	if pdf.GetInfo() != nil {
		t.Fatalf("pdf.GetInfo() %v", pdf.GetInfo())
	}
}

func TestTextColor(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	// create pdf.
	pdf := pdfengine.PdfEngine{}
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	pdf.AddPage()
	err = pdf.AddTTFFont("LiberationSerif", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.SetFont("LiberationSerif", "", 14)
	if err != nil {
		t.Error(err)
		return
	}

	pdf.SetTextColor(255, 0, 2)
	pdf.Br(20)
	pdf.Cell(nil, "a")

	pdf.SetTextColorCMYK(0, 6, 14, 0)
	pdf.Br(20)
	pdf.Cell(nil, "b")

	err = pdf.WritePdf("./test/out/colored_text.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAddHeaderFooter(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	// create pdf.
	pdf := pdfengine.PdfEngine{}
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})

	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Fatal(err)
		return
	}

	pdf.AddHeader(func() {
		pdf.SetY(5)
		pdf.Cell(nil, "header")
	})
	pdf.AddFooter(func() {
		pdf.SetY(825)
		pdf.Cell(nil, "footer")
	})

	pdf.AddPage()
	pdf.SetY(400)
	pdf.Text("page 1 content")
	pdf.AddPage()
	pdf.SetY(400)
	pdf.Text("page 2 content")

	err = pdf.WritePdf("./test/out/header_footer.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestImportPagesFromFile(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	// Primero, crear un PDF simple para posteriormente importarlo
	samplePdfPath := "./test/out/sample_pdf_for_import.pdf"

	// Crear el PDF de prueba
	samplePdf := pdfengine.PdfEngine{}
	samplePdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	samplePdf.AddPage()

	err = samplePdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = samplePdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Error(err)
		return
	}

	samplePdf.SetXY(50, 50)
	samplePdf.Cell(nil, "Página 1")

	samplePdf.AddPage()
	samplePdf.SetXY(50, 50)
	samplePdf.Cell(nil, "Página 2")

	samplePdf.AddPage()
	samplePdf.SetXY(50, 50)
	samplePdf.Cell(nil, "Página 3")

	err = samplePdf.WritePdf(samplePdfPath)
	if err != nil {
		t.Error(err)
		return
	}

	// Ahora importar el PDF creado anteriormente
	pdf := pdfengine.PdfEngine{}
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})

	err = pdf.ImportPagesFromSource(samplePdfPath, "/MediaBox")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}
	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.SetPage(1)
	if err != nil {
		t.Error(err)
		return
	}
	pdf.SetXY(350, 50)
	err = pdf.Cell(&canvas.Rect{W: 20, H: 30}, "Hello World")
	if err != nil {
		t.Error(err)
		return
	}
	err = pdf.SetPage(2)
	if err != nil {
		t.Error(err)
		return
	}
	pdf.SetXY(350, 50)
	err = pdf.Cell(&canvas.Rect{W: 20, H: 30}, "Hello World")
	if err != nil {
		t.Error(err)
		return
	}
	err = pdf.SetPage(3)
	if err != nil {
		t.Error(err)
		return
	}
	pdf.SetXY(350, 50)
	err = pdf.Cell(&canvas.Rect{W: 20, H: 30}, "Hello World")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.WritePdf("./test/out/open-existing-pdf.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}
