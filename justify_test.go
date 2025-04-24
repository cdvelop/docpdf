package docpdf

import (
	"strings"
	"testing"

	"github.com/cdvelop/docpdf/errs"
)

func TestJustify(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := setupDefaultA4PDF(t)
	pdf.AddPage()

	// Test parseTextForJustification
	t.Run("parseTextForJustification", func(t *testing.T) {
		text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
		width := 200.0

		jText, err := parseTextForJustification(pdf, text, width)
		if err != nil {
			t.Errorf("Error analyzing text for justification: %v", err)
			return
		}

		if len(jText.words) != 8 {
			t.Errorf("Incorrect number of words. Expected: 8, Got: %d", len(jText.words))
		}

		if len(jText.spaces) != 7 {
			t.Errorf("Incorrect number of spaces. Expected: 7, Got: %d", len(jText.spaces))
		}

		// Verify that spaces are greater than zero
		for i, space := range jText.spaces {
			if space <= 0 {
				t.Errorf("Space %d is not positive: %f", i, space)
			}
		}
	})

	// Test with empty text
	t.Run("EmptyText", func(t *testing.T) {
		_, err := parseTextForJustification(pdf, "", 200.0)
		if err != errs.EmptyString {
			t.Errorf("Empty text should return errs.EmptyString, but got: %v", err)
		}
	})

	// Test with single word
	t.Run("SingleWord", func(t *testing.T) {
		jText, err := parseTextForJustification(pdf, "Lorem", 200.0)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		if len(jText.words) != 1 {
			t.Errorf("Incorrect number of words. Expected: 1, Got: %d", len(jText.words))
		}
	})

	// Create a PDF with justified text for visual verification
	rect := &Rect{W: 200, H: 100}

	// Non-justified text (left-aligned)
	pdf.SetY(50)
	err = pdf.MultiCell(rect, "This is normal left-aligned text. It should show an irregular margin on the right.")
	if err != nil {
		t.Error(err)
		return
	}

	// Justified text
	pdf.SetY(100)
	opt := cellOption{
		Align: Justify,
	}
	err = pdf.MultiCellWithOption(rect, "This is justified text. It should show uniform margins on both sides except for the last line.", opt)
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
	err = pdf.MultiCellJustified(rect, "This text uses the MultiCellJustified convenience method which internally uses the Justify option.")
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
