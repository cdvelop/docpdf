package pdfengine

import "github.com/cdvelop/docpdf/canvas"

// Config defines the basic configuration for a PDF document.
// It includes settings for unit types, page size, protection, and more.
type Config struct {
	// Unit specifies the unit type to use when composing the document (canvas.UnitPT, canvas.UnitMM, canvas.UnitCM, canvas.UnitIN, canvas.UnitPX)
	Unit int

	// ConversionForUnit is a value used to convert units to points.
	// If this value is not 0, it will be used for unit conversion instead of the default constants.
	// When this is set, the Unit field is ignored.
	// Example: For 300 DPI, set this to 72.0/300.0
	ConversionForUnit float64

	// TrimBox defines the default trim canvas.Box for all pages in the document
	TrimBox canvas.Box

	// canvas.PageSize defines the default page size for all pages in the document
	PageSize canvas.Rect

	// K is a scaling factor (purpose not well-documented)
	K float64

	// Protection contains the settings for PDF document protection
	Protection PdfProtectionConfig
}

// GetUnit returns the unit type from the configuration
func (c Config) GetUnit() int {
	return c.Unit
}

// GetConversionForUnit returns the custom conversion factor from the configuration
func (c Config) GetConversionForUnit() float64 {
	return c.ConversionForUnit
}

// PdfProtectionConfig defines the configuration for PDF document protection.
type PdfProtectionConfig struct {
	// UseProtection determines whether to apply protection to the PDF
	UseProtection bool

	// Permissions specifies the allowed operations on the PDF (PermissionsPrint, PermissionsCopy, etc.)
	Permissions int

	// UserPass is the password required for general users to open the PDF
	UserPass []byte

	// OwnerPass is the password required for owners to remove restrictions
	OwnerPass []byte
}
