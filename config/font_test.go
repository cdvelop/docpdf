package config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cdvelop/docpdf/config"
)

// MockPDFEngine for testing font loading
type mockPDFEngine struct {
	fontLoadError bool
	loadedFont    config.FontFamily
}

func (m *mockPDFEngine) AddFontFamilyConfig(fontFamily config.FontFamily) error {
	m.loadedFont = fontFamily
	if m.fontLoadError {
		return fmt.Errorf("Error loading fonts")
	}
	return nil
}

func (m *mockPDFEngine) SetFont(family string, style string, size any) error {
	return nil
}

func (m *mockPDFEngine) SetTextColor(r, g, b uint8) {}

func (m *mockPDFEngine) SetStrokeColor(r, g, b uint8) {}

func (m *mockPDFEngine) SetLineWidth(width float64) {}

func TestFontConfiguration(t *testing.T) {
	t.Run("Default settings", func(t *testing.T) {
		// Get default text config which includes default font family
		textStyles := config.DefaultTextConfig()

		// Expected default font names from DefaultTextConfig
		expectedFont := config.FontFamily{
			Regular:   "regular.ttf",
			Bold:      "bold.ttf",
			Italic:    "italic.ttf",
			Underline: "",
			Path:      "fonts/",
		}

		if textStyles.GetFontFamily() != expectedFont {
			t.Errorf("got font = %v, want %v", textStyles.GetFontFamily(), expectedFont)
		}
	})

	t.Run("Custom font configuration", func(t *testing.T) {
		textStyles := config.TextStyles{}

		customFont := config.FontFamily{
			Regular: "font.ttf",
			Bold:    "font-bold.ttf",
			Italic:  "font-italic.ttf",
			Path:    "custom/",
		}

		textStyles.SetFontFamily(customFont)

		if textStyles.GetFontFamily() != customFont {
			t.Errorf("got font = %v, want %v", textStyles.GetFontFamily(), customFont)
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

		// Setup mock engine with error flag
		mockEngine := &mockPDFEngine{fontLoadError: true}

		// Create text styles and set custom font
		textStyles := config.TextStyles{}
		textStyles.SetFontFamily(customFont)

		// Try to load fonts, which should trigger an error
		err := textStyles.LoadFonts(mockEngine)

		// Log the error using the custom logger
		if err != nil {
			customLogger("Error loading fonts:", err)
		}

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

		// Setup mock engine without error flag
		mockEngine := &mockPDFEngine{fontLoadError: false}

		// Create text styles with our custom font
		textStyles := config.TextStyles{}
		textStyles.SetFontFamily(oneCustomFont)

		// Load fonts, which should apply the fallback logic inside LoadFonts
		err := textStyles.LoadFonts(mockEngine)

		// Log any errors using the custom logger
		if err != nil {
			customLogger := func(a ...any) {
				logOutput = append(logOutput, a...)
			}
			customLogger("Error loading fonts:", err)
		}

		// The expected result after LoadFonts applies the fallback logic
		expectedFont := config.FontFamily{
			Regular: oneCustomFont.Regular,
			Bold:    oneCustomFont.Regular, // Should fallback to regular
			Italic:  oneCustomFont.Regular, // Should fallback to regular
			Path:    oneCustomFont.Path,
		}

		if len(logOutput) != 0 {
			t.Error("Expected no errors when loading only one font", logOutput)
		}

		// Compare the actual font config set by LoadFonts with the expected one
		// We can check this from the mock engine's loaded font
		if mockEngine.loadedFont != expectedFont {
			t.Errorf("got font = %v, want %v", mockEngine.loadedFont, expectedFont)
		}
	})
}
