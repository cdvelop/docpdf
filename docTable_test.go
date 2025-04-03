package docpdf

import (
	"path/filepath"
	"testing"
)

func TestTableColumnWidths(t *testing.T) {
	// Create output directory if it doesn't exist
	outDir := "test/out"

	// Test cases for different width configurations
	testCases := []struct {
		name      string
		headers   []string
		checkType string // "percent", "fixed", or "auto"
	}{
		{
			name: "Percentage_Widths",
			headers: []string{
				"Col1|W:30%",
				"Col2|W:20%",
				"Col3|W:50%",
			},
			checkType: "percent",
		},
		{
			name: "Fixed_Widths",
			headers: []string{
				"Column 1|W:100",
				"Column 2|W:150",
				"Column 3|W:80",
			},
			checkType: "fixed",
		},
		{
			name: "Auto_Widths",
			headers: []string{
				"Column One",
				"Col 2",
				"Third Column With Long Title",
			},
			checkType: "auto",
		},
	}

	// Run test for each case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new document for this test
			doc := NewDocument(func(a ...any) {
				t.Log(a...)
			})

			// Create table with the test headers
			table := doc.NewTable(tc.headers...)

			// Add some sample data
			table.AddRow("Data 1", "Data 2", "Data 3")
			table.AddRow("Longer data entry", "Short", "Medium length data")

			// Save the column widths before drawing for validation
			actualWidths := make([]float64, len(table.columns))
			for i, col := range table.columns {
				actualWidths[i] = col.width
			}

			// Calculate total width
			totalWidth := 0.0
			for _, width := range actualWidths {
				totalWidth += width
			}

			// Validate width calculation based on check type
			switch tc.checkType {
			case "percent":
				validatePercentageWidths(t, table, doc, actualWidths, totalWidth)
			case "fixed":
				validateFixedWidths(t, table, actualWidths, expectedFixedWidths())
			case "auto":
				validateAutoWidths(t, table, actualWidths)
			}

			// Create a PDF for visual inspection
			outFilePath := filepath.Join(outDir, "table_"+tc.name+".pdf")
			table.Draw()
			err := doc.WritePdf(outFilePath)
			if err != nil {
				t.Fatalf("Error writing PDF: %v", err)
			}
		})
	}

	// Test mixing width types (should ensure consistent approach)
	t.Run("Mixing_Width_Types", func(t *testing.T) {
		doc := NewDocument(func(a ...any) {
			t.Log(a...)
		})

		// Test mixing percentage and fixed widths (should convert all to percentage)
		headers := []string{
			"Col1|W:30%",
			"Col2|W:100", // Fixed width that should be converted to percentage
			"Col3",       // Auto width that should be converted to percentage
		}

		table := doc.NewTable(headers...)
		table.AddRow("Data 1", "Data 2", "Data 3")

		// Calculate total width
		totalWidth := 0.0
		actualWidths := make([]float64, len(table.columns))
		for i, col := range table.columns {
			actualWidths[i] = col.width
			totalWidth += col.width
		}

		// When mixing width types with percentage, all should be treated as percentage
		// So the total width should equal the available width
		if !approximatelyEqual(totalWidth, doc.contentAreaWidth, 0.1) {
			t.Errorf("When mixing width types with percentage, total width (%f) must equal available width (%f)",
				totalWidth, doc.contentAreaWidth)
		}

		// Create a PDF for visual inspection
		outFilePath := filepath.Join(outDir, "table_mixed_widths.pdf")
		table.Draw()
		err := doc.WritePdf(outFilePath)
		if err != nil {
			t.Fatalf("Error writing PDF: %v", err)
		}
	})
}

// Helper function to get expected fixed widths for the test
func expectedFixedWidths() []float64 {
	return []float64{100, 150, 80}
}

// Helper function to validate percentage width calculations
func validatePercentageWidths(t *testing.T, table *docTable, doc *Document, actualWidths []float64, totalWidth float64) {

	// Check that total width equals available width exactly (within a small margin of error)
	if !approximatelyEqual(totalWidth, doc.contentAreaWidth, 0.1) {
		t.Errorf("When using percentage widths, total table width (%f) must equal available width (%f)",
			totalWidth, doc.contentAreaWidth)
	}

	// Verify proportions match the specified percentages
	expectedRatios := []float64{0.3, 0.2, 0.5} // 30%, 20%, 50%
	for i, expectedRatio := range expectedRatios {
		actualRatio := actualWidths[i] / totalWidth
		if !approximatelyEqual(actualRatio, expectedRatio, 0.01) {
			t.Errorf("Column %d ratio: expected %f, got %f", i, expectedRatio, actualRatio)
		}
	}
}

// Helper function to validate fixed width calculations
func validateFixedWidths(t *testing.T, table *docTable, actualWidths []float64, expectedWidths []float64) {
	// Check that each column has exactly the specified width
	for i, expectedWidth := range expectedWidths {
		if !approximatelyEqual(actualWidths[i], expectedWidth, 0.1) {
			t.Errorf("Column %d width: expected %f, got %f", i, expectedWidth, actualWidths[i])
		}
	}
}

// Helper function to validate auto width calculations
func validateAutoWidths(t *testing.T, table *docTable, actualWidths []float64) {
	// Check that longest column title has largest width
	longestColIndex := 2 // "Third Column With Long Title"
	for i := range actualWidths {
		if i != longestColIndex && actualWidths[i] >= actualWidths[longestColIndex] {
			t.Errorf("Column %d width (%f) should be less than column %d width (%f)",
				i, actualWidths[i], longestColIndex, actualWidths[longestColIndex])
		}
	}
}

// Helper function to check if two float values are approximately equal
func approximatelyEqual(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}
