package docpdf

// The units that can be used in the document
const (
	UnitUnset = iota // No units were set, when conversion is called on nothing will happen
	UnitPT           // Points - 1/72 of an inch, traditional unit in PDF documents
	UnitMM           // Millimeters - 1/10 of a centimeter, metric measurement unit
	UnitCM           // Centimeters - 1/100 of a meter, metric measurement unit
	UnitIN           // Inches - Imperial unit equal to 72 points
	UnitPX           // Pixels - screen unit (by default 96 DPI, thus 72/96 = 3/4 point)

	// The math needed to convert units to points
	conversionUnitPT = 1.0
	conversionUnitMM = 72.0 / 25.4
	conversionUnitCM = 72.0 / 2.54
	conversionUnitIN = 72.0
	//We use a dpi of 96 dpi as the default, so we get a conversionUnitPX = 3.0 / 4.0, which comes from 72.0 / 96.0.
	//If you want to change this value, you can change it at config.ConversionForUnit
	//example: If you use dpi at 300.0
	//pdf.Start(docpdf.config{PageSize: *docpdf.PageSizeA4, ConversionForUnit: 72.0 / 300.0 })
	conversionUnitPX = 3.0 / 4.0
)

// The units that can be used in the document (for backward compatibility)
// Deprecated: Use UnitUnset,UnitPT,UnitMM,UnitCM,UnitIN  instead
const (
	Unit_Unset = UnitUnset // No units were set, when conversion is called on nothing will happen
	Unit_PT    = UnitPT    // Points
	Unit_MM    = UnitMM    // Millimeters
	Unit_CM    = UnitCM    // Centimeters
	Unit_IN    = UnitIN    // Inches
	Unit_PX    = UnitPX    // Pixels
)

// config defines the basic configuration for a PDF document.
// It includes settings for unit types, page size, protection, and more.
type config struct {
	// Unit specifies the unit type to use when composing the document (UnitPT, UnitMM, UnitCM, UnitIN, UnitPX)
	Unit int

	// ConversionForUnit is a value used to convert units to points.
	// If this value is not 0, it will be used for unit conversion instead of the default constants.
	// When this is set, the Unit field is ignored.
	// Example: For 300 DPI, set this to 72.0/300.0
	ConversionForUnit float64

	// TrimBox defines the default trim box for all pages in the document
	TrimBox box

	// PageSize defines the default page size for all pages in the document
	PageSize Rect

	// K is a scaling factor (purpose not well-documented)
	K float64

	// Protection contains the settings for PDF document protection
	Protection pdfProtectionConfig
}

// getUnit returns the unit type from the configuration
func (c config) getUnit() int {
	return c.Unit
}

// getConversionForUnit returns the custom conversion factor from the configuration
func (c config) getConversionForUnit() float64 {
	return c.ConversionForUnit
}

// pdfProtectionConfig defines the configuration for PDF document protection.
type pdfProtectionConfig struct {
	// UseProtection determines whether to apply protection to the PDF
	UseProtection bool

	// Permissions specifies the allowed operations on the PDF (PermissionsPrint, PermissionsCopy, etc.)
	Permissions int

	// UserPass is the password required for general users to open the PDF
	UserPass []byte

	// OwnerPass is the password required for owners to remove restrictions
	OwnerPass []byte
}

// unitsToPoints converts a value from the specified unit system to points (1/72 inch).
// This is a utility function for converting measurement units in the PDF document.
//
// Parameters:
//   - unit: Either an integer representing a unit type (UnitPT, UnitMM, etc.) or a unitConfigurator interface
//   - u: The value to convert
//
// Returns:
//   - The equivalent value in points
func unitsToPoints(unit interface{}, u float64) float64 {
	var unitCfg unitConfigurator

	// Determine the unit configuration based on parameter type
	switch t := unit.(type) {
	case int:
		unitCfg = defaultUnitConfig{Unit: t}
	case unitConfigurator:
		unitCfg = t
	default:
		// Default to points if invalid type passed
		unitCfg = defaultUnitConfig{Unit: UnitPT}
	}

	if unitCfg.getConversionForUnit() != 0 {
		return u * unitCfg.getConversionForUnit()
	}
	switch unitCfg.getUnit() {
	case UnitPT:
		return u * conversionUnitPT
	case UnitMM:
		return u * conversionUnitMM
	case UnitCM:
		return u * conversionUnitCM
	case UnitIN:
		return u * conversionUnitIN
	case UnitPX:
		return u * conversionUnitPX
	default:
		return u
	}
}

// unitsToPoints is an internal function that converts units to points using the provided
// unit configuration. It handles custom conversion factors and standard unit types.
//
// Parameters:
//   - unitCfg: The unit configuration that specifies the unit type and any custom conversion factor
//   - u: The value to convert
//
// Returns:
//   - The equivalent value in points
func unitsToPointsCfg(unitCfg unitConfigurator, u float64) float64 {
	if unitCfg.getConversionForUnit() != 0 {
		return u * unitCfg.getConversionForUnit()
	}
	switch unitCfg.getUnit() {
	case UnitPT:
		return u * conversionUnitPT
	case UnitMM:
		return u * conversionUnitMM
	case UnitCM:
		return u * conversionUnitCM
	case UnitIN:
		return u * conversionUnitIN
	case UnitPX:
		return u * conversionUnitPX
	default:
		return u
	}
}

// pointsToUnits converts a value from points to the specified unit system.
// This is a utility function for converting measurement units in the PDF document.
//
// Parameters:
//   - unit: Either an integer representing a unit type (UnitPT, UnitMM, etc.) or a unitConfigurator interface
//   - u: The value in points to convert
//
// Returns:
//   - The equivalent value in the specified unit system
func pointsToUnits(unit interface{}, u float64) float64 {
	var unitCfg unitConfigurator

	// Determine the unit configuration based on parameter type
	switch t := unit.(type) {
	case int:
		unitCfg = defaultUnitConfig{Unit: t}
	case unitConfigurator:
		unitCfg = t
	default:
		// Default to points if invalid type passed
		unitCfg = defaultUnitConfig{Unit: UnitPT}
	}

	if unitCfg.getConversionForUnit() != 0 {
		return u / unitCfg.getConversionForUnit()
	}
	switch unitCfg.getUnit() {
	case UnitPT:
		return u / conversionUnitPT
	case UnitMM:
		return u / conversionUnitMM
	case UnitCM:
		return u / conversionUnitCM
	case UnitIN:
		return u / conversionUnitIN
	case UnitPX:
		return u / conversionUnitPX
	default:
		return u
	}
}

// pointsToUnits is an internal function that converts points to the specified unit system
// using the provided unit configuration. It handles custom conversion factors and standard unit types.
//
// Parameters:
//   - unitCfg: The unit configuration that specifies the unit type and any custom conversion factor
//   - u: The value in points to convert
//
// Returns:
//   - The equivalent value in the specified unit system
func pointsToUnitsCfg(unitCfg unitConfigurator, u float64) float64 {
	if unitCfg.getConversionForUnit() != 0 {
		return u / unitCfg.getConversionForUnit()
	}
	switch unitCfg.getUnit() {
	case UnitPT:
		return u / conversionUnitPT
	case UnitMM:
		return u / conversionUnitMM
	case UnitCM:
		return u / conversionUnitCM
	case UnitIN:
		return u / conversionUnitIN
	case UnitPX:
		return u / conversionUnitPX
	default:
		return u
	}
}

// unitsToPointsVar converts multiple values from the specified unit system to points.
// This is a convenience function to convert multiple values at once.
//
// Parameters:
//   - unit: Either an integer representing a unit type (UnitPT, UnitMM, etc.) or a unitConfigurator interface
//   - u: Pointers to values to convert (modified in place)
func unitsToPointsVar(unit interface{}, u ...*float64) {
	var unitCfg unitConfigurator

	// Determine the unit configuration based on parameter type
	switch t := unit.(type) {
	case int:
		unitCfg = defaultUnitConfig{Unit: t}
	case unitConfigurator:
		unitCfg = t
	default:
		// Default to points if invalid type passed
		unitCfg = defaultUnitConfig{Unit: UnitPT}
	}

	for x := 0; x < len(u); x++ {
		*u[x] = unitsToPointsCfg(unitCfg, *u[x])
	}
}

// unitsToPointsVar is an internal function that converts multiple values from units to points
// using the provided unit configuration.
//
// Parameters:
//   - unitCfg: The unit configuration that specifies the unit type and any custom conversion factor
//   - u: Pointers to values to convert (modified in place)
func unitsToPointsVarCfg(unitCfg unitConfigurator, u ...*float64) {
	for x := 0; x < len(u); x++ {
		*u[x] = unitsToPointsCfg(unitCfg, *u[x])
	}
}

// pointsToUnitsVar converts multiple values from points to the specified unit system.
// This is a convenience function to convert multiple values at once.
//
// Parameters:
//   - unit: Either an integer representing a unit type (UnitPT, UnitMM, etc.) or a unitConfigurator interface
//   - u: Pointers to values to convert (modified in place)
func pointsToUnitsVar(unit interface{}, u ...*float64) {
	var unitCfg unitConfigurator

	// Determine the unit configuration based on parameter type
	switch t := unit.(type) {
	case int:
		unitCfg = defaultUnitConfig{Unit: t}
	case unitConfigurator:
		unitCfg = t
	default:
		// Default to points if invalid type passed
		unitCfg = defaultUnitConfig{Unit: UnitPT}
	}

	for x := 0; x < len(u); x++ {
		*u[x] = pointsToUnitsCfg(unitCfg, *u[x])
	}
}

// pointsToUnitsVar is an internal function that converts multiple values from points to units
// using the provided unit configuration.
//
// Parameters:
//   - unitCfg: The unit configuration that specifies the unit type and any custom conversion factor
//   - u: Pointers to values to convert (modified in place)
func pointsToUnitsVarCfg(unitCfg unitConfigurator, u ...*float64) {
	for x := 0; x < len(u); x++ {
		*u[x] = pointsToUnitsCfg(unitCfg, *u[x])
	}
}

// unitConfigurator is an interface that defines methods for retrieving unit configuration.
// It allows different unit configuration implementations to be used with the conversion functions.
type unitConfigurator interface {
	// getUnit returns the unit type (UnitPT, UnitMM, UnitCM, UnitIN, UnitPX)
	getUnit() int

	// getConversionForUnit returns the custom conversion factor, if any
	getConversionForUnit() float64
}

// getUnit returns the unit type from the defaultUnitConfig
func (d defaultUnitConfig) getUnit() int {
	return d.Unit
}

// getConversionForUnit returns the custom conversion factor from the defaultUnitConfig
func (d defaultUnitConfig) getConversionForUnit() float64 {
	return d.ConversionForUnit
}
