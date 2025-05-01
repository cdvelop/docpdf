package canvas

// canvas.Rect defines a rectangle by its width and height.
// This is used for defining page sizes, content areas, and other rectangular regions in PDF documents.
// The dimensions are stored in the current unit system (points by default, but can be mm, cm, inches, or pixels).
type Rect struct {
	W            float64 // Width of the rectangle
	H            float64 // Height of the rectangle
	unitOverride defaultUnitConfig
}

// defaultUnitConfig is the standard implementation of the unitConfigurator interface.
// It stores the unit type and an optional custom conversion factor.
type defaultUnitConfig struct {
	// Unit specifies the unit type (UnitPT, UnitMM, UnitCM, UnitIN, UnitPX)
	Unit int

	// ConversionForUnit is an optional custom conversion factor
	ConversionForUnit float64
}

// pointsToUnits converts the rectangle's width and height from points to the specified unit system.
// When this is called it is assumed the values of the rectangle are in points.
// The method creates a new canvas.Rect instance with dimensions converted to the specified units.
//
// Parameters:
//   - t: An integer representing the unit type to convert to (UnitPT, UnitMM, UnitCM, UnitIN, UnitPX)
//
// Returns:
//   - A new canvas.Rect pointer with dimensions converted to the specified units
func (rect *Rect) pointsToUnits(t int) (r *Rect) {
	if rect == nil {
		return
	}

	unitCfg := defaultUnitConfig{Unit: t}
	if rect.unitOverride.GetUnit() != UnitUnset {
		unitCfg = rect.unitOverride
	}

	r = &Rect{W: rect.W, H: rect.H}
	pointsToUnitsVar(unitCfg, &r.W, &r.H)
	return
}

// UnitsToPoints converts the rectangle's dimensions to points based on the provided unit information.
// When this is called it is assumed the values of the rectangle are in the specified units.
// The method creates a new canvas.Rect instance with dimensions converted to points.
//
// Parameters:
//   - unit: Either an integer representing a unit type (UnitPT, UnitMM, etc.) or a unitConfigurator interface
//
// Returns:
//   - A new canvas.Rect pointer with dimensions converted to points
func (rect *Rect) UnitsToPoints(unit any) (r *Rect) {
	if rect == nil {
		return
	}

	var unitCfg unitConfigurator

	// Determine the unit configuration based on parameter type
	switch u := unit.(type) {
	case int:
		unitCfg = defaultUnitConfig{Unit: u}
	case unitConfigurator:
		unitCfg = u
	default:
		// Default to points if invalid type passed
		unitCfg = defaultUnitConfig{Unit: UnitPT}
	}

	// Apply unit override if set
	if rect.unitOverride.GetUnit() != UnitUnset {
		unitCfg = rect.unitOverride
	}

	r = &Rect{W: rect.W, H: rect.H}
	UnitsToPointsVar(unitCfg, &r.W, &r.H)
	return
}

// unitConfigurator is an interface that defines methods for retrieving unit configuration.
// It allows different unit configuration implementations to be used with the conversion functions.
type unitConfigurator interface {
	// getUnit returns the unit type (UnitPT, UnitMM, UnitCM, UnitIN, UnitPX)
	GetUnit() int

	// getConversionForUnit returns the custom conversion factor, if any
	GetConversionForUnit() float64
}

// getUnit returns the unit type from the defaultUnitConfig
func (d defaultUnitConfig) GetUnit() int {
	return d.Unit
}

// getConversionForUnit returns the custom conversion factor from the defaultUnitConfig
func (d defaultUnitConfig) GetConversionForUnit() float64 {
	return d.ConversionForUnit
}

// pointsToUnits convierte un valor de puntos al sistema de unidades especificado
// Parámetros:
//   - unit: Entero representando tipo de unidad o una interfaz unitConfigurator
//   - u: El valor en puntos a convertir
// Retorna:
//   - El valor equivalente en el sistema de unidades especificado
func PointsToUnits(unit any, u float64) float64 {
	return PointsToUnitsCfg(getUnitConfigurator(unit), u)
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
func PointsToUnitsCfg(unitCfg unitConfigurator, u float64) float64 {
	if unitCfg.GetConversionForUnit() != 0 {
		return u / unitCfg.GetConversionForUnit()
	}
	switch unitCfg.GetUnit() {
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

// unitsToPointsVar convierte múltiples valores al sistema de unidades especificado a puntos
// Parámetros:
//   - unit: Entero representando tipo de unidad o una interfaz unitConfigurator
//   - u: Punteros a valores a convertir (modificados en el lugar)
func UnitsToPointsVar(unit any, u ...*float64) {
	unitsToPointsVarCfg(getUnitConfigurator(unit), u...)
}

// unitsToPointsVar is an internal function that converts multiple values from units to points
// using the provided unit configuration.
//
// Parameters:
//   - unitCfg: The unit configuration that specifies the unit type and any custom conversion factor
//   - u: Pointers to values to convert (modified in place)
func unitsToPointsVarCfg(unitCfg unitConfigurator, u ...*float64) {
	for x := range u {
		*u[x] = unitsToPointsCfg(unitCfg, *u[x])
	}
}

// pointsToUnitsVar convierte múltiples valores de puntos al sistema de unidades especificado
// Parámetros:
//   - unit: Entero representando tipo de unidad o una interfaz unitConfigurator
//   - u: Punteros a valores a convertir (modificados en el lugar)
func pointsToUnitsVar(unit any, u ...*float64) {
	PointsToUnitsVarCfg(getUnitConfigurator(unit), u...)
}

// pointsToUnitsVar is an internal function that converts multiple values from points to units
// using the provided unit configuration.
//
// Parameters:
//   - unitCfg: The unit configuration that specifies the unit type and any custom conversion factor
//   - u: Pointers to values to convert (modified in place)
func PointsToUnitsVarCfg(unitCfg unitConfigurator, u ...*float64) {
	for x := range u {
		*u[x] = PointsToUnitsCfg(unitCfg, *u[x])
	}
}

// UnitsToPoints is an internal function that converts units to points using the provided
// unit configuration. It handles custom conversion factors and standard unit types.
//
// Parameters:
//   - unitCfg: The unit configuration that specifies the unit type and any custom conversion factor
//   - u: The value to convert
//
// Returns:
//   - The equivalent value in points
func unitsToPointsCfg(unitCfg unitConfigurator, u float64) float64 {
	if unitCfg.GetConversionForUnit() != 0 {
		return u * unitCfg.GetConversionForUnit()
	}
	switch unitCfg.GetUnit() {
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

// UnitsToPoints convierte un valor desde el sistema de unidades especificado a puntos
// Parámetros:
//   - unit: Entero representando tipo de unidad o una interfaz unitConfigurator
//   - u: El valor a convertir
// Retorna:
//   - El valor equivalente en puntos
func UnitsToPoints(unit any, u float64) float64 {
	return unitsToPointsCfg(getUnitConfigurator(unit), u)
}

// getUnitConfigurator extrae la configuración de unidades de diferentes tipos de entrada
func getUnitConfigurator(unit any) unitConfigurator {
	switch t := unit.(type) {
	case int:
		return defaultUnitConfig{Unit: t}
	case unitConfigurator:
		return t
	default:
		return defaultUnitConfig{Unit: UnitPT}
	}
}
