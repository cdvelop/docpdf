package docpdf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDocumentAPIUsage(t *testing.T) {
	// Create a simple document with fileWriter function
	doc := NewDocument(func(filename string, data []byte) error {
		// For testing, we'll write to the specified file
		return os.WriteFile(filename, data, 0644)
	})

	// Setup header and footer with the new API
	doc.SetPageHeader().
		SetLeftText("Header Left").
		SetCenterText("Document Example").
		SetRightText("Confidential").
		// ShowOnFirstPage()

		// Add footer with page numbers in format X/Y
		doc.SetPageFooter().
		SetLeftText("Created: 2023-10-01").
		SetCenterText("footer center example").
		WithPageTotal(Right).
		ShowOnFirstPage()

	// add logo image
	doc.AddImage("test/res/logo.png").Height(35).Inline().Draw()

	// add date and time aligned to the right
	doc.AddText("date: 2024-10-01").AlignRight().Inline().Draw()

	// Add a centered header
	doc.AddHeader1("Example Document").AlignCenter().Draw()

	// Add a level 2 header
	doc.AddHeader2("Section 1: Introduction").Draw()

	// Add normal text
	doc.AddText("This is a normal text paragraph that shows the basic capabilities of the gopdf library. " +
		"We can create documents with different text styles and formats.").Draw()

	// Add text with different styles
	doc.AddText("This text is in bold.").Bold().Draw()

	doc.AddText("This text is in italic.").Italic().Draw()

	// Add right-aligned text (ensuring it's in regular style, not italic)
	doc.AddText("This text is right-aligned.").Regular().AlignRight().Draw()

	// Create and add a bar chart using the new API instead of static image
	barChart := doc.Chart().Bar().
		Title("Monthly Sales").
		Height(320).
		AlignCenter().
		BarWidth(50).
		BarSpacing(10)

	// Add data to the chart
	barChart.AddBar(120, "Jan").
		AddBar(140, "Feb").
		AddBar(160, "Mar").
		AddBar(180, "Apr").
		AddBar(120, "May").
		AddBar(140, "Jun")

	// Configurar el gráfico para mostrar los ejes
	barChart.WithAxis(true, true)

	// Renderizar el gráfico
	barChart.Draw()

	// Add a footnote (in italic by default)
	doc.AddFootnote("This is a footnote.").AlignCenter().Draw()

	// add gopher image as a right-aligned inline image
	doc.AddImage("test/res/gopher-color.png").Height(50).Inline().AlignRight().Draw()
	// Add level 3 header
	doc.AddHeader3("Subsection 1.1: More examples").Draw()

	// Add text with a border
	doc.AddText("This text has a border around it.").WithBorder().Draw()

	// Compare justified vs non-justified
	doc.AddHeader1("Comparison: Normal Text vs Justified Text").Draw()

	doc.AddText("NORMAL TEXT (left-aligned):").Bold().Draw()
	// Normal text (left-aligned)
	const multilineText = "This is a sample text that demonstrates normal text flow. The text continues across multiple lines to show how words wrap naturally at the margins. This creates a simple left-aligned paragraph that is easy to read. When text is not justified, it maintains consistent spacing between words while keeping a ragged right edge."
	doc.AddText(multilineText).Draw()

	// Justified text
	doc.AddText("JUSTIFIED TEXT:").Bold().Draw()
	doc.AddText(multilineText).Justify().Draw()

	// Space between examples
	doc.SpaceBefore(2)
	// Add example of table usage
	doc.AddHeader2("Section 2: Table Examples").Draw()
	doc.AddText("This section demonstrates different table configuration options:").Draw()

	// Define sample data sets that will be reused across all tables
	productData := []map[string]any{
		{"id": "001", "name": "Laptop Pro", "desc": "High-performance laptop", "qty": 2, "price": 1299.99, "discount": 10, "total": 2339.98},
		{"id": "002", "name": "Wireless Mouse", "desc": "Ergonomic mouse", "qty": 5, "price": 24.99, "discount": 5, "total": 118.70},
		{"id": "003", "name": "Monitor 27\"", "desc": "4K UHD display", "qty": 1, "price": 349.99, "discount": 15, "total": 297.49},
		{"id": "004", "name": "USB-C Hub", "desc": "Multi-port adapter", "qty": 3, "price": 39.99, "discount": 0, "total": 119.97},
	}

	// Comprehensive table example with many API features combined
	doc.AddHeader3("1. Comprehensive Table with Multiple Features").Draw()
	doc.AddText("This example shows many table configuration options combined:").Draw()

	// Create a table with extensive formatting options
	comprehensiveTable := doc.NewTable(
		"Code|CC,W:8%",                // Centered header and centered content, 8% width
		"Product|W:15%",               // Default left alignment, 15% width
		"Description|W:25%",           // Default left alignment, 25% width
		"Quantity|HR,CR,S: pcs,W:13%", // Right-aligned header and content with "pcs" suffix, 13% width
		"Price|CR,P:$,W:13%",          // Right-aligned content with "$" prefix, 13% width
		"Discount|HR,CR,S:%,W:13%",    // Right-aligned header with "%" suffix, 13% width
		"Total|CR,P:$,W:13%",          // Right-aligned content with "$" prefix, 13% width
	)

	// Customize header style
	comprehensiveTable.HeaderStyle(CellStyle{
		BorderStyle: BorderStyle{
			Top:      false,
			Left:     false,
			Bottom:   false,
			Right:    false,
			Width:    1.0,
			RGBColor: RGBColor{R: 50, G: 50, B: 150},
		},
		FillColor: RGBColor{R: 220, G: 230, B: 255},
		TextColor: RGBColor{R: 20, G: 20, B: 100},
		Font:      FontBold,
		FontSize:  12,
	})

	// Customize cell style
	comprehensiveTable.CellStyle(CellStyle{
		BorderStyle: BorderStyle{
			Top:      false,
			Left:     true,
			Bottom:   true,
			Right:    true,
			Width:    0.5,
			RGBColor: RGBColor{R: 180, G: 180, B: 220},
		},
		FillColor: RGBColor{R: 255, G: 255, B: 255},
		TextColor: RGBColor{R: 50, G: 50, B: 80},
		Font:      FontRegular,
		FontSize:  11,
	})

	// Add rows with data
	for _, product := range productData {
		comprehensiveTable.AddRow(
			product["id"],
			product["name"],
			product["desc"],
			product["qty"],
			product["price"],
			product["discount"],
			product["total"],
		)
	}

	comprehensiveTable.Draw()

	// Keep only the right-aligned table example
	doc.AddHeader3("2. Right-aligned Table Example").Draw()
	doc.AddText("Table with right alignment:").Draw()

	// Create a right-aligned table with specific column widths
	rightTable := doc.NewTable("Code", "Product", "Price")
	rightTable.AlignRight()

	for _, product := range productData {
		rightTable.AddRow(product["id"], product["name"], product["price"])
	}
	rightTable.Draw()

	doc.AddText("This table is right-aligned.").AlignRight().Draw()

	// add page for checking page header and footer
	doc.AddPage()

	// Create output directory if it doesn't exist
	outDir := "test/out"
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		t.Fatalf("Error creating output directory: %v", err)
	}

	// Set the output file path
	outFilePath := filepath.Join(outDir, "doc_test.pdf")

	// Save the document to the specified location
	err = doc.WritePdf(outFilePath)
	if err != nil {
		t.Fatalf("Error writing PDF: %v", err)
	}

	absPath, _ := filepath.Abs(outFilePath)
	t.Logf("PDF created successfully at: %s", absPath)
}

func TestPageSizeOptions(t *testing.T) {
	// Test directory setup
	outDir := "test/out"
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		t.Fatalf("Error creating output directory: %v", err)
	}

	// Common fileWriter for all tests
	fileWriter := func(filename string, data []byte) error {
		return os.WriteFile(filename, data, 0644)
	}

	// Test 1: Using predefined page size (PageSizeA4)
	t.Run("PredefinedPageSize", func(t *testing.T) {
		doc := NewDocument(fileWriter, PageSizeA4)

		doc.AddText("This document uses predefined PageSizeA4").Bold().AlignCenter().Draw()

		outFilePath := filepath.Join(outDir, "test_predefined_pagesize.pdf")
		err := doc.WritePdf(outFilePath)
		if err != nil {
			t.Fatalf("Error writing PDF: %v", err)
		}
	})

	// Test 2: Using custom PageSize with mm units
	t.Run("CustomPageSizeWithUnits", func(t *testing.T) {
		// Create custom A5 size in mm (148mm x 210mm)
		doc := NewDocument(fileWriter, PageSize{Width: 148, Height: 210, Unit: UnitMM})

		doc.AddText("This document uses custom PageSize (A5 in mm)").Bold().AlignCenter().Draw()

		outFilePath := filepath.Join(outDir, "test_custom_pagesize.pdf")
		err := doc.WritePdf(outFilePath)
		if err != nil {
			t.Fatalf("Error writing PDF: %v", err)
		}
	})

	// Test 3: Using custom PageSize with inches
	t.Run("CustomPageSizeInches", func(t *testing.T) {
		// Create custom size in inches (8.5 x 11 inches - US Letter)
		doc := NewDocument(fileWriter, PageSize{Width: 8.5, Height: 11, Unit: UnitIN})

		doc.AddText("This document uses custom PageSize (Letter in inches)").Bold().AlignCenter().Draw()

		outFilePath := filepath.Join(outDir, "test_custom_pagesize_inches.pdf")
		err := doc.WritePdf(outFilePath)
		if err != nil {
			t.Fatalf("Error writing PDF: %v", err)
		}
	})

	// Test 4: Combining PageSize with other options
	t.Run("CombinedOptions", func(t *testing.T) {
		doc := NewDocument(fileWriter,
			// Custom page size (A4 landscape in mm)
			PageSize{Width: 297, Height: 210, Unit: UnitMM},
			// Custom margins
			Margins{Left: 20, Top: 15, Right: 20, Bottom: 15},
		)

		doc.AddText("This document combines custom PageSize with custom margins").Bold().AlignCenter().Draw()

		outFilePath := filepath.Join(outDir, "test_combined_options.pdf")
		err := doc.WritePdf(outFilePath)
		if err != nil {
			t.Fatalf("Error writing PDF: %v", err)
		}
	})
}
