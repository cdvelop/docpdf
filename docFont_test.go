package docpdf

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cdvelop/docpdf/config"
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

		doc := NewDocument()

		// Updated expected default font names
		expectedFont := config.FontFamily{
			Regular: "regular.ttf",
			Bold:    "bold.ttf",
			Italic:  "italic.ttf",
			Path:    "fonts/",
		}

		if doc.fontConfig.FontFamily != expectedFont {
			t.Errorf("got font = %v, want %v", doc.fontConfig.FontFamily, expectedFont)
		}
	})

	t.Run("Custom font configuration", func(t *testing.T) {
		customFont := config.FontFamily{
			Regular: "font.ttf",
			Bold:    "font-bold.ttf",
			Italic:  "font-italic.ttf",
			Path:    "custom/",
		}

		doc := NewDocument(customFont)

		if doc.fontConfig.FontFamily != customFont {
			t.Errorf("got font = %v, want %v", doc.fontConfig.FontFamily, customFont)
		}
	})
	t.Run("Logger captures errors", func(t *testing.T) {
		var logOutput []any

		// Create a custom logger to capture Log output
		customLogger := func(a ...any) {
			logOutput = append(logOutput, a...)
		}

		// Create a custom font that doesn't exist
		customFont := config.FontFamily{
			Regular: "nonexistent/font.ttf",
			Bold:    "nonexistent/font-bold.ttf",
			Italic:  "nonexistent/font-italic.ttf",
			Path:    "",
		}

		// Create a document with custom logger and the nonexistent font
		NewDocument(customLogger, customFont)

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

		oneCustomFont := config.FontFamily{
			Regular: "LiberationSerif-Regular.ttf", // Use a different name for clarity
			Path:    "pdfengine/test/res/",
		}
		// Create a document with custom logger
		doc := NewDocument(func(a ...any) {
			logOutput = append(logOutput, a...)
		}, oneCustomFont)

		// The expected result after NewDocument applies the fallback logic
		expectedFont := config.FontFamily{
			Regular: oneCustomFont.Regular,
			Bold:    oneCustomFont.Regular, // Should fallback to regular
			Italic:  oneCustomFont.Regular, // Should fallback to regular
			Path:    oneCustomFont.Path,
		}

		if len(logOutput) != 0 {
			t.Error("Expected no errors when loading only one font", logOutput)
		}

		// Compare the actual font config set by NewDocument with the expected one
		if doc.fontConfig.FontFamily != expectedFont {
			t.Errorf("got font = %v, want %v", doc.fontConfig.FontFamily, expectedFont)
		}
	})
}
