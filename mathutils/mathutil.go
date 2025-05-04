package mathutils

import "math"

// AbsInt returns the absolute value of an int.
func AbsInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

// DegreesToRadians returns degrees as radians.
func DegreesToRadians(degrees float64) float64 {
	return degrees * _d2r
}

// RadiansToDegrees translates a radian value to a degree value.
func RadiansToDegrees(value float64) float64 {
	return math.Mod(value, TwoPi) * _r2d
}

// PercentToRadians converts a normalized value (0,1) to radians.
func PercentToRadians(pct float64) float64 {
	return DegreesToRadians(360.0 * pct)
}

// RadianAdd adds a delta to a base in radians.
func RadianAdd(base, delta float64) float64 {
	value := base + delta
	if value > TwoPi {
		return math.Mod(value, TwoPi)
	} else if value < 0 {
		return math.Mod(TwoPi+value, TwoPi)
	}
	return value
}

// DegreesAdd adds a delta to a base in radians.
func DegreesAdd(baseDegrees, deltaDegrees float64) float64 {
	value := baseDegrees + deltaDegrees
	if value > TwoPi {
		return math.Mod(value, 360.0)
	} else if value < 0 {
		return math.Mod(360.0+value, 360.0)
	}
	return value
}

// DegreesToCompass returns the degree value in compass / clock orientation.
func DegreesToCompass(deg float64) float64 {
	return DegreesAdd(deg, -90.0)
}

// CirclePoint returns the absolute alignment.Alignment of a circle diameter point given
// by the radius and the theta.
func CirclePoint(cx, cy int, radius, thetaRadians float64) (x, y int) {
	x = cx + int(radius*math.Sin(thetaRadians))
	y = cy - int(radius*math.Cos(thetaRadians))
	return
}

// RotateCoordinate rotates a coordinate around a given center by a theta in radians.
func RotateCoordinate(cx, cy, x, y int, thetaRadians float64) (rx, ry int) {
	tempX, tempY := float64(x-cx), float64(y-cy)
	rotatedX := tempX*math.Cos(thetaRadians) - tempY*math.Sin(thetaRadians)
	rotatedY := tempX*math.Sin(thetaRadians) + tempY*math.Cos(thetaRadians)
	rx = int(rotatedX) + cx
	ry = int(rotatedY) + cy
	return
}

// Round rounds a float64 to the nearest int.
// eg: 1.5 -> 2, 1.4 -> 1, -1.5 -> -2, -1.4 -> -1
func Round(value float64) int {
	if value < 0.0 {
		value -= 0.5
	} else {
		value += 0.5
	}
	return int(value)
}

// RoundUp rounds up to a given roundTo value.
func RoundUp(value, roundTo float64) float64 {
	if roundTo < 0.000000000000001 {
		return value
	}
	d1 := math.Ceil(value / roundTo)
	return d1 * roundTo
}

// RoundDown rounds down to a given roundTo value.
func RoundDown(value, roundTo float64) float64 {
	if roundTo < 0.000000000000001 {
		return value
	}
	d1 := math.Floor(value / roundTo)
	return d1 * roundTo
}

// Mean returns the mean of a set of values
func Mean(values ...float64) float64 {
	return Sum(values...) / float64(len(values))
}

// MeanInt returns the mean of a set of integer values.
func MeanInt(values ...int) int {
	return SumInt(values...) / len(values)
}

// Sum sums a set of values.
func Sum(values ...float64) float64 {
	var total float64
	for _, v := range values {
		total += v
	}
	return total
}

// SumInt sums a set of values.
func SumInt(values ...int) int {
	var total int
	for _, v := range values {
		total += v
	}
	return total
}

// PercentDifference computes the percentage difference between two values.
// The formula is (v2-v1)/v1.
func PercentDifference(v1, v2 float64) float64 {
	if v1 == 0 {
		return 0
	}
	return (v2 - v1) / v1
}

// GetRoundToForDelta returns a `roundTo` value for a given delta.
func GetRoundToForDelta(delta float64) float64 {
	startingDeltaBound := math.Pow(10.0, 10.0)
	for cursor := startingDeltaBound; cursor > 0; cursor /= 10.0 {
		if delta > cursor {
			return cursor / 10.0
		}
	}

	return 0.0
}

// RoundPlaces rounds an input to a given places.
func RoundPlaces(input float64, places int) (rounded float64) {
	if math.IsNaN(input) {
		return 0.0
	}

	sign := 1.0
	if input < 0 {
		sign = -1
		input *= -1
	}

	precision := math.Pow(10, float64(places))
	digit := input * precision
	_, decimal := math.Modf(digit)

	if decimal >= 0.5 {
		rounded = math.Ceil(digit)
	} else {
		rounded = math.Floor(digit)
	}

	return rounded / precision * sign
}

// RoundToInt convierte un float64 a int redondeando al entero más cercano
func RoundToInt(value float64) int {
	r := RoundPlaces(value, 0)
	return int(r)
}
