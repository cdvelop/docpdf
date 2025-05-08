package pdfengine_test

import (
	"strings"
	"testing"

	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/config"
	"github.com/cdvelop/docpdf/errs"
	"github.com/cdvelop/docpdf/pdfengine"
)

func TestJustify(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := setupDefaultA4PDF(t)
	pdf.AddPage()
	// Test ParseTextForJustification
	t.Run("ParseTextForJustification", func(t *testing.T) {
		text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
		width := 200.0

		jText, err := pdf.ParseTextForJustification(text, width)
		if err != nil {
			t.Errorf("Error analyzing text for justification: %v", err)
			return
		}

		if jText.WordCount() != 8 {
			t.Errorf("Incorrect number of words. Expected: 8, Got: %d", jText.WordCount())
		}

		if jText.SpaceCount() != 7 {
			t.Errorf("Incorrect number of spaces. Expected: 7, Got: %d", jText.SpaceCount())
		}

		// Verify that spaces are greater than zero
		for i, space := range jText.GetSpaces() {
			if space <= 0 {
				t.Errorf("Space %d is not positive: %f", i, space)
			}
		}

		// Verify original string is preserved
		if jText.GetOriginalString() != text {
			t.Errorf("Original string not preserved. Expected: %s, Got: %s", text, jText.GetOriginalString())
		}

		// Verify width is set correctly
		if jText.GetWidth() != width {
			t.Errorf("Width not set correctly. Expected: %f, Got: %f", width, jText.GetWidth())
		}
	})

	// Test with empty text
	t.Run("EmptyText", func(t *testing.T) {
		_, err := pdf.ParseTextForJustification("", 200.0)
		if err != errs.EmptyString {
			t.Errorf("Empty text should return errs.EmptyString, but got: %v", err)
		}
	})
	// Test with single word
	t.Run("SingleWord", func(t *testing.T) {
		jText, err := pdf.ParseTextForJustification("Lorem", 200.0)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		if jText.WordCount() != 1 {
			t.Errorf("Incorrect number of words. Expected: 1, Got: %d", jText.WordCount())
		}

		// Verify original string is preserved
		if jText.GetOriginalString() != "Lorem" {
			t.Errorf("Original string not preserved. Expected: %s, Got: %s", "Lorem", jText.GetOriginalString())
		}
	})

	// Create a PDF with justified text for visual verification
	rect := &canvas.Rect{W: 200, H: 100}

	// Non-justified text (left-aligned)
	pdf.SetY(50)
	err = pdf.MultiCell(rect, "This is normal left-aligned text. It should show an irregular margin on the right.")
	if err != nil {
		t.Error(err)
		return
	}

	// Justified text
	pdf.SetY(100)
	opt := pdfengine.CellOption{
		Align: config.Justify,
	}
	err = pdf.MultiCellWithOption(rect, "This is justified text. It should show uniform canvas.Margins on both sides except for the last line.", opt)
	if err != nil {
		t.Error(err)
		return
	}

	// Justified text with long paragraph
	pdf.SetY(150)
	longText := strings.Repeat("This is a long text that should be justified correctly with uniform spaces. ", 3)
	err = pdf.MultiCellWithOption(rect, longText, opt)
	if err != nil {
		t.Error(err)
		return
	}

	// Convenience method for justifying text
	pdf.SetY(300)
	err = pdf.MultiCellJustified(rect, "This text uses the MultiCellJustified convenience method which internally uses the config.Justify option.")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.WritePdf("./test/out/justify_test.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}
