package pdfengine

import "testing"

func TestCustomTrimSpace(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"", ""},
		{" ", ""},
		{"\n\t  ", ""},
		{"Hola  ", "Hola"},
		{"   Hola", "Hola"},
		{"   Hola Mundo   ", "Hola Mundo"},
	}

	for _, tt := range tests {
		got := customTrimSpace(tt.input)
		if got != tt.output {
			t.Errorf("customTrimSpace(%q) = %q; want %q", tt.input, got, tt.output)
		}
	}
}
