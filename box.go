package docpdf

// unitsToPoints converts the box coordinates to Points.
// When this is called it is assumed the values of the box are in the specified unit system.
// The method creates a new box instance with coordinates converted to points.
//
// Parameters:
//   - unit: Either an integer representing a unit type (UnitPT, UnitMM, etc.) or a unitConfigurator interface
//
// Returns:
//   - A new box pointer with coordinates converted to points
func (b *box) unitsToPoints(unit any) (out *box) {
	if b == nil {
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
	if b.unitOverride.getUnit() != UnitUnset {
		unitCfg = b.unitOverride
	}

	out = &box{
		Left:   b.Left,
		Top:    b.Top,
		Right:  b.Right,
		Bottom: b.Bottom,
	}
	unitsToPointsVar(unitCfg, &out.Left, &out.Top, &out.Right, &out.Bottom)
	return
}
