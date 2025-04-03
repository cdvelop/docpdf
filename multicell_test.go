package docpdf

import (
	"testing"
)

func TestSplitTextWithOptions(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := setupDefaultA4PDF(t)

	var splitTextTests = []struct {
		name string
		in   string
		opts *breakOption
		exp  []string
	}{
		{
			"strict breaks no separator",
			"Lorem ipsum dolor sit amet, consetetur",
			&defaultBreakOption,
			[]string{"Lorem ipsum dol", "or sit amet, conse", "tetur"},
		},
		{
			"no options given",
			"Lorem ipsum dolor sit amet, consetetur",
			nil,
			[]string{"Lorem ipsum dol", "or sit amet, conse", "tetur"},
		},
		{
			"strict breaks with separator",
			"Lorem ipsum dolor sit amet, consetetur",
			&breakOption{
				Separator: "-",
				Mode:      breakModeStrict,
			},
			[]string{"Lorem ipsum d-", "olor sit amet, c-", "onsetetur"},
		},
		{
			"text with possible word-wrap",
			"Lorem ipsum dolor sit amet, consetetur",
			&breakOption{
				BreakIndicator: ' ',
				Mode:           breakModeIndicatorSensitive,
			},
			[]string{"Lorem ipsum", "dolor sit amet,", "consetetur"},
		},
		{
			"text without possible word-wrap",
			"Loremipsumdolorsitamet,consetetur",
			&breakOption{
				BreakIndicator: ' ',
				Mode:           breakModeIndicatorSensitive,
			},
			[]string{"Loremipsumdolo", "rsitamet,consetet", "ur"},
		},
		{
			"text with only empty spaces",
			"                                                ",
			&breakOption{
				BreakIndicator: ' ',
				Mode:           breakModeIndicatorSensitive,
			},
			[]string{"                           ", "                    "},
		},
	}

	for _, tt := range splitTextTests {
		t.Run(tt.name, func(t *testing.T) {
			lines, err := pdf.SplitTextWithOption(tt.in, 100, tt.opts)
			if err != nil {
				t.Fatal(err)
			}
			if len(lines) != len(tt.exp) {
				t.Fatalf("amount of expected and split lines invalid. Expected: %d, result: %d", len(tt.exp), len(lines))
			}
			for i, e := range tt.exp {
				if e != lines[i] {
					t.Fatalf("split text invalid. Expected: '%s', result: '%s'", e, lines[i])
				}
			}
		})
	}
}

// TestMultiCell prueba la funcionalidad básica de la función MultiCell
func TestMultiCell(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := setupDefaultA4PDF(t)
	pdf.AddPage()

	// Texto simple para MultiCell
	rect := &Rect{W: 200, H: 100}
	err = pdf.MultiCell(rect, "This is a test text for MultiCell that should fit within the provided rectangle.")
	if err != nil {
		t.Error(err)
		return
	}

	pdf.SetY(150)
	// Prueba con texto más largo
	err = pdf.MultiCell(rect, "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor. Cras elementum ultrices diam.")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.WritePdf("./test/out/multicell_test.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}

// TestMultiCellWithOption prueba la funcionalidad de MultiCellWithOption
func TestMultiCellWithOption(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := setupDefaultA4PDF(t)
	pdf.AddPage()

	rect := &Rect{W: 200, H: 100}
	opt := cellOption{
		Align: Left | Top,
	}

	// Prueba con opciones básicas
	err = pdf.MultiCellWithOption(rect, "This is a text to test MultiCellWithOption with left and top alignment.", opt)
	if err != nil {
		t.Error(err)
		return
	}

	// Prueba con opciones de ruptura personalizadas
	pdf.SetY(150)
	opt.breakOption = &breakOption{
		Mode:           breakModeIndicatorSensitive,
		BreakIndicator: ' ',
	}
	err = pdf.MultiCellWithOption(rect, "This text should use a space-sensitive algorithm for line breaks.", opt)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.WritePdf("./test/out/multicell_with_option_test.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}

// TestIsFitMultiCell prueba la función IsFitMultiCell
func TestIsFitMultiCell(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := setupDefaultA4PDF(t)

	// Crear un rectángulo lo suficientemente grande para contener el texto
	rectGrande := &Rect{W: 200, H: 100}
	textoCorto := "Short text."

	fits, _, err := pdf.IsFitMultiCell(rectGrande, textoCorto)
	if err != nil {
		t.Error(err)
		return
	}

	if !fits {
		t.Errorf("The text should fit in the large rectangle but IsFitMultiCell returned false")
	}

	// Probar con un rectángulo pequeño y texto largo
	rectPequeno := &Rect{W: 50, H: 20}
	textoLargo := "This is a very long text that definitely should not fit in a small rectangle."

	fits, _, err = pdf.IsFitMultiCell(rectPequeno, textoLargo)
	if err != nil {
		t.Error(err)
		return
	}

	if fits {
		t.Errorf("The long text should not fit in the small rectangle but IsFitMultiCell returned true")
	}
}

// TestIsFitMultiCellWithNewline prueba la función IsFitMultiCellWithNewline
func TestIsFitMultiCellWithNewline(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := setupDefaultA4PDF(t)

	// Texto con saltos de línea
	textoConSaltos := "First line\nSecond line\nThird line"
	rectSuficiente := &Rect{W: 200, H: 100}

	fits, _, err := pdf.IsFitMultiCellWithNewline(rectSuficiente, textoConSaltos)
	if err != nil {
		t.Error(err)
		return
	}

	if !fits {
		t.Errorf("The text with line breaks should fit in the rectangle but IsFitMultiCellWithNewline returned false")
	}

	// Rectángulo insuficiente
	rectInsuficiente := &Rect{W: 200, H: 10}
	fits, _, err = pdf.IsFitMultiCellWithNewline(rectInsuficiente, textoConSaltos)
	if err != nil {
		t.Error(err)
		return
	}

	if fits {
		t.Errorf("The text with line breaks should not fit in the small rectangle but IsFitMultiCellWithNewline returned true")
	}
}
