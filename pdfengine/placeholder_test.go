package pdfengine_test

import (
	"fmt"
	"testing"

	"github.com/cdvelop/docpdf/alignment"
	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/pdfengine"
)

func TestPlaceHolderText(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := pdfengine.PdfEngine{}
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < 5; i++ {
		pdf.AddPage()
		err = pdf.SetFont("LiberationSerif-Regular", "", 14)
		if err != nil {
			t.Log(err.Error())
			return
		}
		pdf.Br(10)
		pdf.SetX(250)
		err := pdf.Text(fmt.Sprintf("%d of ", i+1))
		if err != nil {
			t.Log(err.Error())
			return
		}
		err = pdf.PlaceHolderText("totalnumber", 30) //<-- create PlaceHolder
		if err != nil {
			t.Log(err.Error())
			return
		}
		pdf.Br(20)

		err = pdf.SetFont("LiberationSerif-Regular", "", 11)
		if err != nil {
			t.Log(err.Error())
			return
		}
		pdf.Text("content content content content content contents...")
	}

	err = pdf.FillInPlaceHoldText("totalnumber", fmt.Sprintf("%d", 5), alignment.Left) //<-- fillin text to PlaceHolder
	if err != nil {
		t.Log(err.Error())
		return
	}

	pdf.WritePdf("./test/out/placeholder_text.pdf")
}

func TestPlaceHolderText2(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := pdfengine.PdfEngine{}
	pdf.Start(pdfengine.Config{PageSize: *canvas.PageSizeA4})
	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < 5; i++ {
		pdf.AddPage()
		err = pdf.SetFont("LiberationSerif-Regular", "", 14)
		if err != nil {
			t.Log(err.Error())
			return
		}
		pdf.Br(10)
		pdf.SetX(250)
		pdf.Text("page")
		pagenumberPH := fmt.Sprintf("pagenumber_%d", i)
		err = pdf.PlaceHolderText(pagenumberPH, 20) //<-- create PlaceHolder
		if err != nil {
			t.Log(err.Error())
			return
		}

		err := pdf.Text("of")
		if err != nil {
			t.Log(err.Error())
			return
		}
		err = pdf.PlaceHolderText("totalnumber", 20) //<-- create PlaceHolder
		if err != nil {
			t.Log(err.Error())
			return
		}
		pdf.Br(20)

		err = pdf.SetFont("LiberationSerif-Regular", "", 11)
		if err != nil {
			t.Log(err.Error())
			return
		}
		pdf.Text("content content content content content contents...")

		err = pdf.FillInPlaceHoldText(pagenumberPH, fmt.Sprintf("%d", i+1), alignment.Center) //<-- fillin text to PlaceHolder
		if err != nil {
			t.Log(err.Error())
			return
		}

	}

	err = pdf.FillInPlaceHoldText("totalnumber", fmt.Sprintf("%d", 5), alignment.Center) //<-- fillin text to PlaceHolder
	if err != nil {
		t.Log(err.Error())
		return
	}

	pdf.WritePdf("./test/out/placeholder_text_2.pdf")
}
