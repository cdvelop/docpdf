package docpdf

import (
	"fmt"
	"strings"
	"testing"
)

func TestExtractFontName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Simple TTF path",
			path:     "fonts/RubikBold.ttf",
			expected: "RubikBold",
		},
		{
			name:     "OTF extension",
			path:     "fonts/Arial.otf",
			expected: "Arial",
		},
		{
			name:     "Multiple dots in filename",
			path:     "fonts/Open.Sans.Bold.ttf",
			expected: "OpenSansBold",
		},
		{
			name:     "Deep nested path",
			path:     "assets/fonts/subfolder/Helvetica.ttf",
			expected: "Helvetica",
		},
		{
			name:     "No extension",
			path:     "fonts/ComicSans",
			expected: "ComicSans",
		},
		{
			name:     "Just filename",
			path:     "RubikBold.ttf",
			expected: "RubikBold",
		},
		{
			name:     "Windows style path",
			path:     "fonts\\RubikBold.ttf",
			expected: "RubikBold",
		},
		{
			name:     "Empty path",
			path:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractNameFromPath(tt.path)
			if got != tt.expected {
				t.Errorf("extractNameFromPath() = %v, want %v", got, tt.expected)
			}
		})
	}
}
func TestNewDocument(t *testing.T) {
	t.Run("Default settings", func(t *testing.T) {
		var logOutput []any
		logger := func(a ...any) {
			logOutput = append(logOutput, a...)
		}

		doc := NewDocument(logger)

		if doc == nil {
			t.Fatal("Expected document to be created")
		}

		expectedFont := Font{
			Regular: "Rubik-Regular.ttf",
			Bold:    "Rubik-Bold.ttf",
			Italic:  "Rubik-Italic.ttf",
			Path:    "fonts/",
		}

		if doc.fontConfig.Family != expectedFont {
			t.Errorf("got font = %v, want %v", doc.fontConfig.Family, expectedFont)
		}
	})

	t.Run("Custom font configuration", func(t *testing.T) {
		customFont := Font{
			Regular: "font.ttf",
			Bold:    "font-bold.ttf",
			Italic:  "font-italic.ttf",
			Path:    "custom/",
		}

		doc := NewDocument(func(a ...any) {}, customFont)

		if doc.fontConfig.Family != customFont {
			t.Errorf("got font = %v, want %v", doc.fontConfig.Family, customFont)
		}
	})

	t.Run("Logger captures errors", func(t *testing.T) {
		var logOutput []any
		logger := func(a ...any) {
			logOutput = append(logOutput, a...)
		}

		customFont := Font{
			Regular: "nonexistent/font.ttf",
			Bold:    "nonexistent/font-bold.ttf",
			Italic:  "nonexistent/font-italic.ttf",
		}

		NewDocument(logger, customFont)

		if len(logOutput) == 0 {
			t.Error("Expected logger to capture font loading error")
		}

		errorMsg := fmt.Sprint(logOutput...)
		if !strings.Contains(errorMsg, "Error loading fonts") {
			t.Errorf("Expected error message about font loading, got: %v", errorMsg)
		}
	})

	t.Run("Load only one font", func(t *testing.T) {
		var logOutput []any
		logger := func(a ...any) {
			logOutput = append(logOutput, a...)
		}

		oneCustomFont := Font{
			Regular: "Rubik-Regular.ttf",
			Path:    "fonts/",
		}

		doc := NewDocument(logger, oneCustomFont)

		expectedFont := Font{
			Regular: "Rubik-Regular.ttf",
			Bold:    "Rubik-Regular.ttf",
			Italic:  "Rubik-Regular.ttf",
			Path:    "fonts/",
		}

		if len(logOutput) != 0 {
			t.Error("Expected no errors when loading only one font", logOutput)
		}

		if doc.fontConfig.Family != expectedFont {
			t.Errorf("got font = %v, want %v", doc.fontConfig.Family, expectedFont)
		}

	})
}
