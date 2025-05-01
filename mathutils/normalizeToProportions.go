package mathutils

// NormalizeToProportions convierte un conjunto de números en sus proporciones relativas en el intervalo [0,1].
// Por ejemplo: 4,3,2,1 => 0.4, 0.3, 0.2, 0.1
// Los valores se redondean para evitar errores de punto flotante, pero la suma puede no ser exactamente 1.0.
func NormalizeToProportions(values ...float64) []float64 {
	var total float64
	for _, v := range values {
		total += v
	}

	output := make([]float64, len(values))

	// Si el total es 0, devolver proporciones iguales (o 0 si no hay valores)
	if total == 0 {
		if len(values) > 0 {
			equalProportion := 1.0 / float64(len(values))
			for i := range output {
				output[i] = equalProportion
			}
		}
		return output
	}

	for i, v := range values {
		output[i] = RoundPlaces(v/total, 4) // Usar RoundPlaces en lugar de RoundDown para mayor precisión
	}
	return output
}
