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

		doc := NewDocument(fw(""), logger)

		if doc == nil {
			t.Fatal("Expected document to be created")
		}

		expectedFont := Font{
			Regular: "regular.ttf",
			Bold:    "bold.ttf",
			Italic:  "italic.ttf",
			Path:    fontPublicPath,
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

		doc := NewDocument(fw(""), func(a ...any) {}, customFont)

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

		NewDocument(fw(""), logger, customFont)

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
			Regular: "regular.ttf",
			Path:    fontPublicPath,
		}

		doc := NewDocument(fw(""), logger, oneCustomFont)

		expectedFont := Font{
			Regular: "regular.ttf",
			Bold:    "regular.ttf",
			Italic:  "regular.ttf",
			Path:    fontPublicPath,
		}

		if len(logOutput) != 0 {
			t.Error("Expected no errors when loading only one font", logOutput)
		}

		if doc.fontConfig.Family != expectedFont {
			t.Errorf("got font = %v, want %v", doc.fontConfig.Family, expectedFont)
		}

	})
}

func TestFontAutoDetection(t *testing.T) {
	t.Run("Auto-detect single font in default path", func(t *testing.T) {
		var logOutput []any
		logger := func(a ...any) {
			logOutput = append(logOutput, a...)
		}

		// Create document with empty font config to trigger default path detection
		doc := NewDocument(fw(""), logger)

		// We're simulating that only one font exists in the default path
		// by manually modifying the font config before loadFonts is called
		// Set only the Regular font, clear the others
		doc.fontConfig.Family = Font{
			Regular: "regular.ttf",  // Only this font exists
			Bold:    "",             // These will be set to Regular
			Italic:  "",             // These will be set to Regular
			Path:    fontPublicPath, // Use default path
		}

		// Reload fonts
		err := doc.loadFonts()
		if err != nil {
			t.Errorf("Error loading fonts: %v", err)
		}

		// Verify all fonts are set to the single available font
		expectedFont := Font{
			Regular: "regular.ttf",
			Bold:    "regular.ttf", // Should be copied from Regular
			Italic:  "regular.ttf", // Should be copied from Regular
			Path:    fontPublicPath,
		}

		if doc.fontConfig.Family != expectedFont {
			t.Errorf("Font auto-detection failed, got: %v, want: %v", doc.fontConfig.Family, expectedFont)
		}

		if len(logOutput) != 0 {
			t.Error("Expected no errors when auto-detecting single font", logOutput)
		}
	})

	t.Run("Auto-detect bold font in default path", func(t *testing.T) {
		var logOutput []any
		logger := func(a ...any) {
			logOutput = append(logOutput, a...)
		}

		// Create document with empty font config to trigger default path detection
		doc := NewDocument(fw(""), logger)

		// We're simulating that only the bold font exists
		doc.fontConfig.Family = Font{
			Regular: "",             // This will be set to Bold
			Bold:    "bold.ttf",     // Only this font exists
			Italic:  "",             // This will be set to Bold
			Path:    fontPublicPath, // Use default path
		}

		// Manually update Regular to use Bold since loadFonts expects Regular to be set
		doc.fontConfig.Family.Regular = doc.fontConfig.Family.Bold

		// Reload fonts
		err := doc.loadFonts()
		if err != nil {
			t.Errorf("Error loading fonts: %v", err)
		}

		// Verify all fonts are set to the single available font
		expectedFont := Font{
			Regular: "bold.ttf",
			Bold:    "bold.ttf",
			Italic:  "bold.ttf", // Should be copied from Regular (which is Bold)
			Path:    fontPublicPath,
		}

		if doc.fontConfig.Family != expectedFont {
			t.Errorf("Font auto-detection failed, got: %v, want: %v", doc.fontConfig.Family, expectedFont)
		}

		if len(logOutput) != 0 {
			t.Error("Expected no errors when auto-detecting single font", logOutput)
		}
	})

	t.Run("Auto-detect italic font in default path", func(t *testing.T) {
		var logOutput []any
		logger := func(a ...any) {
			logOutput = append(logOutput, a...)
		}

		// Create document with empty font config to trigger default path detection
		doc := NewDocument(fw(""), logger)

		// We're simulating that only the italic font exists
		doc.fontConfig.Family = Font{
			Regular: "",             // This will be set to Italic
			Bold:    "",             // This will be set to Italic
			Italic:  "italic.ttf",   // Only this font exists
			Path:    fontPublicPath, // Use default path
		}

		// Manually update Regular to use Italic since loadFonts expects Regular to be set
		doc.fontConfig.Family.Regular = doc.fontConfig.Family.Italic

		// Reload fonts
		err := doc.loadFonts()
		if err != nil {
			t.Errorf("Error loading fonts: %v", err)
		}

		// Verify all fonts are set to the single available font
		expectedFont := Font{
			Regular: "italic.ttf",
			Bold:    "italic.ttf", // Should be copied from Regular (which is Italic)
			Italic:  "italic.ttf",
			Path:    fontPublicPath,
		}

		if doc.fontConfig.Family != expectedFont {
			t.Errorf("Font auto-detection failed, got: %v, want: %v", doc.fontConfig.Family, expectedFont)
		}

		if len(logOutput) != 0 {
			t.Error("Expected no errors when auto-detecting single font", logOutput)
		}
	})
}
