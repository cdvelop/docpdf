# docpdf

An open-source Go library for generating PDFs with a minimalist and intuitive API, similar to writing in Word. Optimized to run in the browser with WebAssembly without dependencies.

## Description

docpdf is a Go library that allows you to generate PDF documents with an intuitive and simple API. It is designed to be easy to use, with an approach similar to writing in a word processor like Word. The library is optimized to work in the browser with WebAssembly.

The main focus of this library is to compile into a compact binary size using TinyGo for frontend usage, as standard Go binaries tend to be large. This makes it ideal for web applications where binary size matters.

### TinyGo Compatibility Checklist

The following standard libraries will be replaced or modified as they are not 100% compatible with TinyGo, in order to reduce the binary size:

- [ ] bufio
- [ ] crypto/sha1
- [ ] errors
- [ ] fmt
- [ ] golang.org/x/image
- [ ] io
- [ ] os
- [ ] path/filepath
- [ ] sort
- [ ] strings
- [ ] sync
- [ ] time

## Usage Example:

### Page 1
![image](example/doc_test_P1.jpg)
### Page 2
![image](example/doc_test_P2.jpg)

This example shows the main features of the library:

```go
// Create a document with default settings
doc := docpdf.NewDocument(func(a ...any) {
    // Simple logging function
    fmt.Println(a...)
})

// Configure header with the new fluid API
doc.SetPageHeader().
    SetLeftText("Header Left").
    SetCenterText("Document Example").
    SetRightText("Confidential")
    // ShowOnFirstPage() - Optional to show on first page

// Configure footer with page numbering in X/Y format
doc.SetPageFooter().
    SetLeftText("Created: 2023-10-01").
    SetCenterText("footer center example").
    WithPageTotal(docpdf.Right).
    ShowOnFirstPage()

// Add logo
doc.AddImage("logo.png").Height(35).Inline().Draw()

// Add right-aligned date
doc.AddText("date: 2024-10-01").AlignRight().Inline().Draw()

// Add centered header
doc.AddHeader1("Example Document").AlignCenter().Draw()

// Add level 2 header
doc.AddHeader2("Section 1: Introduction").Draw()

// Add normal text
doc.AddText("This is a normal text paragraph that shows the basic capabilities of the docpdf library. " +
    "We can create documents with different text styles and formats.").Draw()

// Add text with different styles
doc.AddText("This text is in bold.").Bold().Draw()
doc.AddText("This text is in italic.").Italic().Draw()
doc.AddText("This text is right-aligned.").Regular().AlignRight().Draw()

// Add centered chart image
doc.AddImage("barchart.png").Height(150).AlignCenter().Draw()

// Add centered footnote
doc.AddFootnote("This is a footnote.").AlignCenter().Draw()

// Add image as right-aligned inline element
doc.AddImage("gopher-color.png").Height(50).Inline().AlignRight().Draw()

// Add level 3 header
doc.AddHeader3("Subsection 1.1: More examples").Draw()

// Add text with border
doc.AddText("This text has a border around it.").WithBorder().Draw()

// Comparison between normal and justified text
doc.AddHeader1("Comparison: Normal Text vs Justified Text").Draw()

doc.AddText("NORMAL TEXT (left-aligned):").Bold().Draw()
const multilineText = "This is a sample text that demonstrates normal text flow. The text continues across multiple lines to show how words wrap naturally at the margins. This creates a simple left-aligned paragraph that is easy to read. When text is not justified, it maintains consistent spacing between words while keeping a ragged right edge."
doc.AddText(multilineText).Draw()

doc.AddText("JUSTIFIED TEXT:").Bold().Draw()
doc.AddText(multilineText).Justify().Draw()

// Space between examples
doc.SpaceBefore(2)

// Tables
doc.AddHeader2("Section 2: Table Examples").Draw()
doc.AddText("This section demonstrates different table configuration options:").Draw()

// Sample data to be reused across tables
productData := []map[string]any{
    {"id": "001", "name": "Laptop Pro", "desc": "High-performance laptop", "qty": 2, "price": 1299.99, "discount": 10, "total": 2339.98},
    {"id": "002", "name": "Wireless Mouse", "desc": "Ergonomic mouse", "qty": 5, "price": 24.99, "discount": 5, "total": 118.70},
    {"id": "003", "name": "Monitor 27\"", "desc": "4K UHD display", "qty": 1, "price": 349.99, "discount": 15, "total": 297.49},
    {"id": "004", "name": "USB-C Hub", "desc": "Multi-port adapter", "qty": 3, "price": 39.99, "discount": 0, "total": 119.97},
}

// Comprehensive table with multiple features
doc.AddHeader3("1. Comprehensive Table with Multiple Features").Draw()
doc.AddText("This example shows many table configuration options combined:").Draw()

// Create table with extensive formatting options
comprehensiveTable := doc.NewTable(
    "Code|CC,W:8%",                // Centered header and content, 8% width
    "Product|W:15%",               // Default left alignment, 15% width
    "Description|W:25%",           // Default left alignment, 25% width
    "Quantity|HR,CR,S: pcs,W:13%", // Right-aligned header and content with "pcs" suffix, 13% width
    "Price|CR,P:$,W:13%",          // Right-aligned content with "$" prefix, 13% width
    "Discount|HR,CR,S:%,W:13%",    // Right-aligned header with "%" suffix, 13% width
    "Total|CR,P:$,W:13%",          // Right-aligned content with "$" prefix, 13% width
)

// Customize header style
comprehensiveTable.HeaderStyle(docpdf.CellStyle{
    BorderStyle: docpdf.BorderStyle{
        Top:      false,
        Left:     false,
        Bottom:   false,
        Right:    false,
        Width:    1.0,
        RGBColor: docpdf.RGBColor{R: 50, G: 50, B: 150},
    },
    FillColor: docpdf.RGBColor{R: 220, G: 230, B: 255},
    TextColor: docpdf.RGBColor{R: 20, G: 20, B: 100},
    Font:      docpdf.FontBold,
    FontSize:  12,
})

// Customize cell style
comprehensiveTable.CellStyle(docpdf.CellStyle{
    BorderStyle: docpdf.BorderStyle{
        Top:      false,
        Left:     true,
        Bottom:   true,
        Right:    true,
        Width:    0.5,
        RGBColor: docpdf.RGBColor{R: 180, G: 180, B: 220},
    },
    FillColor: docpdf.RGBColor{R: 255, G: 255, B: 255},
    TextColor: docpdf.RGBColor{R: 50, G: 50, B: 80},
    Font:      docpdf.FontRegular,
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

// Right-aligned table example
doc.AddHeader3("2. Right-aligned Table Example").Draw()
doc.AddText("Table with right alignment:").Draw()

// Create right-aligned table
rightTable := doc.NewTable("Code", "Product", "Price")
rightTable.AlignRight()

for _, product := range productData {
    rightTable.AddRow(product["id"], product["name"], product["price"])
}
rightTable.Draw()

doc.AddText("This table is right-aligned.").AlignRight().Draw()


// Save the document
err := doc.WritePdf("example_document.pdf")
if err != nil {
    fmt.Printf("Error saving the PDF: %v\n", err)
}
```

## Acknowledgements

This library would not have been possible without github repositories:
- signintech/gopdf
- phpdave11/gofpdi
- wcharczuk/go-chart
- golang/freetype


## Important

this library is still in development and not all features are implemented yet. please check the issues for more information.
