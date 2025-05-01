package canvas

// UnitsToPoints converts the box coordinates to Points.
// When this is called it is assumed the values of the box are in the specified unit system.
// The method creates a new box instance with coordinates converted to points.
//
// Parameters:
//   - unit: Either an integer representing a unit type (UnitPT, UnitMM, etc.) or a unitConfigurator interface
//
// Returns:
//   - A new box pointer with coordinates converted to points
func (b *Box) UnitsToPoints(unit any) (out *Box) {
	if b == nil {
		return
	}

	var unitCfg unitConfigurator

	// Determinar la configuración de unidades según el tipo de parámetro
	switch u := unit.(type) {
	case int:
		unitCfg = defaultUnitConfig{Unit: u}
	case unitConfigurator:
		unitCfg = u
	default:
		unitCfg = defaultUnitConfig{Unit: UnitPT}
	}

	// Aplicar anulación de unidad si está establecida
	if b.unitOverride.GetUnit() != UnitUnset {
		unitCfg = b.unitOverride
	}

	// Crear variables temporales float64
	left := float64(b.Left)
	top := float64(b.Top)
	right := float64(b.Right)
	bottom := float64(b.Bottom)

	// Aplicar la conversión a las variables temporales
	unitsToPointsVarCfg(unitCfg, &left, &top, &right, &bottom)

	// Crear el box de salida con valores enteros redondeados
	out = &Box{
		Left:   int(left + 0.5),   // Redondeo
		Top:    int(top + 0.5),    // Redondeo
		Right:  int(right + 0.5),  // Redondeo
		Bottom: int(bottom + 0.5), // Redondeo
	}
	return
}
