package docpdf

import (
	"testing"
)

func TestParseHeaderFormat(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedResult tableHeaderFormat
	}{
		{
			name:  "Simple header without options",
			input: "Name",
			expectedResult: tableHeaderFormat{
				HeaderTitle:     "Name",
				HeaderAlignment: Center,
				ColumnAlignment: Left,
				Prefix:          "",
				Suffix:          "",
				Width:           0,
				WidthMode:       widthModeAuto,
			},
		},
		{
			name:  "Header with left alignment",
			input: "Name|HL",
			expectedResult: tableHeaderFormat{
				HeaderTitle:     "Name",
				HeaderAlignment: Left,
				ColumnAlignment: Left,
				Prefix:          "",
				Suffix:          "",
				Width:           0,
				WidthMode:       widthModeAuto,
			},
		},
		{
			name:  "Header with right alignment",
			input: "Price|HR",
			expectedResult: tableHeaderFormat{
				HeaderTitle:     "Price",
				HeaderAlignment: Right,
				ColumnAlignment: Left,
				Prefix:          "",
				Suffix:          "",
				Width:           0,
				WidthMode:       widthModeAuto,
			},
		},
		{
			name:  "Header with column alignment",
			input: "Amount|CR",
			expectedResult: tableHeaderFormat{
				HeaderTitle:     "Amount",
				HeaderAlignment: Center,
				ColumnAlignment: Right,
				Prefix:          "",
				Suffix:          "",
				Width:           0,
				WidthMode:       widthModeAuto,
			},
		},
		{
			name:  "Header with right alignment and right column alignment",
			input: "Price|HR,CR",
			expectedResult: tableHeaderFormat{
				HeaderTitle:     "Price",
				HeaderAlignment: Right,
				ColumnAlignment: Right,
				Prefix:          "",
				Suffix:          "",
				Width:           0,
				WidthMode:       widthModeAuto,
			},
		},
		{
			name:  "Header with prefix",
			input: "Price|HR,CR,P:$",
			expectedResult: tableHeaderFormat{
				HeaderTitle:     "Price",
				HeaderAlignment: Right,
				ColumnAlignment: Right,
				Prefix:          "$",
				Suffix:          "",
				Width:           0,
				WidthMode:       widthModeAuto,
			},
		},
		{
			name:  "Header with suffix",
			input: "Percentage|HC,CC,S:%",
			expectedResult: tableHeaderFormat{
				HeaderTitle:     "Percentage",
				HeaderAlignment: Center,
				ColumnAlignment: Center,
				Prefix:          "",
				Suffix:          "%",
				Width:           0,
				WidthMode:       widthModeAuto,
			},
		},
		{
			name:  "Header with prefix and suffix",
			input: "Balance|HR,CR,P:$,S:USD",
			expectedResult: tableHeaderFormat{
				HeaderTitle:     "Balance",
				HeaderAlignment: Right,
				ColumnAlignment: Right,
				Prefix:          "$",
				Suffix:          "USD",
				Width:           0,
				WidthMode:       widthModeAuto,
			},
		},
		{
			name:  "Header with fixed width",
			input: "Name|HL,CL,W:120",
			expectedResult: tableHeaderFormat{
				HeaderTitle:     "Name",
				HeaderAlignment: Left,
				ColumnAlignment: Left,
				Prefix:          "",
				Suffix:          "",
				Width:           120,
				WidthMode:       widthModeFixed,
			},
		},
		{
			name:  "Header with percentage width",
			input: "Name|HL,CL,W:30%",
			expectedResult: tableHeaderFormat{
				HeaderTitle:     "Name",
				HeaderAlignment: Left,
				ColumnAlignment: Left,
				Prefix:          "",
				Suffix:          "",
				Width:           30,
				WidthMode:       widthModePercent,
			},
		},
		{
			name:  "Complete example with all options",
			input: "Product|HL,CR,P:Item:,S:USD,W:40%",
			expectedResult: tableHeaderFormat{
				HeaderTitle:     "Product",
				HeaderAlignment: Left,
				ColumnAlignment: Right,
				Prefix:          "Item:",
				Suffix:          "USD",
				Width:           40,
				WidthMode:       widthModePercent,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := parseHeaderFormat(tc.input)
			if result.HeaderTitle != tc.expectedResult.HeaderTitle {
				t.Errorf("Expected HeaderTitle to be %s, got %s", tc.expectedResult.HeaderTitle, result.HeaderTitle)
			}
			if result.HeaderAlignment != tc.expectedResult.HeaderAlignment {
				t.Errorf("Expected HeaderAlignment to be %d, got %d", tc.expectedResult.HeaderAlignment, result.HeaderAlignment)
			}
			if result.ColumnAlignment != tc.expectedResult.ColumnAlignment {
				t.Errorf("Expected ColumnAlignment to be %d, got %d", tc.expectedResult.ColumnAlignment, result.ColumnAlignment)
			}
			if result.Prefix != tc.expectedResult.Prefix {
				t.Errorf("Expected Prefix to be %s, got %s", tc.expectedResult.Prefix, result.Prefix)
			}
			if result.Suffix != tc.expectedResult.Suffix {
				t.Errorf("Expected Suffix to be %s, got %s", tc.expectedResult.Suffix, result.Suffix)
			}
			if result.Width != tc.expectedResult.Width {
				t.Errorf("Expected Width to be %f, got %f", tc.expectedResult.Width, result.Width)
			}
			if result.WidthMode != tc.expectedResult.WidthMode {
				t.Errorf("Expected WidthMode to be %v, got %v", tc.expectedResult.WidthMode, result.WidthMode)
			}
		})
	}
}
