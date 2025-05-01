package mathutils

// MinMax returns the minimum and maximum of a given set of values.
func MinMax(values ...float64) (min, max float64) {
	if len(values) == 0 {
		return
	}

	max = values[0]
	min = values[0]
	var value float64
	for index := 1; index < len(values); index++ {
		value = values[index]
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return
}

// MinInt returns the minimum int.
func MinInt(values ...int) (min int) {
	if len(values) == 0 {
		return
	}

	min = values[0]
	var value int
	for index := 1; index < len(values); index++ {
		value = values[index]
		if value < min {
			min = value
		}
	}
	return
}

// MaxInt returns the maximum int.
func MaxInt(values ...int) (max int) {
	if len(values) == 0 {
		return
	}

	max = values[0]
	var value int
	for index := 1; index < len(values); index++ {
		value = values[index]
		if value > max {
			max = value
		}
	}
	return
}
