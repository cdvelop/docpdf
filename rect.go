package docpdf

// pointsToUnits converts the rectangle's width and height from points to the specified unit system.
// When this is called it is assumed the values of the rectangle are in points.
// The method creates a new Rect instance with dimensions converted to the specified units.
//
// Parameters:
//   - t: An integer representing the unit type to convert to (UnitPT, UnitMM, UnitCM, UnitIN, UnitPX)
//
// Returns:
//   - A new Rect pointer with dimensions converted to the specified units
func (rect *Rect) pointsToUnits(t int) (r *Rect) {
	if rect == nil {
		return
	}

	unitCfg := defaultUnitConfig{Unit: t}
	if rect.unitOverride.getUnit() != UnitUnset {
		unitCfg = rect.unitOverride
	}

	r = &Rect{W: rect.W, H: rect.H}
	pointsToUnitsVar(unitCfg, &r.W, &r.H)
	return
}

// unitsToPoints converts the rectangle's dimensions to points based on the provided unit information.
// When this is called it is assumed the values of the rectangle are in the specified units.
// The method creates a new Rect instance with dimensions converted to points.
//
// Parameters:
//   - unit: Either an integer representing a unit type (UnitPT, UnitMM, etc.) or a unitConfigurator interface
//
// Returns:
//   - A new Rect pointer with dimensions converted to points
func (rect *Rect) unitsToPoints(unit any) (r *Rect) {
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
	if rect.unitOverride.getUnit() != UnitUnset {
		unitCfg = rect.unitOverride
	}

	r = &Rect{W: rect.W, H: rect.H}
	unitsToPointsVar(unitCfg, &r.W, &r.H)
	return
}
