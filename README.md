# docpdf

An open-source Go library for generating PDFs with a minimalist and intuitive API, similar to writing in Word. Optimized to run in the browser with WebAssembly without dependencies.

## Description

docpdf is a Go library that allows you to generate PDF documents with an intuitive and simple API. It is designed to be easy to use, with an approach similar to writing in a word processor like Word. The library is optimized to work in the browser with WebAssembly.

The main focus of this library is to compile into a compact binary size using TinyGo for frontend usage, as standard Go binaries tend to be large. This makes it ideal for web applications where binary size matters.

## Writer Interface

The library now uses a `Writer` interface for PDF output, which makes it compatible with WebAssembly environments by removing direct dependencies on the `os` package:


```go
// This allows you to use any type that implements the:
type writer interface {
	Write(p []byte) (n int, err error)
}
// interface to save your PDF documents
```

```go
// Using a file writer (for backend/desktop applications)
fileWrite := &fileWrite{filePath: "output.pdf"}
doc := NewDocument(fileWrite, fmt.Println)

```

### TinyGo Compatibility Checklist

The following standard libraries will be replaced or modified as they are not 100% compatible with TinyGo, in order to reduce the binary size:

- [ ] bufio
- [ ] crypto/sha1
- [ ] errors
- [ ] fmt
- [ ] golang.org/x/image
- [ ] io
- [x] os (removed direct dependency by using Writer interface)
- [ ] path/filepath
- [ ] sort
- [ ] strings
- [ ] sync
- [ ] time

## Page Size Options

The library offers multiple ways to specify page sizes for your documents:

1. **Using predefined page sizes**:
   ```go
   doc := NewDocument(fw, fmt.Println, PageSizeA4)  // Use A4 page size
   doc := NewDocument(fw, fmt.Println, PageSizeLetter)  // Use Letter page size
   ```

2. **Using the new PageSize struct with unit specification**:
   ```go
   // Create an A4 page size (210mm x 297mm) with millimeter units
   doc := NewDocument(fw,fmt.Println, PageSize{Width: 210, Height: 297, Unit: UnitMM})
   
   // Create a US Letter page size (8.5in x 11in) with inch units
   doc := NewDocument(fw,fmt.Println, PageSize{Width: 8.5, Height: 11, Unit: UnitIN})
   ```

3. **Combining page size with other options**:
   ```go
   doc := NewDocument(fw,fmt.Println, 
      PageSize{Width: 210, Height: 297, Unit: UnitMM},  // A4 size
      Margins{Left: 15, Top: 10, Right: 10, Bottom: 10}  // Custom margins
   )
   ```

## Usage Example:

### Page 1
![image](example/doc_test_P1.jpg)
### Page 2
![image](example/doc_test_P2.jpg)

This example shows the main features of the library:

```go
	// Create a simple document with default settings
	doc := NewDocument(fw, func(a ...any) {
		// Simple logger that does nothing for this test
		t.Log(a...)
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
	barChart := doc.AddBarChart().
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
		FillColor: RGBColor{R: 255, G: 255, 255},
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
	err = doc.WritePdfFile(outFilePath)
```

## Acknowledgements

This library would not have been possible without github repositories:
- signintech/gopdf
- phpdave11/gofpdi
- wcharczuk/go-chart
- golang/freetype


## Important

this library is still in development and not all features are implemented yet. please check the issues for more information.
